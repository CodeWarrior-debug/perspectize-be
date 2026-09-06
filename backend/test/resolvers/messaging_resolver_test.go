package resolvers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/auth"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/graphql/model"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/graphql/resolvers"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/realtime"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/repositories"
	portservices "github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeMessaging is a hand-rolled portservices.MessagingService whose behaviour
// each test overrides through the function fields it cares about. Unset fields
// return zero values, which is enough for resolvers that do not call them.
type fakeMessaging struct {
	sendFn              func(ctx context.Context, actor int, in portservices.SendMessageInput) (*domain.Message, error)
	getThreadFn         func(ctx context.Context, actor, threadID int) (*domain.MessageThread, error)
	maxSeqFn            func(ctx context.Context, actor, threadID int) (int64, error)
	listSinceFn         func(ctx context.Context, actor, threadID int, sinceSeq int64) ([]domain.Message, error)
	assertParticipantFn func(ctx context.Context, actor, threadID int) error
}

var _ portservices.MessagingService = (*fakeMessaging)(nil)

func (f *fakeMessaging) CreateThread(ctx context.Context, actor int, ids []int, title *string) (*domain.MessageThread, error) {
	return nil, nil
}

func (f *fakeMessaging) SendMessage(ctx context.Context, actor int, in portservices.SendMessageInput) (*domain.Message, error) {
	if f.sendFn == nil {
		return nil, nil
	}
	return f.sendFn(ctx, actor, in)
}

func (f *fakeMessaging) MarkRead(ctx context.Context, actor, threadID int, seq int64) (*domain.MessageThread, error) {
	return nil, nil
}

func (f *fakeMessaging) AddParticipants(ctx context.Context, actor, threadID int, userIDs []int) (*domain.MessageThread, error) {
	return nil, nil
}

func (f *fakeMessaging) LeaveThread(ctx context.Context, actor, threadID int) error { return nil }

func (f *fakeMessaging) ListThreads(ctx context.Context, actor, limit int, before *time.Time) ([]domain.MessageThread, error) {
	return nil, nil
}

func (f *fakeMessaging) GetHistory(ctx context.Context, actor, threadID, limit int, beforeSeq *int64) ([]domain.Message, error) {
	return nil, nil
}

func (f *fakeMessaging) ListSince(ctx context.Context, actor, threadID int, sinceSeq int64) ([]domain.Message, error) {
	if f.listSinceFn == nil {
		return nil, nil
	}
	return f.listSinceFn(ctx, actor, threadID, sinceSeq)
}

func (f *fakeMessaging) SetTyping(ctx context.Context, actor, threadID int, typing bool) error {
	return nil
}

func (f *fakeMessaging) AssertParticipant(ctx context.Context, actor, threadID int) error {
	if f.assertParticipantFn == nil {
		return nil
	}
	return f.assertParticipantFn(ctx, actor, threadID)
}

func (f *fakeMessaging) GetThread(ctx context.Context, actor, threadID int) (*domain.MessageThread, error) {
	if f.getThreadFn == nil {
		return nil, nil
	}
	return f.getThreadFn(ctx, actor, threadID)
}

func (f *fakeMessaging) MaxSeq(ctx context.Context, actor, threadID int) (int64, error) {
	if f.maxSeqFn == nil {
		return 0, nil
	}
	return f.maxSeqFn(ctx, actor, threadID)
}

// inboxStubMsgRepo / inboxStubThreadRepo are the minimum repository surface the
// Hub touches when fanning a MESSAGE_POSTED envelope out to per-user inboxes.
type inboxStubMsgRepo struct{ msg domain.Message }

