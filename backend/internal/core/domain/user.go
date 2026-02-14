package domain

import "time"

// DeletedUserUsername is the sentinel user that owns orphaned content/perspectives
// after a real user is deleted. This user is seeded by migration 000006.
const DeletedUserUsername = "[deleted]"

// User represents a user who can create perspectives
type User struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsSentinel returns true if this is the system sentinel user.
func (u *User) IsSentinel() bool {
	return u.Username == DeletedUserUsername
}
