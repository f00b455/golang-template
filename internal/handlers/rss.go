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
	cache      map[string]*cacheEntry
	multiCache map[string]*multiCacheEntry
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
		cache:      make(map[string]*cacheEntry),
		multiCache: make(map[string]*multiCacheEntry),
		httpClient: &http.Client{Timeout: requestTimeout},
	}
}

// NewRSSHandlerWithClient creates a new RSSHandler with a custom HTTP client (for testing).
func NewRSSHandlerWithClient(client *http.Client) *RSSHandler {
	return &RSSHandler{
		cfg:        config.Load(),
		cache:      make(map[string]*cacheEntry),
		multiCache: make(map[string]*multiCacheEntry),
		httpClient: client,
	}
}

// GetLatest handles GET /api/rss/spiegel/latest
// @Summary      Get latest SPIEGEL RSS headline
// @Description  Fetches the most recent headline from SPIEGEL RSS feed
// @Tags         rss
// @Accept       json
// @Produce      json
// @Param        filter    query     string  false  "Filter headlines by text (case-insensitive)"
// @Success      200  {object}  shared.RssHeadline
// @Failure      503  {object}  ErrorResponse
// @Router       /rss/spiegel/latest [get]
func (h *RSSHandler) GetLatest(c *gin.Context) {
	filter := c.Query("filter")
	cacheKey := "latest:" + filter

	h.mu.RLock()
	if cached, ok := h.cache[cacheKey]; ok && cached != nil && cached.data != nil && time.Since(cached.timestamp) < cacheTTL {
		headline := *cached.data
		h.mu.RUnlock()
		c.JSON(http.StatusOK, headline)
		return
	}
	h.mu.RUnlock()

	headline, err := h.fetchLatestHeadlineWithFilter(filter)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: "Unable to fetch RSS feed",
		})
		return
	}

	if headline == nil {
		// Return empty response for no match
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	if headline != nil {
		h.mu.Lock()
		h.cache[cacheKey] = &cacheEntry{
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
// @Param        filter   query     string  false  "Filter headlines by text (case-insensitive)"
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
	cacheKey := fmt.Sprintf("top%d:%s", limit, filter)

	h.mu.RLock()
	if cached, ok := h.multiCache[cacheKey]; ok && cached != nil && len(cached.data) > 0 && time.Since(cached.timestamp) < cacheTTL {
		headlines := cached.data
		h.mu.RUnlock()
		c.JSON(http.StatusOK, HeadlinesResponse{Headlines: headlines})
		return
	}
	h.mu.RUnlock()

	headlines, err := h.fetchMultipleHeadlinesWithFilter(limit, filter)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: "Unable to fetch RSS feed",
		})
		return
	}

	h.mu.Lock()
	h.multiCache[cacheKey] = &multiCacheEntry{
		data:      headlines,
		timestamp: time.Now(),
	}
	h.mu.Unlock()

	c.JSON(http.StatusOK, HeadlinesResponse{Headlines: headlines})
}

func (h *RSSHandler) fetchLatestHeadline() (*shared.RssHeadline, error) {
	rssText, err := h.fetchRSSFeed()
	if err != nil {
		return nil, err
	}

	// Find first item in RSS feed
	itemRegex := regexp.MustCompile(`<item[^>]*>([\s\S]*?)</item>`)
	matches := itemRegex.FindStringSubmatch(rssText)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no RSS items found")
	}

	return h.parseRSSItem(matches[1])
}

func (h *RSSHandler) fetchLatestHeadlineWithFilter(filter string) (*shared.RssHeadline, error) {
	if filter == "" {
		return h.fetchLatestHeadline()
	}

	rssText, err := h.fetchRSSFeed()
	if err != nil {
		return nil, err
	}

	// Find all items in RSS feed
	itemRegex := regexp.MustCompile(`<item[^>]*>([\s\S]*?)</item>`)
	matches := itemRegex.FindAllStringSubmatch(rssText, -1)

	// Filter and return first matching item
	filterLower := strings.ToLower(filter)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		headline, err := h.parseRSSItem(match[1])
		if err != nil || headline == nil {
			continue
		}
		if strings.Contains(strings.ToLower(headline.Title), filterLower) {
			return headline, nil
		}
	}

	return nil, nil // No matching headline
}

func (h *RSSHandler) fetchMultipleHeadlines(limit int) ([]shared.RssHeadline, error) {
	rssText, err := h.fetchRSSFeed()
	if err != nil {
		return nil, err
	}

	return h.parseMultipleRSSItems(rssText, limit), nil
}

func (h *RSSHandler) fetchMultipleHeadlinesWithFilter(limit int, filter string) ([]shared.RssHeadline, error) {
	if filter == "" {
		return h.fetchMultipleHeadlines(limit)
	}

	rssText, err := h.fetchRSSFeed()
	if err != nil {
		return nil, err
	}

	var headlines []shared.RssHeadline
	filterLower := strings.ToLower(filter)

	itemRegex := regexp.MustCompile(`<item[^>]*>([\s\S]*?)</item>`)
	// Get all items to filter through
	matches := itemRegex.FindAllStringSubmatch(rssText, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		headline, err := h.parseRSSItem(match[1])
		if err == nil && headline != nil {
			// Apply filter (case-insensitive)
			if strings.Contains(strings.ToLower(headline.Title), filterLower) {
				headlines = append(headlines, *headline)
				if len(headlines) >= limit {
					break
				}
			}
		}
	}

	return headlines, nil
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
	matches := itemRegex.FindAllStringSubmatch(rssText, limit)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		headline, err := h.parseRSSItem(match[1])
		if err == nil && headline != nil {
			headlines = append(headlines, *headline)
			if len(headlines) >= limit {
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

// ResetCache resets both caches (for testing purposes).
func (h *RSSHandler) ResetCache() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.cache = make(map[string]*cacheEntry)
	h.multiCache = make(map[string]*multiCacheEntry)
}
