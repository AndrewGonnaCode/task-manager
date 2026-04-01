package kafka

import (
	"fmt"
	"log"

	"notification-service/internal/models"
	"notification-service/internal/telegram"
)

// TelegramEventHandler handles Kafka events and sends notifications to Telegram
type TelegramEventHandler struct {
	telegramBot *telegram.Bot
}

// NewTelegramEventHandler creates a new Telegram event handler
func NewTelegramEventHandler(telegramBot *telegram.Bot) *TelegramEventHandler {
	return &TelegramEventHandler{
		telegramBot: telegramBot,
	}
}

// HandleEvent processes a task event and sends a notification
func (h *TelegramEventHandler) HandleEvent(event models.TaskEvent) error {
	// Log the event
	log.Printf("Handling event: %s (task_id: %d, correlation_id: %s)",
		event.EventType, event.Task.ID, event.CorrelationID)

	// Format message
	message := telegram.FormatTaskEvent(event)

	// Send to Telegram
	if err := h.telegramBot.SendMessage(message); err != nil {
		return fmt.Errorf("failed to send Telegram notification: %w", err)
	}

	log.Printf("Successfully sent Telegram notification for event %s (task_id: %d)",
		event.EventType, event.Task.ID)

	return nil
}
