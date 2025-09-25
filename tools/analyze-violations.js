#!/usr/bin/env node

/**
 * Analyzes ESLint Clean Code violations across the codebase
 * Generates a report showing which files/functions need refactoring
 */

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

const CLEAN_CODE_RULES = [
  'max-lines-per-function',
  'max-lines',
  'max-statements',
  'complexity',
  'max-depth',
  'max-params'
];

// Configurable buffer size via environment variable (defaults to 10MB)
const BUFFER_SIZE = parseInt(process.env.ESLINT_REPORT_BUFFER_SIZE || '10', 10) * 1024 * 1024;

function runLintWithFormat(format = 'json') {
  try {
    const output = execSync(`pnpm lint --format ${format}`, {
      encoding: 'utf8',
      maxBuffer: BUFFER_SIZE
    return output;
  } catch (error) {
    // ESLint exits with non-zero when there are warnings/errors
    // But we still get the output we need
    return error.stdout || '';
  }
}

function analyzeViolations() {
  console.log('🔍 Analyzing Clean Code violations...\n');
  console.log('⏳ Running ESLint across all packages, this may take a moment...\n');

  // Run ESLint with JSON format
  process.stdout.write('📊 Collecting violations...');
  const jsonOutput = runLintWithFormat('json');
  process.stdout.write(' Done!\n\n');

  if (!jsonOutput) {
    console.log('✅ No violations found!');
    return;
  }

  let results;
  try {
    results = JSON.parse(jsonOutput);
  } catch (error) {
    console.error('Failed to parse ESLint output');
    return;
  }

  // Statistics
  const stats = {
    totalFiles: 0,
    filesWithViolations: 0,
    violationsByRule: {},
    violationsByPackage: {},
    criticalFiles: [] // Files with most violations
  };

  CLEAN_CODE_RULES.forEach(rule => {
    stats.violationsByRule[rule] = {
      count: 0,
      files: []
    };
  });

  // Process results
  results.forEach(file => {
    if (!file.messages || file.messages.length === 0) return;

    stats.totalFiles++;

    const cleanCodeViolations = file.messages.filter(msg =>
      CLEAN_CODE_RULES.includes(msg.ruleId)
    );

    if (cleanCodeViolations.length > 0) {
      stats.filesWithViolations++;

      // Determine package
      const relativePath = path.relative(process.cwd(), file.filePath);
      const packageMatch = relativePath.match(/^(apps|packages)\/([^/]+)/);
      const packageName = packageMatch ? `${packageMatch[1]}/${packageMatch[2]}` : 'root';

      if (!stats.violationsByPackage[packageName]) {
        stats.violationsByPackage[packageName] = {
          count: 0,
          files: []
        };
      }

      stats.violationsByPackage[packageName].count += cleanCodeViolations.length;
      stats.violationsByPackage[packageName].files.push({
        file: relativePath,
        violations: cleanCodeViolations.length
      });

      // Track by rule
      cleanCodeViolations.forEach(violation => {
        stats.violationsByRule[violation.ruleId].count++;
        stats.violationsByRule[violation.ruleId].files.push(relativePath);
      });

      // Track critical files
      if (cleanCodeViolations.length >= 3) {
        stats.criticalFiles.push({
          file: relativePath,
          violations: cleanCodeViolations
        });
      }
    }
  });

  // Generate report
  console.log('\n📊 Clean Code Violation Report\n');
  console.log('═'.repeat(50));

  console.log('\n📈 Summary:');
  console.log(`  Total files analyzed: ${stats.totalFiles}`);
  console.log(`  Files with violations: ${stats.filesWithViolations}`);
  console.log(`  Compliance rate: ${((1 - stats.filesWithViolations / stats.totalFiles) * 100).toFixed(1)}%`);

  console.log('\n📋 Violations by Rule:');
  Object.entries(stats.violationsByRule)
    .sort((a, b) => b[1].count - a[1].count)
    .forEach(([rule, data]) => {
      if (data.count > 0) {
        console.log(`  ${rule}: ${data.count} violations in ${data.files.length} files`);
      }
    });

  console.log('\n📦 Violations by Package:');
  Object.entries(stats.violationsByPackage)
    .sort((a, b) => b[1].count - a[1].count)
    .forEach(([pkg, data]) => {
      console.log(`  ${pkg}: ${data.count} violations in ${data.files.length} files`);
    });

  if (stats.criticalFiles.length > 0) {
    console.log('\n⚠️  Critical Files (3+ violations):');
    stats.criticalFiles
      .sort((a, b) => b.violations.length - a.violations.length)
      .slice(0, 10)
      .forEach(({ file, violations }) => {
        console.log(`  ${file}:`);
        const ruleCount = {};
        violations.forEach(v => {
          ruleCount[v.ruleId] = (ruleCount[v.ruleId] || 0) + 1;
        });
        Object.entries(ruleCount).forEach(([rule, count]) => {
          console.log(`    - ${rule}: ${count} violation(s)`);
        });
      });
  }

  console.log('\n💡 Migration Recommendations:');
  console.log('  1. Start with critical files (3+ violations)');
  console.log('  2. Focus on max-lines-per-function violations first');
  console.log('  3. Refactor package by package');
  console.log('  4. Create separate PRs for each major refactoring');

  console.log('\n📝 Next Steps:');
  console.log('  - Run `pnpm lint` to see all warnings');
  console.log('  - Set ESLINT_CLEAN_CODE_PHASE=3 to test error mode');
  console.log('  - Update .eslintrc files after refactoring each package');

  // Write detailed report to file
  const reportPath = path.join(process.cwd(), 'clean-code-report.json');
  fs.writeFileSync(reportPath, JSON.stringify(stats, null, 2));
  console.log(`\n💾 Detailed report saved to: ${reportPath}`);
}

// Run analysis
analyzeViolations();