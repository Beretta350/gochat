.PHONY: all build run dev test lint fmt clean help
.PHONY: docker-up docker-down docker-logs docker-build docker-run docker-stop

# Variables
APP_NAME=gochat
MAIN_PATH=cmd/main.go
BINARY_PATH=bin/$(APP_NAME)

# Colors
GREEN=\033[0;32m
NC=\033[0m # No Color

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' Makefile | sed 's/## /  /'

# ========================
# Local Development
# ========================

## build: Build the application
build:
	@echo "$(GREEN)Building...$(NC)"
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)Built $(BINARY_PATH)$(NC)"

## run: Run the application locally
run:
	@go run $(MAIN_PATH)

## dev: Run with hot reload (requires air)
dev:
	@air

## test: Run tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	@go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

## lint: Run linter
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@golangci-lint run ./...

## lint-fix: Run linter and fix issues
lint-fix:
	@echo "$(GREEN)Running linter with auto-fix...$(NC)"
	@golangci-lint run --fix ./...

## fmt: Format code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	@go fmt ./...
	@goimports -w .

## tidy: Tidy go modules
tidy:
	@echo "$(GREEN)Tidying modules...$(NC)"
	@go mod tidy

## clean: Clean build artifacts
clean:
	@echo "$(GREEN)Cleaning...$(NC)"
	@rm -rf bin/ tmp/
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Cleaned$(NC)"

## setup: Install development dependencies
setup:
	@echo "$(GREEN)Installing development dependencies...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/air-verse/air@latest
	@echo "$(GREEN)Dependencies installed$(NC)"

## all: Run fmt, lint, test, and build
all: fmt lint test build

# ========================
# Docker - Redis only
# ========================

## redis-up: Start only Redis container
redis-up:
	@echo "$(GREEN)Starting Redis...$(NC)"
	@docker-compose up -d redis redis-commander
	@echo "$(GREEN)Redis started$(NC)"

## redis-down: Stop Redis container
redis-down:
	@echo "$(GREEN)Stopping Redis...$(NC)"
	@docker-compose down
	@echo "$(GREEN)Redis stopped$(NC)"

# ========================
# Docker - Full stack
# ========================

## docker-build: Build Docker image
docker-build:
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker-compose build
	@echo "$(GREEN)Image built$(NC)"

## docker-up: Start all containers (app + redis)
docker-up:
	@echo "$(GREEN)Starting all containers...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)Containers started$(NC)"
	@echo "$(GREEN)App: http://localhost:8080$(NC)"
	@echo "$(GREEN)Redis UI: http://localhost:8081$(NC)"

## docker-down: Stop all containers
docker-down:
	@echo "$(GREEN)Stopping containers...$(NC)"
	@docker-compose down
	@echo "$(GREEN)Containers stopped$(NC)"

## docker-logs: Show Docker logs
docker-logs:
	@docker-compose logs -f

## docker-logs-app: Show only app logs
docker-logs-app:
	@docker-compose logs -f app

## docker-restart: Rebuild and restart all containers
docker-restart:
	@echo "$(GREEN)Restarting containers...$(NC)"
	@docker-compose down
	@docker-compose up -d --build
	@echo "$(GREEN)Containers restarted$(NC)"
