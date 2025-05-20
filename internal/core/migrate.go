package core

import (
	"io"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/sqlite"
)

func NewMigrator(db string) (*dbmate.DB, error) {
	u, err := url.Parse(db)
	if err != nil {
		return nil, err
	}
	m := dbmate.New(u)
	m.Strict = true
	m.Log = io.Discard
	m.SchemaFile = "db/schema.sql"
	m.MigrationsTableName = "schema_migrations"
	m.MigrationsDir = []string{"db/migrations"}
	return m, nil
}
