// Package tasks contains the data model and the "repository"
// — the code that talks to PostgreSQL for task-related work.
//
// Separating this from the HTTP handlers keeps each piece
// small and testable on its own.
package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ── Data model ──────────────────────────────────────────────

// Task is the Go representation of one row in the tasks table.
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTaskInput carries the data the user sends when
// creating a new task.
type CreateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// UpdateTaskInput carries the data the user sends when
// editing an existing task.
type UpdateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

// ValidStatuses is used to validate the status field.
var ValidStatuses = map[string]bool{
	"todo":        true,
	"in_progress": true,
	"completed":   true,
}

// ── Repository ──────────────────────────────────────────────

// Repository provides every database operation for tasks.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a Repository backed by the given pool.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// GetAll returns every task, ordered newest-first.
func (r *Repository) GetAll(ctx context.Context) ([]Task, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, title, description, status, created_at, updated_at
		 FROM tasks
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("querying tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description,
			&t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning task row: %w", err)
		}
		tasks = append(tasks, t)
	}

	// Return an empty slice (not nil) so JSON serialization
	// produces [] instead of null.
	if tasks == nil {
		tasks = []Task{}
	}
	return tasks, nil
}

// GetByID returns a single task or pgx.ErrNoRows.
func (r *Repository) GetByID(ctx context.Context, id string) (*Task, error) {
	var t Task
	err := r.pool.QueryRow(ctx,
		`SELECT id, title, description, status, created_at, updated_at
		 FROM tasks
		 WHERE id = $1`, id).
		Scan(&t.ID, &t.Title, &t.Description,
			&t.Status, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("querying task %s: %w", id, err)
	}
	return &t, nil
}

// Create inserts a new task and returns it.
func (r *Repository) Create(ctx context.Context, input CreateTaskInput) (*Task, error) {
	status := input.Status
	if status == "" {
		status = "todo"
	}

	var t Task
	err := r.pool.QueryRow(ctx,
		`INSERT INTO tasks (title, description, status)
		 VALUES ($1, $2, $3)
		 RETURNING id, title, description, status, created_at, updated_at`,
		input.Title, input.Description, status).
		Scan(&t.ID, &t.Title, &t.Description,
			&t.Status, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("inserting task: %w", err)
	}
	return &t, nil
}

// Update modifies an existing task and returns the updated row.
func (r *Repository) Update(ctx context.Context, id string, input UpdateTaskInput) (*Task, error) {
	var t Task
	err := r.pool.QueryRow(ctx,
		`UPDATE tasks
		 SET title       = COALESCE($1, title),
		     description = COALESCE($2, description),
		     status      = COALESCE($3, status),
		     updated_at  = NOW()
		 WHERE id = $4
		 RETURNING id, title, description, status, created_at, updated_at`,
		input.Title, input.Description, input.Status, id).
		Scan(&t.ID, &t.Title, &t.Description,
			&t.Status, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("updating task %s: %w", id, err)
	}
	return &t, nil
}

// Delete removes a task by ID. Returns true if a row was
// actually deleted, false if the ID did not exist.
func (r *Repository) Delete(ctx context.Context, id string) (bool, error) {
	result, err := r.pool.Exec(ctx,
		`DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return false, fmt.Errorf("deleting task %s: %w", id, err)
	}
	return result.RowsAffected() > 0, nil
}

