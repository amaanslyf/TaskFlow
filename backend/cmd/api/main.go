// TaskFlow API — main entry point.
//
// This file wires everything together:
//   1. Load configuration (environment variables).
//   2. Connect to PostgreSQL.
//   3. Register HTTP routes.
//   4. Start the server with graceful shutdown.
//
// "Graceful shutdown" means that when you press Ctrl-C the
// server finishes any in-flight requests before exiting,
// instead of cutting them off mid-response.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/taskflow/backend/internal/config"
	"github.com/taskflow/backend/internal/database"
	"github.com/taskflow/backend/internal/handlers"
)

func main() {
	// ── 1. Configuration ────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}
	log.Printf("Starting TaskFlow API on port %s", cfg.Port)

	// ── 2. Database ─────────────────────────────────────────
	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to PostgreSQL")

	// ── 3. Routes ───────────────────────────────────────────
	mux := http.NewServeMux()

	taskHandler := handlers.NewTaskHandler(pool)

	// Health probes — used by Kubernetes later.
	mux.HandleFunc("/healthz", handlers.HandleHealthz)
	mux.Handle("/readyz", handlers.HandleReadyz(pool))

	// Task API endpoints.
	// We register two patterns:
	//   /api/v1/tasks   → list all / create
	//   /api/v1/tasks/  → get / update / delete by ID
	mux.HandleFunc("/api/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		// Exact match: no trailing path segments.
		if r.URL.Path == "/api/v1/tasks" || r.URL.Path == "/api/v1/tasks/" {
			// If the path is just /api/v1/tasks/ with nothing after the
			// slash, treat it the same as /api/v1/tasks (the collection).
			trimmed := strings.TrimSuffix(r.URL.Path, "/")
			if trimmed == "/api/v1/tasks" {
				taskHandler.HandleTasks(w, r)
				return
			}
		}
		taskHandler.HandleTasks(w, r)
	})
	mux.HandleFunc("/api/v1/tasks/", func(w http.ResponseWriter, r *http.Request) {
		// Anything under /api/v1/tasks/{something}
		taskHandler.HandleTaskByID(w, r)
	})

	// ── 4. Middleware stack ──────────────────────────────────
	var handler http.Handler = mux
	handler = handlers.CORSMiddleware(cfg.FrontendOrigin, handler)
	handler = handlers.LoggingMiddleware(handler)

	// ── 5. Server with timeouts ─────────────────────────────
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ── 6. Graceful shutdown ────────────────────────────────
	// Start the server in a goroutine so the main goroutine
	// can wait for a shutdown signal.
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until we receive SIGINT (Ctrl-C) or SIGTERM
	// (what Docker / Kubernetes sends to stop a container).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("Received %v — shutting down…", sig)

	// Give in-flight requests 30 seconds to finish.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("Server stopped cleanly")
}

