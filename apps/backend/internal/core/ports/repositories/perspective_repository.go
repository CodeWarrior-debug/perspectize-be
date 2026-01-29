package repositories

import (
	"context"

	"github.com/CodeWarrior-debug/perspectize-be/apps/backend/internal/core/domain"
)

// PerspectiveRepository defines the contract for perspective persistence
type PerspectiveRepository interface {
	Create(ctx context.Context, perspective *domain.Perspective) (*domain.Perspective, error)
	GetByID(ctx context.Context, id int) (*domain.Perspective, error)
	GetByUserAndClaim(ctx context.Context, userID int, claim string) (*domain.Perspective, error)
	Update(ctx context.Context, perspective *domain.Perspective) (*domain.Perspective, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params domain.PerspectiveListParams) (*domain.PaginatedPerspectives, error)
}
