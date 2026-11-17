// ── API Service ──────────────────────────────────────────
// This module centralizes every HTTP call to the Go backend.
//
// Why a separate file?
//   • One place to change if the API URL or shape changes.
//   • Components stay focused on rendering.
//   • Easy to swap with a mock for testing later.

// Default to relative '/api' base path so that browser requests stay same-origin
// and are routed through Nginx (in production) or Vite proxy (in development).
const rawBase = import.meta.env.VITE_API_BASE_URL;
const API_BASE = (rawBase !== undefined && rawBase !== '')
  ? rawBase.replace(/\/+$/, '')
  : '/api';

/**
 * Constructs the final request URL while avoiding duplicate `/api` prefixes.
 *
 * Examples:
 *   API_BASE = '/api', path = '/api/v1/tasks'  --> '/api/v1/tasks'
 *   API_BASE = '/api', path = '/v1/tasks'      --> '/api/v1/tasks'
 *   API_BASE = 'http://localhost:8080', path = '/api/v1/tasks' --> 'http://localhost:8080/api/v1/tasks'
 */
function buildUrl(path) {
  const cleanPath = path.startsWith('/') ? path : `/${path}`;

  // If API_BASE already ends with '/api' and path begins with '/api/',
  // deduplicate the '/api' prefix (e.g. '/api' + '/api/v1/tasks' -> '/api/v1/tasks')
  if (API_BASE.endsWith('/api') && cleanPath.startsWith('/api/')) {
    return `${API_BASE.slice(0, -4)}${cleanPath}`;
  }

  return `${API_BASE}${cleanPath}`;
}

/**
 * Small wrapper around fetch that:
 *   1. Prepends the centralized API base URL (relative /api).
 *   2. Sets JSON headers.
 *   3. Throws on non-2xx responses with the server's error message.
 */
async function request(path, options = {}) {
  const url = buildUrl(path);

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

