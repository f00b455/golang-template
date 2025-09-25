import { describe, it, expect, vi, beforeEach } from 'vitest';
import * as fs from 'fs/promises';
import * as path from 'path';
import {
  collectReports,
  findReportsInDirectory,
  generateReportSummary,
  CollectedReport,
  PackageReports
} from './report-collector';

vi.mock('glob', () => ({
  glob: vi.fn()
}));

vi.mock('fs/promises', () => ({
  default: {},
  readFile: vi.fn(),
  stat: vi.fn(),
  mkdir: vi.fn(),
  copyFile: vi.fn(),
  readdir: vi.fn(),
  rm: vi.fn(),
  writeFile: vi.fn()
}));

describe('report-collector', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('generateReportSummary', () => {
    it('generates summary for collected reports', () => {
      const packages: PackageReports[] = [
        {
          packageName: 'package-a',
          packagePath: '/path/to/package-a',
          reports: [
            {
              type: 'cucumber',
              format: 'json',
              path: '/path/to/package-a/cucumber.json',
              package: 'package-a',
              size: 1024
            },
            {
              type: 'vitest',
              format: 'html',
              path: '/path/to/package-a/vitest.html',
              package: 'package-a',
              size: 2048
            }
          ]
        },
        {
          packageName: 'package-b',
          packagePath: '/path/to/package-b',
          reports: [
            {
              type: 'cucumber',
              format: 'json',
              path: '/path/to/package-b/cucumber.json',
              package: 'package-b',
              size: 512
            }
          ]
        }
      ];

      const summary = generateReportSummary(packages);

      expect(summary).toContain('# Collected Test Reports');
      expect(summary).toContain('Packages:** 2');
      expect(summary).toContain('Reports:** 3');
      expect(summary).toContain('Total Size:**');
      expect(summary).toContain('### package-a');
      expect(summary).toContain('### package-b');
      expect(summary).toContain('cucumber/json');
      expect(summary).toContain('vitest/html');
      expect(summary).toContain('1.0KB');
      expect(summary).toContain('2.0KB');
    });

    it('handles empty packages array', () => {
      const summary = generateReportSummary([]);

      expect(summary).toContain('Packages:** 0');
      expect(summary).toContain('Reports:** 0');
      expect(summary).toContain('Total Size:** 0B');
    });
  });

  describe('collectReports', () => {
    it('collects reports from packages', async () => {
      const glob = await import('glob');
      const mockGlob = glob.glob as any;

      mockGlob.mockResolvedValueOnce([
        '/root/packages/foo/cucumber-report.json',
        '/root/packages/bar/vitest-report.html'
      ]);

      const mockStat = fs.stat as any;
      mockStat.mockResolvedValue({ size: 1000 });

      const mockReadFile = fs.readFile as any;
      mockReadFile.mockResolvedValue('{"cucumber": true}');

      const result = await collectReports({
        rootDir: '/root'
      });

      expect(result).toHaveLength(2);
      expect(result[0].packageName).toBe('bar');
      expect(result[1].packageName).toBe('foo');
      expect(result[0].reports[0].type).toBe('vitest');
      expect(result[1].reports[0].type).toBe('cucumber');
    });
  });

  describe('findReportsInDirectory', () => {
    it('finds reports in a specific directory', async () => {
      const glob = await import('glob');
      const mockGlob = glob.glob as any;

      mockGlob.mockResolvedValueOnce([
        '/dir/cucumber-report.json'
      ]);

      const mockStat = fs.stat as any;
      mockStat.mockResolvedValue({ size: 500 });

      const mockReadFile = fs.readFile as any;
      mockReadFile.mockResolvedValue('{"cucumber": true}');

      const result = await findReportsInDirectory('/dir');

      expect(result).toHaveLength(1);
      expect(result[0].type).toBe('cucumber');
      expect(result[0].path).toBe('/dir/cucumber-report.json');
    });
  });
});