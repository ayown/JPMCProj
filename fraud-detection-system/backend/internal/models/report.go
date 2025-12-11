package models

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	MessageID      *uuid.UUID `json:"message_id,omitempty" db:"message_id"`
	VerificationID *uuid.UUID `json:"verification_id,omitempty" db:"verification_id"`
	ReportType     string     `json:"report_type" db:"report_type"` // FRAUD, FALSE_POSITIVE, FEEDBACK
	Content        string     `json:"content" db:"content"`
	SenderHeader   string     `json:"sender_header" db:"sender_header"`
	Description    string     `json:"description" db:"description"`
	Status         string     `json:"status" db:"status"` // PENDING, REVIEWED, RESOLVED, DISMISSED
	Priority       string     `json:"priority" db:"priority"` // LOW, MEDIUM, HIGH, CRITICAL
	ReviewedBy     *uuid.UUID `json:"reviewed_by,omitempty" db:"reviewed_by"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty" db:"reviewed_at"`
	ReviewNotes    *string    `json:"review_notes,omitempty" db:"review_notes"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

type ReportInput struct {
	MessageID      *uuid.UUID `json:"message_id,omitempty"`
	VerificationID *uuid.UUID `json:"verification_id,omitempty"`
	ReportType     string     `json:"report_type" binding:"required,oneof=FRAUD FALSE_POSITIVE FEEDBACK"`
	Content        string     `json:"content" binding:"required,max=1000"`
	SenderHeader   string     `json:"sender_header" binding:"required"`
	Description    string     `json:"description" binding:"required,max=500"`
}

type ReportResponse struct {
	ID             uuid.UUID  `json:"id"`
	ReportType     string     `json:"report_type"`
	Content        string     `json:"content"`
	SenderHeader   string     `json:"sender_header"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	Priority       string     `json:"priority"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`
	ReviewNotes    *string    `json:"review_notes,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type ReportStats struct {
	TotalReports     int            `json:"total_reports"`
	PendingReports   int            `json:"pending_reports"`
	ResolvedReports  int            `json:"resolved_reports"`
	ByType           map[string]int `json:"by_type"`
	ByPriority       map[string]int `json:"by_priority"`
	Last24Hours      int            `json:"last_24_hours"`
	Last7Days        int            `json:"last_7_days"`
}

