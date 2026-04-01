package models

import (
	"time"
)

// Task represents a todo task
type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title" binding:"required"`
	Description string    `json:"description"`
	Completed   bool      `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTaskInput represents the input for creating a task
type CreateTaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// UpdateTaskInput represents the input for updating a task
type UpdateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed"`
}
