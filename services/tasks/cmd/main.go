package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tasks/internal/app"
	"tasks/internal/config"
	"tasks/internal/logger"
	notification_worker "tasks/internal/notification-worker"
	"tasks/internal/queue"
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

	// broker connection
	broker, err := queue.New(cfg, log, &storage)
	if err != nil {
		log.Error("error while connecting to rabbitmq", "error", err.Error())
		return
	}
	log.Info("rabbitmq connected")

	err = broker.Consume()
	if err != nil {
		log.Error("error while consuming broker", "error", err.Error())
		return
	}

	notificationWorker := notification_worker.New(log, &storage, broker)
	go notificationWorker.Listen()

	app := app.New(log, &storage, broker, &cfg)
	app.MustListen()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	log.Info("stopping application", slog.String("signal", sign.String()))
	app.ConnectionServer.GracefulStop()
	storage.DB.Close()
	broker.GracefulStop()
}
