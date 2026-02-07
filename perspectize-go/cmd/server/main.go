package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// Load .env file
	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") != "production" {
			log.Println("Warning: .env file not found (set APP_ENV=production to suppress)")
		}
	}

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

	// Validate YouTube API key
	if cfg.YouTube.APIKey == "" {
		log.Println("Warning: YOUTUBE_API_KEY is empty — YouTube metadata fetching will fail")
	}

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

	// CORS middleware for frontend dev server
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	// Setup HTTP routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	if os.Getenv("APP_ENV") != "production" {
		http.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	}
	http.Handle("/graphql", corsHandler(srv))

	// Start server with timeouts
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}()

	log.Printf("Server running at http://localhost%s", addr)
	if os.Getenv("APP_ENV") != "production" {
		log.Printf("GraphQL Playground available at http://localhost%s/", addr)
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
