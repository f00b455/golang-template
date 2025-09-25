import { describe, it, expect } from 'vitest';
import {
  generateWikiIndex,
  generateWikiMarkdown,
  formatTimestamp,
  getReportsToRetain,
  parseReportPath,
  buildReportPath,
  type TestReport
} from './index';

describe('Test Reporter Module', () => {
  const mockReports: TestReport[] = [
    {
      runId: 'run-1',
      branch: 'main',
      commitSha: 'abc123',
      timestamp: '2024-01-01T10:00:00Z',
      reportPath: 'reports/main/run-1/index.html',
      status: 'success'
    },
    {
      runId: 'run-2',
      branch: 'main',
      commitSha: 'def456',
      timestamp: '2024-01-02T10:00:00Z',
      reportPath: 'reports/main/run-2/index.html',
      status: 'success'
    },
    {
      runId: 'run-3',
      branch: 'feature-x',
      commitSha: 'ghi789',
      timestamp: '2024-01-03T10:00:00Z',
      reportPath: 'reports/feature-x/run-3/index.html',
      status: 'failure'
    }
  ];

  describe('generateWikiIndex', () => {
    it('should generate index with last run and grouped branches', () => {
      const index = generateWikiIndex(mockReports);

      expect(index.lastRun).toEqual(mockReports[2]); // Most recent by timestamp
      expect(index.recentRuns).toHaveLength(3);
      expect(index.branches).toHaveProperty('main');
      expect(index.branches).toHaveProperty('feature-x');
      expect(index.branches.main).toHaveLength(2);
      expect(index.branches['feature-x']).toHaveLength(1);
    });

    it('should respect max reports per branch config', () => {
      const manyReports: TestReport[] = [];
      for (let i = 0; i < 30; i++) {
        manyReports.push({
          runId: `run-${i}`,
          branch: 'main',
          commitSha: `sha-${i}`,
          timestamp: new Date(2024, 0, i + 1).toISOString(),
          reportPath: `reports/main/run-${i}/index.html`,
          status: 'success'
        });
      }

      const index = generateWikiIndex(manyReports, { maxReportsPerBranch: 10 });
      expect(index.branches.main).toHaveLength(10);
      expect(index.recentRuns).toHaveLength(10);
    });

    it('should handle empty reports', () => {
      const index = generateWikiIndex([]);
      expect(index.lastRun).toBeNull();
      expect(index.recentRuns).toHaveLength(0);
      expect(Object.keys(index.branches)).toHaveLength(0);
    });
  });

  describe('generateWikiMarkdown', () => {
    it('should generate markdown with latest run and history', () => {
      const index = generateWikiIndex(mockReports);
      const markdown = generateWikiMarkdown(index);

      expect(markdown).toContain('# Test Reports');
      expect(markdown).toContain('## Latest Run (run-3)');
      expect(markdown).toContain('Branch: feature-x');
      expect(markdown).toContain('### Recent Runs');
      expect(markdown).toContain('### Branch: main');
      expect(markdown).toContain('### Branch: feature-x');
    });

    it('should handle no reports case', () => {
      const index = generateWikiIndex([]);
      const markdown = generateWikiMarkdown(index);

      expect(markdown).toContain('## No Reports Available');
      expect(markdown).toContain('No test reports have been published yet');
    });
  });

  describe('formatTimestamp', () => {
    it('should format timestamp in UTC', () => {
      const formatted = formatTimestamp('2024-01-01T10:30:45Z', 'UTC');
      expect(formatted).toContain('01/01/2024');
      expect(formatted).toContain('10:30:45');
      expect(formatted).toContain('UTC');
    });

    it('should format timestamp in different timezone', () => {
      const formatted = formatTimestamp('2024-01-01T10:30:45Z', 'America/New_York');
      expect(formatted).toContain('01/01/2024');
      expect(formatted).toContain('05:30:45'); // UTC-5
    });
  });

  describe('getReportsToRetain', () => {
    it('should separate reports to retain and remove', () => {
      const manyReports: TestReport[] = [];
      for (let i = 0; i < 25; i++) {
        manyReports.push({
          runId: `run-${i}`,
          branch: 'main',
          commitSha: `sha-${i}`,
          timestamp: new Date(2024, 0, i + 1).toISOString(),
          reportPath: `reports/main/run-${i}/index.html`,
          status: 'success'
        });
      }

      const { retain, remove } = getReportsToRetain(manyReports, 20);
      expect(retain).toHaveLength(20);
      expect(remove).toHaveLength(5);

      // Most recent should be retained
      expect(retain[0].runId).toBe('run-24');
      // Oldest should be removed
      expect(remove).toContainEqual(expect.objectContaining({ runId: 'run-0' }));
    });

    it('should handle multiple branches', () => {
      const multipleReports: TestReport[] = [
        ...Array(15).fill(null).map((_, i) => ({
          runId: `main-${i}`,
          branch: 'main',
          commitSha: `sha-${i}`,
          timestamp: new Date(2024, 0, i + 1).toISOString(),
          reportPath: `reports/main/main-${i}/index.html`,
          status: 'success' as const
        })),
        ...Array(10).fill(null).map((_, i) => ({
          runId: `feature-${i}`,
          branch: 'feature',
          commitSha: `sha-${i}`,
          timestamp: new Date(2024, 0, i + 20).toISOString(),
          reportPath: `reports/feature/feature-${i}/index.html`,
          status: 'success' as const
        }))
      ];

      const { retain, remove } = getReportsToRetain(multipleReports, 10);

      const mainRetained = retain.filter(r => r.branch === 'main');
      const featureRetained = retain.filter(r => r.branch === 'feature');

      expect(mainRetained).toHaveLength(10);
      expect(featureRetained).toHaveLength(10);
      expect(remove).toHaveLength(5); // 15 main - 10 retained = 5 removed
    });
  });

  describe('parseReportPath', () => {
    it('should parse valid report path', () => {
      const result = parseReportPath('reports/main/run-123/index.html');
      expect(result).toEqual({
        branch: 'main',
        runId: 'run-123'
      });
    });

    it('should parse path without trailing slash', () => {
      const result = parseReportPath('reports/feature-x/run-456');
      expect(result).toEqual({
        branch: 'feature-x',
        runId: 'run-456'
      });
    });

    it('should return null for invalid path', () => {
      expect(parseReportPath('invalid/path')).toBeNull();
      expect(parseReportPath('reports/only-one-part')).toBeNull();
    });
  });

  describe('buildReportPath', () => {
    it('should build report path with default base', () => {
      const path = buildReportPath('main', 'run-123');
      expect(path).toBe('reports/main/run-123/index.html');
    });

    it('should build report path with custom base', () => {
      const path = buildReportPath('feature-x', 'run-456', 'wiki/reports');
      expect(path).toBe('wiki/reports/feature-x/run-456/index.html');
    });
  });

  describe('Helper Functions (via generateWikiIndex)', () => {
    it('should correctly group reports by branch', () => {
      const reportsWithMultipleBranches: TestReport[] = [
        {
          runId: 'run-1',
          branch: 'main',
          commitSha: 'abc123',
          timestamp: '2024-01-01T10:00:00Z',
          reportPath: 'reports/main/run-1/index.html',
          status: 'success'
        },
        {
          runId: 'run-2',
          branch: 'feature-a',
          commitSha: 'def456',
          timestamp: '2024-01-02T10:00:00Z',
          reportPath: 'reports/feature-a/run-2/index.html',
          status: 'success'
        },
        {
          runId: 'run-3',
          branch: 'main',
          commitSha: 'ghi789',
          timestamp: '2024-01-03T10:00:00Z',
          reportPath: 'reports/main/run-3/index.html',
          status: 'success'
        },
        {
          runId: 'run-4',
          branch: 'feature-b',
          commitSha: 'jkl012',
          timestamp: '2024-01-04T10:00:00Z',
          reportPath: 'reports/feature-b/run-4/index.html',
          status: 'failure'
        }
      ];

      const index = generateWikiIndex(reportsWithMultipleBranches);

      // Test groupReportsByBranch via branches property
      expect(Object.keys(index.branches)).toHaveLength(3);
      expect(index.branches).toHaveProperty('main');
      expect(index.branches).toHaveProperty('feature-a');
      expect(index.branches).toHaveProperty('feature-b');
      expect(index.branches.main).toHaveLength(2);
      expect(index.branches['feature-a']).toHaveLength(1);
      expect(index.branches['feature-b']).toHaveLength(1);
    });

    it('should apply branch retention correctly', () => {
      const manyReportsInOneBranch: TestReport[] = [];
      for (let i = 0; i < 15; i++) {
        manyReportsInOneBranch.push({
          runId: `run-${i}`,
          branch: 'main',
          commitSha: `sha-${i}`,
          timestamp: new Date(2024, 0, 15 - i).toISOString(), // Newest first
          reportPath: `reports/main/run-${i}/index.html`,
          status: 'success'
        });
      }

      const index = generateWikiIndex(manyReportsInOneBranch, { maxReportsPerBranch: 5 });

      // Test applyBranchRetention via branches property
      expect(index.branches.main).toHaveLength(5);
      expect(index.branches.main[0].runId).toBe('run-0'); // Most recent
      expect(index.branches.main[4].runId).toBe('run-4'); // 5th most recent
    });
  });

  describe('Helper Functions (via generateWikiMarkdown)', () => {
    it('should generate latest run section correctly', () => {
      const index = generateWikiIndex(mockReports);
      const markdown = generateWikiMarkdown(index);

      // Test generateLatestRunSection
      expect(markdown).toContain('## Latest Run (run-3)');
      expect(markdown).toContain('- **Branch:** feature-x');
      expect(markdown).toContain('- **Status:** failure');
      expect(markdown).toMatch(/- \*\*Time:\*\* .+ UTC/);
    });

    it('should generate recent runs section correctly', () => {
      const index = generateWikiIndex(mockReports);
      const markdown = generateWikiMarkdown(index);

      // Test generateRecentRunsSection
      expect(markdown).toContain('### Recent Runs');
      expect(markdown).toMatch(/- feature-x \/ run-3 → \[Report\]/);
      expect(markdown).toMatch(/- main \/ run-2 → \[Report\]/);
      expect(markdown).toMatch(/- main \/ run-1 → \[Report\]/);
    });

    it('should generate branch sections correctly', () => {
      const index = generateWikiIndex(mockReports);
      const markdown = generateWikiMarkdown(index);

      // Test generateBranchSection
      expect(markdown).toContain('### Branch: main');
      expect(markdown).toContain('### Branch: feature-x');
      // Check that runs are listed under their branch sections
      expect(markdown).toContain('- [run-1](reports/main/run-1/index.html)');
      expect(markdown).toContain('- [run-2](reports/main/run-2/index.html)');
      expect(markdown).toContain('- [run-3](reports/feature-x/run-3/index.html)');
    });
  });

  describe('Helper Functions (via getReportsToRetain)', () => {
    it('should use helper functions for grouping and categorization', () => {
      const reports: TestReport[] = [
        {
          runId: 'run-1',
          branch: 'main',
          commitSha: 'abc123',
          timestamp: '2024-01-01T10:00:00Z',
          reportPath: 'reports/main/run-1/index.html',
          status: 'success'
        },
        {
          runId: 'run-2',
          branch: 'feature',
          commitSha: 'def456',
          timestamp: '2024-01-02T10:00:00Z',
          reportPath: 'reports/feature/run-2/index.html',
          status: 'success'
        },
        {
          runId: 'run-3',
          branch: 'main',
          commitSha: 'ghi789',
          timestamp: '2024-01-03T10:00:00Z',
          reportPath: 'reports/main/run-3/index.html',
          status: 'success'
        }
      ];

      const { retain, remove } = getReportsToRetain(reports, 1);

      // Test groupReportsToMap and categorizeReports
      expect(retain).toHaveLength(2); // One per branch
      expect(remove).toHaveLength(1); // Older main report
      expect(retain.find(r => r.branch === 'main')?.runId).toBe('run-3');
      expect(retain.find(r => r.branch === 'feature')?.runId).toBe('run-2');
      expect(remove[0].runId).toBe('run-1');
    });
  });
});