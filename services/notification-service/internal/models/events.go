package models

import "time"

// EventType represents the type of task event
type EventType string

const (
	TaskCreated EventType = "task.created"
	TaskUpdated EventType = "task.updated"
	TaskDeleted EventType = "task.deleted"
)

// TaskEvent represents a Kafka event message for task operations
type TaskEvent struct {
	EventID       string    `json:"event_id"`        // UUID for event tracking
	EventType     EventType `json:"event_type"`      // task.created, task.updated, task.deleted
	Timestamp     time.Time `json:"timestamp"`       // Event creation time
	CorrelationID string    `json:"correlation_id"`  // Request trace ID
	Task          TaskData  `json:"task"`            // Task payload
}

// TaskData represents the task information in the event
type TaskData struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
