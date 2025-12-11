package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/fraud-detection-system/backend/internal/models"
	"github.com/fraud-detection-system/backend/internal/service"
	"github.com/fraud-detection-system/backend/internal/utils"
)

type VerificationHandler struct {
	verificationService *service.VerificationService
}

func NewVerificationHandler(verificationService *service.VerificationService) *VerificationHandler {
	return &VerificationHandler{
		verificationService: verificationService,
	}
}

// VerifyMessage handles message verification requests
func (h *VerificationHandler) VerifyMessage(c *gin.Context) {
	var req models.VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, utils.ErrBadRequest, err.Error())
		return
	}

	// Get user ID if authenticated (optional)
	var userID *uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(uuid.UUID)
		userID = &id
	}

	result, err := h.verificationService.VerifyMessage(c.Request.Context(), &req, userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, result)
}

// GetVerification handles getting a verification by ID
func (h *VerificationHandler) GetVerification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, utils.ErrBadRequest, "Invalid verification ID")
		return
	}

	result, err := h.verificationService.GetVerificationByID(c.Request.Context(), id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, result)
}

// GetVerificationHistory handles getting verification history
func (h *VerificationHandler) GetVerificationHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, utils.ErrUnauthorized, "User not authenticated")
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	results, err := h.verificationService.GetVerificationHistory(c.Request.Context(), userID.(uuid.UUID), limit, offset)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"verifications": results,
		"limit":         limit,
		"offset":        offset,
		"count":         len(results),
	})
}

// GetStats handles getting verification statistics
func (h *VerificationHandler) GetStats(c *gin.Context) {
	var userID *uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(uuid.UUID)
		userID = &id
	}

	stats, err := h.verificationService.GetStats(c.Request.Context(), userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, stats)
}

