package postgres

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yourorg/perspectize-go/internal/core/domain"
)

// ContentRepository implements the ContentRepository interface using PostgreSQL
type ContentRepository struct {
	db *sqlx.DB
}

// contentRow represents the database row structure for content
type contentRow struct {
	ID          int            `db:"id"`
	Name        string         `db:"name"`
	URL         sql.NullString `db:"url"`
	ContentType string         `db:"content_type"`
	Length      sql.NullInt64  `db:"length"`
	LengthUnits sql.NullString `db:"length_units"`
	Response    []byte         `db:"response"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

// NewContentRepository creates a new PostgreSQL content repository
func NewContentRepository(db *sqlx.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

// Create inserts a new content record into the database
func (r *ContentRepository) Create(ctx context.Context, content *domain.Content) (*domain.Content, error) {
	query := `
		INSERT INTO content (name, url, content_type, length, length_units, response)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	var id int
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query,
		content.Name,
		toNullString(content.URL),
		contentTypeToDBValue(content.ContentType),
		toNullInt64(content.Length),
		toNullString(content.LengthUnits),
		content.Response,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert content: %w", err)
	}

	content.ID = id
	if createdAt.Valid {
		content.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		content.UpdatedAt = updatedAt.Time
	}

	return content, nil
}

// GetByID retrieves a content record by its ID
func (r *ContentRepository) GetByID(ctx context.Context, id int) (*domain.Content, error) {
	query := `SELECT id, name, url, content_type, length, length_units, response, created_at, updated_at
		FROM content WHERE id = $1`

	var row contentRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get content by id: %w", err)
	}

	return rowToDomain(&row), nil
}

// GetByURL retrieves a content record by its URL
func (r *ContentRepository) GetByURL(ctx context.Context, url string) (*domain.Content, error) {
	query := `SELECT id, name, url, content_type, length, length_units, response, created_at, updated_at
		FROM content WHERE url = $1`

	var row contentRow
	err := r.db.GetContext(ctx, &row, query, url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get content by url: %w", err)
	}

	return rowToDomain(&row), nil
}

// rowToDomain converts a database row to a domain Content
func rowToDomain(row *contentRow) *domain.Content {
	content := &domain.Content{
		ID:          row.ID,
		Name:        row.Name,
		ContentType: contentTypeFromDBValue(row.ContentType),
		Response:    row.Response,
	}

	if row.URL.Valid {
		content.URL = &row.URL.String
	}
	if row.Length.Valid {
		length := int(row.Length.Int64)
		content.Length = &length
	}
	if row.LengthUnits.Valid {
		content.LengthUnits = &row.LengthUnits.String
	}
	if row.CreatedAt.Valid {
		content.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		content.UpdatedAt = row.UpdatedAt.Time
	}

	return content
}

// contentTypeToDBValue converts domain ContentType to lowercase for database storage
func contentTypeToDBValue(ct domain.ContentType) string {
	return strings.ToLower(string(ct))
}

// contentTypeFromDBValue converts lowercase database value to domain ContentType
func contentTypeFromDBValue(s string) domain.ContentType {
	return domain.ContentType(strings.ToUpper(s))
}

// toNullString converts a string pointer to sql.NullString
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

