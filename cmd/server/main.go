package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JIIL07/jcloud/internal/server/config"
	"github.com/JIIL07/jcloud/internal/server/cookies"
	"github.com/JIIL07/jcloud/internal/server/logger"
	"github.com/JIIL07/jcloud/internal/server/server"
	"github.com/JIIL07/jcloud/internal/server/static"
	"github.com/JIIL07/jcloud/internal/server/storage"
	jenv "github.com/JIIL07/jcloud/pkg/env"
	jlog "github.com/JIIL07/jcloud/pkg/log"
)

func main() {
	jenv.LoadEnv()
	cfg := config.MustLoad()
	log := logger.NewLogger(cfg.Env)
	s, err := storage.InitDatabase(cfg)
	if err != nil {
		log.Error("Failed to initialize database", jlog.Err(err))
		os.Exit(1)
	}
	defer func() {
		if err := s.CloseDatabase(); err != nil {
			log.Error("Failed to close database", jlog.Err(err))
		}
	}()
	binary, err := static.LoadStatic(cfg.Static.Path)
	if err != nil {
		log.Error("Failed to load static files", jlog.Err(err))
		os.Exit(1)
	}
	cookies.SetNewCookieStore()
	srv := server.New(cfg.Server, s, binary)
	go func() {
		log.Info("Server starting on port :8080")
		if err := srv.Start(); err != nil {
			log.Error("Server failed to start", jlog.Err(err))
			os.Exit(1)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		log.Error("Server shutdown failed", jlog.Err(err))
		os.Exit(1)
	}
	log.Info("Server gracefully stopped")
}
