// nolint:errcheck
package main

import (
	"context"
	"github.com/JIIL07/jcloud/internal/config"
	"github.com/JIIL07/jcloud/internal/lib/cookies"
	"github.com/JIIL07/jcloud/internal/lib/env"
	"github.com/JIIL07/jcloud/internal/lib/slg"
	"github.com/JIIL07/jcloud/internal/logger"
	"github.com/JIIL07/jcloud/internal/server"
	"github.com/JIIL07/jcloud/internal/storage"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// main is the entry point of the program. It loads environment variables, loads
// the configuration file, initializes the logger, initializes the storage, sets
// up a new cookie storage, initializes the server, starts the server in a
// separate goroutine, handles graceful shutdown, and gracefully stops the
// server.
func main() {
	//load env variables
	jenv.LoadEnv()

	//load config file
	cfg := config.MustLoad()

	//init logger
	log := logger.NewLogger(cfg.Env)

	//init storage
	s, err := storage.InitDatabase(cfg)
	if err != nil {
		log.Error("Failed to initialize database", slg.Err(err))
		os.Exit(1)
	}
	defer s.CloseDatabase()

	//setup new cookie storage
	cookies.SetNewCookieStore()

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
