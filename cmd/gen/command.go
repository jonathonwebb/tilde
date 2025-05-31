package gen

import (
	"github.com/jonathonwebb/tilde/cmd/gen/migration"
	"github.com/jonathonwebb/tilde/internal/cli"
)

var Cmd = cli.Command{
	Name:  "gen",
	Usage: "usage: tilde [root flags] gen <command>",
	Help: `usage: tilde [root flags] gen <command>

development template generators.

commands:
  migration   generate a database migration

flags:
  -h, -help   show this help and exit`,
	Commands: []*cli.Command{
		&migration.Cmd,
	},
}
