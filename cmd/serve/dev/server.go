package dev

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/gorilla/websocket"
	"github.com/r3labs/sse/v2"
	"golang.org/x/sync/errgroup"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Server struct {
	Addr string
	Port int

	logger     *slog.Logger
	clients    map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	broadcast  chan []byte
}

type ChangeEvent struct {
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
	Updated []string `json:"updated"`
}

func NewServer(logger *slog.Logger) (*Server, error) {
	s := Server{
		logger:     logger,
		clients:    make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		broadcast:  make(chan []byte),
	}
	go s.run()

	builder, builderErr := esbuild.Context(esbuild.BuildOptions{
		EntryPoints: []string{
			"ui/assets/entrypoints/**/*.js",
			"ui/assets/entrypoints/**/*.jsx",
			"ui/assets/entrypoints/**/*.ts",
			"ui/assets/entrypoints/**/*.tsx",
			"ui/assets/entrypoints/**/*.css",
		},
		Sourcemap:   esbuild.SourceMapInline,
		Bundle:      true,
		Outdir:      "ui/static/build",
		Format:      esbuild.FormatESModule,
		JSX:         esbuild.JSXTransform,
		JSXFactory:  "h",
		JSXFragment: "Fragment",
	})
	if builderErr != nil {
		logger.Error(builderErr.Error())
		// TODO: handle error
	}

	watchErr := builder.Watch(esbuild.WatchOptions{})
	if watchErr != nil {
		logger.Error(watchErr.Error())
		// TODO: handle error
	}

	serveResult, serveErr := builder.Serve(esbuild.ServeOptions{
		Servedir: "ui/static",
		CORS:     esbuild.CORSOptions{Origin: []string{"*"}},
	})
	if serveErr != nil {
		logger.Error(serveErr.Error())
		// TODO: handle error
	}
	// TODO: need to account for multiple hosts, or len 0?
	s.Addr = serveResult.Hosts[0]
	s.Port = int(serveResult.Port)
	logger.Info(fmt.Sprintf("started esbuild server @ %s:%d", serveResult.Hosts[0], serveResult.Port))

	sseClient := sse.NewClient(fmt.Sprintf("http://%s:%d/esbuild", serveResult.Hosts[0], serveResult.Port))

	processPaths := func(paths []string, result *[]string) {
		for _, path := range paths {
			if !strings.HasSuffix(path, ".map") {
				*result = append(*result, path)
			}
		}
	}

	g := new(errgroup.Group)
	g.Go(func() error {
		return sseClient.Subscribe("change", func(msg *sse.Event) {
			var ce ChangeEvent
			added := []string{}
			removed := []string{}
			updated := []string{}

			if len(msg.Data) > 0 {
				err := json.Unmarshal(msg.Data, &ce)
				if err == nil {
					processPaths(ce.Added, &added)
					processPaths(ce.Removed, &removed)
					processPaths(ce.Updated, &updated)
					newCe, err := json.Marshal(ChangeEvent{Added: added, Removed: removed, Updated: updated})
					if err == nil {
						s.Broadcast(newCe)
					}
				}
			}
		})
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *Server) run() {
	for {
		select {
		case client := <-s.register:
			s.logger.Info("dev client connected")
			s.clients[client] = true
		case client := <-s.unregister:
			s.logger.Info("dev client disconnected")
			delete(s.clients, client)
		case msg := <-s.broadcast:
			s.logger.Info("broadcasting")
			for client := range s.clients {
				err := client.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					s.logger.Error(err.Error())
					_ = client.Close()
					delete(s.clients, client)
				}
			}
		}
	}
}

func (s *Server) Upgrade(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	defer func() {
		s.unregister <- conn
		if err := conn.Close(); err != nil {
			s.logger.Error(fmt.Sprintf("error closing connection: %v", err))
		}
	}()

	s.register <- conn

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
	return nil
}

func (s *Server) Broadcast(msg []byte) {
	s.broadcast <- msg
}