func (s inboxStubMsgRepo) Insert(ctx context.Context, m *domain.Message) (*domain.Message, error) {
	return m, nil
}
func (s inboxStubMsgRepo) GetByID(ctx context.Context, id int64) (*domain.Message, error) {
	m := s.msg
	return &m, nil
}
func (s inboxStubMsgRepo) ListHistory(ctx context.Context, threadID, limit int, beforeSeq *int64) ([]domain.Message, error) {
	return nil, nil
}
func (s inboxStubMsgRepo) ListSince(ctx context.Context, threadID int, sinceSeq int64) ([]domain.Message, error) {
	return nil, nil
}
func (s inboxStubMsgRepo) MaxSeq(ctx context.Context, threadID int) (int64, error) { return 0, nil }

type inboxStubThreadRepo struct{ thread domain.MessageThread }

func (inboxStubThreadRepo) CreateThread(ctx context.Context, createdBy int, title *string, ids []int) (*domain.MessageThread, error) {
	return nil, nil
}
func (s inboxStubThreadRepo) GetThread(ctx context.Context, threadID int) (*domain.MessageThread, error) {
	t := s.thread
	return &t, nil
}
func (inboxStubThreadRepo) FindDirectThread(ctx context.Context, a, b int) (*domain.MessageThread, error) {
	return nil, nil
}
func (inboxStubThreadRepo) ListThreadsForUser(ctx context.Context, userID, limit int, before *time.Time) ([]domain.MessageThread, error) {
	return nil, nil
}
func (inboxStubThreadRepo) AddParticipants(ctx context.Context, threadID int, userIDs []int) error {
	return nil
}
func (inboxStubThreadRepo) SetLeft(ctx context.Context, threadID, userID int, at time.Time) error {
	return nil
}
func (inboxStubThreadRepo) SetLastRead(ctx context.Context, threadID, userID int, seq int64) error {
	return nil
}

var (
	_ repositories.MessageRepository = inboxStubMsgRepo{}
	_ repositories.ThreadRepository  = inboxStubThreadRepo{}
)

func authedCtx(userID int) context.Context {
	return auth.WithAuthenticatedUser(context.Background(), &domain.AuthenticatedUser{ID: userID})
}

func TestSendMessageResolver_MapsDomainToModel(t *testing.T) {
	created := time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC)
	fake := &fakeMessaging{sendFn: func(ctx context.Context, actor int, in portservices.SendMessageInput) (*domain.Message, error) {
		assert.Equal(t, 1, actor)
		assert.Equal(t, 2, in.ThreadID)
		assert.Equal(t, "hello", in.Body)
		assert.Equal(t, "n1", in.ClientNonce)
		return &domain.Message{ID: 9, ThreadID: 2, SenderID: 1, Seq: 4, Body: "hello", CreatedAt: created}, nil
	}}
	r := &resolvers.Resolver{Messaging: fake}

	out, err := r.Mutation().SendMessage(authedCtx(1), model.SendMessageInput{ThreadID: "2", Body: "hello", ClientNonce: "n1"})

	require.NoError(t, err)
	assert.Equal(t, "9", out.ID)
	assert.Equal(t, "2", out.ThreadID)
	assert.Equal(t, 4, out.Seq)
	assert.Equal(t, "hello", out.Body)
	assert.Equal(t, "2026-09-06T12:00:00Z", out.CreatedAt)
}

func TestSendMessageResolver_ForbiddenMapsToGraphQLError(t *testing.T) {
	fake := &fakeMessaging{sendFn: func(context.Context, int, portservices.SendMessageInput) (*domain.Message, error) {
		return nil, fmt.Errorf("%w: nope", domain.ErrForbidden)
	}}
	r := &resolvers.Resolver{Messaging: fake}

	_, err := r.Mutation().SendMessage(authedCtx(1), model.SendMessageInput{ThreadID: "2", Body: "x", ClientNonce: "n"})

	assert.ErrorContains(t, err, "access denied")
}

