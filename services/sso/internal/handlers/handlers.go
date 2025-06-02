package handlers

import (
	"log/slog"
	"sso/internal/storage"

	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
)

type SsoServer struct {
	sso.UnimplementedAuthServiceServer
	Log     *slog.Logger
	Storage *storage.Storage
}

// implement routes
