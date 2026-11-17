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

## 🐳 Running TaskFlow with Manual Docker Commands

In this milestone, we run the entire TaskFlow multi-container stack (**React frontend + Go backend + PostgreSQL database**) using individual `docker` commands. We deliberately avoid Docker Compose here to master the underlying primitives: networks, named volumes, port mappings, environment injection, container lifecycles, and internal DNS resolution.

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                                 HOST MACHINE (Windows / macOS / Linux)                  │
│                                                                                        │
│   Browser (Host)                API Client / Tests (Host)        psql / Admin Tool     │
│   http://localhost:3000         http://localhost:8080            localhost:5433        │
└─────────┬───────────────────────────────┬───────────────────────────────┬──────────────┘
          │                               │                               │
   Port   │ Host: 3000                    │ Host: 8080                    │ Host: 5433
  Mapping │ Cont: 80                      │ Cont: 8080                    │ Cont: 5432
          ▼                               ▼                               ▼
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                  Docker User-Defined Bridge Network: taskflow-net                      │
│                                                                                        │
│  ┌───────────────────────┐   HTTP REST   ┌───────────────────────┐   pgx SQL   ┌─────┴─────────────────┐
│  │   taskflow-frontend   │ ────────────▸ │   taskflow-backend    │ ──────────▸ │      taskflow-db      │
│  │ (taskflow-frontend:dev│ (via browser  │ (taskflow-backend:dev)│ taskflow-db:│  (postgres:17-alpine) │
│  │   Nginx on port 80)   │  host port)   │   Go API on port 8080 │    5432     │  Postgres on port 5432│
│  └───────────────────────┘               └───────────────────────┘             └──────────┬────────────┘
│                                                                                           │
└───────────────────────────────────────────────────────────────────────────────────────────┼────┘
                                                                                            │
                                                       Persistent Named Docker Volume:      │
                                                       taskflow-postgres-data ──────────────┘
                                                       (/var/lib/postgresql/data)
```

---

### Step-by-Step Manual Orchestration

#### Step 1: Create the User-Defined Bridge Network

Containers attached to the default bridge network cannot resolve each other by container name. Creating a **user-defined bridge network** activates Docker's internal DNS server, allowing `taskflow-backend` to reach PostgreSQL simply by referring to the hostname `taskflow-db`.

```bash
docker network create taskflow-net
```

*Verify network creation:*
```bash
docker network ls --filter name=taskflow-net
```

---

#### Step 2: Create the Persistent Named Volume

Containers are ephemeral: when a container is removed, its writable layer is destroyed. A **named volume** exists independently on the host storage and preserves PostgreSQL database files across container stops, restarts, and deletions.

```bash
docker volume create taskflow-postgres-data
```

*Verify volume creation:*
```bash
docker volume ls --filter name=taskflow-postgres-data
```

---

#### Step 3: Configure Environment Variables (`.env.docker`)

Create a root-level configuration file named `.env.docker` containing disposable local credentials:

```env
POSTGRES_DB=taskflow
POSTGRES_USER=postgres
POSTGRES_PASSWORD=taskflow_local_dev_only
DATABASE_URL=postgres://postgres:taskflow_local_dev_only@taskflow-db:5432/taskflow?sslmode=disable
PORT=8080
FRONTEND_ORIGIN=http://localhost:3000
```

> **Security Note:** `.env.docker` is excluded from version control via `.gitignore` and is excluded from Docker images via `.dockerignore`.

---

#### Step 4: Run the PostgreSQL Container

Run the official PostgreSQL 17 Alpine image in detached mode (`-d`), attached to `taskflow-net`, with the persistent volume mounted.

> **Port Conflict Note:** If your host machine already runs PostgreSQL locally on port `5432`, map host port `5433` to container port `5432` (`-p 5433:5432`). Inside `taskflow-net`, the container still listens on standard port `5432`.

**Windows PowerShell & macOS/Linux:**
```bash
docker run -d \
  --name taskflow-db \
  --network taskflow-net \
  -p 5433:5432 \
  -e POSTGRES_DB=taskflow \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=taskflow_local_dev_only \
  -v taskflow-postgres-data:/var/lib/postgresql/data \
  --restart unless-stopped \
  postgres:17-alpine
```

*Verify database readiness:*
```bash
docker exec taskflow-db pg_isready -U postgres
```

---

#### Step 5: Apply the Initial Database Schema

Apply `database/schema.sql` into the newly created database container:

**Windows PowerShell:**
```powershell
Get-Content database/schema.sql | docker exec -i taskflow-db psql -U postgres -d taskflow
```

**macOS / Linux (Bash):**
```bash
docker exec -i taskflow-db psql -U postgres -d taskflow < database/schema.sql
```

*Verify the `tasks` table and indexes exist:*
```bash
docker exec taskflow-db psql -U postgres -d taskflow -c "\d tasks"
```

---

#### Step 6: Run the Go Backend Container

Start the Go API container using our multi-stage image `taskflow-backend:dev`. Inject runtime configuration via `--env-file .env.docker`.

**Windows PowerShell & macOS/Linux:**
```bash
docker run -d \
  --name taskflow-backend \
  --network taskflow-net \
  -p 8080:8080 \
  --env-file .env.docker \
  --restart unless-stopped \
  taskflow-backend:dev
