package realtime_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/realtime"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/repositories"
	"github.com/stretchr/testify/assert"
)

// presenceEvent is one recorded PRESENCE_CHANGED fan-out.
type presenceEvent struct {
	userID   int
	state    string
	threadID int
}

// stubNotifier records every PublishEphemeral call under a mutex. It
// structurally satisfies realtime's unexported presenceNotifier interface.
type stubNotifier struct {
	mu     sync.Mutex
	events []presenceEvent
}

func (s *stubNotifier) PublishEphemeral(_ context.Context, env domain.EventEnvelope) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, presenceEvent{userID: env.UserID, state: env.State, threadID: env.ThreadID})
	return nil
}

func (s *stubNotifier) count(state string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, e := range s.events {
		if e.state == state {
			n++
		}
	}
	return n
}

func (s *stubNotifier) withState(state string) []presenceEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []presenceEvent
	for _, e := range s.events {
		if e.state == state {
			out = append(out, e)
		}
	}
	return out
}

// presenceStubThreadRepo returns a fixed set of thread IDs from
// ListThreadsForUser; every other method is a no-op.
type presenceStubThreadRepo struct{ ids []int }

func (presenceStubThreadRepo) CreateThread(_ context.Context, _ int, _ *string, _ []int) (*domain.MessageThread, error) {
	return nil, nil
}
func (presenceStubThreadRepo) GetThread(_ context.Context, _ int) (*domain.MessageThread, error) {
	return nil, nil
}
func (presenceStubThreadRepo) FindDirectThread(_ context.Context, _, _ int) (*domain.MessageThread, error) {
	return nil, nil
}
func (s presenceStubThreadRepo) ListThreadsForUser(_ context.Context, _ int, _ int, _ *time.Time) ([]domain.MessageThread, error) {
	out := make([]domain.MessageThread, 0, len(s.ids))
	for _, id := range s.ids {
		out = append(out, domain.MessageThread{ID: id})
	}
	return out, nil
}
func (presenceStubThreadRepo) AddParticipants(_ context.Context, _ int, _ []int) error { return nil }
func (presenceStubThreadRepo) SetLeft(_ context.Context, _, _ int, _ time.Time) error  { return nil }
func (presenceStubThreadRepo) SetLastRead(_ context.Context, _, _ int, _ int64) error  { return nil }

var _ repositories.ThreadRepository = presenceStubThreadRepo{}

// newSkewableTracker returns a tracker whose clock can be advanced race-free by
// adding to skew (nanoseconds). The grace-window OFFLINE re-check calls
// tracker.IsOnline, which stays true for presenceTTL (45s) after the last
// activity regardless of refcount; advancing the clock past that window (after
// the final Disconnect has stamped its lastSeen) lets the re-check observe the
// user as offline. The real time.AfterFunc / time.Ticker timers inside
// RunPresenceSession are untouched.
func newSkewableTracker(skew *atomic.Int64) *realtime.PresenceTracker {
	tr := realtime.NewPresenceTracker()
	tr.NowFn = func() time.Time { return time.Now().Add(time.Duration(skew.Load())) }
	return tr
}

var testCfg = realtime.PresenceConfig{
	OfflineGrace:      40 * time.Millisecond,
	HeartbeatInterval: 10 * time.Millisecond,
}

func TestPresenceSession_OnlineThenOfflineAfterGrace(t *testing.T) {
	var skew atomic.Int64
	tracker := newSkewableTracker(&skew)
	notifier := &stubNotifier{}
	repo := presenceStubThreadRepo{ids: []int{7}}

	ctx, cancel := context.WithCancel(context.Background())
	go realtime.RunPresenceSession(ctx, tracker, notifier, repo, 42, testCfg)

	time.Sleep(30 * time.Millisecond)
	online := notifier.withState("ONLINE")
	assert.Len(t, online, 1, "exactly one ONLINE published")
	if len(online) == 1 {
		assert.Equal(t, presenceEvent{userID: 42, state: "ONLINE", threadID: 7}, online[0])
	}

	cancel()
	time.Sleep(20 * time.Millisecond) // let the ctx.Done branch run Disconnect
	skew.Add(int64(2 * time.Minute))  // age the last activity past presenceTTL
	time.Sleep(120 * time.Millisecond)

	offline := notifier.withState("OFFLINE")
	assert.Len(t, offline, 1, "exactly one OFFLINE after the grace window")
	if len(offline) == 1 {
		assert.Equal(t, presenceEvent{userID: 42, state: "OFFLINE", threadID: 7}, offline[0])
	}
}

func TestPresenceSession_ReconnectWithinGraceSuppressesOffline(t *testing.T) {
	var skew atomic.Int64
	tracker := newSkewableTracker(&skew)
	notifier := &stubNotifier{}
	repo := presenceStubThreadRepo{ids: []int{7}}

	ctx, cancel := context.WithCancel(context.Background())
	go realtime.RunPresenceSession(ctx, tracker, notifier, repo, 42, testCfg)

	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, 1, notifier.count("ONLINE"))

	cancel()
	time.Sleep(15 * time.Millisecond)
	// A second connection arrives inside the grace window: refcount back to 1.
	tracker.Connect(42)
	skew.Add(int64(2 * time.Minute)) // even with an aged clock, refs>0 => IsOnline
	time.Sleep(120 * time.Millisecond)

	assert.Equal(t, 0, notifier.count("OFFLINE"), "reconnect within grace suppresses OFFLINE")
}

func TestPresenceSession_SecondConnectionKeepsUserOnline(t *testing.T) {
	var skew atomic.Int64
	tracker := newSkewableTracker(&skew)
	notifier := &stubNotifier{}
	repo := presenceStubThreadRepo{ids: []int{7}}

	ctxA, cancelA := context.WithCancel(context.Background())
	ctxB, cancelB := context.WithCancel(context.Background())
	go realtime.RunPresenceSession(ctxA, tracker, notifier, repo, 42, testCfg)
	time.Sleep(15 * time.Millisecond)
	go realtime.RunPresenceSession(ctxB, tracker, notifier, repo, 42, testCfg)

	time.Sleep(40 * time.Millisecond)
	assert.Equal(t, 1, notifier.count("ONLINE"), "only the first Connect publishes ONLINE")

	// Close connection A: not the last connection, so no OFFLINE is scheduled.
	cancelA()
	time.Sleep(20 * time.Millisecond)
	skew.Add(int64(2 * time.Minute))
	time.Sleep(120 * time.Millisecond)
	assert.Equal(t, 0, notifier.count("OFFLINE"), "B still holds a connection")

	// Close connection B: the last one. OFFLINE fires after the grace window.
	cancelB()
	time.Sleep(20 * time.Millisecond)
	skew.Add(int64(2 * time.Minute))
	time.Sleep(120 * time.Millisecond)
	assert.Equal(t, 1, notifier.count("OFFLINE"), "exactly one OFFLINE after the last connection closes")
}

func TestPresenceSession_HeartbeatKeepsLastSeenFresh(t *testing.T) {
	var skew atomic.Int64
	tracker := newSkewableTracker(&skew)
	notifier := &stubNotifier{}
	repo := presenceStubThreadRepo{ids: []int{7}}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go realtime.RunPresenceSession(ctx, tracker, notifier, repo, 42, testCfg)

	// ~4 heartbeat intervals while the session is still open.
	time.Sleep(45 * time.Millisecond)
	assert.True(t, tracker.IsOnline(42), "session still open => online")
	assert.Equal(t, 0, notifier.count("OFFLINE"), "no spurious OFFLINE while the session is alive")
	assert.Equal(t, 1, notifier.count("ONLINE"))
}
