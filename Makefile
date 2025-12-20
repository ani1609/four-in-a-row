.PHONY: help install build run dev test clean docker-up docker-down docker-logs

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Docker commands
docker-up: ## Start Docker services (PostgreSQL & Kafka)
	@echo "Starting Docker services..."
	docker compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "✓ Docker services are ready!"

docker-down: ## Stop Docker services
	docker compose down

docker-logs: ## View Docker service logs
	docker compose logs -f

docker-clean: ## Stop and remove all Docker volumes
	docker compose down -v

# Backend commands
backend-deps: ## Install backend dependencies
	cd backend && go mod download

backend-build: ## Build the backend server
	@echo "Building backend..."
	cd backend && go build -o server .
	cd backend && go build -o consumer ./cmd/consumer
	@echo "✓ Backend built successfully"

backend-run: ## Run the backend server
	cd backend && ./server

backend-consumer: ## Run the Kafka consumer
	cd backend && ./consumer

backend-test: ## Run backend tests
	cd backend && go test -v ./...

# Frontend commands
frontend-deps: ## Install frontend dependencies
	cd frontend && pnpm install

frontend-dev: ## Run frontend development server
	cd frontend && pnpm run dev

frontend-build: ## Build frontend for production
	cd frontend && pnpm run build

# Combined commands
install: backend-deps frontend-deps ## Install all dependencies

build: backend-build frontend-build ## Build both backend and frontend

dev: docker-up backend-build ## Start Docker and build backend
	@echo ""
	@echo "✓ Development environment ready!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Run backend:  make backend-run"
	@echo "  2. Run frontend: make frontend-dev"
	@echo "  3. Open: http://localhost:5173"

test: backend-test ## Run all tests

clean: ## Clean build artifacts
	rm -f backend/server backend/consumer
	rm -rf frontend/dist

# Full development workflow
start: dev backend-run ## Start everything (Docker + Backend)

# Kafka commands
kafka-create-topic: ## Create Kafka topic for analytics
	docker compose exec kafka kafka-topics --create --topic game-analytics --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1 || echo "Topic may already exist"

kafka-list-topics: ## List all Kafka topics
	docker compose exec kafka kafka-topics --list --bootstrap-server localhost:9092

kafka-view-messages: ## View all messages in game-analytics topic
	docker compose exec kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic game-analytics --from-beginning

.DEFAULT_GOAL := help
