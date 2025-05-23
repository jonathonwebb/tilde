package new

import (
	"context"
	"fmt"

	"github.com/jonathonwebb/tilde/internal/core"
	"github.com/jonathonwebb/x/conf"
)

const (
	usage = `usage: tilde [root flags] db new [-h | -help] <name>`
	help  = `usage: tilde [root flags] db new [-h | -help] <name>

New generates a new migration file.

flags:
  -h, --help      show this help and exit`
)

var Cmd = conf.Command{
	Name:  "new",
	Usage: usage,
	Help:  help,
	Action: func(ctx context.Context, e *conf.Env, target any) conf.ExitStatus {
		if len(e.Args) == 0 {
			fmt.Fprintf(e.Stderr, "missing required argument: <name>\n%s\n", usage)
			return conf.ExitUsageError
		}
		if len(e.Args) > 1 {
			fmt.Fprintf(e.Stderr, "invalid argument count: %d\n%s\n", len(e.Args), usage)
			return conf.ExitUsageError
		}

		cfg := target.(*core.Config)
		log := cfg.NewLogger(e.Stderr)
		log.Debug("db new", "cfg", cfg)

		m := core.NewMigrator(nil, "db/migrations", log)
		if err := m.New(e.Args[0]); err != nil {
			log.Error(err.Error())
			return conf.ExitFailure
		}

		return conf.ExitSuccess
	},
}
