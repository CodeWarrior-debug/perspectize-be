package auth

import (
	"context"
	"testing"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForContext(t *testing.T) {
	t.Run("bare context has no authenticated user", func(t *testing.T) {
		user, ok := ForContext(context.Background())
		assert.False(t, ok)
		assert.Nil(t, user)
	})

	t.Run("nil user stored under the key is treated as unauthenticated", func(t *testing.T) {
		ctx := withUser(context.Background(), nil)
		user, ok := ForContext(ctx)
		assert.False(t, ok)
		assert.Nil(t, user)
	})

	t.Run("round-trips a stored user", func(t *testing.T) {
		want := &domain.AuthenticatedUser{ID: 7, ClerkID: "user_abc", Username: "alice", Email: "alice@example.com", Role: domain.UserRoleAdmin}
		user, ok := ForContext(withUser(context.Background(), want))
		require.True(t, ok)
		assert.Equal(t, want, user)
	})

	t.Run("WithAuthenticatedUser uses the same key as withUser", func(t *testing.T) {
		want := &domain.AuthenticatedUser{ID: 9, Username: "bob", Role: domain.UserRoleDefault}
		user, ok := ForContext(WithAuthenticatedUser(context.Background(), want))
		require.True(t, ok)
		assert.Equal(t, want, user)
	})
}

func TestRequireAuth(t *testing.T) {
	t.Run("returns an access-denied error when unauthenticated", func(t *testing.T) {
		user, err := RequireAuth(context.Background())
		assert.Nil(t, user)
		require.Error(t, err)
		assert.Equal(t, "access denied: authentication required", err.Error())
	})

	t.Run("returns the stored user when authenticated", func(t *testing.T) {
		want := &domain.AuthenticatedUser{ID: 7, Username: "alice", Role: domain.UserRoleAdmin}
		user, err := RequireAuth(WithAuthenticatedUser(context.Background(), want))
		require.NoError(t, err)
		assert.Equal(t, want, user)
	})
}
