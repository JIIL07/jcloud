package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	cloudfiles "github.com/JIIL07/cloudFiles-manager/client"
	"github.com/julienschmidt/httprouter"
)

type ServerContext struct {
	Ctx *cloudfiles.FileContext
}

func NewServerContext(db *sql.DB) *ServerContext {
	ctx := &cloudfiles.FileContext{
		DB:   db,
		Info: &cloudfiles.Info{},
	}
	return &ServerContext{
		Ctx: ctx,
	}
}

type config struct {
	port int
}

type application struct {
	config    config
	logger    *log.Logger
	errlogger *log.Logger
	router    *httprouter.Router
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
	logger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errlogger := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	app := &application{
		config:    cfg,
		logger:    logger,
		errlogger: errlogger,
	}

	// create the server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(s),
		IdleTimeout:  45 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	logger.Printf("Server is starting on %s\n", srv.Addr)

	err = srv.ListenAndServe()
	if err != nil {
		logger.Printf("Server failed to start: %v\n", err)
		os.Exit(1)
	}

	return err
}
