package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tasks/internal/app"
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

	app := app.New(log, &storage, &cfg)
	app.MustListen()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	log.Info("stopping application", slog.String("signal", sign.String()))
	app.ConnectionServer.GracefulStop()
	storage.DB.Close()
}
