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

## Project Structure

```
taskflow/
├── backend/
│   ├── cmd/api/main.go          # Application startup & graceful shutdown
│   ├── internal/
│   │   ├── config/              # Env loading & configuration
│   │   ├── database/            # pgxpool connection pool setup
│   │   ├── handlers/            # REST controllers & CORS middleware
│   │   └── tasks/               # Data model & parameterized SQL queries
│   ├── go.mod
│   └── .env.example
├── frontend/
│   ├── src/
│   │   ├── components/          # Dashboard, TaskForm, TaskCard, TaskList, StatusFilter
│   │   ├── services/api.js      # Fetch client wrapper
│   │   ├── App.jsx              # Root component & state management
│   │   └── main.jsx
│   ├── package.json
│   ├── vite.config.js
│   └── .env.example
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
