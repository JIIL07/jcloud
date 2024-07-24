package server

import (
	"context"
	"net"
	"net/http"

	"github.com/JIIL07/cloudFiles-manager/internal/config"
	"github.com/JIIL07/cloudFiles-manager/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	httpServer *http.Server
}

type Connection struct {
	IP net.IP
}

type CurrentUser struct {
	UserData       storage.AboutUser
	Role           string
	NetworkDetails Connection
}

func New(c config.ServerConfig, s *storage.Storage) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              c.Address,
			Handler:           setupRouter(s),
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
