package status

import (
	"context"
	"errors"
	"fmt"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/sqlite"
	"github.com/jonathonwebb/tilde/internal/core"
	"github.com/jonathonwebb/x/conf"
)

var Cmd = conf.Command{
	Name:  "status",
	Usage: `usage: tilde [root flags] db status [-h | -help]`,
	Help: `usage: tilde [root flags] db status [-h | -help]

Status lists pending migrations.

flags:
  -h, --help         show this help and exit`,
	Action: func(ctx context.Context, e *conf.Env, target any) conf.ExitStatus {
		cfg := target.(*core.Config)
		log := cfg.NewLogger(e.Stderr)
		log.Debug("db status", "cfg", cfg)

		m, err := core.NewMigrator(cfg.DB.URI)
		if err != nil {
			log.Error(err.Error())
			return conf.ExitFailure
		}
		_, err = m.Status(true)
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
