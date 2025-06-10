package main

import (
	"tasks/internal/config"
	"tasks/internal/logger"
	"tasks/internal/storage"
)

func main() {
	cfg := config.MustParse()
	log := logger.New(cfg)
	log.Info("config parsed", "config", cfg)

	// storage connection
	storage := storage.New()
	err := storage.Conn()
	if err != nil {
		log.Error("error while connecting to db", "error", err.Error())
		return
	}
	log.Info("db connected")

	for {
	}
}
