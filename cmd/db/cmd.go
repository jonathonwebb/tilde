package db

import (
	"github.com/jonathonwebb/tilde/cmd/db/down"
	"github.com/jonathonwebb/tilde/cmd/db/drop"
	"github.com/jonathonwebb/tilde/cmd/db/dump"
	"github.com/jonathonwebb/tilde/cmd/db/load"
	"github.com/jonathonwebb/tilde/cmd/db/migrate"
	"github.com/jonathonwebb/tilde/cmd/db/new"
	"github.com/jonathonwebb/tilde/cmd/db/status"
	"github.com/jonathonwebb/tilde/cmd/db/up"
	"github.com/jonathonwebb/x/conf"
)

var Cmd = conf.Command{
	Name:  "db",
	Usage: `usage: tilde [root flags] db [-h | -help] [flags] <command>`,
	Help: `usage: tilde [root flags] db [-h | -help] [flags] <command>

Utilities for managing the tilde database.

commands:
  down      reverse most recent migration
  drop      drop the database
  dump      export schema file
  load      apply schema file
  migrate   apply pending migrations
  new       generate new migration
  status    list pending migrations
  up        create the database

flags:
  -h, --help   show this help and exit`,
	Commands: []*conf.Command{
		&down.Cmd,
		&drop.Cmd,
		&dump.Cmd,
		&load.Cmd,
		&migrate.Cmd,
		&new.Cmd,
		&status.Cmd,
		&up.Cmd,
	},
}
