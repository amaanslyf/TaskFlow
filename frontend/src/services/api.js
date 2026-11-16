// ── API Service ──────────────────────────────────────────
// This module centralizes every HTTP call to the Go backend.
//
// Why a separate file?
//   • One place to change if the API URL or shape changes.
//   • Components stay focused on rendering.
//   • Easy to swap with a mock for testing later.

const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

/**
 * Small wrapper around fetch that:
 *   1. Prepends the API base URL.
 *   2. Sets JSON headers.
 *   3. Throws on non-2xx responses with the server's error message.
 */
async function request(path, options = {}) {
  const url = `${API_BASE}${path}`;

  const config = {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  };

  const response = await fetch(url, config);

  // Try to parse the body even on errors — the Go backend
  // always returns JSON error messages.
  const data = await response.json().catch(() => null);

  if (!response.ok) {
    const message = data?.error || `Request failed with status ${response.status}`;
    throw new Error(message);
  }

  return data;
}

// ── Public API ───────────────────────────────────────────

export function getTasks() {
  return request('/api/v1/tasks');
}

export function getTask(id) {
  return request(`/api/v1/tasks/${id}`);
}

export function createTask(task) {
  return request('/api/v1/tasks', {
    method: 'POST',
    body: JSON.stringify(task),
  });
}

export function updateTask(id, task) {
  return request(`/api/v1/tasks/${id}`, {
    method: 'PUT',
    body: JSON.stringify(task),
  });
}

export function deleteTask(id) {
  return request(`/api/v1/tasks/${id}`, {
    method: 'DELETE',
  });
}

