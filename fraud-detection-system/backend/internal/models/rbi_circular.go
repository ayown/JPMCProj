package models

import (
	"time"

	"github.com/google/uuid"
)

type RBICircular struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	CircularNumber  string     `json:"circular_number" db:"circular_number"`
	Title           string     `json:"title" db:"title"`
	Content         string     `json:"content" db:"content"`
	IssuedDate      time.Time  `json:"issued_date" db:"issued_date"`
	EffectiveDate   *time.Time `json:"effective_date,omitempty" db:"effective_date"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty" db:"expiry_date"`
	Category        string     `json:"category" db:"category"` // KYC, SECURITY, COMPLIANCE, etc.
	Keywords        []string   `json:"keywords" db:"keywords"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	SourceURL       string     `json:"source_url" db:"source_url"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type SenderRegistry struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	SenderID        string     `json:"sender_id" db:"sender_id"` // e.g., AX-HDFC
	BankName        string     `json:"bank_name" db:"bank_name"`
	BankCode        string     `json:"bank_code" db:"bank_code"`
	IsVerified      bool       `json:"is_verified" db:"is_verified"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	VerifiedBy      string     `json:"verified_by" db:"verified_by"` // TELECOM, BANK, MANUAL
	TelecomOperator *string    `json:"telecom_operator,omitempty" db:"telecom_operator"`
	RegistrationDate time.Time `json:"registration_date" db:"registration_date"`
	LastVerifiedAt  *time.Time `json:"last_verified_at,omitempty" db:"last_verified_at"`
	ReputationScore float64    `json:"reputation_score" db:"reputation_score"`
	MessageCount    int        `json:"message_count" db:"message_count"`
	FraudReportCount int       `json:"fraud_report_count" db:"fraud_report_count"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type RBIComplianceCheck struct {
	IsCompliant     bool     `json:"is_compliant"`
	MatchedCirculars []string `json:"matched_circulars"`
	Keywords        []string `json:"keywords_found"`
	IsCurrentRequest bool    `json:"is_current_request"`
	Explanation     string   `json:"explanation"`
}

type HeaderVerificationResult struct {
	IsVerified      bool    `json:"is_verified"`
	SenderExists    bool    `json:"sender_exists"`
	IsActive        bool    `json:"is_active"`
	BankName        string  `json:"bank_name"`
	ReputationScore float64 `json:"reputation_score"`
	RiskLevel       string  `json:"risk_level"` // LOW, MEDIUM, HIGH
	Explanation     string  `json:"explanation"`
}

