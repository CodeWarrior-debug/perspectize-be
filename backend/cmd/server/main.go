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
	gqltiming "github.com/CodeWarrior-debug/perspectize/backend/pkg/graphql"
	perfmw "github.com/CodeWarrior-debug/perspectize/backend/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	db, err := database.ConnectGORM(dsn, poolCfg)
	if err != nil {
		log.Fatalf("Failed to connect to database %s: %v", config.SanitizeDSN(dsn), err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Register slow query logger (logs queries >100ms)
	database.RegisterSlowQueryLogger(db)

	// Test connection
	if err := database.PingGORM(context.Background(), db); err != nil {
		log.Fatalf("Database ping failed for %s: %v", config.SanitizeDSN(dsn), err)
	}

	slog.Info("successfully connected to database")

	// Quick query to verify
	var version string
	if err := db.Raw("SELECT version()").Scan(&version).Error; err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	slog.Info("PostgreSQL version", "version", version)

	// Validate YouTube API key
	if cfg.YouTube.APIKey == "" {
		slog.Warn("YOUTUBE_API_KEY is empty — YouTube metadata fetching will fail")
	}

	// Initialize adapters
	youtubeClient := youtube.NewClient(cfg.YouTube.APIKey)
	contentRepo := postgres.NewGormContentRepository(db)
	userRepo := postgres.NewGormUserRepository(db)
	perspectiveRepo := postgres.NewGormPerspectiveRepository(db)

	// Initialize services
	contentService := services.NewContentService(contentRepo, youtubeClient)
	userService := services.NewUserService(userRepo, contentRepo, perspectiveRepo)
	perspectiveService := services.NewPerspectiveService(perspectiveRepo, userRepo)

	// Initialize GraphQL
	resolver := resolvers.NewResolver(contentService, userService, perspectiveService)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	srv.AroundOperations(gqltiming.OperationTimer())

	// Setup chi router
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(perfmw.RequestTimer)  // structured request timing (replaces chi Logger)
	r.Use(middleware.Recoverer) // panic recovery

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
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
	})

	// Health check — liveness probe (M-10)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Ready check — readiness probe with DB ping (M-10)
	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.PingContext(r.Context()) != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("not ready: database unreachable"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	})

	// GraphQL
	r.Handle("/graphql", srv)
	if os.Getenv("APP_ENV") != "production" {
		r.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
		r.Get("/debug/db-stats", database.StatsHandler(sqlDB))
	}

	// Start server with timeouts
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      r, // chi router
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
