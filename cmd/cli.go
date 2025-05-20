package cmd

import (
	"flag"
	"log/slog"

	"github.com/jonathonwebb/tilde/cmd/assets"
	"github.com/jonathonwebb/tilde/cmd/db"
	"github.com/jonathonwebb/tilde/cmd/serve"
	"github.com/jonathonwebb/tilde/cmd/version"
	"github.com/jonathonwebb/tilde/internal/core"
	"github.com/jonathonwebb/x/conf"
)

var CLI = conf.Command{
	Name:  "tilde",
	Usage: `usage: tilde [-h | -help] [flags] <command>`,
	Help: `usage: tilde [-h | -help] [flags] <command>

Utilities for managing the tilde application.

commands:
  assets    compile UI assets
  db        database utils
  serve     start app server
  version   print build info

flags:
  -env=production              app env id ($TLD_ENV)
  -format=text                 log format (text|json) ($TLD_FORMAT)
  -level=info                  log level (debug|info|warn|error) ($TLD_LEVEL)
  -db=sqlite:data/default.db   sqlite database URI ($TLD_DB)
  -migratedir=db/migrations    migration file dir ($TLD_MIGRATE_DIR)
  -h, -help                    show this help and exit`,
	Flags: func(fs *flag.FlagSet, target any) {
		cfg := target.(*core.Config)
		fs.StringVar(&cfg.Env, "env", "production", "")
		fs.TextVar(&cfg.Level, "level", slog.LevelInfo, "")
		fs.TextVar(&cfg.Format, "format", &core.TextFormat, "")
		fs.StringVar(&cfg.DB.URI, "db", "sqlite:data.db", "")
	},
	Vars: map[string]string{
		"env":         "TLD_ENV",
		"level":       "TLD_LEVEL",
		"format":      "TLD_FORMAT",
		"db":          "TLD_DB",
		"migrateddir": "TLD_MIGRATE_DIR",
	},
	Commands: []*conf.Command{&assets.Cmd, &db.Cmd, &serve.Cmd, &version.Cmd},
}
