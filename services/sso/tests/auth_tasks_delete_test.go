package tests

import (
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
	tasks "github.com/ysayonnar/task-contracts/tasks/gen/go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sso/tests/suite"
	"testing"
	"time"
)

func TestDeleteTasksBroker(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword(passwordDefaultLength)

	registrationResponse, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, registrationResponse.GetUserId())

	createCategoryResponse, err := st.TasksClient.CreateCategory(ctx, &tasks.CreateCategoryRequest{Name: gofakeit.HipsterWord()})
	require.NoError(t, err)
	require.NotEmpty(t, createCategoryResponse.GetCategoryId())

	createTaskResponse, err := st.TasksClient.CreateTask(ctx, &tasks.CreateTaskRequest{
		UserId:       registrationResponse.GetUserId(),
		CategoryId:   createCategoryResponse.GetCategoryId(),
		Title:        gofakeit.HipsterWord(),
		Description:  gofakeit.HipsterSentence(5),
		IsNotificate: false,
		Deadline:     timestamppb.New(time.Now().Add(24 * time.Hour)),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, createTaskResponse.GetTaskId())

	getTasksResponse, err := st.TasksClient.GetTasks(ctx, &tasks.GetTasksRequest{UserId: registrationResponse.GetUserId()})
	require.NoError(t, err)
	require.Equal(t, 1, len(getTasksResponse.GetTasks()))

	deleteUserResponse, err := st.AuthClient.Delete(ctx, &sso.DeleteRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.True(t, deleteUserResponse.GetIsDeleted())

	getTasksAfterDeleteResponse, err := st.TasksClient.GetTasks(ctx, &tasks.GetTasksRequest{UserId: registrationResponse.GetUserId()})
	require.NoError(t, err)
	require.Equal(t, 0, len(getTasksAfterDeleteResponse.GetTasks()))
}
