package postgres

// GORM PROTOTYPE — User repository using GORM with hex-clean separate models.
//
// Current sqlx version: 132 lines (user_repository.go)
// This GORM version:    ~75 lines (this file) + ~15 lines (shared in gorm_mappers.go)
//
// Biggest win: No sql.NullTime handling, no QueryRowContext/Scan, no GetContext boilerplate.

import (
	"context"
	"errors"
	"fmt"

	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/domain"
	"gorm.io/gorm"
)

// GORMUserRepository implements UserRepository using GORM.
type GORMUserRepository struct {
	db *gorm.DB
}

// NewGORMUserRepository creates a new GORM-based user repository.
func NewGORMUserRepository(db *gorm.DB) *GORMUserRepository {
	return &GORMUserRepository{db: db}
}

// Create inserts a new user.
func (r *GORMUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	model := userDomainToModel(user)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return r.GetByID(ctx, model.ID)
}

// GetByID retrieves a user by ID.
func (r *GORMUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return userModelToDomain(&model), nil
}

// GetByUsername retrieves a user by username.
func (r *GORMUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return userModelToDomain(&model), nil
}

// GetByEmail retrieves a user by email.
func (r *GORMUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return userModelToDomain(&model), nil
}

// ListAll retrieves all users ordered by username.
func (r *GORMUserRepository) ListAll(ctx context.Context) ([]*domain.User, error) {
	var models []UserModel
	if err := r.db.WithContext(ctx).Order("username ASC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*domain.User, len(models))
	for i := range models {
		users[i] = userModelToDomain(&models[i])
	}
	return users, nil
}
