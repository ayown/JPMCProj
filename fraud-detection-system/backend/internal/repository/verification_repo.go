package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/fraud-detection-system/backend/internal/models"
)

type VerificationRepository struct {
	db *sql.DB
}

func NewVerificationRepository(db *sql.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

// Create creates a new verification record
func (r *VerificationRepository) Create(ctx context.Context, verification *models.Verification) error {
	query := `
		INSERT INTO verifications (id, message_id, user_id, is_fraud, fraud_score, fraud_type,
		                           confidence, model_version, ml_predictions, header_verified,
		                           header_score, rbi_compliant, rbi_verification_result,
		                           explanation, recommendations, processing_time_ms, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`
	_, err := r.db.ExecContext(ctx, query,
		verification.ID, verification.MessageID, verification.UserID, verification.IsFraud,
		verification.FraudScore, verification.FraudType, verification.Confidence, verification.ModelVersion,
		verification.MLPredictions, verification.HeaderVerified, verification.HeaderScore,
		verification.RBICompliant, verification.RBIVerificationResult, verification.Explanation,
		verification.Recommendations, verification.ProcessingTimeMs, verification.CreatedAt, verification.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create verification: %w", err)
	}
	return nil
}

// GetByID retrieves a verification by ID
func (r *VerificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Verification, error) {
	query := `
		SELECT id, message_id, user_id, is_fraud, fraud_score, fraud_type, confidence,
		       model_version, ml_predictions, header_verified, header_score, rbi_compliant,
		       rbi_verification_result, explanation, recommendations, processing_time_ms,
		       created_at, updated_at
		FROM verifications
		WHERE id = $1
	`
	var verification models.Verification
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&verification.ID, &verification.MessageID, &verification.UserID, &verification.IsFraud,
		&verification.FraudScore, &verification.FraudType, &verification.Confidence, &verification.ModelVersion,
		&verification.MLPredictions, &verification.HeaderVerified, &verification.HeaderScore,
		&verification.RBICompliant, &verification.RBIVerificationResult, &verification.Explanation,
		&verification.Recommendations, &verification.ProcessingTimeMs, &verification.CreatedAt, &verification.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("verification not found")
		}
		return nil, fmt.Errorf("failed to get verification: %w", err)
	}
	return &verification, nil
}

// GetByMessageID retrieves a verification by message ID
func (r *VerificationRepository) GetByMessageID(ctx context.Context, messageID uuid.UUID) (*models.Verification, error) {
	query := `
		SELECT id, message_id, user_id, is_fraud, fraud_score, fraud_type, confidence,
		       model_version, ml_predictions, header_verified, header_score, rbi_compliant,
		       rbi_verification_result, explanation, recommendations, processing_time_ms,
		       created_at, updated_at
		FROM verifications
		WHERE message_id = $1
	`
	var verification models.Verification
	err := r.db.QueryRowContext(ctx, query, messageID).Scan(
		&verification.ID, &verification.MessageID, &verification.UserID, &verification.IsFraud,
		&verification.FraudScore, &verification.FraudType, &verification.Confidence, &verification.ModelVersion,
		&verification.MLPredictions, &verification.HeaderVerified, &verification.HeaderScore,
		&verification.RBICompliant, &verification.RBIVerificationResult, &verification.Explanation,
		&verification.Recommendations, &verification.ProcessingTimeMs, &verification.CreatedAt, &verification.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("verification not found")
		}
		return nil, fmt.Errorf("failed to get verification: %w", err)
	}
	return &verification, nil
}

// GetByUserID retrieves verifications by user ID
func (r *VerificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Verification, error) {
	query := `
		SELECT id, message_id, user_id, is_fraud, fraud_score, fraud_type, confidence,
		       model_version, ml_predictions, header_verified, header_score, rbi_compliant,
		       rbi_verification_result, explanation, recommendations, processing_time_ms,
		       created_at, updated_at
		FROM verifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get verifications: %w", err)
	}
	defer rows.Close()

	var verifications []*models.Verification
	for rows.Next() {
		var verification models.Verification
		err := rows.Scan(
			&verification.ID, &verification.MessageID, &verification.UserID, &verification.IsFraud,
			&verification.FraudScore, &verification.FraudType, &verification.Confidence, &verification.ModelVersion,
			&verification.MLPredictions, &verification.HeaderVerified, &verification.HeaderScore,
			&verification.RBICompliant, &verification.RBIVerificationResult, &verification.Explanation,
			&verification.Recommendations, &verification.ProcessingTimeMs, &verification.CreatedAt, &verification.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan verification: %w", err)
		}
		verifications = append(verifications, &verification)
	}

	return verifications, nil
}

// GetStats retrieves verification statistics
func (r *VerificationRepository) GetStats(ctx context.Context, userID *uuid.UUID) (*models.VerificationStats, error) {
	var query string
	var args []interface{}

	if userID != nil {
		query = `
			SELECT 
				COUNT(*) as total,
				COUNT(*) FILTER (WHERE is_fraud = true) as fraud_detected,
				AVG(fraud_score) as avg_fraud_score,
				AVG(processing_time_ms) as avg_processing_time,
				COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '24 hours') as last_24h,
				COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') as last_7d
			FROM verifications
			WHERE user_id = $1
		`
		args = append(args, userID)
	} else {
		query = `
			SELECT 
				COUNT(*) as total,
				COUNT(*) FILTER (WHERE is_fraud = true) as fraud_detected,
				AVG(fraud_score) as avg_fraud_score,
				AVG(processing_time_ms) as avg_processing_time,
				COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '24 hours') as last_24h,
				COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') as last_7d
			FROM verifications
		`
	}

	var stats models.VerificationStats
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&stats.TotalVerifications,
		&stats.FraudDetected,
		&stats.AvgFraudScore,
		&stats.AvgProcessingTime,
		&stats.Last24Hours,
		&stats.Last7Days,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	if stats.TotalVerifications > 0 {
		stats.FraudRate = float64(stats.FraudDetected) / float64(stats.TotalVerifications)
	}

	return &stats, nil
}

// DeleteOldVerifications deletes verifications older than the specified duration
func (r *VerificationRepository) DeleteOldVerifications(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := `DELETE FROM verifications WHERE created_at < $1 AND is_fraud = false`
	result, err := r.db.ExecContext(ctx, query, time.Now().Add(-olderThan))
	if err != nil {
		return 0, fmt.Errorf("failed to delete old verifications: %w", err)
	}
	return result.RowsAffected()
}

