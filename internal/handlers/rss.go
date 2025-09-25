package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/f00b455/golang-template/internal/config"
	"github.com/f00b455/golang-template/pkg/shared"
	"github.com/gin-gonic/gin"
)

const (
	cacheTTL       = 5 * time.Minute
	requestTimeout = 2 * time.Second
)

// RSSHandler handles RSS-related requests.
type RSSHandler struct {
	cfg        *config.Config
	cache      *cacheEntry
	multiCache *multiCacheEntry
	mu         sync.RWMutex
	httpClient *http.Client
}

type cacheEntry struct {
	data      *shared.RssHeadline
	timestamp time.Time
}

type multiCacheEntry struct {
	data      []shared.RssHeadline
	timestamp time.Time
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error" example:"Unable to fetch RSS feed"`
}

// HeadlinesResponse represents the response for multiple headlines.
type HeadlinesResponse struct {
	Headlines []shared.RssHeadline `json:"headlines"`
}

// NewRSSHandler creates a new RSSHandler.
func NewRSSHandler() *RSSHandler {
	return &RSSHandler{
		cfg:        config.Load(),
		cache:      &cacheEntry{},
		multiCache: &multiCacheEntry{},
		httpClient: &http.Client{Timeout: requestTimeout},
	}
}

// NewRSSHandlerWithClient creates a new RSSHandler with a custom HTTP client (for testing).
func NewRSSHandlerWithClient(client *http.Client) *RSSHandler {
	return &RSSHandler{
		cfg:        config.Load(),
		cache:      &cacheEntry{},
		multiCache: &multiCacheEntry{},
		httpClient: client,
	}
}

// GetLatest handles GET /api/rss/spiegel/latest
// @Summary      Get latest SPIEGEL RSS headline
// @Description  Fetches the most recent headline from SPIEGEL RSS feed
// @Tags         rss
// @Accept       json
// @Produce      json
// @Param        filter   query     string  false  "Filter headlines by keyword (case-insensitive)"
// @Success      200  {object}  shared.RssHeadline
// @Failure      503  {object}  ErrorResponse
// @Router       /rss/spiegel/latest [get]
func (h *RSSHandler) GetLatest(c *gin.Context) {
	filter := c.Query("filter")

	// Skip cache if filter is provided for now
	if filter == "" {
		h.mu.RLock()
		if h.cache.data != nil && time.Since(h.cache.timestamp) < cacheTTL {
			headline := *h.cache.data
			h.mu.RUnlock()
			c.JSON(http.StatusOK, headline)
			return
		}
		h.mu.RUnlock()
	}

	headline, err := h.fetchLatestHeadlineWithFilter(filter)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: "Unable to fetch RSS feed",
		})
		return
	}

	if headline == nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: "Unable to fetch RSS feed",
		})
		return
	}

	// Only cache unfiltered results
	if filter == "" {
		h.mu.Lock()
		h.cache = &cacheEntry{
			data:      headline,
			timestamp: time.Now(),
		}
		h.mu.Unlock()
	}

	c.JSON(http.StatusOK, *headline)
}

// GetTop5 handles GET /api/rss/spiegel/top5
// @Summary      Get top N SPIEGEL RSS headlines
// @Description  Fetches the top N headlines from SPIEGEL RSS feed (max 5)
// @Tags         rss
// @Accept       json
// @Produce      json
// @Param        limit    query     int     false  "Number of headlines to fetch (1-5)" minimum(1) maximum(5) default(5)
// @Param        filter   query     string  false  "Filter headlines by keyword (case-insensitive)"
// @Success      200      {object}  HeadlinesResponse
// @Failure      503      {object}  ErrorResponse
// @Router       /rss/spiegel/top5 [get]
func (h *RSSHandler) GetTop5(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 5 {
		limit = 5
	}

	filter := c.Query("filter")

	// Skip cache if filter is provided
	if filter == "" {
		h.mu.RLock()
		if len(h.multiCache.data) > 0 && time.Since(h.multiCache.timestamp) < cacheTTL {
			headlines := h.multiCache.data
			if len(headlines) > limit {
				headlines = headlines[:limit]
			}
			h.mu.RUnlock()
			c.JSON(http.StatusOK, HeadlinesResponse{Headlines: headlines})
			return
		}
		h.mu.RUnlock()
	}

	headlines, err := h.fetchMultipleHeadlinesWithFilter(limit, filter)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: "Unable to fetch RSS feed",
		})
		return
	}

	// Only cache unfiltered results
	if filter == "" && len(headlines) > 0 {
		h.mu.Lock()
		h.multiCache = &multiCacheEntry{
			data:      headlines,
			timestamp: time.Now(),
		}
		h.mu.Unlock()
	}

	c.JSON(http.StatusOK, HeadlinesResponse{Headlines: headlines})
}


