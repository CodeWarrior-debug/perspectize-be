package domain

import "time"

// Identity is the authenticated principal, resolved from a Clerk token by a
// TokenVerifier. It is transport-agnostic (HTTP middleware and the WebSocket
// InitFunc both produce it).
type Identity struct {
	UserID   int
	ClerkID  string
	Username string
}

// EventEnvelope is the decoded pg_notify('thread_events', ...) payload.
type EventEnvelope struct {
	Type      string `json:"type"`
	ThreadID  int    `json:"thread_id"`
	Seq       int64  `json:"seq"`
	MessageID int64  `json:"message_id"`
	// Ephemeral-only fields (absent for MESSAGE_POSTED):
	UserID      int    `json:"user_id,omitempty"`
	Typing      bool   `json:"typing,omitempty"`
	LastReadSeq int64  `json:"last_read_seq,omitempty"`
	Change      string `json:"change,omitempty"`
	State       string `json:"state,omitempty"`
}

type ThreadEvent interface{ isThreadEvent() }

type MessagePostedEvent struct{ Message Message }
type ReadReceiptChangedEvent struct {
	ThreadID    int
	UserID      int
	LastReadSeq int64
}
type TypingChangedEvent struct {
	ThreadID int
	UserID   int
	Typing   bool
}
type ParticipantChangedEvent struct {
	ThreadID int
	UserID   int
	Change   string // "ADDED" | "REMOVED"
}
type PresenceChangedEvent struct {
	ThreadID int
	UserID   int
	State    PresenceState
}
type StreamResetEvent struct{ ThreadID int }

func (MessagePostedEvent) isThreadEvent()      {}
func (ReadReceiptChangedEvent) isThreadEvent() {}
func (TypingChangedEvent) isThreadEvent()      {}
func (ParticipantChangedEvent) isThreadEvent() {}
func (PresenceChangedEvent) isThreadEvent()    {}
func (StreamResetEvent) isThreadEvent()        {}

type InboxEvent struct {
	ThreadID      int
	LastMessageAt time.Time
	LatestSeq     int64
	UnreadCount   int
}
