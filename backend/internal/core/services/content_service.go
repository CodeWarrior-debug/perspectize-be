package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/repositories"
	portservices "github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/services"
)

// ContentService implements business logic for content operations
type ContentService struct {
	repo          repositories.ContentRepository
	youtubeClient portservices.YouTubeClient
}

// NewContentService creates a new content service
func NewContentService(repo repositories.ContentRepository, yt portservices.YouTubeClient) *ContentService {
	return &ContentService{
		repo:          repo,
		youtubeClient: yt,
	}
}

// CreateFromYouTube creates content from a YouTube URL
func (s *ContentService) CreateFromYouTube(ctx context.Context, url string) (*domain.Content, error) {
	// Check if content already exists for this URL
	existing, err := s.repo.GetByURL(ctx, url)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("%w: content with URL %s already exists", domain.ErrAlreadyExists, url)
	}
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("failed to check existing content: %w", err)
	}

	// Extract video ID from URL
	videoID, err := s.youtubeClient.ExtractVideoID(url)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidURL, err)
	}

	// Fetch metadata from YouTube API
	metadata, err := s.youtubeClient.GetVideoMetadata(ctx, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch YouTube metadata: %w", err)
	}

	// Create domain content
	lengthUnits := "seconds"
	content := &domain.Content{
		Name:        metadata.Title,
		URL:         &url,
		ContentType: domain.ContentTypeYouTube,
		Length:      &metadata.Duration,
		LengthUnits: &lengthUnits,
		Response:    metadata.Response,
	}

	// Save to repository
	created, err := s.repo.Create(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to save content: %w", err)
	}

	return created, nil
}

// GetByID retrieves content by ID
func (s *ContentService) GetByID(ctx context.Context, id int) (*domain.Content, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: content id must be a positive integer", domain.ErrInvalidInput)
	}

	content, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}
	return content, nil
}

// ListContent retrieves a paginated list of content
func (s *ContentService) ListContent(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
	if params.First != nil {
		if *params.First < 1 || *params.First > 100 {
			return nil, fmt.Errorf("%w: first must be between 1 and 100", domain.ErrInvalidInput)
		}
	}
	if params.Last != nil {
		if *params.Last < 1 || *params.Last > 100 {
			return nil, fmt.Errorf("%w: last must be between 1 and 100", domain.ErrInvalidInput)
		}
	}

	result, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list content: %w", err)
	}

	return result, nil
}
