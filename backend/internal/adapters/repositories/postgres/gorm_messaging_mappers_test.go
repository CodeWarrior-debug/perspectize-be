package postgres

import (
	"testing"
	"time"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestMessageMapperRoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	d := &domain.Message{
		ID: 7, ThreadID: 3, SenderID: 5, Seq: 42,
		Body: "hello", ClientNonce: "n-1", CreatedAt: now,
	}
	m := messageDomainToModel(d)
	assert.Equal(t, int64(7), m.ID)
	assert.Equal(t, "hello", m.Body)

	back := messageModelToDomain(m)
	assert.Equal(t, d.ThreadID, back.ThreadID)
	assert.Equal(t, d.Seq, back.Seq)
	assert.Equal(t, d.ClientNonce, back.ClientNonce)
	assert.Equal(t, now, back.CreatedAt.UTC())
}

func TestThreadParticipantMapper(t *testing.T) {
	left := time.Now().UTC()
	m := &ThreadParticipantModel{
		ThreadID: 1, UserID: 2, Role: "OWNER", LastReadSeq: 9, LeftAt: &left,
	}
	d := threadParticipantModelToDomain(m)
	assert.Equal(t, domain.ThreadRoleOwner, d.Role)
	assert.Equal(t, int64(9), d.LastReadSeq)
	assert.NotNil(t, d.LeftAt)
}
