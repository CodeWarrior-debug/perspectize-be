package domain

import "time"

// User represents a user who can create perspectives
type User struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
