# Git Hooks Configuration

This directory contains git hooks managed by Husky.

## Pre-push Hook

The `pre-push` hook automatically runs the full pipeline validation before each push:

- ✅ **Lint**: Code style and ESLint checks
- ✅ **Type-check**: TypeScript type validation
- ✅ **Build**: Compilation verification
- ✅ **Tests**: Unit tests with proper mocking
- ✅ **BDD**: Cucumber behavior-driven development tests

## Development Workflow

1. **During development**: Use `pnpm validate:quick` for fast feedback
2. **Before pushing**: `git push` automatically triggers validation
3. **If validation fails**: Fix issues and push again

## Manual Validation

You can also run validation manually:

```bash
# Quick validation (lint + type-check + tests)
pnpm validate:quick

# Full validation (same as pre-push hook)
pnpm validate

# Test pre-push hook manually
pnpm prepush
```

This ensures code quality while keeping the development cycle efficient!