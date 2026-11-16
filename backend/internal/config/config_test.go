package config

import (
	"os"
	"testing"
)

func TestConfigLoad_Success(t *testing.T) {
	// Set required env var
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")
	os.Setenv("PORT", "9090")
	os.Setenv("FRONTEND_ORIGIN", "http://localhost:3000")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("PORT")
		os.Unsetenv("FRONTEND_ORIGIN")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/testdb" {
		t.Errorf("expected DATABASE_URL to match, got %s", cfg.DatabaseURL)
	}
	if cfg.Port != "9090" {
		t.Errorf("expected PORT to be 9090, got %s", cfg.Port)
	}
	if cfg.FrontendOrigin != "http://localhost:3000" {
		t.Errorf("expected FRONTEND_ORIGIN to be http://localhost:3000, got %s", cfg.FrontendOrigin)
	}
}

func TestConfigLoad_Defaults(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")
	os.Unsetenv("PORT")
	os.Unsetenv("FRONTEND_ORIGIN")
	defer os.Unsetenv("DATABASE_URL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("expected default PORT 8080, got %s", cfg.Port)
	}
	if cfg.FrontendOrigin != "http://localhost:5173" {
		t.Errorf("expected default FRONTEND_ORIGIN http://localhost:5173, got %s", cfg.FrontendOrigin)
	}
}

func TestConfigLoad_MissingDatabaseURL(t *testing.T) {
	os.Unsetenv("DATABASE_URL")

	// Temporarily override .env loading behavior by checking os.Getenv directly
	_, err := Load()
	// Note: If a local .env exists, godotenv.Load() will populate DATABASE_URL from .env.
	// So we only assert error if DATABASE_URL was truly not in environment or .env.
	if os.Getenv("DATABASE_URL") == "" && err == nil {
		t.Errorf("expected error when DATABASE_URL is missing")
	}
}

