package handlers

import (
	"gateway/internal/grpc"
	"log/slog"
	"net/http"
)

type Router struct {
	Log     *slog.Logger
	Mux     *http.ServeMux
	Clients *grpc.GrpcClients
}

func InitRouter(logger *slog.Logger, clients *grpc.GrpcClients) Router {
	router := Router{Log: logger, Mux: http.NewServeMux(), Clients: clients}

	router.Mux.HandleFunc("/sso/register", router.SsoRegister)
	router.Mux.HandleFunc("/sso/login", router.SsoLogin)

	return router
}
