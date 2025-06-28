package app

import (
	"fmt"
	"log/slog"
	"net"
	"sso/internal/config"
	"sso/internal/handlers"
	"sso/internal/queue"
	"sso/internal/storage"

	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
	"google.golang.org/grpc"
)

type App struct {
	GrpcServer       handlers.SsoServer
	ConnectionServer *grpc.Server
}

func New(log *slog.Logger, storage *storage.Storage, broker *queue.Broker, cfg *config.Config) App {
	var app App

	app.GrpcServer = handlers.SsoServer{
		Log:     log,
		Storage: storage,
		Cfg:     cfg,
		Broker:  broker,
	}
	app.ConnectionServer = grpc.NewServer()

	return app
}

func (app *App) MustListen() {
	const op = "app.Listen()"
	log := app.GrpcServer.Log.With(slog.String("op", op))

	port := app.GrpcServer.Cfg.GRPC.Port

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	sso.RegisterAuthServiceServer(app.ConnectionServer, &app.GrpcServer)

	log.Info(fmt.Sprintf("listening grpc on port %d", port))
	if err := app.ConnectionServer.Serve(conn); err != nil {
		panic(err)
	}
}
