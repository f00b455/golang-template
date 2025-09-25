# @ts-template/test-reporter

Test reporter module for publishing test results to GitHub Wiki.

## Features

- Pure functions for generating Wiki indexes
- Retention policy management (keep last N reports per branch)
- Markdown generation for Wiki home page
- Support for multiple branches
- Timestamp formatting with timezone support

## Installation

```bash
pnpm add @ts-template/test-reporter
```

## Usage

```typescript
import {
  generateWikiIndex,
  generateWikiMarkdown,
  getReportsToRetain,
  type TestReport
} from '@ts-template/test-reporter';

// Collect test reports
const reports: TestReport[] = [
  {
    runId: 'run-123',
    branch: 'main',
    commitSha: 'abc123def',
    timestamp: new Date().toISOString(),
    reportPath: 'reports/main/run-123/index.html',
    status: 'success'
  }
];

// Generate Wiki index
const index = generateWikiIndex(reports, {
  maxReportsPerBranch: 20,
  basePath: 'reports'
});

// Generate Markdown for Wiki home page
const markdown = generateWikiMarkdown(index, {
  timezone: 'Europe/Berlin'
});

// Determine which reports to keep
const { retain, remove } = getReportsToRetain(reports, 20);
```

## API

### `generateWikiIndex(reports, config)`

Generates a Wiki index from test reports.

### `generateWikiMarkdown(index, config)`

Generates Markdown content for the Wiki home page.

### `getReportsToRetain(reports, maxPerBranch)`

Determines which reports to retain based on retention policy.

### `formatTimestamp(timestamp, timezone)`

Formats a timestamp with timezone support.

### `parseReportPath(path)`

Parses a report path to extract branch and run ID.

### `buildReportPath(branch, runId, basePath)`

Builds a report path from components.

## Testing

```bash
# Run unit tests
pnpm test

# Run BDD tests
pnpm test:cucumber

# Watch mode
pnpm test:watch
```

## GitHub Actions Integration

See `.github/issue-tracking/github-actions-wiki-publish.yml` for example workflow configuration.