# Implement Go Feature Playbook

## When asked to implement a feature in this Go project:

### 1. **Follow BDD-First Approach (RED-GREEN-REFACTOR)**:
   - **RED Phase**: Write feature file FIRST (before any code)
   - Create `.feature` file in `features/` directory
   - Include issue reference: `# Issue: #<number>`
   - Run tests to see them fail: `make test-bdd`

### 2. **Implement Step Definitions**:
   - Create `*_test.go` file with Godog step definitions
   - Map Gherkin steps to Go functions
   - Run tests again (should still fail, no implementation yet)

### 3. **GREEN Phase - Minimal Implementation**:
   - **For API features**:
     - Add handler in `internal/handlers/`
     - Add route in `cmd/api/main.go`
     - Update Swagger docs with `swag init`
   - **For CLI features**:
     - Add command in `cmd/cli/`
     - Use Cobra for command structure
   - Write MINIMAL code to make tests pass

### 4. **REFACTOR Phase**:
   - Improve code quality while keeping tests green
   - Extract common logic to `pkg/`
   - Add proper error handling
   - Follow Go idioms and best practices

### 5. **Testing**:
   - Run unit tests: `make test`
   - Run BDD tests: `make test-bdd`
   - Run linter: `make lint`
   - Full validation: `make validate`

### 6. **Git Workflow**:
   - Commit with descriptive message
   - Push to feature branch
   - Create/update PR
   - Link to issue with "Closes #<number>"

## Go-Specific Guidelines:
- Use table-driven tests for multiple test cases
- Handle all errors explicitly
- Keep functions small (max 60 lines)
- Follow standard Go project layout
- Use interfaces for dependency injection
- Mock external dependencies in tests

## Available Make Commands:
- `make dev` - Start API server in development
- `make build` - Build all binaries
- `make test` - Run unit tests
- `make test-bdd` - Run BDD/Cucumber tests
- `make lint` - Run golangci-lint
- `make validate` - Full validation pipeline