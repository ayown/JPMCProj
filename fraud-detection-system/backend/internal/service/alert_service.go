package service

import (
	"context"
	"fmt"

	"github.com/fraud-detection-system/backend/internal/models"
	"github.com/fraud-detection-system/backend/internal/queue"
	"github.com/fraud-detection-system/backend/internal/utils"
)

type AlertService struct {
	producer *queue.Producer
}

func NewAlertService(producer *queue.Producer) *AlertService {
	return &AlertService{
		producer: producer,
	}
}

// SendFraudAlert sends an alert for detected fraud
func (s *AlertService) SendFraudAlert(ctx context.Context, verification *models.Verification, message *models.Message) error {
	payload := map[string]interface{}{
		"verification_id": verification.ID.String(),
		"message_id":      message.ID.String(),
		"fraud_score":     verification.FraudScore,
		"fraud_type":      verification.FraudType,
		"sender_header":   message.SenderHeader,
		"content":         message.Content,
		"timestamp":       verification.CreatedAt,
	}

	if err := s.producer.PublishAlert(ctx, payload); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to send fraud alert")
		return fmt.Errorf("failed to send fraud alert: %w", err)
	}

	utils.GetLogger().WithField("verification_id", verification.ID).Info("Fraud alert sent")
	return nil
}

// SendHighRiskAlert sends an alert for high-risk messages
func (s *AlertService) SendHighRiskAlert(ctx context.Context, verification *models.Verification) error {
	if verification.FraudScore < 0.7 {
		return nil // Only send alerts for high-risk messages
	}

	payload := map[string]interface{}{
		"verification_id": verification.ID.String(),
		"fraud_score":     verification.FraudScore,
		"risk_level":      "HIGH",
		"timestamp":       verification.CreatedAt,
	}

	if err := s.producer.PublishAlert(ctx, payload); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to send high-risk alert")
		return fmt.Errorf("failed to send high-risk alert: %w", err)
	}

	return nil
}

