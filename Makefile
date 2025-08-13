.PHONY: all dev up-db down build-backend-prod build-frontend-prod test lint format lint-fix clean help setup-dev

# Default target for quick start
all: dev

# --- Development Workflow ---
dev: setup-dev up-db
	@echo "--- Starting development servers ---"
	# Ensure frontend dependencies are installed before running dev server
	@echo "Starting frontend dev server (http://localhost:5173)..."
	cd frontend && npm install && npm run dev & # Run in background, & is crucial
	
	@echo "Starting backend hot-reload with air (http://localhost:8080)..."
	# Ensure Go dependencies are installed for air
	cd backend && go mod tidy && air # Air runs in foreground for logs. Use Ctrl+C to stop.
	
	@echo "Go backend stopped. Remember to run 'make down' to stop the database."
	@echo "You may also need to manually kill the 'npm run dev' process if it's still running."

# Initializes local dev environment: installs Go/JS deps
setup-dev:
	@echo "--- Setting up local development environment ---"
	@if [ ! -f ./.env ]; then \
		echo "Error: ./.env file not found. Please create one based on ./.env.example"; \
		exit 1; \
	fi
	@echo "Copying ./.env to backend/ and frontend/..."
	cp ./.env backend/.env
	cp ./.env frontend/.env
	@echo "Installing Go backend dependencies..."
	cd backend && go mod tidy
	@echo "Installing Air (Go hot-reloader)..."
	cd backend && go install github.com/air-verse/air@latest
	@echo "Installing frontend Node.js dependencies..."
	cd frontend && npm install

# Starts only the database container
up-db:
	@echo "--- Starting PostgreSQL Docker service ---"
	sudo docker-compose up -d --build --wait db # --wait will wait for healthcheck
	@echo "PostgreSQL Docker service started and healthy on port 5432."

# Stops all Docker services
down:
	@echo "--- Stopping all Docker services ---"
	sudo docker-compose down --remove-orphans # Removes stopped containers and networks
	@echo "All Docker services stopped."
	@echo "If you started frontend/backend dev servers via 'make dev', you will need to manually terminate them."

# --- Testing, Linting, Formatting ---

test:
	@echo "--- Running backend tests ---"
	# cd backend && go test -v ./... -failfast # Add -failfast for quicker feedback
	@echo "--- Running frontend tests (if you add them) ---"
	# cd frontend && npm test

lint: format # TODO: why can tsc -b find errors that biome cannot?
	@echo "--- Running comprehensive Go backend linting with golangci-lint ---"
	cd backend && golangci-lint run --config ../.golangci.yml --timeout 5m "./..."
	@echo "--- Linting React frontend ---"
	cd frontend && npm run lint

format:
	@echo "--- Formatting Go backend code ---"
	cd backend && go fmt ./...
	@echo "--- Formatting frontend code ---"
	cd frontend && npm run lint:write
  
lint-fix:
	@echo "--- Running Go backend linting with automatic fixes (if possible) ---"
	cd backend && golangci-lint run --fix --config ../.golangci.yml --timeout 5m "./..."
	@echo "--- Linting React frontend with fixes ---"
	cd frontend && npm run lint:write # Assuming you have a script for this

# --- Clean Up ---
clean:
	@echo "--- Cleaning temporary files and built assets ---"
	rm -rf backend/tmp 
	rm -rf frontend/dist 
	# rm -rf frontend/node_modules # maybe?
	@echo "Local temporary files and built assets removed."
	@echo "To remove Docker volumes, use 'docker-compose down -v'."

# --- Help ---
help:
	@echo "Usage: make [command]"
	@echo ""
	@echo "Commands for Development:"
	@echo "  dev             - Full dev setup: starts DB (Docker), frontend (host:5173), backend (host:8080)."
	@echo "                    (Frontend & Backend processes must be stopped manually with Ctrl+C)"
	@echo "  setup-dev       - Installs Go/JS dependencies and Air for local dev."
	@echo "  up-db           - Starts only the PostgreSQL Docker service in detached mode."
	@echo "  down            - Stops and removes all Docker services (DB)."
	@echo ""
	@echo "General Commands:"
	@echo "  test            - Runs all tests (backend and frontend)."
	@echo "  lint            - Runs linters and formatters (backend and frontend)."
	@echo "  format          - Formats Go backend code using gofmt and goimports."
	@echo "  lint-fix        - Runs Go backend linting with automatic fixes."
	@echo "  clean           - Removes temporary files and built assets."
	@echo "  help            - Displays this help message."
