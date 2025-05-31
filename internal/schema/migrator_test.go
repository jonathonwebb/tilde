package schema_test

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jonathonwebb/tilde/internal/schema"
)

func TestGenerate(t *testing.T) {
	dir := t.TempDir()

	t.Run("without existing migrations", func(t *testing.T) {
		name := "one"
		ts := time.Unix(0, 0)

		schema.NewMigration(t.Context(), dir, name, ts)

		want := `package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jonathonwebb/tilde/internal/schema"
)

var _0000000000_one = schema.Migration{
	Id:   0,
	Desc: "",
	Up: func(ctx context.Context, db *sql.DB, log *slog.Logger) (err error) {
		return fmt.Errorf("up migration 0000000000_one not implemented")
	},
	Down: func(ctx context.Context, db *sql.DB, log *slog.Logger) (err error) {
		return fmt.Errorf("down migration 0000000000_one not implemented")
	},
}
`

		d, err := os.ReadFile(path.Join(dir, fmt.Sprintf("%010d_%s.go", ts.Unix(), name)))
		if err != nil {
			t.Fatal(err)
		}
		got := string(d)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("migration mismatch (-want +got):\n%s", diff)
		}

		want = `package migrations

import "github.com/jonathonwebb/tilde/internal/schema"

var All = []schema.Migration{
	_0000000000_one,
}
`
		d, err = os.ReadFile(path.Join(dir, "all.go"))
		if err != nil {
			t.Fatal(err)
		}
		got = string(d)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("migration index mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("with existing migrations", func(t *testing.T) {
		name := "two"
		ts := time.Unix(1748577600, 0)

		schema.NewMigration(t.Context(), dir, name, ts)

		want := `package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jonathonwebb/tilde/internal/schema"
)

var _1748577600_two = schema.Migration{
	Id:   1748577600,
	Desc: "",
	Up: func(ctx context.Context, db *sql.DB, log *slog.Logger) (err error) {
		return fmt.Errorf("up migration 1748577600_two not implemented")
	},
	Down: func(ctx context.Context, db *sql.DB, log *slog.Logger) (err error) {
		return fmt.Errorf("down migration 1748577600_two not implemented")
	},
}
`

		d, err := os.ReadFile(path.Join(dir, fmt.Sprintf("%010d_%s.go", ts.Unix(), name)))
		if err != nil {
			t.Fatal(err)
		}
		got := string(d)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("migration mismatch (-want +got):\n%s", diff)
		}

		want = `package migrations

import "github.com/jonathonwebb/tilde/internal/schema"

var All = []schema.Migration{
	_0000000000_one,
	_1748577600_two,
}
`
		d, err = os.ReadFile(path.Join(dir, "all.go"))
		if err != nil {
			t.Fatal(err)
		}
		got = string(d)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("migration index mismatch (-want +got):\n%s", diff)
		}

	})
}
