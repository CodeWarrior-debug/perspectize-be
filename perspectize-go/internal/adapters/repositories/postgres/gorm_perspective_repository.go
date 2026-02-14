package postgres

// GORM PROTOTYPE — Perspective repository using GORM with hex-clean separate models.
//
// Current sqlx version: 495 lines (perspective_repository.go)
// This GORM version:    ~195 lines (this file) + ~65 lines (shared in gorm_mappers.go)
//
// What changed:
// - No sql.NullString/NullInt64 — GORM uses Go pointers natively
// - No manual QueryRowContext/Scan — GORM .Create/.First/.Save/.Delete
// - No hand-built dynamic WHERE — GORM chaining with conditional .Where()
// - No hand-built pagination — gorm-cursor-paginator handles cursor encoding/decoding
// - Domain ↔ GORM model mapping still exists (hex-clean cost)

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/domain"
	paginator "github.com/pilagod/gorm-cursor-paginator/v2"
	"gorm.io/gorm"
)

// GORMPerspectiveRepository implements PerspectiveRepository using GORM.
type GORMPerspectiveRepository struct {
	db *gorm.DB
}

// NewGORMPerspectiveRepository creates a new GORM-based perspective repository.
func NewGORMPerspectiveRepository(db *gorm.DB) *GORMPerspectiveRepository {
	return &GORMPerspectiveRepository{db: db}
}

// Create inserts a new perspective.
func (r *GORMPerspectiveRepository) Create(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
	model := perspectiveDomainToModel(p)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, fmt.Errorf("failed to insert perspective: %w", err)
	}

	return r.GetByID(ctx, model.ID)
}

// GetByID retrieves a perspective by ID.
func (r *GORMPerspectiveRepository) GetByID(ctx context.Context, id int) (*domain.Perspective, error) {
	var model PerspectiveModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get perspective by id: %w", err)
	}

	return perspectiveModelToDomain(&model), nil
}

// GetByUserAndClaim retrieves a perspective by user ID and claim.
func (r *GORMPerspectiveRepository) GetByUserAndClaim(ctx context.Context, userID int, claim string) (*domain.Perspective, error) {
	var model PerspectiveModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND claim = ?", userID, claim).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get perspective by user and claim: %w", err)
	}

	return perspectiveModelToDomain(&model), nil
}

// Update updates an existing perspective.
func (r *GORMPerspectiveRepository) Update(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
	model := perspectiveDomainToModel(p)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return nil, fmt.Errorf("failed to update perspective: %w", err)
	}

	return r.GetByID(ctx, p.ID)
}

// Delete removes a perspective by ID.
func (r *GORMPerspectiveRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&PerspectiveModel{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete perspective: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// List retrieves a paginated list of perspectives with dynamic filters and sorting.
func (r *GORMPerspectiveRepository) List(ctx context.Context, params domain.PerspectiveListParams) (*domain.PaginatedPerspectives, error) {
	limit := 10
	if params.First != nil {
		limit = *params.First
	}

	// Build base query with dynamic filters
	query := r.db.WithContext(ctx).Model(&PerspectiveModel{})

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

	// Dynamic ORDER BY — just a string, no CASE WHEN gymnastics
	col := perspectiveGORMSortColumn(params.SortBy)
	dir := "DESC"
	if params.SortOrder == domain.SortOrderAsc {
		dir = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s, id %s", col, dir, dir))

	// Cursor pagination via gorm-cursor-paginator
	p := paginator.New(&paginator.Config{
		Limit: limit,
		Order: paginator.Order(dir),
	})
	if params.After != nil {
		p.SetAfterCursor(*params.After)
	}

	var models []PerspectiveModel
	cursor, err := p.Paginate(query, &models)
	if err != nil {
		return nil, fmt.Errorf("failed to list perspectives: %w", err)
	}

	// Map to domain
	items := make([]*domain.Perspective, len(models))
	for i := range models {
		items[i] = perspectiveModelToDomain(&models[i])
	}

	result := &domain.PaginatedPerspectives{
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

	// Optional total count
	if params.IncludeTotalCount {
		var count int64
		countQuery := r.db.WithContext(ctx).Model(&PerspectiveModel{})
		if params.Filter != nil {
			if params.Filter.UserID != nil {
				countQuery = countQuery.Where("user_id = ?", *params.Filter.UserID)
			}
			if params.Filter.ContentID != nil {
				countQuery = countQuery.Where("content_id = ?", *params.Filter.ContentID)
			}
			if params.Filter.Privacy != nil {
				countQuery = countQuery.Where("privacy = ?", strings.ToLower(string(*params.Filter.Privacy)))
			}
		}
		if err := countQuery.Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count perspectives: %w", err)
		}
		total := int(count)
		result.TotalCount = &total
	}

	return result, nil
}

func perspectiveGORMSortColumn(sortBy domain.PerspectiveSortBy) string {
	switch sortBy {
	case domain.PerspectiveSortByUpdatedAt:
		return "updated_at"
	case domain.PerspectiveSortByClaim:
		return "claim"
	default:
		return "created_at"
	}
}
