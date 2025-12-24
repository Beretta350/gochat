.PHONY: help
.PHONY: docker-up docker-down docker-logs docker-build docker-restart
.PHONY: docker-infra docker-infra-down
.PHONY: docker-api-up docker-api-build docker-api-logs docker-api-restart
.PHONY: dev-api api-build api-test api-lint api-fmt
.PHONY: dev-web web-build web-lint web-test

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
	@echo "  make docker-up           Start ALL services (api + db + redis)"
	@echo "  make docker-down         Stop all services"
	@echo "  make docker-logs         View all logs"
	@echo "  make docker-build        Build all Docker images"
	@echo "  make docker-restart      Rebuild and restart all"
	@echo ""
	@echo "$(YELLOW)Docker - Infrastructure:$(NC)"
	@echo "  make docker-infra        Start PostgreSQL + Redis (for local dev)"
	@echo "  make docker-infra-down   Stop infrastructure"
	@echo ""
	@echo "$(YELLOW)Docker - Backend Only:$(NC)"
	@echo "  make docker-api-up       Start API + infra (no frontend)"
	@echo "  make docker-api-build    Build API Docker image"
	@echo "  make docker-api-logs     View API logs only"
	@echo "  make docker-api-restart  Rebuild and restart API"
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
	@echo "  make web-build           Build frontend"
	@echo "  make web-lint            Lint frontend code"
	@echo "  make web-test            Run frontend tests"

# ========================
# Docker - Full Stack
# ========================

## docker-up: Start all services with Docker Compose
docker-up:
	@echo "$(GREEN)Starting all services...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)Services started:$(NC)"
	@echo "  - API: http://localhost:8080"
	@echo "  - Redis UI: http://localhost:8081"

## docker-down: Stop all services
docker-down:
	@echo "$(GREEN)Stopping all services...$(NC)"
	@docker-compose down

## docker-logs: View logs from all services
docker-logs:
	@docker-compose logs -f

## docker-build: Build all Docker images
docker-build:
	@echo "$(GREEN)Building all services...$(NC)"
	@docker-compose build

## docker-restart: Rebuild and restart all services
docker-restart:
	@echo "$(GREEN)Restarting all services...$(NC)"
	@docker-compose down
	@docker-compose up -d --build
	@echo "$(GREEN)All services restarted$(NC)"

# ========================
# Docker - Infrastructure
# ========================

## docker-infra: Start only PostgreSQL and Redis
docker-infra:
	@echo "$(GREEN)Starting infrastructure...$(NC)"
	@docker-compose up -d postgres redis redis-commander
	@echo "$(GREEN)Infrastructure started:$(NC)"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"
	@echo "  - Redis UI: http://localhost:8081"

## docker-infra-down: Stop infrastructure
docker-infra-down:
	@echo "$(GREEN)Stopping infrastructure...$(NC)"
	@docker-compose down

# ========================
# Docker - Backend Only
# ========================

## docker-api-up: Start API with infrastructure (Docker)
docker-api-up:
	@echo "$(GREEN)Starting API + infrastructure...$(NC)"
	@docker-compose up -d postgres redis redis-commander api
	@echo "$(GREEN)Services started:$(NC)"
	@echo "  - API: http://localhost:8080"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"
	@echo "  - Redis UI: http://localhost:8081"

## docker-api-build: Build API Docker image
docker-api-build:
	@echo "$(GREEN)Building API image...$(NC)"
	@docker-compose build api
	@echo "$(GREEN)API image built$(NC)"

## docker-api-logs: View API logs only
docker-api-logs:
	@docker-compose logs -f api

## docker-api-restart: Rebuild and restart API only
docker-api-restart:
	@echo "$(GREEN)Restarting API...$(NC)"
	@docker-compose up -d --build api
	@echo "$(GREEN)API restarted$(NC)"

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

