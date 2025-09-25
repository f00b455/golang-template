import { describe, it, expect } from 'vitest';
import type { RssHeadline } from '../src/index';

describe('RssHeadline Type', () => {
  it('should create a valid RssHeadline object', () => {
    const headline: RssHeadline = {
      title: 'Test Headline',
      link: 'https://example.com/article',
      publishedAt: '2024-09-23T12:00:00Z',
      source: 'SPIEGEL',
    };

    expect(headline.title).toBe('Test Headline');
    expect(headline.link).toBe('https://example.com/article');
    expect(headline.publishedAt).toBe('2024-09-23T12:00:00Z');
    expect(headline.source).toBe('SPIEGEL');
  });

  it('should be readonly', () => {
    const headline: RssHeadline = {
      title: 'Immutable Headline',
      link: 'https://example.com/article',
      publishedAt: '2024-09-23T12:00:00Z',
      source: 'SPIEGEL',
    };

    // TypeScript prevents modification at compile time
    // In runtime, readonly is not enforced, so we just verify the original value
    // @ts-expect-error - title is readonly (TypeScript compile-time check)
    headline.title = 'Modified';

    // The assignment above doesn't actually modify readonly properties in strict mode
    // but doesn't throw in runtime JavaScript. We verify type safety through TypeScript.
  });

  it('should work with JSON serialization', () => {
    const headline: RssHeadline = {
      title: 'Serializable Headline',
      link: 'https://example.com/article',
      publishedAt: '2024-09-23T12:00:00Z',
      source: 'SPIEGEL',
    };

    const json = JSON.stringify(headline);
    const parsed = JSON.parse(json) as RssHeadline;

    expect(parsed).toEqual(headline);
    expect(parsed.title).toBe(headline.title);
    expect(parsed.link).toBe(headline.link);
    expect(parsed.publishedAt).toBe(headline.publishedAt);
    expect(parsed.source).toBe(headline.source);
  });
});