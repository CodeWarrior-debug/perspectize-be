package youtube

import (
	"context"
	"sync"
	"time"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/services"
)

// cacheEntry holds one cached YouTube API response and when it expires.
type cacheEntry struct {
	metadata  *services.VideoMetadata
	expiresAt time.Time
}

// CachingClient wraps a YouTubeClient with an in-memory, TTL-based cache
// keyed by video ID. It exists to avoid burning YouTube API quota on repeat
// lookups of the same video (re-adding a video someone else already added,
// refreshing metadata, etc).
//
// The cache is in-memory and per-process only — not Postgres-backed, not
// shared across replicas, and reset on every deploy/restart. That's an
// accepted tradeoff; see FEATURE_BACKLOG.md's YouTube Search Proxy entry for
// the fuller quota-sharing picture (search.list is unaffected — this only
// caches GetVideoMetadata).
type CachingClient struct {
	inner services.YouTubeClient
	ttl   time.Duration

	mu    sync.Mutex
	cache map[string]cacheEntry
}

// NewCachingClient wraps inner with a TTL cache. A ttl of zero or less
// disables caching — every call passes straight through to inner.
func NewCachingClient(inner services.YouTubeClient, ttl time.Duration) *CachingClient {
	return &CachingClient{
		inner: inner,
		ttl:   ttl,
		cache: make(map[string]cacheEntry),
	}
}

// GetVideoMetadata returns cached metadata for videoID if present and not
// expired; otherwise it fetches from the wrapped client and caches the result.
func (c *CachingClient) GetVideoMetadata(ctx context.Context, videoID string) (*services.VideoMetadata, error) {
	if c.ttl <= 0 {
		return c.inner.GetVideoMetadata(ctx, videoID)
	}

	c.mu.Lock()
	entry, ok := c.cache[videoID]
	c.mu.Unlock()

	if ok && time.Now().Before(entry.expiresAt) {
		return entry.metadata, nil
	}

	metadata, err := c.inner.GetVideoMetadata(ctx, videoID)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cache[videoID] = cacheEntry{
		metadata:  metadata,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()

	return metadata, nil
}

// ExtractVideoID delegates to the wrapped client — parsing a URL isn't an
// API call, so there's nothing to cache.
func (c *CachingClient) ExtractVideoID(url string) (string, error) {
	return c.inner.ExtractVideoID(url)
}
