package main

import (
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/logger"
	"sso/internal/storage"
)

func main() {
	// config parsing and logger setup
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

	// application initializing
	app := app.New(log, &storage)
	err = app.Listen(cfg.GRPC.Port)
	if err != nil {
		log.Error("error while listening gprc server", "error", err.Error())
		return
	}

	//TODO: graceful
}
