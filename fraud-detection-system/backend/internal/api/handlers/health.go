package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/fraud-detection-system/backend/internal/cache"
	"github.com/fraud-detection-system/backend/internal/database"
	"github.com/fraud-detection-system/backend/internal/service"
)

type HealthHandler struct {
	db       *database.Database
	cache    *cache.RedisCache
	mlClient *service.MLClient
}

func NewHealthHandler(db *database.Database, cache *cache.RedisCache, mlClient *service.MLClient) *HealthHandler {
	return &HealthHandler{
		db:       db,
		cache:    cache,
		mlClient: mlClient,
	}
}

// HealthCheck handles health check requests
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "fraud-detection-api",
	})
}

// ReadinessCheck handles readiness check requests
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	checks := map[string]string{
		"database": "healthy",
		"cache":    "healthy",
		"ml":       "healthy",
	}

	// Check database
	if err := h.db.HealthCheck(); err != nil {
		checks["database"] = "unhealthy"
	}

	// Check cache
	if err := h.cache.HealthCheck(); err != nil {
		checks["cache"] = "unhealthy"
	}

	// Check ML service
	if err := h.mlClient.HealthCheck(c.Request.Context()); err != nil {
		checks["ml"] = "unhealthy"
	}

	// Determine overall status
	status := "ready"
	statusCode := http.StatusOK
	for _, health := range checks {
		if health == "unhealthy" {
			status = "not ready"
			statusCode = http.StatusServiceUnavailable
			break
		}
	}

	c.JSON(statusCode, gin.H{
		"status": status,
		"checks": checks,
	})
}

