package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the notification-service
type Config struct {
	// Kafka configuration
	KafkaBrokers  []string
	KafkaTopic    string
	KafkaGroupID  string
	KafkaClientID string

	// Telegram configuration
	TelegramBotToken string
	TelegramChatID   int64

	// Service configuration
	ServiceName string
	LogLevel    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Parse Telegram Chat ID
	chatIDStr := getEnv("TELEGRAM_CHAT_ID", "")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil && chatIDStr != "" {
		return nil, fmt.Errorf("invalid TELEGRAM_CHAT_ID: %v", err)
	}

	config := &Config{
		// Kafka
		KafkaBrokers:  strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaTopic:    getEnv("KAFKA_TOPIC", "task-events"),
		KafkaGroupID:  getEnv("KAFKA_GROUP_ID", "notification-service-group"),
		KafkaClientID: getEnv("KAFKA_CLIENT_ID", "notification-service"),

		// Telegram
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:   chatID,

		// Service
		ServiceName: getEnv("SERVICE_NAME", "notification-service"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}

	// Validate required fields
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks if all required configuration fields are set
func (c *Config) Validate() error {
	if len(c.KafkaBrokers) == 0 {
		return fmt.Errorf("KAFKA_BROKERS is required")
	}
	if c.KafkaTopic == "" {
		return fmt.Errorf("KAFKA_TOPIC is required")
	}
	if c.KafkaGroupID == "" {
		return fmt.Errorf("KAFKA_GROUP_ID is required")
	}
	if c.TelegramBotToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	if c.TelegramChatID == 0 {
		return fmt.Errorf("TELEGRAM_CHAT_ID is required")
	}

	return nil
}

// Print logs the current configuration (without sensitive data)
func (c *Config) Print() {
	log.Println("=== Notification Service Configuration ===")
	log.Printf("Service Name: %s", c.ServiceName)
	log.Printf("Log Level: %s", c.LogLevel)
	log.Printf("Kafka Brokers: %v", c.KafkaBrokers)
	log.Printf("Kafka Topic: %s", c.KafkaTopic)
	log.Printf("Kafka Group ID: %s", c.KafkaGroupID)
	log.Printf("Telegram Bot Token: %s", maskToken(c.TelegramBotToken))
	log.Printf("Telegram Chat ID: %d", c.TelegramChatID)
	log.Println("==========================================")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// maskToken masks sensitive token for logging
func maskToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "****" + token[len(token)-4:]
}
