package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/fraud-detection-system/backend/internal/models"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create creates a new message
func (r *MessageRepository) Create(ctx context.Context, message *models.Message) error {
	query := `
		INSERT INTO messages (id, user_id, content, sender_header, received_at, message_type,
		                      phone_number, has_links, link_count, extracted_urls, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		message.ID, message.UserID, message.Content, message.SenderHeader, message.ReceivedAt,
		message.MessageType, message.PhoneNumber, message.HasLinks, message.LinkCount,
		pq.Array(message.ExtractedURLs), message.CreatedAt, message.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

// GetByID retrieves a message by ID
func (r *MessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	query := `
		SELECT id, user_id, content, sender_header, received_at, message_type,
		       phone_number, has_links, link_count, extracted_urls, created_at, updated_at
		FROM messages
		WHERE id = $1
	`
	var message models.Message
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&message.ID, &message.UserID, &message.Content, &message.SenderHeader, &message.ReceivedAt,
		&message.MessageType, &message.PhoneNumber, &message.HasLinks, &message.LinkCount,
		pq.Array(&message.ExtractedURLs), &message.CreatedAt, &message.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	return &message, nil
}

// GetByUserID retrieves messages by user ID
func (r *MessageRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Message, error) {
	query := `
		SELECT id, user_id, content, sender_header, received_at, message_type,
		       phone_number, has_links, link_count, extracted_urls, created_at, updated_at
		FROM messages
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(
			&message.ID, &message.UserID, &message.Content, &message.SenderHeader, &message.ReceivedAt,
			&message.MessageType, &message.PhoneNumber, &message.HasLinks, &message.LinkCount,
			pq.Array(&message.ExtractedURLs), &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

// Delete deletes a message
func (r *MessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}

// DeleteOldMessages deletes messages older than the specified duration
func (r *MessageRepository) DeleteOldMessages(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := `DELETE FROM messages WHERE created_at < $1`
	result, err := r.db.ExecContext(ctx, query, time.Now().Add(-olderThan))
	if err != nil {
		return 0, fmt.Errorf("failed to delete old messages: %w", err)
	}
	return result.RowsAffected()
}

