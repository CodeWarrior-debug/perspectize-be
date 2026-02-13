package domain

import (
	"encoding/json"
	"time"
)

// Privacy represents the visibility level of a perspective
type Privacy string

const (
	PrivacyPublic  Privacy = "PUBLIC"
	PrivacyPrivate Privacy = "PRIVATE"
)

// ReviewStatus represents the review state of a perspective
type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "PENDING"
	ReviewStatusApproved ReviewStatus = "APPROVED"
	ReviewStatusRejected ReviewStatus = "REJECTED"
)

// CategorizedRating represents a rating with a category label
type CategorizedRating struct {
	Category string `json:"category"`
	Rating   int    `json:"rating"`
}

// Perspective represents a user's viewpoint on content
type Perspective struct {
	ID        int
	Claim     string // Required, max 255 chars
	UserID    int    // Required, FK to users
	ContentID *int   // Optional, FK to content

	// Optional ratings (0-10000 range enforced by DB domain)
	Quality    *int
	Agreement  *int
	Importance *int
	Confidence *int

	// Optional fields
	Like         *string // Freeform text
	Privacy      Privacy
	Description  *string
	Category     *string
	ReviewStatus *ReviewStatus

	// Array fields
	Parts  []int    // Array of part identifiers
	Labels []string // Array of label strings

	// JSONB field
	CategorizedRatings []CategorizedRating

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RatingMin is the minimum valid rating value
const RatingMin = 0

// RatingMax is the maximum valid rating value
const RatingMax = 10000

// ValidateRating checks if a rating value is within the valid range
func ValidateRating(rating *int) bool {
	if rating == nil {
		return true // nil is valid (optional field)
	}
	return *rating >= RatingMin && *rating <= RatingMax
}

// PerspectiveSortBy represents sortable fields for perspective queries
type PerspectiveSortBy string

const (
	PerspectiveSortByCreatedAt PerspectiveSortBy = "CREATED_AT"
	PerspectiveSortByUpdatedAt PerspectiveSortBy = "UPDATED_AT"
	PerspectiveSortByClaim     PerspectiveSortBy = "CLAIM"
)

// PerspectiveFilter contains filter criteria for perspective queries
type PerspectiveFilter struct {
	UserID    *int
	ContentID *int
	Privacy   *Privacy
}

// PerspectiveListParams contains parameters for paginated perspective queries
type PerspectiveListParams struct {
	First             *int
	After             *string
	Last              *int
	Before            *string
	SortBy            PerspectiveSortBy
	SortOrder         SortOrder
	IncludeTotalCount bool
	Filter            *PerspectiveFilter
}

// PaginatedPerspectives represents a paginated list of perspectives
type PaginatedPerspectives struct {
	Items       []*Perspective
	HasNext     bool
	HasPrev     bool
	StartCursor *string
	EndCursor   *string
	TotalCount  *int
}

// MarshalCategorizedRatings converts CategorizedRatings to JSON for storage
func (p *Perspective) MarshalCategorizedRatings() ([]json.RawMessage, error) {
	if len(p.CategorizedRatings) == 0 {
		return nil, nil
	}
	result := make([]json.RawMessage, len(p.CategorizedRatings))
	for i, cr := range p.CategorizedRatings {
		data, err := json.Marshal(cr)
		if err != nil {
			return nil, err
		}
		result[i] = data
	}
	return result, nil
}
