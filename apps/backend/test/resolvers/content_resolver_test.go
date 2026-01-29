package resolvers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/CodeWarrior-debug/perspectize-be/apps/backend/internal/adapters/graphql/generated"
	"github.com/CodeWarrior-debug/perspectize-be/apps/backend/internal/adapters/graphql/resolvers"
	"github.com/CodeWarrior-debug/perspectize-be/apps/backend/internal/core/domain"
	portservices "github.com/CodeWarrior-debug/perspectize-be/apps/backend/internal/core/ports/services"
	"github.com/CodeWarrior-debug/perspectize-be/apps/backend/internal/core/services"
)

// mockContentRepository implements repositories.ContentRepository for testing
type mockContentRepository struct {
	createFn   func(ctx context.Context, content *domain.Content) (*domain.Content, error)
	getByIDFn  func(ctx context.Context, id int) (*domain.Content, error)
	getByURLFn func(ctx context.Context, url string) (*domain.Content, error)
	listFn     func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error)
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

// mockYouTubeClient implements services.YouTubeClient for testing
type mockYouTubeClient struct {
	getVideoMetadataFn func(ctx context.Context, videoID string) (*portservices.VideoMetadata, error)
}

func (m *mockYouTubeClient) GetVideoMetadata(ctx context.Context, videoID string) (*portservices.VideoMetadata, error) {
	if m.getVideoMetadataFn != nil {
		return m.getVideoMetadataFn(ctx, videoID)
	}
	return nil, fmt.Errorf("not implemented")
}

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

// graphqlResponse represents a generic GraphQL JSON response
type graphqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// setupTestServer creates a test GraphQL server with the given mock dependencies
func setupTestServer(repo *mockContentRepository, ytClient *mockYouTubeClient) *httptest.Server {
	userRepo := &mockUserRepository{}
	contentService := services.NewContentService(repo, ytClient)
	userService := services.NewUserService(userRepo)
	perspectiveService := services.NewPerspectiveService(&mockPerspectiveRepository{}, userRepo)
	resolver := resolvers.NewResolver(contentService, userService, perspectiveService)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	return httptest.NewServer(srv)
}

// executeGraphQL sends a GraphQL query to the test server and returns the response
func executeGraphQL(t *testing.T, server *httptest.Server, query string) graphqlResponse {
	t.Helper()

	body := fmt.Sprintf(`{"query": %s}`, jsonString(query))
	resp, err := http.Post(server.URL, "application/json", strings.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	var result graphqlResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	return result
}

// jsonString encodes a string as a JSON string value
func jsonString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// --- Content Query Tests ---

func TestContentQuery_Success(t *testing.T) {
	url := "https://youtube.com/watch?v=abc123"
	repo := &mockContentRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.Content, error) {
			return &domain.Content{
				ID:          1,
				Name:        "Test Video",
				URL:         &url,
				ContentType: domain.ContentTypeYouTube,
			}, nil
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ contentByID(id: "1") { id name contentType url } }`)

	assert.Empty(t, result.Errors)

	var data struct {
		ContentByID struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ContentType string `json:"contentType"`
			URL         string `json:"url"`
		} `json:"contentByID"`
	}
	err := json.Unmarshal(result.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "1", data.ContentByID.ID)
	assert.Equal(t, "Test Video", data.ContentByID.Name)
	assert.Equal(t, "youtube", data.ContentByID.ContentType)
	assert.Equal(t, url, data.ContentByID.URL)
}

func TestContentQuery_NotFound_ReturnsError(t *testing.T) {
	repo := &mockContentRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.Content, error) {
			return nil, domain.ErrNotFound
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ contentByID(id: "999") { id name } }`)

	// Issue #18 fix: should return an error, not a silent null
	require.NotEmpty(t, result.Errors, "Expected an error when content is not found")
	assert.Contains(t, result.Errors[0].Message, "content not found")
}

func TestContentQuery_InvalidID_NonNumeric(t *testing.T) {
	repo := &mockContentRepository{}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ contentByID(id: "abc") { id name } }`)

	require.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Message, "invalid content ID")
}

func TestContentQuery_InvalidID_Zero(t *testing.T) {
	repo := &mockContentRepository{}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ contentByID(id: "0") { id name } }`)

	require.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Message, "invalid content ID")
}

