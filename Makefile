# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=server
BINARY_PATH=bin/$(BINARY_NAME)

# Build flags
BUILD_FLAGS=-ldflags="-s -w"

.PHONY: all build build-air clean test run dev deps help

# Default target
all: clean deps build

# Build the application
build:
	@echo "Building application..."
	@mkdir -p bin
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_PATH) .
	@echo "Build completed: $(BINARY_PATH)"

# Run the built application
run: build
	@echo "Running application..."
	./$(BINARY_PATH)

# Build for Air hot reloading (outputs to tmp/main)
build-air:
	@echo "Building for Air..."
	@mkdir -p tmp
	$(GOBUILD) -o tmp/main .

# Run in development mode with hot reloading
dev:
	@echo "Starting development server with hot reloading..."
	@command -v air >/dev/null 2>&1 || { echo "Air not found. Installing..."; $(GOGET) -u github.com/air-verse/air@latest; }
	air

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf bin/
	@rm -rf tmp/
	@rm -f coverage.out coverage.html
	@echo "Clean completed"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) -u github.com/air-verse/air@latest
	@echo "Tools installed"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not found. Install it first."; exit 1; }
	golangci-lint run

# Database commands
db-create:
	@echo "🔄 Creating database..."
	@if [ -z "$$DB_NAME" ]; then echo "❌ DB_NAME environment variable not set"; exit 1; fi
	@createdb -h $${DB_HOST:-localhost} -p $${DB_PORT:-5432} -U $${DB_USER:-postgres} $$DB_NAME || echo "Database may already exist"
	@echo "✅ Database created: $$DB_NAME"

db-drop:
	@echo "⚠️  Dropping database..."
	@if [ -z "$$DB_NAME" ]; then echo "❌ DB_NAME environment variable not set"; exit 1; fi
	@dropdb -h $${DB_HOST:-localhost} -p $${DB_PORT:-5432} -U $${DB_USER:-postgres} $$DB_NAME --if-exists
	@echo "✅ Database dropped: $$DB_NAME"

db-reset: db-drop db-create
	@echo "✅ Database reset complete"

migrate:
	@echo "🔄 Running database migrations..."
	@cd migrations && go run migrate.go up

migrate-reset:
	@echo "⚠️  Resetting database with migrations..."
	@cd migrations && go run migrate.go reset

migrate-seed:
	@echo "🌱 Seeding database..."
	@cd migrations && go run migrate.go seed

migrate-models:
	@echo "🔄 Creating/updating models..."
	@cd migrations && go run migrate.go create-models

db-setup: db-create migrate migrate-seed
	@echo "🎉 Database setup complete!"

# Docker commands (for future use)
docker-build:
	@echo "Building Docker image..."
	docker build -t vdt-dashboard-backend .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env vdt-dashboard-backend

# Production build
build-prod:
	@echo "Building for production..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) $(BUILD_FLAGS) -a -installsuffix cgo -o $(BINARY_PATH) .
	@echo "Production build completed: $(BINARY_PATH)"

# Cross-platform builds
build-windows:
	@echo "Building for Windows..."
	@mkdir -p bin
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o bin/$(BINARY_NAME).exe .

build-mac:
	@echo "Building for macOS..."
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o bin/$(BINARY_NAME)-mac .

build-all: build build-windows build-mac
	@echo "All platform builds completed"

# Development workflow
setup: deps install-tools
	@echo "Development environment setup completed"
	@echo "Create your .env file and run 'make dev' to start development"

# Check if .env file exists and show configuration example
check-env:
	@if [ ! -f .env ]; then \
		echo "⚠️  .env file not found"; \
		echo ""; \
		echo "📝 Create a .env file with the following configuration:"; \
		echo ""; \
		echo "# ============================================"; \
		echo "# Server Configuration"; \
		echo "# ============================================"; \
		echo "PORT=8080"; \
		echo "ENVIRONMENT=development"; \
		echo "LOG_LEVEL=info"; \
		echo ""; \
		echo "# ============================================"; \
		echo "# Database Configuration"; \
		echo "# ============================================"; \
		echo "DB_HOST=localhost"; \
		echo "DB_PORT=5432"; \
		echo "DB_USER=postgres"; \
		echo "DB_PASSWORD=postgres"; \
		echo "DB_NAME=vdt_dashboard"; \
		echo ""; \
		echo "# Alternative: Use DATABASE_URL instead of individual DB settings"; \
		echo "# DATABASE_URL=postgres://postgres:postgres@localhost:5432/vdt_dashboard"; \
		echo ""; \
		echo "# ============================================"; \
		echo "# Frontend Configuration"; \
		echo "# ============================================"; \
		echo "FRONTEND_URL=http://localhost:3000"; \
		echo ""; \
		echo "# ============================================"; \
		echo "# Optional: Additional Configuration"; \
		echo "# ============================================"; \
		echo "# JWT_SECRET=your-jwt-secret-key"; \
		echo "# REDIS_URL=redis://localhost:6379"; \
		echo "# MAX_SCHEMA_SIZE=100"; \
		echo ""; \
		echo "💡 Copy this to a .env file in your project root"; \
	else \
		echo "✅ .env file found"; \
		echo "Current configuration:"; \
		echo ""; \
		cat .env; \
	fi

# Help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  build          - Build the application"
	@echo "  build-air      - Build for Air hot reloading"
	@echo "  run            - Build and run the application"
	@echo "  dev            - Run in development mode with hot reloading"
	@echo "  clean          - Clean build artifacts"
	@echo ""
	@echo "Testing:"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo ""
	@echo "Dependencies:"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  deps-update    - Update dependencies"
	@echo "  install-tools  - Install development tools"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo ""
	@echo "Database:"
	@echo "  db-create      - Create database"
	@echo "  db-drop        - Drop database"
	@echo "  db-reset       - Drop and recreate database"
	@echo "  migrate        - Run migrations"
	@echo "  migrate-reset  - Reset database with fresh migrations"
	@echo "  migrate-seed   - Seed database with sample data"
	@echo "  migrate-models - Create/update models using GORM"
	@echo "  db-setup       - Complete database setup (create + migrate + seed)"
	@echo ""
	@echo "Environment:"
	@echo "  setup          - Setup development environment"
	@echo "  check-env      - Check if .env file exists"
	@echo ""
	@echo "Production:"
	@echo "  build-prod     - Build for production"
	@echo "  build-all      - Build for all platforms"
	@echo ""
	@echo "  help           - Show this help message" 