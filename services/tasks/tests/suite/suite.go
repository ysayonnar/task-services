package suite

import (
	"context"
	"net"
	"tasks/internal/config"
	"testing"
	"time"

	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
	tasks "github.com/ysayonnar/task-contracts/tasks/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ssoPort   = "8081"
	tasksPort = "8082"
)

type Suite struct {
	*testing.T
	Cfg         *config.Config
	AuthClient  sso.AuthServiceClient
	TasksClient tasks.TasksServiceClient
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

	ssoConn, err := grpc.DialContext(context.Background(), net.JoinHostPort("localhost", ssoPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc SSO server connecion failed: %v", err)
	}

	tasksConn, err := grpc.DialContext(context.Background(), net.JoinHostPort("localhost", tasksPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc TASKS server connecion failed: %v", err)
	}

	return ctx, &Suite{
		T:           t,
		Cfg:         &cfg,
		AuthClient:  sso.NewAuthServiceClient(ssoConn),
		TasksClient: tasks.NewTasksServiceClient(tasksConn),
	}
}
