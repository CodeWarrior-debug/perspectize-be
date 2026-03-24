package domain

// AuthenticatedUser holds resolved local user info for request context.
// Used by auth middleware and accessed via auth.ForContext(ctx).
type AuthenticatedUser struct {
	ID       int
	ClerkID  string
	Username string
	Email    string
	Role     UserRole
}
