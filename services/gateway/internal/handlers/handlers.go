package handlers

import (
	"gateway/internal/config"
	"gateway/internal/grpc"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Router struct {
	Log     *slog.Logger
	Mux     *http.ServeMux
	Clients *grpc.GrpcClients
	Config  *config.Config
}

func InitRouter(logger *slog.Logger, clients *grpc.GrpcClients, cfg config.Config) Router {
	router := Router{Log: logger, Mux: http.NewServeMux(), Clients: clients, Config: &cfg}

	router.Mux.HandleFunc("/sso/register", router.Register)
	router.Mux.HandleFunc("/sso/login", router.Login)
	router.Mux.HandleFunc("/sso/delete", router.Delete)

	router.Mux.HandleFunc("/tasks/create-task", router.AuthMW(router.CreateTask))
	router.Mux.HandleFunc("/tasks/create-category", router.AuthMW(router.CreateCategory))
	router.Mux.HandleFunc("/tasks/all", router.AuthMW(router.GetTasks))
	router.Mux.HandleFunc("/tasks/category", router.AuthMW(router.GetTasksByCategory))
	router.Mux.HandleFunc("/tasks/delete", router.AuthMW(router.DeleteTask))
	router.Mux.HandleFunc("/tasks/update", router.AuthMW(router.UpdateTask))

	return router
}

func (router *Router) AuthMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.AuthMW"
		log := router.Log.With(slog.String("op", op))

		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(bearer, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return router.Config.Secret, nil
		})

		if err != nil || !token.Valid {
			log.Error("error while parsing jwt", "error", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
