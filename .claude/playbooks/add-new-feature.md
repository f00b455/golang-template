# Add New Feature Playbook

## Steps to add a new feature to the monorepo:

1. **Create shared types/utilities** (if needed):
   - Add to `packages/shared/src/index.ts`
   - Write tests in `packages/shared/__tests__/`
   - Run `pnpm --filter @ts-template/shared test`

2. **Add API endpoints**:
   - Create route in `apps/api/src/routes/`
   - Add Swagger documentation
   - Write unit tests in `apps/api/__tests__/`
   - Add Cucumber feature in `apps/api/features/`
   - Run `pnpm --filter @ts-template/api test`

3. **Update frontend**:
   - Add components in `apps/web/src/components/`
   - Update pages in `apps/web/src/app/`
   - Write unit tests in `apps/web/__tests__/`
   - Add E2E tests in `apps/web/e2e/`
   - Run `pnpm --filter @ts-template/web test`

4. **Integration testing**:
   - Run full test suite: `pnpm test`
   - Run E2E tests: `pnpm test:e2e`
   - Run Cucumber tests: `pnpm test:cucumber`

5. **Create changeset**:
   - Run `pnpm changeset`
   - Describe the changes and select affected packages
   - Commit the changeset file

## Testing Commands:
- `pnpm test` - Run all unit tests
- `pnpm test:e2e` - Run Playwright tests
- `pnpm test:cucumber` - Run BDD tests
- `pnpm lint` - Check code style
- `pnpm type-check` - Verify TypeScript