package drop

// var Cmd = conf.Command{
// 	Name:  "drop",
// 	Usage: `usage: tilde [root flags] db drop [-h | -help]`,
// 	Help: `usage: tilde [root flags] db drop [-h | -help]

// Drops the database.

// flags:
//   -h, --help         show this help and exit`,
// 	Action: func(ctx context.Context, e *conf.Env, target any) conf.ExitStatus {
// 		// cfg := target.(*core.Config)
// 		// log := cfg.NewLogger(e.Stderr)
// 		// log.Debug("db drop", "cfg", cfg)

// 		// m, err := core.NewMigrator(cfg.DB.URI, log)
// 		// if err != nil {
// 		// 	log.Error(err.Error())
// 		// 	return conf.ExitFailure
// 		// }
// 		// err = m.Drop()
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
