package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/fraud-detection-system/backend/internal/models"
)

type RBIRepository struct {
	db *sql.DB
}

func NewRBIRepository(db *sql.DB) *RBIRepository {
	return &RBIRepository{db: db}
}

// CreateCircular creates a new RBI circular
func (r *RBIRepository) CreateCircular(ctx context.Context, circular *models.RBICircular) error {
	query := `
		INSERT INTO rbi_circulars (id, circular_number, title, content, issued_date, effective_date,
		                           expiry_date, category, keywords, is_active, source_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.ExecContext(ctx, query,
		circular.ID, circular.CircularNumber, circular.Title, circular.Content, circular.IssuedDate,
		circular.EffectiveDate, circular.ExpiryDate, circular.Category, pq.Array(circular.Keywords),
		circular.IsActive, circular.SourceURL, circular.CreatedAt, circular.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create circular: %w", err)
	}
	return nil
}

// GetCircularByID retrieves a circular by ID
func (r *RBIRepository) GetCircularByID(ctx context.Context, id uuid.UUID) (*models.RBICircular, error) {
	query := `
		SELECT id, circular_number, title, content, issued_date, effective_date, expiry_date,
		       category, keywords, is_active, source_url, created_at, updated_at
		FROM rbi_circulars
		WHERE id = $1
	`
	var circular models.RBICircular
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&circular.ID, &circular.CircularNumber, &circular.Title, &circular.Content, &circular.IssuedDate,
		&circular.EffectiveDate, &circular.ExpiryDate, &circular.Category, pq.Array(&circular.Keywords),
		&circular.IsActive, &circular.SourceURL, &circular.CreatedAt, &circular.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("circular not found")
		}
		return nil, fmt.Errorf("failed to get circular: %w", err)
	}
	return &circular, nil
}

// GetActiveCirculars retrieves all active circulars
func (r *RBIRepository) GetActiveCirculars(ctx context.Context) ([]*models.RBICircular, error) {
	query := `
		SELECT id, circular_number, title, content, issued_date, effective_date, expiry_date,
		       category, keywords, is_active, source_url, created_at, updated_at
		FROM rbi_circulars
		WHERE is_active = true
		ORDER BY issued_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get circulars: %w", err)
	}
	defer rows.Close()

	var circulars []*models.RBICircular
	for rows.Next() {
		var circular models.RBICircular
		err := rows.Scan(
			&circular.ID, &circular.CircularNumber, &circular.Title, &circular.Content, &circular.IssuedDate,
			&circular.EffectiveDate, &circular.ExpiryDate, &circular.Category, pq.Array(&circular.Keywords),
			&circular.IsActive, &circular.SourceURL, &circular.CreatedAt, &circular.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan circular: %w", err)
		}
		circulars = append(circulars, &circular)
	}

	return circulars, nil
}

// SearchCircularsByKeywords searches circulars by keywords
func (r *RBIRepository) SearchCircularsByKeywords(ctx context.Context, keywords []string) ([]*models.RBICircular, error) {
	query := `
		SELECT id, circular_number, title, content, issued_date, effective_date, expiry_date,
		       category, keywords, is_active, source_url, created_at, updated_at
		FROM rbi_circulars
		WHERE is_active = true AND (
			content ILIKE ANY($1) OR
			title ILIKE ANY($1) OR
			keywords && $2
		)
		ORDER BY issued_date DESC
	`
	
	// Prepare LIKE patterns
	likePatterns := make([]string, len(keywords))
	for i, kw := range keywords {
		likePatterns[i] = "%" + strings.ToLower(kw) + "%"
	}

	rows, err := r.db.QueryContext(ctx, query, pq.Array(likePatterns), pq.Array(keywords))
	if err != nil {
		return nil, fmt.Errorf("failed to search circulars: %w", err)
	}
	defer rows.Close()

	var circulars []*models.RBICircular
	for rows.Next() {
		var circular models.RBICircular
		err := rows.Scan(
			&circular.ID, &circular.CircularNumber, &circular.Title, &circular.Content, &circular.IssuedDate,
			&circular.EffectiveDate, &circular.ExpiryDate, &circular.Category, pq.Array(&circular.Keywords),
			&circular.IsActive, &circular.SourceURL, &circular.CreatedAt, &circular.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan circular: %w", err)
		}
		circulars = append(circulars, &circular)
	}

	return circulars, nil
}

