package resolvers

import (
	"encoding/json"
	"strconv"

	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/adapters/graphql/model"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/domain"
)

// userDomainToModel converts a domain User to a GraphQL model User
func userDomainToModel(u *domain.User) *model.User {
	return &model.User{
		ID:        strconv.Itoa(u.ID),
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// domainToModel converts a domain Content to a GraphQL model Content
func domainToModel(c *domain.Content) *model.Content {
	m := &model.Content{
		ID:          strconv.Itoa(c.ID),
		Name:        c.Name,
		URL:         c.URL,
		ContentType: string(c.ContentType),
		Length:      c.Length,
		LengthUnits: c.LengthUnits,
		CreatedAt:   c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Parse the raw response JSON into a map for GraphQL
	if len(c.Response) > 0 {
		var responseMap map[string]interface{}
		if err := json.Unmarshal(c.Response, &responseMap); err == nil {
			m.Response = responseMap
		}

		// Also extract statistics from the YouTube API response
		var resp struct {
			Items []struct {
				Statistics struct {
					ViewCount    string `json:"viewCount"`
					LikeCount    string `json:"likeCount"`
					CommentCount string `json:"commentCount"`
				} `json:"statistics"`
			} `json:"items"`
		}
		if err := json.Unmarshal(c.Response, &resp); err == nil && len(resp.Items) > 0 {
			stats := resp.Items[0].Statistics
			if v, err := strconv.Atoi(stats.ViewCount); err == nil {
				m.ViewCount = &v
			}
			if v, err := strconv.Atoi(stats.LikeCount); err == nil {
				m.LikeCount = &v
			}
			if v, err := strconv.Atoi(stats.CommentCount); err == nil {
				m.CommentCount = &v
			}
		}
	}

	return m
}

// perspectiveDomainToModel converts a domain Perspective to a GraphQL model Perspective
func perspectiveDomainToModel(p *domain.Perspective) *model.Perspective {
	m := &model.Perspective{
		ID:           strconv.Itoa(p.ID),
		Claim:        p.Claim,
		UserID:       strconv.Itoa(p.UserID),
		Quality:      p.Quality,
		Agreement:    p.Agreement,
		Importance:   p.Importance,
		Confidence:   p.Confidence,
		Like:         p.Like,
		Privacy:      p.Privacy,
		Description:  p.Description,
		Category:     p.Category,
		Parts:        p.Parts,
		Labels:       p.Labels,
		ReviewStatus: p.ReviewStatus,
		CreatedAt:    p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if p.ContentID != nil {
		contentID := strconv.Itoa(*p.ContentID)
		m.ContentID = &contentID
	}

	// Convert categorized ratings
	if len(p.CategorizedRatings) > 0 {
		m.CategorizedRatings = make([]*model.CategorizedRating, len(p.CategorizedRatings))
		for i, cr := range p.CategorizedRatings {
			m.CategorizedRatings[i] = &model.CategorizedRating{
				Category: cr.Category,
				Rating:   cr.Rating,
			}
		}
	}

	return m
}
