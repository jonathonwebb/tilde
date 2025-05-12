package conf

import "log/slog"

type Config struct {
	Level slog.Level
	Env   string
}
