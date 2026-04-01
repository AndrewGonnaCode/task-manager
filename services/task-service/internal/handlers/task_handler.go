package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"task-service/internal/database"
	"task-service/internal/kafka"
	"task-service/internal/models"
)

// TaskHandler handles HTTP requests for task operations
type TaskHandler struct {
	eventPublisher *kafka.EventPublisher
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(eventPublisher *kafka.EventPublisher) *TaskHandler {
	return &TaskHandler{
		eventPublisher: eventPublisher,
	}
}

// GetTasks retrieves all tasks
func (h *TaskHandler) GetTasks(c *gin.Context) {
	var tasks []models.Task

	if result := database.DB.Find(&tasks); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch tasks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tasks,
		"count": len(tasks),
	})
}

// GetTask retrieves a single task by ID
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if result := database.DB.First(&task, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": task,
	})
}

// CreateTask creates a new task
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var input models.CreateTaskInput

	// Validate request body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create task
	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
	}

	if result := database.DB.Create(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create task",
		})
		return
	}

	// Publish Kafka event (don't fail request if Kafka is down)
	correlationID := c.GetString("correlation_id")
	if err := h.eventPublisher.PublishTaskCreated(task, correlationID); err != nil {
		log.Printf("WARNING: Failed to publish task.created event (task_id: %d): %v", task.ID, err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"data":    task,
	})
}

// UpdateTask updates an existing task
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	// Find task
	if result := database.DB.First(&task, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	// Validate request body
	var input models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update fields
	if input.Title != "" {
		task.Title = input.Title
	}
	if input.Description != "" {
		task.Description = input.Description
	}
	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	// Save updates
	if result := database.DB.Save(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update task",
		})
		return
	}

	// Publish Kafka event (don't fail request if Kafka is down)
	correlationID := c.GetString("correlation_id")
	if err := h.eventPublisher.PublishTaskUpdated(task, correlationID); err != nil {
		log.Printf("WARNING: Failed to publish task.updated event (task_id: %d): %v", task.ID, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"data":    task,
	})
}

// DeleteTask deletes a task
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	// Find task
	if result := database.DB.First(&task, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	// Publish Kafka event BEFORE deleting (don't fail request if Kafka is down)
	correlationID := c.GetString("correlation_id")
	if err := h.eventPublisher.PublishTaskDeleted(task, correlationID); err != nil {
		log.Printf("WARNING: Failed to publish task.deleted event (task_id: %d): %v", task.ID, err)
	}

	// Delete task
	if result := database.DB.Delete(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}
