import * as fs from 'fs/promises';
import * as path from 'path';
import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

export interface WikiPublisherOptions {
  wikiPath: string;
  branch: string;
  runId: string;
  commitSha: string;
  maxReportsPerBranch?: number;
  dryRun?: boolean;
}

export interface PublishResult {
  success: boolean;
  reportPath: string;
  filesPublished: number;
  filesDeleted: number;
  error?: string;
}

const DEFAULT_MAX_REPORTS = 20;

/**
 * Check if directory is a git repository
 */
async function isGitRepo(dir: string): Promise<boolean> {
  try {
    await execAsync('git status', { cwd: dir });
    return true;
  } catch {
    return false;
  }
}

/**
 * Clone wiki repository if not exists
 */
export async function ensureWikiRepo(
  wikiUrl: string,
  localPath: string
): Promise<void> {
  try {
    await fs.access(localPath);
    const isRepo = await isGitRepo(localPath);

    if (isRepo) {
      // Pull latest changes
      await execAsync('git pull', { cwd: localPath });
    } else {
      // Directory exists but not a repo, remove and clone
      await fs.rm(localPath, { recursive: true, force: true });
      await execAsync(`git clone ${wikiUrl} ${localPath}`);
    }
  } catch {
    // Directory doesn't exist, clone
    await execAsync(`git clone ${wikiUrl} ${localPath}`);
  }
}

/**
 * Get list of existing run directories for a branch
 */
async function getExistingRuns(
  wikiPath: string,
  branch: string
): Promise<string[]> {
  const branchPath = path.join(wikiPath, 'reports', branch);

  try {
    const entries = await fs.readdir(branchPath, { withFileTypes: true });
    return entries
      .filter(e => e.isDirectory())
      .map(e => e.name)
      .sort((a, b) => b.localeCompare(a));
  } catch {
    return [];
  }
}

/**
 * Clean old reports based on retention policy
 */
async function cleanOldReports(
  wikiPath: string,
  branch: string,
  maxReports: number
): Promise<number> {
  const existingRuns = await getExistingRuns(wikiPath, branch);
  const runsToDelete = existingRuns.slice(maxReports);
  let deletedCount = 0;

  for (const runId of runsToDelete) {
    const runPath = path.join(wikiPath, 'reports', branch, runId);
    try {
      await fs.rm(runPath, { recursive: true, force: true });
      deletedCount++;
    } catch (error) {
      console.error(`Failed to delete ${runPath}: ${error}`);
    }
  }

  return deletedCount;
}

/**
 * Copy reports to wiki repository
 */
async function copyReportsToWiki(
  sourcePath: string,
  wikiPath: string,
  branch: string,
  runId: string
): Promise<number> {
  const targetPath = path.join(wikiPath, 'reports', branch, runId);
  await fs.mkdir(targetPath, { recursive: true });

  let copiedCount = 0;

  try {
    const files = await fs.readdir(sourcePath);

    for (const file of files) {
      const srcFile = path.join(sourcePath, file);
      const destFile = path.join(targetPath, file);

      const stat = await fs.stat(srcFile);
      if (stat.isDirectory()) {
        await copyDir(srcFile, destFile);
      } else {
        await fs.copyFile(srcFile, destFile);
        copiedCount++;
      }
    }
  } catch (error) {
    console.error(`Failed to copy reports: ${error}`);
  }

  return copiedCount;
}

/**
 * Recursively copy directory
 */
async function copyDir(src: string, dest: string): Promise<void> {
  await fs.mkdir(dest, { recursive: true });
  const entries = await fs.readdir(src, { withFileTypes: true });

  for (const entry of entries) {
    const srcPath = path.join(src, entry.name);
    const destPath = path.join(dest, entry.name);

    if (entry.isDirectory()) {
      await copyDir(srcPath, destPath);
    } else {
      await fs.copyFile(srcPath, destPath);
    }
  }
}

/**
 * Update wiki index page
 */
async function updateWikiIndex(
  wikiPath: string,
  branch: string,
  runId: string,
  commitSha: string
): Promise<void> {
  const indexPath = path.join(wikiPath, 'Home.md');
  const timestamp = new Date().toISOString();

  let content: string;
  try {
    content = await fs.readFile(indexPath, 'utf-8');
  } catch {
    content = '# Test Reports\n\nAutomated test reports from CI/CD.\n\n';
  }

  // Remove old latest run section if exists
  const latestRunRegex = /## Latest Run[\s\S]*?(?=##|$)/;
  content = content.replace(latestRunRegex, '');

  // Add new latest run section
  const latestSection = [
    '## Latest Run\n',
    `- **Branch:** ${branch}\n`,
    `- **Run ID:** [${runId}](reports/${branch}/${runId}/index.md)\n`,
    `- **Commit:** ${commitSha}\n`,
    `- **Updated:** ${timestamp}\n`,
    '\n'
  ].join('');

  // Insert after title
  const lines = content.split('\n');
  const titleIndex = lines.findIndex(l => l.startsWith('# '));
  lines.splice(titleIndex + 2, 0, latestSection);

  await fs.writeFile(indexPath, lines.join('\n'));
}

