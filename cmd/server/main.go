package main

import (
	"fmt"
	"os"

	"github.com/JIIL07/cloudFiles-manager/internal/lib/slg"
	"github.com/JIIL07/cloudFiles-manager/internal/logger"
	"github.com/JIIL07/cloudFiles-manager/internal/server"
	"github.com/JIIL07/cloudFiles-manager/internal/storage"
)

func main() {
	cfg := server.MustLoad()

	log := logger.NewLogger(cfg.Env)

	storage, err := storage.InitDatabase(cfg.Database)
	if err != nil {
		log.Error("Failed to initialize database", slg.Err(err))
		os.Exit(1)
	}
	fmt.Println(storage.DB)
	// go func() {
	// 	if err := srv.Start(); err != nil {
	// 		log.Fatalf("Server failed: %v", err)
	// 	}
	// }()

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// <-c

	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// defer cancel()

	// if err := srv.Stop(ctx); err != nil {
	// 	log.Fatalf("Server shutdown failed: %v", err)
	// }

	// log.Println("Server gracefully stopped")
}
