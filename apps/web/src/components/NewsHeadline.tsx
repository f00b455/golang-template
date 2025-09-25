'use client';

import { useEffect, useState, useCallback } from 'react';
import type { RssHeadline } from '@ts-template/shared';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3002';
const REFRESH_INTERVAL = 5 * 60 * 1000; // 5 minutes

type HeadlineState = {
  loading: boolean;
  error: boolean;
  headline: RssHeadline | null;
};

const formatTime = (isoDate: string): string => {
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

  return date.toLocaleDateString('de-DE');
};

const truncateText = (text: string, maxLength: number): string => {
  if (text.length <= maxLength) {
    return text;
  }
  return text.substring(0, maxLength - 3) + '...';
};

export function NewsHeadline(): JSX.Element {
  const [state, setState] = useState<HeadlineState>({
    loading: true,
    error: false,
    headline: null,
  });

  const fetchHeadline = useCallback(async () => {
    try {
      const response = await fetch(`${API_URL}/api/rss/spiegel/latest`);

      if (!response.ok) {
        throw new Error('Failed to fetch headline');
      }

      const data = await response.json() as RssHeadline;

      setState({
        loading: false,
        error: false,
        headline: data,
      });
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Error fetching headline:', error);
      setState({
        loading: false,
        error: true,
        headline: null,
      });
    }
  }, []);

  useEffect(() => {
    fetchHeadline();

    const intervalId = setInterval(fetchHeadline, REFRESH_INTERVAL);

    return () => {
      clearInterval(intervalId);
    };
  }, [fetchHeadline]);

  if (state.loading) {
    return (
      <div
        className="flex items-center space-x-2 text-gray-500"
        role="status"
        aria-live="polite"
        aria-busy="true"
      >
        <div className="animate-pulse flex items-center space-x-2">
          <div className="h-4 w-4 bg-gray-300 rounded-full"></div>
          <div className="h-4 bg-gray-300 rounded w-64"></div>
        </div>
      </div>
    );
  }

  if (state.error || !state.headline) {
    return (
      <div
        className="text-gray-400 text-sm"
        role="status"
        aria-live="polite"
      >
        Gerade keine Schlagzeile verfügbar
      </div>
    );
  }

  const { title, link, publishedAt, source } = state.headline;
  const displayTitle = typeof window !== 'undefined' && window.innerWidth < 640
    ? truncateText(title, 50)
    : title;

  return (
    <div
      className="flex items-center space-x-2"
      role="status"
      aria-live="polite"
      aria-label={`Aktuelle Schlagzeile von ${source}`}
    >
      <span className="inline-flex items-center px-2 py-1 text-xs font-semibold text-red-600 bg-red-100 rounded">
        EILMELDUNG
      </span>
      <a
        href={link}
        target="_blank"
        rel="noopener noreferrer"
        className="text-gray-900 hover:text-red-600 transition-colors duration-200 font-medium text-sm sm:text-base truncate max-w-xs sm:max-w-md lg:max-w-lg xl:max-w-xl"
        aria-label={`${source}: ${title}`}
        title={title}
      >
        {displayTitle}
      </a>
      <span className="text-gray-400 text-xs whitespace-nowrap hidden sm:inline">
        {formatTime(publishedAt)}
      </span>
    </div>
  );
}