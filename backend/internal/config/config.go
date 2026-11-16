// Package config reads environment variables that control
// how the backend behaves.  During development a .env file
// is loaded automatically; in Docker / Kubernetes the real
// environment variables take over and the .env file is not
// needed.
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds every value the application needs at startup.
type Config struct {
	DatabaseURL    string // PostgreSQL connection string
	Port           string // HTTP port (e.g. "8080")
	FrontendOrigin string // Allowed CORS origin (e.g. "http://localhost:5173")
}

// Load reads configuration from the environment.
// It tries to load a .env file first (handy for local
// development), but never fails if the file is absent —
// Docker and Kubernetes will supply real environment
// variables instead.
func Load() (*Config, error) {
	// Best-effort: load .env if it exists.
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // sensible default
	}

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:5173" // Vite default
	}

	return &Config{
		DatabaseURL:    dbURL,
		Port:           port,
		FrontendOrigin: frontendOrigin,
	}, nil
}

