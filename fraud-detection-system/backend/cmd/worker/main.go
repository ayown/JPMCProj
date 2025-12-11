package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/fraud-detection-system/backend/internal/cache"
	"github.com/fraud-detection-system/backend/internal/config"
	"github.com/fraud-detection-system/backend/internal/database"
	"github.com/fraud-detection-system/backend/internal/models"
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
	logger.Info("Starting Worker Service...")

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

	// Create message handler
	messageHandler := func(ctx context.Context, msg *queue.QueueMessage) error {
		logger.WithField("message_id", msg.ID).Info("Processing verification request")

		// Extract verification request from payload
		content, ok := msg.Payload["content"].(string)
		if !ok {
			return fmt.Errorf("invalid content in payload")
		}

		senderHeader, ok := msg.Payload["sender_header"].(string)
		if !ok {
			return fmt.Errorf("invalid sender_header in payload")
		}

		// Create verification request
		req := &models.VerificationRequest{
			Content:      content,
			SenderHeader: senderHeader,
		}

		// Get user ID if present
		var userID *uuid.UUID
		if userIDStr, ok := msg.Payload["user_id"].(string); ok {
			if id, err := uuid.Parse(userIDStr); err == nil {
				userID = &id
			}
		}

		// Perform verification
		_, err := verificationService.VerifyMessage(ctx, req, userID)
		if err != nil {
			logger.WithError(err).Error("Failed to verify message")
			return err
		}

		logger.WithField("message_id", msg.ID).Info("Verification completed")
		return nil
	}

	// Initialize Kafka consumer
	consumer := queue.NewConsumer(cfg, cfg.Kafka.TopicVerification, messageHandler)
	defer consumer.Close()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consumer in a goroutine
	go func() {
		logger.Info("Worker started, waiting for messages...")
		if err := consumer.Start(ctx); err != nil {
			logger.WithError(err).Error("Consumer error")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Worker Service...")
	cancel()

	logger.Info("Worker Service stopped")
}

