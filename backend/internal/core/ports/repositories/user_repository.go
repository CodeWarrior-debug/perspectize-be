package repositories

import (
	"context"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
)

// UserRepository defines the contract for user persistence
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByClerkID(ctx context.Context, clerkID string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	ListAll(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id int) error
	CreateFromClerk(ctx context.Context, clerkID string, username string, email string) (*domain.User, error)
	UpdateByClerkID(ctx context.Context, clerkID string, username string, email string) error
	DeactivateByClerkID(ctx context.Context, clerkID string) error
	// UpdateOnboarding replaces the onboarding JSON for the given user ID.
	UpdateOnboarding(ctx context.Context, userID int, onboarding domain.UserOnboarding) (*domain.User, error)
}
