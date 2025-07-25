package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func IsEmailValid(email string) bool {
	if len(email) < 3 || len(email) > 200 {
		return false
	}

	idx := strings.Index(email, "@")
	return idx > 0 && idx < len(email)-1
}

func GenerateJwt(userId int64, SECRET string, tokenTTL string) (string, error) {
	const op = "token.New()"

	duration, _ := time.ParseDuration(tokenTTL)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(duration).Unix(),
	})

	signedToken, err := token.SignedString([]byte(SECRET))
	if err != nil {
		return "", fmt.Errorf("op: %s, err: %w", op, err)
	}
	return signedToken, nil
}
