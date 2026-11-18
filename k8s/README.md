# TaskFlow — Kubernetes Deployment Guide

This directory contains the production-grade Kubernetes manifests for deploying the complete TaskFlow 3-tier application stack to a Kubernetes cluster (including local environments like Docker Desktop or Kind).

---

## 1. Architecture Overview

TaskFlow on Kubernetes consists of three interconnected components running in the dedicated `taskflow-lab` namespace:

```
[ User Browser ]
       │
       ▼ (kubectl port-forward / Ingress / NodePort)
┌─────────────────────────────────────────────────────────────────────────────┐
│ Kubernetes Namespace: taskflow-lab                                          │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │ Frontend Service (ClusterIP: port 80)                                 │  │
│  │   ▲                                                                   │  │
│  │   │                                                                   │  │
│  │  ┌┴─────────────────────────────┐                                     │  │
│  │  │ Frontend Pod (Nginx + React) │                                     │  │
│  │  │ - Serves static assets       │                                     │  │
│  │  │ - Proxies /api/ ─────────────┼────────┐                            │  │
│  │  └──────────────────────────────┘        │                            │  │
│  └──────────────────────────────────────────┼────────────────────────────┘  │
│                                             ▼                               │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │ Backend Service: "backend" (ClusterIP: port 8080)                     │  │
│  │   ▲                                                                   │  │
│  │   │                                                                   │  │
│  │  ┌┴─────────────────────────────┐                                     │  │
│  │  │ Backend Pod (Go REST API)    │                                     │  │
│  │  │ - Health probes (/healthz,   │                                     │  │
│  │  │   /readyz)                   │                                     │  │
│  │  │ - Non-root user (UID 100)    │                                     │  │
│  │  │ - Queries PostgreSQL ────────┼────────┐                            │  │
│  │  └──────────────────────────────┘        │                            │  │
│  └──────────────────────────────────────────┼────────────────────────────┘  │
│                                             ▼                               │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │ Database Service: "postgres" (ClusterIP: port 5432)                   │  │
│  │   ▲                                                                   │  │
│  │   │                                                                   │  │
│  │  ┌┴─────────────────────────────┐    ┌─────────────────────────────┐  │  │
│  │  │ Database Pod (PostgreSQL 17) │    │ PersistentVolumeClaim (PVC) │  │  │
│  │  │ - Recreate strategy          ├────┤ taskflow-postgres-pvc       │  │  │
│  │  │ - Auto-initializes schema    │    │ (1Gi, ReadWriteOnce)        │  │  │
│  │  └──────────────────────────────┘    └─────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Internal Service Discovery & Zero-Configuration Proxying
1. The **React SPA** runs client-side and makes relative requests to `/api/v1/tasks`.
2. The **Frontend Nginx server** receives `/api/` requests and forwards them to `http://backend:8080`.
3. In Kubernetes, the Backend Service is named `backend`. Kubernetes CoreDNS automatically resolves `backend` to the Backend Service ClusterIP. No environment variables or rebuilds are needed.
4. The **Go Backend** resolves `postgres:5432` to connect to PostgreSQL.

---

## 2. Manifest Structure

| File | Resource Kind | Purpose |
| :--- | :--- | :--- |
| `namespace.yaml` | `Namespace` | Creates the isolated `taskflow-lab` namespace. |
| `secret.example.yaml` | `Secret` | Template for database credentials (safe to commit). |
| `secret.yaml` | `Secret` | Local active secret with working credentials (gitignored). |
| `configmap.yaml` | `ConfigMap` | Non-sensitive backend configuration (`PORT`, `FRONTEND_ORIGIN`). |
| `postgres-init-configmap.yaml` | `ConfigMap` | Mounts `01-schema.sql` into PostgreSQL's `/docker-entrypoint-initdb.d/`. |
| `postgres-pvc.yaml` | `PersistentVolumeClaim` | Requests 1Gi persistent disk storage for database files. |
| `postgres-deployment.yaml` | `Deployment` | Runs PostgreSQL 17 with health checks and `Recreate` strategy. |
| `postgres-service.yaml` | `Service` | Exposes PostgreSQL internally on ClusterIP port 5432. |
| `backend-deployment.yaml` | `Deployment` | Runs Go REST API with GHCR image, non-root user, and health probes. |
| `backend-service.yaml` | `Service` | Exposes Go API internally on ClusterIP port 8080 (DNS name: `backend`). |
| `frontend-deployment.yaml` | `Deployment` | Runs Nginx + React UI with GHCR image and health checks. |
| `frontend-service.yaml` | `Service` | Exposes Frontend UI internally on ClusterIP port 80. |
| `kustomization.yaml` | `Kustomization` | Bundles all resources for declarative, ordered deployment. |

