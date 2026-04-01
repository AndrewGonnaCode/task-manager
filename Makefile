.PHONY: help build test lint up down logs clean restart status

# Default target
help:
	@echo "Available commands:"
	@echo "  make build    - Build all services"
	@echo "  make test     - Run tests for all services"
	@echo "  make lint     - Run linters for all services"
	@echo "  make up       - Start all services with docker-compose"
	@echo "  make down     - Stop all services"
	@echo "  make logs     - View logs from all services"
	@echo "  make restart  - Restart all services"
	@echo "  make status   - Show status of all services"
	@echo "  make clean    - Clean build artifacts and volumes"

# Build all services
build:
	@echo "Building task-service..."
	@cd services/task-service && go build -o ../../bin/task-service ./cmd/main.go
	@echo "Building notification-service..."
	@cd services/notification-service && go build -o ../../bin/notification-service ./cmd/main.go
	@echo "Build completed successfully!"

# Run tests for all services
test:
	@echo "Running tests for task-service..."
	@cd services/task-service && go test -v ./...
	@echo "Running tests for notification-service..."
	@cd services/notification-service && go test -v ./...
	@echo "All tests passed!"

# Run linters
lint:
	@echo "Running linter for task-service..."
	@cd services/task-service && go vet ./...
	@cd services/task-service && gofmt -l .
	@echo "Running linter for notification-service..."
	@cd services/notification-service && go vet ./...
	@cd services/notification-service && gofmt -l .
	@echo "Linting completed!"

# Start all services
up:
	@echo "Starting all services..."
	@docker-compose up -d
	@echo "Services started. Use 'make logs' to view logs."

# Stop all services
down:
	@echo "Stopping all services..."
	@docker-compose down
	@echo "Services stopped."

# View logs
logs:
	@docker-compose logs -f

# Restart all services
restart:
	@echo "Restarting all services..."
	@docker-compose restart
	@echo "Services restarted."

# Show service status
status:
	@docker-compose ps

# Clean build artifacts and volumes
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf services/task-service/task-service
	@rm -rf services/notification-service/notification-service
	@echo "Cleaning Docker volumes (this will delete all data)..."
	@docker-compose down -v
	@echo "Clean completed!"
