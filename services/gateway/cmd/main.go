package main

import (
	"gateway/internal/config"
	"gateway/internal/logger"
)

func main() {
	cfg := config.MustParse()
	log := logger.New(cfg)
	log.Info("config parsed", "config", cfg)

	// TODO: setup application
}
