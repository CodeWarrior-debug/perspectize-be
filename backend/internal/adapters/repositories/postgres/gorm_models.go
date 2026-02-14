package postgres

import (
	"encoding/json"
	"time"
)

// UserModel is the GORM persistence model for users table
type UserModel struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for UserModel
func (UserModel) TableName() string {
	return "users"
}

// ContentModel is the GORM persistence model for content table
type ContentModel struct {
	ID          int             `gorm:"primaryKey;autoIncrement"`
	Name        string          `gorm:"not null"`
	URL         *string         `gorm:"uniqueIndex"`
	ContentType string          `gorm:"column:content_type;not null"`
	Length      *int            `gorm:""`
	LengthUnits *string         `gorm:""`
	Response    json.RawMessage `gorm:"type:jsonb"`

	// Dummy fields for gorm-cursor-paginator sort key validation.
	// These are NOT database columns — SQLRepr provides the actual SQL.
	// The gorm:"-" tag tells GORM to ignore them for queries/migrations.
	ViewCount   int64  `gorm:"-"`
	LikeCount   int64  `gorm:"-"`
	PublishedAt string `gorm:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for ContentModel
func (ContentModel) TableName() string {
	return "content"
}

// PerspectiveModel is the GORM persistence model for perspectives table
type PerspectiveModel struct {
	ID                 int         `gorm:"primaryKey;autoIncrement"`
	Claim              string      `gorm:"not null;size:255"`
	UserID             int         `gorm:"not null"`
	ContentID          *int        `gorm:""`
	Like               *string     `gorm:"column:like"`
	Quality            *int        `gorm:""`
	Agreement          *int        `gorm:""`
	Importance         *int        `gorm:""`
	Confidence         *int        `gorm:""`
	Privacy            *string     `gorm:""`
	Parts              Int64Array  `gorm:"type:integer[]"`
	Category           *string     `gorm:""`
	Labels             StringArray `gorm:"type:text[]"`
	Description        *string     `gorm:""`
	ReviewStatus       *string     `gorm:""`
	CategorizedRatings JSONBArray  `gorm:"type:jsonb[];column:categorized_ratings"`
	CreatedAt          time.Time   `gorm:"autoCreateTime"`
	UpdatedAt          time.Time   `gorm:"autoUpdateTime"`
}

// TableName returns the table name for PerspectiveModel
func (PerspectiveModel) TableName() string {
	return "perspectives"
}
