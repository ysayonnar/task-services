package tests

import (
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
	"sso/tests/suite"
	"testing"
	"time"
)

const (
	passwordDefaultLength = 14
)

func randomFakePassword(length int) string {
	return gofakeit.Password(true, true, true, true, false, length)
}

func TestRegisterLoginDeleteHappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword(passwordDefaultLength)

	registrationResponse, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, registrationResponse.GetUserId())

	loginResponse, err := st.AuthClient.Login(ctx, &sso.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    0,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := loginResponse.GetToken()
	assert.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(st.Cfg.Secret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, registrationResponse.GetUserId(), int64(claims["user_id"].(float64)))

	const deltaSeconds = 1
	tokenTTL, _ := time.ParseDuration(st.Cfg.TokenTTL)
	assert.InDelta(t, loginTime.Add(tokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)

	deleteResponse, err := st.AuthClient.Delete(ctx, &sso.DeleteRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.True(t, deleteResponse.GetIsDeleted())
}
