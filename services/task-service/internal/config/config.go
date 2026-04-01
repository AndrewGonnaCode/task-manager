package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Config holds all configuration for the task-service
type Config struct {
	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Server configuration
	ServerPort string

	// Kafka configuration
	KafkaBrokers  []string
	KafkaTopic    string
	KafkaClientID string

	// Service configuration
	ServiceName string
	LogLevel    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "todo_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Server
		ServerPort: getEnv("SERVER_PORT", "8080"),

		// Kafka
		KafkaBrokers:  strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaTopic:    getEnv("KAFKA_TOPIC", "task-events"),
		KafkaClientID: getEnv("KAFKA_CLIENT_ID", "task-service"),

		// Service
		ServiceName: getEnv("SERVICE_NAME", "task-service"),
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
	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if len(c.KafkaBrokers) == 0 {
		return fmt.Errorf("KAFKA_BROKERS is required")
	}
	if c.KafkaTopic == "" {
		return fmt.Errorf("KAFKA_TOPIC is required")
	}

	return nil
}

// Print logs the current configuration (without sensitive data)
func (c *Config) Print() {
	log.Println("=== Task Service Configuration ===")
	log.Printf("Service Name: %s", c.ServiceName)
	log.Printf("Log Level: %s", c.LogLevel)
	log.Printf("Server Port: %s", c.ServerPort)
	log.Printf("DB Host: %s:%s", c.DBHost, c.DBPort)
	log.Printf("DB Name: %s", c.DBName)
	log.Printf("Kafka Brokers: %v", c.KafkaBrokers)
	log.Printf("Kafka Topic: %s", c.KafkaTopic)
	log.Println("==================================")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
