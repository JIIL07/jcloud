package cloudfiles

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ServerContext struct {
	Ctx *FileContext
}

func NewServerContext(db *sql.DB) *ServerContext {
	ctx := &FileContext{
		DB:   db,
		Info: &Info{},
	}
	return &ServerContext{
		Ctx: ctx,
	}
}

type config struct {
	port int
}

type application struct {
	config config
	logger *slog.Logger
}

func (s *ServerContext) Start() error {
	var cfg config

	port := os.Getenv("PORT")
	intPort, err := strconv.Atoi(port)
	if err != nil {
		intPort = 8080
	}

	// Set the port to run the API on
	cfg.port = intPort

	// create the logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	// create the server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(s),
		IdleTimeout:  45 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("server started", "addr", srv.Addr)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

	return err
}
