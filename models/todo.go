package models

import "time"

type ToDoList struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at,omitempty"`
	Title           string    `json:"title"`
	CompletePercent int       `json:"percent"`
	Deleted         bool      `json:"deleted"`
}

type ToDoItem struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	ListID    uint      `json:"list_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Task      string    `json:"task"`
	Completed bool      `json:"completed"`
	Deleted   bool      `json:"deleted"`
}
