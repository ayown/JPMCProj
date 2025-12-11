package service

import (
	"context"
	"fmt"

	"github.com/fraud-detection-system/backend/internal/models"
	"github.com/fraud-detection-system/backend/internal/repository"
	"github.com/fraud-detection-system/backend/internal/utils"
)

type HeaderVerificationService struct {
	rbiRepo *repository.RBIRepository
}

func NewHeaderVerificationService(rbiRepo *repository.RBIRepository) *HeaderVerificationService {
	return &HeaderVerificationService{
		rbiRepo: rbiRepo,
	}
}

// VerifyHeader verifies the authenticity of a sender header
func (s *HeaderVerificationService) VerifyHeader(ctx context.Context, senderHeader string) (*models.HeaderVerificationResult, error) {
	result := &models.HeaderVerificationResult{
		IsVerified:      false,
		SenderExists:    false,
		IsActive:        false,
		BankName:        "",
		ReputationScore: 0.0,
		RiskLevel:       "HIGH",
		Explanation:     "",
	}

	// Look up sender in registry
	sender, err := s.rbiRepo.GetSenderBySenderID(ctx, senderHeader)
	if err != nil {
		// Sender not found in registry
		result.Explanation = fmt.Sprintf("Sender ID '%s' not found in verified registry. This could be a spoofed header.", senderHeader)
		result.RiskLevel = "CRITICAL"
		return result, nil
	}

	result.SenderExists = true
	result.BankName = sender.BankName
	result.ReputationScore = sender.ReputationScore
	result.IsActive = sender.IsActive

	// Check if sender is verified and active
	if !sender.IsVerified {
		result.Explanation = fmt.Sprintf("Sender ID '%s' exists but is not verified", senderHeader)
		result.RiskLevel = "HIGH"
		return result, nil
	}

	if !sender.IsActive {
		result.Explanation = fmt.Sprintf("Sender ID '%s' is inactive. This could indicate a decommissioned or compromised sender.", senderHeader)
		result.RiskLevel = "CRITICAL"
		return result, nil
	}

	result.IsVerified = true

	// Determine risk level based on reputation score and fraud reports
	if sender.ReputationScore >= 0.8 && sender.FraudReportCount < 5 {
		result.RiskLevel = "LOW"
		result.Explanation = fmt.Sprintf("Sender ID '%s' is verified and has a good reputation (%.2f)", senderHeader, sender.ReputationScore)
	} else if sender.ReputationScore >= 0.5 && sender.FraudReportCount < 20 {
		result.RiskLevel = "MEDIUM"
		result.Explanation = fmt.Sprintf("Sender ID '%s' is verified but has moderate reputation (%.2f) with %d fraud reports", 
			senderHeader, sender.ReputationScore, sender.FraudReportCount)
	} else {
		result.RiskLevel = "HIGH"
		result.Explanation = fmt.Sprintf("Sender ID '%s' is verified but has low reputation (%.2f) with %d fraud reports", 
			senderHeader, sender.ReputationScore, sender.FraudReportCount)
	}

	return result, nil
}

// UpdateSenderStats updates sender statistics after verification
func (s *HeaderVerificationService) UpdateSenderStats(ctx context.Context, senderHeader string, isFraud bool) error {
	err := s.rbiRepo.UpdateSenderStats(ctx, senderHeader, isFraud)
	if err != nil {
		utils.GetLogger().WithError(err).WithField("sender", senderHeader).Error("Failed to update sender stats")
		// Don't fail the operation, just log the error
	}
	return nil
}

