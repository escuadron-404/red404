.PHONY: all dev up down build-backend build-frontend test lint format lint-fix clean help

all: dev

dev: up 
	@echo "Starting frontend dev server..."
	cd frontend && npm run dev & 
	@echo "Starting backend hot-reload with air..."
	cd backend && air 

up:
	docker-compose up -d --build db
	docker-compose up -d --build backend
	@echo "Docker services started. DB: 5432, Backend: 8080"
	@echo "Frontend will run on Vite's default port (e.g., 5173)."

down:
	docker-compose down --remove-orphans
	@echo "All Docker services stopped."

build-backend:
	@echo "Building Go backend production image..."
	docker-compose build backend

build-frontend:
	@echo "Building React frontend production assets..."
	cd frontend && npm run build

test:
	@echo "Running backend tests..."
	cd backend && go test -v ./... 
	@echo "Running frontend tests (if you add them)..."
	# cd frontend && npm test

lint: format 
	@echo "Running comprehensive Go backend linting with golangci-lint..."
	cd backend && golangci-lint run --config ../.golangci.yml --timeout 5m 
	@echo "Linting React frontend..."
	cd frontend && npm run lint

format:
	@echo "Formatting Go backend code with gofmt..."
	cd backend && go fmt ./...
	# maybe? goimports -w ./...

lint-fix:
	@echo "Running Go backend linting with automatic fixes (if possible)..."
	cd backend && golangci-lint run --fix --config ../.golangci.yml --timeout 5m
	@echo "Linting React frontend..."
	cd frontend && npm run lint:write 

clean:
	@echo "Cleaning temporary files and Docker volumes..."
	rm -rf backend/tmp backend/main
	rm -rf frontend/dist
	docker volume rm my-social-app_db_data || true 

help:
	@echo "Usage: make [command]"
	@echo ""
	@echo "Commands:"
	@echo "  dev             - Starts all services and frontend/backend dev servers with hot-reload"
	@echo "  up              - Starts Docker services (DB, Backend) in detached mode"
	@echo "  down            - Stops and removes Docker services"
	@echo "  build-backend   - Builds the production Docker image for the backend"
	@echo "  build-frontend  - Builds the production JavaScript assets for the frontend"
	@echo "  test            - Runs all tests (backend and frontend)"
	@echo "  lint            - Runs linters (including gofmt and golangci-lint) and formatters"
	@echo "  format          - Formats Go backend code using gofmt"
	@echo "  lint-fix        - Runs Go backend linting with automatic fixes (if available)"
	@echo "  clean           - Removes temporary files and Docker volumes"
	@echo "  help            - Displays this help message"
