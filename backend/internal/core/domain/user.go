package domain

import "time"

// DeletedUserUsername is the sentinel user that owns orphaned content/perspectives
// after a real user is deleted. This user is seeded by migration 000006.
const DeletedUserUsername = "[deleted]"

// SystemUserUsername is the sentinel user that owns pre-existing content
// created before user tracking was added. Seeded by migration 000007.
const SystemUserUsername = "[system]"

// User represents a user who can create perspectives
type User struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsSentinel returns true if this is a system sentinel user ([deleted] or [system]).
func (u *User) IsSentinel() bool {
	return u.Username == DeletedUserUsername || u.Username == SystemUserUsername
}
