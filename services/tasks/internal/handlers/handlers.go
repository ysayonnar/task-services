package handlers

import (
	"context"
	"errors"
	"log/slog"
	"tasks/internal/config"
	"tasks/internal/models"
	"tasks/internal/storage"

	tasks "github.com/ysayonnar/task-contracts/tasks/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	taskToCreate := models.Task{
		UserId:       req.GetUserId(),
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

	deletedTaskId, err := server.Storage.DeleteTask(ctx, req.GetUserId(), req.GetTaskId())
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

	foundTasks, err := server.Storage.GetTasksByUserId(ctx, req.GetUserId())
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

	foundTasks, err := server.Storage.GetTasksByUserIdAndCategoryId(ctx, req.GetUserId(), req.GetCategoryId())
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

	var task models.Task
	task.Title = req.GetTitle()
	task.Description = req.GetDescription()
	task.IsNotificate = req.GetIsNotificate()
	task.Deadline = req.GetDeadline().AsTime()
	task.UserId = req.GetUserId()
	task.TaskId = req.GetTaskId()

	taskId, err := server.Storage.UpdateTask(ctx, task)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}

		log.Error("error while updating task")
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &tasks.UpdateTaskResponse{TaskId: taskId}, nil

	return nil, nil
}
