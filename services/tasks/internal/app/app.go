package app

import (
	"fmt"
	"log/slog"
	"net"
	"tasks/internal/config"
	"tasks/internal/handlers"
	"tasks/internal/storage"

	tasks "github.com/ysayonnar/task-contracts/tasks/gen/go"
	"google.golang.org/grpc"
)

type App struct {
	GrpcServer       handlers.TasksServer
	ConnectionServer *grpc.Server
}

func New(log *slog.Logger, storage *storage.Storage, cfg *config.Config) App {
	var app App

	app.GrpcServer = handlers.TasksServer{
		Log:     log,
		Storage: storage,
		Cfg:     cfg,
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

	tasks.RegisterTasksServiceServer(app.ConnectionServer, &app.GrpcServer)

	log.Info(fmt.Sprintf("listening grpc on port %d", port))
	if err := app.ConnectionServer.Serve(conn); err != nil {
		panic(err)
	}
}
