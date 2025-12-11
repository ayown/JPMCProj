package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/fraud-detection-system/backend/internal/config"
	"github.com/fraud-detection-system/backend/internal/utils"
)

type Producer struct {
	writers map[string]*kafka.Writer
	config  *config.Config
}

// NewProducer creates a new Kafka producer
func NewProducer(cfg *config.Config) *Producer {
	writers := make(map[string]*kafka.Writer)

	// Create writers for each topic
	topics := []string{
		cfg.Kafka.TopicVerification,
		cfg.Kafka.TopicReports,
		cfg.Kafka.TopicAlerts,
	}

	for _, topic := range topics {
		writer := &kafka.Writer{
			Addr:         kafka.TCP(cfg.Kafka.Brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    100,
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
			Compression:  kafka.Snappy,
		}
		writers[topic] = writer
	}

	return &Producer{
		writers: writers,
		config:  cfg,
	}
}

// Publish publishes a message to a topic
func (p *Producer) Publish(ctx context.Context, topic string, message *QueueMessage) error {
	writer, ok := p.writers[topic]
	if !ok {
		return fmt.Errorf("writer not found for topic: %s", topic)
	}

	data, err := message.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(message.ID.String()),
		Value: data,
		Time:  time.Now(),
	}

	if err := writer.WriteMessages(ctx, kafkaMsg); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to publish message to Kafka")
		return fmt.Errorf("failed to write message: %w", err)
	}

	utils.GetLogger().WithField("topic", topic).WithField("message_id", message.ID).Info("Message published to Kafka")
	return nil
}

// PublishVerification publishes a verification request
func (p *Producer) PublishVerification(ctx context.Context, payload map[string]interface{}) error {
	message := NewQueueMessage(MessageTypeVerification, payload)
	return p.Publish(ctx, p.config.Kafka.TopicVerification, message)
}

// PublishReport publishes a fraud report
func (p *Producer) PublishReport(ctx context.Context, payload map[string]interface{}) error {
	message := NewQueueMessage(MessageTypeReport, payload)
	return p.Publish(ctx, p.config.Kafka.TopicReports, message)
}

// PublishAlert publishes an alert
func (p *Producer) PublishAlert(ctx context.Context, payload map[string]interface{}) error {
	message := NewQueueMessage(MessageTypeAlert, payload)
	return p.Publish(ctx, p.config.Kafka.TopicAlerts, message)
}

// Close closes all Kafka writers
func (p *Producer) Close() error {
	for _, writer := range p.writers {
		if err := writer.Close(); err != nil {
			return err
		}
	}
	return nil
}

