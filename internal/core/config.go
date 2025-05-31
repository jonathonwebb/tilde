package core

import (
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
)

type Config struct {
	Env          string
	Level        slog.Level
	Format       LogFormat
	DbConnString string

	// assets
	AssetsDir string
	StaticDir string

	// serve
	ServeAddr string
	ServeDev  bool

	// migrate
	DbSchemaVersion SchemaVersion
	DbMigrateSkip   bool
}

func (c Config) LogParams() []any {
	return []any{
		"Env", c.Env,
		"Level", c.Level,
		"Format", c.Format,
		"ServeAddr", c.Format,
		"ServeDev", c.ServeDev,
		"DbString", c.DbConnString,
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

func (c Config) NewLogger(w io.Writer, app string) *slog.Logger {
	return slog.New(c.handler(w)).With(slog.Group("app", "name", app, "env", c.Env))
}

type LogFormat string

var (
	TextFormat LogFormat = "TEXT"
	JSONFormat LogFormat = "JSON"
)

func (l *LogFormat) MarshalText() ([]byte, error) {
	return []byte(*l), nil
}

func (l *LogFormat) UnmarshalText(text []byte) error {
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

type SchemaVersion int64

var (
	SchemaInitial SchemaVersion = -1
	SchemaLatest  SchemaVersion = -2
	SchemaFile    SchemaVersion = -3
)

func (v *SchemaVersion) MarshalText() ([]byte, error) {
	return []byte{}, nil
}

func (v *SchemaVersion) UnmarshalText(text []byte) error {
	lower := strings.ToLower(string(text))
	switch lower {
	case "latest":
		*v = SchemaLatest
	case "initial":
		*v = SchemaInitial
	case "schema":
		*v = SchemaFile
	default:
		n, err := strconv.ParseInt(string(text), 10, 0)
		if err != nil {
			return err
		}
		*v = SchemaVersion(n)
	}

	return nil
}