func TestSendMessageResolver_UnauthenticatedIsForbidden(t *testing.T) {
	r := &resolvers.Resolver{Messaging: &fakeMessaging{}}

	_, err := r.Mutation().SendMessage(context.Background(), model.SendMessageInput{ThreadID: "2", Body: "x", ClientNonce: "n"})

	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestSendMessageResolver_NonNumericThreadIDIsInvalidInput(t *testing.T) {
	r := &resolvers.Resolver{Messaging: &fakeMessaging{}}

	_, err := r.Mutation().SendMessage(authedCtx(1), model.SendMessageInput{ThreadID: "abc", Body: "x", ClientNonce: "n"})

	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestMessageThreadResolver_ReadPointers(t *testing.T) {
	thread := &domain.MessageThread{
		ID: 5,
		Participants: []domain.ThreadParticipant{
			{ThreadID: 5, UserID: 1, LastReadSeq: 7, Role: domain.ThreadRoleOwner},
			{ThreadID: 5, UserID: 2, LastReadSeq: 1, Role: domain.ThreadRoleMember},
		},
	}
	fake := &fakeMessaging{
		getThreadFn: func(context.Context, int, int) (*domain.MessageThread, error) { return thread, nil },
		maxSeqFn:    func(context.Context, int, int) (int64, error) { return 10, nil },
	}
	r := &resolvers.Resolver{Messaging: fake}
	ctx := authedCtx(1)

	got, err := r.Query().MessageThread(ctx, "5")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "5", got.ID)

	latest, err := r.MessageThread().LatestSeq(ctx, got)
	require.NoError(t, err)
	assert.Equal(t, 10, latest)

	mine, err := r.MessageThread().MyLastReadSeq(ctx, got)
	require.NoError(t, err)
	assert.Equal(t, 7, mine)

	unread, err := r.MessageThread().UnreadCount(ctx, got)
	require.NoError(t, err)
	assert.Equal(t, 3, unread)

	parts, err := r.MessageThread().Participants(ctx, got)
	require.NoError(t, err)
	require.Len(t, parts, 2)
	assert.Equal(t, domain.ThreadRoleOwner, parts[0].Role)
	assert.Equal(t, 1, parts[0].SrcUserID)
}

func TestThreadEventsSubscription_ReplaysThenStreamsAndStopsOnCancel(t *testing.T) {
	fake := &fakeMessaging{
		listSinceFn: func(_ context.Context, _, _ int, sinceSeq int64) ([]domain.Message, error) {
			assert.Equal(t, int64(3), sinceSeq)
			return []domain.Message{{ID: 40, ThreadID: 2, SenderID: 1, Seq: 4, Body: "replayed"}}, nil
		},
	}
	hub := realtime.NewHub(nil, nil)
	r := &resolvers.Resolver{Messaging: fake, Hub: hub}

	ctx, cancel := context.WithCancel(authedCtx(1))
	since := 3
	ch, err := r.Subscription().ThreadEvents(ctx, "2", &since)
	require.NoError(t, err)

	// 1. The replayed history arrives first.
	select {
	case evt := <-ch:
		mp, ok := evt.(model.MessagePosted)
		require.True(t, ok, "expected MessagePosted, got %T", evt)
		assert.Equal(t, "replayed", mp.Message.Body)
	case <-time.After(2 * time.Second):
		t.Fatal("no replayed event")
	}

	// 2. Live hub events follow. Retry briefly: the goroutine attaches to the
	// hub before replaying, but the broadcast is non-blocking either way.
	require.Eventually(t, func() bool {
		hub.Broadcast(2, domain.TypingChangedEvent{ThreadID: 2, UserID: 9, Typing: true})
		select {
		case evt := <-ch:
			tc, ok := evt.(model.TypingChanged)
			return ok && tc.UserID == "9"
		case <-time.After(50 * time.Millisecond):
			return false
		}
	}, 2*time.Second, 10*time.Millisecond, "no live typing event")

	// 3. Cancelling the operation context tears the goroutine down and closes
	// the output channel — no leak, no send on a closed channel.
	cancel()
	require.Eventually(t, func() bool {
		select {
		case _, open := <-ch:
			return !open
		default:
			return false
		}
	}, 2*time.Second, 10*time.Millisecond, "subscription channel not closed after cancel")
}

func TestThreadEventsSubscription_HubDropEmitsStreamReset(t *testing.T) {
	hub := realtime.NewHub(nil, nil)
	r := &resolvers.Resolver{Messaging: &fakeMessaging{}, Hub: hub}

	ch, err := r.Subscription().ThreadEvents(authedCtx(1), "2", nil)
	require.NoError(t, err)

	// ResetAll then a forced drop: overflow the hub's 64-slot buffer so the hub
	// closes our upstream channel, which must surface as a StreamReset.
	for i := 0; i < 200; i++ {
		hub.Broadcast(2, domain.TypingChangedEvent{ThreadID: 2, UserID: i})
	}

	deadline := time.After(3 * time.Second)
	for {
		select {
		case evt, open := <-ch:
			if !open {
				t.Fatal("channel closed before a StreamReset was emitted")
			}
			if _, ok := evt.(model.StreamReset); ok {
				return
			}
		case <-deadline:
			t.Fatal("no StreamReset after the hub dropped the subscriber")
		}
	}
}

func TestThreadEventsSubscription_NonParticipantRejected(t *testing.T) {
	fake := &fakeMessaging{assertParticipantFn: func(context.Context, int, int) error {
		return fmt.Errorf("%w: not a participant", domain.ErrForbidden)
	}}
	r := &resolvers.Resolver{Messaging: fake, Hub: realtime.NewHub(nil, nil)}

	_, err := r.Subscription().ThreadEvents(authedCtx(1), "2", nil)

	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestInboxEventsSubscription_DeliversAndClosesOnCancel(t *testing.T) {
	at := time.Date(2026, 9, 6, 8, 30, 0, 0, time.UTC)
	hub := realtime.NewHub(
		inboxStubMsgRepo{msg: domain.Message{ID: 90, ThreadID: 3, SenderID: 12, Seq: 9, CreatedAt: at}},
		inboxStubThreadRepo{thread: domain.MessageThread{
			ID:            3,
			LastMessageAt: at,
			Participants: []domain.ThreadParticipant{
				{ThreadID: 3, UserID: 11, LastReadSeq: 7},
				{ThreadID: 3, UserID: 12, LastReadSeq: 9},
			},
		}},
	)
	r := &resolvers.Resolver{Messaging: &fakeMessaging{}, Hub: hub}

	ctx, cancel := context.WithCancel(authedCtx(11))
	ch, err := r.Subscription().InboxEvents(ctx)
	require.NoError(t, err)

	hub.PublishEnvelope(context.Background(), domain.EventEnvelope{
		Type: "MESSAGE_POSTED", ThreadID: 3, Seq: 9, MessageID: 90,
	})

	select {
	case evt := <-ch:
		require.NotNil(t, evt)
		assert.Equal(t, "3", evt.ThreadID)
		assert.Equal(t, 9, evt.LatestSeq)
		assert.Equal(t, 2, evt.UnreadCount)
		assert.Equal(t, "2026-09-06T08:30:00Z", evt.LastMessageAt)
	case <-time.After(3 * time.Second):
		t.Fatal("no inbox event")
	}

	cancel()
	require.Eventually(t, func() bool {
		select {
		case _, open := <-ch:
			return !open
		default:
			return false
		}
	}, 2*time.Second, 10*time.Millisecond, "inbox channel not closed after cancel")
}

func TestInboxEventsSubscription_UnauthenticatedIsForbidden(t *testing.T) {
	r := &resolvers.Resolver{Messaging: &fakeMessaging{}, Hub: realtime.NewHub(nil, nil)}

	_, err := r.Subscription().InboxEvents(context.Background())

	assert.ErrorIs(t, err, domain.ErrForbidden)
}
