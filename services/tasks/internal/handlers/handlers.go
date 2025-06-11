package handlers

import (
	"log/slog"
	"tasks/internal/config"
	"tasks/internal/storage"

	tasks "github.com/ysayonnar/task-contracts/tasks/gen/go"
)

type TasksServer struct {
	tasks.UnimplementedTasksServiceServer
	Log     *slog.Logger
	Storage *storage.Storage
	Cfg     *config.Config
}

// TODO: implement here
