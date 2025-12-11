package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrNotFound         = errors.New("not found")
	ErrBadRequest       = errors.New("bad request")
	ErrInternalServer   = errors.New("internal server error")
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("expired token")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists       = errors.New("user already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrDatabaseError    = errors.New("database error")
	ErrCacheError       = errors.New("cache error")
	ErrMLServiceError   = errors.New("ML service error")
	ErrKafkaError       = errors.New("kafka error")
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// RespondWithError sends an error response
func RespondWithError(c *gin.Context, statusCode int, err error, message string) {
	c.JSON(statusCode, ErrorResponse{
		Error:   err.Error(),
		Message: message,
		Code:    statusCode,
	})
}

// RespondWithSuccess sends a success response
func RespondWithSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}

// HandleError handles different types of errors
func HandleError(c *gin.Context, err error) {
	switch err {
	case ErrUnauthorized, ErrInvalidToken, ErrExpiredToken:
		RespondWithError(c, http.StatusUnauthorized, err, "Authentication failed")
	case ErrForbidden:
		RespondWithError(c, http.StatusForbidden, err, "Access denied")
	case ErrNotFound, ErrUserNotFound:
		RespondWithError(c, http.StatusNotFound, err, "Resource not found")
	case ErrBadRequest, ErrInvalidCredentials:
		RespondWithError(c, http.StatusBadRequest, err, "Invalid request")
	case ErrUserExists:
		RespondWithError(c, http.StatusConflict, err, "Resource already exists")
	case ErrDatabaseError, ErrCacheError, ErrMLServiceError, ErrKafkaError:
		RespondWithError(c, http.StatusInternalServerError, err, "Service unavailable")
	default:
		RespondWithError(c, http.StatusInternalServerError, ErrInternalServer, "An error occurred")
	}
}

