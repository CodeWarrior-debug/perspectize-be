package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/repositories"
)

// UserService implements business logic for user operations
type UserService struct {
	repo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// emailRegex validates email format
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Create creates a new user with validation
func (s *UserService) Create(ctx context.Context, username, email string) (*domain.User, error) {
	// Validate username
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, fmt.Errorf("%w: username is required", domain.ErrInvalidInput)
	}
	if len(username) > 24 {
		return nil, fmt.Errorf("%w: username must be 24 characters or less", domain.ErrInvalidInput)
	}

	// Validate email
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, fmt.Errorf("%w: email is required", domain.ErrInvalidInput)
	}
	if !emailRegex.MatchString(email) {
		return nil, fmt.Errorf("%w: invalid email format", domain.ErrInvalidInput)
	}

	// Check if username already exists
	existing, err := s.repo.GetByUsername(ctx, username)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("%w: username already taken", domain.ErrAlreadyExists)
	}
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}

	// Check if email already exists
	existing, err = s.repo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("%w: email already registered", domain.ErrAlreadyExists)
	}
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	user := &domain.User{
		Username: username,
		Email:    email,
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return created, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id int) (*domain.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: user id must be a positive integer", domain.ErrInvalidInput)
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetByUsername retrieves a user by username
func (s *UserService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, fmt.Errorf("%w: username is required", domain.ErrInvalidInput)
	}

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// ListAll retrieves all users
func (s *UserService) ListAll(ctx context.Context) ([]*domain.User, error) {
	users, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}
