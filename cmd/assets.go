package cmd

import (
	"context"

	"github.com/jonathonwebb/tilde/internal/assetsapp"
	"github.com/jonathonwebb/x/conf"
)

var assetsCmd = conf.Command{
	Name:  "assets",
	Usage: "usage: tilde [root flags] assets [-h | -help] [flags]",
	Help: `usage: tilde [root flags] assets [-h | -help] [flags]

Assets compiles UI assets.

flags:
  -h, --help   show this help and exit`,
	Action: func(ctx context.Context, e *conf.Env) conf.ExitStatus {
		l := cfg.NewLogger(e.Stderr)
		l.Debug("assets", "config", cfg)
		if err := assetsapp.New(cfg, l).Run(ctx); err != nil {
			l.Error(err.Error())
			return conf.ExitFailure
		}
		return conf.ExitSuccess
	},
}
