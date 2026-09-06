package services

import (
	"context"
	"time"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
)

// EventPublisher is the port the MessagingService uses to emit ephemeral
// realtime events (typing, read receipts, participant changes). The realtime
// Hub (Task 9) implements it. MESSAGE_POSTED is NOT published here — it
// originates from the database trigger.
type EventPublisher interface {
	PublishEphemeral(ctx context.Context, env domain.EventEnvelope) error
}

// SendMessageInput is the payload for MessagingService.SendMessage.
type SendMessageInput struct {
	ThreadID    int
	Body        string
	ClientNonce string
}

// MessagingService is the business-logic port for the messaging feature:
// authorization, rate limiting, idempotent send, 1:1 thread dedup,
// read-pointer updates, and ephemeral event publishing.
type MessagingService interface {
	CreateThread(ctx context.Context, actorUserID int, participantUserIDs []int, title *string) (*domain.MessageThread, error)
	SendMessage(ctx context.Context, actorUserID int, in SendMessageInput) (*domain.Message, error)
	MarkRead(ctx context.Context, actorUserID, threadID int, seq int64) (*domain.MessageThread, error)
	AddParticipants(ctx context.Context, actorUserID, threadID int, userIDs []int) (*domain.MessageThread, error)
	LeaveThread(ctx context.Context, actorUserID, threadID int) error
	ListThreads(ctx context.Context, actorUserID int, limit int, beforeLastMessageAt *time.Time) ([]domain.MessageThread, error)
	GetHistory(ctx context.Context, actorUserID, threadID int, limit int, beforeSeq *int64) ([]domain.Message, error)
	ListSince(ctx context.Context, actorUserID, threadID int, sinceSeq int64) ([]domain.Message, error)
	SetTyping(ctx context.Context, actorUserID, threadID int, typing bool) error
	AssertParticipant(ctx context.Context, actorUserID, threadID int) error
	// GetThread returns a single thread the actor participates in. Enforces
	// participation first. Needed by the T11 messageThread query resolver.
	GetThread(ctx context.Context, actorUserID, threadID int) (*domain.MessageThread, error)
	// MaxSeq returns the highest message seq in the thread. Enforces
	// participation first. Needed by T11 GraphQL field resolvers.
	MaxSeq(ctx context.Context, actorUserID, threadID int) (int64, error)
}
