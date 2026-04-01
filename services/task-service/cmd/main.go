package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"task-service/internal/config"
	"task-service/internal/database"
	"task-service/internal/handlers"
	"task-service/internal/kafka"
	"task-service/internal/middleware"
	"task-service/internal/models"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	log.Println("Hola")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	cfg.Print()

	// Connect to database
	database.Connect()

	// Run auto-migration (GORM will create tables automatically)
	if err := database.DB.AutoMigrate(&models.Task{}); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migration completed")

	// Initialize Kafka producer
	producer, err := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Initialize event publisher
	eventPublisher := kafka.NewEventPublisher(producer)

	// Initialize handlers
	taskHandler := handlers.NewTaskHandler(eventPublisher)

	// Initialize Gin router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": cfg.ServiceName,
			"message": "Task service is running",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Task routes
		tasks := api.Group("/tasks")
		{
			tasks.GET("", taskHandler.GetTasks)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.POST("", taskHandler.CreateTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
		}
	}

	// Setup graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down task-service gracefully...")
		producer.Close()
		os.Exit(0)
	}()

	// Start server
	port := cfg.ServerPort
	log.Printf("Starting %s on port %s...", cfg.ServiceName, port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
