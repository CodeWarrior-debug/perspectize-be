package services

import (
	"context"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
)

// UserService defines the contract for user business logic
type UserService interface {
	// Create creates a new user with validation
	Create(ctx context.Context, username, email string) (*domain.User, error)

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int) (*domain.User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*domain.User, error)

	// ListAll retrieves all users
	ListAll(ctx context.Context) ([]*domain.User, error)
}