/**
 * Update branch-specific index
 */
async function updateBranchIndex(
  wikiPath: string,
  branch: string,
  _runId: string
): Promise<void> {
  const indexPath = path.join(wikiPath, '_index', `${branch}.md`);
  await fs.mkdir(path.dirname(indexPath), { recursive: true });

  const timestamp = new Date().toISOString();
  const existingRuns = await getExistingRuns(wikiPath, branch);

  const lines = [
    `# Test Reports - ${branch}`,
    '',
    '## Recent Runs',
    ''
  ];

  for (const run of existingRuns.slice(0, 20)) {
    const runPath = `../reports/${branch}/${run}/index.md`;
    lines.push(`- [${run}](${runPath}) - ${timestamp}`);
  }

  await fs.writeFile(indexPath, lines.join('\n'));
}

/**
 * Commit and push changes to wiki
 */
async function commitAndPush(
  wikiPath: string,
  message: string,
  dryRun?: boolean
): Promise<void> {
  if (dryRun) {
    console.error('Dry run mode - skipping git operations');
    return;
  }

  try {
    await execAsync('git add .', { cwd: wikiPath });
    await execAsync(
      `git commit -m "${message}"`,
      { cwd: wikiPath }
    );
    await execAsync('git push', { cwd: wikiPath });
  } catch (error) {
    console.error(`Git operations failed: ${error}`);
    throw error;
  }
}

/**
 * Publish reports to GitHub Wiki
 */
export async function publishToWiki(
  reportPath: string,
  options: WikiPublisherOptions
): Promise<PublishResult> {
  const maxReports = options.maxReportsPerBranch || DEFAULT_MAX_REPORTS;

  try {
    // Copy reports to wiki directory
    const filesPublished = await copyReportsToWiki(
      reportPath,
      options.wikiPath,
      options.branch,
      options.runId
    );

    // Clean old reports
    const filesDeleted = await cleanOldReports(
      options.wikiPath,
      options.branch,
      maxReports
    );

    // Update index pages
    await updateWikiIndex(
      options.wikiPath,
      options.branch,
      options.runId,
      options.commitSha
    );

    await updateBranchIndex(
      options.wikiPath,
      options.branch,
      options.runId
    );

    // Commit and push
    const commitMessage = `Update test reports: ${options.branch}/${options.runId}`;
    await commitAndPush(options.wikiPath, commitMessage, options.dryRun);

    return {
      success: true,
      reportPath: `reports/${options.branch}/${options.runId}/`,
      filesPublished,
      filesDeleted
    };
  } catch (error) {
    return {
      success: false,
      reportPath: '',
      filesPublished: 0,
      filesDeleted: 0,
      error: String(error)
    };
  }
}

/**
 * Create fallback entry when no reports available
 */
export async function createFallbackEntry(
  wikiPath: string,
  branch: string,
  runId: string,
  reason: string
): Promise<void> {
  const indexPath = path.join(wikiPath, 'Home.md');
  const timestamp = new Date().toISOString();

  let content: string;
  try {
    content = await fs.readFile(indexPath, 'utf-8');
  } catch {
    content = '# Test Reports\n\nAutomated test reports from CI/CD.\n\n';
  }

  const fallbackSection = [
    `## Run ${runId} - No Reports Available\n`,
    `- **Branch:** ${branch}\n`,
    `- **Reason:** ${reason}\n`,
    `- **Timestamp:** ${timestamp}\n`,
    '\n'
  ].join('');

  content += fallbackSection;
  await fs.writeFile(indexPath, content);
}

/**
 * Check if wiki is accessible
 */
export async function checkWikiAccess(wikiUrl: string): Promise<boolean> {
  try {
    const { stdout } = await execAsync(`git ls-remote ${wikiUrl}`);
    return stdout.length > 0;
  } catch {
    return false;
  }
}

/**
 * Generate wiki URL from repository URL
 */
export function generateWikiUrl(repoUrl: string): string {
  // Convert repo URL to wiki URL
  // https://github.com/owner/repo.git -> https://github.com/owner/repo.wiki.git
  return repoUrl.replace(/\.git$/, '').concat('.wiki.git');
}