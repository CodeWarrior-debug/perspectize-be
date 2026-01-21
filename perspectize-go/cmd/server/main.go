package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yourorg/perspectize-go/internal/config"
	"github.com/yourorg/perspectize-go/pkg/database"
)

func main() {
	// Load .env file (optional - won't error if missing)
	_ = godotenv.Load()

	// Load config
	cfg, err := config.Load("config/config.example.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dsn := cfg.Database.GetDSN()

	// Mask credentials in log output
	if os.Getenv("DATABASE_URL") != "" {
		log.Println("Connecting to database using DATABASE_URL...")
	} else {
		log.Printf("Connecting to database at %s:%d/%s...", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	}

	// Connect to database
	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := database.Ping(context.Background(), db); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Successfully connected to database!")

	// Quick query to verify
	var version string
	if err := db.Get(&version, "SELECT version()"); err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	log.Printf("PostgreSQL version: %s", version)
}
