package postgres

// GORM PROTOTYPE — Side-by-side comparison with current sqlx implementation.
// This file demonstrates the hex-clean GORM pattern: separate GORM model structs
// that live in the adapter layer, keeping domain models free of ORM tags.

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// UserModel is the GORM persistence model for users.
// Domain model (domain.User) has no ORM tags — this maps between them.
type UserModel struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (UserModel) TableName() string { return "users" }

// ContentModel is the GORM persistence model for content.
type ContentModel struct {
	ID          int             `gorm:"primaryKey;autoIncrement"`
	Name        string          `gorm:"not null"`
	URL         *string         `gorm:"uniqueIndex"`
	ContentType string          `gorm:"column:content_type;not null"`
	Length      *int            `gorm:"column:length"`
	LengthUnits *string        `gorm:"column:length_units"`
	Response    json.RawMessage `gorm:"type:jsonb"`
	CreatedAt   time.Time       `gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime"`
}

func (ContentModel) TableName() string { return "content" }

// PerspectiveModel is the GORM persistence model for perspectives.
type PerspectiveModel struct {
	ID                 int              `gorm:"primaryKey;autoIncrement"`
	Claim              string           `gorm:"not null;size:255"`
	UserID             int              `gorm:"not null"`
	ContentID          *int             `gorm:"column:content_id"`
	Like               *string          `gorm:"column:like"`
	Quality            *int             `gorm:"column:quality"`
	Agreement          *int             `gorm:"column:agreement"`
	Importance         *int             `gorm:"column:importance"`
	Confidence         *int             `gorm:"column:confidence"`
	Privacy            *string          `gorm:"column:privacy"`
	Parts              pq.Int64Array    `gorm:"type:integer[]"`
	Category           *string          `gorm:"column:category"`
	Labels             pq.StringArray   `gorm:"type:text[]"`
	Description        *string          `gorm:"column:description"`
	ReviewStatus       *string          `gorm:"column:review_status"`
	CategorizedRatings CategorizedRatingsGORM `gorm:"type:jsonb[];column:categorized_ratings"`
	CreatedAt          time.Time        `gorm:"autoCreateTime"`
	UpdatedAt          time.Time        `gorm:"autoUpdateTime"`
}

func (PerspectiveModel) TableName() string { return "perspectives" }

// CategorizedRatingsGORM handles PostgreSQL jsonb[] serialization for GORM.
type CategorizedRatingsGORM []string

func (c *CategorizedRatingsGORM) Scan(src interface{}) error {
	if src == nil {
		*c = nil
		return nil
	}
	var arr pq.StringArray
	if err := arr.Scan(src); err != nil {
		return err
	}
	*c = CategorizedRatingsGORM(arr)
	return nil
}

func (c CategorizedRatingsGORM) Value() (driver.Value, error) {
	if len(c) == 0 {
		return nil, nil
	}
	return pq.StringArray(c).Value()
}
