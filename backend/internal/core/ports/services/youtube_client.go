package services

import (
	"context"
	"encoding/json"
)

// VideoMetadata contains extracted information from YouTube API response
type VideoMetadata struct {
	Title       string
	Description string
	Duration    int // Duration in seconds
	ChannelName string
	Response    json.RawMessage // Raw API response for storage
}

// YouTubeClient defines the contract for YouTube API interactions
type YouTubeClient interface {
	// GetVideoMetadata fetches video metadata from YouTube API
	GetVideoMetadata(ctx context.Context, videoID string) (*VideoMetadata, error)

	// ExtractVideoID extracts the video ID from a YouTube URL
	ExtractVideoID(url string) (string, error)
}
