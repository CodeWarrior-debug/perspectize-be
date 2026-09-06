package domain

import "time"

type ThreadRole string

const (
	ThreadRoleOwner  ThreadRole = "OWNER"
	ThreadRoleMember ThreadRole = "MEMBER"
)

type PresenceState string

const (
	PresenceOnline  PresenceState = "ONLINE"
	PresenceOffline PresenceState = "OFFLINE"
)

type MessageThread struct {
	ID            int
	Title         *string
	CreatedBy     int
	LastMessageAt time.Time
	CreatedAt     time.Time
	Participants  []ThreadParticipant
}

type ThreadParticipant struct {
	ThreadID    int
	UserID      int
	Role        ThreadRole
	LastReadSeq int64
	Muted       bool
	JoinedAt    time.Time
	LeftAt      *time.Time
}

func (t MessageThread) IsActiveParticipant(userID int) bool {
	for _, p := range t.Participants {
		if p.UserID == userID && p.LeftAt == nil {
			return true
		}
	}
	return false
}

type Message struct {
	ID          int64
	ThreadID    int
	SenderID    int
	Seq         int64
	Body        string
	ClientNonce string
	CreatedAt   time.Time
}
