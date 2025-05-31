package migration

import (
	"context"
	"fmt"
	"time"

	"github.com/jonathonwebb/tilde/internal/cli"
	"github.com/jonathonwebb/tilde/internal/schema"
)

const (
	usage = "usage: tilde [root flags] gen migration <name>"
	help  = `usage: tilde [root flags] gen migration <name>

generate a new migration template with <name>.

flags:
  -h, -help   show this help and exit`
)

var Cmd = cli.Command{
	Name:  "migration",
	Usage: usage,
	Help:  help,
	Action: func(ctx context.Context, e *cli.Env, target any) cli.ExitStatus {
		args := e.Args
		if len(args) != 1 {
			fmt.Fprintf(e.Stderr, "expected 1 <name> arg, got %d\n", len(args))
			fmt.Fprintf(e.Stderr, "%s\n", usage)
			return cli.ExitUsageError
		}

		ts := time.Now().UTC()
		if err := schema.NewMigration(ctx, "internal/migrations", args[0], ts); err != nil {
			fmt.Fprintf(e.Stderr, "failed to generate migration: %v", err)
			return cli.ExitFailure
		}
		return cli.ExitSuccess
	},
}
