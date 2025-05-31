package schema

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	sqlite3 "github.com/mattn/go-sqlite3"
)

type Sqlite3SchemaStore struct {
	instance *sql.DB
	log      *slog.Logger
}

var _ SchemaStore = (*Sqlite3SchemaStore)(nil)

func NewSqlite3SchemaStore(db *sql.DB, log *slog.Logger) *Sqlite3SchemaStore {
	return &Sqlite3SchemaStore{db, log}
}

func (s *Sqlite3SchemaStore) db() *sql.DB {
	return s.instance
}

func (s *Sqlite3SchemaStore) init(ctx context.Context) error {
	if err := s.withTx(ctx, func(tCtx context.Context, tx *sql.Tx) error {
		if _, err := s.instance.ExecContext(tCtx, "CREATE TABLE IF NOT EXISTS schema_lock (id INTEGER PRIMARY KEY)"); err != nil {
			return err
		}
		if _, err := s.instance.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS schema_migrations (id INTEGER PRIMARY KEY, version_id INTEGER UNIQUE NOT NULL, applied_at DATETIME NOT NULL DEFAULT (datetime('now')))"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	s.log.Debug("schema version tables present")
	return nil
}

func (s *Sqlite3SchemaStore) lock(ctx context.Context, pollInterval time.Duration) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			_, err := s.instance.ExecContext(ctx, "INSERT INTO schema_lock (id) VALUES (1)")
			if err == nil {
				s.log.Info("obtained schema write lock")
				return nil
			}

			var sqliteErr sqlite3.Error
			if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint {
				s.log.Info("schema locked for writing")
				continue
			}
			return err
		}
	}
}

func (s *Sqlite3SchemaStore) release(ctx context.Context) error {
	_, err := s.db().ExecContext(ctx, "DELETE FROM schema_lock WHERE id = 1;")
	if err != nil {
		return err
	}
	s.log.Debug("released schema write lock")
	return nil
}

func (s *Sqlite3SchemaStore) state(ctx context.Context) (ids []int64, err error) {
	rows, err := s.db().QueryContext(ctx, `SELECT version_id FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer func() { err = errors.Join(rows.Close()) }()

	for rows.Next() {
		var vid int64
		err = rows.Scan(&vid)
		if err != nil {
			return nil, err
		}
		ids = append(ids, vid)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	s.log.Debug("read migration state", "n", len(ids))
	for i, vid := range ids {
		if i > 0 && vid == ids[i-1] {
			return nil, fmt.Errorf("duplicate version record %d", vid)
		}
		if i > 0 && vid < ids[i-1] {
			return nil, fmt.Errorf("version order mismatch, %d precedes %d", ids[i-1], vid)
		}
	}

	return ids, nil
}

func (s *Sqlite3SchemaStore) commit(ctx context.Context, id int64) error {
	if _, err := s.db().ExecContext(ctx, "INSERT INTO schema_migrations (version_id) VALUES (?)", id); err != nil {
		return err
	}
	s.log.Debug("commit migration", "id", id)
	return nil
}

func (s *Sqlite3SchemaStore) revert(ctx context.Context, id int64) error {
	if _, err := s.db().ExecContext(ctx, "DELETE FROM schema_migrations WHERE version_id >= ?", id); err != nil {
		return err
	}
	s.log.Debug("revert migration", "id", id)
	return nil
}

func (s *Sqlite3SchemaStore) dump(ctx context.Context, w io.Writer) (err error) {
	var stmts []string
	rows, err := s.db().QueryContext(ctx, "SELECT sql FROM sqlite_schema")
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(rows.Close()) }()

	for rows.Next() {
		var stmt sql.NullString
		if err := rows.Scan(&stmt); err != nil {
			return err
		}
		if stmt.Valid {
			el := stmt.String
			if !strings.HasSuffix(el, ";") {
				el += ";"
			}
			stmts = append(stmts, el)
		}
	}

	_, err = w.Write([]byte(strings.Join(stmts, "\n")))
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlite3SchemaStore) load(ctx context.Context, r io.Reader) error {
	stmts := strings.SplitSeq(Schema, ";")
	for stmt := range stmts {
		if _, err := s.db().ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

func (s *Sqlite3SchemaStore) close() error {
	return s.instance.Close()
}

func (s *Sqlite3SchemaStore) withTx(ctx context.Context, fn func(context.Context, *sql.Tx) error) (err error) {
	tx, err := s.db().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, tx.Rollback())
	}()

	err = fn(ctx, tx)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
