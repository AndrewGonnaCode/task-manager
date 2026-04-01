# Notification Service

Kafka consumer service that sends task notifications to Telegram.

## Features

- Consumes task events from Kafka
- Sends formatted notifications to Telegram
- Consumer group support for horizontal scaling
- Graceful shutdown handling
- Startup/shutdown notifications

## How It Works

1. Listens to `task-events` Kafka topic
2. Receives events: `task.created`, `task.updated`, `task.deleted`
3. Formats event data into human-readable messages
4. Sends notifications to configured Telegram chat

## Environment Variables

```bash
# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=task-events
KAFKA_GROUP_ID=notification-service-group
KAFKA_CLIENT_ID=notification-service

# Telegram Configuration
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_CHAT_ID=your_chat_id_here

# Service Configuration
SERVICE_NAME=notification-service
LOG_LEVEL=info
```

## Telegram Setup

### Getting Bot Token

1. Open Telegram and search for `@BotFather`
2. Send `/newbot` command
3. Follow the instructions to create your bot
4. Copy the bot token provided
5. Set it in `TELEGRAM_BOT_TOKEN` environment variable

### Getting Chat ID

1. Open Telegram and search for `@userinfobot`
2. Start a chat and send any message
3. Copy the `Id` shown (this is your chat ID)
4. Set it in `TELEGRAM_CHAT_ID` environment variable

### Testing Your Bot

```bash
# Send a test message using curl
curl -X POST "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/sendMessage" \
  -d chat_id=<YOUR_CHAT_ID> \
  -d text="Test message from notification service"
```

## Event Handling

The service handles three types of events:

### task.created
```
✅ New task created

Title: Complete documentation
Description: Write comprehensive README files
Completed: No
Task ID: 1
Timestamp: 2026-04-01T10:30:00Z
Event ID: uuid-v4
```

### task.updated
```
📝 Task updated

Title: Complete documentation
Description: Write comprehensive README files
Completed: Yes
Task ID: 1
Timestamp: 2026-04-01T11:00:00Z
Event ID: uuid-v4
```

### task.deleted
```
❌ Task deleted

Title: Complete documentation
Description: Write comprehensive README files
Completed: Yes
Task ID: 1
Timestamp: 2026-04-01T12:00:00Z
Event ID: uuid-v4
```

## Running Locally

### With Docker Compose (Recommended)
```bash
cd ../..
docker-compose up -d notification-service
```

### Without Docker
```bash
# Install dependencies
go mod download

# Set environment variables
export TELEGRAM_BOT_TOKEN="your_token"
export TELEGRAM_CHAT_ID="your_chat_id"
export KAFKA_BROKERS="localhost:9092"

# Run the service
go run cmd/main.go
```

## Running Tests
```bash
go test -v ./...
```

## Building
```bash
go build -o notification-service ./cmd/main.go
```

## Docker

### Build Image
```bash
docker build -t notification-service .
```

### Run Container
```bash
docker run \
  -e KAFKA_BROKERS=kafka:29092 \
  -e TELEGRAM_BOT_TOKEN=your_token \
  -e TELEGRAM_CHAT_ID=your_chat_id \
  notification-service
```

## Architecture

The notification-service follows a layered architecture:

- **cmd/** - Application entry point
- **internal/models/** - Event data models
- **internal/kafka/** - Kafka consumer and event handler
- **internal/telegram/** - Telegram bot client and message formatting
- **internal/config/** - Configuration management

## Consumer Group

The service uses Kafka consumer groups for:

- **Load balancing** - Multiple instances can process events in parallel
- **Fault tolerance** - If one instance fails, others continue processing
- **Offset management** - Kafka tracks processed messages automatically

## Scaling

To scale horizontally:

```bash
docker-compose up -d --scale notification-service=3
```

All instances will join the same consumer group and share the workload.

## Troubleshooting

### Bot not receiving messages

1. Check if `TELEGRAM_BOT_TOKEN` is correct
2. Verify `TELEGRAM_CHAT_ID` is your personal chat ID
3. Ensure you've started a conversation with the bot
4. Check service logs for errors

### Events not being consumed

1. Verify Kafka is running and accessible
2. Check `KAFKA_BROKERS` configuration
3. Ensure topic `task-events` exists
4. Check consumer group offset position

### Service crashes on startup

1. Verify all environment variables are set
2. Check Kafka connectivity
3. Verify Telegram credentials
4. Review service logs for specific errors
