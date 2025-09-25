package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/f00b455/golang-template/pkg/shared"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock RSS response for testing
const mockRSSResponse = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>SPIEGEL ONLINE</title>
    <item>
      <title><![CDATA[Headline 1]]></title>
      <link><![CDATA[https://www.spiegel.de/1]]></link>
      <pubDate>Mon, 24 Sep 2023 10:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Headline 2]]></title>
      <link><![CDATA[https://www.spiegel.de/2]]></link>
      <pubDate>Mon, 24 Sep 2023 09:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Headline 3]]></title>
      <link><![CDATA[https://www.spiegel.de/3]]></link>
      <pubDate>Mon, 24 Sep 2023 08:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Headline 4]]></title>
      <link><![CDATA[https://www.spiegel.de/4]]></link>
      <pubDate>Mon, 24 Sep 2023 07:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Headline 5]]></title>
      <link><![CDATA[https://www.spiegel.de/5]]></link>
      <pubDate>Mon, 24 Sep 2023 06:00:00 +0000</pubDate>
    </item>
  </channel>
</rss>`

const mockRSSResponseSmall = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>SPIEGEL ONLINE</title>
    <item>
      <title><![CDATA[Headline 1]]></title>
      <link><![CDATA[https://www.spiegel.de/1]]></link>
      <pubDate>Mon, 24 Sep 2023 10:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Headline 2]]></title>
      <link><![CDATA[https://www.spiegel.de/2]]></link>
      <pubDate>Mon, 24 Sep 2023 09:00:00 +0000</pubDate>
    </item>
  </channel>
</rss>`

// setupMockServer creates a test HTTP server that returns mock RSS data
func setupMockServer(response string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(response))
	}))
}

func TestRSSHandler_GetLatest_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup mock server
	server := setupMockServer(mockRSSResponse, http.StatusOK)
	defer server.Close()

	// Create handler with mock URL
	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL

	// Reset cache to ensure fresh fetch
	handler.ResetCache()

	req := httptest.NewRequest("GET", "/rss/spiegel/latest", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetLatest(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Headline 1", response["title"])
	assert.Equal(t, "https://www.spiegel.de/1", response["link"])
	assert.Equal(t, "SPIEGEL", response["source"])
	assert.NotEmpty(t, response["publishedAt"])
}

func TestRSSHandler_GetLatest_NetworkError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = "http://invalid-url-that-does-not-exist.invalid"
	handler.ResetCache()

	req := httptest.NewRequest("GET", "/rss/spiegel/latest", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetLatest(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Unable to fetch RSS feed", response.Error)
}

func TestRSSHandler_GetLatest_ServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := setupMockServer("", http.StatusInternalServerError)
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	req := httptest.NewRequest("GET", "/rss/spiegel/latest", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetLatest(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Unable to fetch RSS feed", response.Error)
}

func TestRSSHandler_GetTop5_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := setupMockServer(mockRSSResponse, http.StatusOK)
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	req := httptest.NewRequest("GET", "/rss/spiegel/top5", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetTop5(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response HeadlinesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response.Headlines, 5)
	assert.Equal(t, "Headline 1", response.Headlines[0].Title)
	assert.Equal(t, "https://www.spiegel.de/1", response.Headlines[0].Link)
	assert.Equal(t, "SPIEGEL", response.Headlines[0].Source)
}

func TestRSSHandler_GetTop5_WithLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := setupMockServer(mockRSSResponse, http.StatusOK)
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	tests := []struct {
		name     string
		limit    string
		expected int
	}{
		{"limit 2", "2", 2},
		{"limit 3", "3", 3},
		{"limit 5", "5", 5},
		{"invalid limit defaults to 5", "10", 5},
		{"invalid string defaults to 5", "abc", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.ResetCache() // Reset cache for each test

			req := httptest.NewRequest("GET", "/rss/spiegel/top5?limit="+tt.limit, nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.GetTop5(c)

			assert.Equal(t, http.StatusOK, w.Code)

			var response HeadlinesResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response.Headlines, tt.expected)
		})
	}
}

func TestRSSHandler_GetTop5_FewerThan5Items(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := setupMockServer(mockRSSResponseSmall, http.StatusOK)
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	req := httptest.NewRequest("GET", "/rss/spiegel/top5", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetTop5(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response HeadlinesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Headlines, 2)
}

func TestRSSHandler_GetTop5_NetworkError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = "http://invalid-url-that-does-not-exist.invalid"
	handler.ResetCache()

	req := httptest.NewRequest("GET", "/rss/spiegel/top5", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetTop5(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Unable to fetch RSS feed", response.Error)
}

func TestRSSHandler_Cache_Functionality(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Track number of requests to mock server
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockRSSResponse))
	}))
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	// First request - should hit the server
	req1 := httptest.NewRequest("GET", "/rss/spiegel/latest", nil)
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request = req1

	handler.GetLatest(c1)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, 1, requestCount)

	// Second request immediately after - should use cache
	req2 := httptest.NewRequest("GET", "/rss/spiegel/latest", nil)
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = req2

	handler.GetLatest(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, 1, requestCount) // Still only 1 request

	// Responses should be identical
	assert.Equal(t, w1.Body.String(), w2.Body.String())
}

