package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	HandleHealthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", body["status"])
	}
}

func TestHandleHealthz_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	rec := httptest.NewRecorder()

	HandleHealthz(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestExtractID(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/api/v1/tasks/123e4567-e89b-12d3-a456-426614174000", "123e4567-e89b-12d3-a456-426614174000"},
		{"/api/v1/tasks/123e4567-e89b-12d3-a456-426614174000/", "123e4567-e89b-12d3-a456-426614174000"},
		{"/api/v1/tasks", ""},
		{"/api/v1", ""},
	}

	for _, tt := range tests {
		got := extractID(tt.path)
		if got != tt.expected {
			t.Errorf("extractID(%q) = %q, expected %q", tt.path, got, tt.expected)
		}
	}
}

func TestCORSMiddleware(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	cors := CORSMiddleware("http://localhost:5173", dummyHandler)

	// Test OPTIONS preflight
	reqOptions := httptest.NewRequest(http.MethodOptions, "/api/v1/tasks", nil)
	recOptions := httptest.NewRecorder()
	cors.ServeHTTP(recOptions, reqOptions)

	if recOptions.Code != http.StatusNoContent {
		t.Errorf("expected preflight 204, got %d", recOptions.Code)
	}
	if recOptions.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Errorf("expected CORS origin header http://localhost:5173, got %s", recOptions.Header().Get("Access-Control-Allow-Origin"))
	}

	// Test GET request
	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	recGet := httptest.NewRecorder()
	cors.ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusOK {
		t.Errorf("expected GET 200, got %d", recGet.Code)
	}
	if recGet.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Errorf("expected CORS origin header http://localhost:5173, got %s", recGet.Header().Get("Access-Control-Allow-Origin"))
	}
}

