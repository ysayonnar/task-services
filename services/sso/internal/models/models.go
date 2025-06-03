package models

import "time"

type User struct {
	UserId       int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
