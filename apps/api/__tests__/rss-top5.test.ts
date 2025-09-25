import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { FastifyInstance } from 'fastify';
import createServer from '../src/server';
import { resetCache } from '../src/routes/rss';

describe('RSS Top 5 Headlines API', () => {
  let server: FastifyInstance;

  beforeEach(async () => {
    server = await createServer();
    await server.ready();
    resetCache();
  });

  afterEach(async () => {
    await server.close();
    vi.clearAllMocks();
  });

  describe('GET /api/rss/spiegel/top5', () => {
    it('should return up to 5 headlines', async () => {
      const mockRssResponse = `
        <?xml version="1.0" encoding="UTF-8"?>
        <rss version="2.0">
          <channel>
            <title>SPIEGEL ONLINE</title>
            <item>
              <title><![CDATA[Headline 1]]></title>
              <link><![CDATA[https://www.spiegel.de/1]]></link>
              <pubDate>Mon, 24 Sep 2025 10:00:00 +0000</pubDate>
            </item>
            <item>
              <title><![CDATA[Headline 2]]></title>
              <link><![CDATA[https://www.spiegel.de/2]]></link>
              <pubDate>Mon, 24 Sep 2025 09:00:00 +0000</pubDate>
            </item>
            <item>
              <title><![CDATA[Headline 3]]></title>
              <link><![CDATA[https://www.spiegel.de/3]]></link>
              <pubDate>Mon, 24 Sep 2025 08:00:00 +0000</pubDate>
            </item>
            <item>
              <title><![CDATA[Headline 4]]></title>
              <link><![CDATA[https://www.spiegel.de/4]]></link>
              <pubDate>Mon, 24 Sep 2025 07:00:00 +0000</pubDate>
            </item>
            <item>
              <title><![CDATA[Headline 5]]></title>
              <link><![CDATA[https://www.spiegel.de/5]]></link>
              <pubDate>Mon, 24 Sep 2025 06:00:00 +0000</pubDate>
            </item>
            <item>
              <title><![CDATA[Headline 6]]></title>
              <link><![CDATA[https://www.spiegel.de/6]]></link>
              <pubDate>Mon, 24 Sep 2025 05:00:00 +0000</pubDate>
            </item>
          </channel>
        </rss>
      `;

      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        text: async () => mockRssResponse,
      } as Response);

      const response = await server.inject({
        method: 'GET',
        url: '/api/rss/spiegel/top5',
      });

      expect(response.statusCode).toBe(200);
      const data = JSON.parse(response.body);

      expect(data).toHaveProperty('headlines');
      expect(data.headlines).toHaveLength(5);
      expect(data.headlines[0]).toHaveProperty('title', 'Headline 1');
      expect(data.headlines[0]).toHaveProperty('link', 'https://www.spiegel.de/1');
      expect(data.headlines[0]).toHaveProperty('source', 'SPIEGEL');
      expect(data.headlines[0]).toHaveProperty('publishedAt');
    });

    it('should respect limit parameter', async () => {
      const mockRssResponse = `
        <?xml version="1.0" encoding="UTF-8"?>
        <rss version="2.0">
          <channel>
            <title>SPIEGEL ONLINE</title>
            <item>
              <title><![CDATA[Headline 1]]></title>
              <link><![CDATA[https://www.spiegel.de/1]]></link>
              <pubDate>Mon, 24 Sep 2025 10:00:00 +0000</pubDate>
            </item>
            <item>
              <title><![CDATA[Headline 2]]></title>
              <link><![CDATA[https://www.spiegel.de/2]]></link>
              <pubDate>Mon, 24 Sep 2025 09:00:00 +0000</pubDate>
            </item>
            <item>
              <title><![CDATA[Headline 3]]></title>
              <link><![CDATA[https://www.spiegel.de/3]]></link>
              <pubDate>Mon, 24 Sep 2025 08:00:00 +0000</pubDate>
            </item>
          </channel>
        </rss>
      `;

      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        text: async () => mockRssResponse,
      } as Response);

      const response = await server.inject({
        method: 'GET',
        url: '/api/rss/spiegel/top5?limit=2',
      });

      expect(response.statusCode).toBe(200);
      const data = JSON.parse(response.body);

      expect(data.headlines).toHaveLength(2);
    });

    it('should handle fewer than 5 items', async () => {
      const mockRssResponse = `
        <?xml version="1.0" encoding="UTF-8"?>
        <rss version="2.0">
          <channel>
            <title>SPIEGEL ONLINE</title>
            <item>
              <title><![CDATA[Headline 1]]></title>
              <link><![CDATA[https://www.spiegel.de/1]]></link>
              <pubDate>Mon, 24 Sep 2025 10:00:00 +0000</pubDate>
            </item>
            <item>
              <title><![CDATA[Headline 2]]></title>
              <link><![CDATA[https://www.spiegel.de/2]]></link>
              <pubDate>Mon, 24 Sep 2025 09:00:00 +0000</pubDate>
            </item>
          </channel>
        </rss>
      `;

      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        text: async () => mockRssResponse,
      } as Response);

      const response = await server.inject({
        method: 'GET',
        url: '/api/rss/spiegel/top5',
      });

      expect(response.statusCode).toBe(200);
      const data = JSON.parse(response.body);

      expect(data.headlines).toHaveLength(2);
    });

    it('should return error when RSS fetch fails', async () => {
      global.fetch = vi.fn().mockRejectedValue(new Error('Network error'));

      const response = await server.inject({
        method: 'GET',
        url: '/api/rss/spiegel/top5',
      });

      expect(response.statusCode).toBe(503);
      const data = JSON.parse(response.body);
      expect(data).toHaveProperty('error', 'Unable to fetch RSS feed');
    });

    it('should use cache for subsequent requests', async () => {
      const mockRssResponse = `
        <?xml version="1.0" encoding="UTF-8"?>
        <rss version="2.0">
          <channel>
            <title>SPIEGEL ONLINE</title>
            <item>
              <title><![CDATA[Cached Headline]]></title>
              <link><![CDATA[https://www.spiegel.de/cached]]></link>
              <pubDate>Mon, 24 Sep 2025 10:00:00 +0000</pubDate>
            </item>
          </channel>
        </rss>
      `;

      const fetchMock = vi.fn().mockResolvedValue({
        ok: true,
        text: async () => mockRssResponse,
      } as Response);
      global.fetch = fetchMock;

      // First request
      const response1 = await server.inject({
        method: 'GET',
        url: '/api/rss/spiegel/top5',
      });

      expect(response1.statusCode).toBe(200);
      expect(fetchMock).toHaveBeenCalledTimes(1);

      // Second request should use cache
      const response2 = await server.inject({
        method: 'GET',
        url: '/api/rss/spiegel/top5',
      });

      expect(response2.statusCode).toBe(200);
      expect(fetchMock).toHaveBeenCalledTimes(1); // Still only 1 call

      const data1 = JSON.parse(response1.body);
      const data2 = JSON.parse(response2.body);
      expect(data1).toEqual(data2);
    });
  });
});