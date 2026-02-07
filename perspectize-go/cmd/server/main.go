package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/adapters/graphql/generated"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/adapters/graphql/resolvers"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/adapters/repositories/postgres"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/adapters/youtube"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/config"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/services"
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/pkg/database"
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

	// Initialize adapters
	youtubeClient := youtube.NewClient(cfg.YouTube.APIKey)
	contentRepo := postgres.NewContentRepository(db)
	userRepo := postgres.NewUserRepository(db)
	perspectiveRepo := postgres.NewPerspectiveRepository(db)

	// Initialize services
	contentService := services.NewContentService(contentRepo, youtubeClient)
	userService := services.NewUserService(userRepo)
	perspectiveService := services.NewPerspectiveService(perspectiveRepo, userRepo)

	// Initialize GraphQL
	resolver := resolvers.NewResolver(contentService, userService, perspectiveService)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	// Setup HTTP routes
	http.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	http.Handle("/graphql", srv)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server running at http://localhost%s", addr)
	log.Printf("GraphQL Playground available at http://localhost%s/", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
