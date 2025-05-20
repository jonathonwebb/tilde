package assets

import (
	"context"
	"log/slog"

	"github.com/jonathonwebb/tilde/internal/core"
	"github.com/jonathonwebb/x/conf"
)

var Cmd = conf.Command{
	Name:  "assets",
	Usage: "usage: tilde [root flags] assets [-h | -help] [flags]",
	Help: `usage: tilde [root flags] assets [-h | -help] [flags]

Assets compiles UI assets.

flags:
  -h, --help   show this help and exit`,
	Action: func(ctx context.Context, e *conf.Env, target any) conf.ExitStatus {
		cfg := target.(*core.Config)
		app := application{*cfg, cfg.NewLogger(e.Stderr)}
		if err := app.run(); err != nil {
			app.log.Error(err.Error())
			return conf.ExitFailure
		}
		return conf.ExitSuccess
	},
}

type application struct {
	cfg core.Config
	log *slog.Logger
}

func (a *application) run() error {
	return nil
}
