# Migration from TypeScript to Go

This document outlines the migration from the original TypeScript monorepo to the new Go project structure.

## Key Changes

### Project Structure

| TypeScript (Before) | Go (After) | Description |
|---------------------|------------|-------------|
| `apps/api/` | `cmd/api/` | HTTP API server |
| `apps/cli/` | `cmd/cli/` | CLI application |
| `packages/shared/` | `pkg/shared/` | Shared utilities |
| `packages/lib-foo/` | `pkg/core/` | Core business logic |
| N/A | `internal/` | Private application code |
| Various test dirs | `features/` | BDD tests centralized |

### Technology Stack

| Aspect | TypeScript | Go |
|--------|------------|-----|
| **HTTP Framework** | Fastify | Gin |
| **CLI Framework** | CAC | Cobra |
| **Testing** | Vitest | Go testing + Testify |
| **BDD** | Cucumber | Godog |
| **Build Tool** | Turborepo + pnpm | Make + go build |
| **Linting** | ESLint | golangci-lint |
| **Package Manager** | pnpm workspaces | Go modules |

### API Server Migration

#### Fastify → Gin

**Before (TypeScript/Fastify):**
```typescript
export const greetRoute: FastifyPluginAsync = async function (fastify) {
  fastify.get('/greet', {
    schema: { /* OpenAPI schema */ }
  }, async function (request) {
    const { name = 'World' } = request.query;
    return { message: greet(name) };
  });
};
```

**After (Go/Gin):**
```go
func (h *GreetHandler) Greet(c *gin.Context) {
    name := c.DefaultQuery("name", "World")
    message := shared.Greet(name)
    c.JSON(http.StatusOK, GreetResponse{Message: message})
}
```

#### Key Differences:
- **Handlers**: Struct methods instead of plugin functions
- **JSON responses**: Explicit struct types instead of inline objects
- **Swagger**: Annotations in comments instead of schema objects
- **Dependency injection**: Constructor pattern instead of plugin registration

### CLI Application Migration

#### CAC → Cobra

**Before (TypeScript/CAC):**
```typescript
cli
  .command('', 'Generate a colorful hello message')
  .option('--name <name>', 'Name to greet', { default: 'World' })
  .action(async (options) => {
    // Command logic
  });
```

**After (Go/Cobra):**
```go
var rootCmd = &cobra.Command{
    Use:   "hello-cli",
    Short: "A colorful Hello-World CLI application",
    Run:   runHelloCommand,
}

func init() {
    rootCmd.Flags().StringVar(&name, "name", "World", "Name to greet")
}
```

#### Key Differences:
- **Command structure**: Global command variables instead of chaining
- **Options**: Struct fields instead of inline options
- **Action handlers**: Separate functions instead of inline callbacks

### Testing Migration

#### Vitest → Go Testing + Testify

**Before (TypeScript/Vitest):**
```typescript
describe('greet function', () => {
  it('should greet with valid name', () => {
    expect(greet('World')).toBe('Hello, World!');
  });
});
```

**After (Go/Testify):**
```go
func TestGreet(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"valid name", "World", "Hello, World!"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := shared.Greet(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

#### Key Differences:
- **Test structure**: Table-driven tests pattern
- **Assertions**: testify/assert instead of expect/toBe
- **Test organization**: One file per package instead of describe blocks

### BDD Migration

#### Cucumber → Godog

**Feature files remain largely the same**, but step definitions changed:

**Before (TypeScript):**
```typescript
Given('I have the name {string}', function (name: string) {
  this.name = name;
});
```

**After (Go):**
```go
func (ctx *greetFeatureContext) iHaveTheName(name string) error {
    ctx.name = name
    return nil
}
```

### Build System Migration

#### pnpm/Turborepo → Make/Go

**Before (package.json):**
```json
{
  "scripts": {
    "build": "turbo build",
    "dev": "turbo dev",
    "test": "turbo test"
  }
}
```

**After (Makefile):**
```makefile
build: deps
    $(GOBUILD) -o ./bin/api-server $(API_DIR)
    $(GOBUILD) -o ./bin/cli-tool $(CLI_DIR)

dev: deps
    $(GOCMD) run $(API_DIR)/main.go

test: deps
    $(GOTEST) -v -race $(SRC_DIR)
```

## Migration Benefits

### Performance
- **Faster startup**: Go binaries start instantly vs Node.js startup time
- **Lower memory usage**: Go has smaller memory footprint
- **Better concurrency**: Native goroutines vs event loop

### Deployment
- **Single binary**: No need for Node.js runtime or node_modules
- **Static linking**: Self-contained executables
- **Cross-compilation**: Build for different platforms easily

### Development Experience
- **Type safety**: Compile-time type checking
- **Better tooling**: Rich ecosystem of Go tools
- **Explicit error handling**: Less runtime surprises

## Migration Checklist

- [x] Set up Go module structure
- [x] Convert shared utilities (`pkg/shared/`)
- [x] Convert core business logic (`pkg/core/`)
- [x] Migrate API server to Gin
- [x] Migrate CLI application to Cobra
- [x] Convert unit tests to Go testing
- [x] Convert BDD tests to Godog
- [x] Set up CI/CD pipeline
- [x] Update documentation

## Potential Challenges

1. **Learning Curve**: Team needs to learn Go idioms and patterns
2. **Library Ecosystem**: Some TypeScript libraries may not have Go equivalents
3. **JSON Handling**: More verbose than TypeScript's dynamic typing
4. **Error Handling**: Explicit error handling requires more code
5. **Frontend Integration**: Need to maintain API compatibility

## Next Steps

1. **Performance Testing**: Compare performance with TypeScript version
2. **Gradual Migration**: Run both versions in parallel during transition
3. **Team Training**: Provide Go training for development team
4. **Documentation**: Update all development workflows and guides
5. **Monitoring**: Set up monitoring for the new Go services