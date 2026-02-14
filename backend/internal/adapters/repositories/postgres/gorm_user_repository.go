package postgres

import (
	"context"
	"errors"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	repositories "github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/repositories"
	"gorm.io/gorm"
)

// GormUserRepository implements the UserRepository interface using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// Compile-time interface check
var _ repositories.UserRepository = (*GormUserRepository)(nil)

// NewGormUserRepository creates a new GORM-based user repository
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Create inserts a new user record into the database
func (r *GormUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	model := userDomainToModel(user)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	// GORM auto-fills ID, CreatedAt, UpdatedAt
	return userModelToDomain(model), nil
}

// GetByID retrieves a user by their ID
func (r *GormUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return userModelToDomain(&model), nil
}

// GetByUsername retrieves a user by their username
func (r *GormUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).Where("username = ?", username).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return userModelToDomain(&model), nil
}

// GetByEmail retrieves a user by their email address
func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return userModelToDomain(&model), nil
}

// ListAll retrieves all users ordered by username
func (r *GormUserRepository) ListAll(ctx context.Context) ([]*domain.User, error) {
	var models []UserModel

	err := r.db.WithContext(ctx).Order("username ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(models))
	for i := range models {
		users[i] = userModelToDomain(&models[i])
	}

	return users, nil
}
