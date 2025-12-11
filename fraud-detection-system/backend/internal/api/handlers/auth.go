package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/fraud-detection-system/backend/internal/models"
	"github.com/fraud-detection-system/backend/internal/service"
	"github.com/fraud-detection-system/backend/internal/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserRegistration
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, utils.ErrBadRequest, err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, user)
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, utils.ErrBadRequest, err.Error())
		return
	}

	tokens, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, tokens)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, utils.ErrBadRequest, err.Error())
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, tokens)
}

// GetProfile handles getting user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, utils.ErrUnauthorized, "User not authenticated")
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, user)
}

