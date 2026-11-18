# TaskFlow — Kubernetes Gateway API with Traefik Proxy

This directory contains the Gateway API manifests and Helm values for providing external HTTP ingress to the TaskFlow application using **Traefik Proxy** and the **Kubernetes Gateway API**.

---

## 1. Gateway API Architecture Overview

```
                      [ Host Browser / Client ]
                                  │
                                  ▼ (kubectl port-forward :8088 -> Traefik :80)
┌─────────────────────────────────────────────────────────────────────────────┐
│ Namespace: traefik-system                                                   │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │ Service: traefik (ClusterIP: port 80)                                 │  │
│  │   ▲                                                                   │  │
│  │   │                                                                   │  │
│  │  ┌┴─────────────────────────────┐                                     │  │
│  │  │ Pod: traefik                 │ (Entrypoint: web :80)               │  │
│  │  │ - Gateway API Provider       │                                     │  │
│  │  │ - Non-root (UID 65532)       │                                     │  │
│  │  └──────────────┬───────────────┘                                     │  │
│  └─────────────────┼─────────────────────────────────────────────────────┘  │
└────────────────────┼────────────────────────────────────────────────────────┘
                     │
                     │  GatewayClass: traefik-gateway-class
                     │  Gateway:      taskflow-gateway (port 80)
                     │  HTTPRoute:    taskflow-frontend-route (PathPrefix: /)
                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ Namespace: taskflow-lab                                                     │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │ Service: frontend (ClusterIP: port 80)                                │  │
│  │   ▲                                                                   │  │
│  │   │                                                                   │  │
│  │  ┌┴─────────────────────────────┐                                     │  │
│  │  │ Pod: frontend (Nginx + React)│                                     │  │
│  │  │ - Serves static SPA UI (/)   │                                     │  │
│  │  │ - Proxies /api/ ─────────────┼────────┐                            │  │
│  │  └──────────────────────────────┘        │                            │  │
│  └──────────────────────────────────────────┼────────────────────────────┘  │
│                                             ▼                               │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │ Service: backend (ClusterIP: port 8080)                               │  │
│  │   ▲                                                                   │  │
│  │   │                                                                   │  │
│  │  ┌┴─────────────────────────────┐                                     │  │
│  │  │ Pod: backend (Go REST API)   │                                     │  │
│  │  │ - Task endpoints & probes    │                                     │  │
│  │  │ - Queries PostgreSQL ────────┼────────┐                            │  │
│  │  └──────────────────────────────┘        │                            │  │
│  └──────────────────────────────────────────┼────────────────────────────┘  │
│                                             ▼                               │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │ Service: postgres (ClusterIP: port 5432)                              │  │
│  │   ▲                                                                   │  │
│  │   │                                                                   │  │
│  │  ┌┴─────────────────────────────┐    ┌─────────────────────────────┐  │  │
│  │  │ Pod: postgres (PostgreSQL 17)├────┤ PVC: taskflow-postgres-pvc  │  │  │
│  │  └──────────────────────────────┘    └─────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Core Concepts: Gateway API vs. Legacy Ingress

### Why Gateway API Instead of the Retired Ingress NGINX Project?
- **Ingress Limitations**: The original Kubernetes `Ingress` resource (v1) was created in 2015. It combined infrastructure provisioning, TLS, host matching, and routing rules into a single monolithic object. Because it lacked standardization, every ingress controller invented custom vendor annotations (e.g. `nginx.ingress.kubernetes.io/rewrite-target`), making ingress configurations fragile, non-portable, and prone to vendor lock-in. Furthermore, the community `ingress-nginx` project is in maintenance/retirement mode.
- **Gateway API**: Gateway API is the modern, role-oriented, official Kubernetes standard for service networking. It splits networking into clean, decoupled roles:
  1. **Infrastructure Provider** defines `GatewayClass` (e.g. Traefik, Envoy, Istio).
  2. **Cluster / DevOps Operator** provisions the `Gateway` (listeners, ports, TLS certificates).
  3. **Application Developer** creates `HTTPRoute` objects attached to the Gateway to route traffic to their specific Services.

### Gateway API Resource Roles

| Resource | Scope | Role | Responsibility |
| :--- | :--- | :--- | :--- |
| **`GatewayClass`** | Cluster-scoped | Infrastructure | Declares which controller (e.g., `traefik.io/gateway-controller`) implements the gateways. |
| **`Gateway`** | Namespace-scoped | DevOps / Ops | Declares a point where traffic enters the cluster: IP address, listening ports (e.g. port 80), and allowed route namespaces. |
| **`HTTPRoute`** | Namespace-scoped | App Developer | Defines HTTP routing rules (path matches, headers, rewrite rules) that forward traffic from a Gateway to a target Kubernetes Service. |

### How Gateway API Connects to Existing TaskFlow Services
1. Traefik watches the `Gateway` and `HTTPRoute` resources using the standard Gateway API CRDs.
2. The `HTTPRoute` matches any request matching path prefix `/` and targets the existing `frontend` Service on port 80.
3. The frontend's Nginx container serves the static React SPA files. When the browser makes requests to `/api/v1/tasks`, Nginx proxies those requests internally to `http://backend:8080`, which connects to PostgreSQL.
4. **Existing Services, Deployments, and Databases remain 100% untouched.**

