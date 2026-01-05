.PHONY: help
.PHONY: docker-up docker-down docker-logs docker-build docker-restart
.PHONY: docker-infra docker-infra-down
.PHONY: docker-api-up docker-api-build docker-api-logs docker-api-restart
.PHONY: docker-web-up docker-web-build docker-web-logs docker-web-restart
.PHONY: dev-api api-build api-test api-lint api-fmt
.PHONY: dev-web web-install web-build web-lint web-test

# Colors
GREEN=\033[0;32m
YELLOW=\033[0;33m
CYAN=\033[0;36m
NC=\033[0m

## help: Show this help message
help:
	@echo "$(CYAN)GoChat - Full Stack Commands$(NC)"
	@echo ""
	@echo "$(YELLOW)Docker - Full Stack:$(NC)"
	@echo "  make docker-up           Start ALL services (web + api + infra)"
	@echo "  make docker-down         Stop all services"
	@echo "  make docker-logs         View all logs"
	@echo "  make docker-build        Build all Docker images"
	@echo "  make docker-restart      Rebuild and restart all"
	@echo ""
	@echo "$(YELLOW)Docker - Infrastructure:$(NC)"
	@echo "  make docker-infra        Start PostgreSQL + Redis only"
	@echo "  make docker-dev          Start infra + nginx (for local dev)"
	@echo "  make docker-dev-down     Stop dev environment"
	@echo "  make docker-infra-down   Stop infrastructure"
	@echo ""
	@echo "$(YELLOW)Docker - Backend Only:$(NC)"
	@echo "  make docker-api-up       Start API + infra (no frontend)"
	@echo "  make docker-api-build    Build API Docker image"
	@echo "  make docker-api-logs     View API logs only"
	@echo "  make docker-api-restart  Rebuild and restart API"
	@echo ""
	@echo "$(YELLOW)Docker - Frontend Only:$(NC)"
	@echo "  make docker-web-up       Start Web + API + infra"
	@echo "  make docker-web-build    Build Web Docker image"
	@echo "  make docker-web-logs     View Web logs only"
	@echo "  make docker-web-restart  Rebuild and restart Web"
	@echo ""
	@echo "$(YELLOW)Development (Local):$(NC)"
	@echo "  make dev-api             Run API locally with hot reload"
	@echo "  make dev-web             Run frontend dev server"
	@echo ""
	@echo "$(YELLOW)Backend (Go):$(NC)"
	@echo "  make api-build           Build backend binary"
	@echo "  make api-test            Run backend tests"
	@echo "  make api-lint            Lint backend code"
	@echo "  make api-fmt             Format backend code"
	@echo ""
	@echo "$(YELLOW)Frontend (Next.js):$(NC)"
	@echo "  make web-install         Install frontend dependencies"
	@echo "  make web-build           Build frontend"
	@echo "  make web-lint            Lint frontend code"
	@echo "  make web-test            Run frontend tests"

# ========================
# Docker - Full Stack
# ========================

## docker-up: Start all services with Docker Compose (builds if needed)
docker-up:
	@echo "$(GREEN)Starting all services...$(NC)"
	@docker compose up -d
	@echo "$(GREEN)Services started:$(NC)"
	@echo "  - App: http://localhost"
	@echo "  - Redis UI: http://localhost:8081"

## docker-down: Stop all services
docker-down:
	@echo "$(GREEN)Stopping all services...$(NC)"
	@docker compose down

## docker-logs: View logs from all services
docker-logs:
	@docker compose logs -f

## docker-build: Build all Docker images
docker-build:
	@echo "$(GREEN)Building all services...$(NC)"
	@docker compose build

## docker-restart: Rebuild and restart all services
docker-restart:
	@echo "$(GREEN)Restarting all services...$(NC)"
	@docker compose down
	@docker compose up -d --build
	@echo "$(GREEN)All services restarted$(NC)"

# ========================
# Docker - Infrastructure
# ========================

## docker-infra: Start only PostgreSQL and Redis
docker-infra:
	@echo "$(GREEN)Starting infrastructure...$(NC)"
	@docker compose up -d postgres redis redis-commander
	@echo "$(GREEN)Infrastructure started:$(NC)"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"
	@echo "  - Redis UI: http://localhost:8081"

## docker-dev: Start infra + nginx for local development
docker-dev:
	@echo "$(GREEN)Starting dev environment (infra + nginx)...$(NC)"
	@docker compose up -d postgres redis redis-commander
	@docker run -d --rm --name gochat-nginx-dev \
		-p 80:80 \
		-v $(PWD)/nginx-dev.conf:/etc/nginx/conf.d/default.conf:ro \
		--add-host=host.docker.internal:host-gateway \
		nginx:alpine
	@echo "$(GREEN)Dev environment started:$(NC)"
	@echo "  - App: http://localhost (via nginx)"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"
	@echo "  - Redis UI: http://localhost:8081"
	@echo ""
	@echo "$(YELLOW)Now run in separate terminals:$(NC)"
	@echo "  make dev-api"
	@echo "  make dev-web"

## docker-dev-down: Stop dev environment
docker-dev-down:
	@echo "$(GREEN)Stopping dev environment...$(NC)"
	@docker stop gochat-nginx-dev 2>/dev/null || true
	@docker compose down

## docker-infra-down: Stop infrastructure
docker-infra-down:
	@echo "$(GREEN)Stopping infrastructure...$(NC)"
	@docker compose down

# ========================
# Docker - Backend Only
# ========================

## docker-api-up: Start API with infrastructure (Docker)
docker-api-up:
	@echo "$(GREEN)Starting API + infrastructure...$(NC)"
	@docker compose up -d postgres redis redis-commander api
	@echo "$(GREEN)Services started:$(NC)"
	@echo "  - API: internal (use nginx for access)"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"
	@echo "  - Redis UI: http://localhost:8081"

## docker-api-build: Build API Docker image
docker-api-build:
	@echo "$(GREEN)Building API image...$(NC)"
	@docker compose build api
	@echo "$(GREEN)API image built$(NC)"

## docker-api-logs: View API logs only
docker-api-logs:
	@docker compose logs -f api

## docker-api-restart: Rebuild and restart API only
docker-api-restart:
	@echo "$(GREEN)Restarting API...$(NC)"
	@docker compose up -d --build api
	@echo "$(GREEN)API restarted$(NC)"

# ========================
# Docker - Frontend Only
# ========================

## docker-web-up: Start Web + API + infrastructure (Docker)
docker-web-up:
	@echo "$(GREEN)Starting Web + API + infrastructure...$(NC)"
	@docker compose up -d postgres redis redis-commander api web nginx
	@echo "$(GREEN)Services started:$(NC)"
	@echo "  - App: http://localhost"
	@echo "  - Redis UI: http://localhost:8081"

## docker-web-build: Build Web Docker image
docker-web-build:
	@echo "$(GREEN)Building Web image...$(NC)"
	@docker compose build web
	@echo "$(GREEN)Web image built$(NC)"

## docker-web-logs: View Web logs only
docker-web-logs:
	@docker compose logs -f web

## docker-web-restart: Rebuild and restart Web only
docker-web-restart:
	@echo "$(GREEN)Restarting Web...$(NC)"
	@docker compose up -d --build web
	@echo "$(GREEN)Web restarted$(NC)"

# ========================
# Development (Local)
# ========================

## dev-api: Run backend locally with hot reload (requires: make docker-infra)
dev-api:
	@echo "$(GREEN)Starting backend (hot reload)...$(NC)"
	@cd backend && air

## dev-web: Run frontend dev server
dev-web:
	@echo "$(GREEN)Starting frontend dev server...$(NC)"
	@cd frontend && npm run dev

# ========================
# Backend Build & Test
# ========================

## api-build: Build backend binary
api-build:
	@echo "$(GREEN)Building backend...$(NC)"
	@cd backend && go build -o bin/gochat cmd/main.go

## api-test: Run backend tests
api-test:
	@echo "$(GREEN)Running backend tests...$(NC)"
	@cd backend && go test -v ./...

## api-lint: Lint backend code
api-lint:
	@echo "$(GREEN)Linting backend...$(NC)"
	@cd backend && golangci-lint run ./...

## api-fmt: Format backend code
api-fmt:
	@echo "$(GREEN)Formatting backend...$(NC)"
	@cd backend && go fmt ./...

# ========================
# Frontend Commands
# ========================

## web-install: Install frontend dependencies
web-install:
	@echo "$(GREEN)Installing frontend dependencies...$(NC)"
	@cd frontend && npm install

## web-build: Build frontend for production
web-build:
	@echo "$(GREEN)Building frontend...$(NC)"
	@cd frontend && npm run build

## web-lint: Lint frontend code
web-lint:
	@echo "$(GREEN)Linting frontend...$(NC)"
	@cd frontend && npm run lint

## web-test: Run frontend tests
web-test:
	@echo "$(GREEN)Running frontend tests...$(NC)"
	@cd frontend && npm test

