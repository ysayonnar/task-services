package handlers

import (
	"context"
	"encoding/json"
	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type SsoRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SsoRegisterResponse struct {
	UserId int64 `json:"user_id"`
}

type SsoErrorResponse struct {
	Error string `json:"error"`
}

func (router *Router) SsoRegister(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.SsoRegister"
	log := router.Log.With(slog.String("op", op))

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("error while reading body", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var parsedBody SsoRegisterRequest
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		log.Error("error while parsing json", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	// NOTE: additionally, the error needs to be handled, but this can be added later
	response, err := router.Clients.SsoClient.Register(ctx, &sso.RegisterRequest{Email: parsedBody.Email, Password: parsedBody.Password})
	if err != nil {
		var errResponse SsoErrorResponse
		errResponse.Error = err.Error()

		out, err := json.Marshal(errResponse)
		if err != nil {
			log.Error("error while parsing json", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	var jsonRequest SsoRegisterResponse
	jsonRequest.UserId = response.GetUserId()

	out, err := json.Marshal(jsonRequest)
	if err != nil {
		log.Error("error while parsing json", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
