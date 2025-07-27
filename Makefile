.PHONY: build run test clean docker-up docker-down help

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run specific test
test-specific:
	@echo "Usage: make test-specific TEST=TestName"
	@echo "Example: make test-specific TEST=TestCreateDevice"
	go test -v -run $(TEST) ./...

# Run integration tests only
test-integration:
	go test -v ./tests/...

# Run unit tests only (excluding integration tests)
test-unit:
	go test -v ./internal/... ./pkg/...

# Run tests with race detection
test-race:
	go test -race ./...

# Run benchmarks
test-bench:
	go test -bench=. ./...

# Clean build artifacts and test files
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Start Docker services
docker-up:
	docker-compose up -d

# Stop Docker services
docker-down:
	docker-compose down

# Show logs
logs:
	docker-compose logs -f

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate API documentation
docs:
	swag init -g cmd/server/main.go

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run with hot reload (requires air)
dev:
	air

# Create database migration
migrate:
	# Add migration commands here when using a migration tool

# Setup test database
test-db-setup:
	@echo "Setting up test database..."
	@echo "Make sure PostgreSQL is running and accessible"
	@echo "Create test database: CREATE DATABASE iot_platform_test;"

# Run all checks (format, lint, test)
check: fmt lint test

# Help
help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-verbose    - Run tests with verbose output"
	@echo "  test-specific   - Run specific test (TEST=TestName)"
	@echo "  test-integration- Run integration tests only"
	@echo "  test-unit       - Run unit tests only"
	@echo "  test-race       - Run tests with race detection"
	@echo "  test-bench      - Run benchmarks"
	@echo "  clean           - Clean build artifacts and test files"
	@echo "  docker-up       - Start Docker services"
	@echo "  docker-down     - Stop Docker services"
	@echo "  logs            - Show Docker logs"
	@echo "  fmt             - Format code"
	@echo "  lint            - Lint code"
	@echo "  deps            - Install dependencies"
	@echo "  dev             - Run with hot reload (requires air)"
	@echo "  test-db-setup   - Setup test database"
	@echo "  check           - Run all checks (format, lint, test)" 