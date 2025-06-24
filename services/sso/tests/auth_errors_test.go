package tests

import (
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sso/tests/suite"
	"testing"
)

func TestUserAlreadyExists(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword(passwordDefaultLength)

	firstRegistrationResponse, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, firstRegistrationResponse.GetUserId())

	secondRegistrationResponse, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))
	assert.Empty(t, secondRegistrationResponse)
}

func TestInvalidEmailAndPassword(t *testing.T) {
	ctx, st := suite.New(t)

	validPassword := randomFakePassword(passwordDefaultLength)
	invalidEmails := []string{"em", "email.com"}
	for _, invalidEmail := range invalidEmails {
		registerResponse, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
			Email:    invalidEmail,
			Password: validPassword,
		})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Empty(t, registerResponse)

		loginResponse, err := st.AuthClient.Login(ctx, &sso.LoginRequest{
			Email:    invalidEmail,
			Password: validPassword,
			AppId:    0,
		})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Empty(t, loginResponse)

		deleteResponse, err := st.AuthClient.Delete(ctx, &sso.DeleteRequest{
			Email:    invalidEmail,
			Password: validPassword,
		})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Empty(t, deleteResponse)
	}

	validEmail := gofakeit.Email()
	invalidPasswords := []string{randomFakePassword(73), randomFakePassword(7)}
	for _, invalidPassword := range invalidPasswords {
		registerResponse, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
			Email:    validEmail,
			Password: invalidPassword,
		})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Empty(t, registerResponse)

		loginResponse, err := st.AuthClient.Login(ctx, &sso.LoginRequest{
			Email:    validEmail,
			Password: invalidPassword,
			AppId:    0,
		})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Empty(t, loginResponse)

		deleteResponse, err := st.AuthClient.Delete(ctx, &sso.DeleteRequest{
			Email:    validEmail,
			Password: invalidPassword,
		})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Empty(t, deleteResponse)
	}
}
