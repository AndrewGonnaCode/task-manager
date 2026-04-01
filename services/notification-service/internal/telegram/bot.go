package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot wraps the Telegram bot client
type Bot struct {
	api    *tgbotapi.BotAPI
	chatID int64
}

// NewBot creates a new Telegram bot client
func NewBot(token string, chatID int64) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	// Set debug mode based on log level (optional)
	bot.Debug = false

	log.Printf("Telegram bot initialized successfully (@%s)", bot.Self.UserName)

	return &Bot{
		api:    bot,
		chatID: chatID,
	}, nil
}

// SendMessage sends a text message to the configured chat
func (b *Bot) SendMessage(text string) error {
	msg := tgbotapi.NewMessage(b.chatID, text)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true

	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send Telegram message: %w", err)
	}

	log.Printf("Telegram message sent successfully to chat %d", b.chatID)
	return nil
}

// SendFormattedMessage sends a formatted message with HTML parsing
func (b *Bot) SendFormattedMessage(text string) error {
	msg := tgbotapi.NewMessage(b.chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true

	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send formatted Telegram message: %w", err)
	}

	log.Printf("Formatted Telegram message sent successfully to chat %d", b.chatID)
	return nil
}
