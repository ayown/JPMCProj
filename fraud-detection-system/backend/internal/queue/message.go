package queue

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MessageType defines the type of queue message
type MessageType string

const (
	MessageTypeVerification MessageType = "verification"
	MessageTypeReport       MessageType = "report"
	MessageTypeAlert        MessageType = "alert"
)

// QueueMessage represents a message in the queue
type QueueMessage struct {
	ID        uuid.UUID              `json:"id"`
	Type      MessageType            `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp time.Time              `json:"timestamp"`
	Retry     int                    `json:"retry"`
	MaxRetry  int                    `json:"max_retry"`
}

// NewQueueMessage creates a new queue message
func NewQueueMessage(msgType MessageType, payload map[string]interface{}) *QueueMessage {
	return &QueueMessage{
		ID:        uuid.New(),
		Type:      msgType,
		Payload:   payload,
		Timestamp: time.Now(),
		Retry:     0,
		MaxRetry:  3,
	}
}

// ToJSON converts the message to JSON
func (m *QueueMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON creates a message from JSON
func FromJSON(data []byte) (*QueueMessage, error) {
	var msg QueueMessage
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// CanRetry checks if the message can be retried
func (m *QueueMessage) CanRetry() bool {
	return m.Retry < m.MaxRetry
}

// IncrementRetry increments the retry counter
func (m *QueueMessage) IncrementRetry() {
	m.Retry++
}

