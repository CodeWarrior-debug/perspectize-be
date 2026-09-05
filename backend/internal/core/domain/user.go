package domain

import "time"

// DeletedUserUsername is the sentinel user that owns orphaned content/perspectives
// after a real user is deleted. This user is seeded by migration 000006.
const DeletedUserUsername = "[deleted]"

// SystemUserUsername is the sentinel user that owns pre-existing content
// created before user tracking was added. Seeded by migration 000007.
const SystemUserUsername = "[system]"

// CurrentIntroVersion is the checklist-coach intro version shipped with this build.
// Bump only when a material onboarding flow change should force a re-intro.
const CurrentIntroVersion = 1

// UserRole represents the role of a user in the system.
type UserRole string

const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleSentinel UserRole = "SENTINEL"
	UserRoleDefault  UserRole = "DEFAULT"
)

// UserOnboarding is thin persistence for the first-run checklist coach.
// CompletedAt is an ISO-8601 UTC timestamp string when set.
type UserOnboarding struct {
	Version            int     `json:"version"`
	DisplayNextSession bool    `json:"displayNextSession"`
	CompletedAt        *string `json:"completedAt"`
}

// DefaultUserOnboarding returns new-user coach defaults (show intro).
func DefaultUserOnboarding() UserOnboarding {
	return UserOnboarding{
		Version:            0,
		DisplayNextSession: true,
		CompletedAt:        nil,
	}
}

// User represents a user who can create perspectives
type User struct {
	ID          int
	ClerkUserID string
	Username    string
	Email       string
	Role        UserRole
	Active      bool
	Onboarding  UserOnboarding
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// IsSentinel returns true if this is a system sentinel user ([deleted] or [system]).
func (u *User) IsSentinel() bool {
	return u.Role == UserRoleSentinel
}
