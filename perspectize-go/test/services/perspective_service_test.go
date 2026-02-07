package services_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/domain"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/services"
)

// mockPerspectiveRepository implements repositories.PerspectiveRepository for testing
type mockPerspectiveRepository struct {
	createFn            func(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error)
	getByIDFn           func(ctx context.Context, id int) (*domain.Perspective, error)
	getByUserAndClaimFn func(ctx context.Context, userID int, claim string) (*domain.Perspective, error)
	updateFn            func(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error)
	deleteFn            func(ctx context.Context, id int) error
	listFn              func(ctx context.Context, params domain.PerspectiveListParams) (*domain.PaginatedPerspectives, error)
}

func (m *mockPerspectiveRepository) Create(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
	if m.createFn != nil {
		return m.createFn(ctx, p)
	}
	p.ID = 1
	return p, nil
}

func (m *mockPerspectiveRepository) GetByID(ctx context.Context, id int) (*domain.Perspective, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *mockPerspectiveRepository) GetByUserAndClaim(ctx context.Context, userID int, claim string) (*domain.Perspective, error) {
	if m.getByUserAndClaimFn != nil {
		return m.getByUserAndClaimFn(ctx, userID, claim)
	}
	return nil, domain.ErrNotFound
}

func (m *mockPerspectiveRepository) Update(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, p)
	}
	return p, nil
}

func (m *mockPerspectiveRepository) Delete(ctx context.Context, id int) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockPerspectiveRepository) List(ctx context.Context, params domain.PerspectiveListParams) (*domain.PaginatedPerspectives, error) {
	if m.listFn != nil {
		return m.listFn(ctx, params)
	}
	return &domain.PaginatedPerspectives{Items: []*domain.Perspective{}}, nil
}

// mockUserRepoForPerspective implements repositories.UserRepository for perspective tests
type mockUserRepoForPerspective struct {
	getByIDFn func(ctx context.Context, id int) (*domain.User, error)
}

func (m *mockUserRepoForPerspective) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	return user, nil
}

func (m *mockUserRepoForPerspective) GetByID(ctx context.Context, id int) (*domain.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return &domain.User{ID: id, Username: "testuser", Email: "test@example.com"}, nil
}

func (m *mockUserRepoForPerspective) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return nil, domain.ErrNotFound
}

func (m *mockUserRepoForPerspective) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, domain.ErrNotFound
}

// --- Create Tests ---

func TestPerspectiveCreate_Success(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{
		createFn: func(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
			p.ID = 1
			return p, nil
		},
	}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	input := services.CreatePerspectiveInput{
		Claim:  "This is a test claim",
		UserID: 1,
	}

	result, err := svc.Create(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "This is a test claim", result.Claim)
	assert.Equal(t, 1, result.UserID)
	assert.Equal(t, domain.PrivacyPublic, result.Privacy)
}

func TestPerspectiveCreate_WithRatings(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{
		createFn: func(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
			p.ID = 1
			return p, nil
		},
	}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	quality := 8000
	agreement := 5000
	input := services.CreatePerspectiveInput{
		Claim:     "Test claim with ratings",
		UserID:    1,
		Quality:   &quality,
		Agreement: &agreement,
	}

	result, err := svc.Create(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, &quality, result.Quality)
	assert.Equal(t, &agreement, result.Agreement)
}

func TestPerspectiveCreate_ClaimEmpty(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	input := services.CreatePerspectiveInput{
		Claim:  "",
		UserID: 1,
	}

	result, err := svc.Create(context.Background(), input)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "claim is required")
}

func TestPerspectiveCreate_ClaimTooLong(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	longClaim := ""
	for i := 0; i < 256; i++ {
		longClaim += "a"
	}
	input := services.CreatePerspectiveInput{
		Claim:  longClaim,
		UserID: 1,
	}

	result, err := svc.Create(context.Background(), input)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "claim must be 255 characters or less")
}

func TestPerspectiveCreate_UserNotFound(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{}
	userRepo := &mockUserRepoForPerspective{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	input := services.CreatePerspectiveInput{
		Claim:  "Test claim",
		UserID: 999,
	}

	result, err := svc.Create(context.Background(), input)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestPerspectiveCreate_InvalidUserID(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	input := services.CreatePerspectiveInput{
		Claim:  "Test claim",
		UserID: 0,
	}

	result, err := svc.Create(context.Background(), input)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestPerspectiveCreate_RatingTooHigh(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	quality := 10001
	input := services.CreatePerspectiveInput{
		Claim:   "Test claim",
		UserID:  1,
		Quality: &quality,
	}

	result, err := svc.Create(context.Background(), input)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidRating))
}

func TestPerspectiveCreate_RatingNegative(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	agreement := -1
	input := services.CreatePerspectiveInput{
		Claim:     "Test claim",
		UserID:    1,
		Agreement: &agreement,
	}

	result, err := svc.Create(context.Background(), input)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidRating))
}