func (h *RSSHandler) fetchRSSFeed() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", h.cfg.SpiegelRSSURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/rss+xml, application/xml, text/xml")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Golang-Template/1.0)")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("RSS fetch failed: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (h *RSSHandler) parseRSSItem(itemText string) (*shared.RssHeadline, error) {
	titleRegex := regexp.MustCompile(`<title>(.*?)</title>`)
	linkRegex := regexp.MustCompile(`<link>(.*?)</link>`)
	pubDateRegex := regexp.MustCompile(`<pubDate>([^<]+)</pubDate>`)

	titleMatches := titleRegex.FindStringSubmatch(itemText)
	linkMatches := linkRegex.FindStringSubmatch(itemText)

	if len(titleMatches) < 2 || len(linkMatches) < 2 {
		return nil, fmt.Errorf("required RSS fields not found")
	}

	title := h.cleanCDATA(titleMatches[1])
	link := h.cleanCDATA(linkMatches[1])

	publishedAt := time.Now().Format(time.RFC3339)
	if pubDateMatches := pubDateRegex.FindStringSubmatch(itemText); len(pubDateMatches) > 1 {
		if parsed, err := time.Parse(time.RFC1123Z, pubDateMatches[1]); err == nil {
			publishedAt = parsed.Format(time.RFC3339)
		}
	}

	return &shared.RssHeadline{
		Title:       title,
		Link:        link,
		PublishedAt: publishedAt,
		Source:      "SPIEGEL",
	}, nil
}

func (h *RSSHandler) parseMultipleRSSItems(rssText string, limit int) []shared.RssHeadline {
	var headlines []shared.RssHeadline

	itemRegex := regexp.MustCompile(`<item[^>]*>([\s\S]*?)</item>`)
	matches := itemRegex.FindAllStringSubmatch(rssText, -1) // Get all matches

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		headline, err := h.parseRSSItem(match[1])
		if err == nil && headline != nil {
			headlines = append(headlines, *headline)
			if limit > 0 && len(headlines) >= limit {
				break
			}
		}
	}

	return headlines
}

func (h *RSSHandler) cleanCDATA(text string) string {
	text = strings.ReplaceAll(text, "<![CDATA[", "")
	text = strings.ReplaceAll(text, "]]>", "")
	return strings.TrimSpace(text)
}

// fetchLatestHeadlineWithFilter fetches the latest headline matching the filter
func (h *RSSHandler) fetchLatestHeadlineWithFilter(filter string) (*shared.RssHeadline, error) {
	rssText, err := h.fetchRSSFeed()
	if err != nil {
		return nil, err
	}

	// Find all items in RSS feed
	itemRegex := regexp.MustCompile(`<item[^>]*>([\s\S]*?)</item>`)
	matches := itemRegex.FindAllStringSubmatch(rssText, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		headline, err := h.parseRSSItem(match[1])
		if err != nil || headline == nil {
			continue
		}

		// Apply filter if provided
		if filter != "" && !h.matchesFilter(headline.Title, filter) {
			continue
		}

		return headline, nil
	}

	return nil, fmt.Errorf("no RSS items found matching filter")
}

// fetchMultipleHeadlinesWithFilter fetches multiple headlines matching the filter
func (h *RSSHandler) fetchMultipleHeadlinesWithFilter(limit int, filter string) ([]shared.RssHeadline, error) {
	rssText, err := h.fetchRSSFeed()
	if err != nil {
		return nil, err
	}

	// Parse all items first
	allHeadlines := h.parseMultipleRSSItems(rssText, -1) // Get all items

	// Apply filter if provided
	if filter != "" {
		var filtered []shared.RssHeadline
		for _, headline := range allHeadlines {
			if h.matchesFilter(headline.Title, filter) {
				filtered = append(filtered, headline)
				if len(filtered) >= limit {
					break
				}
			}
		}
		return filtered, nil
	}

	// No filter, return up to limit items
	if len(allHeadlines) > limit {
		return allHeadlines[:limit], nil
	}
	return allHeadlines, nil
}

// matchesFilter checks if text contains the filter keyword (case-insensitive)
func (h *RSSHandler) matchesFilter(text, filter string) bool {
	if filter == "" {
		return true
	}
	return strings.Contains(strings.ToLower(text), strings.ToLower(filter))
}

// ResetCache resets both caches (for testing purposes).
func (h *RSSHandler) ResetCache() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.cache = &cacheEntry{}
	h.multiCache = &multiCacheEntry{}
}
