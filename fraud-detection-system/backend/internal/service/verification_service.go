package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/fraud-detection-system/backend/internal/cache"
	"github.com/fraud-detection-system/backend/internal/models"
	"github.com/fraud-detection-system/backend/internal/repository"
	"github.com/fraud-detection-system/backend/internal/utils"
)

type VerificationService struct {
	messageRepo       *repository.MessageRepository
	verificationRepo  *repository.VerificationRepository
	mlClient          *MLClient
	rbiService        *RBIComplianceService
	headerService     *HeaderVerificationService
	cache             *cache.RedisCache
}

func NewVerificationService(
	messageRepo *repository.MessageRepository,
	verificationRepo *repository.VerificationRepository,
	mlClient *MLClient,
	rbiService *RBIComplianceService,
	headerService *HeaderVerificationService,
	cache *cache.RedisCache,
) *VerificationService {
	return &VerificationService{
		messageRepo:      messageRepo,
		verificationRepo: verificationRepo,
		mlClient:         mlClient,
		rbiService:       rbiService,
		headerService:    headerService,
		cache:            cache,
	}
}

// VerifyMessage performs comprehensive fraud verification on a message
func (s *VerificationService) VerifyMessage(ctx context.Context, req *models.VerificationRequest, userID *uuid.UUID) (*models.VerificationResponse, error) {
	startTime := time.Now()

	// Extract message features
	features := s.extractFeatures(req.Content, req.SenderHeader)

	// Create message record
	message := &models.Message{
		ID:            uuid.New(),
		UserID:        userID,
		Content:       req.Content,
		SenderHeader:  req.SenderHeader,
		ReceivedAt:    req.ReceivedAt,
		MessageType:   req.MessageType,
		PhoneNumber:   req.PhoneNumber,
		HasLinks:      features.HasLinks,
		LinkCount:     features.LinkCount,
		ExtractedURLs: features.ExtractedURLs,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to create message")
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Perform ML inference
	mlReq := &models.MLInferenceRequest{
		Content:      req.Content,
		SenderHeader: req.SenderHeader,
		Features:     features,
	}

	mlResp, err := s.mlClient.Predict(ctx, mlReq)
	if err != nil {
		utils.GetLogger().WithError(err).Error("ML prediction failed")
		return nil, fmt.Errorf("ML prediction failed: %w", err)
	}

	// Verify header
	headerResult, err := s.headerService.VerifyHeader(ctx, req.SenderHeader)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Header verification failed")
		// Continue with default values
		headerResult = &models.HeaderVerificationResult{
			IsVerified:      false,
			ReputationScore: 0.0,
			RiskLevel:       "HIGH",
		}
	}

	// Verify RBI compliance
	rbiResult, err := s.rbiService.VerifyCompliance(ctx, req.Content)
	if err != nil {
		utils.GetLogger().WithError(err).Error("RBI compliance check failed")
		// Continue with default values
		rbiResult = &models.RBIComplianceCheck{
			IsCompliant: true,
			Explanation: "Unable to verify RBI compliance",
		}
	}

	// Combine results to determine final fraud score
	finalScore, isFraud := s.calculateFinalScore(mlResp, headerResult, rbiResult)

	// Determine fraud type
	fraudType := ""
	if isFraud && mlResp.FraudType != "" {
		fraudType = mlResp.FraudType
	}

	// Generate explanation
	explanation := s.generateExplanation(mlResp, headerResult, rbiResult, isFraud)

	// Generate recommendations
	recommendations := s.generateRecommendations(isFraud, headerResult, rbiResult)

	// Marshal JSON fields
	mlPredictionsJSON, _ := json.Marshal(mlResp.ModelPredictions)
	rbiResultJSON, _ := json.Marshal(rbiResult)
	recommendationsJSON, _ := json.Marshal(recommendations)

	// Create verification record
	verification := &models.Verification{
		ID:                    uuid.New(),
		MessageID:             message.ID,
		UserID:                userID,
		IsFraud:               isFraud,
		FraudScore:            finalScore,
		FraudType:             &fraudType,
		Confidence:            mlResp.Confidence,
		ModelVersion:          mlResp.ModelVersion,
		MLPredictions:         string(mlPredictionsJSON),
		HeaderVerified:        headerResult.IsVerified,
		HeaderScore:           headerResult.ReputationScore,
		RBICompliant:          rbiResult.IsCompliant,
		RBIVerificationResult: string(rbiResultJSON),
		Explanation:           explanation,
		Recommendations:       string(recommendationsJSON),
		ProcessingTimeMs:      int(time.Since(startTime).Milliseconds()),
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	if err := s.verificationRepo.Create(ctx, verification); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to create verification")
		return nil, fmt.Errorf("failed to create verification: %w", err)
	}

	// Update sender stats
	go s.headerService.UpdateSenderStats(context.Background(), req.SenderHeader, isFraud)

	// Determine risk level
	riskLevel := s.determineRiskLevel(finalScore, headerResult.RiskLevel)

	return &models.VerificationResponse{
		ID:               verification.ID,
		MessageID:        message.ID,
		IsFraud:          isFraud,
		FraudScore:       finalScore,
		FraudType:        &fraudType,
		Confidence:       mlResp.Confidence,
		RiskLevel:        riskLevel,
		HeaderVerified:   headerResult.IsVerified,
		RBICompliant:     rbiResult.IsCompliant,
		Explanation:      explanation,
		Recommendations:  recommendations,
		ModelPredictions: mlResp.ModelPredictions,
		ProcessingTimeMs: verification.ProcessingTimeMs,
		VerifiedAt:       verification.CreatedAt,
	}, nil
}

