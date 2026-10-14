package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New initializes a structured slog.Logger based on the environment.
// In "production" it outputs JSON; otherwise it outputs human-readable text.
func New(env string) *slog.Logger {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if strings.ToLower(env) == "development" || strings.ToLower(env) == "local" {
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