---

## 3. Prerequisites

- A running Kubernetes cluster (Docker Desktop Kubernetes, Kind, Minikube, or K3s).
- `kubectl` configured with cluster access:
  ```powershell
  kubectl cluster-info
  kubectl get nodes
  ```
- Container images published to GHCR or cached on cluster nodes:
  - `ghcr.io/amaanslyf/taskflow-frontend:sha-a94e79db752719a8b06dee5998e70f747addfbf4`
  - `ghcr.io/amaanslyf/taskflow-backend:sha-a94e79db752719a8b06dee5998e70f747addfbf4`
  - `postgres:17-alpine`

---

## 4. Step-by-Step Deployment

### Step 1: Prepare Secret Manifest
`k8s/secret.example.yaml` is tracked in git as a template. Copy it to create your active `k8s/secret.yaml` (which is excluded from Git via `.gitignore`):

```powershell
# Copy the template
cp k8s/secret.example.yaml k8s/secret.yaml
```

If deploying locally with default dev credentials, `k8s/secret.yaml` is already pre-configured:
- `POSTGRES_USER`: `postgres`
- `POSTGRES_PASSWORD`: `taskflow_local_dev_only`
- `POSTGRES_DB`: `taskflow`
- `DATABASE_URL`: `postgres://postgres:taskflow_local_dev_only@postgres:5432/taskflow?sslmode=disable`

### Step 2: Apply Manifests
Deploy all components declaratively using Kustomize:

```powershell
kubectl apply -k k8s/
```

*Alternatively, apply individually in order:*
```powershell
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/postgres-init-configmap.yaml
kubectl apply -f k8s/postgres-pvc.yaml
kubectl apply -f k8s/postgres-service.yaml
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/backend-service.yaml
kubectl apply -f k8s/backend-deployment.yaml
kubectl apply -f k8s/frontend-service.yaml
kubectl apply -f k8s/frontend-deployment.yaml
```

### Step 3: Verify Rollout Status
Wait for all pods to become `Ready (1/1)`:

```powershell
kubectl get pods,pvc,svc -n taskflow-lab -o wide
```

Expected output:
```
NAME                            READY   STATUS    RESTARTS   AGE
pod/backend-5ff8f6cc8b-9nbc8    1/1     Running   0          2m
pod/frontend-567d8794f7-jxxbw   1/1     Running   0          2m
pod/postgres-5bf7f45776-tq4gv   1/1     Running   0          2m

NAME                                          STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS
persistentvolumeclaim/taskflow-postgres-pvc   Bound    pvc-d2a8c9ad-ed63-4bcd-87fc-a15a8209c381   1Gi        RWO            standard

NAME               TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
service/backend    ClusterIP   10.96.69.189   <none>        8080/TCP   2m
service/frontend   ClusterIP   10.96.122.12   <none>        80/TCP     2m
service/postgres   ClusterIP   10.96.225.37   <none>        5432/TCP   2m
```

---

## 5. Accessing Services (Port-Forwarding)

Because all Services use `ClusterIP` (best practice for private internal networking without exposing unnecessary host ports), we use `kubectl port-forward` to access them securely from the host machine:

> [!NOTE]
> If Docker Compose is running on ports 3000, 8080, and 5433, use alternative host ports (e.g. 3001, 8081, 5434) to avoid port collisions.

### 1. Access the React Frontend UI
Forward the frontend service to host port `3001`:
```powershell
kubectl port-forward svc/frontend 3001:80 -n taskflow-lab
```
Open your browser and navigate to:
```
http://localhost:3001
```

