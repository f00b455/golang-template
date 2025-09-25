import { FastifyPluginAsync } from 'fastify';
import type { RssHeadline } from '@ts-template/shared';

const SPIEGEL_RSS_URL = process.env.SPIEGEL_RSS_URL || 'https://www.spiegel.de/schlagzeilen/index.rss';
const CACHE_TTL = 5 * 60 * 1000; // 5 minutes in milliseconds
const REQUEST_TIMEOUT = 2000; // 2 seconds

let cache: { data: RssHeadline | null; timestamp: number } = { data: null, timestamp: 0 };
let multiCache: { data: RssHeadline[]; timestamp: number } = { data: [], timestamp: 0 };

// Export function to reset cache for testing
export const resetCache = (): void => {
  cache = { data: null, timestamp: 0 };
  multiCache = { data: [], timestamp: 0 };
};

const parseRssItem = (itemText: string): RssHeadline | null => {
  const titleMatch = itemText.match(/<title>(.*?)<\/title>/s);
  const linkMatch = itemText.match(/<link>(.*?)<\/link>/s);
  const pubDateMatch = itemText.match(/<pubDate>([^<]+)<\/pubDate>/);

  if (!titleMatch || !linkMatch) {
    return null;
  }

  const title = titleMatch[1]
    ?.replace(/<!\[CDATA\[/, '')
    ?.replace(/\]\]>/, '')
    ?.trim() ?? '';

  const link = linkMatch[1]
    ?.replace(/<!\[CDATA\[/, '')
    ?.replace(/\]\]>/, '')
    ?.trim() ?? '';

  const publishedAt = pubDateMatch?.[1]
    ? new Date(pubDateMatch[1]).toISOString()
    : new Date().toISOString();

  return {
    title,
    link,
    publishedAt,
    source: 'SPIEGEL'
  };
};

const parseMultipleRssItems = (rssText: string, limit: number = 5): RssHeadline[] => {
  const headlines: RssHeadline[] = [];
  const itemMatches = rssText.matchAll(/<item[^>]*>([\s\S]*?)<\/item>/g);

  for (const match of itemMatches) {
    if (headlines.length >= limit) {
      break;
    }

    const itemText = match[1];
    if (itemText) {
      const headline = parseRssItem(itemText);
      if (headline) {
        headlines.push(headline);
      }
    }
  }

  return headlines;
};

const fetchLatestHeadline = async (): Promise<RssHeadline | null> => {
  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), REQUEST_TIMEOUT);

    const response = await fetch(SPIEGEL_RSS_URL, {
      signal: controller.signal,
      headers: {
        'Accept': 'application/rss+xml, application/xml, text/xml',
        'User-Agent': 'Mozilla/5.0 (compatible; TypeScript-Template/1.0)'
      }
    });

    clearTimeout(timeoutId);

    if (!response.ok) {
      throw new Error(`RSS fetch failed: ${response.status}`);
    }

    const rssText = await response.text();

    // Extract first item from RSS feed
    const firstItemMatch = rssText.match(/<item[^>]*>([\s\S]*?)<\/item>/);

    if (!firstItemMatch) {
      return null;
    }

    return parseRssItem(firstItemMatch[1] ?? '');
  } catch (error) {
    console.error('Failed to fetch RSS:', error);
    return null;
  }
};

const fetchMultipleHeadlines = async (limit: number = 5): Promise<RssHeadline[]> => {
  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), REQUEST_TIMEOUT);

    const response = await fetch(SPIEGEL_RSS_URL, {
      signal: controller.signal,
      headers: {
        'Accept': 'application/rss+xml, application/xml, text/xml',
        'User-Agent': 'Mozilla/5.0 (compatible; TypeScript-Template/1.0)'
      }
    });

    clearTimeout(timeoutId);

    if (!response.ok) {
      throw new Error(`RSS fetch failed: ${response.status}`);
    }

    const rssText = await response.text();
    return parseMultipleRssItems(rssText, limit);
  } catch (error) {
    console.error('Failed to fetch RSS:', error);
    return [];
  }
};

export const rssRoute: FastifyPluginAsync = async function (fastify) {
  fastify.get<{ Reply: RssHeadline | { error: string } }>(
    '/rss/spiegel/latest',
    {
      schema: {
        description: 'Get the latest SPIEGEL RSS headline',
        tags: ['rss'],
        response: {
          200: {
            type: 'object',
            properties: {
              title: {
                type: 'string',
                description: 'Headline title',
              },
              link: {
                type: 'string',
                description: 'URL to the article',
              },
              publishedAt: {
                type: 'string',
                description: 'ISO 8601 date string',
              },
              source: {
                type: 'string',
                description: 'News source',
              },
            },
            required: ['title', 'link', 'publishedAt', 'source'],
          },
          503: {
            type: 'object',
            properties: {
              error: {
                type: 'string',
                description: 'Error message',
              },
            },
          },
        },
      },
    },
    async function (_request, reply) {
      // Check cache first
      const now = Date.now();
      if (cache.data && (now - cache.timestamp) < CACHE_TTL) {
        return cache.data;
      }

      // Fetch fresh data
      const headline = await fetchLatestHeadline();

      if (!headline) {
        reply.code(503);
        return { error: 'Unable to fetch RSS feed' };
      }

      // Update cache
      cache = { data: headline, timestamp: now };

      return headline;
    }
  );

  fastify.get<{
    Querystring: { limit?: string };
    Reply: { headlines: RssHeadline[] } | { error: string };
  }>(
    '/rss/spiegel/top5',
    {
      schema: {
        description: 'Get the top N SPIEGEL RSS headlines (max 5)',
        tags: ['rss'],
        querystring: {
          type: 'object',
          properties: {
            limit: {
              type: 'string',
              description: 'Number of headlines to fetch (1-5)',
              pattern: '^[1-5]$'
            },
          },
        },
        response: {
          200: {
            type: 'object',
            properties: {
              headlines: {
                type: 'array',
                items: {
                  type: 'object',
                  properties: {
                    title: {
                      type: 'string',
                      description: 'Headline title',
                    },
                    link: {
                      type: 'string',
                      description: 'URL to the article',
                    },
                    publishedAt: {
                      type: 'string',
                      description: 'ISO 8601 date string',
                    },
                    source: {
                      type: 'string',
                      description: 'News source',
                    },
                  },
                  required: ['title', 'link', 'publishedAt', 'source'],
                },
              },
            },
            required: ['headlines'],
          },
          503: {
            type: 'object',
            properties: {
              error: {
                type: 'string',
                description: 'Error message',
              },
            },
          },
        },
      },
    },
    async function (request, reply) {
      const limit = request.query.limit ? parseInt(request.query.limit, 10) : 5;

      // Check cache first
      const now = Date.now();
      if (multiCache.data.length > 0 && (now - multiCache.timestamp) < CACHE_TTL) {
        return { headlines: multiCache.data.slice(0, limit) };
      }

      // Fetch fresh data
      const headlines = await fetchMultipleHeadlines(5);

      if (headlines.length === 0) {
        reply.code(503);
        return { error: 'Unable to fetch RSS feed' };
      }

      // Update cache
      multiCache = { data: headlines, timestamp: now };

      return { headlines: headlines.slice(0, limit) };
    }
  );
};