package config

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ValidateDatabaseURL validates the format of a DATABASE_URL connection string
func ValidateDatabaseURL(raw string) error {
	if raw == "" {
		return nil // Empty is valid (not using DATABASE_URL)
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check scheme
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return fmt.Errorf("invalid scheme %q, must be postgres:// or postgresql://", parsed.Scheme)
	}

	// Check hostname present
	if parsed.Host == "" {
		return fmt.Errorf("missing hostname")
	}

	// Check database name present (path without leading slash)
	dbName := strings.TrimPrefix(parsed.Path, "/")
	if dbName == "" {
		return fmt.Errorf("missing database name in path")
	}

	return nil
}

// SanitizeDSN removes password from DSN for safe logging
func SanitizeDSN(dsn string) string {
	// Try parsing as URL first
	parsed, err := url.Parse(dsn)
	if err == nil && parsed.Scheme != "" {
		// URL format: postgres://user:password@host/db
		if parsed.User != nil {
			username := parsed.User.Username()
			parsed.User = url.User(username) // Remove password
		}
		return parsed.String()
	}

	// Key-value format: host=x port=y user=u password=p dbname=d
	// Use regex to mask password
	re := regexp.MustCompile(`password=[^\s]+`)
	return re.ReplaceAllString(dsn, "password=***")
}
