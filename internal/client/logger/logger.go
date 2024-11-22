package logger

import (
	"log/slog"
	"os"
)

func NewClientLogger(f *os.File) *slog.Logger {
	return slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
