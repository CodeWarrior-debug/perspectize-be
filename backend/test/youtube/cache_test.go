package youtube_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/youtube"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeYouTubeClient counts calls so tests can assert cache hits vs. misses
// without hitting the network.
type fakeYouTubeClient struct {
	calls int
}

func (f *fakeYouTubeClient) GetVideoMetadata(ctx context.Context, videoID string) (*services.VideoMetadata, error) {
	f.calls++
	return &services.VideoMetadata{
		Title:    "Video " + videoID,
		Duration: f.calls, // changes per call so tests can detect a re-fetch
		Response: json.RawMessage(`{}`),
	}, nil
}

func (f *fakeYouTubeClient) ExtractVideoID(url string) (string, error) {
	return url, nil
}

func TestCachingClient_CachesWithinTTL(t *testing.T) {
	inner := &fakeYouTubeClient{}
	client := youtube.NewCachingClient(inner, time.Minute)

	first, err := client.GetVideoMetadata(context.Background(), "abc123")
	require.NoError(t, err)

	second, err := client.GetVideoMetadata(context.Background(), "abc123")
	require.NoError(t, err)

	assert.Equal(t, 1, inner.calls, "second lookup should be served from cache, not the wrapped client")
	assert.Equal(t, first.Duration, second.Duration)
}

func TestCachingClient_RefetchesAfterExpiry(t *testing.T) {
	inner := &fakeYouTubeClient{}
	client := youtube.NewCachingClient(inner, time.Millisecond)

	_, err := client.GetVideoMetadata(context.Background(), "abc123")
	require.NoError(t, err)

	time.Sleep(5 * time.Millisecond)

	_, err = client.GetVideoMetadata(context.Background(), "abc123")
	require.NoError(t, err)

	assert.Equal(t, 2, inner.calls, "expired entry should trigger a fresh fetch")
}

func TestCachingClient_ZeroTTLDisablesCaching(t *testing.T) {
	inner := &fakeYouTubeClient{}
	client := youtube.NewCachingClient(inner, 0)

	_, err := client.GetVideoMetadata(context.Background(), "abc123")
	require.NoError(t, err)
	_, err = client.GetVideoMetadata(context.Background(), "abc123")
	require.NoError(t, err)

	assert.Equal(t, 2, inner.calls, "ttl<=0 should pass every call through")
}

func TestCachingClient_DifferentVideoIDsCachedSeparately(t *testing.T) {
	inner := &fakeYouTubeClient{}
	client := youtube.NewCachingClient(inner, time.Minute)

	_, err := client.GetVideoMetadata(context.Background(), "video1")
	require.NoError(t, err)
	_, err = client.GetVideoMetadata(context.Background(), "video2")
	require.NoError(t, err)

	assert.Equal(t, 2, inner.calls, "distinct video IDs should each cause a fetch")
}

func TestCachingClient_ExtractVideoIDDelegates(t *testing.T) {
	inner := &fakeYouTubeClient{}
	client := youtube.NewCachingClient(inner, time.Minute)

	id, err := client.ExtractVideoID("https://youtu.be/xyz")
	require.NoError(t, err)
	assert.Equal(t, "https://youtu.be/xyz", id)
}
