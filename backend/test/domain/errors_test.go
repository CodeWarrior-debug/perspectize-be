package domain_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestDomainErrors_AreDefined(t *testing.T) {
	assert.NotNil(t, domain.ErrNotFound)
	assert.NotNil(t, domain.ErrAlreadyExists)
	assert.NotNil(t, domain.ErrInvalidInput)
	assert.NotNil(t, domain.ErrInvalidURL)
	assert.NotNil(t, domain.ErrYouTubeAPI)
}

func TestDomainErrors_Messages(t *testing.T) {
	assert.Equal(t, "resource not found", domain.ErrNotFound.Error())
	assert.Equal(t, "resource already exists", domain.ErrAlreadyExists.Error())
	assert.Equal(t, "invalid input", domain.ErrInvalidInput.Error())
	assert.Equal(t, "invalid URL", domain.ErrInvalidURL.Error())
	assert.Equal(t, "youtube API error", domain.ErrYouTubeAPI.Error())
}

func TestDomainErrors_AreDistinct(t *testing.T) {
	errs := []error{
		domain.ErrNotFound,
		domain.ErrAlreadyExists,
		domain.ErrInvalidInput,
		domain.ErrInvalidURL,
		domain.ErrYouTubeAPI,
	}

	for i, err1 := range errs {
		for j, err2 := range errs {
			if i != j {
				assert.False(t, errors.Is(err1, err2),
					"Expected %v and %v to be distinct errors", err1, err2)
			}
		}
	}
}

func TestDomainErrors_WrappedErrorsAreDetectable(t *testing.T) {
	wrapped := fmt.Errorf("something failed: %w", domain.ErrNotFound)
	assert.True(t, errors.Is(wrapped, domain.ErrNotFound))
	assert.False(t, errors.Is(wrapped, domain.ErrAlreadyExists))
}

func TestDomainErrors_DoubleWrappedErrorsAreDetectable(t *testing.T) {
	inner := fmt.Errorf("repo error: %w", domain.ErrNotFound)
	outer := fmt.Errorf("service error: %w", inner)

	assert.True(t, errors.Is(outer, domain.ErrNotFound))
	assert.Contains(t, outer.Error(), "service error")
	assert.Contains(t, outer.Error(), "repo error")
	assert.Contains(t, outer.Error(), "resource not found")
}

func TestErrInvalidInput_CanBeWrapped(t *testing.T) {
	wrapped := fmt.Errorf("%w: id must be positive", domain.ErrInvalidInput)
	assert.True(t, errors.Is(wrapped, domain.ErrInvalidInput))
	assert.Contains(t, wrapped.Error(), "invalid input")
	assert.Contains(t, wrapped.Error(), "id must be positive")
}
