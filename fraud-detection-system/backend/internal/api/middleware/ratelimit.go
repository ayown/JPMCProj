package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/fraud-detection-system/backend/internal/cache"
	"github.com/fraud-detection-system/backend/internal/utils"
)

// RateLimitMiddleware implements rate limiting using Redis
func RateLimitMiddleware(cache *cache.RedisCache, requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP or user ID)
		identifier := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			identifier = fmt.Sprintf("user:%v", userID)
		}

		key := fmt.Sprintf("ratelimit:%s", identifier)
		ctx := c.Request.Context()

		// Increment counter
		count, err := cache.Increment(ctx, key)
		if err != nil {
			utils.GetLogger().WithError(err).Error("Failed to increment rate limit counter")
			// Don't block on error, just log it
			c.Next()
			return
		}

		// Set expiry on first request
		if count == 1 {
			if err := cache.SetExpire(ctx, key, time.Minute); err != nil {
				utils.GetLogger().WithError(err).Error("Failed to set rate limit expiry")
			}
		}

		// Check if limit exceeded
		if count > int64(requestsPerMinute) {
			utils.RespondWithError(c, http.StatusTooManyRequests, 
				fmt.Errorf("rate limit exceeded"), 
				fmt.Sprintf("Rate limit exceeded. Maximum %d requests per minute allowed", requestsPerMinute))
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", requestsPerMinute-int(count)))

		c.Next()
	}
}

