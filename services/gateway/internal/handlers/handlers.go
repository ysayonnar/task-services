package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Router struct {
	Log *slog.Logger
	Mux *http.ServeMux
}

func InitRouter(logger *slog.Logger) Router {
	router := Router{Log: logger, Mux: http.NewServeMux()}

	router.Mux.HandleFunc("/echo", router.Echo)

	return router
}

func (router *Router) Echo(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.Echo"
	log := router.Log.With(slog.String("op", op))

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("error while parsing body", "error", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, string(body))
}
