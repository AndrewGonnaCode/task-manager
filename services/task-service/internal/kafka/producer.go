package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// Producer wraps a Kafka producer client
type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()

	// Producer configuration
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas
	config.Producer.Retry.Max = 3
	config.Producer.Compression = sarama.CompressionSnappy

	// Set client ID for better debugging
	config.ClientID = "task-service-producer"

	// Create sync producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	log.Printf("Kafka producer initialized successfully (brokers: %v, topic: %s)", brokers, topic)

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

// PublishEvent publishes an event to Kafka with the given key
func (p *Producer) PublishEvent(key string, event interface{}) error {
	// Marshal event to JSON
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create Kafka message
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(eventBytes),
	}

	// Send message
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send Kafka message: %w", err)
	}

	log.Printf("Event published successfully (partition: %d, offset: %d, key: %s)", partition, offset, key)
	return nil
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	if p.producer != nil {
		log.Println("Closing Kafka producer...")
		return p.producer.Close()
	}
	return nil
}
