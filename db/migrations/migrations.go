package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

var migrations = []migration{
	_001_CreateUsers,
}

type migration struct {
	name string
	up   func(context.Context, *sql.DB, *slog.Logger) error
	down func(context.Context, *sql.DB, *slog.Logger) error
}

type migrator struct {
	db  *sql.DB
	log *slog.Logger
}

func (m *migrator) ensureVersionTables(ctx context.Context) (err error) {
	return withTx(ctx, m.db, m.log, []string{
		"CREATE TABLE IF NOT EXISTS schema_lock (" +
			"id INTEGER PRIMARY KEY" +
			");",
		"CREATE TABLE IF NOT EXISTS schema_migrations (" +
			"id INTEGER PRIMARY KEY" +
			",version_id INTEGER NOT NULL" +
			",applied_at DATETIME NOT NULL DEFAULT (datetime('now'))" +
			");",
	})
}

func (m *migrator) up(ctx context.Context) error {
	all, err := allMigrations()
	if err != nil {
		return err
	}

	if err := m.ensureVersionTables(ctx); err != nil {
		return err
	}

	if err := m.lock(ctx); err != nil {
		return err
	}
	defer m.release(ctx)

	latest, err := m.latest(ctx)

	for id, mi := range all {
		if id > latest || errors.Is(err, sql.ErrNoRows) {
			if err := mi.up(ctx, m.db, m.log); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *migrator) down(ctx context.Context) error {
	return nil
}

// returns sql.ErrNoRows if no migrations have been applied
func (m *migrator) latest(ctx context.Context) (int64, error) {
	var val int64
	if err := m.db.QueryRowContext(ctx, "SELECT version_id FROM schema_migrations ORDER BY version_id DESC LIMIT 1").Scan(&val); err != nil {
		return 0, err
	}
	return val, nil
}

func (m *migrator) lock(ctx context.Context) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// TODO: poll?
	var lock_id int64
	err = tx.QueryRowContext(ctx, "SELECT id FROM schema_lock").Scan(&lock_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("database is already locked for migrations")
		}
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO schema_lock DEFAULT VALUES;")
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *migrator) release(ctx context.Context) error {
	if _, err := m.db.ExecContext(ctx, "DELETE FROM schema_lock;"); err != nil {
		return err
	}
	return nil
}

func allMigrations() (map[int64]*migration, error) {
	all := make(map[int64]*migration)

	for _, m := range migrations {
		parts := strings.SplitN(m.name, "_", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid migration name %q: separator '_' is required", m.name)
		}
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid migration id %q: non-integer", m.name)
		}
		if id < 0 {
			return nil, fmt.Errorf("invalid migration id %q: negative", m.name)
		}
		if _, ok := all[int64(id)]; ok {
			return nil, fmt.Errorf("invalid migration id %q: duplicate", m.name)
		}
	}

	return all, nil
}

func withTx(ctx context.Context, db *sql.DB, log *slog.Logger, stmts []string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin: %v", err)
	}
	defer tx.Rollback()

	for _, stmt := range stmts {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("exec: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %v", err)
	}

	return nil
}
