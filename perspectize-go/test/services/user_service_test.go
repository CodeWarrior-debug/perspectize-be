package services_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourorg/perspectize-go/internal/core/domain"
	"github.com/yourorg/perspectize-go/internal/core/services"
)

// mockUserRepository implements repositories.UserRepository for testing
type mockUserRepository struct {
	createFn        func(ctx context.Context, user *domain.User) (*domain.User, error)
	getByIDFn       func(ctx context.Context, id int) (*domain.User, error)
	getByUsernameFn func(ctx context.Context, username string) (*domain.User, error)
	getByEmailFn    func(ctx context.Context, email string) (*domain.User, error)
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	user.ID = 1
	return user, nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	if m.getByUsernameFn != nil {
		return m.getByUsernameFn(ctx, username)
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return nil, domain.ErrNotFound
}

// --- Create Tests ---

func TestCreate_Success(t *testing.T) {
	repo := &mockUserRepository{
		createFn: func(ctx context.Context, user *domain.User) (*domain.User, error) {
			user.ID = 1
			return user, nil
		},
	}

	svc := services.NewUserService(repo)
	result, err := svc.Create(context.Background(), "testuser", "test@example.com")

	require.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestCreate_UsernameEmpty(t *testing.T) {
	repo := &mockUserRepository{}
	svc := services.NewUserService(repo)

	result, err := svc.Create(context.Background(), "", "test@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username is required")
}

func TestCreate_UsernameWhitespace(t *testing.T) {
	repo := &mockUserRepository{}
	svc := services.NewUserService(repo)

	result, err := svc.Create(context.Background(), "   ", "test@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username is required")
}

func TestCreate_UsernameTooLong(t *testing.T) {
	repo := &mockUserRepository{}
	svc := services.NewUserService(repo)

	// Username with 25 characters (limit is 24)
	result, err := svc.Create(context.Background(), "abcdefghijklmnopqrstuvwxy", "test@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username must be 24 characters or less")
}

func TestCreate_EmailEmpty(t *testing.T) {
	repo := &mockUserRepository{}
	svc := services.NewUserService(repo)

	result, err := svc.Create(context.Background(), "testuser", "")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "email is required")
}

func TestCreate_EmailInvalid(t *testing.T) {
	testCases := []string{
		"notanemail",
		"@example.com",
		"test@",
		"test@.com",
		"test@com",
	}

	repo := &mockUserRepository{}
	svc := services.NewUserService(repo)

	for _, email := range testCases {
		t.Run(email, func(t *testing.T) {
			result, err := svc.Create(context.Background(), "testuser", email)

			assert.Nil(t, result)
			require.Error(t, err)
			assert.True(t, errors.Is(err, domain.ErrInvalidInput))
			assert.Contains(t, err.Error(), "invalid email format")
		})
	}
}

func TestCreate_UsernameAlreadyExists(t *testing.T) {
	repo := &mockUserRepository{
		getByUsernameFn: func(ctx context.Context, username string) (*domain.User, error) {
			return &domain.User{ID: 1, Username: username}, nil
		},
	}

	svc := services.NewUserService(repo)
	result, err := svc.Create(context.Background(), "existinguser", "new@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrAlreadyExists))
	assert.Contains(t, err.Error(), "username already taken")
}

func TestCreate_EmailAlreadyExists(t *testing.T) {
	repo := &mockUserRepository{
		getByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{ID: 1, Email: email}, nil
		},
	}

	svc := services.NewUserService(repo)
	result, err := svc.Create(context.Background(), "newuser", "existing@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrAlreadyExists))
	assert.Contains(t, err.Error(), "email already registered")
}

func TestCreate_RepositoryCreateError(t *testing.T) {
	repo := &mockUserRepository{
		createFn: func(ctx context.Context, user *domain.User) (*domain.User, error) {
			return nil, fmt.Errorf("database error")
		},
	}

	svc := services.NewUserService(repo)
	result, err := svc.Create(context.Background(), "testuser", "test@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
}

func TestCreate_GetByUsernameUnexpectedError(t *testing.T) {
	repo := &mockUserRepository{
		getByUsernameFn: func(ctx context.Context, username string) (*domain.User, error) {
			return nil, fmt.Errorf("unexpected database error")
		},
	}

	svc := services.NewUserService(repo)
	result, err := svc.Create(context.Background(), "testuser", "test@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check username")
}

func TestCreate_GetByEmailUnexpectedError(t *testing.T) {
	repo := &mockUserRepository{
		getByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, fmt.Errorf("unexpected database error")
		},
	}

	svc := services.NewUserService(repo)
	result, err := svc.Create(context.Background(), "testuser", "test@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check email")
}

// --- GetByID Tests ---

func TestUserGetByID_Success(t *testing.T) {
	expected := &domain.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			assert.Equal(t, 1, id)
			return expected, nil
		},
	}

	svc := services.NewUserService(repo)
	result, err := svc.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestUserGetByID_NotFound(t *testing.T) {
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := services.NewUserService(repo)
	result, err := svc.GetByID(context.Background(), 999)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestUserGetByID_InvalidID_Zero(t *testing.T) {
	repo := &mockUserRepository{}
	svc := services.NewUserService(repo)

	result, err := svc.GetByID(context.Background(), 0)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "user id must be a positive integer")
}

func TestUserGetByID_InvalidID_Negative(t *testing.T) {
	repo := &mockUserRepository{}
	svc := services.NewUserService(repo)

	result, err := svc.GetByID(context.Background(), -5)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

// --- GetByUsername Tests ---

func TestGetByUsername_Success(t *testing.T) {
	expected := &domain.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	repo := &mockUserRepository{
		getByUsernameFn: func(ctx context.Context, username string) (*domain.User, error) {
			assert.Equal(t, "testuser", username)
			return expected, nil
		},
	}

	svc := services.NewUserService(repo)
	result, err := svc.GetByUsername(context.Background(), "testuser")

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetByUsername_NotFound(t *testing.T) {
	repo := &mockUserRepository{
		getByUsernameFn: func(ctx context.Context, username string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := services.NewUserService(repo)
	result, err := svc.GetByUsername(context.Background(), "nonexistent")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestGetByUsername_Empty(t *testing.T) {
	repo := &mockUserRepository{}
	svc := services.NewUserService(repo)

	result, err := svc.GetByUsername(context.Background(), "")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username is required")
}

func TestGetByUsername_Whitespace(t *testing.T) {
	repo := &mockUserRepository{}
	svc := services.NewUserService(repo)

	result, err := svc.GetByUsername(context.Background(), "   ")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username is required")
}

// --- NewUserService Tests ---

func TestNewUserService(t *testing.T) {
	repo := &mockUserRepository{}
	svc := services.NewUserService(repo)

	assert.NotNil(t, svc)
}
