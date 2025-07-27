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

# Clean build artifacts
clean:
	rm -rf bin/

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

# Help
help:
	@echo "Available commands:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  docker-up  - Start Docker services"
	@echo "  docker-down- Stop Docker services"
	@echo "  logs       - Show Docker logs"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code"
	@echo "  deps       - Install dependencies"
	@echo "  dev        - Run with hot reload (requires air)" 