// CreateSenderRegistry creates a new sender registry entry
func (r *RBIRepository) CreateSenderRegistry(ctx context.Context, sender *models.SenderRegistry) error {
	query := `
		INSERT INTO sender_registry (id, sender_id, bank_name, bank_code, is_verified, is_active,
		                             verified_by, telecom_operator, registration_date, last_verified_at,
		                             reputation_score, message_count, fraud_report_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	_, err := r.db.ExecContext(ctx, query,
		sender.ID, sender.SenderID, sender.BankName, sender.BankCode, sender.IsVerified, sender.IsActive,
		sender.VerifiedBy, sender.TelecomOperator, sender.RegistrationDate, sender.LastVerifiedAt,
		sender.ReputationScore, sender.MessageCount, sender.FraudReportCount, sender.CreatedAt, sender.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create sender registry: %w", err)
	}
	return nil
}

// GetSenderBySenderID retrieves a sender by sender ID
func (r *RBIRepository) GetSenderBySenderID(ctx context.Context, senderID string) (*models.SenderRegistry, error) {
	query := `
		SELECT id, sender_id, bank_name, bank_code, is_verified, is_active, verified_by,
		       telecom_operator, registration_date, last_verified_at, reputation_score,
		       message_count, fraud_report_count, created_at, updated_at
		FROM sender_registry
		WHERE sender_id = $1
	`
	var sender models.SenderRegistry
	err := r.db.QueryRowContext(ctx, query, senderID).Scan(
		&sender.ID, &sender.SenderID, &sender.BankName, &sender.BankCode, &sender.IsVerified,
		&sender.IsActive, &sender.VerifiedBy, &sender.TelecomOperator, &sender.RegistrationDate,
		&sender.LastVerifiedAt, &sender.ReputationScore, &sender.MessageCount, &sender.FraudReportCount,
		&sender.CreatedAt, &sender.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sender not found")
		}
		return nil, fmt.Errorf("failed to get sender: %w", err)
	}
	return &sender, nil
}

// UpdateSenderStats updates sender statistics
func (r *RBIRepository) UpdateSenderStats(ctx context.Context, senderID string, isFraud bool) error {
	query := `
		UPDATE sender_registry
		SET message_count = message_count + 1,
		    fraud_report_count = fraud_report_count + $2,
		    reputation_score = CASE 
		        WHEN $2 > 0 THEN GREATEST(reputation_score - 0.1, 0)
		        ELSE LEAST(reputation_score + 0.01, 1.0)
		    END,
		    updated_at = NOW()
		WHERE sender_id = $1
	`
	fraudIncrement := 0
	if isFraud {
		fraudIncrement = 1
	}
	_, err := r.db.ExecContext(ctx, query, senderID, fraudIncrement)
	if err != nil {
		return fmt.Errorf("failed to update sender stats: %w", err)
	}
	return nil
}

// GetVerifiedSenders retrieves all verified senders
func (r *RBIRepository) GetVerifiedSenders(ctx context.Context) ([]*models.SenderRegistry, error) {
	query := `
		SELECT id, sender_id, bank_name, bank_code, is_verified, is_active, verified_by,
		       telecom_operator, registration_date, last_verified_at, reputation_score,
		       message_count, fraud_report_count, created_at, updated_at
		FROM sender_registry
		WHERE is_verified = true AND is_active = true
		ORDER BY reputation_score DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get senders: %w", err)
	}
	defer rows.Close()

	var senders []*models.SenderRegistry
	for rows.Next() {
		var sender models.SenderRegistry
		err := rows.Scan(
			&sender.ID, &sender.SenderID, &sender.BankName, &sender.BankCode, &sender.IsVerified,
			&sender.IsActive, &sender.VerifiedBy, &sender.TelecomOperator, &sender.RegistrationDate,
			&sender.LastVerifiedAt, &sender.ReputationScore, &sender.MessageCount, &sender.FraudReportCount,
			&sender.CreatedAt, &sender.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sender: %w", err)
		}
		senders = append(senders, &sender)
	}

	return senders, nil
}

