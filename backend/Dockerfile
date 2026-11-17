# ============================================================
# TaskFlow Backend — Multi-Stage Dockerfile
# ============================================================
# Multi-stage builds divide the container creation into stages:
#   Stage 1 (builder): Contains the heavy Go compiler and tools
#                      needed to build the source code into an
#                      executable binary.
#   Stage 2 (runtime): Contains only a minimal Linux OS and the
#                      compiled binary. The Go compiler, source
#                      code, and build tools are discarded, leaving
#                      a tiny, secure, production-ready image.
# ============================================================

# ── Stage 1: Build the Go binary ────────────────────────────
FROM golang:1.22-alpine AS builder

# Set the working directory inside the builder container.
WORKDIR /build

# Copy dependency definition files first.
# Docker caches image layers. By copying go.mod and go.sum
# before the rest of the code, Docker will only re-run
# "go mod download" if dependencies change — not every time
# you edit a .go source file.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the backend source code into the builder.
COPY . .

# Compile the application into a standalone static binary.
# Flags explained:
#   CGO_ENABLED=0 : Disables C bindings, generating a purely static
#                   binary that runs on any Linux distribution
#                   without requiring external C libraries (glibc/musl).
#   GOOS=linux    : Targets the Linux operating system kernel.
#   GOARCH=amd64  : Targets 64-bit x86 architecture.
#   -ldflags="-s -w" : Strips debug information and symbol tables,
#                      dramatically shrinking the binary size without
#                      altering application performance.
#   -o /app/api   : Output file destination for the compiled binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/api \
    ./cmd/api

# ── Stage 2: Minimal Runtime Image ──────────────────────────
FROM alpine:3.21 AS runner

# Install root CA certificates so the Go application can establish
# secure TLS/HTTPS connections (e.g., to cloud-managed PostgreSQL).
RUN apk --no-cache add ca-certificates tzdata

# Security best practice: Never run container processes as root.
# We create a dedicated non-root system user and group.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the operational directory inside the runtime container.
WORKDIR /app

# Copy ONLY the compiled binary from the builder stage.
# Everything else (Go SDK, compiler, caches, source files) is discarded.
COPY --from=builder /app/api /app/api

# Ensure the non-root user owns the executable.
RUN chown -R appuser:appgroup /app

# Switch execution context to the non-privileged user.
USER appuser

# Document that this container listens on port 8080.
EXPOSE 8080

# Container healthcheck: Periodically queries the /healthz endpoint.
# If the process crashes or hangs, Docker marks the container "unhealthy".
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

# Execute the application binary when the container launches.
ENTRYPOINT ["/app/api"]

