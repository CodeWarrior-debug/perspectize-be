package services_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	portservices "github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/services"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockUserRepository implements repositories.UserRepository for testing
type mockUserRepository struct {
	createFn        func(ctx context.Context, user *domain.User) (*domain.User, error)
	getByIDFn       func(ctx context.Context, id int) (*domain.User, error)
	getByUsernameFn func(ctx context.Context, username string) (*domain.User, error)
	getByEmailFn    func(ctx context.Context, email string) (*domain.User, error)
	listAllFn       func(ctx context.Context) ([]*domain.User, error)
	updateFn        func(ctx context.Context, user *domain.User) (*domain.User, error)
	deleteFn        func(ctx context.Context, id int) error
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

func (m *mockUserRepository) ListAll(ctx context.Context) ([]*domain.User, error) {
	if m.listAllFn != nil {
		return m.listAllFn(ctx)
	}
	return []*domain.User{}, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, user)
	}
	return user, nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id int) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

// mockContentRepoForUser implements repositories.ContentRepository for user tests
type mockContentRepoForUser struct {
	reassignByUserFn func(ctx context.Context, fromUserID, toUserID int) error
}

func (m *mockContentRepoForUser) Create(ctx context.Context, content *domain.Content) (*domain.Content, error) {
	return content, nil
}
func (m *mockContentRepoForUser) GetByID(ctx context.Context, id int) (*domain.Content, error) {
	return nil, domain.ErrNotFound
}
func (m *mockContentRepoForUser) GetByURL(ctx context.Context, url string) (*domain.Content, error) {
	return nil, domain.ErrNotFound
}
func (m *mockContentRepoForUser) List(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
	return &domain.PaginatedContent{Items: []*domain.Content{}}, nil
}
func (m *mockContentRepoForUser) ReassignByUser(ctx context.Context, fromUserID, toUserID int) error {
	if m.reassignByUserFn != nil {
		return m.reassignByUserFn(ctx, fromUserID, toUserID)
	}
	return nil
}

// mockPerspectiveRepoForUser implements repositories.PerspectiveRepository for user tests
type mockPerspectiveRepoForUser struct {
	reassignByUserFn func(ctx context.Context, fromUserID, toUserID int) error
}

func (m *mockPerspectiveRepoForUser) Create(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
	return p, nil
}
func (m *mockPerspectiveRepoForUser) GetByID(ctx context.Context, id int) (*domain.Perspective, error) {
	return nil, domain.ErrNotFound
}
func (m *mockPerspectiveRepoForUser) Update(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
	return p, nil
}
func (m *mockPerspectiveRepoForUser) Delete(ctx context.Context, id int) error {
	return nil
}
func (m *mockPerspectiveRepoForUser) List(ctx context.Context, params domain.PerspectiveListParams) (*domain.PaginatedPerspectives, error) {
	return &domain.PaginatedPerspectives{Items: []*domain.Perspective{}}, nil
}
func (m *mockPerspectiveRepoForUser) ReassignByUser(ctx context.Context, fromUserID, toUserID int) error {
	if m.reassignByUserFn != nil {
		return m.reassignByUserFn(ctx, fromUserID, toUserID)
	}
	return nil
}

// newTestUserService creates a UserService with default mocks for content/perspective repos
func newTestUserService(repo *mockUserRepository) *services.UserService {
	return services.NewUserService(repo, &mockContentRepoForUser{}, &mockPerspectiveRepoForUser{})
}

// newTestUserServiceFull creates a UserService with explicit content/perspective repo mocks
func newTestUserServiceFull(repo *mockUserRepository, contentRepo *mockContentRepoForUser, perspectiveRepo *mockPerspectiveRepoForUser) *services.UserService {
	return services.NewUserService(repo, contentRepo, perspectiveRepo)
}

// --- Create Tests ---

func TestCreate_Success(t *testing.T) {
	repo := &mockUserRepository{
		createFn: func(ctx context.Context, user *domain.User) (*domain.User, error) {
			user.ID = 1
			return user, nil
		},
	}

	svc := newTestUserService(repo)
	result, err := svc.Create(context.Background(), "testuser", "test@example.com")

	require.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, "test@example.com", result.Email)
	assert.True(t, result.Active)
}

func TestCreate_UsernameEmpty(t *testing.T) {
	repo := &mockUserRepository{}
	svc := newTestUserService(repo)

	result, err := svc.Create(context.Background(), "", "test@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username is required")
}

func TestCreate_UsernameWhitespace(t *testing.T) {
	repo := &mockUserRepository{}
	svc := newTestUserService(repo)

	result, err := svc.Create(context.Background(), "   ", "test@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username is required")
}

