package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/graphql/generated"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/graphql/resolvers"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/repositories/postgres"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/youtube"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/config"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/services"
	"github.com/CodeWarrior-debug/perspectize/backend/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") != "production" {
			slog.Warn(".env file not found", "hint", "set APP_ENV=production to suppress")
		}
	}

	// Load config (path from env or default)
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.example.json"
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate DATABASE_URL if set
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		if err := config.ValidateDatabaseURL(dbURL); err != nil {
			log.Fatalf("Invalid DATABASE_URL: %v", err)
		}
	}

	dsn := cfg.Database.GetDSN()

	// Mask credentials in log output
	if os.Getenv("DATABASE_URL") != "" {
		slog.Info("connecting to database using DATABASE_URL")
	} else {
		slog.Info("connecting to database", "host", cfg.Database.Host, "port", cfg.Database.Port, "name", cfg.Database.Name)
	}

	// Connect to database with configurable pool
	poolCfg := database.PoolConfigFromEnv()
	db, err := database.Connect(dsn, poolCfg)
	if err != nil {
		log.Fatalf("Failed to connect to database %s: %v", config.SanitizeDSN(dsn), err)
	}
	defer db.Close()

	// Test connection
	if err := database.Ping(context.Background(), db); err != nil {
		log.Fatalf("Database ping failed for %s: %v", config.SanitizeDSN(dsn), err)
	}

	slog.Info("successfully connected to database")

	// Quick query to verify
	var version string
	if err := db.Get(&version, "SELECT version()"); err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	slog.Info("PostgreSQL version", "version", version)

	// Validate YouTube API key
	if cfg.YouTube.APIKey == "" {
		slog.Warn("YOUTUBE_API_KEY is empty — YouTube metadata fetching will fail")
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
		slog.Info("shutting down gracefully")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("shutdown failed", "error", err)
		}
	}()

	slog.Info("server running", "addr", addr)
	if os.Getenv("APP_ENV") != "production" {
		slog.Info("GraphQL Playground available", "url", fmt.Sprintf("http://localhost%s/", addr))
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
