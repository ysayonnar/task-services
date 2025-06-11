package models

import "time"

type Task struct {
	TaskId       int64
	UserId       int64
	Title        string
	Description  string
	Deadline     time.Time
	IsNotificate bool
	CreatedAt    time.Time
}

type Category struct {
	CategoryId int64
	Name       string
}
