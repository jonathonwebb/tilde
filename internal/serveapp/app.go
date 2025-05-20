package serveapp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/jonathonwebb/tilde/internal/core"
	"github.com/jonathonwebb/tilde/internal/serveapp/dev"
)

type App struct {
	cfg           *core.Config
	logger        *slog.Logger
	templateFuncs template.FuncMap
	templateCache map[string]*template.Template
}

func New(cfg *core.Config, logger *slog.Logger) *App {
	return &App{
		cfg:           cfg,
		logger:        logger,
		templateCache: map[string]*template.Template{},
	}
}

func interceptListings(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) Run(ctx context.Context) error {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("ui/static"))
	mux.Handle("/public/", http.StripPrefix("/public", interceptListings(fileServer)))

	mux.Handle("GET /{$}", http.HandlerFunc(a.rootHandler))
	mux.Handle("GET /login", http.HandlerFunc(a.loginHandler))
	mux.Handle("GET /settings", http.HandlerFunc(a.settingsHandler))

	mux.Handle("GET /health", http.HandlerFunc(a.healthHandler))

	var devServer *dev.Server
	var assetMeta map[string]struct {
		Path string `json:"path"`
		SRI  string `json:"sri"`
	}

	if a.cfg.Serve.Dev {
		devServer = dev.NewServer(a.logger)
		mux.HandleFunc("/dev", func(w http.ResponseWriter, r *http.Request) {
			err := devServer.Upgrade(w, r)
			if err != nil {
				w.WriteHeader(http.StatusOK)
			}
		})
	} else {
		metaJson, err := os.ReadFile("ui/static/build/meta.json")
		if err != nil {
			return err
		}
		json.Unmarshal(metaJson, &assetMeta)
	}

	a.templateFuncs = template.FuncMap{
		"resource": func(path string) (string, error) {
			if devServer != nil {
				base := fmt.Sprintf("http://%s:%d", devServer.Addr, devServer.Port)
				s, err := url.JoinPath(base, path)
				if err != nil {
					return "", err
				}
				return s, nil
			} else {
				s, ok := assetMeta[path]
				if !ok {
					return "", fmt.Errorf("asset not found")
				}
				return s.Path, nil
			}
		},
	}

	s := &http.Server{
		Addr:         a.cfg.Serve.Addr,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      mux,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	return s.ListenAndServe()
}

func (a *App) rootHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"ui/html/base.tmpl",
		"ui/html/pages/root.tmpl",
	}
	ts, err := template.New("root").Funcs(a.templateFuncs).ParseFiles(files...)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"ui/html/base.tmpl",
		"ui/html/pages/login.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (a *App) settingsHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"ui/html/base.tmpl",
		"ui/html/pages/settings.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	// if isTerminating.Load() {
	// 	http.Error(w, "terminating", http.StatusServiceUnavailable)
	// 	return
	// }
	fmt.Fprintln(w, "ready")
}
