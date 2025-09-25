# Clean Code Migration Plan

## Overview
This document outlines the phased approach to migrate the codebase to comply with Clean Code principles enforced through ESLint rules.

## Current State Assessment

Run the following command to assess current violations:
```bash
pnpm lint:report
```

This generates a detailed report showing:
- Files with violations
- Violation count by rule type
- Package-level statistics
- Critical files requiring immediate attention

## Migration Phases

### Phase 1: Assessment (Current)
**Duration**: 1-2 weeks
**ESLint Mode**: `warn`
**Environment**: `ESLINT_CLEAN_CODE_PHASE=1` (default)

**Goals**:
- Identify all violations across the codebase
- Estimate refactoring effort
- Prioritize packages for migration
- Create refactoring tickets/issues

**Actions**:
1. Run `pnpm lint:report` to generate violation report
2. Review critical files (3+ violations)
3. Create GitHub issues for each package refactoring
4. Document common patterns that need refactoring

### Phase 2: Gradual Refactoring
**Duration**: 4-6 weeks
**ESLint Mode**: `warn`
**Environment**: `ESLINT_CLEAN_CODE_PHASE=2`

**Goals**:
- Refactor package by package
- Maintain CI/CD stability
- Learn and document best practices
- Update team coding standards

**Actions**:
1. **Week 1-2**: Refactor shared packages
   - `packages/shared`
   - `packages/lib-foo`
   - Focus on pure functions and utility modules

2. **Week 3-4**: Refactor application code
   - `apps/api`
   - `apps/web`
   - `apps/cli`
   - Extract business logic to smaller functions

3. **Week 5-6**: Refactor test files
   - Update test structure
   - Extract test helpers
   - Consolidate test utilities

### Phase 3: Full Enforcement
**Duration**: Ongoing
**ESLint Mode**: `error`
**Environment**: `ESLINT_CLEAN_CODE_PHASE=3`

**Goals**:
- Enforce all rules in CI/CD
- Prevent new violations
- Maintain code quality standards

**Actions**:
1. Set `ESLINT_CLEAN_CODE_PHASE=3` in CI environment
2. Update pre-push hooks to use error mode
3. Monitor and address any remaining violations
4. Document lessons learned

## Refactoring Strategies

### For Large Functions (max-lines-per-function)
```typescript
// Before: 40+ line function
function processUserData(userData) {
  // validation logic (10 lines)
  // transformation logic (15 lines)
  // business logic (10 lines)
  // persistence logic (10 lines)
}

// After: Multiple focused functions
function validateUserData(userData) { /* 8 lines */ }
function transformUserData(userData) { /* 10 lines */ }
function applyBusinessRules(userData) { /* 8 lines */ }
function persistUserData(userData) { /* 8 lines */ }

function processUserData(userData) {
  const validated = validateUserData(userData);
  const transformed = transformUserData(validated);
  const processed = applyBusinessRules(transformed);
  return persistUserData(processed);
}
```

### For Large Files (max-lines)
- Extract related functions to separate modules
- Create barrel exports for cleaner imports
- Split classes into smaller, focused classes
- Use composition over inheritance

### For Complex Functions (complexity, max-depth)
- Use early returns/guard clauses
- Extract conditional logic to separate functions
- Replace nested if-else with switch or lookup tables
- Use functional programming patterns (map, filter, reduce)

### For Many Parameters (max-params)
```typescript
// Before: 5+ parameters
function createUser(name, email, age, role, department) { }

// After: Use options object
interface CreateUserOptions {
  name: string;
  email: string;
  age: number;
  role: string;
  department: string;
}

function createUser(options: CreateUserOptions) { }
```

## Common Patterns

### Extract Helper Functions
```typescript
// Before: One large function
function processOrder(order) {
  // 50+ lines of mixed logic
}

// After: Multiple helpers
function processOrder(order) {
  const validatedOrder = validateOrder(order);
  const pricedOrder = calculatePricing(validatedOrder);
  const taxedOrder = applyTaxes(pricedOrder);
  return submitOrder(taxedOrder);
}
```

### Use Composition
```typescript
// Compose small functions
const processOrder = pipe(
  validateOrder,
  calculatePricing,
  applyTaxes,
  submitOrder
);
```

### Extract Constants
```typescript
// Move magic numbers and strings to constants
const MAX_RETRY_ATTEMPTS = 3;
const DEFAULT_TIMEOUT_MS = 5000;
const API_BASE_URL = 'https://api.example.com';
```

## Monitoring Progress

### Weekly Metrics
Run weekly to track progress:
```bash
# Check current violations
pnpm lint:report

# Test with error mode
ESLINT_CLEAN_CODE_PHASE=3 pnpm lint

# Check specific package
pnpm --filter @ts-template/shared lint
```

### Success Criteria
- Phase 1: Complete assessment, migration plan approved
- Phase 2: <10% files with violations
- Phase 3: 0 violations in error mode

## Team Guidelines

### Code Review Checklist
- [ ] Functions under 20 lines?
- [ ] Files under 300 lines?
- [ ] Cyclomatic complexity under 10?
- [ ] Max 4 parameters per function?
- [ ] Max 3 levels of nesting?

### When to Grant Exceptions
Some files may legitimately need exceptions:
- Configuration files with many options
- Test files with extensive setup
- Generated code
- Third-party integrations

Document exceptions in code:
```typescript
/* eslint-disable max-lines */
// Reason: Configuration file requires all options in one place
```

## Resources

- [Clean Code by Robert C. Martin](https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882)
- [Refactoring by Martin Fowler](https://refactoring.com/)
- [ESLint Rules Documentation](https://eslint.org/docs/rules/)

## FAQ

**Q: What if a function legitimately needs more than 20 lines?**
A: First try to extract helper functions. If truly necessary, document why with a comment and consider disabling the rule for that specific function.

**Q: How do we handle legacy code?**
A: Add to technical debt backlog and refactor opportunistically when touching the code for features or bugs.

**Q: Should we refactor test files?**
A: Yes, but with more lenient limits (50 lines for test functions, 500 lines for test files).

**Q: What about generated code?**
A: Add generated files to `.eslintignore` or use `ignorePatterns` in ESLint config.