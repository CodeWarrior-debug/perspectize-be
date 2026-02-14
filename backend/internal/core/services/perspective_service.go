package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/repositories"
	portservices "github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/services"
)

// PerspectiveService implements business logic for perspective operations
type PerspectiveService struct {
	repo     repositories.PerspectiveRepository
	userRepo repositories.UserRepository
}

// NewPerspectiveService creates a new perspective service
func NewPerspectiveService(repo repositories.PerspectiveRepository, userRepo repositories.UserRepository) *PerspectiveService {
	return &PerspectiveService{
		repo:     repo,
		userRepo: userRepo,
	}
}

// Create creates a new perspective with validation
func (s *PerspectiveService) Create(ctx context.Context, input portservices.CreatePerspectiveInput) (*domain.Perspective, error) {
	// Validate claim
	claim := strings.TrimSpace(input.Claim)
	if claim == "" {
		return nil, fmt.Errorf("%w: claim is required", domain.ErrInvalidInput)
	}
	if len(claim) > 255 {
		return nil, fmt.Errorf("%w: claim must be 255 characters or less", domain.ErrInvalidInput)
	}

	// Validate user exists
	if input.UserID <= 0 {
		return nil, fmt.Errorf("%w: user_id must be a positive integer", domain.ErrInvalidInput)
	}
	_, err := s.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: user with id %d not found", domain.ErrNotFound, input.UserID)
		}
		return nil, fmt.Errorf("failed to validate user: %w", err)
	}

	// Validate ratings are in range
	if !domain.ValidateRating(input.Quality) {
		return nil, fmt.Errorf("%w: quality %d", domain.ErrInvalidRating, *input.Quality)
	}
	if !domain.ValidateRating(input.Agreement) {
		return nil, fmt.Errorf("%w: agreement %d", domain.ErrInvalidRating, *input.Agreement)
	}
	if !domain.ValidateRating(input.Importance) {
		return nil, fmt.Errorf("%w: importance %d", domain.ErrInvalidRating, *input.Importance)
	}
	if !domain.ValidateRating(input.Confidence) {
		return nil, fmt.Errorf("%w: confidence %d", domain.ErrInvalidRating, *input.Confidence)
	}

	// Validate categorized ratings
	for _, cr := range input.CategorizedRatings {
		rating := cr.Rating
		if !domain.ValidateRating(&rating) {
			return nil, fmt.Errorf("%w: categorized rating for '%s' is %d", domain.ErrInvalidRating, cr.Category, rating)
		}
	}

	// Check for duplicate claim by same user
	existing, err := s.repo.GetByUserAndClaim(ctx, input.UserID, claim)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("%w: '%s'", domain.ErrDuplicateClaim, claim)
	}
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("failed to check existing claim: %w", err)
	}

	// Set default privacy
	privacy := domain.PrivacyPublic
	if input.Privacy != nil {
		privacy = *input.Privacy
	}

	perspective := &domain.Perspective{
		Claim:              claim,
		UserID:             input.UserID,
		ContentID:          input.ContentID,
		Quality:            input.Quality,
		Agreement:          input.Agreement,
		Importance:         input.Importance,
		Confidence:         input.Confidence,
		Like:               input.Like,
		Privacy:            privacy,
		Description:        input.Description,
		Category:           input.Category,
		Parts:              input.Parts,
		Labels:             input.Labels,
		CategorizedRatings: input.CategorizedRatings,
	}

	created, err := s.repo.Create(ctx, perspective)
	if err != nil {
		return nil, fmt.Errorf("failed to create perspective: %w", err)
	}

	return created, nil
}

// GetByID retrieves a perspective by ID
func (s *PerspectiveService) GetByID(ctx context.Context, id int) (*domain.Perspective, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: perspective id must be a positive integer", domain.ErrInvalidInput)
	}

	perspective, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get perspective: %w", err)
	}
	return perspective, nil
}

