# GitHub Actions Workflow Update for Coverage Check

## Required Changes to `.github/workflows/go.yml`

To enable the 95% test coverage validation in the CI/CD pipeline, the following changes need to be made to the `test` job in `.github/workflows/go.yml`:

### Update the Test Job

Replace the existing test step (around line 55-56) with the following:

```yaml
    - name: Run unit tests with coverage for production code
      run: go test -v -race -coverprofile=coverage.out -covermode=atomic ./pkg/...

    - name: Check test coverage threshold
      run: |
        chmod +x scripts/check-coverage.sh
        ./scripts/check-coverage.sh 95.0 coverage.out
```

### Complete Updated Test Job

Here's the complete updated `test` job configuration:

```yaml
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: ${{ env.GO_VERSION }}

    - name: Install dependencies
      run: go mod download

    - name: Install Hugo
      run: |
        bash scripts/install-hugo.sh
        ./bin/hugo version

    - name: Build binaries for BDD tests
      run: |
        go build -o bin/api-server ./cmd/api
        go build -o bin/cli-tool ./cmd/cli
        go build -o bin/web-server ./cmd/web

    - name: Run unit tests with coverage for production code
      run: go test -v -race -coverprofile=coverage.out -covermode=atomic ./pkg/...

    - name: Check test coverage threshold
      run: |
        chmod +x scripts/check-coverage.sh
        ./scripts/check-coverage.sh 95.0 coverage.out

    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v4
      with:
        file: ./coverage.out
        token: ${{ secrets.CODECOV_TOKEN }}
      if: always()
```

## Key Changes Explained

1. **Coverage Mode**: Added `-covermode=atomic` for accurate coverage in concurrent code
2. **Coverage Check**: Added a new step that runs the coverage validation script
3. **Threshold**: Set to 95.0% as required
4. **Script Permissions**: Ensures the script is executable before running

## How Coverage Check Works

The coverage check script (`scripts/check-coverage.sh`):
- Tests only production packages in `pkg/` directory
- Calculates coverage percentage from the coverage profile
- Compares against the 95% threshold
- Fails the pipeline if coverage is below threshold
- Provides detailed reporting by package
- Shows files with lowest coverage for improvement guidance

**Note**: The coverage check focuses exclusively on `pkg/` packages as they contain the core production code. Other directories (cmd/, internal/, etc.) are tested separately but not included in the coverage threshold validation.

## Testing the Changes Locally

Before pushing, you can test the coverage check locally:

```bash
# Run tests with coverage check
make test-coverage-check

# Or run the validation pipeline
make validate
```

## Coverage Scope

The coverage check focuses on production packages:
- **Included**: All packages under `pkg/` directory (core business logic)
- **Excluded**: Non-production code:
  - `cmd/` - Entry point files with minimal logic
  - `internal/` - Internal application code
  - `docs/` - Documentation and generated files
  - `scripts/` - Build and deployment scripts
  - `features/` - BDD test files
  - Other test and generated files

This approach ensures the 95% threshold applies to the most critical production code where high test coverage provides the most value.

## Benefits

- Ensures 95% minimum test coverage for production code
- Provides clear visibility into coverage gaps
- Fails fast in CI/CD if coverage drops
- Detailed reporting helps identify areas needing tests
- Maintains high code quality standards