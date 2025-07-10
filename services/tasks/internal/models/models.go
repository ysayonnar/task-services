package models

import "time"

type Task struct {
	TaskId       int64     `json:"task_id,omitempty"`
	UserId       int64     `json:"user_id,omitempty"`
	Title        string    `json:"title,omitempty"`
	Description  string    `json:"description,omitempty"`
	Deadline     time.Time `json:"deadline,omitempty"`
	IsNotificate bool      `json:"is_notificate,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type Category struct {
	CategoryId int64  `json:"category_id,omitempty"`
	Name       string `json:"name,omitempty"`
}
