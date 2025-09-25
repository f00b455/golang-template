#!/usr/bin/env node
import * as fs from 'fs/promises';
import * as path from 'path';
import { collectReports, copyReports, mergeReports } from './report-collector';
import {
  cucumberJsonToMarkdown,
  generateIndexPage,
  readCucumberReport,
  writeMarkdownReport
} from './markdown-formatter';
import {
  publishToWiki,
  ensureWikiRepo,
  checkWikiAccess,
  generateWikiUrl,
  createFallbackEntry
} from './wiki-publisher';

interface CliOptions {
  command: 'collect' | 'format' | 'publish' | 'all';
  rootDir?: string;
  outputDir?: string;
  wikiUrl?: string;
  wikiPath?: string;
  branch?: string;
  runId?: string;
  commitSha?: string;
  maxReports?: number;
  dryRun?: boolean;
  verbose?: boolean;
}

/**
 * Parse command line arguments
 */
function parseArgs(args: string[]): CliOptions {
  const options: CliOptions = {
    command: 'all',
    rootDir: process.cwd(),
    outputDir: './test-reports',
    branch: process.env.GITHUB_REF_NAME || 'main',
    runId: process.env.GITHUB_RUN_ID || `run-${Date.now()}`,
    commitSha: process.env.GITHUB_SHA || 'unknown',
    maxReports: 20,
    dryRun: false,
    verbose: false
  };

  for (let i = 2; i < args.length; i++) {
    const arg = args[i];
    const value = args[i + 1];

    switch (arg) {
      case 'collect':
      case 'format':
      case 'publish':
      case 'all':
        options.command = arg as CliOptions['command'];
        break;
      case '--root':
      case '-r':
        if (value) options.rootDir = value;
        i++;
        break;
      case '--output':
      case '-o':
        if (value) options.outputDir = value;
        i++;
        break;
      case '--wiki-url':
        if (value) options.wikiUrl = value;
        i++;
        break;
      case '--wiki-path':
        if (value) options.wikiPath = value;
        i++;
        break;
      case '--branch':
      case '-b':
        if (value) options.branch = value;
        i++;
        break;
      case '--run-id':
        if (value) options.runId = value;
        i++;
        break;
      case '--commit':
      case '-c':
        if (value) options.commitSha = value;
        i++;
        break;
      case '--max-reports':
        if (value) options.maxReports = parseInt(value, 10);
        i++;
        break;
      case '--dry-run':
        options.dryRun = true;
        break;
      case '--verbose':
      case '-v':
        options.verbose = true;
        break;
      case '--help':
      case '-h':
        printHelp();
        process.exit(0);
    }
  }

  return options;
}

/**
 * Print help message
 */
function printHelp(): void {
  console.log(`
Test Reporter CLI

Usage: test-reporter [command] [options]

Commands:
  collect    Collect test reports from packages
  format     Format reports to markdown
  publish    Publish reports to GitHub Wiki
  all        Run all steps (default)

Options:
  --root, -r <dir>      Root directory to search (default: cwd)
  --output, -o <dir>    Output directory for reports (default: ./test-reports)
  --wiki-url <url>      GitHub Wiki repository URL
  --wiki-path <path>    Local path to wiki repository
  --branch, -b <name>   Branch name (default: from env or 'main')
  --run-id <id>         CI run ID (default: from env or timestamp)
  --commit, -c <sha>    Commit SHA (default: from env or 'unknown')
  --max-reports <n>     Max reports per branch (default: 20)
  --dry-run             Don't push changes to wiki
  --verbose, -v         Enable verbose output
  --help, -h            Show this help message

Environment Variables:
  GITHUB_REPOSITORY     Repository name (owner/repo)
  GITHUB_REF_NAME       Branch name
  GITHUB_RUN_ID         CI run ID
  GITHUB_SHA            Commit SHA
  WIKI_ACCESS_TOKEN     GitHub token for wiki access
`);
}

/**
 * Log message if verbose
 */
function log(message: string, verbose?: boolean): void {
  if (verbose) {
    console.log(`[test-reporter] ${message}`);
  }
}

/**
 * Collect reports from packages
 */
async function collectCommand(options: CliOptions): Promise<void> {
  log(`Collecting reports from ${options.rootDir}`, options.verbose);

  const packages = await collectReports({
    rootDir: options.rootDir!
  });

  log(`Found ${packages.length} packages with reports`, options.verbose);

  const allReports = packages.flatMap(pkg => pkg.reports);
  const outputDir = path.join(options.outputDir!, 'collected');

  const copied = await copyReports(allReports, outputDir);
  log(`Copied ${copied.length} reports to ${outputDir}`, options.verbose);

  // Merge Cucumber reports
  const cucumberReports = allReports.filter(
    r => r.type === 'cucumber' && r.format === 'json'
  );

  if (cucumberReports.length > 0) {
    const mergedPath = path.join(outputDir, 'merged-cucumber.json');
    await mergeReports(cucumberReports, mergedPath);
    log(`Merged ${cucumberReports.length} Cucumber reports`, options.verbose);
  }

  console.log(`✅ Collected ${allReports.length} reports from ${packages.length} packages`);
}

