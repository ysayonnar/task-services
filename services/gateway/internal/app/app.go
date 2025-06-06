package app

import (
	"gateway/internal/config"
	"log/slog"
	"net/http"
)

type App struct {
	Server *http.Server
	Log    *slog.Logger
	Config *config.Config
}

func New(logger *slog.Logger, cfg config.Config) App {
	//TODO: setup server from cfg
	srv := http.Server{}
	return App{Server: &srv, Log: logger, Config: &cfg}
}

func (app *App) MustListen() {
	//TODO: listen
}
