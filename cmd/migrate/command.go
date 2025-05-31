package migrate

import (
	"context"
	"flag"

	"github.com/jonathonwebb/tilde/internal/cli"
	"github.com/jonathonwebb/tilde/internal/core"
)

var Cmd = cli.Command{
	Name:  "migrate",
	Usage: "usage: tilde [root flags] migrate [-h] [flags]",
	Help: `usage: tilde [root flags] migrate [-h] [flags]

update the database schema to the latest version.

flags:
  -skip        initialize without migrating
  -to=latest   version target (initial|latest|schema|uint64)
  -h, -help    show this help and exit`,
	Flags: func(fs *flag.FlagSet, target any) {
		cfg := target.(*core.Config)
		fs.TextVar(&cfg.DbSchemaVersion, "to", &core.SchemaLatest, "latest")
		fs.BoolVar(&cfg.DbMigrateSkip, "skip", false, "")
	},
	Action: func(ctx context.Context, e *cli.Env, target any) cli.ExitStatus {
		cfg := target.(*core.Config)

		if err := run(ctx, e.Stderr, cfg); err != nil {
			return cli.ExitFailure
		}

		return cli.ExitSuccess
	},
}