func TestCreate_UsernameTooLong(t *testing.T) {
	repo := &mockUserRepository{}
	svc := newTestUserService(repo)

	// Username with 25 characters (limit is 24)
	result, err := svc.Create(context.Background(), "abcdefghijklmnopqrstuvwxy", "test@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username must be 24 characters or less")
}

func TestCreate_EmailEmpty_Succeeds(t *testing.T) {
	repo := &mockUserRepository{
		createFn: func(ctx context.Context, user *domain.User) (*domain.User, error) {
			user.ID = 1
			return user, nil
		},
	}

	svc := newTestUserService(repo)
	result, err := svc.Create(context.Background(), "testuser", "")

	require.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, "", result.Email)
	assert.True(t, result.Active)
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
	svc := newTestUserService(repo)

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

	svc := newTestUserService(repo)
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

	svc := newTestUserService(repo)
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

	svc := newTestUserService(repo)
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

	svc := newTestUserService(repo)
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

	svc := newTestUserService(repo)
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

	svc := newTestUserService(repo)
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

	svc := newTestUserService(repo)
	result, err := svc.GetByID(context.Background(), 999)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestUserGetByID_InvalidID_Zero(t *testing.T) {
	repo := &mockUserRepository{}
	svc := newTestUserService(repo)

	result, err := svc.GetByID(context.Background(), 0)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "user id must be a positive integer")
}

func TestUserGetByID_InvalidID_Negative(t *testing.T) {
	repo := &mockUserRepository{}
	svc := newTestUserService(repo)

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

	svc := newTestUserService(repo)
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

	svc := newTestUserService(repo)
	result, err := svc.GetByUsername(context.Background(), "nonexistent")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestGetByUsername_Empty(t *testing.T) {
	repo := &mockUserRepository{}
	svc := newTestUserService(repo)

	result, err := svc.GetByUsername(context.Background(), "")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username is required")
}

func TestGetByUsername_Whitespace(t *testing.T) {
	repo := &mockUserRepository{}
	svc := newTestUserService(repo)

	result, err := svc.GetByUsername(context.Background(), "   ")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username is required")
}

// --- NewUserService Tests ---

func TestNewUserService(t *testing.T) {
	repo := &mockUserRepository{}
	svc := newTestUserService(repo)

	assert.NotNil(t, svc)
}

// --- Update Tests ---

func TestUpdate_Success(t *testing.T) {
	newUsername := "updateduser"
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{ID: 2, Username: "olduser", Email: "old@example.com"}, nil
		},
		updateFn: func(ctx context.Context, user *domain.User) (*domain.User, error) {
			return user, nil
		},
	}

	svc := newTestUserService(repo)
	result, err := svc.Update(context.Background(), portservices.UpdateUserInput{
		ID:       2,
		Username: &newUsername,
	})

	require.NoError(t, err)
	assert.Equal(t, "updateduser", result.Username)
	assert.Equal(t, "old@example.com", result.Email)
}

func TestUpdate_SentinelUserBlocked(t *testing.T) {
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{ID: 1, Username: domain.DeletedUserUsername, Email: "deleted@system.internal"}, nil
		},
	}

	newUsername := "hacker"
	svc := newTestUserService(repo)
	result, err := svc.Update(context.Background(), portservices.UpdateUserInput{
		ID:       1,
		Username: &newUsername,
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrSentinelUser))
}

func TestUpdate_InvalidID(t *testing.T) {
	repo := &mockUserRepository{}
	svc := newTestUserService(repo)

	result, err := svc.Update(context.Background(), portservices.UpdateUserInput{ID: 0})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestUpdate_UsernameAlreadyTaken(t *testing.T) {
	takenName := "taken"
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{ID: 2, Username: "olduser", Email: "old@example.com"}, nil
		},
		getByUsernameFn: func(ctx context.Context, username string) (*domain.User, error) {
			if username == "taken" {
				return &domain.User{ID: 3, Username: "taken"}, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	svc := newTestUserService(repo)
	result, err := svc.Update(context.Background(), portservices.UpdateUserInput{
		ID:       2,
		Username: &takenName,
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrAlreadyExists))
}

// --- Delete Tests ---

func TestDelete_Success(t *testing.T) {
	sentinelUser := &domain.User{ID: 1, Username: domain.DeletedUserUsername}
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{ID: 2, Username: "testuser", Email: "test@example.com"}, nil
		},
		getByUsernameFn: func(ctx context.Context, username string) (*domain.User, error) {
			if username == domain.DeletedUserUsername {
				return sentinelUser, nil
			}
			return nil, domain.ErrNotFound
		},
		deleteFn: func(ctx context.Context, id int) error {
			assert.Equal(t, 2, id)
			return nil
		},
	}

	contentRepo := &mockContentRepoForUser{}
	perspectiveRepo := &mockPerspectiveRepoForUser{}
	svc := newTestUserServiceFull(repo, contentRepo, perspectiveRepo)

	err := svc.Delete(context.Background(), 2)
	require.NoError(t, err)
}

