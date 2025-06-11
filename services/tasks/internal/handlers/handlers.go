package handlers

import (
	"context"
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

func (server *TasksServer) CreateTask(ctx context.Context, req *tasks.CreateTaskRequest) (*tasks.CreateTaskResponse, error) {
	const op = "handlers.CreateTask"
	log := server.Log.With(slog.String("op", op))

	return nil, nil
}

func (server *TasksServer) DeleteTask(ctx context.Context, req *tasks.DeleteTaskRequest) (*tasks.DeleteTaskResponse, error) {
	const op = "handlers.DeleteTask"
	log := server.Log.With(slog.String("op", op))
	_ = log

	return nil, nil
}

func (server *TasksServer) GetTasks(ctx context.Context, req *tasks.GetTasksRequest) (*tasks.GetTasksResponse, error) {
	const op = "handlers.GetTasks"
	log := server.Log.With(slog.String("op", op))
	_ = log

	return nil, nil
}

func (server *TasksServer) GetTasksByCategory(ctx context.Context, req *tasks.GetTasksByCategoryRequest) (*tasks.GetTasksByCategoryResponse, error) {
	const op = "handlers.GetTasksByCategory"
	log := server.Log.With(slog.String("op", op))
	_ = log

	return nil, nil
}

func (server *TasksServer) UpdateTask(ctx context.Context, req *tasks.UpdateTaskRequest) (*tasks.UpdateTaskResponse, error) {
	const op = "handlers.UpdateTask"
	log := server.Log.With(slog.String("op", op))
	_ = log

	return nil, nil
}
