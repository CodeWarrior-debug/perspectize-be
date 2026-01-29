package resolvers

import (
	"encoding/json"
	"strconv"

	"github.com/yourorg/perspectize-go/internal/adapters/graphql/model"
	"github.com/yourorg/perspectize-go/internal/core/domain"
)

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
