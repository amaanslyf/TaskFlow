# 📋 TaskFlow

A simple full-stack task management application built to **learn DevOps step by step**.

> **Stage 1 — The Application.**  This is the first milestone: a working full-stack app with a React frontend, a Go backend, and a PostgreSQL database. In later stages we will add Docker, Docker Compose, CI/CD, Kubernetes, Helm, Prometheus, Grafana, and more — one concept at a time.

---

## What does TaskFlow do?

TaskFlow lets you **create, view, edit, and delete tasks**. Each task has a title, an optional description, and a status (`todo`, `in_progress`, or `completed`). Think of it as a minimal Trello or Jira.

---

## Architecture overview

```
┌──────────────┐        HTTP / JSON       ┌──────────────┐        SQL        ┌──────────────┐
│              │  ────────────────────▸   │              │ ──────────────▸   │              │
│   React UI   │                          │   Go API     │                   │  PostgreSQL  │
│  (Vite dev   │  ◂────────────────────   │  (net/http)  │ ◂──────────────   │   Database   │
│   server)    │                          │              │                   │              │
│  port 5173   │                          │  port 8080   │                   │  port 5432   │
└──────────────┘                          └──────────────┘                   └──────────────┘
```

### What each piece does

| Layer | Technology | Role |
| :--- | :--- | :--- |
| **Frontend** | React + Vite | Shows the user interface in the browser. Sends HTTP requests to the backend whenever the user creates, edits, or deletes a task. |
| **Backend** | Go (`net/http`) | Receives HTTP requests, validates them, runs SQL queries against PostgreSQL, and sends JSON responses back to the frontend. |
| **Database** | PostgreSQL | Stores all the tasks permanently. Even if you restart the backend, data persists. |

---

## Environment Variables Explained

TaskFlow strictly adheres to the **Twelve-Factor App** configuration principles. All configuration is driven via environment variables:

| Variable | Component | Default / Example | Purpose |
| :--- | :--- | :--- | :--- |
| `DATABASE_URL` | Backend | `postgres://postgres:postgres@localhost:5432/taskflow?sslmode=disable` | Complete connection string for PostgreSQL including username, password, host, port, database name, and SSL mode. |
| `PORT` | Backend | `8080` | Port the Go HTTP server listens on. |
| `FRONTEND_ORIGIN` | Backend | `http://localhost:5173` | Allowed origin for Cross-Origin Resource Sharing (CORS). Tells the browser to permit requests originating from the React dev server. |
| `VITE_API_BASE_URL` | Frontend | `http://localhost:8080` | Base URL used by the frontend `api.js` client to make REST requests to the Go backend. |

---

## Prerequisites

Install these tools before beginning:

