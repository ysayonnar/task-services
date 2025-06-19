package handlers

import (
	"context"
	"encoding/json"
	"gateway/internal/utils"
	"io"
	"log/slog"
	"net/http"
	"time"

	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
)

func (router *Router) Register(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.Register"
	log := router.Log.With(slog.String("op", op))

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var dto sso.RegisterRequest
	err := utils.ReadJSON(w, r, &dto)
	if err != nil {
		log.Error("error while reading json", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	response, err := router.Clients.SsoClient.Register(ctx, &dto)
	if err != nil {
		var errResponse ErrorResponse
		errResponse.Error = err.Error()

		utils.WriteJSON(w, errResponse, http.StatusOK)
		return
	}

	utils.WriteJSON(w, response, http.StatusOK)
}

func (router *Router) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.Login"
	log := router.Log.With(slog.String("op", op))

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var dto sso.LoginRequest
	err := utils.ReadJSON(w, r, &dto)
	if err != nil {
		log.Error("error while reading json", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	response, err := router.Clients.SsoClient.Login(ctx, &dto)
	if err != nil {
		var errResponse ErrorResponse
		errResponse.Error = err.Error()

		utils.WriteJSON(w, errResponse, http.StatusOK)
		return
	}

	utils.WriteJSON(w, response, http.StatusOK)
}

func (router *Router) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.DELETE"
	log := router.Log.With(slog.String("op", op))

	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var dto sso.DeleteRequest
	err := utils.ReadJSON(w, r, &dto)
	if err != nil {
		log.Error("error while reading json", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	response, err := router.Clients.SsoClient.Delete(ctx, &dto)
	if err != nil {
		var errResponse ErrorResponse
		errResponse.Error = err.Error()

		utils.WriteJSON(w, errResponse, http.StatusOK)
		return
	}

	utils.WriteJSON(w, response, http.StatusOK)
}
