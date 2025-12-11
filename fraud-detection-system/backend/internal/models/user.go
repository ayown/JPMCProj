package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	FullName     string     `json:"full_name" db:"full_name"`
	PhoneNumber  string     `json:"phone_number" db:"phone_number"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	IsVerified   bool       `json:"is_verified" db:"is_verified"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at" db:"last_login_at"`
}

type UserRegistration struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FullName    string `json:"full_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID          uuid.UUID  `json:"id"`
	Email       string     `json:"email"`
	FullName    string     `json:"full_name"`
	PhoneNumber string     `json:"phone_number"`
	IsActive    bool       `json:"is_active"`
	IsVerified  bool       `json:"is_verified"`
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

