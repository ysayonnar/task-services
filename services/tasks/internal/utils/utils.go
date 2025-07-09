package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
)

func ParseUserId(bearer string, secret string) (int64, error) {
	const op = "utils.ParseUserId"

	token, err := jwt.Parse(bearer, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("op: %s, err: %w", op, err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("op: %s, err: %w", op, err)
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return int64(userId), nil
}
