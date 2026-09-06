package realtime_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/realtime"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/repositories"
	portservices "github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubMsgRepo is a no-op MessageRepository that returns a fixed message from GetByID.
type stubMsgRepo struct{ msg domain.Message }

func (s stubMsgRepo) Insert(ctx context.Context, m *domain.Message) (*domain.Message, error) {
	return m, nil
}
func (s stubMsgRepo) GetByID(ctx context.Context, id int64) (*domain.Message, error) {
	m := s.msg
	return &m, nil
}
func (s stubMsgRepo) ListHistory(ctx context.Context, t int, l int, b *int64) ([]domain.Message, error) {
	return nil, nil
}
func (s stubMsgRepo) ListSince(ctx context.Context, t int, s2 int64) ([]domain.Message, error) {
	return nil, nil
}
func (s stubMsgRepo) MaxSeq(ctx context.Context, t int) (int64, error) { return 0, nil }

// stubThreadRepo is a ThreadRepository whose GetThread returns a caller-supplied
// thread; the Hub's inbox fan-out uses it to resolve participants.
type stubThreadRepo struct{ thread *domain.MessageThread }

func (stubThreadRepo) CreateThread(ctx context.Context, createdBy int, title *string, participantUserIDs []int) (*domain.MessageThread, error) {
	return nil, nil
}
func (s stubThreadRepo) GetThread(ctx context.Context, threadID int) (*domain.MessageThread, error) {
	if s.thread == nil {
		return nil, nil
	}
	t := *s.thread
	return &t, nil
}
func (stubThreadRepo) FindDirectThread(ctx context.Context, userA, userB int) (*domain.MessageThread, error) {
	return nil, nil
}
func (stubThreadRepo) ListThreadsForUser(ctx context.Context, userID int, limit int, beforeLastMessageAt *time.Time) ([]domain.MessageThread, error) {
	return nil, nil
}
func (stubThreadRepo) AddParticipants(ctx context.Context, threadID int, userIDs []int) error {
	return nil
}
func (stubThreadRepo) SetLeft(ctx context.Context, threadID, userID int, at time.Time) error {
	return nil
}
func (stubThreadRepo) SetLastRead(ctx context.Context, threadID, userID int, seq int64) error {
	return nil
}

var (
	_ repositories.MessageRepository = stubMsgRepo{}
	_ repositories.ThreadRepository  = stubThreadRepo{}
	_ portservices.EventPublisher    = (*realtime.Hub)(nil)
)

func newHub(msg domain.Message) *realtime.Hub {
	return realtime.NewHub(stubMsgRepo{msg: msg}, stubThreadRepo{})
}

func TestHub_InboxFanout(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	thread := &domain.MessageThread{
		ID:            1,
		LastMessageAt: now,
		Participants: []domain.ThreadParticipant{
			{ThreadID: 1, UserID: 11, LastReadSeq: 2},
			{ThreadID: 1, UserID: 12, LastReadSeq: 5},
			// A departed participant must not receive inbox events.
			{ThreadID: 1, UserID: 13, LastReadSeq: 0, LeftAt: &now},
		},
	}
	hub := realtime.NewHub(
		stubMsgRepo{msg: domain.Message{ID: 10, ThreadID: 1, Seq: 5, Body: "hi", CreatedAt: now}},
		stubThreadRepo{thread: thread},
	)

	inbox11, unsub11 := hub.SubscribeInbox(11)
	defer unsub11()
	inbox12, unsub12 := hub.SubscribeInbox(12)
	defer unsub12()
	inbox13, unsub13 := hub.SubscribeInbox(13)
	defer unsub13()

	// The per-thread ThreadEvent fan-out must be unaffected.
	threadCh, unsubThread := hub.Subscribe(1)
	defer unsubThread()

	hub.PublishEnvelope(context.Background(), domain.EventEnvelope{
		Type: "MESSAGE_POSTED", ThreadID: 1, Seq: 5, MessageID: 10,
	})

	select {
	case evt := <-threadCh:
		_, ok := evt.(domain.MessagePostedEvent)
		assert.True(t, ok, "thread fan-out still delivers MessagePostedEvent")
	case <-time.After(time.Second):
		t.Fatal("no thread event")
	}

	select {
	case e := <-inbox11:
		assert.Equal(t, 1, e.ThreadID)
		assert.Equal(t, int64(5), e.LatestSeq)
		assert.Equal(t, 3, e.UnreadCount, "seq 5 minus lastRead 2")
		assert.Equal(t, now, e.LastMessageAt)
	case <-time.After(time.Second):
		t.Fatal("participant 11 got no inbox event")
	}

	select {
	case e := <-inbox12:
		assert.Equal(t, 0, e.UnreadCount, "caught-up participant has no unread")
	case <-time.After(time.Second):
		t.Fatal("participant 12 got no inbox event")
	}

	select {
	case e := <-inbox13:
		t.Fatalf("departed participant received an inbox event: %+v", e)
	case <-time.After(50 * time.Millisecond):
	}
}

