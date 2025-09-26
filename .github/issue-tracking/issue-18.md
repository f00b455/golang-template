# Issue #18: Fix Swagger documentation build error for CSV export endpoint

**Issue URL**: https://github.com/f00b455/golang-template/issues/18
**Created**: 2025-09-26T06:25:41Z
**Assignee**: Unassigned

## Description
## Problem
The build is failing with a Swagger documentation generation error after merging PR #12 (RSS Export feature).

## Error Message
```
2025/09/26 08:24:28 ParseComment error in file /Users/misha/dev/space/claude-experiments/golang-template/internal/handlers/rss.go for comment: '// @Produce      json,csv': csv produce type can't be accepted
make: *** [docs] Error 1
```

## Root Cause
The Swagger annotation in `internal/handlers/rss.go` line 379 uses:
```go
// @Produce      json,csv
```

However, `csv` is not a valid MIME type for Swagger. It should be `text/csv` or the annotation should be split into separate lines.

## Suggested Fix
Replace the annotation with:
```go
// @Produce      json
// @Produce      text/csv
```

## Impact
- Build fails when running `make build`
- Swagger documentation cannot be generated
- Application cannot be compiled and started

## Priority
High - This is blocking the build process

## Work Log
- Branch created: issue-18-fix-swagger-documentation-build-error-for-csv-expo
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
