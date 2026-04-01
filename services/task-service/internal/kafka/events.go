package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"task-service/internal/models"
)

// EventPublisher provides methods to publish different types of task events
type EventPublisher struct {
	producer *Producer
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher(producer *Producer) *EventPublisher {
	return &EventPublisher{
		producer: producer,
	}
}

// PublishTaskCreated publishes a task.created event
func (ep *EventPublisher) PublishTaskCreated(task models.Task, correlationID string) error {
	event := models.TaskEvent{
		EventID:       uuid.New().String(),
		EventType:     models.TaskCreated,
		Timestamp:     time.Now(),
		CorrelationID: correlationID,
		Task:          models.TaskDataFromTask(task),
	}

	key := fmt.Sprintf("%d", task.ID)
	if err := ep.producer.PublishEvent(key, event); err != nil {
		log.Printf("ERROR: Failed to publish task.created event: %v", err)
		return err
	}

	log.Printf("Published task.created event (task_id: %d, event_id: %s)", task.ID, event.EventID)
	return nil
}

// PublishTaskUpdated publishes a task.updated event
func (ep *EventPublisher) PublishTaskUpdated(task models.Task, correlationID string) error {
	event := models.TaskEvent{
		EventID:       uuid.New().String(),
		EventType:     models.TaskUpdated,
		Timestamp:     time.Now(),
		CorrelationID: correlationID,
		Task:          models.TaskDataFromTask(task),
	}

	key := fmt.Sprintf("%d", task.ID)
	if err := ep.producer.PublishEvent(key, event); err != nil {
		log.Printf("ERROR: Failed to publish task.updated event: %v", err)
		return err
	}

	log.Printf("Published task.updated event (task_id: %d, event_id: %s)", task.ID, event.EventID)
	return nil
}

// PublishTaskDeleted publishes a task.deleted event
func (ep *EventPublisher) PublishTaskDeleted(task models.Task, correlationID string) error {
	event := models.TaskEvent{
		EventID:       uuid.New().String(),
		EventType:     models.TaskDeleted,
		Timestamp:     time.Now(),
		CorrelationID: correlationID,
		Task:          models.TaskDataFromTask(task),
	}

	key := fmt.Sprintf("%d", task.ID)
	if err := ep.producer.PublishEvent(key, event); err != nil {
		log.Printf("ERROR: Failed to publish task.deleted event: %v", err)
		return err
	}

	log.Printf("Published task.deleted event (task_id: %d, event_id: %s)", task.ID, event.EventID)
	return nil
}
