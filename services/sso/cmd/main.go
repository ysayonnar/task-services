package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/logger"
	"sso/internal/queue"
	"sso/internal/storage"
	"syscall"
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

	// broker initialization
	broker, err := queue.New(cfg)
	if err != nil {
		log.Error("error while connecting to rabbitmq", "error", err.Error())
		return
	}
	log.Info("rabbitmq connected")

	// application initializing
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
