package postgres

import (
	"encoding/json"
	"strings"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
)

// userModelToDomain converts a GORM UserModel to domain.User
func userModelToDomain(m *UserModel) *domain.User {
	if m == nil {
		return nil
	}
	return &domain.User{
		ID:        m.ID,
		Username:  m.Username,
		Email:     m.Email,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// userDomainToModel converts a domain.User to GORM UserModel
func userDomainToModel(u *domain.User) *UserModel {
	if u == nil {
		return nil
	}
	return &UserModel{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		// CreatedAt and UpdatedAt are managed by GORM
	}
}

// contentModelToDomain converts a GORM ContentModel to domain.Content
func contentModelToDomain(m *ContentModel) *domain.Content {
	if m == nil {
		return nil
	}
	return &domain.Content{
		ID:            m.ID,
		Name:          m.Name,
		URL:           m.URL,
		ContentType:   domain.ContentType(strings.ToUpper(m.ContentType)),
		AddedByUserID: m.AddedByUserID,
		Length:        m.Length,
		LengthUnits:   m.LengthUnits,
		Response:      m.Response,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// contentDomainToModel converts a domain.Content to GORM ContentModel
func contentDomainToModel(c *domain.Content) *ContentModel {
	if c == nil {
		return nil
	}
	return &ContentModel{
		ID:            c.ID,
		Name:          c.Name,
		URL:           c.URL,
		ContentType:   strings.ToLower(string(c.ContentType)),
		AddedByUserID: c.AddedByUserID,
		Length:        c.Length,
		LengthUnits:   c.LengthUnits,
		Response:      c.Response,
		// CreatedAt and UpdatedAt are managed by GORM
	}
}

// perspectiveModelToDomain converts a GORM PerspectiveModel to domain.Perspective
func perspectiveModelToDomain(m *PerspectiveModel) *domain.Perspective {
	if m == nil {
		return nil
	}

	p := &domain.Perspective{
		ID:          m.ID,
		UserID:      m.UserID,
		ContentID:   m.ContentID,
		Like:        m.Like,
		Quality:     m.Quality,
		Agreement:   m.Agreement,
		Importance:  m.Importance,
		Confidence:  m.Confidence,
		Category:    m.Category,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}

	// Privacy: default to PUBLIC if nil
	if m.Privacy != nil {
		p.Privacy = domain.Privacy(strings.ToUpper(*m.Privacy))
	} else {
		p.Privacy = domain.PrivacyPublic
	}

	// ReviewStatus: convert pointer with ToUpper
	if m.ReviewStatus != nil {
		rs := domain.ReviewStatus(strings.ToUpper(*m.ReviewStatus))
		p.ReviewStatus = &rs
	}

	// Parts: convert int64 to int
	if len(m.Parts) > 0 {
		p.Parts = make([]int, len(m.Parts))
		for i, v := range m.Parts {
			p.Parts[i] = int(v)
		}
	}

	// Labels: direct copy
	if len(m.Labels) > 0 {
		p.Labels = m.Labels
	}

	// CategorizedRatings: unmarshal from JSONBArray strings
	if len(m.CategorizedRatings) > 0 {
		p.CategorizedRatings = make([]domain.CategorizedRating, 0, len(m.CategorizedRatings))
		for _, jsonStr := range m.CategorizedRatings {
			var cr domain.CategorizedRating
			if err := json.Unmarshal([]byte(jsonStr), &cr); err != nil {
				// Skip invalid JSON - same behavior as sqlx implementation
				continue
			}
			p.CategorizedRatings = append(p.CategorizedRatings, cr)
		}
	}

	return p
}

// perspectiveDomainToModel converts a domain.Perspective to GORM PerspectiveModel
func perspectiveDomainToModel(p *domain.Perspective) *PerspectiveModel {
	if p == nil {
		return nil
	}

	m := &PerspectiveModel{
		ID:          p.ID,
		UserID:      p.UserID,
		ContentID:   p.ContentID,
		Like:        p.Like,
		Quality:     p.Quality,
		Agreement:   p.Agreement,
		Importance:  p.Importance,
		Confidence:  p.Confidence,
		Category:    p.Category,
		Description: p.Description,
	}

	// Privacy: ToLower
	privacy := strings.ToLower(string(p.Privacy))
	m.Privacy = &privacy

	// ReviewStatus: ToLower pointer
	if p.ReviewStatus != nil {
		rs := strings.ToLower(string(*p.ReviewStatus))
		m.ReviewStatus = &rs
	}

	// Parts: convert int to int64
	if len(p.Parts) > 0 {
		m.Parts = make(Int64Array, len(p.Parts))
		for i, v := range p.Parts {
			m.Parts[i] = int64(v)
		}
	}

	// Labels: direct copy
	if len(p.Labels) > 0 {
		m.Labels = p.Labels
	}

	// CategorizedRatings: marshal to JSONBArray strings
	if len(p.CategorizedRatings) > 0 {
		m.CategorizedRatings = make(JSONBArray, len(p.CategorizedRatings))
		for i, cr := range p.CategorizedRatings {
			data, err := json.Marshal(cr)
			if err != nil {
				// Skip invalid data - same as sqlx implementation
				continue
			}
			m.CategorizedRatings[i] = string(data)
		}
	}

	return m
}
