package services

import (
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	max    int
	window time.Duration
	NowFn  func() time.Time

	mu   sync.Mutex
	hits map[string][]time.Time
}

func NewSlidingWindowLimiter(max int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		max: max, window: window, NowFn: time.Now,
		hits: make(map[string][]time.Time),
	}
}

func (l *SlidingWindowLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.NowFn()
	cutoff := now.Add(-l.window)
	kept := l.hits[key][:0]
	for _, ts := range l.hits[key] {
		if ts.After(cutoff) {
			kept = append(kept, ts)
		}
	}
	if len(kept) >= l.max {
		l.hits[key] = kept
		return false
	}
	l.hits[key] = append(kept, now)
	return true
}
