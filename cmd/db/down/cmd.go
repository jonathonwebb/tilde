package down

// var Cmd = conf.Command{
// 	Name:  "down",
// 	Usage: `usage: tilde [root flags] db down [-h | -help]`,
// 	Help: `usage: tilde [root flags] db down [-h | -help]

// Down reverses the most recent migration.

// flags:
//   -h, --help         show this help and exit`,
// 	Action: func(ctx context.Context, e *conf.Env, target any) conf.ExitStatus {
// 		// cfg := target.(*core.Config)
// 		// log := cfg.NewLogger(e.Stderr)
// 		// log.Debug("db down", "cfg", cfg)

// 		// m, err := core.NewMigrator(cfg.DB.URI, log)
// 		// if err != nil {
// 		// 	log.Error(err.Error())
// 		// 	return conf.ExitFailure
// 		// }
// 		// err = m.Rollback()
// 		// if err != nil {
// 		// 	if errors.Is(err, dbmate.ErrInvalidURL) {
// 		// 		err = fmt.Errorf("invalid database URI: %q", cfg.DB.URI)
// 		// 	}
// 		// 	log.Error(err.Error())
// 		// 	return conf.ExitFailure
// 		// }

// 		return conf.ExitSuccess
// 	},
// }
