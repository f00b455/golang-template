package handlers

import (
	"bytes"
	"context"
	"encoding/csv"
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
	// maxFetchItems defines how many RSS items to fetch from the feed.
	// We fetch 250 items to ensure we have enough data for the 200 item limit,
	// while accounting for potential filtering. This provides a buffer for
	// filtered results while keeping memory usage manageable.
	maxFetchItems = 250
	// maxReturnItems defines the maximum number of items to return in the API response.
	// Increased to 200 to support displaying more news items in the terminal UI.
	maxReturnItems = 200
	// defaultReturnItems defines the default number of items when no limit is specified.
	// Kept at 5 for backward compatibility.
	defaultReturnItems = 5
	// maxFilterLength is the maximum allowed length for filter parameters to prevent DoS
	maxFilterLength = 100
	// maxExportItems is the maximum number of items allowed in export to prevent resource exhaustion
	maxExportItems = 1000
)

// RSSHandler handles RSS-related requests.
type RSSHandler struct {
	cfg         *config.Config
	cache       *cacheEntry
	multiCache  *multiCacheEntry
	mu          sync.RWMutex
	httpClient  *http.Client
	fetchMutex  sync.Mutex // Prevents concurrent RSS fetches
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
	Headlines  []shared.RssHeadline `json:"headlines"`
	TotalCount int                  `json:"totalCount,omitempty"`
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
// @Success      200  {object}  shared.RssHeadline
// @Failure      503  {object}  ErrorResponse
// @Router       /rss/spiegel/latest [get]
func (h *RSSHandler) GetLatest(c *gin.Context) {
	h.mu.RLock()
	if h.cache.data != nil && time.Since(h.cache.timestamp) < cacheTTL {
		headline := *h.cache.data
		h.mu.RUnlock()
		c.JSON(http.StatusOK, headline)
		return
	}
	h.mu.RUnlock()

	headline, err := h.fetchLatestHeadline()
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

	h.mu.Lock()
	h.cache = &cacheEntry{
		data:      headline,
		timestamp: time.Now(),
	}
	h.mu.Unlock()

	c.JSON(http.StatusOK, *headline)
}

// GetTop5 handles GET /api/rss/spiegel/top5
// @Summary      Get top N SPIEGEL RSS headlines
// @Description  Fetches the top N headlines from SPIEGEL RSS feed (max 200)
// @Tags         rss
// @Accept       json
// @Produce      json
// @Param        limit    query     int     false  "Number of headlines to fetch (1-200)" minimum(1) maximum(200) default(5)
// @Param        filter   query     string  false  "Filter headlines by keyword"
// @Success      200      {object}  HeadlinesResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      503      {object}  ErrorResponse
// @Router       /rss/spiegel/top5 [get]
func (h *RSSHandler) GetTop5(c *gin.Context) {
	limit := h.parseLimit(c)
	filterKeyword := c.Query("filter")

	// Validate filter parameter
	if err := h.validateFilter(filterKeyword); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Try to get headlines from cache
	headlines, totalCount := h.getCachedHeadlines()
	if headlines == nil {
		// Cache miss - fetch from RSS feed
		var err error
		headlines, err = h.fetchAndCacheHeadlines()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Error: "Unable to fetch RSS feed",
			})
			return
		}
		totalCount = len(headlines)
	}

	// Apply filter and limit
	headlines = h.applyFilterAndLimit(headlines, filterKeyword, limit)

	c.JSON(http.StatusOK, HeadlinesResponse{
		Headlines:  headlines,
		TotalCount: totalCount,
	})
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

func (h *RSSHandler) fetchMultipleHeadlines(limit int) ([]shared.RssHeadline, error) {
	rssText, err := h.fetchRSSFeed()
	if err != nil {
		return nil, err
	}

	return h.parseMultipleRSSItems(rssText, limit), nil
}

func (h *RSSHandler) fetchRSSFeed() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", h.cfg.SpiegelRSSURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/rss+xml, application/xml, text/xml")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Golang-Template/1.0)")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("request timeout after %v", requestTimeout)
		}
		return "", fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("RSS fetch failed with status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
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
	matches := h.extractRSSItems(rssText, limit)
	return h.processRSSMatches(matches, limit)
}

// extractRSSItems finds RSS item matches in the text
func (h *RSSHandler) extractRSSItems(rssText string, limit int) [][]string {
	itemRegex := regexp.MustCompile(`<item[^>]*>([\s\S]*?)</item>`)
	maxMatches := limit + (limit / 5) // Add 20% buffer for invalid items
	return itemRegex.FindAllStringSubmatch(rssText, maxMatches)
}

// processRSSMatches converts regex matches to RssHeadline objects
func (h *RSSHandler) processRSSMatches(matches [][]string, limit int) []shared.RssHeadline {
	headlines := make([]shared.RssHeadline, 0, limit)

	for i := 0; i < len(matches) && len(headlines) < limit; i++ {
		if len(matches[i]) < 2 {
			continue
		}

		if headline := h.parseItemSafe(matches[i][1]); headline != nil {
			headlines = append(headlines, *headline)
		}
	}

	return headlines
}

