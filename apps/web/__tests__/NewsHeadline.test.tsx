import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { NewsHeadline } from '@/components/NewsHeadline';

// Mock fetch globally
global.fetch = vi.fn();

// Mock window.innerWidth
Object.defineProperty(window, 'innerWidth', {
  writable: true,
  configurable: true,
  value: 1024,
});

describe('NewsHeadline Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('should display loading state initially', () => {
    vi.useFakeTimers();

    // Mock a pending promise that never resolves
    (global.fetch as ReturnType<typeof vi.fn>).mockImplementation(() =>
      new Promise(() => {}) // Never resolves to keep loading state
    );

    render(<NewsHeadline />);

    const loadingElement = screen.getByRole('status', { busy: true });
    expect(loadingElement).toBeInTheDocument();
    expect(loadingElement).toHaveAttribute('aria-busy', 'true');
    expect(loadingElement).toHaveAttribute('aria-live', 'polite');

    vi.useRealTimers();
  });

  it('should display headline when API returns data', async () => {
    const mockHeadline = {
      title: 'Breaking News: Test Headline',
      link: 'https://spiegel.de/article',
      publishedAt: new Date().toISOString(),
      source: 'SPIEGEL',
    };

    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue(mockHeadline),
    } as unknown as Response);

    render(<NewsHeadline />);

    await waitFor(() => {
      expect(screen.getByText(mockHeadline.title)).toBeInTheDocument();
    });

    const link = screen.getByRole('link', { name: /SPIEGEL: Breaking News: Test Headline/i });
    expect(link).toHaveAttribute('href', mockHeadline.link);
    expect(link).toHaveAttribute('target', '_blank');
    expect(link).toHaveAttribute('rel', 'noopener noreferrer');
  });

  it('should display error message when API fails', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Network error'));

    render(<NewsHeadline />);

    await waitFor(() => {
      expect(screen.getByText('Gerade keine Schlagzeile verfügbar')).toBeInTheDocument();
    });

    const statusElement = screen.getByRole('status');
    expect(statusElement).toHaveAttribute('aria-live', 'polite');
  });

  it('should display error message when API returns non-ok response', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 503,
    } as unknown as Response);

    render(<NewsHeadline />);

    await waitFor(() => {
      expect(screen.getByText('Gerade keine Schlagzeile verfügbar')).toBeInTheDocument();
    });
  });

  it('should include EILMELDUNG badge', async () => {
    const mockHeadline = {
      title: 'Breaking News',
      link: 'https://spiegel.de/article',
      publishedAt: new Date().toISOString(),
      source: 'SPIEGEL',
    };

    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue(mockHeadline),
    } as unknown as Response);

    render(<NewsHeadline />);

    await waitFor(() => {
      expect(screen.getByText('EILMELDUNG')).toBeInTheDocument();
    });

    const badge = screen.getByText('EILMELDUNG');
    expect(badge).toHaveClass('text-red-600');
    expect(badge).toHaveClass('bg-red-100');
  });

  it('should have proper accessibility attributes', async () => {
    const mockHeadline = {
      title: 'Accessible Headline',
      link: 'https://spiegel.de/article',
      publishedAt: new Date().toISOString(),
      source: 'SPIEGEL',
    };

    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue(mockHeadline),
    } as unknown as Response);

    render(<NewsHeadline />);

    await waitFor(() => {
      expect(screen.getByText('Accessible Headline')).toBeInTheDocument();
    });

    const container = screen.getByRole('status');
    expect(container).toHaveAttribute('aria-live', 'polite');
    expect(container).toHaveAttribute('aria-label', 'Aktuelle Schlagzeile von SPIEGEL');

    const link = screen.getByRole('link');
    expect(link).toHaveAttribute('aria-label', 'SPIEGEL: Accessible Headline');
  });
});