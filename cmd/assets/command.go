package assets

import (
	"context"

	"github.com/jonathonwebb/tilde/internal/cli"
	"github.com/jonathonwebb/tilde/internal/core"
)

const (
	usage = "usage: tilde [root flags] assets"
	help  = `usage: tilde [root flags] assets

compile frontend assets.

flags:
  -h, -help   show this help and exit`
)

var Cmd = cli.Command{
	Name:  "assets",
	Usage: usage,
	Help:  help,
	Action: func(ctx context.Context, e *cli.Env, target any) cli.ExitStatus {
		cfg := target.(*core.Config)
		if len(e.Args) != 0 {
			e.PrintUsageErr(usage, "expected 0 args, but got %d", len(e.Args))
			return cli.ExitUsageError
		}
		if err := run(ctx, e.Stderr, cfg); err != nil {
			return cli.ExitFailure
		}
		return cli.ExitSuccess
	},
}
