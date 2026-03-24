package services

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles JWT token generation and validation.
type AuthService struct {
	jwtSecret      []byte
	accessTokenTTL time.Duration
}

// NewAuthService creates a new AuthService with the given secret and token TTL.
// Logs a warning if the secret is shorter than 32 bytes.
func NewAuthService(jwtSecret []byte, accessTokenTTL time.Duration) *AuthService {
	if len(jwtSecret) < 32 {
		slog.Warn("JWT secret is shorter than 32 bytes — this is insecure for production",
			"length", len(jwtSecret))
	}
	return &AuthService{
		jwtSecret:      jwtSecret,
		accessTokenTTL: accessTokenTTL,
	}
}

// GenerateAccessToken creates a signed JWT with the user's ID and email.
func (s *AuthService) GenerateAccessToken(userID int, email string) (string, error) {
	now := time.Now()
	claims := domain.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "perspectize",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken parses and validates a JWT string, returning the claims if valid.
func (s *AuthService) ValidateToken(tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
