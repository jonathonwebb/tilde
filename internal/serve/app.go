package serve

import (
	"context"
	"log/slog"

	"github.com/jonathonwebb/tilde/internal/conf"
)

type App struct {
	logger *slog.Logger
	config conf.Config
}

func NewApp(logger *slog.Logger, config conf.Config) *App {
	return &App{logger: logger, config: config}
}

func (s *App) Run(ctx context.Context) error {
	s.logger.Info("serve app", "config", s.config)
	return nil
}