// parseItemSafe safely parses an RSS item, returning nil on error
func (h *RSSHandler) parseItemSafe(itemText string) *shared.RssHeadline {
	headline, err := h.parseRSSItem(itemText)
	if err != nil {
		return nil
	}
	return headline
}

func (h *RSSHandler) cleanCDATA(text string) string {
	text = strings.ReplaceAll(text, "<![CDATA[", "")
	text = strings.ReplaceAll(text, "]]>", "")
	return strings.TrimSpace(text)
}

// parseLimit extracts and validates the limit parameter from the request.
func (h *RSSHandler) parseLimit(c *gin.Context) int {
	limitStr := c.DefaultQuery("limit", strconv.Itoa(defaultReturnItems))
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return defaultReturnItems
	}
	if limit > maxReturnItems {
		return maxReturnItems
	}
	return limit
}

// validateFilter validates the filter parameter.
func (h *RSSHandler) validateFilter(filter string) error {
	if len(filter) > maxFilterLength {
		return fmt.Errorf("filter parameter too long (max %d characters)", maxFilterLength)
	}
	return nil
}

// getCachedHeadlines retrieves headlines from cache if available.
func (h *RSSHandler) getCachedHeadlines() ([]shared.RssHeadline, int) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.multiCache.data) > 0 && time.Since(h.multiCache.timestamp) < cacheTTL {
		// Return a copy to avoid race conditions
		headlines := make([]shared.RssHeadline, len(h.multiCache.data))
		copy(headlines, h.multiCache.data)
		return headlines, len(headlines)
	}
	return nil, 0
}

// fetchAndCacheHeadlines fetches headlines from RSS feed and updates the cache.
func (h *RSSHandler) fetchAndCacheHeadlines() ([]shared.RssHeadline, error) {
	// Prevent concurrent RSS fetches to avoid overwhelming the server
	h.fetchMutex.Lock()
	defer h.fetchMutex.Unlock()

	// Double-check cache after acquiring lock
	headlines, _ := h.getCachedHeadlines()
	if headlines != nil {
		return headlines, nil
	}

	// Fetch headlines from RSS feed
	headlines, err := h.fetchMultipleHeadlines(maxFetchItems)
	if err != nil || len(headlines) == 0 {
		return nil, err
	}

	// Make a copy to avoid data races when reading from cache
	headlinesCopy := make([]shared.RssHeadline, len(headlines))
	copy(headlinesCopy, headlines)

	h.mu.Lock()
	h.multiCache = &multiCacheEntry{
		data:      headlinesCopy,
		timestamp: time.Now(),
	}
	h.mu.Unlock()

	return headlines, nil
}

// applyFilterAndLimit applies the filter keyword and limit to headlines.
func (h *RSSHandler) applyFilterAndLimit(headlines []shared.RssHeadline, filter string, limit int) []shared.RssHeadline {
	// Early return for common case
	if filter == "" && len(headlines) <= limit {
		return headlines
	}

	if filter != "" {
		headlines = h.filterHeadlines(headlines, filter)
	}
	if len(headlines) > limit {
		headlines = headlines[:limit]
	}
	return headlines
}

// filterHeadlines filters headlines based on a keyword (case-insensitive).
func (h *RSSHandler) filterHeadlines(headlines []shared.RssHeadline, keyword string) []shared.RssHeadline {
	if keyword == "" {
		return headlines
	}

	keyword = strings.ToLower(keyword)
	// Pre-allocate with estimated capacity (assuming ~30% match rate)
	estimatedCapacity := len(headlines) / 3
	if estimatedCapacity < 1 {
		estimatedCapacity = 1
	}
	filtered := make([]shared.RssHeadline, 0, estimatedCapacity)

	for _, headline := range headlines {
		if strings.Contains(strings.ToLower(headline.Title), keyword) {
			filtered = append(filtered, headline)
		}
	}

	return filtered
}

// ExportHeadlines handles GET /api/rss/spiegel/export
// @Summary      Export SPIEGEL RSS headlines
// @Description  Exports RSS headlines in CSV or JSON format
// @Tags         rss
// @Accept       json
// @Produce      json
// @Produce      text/csv
// @Param        format   query     string  true   "Export format (json or csv)"
// @Param        filter   query     string  false  "Filter headlines by keyword"
// @Param        limit    query     int     false  "Number of headlines to export (1-1000)" minimum(1) maximum(1000)
// @Success      200      {object}  object
// @Failure      400      {object}  ErrorResponse
// @Failure      503      {object}  ErrorResponse
// @Router       /rss/spiegel/export [get]
// validateExportFormat checks if the export format is valid
func (h *RSSHandler) validateExportFormat(format string) error {
	if format == "" {
		return fmt.Errorf("missing format parameter")
	}
	if format != "json" && format != "csv" {
		return fmt.Errorf("invalid format parameter: must be 'json' or 'csv'")
	}
	return nil
}

