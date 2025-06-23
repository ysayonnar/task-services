package suite

import (
	"context"
	"net"
	"sso/internal/config"
	"testing"
	"time"

	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ssoPort = "8081"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient sso.AuthServiceClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustParseByPath("../config/config.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	conn, err := grpc.DialContext(context.Background(), net.JoinHostPort("localhost", ssoPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connecion failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        &cfg,
		AuthClient: sso.NewAuthServiceClient(conn),
	}
}
