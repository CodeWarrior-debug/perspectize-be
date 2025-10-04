package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yourorg/perspectize-go/internal/config"
)

// TestLoad_RealConfigFile tests loading the actual config.example.json
func TestLoad_RealConfigFile(t *testing.T) {
	// Load the ACTUAL config file from the project
	configPath := "../../config/config.example.json"

	cfg, err := config.Load(configPath)
	assert.NoError(t, err, "Should load real config.example.json without errors")
	assert.NotNil(t, cfg)

	// Verify the actual structure and values in config.example.json
	assert.Equal(t, 8080, cfg.Server.Port, "Server port should match config.example.json")
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "testdb", cfg.Database.Name)
	assert.Equal(t, "testuser", cfg.Database.User)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
	assert.Equal(t, "", cfg.YouTube.APIKey, "API key should be empty in example config")
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
}

// TestLoad_RealConfigWithEnvOverrides tests that env vars override the real config
func TestLoad_RealConfigWithEnvOverrides(t *testing.T) {
	configPath := "../../config/config.example.json"

	// Set environment variables
	os.Setenv("DATABASE_PASSWORD", "secret123")
	os.Setenv("YOUTUBE_API_KEY", "yt_key_456")
	defer func() {
		os.Unsetenv("DATABASE_PASSWORD")
		os.Unsetenv("YOUTUBE_API_KEY")
	}()

	cfg, err := config.Load(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Verify environment variables override empty values from config.example.json
	assert.Equal(t, "secret123", cfg.Database.Password, "DATABASE_PASSWORD env var should override config")
	assert.Equal(t, "yt_key_456", cfg.YouTube.APIKey, "YOUTUBE_API_KEY env var should override config")
}

// TestLoad_InvalidPath tests error handling for missing file
func TestLoad_InvalidPath(t *testing.T) {
	cfg, err := config.Load("/nonexistent/path/config.json")
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to open config file")
}

// TestLoad_InvalidJSON tests error handling for malformed JSON
// Note: This uses a temp file since we don't want to commit invalid JSON to the repo
func TestLoad_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid.json")

	invalidJSON := `{ "server": { "port": invalid } }`
	err := os.WriteFile(configPath, []byte(invalidJSON), 0644)
	assert.NoError(t, err)

	cfg, err := config.Load(configPath)
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to parse config file")
}

// TestServerConfig_GetAddr tests the server address helper
func TestServerConfig_GetAddr(t *testing.T) {
	cfg := &config.ServerConfig{
		Host: "0.0.0.0",
		Port: 3000,
	}

	addr := cfg.GetAddr()
	assert.Equal(t, "0.0.0.0:3000", addr)
}

// TestDatabaseConfig_GetDSN tests the database connection string generation
func TestDatabaseConfig_GetDSN(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "db.example.com",
		Port:     5432,
		User:     "myuser",
		Password: "mypassword",
		Name:     "mydb",
		SSLMode:  "require",
	}

	dsn := cfg.GetDSN()
	expected := "host=db.example.com port=5432 user=myuser password=mypassword dbname=mydb sslmode=require"
	assert.Equal(t, expected, dsn)
}

// TestDatabaseConfig_GetDSN_FromRealConfig tests DSN generation with actual config values
func TestDatabaseConfig_GetDSN_FromRealConfig(t *testing.T) {
	configPath := "../../config/config.example.json"

	// Set password via env var (like production)
	os.Setenv("DATABASE_PASSWORD", "testpass")
	defer os.Unsetenv("DATABASE_PASSWORD")

	cfg, err := config.Load(configPath)
	assert.NoError(t, err)

	dsn := cfg.Database.GetDSN()

	// Verify DSN contains values from real config
	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "port=5432")
	assert.Contains(t, dsn, "user=testuser")
	assert.Contains(t, dsn, "password=testpass")
	assert.Contains(t, dsn, "dbname=testdb")
	assert.Contains(t, dsn, "sslmode=disable")

	// Verify it's a valid PostgreSQL connection string format
	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	assert.Equal(t, expected, dsn)
}
