package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"notification-service/internal/models"
)

// EventHandler defines the interface for handling events
type EventHandler interface {
	HandleEvent(event models.TaskEvent) error
}

// Consumer wraps a Kafka consumer group
type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	handler       EventHandler
	topic         string
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, groupID, topic string, handler EventHandler) (*Consumer, error) {
	config := sarama.NewConfig()

	// Consumer configuration
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	// Set client ID
	config.ClientID = "notification-service-consumer"

	// Create consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	log.Printf("Kafka consumer group created (brokers: %v, group: %s, topic: %s)", brokers, groupID, topic)

	return &Consumer{
		consumerGroup: consumerGroup,
		handler:       handler,
		topic:         topic,
	}, nil
}

// Start starts consuming messages from Kafka
func (c *Consumer) Start(ctx context.Context) error {
	consumerHandler := &consumerGroupHandler{
		handler: c.handler,
	}

	// Start consuming
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping consumer...")
			return c.consumerGroup.Close()
		default:
			err := c.consumerGroup.Consume(ctx, []string{c.topic}, consumerHandler)
			if err != nil {
				log.Printf("Error consuming from Kafka: %v", err)
				return err
			}
		}
	}
}

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	log.Println("Closing Kafka consumer...")
	return c.consumerGroup.Close()
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	handler EventHandler
}

// Setup is called at the beginning of a new session
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session setup")
	return nil
}

// Cleanup is called at the end of a session
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session cleanup")
	return nil
}

// ConsumeClaim processes messages from a partition
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// Parse event
		var event models.TaskEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("ERROR: Failed to unmarshal event: %v (message: %s)", err, string(message.Value))
			// Mark message as consumed even on parse error to avoid blocking
			session.MarkMessage(message, "")
			continue
		}

		log.Printf("Received event: %s (task_id: %d, event_id: %s)",
			event.EventType, event.Task.ID, event.EventID)

		// Handle event
		if err := h.handler.HandleEvent(event); err != nil {
			log.Printf("ERROR: Failed to handle event: %v", err)
			// Don't mark message as consumed on handler error - will retry
			continue
		}

		// Mark message as successfully processed
		session.MarkMessage(message, "")
		log.Printf("Event processed successfully (event_id: %s)", event.EventID)
	}

	return nil
}
