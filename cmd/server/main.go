// nolint:errcheck
package main

import (
	"context"
	"github.com/JIIL07/jcloud/internal/server/config"
	"github.com/JIIL07/jcloud/internal/server/cookies"
	"github.com/JIIL07/jcloud/internal/server/logger"
	"github.com/JIIL07/jcloud/internal/server/server"
	"github.com/JIIL07/jcloud/internal/server/static"
	"github.com/JIIL07/jcloud/internal/server/storage"
	"github.com/JIIL07/jcloud/pkg/env"
	"github.com/JIIL07/jcloud/pkg/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// main is the entry point of the program. It loads environment variables, loads
// the configuration file, initializes the log, initializes the storage, sets
// up a new cookie storage, initializes the server, starts the server in a
// separate goroutine, handles graceful shutdown, and gracefully stops the
// server.
func main() {
	//load env variables
	jenv.LoadEnv()

	//load config file
	cfg := config.MustLoad()

	//init log
	log := logger.NewLogger(cfg.Env)

	//init storage
	s, err := storage.InitDatabase(cfg)
	if err != nil {
		log.Error("Failed to initialize database", jlog.Err(err))
		os.Exit(1)
	}
	defer s.CloseDatabase()

	binary, err := static.LoadStatic(cfg.Static.Path)
	if err != nil {
		log.Error("Failed to load static files", jlog.Err(err))
		os.Exit(1)
	}

	//setup new cookie storage
	cookies.SetNewCookieStore()

	//init server
	srv := server.New(cfg.Server, s, binary)

	//start server
	go func() {
		log.Info("Server starting on port :8080")
		if err := srv.Start(); err != nil {
			log.Error("Server failed to start", jlog.Err(err))
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
		log.Error("Server shutdown failed", jlog.Err(err))
		os.Exit(1)
	}

	log.Info("Server gracefully stopped")
}
