package logger

import (
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDebug = "debug"
	envProd  = "prod"
)

func NewLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDebug:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func NewClientLogger(f *os.File) *slog.Logger {
	return slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
