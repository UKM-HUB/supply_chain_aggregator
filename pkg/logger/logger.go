package logger

import (
	"log/slog"
	"os"
)

// New creates a structured JSON logger.
// In "development" environment it uses human-readable text output instead.
func New(service, environment string) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{Level: slog.LevelDebug}

	if environment == "development" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler).With(
		slog.String("service", service),
		slog.String("env", environment),
	)
}
