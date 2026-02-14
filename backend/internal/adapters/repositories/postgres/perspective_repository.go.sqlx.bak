package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/jmoiron/sqlx"
)

// JSONBArray is a custom type for PostgreSQL jsonb[] columns.
// It stores JSON data as strings and properly handles NULL values.
type JSONBArray []string

// Scan implements the sql.Scanner interface for reading jsonb[] from PostgreSQL
func (a *JSONBArray) Scan(src interface{}) error {
	if src == nil {
		*a = nil
		return nil
	}

	// PostgreSQL returns jsonb[] as []byte containing array of bytea
	// We need to scan through StringArray
	var strArray StringArray
	if err := strArray.Scan(src); err != nil {
		return err
	}
	*a = JSONBArray(strArray)
	return nil
}

// Value implements the driver.Valuer interface for writing jsonb[] to PostgreSQL
func (a JSONBArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return StringArray(a).Value()
}

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
	Parts              Int64Array     `db:"parts"`
	Category           sql.NullString `db:"category"`
	Labels             StringArray    `db:"labels"`
	Description        sql.NullString `db:"description"`
	ReviewStatus       sql.NullString `db:"review_status"`
	CategorizedRatings JSONBArray     `db:"categorized_ratings"`
	CreatedAt          sql.NullTime   `db:"created_at"`
	UpdatedAt          sql.NullTime   `db:"updated_at"`
}

// NewPerspectiveRepository creates a new PostgreSQL perspective repository
func NewPerspectiveRepository(db *sqlx.DB) *PerspectiveRepository {
	return &PerspectiveRepository{db: db}
}

// Create inserts a new perspective record into the database
func (r *PerspectiveRepository) Create(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
	// Marshal categorized ratings to JSON strings for jsonb[] column
	var categorizedRatings JSONBArray
	if len(p.CategorizedRatings) > 0 {
		categorizedRatings = make(JSONBArray, len(p.CategorizedRatings))
		for i, cr := range p.CategorizedRatings {
			data, err := json.Marshal(cr)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal categorized rating: %w", err)
			}
			categorizedRatings[i] = string(data)
		}
	}

	query := `
		INSERT INTO perspectives (
			claim, user_id, content_id, "like", quality, agreement, importance, confidence,
			privacy, parts, category, labels, description, review_status, categorized_ratings
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15::jsonb[]
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
		privacyToDBValue(p.Privacy),
		intSliceToInt64Array(p.Parts),
		toNullString(p.Category),
		StringArray(p.Labels),
		toNullString(p.Description),
		reviewStatusToDBValue(p.ReviewStatus),
		categorizedRatings,
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
	// Marshal categorized ratings to JSON strings for jsonb[] column
	var categorizedRatings JSONBArray
	if len(p.CategorizedRatings) > 0 {
		categorizedRatings = make(JSONBArray, len(p.CategorizedRatings))
		for i, cr := range p.CategorizedRatings {
			data, err := json.Marshal(cr)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal categorized rating: %w", err)
			}
			categorizedRatings[i] = string(data)
		}
	}

	query := `
		UPDATE perspectives SET
			claim = $1, content_id = $2, "like" = $3, quality = $4, agreement = $5,
			importance = $6, confidence = $7, privacy = $8, parts = $9, category = $10,
			labels = $11, description = $12, review_status = $13, categorized_ratings = $14::jsonb[]
		WHERE id = $15`

	_, err := r.db.ExecContext(ctx, query,
		p.Claim,
		toNullInt64FromIntPtr(p.ContentID),
		toNullString(p.Like),
		toNullInt64FromIntPtr(p.Quality),
		toNullInt64FromIntPtr(p.Agreement),
		toNullInt64FromIntPtr(p.Importance),
		toNullInt64FromIntPtr(p.Confidence),
		privacyToDBValue(p.Privacy),
		intSliceToInt64Array(p.Parts),
		toNullString(p.Category),
		StringArray(p.Labels),
		toNullString(p.Description),
		reviewStatusToDBValue(p.ReviewStatus),
		categorizedRatings,
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
			args = append(args, privacyToDBValue(*params.Filter.Privacy))
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
				countArgs = append(countArgs, privacyToDBValue(*params.Filter.Privacy))
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
		p.Privacy = privacyFromDBValue(row.Privacy.String)
	}
	if row.Category.Valid {
		p.Category = &row.Category.String
	}
	if row.Description.Valid {
		p.Description = &row.Description.String
	}
	p.ReviewStatus = reviewStatusFromDBValue(row.ReviewStatus)

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

	// Parse categorized ratings from JSON strings
	if len(row.CategorizedRatings) > 0 {
		p.CategorizedRatings = make([]domain.CategorizedRating, 0, len(row.CategorizedRatings))
		for _, jsonStr := range row.CategorizedRatings {
			var cr domain.CategorizedRating
			if err := json.Unmarshal([]byte(jsonStr), &cr); err != nil {
				slog.Warn("failed to parse categorized rating JSON", "error", err)
				continue
			}
			p.CategorizedRatings = append(p.CategorizedRatings, cr)
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

// privacyToDBValue converts domain Privacy to lowercase for database storage
func privacyToDBValue(p domain.Privacy) string {
	return strings.ToLower(string(p))
}

// privacyFromDBValue converts lowercase database value to domain Privacy
func privacyFromDBValue(s string) domain.Privacy {
	return domain.Privacy(strings.ToUpper(s))
}

// reviewStatusToDBValue converts domain ReviewStatus to lowercase for database storage
func reviewStatusToDBValue(rs *domain.ReviewStatus) sql.NullString {
	if rs == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: strings.ToLower(string(*rs)), Valid: true}
}

// reviewStatusFromDBValue converts lowercase database value to domain ReviewStatus
func reviewStatusFromDBValue(s sql.NullString) *domain.ReviewStatus {
	if !s.Valid {
		return nil
	}
	rs := domain.ReviewStatus(strings.ToUpper(s.String))
	return &rs
}

// toNullStringFromReviewStatus converts a ReviewStatus pointer to sql.NullString
// Deprecated: Use reviewStatusToDBValue instead
func toNullStringFromReviewStatus(s *domain.ReviewStatus) sql.NullString {
	return reviewStatusToDBValue(s)
}

// perspectiveSortColumnName maps a domain sort field to a safe SQL column name
func perspectiveSortColumnName(sortBy domain.PerspectiveSortBy) string {
	switch sortBy {
	case domain.PerspectiveSortByUpdatedAt:
		return "updated_at"
	case domain.PerspectiveSortByClaim:
		return "claim"
	case domain.PerspectiveSortByCreatedAt:
		return "created_at"
	default:
		return "created_at"
	}
}

// intSliceToInt64Array converts []int to Int64Array for database storage
func intSliceToInt64Array(ints []int) Int64Array {
	if ints == nil {
		return nil
	}
	result := make(Int64Array, len(ints))
	for i, v := range ints {
		result[i] = int64(v)
	}
	return result
}
