import * as fs from 'fs/promises';
import * as path from 'path';

export interface CucumberReport {
  features: Feature[];
  stats: {
    passed: number;
    failed: number;
    skipped: number;
    total: number;
    duration: number;
  };
}

export interface Feature {
  name: string;
  description: string;
  keyword: string;
  tags: Array<{ name: string }>;
  scenarios: Scenario[];
}

export interface Scenario {
  name: string;
  keyword: string;
  tags: Array<{ name: string }>;
  steps: Step[];
  status: 'passed' | 'failed' | 'skipped';
  duration?: number;
}

export interface Step {
  name: string;
  keyword: string;
  result?: {
    status: 'passed' | 'failed' | 'skipped' | 'pending';
    duration?: number;
    error_message?: string;
  };
}

export interface MarkdownOptions {
  runId: string;
  branch: string;
  commitSha: string;
  timestamp: string;
  basePath?: string;
  includeStats?: boolean;
  includeErrors?: boolean;
}

/**
 * Format duration from nanoseconds to readable string
 */
export function formatDuration(nanoSeconds?: number): string {
  if (!nanoSeconds) return '0ms';
  const ms = Math.round(nanoSeconds / 1000000);
  if (ms < 1500) return `${ms}ms`; // Show ms for anything under 1.5 seconds
  const seconds = (ms / 1000).toFixed(2);
  return `${seconds}s`;
}

/**
 * Calculate stats from features
 */
export function calculateStats(features: Feature[]): CucumberReport['stats'] {
  let passed = 0;
  let failed = 0;
  let skipped = 0;
  let duration = 0;

  for (const feature of features) {
    for (const scenario of feature.scenarios) {
      if (scenario.status === 'passed') passed++;
      else if (scenario.status === 'failed') failed++;
      else skipped++;
      duration += scenario.duration || 0;
    }
  }

  return {
    passed,
    failed,
    skipped,
    total: passed + failed + skipped,
    duration
  };
}

/**
 * Generate stats section
 */
function generateStatsSection(stats: CucumberReport['stats']): string[] {
  const lines: string[] = ['## Test Statistics', ''];
  const passRate = stats.total > 0
    ? Math.round((stats.passed / stats.total) * 100)
    : 0;

  lines.push(`- **Total Scenarios:** ${stats.total}`);
  lines.push(`- **Passed:** ${stats.passed} ✅`);
  lines.push(`- **Failed:** ${stats.failed} ❌`);
  lines.push(`- **Skipped:** ${stats.skipped} ⏭️`);
  lines.push(`- **Pass Rate:** ${passRate}%`);
  lines.push(`- **Duration:** ${formatDuration(stats.duration)}`);
  lines.push('');

  return lines;
}

/**
 * Generate feature section
 */
function generateFeatureSection(
  feature: Feature,
  includeErrors: boolean
): string[] {
  const lines: string[] = [`### ${feature.keyword}: ${feature.name}`, ''];

  if (feature.description) {
    lines.push(feature.description, '');
  }

  if (feature.tags.length > 0) {
    const tags = feature.tags.map(t => t.name).join(' ');
    lines.push(`**Tags:** ${tags}`, '');
  }

  for (const scenario of feature.scenarios) {
    lines.push(...generateScenarioSection(scenario, includeErrors));
  }

  return lines;
}

/**
 * Generate scenario section
 */
function generateScenarioSection(
  scenario: Scenario,
  includeErrors: boolean
): string[] {
  const statusIcon =
    scenario.status === 'passed' ? '✅' :
    scenario.status === 'failed' ? '❌' : '⏭️';

  const lines: string[] = [
    `#### ${scenario.keyword}: ${scenario.name} ${statusIcon}`,
    ''
  ];

  if (scenario.duration) {
    lines.push(`Duration: ${formatDuration(scenario.duration)}`, '');
  }

  const failedSteps = scenario.steps.filter(
    s => s.result?.status === 'failed'
  );

  if (failedSteps.length > 0 && includeErrors) {
    lines.push('**Failed Steps:**', '');
    for (const step of failedSteps) {
      lines.push(`- ${step.keyword} ${step.name}`);
      if (step.result?.error_message) {
        lines.push('  ```');
        lines.push(`  ${step.result.error_message}`);
        lines.push('  ```');
      }
    }
    lines.push('');
  }

  return lines;
}

/**
 * Generate summary section
 */
function generateSummarySection(
  options: MarkdownOptions,
  stats: CucumberReport['stats']
): string[] {
  const status = stats.failed === 0 ? 'SUCCESS ✅' : 'FAILED ❌';

  return [
    `# Test Report - ${options.runId}`,
    '',
    `## Summary`,
    '',
    `- Status: ${status}`,
    `- Branch: ${options.branch}`,
    `- Commit: ${options.commitSha}`,
    `- Timestamp: ${options.timestamp}`,
    `- Duration: ${formatDuration(stats.duration)}`,
    ''
  ];
}

/**
 * Convert Cucumber JSON report to markdown
 */
