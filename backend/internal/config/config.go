package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	YouTube  YouTubeConfig  `json:"youtube"`
	Logging  LoggingConfig  `json:"logging"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password,omitempty"` // Will be overridden by env var
	SSLMode  string `json:"sslmode"`
}

// YouTubeConfig holds YouTube API configuration
type YouTubeConfig struct {
	APIKey string `json:"api_key"` // Will be overridden by env var

	// CacheTTLSeconds controls how long a fetched video's metadata is kept in
	// the in-memory YouTube response cache before it's re-fetched from the
	// API. Default is 6 hours — see YOUTUBE_API_CACHE_TTL_SECONDS in .env.example.
	CacheTTLSeconds int `json:"cache_ttl_seconds"`
}

// DefaultYouTubeCacheTTLSeconds is 6 hours, expressed in seconds.
const DefaultYouTubeCacheTTLSeconds = 6 * 60 * 60

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// Load reads configuration from file and environment variables.
// If the config file is missing, returns sensible defaults (production uses env vars).
func Load(configPath string) (*Config, error) {
	cfg := Config{
		Server:  ServerConfig{Port: 8080, Host: ""},
		YouTube: YouTubeConfig{CacheTTLSeconds: DefaultYouTubeCacheTTLSeconds},
	}

	// Read config file (optional in production where env vars provide all config)
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// No config file — use defaults + env vars
		} else {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
	} else {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with environment variables (for secrets)
	if dbPassword := os.Getenv("DATABASE_PASSWORD"); dbPassword != "" {
		cfg.Database.Password = dbPassword
	}

	if ytAPIKey := os.Getenv("YOUTUBE_API_KEY"); ytAPIKey != "" {
		cfg.YouTube.APIKey = ytAPIKey
	}

	// Unlike getEnvInt (security.go), 0 is a valid value here — it disables
	// the YouTube response cache entirely — so parse directly rather than
	// treating 0 as "unset". An unset or invalid value falls back to
	// whatever's already in cfg.YouTube.CacheTTLSeconds (config file value,
	// or DefaultYouTubeCacheTTLSeconds set above).
	if ttlStr := os.Getenv("YOUTUBE_API_CACHE_TTL_SECONDS"); ttlStr != "" {
		if v, err := strconv.Atoi(ttlStr); err == nil && v >= 0 {
			cfg.YouTube.CacheTTLSeconds = v
		}
	}

	return &cfg, nil
}

// GetAddr returns the server address in host:port format
func (c *ServerConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetDSN returns the PostgreSQL connection string (Data Source Name)
// Prefers DATABASE_URL env var if set (for hosted databases like Sevalla)
func (c *DatabaseConfig) GetDSN() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}
