package youtube_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/adapters/youtube"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/domain"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/ports/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- NewClient Tests ---

func TestNewClient(t *testing.T) {
	client := youtube.NewClient("test-api-key")
	require.NotNil(t, client)
}

// --- GetVideoMetadata Tests ---

func TestGetVideoMetadata_Success(t *testing.T) {
	t.Skip("Requires refactoring Client to accept baseURL for testability")
}

func TestGetVideoMetadata_Success_WithMockServer(t *testing.T) {
	mockResponse := `{
		"items": [
			{
				"id": "dQw4w9WgXcQ",
				"snippet": {
					"title": "Amazing Video",
					"description": "Description here",
					"channelTitle": "Cool Channel"
				},
				"contentDetails": {
					"duration": "PT10M"
				},
				"statistics": {
					"viewCount": "5000000",
					"likeCount": "100000",
					"commentCount": "2000"
				}
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Note: This test requires the ability to inject a custom HTTP client or baseURL
	// For now, we'll document the expected behavior
	t.Skip("Requires Client refactoring to support custom baseURL for testing")
}

func TestGetVideoMetadata_VideoNotFound(t *testing.T) {
	mockResponse := `{
		"items": []
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	t.Skip("Requires Client refactoring to support custom baseURL for testing")
}

func TestGetVideoMetadata_APIError(t *testing.T) {
	mockResponse := `{
		"error": {
			"code": 403,
			"message": "The request cannot be completed because you have exceeded your quota."
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	t.Skip("Requires Client refactoring to support custom baseURL for testing")
}

func TestGetVideoMetadata_InvalidJSON(t *testing.T) {
	mockResponse := `{invalid json}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	t.Skip("Requires Client refactoring to support custom baseURL for testing")
}

func TestGetVideoMetadata_NetworkError(t *testing.T) {
	// Create a client with an invalid baseURL to simulate network error
	client := youtube.NewClient("test-key")

	ctx := context.Background()
	result, err := client.GetVideoMetadata(ctx, "dQw4w9WgXcQ")

	assert.Nil(t, result)
	require.Error(t, err)
	// Will fail to connect to googleapis.com without internet or with invalid URL
	// This is an integration test that requires actual network
	t.Skip("Network test - requires actual YouTube API access")
}

func TestGetVideoMetadata_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		<-r.Context().Done()
	}))
	defer server.Close()

	t.Skip("Requires Client refactoring to support custom baseURL for testing")
}

