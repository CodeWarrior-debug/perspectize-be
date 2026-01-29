package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/yourorg/perspectize-go/internal/core/domain"
)

// PerspectiveRepository implements the PerspectiveRepository interface using PostgreSQL
type PerspectiveRepository struct {
	db *sqlx.DB
}

// perspectiveRow represents the database row structure for perspectives
type perspectiveRow struct {
	ID                 int            `db:"id"`
	Claim              string         `db:"claim"`
	UserID             int            `db:"user_id"`
	ContentID          sql.NullInt64  `db:"content_id"`
	Like               sql.NullString `db:"like"`
	Quality            sql.NullInt64  `db:"quality"`
	Agreement          sql.NullInt64  `db:"agreement"`
	Importance         sql.NullInt64  `db:"importance"`
	Confidence         sql.NullInt64  `db:"confidence"`
	Privacy            sql.NullString `db:"privacy"`
	Parts              pq.Int64Array  `db:"parts"`
	Category           sql.NullString `db:"category"`
	Labels             pq.StringArray `db:"labels"`
	Description        sql.NullString `db:"description"`
	ReviewStatus       sql.NullString `db:"review_status"`
	CategorizedRatings [][]byte       `db:"categorized_ratings"`
	CreatedAt          sql.NullTime   `db:"created_at"`
	UpdatedAt          sql.NullTime   `db:"updated_at"`
}

// NewPerspectiveRepository creates a new PostgreSQL perspective repository
func NewPerspectiveRepository(db *sqlx.DB) *PerspectiveRepository {
	return &PerspectiveRepository{db: db}
}