export function cucumberJsonToMarkdown(
  jsonReport: unknown[],
  options: MarkdownOptions
): string {
  const features: Feature[] = jsonReport.map(parseFeature);
  const stats = calculateStats(features);

  const lines: string[] = [];

  lines.push(...generateSummarySection(options, stats));

  if (options.includeStats !== false) {
    lines.push(...generateStatsSection(stats));
  }

  lines.push('## Features', '');

  for (const feature of features) {
    lines.push(...generateFeatureSection(
      feature,
      options.includeErrors !== false
    ));
  }

  lines.push('---');
  lines.push(`*Generated on ${new Date().toISOString()}*`);

  return lines.join('\n');
}

/**
 * Parse feature from Cucumber JSON
 */
function parseFeature(featureData: unknown): Feature {
  const data = featureData as Record<string, unknown>;
  const scenarios: Scenario[] = [];

  const elements = (data.elements as unknown[]) || [];
  for (const element of elements) {
    const el = element as Record<string, unknown>;
    if (el.type === 'scenario') {
      scenarios.push(parseScenario(element));
    }
  }

  return {
    name: (data.name as string) || 'Unnamed Feature',
    description: (data.description as string) || '',
    keyword: (data.keyword as string) || 'Feature',
    tags: (data.tags as Array<{ name: string }>) || [],
    scenarios
  };
}

/**
 * Parse scenario from Cucumber JSON
 */
function parseScenario(scenarioData: unknown): Scenario {
  const data = scenarioData as Record<string, unknown>;
  const steps: Step[] = ((data.steps as unknown[]) || []).map(parseStep);
  const status = getScenarioStatus(steps);
  const duration = calculateStepsDuration(steps);

  return {
    name: (data.name as string) || 'Unnamed Scenario',
    keyword: (data.keyword as string) || 'Scenario',
    tags: (data.tags as Array<{ name: string }>) || [],
    steps,
    status,
    duration
  };
}

function getScenarioStatus(steps: Step[]): 'passed' | 'failed' | 'skipped' {
  const hasFailedStep = steps.some(s => s.result?.status === 'failed');
  const hasSkippedStep = steps.some(s => s.result?.status === 'skipped');
  return hasFailedStep ? 'failed' : hasSkippedStep ? 'skipped' : 'passed';
}

function calculateStepsDuration(steps: Step[]): number {
  return steps.reduce((sum, step) => sum + (step.result?.duration || 0), 0);
}

/**
 * Parse step from Cucumber JSON
 */
function parseStep(stepData: unknown): Step {
  const data = stepData as Record<string, unknown>;
  return {
    name: (data.name as string) || '',
    keyword: (data.keyword as string) || '',
    result: (data.result as Step['result']) || {
      status: 'pending',
      duration: 0
    }
  };
}

/**
 * Generate index page for all reports
 */
export function generateIndexPage(
  reports: ReportEntry[],
  options?: IndexOptions
): string {
  const lines = generateHeader(options?.title);
  const groupedByBranch = groupReportsByBranch(reports);

  for (const [branch, branchReports] of groupedByBranch) {
    lines.push(...generateBranchSection(branch, branchReports, options));
  }

  lines.push(...generateFooter());
  return lines.join('\n');
}

interface ReportEntry {
  runId: string;
  branch: string;
  timestamp: string;
  path: string;
  status: string;
}

interface IndexOptions {
  title?: string;
  maxReports?: number;
}

function generateHeader(title?: string): string[] {
  return [
    title || '# Test Reports Index',
    '',
    'Automated test reports from CI/CD pipeline.',
    ''
  ];
}

function groupReportsByBranch(reports: ReportEntry[]): Map<string, ReportEntry[]> {
  const grouped = new Map<string, ReportEntry[]>();
  for (const report of reports) {
    if (!grouped.has(report.branch)) {
      grouped.set(report.branch, []);
    }
    grouped.get(report.branch)!.push(report);
  }
  return grouped;
}

function generateBranchSection(
  branch: string,
  branchReports: ReportEntry[],
  options?: IndexOptions
): string[] {
  const lines: string[] = [`## Branch: ${branch}`, ''];
  const sortedReports = sortAndLimitReports(branchReports, options?.maxReports);

  for (const report of sortedReports) {
    const statusIcon = report.status === 'success' ? '✅' : '❌';
    lines.push(
      `- [${report.runId}](${report.path}) - ${report.timestamp} ${statusIcon}`
    );
  }
  lines.push('');
  return lines;
}

function sortAndLimitReports(reports: ReportEntry[], maxReports = 20): ReportEntry[] {
  return reports
    .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
    .slice(0, maxReports);
}

function generateFooter(): string[] {
  return [
    '---',
    `*Last updated: ${new Date().toISOString()}*`
  ];
}

/**
 * Read and parse Cucumber JSON report
 */
export async function readCucumberReport(
  filePath: string
): Promise<unknown[]> {
  const content = await fs.readFile(filePath, 'utf-8');
  return JSON.parse(content);
}

/**
 * Write markdown report to file
 */
export async function writeMarkdownReport(
  content: string,
  outputPath: string
): Promise<void> {
  await fs.mkdir(path.dirname(outputPath), { recursive: true });
  await fs.writeFile(outputPath, content, 'utf-8');
}