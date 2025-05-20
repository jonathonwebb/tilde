package core

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type Config struct {
	Env    string
	Level  slog.Level
	Format logFormat
	Serve  struct {
		Addr string
		Dev  bool
	}
	DB struct {
		URI           string
		Schema        string
		MigrationName string
		Migrate       bool
	}
}

func (c Config) handlerOpts() *slog.HandlerOptions {
	return &slog.HandlerOptions{Level: c.Level}
}

func (c Config) handler(w io.Writer) slog.Handler {
	switch c.Format {
	case JSONFormat:
		return slog.NewJSONHandler(w, c.handlerOpts())
	default:
		return slog.NewTextHandler(w, c.handlerOpts())
	}
}

func (c Config) NewLogger(w io.Writer) *slog.Logger {
	return slog.New(c.handler(w))
}

type logFormat string

var (
	TextFormat logFormat = "TEXT"
	JSONFormat logFormat = "JSON"
)

func (l *logFormat) MarshalText() ([]byte, error) {
	return []byte(*l), nil
}

func (l *logFormat) UnmarshalText(text []byte) error {
	lower := strings.ToLower(string(text))
	switch lower {
	case "text":
		*l = TextFormat
	case "json":
		*l = JSONFormat
	default:
		return fmt.Errorf("expected one of: text, json")
	}
	return nil
}
