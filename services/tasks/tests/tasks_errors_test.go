package tests

import (
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
	tasks "github.com/ysayonnar/task-contracts/tasks/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"tasks/tests/suite"
	"testing"
	"time"
)

func TestCategoryAlreadyExistsOrNotFound(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword(passwordDefaultLength)
	registrationResponse, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, registrationResponse.GetUserId())

	loginResponse, err := st.AuthClient.Login(ctx, &sso.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, loginResponse.GetToken())

	categoryName := gofakeit.HipsterWord()
	createCatResp1, err := st.TasksClient.CreateCategory(ctx, &tasks.CreateCategoryRequest{
		Token: loginResponse.GetToken(),
		Name:  categoryName,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, createCatResp1.GetCategoryId())

	createCatResp2, err := st.TasksClient.CreateCategory(ctx, &tasks.CreateCategoryRequest{
		Token: loginResponse.GetToken(),
		Name:  categoryName,
	})
	require.Error(t, err)
	require.Equal(t, codes.AlreadyExists, status.Code(err))
	assert.Empty(t, createCatResp2)

	taskTitle := gofakeit.HipsterWord()
	taskDesc := gofakeit.HipsterSentence(descriptionDefaultLength)
	deadline := timestamppb.New(time.Now().Add(24 * time.Hour))

	createTaskResponse, err := st.TasksClient.CreateTask(ctx, &tasks.CreateTaskRequest{
		Token:        loginResponse.GetToken(),
		CategoryId:   123456789,
		Title:        taskTitle,
		Description:  taskDesc,
		IsNotificate: false,
		Deadline:     deadline,
	})
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Empty(t, createTaskResponse)
}
