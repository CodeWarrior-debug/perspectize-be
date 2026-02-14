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

// GormContentRepository implements the ContentRepository interface using GORM
type GormContentRepository struct {
	db *gorm.DB
}

// Compile-time interface check
var _ repositories.ContentRepository = (*GormContentRepository)(nil)

// NewGormContentRepository creates a new GORM content repository
func NewGormContentRepository(db *gorm.DB) *GormContentRepository {
	return &GormContentRepository{db: db}
}

// Create inserts a new content record into the database
func (r *GormContentRepository) Create(ctx context.Context, content *domain.Content) (*domain.Content, error) {
	model := contentDomainToModel(content)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, fmt.Errorf("failed to insert content: %w", err)
	}

	return contentModelToDomain(model), nil
}

// GetByID retrieves a content record by its ID
func (r *GormContentRepository) GetByID(ctx context.Context, id int) (*domain.Content, error) {
	var model ContentModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get content by id: %w", err)
	}

	return contentModelToDomain(&model), nil
}

// GetByURL retrieves a content record by its URL
func (r *GormContentRepository) GetByURL(ctx context.Context, url string) (*domain.Content, error) {
	var model ContentModel
	err := r.db.WithContext(ctx).Where("url = ?", url).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get content by url: %w", err)
	}

	return contentModelToDomain(&model), nil
}

// List retrieves a paginated list of content using cursor-based pagination
func (r *GormContentRepository) List(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
	limit := 10
	if params.First != nil {
		limit = *params.First
	}

	col := sortColumnName(params.SortBy)
	dir := sortDirection(params.SortOrder)

	// Start query with context
	query := r.db.WithContext(ctx).Model(&ContentModel{})

	// Apply filters via GORM chaining
	if params.Filter != nil {
		if params.Filter.ContentType != nil {
			query = query.Where("content_type = ?", strings.ToLower(string(*params.Filter.ContentType)))
		}
		if params.Filter.MinLengthSeconds != nil {
			query = query.Where("length >= ?", *params.Filter.MinLengthSeconds)
		}
		if params.Filter.MaxLengthSeconds != nil {
			query = query.Where("length <= ?", *params.Filter.MaxLengthSeconds)
		}
		if params.Filter.Search != nil && *params.Filter.Search != "" {
			query = query.Where("name ILIKE ?", "%"+*params.Filter.Search+"%")
		}
	}

	// Total count (before cursor/limit — respects filters only)
	var totalCountInt *int
	if params.IncludeTotalCount {
		var count int64
		if err := query.Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count content: %w", err)
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

	// Dynamic ORDER BY with JSONB support
	nullsClause := ""
	if isJSONBSort(params.SortBy) {
		nullsClause = " NULLS LAST"
	}
	orderClause := col + " " + dir + nullsClause + ", id " + dir
	query = query.Order(gorm.Expr(orderClause))

	// Fetch limit+1 for hasNextPage detection
	query = query.Limit(limit + 1)
	var models []ContentModel
	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list content: %w", err)
	}

	// Build PaginatedContent result
	hasNext := len(models) > limit
	if hasNext {
		models = models[:limit]
	}

	hasPrev := params.After != nil

	items := make([]*domain.Content, len(models))
	for i := range models {
		items[i] = contentModelToDomain(&models[i])
	}

	result := &domain.PaginatedContent{
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
