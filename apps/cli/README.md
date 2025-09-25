# Hello-World CLI

A colorful command-line interface application that generates personalized greeting messages with visual effects.

## Features

- Personalized greetings with custom names
- 5-second animated spinner with "Ready!" completion message
- 3-second progress bar from 0% to 100%
- Colorful gradient text output
- Decorative box framing with title
- Integration with `@ts-template/lib-foo` for greeting generation

## Installation

```bash
pnpm install
pnpm build
```

## Usage

```bash
# Default greeting (uses "World")
pnpm start

# Custom name greeting
pnpm start --name "Alice"

# Direct execution after build
node dist/index.js --name "Bob"
```

## Development

```bash
# Run in development mode with watch
pnpm dev

# Run tests
pnpm test

# Run BDD tests
pnpm test:cucumber

# Lint code
pnpm lint

# Type checking
pnpm type-check
```

## Architecture

This CLI application:
- Uses `cac` for command-line argument parsing
- Integrates with `@ts-template/lib-foo` for the core greeting logic
- Implements visual effects with `nanospinner` and `cli-progress`
- Creates colorful output using `picocolors` and `gradient-string`
- Frames output in decorative boxes using `boxen`

## Testing

The application includes:
- Unit tests using Vitest
- BDD tests using Cucumber with feature files
- Full integration testing with lib-foo package

## Issue Reference

This CLI was implemented as part of Issue #14 to create a colorful Hello-World CLI application.