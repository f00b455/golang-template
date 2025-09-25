# Debug Issues Playbook

## Common debugging steps for this monorepo:

### Build Issues:
1. Check TypeScript compilation: `pnpm type-check`
2. Verify dependencies are installed: `pnpm install`
3. Clean build cache: `pnpm clean && pnpm install`
4. Check workspace references in tsconfig files

### Test Failures:
1. Run specific package tests: `pnpm --filter @ts-template/PACKAGE test`
2. Check test database (always use test database!)
3. Verify mock implementations
4. Check test environment setup in vitest.config.ts

### API Issues:
1. Check Fastify server logs
2. Verify API routes in `apps/api/src/routes/`
3. Test endpoints with Swagger UI at `/documentation`
4. Check CORS configuration for frontend calls

### Frontend Issues:
1. Check Next.js dev server output
2. Verify component imports from shared package
3. Check browser console for runtime errors
4. Test with different browsers in Playwright

### Monorepo Issues:
1. Verify workspace configuration in `pnpm-workspace.yaml`
2. Check package dependencies in individual package.json files
3. Ensure Turborepo cache is not corrupted: `turbo clean`
4. Verify package names match workspace references

## Debugging Commands:
- `pnpm --filter PACKAGE dev` - Run specific package in dev mode
- `turbo build --dry-run` - Show build dependencies
- `pnpm why PACKAGE` - Show why package is installed