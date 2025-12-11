package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/fraud-detection-system/backend/internal/models"
	"github.com/fraud-detection-system/backend/internal/repository"
	"github.com/fraud-detection-system/backend/internal/utils"
)

type ReportHandler struct {
	reportRepo *repository.ReportRepository
}

func NewReportHandler(reportRepo *repository.ReportRepository) *ReportHandler {
	return &ReportHandler{
		reportRepo: reportRepo,
	}
}

// CreateReport handles creating a fraud report
func (h *ReportHandler) CreateReport(c *gin.Context) {
	var req models.ReportInput
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

	// Determine priority based on report type
	priority := "MEDIUM"
	if req.ReportType == "FRAUD" {
		priority = "HIGH"
	}

	report := &models.Report{
		ID:             uuid.New(),
		UserID:         userID,
		MessageID:      req.MessageID,
		VerificationID: req.VerificationID,
		ReportType:     req.ReportType,
		Content:        req.Content,
		SenderHeader:   req.SenderHeader,
		Description:    req.Description,
		Status:         "PENDING",
		Priority:       priority,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.reportRepo.Create(c.Request.Context(), report); err != nil {
		utils.HandleError(c, err)
		return
	}

	response := &models.ReportResponse{
		ID:           report.ID,
		ReportType:   report.ReportType,
		Content:      report.Content,
		SenderHeader: report.SenderHeader,
		Description:  report.Description,
		Status:       report.Status,
		Priority:     report.Priority,
		CreatedAt:    report.CreatedAt,
	}

	utils.RespondWithSuccess(c, http.StatusCreated, response)
}

// GetReport handles getting a report by ID
func (h *ReportHandler) GetReport(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, utils.ErrBadRequest, "Invalid report ID")
		return
	}

	report, err := h.reportRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	response := &models.ReportResponse{
		ID:           report.ID,
		ReportType:   report.ReportType,
		Content:      report.Content,
		SenderHeader: report.SenderHeader,
		Description:  report.Description,
		Status:       report.Status,
		Priority:     report.Priority,
		ReviewedAt:   report.ReviewedAt,
		ReviewNotes:  report.ReviewNotes,
		CreatedAt:    report.CreatedAt,
	}

	utils.RespondWithSuccess(c, http.StatusOK, response)
}

// GetUserReports handles getting reports by user
func (h *ReportHandler) GetUserReports(c *gin.Context) {
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

	reports, err := h.reportRepo.GetByUserID(c.Request.Context(), userID.(uuid.UUID), limit, offset)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	responses := make([]*models.ReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = &models.ReportResponse{
			ID:           report.ID,
			ReportType:   report.ReportType,
			Content:      report.Content,
			SenderHeader: report.SenderHeader,
			Description:  report.Description,
			Status:       report.Status,
			Priority:     report.Priority,
			ReviewedAt:   report.ReviewedAt,
			ReviewNotes:  report.ReviewNotes,
			CreatedAt:    report.CreatedAt,
		}
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"reports": responses,
		"limit":   limit,
		"offset":  offset,
		"count":   len(responses),
	})
}

// GetReportStats handles getting report statistics
func (h *ReportHandler) GetReportStats(c *gin.Context) {
	stats, err := h.reportRepo.GetStats(c.Request.Context())
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, stats)
}

