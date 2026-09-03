package youtube

import (
	"context"
	"log/slog"
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
//
// Every outcome is logged at Info level with a stable "event" field
// (cache_disabled/cache_hit/cache_miss/cache_store) precisely so this can be
// verified directly in logs, not just inferred from a drop in quota errors —
// e.g. `grep 'youtube cache' server.log` or the equivalent query in your log
// viewer (Sevalla's log tab in production; this app isn't on Kubernetes, so
// there's no kubectl to exec into — logs are the direct signal here).
func (c *CachingClient) GetVideoMetadata(ctx context.Context, videoID string) (*services.VideoMetadata, error) {
	if c.ttl <= 0 {
		slog.Info("youtube cache", "event", "cache_disabled", "videoID", videoID)
		return c.inner.GetVideoMetadata(ctx, videoID)
	}

	c.mu.Lock()
	entry, ok := c.cache[videoID]
	c.mu.Unlock()

	if ok && time.Now().Before(entry.expiresAt) {
		slog.Info("youtube cache", "event", "cache_hit", "videoID", videoID, "expiresIn", time.Until(entry.expiresAt).String())
		return entry.metadata, nil
	}

	slog.Info("youtube cache", "event", "cache_miss", "videoID", videoID)

	metadata, err := c.inner.GetVideoMetadata(ctx, videoID)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cache[videoID] = cacheEntry{
		metadata:  metadata,
		expiresAt: time.Now().Add(c.ttl),
	}
	size := len(c.cache)
	c.mu.Unlock()

	slog.Info("youtube cache", "event", "cache_store", "videoID", videoID, "ttl", c.ttl.String(), "cacheSize", size)

	return metadata, nil
}

// ExtractVideoID delegates to the wrapped client — parsing a URL isn't an
// API call, so there's nothing to cache.
func (c *CachingClient) ExtractVideoID(url string) (string, error) {
	return c.inner.ExtractVideoID(url)
}