```

*Verify health and readiness probes:*
```bash
# Check Docker container status and built-in HEALTHCHECK
docker ps --filter name=taskflow-backend

# Test HTTP probes from host
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

---

#### Step 7: Run the React Frontend Container

Start the Nginx static frontend container using `taskflow-frontend:dev`. Map host port `3000` to Nginx container port `80`.

**Windows PowerShell & macOS/Linux:**
```bash
docker run -d \
  --name taskflow-frontend \
  --network taskflow-net \
  -p 3000:80 \
  --restart unless-stopped \
  taskflow-frontend:dev
```

Open **http://localhost:3000** in your browser to interact with the full-stack application!

---

### Inspecting and Monitoring Containers

#### Check Container Status
```bash
# Show running containers with port mappings
docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}"

# Show all containers including stopped ones
docker ps -a
```

#### Inspect Live Logs
```bash
# Follow Go backend logs
docker logs -f taskflow-backend

# Follow PostgreSQL logs
docker logs -f taskflow-db

# Follow Nginx frontend access logs
docker logs -f taskflow-frontend

# Inspect the last 50 lines of logs
docker logs --tail 50 taskflow-backend
```

#### Monitor Resource Utilization
```bash
# Snapshot of CPU, memory, and network I/O per container
docker stats --no-stream
```

#### Inspect Network and DNS Members
```bash
# View all containers attached to taskflow-net and their internal IP addresses
docker network inspect taskflow-net
```

#### Inspect Persistent Volume
```bash
# View volume mountpoint and host storage path
docker volume inspect taskflow-postgres-data
```

---

### Verifying End-to-End Application Health

#### 1. REST API Integration Test (from Host)

**Windows PowerShell:**
```powershell
# Create a task
$task = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/tasks" -Method Post -Body '{"title":"Dockerized Task","description":"Running inside Docker","status":"todo"}' -ContentType "application/json"
$taskId = $task.id

# Verify task exists in PostgreSQL container via SQL
docker exec taskflow-db psql -U postgres -d taskflow -c "SELECT id, title, status FROM tasks WHERE id = '$taskId';"

# Update task status to in_progress
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/tasks/$taskId" -Method Put -Body '{"status":"in_progress"}' -ContentType "application/json"

# Delete task
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/tasks/$taskId" -Method Delete
```

**macOS / Linux (`curl`):**
```bash
# Create a task
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Dockerized Task","description":"Running inside Docker","status":"todo"}'

# Verify list of tasks
curl http://localhost:8080/api/v1/tasks
```

#### 2. Run Automated Pytest Suite Against Containers
```bash
pytest -c tests/api/pytest.ini tests/api -v
```

#### 3. Run Automated Playwright E2E Tests Against Containers
```bash
cd frontend
npx playwright test
```

---

### Verifying Data Persistence

To prove that database records persist when containers are stopped or removed:

```bash
# 1. Insert a test record via psql into taskflow-db
docker exec taskflow-db psql -U postgres -d taskflow -c "INSERT INTO tasks (title, status) VALUES ('Persistence Test Task', 'in_progress');"

# 2. Stop and restart the PostgreSQL container
docker stop taskflow-db
docker start taskflow-db

# 3. Query the database to verify the record survived
docker exec taskflow-db psql -U postgres -d taskflow -c "SELECT title, status FROM tasks WHERE title = 'Persistence Test Task';"
```

---

### Container Lifecycle & Teardown Commands

#### Stop Running Containers
```bash
docker stop taskflow-frontend taskflow-backend taskflow-db
```

#### Restart Stopped Containers
```bash
docker start taskflow-db taskflow-backend taskflow-frontend
```

#### Remove Containers (Preserving Data Volume)
```bash
# Stop and remove containers
docker stop taskflow-frontend taskflow-backend taskflow-db
docker rm taskflow-frontend taskflow-backend taskflow-db
```
*Note: Your tasks remain completely safe inside `taskflow-postgres-data`.*

#### Full Cleanup (Destructive: Removes Containers, Network, and Volume)
```bash
# Stop and remove containers
docker rm -f taskflow-frontend taskflow-backend taskflow-db

# Remove user-defined network
docker network rm taskflow-net

# Remove persistent volume (WARNING: All database records will be permanently deleted)
docker volume rm taskflow-postgres-data
```

---

### Deep-Dive: Core Docker Concepts for Beginners

