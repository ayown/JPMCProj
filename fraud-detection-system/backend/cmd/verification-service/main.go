package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/fraud-detection-system/backend/internal/api/handlers"
	"github.com/fraud-detection-system/backend/internal/api/middleware"
	"github.com/fraud-detection-system/backend/internal/cache"
	"github.com/fraud-detection-system/backend/internal/config"
	"github.com/fraud-detection-system/backend/internal/database"
	"github.com/fraud-detection-system/backend/internal/repository"
	"github.com/fraud-detection-system/backend/internal/service"
	"github.com/fraud-detection-system/backend/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	utils.InitLogger(cfg.App.LogLevel)
	logger := utils.GetLogger()
	logger.Info("Starting Verification Service...")

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()
	logger.Info("Database connected")

	// Initialize cache
	redisCache, err := cache.NewRedisCache(cfg)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to Redis")
	}
	defer redisCache.Close()
	logger.Info("Redis connected")

	// Initialize repositories
	messageRepo := repository.NewMessageRepository(db.DB)
	verificationRepo := repository.NewVerificationRepository(db.DB)
	rbiRepo := repository.NewRBIRepository(db.DB)

	// Initialize services
	mlClient := service.NewMLClient(cfg)
	rbiService := service.NewRBIComplianceService(rbiRepo)
	headerService := service.NewHeaderVerificationService(rbiRepo)
	verificationService := service.NewVerificationService(
		messageRepo,
		verificationRepo,
		mlClient,
		rbiService,
		headerService,
		redisCache,
	)

	// Initialize handlers
	verificationHandler := handlers.NewVerificationHandler(verificationService)

	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware(cfg))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Verification routes
	v1 := router.Group("/api/v1/verify")
	v1.Use(middleware.OptionalAuthMiddleware(cfg))
	v1.Use(middleware.RateLimitMiddleware(redisCache, 100))
	{
		v1.POST("", verificationHandler.VerifyMessage)
		v1.GET("/:id", verificationHandler.GetVerification)
		v1.GET("/stats", verificationHandler.GetStats)
	}

	// Protected routes
	protected := router.Group("/api/v1/verify")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/history", verificationHandler.GetVerificationHistory)
	}

	// Create HTTP server
	port := os.Getenv("VERIFICATION_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start server
	go func() {
		logger.WithField("address", addr).Info("Verification Service started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Verification Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	logger.Info("Verification Service stopped")
}

