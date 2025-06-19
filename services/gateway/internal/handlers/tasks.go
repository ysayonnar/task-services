package handlers

import (
	"context"
	"gateway/internal/utils"
	tasks "github.com/ysayonnar/task-contracts/tasks/gen/go"
	"log/slog"
	"net/http"
	"time"
)

func (router *Router) CreateTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.CreateTask"
	log := router.Log.With(slog.String("op", op))

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var dto tasks.CreateTaskRequest
	err := utils.ReadJSON(w, r, &dto)
	if err != nil {
		log.Error("error while reading json", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	response, err := router.Clients.TasksClient.CreateTask(ctx, &dto)
	if err != nil {
		var errResponse ErrorResponse
		errResponse.Error = err.Error()

		utils.WriteJSON(w, errResponse, http.StatusOK)
		return
	}

	utils.WriteJSON(w, response, http.StatusOK)
}

func (router *Router) CreateCategory(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.CreateCategory"
	log := router.Log.With(slog.String("op", op))

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var dto tasks.CreateCategoryRequest
	err := utils.ReadJSON(w, r, &dto)
	if err != nil {
		log.Error("error while reading json", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	response, err := router.Clients.TasksClient.CreateCategory(ctx, &dto)
	if err != nil {
		var errResponse ErrorResponse
		errResponse.Error = err.Error()

		utils.WriteJSON(w, errResponse, http.StatusOK)
		return
	}

	utils.WriteJSON(w, response, http.StatusOK)
}

func (router *Router) DeleteTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.DeleteTask"
	log := router.Log.With(slog.String("op", op))

	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var dto tasks.DeleteTaskRequest
	err := utils.ReadJSON(w, r, &dto)
	if err != nil {
		log.Error("error while reading json", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	response, err := router.Clients.TasksClient.DeleteTask(ctx, &dto)
	if err != nil {
		var errResponse ErrorResponse
		errResponse.Error = err.Error()

		utils.WriteJSON(w, errResponse, http.StatusOK)
		return
	}

	utils.WriteJSON(w, response, http.StatusOK)
}

func (router *Router) GetTasks(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.GetTasks"
	log := router.Log.With(slog.String("op", op))

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var dto tasks.GetTasksRequest
	err := utils.ReadJSON(w, r, &dto)
	if err != nil {
		log.Error("error while reading json", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	response, err := router.Clients.TasksClient.GetTasks(ctx, &dto)
	if err != nil {
		var errResponse ErrorResponse
		errResponse.Error = err.Error()

		utils.WriteJSON(w, errResponse, http.StatusOK)
		return
	}

	utils.WriteJSON(w, response, http.StatusOK)
}

func (router *Router) GetTasksByCategory(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.GetTasksByCategory"
	log := router.Log.With(slog.String("op", op))

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var dto tasks.GetTasksByCategoryRequest
	err := utils.ReadJSON(w, r, &dto)
	if err != nil {
		log.Error("error while reading json", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	response, err := router.Clients.TasksClient.GetTasksByCategory(ctx, &dto)
	if err != nil {
		var errResponse ErrorResponse
		errResponse.Error = err.Error()

		utils.WriteJSON(w, errResponse, http.StatusOK)
		return
	}

	utils.WriteJSON(w, response, http.StatusOK)
}

func (router *Router) UpdateTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.UpdateTask"
	log := router.Log.With(slog.String("op", op))

	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var dto tasks.UpdateTaskRequest
	err := utils.ReadJSON(w, r, &dto)
	if err != nil {
		log.Error("error while reading json", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	response, err := router.Clients.TasksClient.UpdateTask(ctx, &dto)
	if err != nil {
		var errResponse ErrorResponse
		errResponse.Error = err.Error()

		utils.WriteJSON(w, errResponse, http.StatusOK)
		return
	}

	utils.WriteJSON(w, response, http.StatusOK)
}