func TestContentQuery_InvalidID_Negative(t *testing.T) {
	repo := &mockContentRepository{}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ contentByID(id: "-1") { id name } }`)

	require.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Message, "invalid content ID")
}

func TestContentQuery_DatabaseError(t *testing.T) {
	repo := &mockContentRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.Content, error) {
			return nil, fmt.Errorf("connection refused")
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ contentByID(id: "1") { id name } }`)

	require.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Message, "failed to get content")
}

// --- CreateContentFromYouTube Mutation Tests ---

func TestCreateContentFromYouTube_Success(t *testing.T) {
	metadata := &portservices.VideoMetadata{
		Title:       "Amazing Video",
		Description: "Description",
		Duration:    600,
		ChannelName: "Channel",
		Response:    json.RawMessage(`{"items":[]}`),
	}

	repo := &mockContentRepository{
		getByURLFn: func(ctx context.Context, url string) (*domain.Content, error) {
			return nil, domain.ErrNotFound
		},
		createFn: func(ctx context.Context, content *domain.Content) (*domain.Content, error) {
			content.ID = 42
			return content, nil
		},
	}

	ytClient := &mockYouTubeClient{
		getVideoMetadataFn: func(ctx context.Context, videoID string) (*portservices.VideoMetadata, error) {
			return metadata, nil
		},
	}

	server := setupTestServer(repo, ytClient)
	defer server.Close()

	result := executeGraphQL(t, server, `mutation { createContentFromYouTube(input: { url: "https://www.youtube.com/watch?v=dQw4w9WgXcQ" }) { id name contentType } }`)

	assert.Empty(t, result.Errors)

	var data struct {
		CreateContentFromYouTube struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ContentType string `json:"contentType"`
		} `json:"createContentFromYouTube"`
	}
	err := json.Unmarshal(result.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "42", data.CreateContentFromYouTube.ID)
	assert.Equal(t, "Amazing Video", data.CreateContentFromYouTube.Name)
	assert.Equal(t, "youtube", data.CreateContentFromYouTube.ContentType)
}

func TestCreateContentFromYouTube_AlreadyExists(t *testing.T) {
	existingURL := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	repo := &mockContentRepository{
		getByURLFn: func(ctx context.Context, url string) (*domain.Content, error) {
			return &domain.Content{ID: 1, URL: &existingURL}, nil
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `mutation { createContentFromYouTube(input: { url: "https://www.youtube.com/watch?v=dQw4w9WgXcQ" }) { id } }`)

	require.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Message, "content already exists")
}

func TestCreateContentFromYouTube_InvalidURL(t *testing.T) {
	repo := &mockContentRepository{
		getByURLFn: func(ctx context.Context, url string) (*domain.Content, error) {
			return nil, domain.ErrNotFound
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `mutation { createContentFromYouTube(input: { url: "not-a-youtube-url" }) { id } }`)

	require.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Message, "invalid YouTube URL")
}

// --- Paginated Content Query Tests ---

func TestPaginatedContentQuery_DefaultPagination(t *testing.T) {
	repo := &mockContentRepository{
		listFn: func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
			// Verify default values - GraphQL passes the schema default (10), not nil
			require.NotNil(t, params.First)
			assert.Equal(t, 10, *params.First)
			assert.Equal(t, domain.ContentSortByCreatedAt, params.SortBy)
			assert.Equal(t, domain.SortOrderDesc, params.SortOrder)
			assert.False(t, params.IncludeTotalCount)

			url := "https://youtube.com/watch?v=abc123"
			return &domain.PaginatedContent{
				Items: []*domain.Content{
					{ID: 1, Name: "Video 1", URL: &url, ContentType: domain.ContentTypeYouTube},
					{ID: 2, Name: "Video 2", URL: &url, ContentType: domain.ContentTypeYouTube},
				},
				HasNext: false,
				HasPrev: false,
			}, nil
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ content { items { id name } pageInfo { hasNextPage hasPreviousPage } } }`)

	assert.Empty(t, result.Errors)

	var data struct {
		Content struct {
			Items []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"items"`
			PageInfo struct {
				HasNextPage     bool `json:"hasNextPage"`
				HasPreviousPage bool `json:"hasPreviousPage"`
			} `json:"pageInfo"`
		} `json:"content"`
	}
	err := json.Unmarshal(result.Data, &data)
	require.NoError(t, err)

	assert.Len(t, data.Content.Items, 2)
	assert.Equal(t, "1", data.Content.Items[0].ID)
	assert.Equal(t, "Video 1", data.Content.Items[0].Name)
	assert.False(t, data.Content.PageInfo.HasNextPage)
	assert.False(t, data.Content.PageInfo.HasPreviousPage)
}

