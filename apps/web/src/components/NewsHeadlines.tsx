'use client';

import { useEffect, useState, useCallback } from 'react';
import type { RssHeadline } from '@ts-template/shared';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3002';
const REFRESH_INTERVAL = 5 * 60 * 1000; // 5 minutes
const NEWS_LIMIT = 5;

type HeadlinesState = {
  loading: boolean;
  error: boolean;
  headlines: RssHeadline[];
};

const formatDateTime = (isoDate: string): string => {
  const date = new Date(isoDate);

  // Format for Europe/Berlin timezone
  const options: Intl.DateTimeFormatOptions = {
    timeZone: 'Europe/Berlin',
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  };

  const formatter = new Intl.DateTimeFormat('de-DE', options);
  const parts = formatter.formatToParts(date);

  const day = parts.find(p => p.type === 'day')?.value || '';
  const month = parts.find(p => p.type === 'month')?.value || '';
  const year = parts.find(p => p.type === 'year')?.value || '';
  const hour = parts.find(p => p.type === 'hour')?.value || '';
  const minute = parts.find(p => p.type === 'minute')?.value || '';

  return `${day}.${month}.${year} ${hour}:${minute}`;
};

const formatRelativeTime = (isoDate: string): string => {
  const date = new Date(isoDate);
  const now = new Date();
  const diffMinutes = Math.floor((now.getTime() - date.getTime()) / 60000);

  if (diffMinutes < 1) {
    return 'gerade eben';
  }
  if (diffMinutes === 1) {
    return 'vor 1 Minute';
  }
  if (diffMinutes < 60) {
    return `vor ${diffMinutes} Minuten`;
  }

  const diffHours = Math.floor(diffMinutes / 60);
  if (diffHours === 1) {
    return 'vor 1 Stunde';
  }
  if (diffHours < 24) {
    return `vor ${diffHours} Stunden`;
  }

  // For dates older than 24h, show the formatted date
  return formatDateTime(isoDate);
};

export function NewsHeadlines(): JSX.Element {
  const [state, setState] = useState<HeadlinesState>({
    loading: true,
    error: false,
    headlines: [],
  });

  const fetchHeadlines = useCallback(async () => {
    try {
      const response = await fetch(`${API_URL}/api/rss/spiegel/top5?limit=${NEWS_LIMIT}`);

      if (!response.ok) {
        throw new Error('Failed to fetch headlines');
      }

      const data = await response.json() as { headlines: RssHeadline[] };

      setState({
        loading: false,
        error: false,
        headlines: data.headlines,
      });
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Error fetching headlines:', error);
      setState({
        loading: false,
        error: true,
        headlines: [],
      });
    }
  }, []);

  useEffect(() => {
    fetchHeadlines();

    const intervalId = setInterval(fetchHeadlines, REFRESH_INTERVAL);

    return () => {
      clearInterval(intervalId);
    };
  }, [fetchHeadlines]);

  if (state.loading) {
    return (
      <div
        className="space-y-2"
        role="status"
        aria-live="polite"
        aria-busy="true"
        data-testid="news-headlines"
      >
        <div className="animate-pulse space-y-2">
          {Array.from({ length: 3 }).map((_, index) => (
            <div key={index} className="flex items-center justify-between">
              <div className="flex-1">
                <div className="h-4 bg-gray-300 dark:bg-gray-700 rounded w-3/4 mb-1"></div>
                <div className="h-3 bg-gray-200 dark:bg-gray-600 rounded w-32"></div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (state.error || state.headlines.length === 0) {
    return (
      <div
        className="text-gray-500 dark:text-gray-400 text-sm"
        role="status"
        aria-live="polite"
        data-testid="news-fallback"
      >
        Gerade keine Schlagzeilen verfügbar.
      </div>
    );
  }

  return (
    <div
      className="space-y-2"
      role="region"
      aria-label="Aktuelle SPIEGEL Schlagzeilen"
      aria-live="polite"
      data-testid="news-headlines"
    >
      <ul className="space-y-2">
        {state.headlines.map((headline, index) => (
          <li
            key={`${headline.link}-${index}`}
            className="flex items-start space-x-2 group"
            data-testid="news-headline-item"
          >
            <span className="text-red-600 dark:text-red-500 font-bold text-xs mt-0.5">
              {index + 1}.
            </span>
            <div className="flex-1 min-w-0">
              <a
                href={headline.link}
                target="_blank"
                rel="noopener noreferrer"
                className="block text-sm text-gray-900 dark:text-gray-100 hover:text-red-600 dark:hover:text-red-400 transition-colors duration-200 truncate"
                aria-label={`SPIEGEL: ${headline.title}`}
                title={headline.title}
                data-testid="headline-link"
              >
                <span className="font-medium" data-testid="headline-title">
                  {headline.title}
                </span>
              </a>
              <time
                className="text-xs text-gray-500 dark:text-gray-400"
                dateTime={headline.publishedAt}
                title={formatDateTime(headline.publishedAt)}
                data-testid="headline-date"
              >
                {formatRelativeTime(headline.publishedAt)}
              </time>
            </div>
          </li>
        ))}
      </ul>
      <div className="pt-2 border-t border-gray-200 dark:border-gray-700">
        <a
          href="https://www.spiegel.de"
          target="_blank"
          rel="noopener noreferrer"
          className="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
        >
          Mehr auf SPIEGEL.de →
        </a>
      </div>
    </div>
  );
}