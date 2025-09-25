import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import Fastify from 'fastify';
import { rssRoute, resetCache } from './rss.js';

// Mock fetch globally
global.fetch = vi.fn();

describe('RSS Route', () => {
  let server: ReturnType<typeof Fastify>;

  beforeEach(async () => {
    server = Fastify();
    await server.register(rssRoute);
    vi.clearAllMocks();
    resetCache(); // Reset cache between tests
  });

  afterEach(async () => {
    await server.close();
    vi.resetAllMocks();
  });

  describe('GET /rss/spiegel/latest', () => {
    it('should return RSS headline when fetch succeeds', async () => {
      const mockRssResponse = `<?xml version="1.0" encoding="UTF-8"?>
        <rss version="2.0">
          <channel>
            <item>
              <title>Test Headline</title>
              <link>https://example.com/article</link>
              <pubDate>Mon, 23 Sep 2024 12:00:00 +0200</pubDate>
            </item>
          </channel>
        </rss>`;

      (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: vi.fn().mockResolvedValue(mockRssResponse),
      } as unknown as Response);

      const response = await server.inject({
        method: 'GET',
        url: '/rss/spiegel/latest',
      });

      expect(response.statusCode).toBe(200);
      const data = JSON.parse(response.body);
      expect(data).toMatchObject({
        title: 'Test Headline',
        link: 'https://example.com/article',
        source: 'SPIEGEL',
      });
      expect(data.publishedAt).toBeDefined();
    });

    it('should return 503 error when RSS fetch fails', async () => {
      (global.fetch as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Network error'));

      const response = await server.inject({
        method: 'GET',
        url: '/rss/spiegel/latest',
      });

      expect(response.statusCode).toBe(503);
      const data = JSON.parse(response.body);
      expect(data).toEqual({ error: 'Unable to fetch RSS feed' });
    });

    it('should return 503 when RSS response is not ok', async () => {
      (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: false,
        status: 404,
      } as unknown as Response);

      const response = await server.inject({
        method: 'GET',
        url: '/rss/spiegel/latest',
      });

      expect(response.statusCode).toBe(503);
      const data = JSON.parse(response.body);
      expect(data).toEqual({ error: 'Unable to fetch RSS feed' });
    });

    it('should return 503 when RSS has no items', async () => {
      const mockRssResponse = `<?xml version="1.0" encoding="UTF-8"?>
        <rss version="2.0">
          <channel>
          </channel>
        </rss>`;

      (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: vi.fn().mockResolvedValue(mockRssResponse),
      } as unknown as Response);

      const response = await server.inject({
        method: 'GET',
        url: '/rss/spiegel/latest',
      });

      expect(response.statusCode).toBe(503);
      const data = JSON.parse(response.body);
      expect(data).toEqual({ error: 'Unable to fetch RSS feed' });
    });

    it('should handle CDATA in RSS titles', async () => {
      const mockRssResponse = `<?xml version="1.0" encoding="UTF-8"?>
        <rss version="2.0">
          <channel>
            <item>
              <title><![CDATA[Breaking News: Test]]></title>
              <link><![CDATA[https://example.com/breaking]]></link>
              <pubDate>Mon, 23 Sep 2024 12:00:00 +0200</pubDate>
            </item>
          </channel>
        </rss>`;

      (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: vi.fn().mockResolvedValue(mockRssResponse),
      } as unknown as Response);

      const response = await server.inject({
        method: 'GET',
        url: '/rss/spiegel/latest',
      });

      expect(response.statusCode).toBe(200);
      const data = JSON.parse(response.body);
      expect(data.title).toBe('Breaking News: Test');
      expect(data.link).toBe('https://example.com/breaking');
    });

    it('should use cache within TTL window', async () => {
      const mockRssResponse = `<?xml version="1.0" encoding="UTF-8"?>
        <rss version="2.0">
          <channel>
            <item>
              <title>Cached Headline</title>
              <link>https://example.com/cached</link>
              <pubDate>Mon, 23 Sep 2024 12:00:00 +0200</pubDate>
            </item>
          </channel>
        </rss>`;

      (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        text: vi.fn().mockResolvedValue(mockRssResponse),
      } as unknown as Response);

      // First request
      const response1 = await server.inject({
        method: 'GET',
        url: '/rss/spiegel/latest',
      });

      // Second request (should use cache)
      const response2 = await server.inject({
        method: 'GET',
        url: '/rss/spiegel/latest',
      });

      expect(response1.statusCode).toBe(200);
      expect(response2.statusCode).toBe(200);
      expect(global.fetch).toHaveBeenCalledTimes(1); // Only called once due to cache

      const data1 = JSON.parse(response1.body);
      const data2 = JSON.parse(response2.body);
      expect(data1).toEqual(data2);
    });

    it('should handle timeout correctly', async () => {
      const timeoutError = new Error('Request timeout');
      timeoutError.name = 'AbortError';

      (global.fetch as ReturnType<typeof vi.fn>).mockRejectedValue(timeoutError);

      const response = await server.inject({
        method: 'GET',
        url: '/rss/spiegel/latest',
      });

      expect(response.statusCode).toBe(503);
      const data = JSON.parse(response.body);
      expect(data).toEqual({ error: 'Unable to fetch RSS feed' });
    });
  });
});