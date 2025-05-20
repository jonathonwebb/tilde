package serve

import (
	"context"
	"flag"

	"github.com/jonathonwebb/tilde/internal/core"
	"github.com/jonathonwebb/x/conf"
)

var Cmd = conf.Command{
	Name:  "serve",
	Usage: `usage: tilde [root flags] serve [-h | -help] [flags]`,
	Help: `usage: tilde [root flags] serve [-h | -help] [flags]

Serve starts the tilde application server.

flags:
  -listen=:0   listener address
  -dev         enable dev server
  -h, --help   show this help and exit`,
	Flags: func(fs *flag.FlagSet, cfg any) {
		if cfg, ok := cfg.(*core.Config); ok {
			fs.StringVar(&cfg.Serve.Addr, "listen", ":0", "")
			fs.BoolVar(&cfg.Serve.Dev, "dev", false, "")
		}
	},
	Vars: map[string]string{
		"listen": "TLD_LISTEN",
	},
	Action: func(ctx context.Context, e *conf.Env, target any) conf.ExitStatus {
		cfg := target.(*core.Config)
		if err := run(ctx, e.Stderr, *cfg); err != nil {
			return conf.ExitFailure
		}
		return conf.ExitSuccess
	},
}
