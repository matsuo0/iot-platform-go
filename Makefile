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

# CI/CD commands
ci-test:
	go test -v -race -coverprofile=coverage.out ./internal/... ./pkg/...
	go test -v -race ./tests/...
	go tool cover -html=coverage.out -o coverage.html

ci-build:
	go build -v -o bin/server cmd/server/main.go

ci-lint:
	golangci-lint run --timeout=5m

# Docker commands
docker-build:
	docker build -t iot-platform-go .

docker-run:
	docker run -p 8080:8080 iot-platform-go

# Dependencies management
deps-update:
	go get -u ./...
	go mod tidy

deps-check:
	go list -u -m all

deps-audit:
	govulncheck ./...

deps-clean:
	go clean -modcache

# Security scan
security-scan:
	trivy fs --format sarif --output trivy-results.sarif .
	govulncheck ./...

# Local CI simulation
ci-local: fmt lint test-coverage build

# GitHub Actions helpers
actions-test:
	@echo "Testing GitHub Actions locally..."
	@echo "Use act to run GitHub Actions locally:"
	@echo "  act -j test"
	@echo "  act -j build"
	@echo "  act -j security"

actions-validate:
	@echo "Validating GitHub Actions workflows..."
	@for file in .github/workflows/*.yml; do \
		echo "Validating $$file"; \
		yamllint "$$file" || echo "Warning: yamllint not installed"; \
	done

# Performance testing
bench:
	go test -bench=. -benchmem ./internal/...

bench-cpu:
	go test -bench=. -cpuprofile=cpu.prof ./internal/...

bench-memory:
	go test -bench=. -memprofile=mem.prof ./internal/...

# Documentation
docs-serve:
	@echo "Starting documentation server..."
	@echo "Visit http://localhost:6060/pkg/iot-platform-go/"
	godoc -http=:6060

# Release helpers
release-patch:
	@echo "Creating patch release..."
	@git tag -a v$$(semver bump patch) -m "Release v$$(semver bump patch)"
	@git push origin --tags

release-minor:
	@echo "Creating minor release..."
	@git tag -a v$$(semver bump minor) -m "Release v$$(semver bump minor)"
	@git push origin --tags

release-major:
	@echo "Creating major release..."
	@git tag -a v$$(semver bump major) -m "Release v$$(semver bump major)"
	@git push origin --tags

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
	@echo "  deps-update     - Update all dependencies"
	@echo "  deps-check      - Check for outdated dependencies"
	@echo "  deps-audit      - Audit dependencies for vulnerabilities"
	@echo "  deps-clean      - Clean module cache"
	@echo "  dev             - Run with hot reload (requires air)"
	@echo "  test-db-setup   - Setup test database"
	@echo "  check           - Run all checks (format, lint, test)"
	@echo "  ci-test         - Run CI tests"
	@echo "  ci-build        - Run CI build"
	@echo "  ci-lint         - Run CI linting"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  security-scan   - Run security scan"
	@echo "  ci-local        - Run local CI simulation"
	@echo "  actions-test    - Test GitHub Actions locally"
	@echo "  actions-validate- Validate GitHub Actions workflows"
	@echo "  bench           - Run benchmarks"
	@echo "  bench-cpu       - Run CPU profiling"
	@echo "  bench-memory    - Run memory profiling"
	@echo "  docs-serve      - Serve documentation"
	@echo "  release-patch   - Create patch release"
	@echo "  release-minor   - Create minor release"
	@echo "  release-major   - Create major release" 