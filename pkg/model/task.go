package model

import (
	"time"
)

// Task - タスクを表す構造体
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      *string    `json:"status" gorm:"default:'pending'"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Deadline    *time.Time `json:"deadline"`
}
