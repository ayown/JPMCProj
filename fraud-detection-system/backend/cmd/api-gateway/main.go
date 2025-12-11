package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fraud-detection-system/backend/internal/api/handlers"
	"github.com/fraud-detection-system/backend/internal/api/routes"
	"github.com/fraud-detection-system/backend/internal/cache"
	"github.com/fraud-detection-system/backend/internal/config"
	"github.com/fraud-detection-system/backend/internal/database"
	"github.com/fraud-detection-system/backend/internal/queue"
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
	logger.Info("Starting API Gateway...")

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

	// Initialize Kafka producer
	producer := queue.NewProducer(cfg)
	defer producer.Close()
	logger.Info("Kafka producer initialized")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	messageRepo := repository.NewMessageRepository(db.DB)
	verificationRepo := repository.NewVerificationRepository(db.DB)
	reportRepo := repository.NewReportRepository(db.DB)
	rbiRepo := repository.NewRBIRepository(db.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)
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
	healthHandler := handlers.NewHealthHandler(db, redisCache, mlClient)
	authHandler := handlers.NewAuthHandler(authService)
	verificationHandler := handlers.NewVerificationHandler(verificationService)
	reportHandler := handlers.NewReportHandler(reportRepo)

	// Setup router
	router := routes.SetupRouter(&routes.RouterConfig{
		Config:              cfg,
		Cache:               redisCache,
		HealthHandler:       healthHandler,
		AuthHandler:         authHandler,
		VerificationHandler: verificationHandler,
		ReportHandler:       reportHandler,
	})

	// Create HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.WithField("address", addr).Info("API Gateway started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down API Gateway...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	logger.Info("API Gateway stopped")
}

