package migrate

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/jonathonwebb/tilde/internal/core"
	"github.com/jonathonwebb/tilde/internal/migrations"
	"github.com/jonathonwebb/tilde/internal/schema"
)

func run(ctx context.Context, w io.Writer, cfg *core.Config) (err error) {
	log := cfg.NewLogger(w, "migrate")
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()

	db, err := sql.Open("sqlite3", cfg.DbConnString)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, db.Close())
	}()

	store := schema.NewSqlite3SchemaStore(db, log)
	m := &schema.Migrator{
		Store:   store,
		Log:     log,
		Sources: migrations.All,
	}
	// return m.Init(ctx)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	switch cfg.DbSchemaVersion {
	case core.SchemaInitial:
		err = m.ApplyInitial(ctx)
	case core.SchemaLatest:
		err = m.ApplyLatest(ctx)
	case core.SchemaFile:
		err = m.Load(ctx, strings.NewReader(schema.Schema))
	default:
		err = m.Apply(ctx, int64(cfg.DbSchemaVersion))
	}

	if err != nil {
		return err
	}

	// app := &application{
	// 	log:        log,
	// 	store:      schema.NewSqlite3SchemaStore(conn),
	// 	schema:     strings.NewReader(db.Schema),
	// 	migrations: db.Migrations,
	// }

	// if err := app.migrate(ctx, cfg.DbSchemaVersion); err != nil {
	// 	return err
	// }

	return nil
}

// type application struct {
// 	log        *slog.Logger
// 	store      schema.SchemaStore
// 	migrations map[uint64]schema.Migration
// 	schema     io.Reader
// }

// func (app *application) newMigrator() *schema.Migrator {
// 	return &schema.Migrator{Store: app.store, Sources: app.migrations}
// }

// func (app *application) migrate(ctx context.Context, v core.SchemaVersion) error {
// 	m := app.newMigrator()
// 	if err := m.Init(ctx); err != nil {
// 		return err
// 	}

// 	switch v {
// 	case core.SchemaInitial:
// 		return m.ApplyInitial(ctx)
// 	case core.SchemaLatest:
// 		return m.ApplyLatest(ctx)
// 	case core.SchemaFile:
// 		return m.Load(ctx, strings.NewReader(db.Schema))
// 	default:
// 		return m.Apply(ctx, int64(v))
// 	}
// }
