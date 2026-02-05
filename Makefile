.PHONY: help dev dev-backend dev-frontend build test lint clean docker-up docker-down

# Default target
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ============================================================
# Development
# ============================================================

dev: ## Run both backend and frontend in development mode
	@echo "Starting backend and frontend..."
	@$(MAKE) dev-backend &
	@$(MAKE) dev-frontend

dev-backend: ## Run backend in development mode
	cd backend && APP_ENV=development go run server.go

dev-frontend: ## Run frontend in development mode
	cd frontend && npm run dev

install: ## Install all dependencies
	cd backend && go mod download
	cd frontend && npm ci --legacy-peer-deps

# ============================================================
# Build
# ============================================================

build: build-backend build-frontend ## Build both backend and frontend

build-backend: ## Build backend binary
	cd backend && CGO_ENABLED=0 go build -ldflags="-s -w" -o server .

build-frontend: ## Build frontend for production
	cd frontend && VITE_API_URL=/graphql npm run build

# ============================================================
# Test
# ============================================================

test: test-backend test-frontend ## Run all tests

test-backend: ## Run backend tests with coverage
	cd backend && go test -v -race -coverprofile=coverage.out ./...
	cd backend && go tool cover -func=coverage.out

test-frontend: ## Run frontend linter and type check
	cd frontend && npm run lint
	cd frontend && npx tsc --noEmit

# ============================================================
# Code Quality
# ============================================================

lint: lint-backend lint-frontend ## Run all linters

lint-backend: ## Run Go linters
	cd backend && go vet ./...

lint-frontend: ## Run ESLint
	cd frontend && npm run lint

# ============================================================
# Docker
# ============================================================

docker-up: ## Start all services with Docker Compose
	docker compose up --build -d

docker-down: ## Stop all services
	docker compose down

docker-logs: ## View container logs
	docker compose logs -f

docker-ps: ## Show running containers
	docker compose ps

# ============================================================
# Utilities
# ============================================================

clean: ## Clean build artifacts
	rm -f backend/server backend/coverage.out
	rm -rf frontend/dist frontend/node_modules

health: ## Check server health
	@curl -s http://localhost:4000/health | python3 -m json.tool 2>/dev/null || echo "Server not running"

metrics: ## View server metrics
	@curl -s http://localhost:4000/metrics | python3 -m json.tool 2>/dev/null || echo "Server not running"

token-admin: ## Generate an admin JWT token
	@curl -s -X POST http://localhost:4000/auth/token \
		-H "Content-Type: application/json" \
		-d '{"user_id":"admin-1","email":"admin@example.com","name":"Admin User","role":"admin"}' \
		| python3 -m json.tool 2>/dev/null || echo "Server not running"

token-user: ## Generate a user JWT token
	@curl -s -X POST http://localhost:4000/auth/token \
		-H "Content-Type: application/json" \
		-d '{"user_id":"user-1","email":"user@example.com","name":"Test User","role":"user"}' \
		| python3 -m json.tool 2>/dev/null || echo "Server not running"
