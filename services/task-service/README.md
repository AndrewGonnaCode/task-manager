# Task Service

REST API service for managing tasks with Kafka event publishing.

## Features

- CRUD operations for tasks (Create, Read, Update, Delete)
- PostgreSQL database with GORM
- Kafka event publishing on task operations
- Correlation ID tracking for distributed tracing
- Health check endpoint
- Docker containerization

## API Endpoints

### Health Check
```bash
GET /health
```

### Task Operations

#### Get All Tasks
```bash
GET /api/v1/tasks
```

Response:
```json
{
  "data": [
    {
      "id": 1,
      "title": "Complete project",
      "description": "Finish microservices implementation",
      "completed": false,
      "created_at": "2026-04-01T10:00:00Z",
      "updated_at": "2026-04-01T10:00:00Z"
    }
  ],
  "count": 1
}
```

#### Get Single Task
```bash
GET /api/v1/tasks/:id
```

#### Create Task
```bash
POST /api/v1/tasks
Content-Type: application/json

{
  "title": "New task",
  "description": "Task description"
}
```

#### Update Task
```bash
PUT /api/v1/tasks/:id
Content-Type: application/json

{
  "title": "Updated title",
  "description": "Updated description",
  "completed": true
}
```

#### Delete Task
```bash
DELETE /api/v1/tasks/:id
```

## Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5435
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todo_db
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=task-events
KAFKA_CLIENT_ID=task-service

# Service Configuration
SERVICE_NAME=task-service
LOG_LEVEL=info
```

## Kafka Events

The service publishes events to Kafka on every task operation:

- **task.created** - When a task is created
- **task.updated** - When a task is updated
- **task.deleted** - When a task is deleted

Event format:
```json
{
  "event_id": "uuid-v4",
  "event_type": "task.created",
  "timestamp": "2026-04-01T10:30:00Z",
  "correlation_id": "request-trace-id",
  "task": {
    "id": 1,
    "title": "Task title",
    "description": "Task description",
    "completed": false,
    "created_at": "2026-04-01T10:30:00Z",
    "updated_at": "2026-04-01T10:30:00Z"
  }
}
```

## Running Locally

### With Docker Compose (Recommended)
```bash
cd ../..
docker-compose up -d task-service
```

### Without Docker
```bash
# Install dependencies
go mod download

# Run the service
go run cmd/main.go
```

## Running Tests
```bash
go test -v ./...
```

## Building
```bash
go build -o task-service ./cmd/main.go
```

## Docker

### Build Image
```bash
docker build -t task-service .
```

### Run Container
```bash
docker run -p 8080:8080 \
  -e DB_HOST=localhost \
  -e KAFKA_BROKERS=localhost:9092 \
  task-service
```

## Architecture

The task-service follows a layered architecture:

- **cmd/** - Application entry point
- **internal/models/** - Data models and domain logic
- **internal/handlers/** - HTTP request handlers
- **internal/database/** - Database connection and operations
- **internal/kafka/** - Kafka producer and event publishing
- **internal/middleware/** - HTTP middleware (CORS, logging)
- **internal/config/** - Configuration management

## Resilience

- Kafka failures don't affect API responses (logged but not blocking)
- Database connection pooling
- Health checks for monitoring
- Graceful shutdown on SIGINT/SIGTERM
