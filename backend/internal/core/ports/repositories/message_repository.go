package repositories

import (
	"context"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
)

// MessageRepository is the port for message storage operations.
type MessageRepository interface {
	Insert(ctx context.Context, m *domain.Message) (*domain.Message, error)
	GetByID(ctx context.Context, id int64) (*domain.Message, error)
	ListHistory(ctx context.Context, threadID int, limit int, beforeSeq *int64) ([]domain.Message, error)
	ListSince(ctx context.Context, threadID int, sinceSeq int64) ([]domain.Message, error)
	MaxSeq(ctx context.Context, threadID int) (int64, error)
	// CountSince returns how many messages in the thread have seq strictly
	// greater than sinceSeq. Counting rows (rather than subtracting seqs) stays
	// correct when pruning has left gaps in the sequence.
	CountSince(ctx context.Context, threadID int, sinceSeq int64) (int, error)
}
