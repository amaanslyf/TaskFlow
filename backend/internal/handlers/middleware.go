// Package handlers – middleware.go
// Contains reusable HTTP middleware for the TaskFlow API.
package handlers

import (
	"log"
	"net/http"
	"time"
)

// CORSMiddleware adds the Cross-Origin Resource Sharing
// headers that let the React frontend (running on a different
// port) talk to this API.
//
// Without CORS the browser blocks the request because the
// frontend origin (e.g. localhost:5173) differs from the API
// origin (e.g. localhost:8080).
func CORSMiddleware(allowedOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods",
			"GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Authorization")

		// Preflight requests (OPTIONS) are the browser's way of
		// asking "is this cross-origin call allowed?" — we just
		// say yes and return immediately.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs every incoming request with its
// method, path, and how long it took.  Very useful for
// debugging and later for understanding what Prometheus
// metrics look like.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

