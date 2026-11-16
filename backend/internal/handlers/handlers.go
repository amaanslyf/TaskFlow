// Package handlers contains the HTTP handlers (controllers)
// for the TaskFlow REST API.
//
// Each handler:
//   1. Reads and validates the HTTP request.
//   2. Calls the repository to talk to PostgreSQL.
//   3. Writes a JSON response with the correct status code.
//
// This separation means the handler never writes SQL and the
// repository never touches HTTP — each layer has one job.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/taskflow/backend/internal/tasks"
)

// uuidRegex validates UUID format before we hit the database.
var uuidRegex = regexp.MustCompile(
	`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// TaskHandler groups the HTTP handlers and holds a reference
// to the repository they all share.
type TaskHandler struct {
	repo *tasks.Repository
}

// NewTaskHandler creates a TaskHandler backed by the given pool.
func NewTaskHandler(pool *pgxpool.Pool) *TaskHandler {
	return &TaskHandler{
		repo: tasks.NewRepository(pool),
	}
}

// ── Helpers ─────────────────────────────────────────────────

// jsonResponse writes a JSON body with the given status code.
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("ERROR encoding response: %v", err)
		}
	}
}

// jsonError writes a JSON error body.
func jsonError(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}

// extractID pulls the task ID from the URL path.
// For Go 1.22+ net/http we can use the new PathValue, but to
// keep compatibility wider we do a simple manual parse.
func extractID(path string) string {
	// Expected: /api/v1/tasks/{id}
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) < 5 {
		return ""
	}
	return parts[4]
}

// ── Handlers ────────────────────────────────────────────────

// HandleTasks routes /api/v1/tasks to the right method.
func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// HandleTaskByID routes /api/v1/tasks/{id} to the right method.
func (h *TaskHandler) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" || !uuidRegex.MatchString(id) {
		jsonError(w, http.StatusBadRequest, "invalid or missing task ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, id)
	case http.MethodPut:
		h.updateTask(w, r, id)
	case http.MethodDelete:
		h.deleteTask(w, r, id)
	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// ── Individual handlers ─────────────────────────────────────

func (h *TaskHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	taskList, err := h.repo.GetAll(r.Context())
	if err != nil {
		log.Printf("ERROR listing tasks: %v", err)
		jsonError(w, http.StatusInternalServerError, "failed to retrieve tasks")
		return
	}
	jsonResponse(w, http.StatusOK, taskList)
}

func (h *TaskHandler) getTask(w http.ResponseWriter, r *http.Request, id string) {
	task, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("ERROR getting task %s: %v", id, err)
		jsonError(w, http.StatusInternalServerError, "failed to retrieve task")
		return
	}
	if task == nil {
		jsonError(w, http.StatusNotFound, "task not found")
		return
	}
	jsonResponse(w, http.StatusOK, task)
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var input tasks.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// Validate required field.
	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		jsonError(w, http.StatusBadRequest, "title is required")
		return
	}

	// Validate status if provided.
	if input.Status != "" && !tasks.ValidStatuses[input.Status] {
		jsonError(w, http.StatusBadRequest,
			"status must be one of: todo, in_progress, completed")
		return
	}

	task, err := h.repo.Create(r.Context(), input)
	if err != nil {
		log.Printf("ERROR creating task: %v", err)
		jsonError(w, http.StatusInternalServerError, "failed to create task")
		return
	}
	jsonResponse(w, http.StatusCreated, task)
}

func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request, id string) {
	var input tasks.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// Validate title if provided.
	if input.Title != nil {
		trimmed := strings.TrimSpace(*input.Title)
		if trimmed == "" {
			jsonError(w, http.StatusBadRequest, "title cannot be empty")
			return
		}
		input.Title = &trimmed
	}

	// Validate status if provided.
	if input.Status != nil && !tasks.ValidStatuses[*input.Status] {
		jsonError(w, http.StatusBadRequest,
			"status must be one of: todo, in_progress, completed")
		return
	}

	task, err := h.repo.Update(r.Context(), id, input)
	if err != nil {
		log.Printf("ERROR updating task %s: %v", id, err)
		jsonError(w, http.StatusInternalServerError, "failed to update task")
		return
	}
	if task == nil {
		jsonError(w, http.StatusNotFound, "task not found")
		return
	}
	jsonResponse(w, http.StatusOK, task)
}

func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request, id string) {
	deleted, err := h.repo.Delete(r.Context(), id)
	if err != nil {
		log.Printf("ERROR deleting task %s: %v", id, err)
		jsonError(w, http.StatusInternalServerError, "failed to delete task")
		return
	}
	if !deleted {
		jsonError(w, http.StatusNotFound, "task not found")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "task deleted"})
}

// ── Health probes ───────────────────────────────────────────

// HandleHealthz responds to GET /healthz.
// It always returns 200 — it only checks that the process is
// running and can serve HTTP.
func HandleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

// HandleReadyz responds to GET /readyz.
// It pings PostgreSQL and returns 200 only if the database is
// reachable — this tells Kubernetes whether the backend is
// truly ready to accept traffic.
func HandleReadyz(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := pool.Ping(r.Context()); err != nil {
			log.Printf("readyz: database ping failed: %v", err)
			jsonError(w, http.StatusServiceUnavailable,
				"database connection failed")
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}

