package migrations

import (
	"context"
	"database/sql"
	"log/slog"
)

var _001_CreateUsers = migration{
	name: "001_CreateUsers",
	up: func(ctx context.Context, db *sql.DB, log *slog.Logger) error {
		return withTx(ctx, db, log, []string{
			"CREATE TABLE IF NOT EXISTS users (" +
				"id INTEGER PRIMARY KEY" +
				",username TEXT NOT NULL" +
				",created_at DATETIME NOT NULL DEFAULT (datetime('now'))" +
				",updated_at DATETIME NOT NULL DEFAULT (datetime('now'))" +
				");",
		})
	},
	down: func(ctx context.Context, db *sql.DB, log *slog.Logger) error {
		return withTx(ctx, db, log, []string{
			"DROP TABLE IF EXISTS users;",
		})
	},
}
