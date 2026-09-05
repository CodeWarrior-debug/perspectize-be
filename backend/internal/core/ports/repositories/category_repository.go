package repositories

import (
	"context"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
)

// CategoryRepository defines the contract for category persistence
type CategoryRepository interface {
	// Upsert inserts or updates a category by wikidata_qid
	Upsert(ctx context.Context, category *domain.Category) (*domain.Category, error)

	// GetByID fetches a category by its primary key
	GetByID(ctx context.Context, id int) (*domain.Category, error)
}
