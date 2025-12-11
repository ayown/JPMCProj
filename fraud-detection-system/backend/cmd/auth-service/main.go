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
	logger.Info("Starting Auth Service...")

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
	userRepo := repository.NewUserRepository(db.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware(cfg))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Auth routes
	v1 := router.Group("/api/v1/auth")
	v1.Use(middleware.RateLimitMiddleware(redisCache, 50))
	{
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)
		v1.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/profile", authHandler.GetProfile)
	}

	// Create HTTP server
	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8081"
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
		logger.WithField("address", addr).Info("Auth Service started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Auth Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	logger.Info("Auth Service stopped")
}

