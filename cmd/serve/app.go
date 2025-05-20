package serve

import (
	"context"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/jonathonwebb/tilde/internal/core"
)

type application struct {
	log       *slog.Logger
	templates map[string]*template.Template
}

func run(ctx context.Context, w io.Writer, cfg core.Config) error {
	log := cfg.NewLogger(w)
	log.Debug("serve", "cfg", cfg)

	drivers, err := drivers()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	templates, err := pages(drivers)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("templates", "templates", templates)

	app := &application{
		log:       log,
		templates: templates,
	}

	log.Info("starting server", "addr", cfg.Serve.Addr)
	err = http.ListenAndServe(cfg.Serve.Addr, app.routes())
	log.Error(err.Error())
	return err
}

func drivers() (*template.Template, error) {
	drivers, err := template.ParseGlob("ui/html/layouts/*.html")
	if err != nil {
		return nil, err
	}
	// drivers, err = drivers.ParseGlob("ui/html/partials/*.html")
	// if err != nil {
	// 	return nil, err
	// }
	return drivers, nil
}

func pages(drivers *template.Template) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		clone, err := drivers.Clone()
		if err != nil {
			return nil, err
		}
		tmpl, err := clone.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = tmpl
	}
	return cache, nil
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		tmpl, ok := app.templates["root.html"]
		if !ok {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, "base.html", nil)
	})
	return app.logReq(app.sharedHeaders(mux))
}
