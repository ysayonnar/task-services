package grpc

import (
	"fmt"

	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
	tasks "github.com/ysayonnar/task-contracts/tasks/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClients struct {
	SsoClient   sso.AuthServiceClient
	TasksClient tasks.TasksServiceClient
	// NOTE: add services here
}

func NewGrpcClients() (*GrpcClients, error) {
	const op = "grpc.NewGrpcClients"

	ssoConn, err := grpc.Dial("sso:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	tasksConn, err := grpc.Dial("tasks:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	// NOTE: add similar connections here

	return &GrpcClients{
		SsoClient:   sso.NewAuthServiceClient(ssoConn),
		TasksClient: tasks.NewTasksServiceClient(tasksConn),
	}, nil
}