| Tool | Recommended Version | Download Link |
| :--- | :--- | :--- |
| **Node.js** | 18+ (Node 24 tested) | [nodejs.org](https://nodejs.org/) |
| **npm** | Bundled with Node | — |
| **Go** | 1.22+ (Go 1.27 tested) | [go.dev/dl](https://go.dev/dl/) |
| **PostgreSQL** | 14+ (PostgreSQL 17 tested) | [postgresql.org/download](https://www.postgresql.org/download/) |

Verify your local installation:

**Windows (PowerShell):**
```powershell
node --version
npm --version
go version
psql --version
```

**macOS / Linux:**
```bash
node -v
npm -v
go version
psql --version
```

---

## Step-by-Step Local Setup

### 1. Start PostgreSQL Service

Ensure your PostgreSQL service is running:

- **Windows (PowerShell as Admin):**
  ```powershell
  # Check service status
  Get-Service postgresql*

  # Start service if stopped (e.g., for PostgreSQL 17)
  Start-Service postgresql-x64-17
  ```
- **macOS (Homebrew):**
  ```bash
  brew services start postgresql@14
  ```
- **Linux (systemd):**
  ```bash
  sudo systemctl start postgresql
  ```

---

### 2. Create the Database & Apply Schema

Connect to PostgreSQL and create the `taskflow` database:

**Windows (PowerShell):**
```powershell
# Create database
psql -U postgres -c "CREATE DATABASE taskflow;"

# Apply schema file
psql -U postgres -d taskflow -f "database/schema.sql"

# Verify table creation
psql -U postgres -d taskflow -c "\d tasks"
```

**macOS / Linux (Bash):**
```bash
# Create database
createdb -U postgres taskflow || psql -U postgres -c "CREATE DATABASE taskflow;"

# Apply schema file
psql -U postgres -d taskflow -f database/schema.sql

# Verify table creation
psql -U postgres -d taskflow -c "\d tasks"
```

---

### 3. Configure the Backend

In the `backend/` directory, copy the example environment file and adjust your database password:

**Windows (PowerShell):**
```powershell
cd backend
Copy-Item .env.example .env
```

**macOS / Linux (Bash):**
```bash
cd backend
cp .env.example .env
```

Edit `backend/.env` with your PostgreSQL credentials:
```env
DATABASE_URL=postgres://postgres:YOUR_PASSWORD@localhost:5432/taskflow?sslmode=disable
PORT=8080
FRONTEND_ORIGIN=http://localhost:5173
```

---

### 4. Configure the Frontend

In the `frontend/` directory, create `.env`:

**Windows (PowerShell):**
```powershell
cd frontend
Copy-Item .env.example .env
```

**macOS / Linux (Bash):**
```bash
cd frontend
cp .env.example .env
```

Ensure `frontend/.env` points to your Go API:
```env
VITE_API_BASE_URL=http://localhost:8080
```

---

### 5. Install Dependencies & Build

**Backend (Go):**
```bash
cd backend
go mod tidy
go test -v ./...
go build ./cmd/api
```

**Frontend (React/Vite):**
```bash
cd frontend
npm install
npm run build
```

---

## Running the Application

You will need **two terminal sessions**:

### Terminal 1 — Start the Go Backend

```bash
cd backend
go run ./cmd/api
```

Expected startup output:
```
Starting TaskFlow API on port 8080
Connected to PostgreSQL
```

### Terminal 2 — Start the React Frontend

```bash
cd frontend
npm run dev
```

Vite will start the dev server and open your browser at **http://localhost:5173**.

---

## API Endpoints & Verification

All endpoints consume and return JSON.

### Health Probes
- `GET /healthz` — Liveness probe (verifies process is running).
- `GET /readyz` — Readiness probe (pings PostgreSQL).

### Tasks CRUD
- `GET /api/v1/tasks` — List all tasks.
- `GET /api/v1/tasks/{id}` — Get single task by UUID.
- `POST /api/v1/tasks` — Create a task (`title` required, optional `description`, `status`).
- `PUT /api/v1/tasks/{id}` — Update a task's `title`, `description`, or `status`.
- `DELETE /api/v1/tasks/{id}` — Delete a task.

---

### Verification Commands

#### Using Windows PowerShell:
```powershell
# 1. Healthz & Readyz
Invoke-RestMethod -Uri "http://localhost:8080/healthz"
Invoke-RestMethod -Uri "http://localhost:8080/readyz"

# 2. Create Task
$body = @{ title = "Learn Kubernetes"; description = "Deploy TaskFlow to k8s"; status = "todo" } | ConvertTo-Json
$task = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/tasks" -Method Post -Body $body -ContentType "application/json"
$taskId = $task.id
Write-Host "Created Task ID: $taskId"

# 3. List All Tasks
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/tasks"

# 4. Update Task Status
$update = @{ status = "in_progress" } | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/tasks/$taskId" -Method Put -Body $update -ContentType "application/json"

# 5. Delete Task
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/tasks/$taskId" -Method Delete
```

#### Using macOS / Linux (`curl`):
```bash
# 1. Healthz & Readyz
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz

# 2. Create Task
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Kubernetes", "description": "Deploy TaskFlow to k8s", "status": "todo"}'

# 3. List All Tasks
curl http://localhost:8080/api/v1/tasks

# 4. Update Task Status (replace <ID>)
curl -X PUT http://localhost:8080/api/v1/tasks/<ID> \
  -H "Content-Type: application/json" \
  -d '{"status": "in_progress"}'

# 5. Delete Task (replace <ID>)
curl -X DELETE http://localhost:8080/api/v1/tasks/<ID>
```

---

## Automated Testing

TaskFlow implements a multi-tier testing strategy following the **Testing Pyramid**:

```
            / \
           /   \       E2E Tests (Playwright)
          / E2E \      Browser interaction & user workflows
         /-------\
        /   API   \    API Integration Tests (pytest)
       /Integration\   HTTP contracts, status codes & DB persistence
      /-------------\
     /     Unit      \ Unit Tests (Go testing)
    /     Testing     \ Business logic, validation, config & middleware
   /-------------------\
```

### Understanding Test Types

| Test Level | Scope | Framework | Execution Speed | What it Catches |
| :--- | :--- | :--- | :--- | :--- |
| **Unit Testing** | Individual Go functions and handlers | Go `testing`, `httptest` | ~0.01s (Fastest) | Logic flaws, invalid status rejection, faulty defaults, parsing bugs |
| **API Integration** | Full HTTP request/response cycle against real backend & PostgreSQL | `pytest` + `requests` | ~0.5s (Fast) | SQL errors, DB constraints, HTTP status codes, schema validation |
| **Browser E2E** | Live UI rendered in headless browser | `Playwright` (Chromium) | ~30s (Comprehensive) | DOM regressions, broken click handlers, form flow, dialog handling |

---

### A. Go Backend Unit Testing

Unit tests run directly in Go without requiring an external database or network connection.

#### Execution Commands:

```bash
cd backend

# Run all unit tests
go test ./...

# Run with verbose test logs
go test -v ./...

# Run with statement coverage report
go test -cover ./...
```

**Windows PowerShell:**
```powershell
cd backend
go test -v -cover ./...
```

---

### B. Python API Integration Testing (pytest)

The API test suite validates all REST endpoints against the live Go backend and PostgreSQL database. It verifies JSON schemas, HTTP status codes, and input validation without polluting existing user data.

#### Prerequisites:
The Go backend and PostgreSQL **must be running** before executing these tests:
- PostgreSQL service running
- Go backend running: `go run ./cmd/api` (port 8080)

#### 1. Create and Activate Virtual Environment:

**Windows (PowerShell):**
```powershell
# From project root
python -m venv tests/.venv

# Activate virtual environment
.\tests\.venv\Scripts\Activate.ps1
```

*(If script execution is disabled in PowerShell, run `Set-ExecutionPolicy -Scope Process -ExecutionPolicy RemoteSigned`)*

**macOS / Linux (Bash):**
```bash
# From project root
python3 -m venv tests/.venv

# Activate virtual environment
source tests/.venv/bin/activate
```

#### 2. Install Dependencies:
```bash
pip install -r tests/api/requirements.txt
```

#### 3. Execute API Tests:
```bash
# Run entire API test suite
pytest -c tests/api/pytest.ini tests/api

# Run with extra verbose output
pytest -c tests/api/pytest.ini tests/api -v

# Run a specific test file
pytest -c tests/api/pytest.ini tests/api/test_health.py
pytest -c tests/api/pytest.ini tests/api/test_tasks.py
```

*Note: You can override the backend URL using the `API_BASE_URL` environment variable (default: `http://localhost:8080`).*

---

### C. Playwright Browser End-to-End Testing

Playwright drives a headless Chromium browser to execute user flows: viewing the dashboard, filling the task creation form, editing tasks, toggling status badges, filtering, and deleting records.

#### Prerequisites:
- PostgreSQL and the Go backend **must be running** (port 8080).
- Playwright's `webServer` configuration automatically starts the Vite frontend dev server at `http://localhost:5173`.

#### 1. Install Dependencies & Browser Binaries:

**Windows PowerShell & macOS/Linux:**
```bash
cd frontend

# Install dev dependencies (including @playwright/test)
npm install

# Download Chromium browser binary (one-time setup)
npx playwright install chromium
```

#### 2. Execute Browser Tests:
```bash
cd frontend

# Run all E2E tests in headless Chromium
npx playwright test

# Or using the npm script shortcut:
npm run test:e2e
```

#### 3. View Reports & Trace Artifacts:
```bash
cd frontend

# Open the HTML test report
npx playwright show-report

# Inspect traces and failure screenshots (if any test fails)
npx playwright show-trace test-results/<test-folder>/trace.zip
```

---

## Docker & Containerization Fundamentals

In this phase, we package the backend and frontend into self-contained, immutable **Docker images**.

```
┌────────────────────────────────────────────────────────────────────────┐
│                        Docker Image Architecture                       │
├──────────────────────────────────┬─────────────────────────────────────┤
│      taskflow-backend:dev        │        taskflow-frontend:dev        │
│          (~10.8 MB)              │             (~21.1 MB)              │
├──────────────────────────────────┼─────────────────────────────────────┤
│ ┌──────────────────────────────┐ │ ┌─────────────────────────────────┐ │
│ │ Stage 1: Build (Go compiler) │ │ │ Stage 1: Build (Node.js & Vite) │ │
│ │ Compiles static Linux binary │ │ │ Compiles JSX to static HTML/CSS │ │
│ └──────────────┬───────────────┘ │ └───────────────┬─────────────────┘ │
│                │ [binary only]   │                 │ [dist/ only]      │
│ ┌──────────────▼───────────────┐ │ ┌───────────────▼─────────────────┐ │
│ │ Stage 2: Runtime (Alpine)    │ │ │ Stage 2: Runtime (Nginx Alpine) │ │
│ │ Non-root appuser, port 8080  │ │ │ High-speed static server, p. 80 │ │
│ └──────────────────────────────┘ │ └─────────────────────────────────┘ │
└──────────────────────────────────┴─────────────────────────────────────┘
```

### Core Concepts Explained for Beginners

1. **What is Docker?**  
   Docker is a containerization platform that packages an application and all of its dependencies (system libraries, settings, runtime) together into a standard unit called an **image**. It guarantees that the application runs identically on your laptop, a teammate's computer, a CI/CD server, or a cloud Kubernetes cluster.
2. **Image vs. Container**:  
   - An **Image** is a read-only blueprint/template (similar to a class in OOP or an ISO installer).
   - A **Container** is a running instance of an image (similar to an object instance or a running process).
3. **What is a Dockerfile?**  
   A plain text file containing sequential instructions (`FROM`, `COPY`, `RUN`, `EXPOSE`, `CMD`) that Docker executes to assemble an image.
4. **What are Docker Layers?**  
   Each instruction in a Dockerfile creates a read-only filesystem layer. Docker caches layers: if a layer and its inputs haven't changed, Docker reuses the cached layer, dramatically speeding up subsequent builds.
5. **What is a Multi-Stage Build?**  
   A multi-stage build uses multiple `FROM` instructions in a single Dockerfile.
   - Heavy compilation tools (Go compiler ~800MB, Node.js ~1.1GB) are used in temporary **builder** stages.
   - The final **runtime** stage copies *only* the compiled artifacts (`api` executable or `dist/` HTML/JS files) into a featherweight base image (Alpine Linux ~7MB or Nginx Alpine ~20MB).
   - **Result**: Drastically smaller images (10-20MB instead of 1GB), faster deployments, and a substantially smaller attack surface (no shell scripts, package managers, or compilers in production).
6. **Why Nginx instead of Node.js for Production Frontend?**  
   Node.js is only needed at build-time to compile JSX and bundle assets. The output is purely static files. Nginx uses ~10MB RAM, handles thousands of concurrent requests, and optimizes gzip compression and client-side caching far better than a Node runtime.
7. **Why .dockerignore is Critical**:  
   Similar to `.gitignore`, `.dockerignore` prevents files like `.git`, `node_modules`, local test artifacts, and sensitive `.env` files from being sent to the Docker daemon. This keeps the build context small, accelerates builds, and prevents credentials from being baked into image layers.

---

### Building the Docker Images

#### 1. Build the Backend Image:

**Windows PowerShell:**
```powershell
docker build -t taskflow-backend:dev -f backend/Dockerfile backend
```

**macOS / Linux:**
```bash
docker build -t taskflow-backend:dev -f backend/Dockerfile backend
```

#### 2. Build the Frontend Image:

The frontend requires the API base URL to be injected at build time via a Docker build argument:

**Windows PowerShell:**
```powershell
docker build --build-arg VITE_API_BASE_URL=http://localhost:8080 -t taskflow-frontend:dev -f frontend/Dockerfile frontend
```

**macOS / Linux:**
```bash
docker build --build-arg VITE_API_BASE_URL=http://localhost:8080 -t taskflow-frontend:dev -f frontend/Dockerfile frontend
```

---

### Inspecting Docker Images

To view image sizes and layer history:

```bash
# List local images and verify size
docker images taskflow-backend:dev
docker images taskflow-frontend:dev

# Inspect the layer-by-layer build history of the backend
docker history taskflow-backend:dev

# Inspect the layer-by-layer build history of the frontend
docker history taskflow-frontend:dev

# Inspect the image metadata (entrypoint, healthcheck, exposed ports, non-root user)
docker image inspect taskflow-backend:dev
docker image inspect taskflow-frontend:dev
```

---

## Project Structure

```
taskflow/
├── backend/
│   ├── cmd/api/main.go          # Application startup & graceful shutdown
│   ├── internal/
│   │   ├── config/              # Env loading & config_test.go
│   │   ├── database/            # pgxpool connection pool setup
│   │   ├── handlers/            # REST controllers, middleware & handlers_test.go
│   │   └── tasks/               # Data model, parameterized SQL & tasks_test.go
│   ├── Dockerfile               # Multi-stage build: Go compiler -> Alpine runtime
│   ├── .dockerignore            # Excludes secrets, binaries, and test caches
│   ├── go.mod
│   ├── go.sum
│   └── .env.example
├── frontend/
│   ├── e2e/
│   │   └── taskflow.spec.js     # Playwright E2E browser test suite
│   ├── src/
│   │   ├── components/          # Dashboard, TaskForm, TaskCard, TaskList, StatusFilter
│   │   ├── services/api.js      # Fetch client wrapper
│   │   ├── App.jsx              # Root component & state management
│   │   └── main.jsx
│   ├── Dockerfile               # Multi-stage build: Node.js builder -> Nginx Alpine
│   ├── .dockerignore            # Excludes node_modules, dist, and test artifacts
│   ├── nginx.conf               # Nginx SPA fallback routing configuration
│   ├── package.json
│   ├── vite.config.js
│   ├── playwright.config.js     # Playwright test configuration
│   └── .env.example
├── tests/
│   ├── api/
│   │   ├── conftest.py          # pytest fixtures, APIClient & cleanup logic
│   │   ├── test_health.py       # Liveness / readiness probe tests
│   │   ├── test_tasks.py        # Complete task CRUD & validation test suite
│   │   ├── pytest.ini           # pytest settings & testpaths
│   │   └── requirements.txt     # Python test dependencies
│   └── .venv/                   # Python virtual environment (ignored in git)
├── database/
│   └── schema.sql               # PostgreSQL DDL
├── .gitignore
└── README.md
```

---

## Upcoming Stages

- 🐳 **Docker**: Containerize frontend and backend with multi-stage builds.
- 🐙 **Docker Compose**: Orchestrate React, Go, and PostgreSQL with a single command.
- 🔄 **GitHub Actions**: Automated linting, test suites, and container builds.
- ☸️ **Kubernetes & Helm**: Deploy declarative Pods, Services, Ingress, and ConfigMaps.
- 📊 **Prometheus & Grafana**: Export and visualize real-time application metrics.
