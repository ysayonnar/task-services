package app

import (
	"fmt"
	"log/slog"
	"net"
	"sso/internal/handlers"
	"sso/internal/storage"

	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
	"google.golang.org/grpc"
)

type App struct {
	Log        *slog.Logger
	GrpcServer handlers.SsoServer
}

func New(log *slog.Logger, storage *storage.Storage) App {
	var app App

	app.GrpcServer = handlers.SsoServer{
		Log:     log,
		Storage: storage,
	}

	app.Log = log

	return app
}

func (app *App) Listen(port int) error {
	const op = "app.Listen()"

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	server := grpc.NewServer()
	sso.RegisterAuthServiceServer(server, app.GrpcServer)

	app.Log.Info(fmt.Sprintf("listening grpc on port %d", port))
	if err := server.Serve(conn); err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	return nil
}