func TestPaginatedContentQuery_WithFirstParameter(t *testing.T) {
	repo := &mockContentRepository{
		listFn: func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
			require.NotNil(t, params.First)
			assert.Equal(t, 5, *params.First)

			url := "https://youtube.com/watch?v=abc123"
			items := make([]*domain.Content, 5)
			for i := 0; i < 5; i++ {
				items[i] = &domain.Content{ID: i + 1, Name: fmt.Sprintf("Video %d", i+1), URL: &url, ContentType: domain.ContentTypeYouTube}
			}
			endCursor := "cursor123"
			return &domain.PaginatedContent{
				Items:     items,
				HasNext:   true,
				HasPrev:   false,
				EndCursor: &endCursor,
			}, nil
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ content(first: 5) { items { id } pageInfo { hasNextPage endCursor } } }`)

	assert.Empty(t, result.Errors)

	var data struct {
		Content struct {
			Items []struct {
				ID string `json:"id"`
			} `json:"items"`
			PageInfo struct {
				HasNextPage bool    `json:"hasNextPage"`
				EndCursor   *string `json:"endCursor"`
			} `json:"pageInfo"`
		} `json:"content"`
	}
	err := json.Unmarshal(result.Data, &data)
	require.NoError(t, err)

	assert.Len(t, data.Content.Items, 5)
	assert.True(t, data.Content.PageInfo.HasNextPage)
	assert.NotNil(t, data.Content.PageInfo.EndCursor)
}

func TestPaginatedContentQuery_WithTotalCount(t *testing.T) {
	repo := &mockContentRepository{
		listFn: func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
			assert.True(t, params.IncludeTotalCount)

			totalCount := 42
			return &domain.PaginatedContent{
				Items:      []*domain.Content{},
				HasNext:    false,
				HasPrev:    false,
				TotalCount: &totalCount,
			}, nil
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ content(includeTotalCount: true) { totalCount items { id } } }`)

	assert.Empty(t, result.Errors)

	var data struct {
		Content struct {
			TotalCount *int `json:"totalCount"`
			Items      []struct {
				ID string `json:"id"`
			} `json:"items"`
		} `json:"content"`
	}
	err := json.Unmarshal(result.Data, &data)
	require.NoError(t, err)

	require.NotNil(t, data.Content.TotalCount)
	assert.Equal(t, 42, *data.Content.TotalCount)
}

func TestPaginatedContentQuery_WithSorting(t *testing.T) {
	repo := &mockContentRepository{
		listFn: func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
			assert.Equal(t, domain.ContentSortByName, params.SortBy)
			assert.Equal(t, domain.SortOrderAsc, params.SortOrder)

			return &domain.PaginatedContent{
				Items:   []*domain.Content{},
				HasNext: false,
				HasPrev: false,
			}, nil
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ content(sortBy: NAME, sortOrder: ASC) { items { id } } }`)

	assert.Empty(t, result.Errors)
}

func TestPaginatedContentQuery_WithAfterCursor(t *testing.T) {
	repo := &mockContentRepository{
		listFn: func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
			require.NotNil(t, params.After)
			assert.Equal(t, "someCursor123", *params.After)

			url := "https://youtube.com/watch?v=abc123"
			return &domain.PaginatedContent{
				Items: []*domain.Content{
					{ID: 11, Name: "Video 11", URL: &url, ContentType: domain.ContentTypeYouTube},
				},
				HasNext: false,
				HasPrev: true,
			}, nil
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ content(after: "someCursor123") { items { id } pageInfo { hasPreviousPage } } }`)

	assert.Empty(t, result.Errors)

	var data struct {
		Content struct {
			Items []struct {
				ID string `json:"id"`
			} `json:"items"`
			PageInfo struct {
				HasPreviousPage bool `json:"hasPreviousPage"`
			} `json:"pageInfo"`
		} `json:"content"`
	}
	err := json.Unmarshal(result.Data, &data)
	require.NoError(t, err)

	assert.Len(t, data.Content.Items, 1)
	assert.Equal(t, "11", data.Content.Items[0].ID)
	assert.True(t, data.Content.PageInfo.HasPreviousPage)
}

