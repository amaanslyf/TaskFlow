# TaskFlow Helm Chart

TaskFlow is a production-grade 3-tier task management application packaged as a Kubernetes Helm chart. This chart automates the deployment, scaling, configuration, and traffic routing for the TaskFlow web application on Kubernetes.

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [What is Helm? (Beginner-Friendly Concepts)](#what-is-helm-beginner-friendly-concepts)
3. [Chart Structure](#chart-structure)
4. [Prerequisites](#prerequisites)
5. [Configuration & Values](#configuration--values)
   - [Environment Profiles (Development, Staging, Production)](#environment-profiles)
   - [Database Secrets Management](#database-secrets-management)
   - [Persistent Storage & Data Safety](#persistent-storage--data-safety)
6. [Gateway API & Ingress Routing](#gateway-api--ingress-routing)
7. [Installation & Upgrade Guide](#installation--upgrade-guide)
   - [Fresh Installation](#fresh-installation)
   - [Safe Migration from Existing Manifests](#safe-migration-from-existing-manifests)
   - [Upgrading Releases](#upgrading-releases)
   - [Uninstalling](#uninstalling)
8. [Troubleshooting & Verification](#troubleshooting--verification)

---

## Architecture Overview

The TaskFlow Helm chart deploys the following microservices and resources:

```
[ External HTTP Traffic (Port 80) ]
               │
               ▼
   [ Gateway API / Traefik Proxy ]
               │
               ▼ (HTTPRoute: taskflow-frontend-route)
    [ Frontend Pod (React + Nginx) ]
        ├── Serves React Single Page Application (SPA)
        └── Proxies /api/ requests internally to backend:8080
               │
               ▼ (ClusterIP Service: backend:8080)
    [ Backend Pod (Go REST API) ]
        ├── /healthz & /readyz liveness/readiness probes
        └── Connects via TCP to postgres:5432
               │
               ▼ (ClusterIP Service: postgres:5432)
    [ Database Pod (PostgreSQL 17) ]
        ├── Bound to PersistentVolumeClaim (taskflow-postgres-pvc)
        └── Initialized with schema via taskflow-postgres-init ConfigMap
```

### Components

| Component | Technology | Default Image Tag / Spec | Service Port |
| :--- | :--- | :--- | :--- |
| **Frontend** | React, Vite, Nginx Alpine | `ghcr.io/amaanslyf/taskflow-frontend:sha-a94e79db752719a8b06dee5998e70f747addfbf4` | `80` (HTTP) |
| **Backend** | Go REST API (`net/http`) | `ghcr.io/amaanslyf/taskflow-backend:sha-a94e79db752719a8b06dee5998e70f747addfbf4` | `8080` (HTTP) |
| **PostgreSQL** | PostgreSQL 17 Alpine | `postgres:17-alpine` | `5432` (TCP) |
| **Gateway** | Gateway API (`gateway.networking.k8s.io/v1`) | Managed by Traefik Proxy Controller | `80` (Gateway) |

---

## What is Helm? (Beginner-Friendly Concepts)

If you are new to Helm, think of Helm as the **package manager for Kubernetes** (similar to `apt` on Ubuntu, `brew` on macOS, or `npm` in Node.js).

### Core Concepts

1. **Chart**: A bundle of YAML template files and configuration files that describe a related set of Kubernetes resources. This directory (`helm/taskflow/`) is a Helm chart.
2. **`Chart.yaml`**: The manifest file containing metadata about the chart, such as its name, version (`0.1.0`), application version (`1.0.0`), description, and maintainer details.
3. **`values.yaml`**: The default configuration values for the chart. Instead of hardcoding values like replica counts, ports, or image tags directly into Kubernetes YAML manifests, templates reference variables like `{{ .Values.frontend.replicaCount }}`.
4. **Templates (`templates/`)**: Directory containing parameterized Kubernetes manifest templates written in Go templating syntax. When you run Helm, it combines `values.yaml` with these templates to produce standard Kubernetes YAML manifests.
5. **Release**: A running instance of a chart in a Kubernetes cluster. You can install the same chart multiple times with different release names (e.g., `taskflow-dev` and `taskflow-staging`).
6. **`helm install` vs `helm upgrade`**:
   - `helm install <release-name> <chart-path>`: Creates new Kubernetes resources from the chart templates for the first time.
   - `helm upgrade <release-name> <chart-path>`: Updates an existing release when configuration values, templates, or image tags change.
   - `helm upgrade --install`: A safe idempotency command that installs if the release does not exist or upgrades if it already does.

---

## Chart Structure

```
helm/taskflow/
├── .helmignore                      # Files and directories ignored by Helm packaging
├── Chart.yaml                       # Chart metadata and semantic versioning
├── README.md                        # Documentation and user guide
├── values.yaml                      # Default configuration values
├── values.schema.json               # JSON Schema (Draft 7) for strict values validation
├── environments/                    # Environment-specific configuration overlays
│   ├── development.yaml             # Minimal resources, single replicas, DEBUG logging
│   ├── staging.yaml                 # Multi-replica backend, moderate resource limits
│   └── production.yaml              # High availability, strict resource requests/limits, INFO logging
└── templates/                       # Kubernetes YAML templates
    ├── _helpers.tpl                 # Reusable template partials (labels, names, selector labels)
    ├── NOTES.txt                    # Post-installation instructions displayed in terminal
    ├── configmap.yaml               # Application environment ConfigMap (PORT, LOG_LEVEL, etc.)
    ├── postgres-init-configmap.yaml # Database schema initialization script (init.sql)
    ├── postgres-pvc.yaml            # PersistentVolumeClaim for database storage
    ├── postgres-service.yaml        # ClusterIP service for PostgreSQL
    ├── postgres-deployment.yaml     # Deployment specification for PostgreSQL
    ├── backend-service.yaml         # ClusterIP service for Go backend
    ├── backend-deployment.yaml      # Deployment specification for Go backend
    ├── frontend-service.yaml        # ClusterIP service for React frontend
    ├── frontend-deployment.yaml     # Deployment specification for React frontend
    ├── gateway.yaml                 # Gateway API Gateway resource (Traefik GatewayClass)
    └── httproute.yaml               # Gateway API HTTPRoute routing port 80 to frontend service
```

---

## Prerequisites

Before deploying the TaskFlow Helm chart, ensure your local or remote Kubernetes environment meets the following requirements:

1. **Kubernetes Cluster**: Kubernetes `v1.28+` (Tested on Docker Desktop Kubernetes / Kind `v1.36.1`).
2. **Helm CLI**: Helm `v3.12+` installed and available in your `PATH`.
3. **Gateway API CRDs**: Standard channel CRDs (`v1.0.0+` or `v1.7.0`) installed in the cluster.
4. **Gateway Controller**: Traefik Proxy installed (e.g., in namespace `traefik-system`) with GatewayClass `traefik-gateway-class`.
5. **Pre-existing Kubernetes Secret**: A Secret containing the database credentials must exist prior to deploying the chart (see [Database Secrets Management](#database-secrets-management)).

---

## Configuration & Values

All configurable parameters are documented in `helm/taskflow/values.yaml` and validated by `helm/taskflow/values.schema.json`.

### Key Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `global.namespace` | `string` | `taskflow-lab` | Target namespace for application deployment |
| `frontend.replicaCount` | `integer` | `1` | Number of frontend Pod replicas |
| `frontend.image.repository` | `string` | `ghcr.io/amaanslyf/taskflow-frontend` | Container image repository |
| `frontend.image.tag` | `string` | `sha-a94e79db752719a8b06dee5998e70f747addfbf4` | Container image tag |
| `backend.replicaCount` | `integer` | `1` | Number of backend Pod replicas |
| `backend.image.repository` | `string` | `ghcr.io/amaanslyf/taskflow-backend` | Container image repository |
| `backend.image.tag` | `string` | `sha-a94e79db752719a8b06dee5998e70f747addfbf4` | Container image tag |
| `backend.config.logLevel` | `string` | `INFO` | Application log level (`DEBUG`, `INFO`, `WARN`, `ERROR`) |
| `postgres.storage.size` | `string` | `1Gi` | Persistent disk allocation for database |
| `postgres.persistence.preserveOnDelete` | `boolean` | `true` | Retain PVC if Helm release is deleted (`helm.sh/resource-policy: keep`) |
| `secret.existingSecret` | `string` | `taskflow-secrets` | Name of pre-created Kubernetes Secret containing database credentials |
| `gateway.enabled` | `boolean` | `true` | Deploy Gateway API Gateway and HTTPRoute |
| `gateway.gatewayClassName` | `string` | `traefik-gateway-class` | Name of the cluster GatewayClass |

---

### Environment Profiles

We provide pre-tuned value files in `helm/taskflow/environments/`:

- **Development (`environments/development.yaml`)**:
  - Replicas: 1 for frontend, 1 for backend.
  - Log Level: `DEBUG`.
  - Minimal resource requests/limits to conserve host CPU and memory.
- **Staging (`environments/staging.yaml`)**:
  - Replicas: 2 for frontend, 2 for backend.
  - Log Level: `INFO`.
  - Realistic resource allocations matching staging workload needs.
- **Production (`environments/production.yaml`)**:
  - Replicas: 3 for frontend, 3 for backend.
  - Log Level: `INFO`.
  - Conservative memory and CPU limits, pod anti-affinity ready.

---

### Database Secrets Management

> [!IMPORTANT]
> **Zero Secrets in Git**: Plaintext passwords or sensitive database URLs must never be committed to Git or stored in `values.yaml`.

This chart expects a pre-created Kubernetes Secret named `taskflow-secrets` (or custom name configured in `secret.existingSecret`).

The secret must contain the following keys:
- `POSTGRES_USER`: The PostgreSQL database username.
- `POSTGRES_PASSWORD`: The PostgreSQL database password.
- `POSTGRES_DB`: The database name (e.g., `taskflow`).
- `DATABASE_URL`: Full connection string (e.g., `postgres://<user>:<password>@postgres:5432/<db>?sslmode=disable`).

To create this secret manually in your namespace:
```powershell
kubectl create secret generic taskflow-secrets `
  --namespace taskflow-lab `
  --from-literal=POSTGRES_USER="taskflow_user" `
  --from-literal=POSTGRES_PASSWORD="YourStrongPasswordHere" `
  --from-literal=POSTGRES_DB="taskflow" `
  --from-literal=DATABASE_URL="postgres://taskflow_user:YourStrongPasswordHere@postgres:5432/taskflow?sslmode=disable"
```

---

### Persistent Storage & Data Safety

The PostgreSQL database relies on a `PersistentVolumeClaim` (PVC) named `taskflow-postgres-pvc` backed by a host-path PersistentVolume (`rancher.io/local-path`).

#### Understanding Storage Protection Layers:
1. **Helm Resource Policy (`helm.sh/resource-policy: keep`)**:
   - Declared on `PersistentVolumeClaim/taskflow-postgres-pvc`.
   - Tells Helm **not** to issue a Kubernetes API `DELETE` call for this PVC when running `helm uninstall taskflow`. The PVC remains bound in Kubernetes after the release is uninstalled.
2. **Kubernetes PV Reclaim Policy (`persistentVolumeReclaimPolicy`)**:
   - The cluster's dynamically provisioned PV currently has `persistentVolumeReclaimPolicy: Delete`.
   - **Crucial Distinction**: If an administrator or script deletes the Kubernetes PVC via `kubectl delete pvc taskflow-postgres-pvc`, Kubernetes will immediately delete the PV and wipe the physical directory on disk (`/var/local-path-provisioner/...`)!
   - Helm's `resource-policy: keep` only protects against `helm uninstall`. It does **not** alter the PV reclaim policy in Kubernetes. To protect against accidental `kubectl delete pvc`, the PV's reclaim policy can be patched to `Retain`.
3. **Explicit Creation vs. Reuse**:
   - Setting `postgres.persistence.existingClaim: ""` (default): Helm renders and manages the PVC (`taskflow-postgres-pvc`).
   - Setting `postgres.persistence.existingClaim: "<name>"`: Helm skips PVC creation entirely and attaches the database deployment directly to the external PVC.

---

## Gateway API & Ingress Routing

The chart implements modern **Kubernetes Gateway API** (`gateway.networking.k8s.io/v1`) using Traefik:

- **Gateway (`templates/gateway.yaml`)**:
  - Name: `taskflow-gateway`
  - GatewayClass: `traefik-gateway-class`
  - Listener: HTTP on Port `80` allowing routes from the same namespace (`Same`).
- **HTTPRoute (`templates/httproute.yaml`)**:
  - Name: `taskflow-frontend-route`
  - Parent: `taskflow-gateway`
  - Match Rule: PathPrefix `/` forwards all incoming web traffic to the `frontend` Service on Port `80`.
  - Internal API Proxying: The React frontend uses an embedded Nginx reverse proxy configured to forward `/api/` calls directly to `http://backend:8080/api/` inside the Kubernetes network.

---

## Installation & Upgrade Guide

### 1. Validate the Chart

Always lint and dry-run template rendering before installing:

```powershell
# Lint the chart against syntax errors and JSON schema validation
helm lint .\helm\taskflow

# Render templates locally to inspect generated Kubernetes YAML
helm template taskflow .\helm\taskflow --namespace taskflow-lab

# Test with environment-specific values
helm template taskflow .\helm\taskflow -f .\helm\taskflow\environments\development.yaml --namespace taskflow-lab
```

---

### 2. Fresh Installation

To install TaskFlow into a fresh cluster or clean namespace:

```powershell
# 1. Create namespace if not present
kubectl create namespace taskflow-lab --dry-run=client -o yaml | kubectl apply -f -

# 2. Ensure database secret exists
kubectl create secret generic taskflow-secrets `
  --namespace taskflow-lab `
  --from-literal=POSTGRES_USER="taskflow" `
  --from-literal=POSTGRES_PASSWORD="securepassword" `
  --from-literal=POSTGRES_DB="taskflow" `
  --from-literal=DATABASE_URL="postgres://taskflow:securepassword@postgres:5432/taskflow?sslmode=disable"

# 3. Install the Helm release
helm install taskflow .\helm\taskflow `
  --namespace taskflow-lab `
  --values .\helm\taskflow\environments\development.yaml
```

---

### 3. Safe Migration from Existing Manifests

If TaskFlow is already running in your cluster (deployed via standalone manifests / Kustomize) and you want Helm to take over management of all resources **without downtime and without deleting or recreating the PostgreSQL PVC**:

The live namespace contains 10 application resources rendered by the chart:
- Deployments: `frontend`, `backend`, `postgres`
- Services: `frontend`, `backend`, `postgres`
- ConfigMaps: `taskflow-config`, `taskflow-postgres-init`
- PersistentVolumeClaim: `taskflow-postgres-pvc`
- Gateway: `taskflow-gateway`
- HTTPRoute: `taskflow-frontend-route`

*(Note: Secret `taskflow-secrets` and cluster GatewayClass `traefik-gateway-class` remain externally managed).*

#### Recommended Method: Helm v4 Native Adoption (`--take-ownership`)

Helm v4 provides a native `--take-ownership` flag that allows Helm to adopt existing Kubernetes resources into the release without manual labeling:

```powershell
# Step 1: Verify current workloads and PVC status
kubectl get pods,svc,pvc,gateway,httproute -n taskflow-lab

# Step 2: Validate with server-side dry run
helm upgrade --install taskflow .\helm\taskflow `
  --namespace taskflow-lab `
  --values .\helm\taskflow\environments\development.yaml `
  --take-ownership `
  --dry-run=server

# Step 3: Execute the live adoption
helm upgrade --install taskflow .\helm\taskflow `
  --namespace taskflow-lab `
  --values .\helm\taskflow\environments\development.yaml `
  --take-ownership
```

#### Alternative Method: Explicit Manifest Annotation

If adopting resources explicitly or in environments requiring pre-approved change scripts:

```powershell
# Annotate and label all 10 resources for release 'taskflow' in namespace 'taskflow-lab':
$resources = @(
  "deployment/frontend", "deployment/backend", "deployment/postgres",
  "service/frontend", "service/backend", "service/postgres",
  "configmap/taskflow-config", "configmap/taskflow-postgres-init",
  "pvc/taskflow-postgres-pvc",
  "gateway.gateway.networking.k8s.io/taskflow-gateway",
  "httproute.gateway.networking.k8s.io/taskflow-frontend-route"
)

foreach ($r in $resources) {
  kubectl annotate $r -n taskflow-lab meta.helm.sh/release-name=taskflow --overwrite
  kubectl annotate $r -n taskflow-lab meta.helm.sh/release-namespace=taskflow-lab --overwrite
  kubectl label $r -n taskflow-lab app.kubernetes.io/managed-by=Helm --overwrite
}

# Ensure PVC keeps retention policy
kubectl annotate pvc taskflow-postgres-pvc -n taskflow-lab helm.sh/resource-policy=keep --overwrite

# Then execute standard upgrade
helm upgrade --install taskflow .\helm\taskflow `
  --namespace taskflow-lab `
  --values .\helm\taskflow\environments\development.yaml
```

---

### 4. Upgrading Releases

When upgrading images or application settings:

```powershell
# Upgrade using production values
helm upgrade taskflow .\helm\taskflow `
  --namespace taskflow-lab `
  --values .\helm\taskflow\environments\production.yaml

# Rollback to the previous revision if needed
helm rollback taskflow 1 --namespace taskflow-lab
```

---

### 5. Uninstalling

```powershell
# Uninstall the release
helm uninstall taskflow --namespace taskflow-lab

# Note: The database PVC (taskflow-postgres-pvc) will be PRESERVED due to 'helm.sh/resource-policy: keep'.
```

---

## Troubleshooting & Verification

### Check Status of Helm Release
```powershell
helm status taskflow -n taskflow-lab
helm list -n taskflow-lab
```

### Inspect Pods & Services
```powershell
kubectl get pods,svc,gateway,httproute -n taskflow-lab
```

### Inspect Gateway API Listener & Routing
```powershell
kubectl describe gateway taskflow-gateway -n taskflow-lab
kubectl describe httproute taskflow-frontend-route -n taskflow-lab
```

### Test External Traffic
```powershell
# Send request to Traefik Gateway Port 80
curl -i http://localhost/healthz
curl -i http://localhost/api/v1/tasks
```