/**
 * Format reports to markdown
 */
async function formatCommand(options: CliOptions): Promise<void> {
  log(`Formatting reports in ${options.outputDir}`, options.verbose);

  const collectedDir = path.join(options.outputDir!, 'collected');
  const formattedDir = path.join(options.outputDir!, 'formatted');

  await fs.mkdir(formattedDir, { recursive: true });

  // Look for merged Cucumber report
  const mergedPath = path.join(collectedDir, 'merged-cucumber.json');

  try {
    const cucumberData = await readCucumberReport(mergedPath);

    const markdown = cucumberJsonToMarkdown(cucumberData, {
      runId: options.runId!,
      branch: options.branch!,
      commitSha: options.commitSha!,
      timestamp: new Date().toISOString()
    });

    const outputPath = path.join(formattedDir, 'report.md');
    await writeMarkdownReport(markdown, outputPath);

    log(`Generated markdown report at ${outputPath}`, options.verbose);

    // Generate index page
    const indexContent = generateIndexPage([{
      runId: options.runId!,
      branch: options.branch!,
      timestamp: new Date().toISOString(),
      path: `reports/${options.branch}/${options.runId}/report.md`,
      status: 'success'
    }]);

    const indexPath = path.join(formattedDir, 'index.md');
    await fs.writeFile(indexPath, indexContent);

    console.log(`✅ Formatted reports to markdown`);
  } catch (error) {
    console.error(`Failed to format reports: ${error}`);

    // Create fallback
    const fallbackContent = [
      `# Test Report - ${options.runId}`,
      '',
      '## No Reports Available',
      '',
      `- **Branch:** ${options.branch}`,
      `- **Run ID:** ${options.runId}`,
      `- **Commit:** ${options.commitSha}`,
      `- **Timestamp:** ${new Date().toISOString()}`,
      `- **Reason:** Failed to process test reports`,
      ''
    ].join('\n');

    const fallbackPath = path.join(formattedDir, 'report.md');
    await fs.writeFile(fallbackPath, fallbackContent);

    console.log(`⚠️ Created fallback report due to processing error`);
  }
}

/**
 * Publish reports to wiki
 */
async function publishCommand(options: CliOptions): Promise<void> {
  if (!options.wikiUrl && !options.wikiPath) {
    // Try to generate wiki URL from repository
    const repo = process.env.GITHUB_REPOSITORY;
    if (repo) {
      options.wikiUrl = generateWikiUrl(
        `https://github.com/${repo}.git`
      );
      log(`Generated wiki URL: ${options.wikiUrl}`, options.verbose);
    } else {
      console.error('❌ Wiki URL or path required for publishing');
      process.exit(1);
    }
  }

  if (!options.wikiPath) {
    options.wikiPath = path.join(options.outputDir!, 'wiki');
  }

  log(`Publishing to wiki at ${options.wikiPath}`, options.verbose);

  // Ensure wiki repository exists
  if (options.wikiUrl) {
    const accessible = await checkWikiAccess(options.wikiUrl);
    if (!accessible) {
      console.error(`❌ Cannot access wiki at ${options.wikiUrl}`);
      console.error('Make sure the wiki is enabled and you have access');

      // Create fallback entry
      await createFallbackEntry(
        options.wikiPath,
        options.branch!,
        options.runId!,
        'Wiki not accessible'
      );
      return;
    }

    await ensureWikiRepo(options.wikiUrl, options.wikiPath);
  }

  const formattedDir = path.join(options.outputDir!, 'formatted');

  const result = await publishToWiki(formattedDir, {
    wikiPath: options.wikiPath,
    branch: options.branch!,
    runId: options.runId!,
    commitSha: options.commitSha!,
    maxReportsPerBranch: options.maxReports ?? 20,
    dryRun: options.dryRun ?? false
  });

  if (result.success) {
    console.log(`✅ Published ${result.filesPublished} files to wiki`);
    console.log(`   Cleaned ${result.filesDeleted} old reports`);
    console.log(`   Report path: ${result.reportPath}`);
  } else {
    console.error(`❌ Failed to publish: ${result.error}`);
    process.exit(1);
  }
}

/**
 * Run all commands
 */
async function runAll(options: CliOptions): Promise<void> {
  await collectCommand(options);
  await formatCommand(options);
  await publishCommand(options);
}

/**
 * Main entry point
 */
async function main(): Promise<void> {
  const options = parseArgs(process.argv);

  log('Starting test-reporter', options.verbose);
  log(`Command: ${options.command}`, options.verbose);

  try {
    switch (options.command) {
      case 'collect':
        await collectCommand(options);
        break;
      case 'format':
        await formatCommand(options);
        break;
      case 'publish':
        await publishCommand(options);
        break;
      case 'all':
        await runAll(options);
        break;
    }

    log('Completed successfully', options.verbose);
  } catch (error) {
    console.error(`❌ Error: ${error}`);
    process.exit(1);
  }
}

// Run if called directly
if (require.main === module) {
  main().catch(error => {
    console.error(error);
    process.exit(1);
  });
}

export { main, parseArgs };
export type { CliOptions };