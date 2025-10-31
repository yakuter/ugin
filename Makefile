.PHONY: help build run clean test lint fmt vet docker-build docker-run

# Default target
help:
	@echo "Available targets:"
	@echo "  build              - Build the application binary"
	@echo "  run                - Build and run the application"
	@echo "  clean              - Remove build artifacts"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  lint               - Run linter (requires golangci-lint)"
	@echo "  fmt                - Format code"
	@echo "  vet                - Run go vet"
	@echo "  build-image        - Build Docker image"
	@echo "  run-app-mysql      - Run with MySQL using docker-compose"
	@echo "  clean-app-mysql    - Stop MySQL docker-compose"
	@echo "  run-app-postgres   - Run with PostgreSQL using docker-compose"
	@echo "  clean-app-postgres - Stop PostgreSQL docker-compose"

# Build variables
BINARY_NAME=ugin
BINARY_PATH=./bin/$(BINARY_NAME)
CMD_PATH=./cmd/ugin
MAIN_FILE=$(CMD_PATH)/main.go

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(MAIN_FILE)
	@echo "Build complete: $(BINARY_PATH)"

# Build with optimizations (smaller binary)
build-prod:
	@echo "Building $(BINARY_NAME) for production..."
	@mkdir -p bin
	@CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY_PATH) $(MAIN_FILE)
	@echo "Production build complete: $(BINARY_PATH)"

# Run the application
run: build
	@echo "Starting $(BINARY_NAME)..."
	@$(BINARY_PATH)

# Run without building (development)
run-dev:
	@echo "Running in development mode..."
	@go run $(MAIN_FILE)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin
	@rm -f ugin ugin.db ugin.log ugin.db.log ugin.access.log
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run all checks
check: fmt vet test
	@echo "All checks passed!"

# Docker targets
build-image:
	@echo "Building Docker image..."
	@docker build -t ugin -f containers/images/Dockerfile .

run-app-mysql:
	@echo "Starting application with MySQL..."
	@docker-compose -f containers/composes/dc.mysql.yml up

clean-app-mysql:
	@echo "Stopping MySQL containers..."
	@docker-compose -f containers/composes/dc.mysql.yml down

run-app-postgres:
	@echo "Starting application with PostgreSQL..."
	@docker-compose -f containers/composes/dc.postgres.yml up

clean-app-postgres:
	@echo "Stopping PostgreSQL containers..."
	@docker-compose -f containers/composes/dc.postgres.yml down

# Install dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
