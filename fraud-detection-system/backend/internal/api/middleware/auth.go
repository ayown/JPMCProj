package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/fraud-detection-system/backend/internal/config"
	"github.com/fraud-detection-system/backend/internal/utils"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondWithError(c, http.StatusUnauthorized, utils.ErrUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondWithError(c, http.StatusUnauthorized, utils.ErrUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(token, cfg.JWT.Secret)
		if err != nil {
			utils.RespondWithError(c, http.StatusUnauthorized, utils.ErrInvalidToken, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens but doesn't require them
func OptionalAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				claims, err := utils.ValidateToken(token, cfg.JWT.Secret)
				if err == nil {
					c.Set("user_id", claims.UserID)
					c.Set("user_email", claims.Email)
				}
			}
		}
		c.Next()
	}
}

