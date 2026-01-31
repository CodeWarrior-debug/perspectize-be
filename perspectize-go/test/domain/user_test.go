package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yourorg/perspectize-go/internal/core/domain"
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
