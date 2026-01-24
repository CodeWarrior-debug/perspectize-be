package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
		string(content.ContentType),
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
		ContentType: domain.ContentType(row.ContentType),
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
