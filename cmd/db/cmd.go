package db

import (
	dbNew "github.com/jonathonwebb/tilde/cmd/db/new"
	dbUp "github.com/jonathonwebb/tilde/cmd/db/up"
	"github.com/jonathonwebb/x/conf"
)

var Cmd = conf.Command{
	Name:  "db",
	Usage: `usage: tilde [root flags] db [-h | -help] [flags] <command>`,
	Help: `usage: tilde [root flags] db [-h | -help] [flags] <command>

Utilities for managing the tilde database.

commands:
  down      reverse applied migrations
  drop      drop the database
  dump      export schema file
  load      apply schema file
  new       generate new migration
  up        apply pending migrations

flags:
  -h, --help   show this help and exit`,
	Commands: []*conf.Command{
		&dbNew.Cmd,
		&dbUp.Cmd,
		// &down.Cmd,
		// &drop.Cmd,
		// &dump.Cmd,
		// &load.Cmd,
		// &migrate.Cmd,
		// &new.Cmd,
		// &status.Cmd,
		// &up.Cmd,
	},
}
