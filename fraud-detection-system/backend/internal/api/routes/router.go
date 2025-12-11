package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/fraud-detection-system/backend/internal/api/handlers"
	"github.com/fraud-detection-system/backend/internal/api/middleware"
	"github.com/fraud-detection-system/backend/internal/cache"
	"github.com/fraud-detection-system/backend/internal/config"
)

type RouterConfig struct {
	Config              *config.Config
	Cache               *cache.RedisCache
	HealthHandler       *handlers.HealthHandler
	AuthHandler         *handlers.AuthHandler
	VerificationHandler *handlers.VerificationHandler
	ReportHandler       *handlers.ReportHandler
}

// SetupRouter sets up the Gin router with all routes
func SetupRouter(cfg *RouterConfig) *gin.Engine {
	// Set Gin mode
	if cfg.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware(cfg.Config))

	// Health check endpoints (no auth required)
	router.GET("/health", cfg.HealthHandler.HealthCheck)
	router.GET("/ready", cfg.HealthHandler.ReadinessCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (no auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", cfg.AuthHandler.Register)
			auth.POST("/login", cfg.AuthHandler.Login)
			auth.POST("/refresh", cfg.AuthHandler.RefreshToken)
		}

		// Verification routes (optional auth)
		verify := v1.Group("/verify")
		verify.Use(middleware.OptionalAuthMiddleware(cfg.Config))
		verify.Use(middleware.RateLimitMiddleware(cfg.Cache, 100))
		{
			verify.POST("", cfg.VerificationHandler.VerifyMessage)
			verify.GET("/:id", cfg.VerificationHandler.GetVerification)
			verify.GET("/stats", cfg.VerificationHandler.GetStats)
		}

		// Protected routes (auth required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.Config))
		protected.Use(middleware.RateLimitMiddleware(cfg.Cache, 100))
		{
			// User profile
			protected.GET("/profile", cfg.AuthHandler.GetProfile)

			// Verification history
			protected.GET("/verify/history", cfg.VerificationHandler.GetVerificationHistory)

			// Reports
			reports := protected.Group("/reports")
			{
				reports.POST("", cfg.ReportHandler.CreateReport)
				reports.GET("/:id", cfg.ReportHandler.GetReport)
				reports.GET("", cfg.ReportHandler.GetUserReports)
				reports.GET("/stats", cfg.ReportHandler.GetReportStats)
			}
		}
	}

	return router
}

