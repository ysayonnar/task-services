package main

import (
	"log/slog"
	"notification/internal/cache"
	"notification/internal/config"
	"notification/internal/logger"
	"notification/internal/queue"
	"os"
	"os/signal"
	"syscall"
)

// Будет Redis для хранения пар user_id : email, при обработке запроса на уведомление, хожу в редис, беру email.
// Если email нет в Redis, кидаю http запрос в gateway на получение email по user_id.
// При регистрации пользователя тоже ловлю сигнал и добавляю в Redis

func main() {
	cfg := config.MustParse()

	log := logger.New(cfg)
	log.Info("config parsed")
	log.Info("logger set up")

	redisClient := cache.ConnectRedis()

	//TODO: setup SMTP

	broker, err := queue.New(cfg, log, redisClient)
	if err != nil {
		log.Error("error while connecting to the broker", "error", err.Error())
		return
	}

	err = broker.Consume()
	if err != nil {
		log.Error("error while listening from broker", "error", err.Error())
		return
	}

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	log.Info("stopping application", slog.String("signal", sign.String()))
	broker.GracefulStop()
	redisClient.Close()
}
