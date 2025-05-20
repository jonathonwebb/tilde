package new

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
	Name:  "new",
	Usage: `usage: tilde [root flags] db new [-h | -help] [flags]`,
	Help: `usage: tilde [root flags] db new [-h | -help] [flags]

New generates a new migration file.

flags:
  -name=unnamed   migration name
  -h, --help      show this help and exit`,
	Flags: func(fs *flag.FlagSet, target any) {
		cfg := target.(*core.Config)
		fs.StringVar(&cfg.DB.MigrationName, "name", "unnamed", "")
	},
	Action: func(ctx context.Context, e *conf.Env, target any) conf.ExitStatus {
		cfg := target.(*core.Config)
		log := cfg.NewLogger(e.Stderr)
		log.Debug("db new", "cfg", cfg)

		m, err := core.NewMigrator(cfg.DB.URI)
		if err != nil {
			log.Error(err.Error())
			return conf.ExitFailure
		}

		err = m.NewMigration(cfg.DB.MigrationName)
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
