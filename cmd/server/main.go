package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/azzamdhx/moneybro/backend/internal/config"
	"github.com/azzamdhx/moneybro/backend/internal/cron"
	"github.com/azzamdhx/moneybro/backend/internal/database"
	"github.com/azzamdhx/moneybro/backend/internal/graph"
	"github.com/azzamdhx/moneybro/backend/internal/middleware"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
	"github.com/azzamdhx/moneybro/backend/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	rdb, err := database.NewRedis(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	repos := repository.NewRepositories(db)

	svc := services.NewServices(services.Config{
		Repos:             repos,
		Redis:             rdb,
		JWTSecret:         cfg.JWTSecret,
		ResendAPIKey:      cfg.ResendAPIKey,
		FrontendURL:       cfg.FrontendURL,
		EmailTemplatesDir: cfg.EmailTemplatesDir,
	})

	cronScheduler := cron.NewScheduler(svc.Notification)
	cronScheduler.Start()
	defer cronScheduler.Stop()

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.CORS(cfg.FrontendURL))
	r.Use(middleware.RateLimit(100, time.Minute))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: graph.NewResolver(svc),
	}))

	r.Handle("/graphql", middleware.Auth(cfg.JWTSecret)(srv))

	if cfg.Env == "development" {
		r.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))
	}

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
