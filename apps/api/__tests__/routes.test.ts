import { describe, it, expect, vi } from 'vitest';

// Mock the shared package
vi.mock('@ts-template/shared', () => ({
  greet: vi.fn((name: string) => `Hello, ${name}!`)
}));

describe('API Routes', () => {
  it('should have basic test structure', () => {
    expect(true).toBe(true);
  });
});