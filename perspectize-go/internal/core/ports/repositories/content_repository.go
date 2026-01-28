package repositories

import (
	"context"

	"github.com/yourorg/perspectize-go/internal/core/domain"
)

// ContentRepository defines the contract for content persistence
type ContentRepository interface {
	Create(ctx context.Context, content *domain.Content) (*domain.Content, error)
	GetByID(ctx context.Context, id int) (*domain.Content, error)
	GetByURL(ctx context.Context, url string) (*domain.Content, error)
}
