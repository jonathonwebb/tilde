package serve

import (
	"context"
	"flag"

	"github.com/jonathonwebb/tilde/internal/cli"
	"github.com/jonathonwebb/tilde/internal/core"
)

var Cmd = cli.Command{
	Name:  "serve",
	Usage: `usage: tilde [root flags] serve [-h | -help] [flags]`,
	Help: `usage: tilde [root flags] serve [-h | -help] [flags]

starts the tilde application server.

flags:
  -addr=:0     listener address ($TLD_ADDR)
  -dev         enable dev server
  -h, --help   show this help and exit`,
	Flags: func(fs *flag.FlagSet, cfg any) {
		if cfg, ok := cfg.(*core.Config); ok {
			fs.StringVar(&cfg.ServeAddr, "addr", ":0", "")
			fs.BoolVar(&cfg.ServeDev, "dev", false, "")
		}
	},
	Vars: map[string]string{
		"addr": "TLD_ADDR",
	},
	Action: func(ctx context.Context, e *cli.Env, target any) cli.ExitStatus {
		cfg := target.(*core.Config)
		if err := run(ctx, e.Stderr, cfg); err != nil {
			return cli.ExitFailure
		}
		return cli.ExitSuccess
	},
}