func TestPerspectiveCreate_DuplicateClaim(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{
		getByUserAndClaimFn: func(ctx context.Context, userID int, claim string) (*domain.Perspective, error) {
			return &domain.Perspective{ID: 1, Claim: claim, UserID: userID}, nil
		},
	}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	input := services.CreatePerspectiveInput{
		Claim:  "Existing claim",
		UserID: 1,
	}

	result, err := svc.Create(context.Background(), input)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrDuplicateClaim))
}

func TestPerspectiveCreate_RepositoryError(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{
		createFn: func(ctx context.Context, p *domain.Perspective) (*domain.Perspective, error) {
			return nil, fmt.Errorf("database error")
		},
	}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	input := services.CreatePerspectiveInput{
		Claim:  "Test claim",
		UserID: 1,
	}

	result, err := svc.Create(context.Background(), input)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create perspective")
}

// --- GetByID Tests ---

func TestPerspectiveGetByID_Success(t *testing.T) {
	expected := &domain.Perspective{
		ID:      1,
		Claim:   "Test claim",
		UserID:  1,
		Privacy: domain.PrivacyPublic,
	}

	perspectiveRepo := &mockPerspectiveRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.Perspective, error) {
			return expected, nil
		},
	}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	result, err := svc.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestPerspectiveGetByID_NotFound(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.Perspective, error) {
			return nil, domain.ErrNotFound
		},
	}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	result, err := svc.GetByID(context.Background(), 999)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestPerspectiveGetByID_InvalidID(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	result, err := svc.GetByID(context.Background(), 0)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

// --- Delete Tests ---

func TestPerspectiveDelete_Success(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.Perspective, error) {
			return &domain.Perspective{ID: id}, nil
		},
		deleteFn: func(ctx context.Context, id int) error {
			return nil
		},
	}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	err := svc.Delete(context.Background(), 1)

	require.NoError(t, err)
}

func TestPerspectiveDelete_NotFound(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.Perspective, error) {
			return nil, domain.ErrNotFound
		},
	}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	err := svc.Delete(context.Background(), 999)

	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestPerspectiveDelete_InvalidID(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	err := svc.Delete(context.Background(), 0)

	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

// --- ListPerspectives Tests ---

func TestPerspectiveList_Success(t *testing.T) {
	expected := &domain.PaginatedPerspectives{
		Items: []*domain.Perspective{
			{ID: 1, Claim: "Claim 1", UserID: 1},
			{ID: 2, Claim: "Claim 2", UserID: 1},
		},
		HasNext: false,
		HasPrev: false,
	}

	perspectiveRepo := &mockPerspectiveRepository{
		listFn: func(ctx context.Context, params domain.PerspectiveListParams) (*domain.PaginatedPerspectives, error) {
			return expected, nil
		},
	}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	result, err := svc.ListPerspectives(context.Background(), domain.PerspectiveListParams{})

	require.NoError(t, err)
	assert.Equal(t, 2, len(result.Items))
}

func TestPerspectiveList_InvalidFirst(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	first := 0
	result, err := svc.ListPerspectives(context.Background(), domain.PerspectiveListParams{First: &first})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestPerspectiveList_FirstTooLarge(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)
	first := 101
	result, err := svc.ListPerspectives(context.Background(), domain.PerspectiveListParams{First: &first})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

// --- NewPerspectiveService Tests ---

func TestNewPerspectiveService(t *testing.T) {
	perspectiveRepo := &mockPerspectiveRepository{}
	userRepo := &mockUserRepoForPerspective{}

	svc := services.NewPerspectiveService(perspectiveRepo, userRepo)

	assert.NotNil(t, svc)
}

// --- ValidateRating Tests ---

func TestValidateRating_Valid(t *testing.T) {
	testCases := []int{0, 5000, 10000}

	for _, rating := range testCases {
		t.Run(fmt.Sprintf("rating_%d", rating), func(t *testing.T) {
			assert.True(t, domain.ValidateRating(&rating))
		})
	}
}

func TestValidateRating_Invalid(t *testing.T) {
	testCases := []int{-1, 10001, -100}

	for _, rating := range testCases {
		t.Run(fmt.Sprintf("rating_%d", rating), func(t *testing.T) {
			assert.False(t, domain.ValidateRating(&rating))
		})
	}
}

func TestValidateRating_Nil(t *testing.T) {
	assert.True(t, domain.ValidateRating(nil))
}
