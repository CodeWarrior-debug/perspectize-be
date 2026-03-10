package domain

import "github.com/golang-jwt/jwt/v5"

// Claims represents the JWT claims stored in authentication tokens.
// Deprecated: Use AuthenticatedUser and the Clerk auth middleware instead.
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
