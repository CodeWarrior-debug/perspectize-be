package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	portservices "github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/services"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockContentRepository implements repositories.ContentRepository for testing
type mockContentRepository struct {
	createFn               func(ctx context.Context, content *domain.Content) (*domain.Content, error)
	getByIDFn              func(ctx context.Context, id int) (*domain.Content, error)
	getByURLFn             func(ctx context.Context, url string) (*domain.Content, error)
	listFn                 func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error)
	updateAttributionFn    func(ctx context.Context, contentID, newUserID int) (*domain.Content, error)
}

func (m *mockContentRepository) Create(ctx context.Context, content *domain.Content) (*domain.Content, error) {
	if m.createFn != nil {
		return m.createFn(ctx, content)
	}
	return content, nil
}

func (m *mockContentRepository) GetByID(ctx context.Context, id int) (*domain.Content, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *mockContentRepository) GetByURL(ctx context.Context, url string) (*domain.Content, error) {
	if m.getByURLFn != nil {
		return m.getByURLFn(ctx, url)
	}
	return nil, domain.ErrNotFound
}

func (m *mockContentRepository) List(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
	if m.listFn != nil {
		return m.listFn(ctx, params)
	}
	return &domain.PaginatedContent{Items: []*domain.Content{}}, nil
}

func (m *mockContentRepository) ReassignByUser(ctx context.Context, fromUserID, toUserID int) error {
	return nil
}

func (m *mockContentRepository) UpdateAttribution(ctx context.Context, contentID, newUserID int) (*domain.Content, error) {
	if m.updateAttributionFn != nil {
		return m.updateAttributionFn(ctx, contentID, newUserID)
	}
	return nil, domain.ErrNotFound
}

// mockYouTubeClient implements services.YouTubeClient for testing
type mockYouTubeClient struct {
	getVideoMetadataFn func(ctx context.Context, videoID string) (*portservices.VideoMetadata, error)
	extractVideoIDFn   func(url string) (string, error)
}

func (m *mockYouTubeClient) GetVideoMetadata(ctx context.Context, videoID string) (*portservices.VideoMetadata, error) {
	if m.getVideoMetadataFn != nil {
		return m.getVideoMetadataFn(ctx, videoID)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockYouTubeClient) ExtractVideoID(url string) (string, error) {
	if m.extractVideoIDFn != nil {
		return m.extractVideoIDFn(url)
	}
	return "", fmt.Errorf("could not extract video ID")
}

// --- GetByID Tests ---

func TestGetByID_Success(t *testing.T) {
	url := "https://youtube.com/watch?v=abc123"
	expected := &domain.Content{
		ID:          1,
		Name:        "Test Video",
		URL:         &url,
		ContentType: domain.ContentTypeYouTube,
	}

	repo := &mockContentRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.Content, error) {
			assert.Equal(t, 1, id)
			return expected, nil
		},
	}

	svc := services.NewContentService(repo, &mockYouTubeClient{})
	result, err := svc.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &mockContentRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.Content, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := services.NewContentService(repo, &mockYouTubeClient{})
	result, err := svc.GetByID(context.Background(), 999)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestGetByID_InvalidID_Zero(t *testing.T) {
	repo := &mockContentRepository{}
	svc := services.NewContentService(repo, &mockYouTubeClient{})

	result, err := svc.GetByID(context.Background(), 0)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "content id must be a positive integer")
}

func TestGetByID_InvalidID_Negative(t *testing.T) {
	repo := &mockContentRepository{}
	svc := services.NewContentService(repo, &mockYouTubeClient{})

	result, err := svc.GetByID(context.Background(), -5)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestGetByID_RepositoryError(t *testing.T) {
	repo := &mockContentRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.Content, error) {
			return nil, fmt.Errorf("database connection failed")
		},
	}

	svc := services.NewContentService(repo, &mockYouTubeClient{})
	result, err := svc.GetByID(context.Background(), 1)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get content")
}

// --- CreateFromYouTube Tests ---

