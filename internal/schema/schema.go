package schema

const (
	SchemaVersion = 2
	Schema        = `CREATE TABLE schema_lock (id INTEGER PRIMARY KEY);
CREATE TABLE schema_migrations (id INTEGER PRIMARY KEY, version_id INTEGER UNIQUE NOT NULL, applied_at DATETIME NOT NULL DEFAULT (datetime('now')));
CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT UNIQUE NOT NULL);
CREATE TABLE orgs (id INTEGER PRIMARY KEY, name TEXT UNIQUE NOT NULL);`
)
