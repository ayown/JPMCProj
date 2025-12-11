package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fraud-detection-system/backend/internal/models"
	"github.com/fraud-detection-system/backend/internal/repository"
	"github.com/fraud-detection-system/backend/internal/utils"
)

type RBIComplianceService struct {
	rbiRepo *repository.RBIRepository
}

func NewRBIComplianceService(rbiRepo *repository.RBIRepository) *RBIComplianceService {
	return &RBIComplianceService{
		rbiRepo: rbiRepo,
	}
}

// VerifyCompliance checks if a message complies with RBI regulations
func (s *RBIComplianceService) VerifyCompliance(ctx context.Context, content string) (*models.RBIComplianceCheck, error) {
	result := &models.RBIComplianceCheck{
		IsCompliant:      true,
		MatchedCirculars: []string{},
		Keywords:         []string{},
		IsCurrentRequest: false,
	}

	// Extract potential keywords from content
	keywords := s.extractKeywords(content)
	if len(keywords) == 0 {
		result.Explanation = "No regulatory keywords found in message"
		return result, nil
	}

	result.Keywords = keywords

	// Search for matching circulars
	circulars, err := s.rbiRepo.SearchCircularsByKeywords(ctx, keywords)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to search RBI circulars")
		// Don't fail the verification, just log the error
		result.Explanation = "Unable to verify RBI compliance due to service error"
		return result, nil
	}

	// Check if any circulars match
	if len(circulars) > 0 {
		result.IsCurrentRequest = true
		for _, circular := range circulars {
			result.MatchedCirculars = append(result.MatchedCirculars, circular.CircularNumber)
			
			// Check if circular is currently valid
			now := time.Now()
			if circular.ExpiryDate != nil && circular.ExpiryDate.Before(now) {
				result.IsCompliant = false
				result.Explanation = fmt.Sprintf("Message references expired RBI circular %s (expired on %s)", 
					circular.CircularNumber, circular.ExpiryDate.Format("2006-01-02"))
				break
			}
		}

		if result.IsCompliant && result.IsCurrentRequest {
			result.Explanation = "Message references valid RBI circulars"
		}
	} else {
		// No matching circulars found, but message contains regulatory keywords
		// This could indicate a fake regulatory request
		if s.hasUrgentRegulatoryLanguage(content) {
			result.IsCompliant = false
			result.Explanation = "Message uses regulatory language but does not match any official RBI circulars"
		} else {
			result.Explanation = "Message contains regulatory keywords but no specific circular references"
		}
	}

	return result, nil
}

// extractKeywords extracts potential regulatory keywords from content
func (s *RBIComplianceService) extractKeywords(content string) []string {
	keywords := []string{}
	lowerContent := strings.ToLower(content)

	regulatoryTerms := []string{
		"kyc", "know your customer", "rbi", "reserve bank",
		"compliance", "verification", "mandatory", "regulation",
		"circular", "directive", "guideline", "policy",
		"pan card", "aadhaar", "identity verification",
	}

	for _, term := range regulatoryTerms {
		if strings.Contains(lowerContent, term) {
			keywords = append(keywords, term)
		}
	}

	return keywords
}

// hasUrgentRegulatoryLanguage checks if content has urgent regulatory language
func (s *RBIComplianceService) hasUrgentRegulatoryLanguage(content string) bool {
	urgentPhrases := []string{
		"mandatory kyc",
		"account will be blocked",
		"immediate verification required",
		"rbi directive",
		"comply within",
		"failure to comply",
		"account suspension",
	}

	lowerContent := strings.ToLower(content)
	for _, phrase := range urgentPhrases {
		if strings.Contains(lowerContent, phrase) {
			return true
		}
	}

	return false
}

// GetActiveCirculars retrieves all active RBI circulars
func (s *RBIComplianceService) GetActiveCirculars(ctx context.Context) ([]*models.RBICircular, error) {
	return s.rbiRepo.GetActiveCirculars(ctx)
}

