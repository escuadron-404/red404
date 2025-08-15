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

	"github.com/escuadron-404/red404/backend/config"
	"github.com/escuadron-404/red404/backend/internal/handlers"
	"github.com/escuadron-404/red404/backend/internal/migration"
	"github.com/escuadron-404/red404/backend/internal/repositories"
	"github.com/escuadron-404/red404/backend/internal/routes"
	"github.com/escuadron-404/red404/backend/internal/services"
	"github.com/escuadron-404/red404/backend/pkg/database"
	"github.com/escuadron-404/red404/backend/pkg/middleware"
	"github.com/escuadron-404/red404/backend/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	dbConfig := &database.Config{
		DBHost:     cfg.DBHost,
		DBPort:     cfg.DBPort,
		DBUser:     cfg.DBUser,
		DBPassword: cfg.DBPassword,
		DBName:     cfg.DBName,
		DBSSLMode:  cfg.DBSSLMode,
	}

	db, err := database.NewDB(dbConfig)
	if err != nil {
		log.Printf("Failed to connect to database: %v\n", err)
		return
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db.Pool); err != nil {
		log.Printf("Failed to run migrations: %v\n", err)
		return
	}

	// Initialize utilities
	validate := validator.New()
	jwtUtil := utils.NewJWTUtil(cfg.JWTSecret, cfg.JWTExpirationHours)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db.Pool)

	// Initialize services
	userService := services.NewUserService(userRepo, validate)
	authService := services.NewAuthService(userRepo, validate, jwtUtil)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, validate)
	authHandler := handlers.NewAuthHandler(authService, validate)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtUtil)

	// Setup routes using the routes package
	mux := routes.SetupRoutes(userHandler, authHandler, authMiddleware)

	// Wrap mux with CORS
	//corsHandler := middleware.NewCORS().Handler(mux)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: mux,
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Listen for Ctrl+C and termination signals

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server failed to start: %v\n", err)
			// FIX: Instead of os.Exit(1), signal the main goroutine to shut down.
			quit <- syscall.SIGTERM
		}
	}()

	<-quit // Block until a signal is received
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // This defer will now correctly run

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
		// FIX: Replace os.Exit(1) with return.
		return
	}

	log.Println("Server exited")
}

func runMigrations(pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
	);`

	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	migrationService := migration.NewMigration(pool)
	if err := migrationService.Up(); err != nil {
		return err
	}

	return nil
}