### 2. Access the Backend API Directly
Forward the backend service to host port `8081`:
```powershell
kubectl port-forward svc/backend 8081:8080 -n taskflow-lab
```
Test health probes:
```powershell
curl http://localhost:8081/healthz
curl http://localhost:8081/readyz
```

### 3. Access PostgreSQL Directly
Forward PostgreSQL to host port `5434`:
```powershell
kubectl port-forward svc/postgres 5434:5432 -n taskflow-lab
```
Connect via `psql`:
```powershell
psql -h localhost -p 5434 -U postgres -d taskflow
```

---

## 6. Health Probes & Verification

### Probes Implementation Details
- **Backend Liveness Probe (`GET /healthz`)**:
  Validates that the Go HTTP server process is running and actively responding to HTTP requests.
- **Backend Readiness Probe (`GET /readyz`)**:
  Executes an active `pool.Ping(ctx)` to PostgreSQL. If the database is initializing, restarting, or network is disconnected, `/readyz` returns HTTP 503 (`Service Unavailable`), preventing Kubernetes from routing incoming traffic to the backend until PostgreSQL is ready.
- **PostgreSQL Probes (`pg_isready`)**:
  Executes `pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"` inside the container to ensure the database engine is accepting client connections.
- **Frontend Probes (`GET /`)**:
  Verifies Nginx is serving the SPA index page on port 80.

### Testing CRUD Operations via Frontend Proxy
Through the port-forward (`http://localhost:3001`):

#### 1. List Tasks (GET)
```bash
curl -X GET http://localhost:3001/api/v1/tasks
```

#### 2. Create a Task (POST)
```bash
curl -X POST http://localhost:3001/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Deploy to Kubernetes", "description": "Complete manifests and verify persistence", "status": "in_progress"}'
```

#### 3. Update a Task (PUT)
```bash
curl -X PUT http://localhost:3001/api/v1/tasks/<TASK_ID> \
  -H "Content-Type: application/json" \
  -d '{"title": "Deploy to Kubernetes", "description": "Completed manifests and verified persistence", "status": "completed"}'
```

#### 4. Delete a Task (DELETE)
```bash
curl -X DELETE http://localhost:3001/api/v1/tasks/<TASK_ID>
```

---

## 7. Verifying PostgreSQL Data Persistence

PostgreSQL data is stored in the `PersistentVolumeClaim` (`taskflow-postgres-pvc`), backed by the cluster's local storage provider. The PVC lifecycle is independent of the PostgreSQL Pod.

To verify persistence:

1. Create a task via API:
   ```bash
   curl -X POST http://localhost:3001/api/v1/tasks \
     -H "Content-Type: application/json" \
     -d '{"title": "Persistence Test Task", "description": "Must survive Pod deletion", "status": "todo"}'
   ```

2. Delete the running PostgreSQL Pod:
   ```powershell
   kubectl delete pod -l app.kubernetes.io/name=postgres -n taskflow-lab
   ```

3. Watch the Kubernetes Deployment automatically spawn a replacement Pod:
   ```powershell
   kubectl get pods -l app.kubernetes.io/name=postgres -n taskflow-lab -w
   ```

4. Inspect the new Pod's logs:
   ```powershell
   kubectl logs -l app.kubernetes.io/name=postgres -n taskflow-lab
   ```
   You will observe:
   ```
   PostgreSQL Database directory appears to contain a database; Skipping initialization
   LOG: database system is ready to accept connections
   ```

5. Fetch the tasks again:
   ```bash
   curl -X GET http://localhost:3001/api/v1/tasks
   ```
   The `Persistence Test Task` is still present, proving that persistent data was completely preserved!

---

## 8. Teardown and Cleanup

To delete all TaskFlow Kubernetes resources from the cluster:

```powershell
# Delete all resources defined in k8s manifests
kubectl delete -k k8s/

# Or delete the entire namespace (deletes all resources within it)
kubectl delete namespace taskflow-lab
```

> [!NOTE]
> Deleting Kubernetes resources does NOT affect or delete the Docker Compose containers or volumes (`taskflow-postgres-data`). The two environments are completely isolated.
