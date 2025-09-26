package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/f00b455/golang-template/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetTop200Items tests the endpoint with 200 items limit
func TestGetTop200Items(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		limit          string
		expectedItems  int
		expectedStatus int
		description    string
	}{
		{
			name:           "Request 200 items",
			limit:          "200",
			expectedItems:  200,
			expectedStatus: http.StatusOK,
			description:    "Should return up to 200 items when limit=200",
		},
		{
			name:           "Request 100 items",
			limit:          "100",
			expectedItems:  100,
			expectedStatus: http.StatusOK,
			description:    "Should return 100 items when limit=100",
		},
		{
			name:           "Request 50 items",
			limit:          "50",
			expectedItems:  50,
			expectedStatus: http.StatusOK,
			description:    "Should return 50 items when limit=50",
		},
		{
			name:           "Default without limit",
			limit:          "",
			expectedItems:  5,
			expectedStatus: http.StatusOK,
			description:    "Should maintain backward compatibility with 5 items by default",
		},
		{
			name:           "Request more than max (201)",
			limit:          "201",
			expectedItems:  200,
			expectedStatus: http.StatusOK,
			description:    "Should cap at 200 items when requesting more",
		},
		{
			name:           "Request negative number",
			limit:          "-1",
			expectedItems:  5,
			expectedStatus: http.StatusOK,
			description:    "Should return default 5 items for invalid negative limit",
		},
		{
			name:           "Request with invalid string",
			limit:          "abc",
			expectedItems:  5,
			expectedStatus: http.StatusOK,
			description:    "Should return default 5 items for invalid string limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock HTTP client with large RSS feed
			mockClient := testutil.CreateMockHTTPClient(t, generateLargeRSSFeed(250))
			handler := NewRSSHandlerWithClient(mockClient)

			// Setup request
			router := gin.New()
			router.GET("/api/rss/spiegel/top5", handler.GetTop5)

			url := "/api/rss/spiegel/top5"
			if tt.limit != "" {
				url = fmt.Sprintf("%s?limit=%s", url, tt.limit)
			}

			req, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			// Execute request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify response
			assert.Equal(t, tt.expectedStatus, w.Code, tt.description)

			if w.Code == http.StatusOK {
				var response HeadlinesResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// Check number of items returned
				actualItems := len(response.Headlines)
				assert.LessOrEqual(t, actualItems, tt.expectedItems,
					"Should not return more than requested items")

				// For smaller limits, should match exactly if feed has enough items
				if tt.expectedItems <= 200 {
					assert.Equal(t, tt.expectedItems, actualItems,
						fmt.Sprintf("%s: expected %d items, got %d",
							tt.description, tt.expectedItems, actualItems))
				}
			}
		})
	}
}

// TestGetTop200WithFiltering tests filtering with large datasets
func TestGetTop200WithFiltering(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create RSS feed with 200 items, some containing "technology"
	mockClient := testutil.CreateMockHTTPClient(t, generateRSSFeedWithKeywords(200, "technology", 50))
	handler := NewRSSHandlerWithClient(mockClient)

	router := gin.New()
	router.GET("/api/rss/spiegel/top5", handler.GetTop5)

	// Request 200 items with filter
	req, err := http.NewRequest("GET", "/api/rss/spiegel/top5?limit=200&filter=technology", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response HeadlinesResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Should return filtered items
	for _, item := range response.Headlines {
		assert.Contains(t, strings.ToLower(item.Title), "technology",
			"All returned items should contain the filter keyword")
	}

	// Should have some results but not all 200
	assert.Greater(t, len(response.Headlines), 0, "Should have some filtered results")
	assert.Less(t, len(response.Headlines), 200, "Should have fewer items after filtering")
}

// TestPerformanceWith200Items tests that the endpoint performs well with 200 items
func TestPerformanceWith200Items(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock with 250 items (to ensure we can return 200)
	mockClient := testutil.CreateMockHTTPClient(t, generateLargeRSSFeed(250))
	handler := NewRSSHandlerWithClient(mockClient)

	router := gin.New()
	router.GET("/api/rss/spiegel/top5", handler.GetTop5)

	// Measure time to fetch 200 items
	start := time.Now()

	req, err := http.NewRequest("GET", "/api/rss/spiegel/top5?limit=200", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	elapsed := time.Since(start)

	assert.Equal(t, http.StatusOK, w.Code)

	// Performance assertion - should complete within reasonable time
	assert.Less(t, elapsed, 2*time.Second,
		"Fetching 200 items should complete within 2 seconds")

	var response HeadlinesResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 200, len(response.Headlines),
		"Should return exactly 200 items")
}

// TestCachingWith200Items tests that caching works correctly with large datasets
func TestCachingWith200Items(t *testing.T) {
	gin.SetMode(gin.TestMode)

	callCount := 0
	mockClient := &http.Client{
		Transport: &testutil.MockTransport{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				callCount++
				return &http.Response{
					StatusCode: 200,
					Body:       testutil.CreateReadCloser(generateLargeRSSFeed(200)),
					Header:     make(http.Header),
				}, nil
			},
		},
	}

	handler := NewRSSHandlerWithClient(mockClient)
	router := gin.New()
	router.GET("/api/rss/spiegel/top5", handler.GetTop5)

	// First request - should hit the feed
	req1, _ := http.NewRequest("GET", "/api/rss/spiegel/top5?limit=200", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, 1, callCount, "First request should fetch from feed")

	// Second request immediately - should use cache
	req2, _ := http.NewRequest("GET", "/api/rss/spiegel/top5?limit=200", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, 1, callCount, "Second request should use cache")

	// Verify both responses are identical
	assert.Equal(t, w1.Body.String(), w2.Body.String(),
		"Cached response should be identical")
}

