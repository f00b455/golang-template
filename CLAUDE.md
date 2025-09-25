# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Clean Code Principles

### Code Quality Standards:

- **Single Responsibility**: Each function/struct should have one reason to change
- **Pure Functions**: Prefer functions without side effects when possible
- **Explicit Types**: Use Go's type system fully - avoid `interface{}` where possible
- **Descriptive Naming**: Use clear, searchable names for variables and functions
- **Small Functions**: Keep functions focused and under 20 lines when possible (enforced by golangci-lint)
- **No Comments for What**: Code should be self-documenting; comments explain why, not what

### Golangci-lint Clean Code Enforcement:

The codebase uses golangci-lint rules to automatically enforce Clean Code principles:

**Function Complexity Limits:**
- **Cyclomatic complexity**: 10 max (number of code paths)
- **Cognitive complexity**: 10 max (mental effort to understand)
- **Function length**: 60 lines max (enforces Single Responsibility)
- **Rationale**: Smaller functions are easier to test, understand, and maintain

**File and Package Organization:**
- **File length**: 500 lines max (promotes modular design)
- **Package cohesion**: Related functionality in same package
- **Import organization**: Standard library, third-party, local imports
- **Rationale**: Well-organized code is easier to navigate and maintain

**Code Quality Rules:**
- **Error handling**: All errors must be handled or explicitly ignored
- **Naming conventions**: Go standard naming (camelCase, PascalCase)
- **No magic numbers**: Use named constants for numeric literals
- **Interface segregation**: Small, focused interfaces

To run linting:
```bash
make lint                    # Run all linters
./bin/golangci-lint run     # Run linter directly
make lint-fix               # Auto-fix linting issues where possible
```

### Pure Function Principles:

- **Deterministic**: Same input always produces same output
- **No Side Effects**: Don't modify external state or global variables
- **No External Dependencies**: Don't rely on external mutable state
- **Immutable Parameters**: Don't mutate input parameters (use pointers carefully)
- **Referential Transparency**: Function calls can be replaced with their return values
- **Predictable**: Easy to test, debug, and reason about
- **Examples of Pure Functions**:

  ```go
  // ✅ Pure - deterministic, no side effects
  func Add(a, b int) int { return a + b }
  func FormatUser(user User) string { return fmt.Sprintf("%s (%s)", user.Name, user.Email) }

  // ❌ Impure - side effects, external dependencies
  func LogAndAdd(a, b int) int {
    log.Println("Adding numbers") // Side effect
    return a + b
  }
  func GetCurrentUser() User { return database.CurrentUser } // External dependency
  ```

- **When to Use Pure Functions**: Data transformations, calculations, formatting, validation
- **When Impure is Acceptable**: I/O operations, API calls, database queries, logging

### Testing Principles:

- **Test-Driven Development**: Write tests first when adding new features
- **Mock External Dependencies**: Always use mocks for external services/APIs
- **Test Database**: Always use test database, never production data
- **Arrange-Act-Assert**: Structure tests clearly with setup, execution, and verification
- **Descriptive Test Names**: Test names should describe the behavior being tested
- **Table-driven tests**: Use Go's table-driven test pattern for multiple test cases
- **BDD for Packages**: Every package must have feature files describing user stories and business requirements

### Architecture Guidelines:

- **Dependency Injection**: Prefer injecting dependencies over tight coupling
- **Error Handling**: Use explicit error types and handle all error cases
- **Immutability**: Prefer immutable data structures and operations
- **Separation of Concerns**: Keep business logic separate from framework code
- **API Design**: RESTful endpoints with proper HTTP status codes
- **Type Safety**: Leverage Go's compile-time type checking

### Performance Considerations:

- **Memory Allocation**: Minimize heap allocations in hot paths
- **Goroutine Management**: Use worker pools for concurrent processing
- **Database Queries**: Use efficient queries and avoid N+1 problems
- **Caching**: Implement appropriate caching strategies
- **Profiling**: Use Go's built-in profiling tools (pprof)

## Project Structure

This is a pure Go project with a clean architecture for developing APIs, CLI applications, and reusable packages:

```
golang-template/
├── cmd/                  # Application entry points
│   ├── api/              # REST API server
│   └── cli/              # CLI application
├── internal/             # Private application code
│   ├── config/           # Configuration management
│   ├── handlers/         # HTTP request handlers
│   └── middleware/       # HTTP middleware
├── pkg/                  # Public reusable packages
│   ├── core/             # Core business logic
│   └── shared/           # Shared utilities and types
├── features/             # BDD feature files and tests
├── scripts/              # Build and deployment scripts
├── docs/                 # Documentation
├── bin/                  # Compiled binaries
├── .github/              # CI/CD workflows
└── .husky/               # Git hooks
```

## Development Commands

### Primary Commands (run these from project root):

- `make dev` - Start API server in development mode
- `make build` - Build all binaries (API server and CLI tool)
- `make test` - Run unit tests across all packages
- `make test-bdd` - Run BDD/Cucumber tests
- `make lint` - Run golangci-lint across all packages
- `make format` - Format code with gofmt and goimports
- `make clean` - Clean build artifacts and binaries

### Development Workflow:

**AUTOMATED PRE-PUSH VALIDATION**: This repo uses Husky to automatically validate code before each push.
The pre-push hook runs: `make validate` (lint + test + test-bdd + build)

