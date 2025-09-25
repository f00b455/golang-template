import * as fs from 'fs/promises';
import * as path from 'path';
import { glob } from 'glob';

export interface CollectorOptions {
  rootDir: string;
  patterns?: string[];
  excludePatterns?: string[];
  recursive?: boolean;
}

export interface CollectedReport {
  type: 'cucumber' | 'vitest' | 'playwright' | 'unknown';
  format: 'json' | 'html' | 'xml';
  path: string;
  package: string;
  size: number;
}

export interface PackageReports {
  packageName: string;
  packagePath: string;
  reports: CollectedReport[];
}

const DEFAULT_PATTERNS = [
  '**/cucumber-report.json',
  '**/cucumber-report.html',
  '**/vitest-report.json',
  '**/vitest-report.html',
  '**/playwright-report/index.html',
  '**/coverage/lcov-report/index.html'
];

const DEFAULT_EXCLUDE = [
  '**/node_modules/**',
  '**/dist/**',
  '**/.git/**'
];

/**
 * Detect report type from file path and content
 */
async function detectReportType(
  filePath: string
): Promise<CollectedReport['type']> {
  const basename = path.basename(filePath);

  if (basename.includes('cucumber')) return 'cucumber';
  if (basename.includes('vitest')) return 'vitest';
  if (basename.includes('playwright')) return 'playwright';

  try {
    const content = await fs.readFile(filePath, 'utf-8');
    const sample = content.slice(0, 1000);

    if (sample.includes('"cucumber"') || sample.includes('"feature"')) {
      return 'cucumber';
    }
    if (sample.includes('"vitest"') || sample.includes('"testResults"')) {
      return 'vitest';
    }
    if (sample.includes('playwright')) {
      return 'playwright';
    }
  } catch {
    // If we can't read the file, leave as unknown
  }

  return 'unknown';
}

/**
 * Detect report format from file extension
 */
function detectReportFormat(
  filePath: string
): CollectedReport['format'] {
  const ext = path.extname(filePath).toLowerCase();

  switch (ext) {
    case '.json':
      return 'json';
    case '.html':
      return 'html';
    case '.xml':
      return 'xml';
    default:
      return 'html';
  }
}

/**
 * Find package name from path
 */
function findPackageName(
  filePath: string,
  rootDir: string
): string {
  const relativePath = path.relative(rootDir, filePath);
  const parts = relativePath.split(path.sep);

  // Look for packages/ or apps/ directory
  const packagesIndex = parts.findIndex(p => p === 'packages' || p === 'apps');
  if (packagesIndex >= 0 && parts[packagesIndex + 1]) {
    return parts[packagesIndex + 1] as string;
  }

  // Fallback to first directory
  return parts[0] ?? 'root';
}

/**
 * Get file stats safely
 */
async function getFileStats(
  filePath: string
): Promise<{ size: number; exists: boolean }> {
  try {
    const stats = await fs.stat(filePath);
    return { size: stats.size, exists: true };
  } catch {
    return { size: 0, exists: false };
  }
}

/**
 * Collect test reports from file system
 */
export async function collectReports(
  options: CollectorOptions
): Promise<PackageReports[]> {
  const patterns = options.patterns || DEFAULT_PATTERNS;
  const excludePatterns = options.excludePatterns || DEFAULT_EXCLUDE;

  const allFiles: string[] = [];

  for (const pattern of patterns) {
    const matches = await glob(pattern, {
      cwd: options.rootDir,
      absolute: true,
      ignore: excludePatterns,
      nodir: true
    });
    if (matches) allFiles.push(...matches);
  }

  const uniqueFiles = Array.from(new Set(allFiles));
  const reports: CollectedReport[] = [];

  for (const file of uniqueFiles) {
    const { size, exists } = await getFileStats(file);
    if (!exists) continue;

    const type = await detectReportType(file);
    const format = detectReportFormat(file);
    const packageName = findPackageName(file, options.rootDir);

    reports.push({
      type,
      format,
      path: file,
      package: packageName,
      size
    });
  }

  return groupReportsByPackage(reports);
}

/**
 * Group reports by package
 */
