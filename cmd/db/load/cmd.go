package load

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/sqlite"
	"github.com/jonathonwebb/tilde/internal/core"
	"github.com/jonathonwebb/x/conf"
)

var Cmd = conf.Command{
	Name:  "load",
	Usage: `usage: tilde [root flags] db load [-h | -help] [flags]`,
	Help: `usage: tilde [root flags] db load [-h | -help] [flags]

Loads the schema file into the database.

flags:
  -f=db/schema.sql   schema file path
  -h, --help         show this help and exit`,
	Flags: func(fs *flag.FlagSet, target any) {
		cfg := target.(*core.Config)
		fs.StringVar(&cfg.DB.Schema, "f", "db/schema.sql", "")
	},
	Vars: map[string]string{
		"f": "TILDE_SCHEMA",
	},
	Action: func(ctx context.Context, e *conf.Env, target any) conf.ExitStatus {
		cfg := target.(*core.Config)
		log := cfg.NewLogger(e.Stderr)
		log.Debug("db load", "cfg", cfg)

		m, err := core.NewMigrator(cfg.DB.URI)
		if err != nil {
			log.Error(err.Error())
			return conf.ExitFailure
		}
		err = m.LoadSchema()
		if err != nil {
			if errors.Is(err, dbmate.ErrInvalidURL) {
				err = fmt.Errorf("invalid database URI: %q", cfg.DB.URI)
			}
			log.Error(err.Error())
			return conf.ExitFailure
		}

		return conf.ExitSuccess
	},
}