func TestHub_InboxUnsubscribeClosesChannel(t *testing.T) {
	hub := newHub(domain.Message{})
	ch, unsub := hub.SubscribeInbox(42)
	unsub()
	_, open := <-ch
	assert.False(t, open, "inbox channel closed on unsubscribe")
	assert.NotPanics(t, unsub, "unsubscribe is idempotent")
}

func TestHub_InboxFanoutToleratesMissingThread(t *testing.T) {
	// stubThreadRepo with no thread returns (nil, nil): the fan-out must skip
	// rather than panic, and the thread fan-out must still happen.
	hub := newHub(domain.Message{ID: 1, ThreadID: 7, Seq: 1})
	ch, unsub := hub.Subscribe(7)
	defer unsub()

	assert.NotPanics(t, func() {
		hub.PublishEnvelope(context.Background(), domain.EventEnvelope{Type: "MESSAGE_POSTED", ThreadID: 7, MessageID: 1})
	})
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("thread fan-out did not happen")
	}
}

func TestHub_FanOutMessagePosted(t *testing.T) {
	hub := newHub(domain.Message{ID: 10, ThreadID: 1, Seq: 3, Body: "hi"})
	ch, unsub := hub.Subscribe(1)
	defer unsub()

	hub.PublishEnvelope(context.Background(), domain.EventEnvelope{Type: "MESSAGE_POSTED", ThreadID: 1, Seq: 3, MessageID: 10})

	select {
	case evt := <-ch:
		mp, ok := evt.(domain.MessagePostedEvent)
		require.True(t, ok)
		assert.Equal(t, "hi", mp.Message.Body)
	case <-time.After(time.Second):
		t.Fatal("no event")
	}
}

func TestHub_TypingEnvelopeBecomesTypingEvent(t *testing.T) {
	hub := newHub(domain.Message{})
	ch, unsub := hub.Subscribe(2)
	defer unsub()

	hub.PublishEnvelope(context.Background(), domain.EventEnvelope{Type: "TYPING_CHANGED", ThreadID: 2, UserID: 7, Typing: true})

	evt := <-ch
	tc, ok := evt.(domain.TypingChangedEvent)
	require.True(t, ok)
	assert.Equal(t, 7, tc.UserID)
	assert.True(t, tc.Typing)
}

func TestHub_UnsubscribeStopsDelivery(t *testing.T) {
	hub := newHub(domain.Message{})
	ch, unsub := hub.Subscribe(3)
	unsub()
	_, open := <-ch
	assert.False(t, open, "channel closed on unsubscribe")

	// Idempotent: a second unsub must not panic.
	assert.NotPanics(t, unsub)
}

func TestHub_SlowConsumerDropped(t *testing.T) {
	hub := newHub(domain.Message{ID: 1, ThreadID: 4})
	ch, _ := hub.Subscribe(4)

	// Never drain ch. Publish 65 message events (buffer is 64) — the 65th forces a drop.
	for i := 0; i < 65; i++ {
		hub.PublishEnvelope(context.Background(), domain.EventEnvelope{Type: "MESSAGE_POSTED", ThreadID: 4, MessageID: 1})
	}

	// Drain what buffered, then expect a close. Reaching the end means the
	// drop path closed the channel (otherwise this range would block forever).
	done := make(chan struct{})
	go func() {
		for range ch {
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("channel was not closed by the drop path")
	}
}

func TestHub_ResetAllSendsStreamReset(t *testing.T) {
	hub := newHub(domain.Message{})
	ch, unsub := hub.Subscribe(5)
	defer unsub()

	hub.ResetAll()

	evt := <-ch
	_, ok := evt.(domain.StreamResetEvent)
	assert.True(t, ok)
}

func TestHub_ConcurrentPublishSubscribeUnsub(t *testing.T) {
	hub := newHub(domain.Message{ID: 1, ThreadID: 1})
	var wg sync.WaitGroup

	// Publishers.
	for p := 0; p < 8; p++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				hub.PublishEnvelope(context.Background(), domain.EventEnvelope{Type: "TYPING_CHANGED", ThreadID: 1, UserID: i})
			}
		}()
	}

	// Subscribers that churn: subscribe, maybe drain a little, unsubscribe.
	for s := 0; s < 8; s++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				ch, unsub := hub.Subscribe(1)
				select {
				case <-ch:
				default:
				}
				unsub()
			}
		}()
	}

	// A resetter running in parallel.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			hub.ResetAll()
		}
	}()

	wg.Wait()
}

func TestHub_PublishEphemeralSatisfiesPort(t *testing.T) {
	var pub portservices.EventPublisher = newHub(domain.Message{})
	err := pub.PublishEphemeral(context.Background(), domain.EventEnvelope{Type: "PRESENCE_CHANGED", ThreadID: 9, UserID: 1, State: "ONLINE"})
	assert.NoError(t, err)
}