**Development Commands**:
- `make validate-quick` - Quick validation (lint + test) - **Use during development**
- `make validate` - Full pipeline validation (lint + test + test-bdd + build) - **Runs automatically on push**

**Efficient Dev Cycle**:
1. Code freely and commit often (no pre-commit validation)
2. Run `make validate-quick` during development for fast feedback
3. When ready to push: `git push` - automatic validation prevents broken CI/CD
4. If validation fails, fix issues and push again

### Package-Specific Commands:

- `go run ./cmd/api` - Run API server directly
- `go run ./cmd/cli --help` - Run CLI tool with help
- `go test ./pkg/core/...` - Test core package only
- `go test -race ./...` - Run all tests with race detection

## Architecture

### Project Organization:

- **cmd/**: Application entry points following Go project layout standards
- **internal/**: Private application code, not importable by other projects
- **pkg/**: Public packages that can be imported by other projects
- **Makefile**: Centralized build, test, and validation commands

### API Server (cmd/api):

- **Gin**: Fast HTTP router and middleware framework
- **Swagger**: API documentation (available at /documentation)
- **Graceful shutdown**: Proper cleanup on termination signals
- **CORS**: Configured for frontend integration
- **Structured logging**: JSON logging with levels

### CLI Tool (cmd/cli):

- **Cobra**: Command-line interface framework
- **Progress bars**: Visual feedback for long operations
- **Colorized output**: Enhanced user experience
- **Configuration**: Support for config files and environment variables

### Core Packages:

- **pkg/core/**: Business logic and domain models
- **pkg/shared/**: Common utilities and helper functions
- **internal/**: Application-specific implementations

### Testing Strategy:

- **Unit Tests**: Standard Go testing with testify for assertions
- **BDD Tests**: Godog (Cucumber for Go) with Gherkin syntax
- **Test doubles**: Mocks for external dependencies
- **Table-driven tests**: Go idiom for testing multiple scenarios
- **Always use test database** (important requirement)

### CI/CD:

- **GitHub Actions**: Separate jobs for linting, testing, building
- **Automated validation**: Pre-push hooks prevent broken builds
- **Build artifacts**: Compiled binaries for multiple platforms
- **Code coverage**: Coverage reports and thresholds

## BDD Requirements for Packages

### Mandatory Feature Files:

Every package in `pkg/` MUST include:

- `.feature` files in the `features/` directory
- Gherkin scenarios describing user stories and business requirements
- Step definitions in Go test files
- Integration with `make test-bdd` command
- **GitHub Issue references** (see below)

### Feature File Structure with Issue References:

```gherkin
# Issue: #<ISSUE_NUMBER>
# URL: https://github.com/<OWNER>/<REPO>/issues/<ISSUE_NUMBER>
@pkg(<pkg>) @issue-<ISSUE_NUMBER>
Feature: [Feature Name]
  As a [user type]
  I want to [goal]
  So that [benefit]

  Background:
    Given [common setup]

  @happy-path
  Scenario: [Scenario description]
    Given [precondition]
    When [action]
    Then [expected result]
```

### Issue Reference Requirements:

**EVERY feature file in the project MUST contain:**

1. **Header Comments** (first 2 lines):
   - `# Issue: #<number>` - GitHub issue number reference
   - `# URL: https://github.com/<owner>/<repo>/issues/<number>` - Full GitHub issue URL

2. **Tags** (on the Feature line):
   - `@issue-<number>` - Issue tag for filtering and tracking
   - `@pkg(<package-name>)` - Package identifier tag (e.g., @pkg(core), @pkg(shared))

3. **CI Verification**:
   - The CI pipeline automatically verifies all feature files have proper issue references
   - Missing references will cause the CI to fail with specific error messages
   - This ensures full traceability between user stories/issues and BDD tests

4. **Step Definitions Implementation**:
   - When creating new feature files, step definitions MUST be implemented in Go
   - Use Godog framework for step definitions:

   ```go
   func (ctx *FeatureContext) iHaveTheInput(input string) error {
       ctx.input = input
       return nil
   }

   func InitializeScenario(ctx *godog.ScenarioContext) {
       featureCtx := &FeatureContext{}
       ctx.Step(`^I have the input "([^"]*)"$`, featureCtx.iHaveTheInput)
   }
   ```

### Coverage Requirements:

- **One feature file per user story/issue minimum**
- Cover all public API functions
- Test pure function properties (determinism, immutability)
- Include error handling scenarios
- Verify integration points between packages

### Example Package BDD Structure:

```
pkg/core/
├── features/
│   ├── foo-processing.feature      # Core processing functionality
│   ├── foo-greeting.feature        # Integration features
│   └── foo-data-operations.feature # Data transformation features
├── foo_test.go                     # Unit tests
├── foo_bdd_test.go                 # BDD step definitions
└── foo.go                          # Implementation
```

### Running BDD Tests:

- All packages: `make test-bdd` (from root)
- Specific package: `go test -v ./pkg/core/`
- With coverage: `make test-cover`
- CI pipeline includes BDD test execution

## Important Notes:

- Always use test database for integration tests
- Use mocks in tests as specified in user requirements
- Go modules enabled with semantic versioning
- Follow Go project layout standards
- **BDD is mandatory for all public packages**
- Use `gofmt` and `goimports` for consistent formatting
- Handle all errors explicitly - no silent failures