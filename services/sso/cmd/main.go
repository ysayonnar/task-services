package main

import (
	"sso/internal/config"
	"sso/internal/logger"
	"sso/internal/storage"
)

func main() {
	cfg := config.Parse()

	log := logger.New(cfg)
	log.Info("config parsed", "config", cfg)

	storage := storage.New()
	err := storage.Conn()
	if err != nil {
		log.Error("error while connecting to db", "error", err.Error())
		return
	}

	log.Info("db connected")

	for {
	}

	//TODO: инициализировать приложение
}
