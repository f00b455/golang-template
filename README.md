# Golang Template

A clean, modern Go project template featuring:

- **HTTP API Server** with Gin framework
- **CLI Application** with Cobra framework
- **Clean Architecture** with separated concerns
- **Comprehensive Testing** with unit tests and BDD
- **CI/CD Pipeline** with GitHub Actions
- **Code Quality** with golangci-lint

## Project Structure

```
golang-template/
├── cmd/
│   ├── api/          # HTTP API server
│   └── cli/          # CLI application
├── pkg/
│   ├── shared/       # Shared utilities and types
│   └── core/         # Core business logic
├── internal/
│   ├── config/       # Configuration management
│   ├── handlers/     # HTTP handlers
│   └── middleware/   # HTTP middleware
├── features/         # BDD test features
├── bin/             # Built binaries
└── docs/            # Documentation
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/f00b455/golang-template.git
cd golang-template
```

2. Install dependencies:
```bash
make deps
# or
go mod download
```

### Development

#### Running the API Server

```bash
# Development mode
make dev
# or
go run cmd/api/main.go

# Build and run
make build
./bin/api-server
```

The API server will start on `http://localhost:3002` with:
- API endpoints at `/api/`
- Swagger documentation at `/documentation/`

#### Running the CLI Tool

```bash
# Development mode
go run cmd/cli/main.go --name "Your Name"

# Build and run
make build
./bin/cli-tool --name "Your Name"
```

### Available Commands

```bash
# Development
make dev              # Run API server in development mode
make build            # Build all binaries
make clean            # Clean build artifacts

# Testing
make test             # Run unit tests
make test-cover       # Run tests with coverage
make test-bdd         # Run BDD tests

# Code Quality
make lint             # Run linter
make format           # Format code

# Validation
make validate         # Full validation pipeline
make validate-quick   # Quick validation (lint + test)

# Setup
make setup            # Install dev tools and dependencies
```

## API Endpoints

### Greet API

- **GET** `/api/greet?name=World` - Get greeting message

### RSS API

- **GET** `/api/rss/spiegel/latest` - Get latest SPIEGEL headline
- **GET** `/api/rss/spiegel/top5?limit=3` - Get top N headlines (max 5)

## CLI Usage

```bash
# Basic greeting
./bin/cli-tool

# Custom name
./bin/cli-tool --name "Alice"

# Help
./bin/cli-tool --help
```

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### BDD Tests

```bash
# Run BDD features
go test ./features/...

# Run specific feature
go test ./features/ -godog.tags="@issue-1"
```

## Architecture

### Clean Code Principles

- **Single Responsibility**: Each package/function has one purpose
- **Dependency Injection**: Handlers receive their dependencies
- **Pure Functions**: Business logic is stateless where possible
- **Explicit Error Handling**: All errors are handled explicitly
- **Immutable Data**: Prefer immutable data structures

### Package Organization

- `cmd/` - Application entry points
- `pkg/` - Public APIs and libraries
- `internal/` - Private application code
- `features/` - BDD test specifications

### Testing Strategy

- **Unit Tests**: Standard Go testing with testify
- **BDD Tests**: Godog for behavior-driven development
- **Mocking**: Always use mocks for external dependencies
- **Test Database**: Use test database for integration tests

## Configuration

Set environment variables:

```bash
PORT=3002                    # API server port
ENV=development             # Environment (development/production)
SPIEGEL_RSS_URL=https://...  # RSS feed URL
GO_ENV=test                 # For testing (shorter delays)
```

## CI/CD

The project includes a GitHub Actions workflow that:

1. **Lints** code with golangci-lint
2. **Tests** with unit tests and BDD tests
3. **Builds** both API server and CLI binaries
4. **Uploads** coverage to Codecov

## Contributing

1. Follow the existing code style
2. Write tests for new features
3. Update BDD features for new user stories
4. Run `make validate` before submitting PRs

## License

MIT License - see LICENSE file for details.# Trigger pipeline workflow
