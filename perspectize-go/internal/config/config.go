package config

import (
	"encoding/json"
	"fmt"
	"os"
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
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// Load reads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	// Read config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Parse JSON
	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override with environment variables (for secrets)
	if dbPassword := os.Getenv("DATABASE_PASSWORD"); dbPassword != "" {
		cfg.Database.Password = dbPassword
	}

	if ytAPIKey := os.Getenv("YOUTUBE_API_KEY"); ytAPIKey != "" {
		cfg.YouTube.APIKey = ytAPIKey
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