func TestCreateFromYouTube_Success(t *testing.T) {
	videoURL := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	metadata := &portservices.VideoMetadata{
		Title:       "Test Video Title",
		Description: "A great video",
		Duration:    300,
		ChannelName: "Test Channel",
		Response:    json.RawMessage(`{"items":[]}`),
	}

	repo := &mockContentRepository{
		getByURLFn: func(ctx context.Context, url string) (*domain.Content, error) {
			return nil, domain.ErrNotFound // URL does not exist yet
		},
		createFn: func(ctx context.Context, content *domain.Content) (*domain.Content, error) {
			content.ID = 1
			return content, nil
		},
	}

	ytClient := &mockYouTubeClient{
		extractVideoIDFn: func(url string) (string, error) {
			return "dQw4w9WgXcQ", nil
		},
		getVideoMetadataFn: func(ctx context.Context, videoID string) (*portservices.VideoMetadata, error) {
			assert.Equal(t, "dQw4w9WgXcQ", videoID)
			return metadata, nil
		},
	}

	svc := services.NewContentService(repo, ytClient)

	result, err := svc.CreateFromYouTube(context.Background(), videoURL, 42)

	require.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Test Video Title", result.Name)
	assert.Equal(t, domain.ContentTypeYouTube, result.ContentType)
	assert.Equal(t, &videoURL, result.URL)
	assert.Equal(t, 42, result.AddedByUserID)
	require.NotNil(t, result.Length)
	assert.Equal(t, 300, *result.Length)
	require.NotNil(t, result.LengthUnits)
	assert.Equal(t, "seconds", *result.LengthUnits)
}

func TestCreateFromYouTube_AlreadyExists(t *testing.T) {
	existingURL := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	existing := &domain.Content{
		ID:   1,
		Name: "Existing Video",
		URL:  &existingURL,
	}

	repo := &mockContentRepository{
		getByURLFn: func(ctx context.Context, url string) (*domain.Content, error) {
			return existing, nil // URL already exists
		},
	}

	svc := services.NewContentService(repo, &mockYouTubeClient{})

	result, err := svc.CreateFromYouTube(context.Background(), existingURL, 1)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrAlreadyExists))
}

func TestCreateFromYouTube_InvalidURL(t *testing.T) {
	repo := &mockContentRepository{
		getByURLFn: func(ctx context.Context, url string) (*domain.Content, error) {
			return nil, domain.ErrNotFound
		},
	}

	ytClient := &mockYouTubeClient{
		extractVideoIDFn: func(url string) (string, error) {
			return "", fmt.Errorf("could not extract video ID")
		},
	}

	svc := services.NewContentService(repo, ytClient)

	result, err := svc.CreateFromYouTube(context.Background(), "not-a-valid-url", 1)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidURL))
}

func TestCreateFromYouTube_YouTubeAPIError(t *testing.T) {
	repo := &mockContentRepository{
		getByURLFn: func(ctx context.Context, url string) (*domain.Content, error) {
			return nil, domain.ErrNotFound
		},
	}

	ytClient := &mockYouTubeClient{
		extractVideoIDFn: func(url string) (string, error) {
			return "abc123", nil
		},
		getVideoMetadataFn: func(ctx context.Context, videoID string) (*portservices.VideoMetadata, error) {
			return nil, fmt.Errorf("%w: status 403", domain.ErrYouTubeAPI)
		},
	}

	svc := services.NewContentService(repo, ytClient)

	result, err := svc.CreateFromYouTube(context.Background(), "https://youtube.com/watch?v=abc123", 1)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch YouTube metadata")
}

func TestCreateFromYouTube_RepositoryCreateError(t *testing.T) {
	metadata := &portservices.VideoMetadata{
		Title:    "Video",
		Duration: 60,
		Response: json.RawMessage(`{}`),
	}

	repo := &mockContentRepository{
		getByURLFn: func(ctx context.Context, url string) (*domain.Content, error) {
			return nil, domain.ErrNotFound
		},
		createFn: func(ctx context.Context, content *domain.Content) (*domain.Content, error) {
			return nil, fmt.Errorf("database write error")
		},
	}

	ytClient := &mockYouTubeClient{
		extractVideoIDFn: func(url string) (string, error) {
			return "abc123", nil
		},
		getVideoMetadataFn: func(ctx context.Context, videoID string) (*portservices.VideoMetadata, error) {
			return metadata, nil
		},
	}

	svc := services.NewContentService(repo, ytClient)

	result, err := svc.CreateFromYouTube(context.Background(), "https://youtube.com/watch?v=abc123", 1)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save content")
}

