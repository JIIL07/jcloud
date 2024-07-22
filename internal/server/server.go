package server

import (
	"context"
	"net/http"

	"github.com/JIIL07/cloudFiles-manager/internal/storage"

	_ "github.com/mattn/go-sqlite3"
)

func New(config ServerConfig, dbConfig storage.DBConfig) (*Server, error) {
	storage, err := storage.InitDatabase(dbConfig)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		httpServer: &http.Server{
			Addr:         config.Address,
			Handler:      setupRouter(storage.DB),
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			IdleTimeout:  config.IdleTimeout,
		},
		db: storage.DB,
	}

	return srv, nil
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.httpServer.SetKeepAlivesEnabled(false)
	return s.httpServer.Close()
}
