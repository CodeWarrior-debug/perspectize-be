package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/repositories"
	"gorm.io/gorm"
)

// GormPerspectiveRepository implements the PerspectiveRepository interface using GORM
type GormPerspectiveRepository struct {
	db *gorm.DB
}

// Compile-time interface check
var _ repositories.PerspectiveRepository = (*GormPerspectiveRepository)(nil)

// NewGormPerspectiveRepository creates a new GORM perspective repository
func NewGormPerspectiveRepository(db *gorm.DB) *GormPerspectiveRepository {
	return &GormPerspectiveRepository{db: db}
}

// Create inserts a new perspective record into the database
func (r *GormPerspectiveRepository) Create(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
	model := perspectiveDomainToModel(p)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, fmt.Errorf("failed to insert perspective: %w", err)
	}

	// Fetch fresh record with DB-generated timestamps
	return r.GetByID(ctx, model.ID)
}

// GetByID retrieves a perspective by its ID
func (r *GormPerspectiveRepository) GetByID(ctx context.Context, id int) (*domain.Perspective, error) {
	var model PerspectiveModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get perspective by id: %w", err)
	}

	return perspectiveModelToDomain(&model), nil
}

// GetByUserAndClaim retrieves a perspective by user ID and claim text
func (r *GormPerspectiveRepository) GetByUserAndClaim(ctx context.Context, userID int, claim string) (*domain.Perspective, error) {
	var model PerspectiveModel
	err := r.db.WithContext(ctx).Where("user_id = ? AND claim = ?", userID, claim).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get perspective by user and claim: %w", err)
	}

	return perspectiveModelToDomain(&model), nil
}

// Update updates an existing perspective
func (r *GormPerspectiveRepository) Update(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
	model := perspectiveDomainToModel(p)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return nil, fmt.Errorf("failed to update perspective: %w", err)
	}

	// Fetch fresh record with updated timestamps
	return r.GetByID(ctx, model.ID)
}

// Delete removes a perspective by ID
func (r *GormPerspectiveRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&PerspectiveModel{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete perspective: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// List retrieves a paginated list of perspectives
func (r *GormPerspectiveRepository) List(ctx context.Context, params domain.PerspectiveListParams) (*domain.PaginatedPerspectives, error) {
	limit := 10
	if params.First != nil {
		limit = *params.First
	}

	col := perspectiveSortColumnName(params.SortBy)
	dir := sortDirection(params.SortOrder)

	// Start query with context
	query := r.db.WithContext(ctx).Model(&PerspectiveModel{})

	// Apply filters via GORM chaining
	if params.Filter != nil {
		if params.Filter.UserID != nil {
			query = query.Where("user_id = ?", *params.Filter.UserID)
		}
		if params.Filter.ContentID != nil {
			query = query.Where("content_id = ?", *params.Filter.ContentID)
		}
		if params.Filter.Privacy != nil {
			query = query.Where("privacy = ?", strings.ToLower(string(*params.Filter.Privacy)))
		}
	}

	// Total count (before cursor/limit — respects filters only)
	var totalCountInt *int
	if params.IncludeTotalCount {
		var count int64
		if err := query.Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count perspectives: %w", err)
		}
		countInt := int(count)
		totalCountInt = &countInt
	}

	// Cursor pagination
	if params.After != nil {
		cursorID, err := decodeCursor(*params.After)
		if err != nil {
			return nil, fmt.Errorf("invalid after cursor: %w", err)
		}
		if dir == "DESC" {
			query = query.Where("id < ?", cursorID)
		} else {
			query = query.Where("id > ?", cursorID)
		}
	}

	// Dynamic ORDER BY (no JSONB sorts for perspectives)
	orderClause := col + " " + dir + ", id " + dir
	query = query.Order(orderClause)

	// Fetch limit+1 for hasNextPage detection
	query = query.Limit(limit + 1)
	var models []PerspectiveModel
	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list perspectives: %w", err)
	}

	// Build PaginatedPerspectives result
	hasNext := len(models) > limit
	if hasNext {
		models = models[:limit]
	}

	hasPrev := params.After != nil

	items := make([]*domain.Perspective, len(models))
	for i := range models {
		items[i] = perspectiveModelToDomain(&models[i])
	}

	result := &domain.PaginatedPerspectives{
		Items:      items,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
		TotalCount: totalCountInt,
	}

	if len(items) > 0 {
		start := encodeCursor(items[0].ID)
		end := encodeCursor(items[len(items)-1].ID)
		result.StartCursor = &start
		result.EndCursor = &end
	}

	return result, nil
}
