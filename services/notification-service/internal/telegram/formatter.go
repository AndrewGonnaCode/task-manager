package telegram

import (
	"fmt"
	"time"

	"notification-service/internal/models"
)

// FormatTaskEvent formats a TaskEvent into a human-readable message
func FormatTaskEvent(event models.TaskEvent) string {
	var emoji string
	var action string

	switch event.EventType {
	case models.TaskCreated:
		emoji = "✅"
		action = "New task created"
	case models.TaskUpdated:
		emoji = "📝"
		action = "Task updated"
	case models.TaskDeleted:
		emoji = "❌"
		action = "Task deleted"
	default:
		emoji = "ℹ️"
		action = "Task event"
	}

	completedStatus := "No"
	if event.Task.Completed {
		completedStatus = "Yes"
	}

	message := fmt.Sprintf(
		"%s *%s*\n\n"+
			"*Title:* %s\n"+
			"*Description:* %s\n"+
			"*Completed:* %s\n"+
			"*Task ID:* %d\n"+
			"*Timestamp:* %s\n"+
			"*Event ID:* `%s`",
		emoji,
		action,
		escapeMarkdown(event.Task.Title),
		escapeMarkdown(event.Task.Description),
		completedStatus,
		event.Task.ID,
		event.Timestamp.Format(time.RFC3339),
		event.EventID,
	)

	return message
}

// FormatTaskEventHTML formats a TaskEvent into HTML format
func FormatTaskEventHTML(event models.TaskEvent) string {
	var emoji string
	var action string

	switch event.EventType {
	case models.TaskCreated:
		emoji = "✅"
		action = "New task created"
	case models.TaskUpdated:
		emoji = "📝"
		action = "Task updated"
	case models.TaskDeleted:
		emoji = "❌"
		action = "Task deleted"
	default:
		emoji = "ℹ️"
		action = "Task event"
	}

	completedStatus := "No"
	if event.Task.Completed {
		completedStatus = "Yes"
	}

	message := fmt.Sprintf(
		"%s <b>%s</b>\n\n"+
			"<b>Title:</b> %s\n"+
			"<b>Description:</b> %s\n"+
			"<b>Completed:</b> %s\n"+
			"<b>Task ID:</b> %d\n"+
			"<b>Timestamp:</b> %s\n"+
			"<b>Event ID:</b> <code>%s</code>",
		emoji,
		action,
		escapeHTML(event.Task.Title),
		escapeHTML(event.Task.Description),
		completedStatus,
		event.Task.ID,
		event.Timestamp.Format(time.RFC3339),
		event.EventID,
	)

	return message
}

// escapeMarkdown escapes special characters for Markdown
func escapeMarkdown(text string) string {
	if text == "" {
		return "_empty_"
	}
	// Simple escape - for production, use a proper markdown escaper
	return text
}

// escapeHTML escapes special characters for HTML
func escapeHTML(text string) string {
	if text == "" {
		return "<i>empty</i>"
	}
	// Basic HTML escaping
	replacer := map[rune]string{
		'<':  "&lt;",
		'>':  "&gt;",
		'&':  "&amp;",
		'"':  "&quot;",
		'\'': "&#39;",
	}

	result := ""
	for _, char := range text {
		if escaped, ok := replacer[char]; ok {
			result += escaped
		} else {
			result += string(char)
		}
	}

	return result
}
