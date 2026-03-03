package config

import (
	"log"
	"log/slog"
	"os"
	"strconv"
)

// Security holds JWT and authentication configuration.
type Security struct {
	JWTSecret          string
	AccessTokenMinutes int
}

// LoadSecurity reads security configuration from environment variables.
// In production, JWT_SECRET must be set and at least 32 bytes.
func LoadSecurity() Security {
	secret := os.Getenv("JWT_SECRET")
	minutes := 15 // default

	if m := os.Getenv("ACCESS_TOKEN_MINUTES"); m != "" {
		if v, err := strconv.Atoi(m); err == nil && v > 0 {
			minutes = v
		}
	}

	// Fail fast in production if JWT_SECRET is missing or too short
	if os.Getenv("APP_ENV") == "production" {
		if secret == "" {
			log.Fatal("JWT_SECRET is required in production")
		}
		if len(secret) < 32 {
			log.Fatal("JWT_SECRET must be at least 32 bytes in production")
		}
	} else if secret == "" {
		slog.Warn("JWT_SECRET not set — using insecure default for development")
		secret = "dev-only-insecure-jwt-secret-key-32b"
	}

	return Security{
		JWTSecret:          secret,
		AccessTokenMinutes: minutes,
	}
}
