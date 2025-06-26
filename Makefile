# Rural Health Management System Makefile

.PHONY: help build run dev test clean seed migrate

# Default target
help:
	@echo "Available commands:"
	@echo "  build     - Build the application"
	@echo "  run       - Run the application"
	@echo "  dev       - Run in development mode with auto-reload"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build artifacts"
	@echo "  seed      - Seed the database with sample data"
	@echo "  deps      - Install dependencies"

# Build the application
build:
	go build -o bin/rural-health-api main.go

# Run the application
run:
	go run main.go

# Run in development mode (requires air for auto-reload)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Installing air for hot reload..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf tmp/
	go clean

# Seed the database with sample data
seed:
	go run cmd/seed/main.go

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run the application with specific port
run-port:
	PORT=8080 go run main.go

# Check for security vulnerabilities
security:
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "Installing gosec..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Lint code (requires golangci-lint)
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it from https://golangci-lint.run/"; \
	fi

# All checks
check: fmt vet lint test

# Docker commands
docker-build:
	docker-compose build --no-cache

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v --remove-orphans
	docker system prune -f

docker-restart:
	docker-compose down
	docker-compose up -d

docker-seed:
	docker-compose exec api go run cmd/seed/main.go

docker-shell:
	docker-compose exec api sh

# Complete Docker setup
docker-setup:
	@echo "🏥 Setting up Rural Health Management System with Docker..."
	docker-compose build --no-cache
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 10
	docker-compose exec api go run cmd/seed/main.go
	@echo "🎉 Setup complete! API available at http://localhost:3000"
