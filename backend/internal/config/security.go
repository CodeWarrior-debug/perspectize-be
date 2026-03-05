package config

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// Security holds JWT, authentication, rate limiting, and CORS configuration.
type Security struct {
	JWTSecret          string
	AccessTokenMinutes int
	ClerkSecretKey     string
	RateLimitPerMin    int      // Default 100
	CORSOrigins        []string // Explicit origins
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

	clerkKey := os.Getenv("CLERK_SECRET_KEY")
	if os.Getenv("APP_ENV") == "production" && clerkKey == "" {
		log.Fatal("CLERK_SECRET_KEY is required in production")
	} else if clerkKey == "" {
		slog.Warn("CLERK_SECRET_KEY not set — Clerk auth will not work")
	}

	return Security{
		JWTSecret:          secret,
		AccessTokenMinutes: minutes,
		ClerkSecretKey:     clerkKey,
		RateLimitPerMin:    getEnvInt("RATE_LIMIT_PER_MIN", 100),
		CORSOrigins:        getEnvStringSlice("CORS_ORIGINS", []string{"*"}),
	}
}

// getEnvInt reads an integer from an environment variable with a default.
func getEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	v, err := strconv.Atoi(val)
	if err != nil || v <= 0 {
		return defaultValue
	}
	return v
}

// getEnvStringSlice reads a comma-separated string list from an environment variable.
func getEnvStringSlice(key string, defaultValue []string) []string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return strings.Split(val, ",")
}
