package cmd

import (
	"context"
	"log/slog"

	"github.com/jonathonwebb/tilde/internal/cli"
	"github.com/jonathonwebb/tilde/internal/serve"
)

const (
	serveUsage = `usage: tilde [root flags] serve [-h] [flags]`
	serveHelp  = `usage: tilde [root flags] serve [-h] [flags]

Serve starts the application server.

flags:
  -h, --help   show this help and exit`
)

var serveCmd = cli.Command{
	Name:  "serve",
	Usage: serveUsage,
	Help:  serveHelp,
	Action: func(ctx context.Context, e *cli.Env) cli.ExitStatus {
		logger := slog.New(slog.NewTextHandler(e.Stderr, &slog.HandlerOptions{Level: cfg.Level}))
		if err := serve.NewApp(logger, cfg).Run(ctx); err != nil {
			logger.Error(err.Error())
			return cli.ExitFailure
		}
		return cli.ExitSuccess
	},
}