func TestDelete_SentinelUserBlocked(t *testing.T) {
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{ID: 1, Username: domain.DeletedUserUsername}, nil
		},
	}

	svc := newTestUserService(repo)
	err := svc.Delete(context.Background(), 1)

	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrDeleteSentinel))
}

func TestDelete_InvalidID(t *testing.T) {
	repo := &mockUserRepository{}
	svc := newTestUserService(repo)

	err := svc.Delete(context.Background(), 0)

	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestDelete_UserNotFound(t *testing.T) {
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := newTestUserService(repo)
	err := svc.Delete(context.Background(), 999)

	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestDelete_ReassignContentFails(t *testing.T) {
	sentinelUser := &domain.User{ID: 1, Username: domain.DeletedUserUsername}
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{ID: 2, Username: "testuser"}, nil
		},
		getByUsernameFn: func(ctx context.Context, username string) (*domain.User, error) {
			if username == domain.DeletedUserUsername {
				return sentinelUser, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	contentRepo := &mockContentRepoForUser{
		reassignByUserFn: func(ctx context.Context, fromUserID, toUserID int) error {
			return fmt.Errorf("database error")
		},
	}
	perspectiveRepo := &mockPerspectiveRepoForUser{}
	svc := newTestUserServiceFull(repo, contentRepo, perspectiveRepo)

	err := svc.Delete(context.Background(), 2)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to reassign content")
}

func TestDelete_ReassignPerspectivesFails(t *testing.T) {
	sentinelUser := &domain.User{ID: 1, Username: domain.DeletedUserUsername}
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{ID: 2, Username: "testuser"}, nil
		},
		getByUsernameFn: func(ctx context.Context, username string) (*domain.User, error) {
			if username == domain.DeletedUserUsername {
				return sentinelUser, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	contentRepo := &mockContentRepoForUser{}
	perspectiveRepo := &mockPerspectiveRepoForUser{
		reassignByUserFn: func(ctx context.Context, fromUserID, toUserID int) error {
			return fmt.Errorf("database error")
		},
	}
	svc := newTestUserServiceFull(repo, contentRepo, perspectiveRepo)

	err := svc.Delete(context.Background(), 2)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to reassign perspectives")
}

// --- Reserved Username Tests ---

func TestCreate_ReservedDeletedUsername(t *testing.T) {
	repo := &mockUserRepository{}
	svc := newTestUserService(repo)

	result, err := svc.Create(context.Background(), domain.DeletedUserUsername, "test@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username is reserved")
}

func TestCreate_ReservedSystemUsername(t *testing.T) {
	repo := &mockUserRepository{}
	svc := newTestUserService(repo)

	result, err := svc.Create(context.Background(), domain.SystemUserUsername, "test@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username is reserved")
}

func TestUpdate_ReservedDeletedUsername(t *testing.T) {
	reserved := domain.DeletedUserUsername
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{ID: 2, Username: "normaluser", Email: "normal@example.com"}, nil
		},
	}

	svc := newTestUserService(repo)
	result, err := svc.Update(context.Background(), portservices.UpdateUserInput{
		ID:       2,
		Username: &reserved,
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username is reserved")
}

func TestUpdate_ReservedSystemUsername(t *testing.T) {
	reserved := domain.SystemUserUsername
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{ID: 2, Username: "normaluser", Email: "normal@example.com"}, nil
		},
	}

	svc := newTestUserService(repo)
	result, err := svc.Update(context.Background(), portservices.UpdateUserInput{
		ID:       2,
		Username: &reserved,
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username is reserved")
}

func TestUpdate_SystemSentinelBlocked(t *testing.T) {
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{ID: 1, Username: domain.SystemUserUsername, Email: "system@system.internal"}, nil
		},
	}

	newUsername := "hacker"
	svc := newTestUserService(repo)
	result, err := svc.Update(context.Background(), portservices.UpdateUserInput{
		ID:       1,
		Username: &newUsername,
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrSentinelUser))
}

func TestDelete_SystemSentinelBlocked(t *testing.T) {
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{ID: 1, Username: domain.SystemUserUsername}, nil
		},
	}

	svc := newTestUserService(repo)
	err := svc.Delete(context.Background(), 1)

	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrDeleteSentinel))
}