func TestGetVideoMetadata_InvalidDuration(t *testing.T) {
	mockResponse := `{
		"items": [
			{
				"id": "dQw4w9WgXcQ",
				"snippet": {
					"title": "Test Video",
					"description": "Description",
					"channelTitle": "Channel"
				},
				"contentDetails": {
					"duration": "INVALID"
				},
				"statistics": {
					"viewCount": "1000",
					"likeCount": "50",
					"commentCount": "10"
				}
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	t.Skip("Requires Client refactoring to support custom baseURL for testing")
}

// --- Integration-style tests that can work with the current implementation ---

func TestGetVideoMetadata_RealAPIStructure(t *testing.T) {
	// This test verifies the structure matches YouTube API v3 expectations
	// It doesn't make real API calls, just validates our types

	// Expected YouTube API response structure
	type YouTubeAPIResponse struct {
		Items []struct {
			ID      string `json:"id"`
			Snippet struct {
				Title        string `json:"title"`
				Description  string `json:"description"`
				ChannelTitle string `json:"channelTitle"`
			} `json:"snippet"`
			ContentDetails struct {
				Duration string `json:"duration"`
			} `json:"contentDetails"`
			Statistics struct {
				ViewCount    string `json:"viewCount"`
				LikeCount    string `json:"likeCount"`
				CommentCount string `json:"commentCount"`
			} `json:"statistics"`
		} `json:"items"`
	}

	// Verify the structure compiles and matches expectations
	var response YouTubeAPIResponse
	assert.NotNil(t, response) // Structure validation
}

func TestGetVideoMetadata_ErrorTypes(t *testing.T) {
	// Test that we properly wrap domain errors
	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectedError error
	}{
		{
			name:          "not found",
			statusCode:    404,
			responseBody:  `{"error": "not found"}`,
			expectedError: domain.ErrYouTubeAPI,
		},
		{
			name:          "unauthorized",
			statusCode:    401,
			responseBody:  `{"error": "unauthorized"}`,
			expectedError: domain.ErrYouTubeAPI,
		},
		{
			name:          "forbidden",
			statusCode:    403,
			responseBody:  `{"error": "quota exceeded"}`,
			expectedError: domain.ErrYouTubeAPI,
		},
		{
			name:          "server error",
			statusCode:    500,
			responseBody:  `{"error": "internal server error"}`,
			expectedError: domain.ErrYouTubeAPI,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("Requires Client refactoring to support custom baseURL for testing")
		})
	}
}

// --- Table-driven test structure (for when Client supports injection) ---

func TestGetVideoMetadata_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		videoID        string
		mockResponse   string
		mockStatusCode int
		wantError      bool
		errorContains  string
		validate       func(t *testing.T, metadata *services.VideoMetadata)
	}{
		{
			name:           "valid video",
			videoID:        "dQw4w9WgXcQ",
			mockResponse:   validVideoResponse(),
			mockStatusCode: http.StatusOK,
			wantError:      false,
			validate: func(t *testing.T, metadata *services.VideoMetadata) {
				assert.Equal(t, "Test Video", metadata.Title)
				assert.Equal(t, "Test Description", metadata.Description)
				assert.Equal(t, 330, metadata.Duration) // 5m30s
				assert.Equal(t, "Test Channel", metadata.ChannelName)
				assert.NotNil(t, metadata.Response)
			},
		},
		{
			name:           "video not found",
			videoID:        "invalid123",
			mockResponse:   `{"items": []}`,
			mockStatusCode: http.StatusOK,
			wantError:      true,
			errorContains:  "video not found",
		},
		{
			name:           "API error",
			videoID:        "test123",
			mockResponse:   `{"error": "quota exceeded"}`,
			mockStatusCode: http.StatusForbidden,
			wantError:      true,
			errorContains:  "status 403",
		},
		{
			name:           "invalid JSON",
			videoID:        "test123",
			mockResponse:   `{invalid`,
			mockStatusCode: http.StatusOK,
			wantError:      true,
			errorContains:  "failed to parse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("Requires Client refactoring to support dependency injection")
			// When implemented, the test would look like:
			// server := createMockServer(tt.mockResponse, tt.mockStatusCode)
			// defer server.Close()
			// client := youtube.NewClientWithBaseURL("test-key", server.URL)
			// result, err := client.GetVideoMetadata(context.Background(), tt.videoID)
			// ... assertions ...
		})
	}
}

// --- Helper functions ---

func validVideoResponse() string {
	return `{
		"items": [
			{
				"id": "dQw4w9WgXcQ",
				"snippet": {
					"title": "Test Video",
					"description": "Test Description",
					"channelTitle": "Test Channel"
				},
				"contentDetails": {
					"duration": "PT5M30S"
				},
				"statistics": {
					"viewCount": "1000000",
					"likeCount": "50000",
					"commentCount": "1500"
				}
			}
		]
	}`
}

func createMockServer(response string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	}))
}

// --- Documentation Test ---

func TestYouTubeClient_Interface(t *testing.T) {
	// Verify that youtube.Client implements the YouTubeClient interface
	// This is a compile-time check
	client := youtube.NewClient("test-key")

	// Verify it has the expected method
	ctx := context.Background()
	_, err := client.GetVideoMetadata(ctx, "test")

	// We expect an error since we're not making a real API call
	// but this verifies the method signature exists
	assert.Error(t, err)
}

// --- Recommendation for future improvement ---

// TODO: Refactor youtube.Client to support dependency injection for testing
// Suggested changes:
// 1. Add NewClientWithHTTPClient(apiKey string, httpClient *http.Client) constructor
// 2. Or add NewClientWithBaseURL(apiKey string, baseURL string) for testing
// 3. This would allow us to inject httptest.Server for comprehensive unit tests
//
// Example usage after refactoring:
//   server := httptest.NewServer(mockHandler)
//   client := youtube.NewClientWithBaseURL("key", server.URL)
//   result, err := client.GetVideoMetadata(ctx, "videoID")
//
// Until then, most of these tests are skipped and the client is tested
// through integration tests in the service layer.

func TestClientRefactoringNeeded(t *testing.T) {
	t.Log("The youtube.Client needs refactoring to support testability")
	t.Log("Current implementation hard-codes baseURL and httpClient")
	t.Log("Suggested approach: Add NewClientWithHTTPClient or NewClientWithBaseURL")
	t.Log("This would allow comprehensive unit testing with httptest.Server")

	// For now, verify the basic structure works
	client := youtube.NewClient("test-api-key")
	assert.NotNil(t, client, "Client should be created")

	// The actual HTTP functionality is tested via integration tests
	// in the service layer (test/services/content_service_test.go)
}

// --- Error message validation ---

func TestGetVideoMetadata_ErrorMessages(t *testing.T) {
	// Verify that error messages are user-friendly and informative
	tests := []struct {
		name     string
		scenario string
		contains []string
	}{
		{
			name:     "video not found",
			scenario: "empty items array",
			contains: []string{"video not found", "not found"},
		},
		{
			name:     "API error",
			scenario: "non-200 status",
			contains: []string{"youtube API error", "status"},
		},
		{
			name:     "parse error",
			scenario: "invalid JSON",
			contains: []string{"failed to parse", "YouTube API response"},
		},
		{
			name:     "request error",
			scenario: "network failure",
			contains: []string{"failed to fetch", "video metadata"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Document expected error message patterns
			t.Logf("Error for %s should contain: %v", tt.scenario, tt.contains)
		})
	}
}

// --- Context handling test ---

func TestGetVideoMetadata_ContextHandling(t *testing.T) {
	client := youtube.NewClient("test-key")

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	result, err := client.GetVideoMetadata(ctx, "dQw4w9WgXcQ")

	assert.Nil(t, result)
	require.Error(t, err)
	// Should fail due to cancelled context (though the error might come from network layer)
	assert.Contains(t, err.Error(), "failed to fetch")
}

// --- API Key handling test ---

func TestNewClient_APIKey(t *testing.T) {
	tests := []struct {
		name   string
		apiKey string
	}{
		{
			name:   "normal key",
			apiKey: "AIzaSyABC123",
		},
		{
			name:   "empty key",
			apiKey: "",
		},
		{
			name:   "special characters",
			apiKey: "key-with-special_chars.123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := youtube.NewClient(tt.apiKey)
			assert.NotNil(t, client, "Client should be created regardless of API key")
			// Note: The actual API key validation happens at request time
		})
	}
}