// extractFeatures extracts features from message content
func (s *VerificationService) extractFeatures(content, senderHeader string) models.MessageFeatures {
	urls := utils.ExtractURLs(content)
	phoneNumbers := utils.ExtractPhoneNumbers(content)

	return models.MessageFeatures{
		Content:          content,
		SenderHeader:     senderHeader,
		MessageLength:    len(content),
		HasLinks:         len(urls) > 0,
		LinkCount:        len(urls),
		ExtractedURLs:    urls,
		HasPhoneNumber:   len(phoneNumbers) > 0,
		PhoneNumberCount: len(phoneNumbers),
		HasUrgentWords:   utils.HasUrgentWords(content),
		UrgentWordCount:  utils.CountUrgentWords(content),
		SpecialCharRatio: utils.CalculateSpecialCharRatio(content),
		CapitalRatio:     utils.CalculateCapitalRatio(content),
		NumberRatio:      utils.CalculateNumberRatio(content),
		HasKYCKeywords:   utils.HasKYCKeywords(content),
		HasBankNames:     utils.HasBankNames(content),
	}
}

// calculateFinalScore combines all verification results
func (s *VerificationService) calculateFinalScore(mlResp *models.MLInferenceResponse, headerResult *models.HeaderVerificationResult, rbiResult *models.RBIComplianceCheck) (float64, bool) {
	// Start with ML score (weight: 0.6)
	finalScore := mlResp.FraudScore * 0.6

	// Add header verification score (weight: 0.25)
	if !headerResult.IsVerified {
		finalScore += 0.25
	} else {
		// Inverse of reputation score
		finalScore += (1.0 - headerResult.ReputationScore) * 0.25
	}

	// Add RBI compliance score (weight: 0.15)
	if !rbiResult.IsCompliant {
		finalScore += 0.15
	}

	// Determine if it's fraud (threshold: 0.5)
	isFraud := finalScore >= 0.5

	return finalScore, isFraud
}

// generateExplanation generates a human-readable explanation
func (s *VerificationService) generateExplanation(mlResp *models.MLInferenceResponse, headerResult *models.HeaderVerificationResult, rbiResult *models.RBIComplianceCheck, isFraud bool) string {
	var parts []string

	if isFraud {
		parts = append(parts, "⚠️ This message has been flagged as potentially fraudulent.")
	} else {
		parts = append(parts, "✓ This message appears to be legitimate.")
	}

	// ML explanation
	if mlResp.Explanation != "" {
		parts = append(parts, fmt.Sprintf("ML Analysis: %s", mlResp.Explanation))
	}

	// Header explanation
	parts = append(parts, fmt.Sprintf("Sender Verification: %s", headerResult.Explanation))

	// RBI explanation
	if rbiResult.Explanation != "" {
		parts = append(parts, fmt.Sprintf("RBI Compliance: %s", rbiResult.Explanation))
	}

	return strings.Join(parts, " ")
}

