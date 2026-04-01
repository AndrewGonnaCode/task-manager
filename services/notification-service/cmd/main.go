package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"notification-service/internal/config"
	"notification-service/internal/kafka"
	"notification-service/internal/telegram"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	cfg.Print()

	// Initialize Telegram bot
	bot, err := telegram.NewBot(cfg.TelegramBotToken, cfg.TelegramChatID)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}

	// Send startup notification
	startupMessage := "🚀 *Notification Service Started*\n\nThe notification service is now running and ready to process task events."
	if err := bot.SendMessage(startupMessage); err != nil {
		log.Printf("Warning: Failed to send startup notification: %v", err)
	}

	// Create event handler
	eventHandler := kafka.NewTelegramEventHandler(bot)

	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaGroupID,
		cfg.KafkaTopic,
		eventHandler,
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consumer in a goroutine
	go func() {
		log.Printf("Starting %s...", cfg.ServiceName)
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	<-sigterm

	// Send shutdown notification
	shutdownMessage := "⏹️ *Notification Service Stopped*\n\nThe notification service has been gracefully shut down."
	if err := bot.SendMessage(shutdownMessage); err != nil {
		log.Printf("Warning: Failed to send shutdown notification: %v", err)
	}

	// Cancel context to stop consumer
	log.Println("Shutting down notification-service gracefully...")
	cancel()

	log.Println("Notification service stopped successfully")
}