func TestCreateFromYouTube_GetByURLUnexpectedError(t *testing.T) {
	repo := &mockContentRepository{
		getByURLFn: func(ctx context.Context, url string) (*domain.Content, error) {
			return nil, fmt.Errorf("unexpected database error")
		},
	}

	svc := services.NewContentService(repo, &mockYouTubeClient{})

	result, err := svc.CreateFromYouTube(context.Background(), "https://youtube.com/watch?v=abc123", 1)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check existing content")
}

// --- UpdateAttribution Tests ---

func TestUpdateAttribution_Success(t *testing.T) {
	url := "https://youtube.com/watch?v=abc123"
	expected := &domain.Content{
		ID:            1,
		Name:          "Test Video",
		URL:           &url,
		ContentType:   domain.ContentTypeYouTube,
		AddedByUserID: 42, // new user ID
	}

	repo := &mockContentRepository{
		updateAttributionFn: func(ctx context.Context, contentID, newUserID int) (*domain.Content, error) {
			assert.Equal(t, 1, contentID)
			assert.Equal(t, 42, newUserID)
			return expected, nil
		},
	}

	svc := services.NewContentService(repo, &mockYouTubeClient{})
	result, err := svc.UpdateAttribution(context.Background(), 1, 42)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, 42, result.AddedByUserID)
}

func TestUpdateAttribution_ContentIDZero(t *testing.T) {
	repo := &mockContentRepository{}
	svc := services.NewContentService(repo, &mockYouTubeClient{})

	result, err := svc.UpdateAttribution(context.Background(), 0, 42)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "content id must be a positive integer")
}

func TestUpdateAttribution_ContentIDNegative(t *testing.T) {
	repo := &mockContentRepository{}
	svc := services.NewContentService(repo, &mockYouTubeClient{})

	result, err := svc.UpdateAttribution(context.Background(), -5, 42)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "content id must be a positive integer")
}

func TestUpdateAttribution_NewUserIDZero(t *testing.T) {
	repo := &mockContentRepository{}
	svc := services.NewContentService(repo, &mockYouTubeClient{})

	result, err := svc.UpdateAttribution(context.Background(), 1, 0)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "user id must be a positive integer")
}

func TestUpdateAttribution_NewUserIDNegative(t *testing.T) {
	repo := &mockContentRepository{}
	svc := services.NewContentService(repo, &mockYouTubeClient{})

	result, err := svc.UpdateAttribution(context.Background(), 1, -10)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "user id must be a positive integer")
}

func TestUpdateAttribution_ContentNotFound(t *testing.T) {
	repo := &mockContentRepository{
		updateAttributionFn: func(ctx context.Context, contentID, newUserID int) (*domain.Content, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := services.NewContentService(repo, &mockYouTubeClient{})
	result, err := svc.UpdateAttribution(context.Background(), 999, 42)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestUpdateAttribution_RepositoryError(t *testing.T) {
	repo := &mockContentRepository{
		updateAttributionFn: func(ctx context.Context, contentID, newUserID int) (*domain.Content, error) {
			return nil, fmt.Errorf("database connection failed")
		},
	}

	svc := services.NewContentService(repo, &mockYouTubeClient{})
	result, err := svc.UpdateAttribution(context.Background(), 1, 42)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update content attribution")
}

func TestUpdateAttribution_BothIDsInvalid(t *testing.T) {
	repo := &mockContentRepository{}
	svc := services.NewContentService(repo, &mockYouTubeClient{})

	result, err := svc.UpdateAttribution(context.Background(), 0, 0)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	// Error should be for contentID since it's checked first
	assert.Contains(t, err.Error(), "content id must be a positive integer")
}

// --- NewContentService Tests ---

func TestNewContentService(t *testing.T) {
	repo := &mockContentRepository{}
	ytClient := &mockYouTubeClient{}

	svc := services.NewContentService(repo, ytClient)

	assert.NotNil(t, svc)
}
