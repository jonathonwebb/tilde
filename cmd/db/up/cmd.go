package up

import (
	"context"
	"database/sql"
	"flag"
	"time"

	"github.com/jonathonwebb/tilde/internal/core"
	"github.com/jonathonwebb/x/conf"
)

var Cmd = conf.Command{
	Name:  "up",
	Usage: `usage: tilde [root flags] db up [-h | -help] [flags]`,
	Help: `usage: tilde [root flags] db up [-h | -help] [flags]

Up applies pending migrations.

flags:
  -to          target migration id
  -h, --help   show this help and exit`,
	Flags: func(fs *flag.FlagSet, target any) {
		cfg := target.(*core.Config)
		fs.StringVar(&cfg.DB.MigrateTo, "to", "", "")
		fs.DurationVar(&cfg.DB.MigrateTimeout, "timeout", 30*time.Second, "")
	},
	Action: func(ctx context.Context, e *conf.Env, target any) conf.ExitStatus {
		cfg := target.(*core.Config)
		log := cfg.NewLogger(e.Stderr)
		log.Debug("db up", "cfg", cfg)

		db, err := sql.Open("sqlite3", cfg.DB.URI)
		if err != nil {
			log.Error(err.Error())
			return conf.ExitFailure
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.Error(err.Error())
			}
		}()

		m := core.NewMigrator(db, "db/migrations", log)

		ctx, cancel := context.WithTimeout(ctx, cfg.DB.MigrateTimeout)
		defer cancel()

		if err := m.Up(ctx, cfg.DB.MigrateTo); err != nil {
			log.Error(err.Error())
			return conf.ExitFailure
		}

		// cfg := target.(*core.Config)
		// log := cfg.NewLogger(e.Stderr)
		// log.Debug("db up", "cfg", cfg)

		// m, err := core.NewMigrator(cfg.DB.URI, log)
		// if err != nil {
		// 	log.Error(err.Error())
		// 	return conf.ExitFailure
		// }
		// err = m.Create()
		// if err != nil {
		// 	if errors.Is(err, dbmate.ErrInvalidURL) {
		// 		err = fmt.Errorf("invalid database URI: %q", cfg.DB.URI)
		// 	}
		// 	log.Error(err.Error())
		// 	return conf.ExitFailure
		// }

		return conf.ExitSuccess
	},
}
