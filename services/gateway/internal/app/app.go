package app

import (
	"fmt"
	"gateway/internal/config"
	"gateway/internal/grpc"
	"gateway/internal/handlers"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	Server *http.Server
	Log    *slog.Logger
	Config *config.Config
}

func New(logger *slog.Logger, cfg config.Config, clients *grpc.GrpcClients) App {
	router := handlers.InitRouter(logger, clients)

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  time.Second * time.Duration(cfg.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(cfg.WriteTimeout),
		Handler:      router.Mux,
	}

	return App{Server: &srv, Log: logger, Config: &cfg}
}

func (app *App) MustListen() {
	err := app.Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
