package main

import (
	"notification/internal/config"
	"notification/internal/logger"
)

// Будет Redis для хранения пар user_id : email, при обработке запроса на уведомление, хожу в редис, беру email.
// Если email нет в Redis, кидаю http запрос в gateway на получение email по user_id.
// При регистрации пользователя тоже ловлю сигнал и добавляю в Redis

func main() {
	cfg := config.MustParse()

	log := logger.New(cfg)
	log.Info("config parsed")
	log.Info("logger set up")

	for {

	}

	//TODO: connect to redis

	//TODO: setup SMTP

	//TODO: connect to broker

	//TODO: listen broker
}
