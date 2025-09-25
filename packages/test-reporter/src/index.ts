export interface TestReport {
  runId: string;
  branch: string;
  commitSha: string;
  timestamp: string;
  reportPath: string;
  status: 'success' | 'failure' | 'pending';
}

export interface WikiIndex {
  lastRun: TestReport | null;
  recentRuns: TestReport[];
  branches: Record<string, TestReport[]>;
}

export interface WikiConfig {
  maxReportsPerBranch: number;
  basePath: string;
  timezone?: string;
}

const DEFAULT_CONFIG: WikiConfig = {
  maxReportsPerBranch: 20,
  basePath: 'reports',
  timezone: 'UTC'
};

/**
 * Helper function to group reports by branch
 */
function groupReportsByBranch(reports: TestReport[]): Record<string, TestReport[]> {
  const branchGroups: Record<string, TestReport[]> = {};
  for (const report of reports) {
    if (!branchGroups[report.branch]) {
      branchGroups[report.branch] = [];
    }
    branchGroups[report.branch]?.push(report);
  }
  return branchGroups;
}

/**
 * Helper function to apply retention policy to branches
 */
function applyBranchRetention(
  branchGroups: Record<string, TestReport[]>,
  maxReports: number
): Record<string, TestReport[]> {
  const prunedBranches: Record<string, TestReport[]> = {};
  for (const [branch, branchReports] of Object.entries(branchGroups)) {
    prunedBranches[branch] = branchReports.slice(0, maxReports);
  }
  return prunedBranches;
}

/**
 * Pure function to generate Wiki index from test reports
 */
export function generateWikiIndex(reports: TestReport[], config: Partial<WikiConfig> = {}): WikiIndex {
  const finalConfig = { ...DEFAULT_CONFIG, ...config };

  const sortedReports = [...reports].sort((a, b) =>
    new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
  );

  const branchGroups = groupReportsByBranch(sortedReports);
  const prunedBranches = applyBranchRetention(branchGroups, finalConfig.maxReportsPerBranch);

  return {
    lastRun: sortedReports[0] || null,
    recentRuns: sortedReports.slice(0, finalConfig.maxReportsPerBranch),
    branches: prunedBranches
  };
}

/**
 * Generate markdown section for latest run
 */
function generateLatestRunSection(run: TestReport | null, timezone?: string): string[] {
  const lines: string[] = [];
  if (run) {
    lines.push(`## Latest Run (${run.runId})`, '');
    lines.push(`- **Branch:** ${run.branch}`);
    lines.push(`- **Run:** [${run.runId}](${run.reportPath})`);
    lines.push(`- **Commit:** ${run.commitSha}`);
    lines.push(`- **Time:** ${formatTimestamp(run.timestamp, timezone)}`);
    lines.push(`- **Status:** ${run.status}`, '');
  } else {
    lines.push('## No Reports Available', '');
    lines.push('No test reports have been published yet.', '');
  }
  return lines;
}

/**
 * Generate markdown section for recent runs
 */
function generateRecentRunsSection(runs: TestReport[]): string[] {
  if (runs.length === 0) return [];
  const lines: string[] = ['### Recent Runs', ''];
  for (const run of runs) {
    lines.push(`- ${run.branch} / ${run.runId} → [Report](${run.reportPath})`);
  }
  lines.push('');
  return lines;
}

/**
 * Generate markdown section for a single branch
 */
function generateBranchSection(branch: string, runs: TestReport[], timezone?: string): string[] {
  if (runs.length === 0) return [];
  const lines: string[] = [`### Branch: ${branch}`, ''];
  for (const run of runs) {
    const time = formatTimestamp(run.timestamp, timezone);
    lines.push(`- [${run.runId}](${run.reportPath}) - ${time} (${run.status})`);
  }
  lines.push('');
  return lines;
}

/**
 * Pure function to generate Markdown for Wiki home page
 */
export function generateWikiMarkdown(index: WikiIndex, config: Partial<WikiConfig> = {}): string {
  const finalConfig = { ...DEFAULT_CONFIG, ...config };
  const lines: string[] = ['# Test Reports', ''];

  lines.push(...generateLatestRunSection(index.lastRun, finalConfig.timezone));
  lines.push(...generateRecentRunsSection(index.recentRuns));

  for (const [branch, runs] of Object.entries(index.branches)) {
    lines.push(...generateBranchSection(branch, runs, finalConfig.timezone));
  }

  return lines.join('\n');
}

/**
 * Pure function to format timestamp
 */
export function formatTimestamp(timestamp: string, timezone = 'UTC'): string {
  const date = new Date(timestamp);
  const options: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    timeZone: timezone,
    timeZoneName: 'short'
  };

  return date.toLocaleString('en-US', options);
}

/**
 * Helper to group reports by branch using Map
 */
function groupReportsToMap(reports: TestReport[]): Map<string, TestReport[]> {
  const branches = new Map<string, TestReport[]>();
  for (const report of reports) {
    if (!branches.has(report.branch)) {
      branches.set(report.branch, []);
    }
    branches.get(report.branch)!.push(report);
  }
  return branches;
}

/**
 * Helper to categorize reports based on retention policy
 */
function categorizeReports(
  branches: Map<string, TestReport[]>,
  maxPerBranch: number
): { retain: TestReport[]; remove: TestReport[] } {
  const retain: TestReport[] = [];
  const remove: TestReport[] = [];

  for (const [, branchReports] of branches) {
    branchReports.forEach((report, index) => {
      (index < maxPerBranch ? retain : remove).push(report);
    });
  }

  return { retain, remove };
}

/**
 * Pure function to determine reports to retain
 */
export function getReportsToRetain(
  reports: TestReport[],
  maxPerBranch: number
): { retain: TestReport[]; remove: TestReport[] } {
  const sorted = [...reports].sort((a, b) =>
    new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
  );

  const branches = groupReportsToMap(sorted);
  return categorizeReports(branches, maxPerBranch);
}

/**
 * Pure function to parse report path
 */
export function parseReportPath(path: string): { branch: string; runId: string } | null {
  // Expected format: reports/{branch}/{runId}/
  const match = path.match(/reports\/([^/]+)\/([^/]+)\/?/);
  if (match && match[1] && match[2]) {
    return {
      branch: match[1],
      runId: match[2]
    };
  }
  return null;
}

/**
 * Pure function to build report path
 */
export function buildReportPath(branch: string, runId: string, basePath = 'reports'): string {
  return `${basePath}/${branch}/${runId}/index.html`;
}

// Export modules
export * from './markdown-formatter';
export * from './report-collector';
export * from './wiki-publisher';