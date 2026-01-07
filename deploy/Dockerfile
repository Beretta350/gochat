# ===========================================
# Stage 1: Build Go backend
# ===========================================
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source
COPY backend/ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o gochat ./cmd/main.go

# ===========================================
# Stage 2: Build Next.js frontend
# ===========================================
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# Copy package files
COPY frontend/package.json frontend/package-lock.json* ./

# Install dependencies
RUN npm ci

# Copy source
COPY frontend/ ./

# Build
RUN npm run build

# ===========================================
# Stage 3: Production image
# ===========================================
FROM alpine:3.19

# Install nginx, nodejs, and supervisor
RUN apk add --no-cache nginx nodejs supervisor

# Create directories
RUN mkdir -p /app/backend /app/frontend /run/nginx /var/log/supervisor

# Copy Go binary
COPY --from=backend-builder /app/gochat /app/backend/gochat

# Copy Next.js standalone build
COPY --from=frontend-builder /app/.next/standalone /app/frontend/
COPY --from=frontend-builder /app/.next/static /app/frontend/.next/static
COPY --from=frontend-builder /app/public /app/frontend/public

# Copy nginx config
COPY deploy/nginx.conf /etc/nginx/http.d/default.conf

# Copy supervisor config
COPY deploy/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Expose port
EXPOSE 8080

# Start supervisor (manages nginx, node, go)
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]

