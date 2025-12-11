package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/fraud-detection-system/backend/internal/config"
	"github.com/fraud-detection-system/backend/internal/models"
	"github.com/fraud-detection-system/backend/internal/repository"
	"github.com/fraud-detection-system/backend/internal/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *models.UserRegistration) (*models.UserResponse, error) {
	// Check if user already exists
	exists, err := s.userRepo.Exists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, utils.ErrUserExists
	}

	// Validate password
	if !utils.ValidatePassword(req.Password) {
		return nil, fmt.Errorf("password must be at least 8 characters and contain uppercase, lowercase, number, and special character")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		PhoneNumber:  req.PhoneNumber,
		IsActive:     true,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		IsActive:    user.IsActive,
		IsVerified:  user.IsVerified,
		CreatedAt:   user.CreatedAt,
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *models.UserLogin) (*models.TokenPair, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Email, s.config.JWT.Secret, s.config.JWT.Expiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, s.config.JWT.Secret, s.config.JWT.RefreshTokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to update last login")
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWT.Expiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// RefreshToken refreshes an access token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(refreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, utils.ErrInvalidToken
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, utils.ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Generate new tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Email, s.config.JWT.Secret, s.config.JWT.Expiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, s.config.JWT.Secret, s.config.JWT.RefreshTokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.config.JWT.Expiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, utils.ErrUserNotFound
	}

	return &models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		IsActive:    user.IsActive,
		IsVerified:  user.IsVerified,
		CreatedAt:   user.CreatedAt,
		LastLoginAt: user.LastLoginAt,
	}, nil
}

