package main

import (
	"gateway/internal/app"
	"gateway/internal/config"
	"gateway/internal/grpc"
	"gateway/internal/logger"
	"log/slog"
)

func main() {
	cfg := config.MustParse()
	log := logger.New(cfg)
	log.Info("config parsed", "config", cfg)

	grpcClients, err := grpc.NewGrpcClients()
	if err != nil {
		log.Error("error while connecting to grpc", "error", err.Error())
		return
	}

	app := app.New(log, cfg, grpcClients)
	log.Info("server started", slog.Int("port", cfg.Port))
	app.MustListen()

}
