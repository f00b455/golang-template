# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Clean Code Principles

### Code Quality Standards:

- **Single Responsibility**: Each function/class should have one reason to change
- **Pure Functions**: Prefer functions without side effects when possible
- **Explicit Types**: Use TypeScript's type system fully - avoid `any`
- **Descriptive Naming**: Use clear, searchable names for variables and functions
- **Small Functions**: Keep functions focused and under 20 lines when possible (enforced by ESLint)
- **No Comments for What**: Code should be self-documenting; comments explain why, not what

### ESLint Clean Code Enforcement:

The codebase uses ESLint rules to automatically enforce Clean Code principles:

**Function Length Limits:**
- **Regular code**: 20 lines max per function (enforces Single Responsibility)
- **Test files**: 50 lines max (tests need more setup/assertions)
- **Rationale**: Smaller functions are easier to test, understand, and maintain

**File Length Limits:**
- **Regular files**: 300 lines max (promotes modular design)
- **Test files**: 500 lines max (test suites can be larger)
- **Rationale**: Prevents monolithic files that are hard to navigate

**Complexity Limits:**
- **Cyclomatic complexity**: 10 max (number of code paths)
- **Nesting depth**: 3 levels max (encourages early returns)
- **Parameters**: 4 max per function (use objects for complex data)
- **Statements**: 15 max per function (complements line limit)

**Phase-Based Migration:**
Set `ESLINT_CLEAN_CODE_PHASE` environment variable:
- **Phase 1** (default): `warn` - Assessment phase to identify violations
- **Phase 2**: `warn` - Progressive refactoring with CI stability
- **Phase 3**: `error` - Full enforcement, violations break the build

To assess current violations:
```bash
# Run with default Phase 1 (warnings)
pnpm lint

# Generate detailed report
pnpm lint:report

# Run with Phase 3 enforcement
ESLINT_CLEAN_CODE_PHASE=3 pnpm lint
```

### Pure Function Principles:

- **Deterministic**: Same input always produces same output
- **No Side Effects**: Don't modify external state, DOM, or global variables
- **No External Dependencies**: Don't rely on external mutable state
- **Immutable Parameters**: Don't mutate input parameters
- **Referential Transparency**: Function calls can be replaced with their return values
- **Predictable**: Easy to test, debug, and reason about
- **Examples of Pure Functions**:

  ```typescript
  // ✅ Pure - deterministic, no side effects
  const add = (a: number, b: number): number => a + b;
  const formatUser = (user: User): string => `${user.name} (${user.email})`;

  // ❌ Impure - side effects, external dependencies
  const logAndAdd = (a: number, b: number): number => {
    console.log('Adding numbers'); // Side effect
    return a + b;
  };
  const getCurrentUser = (): User => database.users.current; // External dependency
  ```

- **When to Use Pure Functions**: Data transformations, calculations, formatting, validation
- **When Impure is Acceptable**: I/O operations, API calls, database queries, logging

### Testing Principles:

- **Test-Driven Development**: Write tests first when adding new features
- **Mock External Dependencies**: Always use mocks for external services/APIs
- **Test Database**: Always use test database, never production data
- **Arrange-Act-Assert**: Structure tests clearly with setup, execution, and verification
- **Descriptive Test Names**: Test names should describe the behavior being tested
- **BDD for Libraries**: Every library must have feature files describing user stories and business requirements

### Architecture Guidelines:

- **Dependency Injection**: Prefer injecting dependencies over tight coupling
- **Error Handling**: Use explicit error types and handle all error cases
- **Immutability**: Prefer immutable data structures and operations
- **Separation of Concerns**: Keep business logic separate from framework code
- **API Design**: RESTful endpoints with proper HTTP status codes
- **Type Safety**: Leverage TypeScript's strict mode for compile-time safety

### Performance Considerations:

- **Bundle Size**: Keep frontend bundle size minimal
- **Database Queries**: Use efficient queries and avoid N+1 problems
- **Caching**: Implement appropriate caching strategies
- **Lazy Loading**: Load resources only when needed

## Project Structure

This is a TypeScript monorepo using pnpm workspaces and Turborepo for developing core functionalities, MCP servers, CLI apps, and libraries:

```
ts-pure-template/
├── apps/
│   ├── web/          # Next.js frontend application
│   ├── api/          # Fastify backend API
│   └── cli/          # CLI applications
├── packages/
│   ├── shared/       # Shared TypeScript utilities and types
│   ├── core/         # Core business logic and functionalities
│   ├── mcp-*/        # MCP (Model Context Protocol) servers
│   └── lib-*/        # Reusable libraries
├── tools/            # Build tools and development utilities
├── .changeset/       # Changesets configuration for versioning
├── .claude/          # Claude Code-Playbooks
└── .github/          # CI/CD workflows
```

## Development Commands

### Root Commands (run these from project root):

- `pnpm dev` - Start all applications in development mode
- `pnpm build` - Build all packages and applications
- `pnpm test` - Run unit tests across all packages
- `pnpm test:e2e` - Run Playwright E2E tests
- `pnpm test:cucumber` - Run Cucumber BDD tests
- `pnpm lint` - Lint all packages
- `pnpm type-check` - TypeScript type checking
- `pnpm format` - Format code with Prettier
- `pnpm changeset` - Create a changeset for versioning

### Development Workflow:

**AUTOMATED PRE-PUSH VALIDATION**: This repo uses Husky to automatically validate code before each push.
The pre-push hook runs: `pnpm validate` (lint + type-check + build + test:run + test:cucumber)

