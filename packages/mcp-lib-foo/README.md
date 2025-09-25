# MCP Server for lib-foo

MCP (Model Context Protocol) Server that exposes lib-foo functionality for use with Claude Code and other MCP-compatible clients.

## Features

This MCP server exposes the following tools from the lib-foo package:

- **fooProcess**: Process text with prefix and optional suffix
- **fooGreet**: Create greetings with foo processing
- **createFooProcessor**: Create a configured foo processor
- **fooTransform**: Transform arrays of strings with custom transformers
- **fooFilter**: Filter arrays based on predicates

## Installation

### Local Development

```bash
# Install dependencies
pnpm install

# Build the package
pnpm build

# Run tests
pnpm test
pnpm test:cucumber

# Start the server
pnpm start
```

### Docker Deployment

#### Build and Run with Docker

```bash
# Build the Docker image
docker build -t mcp-lib-foo:latest -f packages/mcp-lib-foo/Dockerfile .

# Run the container
docker run -it --rm \
  -e NODE_ENV=production \
  -e MCP_SERVER_NAME=mcp-lib-foo \
  -p 3010:3000 \
  mcp-lib-foo:latest
```

#### Using Docker Compose

```bash
# Start the server
cd packages/mcp-lib-foo
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the server
docker-compose down

# Run tests in Docker
docker-compose --profile test up mcp-lib-foo-test
```

## Configuration

### Environment Variables

- `NODE_ENV`: Environment mode (development/production/test)
- `MCP_SERVER_NAME`: Server identifier (default: mcp-lib-foo)
- `MCP_SERVER_VERSION`: Server version (default: 0.0.0)
- `LOG_LEVEL`: Logging level (info/debug/error)

### Claude Code Integration

To use this MCP server with Claude Code, add it to your MCP configuration:

```json
{
  "mcpServers": {
    "mcp-lib-foo": {
      "command": "node",
      "args": ["/path/to/mcp-lib-foo/dist/index.js"],
      "env": {
        "NODE_ENV": "production"
      }
    }
  }
}
```

Or with Docker:

```json
{
  "mcpServers": {
    "mcp-lib-foo": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "mcp-lib-foo:latest"],
      "env": {
        "NODE_ENV": "production"
      }
    }
  }
}
```

## API Documentation

### Available Tools

#### fooProcess

Process text with prefix and optional suffix.

**Parameters:**
- `input` (string, required): The input text to process
- `prefix` (string, required): Prefix to add to the input
- `suffix` (string, optional): Suffix to add to the input

**Example:**
```json
{
  "tool": "fooProcess",
  "arguments": {
    "input": "test",
    "prefix": "pre-",
    "suffix": "-post"
  }
}
```

#### fooGreet

Create greetings with foo processing.

**Parameters:**
- `name` (string, required): The name to greet
- `prefix` (string, required): Prefix for the greeting
- `suffix` (string, optional): Suffix for the greeting

**Example:**
```json
{
  "tool": "fooGreet",
  "arguments": {
    "name": "Alice",
    "prefix": "***",
    "suffix": "!!!"
  }
}
```

#### createFooProcessor

Create a configured foo processor.

**Parameters:**
- `prefix` (string, required): Default prefix for the processor
- `suffix` (string, optional): Default suffix for the processor
- `action` (string, required): Action to perform ("process" | "greetWithFoo")
- `input` (string, required): Input for the selected action

**Example:**
```json
{
  "tool": "createFooProcessor",
  "arguments": {
    "prefix": "[",
    "suffix": "]",
    "action": "process",
    "input": "data"
  }
}
```

#### fooTransform

Transform an array of strings.

**Parameters:**
- `data` (array of strings, required): Array to transform
- `transformType` (string, required): Type of transformation ("uppercase" | "lowercase" | "reverse" | "trim")

**Example:**
```json
{
  "tool": "fooTransform",
  "arguments": {
    "data": ["hello", "world"],
    "transformType": "uppercase"
  }
}
```

#### fooFilter

Filter an array based on a predicate.

**Parameters:**
- `items` (array, required): Array of items to filter
- `filterType` (string, required): Type of filter ("truthy" | "falsy" | "nonEmpty" | "numeric" | "string")
- `customPredicate` (string, optional): Custom predicate for string filtering (e.g., "minLength:5")

**Example:**
```json
{
  "tool": "fooFilter",
  "arguments": {
    "items": [1, "test", 0, "", "valid", null],
    "filterType": "truthy"
  }
}
```

## Development

### Project Structure

```
packages/mcp-lib-foo/
├── src/
│   ├── index.ts          # Main MCP server implementation
│   └── index.test.ts     # Unit tests
├── features/             # BDD feature files
│   ├── mcp-tools.feature
│   ├── mcp-server.feature
│   └── step_definitions/
│       └── mcp.steps.ts
├── Dockerfile            # Multi-stage Docker build
├── docker-compose.yml    # Local development setup
├── package.json
├── tsconfig.json
├── tsup.config.ts
├── vitest.config.ts
└── cucumber.cjs
```

### Testing

The package includes comprehensive testing:

- **Unit Tests**: Run with `pnpm test`
- **BDD Tests**: Run with `pnpm test:cucumber`
- **Docker Tests**: Run with `docker-compose --profile test up`

### Building

```bash
# TypeScript build
pnpm build

# Docker build
docker build -t mcp-lib-foo:latest -f Dockerfile ../..
```

## Troubleshooting

### Server Won't Start

1. Check that all dependencies are installed: `pnpm install`
2. Ensure the package is built: `pnpm build`
3. Verify environment variables are set correctly
4. Check Docker logs if using containerized deployment

### Connection Issues with Claude Code

1. Verify the server is running: `docker ps` or check process list
2. Check the MCP configuration in Claude Code settings
3. Ensure the correct path to the server executable
4. Review server logs for error messages

### Docker Build Failures

1. Ensure you're building from the repository root
2. Check that all workspace dependencies are present
3. Verify Docker daemon is running
4. Clear Docker cache: `docker builder prune`

## License

Private package - see repository license.