| Concept | Explanation | TaskFlow Example |
| :--- | :--- | :--- |
| **Image vs. Container** | An **Image** is a read-only, immutable template containing code, libraries, and runtime files. A **Container** is a running, isolated OS process executing that image with a thin writable layer. | `taskflow-backend:dev` is the image (~10.8 MB); `taskflow-backend` is the running container process (~5.8 MB RAM). |
| **Container ID vs. Name** | A **Container ID** is a unique 64-character SHA256 hash (e.g. `030c97a1025f`). A **Container Name** is a human-readable alias assigned via `--name` (e.g. `taskflow-backend`). Both reference the container in CLI commands. | You can run `docker logs 030c97a1025f` or `docker logs taskflow-backend` interchangeably. |
| **Docker Daemon vs. CLI** | The **Docker CLI** (`docker`) is a client terminal tool that sends REST API commands to the background **Docker Daemon** (`dockerd` / Docker Desktop Engine), which actually manages containers, storage, and virtual networks. | Typing `docker run ...` in PowerShell communicates via a Windows named pipe (`\\.\pipe\docker_engine`) to the Docker daemon. |
| **Port Mapping (`-p host:cont`)** | Containers run inside isolated virtual network namespaces. `-p <host_port>:<container_port>` instructs Docker to configure Linux `iptables` / Windows NAT forwarding so incoming traffic on the host port is routed into the container's listening port. | `-p 5433:5432` forwards host requests on port `5433` into PostgreSQL listening on port `5432`. |
| **User-Defined Bridge Network** | A private virtual software switch created on the host. Containers on the same bridge network can communicate over any port without publishing those ports to the host machine. | `taskflow-net` connects the backend and database securely. |
| **Docker Internal DNS** | User-defined networks feature an embedded DNS server at `127.0.0.11`. Docker automatically registers container names as hostnames. Containers resolve peer IP addresses dynamically without hardcoded IPs. | Go backend connects to `taskflow-db:5432`. Docker DNS resolves `taskflow-db` to `172.18.0.2`. |
| **`localhost` Inside a Container** | Inside a container, `localhost` (127.0.0.1) refers strictly to **that container's own network namespace**, NOT the host computer and NOT other containers. | If `taskflow-backend` tries to connect to `localhost:5432`, it fails because PostgreSQL is running in a different container. It must connect to `taskflow-db:5432`. |
| **Named Volumes vs. Writable Layer** | A container's writable layer is deleted when the container is removed. A **Named Volume** is a Docker-managed directory on the host filesystem that bypasses the container storage driver and persists indefinitely until explicitly deleted via `docker volume rm`. | `taskflow-postgres-data` stores `/var/lib/postgresql/data`. Even if `taskflow-db` is deleted and recreated, database tables and records remain intact. |
| **Runtime vs. Build-Time Variables** | Build-time variables (`--build-arg` in `docker build`) bake values directly into static frontend assets (e.g., Vite compiling `VITE_API_BASE_URL` into JavaScript bundles). Runtime variables (`-e` or `--env-file` in `docker run`) inject dynamic values into backend processes at container startup. | `VITE_API_BASE_URL` is baked into frontend HTML/JS during `docker build`. `DATABASE_URL` is injected into the Go backend process at container startup via `--env-file .env.docker`. |
| **Docker Container Health Status** | Docker checks container health via `HEALTHCHECK`. Status progresses from `starting` to `healthy` (probe exits 0) or `unhealthy` (probe fails consecutive retries). | `taskflow-backend` runs `wget -q -O - http://localhost:8080/healthz` every 15s. `docker ps` displays `Up (healthy)`. |
| **Docker HEALTHCHECK vs. Kubernetes Probes** | Docker's `HEALTHCHECK` informs local container engines whether a single container is healthy. Kubernetes ignores Dockerfile healthchecks and instead uses declarative **Liveness** (restarts dead pods) and **Readiness** (routes ingress traffic only when ready) probes configured in Pod YAML manifests. | Our `/healthz` and `/readyz` endpoints serve both Docker's local `HEALTHCHECK` and future Kubernetes probes. |

---

---

## 🐙 Multi-Container Orchestration with Docker Compose

In this milestone, we graduate from manual `docker run` commands to **declarative infrastructure-as-code** using modern **Docker Compose** (`compose.yaml`).

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                              DOCKER COMPOSE TOPOLOGY                                    │
│                                                                                        │
│   Browser (Host)                API Client / Tests (Host)        psql / Admin Tool     │
│   http://localhost:3000         http://localhost:8080            localhost:5433        │
└─────────┬───────────────────────────────┬───────────────────────────────┬──────────────┘
          │                               │                               │
   Port   │ Host: 3000                    │ Host: 8080                    │ Host: 5433
  Mapping │ Cont: 80                      │ Cont: 8080                    │ Cont: 5432
          ▼                               ▼                               ▼
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                 Compose Default Bridge Network: taskflow_default                       │
│                                                                                        │
│  ┌───────────────────────┐   HTTP REST   ┌───────────────────────┐   pgx SQL   ┌─────┴─────────────────┐
│  │       frontend        │ ────────────▸ │        backend        │ ──────────▸ │       postgres        │
│  │ (taskflow-frontend:dev│ (via browser  │ (taskflow-backend:dev)│  postgres:  │  (postgres:17-alpine) │
│  │   Nginx on port 80)   │  host port)   │   Go API on port 8080 │    5432     │  Postgres on port 5432│
│  └──────────┬────────────┘               └──────────┬────────────┘             └──────────┬────────────┘
│             │                                       │                                     │
│             │ depends_on (service_healthy)          │ depends_on (service_healthy)        │
│             └───────────────────────────────────────┴─────────────────────────────────────┘
└───────────────────────────────────────────────────────────────────────────────────────────┼────┘
                                                                                            │
                                                       Persistent External Named Volume:    │
                                                       taskflow-postgres-data ──────────────┘
                                                       (/var/lib/postgresql/data)
