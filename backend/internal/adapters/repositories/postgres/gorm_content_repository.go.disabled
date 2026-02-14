package postgres

// GORM PROTOTYPE — Content repository using GORM with hex-clean separate models.
//
// Current sqlx version: 364 lines (content_repository.go)
// This GORM version:    ~155 lines (this file) + ~25 lines (shared in gorm_mappers.go)
//
// Key win: Dynamic ORDER BY with JSONB path expressions is just a string passed to .Order().
// No CASE WHEN workaround needed (unlike sqlc).

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/domain"
	paginator "github.com/pilagod/gorm-cursor-paginator/v2"
	"gorm.io/gorm"
)

// GORMContentRepository implements ContentRepository using GORM.
type GORMContentRepository struct {
	db *gorm.DB
}

// NewGORMContentRepository creates a new GORM-based content repository.
func NewGORMContentRepository(db *gorm.DB) *GORMContentRepository {
	return &GORMContentRepository{db: db}
}

// Create inserts a new content record.
func (r *GORMContentRepository) Create(ctx context.Context, content *domain.Content) (*domain.Content, error) {
	model := contentDomainToModel(content)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, fmt.Errorf("failed to insert content: %w", err)
	}

	content.ID = model.ID
	content.CreatedAt = model.CreatedAt
	content.UpdatedAt = model.UpdatedAt
	return content, nil
}

// GetByID retrieves content by ID.
func (r *GORMContentRepository) GetByID(ctx context.Context, id int) (*domain.Content, error) {
	var model ContentModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get content by id: %w", err)
	}
	return contentModelToDomain(&model), nil
}

// GetByURL retrieves content by URL.
func (r *GORMContentRepository) GetByURL(ctx context.Context, url string) (*domain.Content, error) {
	var model ContentModel
	if err := r.db.WithContext(ctx).Where("url = ?", url).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get content by url: %w", err)
	}
	return contentModelToDomain(&model), nil
}

// List retrieves a paginated, sorted, filtered list of content.
func (r *GORMContentRepository) List(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
	limit := 10
	if params.First != nil {
		limit = *params.First
	}

	// Build base query with dynamic filters
	query := r.db.WithContext(ctx).Model(&ContentModel{})

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

	// Dynamic ORDER BY — JSONB path expressions work naturally
	col := contentGORMSortColumn(params.SortBy)
	dir := "DESC"
	if params.SortOrder == domain.SortOrderAsc {
		dir = "ASC"
	}
	nullsClause := ""
	if isJSONBSort(params.SortBy) {
		nullsClause = " NULLS LAST"
	}
	query = query.Order(fmt.Sprintf("%s %s%s, id %s", col, dir, nullsClause, dir))

	// Cursor pagination
	p := paginator.New(&paginator.Config{
		Limit: limit,
		Order: paginator.Order(dir),
	})
	if params.After != nil {
		p.SetAfterCursor(*params.After)
	}

	var models []ContentModel
	cursor, err := p.Paginate(query, &models)
	if err != nil {
		return nil, fmt.Errorf("failed to list content: %w", err)
	}

	items := make([]*domain.Content, len(models))
	for i := range models {
		items[i] = contentModelToDomain(&models[i])
	}

	result := &domain.PaginatedContent{
		Items:   items,
		HasNext: cursor.After != nil,
		HasPrev: params.After != nil,
	}
	if cursor.After != nil {
		after := *cursor.After
		result.EndCursor = &after
	}
	if cursor.Before != nil {
		before := *cursor.Before
		result.StartCursor = &before
	}

	// Optional total count (respects filters, not cursor)
	if params.IncludeTotalCount {
		var count int64
		countQuery := r.db.WithContext(ctx).Model(&ContentModel{})
		if params.Filter != nil {
			if params.Filter.ContentType != nil {
				countQuery = countQuery.Where("content_type = ?", strings.ToLower(string(*params.Filter.ContentType)))
			}
			if params.Filter.MinLengthSeconds != nil {
				countQuery = countQuery.Where("length >= ?", *params.Filter.MinLengthSeconds)
			}
			if params.Filter.MaxLengthSeconds != nil {
				countQuery = countQuery.Where("length <= ?", *params.Filter.MaxLengthSeconds)
			}
			if params.Filter.Search != nil && *params.Filter.Search != "" {
				countQuery = countQuery.Where("name ILIKE ?", "%"+*params.Filter.Search+"%")
			}
		}
		if err := countQuery.Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count content: %w", err)
		}
		total := int(count)
		result.TotalCount = &total
	}

	return result, nil
}

// contentGORMSortColumn maps domain sort field to SQL column — same as current, just cleaner context.
func contentGORMSortColumn(sortBy domain.ContentSortBy) string {
	switch sortBy {
	case domain.ContentSortByUpdatedAt:
		return "updated_at"
	case domain.ContentSortByName:
		return "name"
	case domain.ContentSortByViewCount:
		return "(response->'items'->0->'statistics'->>'viewCount')::BIGINT"
	case domain.ContentSortByLikeCount:
		return "(response->'items'->0->'statistics'->>'likeCount')::BIGINT"
	case domain.ContentSortByPublishedAt:
		return "response->'items'->0->'snippet'->>'publishedAt'"
	default:
		return "created_at"
	}
}
