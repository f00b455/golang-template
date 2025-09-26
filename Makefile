# Golang Template Makefile

.PHONY: build dev test test-bdd lint format clean install deps

# Go settings
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build settings
BINARY_NAME=golang-template
BINARY_PATH=./bin/$(BINARY_NAME)
API_BINARY=api-server
CLI_BINARY=cli-tool
WEB_BINARY=web-server

# Directories
SRC_DIR=./...
API_DIR=./cmd/api
CLI_DIR=./cmd/cli
WEB_DIR=./cmd/web
PKG_DIR=./pkg/...
INTERNAL_DIR=./internal/...

# Default target
all: test build

# Install dependencies
deps:
	$(GOMOD) tidy
	$(GOMOD) download

# Generate Swagger documentation
docs: deps
	@echo "Generating Swagger documentation..."
	@GOPATH_BIN=$$(go env GOPATH)/bin; \
	if [ ! -f $$GOPATH_BIN/swag ]; then \
		echo "Installing swag..." && $(GOGET) github.com/swaggo/swag/cmd/swag; \
	fi; \
	$$GOPATH_BIN/swag init -g cmd/api/main.go -o docs/

# Build all binaries
build: docs deps
	@echo "Building API server..."
	$(GOBUILD) -o ./bin/$(API_BINARY) $(API_DIR)
	@echo "Building CLI tool..."
	$(GOBUILD) -o ./bin/$(CLI_BINARY) $(CLI_DIR)
	@echo "Building web server..."
	$(GOBUILD) -o ./bin/$(WEB_BINARY) $(WEB_DIR)

# Development mode (with auto-reload would require additional tooling)
dev: deps
	@echo "Running API server in development mode..."
	$(GOCMD) run $(API_DIR)/main.go

# Run web server in development mode
dev-web: deps
	@echo "Running web server in development mode..."
	$(GOCMD) run $(WEB_DIR)/main.go

# Install Hugo
install-hugo:
	@echo "Installing Hugo..."
	@bash scripts/install-hugo.sh

# Create Hugo site (if not exists)
hugo-init: install-hugo
	@if [ ! -f site/hugo.toml ]; then \
		echo "Creating Hugo site..."; \
		./bin/hugo new site site --force; \
	else \
		echo "Hugo site already exists"; \
	fi

# Build Hugo site
hugo-build: install-hugo
	@echo "Building Hugo site..."
	./bin/hugo -s site

# Serve Hugo site (port 1313)
hugo-serve: install-hugo
	@echo "Starting Hugo server on port 1313..."
	./bin/hugo server -s site -p 1313

# Run both API server and Hugo server
dev-full:
	@echo "Starting both API server (port 3002) and Hugo server (port 1313)..."
	@trap 'kill %1; kill %2' INT; \
	$(GOCMD) run $(API_DIR)/main.go & \
	./bin/hugo server -s site -p 1313 & \
	wait

# Clean Hugo build
hugo-clean:
	@echo "Cleaning Hugo build..."
	rm -rf site/public/

# Run unit tests
test: deps
	@echo "Running unit tests..."
	$(GOTEST) -v -race $(SRC_DIR)


# Run BDD tests (Cucumber with Godog)
test-bdd: deps build
	@echo "Running BDD tests..."
	@echo "Building binaries first..."
	@echo "Running shared package BDD tests..."
	$(GOCMD) test -v ./features/ -run TestSharedFeatures
	@echo "Running API endpoint BDD tests..."
	$(GOCMD) test -v ./features/ -run TestAPIFeatures

# Run all BDD tests
test-bdd-all: deps build
	@echo "Running all BDD feature tests..."
	$(GOCMD) test -v ./features/...

# Run CLI BDD tests
test-bdd-cli: deps build
	@echo "Running CLI BDD tests..."
	$(GOCMD) test -v ./features/ -run TestCLIFeatures

# Run Web BDD tests
test-bdd-web: deps build
	@echo "Running Web BDD tests..."
	$(GOCMD) test -v ./features/ -run TestWebFeatures

# Lint the code
lint: deps
	@echo "Running linter..."
	@if [ -f ./bin/golangci-lint ]; then \
		./bin/golangci-lint run --no-config; \
	else \
		echo "Installing golangci-lint locally..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin && \
		./bin/golangci-lint run --no-config; \
	fi


# Clean build artifacts
clean: hugo-clean
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf ./bin/
	rm -f coverage.out coverage.html

# Install tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) golang.org/x/tools/cmd/goimports
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin

# Full validation pipeline (equivalent to npm run validate)
validate: lint test-coverage-check test-bdd build
	@echo "✅ All validation checks passed!"
	@echo "🎉 Safe to push to remote repository"

# Quick validation (equivalent to npm run validate:quick)
validate-quick: lint test
	@echo "✅ Quick validation checks passed!"

# Format code with gofmt and goimports
format:
	@echo "Formatting code..."
	$(GOFMT) $(SRC_DIR)
	@which goimports > /dev/null || $(GOGET) golang.org/x/tools/cmd/goimports
	goimports -w .

# Run tests with coverage
test-cover: deps
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out $(SRC_DIR)
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	$(GOCMD) tool cover -func=coverage.out

