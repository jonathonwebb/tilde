package cmd

import (
	"flag"
	"log/slog"

	"github.com/jonathonwebb/tilde/internal/core"
	"github.com/jonathonwebb/x/conf"
)

var cfg core.Config

var CLI = conf.Command{
	Name:  "tilde",
	Usage: `usage: tilde [-h | -help] [flags] <command>`,
	Help: `usage: tilde [-h | -help] [flags] <command>

Utilities for managing the tilde application.

commands:
  assets    compile UI assets
  serve     start app server
  version   print build info

flags:
  -env=production   app env id ($TLD_ENV)
  -format=text      log format (text|json) ($TLD_FORMAT)
  -level=info       log level (debug|info|warn|error) ($TLD_LEVEL)
  -h, -help         show this help and exit`,
	Flags: func(fs *flag.FlagSet) {
		fs.StringVar(&cfg.Env, "env", "production", "")
		fs.TextVar(&cfg.Level, "level", slog.LevelInfo, "")
		fs.TextVar(&cfg.Format, "format", &core.TextFormat, "")
	},
	Vars: map[string]string{
		"env":    "TLD_ENV",
		"level":  "TLD_LEVEL",
		"format": "TLD_FORMAT",
	},
	Commands: []*conf.Command{&assetsCmd, &serveCmd, &versionCmd},
}
