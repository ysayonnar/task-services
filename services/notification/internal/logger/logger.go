package logger

import (
	"log/slog"
	"notification/internal/config"
	"os"
)

func New(cfg config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{}

	switch cfg.Env {
	case "dev":
		opts.Level = slog.LevelDebug
	case "prod":
		opts.Level = slog.LevelInfo
	default:
		opts.Level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, opts)

	return slog.New(handler)
}
