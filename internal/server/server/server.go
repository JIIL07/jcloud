package server

import (
	"context"
	"github.com/JIIL07/jcloud/internal/server/config"
	"github.com/JIIL07/jcloud/internal/server/routes"
	"github.com/JIIL07/jcloud/internal/server/static"
	"github.com/JIIL07/jcloud/internal/server/storage"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	httpServer *http.Server
}

func New(config config.ServerConfig, storage *storage.Storage, binary *static.Files) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              config.Address,
			Handler:           routes.SetupRouter(binary, storage),
			ReadTimeout:       config.ReadTimeout,
			WriteTimeout:      config.WriteTimeout,
			IdleTimeout:       config.IdleTimeout,
			ReadHeaderTimeout: config.ReadTimeout,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.httpServer.SetKeepAlivesEnabled(false)
	return s.httpServer.Shutdown(ctx)
}
