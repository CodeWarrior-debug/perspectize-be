package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/domain"
)

// UserRepository implements the UserRepository interface using PostgreSQL
type UserRepository struct {
	db *sqlx.DB
}

// userRow represents the database row structure for users
type userRow struct {
	ID        int          `db:"id"`
	Username  string       `db:"username"`
	Email     string       `db:"email"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user record into the database
func (r *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (username, email)
		VALUES ($1, $2)
		RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	// Fetch the complete user record with timestamps
	return r.GetByID(ctx, id)
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1`

	var row userRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return userRowToDomain(&row), nil
}

// GetByUsername retrieves a user by their username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE username = $1`

	var row userRow
	err := r.db.GetContext(ctx, &row, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return userRowToDomain(&row), nil
}

// GetByEmail retrieves a user by their email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE email = $1`

	var row userRow
	err := r.db.GetContext(ctx, &row, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return userRowToDomain(&row), nil
}

// ListAll retrieves all users ordered by username
func (r *UserRepository) ListAll(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users ORDER BY username ASC`

	var rows []userRow
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Return empty slice (not nil) if no users found
	users := make([]*domain.User, 0, len(rows))
	for i := range rows {
		users = append(users, userRowToDomain(&rows[i]))
	}

	return users, nil
}

// userRowToDomain converts a database row to a domain User
func userRowToDomain(row *userRow) *domain.User {
	user := &domain.User{
		ID:       row.ID,
		Username: row.Username,
		Email:    row.Email,
	}

	if row.CreatedAt.Valid {
		user.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		user.UpdatedAt = row.UpdatedAt.Time
	}

	return user
}
