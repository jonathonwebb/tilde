package root

import (
	"flag"
	"log/slog"

	"github.com/jonathonwebb/tilde/cmd/assets"
	"github.com/jonathonwebb/tilde/cmd/gen"
	"github.com/jonathonwebb/tilde/cmd/migrate"
	"github.com/jonathonwebb/tilde/cmd/serve"
	"github.com/jonathonwebb/tilde/cmd/version"
	"github.com/jonathonwebb/tilde/internal/cli"
	"github.com/jonathonwebb/tilde/internal/core"
)

var Cmd = cli.Command{
	Name:  "tilde",
	Usage: "usage: tilde [-h] [flags] <command>",
	Help: `usage: tilde [-h] [flags] <command>

utils for managing the tilde application.

commands:
  assets    compile frontend assets
  gen       generate dev templates
  migrate   update database schema
  serve     start app server
  version   print version info

flags:
  -assets=ui/assets   assets src dir ($TLD_ASSETS)
  -db=data.db         sqlite connection string ($TLD_DB)
  -env=production     app env id ($TLD_ENV)
  -format=text        log format (text|json) ($TLD_FMT)
  -level=info         log level (debug|info|warn|error) ($TLD_LVL)
  -public=ui/static   public asset dir ($TLD_PUBLIC)
  -h, -help           show this help and exit`,
	Flags: func(fs *flag.FlagSet, target any) {
		cfg := target.(*core.Config)
		fs.StringVar(&cfg.AssetsDir, "assets", "ui/assets", "")
		fs.StringVar(&cfg.DbConnString, "db", "data.db", "")
		fs.StringVar(&cfg.Env, "env", "production", "")
		fs.TextVar(&cfg.Level, "level", slog.LevelInfo, "")
		fs.TextVar(&cfg.Format, "format", &core.TextFormat, "")
		fs.StringVar(&cfg.StaticDir, "public", "ui/static", "")
	},
	Vars: map[string]string{
		"assets": "TLD_ASSETS",
		"db":     "TLD_DB",
		"env":    "TLD_ENV",
		"format": "TLD_FMT",
		"level":  "TLD_LVL",
		"public": "TLD_PUBLIC",
	},
	Commands: []*cli.Command{&assets.Cmd, &gen.Cmd, &migrate.Cmd, &serve.Cmd, &version.Cmd},
}
