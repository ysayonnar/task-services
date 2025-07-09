package tests

import (
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
	tasks "github.com/ysayonnar/task-contracts/tasks/gen/go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"tasks/tests/suite"
	"testing"
	"time"
)

const (
	passwordDefaultLength    = 14
	descriptionDefaultLength = 5
)

func randomFakePassword(length int) string {
	return gofakeit.Password(true, true, true, true, false, length)
}

func TestTasksHappyPath(t *testing.T) {
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
	createCategoryResponse, err := st.TasksClient.CreateCategory(ctx, &tasks.CreateCategoryRequest{
		Token: loginResponse.GetToken(),
		Name:  categoryName,
	})
	require.NoError(t, err)
	require.NotEmpty(t, createCategoryResponse.GetCategoryId())

	taskTitle := gofakeit.HipsterWord()
	taskDesc := gofakeit.HipsterSentence(descriptionDefaultLength)
	deadline := timestamppb.New(time.Now().Add(24 * time.Hour))

	createTaskResponse, err := st.TasksClient.CreateTask(ctx, &tasks.CreateTaskRequest{
		Token:        loginResponse.GetToken(),
		CategoryId:   createCategoryResponse.GetCategoryId(),
		Title:        taskTitle,
		Description:  taskDesc,
		IsNotificate: false,
		Deadline:     deadline,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, createTaskResponse.GetTaskId())

	getTasksResponse, err := st.TasksClient.GetTasks(ctx, &tasks.GetTasksRequest{Token: loginResponse.GetToken()})
	require.NoError(t, err)
	require.Equal(t, 1, len(getTasksResponse.GetTasks()))
	assert.Equal(t, createCategoryResponse.GetCategoryId(), getTasksResponse.GetTasks()[0].CategoryId)
	assert.Equal(t, taskTitle, getTasksResponse.GetTasks()[0].Title)
	assert.Equal(t, taskDesc, getTasksResponse.GetTasks()[0].Description)
	assert.Equal(t, false, getTasksResponse.GetTasks()[0].IsNotificate)
	assert.Equal(t, deadline.GetSeconds(), getTasksResponse.GetTasks()[0].Deadline.GetSeconds())

	newTitle := gofakeit.HipsterWord()
	newDesc := gofakeit.HipsterSentence(descriptionDefaultLength)
	updateTaskResponse, err := st.TasksClient.UpdateTask(ctx, &tasks.UpdateTaskRequest{
		TaskId:       createTaskResponse.GetTaskId(),
		Token:        loginResponse.GetToken(),
		CategoryId:   createCategoryResponse.GetCategoryId(),
		Title:        newTitle,
		Description:  newDesc,
		IsNotificate: true,
		Deadline:     deadline,
	})
	require.NoError(t, err)
	require.NotEmpty(t, updateTaskResponse.GetTaskId())
	assert.Equal(t, createTaskResponse.GetTaskId(), updateTaskResponse.GetTaskId())

	getTasksByCategoryResponse, err := st.TasksClient.GetTasksByCategory(ctx, &tasks.GetTasksByCategoryRequest{
		Token:      loginResponse.GetToken(),
		CategoryId: createCategoryResponse.GetCategoryId(),
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(getTasksResponse.GetTasks()))
	assert.Equal(t, createCategoryResponse.GetCategoryId(), getTasksByCategoryResponse.GetTasks()[0].CategoryId)
	assert.Equal(t, newTitle, getTasksByCategoryResponse.GetTasks()[0].Title)
	assert.Equal(t, newDesc, getTasksByCategoryResponse.GetTasks()[0].Description)
	assert.Equal(t, true, getTasksByCategoryResponse.GetTasks()[0].IsNotificate)
	assert.Equal(t, deadline.GetSeconds(), getTasksByCategoryResponse.GetTasks()[0].Deadline.GetSeconds())

	deleteTaskResponse, err := st.TasksClient.DeleteTask(ctx, &tasks.DeleteTaskRequest{
		Token:  loginResponse.GetToken(),
		TaskId: createTaskResponse.GetTaskId(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, deleteTaskResponse.GetTaskId())
	assert.Equal(t, createTaskResponse.GetTaskId(), deleteTaskResponse.GetTaskId())
}
