package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	handler.cache = &cacheEntry{
		data:      testHeadline,
		timestamp: time.Now(),
	}
	handler.multiCache = &multiCacheEntry{
		data:      []shared.RssHeadline{*testHeadline},
		timestamp: time.Now(),
	}

	// Verify cache has data
	assert.NotNil(t, handler.cache.data)
	assert.NotEmpty(t, handler.multiCache.data)

	// Reset cache
	handler.ResetCache()

	// Verify cache is empty
	assert.Nil(t, handler.cache.data)
	assert.Empty(t, handler.multiCache.data)
}