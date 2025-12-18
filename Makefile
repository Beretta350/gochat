.PHONY: all build run dev test lint fmt clean docker-up docker-down docker-logs help

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

## build: Build the application
build:
	@echo "$(GREEN)Building...$(NC)"
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)Built $(BINARY_PATH)$(NC)"

## run: Run the application
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
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Cleaned$(NC)"

## docker-up: Start Docker containers
docker-up:
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)Containers started$(NC)"

## docker-down: Stop Docker containers
docker-down:
	@echo "$(GREEN)Stopping Docker containers...$(NC)"
	@docker-compose down
	@echo "$(GREEN)Containers stopped$(NC)"

## docker-logs: Show Docker logs
docker-logs:
	@docker-compose logs -f

## docker-restart: Restart Docker containers
docker-restart: docker-down docker-up

## setup: Install development dependencies
setup:
	@echo "$(GREEN)Installing development dependencies...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/air-verse/air@latest
	@echo "$(GREEN)Dependencies installed$(NC)"

## all: Run fmt, lint, test, and build
all: fmt lint test build
