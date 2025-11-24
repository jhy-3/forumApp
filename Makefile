# Detect compose command (Podman preferred, then Docker)
COMPOSE_CMD := $(shell \
	if command -v podman-compose > /dev/null 2>&1; then \
		echo "podman-compose"; \
	elif command -v docker-compose > /dev/null 2>&1; then \
		echo "docker-compose"; \
	else \
		echo "docker compose"; \
	fi)

.PHONY: help start stop restart logs clean dev-db build

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

start: ## Start all services (database, backend, frontend)
	$(COMPOSE_CMD) up -d
	@echo "Services started!"
	@echo "Frontend: http://localhost:3000"
	@echo "Backend:  http://localhost:8080"
	@echo "Database: localhost:3306"

stop: ## Stop all services
	$(COMPOSE_CMD) down

restart: stop start ## Restart all services

logs: ## Show logs from all services
	$(COMPOSE_CMD) logs -f

clean: ## Stop services and remove volumes
	$(COMPOSE_CMD) down -v

dev-db: ## Start only database for local development
	$(COMPOSE_CMD) -f docker-compose.dev.yml up -d
	@echo "Database started on localhost:3306"

build: ## Rebuild all containers
	$(COMPOSE_CMD) build

backend-dev: ## Run backend locally (requires dev-db)
	cd backend && go run main.go

frontend-dev: ## Run frontend locally
	cd frontend && npm start

test-api: ## Test backend API health
	@curl -s http://localhost:8080/health | jq '.' || echo "Backend not responding"

test-register: ## Test user registration
	@curl -s -X POST http://localhost:8080/api/users/register \
		-H "Content-Type: application/json" \
		-d '{"username":"testuser","email":"test@example.com","password":"password123"}' | jq '.'

db-shell: ## Open MariaDB shell
	$(COMPOSE_CMD) exec database mysql -uforumuser -pforumpass123 forumdb

check-engine: ## Check which container engine is being used
	@echo "Container engine: $(COMPOSE_CMD)"