**Development Commands**:
- `pnpm validate:quick` - Quick validation (lint + type-check + tests) - **Use during development**
- `pnpm validate` - Full pipeline validation (lint + type-check + build + tests + cucumber) - **Runs automatically on push**

**Efficient Dev Cycle**:
1. Code freely and commit often (no pre-commit validation)
2. Run `pnpm validate:quick` during development for fast feedback
3. When ready to push: `git push` - automatic validation prevents broken CI/CD
4. If validation fails, fix issues and push again

### Package-Specific Commands:

- `pnpm --filter @ts-template/web dev` - Run Next.js in dev mode
- `pnpm --filter @ts-template/api dev` - Run Fastify API in dev mode
- `pnpm --filter @ts-template/shared test` - Test shared package only

## Architecture

### Monorepo Setup:

- **pnpm**: Package manager with workspace support
- **Turborepo**: Build system and task runner for monorepos
- **Changesets**: Versioning and changelog management

### Frontend (apps/web):

- **Next.js 14**: React framework with App Router
- **Vitest**: Unit testing framework
- **Playwright**: End-to-end testing
- Uses shared package via workspace protocol

### Backend (apps/api):

- **Fastify**: Fast Node.js web framework
- **Swagger**: API documentation (available at /documentation)
- **Vitest**: Unit testing
- **Cucumber**: BDD testing with Gherkin syntax
- CORS configured for Next.js frontend

### Shared Package (packages/shared):

- Common TypeScript utilities and types
- Consumed by both web and API applications
- Includes comprehensive unit tests

### Testing Strategy:

- **Unit Tests**: Vitest with mocks for external dependencies
- **E2E Tests**: Playwright for full application flows
- **BDD Tests**: Cucumber with Gherkin syntax for business requirements
- **Library BDD Requirement**: Every library package MUST have feature files
- **Feature File Coverage**: Each user story/issue requires at least one .feature file
- **Always use test database** (important requirement)

### CI/CD:

- GitHub Actions with separate jobs for linting, testing, building
- Automated versioning with Changesets (private packages)
- Coverage reporting with Codecov
- Release PRs for version management (no npm publishing)

## BDD Requirements for Libraries

### Mandatory Feature Files:

Every library package (packages/lib-_, packages/shared) AND application package (apps/_) MUST include:

- `.feature` files in the `features/` directory
- Gherkin scenarios describing user stories and business requirements
- Step definitions in `features/step_definitions/`
- Cucumber configuration (`cucumber.js`)
- `test:cucumber` script in package.json
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

**EVERY feature file in the monorepo MUST contain:**

1. **Header Comments** (first 2 lines):
   - `# Issue: #<number>` - GitHub issue number reference
   - `# URL: https://github.com/<owner>/<repo>/issues/<number>` - Full GitHub issue URL

2. **Tags** (on the Feature line):
   - `@issue-<number>` - Issue tag for filtering and tracking
   - `@pkg(<package-name>)` - Package identifier tag (e.g., @pkg(lib-foo), @pkg(api), @pkg(web))

3. **CI Verification**:
   - The CI pipeline automatically verifies all feature files have proper issue references
   - Missing references will cause the CI to fail with specific error messages
   - This ensures full traceability between user stories/issues and BDD tests

4. **Example for Issue #7**:

   ```gherkin
   # Issue: #7
   # URL: https://github.com/f00b455/ts-pure-template/issues/7
   @pkg(lib-foo) @issue-7
   Feature: lib-foo – Processing Functions
   ```

5. **Step Definitions with Dummy Implementations**:
   - When creating new feature files, step definitions MUST be implemented
   - Unimplemented steps should throw an exception with a clear message:

   ```typescript
   function unimplemented(step: string): never {
     throw new Error(`UNIMPLEMENTED_STEP: ${step} — please implement.`);
   }

   Given('<condition>', () => unimplemented('Given <condition>'));
   ```

### Coverage Requirements:

- **One feature file per user story/issue minimum**
- Cover all public API functions
- Test pure function properties (determinism, immutability)
- Include error handling scenarios
- Verify integration points between libraries

### Example Library BDD Structure:

```
packages/lib-foo/
├── features/
│   ├── foo-processing.feature      # Core processing functionality
│   ├── foo-greeting.feature        # Integration features
│   ├── foo-data-operations.feature # Data transformation features
│   └── step_definitions/
│       └── foo.steps.ts            # Package-specific step implementations
├── cucumber.cjs                    # Cucumber configuration
└── package.json                    # Includes test:cucumber script
```

### Shared Cucumber Steps:

To avoid duplication, common step definitions are centralized in `packages/cucumber-shared/`:

```
packages/cucumber-shared/
├── src/
│   ├── index.ts                    # Re-exports all steps and types
│   └── steps/
│       ├── common.steps.ts         # Generic Given/When/Then steps
│       ├── assertions.steps.ts     # Common assertion patterns
│       ├── data-operations.steps.ts # Array and data manipulation steps
│       └── api.steps.ts            # API testing common steps
├── package.json
└── tsup.config.ts                  # Build configuration
```

**Using Shared Steps:**
- All cucumber configurations automatically include shared steps
- Package-specific steps extend or override shared functionality
- Shared steps use a common World context interface
- Packages can extend the World interface for local state

### Running BDD Tests:

- Individual package: `pnpm --filter @ts-template/lib-foo test:cucumber`
- All packages: `pnpm test:cucumber` (from root)
- CI pipeline includes Cucumber test execution

## Important Notes:

- Always use test database for integration tests
- Use mocks in tests as specified in user requirements
- TypeScript strict mode enabled across all packages
- Shared tsconfig.base.json for consistent configuration
- **BDD is mandatory for all library packages**
