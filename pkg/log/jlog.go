package jlog

import (
	"log/slog"
)

func Info(msg string) slog.Attr {
	return slog.Attr{
		Key:   "info",
		Value: slog.StringValue(msg),
	}
}
func Debug(msg string) slog.Attr {
	return slog.Attr{
		Key:   "debug",
		Value: slog.StringValue(msg),
	}
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
