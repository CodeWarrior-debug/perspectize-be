package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/domain"
)

func TestPerspectiveStruct(t *testing.T) {
	now := time.Now()
	contentID := 1
	quality := 8000
	description := "Test description"

	perspective := domain.Perspective{
		ID:          1,
		Claim:       "Test claim",
		UserID:      1,
		ContentID:   &contentID,
		Quality:     &quality,
		Privacy:     domain.PrivacyPublic,
		Description: &description,
		Parts:       []int{1, 2, 3},
		Labels:      []string{"label1", "label2"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, 1, perspective.ID)
	assert.Equal(t, "Test claim", perspective.Claim)
	assert.Equal(t, 1, perspective.UserID)
	assert.Equal(t, &contentID, perspective.ContentID)
	assert.Equal(t, &quality, perspective.Quality)
	assert.Equal(t, domain.PrivacyPublic, perspective.Privacy)
	assert.Equal(t, &description, perspective.Description)
	assert.Equal(t, []int{1, 2, 3}, perspective.Parts)
	assert.Equal(t, []string{"label1", "label2"}, perspective.Labels)
	assert.Equal(t, now, perspective.CreatedAt)
	assert.Equal(t, now, perspective.UpdatedAt)
}

func TestPerspectiveZeroValue(t *testing.T) {
	var perspective domain.Perspective

	assert.Equal(t, 0, perspective.ID)
	assert.Equal(t, "", perspective.Claim)
	assert.Equal(t, 0, perspective.UserID)
	assert.Nil(t, perspective.ContentID)
	assert.Nil(t, perspective.Quality)
	assert.Equal(t, domain.Privacy(""), perspective.Privacy)
	assert.True(t, perspective.CreatedAt.IsZero())
	assert.True(t, perspective.UpdatedAt.IsZero())
}

func TestPrivacyConstants(t *testing.T) {
	assert.Equal(t, domain.Privacy("PUBLIC"), domain.PrivacyPublic)
	assert.Equal(t, domain.Privacy("PRIVATE"), domain.PrivacyPrivate)
}

func TestReviewStatusConstants(t *testing.T) {
	assert.Equal(t, domain.ReviewStatus("PENDING"), domain.ReviewStatusPending)
	assert.Equal(t, domain.ReviewStatus("APPROVED"), domain.ReviewStatusApproved)
	assert.Equal(t, domain.ReviewStatus("REJECTED"), domain.ReviewStatusRejected)
}

func TestRatingConstants(t *testing.T) {
	assert.Equal(t, 0, domain.RatingMin)
	assert.Equal(t, 10000, domain.RatingMax)
}

func TestCategorizedRating(t *testing.T) {
	cr := domain.CategorizedRating{
		Category: "accuracy",
		Rating:   7500,
	}

	assert.Equal(t, "accuracy", cr.Category)
	assert.Equal(t, 7500, cr.Rating)
}

func TestPerspectiveSortByConstants(t *testing.T) {
	assert.Equal(t, domain.PerspectiveSortBy("CREATED_AT"), domain.PerspectiveSortByCreatedAt)
	assert.Equal(t, domain.PerspectiveSortBy("UPDATED_AT"), domain.PerspectiveSortByUpdatedAt)
	assert.Equal(t, domain.PerspectiveSortBy("CLAIM"), domain.PerspectiveSortByClaim)
}

func TestPerspectiveWithCategorizedRatings(t *testing.T) {
	perspective := domain.Perspective{
		ID:     1,
		Claim:  "Test",
		UserID: 1,
		CategorizedRatings: []domain.CategorizedRating{
			{Category: "accuracy", Rating: 8000},
			{Category: "clarity", Rating: 9000},
		},
	}

	assert.Len(t, perspective.CategorizedRatings, 2)
	assert.Equal(t, "accuracy", perspective.CategorizedRatings[0].Category)
	assert.Equal(t, 8000, perspective.CategorizedRatings[0].Rating)
	assert.Equal(t, "clarity", perspective.CategorizedRatings[1].Category)
	assert.Equal(t, 9000, perspective.CategorizedRatings[1].Rating)
}