// toNullInt64 converts an int pointer to sql.NullInt64
func toNullInt64(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

// encodeCursor encodes a content ID as an opaque base64 cursor
func encodeCursor(id int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor:%d", id)))
}

// decodeCursor decodes an opaque base64 cursor back to a content ID
func decodeCursor(cursor string) (int, error) {
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, fmt.Errorf("invalid cursor: %w", err)
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 || parts[0] != "cursor" {
		return 0, fmt.Errorf("invalid cursor format")
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid cursor id: %w", err)
	}
	return id, nil
}

// sortColumnName maps a domain sort field to a safe SQL column name
func sortColumnName(sortBy domain.ContentSortBy) string {
	switch sortBy {
	case domain.ContentSortByUpdatedAt:
		return "updated_at"
	case domain.ContentSortByName:
		return "name"
	case domain.ContentSortByCreatedAt:
		return "created_at"
	default:
		return "created_at"
	}
}

// sortDirection returns a safe SQL sort direction string
func sortDirection(order domain.SortOrder) string {
	if order == domain.SortOrderAsc {
		return "ASC"
	}
	return "DESC"
}

// List retrieves a paginated list of content using cursor-based pagination
func (r *ContentRepository) List(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
	limit := 10
	if params.First != nil {
		limit = *params.First
	}

	col := sortColumnName(params.SortBy)
	dir := sortDirection(params.SortOrder)

	// Build query
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
		if params.Filter.ContentType != nil {
			conditions = append(conditions, fmt.Sprintf("content_type = $%d", argIdx))
			args = append(args, contentTypeToDBValue(*params.Filter.ContentType))
			argIdx++
		}
		if params.Filter.MinLengthSeconds != nil {
			conditions = append(conditions, fmt.Sprintf("length >= $%d", argIdx))
			args = append(args, *params.Filter.MinLengthSeconds)
			argIdx++
		}
		if params.Filter.MaxLengthSeconds != nil {
			conditions = append(conditions, fmt.Sprintf("length <= $%d", argIdx))
			args = append(args, *params.Filter.MaxLengthSeconds)
			argIdx++
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Fetch limit+1 to determine hasNextPage
	query := fmt.Sprintf(
		`SELECT id, name, url, content_type, length, length_units, response, created_at, updated_at
		FROM content %s
		ORDER BY %s %s, id %s
		LIMIT $%d`,
		whereClause, col, dir, dir, argIdx,
	)
	args = append(args, limit+1)

	var rows []contentRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to list content: %w", err)
	}

	hasNext := len(rows) > limit
	if hasNext {
		rows = rows[:limit]
	}

	hasPrev := params.After != nil

	items := make([]*domain.Content, len(rows))
	for i := range rows {
		items[i] = rowToDomain(&rows[i])
	}

	conn := &domain.PaginatedContent{
		Items:   items,
		HasNext: hasNext,
		HasPrev: hasPrev,
	}

	if len(items) > 0 {
		start := encodeCursor(items[0].ID)
		end := encodeCursor(items[len(items)-1].ID)
		conn.StartCursor = &start
		conn.EndCursor = &end
	}

	// Optional total count (respects filters but not cursor pagination)
	if params.IncludeTotalCount {
		var countConditions []string
		var countArgs []interface{}
		countArgIdx := 1

		if params.Filter != nil {
			if params.Filter.ContentType != nil {
				countConditions = append(countConditions, fmt.Sprintf("content_type = $%d", countArgIdx))
				countArgs = append(countArgs, contentTypeToDBValue(*params.Filter.ContentType))
				countArgIdx++
			}
			if params.Filter.MinLengthSeconds != nil {
				countConditions = append(countConditions, fmt.Sprintf("length >= $%d", countArgIdx))
				countArgs = append(countArgs, *params.Filter.MinLengthSeconds)
				countArgIdx++
			}
			if params.Filter.MaxLengthSeconds != nil {
				countConditions = append(countConditions, fmt.Sprintf("length <= $%d", countArgIdx))
				countArgs = append(countArgs, *params.Filter.MaxLengthSeconds)
				countArgIdx++
			}
		}

		countQuery := "SELECT COUNT(*) FROM content"
		if len(countConditions) > 0 {
			countQuery += " WHERE " + strings.Join(countConditions, " AND ")
		}

		var count int
		if err := r.db.GetContext(ctx, &count, countQuery, countArgs...); err != nil {
			return nil, fmt.Errorf("failed to count content: %w", err)
		}
		conn.TotalCount = &count
	}

	return conn, nil
}
