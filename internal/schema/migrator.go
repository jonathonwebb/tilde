package schema

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"time"
)

type Migration struct {
	Id   uint64
	Desc string
	Up   func(context.Context, *sql.DB, *slog.Logger) error
	Down func(context.Context, *sql.DB, *slog.Logger) error
}

type SchemaStore interface {
	db() *sql.DB
	init(context.Context) error
	lock(context.Context, time.Duration) error
	release(context.Context) error
	state(context.Context) ([]int64, error)
	commit(context.Context, int64) error
	revert(context.Context, int64) error
	dump(context.Context, io.Writer) error
	load(context.Context, io.Reader) error
	close() error
}

type Migrator struct {
	Store   SchemaStore
	Log     *slog.Logger
	Sources []Migration
}

var (
	migrationTmpl = template.Must(template.New("migration").Parse(`package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jonathonwebb/tilde/internal/schema"
)

var _{{.Id}}_{{.Name}} = schema.Migration{
	Id:   {{.IdNum}},
	Desc: "",
	Up: func(ctx context.Context, db *sql.DB, log *slog.Logger) (err error) {
		return fmt.Errorf("up migration {{.Id}}_{{.Name}} not implemented")
	},
	Down: func(ctx context.Context, db *sql.DB, log *slog.Logger) (err error) {
		return fmt.Errorf("down migration {{.Id}}_{{.Name}} not implemented")
	},
}
`))
	allTmpl = template.Must(template.New("allMigrations").Parse(`package migrations

import "github.com/jonathonwebb/tilde/internal/schema"

var All = []schema.Migration{
{{- range . }}
	{{.}},
{{- end }}
}
`))
)

func NewMigration(ctx context.Context, dir, migrationName string, ts time.Time) error {
	id := ts.Unix()
	label := fmt.Sprintf("%010d_%s", id, migrationName)

	p := path.Join(dir, fmt.Sprintf("%s.go", label))
	f, err := os.Create(p)
	if err != nil {
		return err
	}

	err = migrationTmpl.Execute(f, struct {
		Id    string
		IdNum int64
		Name  string
	}{fmt.Sprintf("%010d", id), id, migrationName})
	f.Close()
	if err != nil {
		return err
	}

	all := []string{}
	dirents, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, dirent := range dirents {
		name := dirent.Name()
		isReg := dirent.Type().IsRegular()
		isGo := filepath.Ext(name) == ".go"

		if isReg && isGo && name != "all.go" {
			base := strings.TrimSuffix(name, ".go")
			all = append(all, fmt.Sprintf("_%s", base))
		}
	}

	p = path.Join(dir, "all.go")
	f, err = os.Create(p)
	if err != nil {
		return err
	}

	err = allTmpl.Execute(f, all)
	f.Close()
	if err != nil {
		return err
	}

	return nil
}

func (m *Migrator) migrationIds() []int64 {
	ids := make([]int64, 0, len(m.Sources))
	for id := range m.Sources {
		ids = append(ids, int64(id))
	}
	slices.Sort(ids)
	return ids
}

func (m *Migrator) Init(ctx context.Context) error {
	return m.Store.init(ctx)
}

func (m *Migrator) Apply(ctx context.Context, v int64) (err error) {
	local := m.migrationIds()
	m.Log.Debug("read local migrations", "n", len(local))
	if !slices.Contains(local, v) {
		return fmt.Errorf("unknown version: %d", v)
	}

	if err := m.Init(ctx); err != nil {
		return fmt.Errorf("init store: %v", err)
	}

	shouldRelease := true
	if err := m.Store.lock(ctx, 1*time.Second); err != nil {
		return fmt.Errorf("lock store: %v", err)
	}
	defer func() {
		if shouldRelease {
			if rlErr := m.Store.release(ctx); rlErr != nil {
				err = errors.Join(err, fmt.Errorf("release store: %v", rlErr))
			}
		}
	}()

	remote, err := m.Store.state(ctx)
	if err != nil {
		return fmt.Errorf("get store state: %v", err)
	}
	var latest int64 = -1
	if len(remote) > 0 {
		latest = remote[len(remote)-1]
	}

	if latest == v {
		// up-to-date
		return nil
	}

	if latest < v {
		// migrate up
		for _, id := range local {
			if id > latest {
				src := m.Sources[uint64(id)]
				if err := src.Up(ctx, m.Store.db(), m.Log); err != nil {
					shouldRelease = false
					return err
				}
				if err := m.Store.commit(ctx, id); err != nil {
					shouldRelease = false
					return err
				}
			}
		}
	} else {
		// migrate down
		for i := len(remote) - 1; i >= 0; i++ {
			id := remote[i]
			if id > v {
				src := m.Sources[uint64(id)]
				if err := src.Down(ctx, m.Store.db(), m.Log); err != nil {
					shouldRelease = false
					return err
				}
				if err := m.Store.revert(ctx, id); err != nil {
					shouldRelease = false
					return err
				}
			}
		}
	}

	return nil
}

func (m *Migrator) ApplyLatest(ctx context.Context) error {
	local := m.migrationIds()
	if len(local) > 0 {
		return m.Apply(ctx, local[len(local)-1])
	}

	return nil
}

func (m *Migrator) ApplyInitial(ctx context.Context) error {
	return m.Apply(ctx, -1)
}

func (m *Migrator) Load(ctx context.Context, r io.Reader) error {
	return m.Store.load(ctx, strings.NewReader(Schema))
}

var schemaTmpl = template.Must(template.New("schema").Parse(`package schema

const (
	SchemaVersion = {{.SchemaVersion}}
	Schema        = ` + "`{{.Schema}}`" + `
)`,
))

func (m *Migrator) Dump(ctx context.Context, dir string, w io.Writer) error {
	remote, err := m.Store.state(ctx)
	if err != nil {
		return fmt.Errorf("get store state: %v", err)
	}

	var latest int64 = -1
	var b strings.Builder
	if len(remote) > 0 {
		latest = remote[len(remote)-1]
		m.Store.dump(ctx, &b)
	}

	p := path.Join(dir, "schema.go")
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()

	return schemaTmpl.Execute(f, struct {
		SchemaVersion int64
		Schema        string
	}{latest, b.String()})
}

func (m *Migrator) Close() error {
	return m.Store.close()
}
