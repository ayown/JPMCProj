package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/fraud-detection-system/backend/internal/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// Create creates a new report
func (r *ReportRepository) Create(ctx context.Context, report *models.Report) error {
	query := `
		INSERT INTO reports (id, user_id, message_id, verification_id, report_type, content,
		                     sender_header, description, status, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		report.ID, report.UserID, report.MessageID, report.VerificationID, report.ReportType,
		report.Content, report.SenderHeader, report.Description, report.Status, report.Priority,
		report.CreatedAt, report.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create report: %w", err)
	}
	return nil
}

// GetByID retrieves a report by ID
func (r *ReportRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Report, error) {
	query := `
		SELECT id, user_id, message_id, verification_id, report_type, content, sender_header,
		       description, status, priority, reviewed_by, reviewed_at, review_notes,
		       created_at, updated_at
		FROM reports
		WHERE id = $1
	`
	var report models.Report
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&report.ID, &report.UserID, &report.MessageID, &report.VerificationID, &report.ReportType,
		&report.Content, &report.SenderHeader, &report.Description, &report.Status, &report.Priority,
		&report.ReviewedBy, &report.ReviewedAt, &report.ReviewNotes, &report.CreatedAt, &report.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("report not found")
		}
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	return &report, nil
}

// GetByUserID retrieves reports by user ID
func (r *ReportRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Report, error) {
	query := `
		SELECT id, user_id, message_id, verification_id, report_type, content, sender_header,
		       description, status, priority, reviewed_by, reviewed_at, review_notes,
		       created_at, updated_at
		FROM reports
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports: %w", err)
	}
	defer rows.Close()

	var reports []*models.Report
	for rows.Next() {
		var report models.Report
		err := rows.Scan(
			&report.ID, &report.UserID, &report.MessageID, &report.VerificationID, &report.ReportType,
			&report.Content, &report.SenderHeader, &report.Description, &report.Status, &report.Priority,
			&report.ReviewedBy, &report.ReviewedAt, &report.ReviewNotes, &report.CreatedAt, &report.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan report: %w", err)
		}
		reports = append(reports, &report)
	}

	return reports, nil
}

// GetByStatus retrieves reports by status
func (r *ReportRepository) GetByStatus(ctx context.Context, status string, limit, offset int) ([]*models.Report, error) {
	query := `
		SELECT id, user_id, message_id, verification_id, report_type, content, sender_header,
		       description, status, priority, reviewed_by, reviewed_at, review_notes,
		       created_at, updated_at
		FROM reports
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports: %w", err)
	}
	defer rows.Close()

	var reports []*models.Report
	for rows.Next() {
		var report models.Report
		err := rows.Scan(
			&report.ID, &report.UserID, &report.MessageID, &report.VerificationID, &report.ReportType,
			&report.Content, &report.SenderHeader, &report.Description, &report.Status, &report.Priority,
			&report.ReviewedBy, &report.ReviewedAt, &report.ReviewNotes, &report.CreatedAt, &report.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan report: %w", err)
		}
		reports = append(reports, &report)
	}

	return reports, nil
}

// Update updates a report
func (r *ReportRepository) Update(ctx context.Context, report *models.Report) error {
	query := `
		UPDATE reports
		SET status = $2, priority = $3, reviewed_by = $4, reviewed_at = $5,
		    review_notes = $6, updated_at = $7
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		report.ID, report.Status, report.Priority, report.ReviewedBy,
		report.ReviewedAt, report.ReviewNotes, report.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update report: %w", err)
	}
	return nil
}

// GetStats retrieves report statistics
func (r *ReportRepository) GetStats(ctx context.Context) (*models.ReportStats, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'PENDING') as pending,
			COUNT(*) FILTER (WHERE status = 'RESOLVED') as resolved,
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '24 hours') as last_24h,
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') as last_7d
		FROM reports
	`
	var stats models.ReportStats
	stats.ByType = make(map[string]int)
	stats.ByPriority = make(map[string]int)

	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalReports,
		&stats.PendingReports,
		&stats.ResolvedReports,
		&stats.Last24Hours,
		&stats.Last7Days,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	// Get by type
	typeQuery := `SELECT report_type, COUNT(*) FROM reports GROUP BY report_type`
	rows, err := r.db.QueryContext(ctx, typeQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get type stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reportType string
		var count int
		if err := rows.Scan(&reportType, &count); err != nil {
			return nil, fmt.Errorf("failed to scan type stats: %w", err)
		}
		stats.ByType[reportType] = count
	}

	// Get by priority
	priorityQuery := `SELECT priority, COUNT(*) FROM reports GROUP BY priority`
	rows, err = r.db.QueryContext(ctx, priorityQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get priority stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var priority string
		var count int
		if err := rows.Scan(&priority, &count); err != nil {
			return nil, fmt.Errorf("failed to scan priority stats: %w", err)
		}
		stats.ByPriority[priority] = count
	}

	return &stats, nil
}