```

---

### Core Concepts Explained for Beginners

1. **What is Docker Compose?**  
   Docker Compose is a declarative tool for defining and running multi-container Docker applications. Instead of running 10 different shell commands with complex arguments, you describe your entire system in a single YAML file (`compose.yaml`) and start everything with `docker compose up`.
2. **What is a Compose Service?**  
   A service is a computing component of your application (e.g. `postgres`, `backend`, `frontend`). Compose builds or pulls the image, creates an isolated container for each service, and automatically wires up networking and volumes.
3. **Automatic Network & Service Discovery:**  
   Compose automatically creates an isolated bridge network for the project (`taskflow_default`). Each service is automatically reachable by its **service name** (e.g., `postgres` is resolvable as a hostname at port `5432`). You never need to hardcode IP addresses.
4. **Orderly Startup with `depends_on` and `condition: service_healthy`:**  
   Applications have startup dependencies. The backend cannot work without PostgreSQL, and the frontend needs the backend. With `condition: service_healthy`, Compose does not merely wait for a container process to spawn—it waits for PostgreSQL's health check (`pg_isready`) to succeed before launching the backend, and waits for the backend's `/healthz` probe before launching the frontend.
5. **Preserving External Named Volumes:**  
   By marking `postgres_data` with `external: true` and pointing it to `taskflow-postgres-data`, Compose mounts our existing volume without creating a new prefixed volume or resetting the database. `docker compose down` will destroy the containers and network, but will **never delete the external volume**.
6. **PostgreSQL Initialization Semantics:**  
   Mounting `./database/schema.sql` into `/docker-entrypoint-initdb.d/001-schema.sql` ensures that if the database volume is ever completely fresh, PostgreSQL automatically initializes the tables. If an existing database already exists in the volume, PostgreSQL logs `PostgreSQL Database directory appears to contain a database; Skipping initialization` and leaves your existing data completely untouched!

---

### Compose Configuration File (`compose.yaml`)

```yaml
services:
  postgres:
    image: postgres:17-alpine
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    ports:
      - "5433:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./database/schema.sql:/docker-entrypoint-initdb.d/001-schema.sql:ro
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U \"$${POSTGRES_USER}\" -d \"$${POSTGRES_DB}\""]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    image: taskflow-backend:dev
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: ${DATABASE_URL}
      PORT: ${PORT}
      FRONTEND_ORIGIN: ${FRONTEND_ORIGIN}
    depends_on:
      postgres:
        condition: service_healthy

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
      args:
        VITE_API_BASE_URL: ${VITE_API_BASE_URL:-http://localhost:8080}
    image: taskflow-frontend:dev
    restart: unless-stopped
    ports:
      - "3000:80"
    depends_on:
      backend:
        condition: service_healthy

volumes:
  postgres_data:
    external: true
    name: taskflow-postgres-data
```

---

### Step-by-Step Operations

### Docker Compose Health Checks and Lifecycle Management

#### 1. Service Health Monitoring Architecture

All three services implement native Docker health checks to ensure reliability and automated startup coordination:

| Service | Probe Mechanism | Command | Key Parameters |
| :--- | :--- | :--- | :--- |
| **`postgres`** | Database Ping | `pg_isready -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"` | `interval: 10s`, `timeout: 5s`, `retries: 5`, `start-period: 10s` |
| **`backend`** | HTTP Probe | `wget -q -O - http://localhost:8080/healthz \|\| exit 1` | `interval: 30s`, `timeout: 5s`, `retries: 3`, `start-period: 5s` |
| **`frontend`** | Nginx Web Ping | `wget -q -O /dev/null http://127.0.0.1/ \|\| exit 1` | `interval: 30s`, `timeout: 5s`, `retries: 3`, `start-period: 5s` |

##### Understanding Health Check Parameters:
- **`interval`**: How often Docker runs the check (e.g., `30s` = every 30 seconds).
- **`timeout`**: Maximum allowed execution time before the probe is marked failed (e.g., `5s` = fails if probe hangs for >5 seconds).
- **`start-period`**: Initialization grace period. Probes that fail during this initial window do NOT count towards the retry limit, giving slow services time to boot.
- **`retries`**: Number of consecutive failures needed before the container transitions to `unhealthy`.

##### Expected `docker compose ps` Output:
```powershell
docker compose --env-file .env.docker ps
```
```
NAME                  IMAGE                   COMMAND                  SERVICE    STATUS                    PORTS
taskflow-backend-1    taskflow-backend:dev    "/app/api"               backend    Up (healthy)              0.0.0.0:8080->8080/tcp
taskflow-frontend-1   taskflow-frontend:dev   "/docker-entrypoint.…"   frontend   Up (healthy)              0.0.0.0:3000->80/tcp
taskflow-postgres-1   postgres:17-alpine      "docker-entrypoint.s…"   postgres   Up (healthy)              0.0.0.0:5433->5432/tcp
```

---

#### 2. Health-Based Dependency Ordering (`depends_on: condition: service_healthy`)

```
┌─────────────────────────┐
│   postgres container    │
│  (runs pg_isready loop) │
└───────────┬─────────────┘
            │ becomes healthy
            ▼
┌─────────────────────────┐
│    backend container    │
│ (connects & starts API) │
│   (evaluates /healthz)  │
└───────────┬─────────────┘
            │ becomes healthy
            ▼
┌─────────────────────────┐
│   frontend container    │
│ (starts serving Nginx)  │
└─────────────────────────┘
```

- **Why a running PostgreSQL container is not necessarily ready**: A container transition from "created" to "running" only means the Linux kernel started PID 1. PostgreSQL might still be replaying write-ahead logs (WAL), allocating buffer pools, or binding sockets. If the Go backend tries to connect immediately, connection refused errors crash the backend.
- **Why `condition: service_healthy` is essential**: Traditional short-form `depends_on: [postgres]` only waits for the postgres container process to start. `condition: service_healthy` halts dependent container launches until the upstream health probe actually returns exit code `0`.
- **What happens on initial failure**: If PostgreSQL fails its health checks during startup, Compose cancels the startup pipeline and does NOT launch `backend` or `frontend`.
- **Critical Limitation — Startup Only, Not Continuous Healing**:  
  `depends_on` is strictly a **startup-order coordinator**. If PostgreSQL crashes or becomes `unhealthy` 2 hours later, Docker Compose will **NOT** automatically restart or stop the backend. Continuous self-healing, automatic restarts, and traffic shedding are responsibilities of production orchestrators like Kubernetes.

---

#### 3. Deep-Dive: Lifecycle Commands Compared

| Command | Action on Containers | Action on Images | Action on Network/Volumes | When to Use in TaskFlow |
| :--- | :--- | :--- | :--- | :--- |
| **`docker compose start`** | Starts *existing, stopped* containers. | Does NOT build or pull. | Uses existing network and volumes. | Fast wake-up when containers were paused via `docker compose stop` and no configuration or code has changed. |
| **`docker compose up -d`** | Creates missing containers; recreates containers whose configuration or image changed; reuses unchanged containers. | Pulls or builds if missing. | Creates project network if missing; attaches volumes. | **Primary daily command**: Declaratively ensures desired state matches `compose.yaml`. |
| **`docker compose up --build -d`** | Forces recompilation and rebuild of images before starting or recreating containers. | Rebuilds `taskflow-backend:dev` and `taskflow-frontend:dev`. | Reuses existing volumes and networks. | Use whenever you modify frontend React code or Go backend source code. |
| **`docker compose stop`** | Sends `SIGTERM` (then `SIGKILL`) to halt processes. Containers enter `Exited (0)`. | None. | Containers, networks, and volumes remain completely preserved. | Use when pausing work to save laptop CPU/RAM without deleting container state. |
| **`docker compose down`** | Stops and deletes service containers and removes project network `taskflow_default`. | None. | **Preserves external named volumes** (`taskflow-postgres-data`). | Clean teardown at the end of the day. |

> ⚠️ **DANGER — Never run `docker compose down -v`:**  
> The `-v` flag instructs Compose to delete all attached volumes. While `external: true` protects `taskflow-postgres-data` from default teardown, avoiding `-v` is an essential DevOps best practice to prevent accidental production data loss.

---

#### 4. Practical Workflow Commands (Windows PowerShell)

```powershell
# 1. Validate configuration without starting
docker compose --env-file .env.docker config --quiet

# 2. Build and launch all services in background
docker compose --env-file .env.docker up --build -d

# 3. Check health and status of all services
docker compose --env-file .env.docker ps

# 4. Rebuild and update ONLY the frontend service
docker compose --env-file .env.docker up --build -d --no-deps frontend

# 5. Follow live logs of a specific service
docker compose --env-file .env.docker logs -f backend
docker compose --env-file .env.docker logs -f postgres

# 6. Execute SQL queries directly in the database container
docker compose --env-file .env.docker exec postgres psql -U postgres -d taskflow -c "\d tasks"

# 7. Pause the stack (saves battery/RAM without deleting containers)
docker compose --env-file .env.docker stop

# 8. Resume paused stack
docker compose --env-file .env.docker start

# 9. Clean teardown (preserves database volume)
docker compose --env-file .env.docker down
```

---

## 🔄 Continuous Integration with GitHub Actions

In this milestone, we introduce **Continuous Integration (CI)** using **GitHub Actions**. Every commit pushed to `main` and every Pull Request targeting `main` automatically triggers an end-to-end automated validation pipeline before allowing production Docker images to be built.

```
                  GitHub Trigger: Push / Pull Request on main
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                           PARALLEL STAGE: AUTOMATED QUALITY GATES                       │
├─────────────────────┬─────────────────────┬──────────────────────┬──────────────────────┤
│       go-unit       │   frontend-build    │   api-integration    │     browser-e2e      │
│  (Go 1.27.4 Runner) │  (Node.js 22 Runner)│  (pytest + PG 17)    │ (Playwright + PG 17) │
│                     │                     │                      │                      │
│ • go test -race     │ • npm ci            │ • Service Container  │ • Service Container  │
│ • go vet            │ • npm run build     │ • Apply schema.sql   │ • Apply schema.sql   │
│ • Memory safety     │ • Static asset lint │ • Go API in background│ • Go API + Vite dev  │
│ • Unit coverage     │ • Tree-shaking test │ • 19 pytest cases    │ • 6 Chromium flows   │
└──────────┬──────────┴──────────┬──────────┴──────────┬───────────┴──────────┬───────────┘
           │                     │                     │                      │
           └─────────────────────┴──────────┬──────────┴──────────────────────┘
                                            │
                                            ▼ (all 4 jobs must succeed)
                               ┌─────────────────────────┐
                               │      docker-build       │
                               │  (Production Artifacts) │
                               │                         │
                               │ • taskflow-backend:ci   │
                               │ • taskflow-frontend:ci  │
                               │ • Multi-stage verify    │
                               └─────────────────────────┘
```

---

### Core CI Concepts for Beginners

1. **What is Continuous Integration (CI)?**  
   CI is the automated practice of testing and building application code on every single commit. Instead of discovering bugs during release week or on production servers, defects are caught within 2 minutes of pushing code.
2. **What is GitHub Actions?**  
   An event-driven automation platform built directly into GitHub. It automatically triggers workflows in response to repository events like git pushes, pull requests, releases, or cron schedules.
3. **Workflows, Jobs, and Steps:**  
   - **Workflow**: The automated process defined in a YAML file (`.github/workflows/ci.yml`).
   - **Job**: A set of steps that execute on an independent virtual machine (**Runner**). Jobs execute **in parallel** by default unless configured with dependencies.
   - **Step**: An individual task inside a job (e.g., checking out code, setting up a language, executing a test command).
4. **What are GitHub-Hosted Runners?**  
   Clean, ephemeral virtual machines provided by GitHub (e.g., `ubuntu-latest`). Each job gets its own isolated VM that is automatically destroyed when the job completes.
5. **What are Service Containers?**  
   Docker containers configured in GitHub Actions to provide supporting services (like PostgreSQL 17) alongside the runner VM. The container runs on a Docker bridge network with ports published to `localhost` on the Ubuntu runner.
6. **Why Run Jobs in Parallel?**  
   Running `go-unit`, `frontend-build`, `api-integration`, and `browser-e2e` concurrently slashes pipeline runtime from ~15 minutes down to ~2–3 minutes.
7. **Quality Gates with `needs`:**  
   The `docker-build` job uses `needs: [go-unit, frontend-build, api-integration, browser-e2e]`. If even one test fails, GitHub Actions halts the pipeline immediately and will **not** build Docker images from flawed code.

---

### The TaskFlow CI Pipeline Jobs Explained

#### Job 1: `go-unit` (Go Unit Testing & Code Quality)
- **Role**: Validates Go business logic, validation rules, CORS headers, and configuration loaders in total isolation without requiring PostgreSQL.
- **Commands**:
  ```bash
  go test -v -race -cover ./...
  go vet ./...
  ```
- **Why `-race` Matters**: The Go race detector instruments memory accesses to catch multi-threaded data races that could cause memory corruption in concurrent web servers.
- **Why `go vet` Matters**: Analyzes Go source code for suspicious constructs, dead code, printf formatting mismatches, and unreachable statements.

#### Job 2: `frontend-build` (Vite SPA Compilation)
- **Role**: Validates that all JSX components, CSS imports, and dependencies compile into a valid production bundle.
- **Why `npm ci` instead of `npm install`?**  
  `npm ci` strictly adheres to `package-lock.json`, bypasses user prompts, deletes existing `node_modules`, and guarantees deterministic dependency trees across every CI build.
- **Commands**:
  ```bash
  npm ci
  npm run build
  ```

#### Job 3: `api-integration` (REST API & PostgreSQL Testing)
- **Role**: Verifies REST endpoints against a live PostgreSQL 17 database using `pytest`.
- **How It Works**:
  1. GitHub Actions launches a `postgres:17-alpine` service container on `localhost:5432` with health checks.
  2. The runner applies `database/schema.sql`.
  3. The Go backend launches in the background using CI-disposable credentials (`ci_only_password`).
  4. The job polls `GET http://localhost:8080/healthz` until healthy.
  5. Pytest executes 19 integration tests validating CRUD operations, UUID parsing, 404 responses, and status constraints.
  6. Backend logs are automatically printed if tests fail.

#### Job 4: `browser-e2e` (Playwright Chromium End-to-End Testing)
- **Role**: Simulates actual user interactions in a real headless browser.
- **How It Works**:
  1. Boots an independent PostgreSQL 17 service container and applies the schema.
  2. Starts the Go backend on `localhost:8080`.
  3. Installs Playwright Chromium with system dependencies: `npx playwright install --with-deps chromium`.
  4. Playwright's `webServer` launches the Vite development server on `localhost:5173`.
  5. Chromium tests user flows: creating tasks, editing descriptions, clicking status badges, filtering, and deleting records.
  6. On failure, Playwright HTML trace reports are automatically uploaded as GitHub Actions artifacts for debugging.

#### Jobs 5A & 5B: Mutually Exclusive Docker Build & Publish Quality Gates
To maintain strict **least privilege**, the pipeline separates image building into two mutually exclusive jobs:

- **`docker-build-pr` (Pull Request Verification)**:
  - **Trigger**: `if: github.event_name == 'pull_request'`
  - **Permissions**: `contents: read` (strictly read-only; no `packages: write`).
  - **Action**: Builds both images via Docker Buildx and loads them locally (`push: false`, `load: true`) to verify compilation and Dockerfile validity.
- **`docker-publish-main` (Production GHCR Publishing)**:
  - **Trigger**: `if: github.event_name == 'push' && github.ref == 'refs/heads/main'`
  - **Permissions**: `contents: read` and `packages: write`.
  - **Action**: Authenticates via `GITHUB_TOKEN`, builds both images with BuildKit layer caching, and publishes them with dual tags (`main` and `sha-<commit>`) to `ghcr.io`.
- **Quality Gate Dependency**: Both jobs declare `needs: [go-unit, frontend-build, api-integration, browser-e2e]`. If any test fails, Docker builds never execute.

---

### How to Inspect and Troubleshoot CI Workflow Runs

1. Navigate to the **Actions** tab in your GitHub repository (`https://github.com/<owner>/TaskFlow/actions`).
2. Click on the latest commit workflow run to view the interactive DAG pipeline graph.
3. Click any individual job (e.g. `api-integration` or `browser-e2e`) to inspect live step logs.
4. **Workflow Run Summary**: After a successful `main` branch push, inspect the **Job Summary** at the bottom of the `docker-publish-main` page for the published image repositories, tags, and SHA-256 digests.
5. **Troubleshooting Tips**:
   - **Backend Failure**: If the backend fails to connect to PostgreSQL or exits unexpectedly, the workflow runs `cat backend.log` to print Go panic traces or database connection errors.
   - **Playwright Test Failure**: Download the `playwright-artifacts` ZIP archive at the bottom of the workflow page to view DOM screenshots, videos, and Playwright HTML traces.

---

## 📦 GitHub Container Registry (GHCR)

TaskFlow automatically builds, version-tags, and publishes production-ready OCI container images to **GitHub Container Registry (GHCR)** upon every push to `main`.

### 1. What is a Container Registry?
A container registry is a centralized, secure service for storing, versioning, and distributing Docker and OCI-compliant container images. Just as GitHub stores and versions Git source code commits, a container registry stores binary image layers and cryptographic image manifests.

**GitHub Container Registry (`ghcr.io`)** is GitHub's native registry built into GitHub Packages. Unlike generic third-party registries, GHCR provides:
- Seamless integration with GitHub repositories and GitHub Actions.
- Native authentication using ephemeral, scoped repository tokens (`GITHUB_TOKEN`).
- Granular package permissions tied directly to repository collaborators and organizations.

### 2. Published TaskFlow Image Repositories

| Component | Target Image Repository | Published Tags |
| :--- | :--- | :--- |
| **Backend API** | `ghcr.io/amaanslyf/taskflow-backend` | `:main`, `:sha-<full-commit-sha>` |
| **Frontend SPA** | `ghcr.io/amaanslyf/taskflow-frontend` | `:main`, `:sha-<full-commit-sha>` |

> **Naming Rule**: All container repository paths on GHCR must be strictly lowercase (`ghcr.io/<owner>/<repo>`). The workflow normalizes the repository owner dynamically using explicit step outputs (`steps.image-names.outputs.backend` and `steps.image-names.outputs.frontend`).

### 3. Image Tagging Strategy: Mutable Tags vs Immutable Digests

TaskFlow implements an automated dual-tagging strategy:

```
                                  ┌───► :main (Mutable branch pointer)
ghcr.io/amaanslyf/taskflow-backend ───┼───► :sha-e9c5747... (Unique per commit by convention)
                                  └───► @sha256:7f3b... (Cryptographically immutable digest)
```

1. **Mutable Branch Tag (`:main`)**:
   - Updated on every push to `main`.
   - Useful for quick developer local testing to pull the latest build without checking commit hashes.
   - **Limitation**: Not repeatable over time—`:main` points to different code today than tomorrow.
2. **Commit SHA Tag (`:sha-<full-commit-sha>`)**:
   - Pinned to the exact 40-character Git commit SHA that produced the build.
   - **Important Security Distinction**: A Git SHA tag is unique per commit and immutable *by convention*, but standard container registries do not automatically prevent a publisher with `packages: write` from overwriting a tag.
3. **Cryptographic Image Digest (`@sha256:<digest>`)**:
   - The SHA-256 hash of the image manifest itself.
   - **True Immutability**: Any alteration to a single byte in any layer produces an entirely different digest.
   - **Kubernetes Production Best Practice**: Production Kubernetes manifests can reference images directly by digest:
     ```yaml
     image: ghcr.io/amaanslyf/taskflow-backend@sha256:7f3b25...
     ```
     This guarantees tamper-proof, deterministic container execution regardless of any registry tag changes.

### 4. Docker Buildx and BuildKit Acceleration

The workflow introduces `docker/setup-buildx-action@v4` and `docker/build-push-action@v7`, replacing legacy `docker build` with Moby **BuildKit**:
- **Parallel Stage Execution**: Non-dependent multi-stage build layers are compiled simultaneously.
- **GitHub Actions Layer Caching (`type=gha`)**: Unchanged dependency layers (such as Go modules in `go.sum` and NPM packages in `package-lock.json`) are cached in GitHub Actions cache storage.
- **Isolated Cache Scopes**:
  - `taskflow-backend`: Isolated Go build cache.
  - `taskflow-frontend`: Isolated Node.js and Vite build cache.
  - Separating scopes prevents cache eviction caused by one service overwriting another's cache entries.

### 5. Authentication & Strict Least-Privilege Security Model

To prevent privilege escalation, the workflow splits container operations into **two mutually exclusive jobs**:

1. **`docker-build-pr` (Pull Requests)**:
   - Granted **strictly `contents: read`** permissions.
   - Pull requests—including same-repository PRs—never receive `packages: write` tokens.
   - Compiles and validates both Docker images in BuildKit (`load: true`, `push: false`) without authenticating to GHCR.
2. **`docker-publish-main` (Pushes to `main`)**:
   - Granted `contents: read` and `packages: write`.
   - Authenticates using the ephemeral, scoped `GITHUB_TOKEN` (`secrets.GITHUB_TOKEN`).
   - Pushes both images to GHCR with dual tags (`main` and `sha-<commit>`).

### 6. Package Visibility and Permissions

It is a common misconception that packages automatically share the exact visibility of their source repository. In GHCR:
- **Permission Inheritance**: Packages inherit access permissions (e.g. read/write rights of repository collaborators and teams).
- **Visibility Separation**: Package visibility is managed independently from repository visibility.
- **Default Visibility**: Newly published packages on GHCR **default to private**, even if the source repository is public.
- **Making a Package Public**: Requires an explicit, intentional manual setting change by a repository admin:
  1. Navigate to **GitHub Profile / Organization** -> **Packages** (or repository sidebar -> **Packages**).
  2. Click **taskflow-backend** or **taskflow-frontend**.
  3. Go to **Package settings** -> **Danger Zone** -> **Change package visibility** -> **Public**.
- **Public vs. Private Access**:
  - **Public Packages**: Anyone can run `docker pull ghcr.io/...` anonymously without credentials.
  - **Private Packages**: Pulling requires an authenticated Docker client with a Personal Access Token (`read:packages` scope) or a Kubernetes `imagePullSecret`.

### 7. How to Pull and Run Published Images Locally

#### Authenticate with GHCR (Required for Private Packages)
Generate a GitHub Personal Access Token (classic) with the `read:packages` scope, then authenticate via PowerShell:

```powershell
# Set your token into a PowerShell variable (or read from secure storage)
$env:CR_PAT = "your_personal_access_token"

# Log in to GitHub Container Registry
$env:CR_PAT | docker login ghcr.io -u amaanslyf --password-stdin
```

#### Pull Published Images
```powershell
# Pull the latest main branch images
docker pull ghcr.io/amaanslyf/taskflow-backend:main
docker pull ghcr.io/amaanslyf/taskflow-frontend:main

# Pull an exact commit SHA build
docker pull ghcr.io/amaanslyf/taskflow-backend:sha-<commit-sha>
docker pull ghcr.io/amaanslyf/taskflow-frontend:sha-<commit-sha>

# Pull by cryptographic digest (true immutability)
docker pull ghcr.io/amaanslyf/taskflow-backend@sha256:<digest>
```

#### Inspect Published Image Metadata
```powershell
# Inspect image architecture, user, and OCI labels generated by Buildx
docker image inspect ghcr.io/amaanslyf/taskflow-backend:main --format '{{json .Config.Labels}}'
```

---

> **⚠️ Frontend Runtime Configuration Notice**: The frontend production image is currently built with build argument `VITE_API_BASE_URL=http://localhost:8080`. This configuration is designed for local development and Docker Compose testing where the browser accesses the host machine. In subsequent Kubernetes and Helm deployment stages, we will introduce runtime configuration injection or Nginx reverse-proxy routing so external browser clients can connect to remote cluster endpoints.

---

## Project Structure

```
taskflow/
├── .github/
│   └── workflows/
│       └── ci.yml               # GitHub Actions CI multi-job pipeline
├── compose.yaml                 # Declarative Docker Compose orchestration
├── .env.docker                  # Local Docker environment configuration
├── backups/                     # Database backups (ignored in Git)
│   └── taskflow_backup_*.sql
├── backend/
│   ├── cmd/api/main.go          # Application startup & graceful shutdown
│   ├── internal/
│   │   ├── config/              # Env loading & config_test.go
│   │   ├── database/            # pgxpool connection pool setup
│   │   ├── handlers/            # REST controllers, middleware & handlers_test.go
│   │   └── tasks/               # Data model, parameterized SQL & tasks_test.go
│   ├── Dockerfile               # Multi-stage: Go 1.27.4 Alpine -> Alpine 3.21 Runner
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
│   ├── Dockerfile               # Multi-stage: Node.js 22 builder -> Nginx 1.27 Alpine
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
