package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/f00b455/golang-template/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestRSSHandler_FilteringSearchesEntireFeed verifies the fix for issue #15
// It ensures that filtering searches through more than just the first 5 items
func TestRSSHandler_FilteringSearchesEntireFeed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Use the shared mock transport that creates a feed with keywords in items 11-20
	mockTransport := testutil.NewLargeMockRSSTransport("special-keyword", 11, 20)
	mockClient := &http.Client{
		Transport: mockTransport,
	}

	handler := NewRSSHandlerWithClient(mockClient)
	handler.cfg.SpiegelRSSURL = "http://test.example.com/rss"
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

	// Use shared mock transport with rare-term in items 31-35
	mockTransport := testutil.NewLargeMockRSSTransport("rare-term", 31, 35)
	mockClient := &http.Client{
		Transport: mockTransport,
	}

	handler := NewRSSHandlerWithClient(mockClient)
	handler.cfg.SpiegelRSSURL = "http://test.example.com/rss"
	handler.ResetCache()

	// First request without filter - should cache items
	req1 := httptest.NewRequest("GET", "/rss/spiegel/top5", nil)
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request = req1
	handler.GetTop5(c1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second request with filter - should use cache
	req2 := httptest.NewRequest("GET", "/rss/spiegel/top5?filter=rare-term", nil)
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = req2
	handler.GetTop5(c2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var response HeadlinesResponse
	err := json.Unmarshal(w2.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should find the rare-term items even though they're after item 30
	assert.Greater(t, len(response.Headlines), 0, "Should find items with rare-term from cache")
	for _, headline := range response.Headlines {
		assert.Contains(t, strings.ToLower(headline.Title), "rare-term")
	}
}