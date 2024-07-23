package server

import (
	"context"
	"net/http"

	"github.com/JIIL07/cloudFiles-manager/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

func New(c config.ServerConfig) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              c.Address,
			Handler:           setupRouter(),
			ReadTimeout:       c.ReadTimeout,
			WriteTimeout:      c.WriteTimeout,
			IdleTimeout:       c.IdleTimeout,
			ReadHeaderTimeout: c.ReadTimeout,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.httpServer.SetKeepAlivesEnabled(false)
	return s.httpServer.Close()
}