// prepareExportData fetches and filters headlines for export
func (h *RSSHandler) prepareExportData(filterKeyword string, limit int) ([]shared.RssHeadline, error) {
	headlines, _ := h.getCachedHeadlines()
	if headlines == nil {
		var err error
		headlines, err = h.fetchAndCacheHeadlines()
		if err != nil {
			return nil, err
		}
	}

	// Apply filter
	if filterKeyword != "" {
		headlines = h.filterHeadlines(headlines, filterKeyword)
	}

	// Apply limit
	if limit > 0 && len(headlines) > limit {
		headlines = headlines[:limit]
	}

	return headlines, nil
}

// generateExportFilename creates a filename for export with optional filter
func (h *RSSHandler) generateExportFilename(format, filter string) string {
	timestamp := time.Now().Format("20060102_150405")
	if filter != "" {
		return fmt.Sprintf("rss_export_%s_%s.%s", filter, timestamp, format)
	}
	return fmt.Sprintf("rss_export_%s.%s", timestamp, format)
}

func (h *RSSHandler) ExportHeadlines(c *gin.Context) {
	params, err := h.validateExportParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	headlines, err := h.prepareExportData(params.filter, params.limit)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{Error: "Unable to fetch RSS feed"})
		return
	}

	h.performExport(c, headlines, params)
}

// exportParams holds validated export parameters
type exportParams struct {
	format string
	filter string
	limit  int
}

// validateExportParams validates all export parameters
func (h *RSSHandler) validateExportParams(c *gin.Context) (*exportParams, error) {
	format := c.Query("format")
	if err := h.validateExportFormat(format); err != nil {
		return nil, err
	}

	filter := c.Query("filter")
	if err := h.validateFilter(filter); err != nil {
		return nil, err
	}

	limit, err := h.validateAndParseExportLimit(c)
	if err != nil {
		return nil, err
	}

	return &exportParams{
		format: format,
		filter: filter,
		limit:  limit,
	}, nil
}

// validateAndParseExportLimit validates and parses the export limit
func (h *RSSHandler) validateAndParseExportLimit(c *gin.Context) (int, error) {
	limitStr := c.Query("limit")
	if limitStr == "" {
		return maxExportItems, nil
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return maxExportItems, nil
	}

	if limit > maxExportItems {
		return 0, fmt.Errorf("limit exceeds maximum allowed value of %d", maxExportItems)
	}

	return limit, nil
}

// performExport executes the actual export based on format
func (h *RSSHandler) performExport(c *gin.Context, headlines []shared.RssHeadline, params *exportParams) {
	filename := h.generateExportFilename(params.format, params.filter)

	if params.format == "json" {
		h.exportAsJSON(c, headlines, params.filter, filename)
	} else {
		h.exportAsCSV(c, headlines, filename)
	}
}

func (h *RSSHandler) exportAsJSON(c *gin.Context, headlines []shared.RssHeadline, filter, filename string) {
	response := struct {
		ExportDate    string               `json:"export_date"`
		TotalItems    int                  `json:"total_items"`
		FilterApplied string               `json:"filter_applied,omitempty"`
		Headlines     []shared.RssHeadline `json:"headlines"`
	}{
		ExportDate: time.Now().Format(time.RFC3339),
		TotalItems: len(headlines),
		Headlines:  headlines,
	}

	if filter != "" {
		response.FilterApplied = filter
	}

	// Set security headers
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy", "default-src 'none'")
	c.JSON(http.StatusOK, response)
}

func (h *RSSHandler) exportAsCSV(c *gin.Context, headlines []shared.RssHeadline, filename string) {
	// Build CSV content in memory to calculate Content-Length
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{"Title", "Link", "Published_At", "Source"}
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to write CSV headers",
		})
		return
	}

	// Write data rows with sanitization
	for _, headline := range headlines {
		row := []string{
			h.sanitizeCSVField(headline.Title),
			h.sanitizeCSVField(headline.Link),
			h.sanitizeCSVField(headline.PublishedAt),
			h.sanitizeCSVField(headline.Source),
		}
		if err := writer.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to write CSV row",
			})
			return
		}
	}

	writer.Flush()

	// Check for any errors in CSV writer
	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to generate CSV",
		})
		return
	}

	// Set headers including Content-Length
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy", "default-src 'none'")

	// Write the response
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}


// sanitizeCSVField protects against CSV injection by sanitizing field values.
// It prefixes potentially dangerous characters with a single quote to neutralize
// formula injection attempts.
func (h *RSSHandler) sanitizeCSVField(field string) string {
	if field == "" {
		return field
	}

	// Check if the field starts with a potentially dangerous character
	// These characters can trigger formula execution in spreadsheet applications
	dangerousChars := []rune{'=', '+', '-', '@', '\t', '\r'}
	firstChar := rune(field[0])

	for _, dangerous := range dangerousChars {
		if firstChar == dangerous {
			// Prefix with single quote to neutralize formula injection
			return "'" + field
		}
	}

	return field
}

// ResetCache resets both caches (for testing purposes).
func (h *RSSHandler) ResetCache() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.cache = &cacheEntry{}
	h.multiCache = &multiCacheEntry{}
}
