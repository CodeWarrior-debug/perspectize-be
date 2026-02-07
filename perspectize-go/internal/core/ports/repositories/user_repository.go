package repositories

import (
	"context"

	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/domain"
)

// UserRepository defines the contract for user persistence
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
