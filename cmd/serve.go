package cmd

import (
	"context"
	"flag"

	"github.com/jonathonwebb/tilde/internal/serveapp"
	"github.com/jonathonwebb/x/conf"
)

var serveCmd = conf.Command{
	Name:  "serve",
	Usage: `usage: tilde [root flags] serve [-h | -help] [flags]`,
	Help: `usage: tilde [root flags] serve [-h | -help] [flags]

Serve starts the tilde application server.

flags:
  -listen=:0   listener address
  -dev         enable dev server
  -h, --help   show this help and exit`,
	Flags: func(fs *flag.FlagSet) {
		fs.StringVar(&cfg.Serve.Addr, "listen", ":0", "")
		fs.BoolVar(&cfg.Serve.Dev, "dev", false, "")
	},
	Vars: map[string]string{
		"listen": "TLD_LISTEN",
	},
	Action: func(ctx context.Context, e *conf.Env) conf.ExitStatus {
		l := cfg.NewLogger(e.Stderr)
		l.Debug("serve", "config", cfg)
		if err := serveapp.New(cfg, l).Run(ctx); err != nil {
			l.Error(err.Error())
			return conf.ExitFailure
		}
		return conf.ExitSuccess
	},
}
