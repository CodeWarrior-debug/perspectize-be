package domain_test

import (
	"testing"
	"time"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserStruct(t *testing.T) {
	now := time.Now()

	user := domain.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, now, user.CreatedAt)
	assert.Equal(t, now, user.UpdatedAt)
}

func TestUserZeroValue(t *testing.T) {
	var user domain.User

	assert.Equal(t, 0, user.ID)
	assert.Equal(t, "", user.Username)
	assert.Equal(t, "", user.Email)
	assert.True(t, user.CreatedAt.IsZero())
	assert.True(t, user.UpdatedAt.IsZero())
}

func TestIsSentinel_RoleBased(t *testing.T) {
	sentinel := &domain.User{Role: domain.UserRoleSentinel, Username: "[deleted]"}
	assert.True(t, sentinel.IsSentinel())

	admin := &domain.User{Role: domain.UserRoleAdmin, Username: "admin"}
	assert.False(t, admin.IsSentinel())

	regular := &domain.User{Role: domain.UserRoleDefault, Username: "alice"}
	assert.False(t, regular.IsSentinel())
}

func TestUserRoleConstants(t *testing.T) {
	assert.Equal(t, domain.UserRole("ADMIN"), domain.UserRoleAdmin)
	assert.Equal(t, domain.UserRole("SENTINEL"), domain.UserRoleSentinel)
	assert.Equal(t, domain.UserRole("DEFAULT"), domain.UserRoleDefault)
}
