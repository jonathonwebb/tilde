package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jonathonwebb/tilde/internal/cli"
	"github.com/jonathonwebb/tilde/internal/conf"
	"github.com/spf13/pflag"
)

var (
	levelFlag levelValue
	jsonFlag  bool

	cfg conf.Config
)

const (
	rootUsage = `usage: tilde [-h] [flags] command`
	rootHelp  = `usage: tilde [-h] [flags] command

Utilities for managing the tilde application.

commands:
  serve     start app server
  version   print build info

flags:
  -C, --config path   config file path [$TLD_CONFIG]
  -e, --env id        app env id (default: production) [$TLD_ENV]
  -l, --level name    log level name (default: info) [$TLD_LEVEL]
  -j, --json          format logs as json [$TLD_JSON]
  -h, --help          show this help and exit`
)

var rootCmd = cli.Command{
	Name:  "tilde",
	Usage: rootUsage,
	Help:  rootHelp,
}

func Execute(ctx context.Context, env *cli.Env) cli.ExitStatus {
	rootCmd.SetFlags(func(fs *pflag.FlagSet) {
		fs.VarP(&levelFlag, "level", "l", "log level `name`")
		rootCmd.MapEnvVar("level", "TLD_LEVEL")
		rootCmd.MapConfigKey("level", "level")

		fs.StringP("config", "C", "", "config file `path`")
		rootCmd.MapEnvVar("config", "TLD_CONFIG")
		rootCmd.SetConfigFlag("config")

		fs.StringVarP(&cfg.Env, "env", "e", "production", "app env `id`")
		rootCmd.MapEnvVar("env", "TLD_ENV")
		rootCmd.MapConfigKey("env", "env")
	})

	rootCmd.AddCommand(&serveCmd)
	rootCmd.AddCommand(&versionCmd)

	return rootCmd.Execute(ctx, env)
}

type levelValue slog.Level

func (lv *levelValue) String() string {
	return slog.Level(*lv).String()
}

func (lv *levelValue) Set(val string) error {
	lower := strings.ToLower(val)
	switch lower {
	case "debug":
		*lv = levelValue(slog.LevelDebug)
		return nil
	case "info":
		*lv = levelValue(slog.LevelInfo)
		return nil
	case "warn":
		*lv = levelValue(slog.LevelWarn)
		return nil
	case "error":
		*lv = levelValue(slog.LevelError)
		return nil
	default:
		return fmt.Errorf("unknown level %q, expected one of: debug, info, warn, error", lower)
	}
}

func (lv *levelValue) Type() string {
	return "level"
}
