package services

import (
	"context"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
)

// ContentService defines the contract for content business logic
type ContentService interface {
	// CreateFromYouTube creates content from a YouTube URL
	CreateFromYouTube(ctx context.Context, url string) (*domain.Content, error)

	// GetByID retrieves content by ID
	GetByID(ctx context.Context, id int) (*domain.Content, error)

	// ListContent retrieves a paginated list of content
	ListContent(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error)
}
