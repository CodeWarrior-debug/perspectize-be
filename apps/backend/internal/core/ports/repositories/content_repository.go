package repositories

import (
	"context"

	"github.com/CodeWarrior-debug/perspectize-be/apps/backend/internal/core/domain"
)

// ContentRepository defines the contract for content persistence
type ContentRepository interface {
	Create(ctx context.Context, content *domain.Content) (*domain.Content, error)
	GetByID(ctx context.Context, id int) (*domain.Content, error)
	GetByURL(ctx context.Context, url string) (*domain.Content, error)
	List(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error)
}