function groupReportsByPackage(
  reports: CollectedReport[]
): PackageReports[] {
  const grouped = new Map<string, PackageReports>();

  for (const report of reports) {
    if (!grouped.has(report.package)) {
      grouped.set(report.package, {
        packageName: report.package,
        packagePath: path.dirname(report.path),
        reports: []
      });
    }
    grouped.get(report.package)!.reports.push(report);
  }

  return Array.from(grouped.values())
    .sort((a, b) => a.packageName.localeCompare(b.packageName));
}

/**
 * Find reports in specific directory
 */
export async function findReportsInDirectory(
  directory: string,
  patterns?: string[]
): Promise<CollectedReport[]> {
  const results = await collectReports({
    rootDir: directory,
    ...(patterns && { patterns }),
    recursive: true
  });

  return results.flatMap(pkg => pkg.reports);
}

/**
 * Copy report files to destination
 */
export async function copyReports(
  reports: CollectedReport[],
  destinationDir: string
): Promise<Array<{ source: string; destination: string }>> {
  const copied: Array<{ source: string; destination: string }> = [];

  await fs.mkdir(destinationDir, { recursive: true });

  for (const report of reports) {
    const basename = path.basename(report.path);
    const packageDir = path.join(destinationDir, report.package);
    await fs.mkdir(packageDir, { recursive: true });

    const destination = path.join(packageDir, basename);

    try {
      await fs.copyFile(report.path, destination);
      copied.push({ source: report.path, destination });
    } catch (error) {
      console.error(`Failed to copy ${report.path}: ${error}`);
    }
  }

  return copied;
}

/**
 * Clean old reports from directory
 */
export async function cleanOldReports(
  directory: string,
  keepCount: number
): Promise<string[]> {
  try {
    const entries = await fs.readdir(directory, { withFileTypes: true });

    const directories = entries
      .filter(e => e.isDirectory())
      .map(e => ({
        name: e.name,
        path: path.join(directory, e.name)
      }));

    // Sort by name (assuming run IDs include timestamp)
    directories.sort((a, b) => b.name.localeCompare(a.name));

    const toDelete = directories.slice(keepCount);
    const deleted: string[] = [];

    for (const dir of toDelete) {
      await fs.rm(dir.path, { recursive: true, force: true });
      deleted.push(dir.path);
    }

    return deleted;
  } catch {
    return [];
  }
}

/**
 * Merge multiple reports of same type
 */
export async function mergeReports(
  reports: CollectedReport[],
  outputPath: string
): Promise<void> {
  const cucumberReports = reports.filter(
    r => r.type === 'cucumber' && r.format === 'json'
  );

  if (cucumberReports.length > 0) {
    const merged: unknown[] = [];

    for (const report of cucumberReports) {
      try {
        const content = await fs.readFile(report.path, 'utf-8');
        const data = JSON.parse(content) as unknown;
        if (Array.isArray(data)) {
          merged.push(...data);
        }
      } catch (error) {
        console.error(`Failed to merge ${report.path}: ${error}`);
      }
    }

    await fs.mkdir(path.dirname(outputPath), { recursive: true });
    await fs.writeFile(outputPath, JSON.stringify(merged, null, 2));
  }
}

/**
 * Generate report summary
 */
export function generateReportSummary(
  packages: PackageReports[]
): string {
  const lines: string[] = ['# Collected Test Reports', ''];

  const totalReports = packages.reduce(
    (sum, pkg) => sum + pkg.reports.length,
    0
  );
  const totalSize = packages.reduce(
    (sum, pkg) => pkg.reports.reduce(
      (pkgSum, r) => pkgSum + r.size,
      sum
    ),
    0
  );

  lines.push(`## Summary`);
  lines.push(`- **Packages:** ${packages.length}`);
  lines.push(`- **Reports:** ${totalReports}`);
  lines.push(`- **Total Size:** ${formatBytes(totalSize)}`);
  lines.push('');

  for (const pkg of packages) {
    lines.push(`### ${pkg.packageName}`);
    for (const report of pkg.reports) {
      lines.push(
        `- ${report.type}/${report.format} - ${formatBytes(report.size)}`
      );
    }
    lines.push('');
  }

  return lines.join('\n');
}

/**
 * Format bytes to human readable
 */
function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes}B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)}KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)}MB`;
}