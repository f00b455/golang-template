import { describe, it, expect } from 'vitest';
import {
  cucumberJsonToMarkdown,
  generateIndexPage,
  formatDuration,
  calculateStats,
  MarkdownOptions
} from './markdown-formatter';

describe('markdown-formatter', () => {
  describe('formatDuration', () => {
    it('formats nanoseconds to milliseconds', () => {
      expect(formatDuration(500000)).toBe('1ms');
      expect(formatDuration(999999999)).toBe('1000ms');
    });

    it('formats to seconds for large values', () => {
      expect(formatDuration(1500000000)).toBe('1.50s');
      expect(formatDuration(5000000000)).toBe('5.00s');
    });

    it('handles undefined input', () => {
      expect(formatDuration(undefined)).toBe('0ms');
    });
  });

  describe('cucumberJsonToMarkdown', () => {
    it('generates markdown from cucumber report', () => {
      const report = [
        {
          name: 'Test Feature',
          keyword: 'Feature',
          description: 'A test feature',
          tags: [{ name: '@test' }],
          elements: [
            {
              type: 'scenario',
              name: 'Test Scenario',
              keyword: 'Scenario',
              tags: [],
              steps: [
                {
                  name: 'a test step',
                  keyword: 'Given ',
                  result: {
                    status: 'passed',
                    duration: 1000000
                  }
                }
              ]
            }
          ]
        }
      ];

      const options: MarkdownOptions = {
        runId: 'run-123',
        branch: 'main',
        commitSha: 'abc123',
        timestamp: '2024-01-01T00:00:00Z'
      };

      const markdown = cucumberJsonToMarkdown(report, options);

      expect(markdown).toContain('# Test Report - run-123');
      expect(markdown).toContain('Branch: main');
      expect(markdown).toContain('Commit: abc123');
      expect(markdown).toContain('Feature: Test Feature');
      expect(markdown).toContain('Scenario: Test Scenario');
      expect(markdown).toContain('✅');
    });

    it('handles failed scenarios', () => {
      const report = [
        {
          elements: [
            {
              type: 'scenario',
              name: 'Failed Scenario',
              steps: [
                {
                  name: 'a failed step',
                  keyword: 'When ',
                  result: {
                    status: 'failed',
                    error_message: 'Step failed'
                  }
                }
              ]
            }
          ]
        }
      ];

      const options: MarkdownOptions = {
        runId: 'run-456',
        branch: 'feature',
        commitSha: 'def456',
        timestamp: '2024-01-01T00:00:00Z'
      };

      const markdown = cucumberJsonToMarkdown(report, options);

      expect(markdown).toContain('FAILED ❌');
      expect(markdown).toContain('Failed:**');
      expect(markdown).toContain('Step failed');
    });
  });

  describe('generateIndexPage', () => {
    it('generates index page with reports grouped by branch', () => {
      const reports = [
        {
          runId: 'run-1',
          branch: 'main',
          timestamp: '2024-01-01T00:00:00Z',
          path: 'reports/main/run-1/',
          status: 'success'
        },
        {
          runId: 'run-2',
          branch: 'feature',
          timestamp: '2024-01-02T00:00:00Z',
          path: 'reports/feature/run-2/',
          status: 'failed'
        }
      ];

      const index = generateIndexPage(reports);

      expect(index).toContain('# Test Reports Index');
      expect(index).toContain('## Branch: main');
      expect(index).toContain('## Branch: feature');
      expect(index).toContain('[run-1]');
      expect(index).toContain('[run-2]');
      expect(index).toContain('✅');
      expect(index).toContain('❌');
    });

    it('respects maxReports option', () => {
      const reports = Array.from({ length: 25 }, (_, i) => ({
        runId: `run-${i}`,
        branch: 'main',
        timestamp: new Date(2024, 0, i + 1).toISOString(),
        path: `reports/main/run-${i}/`,
        status: 'success'
      }));

      const index = generateIndexPage(reports, { maxReports: 5 });
      const matches = index.match(/\[run-\d+\]/g); // Only match the link text, not the URL

      expect(matches?.length).toBe(5);
    });
  });
});