// TestExportLimitValidation tests that export limit is properly validated
func TestExportLimitValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		limit          string
		expectedStatus int
		description    string
	}{
		{
			name:           "Valid limit within range",
			limit:          "500",
			expectedStatus: http.StatusOK,
			description:    "Should accept limit within allowed range",
		},
		{
			name:           "Maximum allowed limit",
			limit:          "1000",
			expectedStatus: http.StatusOK,
			description:    "Should accept maximum limit of 1000",
		},
		{
			name:           "Exceeds maximum limit",
			limit:          "1001",
			expectedStatus: http.StatusBadRequest,
			description:    "Should reject limit exceeding 1000",
		},
		{
			name:           "Negative limit",
			limit:          "-1",
			expectedStatus: http.StatusOK,
			description:    "Should use default for invalid negative limit",
		},
		{
			name:           "Zero limit",
			limit:          "0",
			expectedStatus: http.StatusOK,
			description:    "Should use default for zero limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutil.CreateMockHTTPClient(t, generateLargeRSSFeed(250))
			handler := NewRSSHandlerWithClient(mockClient)
			router := gin.New()
			router.GET("/api/rss/spiegel/export", handler.ExportHeadlines)

			url := fmt.Sprintf("/api/rss/spiegel/export?format=json&limit=%s", tt.limit)
			req, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, tt.description)

			if w.Code == http.StatusBadRequest {
				var response ErrorResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response.Error, "limit exceeds maximum")
			}
		})
	}
}

// TestConcurrentRequests tests that the handler properly handles concurrent requests
func TestConcurrentRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClient := testutil.CreateMockHTTPClient(t, generateLargeRSSFeed(200))
	handler := NewRSSHandlerWithClient(mockClient)
	router := gin.New()
	router.GET("/api/rss/spiegel/top5", handler.GetTop5)

	// Number of concurrent requests
	numRequests := 20
	results := make(chan int, numRequests)
	errors := make(chan error, numRequests)

	// Launch concurrent requests
	for i := 0; i < numRequests; i++ {
		go func(id int) {
			req, err := http.NewRequest("GET", "/api/rss/spiegel/top5?limit=50", nil)
			if err != nil {
				errors <- err
				return
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			results <- w.Code
		}(i)
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		select {
		case code := <-results:
			assert.Equal(t, http.StatusOK, code,
				"All concurrent requests should succeed")
		case err := <-errors:
			t.Fatalf("Error in concurrent request: %v", err)
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent requests")
		}
	}
}

// Helper function to generate a large RSS feed
func generateLargeRSSFeed(itemCount int) string {
	var items strings.Builder
	for i := 1; i <= itemCount; i++ {
		items.WriteString(fmt.Sprintf(`
		<item>
			<title>News Item %d: Important Headlines Today</title>
			<link>https://example.com/news/%d</link>
			<description>This is the description for news item number %d</description>
			<pubDate>Mon, 02 Jan 2006 15:04:05 -0700</pubDate>
		</item>`, i, i, i))
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<rss version="2.0">
		<channel>
			<title>Test RSS Feed</title>
			<link>https://example.com</link>
			<description>Test feed with %d items</description>
			%s
		</channel>
	</rss>`, itemCount, items.String())
}

// Helper function to generate RSS feed with specific keywords
func generateRSSFeedWithKeywords(totalItems int, keyword string, keywordCount int) string {
	var items strings.Builder
	for i := 1; i <= totalItems; i++ {
		title := fmt.Sprintf("News Item %d: Important Headlines Today", i)
		if i <= keywordCount {
			title = fmt.Sprintf("News Item %d: %s News Update", i, keyword)
		}
		items.WriteString(fmt.Sprintf(`
		<item>
			<title>%s</title>
			<link>https://example.com/news/%d</link>
			<description>Description for item %d</description>
			<pubDate>Mon, 02 Jan 2006 15:04:05 -0700</pubDate>
		</item>`, title, i, i))
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<rss version="2.0">
		<channel>
			<title>Test RSS Feed</title>
			<link>https://example.com</link>
			<description>Test feed with keywords</description>
			%s
		</channel>
	</rss>`, items.String())
}