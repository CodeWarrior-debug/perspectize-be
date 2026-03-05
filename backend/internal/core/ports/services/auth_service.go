package services

import "github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"

// AuthService defines the interface for JWT authentication operations.
type AuthService interface {
	// GenerateAccessToken creates a signed JWT for the given user.
	GenerateAccessToken(userID int, email string) (string, error)

	// ValidateToken parses and validates a JWT, returning the claims if valid.
	ValidateToken(tokenString string) (*domain.Claims, error)
}
