package resolvers

import (
	"encoding/json"
	"strconv"

	"github.com/CodeWarrior-debug/perspectize-be/apps/backend/internal/adapters/graphql/model"
	"github.com/CodeWarrior-debug/perspectize-be/apps/backend/internal/core/domain"
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
		ID:          strconv.Itoa(p.ID),
		Claim:       p.Claim,
		UserID:      strconv.Itoa(p.UserID),
		Quality:     p.Quality,
		Agreement:   p.Agreement,
		Importance:  p.Importance,
		Confidence:  p.Confidence,
		Like:        p.Like,
		Privacy:     privacyDomainToModel(p.Privacy),
		Description: p.Description,
		Category:    p.Category,
		Parts:       p.Parts,
		Labels:      p.Labels,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if p.ContentID != nil {
		contentID := strconv.Itoa(*p.ContentID)
		m.ContentID = &contentID
	}

	if p.ReviewStatus != nil {
		status := reviewStatusDomainToModel(*p.ReviewStatus)
		m.ReviewStatus = &status
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

// privacyDomainToModel converts domain Privacy to GraphQL model Privacy
func privacyDomainToModel(p domain.Privacy) model.Privacy {
	switch p {
	case domain.PrivacyPrivate:
		return model.PrivacyPrivate
	default:
		return model.PrivacyPublic
	}
}

// privacyModelToDomain converts GraphQL model Privacy to domain Privacy
func privacyModelToDomain(p *model.Privacy) *domain.Privacy {
	if p == nil {
		return nil
	}
	var dp domain.Privacy
	switch *p {
	case model.PrivacyPrivate:
		dp = domain.PrivacyPrivate
	default:
		dp = domain.PrivacyPublic
	}
	return &dp
}

// reviewStatusDomainToModel converts domain ReviewStatus to GraphQL model ReviewStatus
func reviewStatusDomainToModel(s domain.ReviewStatus) model.ReviewStatus {
	switch s {
	case domain.ReviewStatusApproved:
		return model.ReviewStatusApproved
	case domain.ReviewStatusRejected:
		return model.ReviewStatusRejected
	default:
		return model.ReviewStatusPending
	}
}

// reviewStatusModelToDomain converts GraphQL model ReviewStatus to domain ReviewStatus
func reviewStatusModelToDomain(s *model.ReviewStatus) *domain.ReviewStatus {
	if s == nil {
		return nil
	}
	var ds domain.ReviewStatus
	switch *s {
	case model.ReviewStatusApproved:
		ds = domain.ReviewStatusApproved
	case model.ReviewStatusRejected:
		ds = domain.ReviewStatusRejected
	default:
		ds = domain.ReviewStatusPending
	}
	return &ds
}