// generateRecommendations generates actionable recommendations
func (s *VerificationService) generateRecommendations(isFraud bool, headerResult *models.HeaderVerificationResult, rbiResult *models.RBIComplianceCheck) []string {
	recommendations := []string{}

	if isFraud {
		recommendations = append(recommendations, "Do not click on any links in this message")
		recommendations = append(recommendations, "Do not share personal or financial information")
		recommendations = append(recommendations, "Do not call any phone numbers provided in the message")
		
		if !headerResult.IsVerified {
			recommendations = append(recommendations, "The sender is not verified - this is likely a spoofed message")
		}
		
		if !rbiResult.IsCompliant {
			recommendations = append(recommendations, "This message references fake or expired regulatory requirements")
		}
		
		recommendations = append(recommendations, "Report this message to your bank and regulatory authorities")
		recommendations = append(recommendations, "Delete this message immediately")
	} else {
		recommendations = append(recommendations, "This message appears legitimate, but always verify through official channels")
		recommendations = append(recommendations, "Contact your bank directly using their official contact information")
		recommendations = append(recommendations, "Never share OTPs, passwords, or CVV numbers")
	}

	return recommendations
}

// determineRiskLevel determines the overall risk level
func (s *VerificationService) determineRiskLevel(fraudScore float64, headerRisk string) string {
	if fraudScore >= 0.8 || headerRisk == "CRITICAL" {
		return "CRITICAL"
	} else if fraudScore >= 0.6 || headerRisk == "HIGH" {
		return "HIGH"
	} else if fraudScore >= 0.4 || headerRisk == "MEDIUM" {
		return "MEDIUM"
	}
	return "LOW"
}

// GetVerificationByID retrieves a verification by ID
func (s *VerificationService) GetVerificationByID(ctx context.Context, id uuid.UUID) (*models.VerificationResponse, error) {
	verification, err := s.verificationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, utils.ErrNotFound
	}

	// Parse JSON fields
	var modelPredictions map[string]interface{}
	json.Unmarshal([]byte(verification.MLPredictions), &modelPredictions)

	var recommendations []string
	json.Unmarshal([]byte(verification.Recommendations), &recommendations)

	riskLevel := s.determineRiskLevel(verification.FraudScore, "")

	return &models.VerificationResponse{
		ID:               verification.ID,
		MessageID:        verification.MessageID,
		IsFraud:          verification.IsFraud,
		FraudScore:       verification.FraudScore,
		FraudType:        verification.FraudType,
		Confidence:       verification.Confidence,
		RiskLevel:        riskLevel,
		HeaderVerified:   verification.HeaderVerified,
		RBICompliant:     verification.RBICompliant,
		Explanation:      verification.Explanation,
		Recommendations:  recommendations,
		ModelPredictions: modelPredictions,
		ProcessingTimeMs: verification.ProcessingTimeMs,
		VerifiedAt:       verification.CreatedAt,
	}, nil
}

// GetVerificationHistory retrieves verification history for a user
func (s *VerificationService) GetVerificationHistory(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.VerificationResponse, error) {
	verifications, err := s.verificationRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.VerificationResponse, len(verifications))
	for i, v := range verifications {
		var modelPredictions map[string]interface{}
		json.Unmarshal([]byte(v.MLPredictions), &modelPredictions)

		var recommendations []string
		json.Unmarshal([]byte(v.Recommendations), &recommendations)

		riskLevel := s.determineRiskLevel(v.FraudScore, "")

		responses[i] = &models.VerificationResponse{
			ID:               v.ID,
			MessageID:        v.MessageID,
			IsFraud:          v.IsFraud,
			FraudScore:       v.FraudScore,
			FraudType:        v.FraudType,
			Confidence:       v.Confidence,
			RiskLevel:        riskLevel,
			HeaderVerified:   v.HeaderVerified,
			RBICompliant:     v.RBICompliant,
			Explanation:      v.Explanation,
			Recommendations:  recommendations,
			ModelPredictions: modelPredictions,
			ProcessingTimeMs: v.ProcessingTimeMs,
			VerifiedAt:       v.CreatedAt,
		}
	}

	return responses, nil
}

// GetStats retrieves verification statistics
func (s *VerificationService) GetStats(ctx context.Context, userID *uuid.UUID) (*models.VerificationStats, error) {
	return s.verificationRepo.GetStats(ctx, userID)
}

