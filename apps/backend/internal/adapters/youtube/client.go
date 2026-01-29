package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/CodeWarrior-debug/perspectize-be/apps/backend/internal/core/domain"
	"github.com/CodeWarrior-debug/perspectize-be/apps/backend/internal/core/ports/services"
)

// Client implements the YouTubeClient interface for YouTube Data API v3
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new YouTube API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		baseURL:    "https://www.googleapis.com/youtube/v3",
	}
}

// YouTubeAPIResponse represents the response from YouTube Data API
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

// GetVideoMetadata fetches video metadata from YouTube Data API
func (c *Client) GetVideoMetadata(ctx context.Context, videoID string) (*services.VideoMetadata, error) {
	endpoint := fmt.Sprintf("%s/videos?part=snippet,statistics,contentDetails&id=%s&key=%s",
		c.baseURL,
		url.QueryEscape(videoID),
		url.QueryEscape(c.apiKey),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video metadata: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d: %s", domain.ErrYouTubeAPI, resp.StatusCode, string(body))
	}

	var apiResponse YouTubeAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse YouTube API response: %w", err)
	}

	if len(apiResponse.Items) == 0 {
		return nil, fmt.Errorf("%w: video not found: %s", domain.ErrNotFound, videoID)
	}

	item := apiResponse.Items[0]

	duration, err := ParseISO8601Duration(item.ContentDetails.Duration)
	if err != nil {
		duration = 0 // Default to 0 if parsing fails
	}

	return &services.VideoMetadata{
		Title:       item.Snippet.Title,
		Description: item.Snippet.Description,
		Duration:    duration,
		ChannelName: item.Snippet.ChannelTitle,
		Response:    body,
	}, nil
}
