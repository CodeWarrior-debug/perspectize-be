package postgres

// GORM PROTOTYPE — Domain ↔ GORM model mappers.
// These replace the current rowToDomain / toNullString / toNullInt64 helpers.
// With GORM's pointer-based nullability, mappers are simpler than sqlx's sql.Null* approach.

import (
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/domain"
)

// --- User mappers ---

func userModelToDomain(m *UserModel) *domain.User {
	return &domain.User{
		ID:        m.ID,
		Username:  m.Username,
		Email:     m.Email,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func userDomainToModel(u *domain.User) *UserModel {
	return &UserModel{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}

// --- Content mappers ---

func contentModelToDomain(m *ContentModel) *domain.Content {
	return &domain.Content{
		ID:          m.ID,
		Name:        m.Name,
		URL:         m.URL,
		ContentType: domain.ContentType(strings.ToUpper(m.ContentType)),
		Length:      m.Length,
		LengthUnits: m.LengthUnits,
		Response:    m.Response,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func contentDomainToModel(c *domain.Content) *ContentModel {
	return &ContentModel{
		ID:          c.ID,
		Name:        c.Name,
		URL:         c.URL,
		ContentType: strings.ToLower(string(c.ContentType)),
		Length:      c.Length,
		LengthUnits: c.LengthUnits,
		Response:    c.Response,
	}
}

// --- Perspective mappers ---

func perspectiveModelToDomain(m *PerspectiveModel) *domain.Perspective {
	p := &domain.Perspective{
		ID:        m.ID,
		Claim:     m.Claim,
		UserID:    m.UserID,
		ContentID: m.ContentID,
		Like:      m.Like,
		Quality:   m.Quality,
		Agreement: m.Agreement,
		Importance: m.Importance,
		Confidence: m.Confidence,
		Privacy:    domain.PrivacyPublic,
		Category:   m.Category,
		Description: m.Description,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}

	if m.Privacy != nil {
		p.Privacy = domain.Privacy(strings.ToUpper(*m.Privacy))
	}
	if m.ReviewStatus != nil {
		rs := domain.ReviewStatus(strings.ToUpper(*m.ReviewStatus))
		p.ReviewStatus = &rs
	}

	// Convert int64 array to int slice
	if len(m.Parts) > 0 {
		p.Parts = make([]int, len(m.Parts))
		for i, v := range m.Parts {
			p.Parts[i] = int(v)
		}
	}
	if len(m.Labels) > 0 {
		p.Labels = m.Labels
	}

	// Parse categorized ratings from JSON strings
	if len(m.CategorizedRatings) > 0 {
		p.CategorizedRatings = make([]domain.CategorizedRating, 0, len(m.CategorizedRatings))
		for _, jsonStr := range m.CategorizedRatings {
			var cr domain.CategorizedRating
			if err := json.Unmarshal([]byte(jsonStr), &cr); err != nil {
				slog.Warn("failed to parse categorized rating JSON", "error", err)
				continue
			}
			p.CategorizedRatings = append(p.CategorizedRatings, cr)
		}
	}

	return p
}

func perspectiveDomainToModel(p *domain.Perspective) *PerspectiveModel {
	m := &PerspectiveModel{
		ID:        p.ID,
		Claim:     p.Claim,
		UserID:    p.UserID,
		ContentID: p.ContentID,
		Like:      p.Like,
		Quality:   p.Quality,
		Agreement: p.Agreement,
		Importance: p.Importance,
		Confidence: p.Confidence,
		Category:   p.Category,
		Description: p.Description,
	}

	privStr := strings.ToLower(string(p.Privacy))
	m.Privacy = &privStr

	if p.ReviewStatus != nil {
		rs := strings.ToLower(string(*p.ReviewStatus))
		m.ReviewStatus = &rs
	}

	if len(p.Parts) > 0 {
		m.Parts = make(pq.Int64Array, len(p.Parts))
		for i, v := range p.Parts {
			m.Parts[i] = int64(v)
		}
	}
	if len(p.Labels) > 0 {
		m.Labels = p.Labels
	}

	// Marshal categorized ratings to JSON strings
	if len(p.CategorizedRatings) > 0 {
		m.CategorizedRatings = make(CategorizedRatingsGORM, len(p.CategorizedRatings))
		for i, cr := range p.CategorizedRatings {
			data, _ := json.Marshal(cr)
			m.CategorizedRatings[i] = string(data)
		}
	}

	return m
}
