.PHONY: help deps fmt build run test clean docker-up docker-down docker-logs

help: ## Show available commands
	@echo 'Backend Makefile Commands:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Dependencies
deps: ## Install Go dependencies
	go mod download
	go mod tidy

# Formatting
fmt: ## Format Go code
	@echo "Formatting code..."
	gofmt -w .
	@echo "✓ Code formatted"

# Build
build: fmt ## Build all binaries
	@echo "Building binaries..."
	go build -o bin/server .
	go build -o bin/kafka-consumer ./cmd/kafka-consumer
	go build -o bin/redis-consumer ./cmd/redis-consumer
	@echo "✓ All binaries built in bin/"

# Run
run: ## Run main server
	./bin/server

run-kafka-consumer: ## Run Kafka analytics consumer
	./bin/kafka-consumer

run-redis-consumer: ## Run Redis analytics consumer
	./bin/redis-consumer

# Docker
docker-up: ## Start all Docker services (PostgreSQL, Zookeeper, Kafka, Redis)
	@echo "Starting Docker services..."
	docker compose up -d postgres zookeeper kafka redis
	@echo "Waiting for services..."
	@sleep 10
	@echo "✓ Services ready:"
	@docker ps --filter "name=4-in-a-row" --format "  - {{.Names}}: {{.Status}}"

docker-down: ## Stop Docker services
	docker compose down

docker-logs: ## View Docker logs
	docker compose logs -f

docker-clean: ## Stop and remove volumes
	docker compose down -v

# Clean
clean: ## Clean build artifacts
	rm -rf bin/
	@echo "✓ Cleaned build artifacts"
