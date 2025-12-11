package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/fraud-detection-system/backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, phone_number, is_active, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FullName, user.PhoneNumber,
		user.IsActive, user.IsVerified, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, phone_number, is_active, is_verified,
		       created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`
	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.PhoneNumber,
		&user.IsActive, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, phone_number, is_active, is_verified,
		       created_at, updated_at, last_login_at
		FROM users
		WHERE email = $1
	`
	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.PhoneNumber,
		&user.IsActive, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = $2, full_name = $3, phone_number = $4, is_active = $5,
		    is_verified = $6, updated_at = $7, last_login_at = $8
		WHERE id = $1
	`
	user.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.FullName, user.PhoneNumber,
		user.IsActive, user.IsVerified, user.UpdatedAt, user.LastLoginAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET last_login_at = $1 WHERE id = $2`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// Exists checks if a user exists by email
func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return exists, nil
}

