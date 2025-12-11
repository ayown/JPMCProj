package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/fraud-detection-system/backend/internal/config"
	"github.com/fraud-detection-system/backend/internal/utils"
)

type MessageHandler func(ctx context.Context, message *QueueMessage) error

type Consumer struct {
	reader  *kafka.Reader
	handler MessageHandler
	config  *config.Config
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.Config, topic string, handler MessageHandler) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Kafka.Brokers,
		Topic:          topic,
		GroupID:        cfg.Kafka.ConsumerGroup,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.FirstOffset,
		MaxWait:        500 * time.Millisecond,
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
		config:  cfg,
	}
}

// Start starts consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	utils.GetLogger().Info("Starting Kafka consumer")

	for {
		select {
		case <-ctx.Done():
			utils.GetLogger().Info("Stopping Kafka consumer")
			return nil
		default:
			if err := c.processMessage(ctx); err != nil {
				utils.GetLogger().WithError(err).Error("Failed to process message")
				// Continue processing other messages
			}
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context) error {
	kafkaMsg, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch message: %w", err)
	}

	message, err := FromJSON(kafkaMsg.Value)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to unmarshal message")
		// Commit the message to skip it
		_ = c.reader.CommitMessages(ctx, kafkaMsg)
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	utils.GetLogger().WithField("message_id", message.ID).Info("Processing message")

	// Process the message
	if err := c.handler(ctx, message); err != nil {
		utils.GetLogger().WithError(err).WithField("message_id", message.ID).Error("Handler failed to process message")
		
		// Check if we can retry
		if message.CanRetry() {
			message.IncrementRetry()
			utils.GetLogger().WithField("message_id", message.ID).WithField("retry", message.Retry).Info("Retrying message")
			// In a real implementation, you might want to republish to a retry topic
			// For now, we'll just commit and move on
		}
		
		// Commit the message even if processing failed to avoid infinite retries
		if err := c.reader.CommitMessages(ctx, kafkaMsg); err != nil {
			return fmt.Errorf("failed to commit message: %w", err)
		}
		return err
	}

	// Commit the message
	if err := c.reader.CommitMessages(ctx, kafkaMsg); err != nil {
		return fmt.Errorf("failed to commit message: %w", err)
	}

	utils.GetLogger().WithField("message_id", message.ID).Info("Message processed successfully")
	return nil
}

// Close closes the Kafka reader
func (c *Consumer) Close() error {
	return c.reader.Close()
}