func TestPaginatedContentQuery_InvalidFirstParameter(t *testing.T) {
	repo := &mockContentRepository{}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	// first: 0 is invalid (must be 1-100)
	result := executeGraphQL(t, server, `{ content(first: 0) { items { id } } }`)

	require.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Message, "invalid pagination parameters")
}

func TestPaginatedContentQuery_FirstTooLarge(t *testing.T) {
	repo := &mockContentRepository{}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	// first: 101 exceeds max of 100
	result := executeGraphQL(t, server, `{ content(first: 101) { items { id } } }`)

	require.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Message, "invalid pagination parameters")
}

func TestPaginatedContentQuery_RepositoryError(t *testing.T) {
	repo := &mockContentRepository{
		listFn: func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
			return nil, fmt.Errorf("database connection failed")
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ content { items { id } } }`)

	require.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Message, "failed to list content")
}

func TestPaginatedContentQuery_WithContentTypeFilter(t *testing.T) {
	repo := &mockContentRepository{
		listFn: func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
			require.NotNil(t, params.Filter)
			require.NotNil(t, params.Filter.ContentType)
			assert.Equal(t, domain.ContentTypeYouTube, *params.Filter.ContentType)

			url := "https://youtube.com/watch?v=abc123"
			return &domain.PaginatedContent{
				Items: []*domain.Content{
					{ID: 1, Name: "YouTube Video", URL: &url, ContentType: domain.ContentTypeYouTube},
				},
				HasNext: false,
				HasPrev: false,
			}, nil
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ content(filter: { contentType: YOUTUBE }) { items { id name contentType } } }`)

	assert.Empty(t, result.Errors)

	var data struct {
		Content struct {
			Items []struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				ContentType string `json:"contentType"`
			} `json:"items"`
		} `json:"content"`
	}
	err := json.Unmarshal(result.Data, &data)
	require.NoError(t, err)

	assert.Len(t, data.Content.Items, 1)
	assert.Equal(t, "YouTube Video", data.Content.Items[0].Name)
	assert.Equal(t, "youtube", data.Content.Items[0].ContentType)
}

func TestPaginatedContentQuery_WithFilterAndTotalCount(t *testing.T) {
	repo := &mockContentRepository{
		listFn: func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
			require.NotNil(t, params.Filter)
			require.NotNil(t, params.Filter.ContentType)
			assert.True(t, params.IncludeTotalCount)

			totalCount := 5
			return &domain.PaginatedContent{
				Items:      []*domain.Content{},
				HasNext:    false,
				HasPrev:    false,
				TotalCount: &totalCount,
			}, nil
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ content(filter: { contentType: YOUTUBE }, includeTotalCount: true) { totalCount items { id } } }`)

	assert.Empty(t, result.Errors)

	var data struct {
		Content struct {
			TotalCount *int `json:"totalCount"`
			Items      []struct {
				ID string `json:"id"`
			} `json:"items"`
		} `json:"content"`
	}
	err := json.Unmarshal(result.Data, &data)
	require.NoError(t, err)

	require.NotNil(t, data.Content.TotalCount)
	assert.Equal(t, 5, *data.Content.TotalCount)
}

func TestPaginatedContentQuery_NoFilter(t *testing.T) {
	repo := &mockContentRepository{
		listFn: func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
			assert.Nil(t, params.Filter)

			return &domain.PaginatedContent{
				Items:   []*domain.Content{},
				HasNext: false,
				HasPrev: false,
			}, nil
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ content { items { id } } }`)

	assert.Empty(t, result.Errors)
}

