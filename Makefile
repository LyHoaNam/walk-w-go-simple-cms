.PHONY: help run build test clean docker-up docker-down docker-logs migrate

# Variables
APP_NAME=simple-golang-api
DOCKER_COMPOSE=docker-compose

help: ## Hiển thị trợ giúp
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Chạy ứng dụng local
	@echo "Starting application..."
	go run cmd/api/main.go

build: ## Build ứng dụng
	@echo "Building application..."
	go build -o bin/$(APP_NAME) cmd/api/main.go

test: ## Chạy tests
	@echo "Running tests..."
	go test -v ./...

clean: ## Xóa build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	go clean

docker-up: ## Khởi động Docker containers
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d

docker-down: ## Dừng Docker containers
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

docker-logs: ## Xem logs của Docker containers
	$(DOCKER_COMPOSE) logs -f

docker-rebuild: ## Rebuild và restart Docker containers
	@echo "Rebuilding Docker containers..."
	$(DOCKER_COMPOSE) up -d --build

migrate:
	@echo "Running migrations..."
	@if [ -f .env ]; then \
		export $$(cat .env | xargs) && \
		mysql -h$$DB_HOST -P$$DB_PORT -u$$DB_USER -p$$DB_PASSWORD < migrations/001_create_users_table.sql; \
	else \
		echo "Error: .env file not found. Please copy .env.example to .env first."; \
	fi

init: ## Khởi tạo project (copy .env, tải dependencies)
	@echo "Initializing project..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from .env.example"; \
	fi
	go mod download
	@echo "Project initialized successfully!"

tidy: ## Tidy go modules
	go mod tidy
