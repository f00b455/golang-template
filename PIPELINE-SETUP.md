# Pipeline Validation Setup Guide

This guide shows multiple ways to ensure your pipeline runs successfully after every change.

## Option 1: Manual Pre-commit Script (Simple)

Create a simple script to run before commits:

```bash
# Create the script
cat > validate-pipeline.sh << 'EOF'
#!/bin/bash
set -e

echo "🔍 Running pipeline validation..."

echo "📝 Linting..."
pnpm lint

echo "🔨 Building..."
pnpm build

echo "🧪 Running tests..."
pnpm test:run

echo "🥒 Running BDD tests..."
pnpm test:cucumber

echo "✅ All checks passed!"
EOF

chmod +x validate-pipeline.sh

# Run before each commit:
./validate-pipeline.sh && git commit -m "your message"
```

## Option 2: Husky + lint-staged (Recommended)

Install and configure Husky for automatic pre-commit hooks:

```bash
# Install Husky
pnpm add -D husky lint-staged

# Initialize Husky
pnpm exec husky init

# Add pre-commit hook
echo 'pnpm lint-staged' > .husky/pre-commit

# Configure lint-staged in package.json
```

Add to your `package.json`:

```json
{
  "lint-staged": {
    "*.{ts,tsx}": [
      "eslint --fix",
      "prettier --write"
    ],
    "packages/**/*.{ts,tsx}": [
      "bash -c 'pnpm lint'",
      "bash -c 'pnpm build'",
      "bash -c 'pnpm test:run'"
    ]
  }
}
```

## Option 3: GitHub Actions CI/CD (Automated)

Create `.github/workflows/ci.yml`:

```yaml
name: CI Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  validate:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - uses: pnpm/action-setup@v2
      with:
        version: 8

    - uses: actions/setup-node@v3
      with:
        node-version: '20'
        cache: 'pnpm'

    - name: Install dependencies
      run: pnpm install --frozen-lockfile

    - name: Lint
      run: pnpm lint

    - name: Build
      run: pnpm build

    - name: Test
      run: pnpm test:run

    - name: BDD Tests
      run: pnpm test:cucumber
```

## Option 4: VS Code Task Runner

Create `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Validate Pipeline",
      "type": "shell",
      "command": "pnpm lint && pnpm build && pnpm test:run && pnpm test:cucumber",
      "group": {
        "kind": "test",
        "isDefault": true
      },
      "presentation": {
        "reveal": "always",
        "panel": "new"
      },
      "problemMatcher": "$tsc"
    }
  ]
}
```

Then use `Cmd+Shift+P` → "Tasks: Run Test Task" or bind it to a keyboard shortcut.

## Option 5: Package.json Scripts

Add convenient scripts to your root `package.json`:

```json
{
  "scripts": {
    "validate": "pnpm lint && pnpm build && pnpm test:run && pnpm test:cucumber",
    "validate:quick": "pnpm lint && pnpm test:run",
    "precommit": "pnpm validate",
    "prepush": "pnpm validate"
  }
}
```

Then run:
```bash
pnpm validate  # Before committing
pnpm validate:quick  # For quick checks during development
```

## Option 6: Git Aliases

Add to your `~/.gitconfig`:

```ini
[alias]
    # Commit with validation
    cv = !pnpm validate && git commit
    # Add and commit with validation
    acv = !git add -A && pnpm validate && git commit
```

Usage:
```bash
git cv -m "your commit message"
git acv -m "add all and commit with validation"
```

## Option 7: Turbo Watch Mode (Development)

For continuous validation during development:

```bash
# Add to package.json
{
  "scripts": {
    "dev:validate": "turbo run lint build test --watch"
  }
}
```

## Recommended Setup

For the best developer experience, combine:

1. **GitHub Actions** - Catches issues in CI/CD
2. **Husky + lint-staged** - Prevents bad commits locally
3. **VS Code Tasks** - Quick validation during development
4. **Package.json scripts** - Convenient manual validation

## Quick Start (Minimal Setup)

```bash
# 1. Add validation script to package.json
pnpm pkg set scripts.validate="pnpm lint && pnpm build && pnpm test:run && pnpm test:cucumber"

# 2. Create simple Git alias
git config --local alias.validate '!pnpm validate && git'

# 3. Use it
git validate commit -m "your message"
```

## Tips

- Use `pnpm validate:quick` during development for faster feedback
- Configure your IDE to run linting on save
- Use `turbo --filter` to validate only changed packages
- Set up branch protection rules on GitHub to require CI passing

## CLAUDE.md Integration

Add this to your CLAUDE.md for AI assistance:

```markdown
### Pre-Commit Workflow:

**IMPORTANT**: Before committing any changes, ALWAYS run:
1. `pnpm lint` - Check code style
2. `pnpm build` - Ensure compilation
3. `pnpm test:run` - Verify unit tests
4. `pnpm test:cucumber` - Validate BDD tests

Or simply run: `pnpm validate`
```