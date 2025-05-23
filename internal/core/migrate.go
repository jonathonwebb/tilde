package core

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type migrationSource struct {
	name string
	up   string
	down string
}

type Migrator struct {
	db      *sql.DB
	dir     string
	pattern string
	log     *slog.Logger
}

func NewMigrator(db *sql.DB, dir string, log *slog.Logger) *Migrator {
	pattern := path.Join(dir, "*.sql")
	return &Migrator{
		db:      db,
		dir:     dir,
		pattern: pattern,
		log:     log,
	}
}

func (m *Migrator) Create(name string, content []byte) error {
	p := filepath.Join(m.dir, name)
	if _, err := os.Stat(p); !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("migration %s already exists", name)
	}
	return os.WriteFile(filepath.Join(m.dir, name), content, 0644)
}

func (m *Migrator) sources() (map[int]migrationSource, error) {
	files, err := filepath.Glob(m.pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid migration glob %q: %v", m.pattern, err)
	}

	migrations := map[int]migrationSource{}
	for _, p := range files {
		base := filepath.Base(p)
		ext := filepath.Ext(p)
		if ext != ".sql" {
			return nil, fmt.Errorf("invalid migration filename %q: \".sql\" extension is required", base)
		}
		name := strings.TrimSuffix(base, ext)
		parts := strings.SplitN(name, "_", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid migration name %q: separator '_' is required", name)
		}
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid migration id %q: non-integer", name)
		}
		if id < 0 {
			return nil, fmt.Errorf("invalid migration id %q: negative", name)
		}
		if _, ok := migrations[id]; ok {
			return nil, fmt.Errorf("invalid migration id %q: duplicate", name)
		}

		f, err := os.Open(p)
		if err != nil {
			return nil, fmt.Errorf("failed to open migration file: %v", err)
		}

		var section string
		var up, down strings.Builder

		s := bufio.NewScanner(f)
		for s.Scan() {
			line := s.Text()
			trimmed := strings.TrimSpace(line)
			switch trimmed {
			case "-- up":
				section = "up"
				continue
			case "-- down":
				section = "down"
				continue
			}

			switch section {
			case "up":
				up.WriteString(line)
				up.WriteRune('\n')
			case "down":
				down.WriteString(line)
				down.WriteRune('\n')
			}
		}

		if err := s.Err(); err != nil {
			return nil, err
		}

		if strings.TrimSpace(up.String()) == "" {
			return nil, fmt.Errorf("failed to parse migration file %q: missing up migration", p)
		}

		migrations[id] = migrationSource{name: p, up: strings.TrimSpace(up.String()), down: strings.TrimSpace(down.String())}
	}

	return migrations, nil
}

var (
	migrateTmpl = template.Must(template.New("migration").Parse("-- {{ .Id }}\n\n-- up\n\n-- down"))
)

func (m *Migrator) New(name string) error {
	if name == "" {
		return fmt.Errorf("migration name required")
	}

	t := time.Now().UTC().Format("20060102150405")
	id := fmt.Sprintf("%s_%s", t, name)

	var b bytes.Buffer
	if err := migrateTmpl.Execute(&b, struct{ Id string }{id}); err != nil {
		return err
	}
	up := fmt.Sprintf("%s_up.sql", id)
	if err := m.Create(up, b.Bytes()); err != nil {
		return err
	}
	m.log.Info("created migration file", "name", up)

	return nil
}

type migration_entry struct {
	id         int64
	version_id int64
	applied_at time.Time
}

const (
	migrateSetupQuery = `CREATE TABLE IF NOT EXISTS schema_migrations (
	 id INTEGER PRIMARY KEY AUTOINCREMENT
	,version_id INTEGER NOT NULL
	,applied_at DATETIME NOT NULL DEFAULT (datetime('now'))
);`
)

func (m *Migrator) ensureMigrationTable(ctx context.Context) (err error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin: %v", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, migrateSetupQuery); err != nil {
		return fmt.Errorf("create table: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %v", err)
	}

	return nil
}

func (m *Migrator) Up(ctx context.Context, to string) error {
	if err := m.ensureMigrationTable(ctx); err != nil {
		return fmt.Errorf("migrations table: %v", err)
	}

	sources, err := m.sources()
	if err != nil {
		return fmt.Errorf("load migration sources: %v", err)
	}
	if len(sources) == 0 {
		return nil
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("transaction start failed: %v", err)
	}
	defer tx.Rollback()

	var ids []int
	for id := range sources {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	last := -1

	var row migration_entry
	if err := tx.QueryRowContext(ctx, "SELECT * FROM schema_migrations ORDER BY version_id DESC LIMIT 1").Scan(
		&row.id, &row.version_id, &row.applied_at,
	); err != nil {
		if err != sql.ErrNoRows {
			return fmt.Errorf("querying migrations failed: %v", err)
		}
	} else {
		last = int(row.version_id)
	}

	for _, id := range ids {
		if id > last {
			source := sources[id]
			if _, err := tx.ExecContext(ctx, source.up); err != nil {
				return fmt.Errorf("applying %q: %v", source.name, err)
			}
			if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version_id) VALUES (?)", id); err != nil {
				return fmt.Errorf("recording %q: %v", source.name, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %v", err)
	}

	return nil
}
