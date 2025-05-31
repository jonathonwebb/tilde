package serve

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/jonathonwebb/tilde/internal/core"
)

func run(ctx context.Context, w io.Writer, cfg *core.Config) (err error) {
	log := cfg.NewLogger(w, "serve")
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()

	app := &application{
		log: log,
	}

	log.Info("starting server", "addr", cfg.ServeAddr, "dev", cfg.ServeDev)
	return http.ListenAndServe(cfg.ServeAddr, app.handlers())
}

type application struct {
	log *slog.Logger
}

func (app *application) handlers() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, world!"))
	})
	return m
}
