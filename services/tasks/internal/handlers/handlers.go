package handlers

import (
	"context"
	"errors"
	tasks "github.com/ysayonnar/task-contracts/tasks/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"tasks/internal/config"
	"tasks/internal/models"
	"tasks/internal/queue"
	"tasks/internal/storage"
	"tasks/internal/utils"
)

type TasksServer struct {
	tasks.UnimplementedTasksServiceServer
	Log     *slog.Logger
	Storage *storage.Storage
	Cfg     *config.Config
	Broker  *queue.Broker
}

func (server *TasksServer) CreateTask(ctx context.Context, req *tasks.CreateTaskRequest) (*tasks.CreateTaskResponse, error) {
	const op = "handlers.CreateTask"
	log := server.Log.With(slog.String("op", op))

	userId, err := utils.ParseUserId(req.GetToken(), server.Cfg.Secret)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid bearer token")
	}

	taskToCreate := models.Task{
		UserId:       userId,
		Title:        req.GetTitle(),
		Description:  req.GetDescription(),
		Deadline:     req.GetDeadline().AsTime(),
		IsNotificate: req.GetIsNotificate(),
	}

	taskId, err := server.Storage.InsertTask(ctx, taskToCreate, req.GetCategoryId())
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			return nil, status.Error(codes.InvalidArgument, "category with such id doesn't exist")
		}

		log.Error("error while inserting task", "error", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &tasks.CreateTaskResponse{TaskId: taskId}, nil
}

func (server *TasksServer) CreateCategory(ctx context.Context, req *tasks.CreateCategoryRequest) (*tasks.CreateCategoryResponse, error) {
	const op = "handlers.CreateCategory"
	log := server.Log.With(slog.String("op", op))

	_, err := utils.ParseUserId(req.GetToken(), server.Cfg.Secret)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid bearer token")
	}

	categoryId, err := server.Storage.InsertCategory(ctx, req.GetName())
	if err != nil {
		if errors.Is(err, storage.ErrCategoryAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "category already exists")
		}

		log.Error("error while inserting category", "error", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &tasks.CreateCategoryResponse{CategoryId: categoryId}, nil
}

func (server *TasksServer) DeleteTask(ctx context.Context, req *tasks.DeleteTaskRequest) (*tasks.DeleteTaskResponse, error) {
	const op = "handlers.DeleteTask"
	log := server.Log.With(slog.String("op", op))

	userId, err := utils.ParseUserId(req.GetToken(), server.Cfg.Secret)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid bearer token")
	}

	deletedTaskId, err := server.Storage.DeleteTask(ctx, userId, req.GetTaskId())
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			return nil, status.Error(codes.NotFound, "task with such id doesn't exist")
		}

		log.Error("error while deleting task", "error", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &tasks.DeleteTaskResponse{TaskId: deletedTaskId}, nil
}

func (server *TasksServer) GetTasks(ctx context.Context, req *tasks.GetTasksRequest) (*tasks.GetTasksResponse, error) {
	const op = "handlers.GetTasks"
	log := server.Log.With(slog.String("op", op))

	userId, err := utils.ParseUserId(req.GetToken(), server.Cfg.Secret)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid bearer token")
	}

	foundTasks, err := server.Storage.GetTasksByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			return nil, status.Error(codes.NotFound, "tasks not found")
		}

		log.Error("error while finding tasks", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &tasks.GetTasksResponse{Tasks: foundTasks}, nil
}

func (server *TasksServer) GetTasksByCategory(ctx context.Context, req *tasks.GetTasksByCategoryRequest) (*tasks.GetTasksByCategoryResponse, error) {
	const op = "handlers.GetTasksByCategory"
	log := server.Log.With(slog.String("op", op))

	userId, err := utils.ParseUserId(req.GetToken(), server.Cfg.Secret)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid bearer token")
	}

	foundTasks, err := server.Storage.GetTasksByUserIdAndCategoryId(ctx, userId, req.GetCategoryId())
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			return nil, status.Error(codes.NotFound, "tasks not found")
		}

		log.Error("error while finding tasks", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &tasks.GetTasksByCategoryResponse{Tasks: foundTasks}, nil
}

func (server *TasksServer) UpdateTask(ctx context.Context, req *tasks.UpdateTaskRequest) (*tasks.UpdateTaskResponse, error) {
	const op = "handlers.UpdateTask"
	log := server.Log.With(slog.String("op", op))

	userId, err := utils.ParseUserId(req.GetToken(), server.Cfg.Secret)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid bearer token")
	}

	taskToUpdate := models.Task{
		TaskId:       req.GetTaskId(),
		UserId:       userId,
		Title:        req.GetTitle(),
		Description:  req.GetDescription(),
		Deadline:     req.GetDeadline().AsTime(),
		IsNotificate: req.GetIsNotificate(),
	}

	taskId, err := server.Storage.UpdateTask(ctx, taskToUpdate)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}

		log.Error("error while updating task")
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &tasks.UpdateTaskResponse{TaskId: taskId}, nil
}