---

## 3. Helm Fundamentals

### What is Helm?
**Helm** is the package manager for Kubernetes (analogous to `npm` for Node.js, `pip` for Python, or `apt` for Debian/Ubuntu).

### What is a Helm Chart?
A **Helm Chart** is a collection of versioned, templated YAML files that describe a related set of Kubernetes resources. Instead of maintaining hundreds of raw manifest lines for deployments, RBAC cluster roles, service accounts, and services, a Helm chart packages them together and allows configuring them using a single `values.yaml` file.

### Key Helm Commands

| Command | Purpose | Example |
| :--- | :--- | :--- |
| `helm repo add` | Registers a remote chart repository. | `helm repo add traefik https://traefik.github.io/charts` |
| `helm repo update` | Fetches the latest chart metadata. | `helm repo update` |
| `helm search repo` | Lists available chart versions. | `helm search repo traefik/traefik --versions` |
| `helm install` | Deploys a new chart release into the cluster. | `helm install traefik traefik/traefik --version 42.4.1 ...` |
| `helm upgrade` | Applies changes or updates to an existing release. | `helm upgrade traefik traefik/traefik -f values.yaml` |
| `helm list` | Lists all installed releases across namespaces. | `helm list -A` |
| `helm uninstall` | Deletes a release and its resources. | `helm uninstall traefik -n traefik-system` |

---

## 4. Manifest Inventory in `k8s/gateway/`

| File | Resource Kind | Description |
| :--- | :--- | :--- |
| `traefik-values.yaml` | Helm Values | Pinned Helm configuration for Traefik with Gateway API provider enabled, ClusterIP service, and non-root securityContext. |
| `gatewayclass.yaml` | `GatewayClass` | Declares `traefik-gateway-class` referencing controller `traefik.io/gateway-controller`. |
| `gateway.yaml` | `Gateway` | Declares `taskflow-gateway` in namespace `taskflow-lab` with an HTTP listener on port 80. |
| `httproute.yaml` | `HTTPRoute` | Attaches to `taskflow-gateway`, matching `PathPrefix: /` and forwarding to Service `frontend:80`. |
| `kustomization.yaml` | `Kustomization` | Bundles Gateway API manifests for declarative one-command application. |
| `README.md` | Documentation | This comprehensive educational guide. |

---

## 5. Step-by-Step Installation & Deployment

### Step 1: Install Standard Gateway API CRDs
The Gateway API Custom Resource Definitions must be installed before starting the Traefik Gateway provider. We install standard channel release **v1.7.0**:

```powershell
kubectl apply --server-side -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.7.0/standard-install.yaml
```

Verify the CRDs are present:
```powershell
kubectl get crds | findstr "gateway.networking.k8s.io"
```

### Step 2: Install Traefik with Pinned Helm Chart
Add the repository and install the pinned chart version (`42.4.1` with Traefik Proxy `v3.7.19`) into the `traefik-system` namespace:

```powershell
# Add repository
helm repo add traefik https://traefik.github.io/charts
helm repo update

# Install pinned release with custom values
helm install traefik traefik/traefik `
  --version 42.4.1 `
  --namespace traefik-system `
  --create-namespace `
  -f k8s/gateway/traefik-values.yaml
```

Wait for the Traefik Pod to become `Ready (1/1)`:
```powershell
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=traefik -n traefik-system --timeout=90s
```

### Step 3: Apply Gateway API Resources
Deploy the `GatewayClass`, `Gateway`, and `HTTPRoute`:

```powershell
kubectl apply -k k8s/gateway/
```

Verify resource acceptance and status:
```powershell
# Check GatewayClass (should show ACCEPTED: True)
kubectl get gatewayclass

# Check Gateway (should show PROGRAMMED: True)
kubectl get gateway -n taskflow-lab

# Check HTTPRoute
kubectl get httproute -n taskflow-lab
```

---

## 6. Accessing & Testing the Gateway

Because the Traefik Service is deployed as `ClusterIP` (avoiding host port binding or external load balancer dependencies), access the gateway via `kubectl port-forward`:

```powershell
# Forward local port 8088 to Traefik's port 80:
kubectl port-forward svc/traefik 8088:80 -n traefik-system
```

Open a browser or terminal and test:

### 1. Test Frontend SPA UI
Navigate to:
```
http://localhost:8088
```
Expected: The TaskFlow React application UI loads successfully.

### 2. Test API CRUD Operations
```bash
# List tasks
curl -X GET http://localhost:8088/api/v1/tasks

# Create a task
curl -X POST http://localhost:8088/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Test via Gateway API", "description": "Routing verified", "status": "todo"}'

# Delete a task
curl -X DELETE http://localhost:8088/api/v1/tasks/<TASK_ID>
```

---

## 7. Teardown Instructions

To remove Gateway API resources and Traefik:

```powershell
# Delete TaskFlow Gateway and HTTPRoute
kubectl delete -k k8s/gateway/

# Uninstall Traefik Helm release
helm uninstall traefik -n traefik-system
kubectl delete namespace traefik-system

# Optionally delete Gateway API CRDs
kubectl delete -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.7.0/standard-install.yaml
```
*(TaskFlow deployments, services, and PostgreSQL PVC in `taskflow-lab` remain completely untouched).*