// Update updates an existing perspective
func (s *PerspectiveService) Update(ctx context.Context, input portservices.UpdatePerspectiveInput) (*domain.Perspective, error) {
	if input.ID <= 0 {
		return nil, fmt.Errorf("%w: perspective id must be a positive integer", domain.ErrInvalidInput)
	}

	// Get existing perspective
	existing, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get perspective: %w", err)
	}

	// Update claim if provided
	if input.Claim != nil {
		claim := strings.TrimSpace(*input.Claim)
		if claim == "" {
			return nil, fmt.Errorf("%w: claim cannot be empty", domain.ErrInvalidInput)
		}
		if len(claim) > 255 {
			return nil, fmt.Errorf("%w: claim must be 255 characters or less", domain.ErrInvalidInput)
		}
		// Check for duplicate if claim changed
		if claim != existing.Claim {
			dup, err := s.repo.GetByUserAndClaim(ctx, existing.UserID, claim)
			if err == nil && dup != nil {
				return nil, fmt.Errorf("%w: '%s'", domain.ErrDuplicateClaim, claim)
			}
			if err != nil && !errors.Is(err, domain.ErrNotFound) {
				return nil, fmt.Errorf("failed to check existing claim: %w", err)
			}
		}
		existing.Claim = claim
	}

	// Validate and update ratings
	if input.Quality != nil {
		if !domain.ValidateRating(input.Quality) {
			return nil, fmt.Errorf("%w: quality %d", domain.ErrInvalidRating, *input.Quality)
		}
		existing.Quality = input.Quality
	}
	if input.Agreement != nil {
		if !domain.ValidateRating(input.Agreement) {
			return nil, fmt.Errorf("%w: agreement %d", domain.ErrInvalidRating, *input.Agreement)
		}
		existing.Agreement = input.Agreement
	}
	if input.Importance != nil {
		if !domain.ValidateRating(input.Importance) {
			return nil, fmt.Errorf("%w: importance %d", domain.ErrInvalidRating, *input.Importance)
		}
		existing.Importance = input.Importance
	}
	if input.Confidence != nil {
		if !domain.ValidateRating(input.Confidence) {
			return nil, fmt.Errorf("%w: confidence %d", domain.ErrInvalidRating, *input.Confidence)
		}
		existing.Confidence = input.Confidence
	}

	// Validate categorized ratings if provided
	if input.CategorizedRatings != nil {
		for _, cr := range input.CategorizedRatings {
			rating := cr.Rating
			if !domain.ValidateRating(&rating) {
				return nil, fmt.Errorf("%w: categorized rating for '%s' is %d", domain.ErrInvalidRating, cr.Category, rating)
			}
		}
		existing.CategorizedRatings = input.CategorizedRatings
	}

	// Update optional fields
	if input.ContentID != nil {
		existing.ContentID = input.ContentID
	}
	if input.Like != nil {
		existing.Like = input.Like
	}
	if input.Privacy != nil {
		existing.Privacy = *input.Privacy
	}
	if input.Description != nil {
		existing.Description = input.Description
	}
	if input.Category != nil {
		existing.Category = input.Category
	}
	if input.ReviewStatus != nil {
		existing.ReviewStatus = input.ReviewStatus
	}
	if input.Parts != nil {
		existing.Parts = input.Parts
	}
	if input.Labels != nil {
		existing.Labels = input.Labels
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update perspective: %w", err)
	}

	return updated, nil
}

// Delete removes a perspective by ID
func (s *PerspectiveService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("%w: perspective id must be a positive integer", domain.ErrInvalidInput)
	}

	// Verify perspective exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get perspective: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete perspective: %w", err)
	}

	return nil
}

// ListPerspectives retrieves a paginated list of perspectives
func (s *PerspectiveService) ListPerspectives(ctx context.Context, params domain.PerspectiveListParams) (*domain.PaginatedPerspectives, error) {
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
		return nil, fmt.Errorf("failed to list perspectives: %w", err)
	}

	return result, nil
}