# Run tests with coverage and validate threshold (95% minimum for production code)
test-coverage-check: deps
	@echo "Running tests with coverage validation..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic $(SRC_DIR)
	@chmod +x scripts/check-coverage.sh
	@./scripts/check-coverage.sh 95.0 coverage.out

# Auto-fix linting issues where possible
lint-fix: deps
	@echo "Auto-fixing linting issues..."
	@if [ -f ./bin/golangci-lint ]; then \
		./bin/golangci-lint run --fix --no-config; \
	else \
		echo "Installing golangci-lint locally..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin && \
		./bin/golangci-lint run --fix --no-config; \
	fi

# Run specific package tests
test-pkg: deps
	@if [ -z "$(PKG)" ]; then \
		echo "Usage: make test-pkg PKG=./pkg/core"; \
		exit 1; \
	fi
	@echo "Running tests for package: $(PKG)"
	$(GOTEST) -v -race $(PKG)/...

# Build with version information
build-with-version: deps
	@echo "Building with version information..."
	@VERSION=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev"); \
	BUILD_TIME=$$(date -u '+%Y-%m-%d_%H:%M:%S'); \
	$(GOBUILD) -ldflags "-X main.Version=$$VERSION -X main.BuildTime=$$BUILD_TIME" -o ./bin/$(API_BINARY) $(API_DIR); \
	$(GOBUILD) -ldflags "-X main.Version=$$VERSION -X main.BuildTime=$$BUILD_TIME" -o ./bin/$(CLI_BINARY) $(CLI_DIR)

# Run benchmarks
bench: deps
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem $(SRC_DIR)

# Generate mocks (if using mockery or similar)
generate-mocks:
	@echo "Generating mocks..."
	@which mockery > /dev/null || (echo "Installing mockery..." && $(GOGET) github.com/vektra/mockery/v2/.../mockery)
	go generate ./...

# Security scan
security: deps
	@echo "Running security scan..."
	@which gosec > /dev/null || (echo "Installing gosec..." && $(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec)
	gosec ./...

# Dependency check
deps-check:
	@echo "Checking for outdated dependencies..."
	$(GOCMD) list -u -m all

# Docker build (if Dockerfile exists)
docker-build:
	@if [ -f Dockerfile ]; then \
		echo "Building Docker image..."; \
		docker build -t golang-template .; \
	else \
		echo "No Dockerfile found"; \
	fi

# Development setup
setup: install-tools deps
	@echo "Development environment setup complete!"

# Help target
help:
	@echo "🔨 Golang Template - Available Make Targets:"
	@echo ""
	@echo "📦 Build Commands:"
	@echo "  build              - Build all binaries (API server and CLI tool)"
	@echo "  build-with-version - Build with Git version and build time info"
	@echo "  docs               - Generate Swagger API documentation"
	@echo "  clean              - Clean build artifacts and binaries"
	@echo ""
	@echo "🚀 Development Commands:"
	@echo "  dev                - Run API server in development mode"
	@echo "  dev-web            - Run web server in development mode"
	@echo "  dev-full           - Run both API server and Hugo site concurrently"
	@echo "  setup              - Setup development environment (install tools + deps)"
	@echo ""
	@echo "🌐 Hugo Commands:"
	@echo "  install-hugo       - Install Hugo binary"
	@echo "  hugo-init          - Initialize Hugo site (if not exists)"
	@echo "  hugo-build         - Build Hugo static site"
	@echo "  hugo-serve         - Serve Hugo site on port 1313"
	@echo "  hugo-clean         - Clean Hugo build artifacts"
	@echo ""
	@echo "🧪 Testing Commands:"
	@echo "  test               - Run unit tests"
	@echo "  test-cover         - Run tests with coverage report"
	@echo "  test-coverage-check - Run tests with 95% coverage validation"
	@echo "  test-bdd           - Run BDD/Cucumber tests"
	@echo "  test-pkg PKG=path  - Run tests for specific package"
	@echo "  bench              - Run benchmarks"
	@echo ""
	@echo "🔍 Quality Assurance:"
	@echo "  lint               - Run golangci-lint"
	@echo "  lint-fix           - Auto-fix linting issues where possible"
	@echo "  format             - Format code with gofmt and goimports"
	@echo "  security           - Run security scan with gosec"
	@echo ""
	@echo "✅ Validation Pipelines:"
	@echo "  validate           - Full validation (lint + test + test-bdd + build)"
	@echo "  validate-quick     - Quick validation (lint + test)"
	@echo ""
	@echo "🔧 Maintenance Commands:"
	@echo "  deps               - Install Go dependencies"
	@echo "  deps-check         - Check for outdated dependencies"
	@echo "  install-tools      - Install development tools"
	@echo "  generate-mocks     - Generate mocks (requires mockery)"
	@echo "  docker-build       - Build Docker image (if Dockerfile exists)"
	@echo ""
	@echo "❓ Help:"
	@echo "  help               - Show this help message"
	@echo ""
	@echo "💡 Examples:"
	@echo "  make validate-quick          # Quick check during development"
	@echo "  make test-pkg PKG=./pkg/core # Test only core package"
	@echo "  make test-cover              # Generate coverage report"