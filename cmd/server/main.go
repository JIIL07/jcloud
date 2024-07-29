package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JIIL07/cloudFiles-manager/internal/config"
	jenv "github.com/JIIL07/cloudFiles-manager/internal/lib/env"
	"github.com/JIIL07/cloudFiles-manager/internal/lib/slg"
	"github.com/JIIL07/cloudFiles-manager/internal/logger"
	"github.com/JIIL07/cloudFiles-manager/internal/server"
	"github.com/JIIL07/cloudFiles-manager/internal/storage"
)

func main() {
	//load env variables
	jenv.LoadEnv()

	//load config file
	cfg := config.MustLoad()

	//init logger
	log := logger.NewLogger(cfg.Env)

	//init storage
	s, err := storage.InitDatabase(cfg.Database)
	if err != nil {
		log.Error("Failed to initialize database", slg.Err(err))
		os.Exit(1)
	}
	defer s.CloseDatabase()

	//load .env
	jenv.LoadEnv()

	//init server
	srv := server.New(cfg.Server, s)

	//start server
	go func() {
		log.Info("Server starting on port :8080")
		if err := srv.Start(); err != nil {
			log.Error("Server failed to start", slg.Err(err))
			os.Exit(1)
		}
	}()

	//graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		log.Error("Server shutdown failed", slg.Err(err))
		os.Exit(1)
	}

	log.Info("Server gracefully stopped")
}
