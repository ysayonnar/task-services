package main

import (
	"gateway/internal/app"
	"gateway/internal/config"
	"gateway/internal/logger"
	"log/slog"
)

func main() {
	cfg := config.MustParse()
	log := logger.New(cfg)
	log.Info("config parsed", "config", cfg)

	app := app.New(log, cfg)
	app.MustListen()

	log.Info("server started", slog.Int("port", cfg.Port))
}
