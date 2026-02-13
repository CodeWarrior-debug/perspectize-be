package domain

import (
	"encoding/json"
	"time"
)

// ContentType represents the type of content
type ContentType string

const (
	ContentTypeYouTube ContentType = "YOUTUBE"
)

// Content represents a media item that users create perspectives on
type Content struct {
	ID          int
	Name        string
	URL         *string
	ContentType ContentType
	Length      *int
	LengthUnits *string
	Response    json.RawMessage
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
