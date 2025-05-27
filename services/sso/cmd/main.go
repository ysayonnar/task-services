package main

import (
	"sso/internal/config"
	"sso/internal/logger"
)

func main() {
	cfg := config.Parse()

	log := logger.New(cfg)

	log.Info("config parsed", "config", cfg)

	//TODO: подключиться к бд

	//TODO: инициализировать приложение
}
