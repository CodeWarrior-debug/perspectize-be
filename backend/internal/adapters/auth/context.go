package auth

import (
	"context"
	"fmt"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
)

type contextKey string

const authContextKey contextKey = "authenticated_user"

// ForContext extracts the authenticated user from the request context.
// Returns the user and true if authenticated, nil and false otherwise.
func ForContext(ctx context.Context) (*domain.AuthenticatedUser, bool) {
	user, ok := ctx.Value(authContextKey).(*domain.AuthenticatedUser)
	if !ok || user == nil {
		return nil, false
	}
	return user, true
}

// RequireAuth extracts the authenticated user or returns an error.
func RequireAuth(ctx context.Context) (*domain.AuthenticatedUser, error) {
	user, ok := ForContext(ctx)
	if !ok {
		return nil, fmt.Errorf("access denied: authentication required")
	}
	return user, nil
}

// withUser stores an authenticated user in the context.
func withUser(ctx context.Context, user *domain.AuthenticatedUser) context.Context {
	return context.WithValue(ctx, authContextKey, user)
}

// WithAuthenticatedUser stores an authenticated user in the context.
// Exported for use in tests that need to simulate authenticated requests.
func WithAuthenticatedUser(ctx context.Context, user *domain.AuthenticatedUser) context.Context {
	return context.WithValue(ctx, authContextKey, user)
}
