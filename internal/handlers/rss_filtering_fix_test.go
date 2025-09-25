package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestRSSHandler_FilteringSearchesEntireFeed verifies the fix for issue #15
// It ensures that filtering searches through more than just the first 5 items
func TestRSSHandler_FilteringSearchesEntireFeed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock RSS feed with 60 items where the keyword only appears after item 5
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>SPIEGEL ONLINE</title>`)

	// First 10 items without the special keyword
	for i := 1; i <= 10; i++ {
		builder.WriteString(fmt.Sprintf(`
    <item>
      <title><![CDATA[Regular Article %d]]></title>
      <link><![CDATA[https://www.spiegel.de/%d]]></link>
      <pubDate>Mon, 24 Sep 2023 %02d:00:00 +0000</pubDate>
    </item>`, i, i, 23-i))
	}

	// Items 11-20 contain "special-keyword"
	for i := 11; i <= 20; i++ {
		builder.WriteString(fmt.Sprintf(`
    <item>
      <title><![CDATA[Article with special-keyword %d]]></title>
      <link><![CDATA[https://www.spiegel.de/%d]]></link>
      <pubDate>Mon, 24 Sep 2023 %02d:00:00 +0000</pubDate>
    </item>`, i, i, 23-i%24))
	}

	// Items 21-60 are regular articles
	for i := 21; i <= 60; i++ {
		builder.WriteString(fmt.Sprintf(`
    <item>
      <title><![CDATA[Other News %d]]></title>
      <link><![CDATA[https://www.spiegel.de/%d]]></link>
      <pubDate>Mon, 23 Sep 2023 %02d:00:00 +0000</pubDate>
    </item>`, i, i, 23-i%24))
	}

	builder.WriteString(`
  </channel>
</rss>`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(builder.String()))
	}))
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	// Test filtering for keyword that only appears after item 10
	req := httptest.NewRequest("GET", "/rss/spiegel/top5?filter=special-keyword", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetTop5(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response HeadlinesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should find results even though they're not in the first 5 items
	assert.Greater(t, len(response.Headlines), 0, "Should find headlines with special-keyword even though they're after item 10")

	// Verify all returned headlines contain the keyword
	for _, headline := range response.Headlines {
		assert.Contains(t, strings.ToLower(headline.Title), "special-keyword",
			"All filtered headlines should contain the search keyword")
	}

	// Should return at most 5 results
	assert.LessOrEqual(t, len(response.Headlines), 5, "Should return at most 5 filtered results")

	// TotalCount should reflect that we fetched many items (at least 50)
	assert.GreaterOrEqual(t, response.TotalCount, 50, "TotalCount should show we fetched at least 50 items")
}

// TestRSSHandler_CacheStoresMoreItems verifies that the cache stores more items for better filtering
func TestRSSHandler_CacheStoresMoreItems(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock RSS feed with 60 items
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>SPIEGEL ONLINE</title>`)

	for i := 1; i <= 60; i++ {
		keyword := ""
		if i > 30 && i <= 35 {
			keyword = " with rare-term"
		}
		builder.WriteString(fmt.Sprintf(`
    <item>
      <title><![CDATA[Article %d%s]]></title>
      <link><![CDATA[https://www.spiegel.de/%d]]></link>
      <pubDate>Mon, 24 Sep 2023 %02d:00:00 +0000</pubDate>
    </item>`, i, keyword, i, 23-i%24))
	}

	builder.WriteString(`
  </channel>
</rss>`)

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(builder.String()))
	}))
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	// First request without filter - should cache items
	req1 := httptest.NewRequest("GET", "/rss/spiegel/top5", nil)
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request = req1
	handler.GetTop5(c1)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, 1, requestCount, "First request should fetch from server")

	// Second request with filter - should use cache
	req2 := httptest.NewRequest("GET", "/rss/spiegel/top5?filter=rare-term", nil)
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = req2
	handler.GetTop5(c2)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, 1, requestCount, "Second request should use cache, not fetch again")

	var response HeadlinesResponse
	err := json.Unmarshal(w2.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should find the rare-term items even though they're after item 30
	assert.Greater(t, len(response.Headlines), 0, "Should find items with rare-term from cache")
	for _, headline := range response.Headlines {
		assert.Contains(t, strings.ToLower(headline.Title), "rare-term")
	}
}