func TestRSSHandler_ResetCache(t *testing.T) {
	handler := NewRSSHandler()

	// Create test headline
	testHeadline := &shared.RssHeadline{
		Title:       "Test Headline",
		Link:        "https://test.com",
		PublishedAt: time.Now().Format(time.RFC3339),
		Source:      "TEST",
	}

	// Manually set cache data
	handler.cache["latest:test"] = &cacheEntry{
		data:      testHeadline,
		timestamp: time.Now(),
	}
	handler.multiCache["top5:test"] = &multiCacheEntry{
		data:      []shared.RssHeadline{*testHeadline},
		timestamp: time.Now(),
	}

	// Verify cache has data
	assert.Len(t, handler.cache, 1)
	assert.Len(t, handler.multiCache, 1)

	// Reset cache
	handler.ResetCache()

	// Verify cache is empty
	assert.Empty(t, handler.cache)
	assert.Empty(t, handler.multiCache)
}

// Test cases for filtering functionality

const mockRSSResponseWithFilter = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>SPIEGEL ONLINE</title>
    <item>
      <title><![CDATA[Breaking News: Major Tech Announcement]]></title>
      <link><![CDATA[https://www.spiegel.de/tech1]]></link>
      <pubDate>Mon, 24 Sep 2023 10:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Politics Today: Election Updates]]></title>
      <link><![CDATA[https://www.spiegel.de/politics1]]></link>
      <pubDate>Mon, 24 Sep 2023 09:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Tech Giants Report Earnings]]></title>
      <link><![CDATA[https://www.spiegel.de/tech2]]></link>
      <pubDate>Mon, 24 Sep 2023 08:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Sports: Championship Results]]></title>
      <link><![CDATA[https://www.spiegel.de/sports1]]></link>
      <pubDate>Mon, 24 Sep 2023 07:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Breaking: Economic Data Released]]></title>
      <link><![CDATA[https://www.spiegel.de/economy1]]></link>
      <pubDate>Mon, 24 Sep 2023 06:00:00 +0000</pubDate>
    </item>
  </channel>
</rss>`

func TestRSSHandler_GetLatest_WithFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := setupMockServer(mockRSSResponseWithFilter, http.StatusOK)
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	tests := []struct {
		name        string
		filter      string
		expectTitle string
		expectEmpty bool
	}{
		{"filter tech", "tech", "Breaking News: Major Tech Announcement", false},
		{"filter breaking", "breaking", "Breaking News: Major Tech Announcement", false},
		{"filter politics", "politics", "Politics Today: Election Updates", false},
		{"filter nonexistent", "xyz123nonexistent", "", true},
		{"case insensitive TECH", "TECH", "Breaking News: Major Tech Announcement", false},
		{"empty filter returns first", "", "Breaking News: Major Tech Announcement", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.ResetCache()

			req := httptest.NewRequest("GET", "/rss/spiegel/latest?filter="+tt.filter, nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.GetLatest(c)

			assert.Equal(t, http.StatusOK, w.Code)

			if tt.expectEmpty {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Empty(t, response)
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectTitle, response["title"])
			}
		})
	}
}

func TestRSSHandler_GetTop5_WithFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := setupMockServer(mockRSSResponseWithFilter, http.StatusOK)
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	tests := []struct {
		name          string
		filter        string
		limit         string
		expectedCount int
		checkTitles   bool
	}{
		{"filter tech returns 2", "tech", "5", 2, true},
		{"filter breaking returns 2", "breaking", "5", 2, true},
		{"filter with limit 1", "breaking", "1", 1, true},
		{"case insensitive", "TECH", "5", 2, true},
		{"no matches", "xyz123nonexistent", "5", 0, false},
		{"empty filter returns all", "", "5", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.ResetCache()

			url := "/rss/spiegel/top5?filter=" + tt.filter
			if tt.limit != "" {
				url += "&limit=" + tt.limit
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.GetTop5(c)

			assert.Equal(t, http.StatusOK, w.Code)

			var response HeadlinesResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response.Headlines, tt.expectedCount)

			if tt.checkTitles && tt.filter != "" && tt.expectedCount > 0 {
				filterLower := strings.ToLower(tt.filter)
				for _, headline := range response.Headlines {
					assert.Contains(t, strings.ToLower(headline.Title), filterLower)
				}
			}
		})
	}
}

func TestRSSHandler_FilterAppliedBeforeLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create RSS with many tech items
	mockRSS := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>SPIEGEL ONLINE</title>`

	// Add 10 tech headlines
	for i := 1; i <= 10; i++ {
		mockRSS += fmt.Sprintf(`
    <item>
      <title><![CDATA[Tech News %d: Innovation Update]]></title>
      <link><![CDATA[https://www.spiegel.de/tech%d]]></link>
      <pubDate>Mon, 24 Sep 2023 %02d:00:00 +0000</pubDate>
    </item>`, i, i, 10-i)
	}

	// Add 10 non-tech headlines
	for i := 1; i <= 10; i++ {
		mockRSS += fmt.Sprintf(`
    <item>
      <title><![CDATA[Sports News %d: Game Results]]></title>
      <link><![CDATA[https://www.spiegel.de/sports%d]]></link>
      <pubDate>Mon, 24 Sep 2023 %02d:00:00 +0000</pubDate>
    </item>`, i, i, 10-i)
	}

	mockRSS += `
  </channel>
</rss>`

	server := setupMockServer(mockRSS, http.StatusOK)
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	// Request filtered tech headlines with limit 5
	req := httptest.NewRequest("GET", "/rss/spiegel/top5?filter=tech&limit=5", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetTop5(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response HeadlinesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should get exactly 5 tech headlines
	assert.Len(t, response.Headlines, 5)

	// All should contain "tech"
	for _, headline := range response.Headlines {
		assert.Contains(t, strings.ToLower(headline.Title), "tech")
	}
}