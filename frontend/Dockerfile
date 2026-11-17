# ============================================================
# TaskFlow Frontend — Multi-Stage Dockerfile
# ============================================================
# Why multi-stage for a frontend?
#   1. Node.js is only needed at BUILD TIME (to compile JSX,
#      bundle JavaScript modules, and optimize CSS).
#   2. The resulting output is purely STATIC files (HTML, CSS, JS).
#   3. In production, we don't need Node.js at all! We use
#      Nginx, an ultra-fast, lightweight web server written in C.
#
# Vite Environment Variables & Docker Networking:
#   Vite injects "import.meta.env.VITE_*" variables into the
#   JavaScript bundle at COMPILE TIME, replacing references with
#   raw string literals.
#   Importantly: This code executes inside the END-USER'S BROWSER,
#   not inside the container. The browser runs on the user's host
#   machine, which cannot resolve Docker's internal virtual bridge
#   names like "http://backend:8080". Therefore, VITE_API_BASE_URL
#   must point to a URL accessible from the client's browser
#   (default: "http://localhost:8080").
# ============================================================

# ── Stage 1: Build static assets with Node.js ──────────────
FROM node:22-alpine AS builder

WORKDIR /app

# Copy dependency manifests first for layer caching.
COPY package.json package-lock.json ./

# "npm ci" performs a clean install strictly based on package-lock.json.
# Faster and more deterministic than "npm install" in CI/CD.
RUN npm ci

# Declare build argument for API URL with default fallback.
ARG VITE_API_BASE_URL=http://localhost:8080
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

# Copy source code and configuration files.
COPY . .

# Run Vite build to generate the production bundle in /app/dist.
RUN npm run build

# ── Stage 2: Serve static files with Nginx Alpine ───────────
FROM nginx:1.27-alpine AS runner

# Copy custom Nginx configuration for React Single-Page Application routing.
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy compiled static assets from the builder stage into Nginx web root.
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose HTTP port 80.
EXPOSE 80

# Healthcheck: Periodically verify Nginx is actively serving the static React homepage.
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -q -O /dev/null http://127.0.0.1/ || exit 1

# Start Nginx in the foreground ("daemon off" prevents the container from exiting immediately).
CMD ["nginx", "-g", "daemon off;"]