func TestPaginatedContentQuery_WithMinLengthFilter(t *testing.T) {
	repo := &mockContentRepository{
		listFn: func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
			require.NotNil(t, params.Filter)
			require.NotNil(t, params.Filter.MinLengthSeconds)
			assert.Equal(t, 300, *params.Filter.MinLengthSeconds)
			assert.Nil(t, params.Filter.MaxLengthSeconds)

			url := "https://youtube.com/watch?v=abc123"
			length := 600
			return &domain.PaginatedContent{
				Items: []*domain.Content{
					{ID: 1, Name: "Long Video", URL: &url, ContentType: domain.ContentTypeYouTube, Length: &length},
				},
				HasNext: false,
				HasPrev: false,
			}, nil
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ content(filter: { minLengthSeconds: 300 }) { items { id name length } } }`)

	assert.Empty(t, result.Errors)

	var data struct {
		Content struct {
			Items []struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Length *int   `json:"length"`
			} `json:"items"`
		} `json:"content"`
	}
	err := json.Unmarshal(result.Data, &data)
	require.NoError(t, err)

	assert.Len(t, data.Content.Items, 1)
	require.NotNil(t, data.Content.Items[0].Length)
	assert.Equal(t, 600, *data.Content.Items[0].Length)
}

func TestPaginatedContentQuery_WithMaxLengthFilter(t *testing.T) {
	repo := &mockContentRepository{
		listFn: func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
			require.NotNil(t, params.Filter)
			require.NotNil(t, params.Filter.MaxLengthSeconds)
			assert.Equal(t, 180, *params.Filter.MaxLengthSeconds)
			assert.Nil(t, params.Filter.MinLengthSeconds)

			url := "https://youtube.com/watch?v=abc123"
			length := 120
			return &domain.PaginatedContent{
				Items: []*domain.Content{
					{ID: 1, Name: "Short Video", URL: &url, ContentType: domain.ContentTypeYouTube, Length: &length},
				},
				HasNext: false,
				HasPrev: false,
			}, nil
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ content(filter: { maxLengthSeconds: 180 }) { items { id name length } } }`)

	assert.Empty(t, result.Errors)

	var data struct {
		Content struct {
			Items []struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Length *int   `json:"length"`
			} `json:"items"`
		} `json:"content"`
	}
	err := json.Unmarshal(result.Data, &data)
	require.NoError(t, err)

	assert.Len(t, data.Content.Items, 1)
	require.NotNil(t, data.Content.Items[0].Length)
	assert.Equal(t, 120, *data.Content.Items[0].Length)
}

func TestPaginatedContentQuery_WithMinMaxLengthFilter(t *testing.T) {
	repo := &mockContentRepository{
		listFn: func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error) {
			require.NotNil(t, params.Filter)
			require.NotNil(t, params.Filter.MinLengthSeconds)
			require.NotNil(t, params.Filter.MaxLengthSeconds)
			assert.Equal(t, 120, *params.Filter.MinLengthSeconds)
			assert.Equal(t, 300, *params.Filter.MaxLengthSeconds)

			url := "https://youtube.com/watch?v=abc123"
			length := 200
			return &domain.PaginatedContent{
				Items: []*domain.Content{
					{ID: 1, Name: "Medium Video", URL: &url, ContentType: domain.ContentTypeYouTube, Length: &length},
				},
				HasNext: false,
				HasPrev: false,
			}, nil
		},
	}

	server := setupTestServer(repo, &mockYouTubeClient{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ content(filter: { minLengthSeconds: 120, maxLengthSeconds: 300 }) { items { id length } } }`)

	assert.Empty(t, result.Errors)

	var data struct {
		Content struct {
			Items []struct {
				ID     string `json:"id"`
				Length *int   `json:"length"`
			} `json:"items"`
		} `json:"content"`
	}
	err := json.Unmarshal(result.Data, &data)
	require.NoError(t, err)

	assert.Len(t, data.Content.Items, 1)
}

// --- NewResolver Tests ---

func TestNewResolver(t *testing.T) {
	repo := &mockContentRepository{}
	ytClient := &mockYouTubeClient{}
	userRepo := &mockUserRepository{}
	contentService := services.NewContentService(repo, ytClient)
	userService := services.NewUserService(userRepo)
	perspectiveService := services.NewPerspectiveService(&mockPerspectiveRepository{}, userRepo)

	resolver := resolvers.NewResolver(contentService, userService, perspectiveService)

	assert.NotNil(t, resolver)
	assert.Equal(t, contentService, resolver.ContentService)
	assert.Equal(t, userService, resolver.UserService)
	assert.Equal(t, perspectiveService, resolver.PerspectiveService)
}