// Create inserts a new perspective record into the database
func (r *PerspectiveRepository) Create(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
	// Marshal categorized ratings
	var categorizedRatings [][]byte
	if len(p.CategorizedRatings) > 0 {
		categorizedRatings = make([][]byte, len(p.CategorizedRatings))
		for i, cr := range p.CategorizedRatings {
			data, err := json.Marshal(cr)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal categorized rating: %w", err)
			}
			categorizedRatings[i] = data
		}
	}

	query := `
		INSERT INTO perspectives (
			claim, user_id, content_id, "like", quality, agreement, importance, confidence,
			privacy, parts, category, labels, description, review_status, categorized_ratings
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		) RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		p.Claim,
		p.UserID,
		toNullInt64FromIntPtr(p.ContentID),
		toNullString(p.Like),
		toNullInt64FromIntPtr(p.Quality),
		toNullInt64FromIntPtr(p.Agreement),
		toNullInt64FromIntPtr(p.Importance),
		toNullInt64FromIntPtr(p.Confidence),
		string(p.Privacy),
		pq.Array(p.Parts),
		toNullString(p.Category),
		pq.Array(p.Labels),
		toNullString(p.Description),
		toNullStringFromReviewStatus(p.ReviewStatus),
		pq.Array(categorizedRatings),
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to insert perspective: %w", err)
	}

	return r.GetByID(ctx, id)
}

// GetByID retrieves a perspective by its ID
func (r *PerspectiveRepository) GetByID(ctx context.Context, id int) (*domain.Perspective, error) {
	query := `SELECT id, claim, user_id, content_id, "like", quality, agreement, importance,
		confidence, privacy, parts, category, labels, description, review_status,
		categorized_ratings, created_at, updated_at
		FROM perspectives WHERE id = $1`

	var row perspectiveRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get perspective by id: %w", err)
	}

	return perspectiveRowToDomain(&row), nil
}

// GetByUserAndClaim retrieves a perspective by user ID and claim text
func (r *PerspectiveRepository) GetByUserAndClaim(ctx context.Context, userID int, claim string) (*domain.Perspective, error) {
	query := `SELECT id, claim, user_id, content_id, "like", quality, agreement, importance,
		confidence, privacy, parts, category, labels, description, review_status,
		categorized_ratings, created_at, updated_at
		FROM perspectives WHERE user_id = $1 AND claim = $2`

	var row perspectiveRow
	err := r.db.GetContext(ctx, &row, query, userID, claim)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get perspective by user and claim: %w", err)
	}

	return perspectiveRowToDomain(&row), nil
}

// Update updates an existing perspective
func (r *PerspectiveRepository) Update(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
	// Marshal categorized ratings
	var categorizedRatings [][]byte
	if len(p.CategorizedRatings) > 0 {
		categorizedRatings = make([][]byte, len(p.CategorizedRatings))
		for i, cr := range p.CategorizedRatings {
			data, err := json.Marshal(cr)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal categorized rating: %w", err)
			}
			categorizedRatings[i] = data
		}
	}

	query := `
		UPDATE perspectives SET
			claim = $1, content_id = $2, "like" = $3, quality = $4, agreement = $5,
			importance = $6, confidence = $7, privacy = $8, parts = $9, category = $10,
			labels = $11, description = $12, review_status = $13, categorized_ratings = $14
		WHERE id = $15`

	_, err := r.db.ExecContext(ctx, query,
		p.Claim,
		toNullInt64FromIntPtr(p.ContentID),
		toNullString(p.Like),
		toNullInt64FromIntPtr(p.Quality),
		toNullInt64FromIntPtr(p.Agreement),
		toNullInt64FromIntPtr(p.Importance),
		toNullInt64FromIntPtr(p.Confidence),
		string(p.Privacy),
		pq.Array(p.Parts),
		toNullString(p.Category),
		pq.Array(p.Labels),
		toNullString(p.Description),
		toNullStringFromReviewStatus(p.ReviewStatus),
		pq.Array(categorizedRatings),
		p.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update perspective: %w", err)
	}

	return r.GetByID(ctx, p.ID)
}

// Delete removes a perspective by ID
func (r *PerspectiveRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM perspectives WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete perspective: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// List retrieves a paginated list of perspectives
func (r *PerspectiveRepository) List(ctx context.Context, params domain.PerspectiveListParams) (*domain.PaginatedPerspectives, error) {
	limit := 10
	if params.First != nil {
		limit = *params.First
	}

	col := perspectiveSortColumnName(params.SortBy)
	dir := sortDirection(params.SortOrder)

	var conditions []string
	var args []interface{}
	argIdx := 1

	if params.After != nil {
		cursorID, err := decodeCursor(*params.After)
		if err != nil {
			return nil, fmt.Errorf("invalid after cursor: %w", err)
		}
		if dir == "DESC" {
			conditions = append(conditions, fmt.Sprintf("id < $%d", argIdx))
		} else {
			conditions = append(conditions, fmt.Sprintf("id > $%d", argIdx))
		}
		args = append(args, cursorID)
		argIdx++
	}

	// Apply filters
	if params.Filter != nil {
		if params.Filter.UserID != nil {
			conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIdx))
			args = append(args, *params.Filter.UserID)
			argIdx++
		}
		if params.Filter.ContentID != nil {
			conditions = append(conditions, fmt.Sprintf("content_id = $%d", argIdx))
			args = append(args, *params.Filter.ContentID)
			argIdx++
		}
		if params.Filter.Privacy != nil {
			conditions = append(conditions, fmt.Sprintf("privacy = $%d", argIdx))
			args = append(args, string(*params.Filter.Privacy))
			argIdx++
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(
		`SELECT id, claim, user_id, content_id, "like", quality, agreement, importance,
			confidence, privacy, parts, category, labels, description, review_status,
			categorized_ratings, created_at, updated_at
		FROM perspectives %s
		ORDER BY %s %s, id %s
		LIMIT $%d`,
		whereClause, col, dir, dir, argIdx,
	)
	args = append(args, limit+1)

	var rows []perspectiveRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to list perspectives: %w", err)
	}

	hasNext := len(rows) > limit
	if hasNext {
		rows = rows[:limit]
	}

	hasPrev := params.After != nil

	items := make([]*domain.Perspective, len(rows))
	for i := range rows {
		items[i] = perspectiveRowToDomain(&rows[i])
	}

	result := &domain.PaginatedPerspectives{
		Items:   items,
		HasNext: hasNext,
		HasPrev: hasPrev,
	}

	if len(items) > 0 {
		start := encodeCursor(items[0].ID)
		end := encodeCursor(items[len(items)-1].ID)
		result.StartCursor = &start
		result.EndCursor = &end
	}

	// Optional total count
	if params.IncludeTotalCount {
		var countConditions []string
		var countArgs []interface{}
		countArgIdx := 1

		if params.Filter != nil {
			if params.Filter.UserID != nil {
				countConditions = append(countConditions, fmt.Sprintf("user_id = $%d", countArgIdx))
				countArgs = append(countArgs, *params.Filter.UserID)
				countArgIdx++
			}
			if params.Filter.ContentID != nil {
				countConditions = append(countConditions, fmt.Sprintf("content_id = $%d", countArgIdx))
				countArgs = append(countArgs, *params.Filter.ContentID)
				countArgIdx++
			}
			if params.Filter.Privacy != nil {
				countConditions = append(countConditions, fmt.Sprintf("privacy = $%d", countArgIdx))
				countArgs = append(countArgs, string(*params.Filter.Privacy))
				countArgIdx++
			}
		}

		countQuery := "SELECT COUNT(*) FROM perspectives"
		if len(countConditions) > 0 {
			countQuery += " WHERE " + strings.Join(countConditions, " AND ")
		}

		var count int
		if err := r.db.GetContext(ctx, &count, countQuery, countArgs...); err != nil {
			return nil, fmt.Errorf("failed to count perspectives: %w", err)
		}
		result.TotalCount = &count
	}

	return result, nil
}

// perspectiveRowToDomain converts a database row to a domain Perspective
func perspectiveRowToDomain(row *perspectiveRow) *domain.Perspective {
	p := &domain.Perspective{
		ID:      row.ID,
		Claim:   row.Claim,
		UserID:  row.UserID,
		Privacy: domain.PrivacyPublic,
	}

	if row.ContentID.Valid {
		contentID := int(row.ContentID.Int64)
		p.ContentID = &contentID
	}
	if row.Like.Valid {
		p.Like = &row.Like.String
	}
	if row.Quality.Valid {
		quality := int(row.Quality.Int64)
		p.Quality = &quality
	}
	if row.Agreement.Valid {
		agreement := int(row.Agreement.Int64)
		p.Agreement = &agreement
	}
	if row.Importance.Valid {
		importance := int(row.Importance.Int64)
		p.Importance = &importance
	}
	if row.Confidence.Valid {
		confidence := int(row.Confidence.Int64)
		p.Confidence = &confidence
	}
	if row.Privacy.Valid {
		p.Privacy = domain.Privacy(row.Privacy.String)
	}
	if row.Category.Valid {
		p.Category = &row.Category.String
	}
	if row.Description.Valid {
		p.Description = &row.Description.String
	}
	if row.ReviewStatus.Valid {
		status := domain.ReviewStatus(row.ReviewStatus.String)
		p.ReviewStatus = &status
	}

	// Convert arrays
	if len(row.Parts) > 0 {
		p.Parts = make([]int, len(row.Parts))
		for i, v := range row.Parts {
			p.Parts[i] = int(v)
		}
	}
	if len(row.Labels) > 0 {
		p.Labels = row.Labels
	}

	// Parse categorized ratings
	if len(row.CategorizedRatings) > 0 {
		p.CategorizedRatings = make([]domain.CategorizedRating, 0, len(row.CategorizedRatings))
		for _, data := range row.CategorizedRatings {
			var cr domain.CategorizedRating
			if err := json.Unmarshal(data, &cr); err == nil {
				p.CategorizedRatings = append(p.CategorizedRatings, cr)
			}
		}
	}

	if row.CreatedAt.Valid {
		p.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		p.UpdatedAt = row.UpdatedAt.Time
	}

	return p
}

// toNullInt64FromIntPtr converts an int pointer to sql.NullInt64
func toNullInt64FromIntPtr(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

// toNullStringFromReviewStatus converts a ReviewStatus pointer to sql.NullString
func toNullStringFromReviewStatus(s *domain.ReviewStatus) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(*s), Valid: true}
}

// perspectiveSortColumnName maps a domain sort field to a safe SQL column name
func perspectiveSortColumnName(sortBy domain.PerspectiveSortBy) string {
	switch sortBy {
	case domain.PerspectiveSortByUpdatedAt:
		return "updated_at"
	case domain.PerspectiveSortByClaim:
		return "claim"
	default:
		return "created_at"
	}
}
