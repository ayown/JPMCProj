package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/fraud-detection-system/backend/internal/utils"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Log request
		duration := time.Since(startTime)
		logger := utils.GetLogger().WithFields(map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration.Milliseconds(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		if len(c.Errors) > 0 {
			logger.Error(c.Errors.String())
		} else {
			logger.Info("HTTP request processed")
		}
	}
}

