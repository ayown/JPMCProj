package models

import (
	"time"

	"github.com/google/uuid"
)

type Verification struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	MessageID             uuid.UUID  `json:"message_id" db:"message_id"`
	UserID                *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	IsFraud               bool       `json:"is_fraud" db:"is_fraud"`
	FraudScore            float64    `json:"fraud_score" db:"fraud_score"`
	FraudType             *string    `json:"fraud_type,omitempty" db:"fraud_type"`
	Confidence            float64    `json:"confidence" db:"confidence"`
	ModelVersion          string     `json:"model_version" db:"model_version"`
	MLPredictions         string     `json:"ml_predictions" db:"ml_predictions"` // JSON
	HeaderVerified        bool       `json:"header_verified" db:"header_verified"`
	HeaderScore           float64    `json:"header_score" db:"header_score"`
	RBICompliant          bool       `json:"rbi_compliant" db:"rbi_compliant"`
	RBIVerificationResult string     `json:"rbi_verification_result" db:"rbi_verification_result"` // JSON
	Explanation           string     `json:"explanation" db:"explanation"`
	Recommendations       string     `json:"recommendations" db:"recommendations"` // JSON
	ProcessingTimeMs      int        `json:"processing_time_ms" db:"processing_time_ms"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
}

type VerificationRequest struct {
	Content      string     `json:"content" binding:"required,max=1000"`
	SenderHeader string     `json:"sender_header" binding:"required"`
	ReceivedAt   *time.Time `json:"received_at,omitempty"`
	MessageType  string     `json:"message_type" binding:"omitempty,oneof=SMS WhatsApp Email"`
	PhoneNumber  *string    `json:"phone_number,omitempty"`
}

type VerificationResponse struct {
	ID                uuid.UUID              `json:"id"`
	MessageID         uuid.UUID              `json:"message_id"`
	IsFraud           bool                   `json:"is_fraud"`
	FraudScore        float64                `json:"fraud_score"`
	FraudType         *string                `json:"fraud_type,omitempty"`
	Confidence        float64                `json:"confidence"`
	RiskLevel         string                 `json:"risk_level"` // LOW, MEDIUM, HIGH, CRITICAL
	HeaderVerified    bool                   `json:"header_verified"`
	RBICompliant      bool                   `json:"rbi_compliant"`
	Explanation       string                 `json:"explanation"`
	Recommendations   []string               `json:"recommendations"`
	ModelPredictions  map[string]interface{} `json:"model_predictions"`
	ProcessingTimeMs  int                    `json:"processing_time_ms"`
	VerifiedAt        time.Time              `json:"verified_at"`
}

type MLInferenceRequest struct {
	Content      string          `json:"content"`
	SenderHeader string          `json:"sender_header"`
	Features     MessageFeatures `json:"features"`
}

type MLInferenceResponse struct {
	IsFraud           bool                   `json:"is_fraud"`
	FraudScore        float64                `json:"fraud_score"`
	FraudType         string                 `json:"fraud_type"`
	Confidence        float64                `json:"confidence"`
	ModelPredictions  map[string]interface{} `json:"model_predictions"`
	Explanation       string                 `json:"explanation"`
	InferenceTimeMs   int                    `json:"inference_time_ms"`
	ModelVersion      string                 `json:"model_version"`
}

type VerificationStats struct {
	TotalVerifications int     `json:"total_verifications"`
	FraudDetected      int     `json:"fraud_detected"`
	FraudRate          float64 `json:"fraud_rate"`
	AvgFraudScore      float64 `json:"avg_fraud_score"`
	AvgProcessingTime  float64 `json:"avg_processing_time"`
	Last24Hours        int     `json:"last_24_hours"`
	Last7Days          int     `json:"last_7_days"`
}

