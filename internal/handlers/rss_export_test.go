package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRSSHandler_ExportHeadlines_JSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := SetupMockServer(MockRSSResponse, http.StatusOK)
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	tests := []struct {
		name           string
		format         string
		limit          string
		expectedStatus int
		expectedType   string
	}{
		{
			name:           "Export JSON format",
			format:         "json",
			limit:          "5",
			expectedStatus: http.StatusOK,
			expectedType:   "application/json",
		},
		{
			name:           "Export JSON with limit 3",
			format:         "json",
			limit:          "3",
			expectedStatus: http.StatusOK,
			expectedType:   "application/json",
		},
		{
			name:           "Export JSON no limit",
			format:         "json",
			limit:          "",
			expectedStatus: http.StatusOK,
			expectedType:   "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.ResetCache()

			url := "/rss/spiegel/export?format=" + tt.format
			if tt.limit != "" {
				url += "&limit=" + tt.limit
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.ExportHeadlines(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Header().Get("Content-Type"), tt.expectedType)
			assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment; filename=")
			assert.Contains(t, w.Header().Get("Content-Disposition"), ".json")

			// Verify JSON response structure
			var response HeadlinesResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.limit != "" {
				limit, _ := strconv.Atoi(tt.limit)
				assert.LessOrEqual(t, len(response.Headlines), limit)
			}
		})
	}
}

func TestRSSHandler_ExportHeadlines_CSV(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := SetupMockServer(MockRSSResponse, http.StatusOK)
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	tests := []struct {
		name           string
		format         string
		limit          string
		expectedStatus int
		expectedType   string
	}{
		{
			name:           "Export CSV format",
			format:         "csv",
			limit:          "5",
			expectedStatus: http.StatusOK,
			expectedType:   "text/csv",
		},
		{
			name:           "Export CSV with limit 2",
			format:         "csv",
			limit:          "2",
			expectedStatus: http.StatusOK,
			expectedType:   "text/csv",
		},
		{
			name:           "Export CSV no limit",
			format:         "csv",
			limit:          "",
			expectedStatus: http.StatusOK,
			expectedType:   "text/csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.ResetCache()

			url := "/rss/spiegel/export?format=" + tt.format
			if tt.limit != "" {
				url += "&limit=" + tt.limit
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.ExportHeadlines(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Header().Get("Content-Type"), tt.expectedType)
			assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment; filename=")
			assert.Contains(t, w.Header().Get("Content-Disposition"), ".csv")

			// Verify CSV structure
			csvContent := w.Body.String()
			assert.Contains(t, csvContent, "Title,Link,Published_At,Source")
			lines := strings.Split(csvContent, "\n")

			if tt.limit != "" {
				limit, _ := strconv.Atoi(tt.limit)
				// +1 for header, +1 for empty line at end
				assert.LessOrEqual(t, len(lines)-2, limit)
			}
		})
	}
}

func TestRSSHandler_ExportHeadlines_Errors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := SetupMockServer(MockRSSResponse, http.StatusOK)
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	tests := []struct {
		name           string
		format         string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Invalid format",
			format:         "xml",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid format parameter: must be 'json' or 'csv'",
		},
		{
			name:           "Missing format",
			format:         "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "missing format parameter",
		},
		{
			name:           "Invalid format with special chars",
			format:         "invalid_format",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid format parameter: must be 'json' or 'csv'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/rss/spiegel/export"
			if tt.format != "" {
				url += "?format=" + tt.format
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.ExportHeadlines(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedError, response.Error)
		})
	}
}

func TestRSSHandler_ExportHeadlines_WithFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock RSS response with varied titles for filtering
	mockRSSWithVariedTitles := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>SPIEGEL ONLINE</title>
    <item>
      <title><![CDATA[Politik: Neue Gesetzgebung]]></title>
      <link><![CDATA[https://www.spiegel.de/1]]></link>
      <pubDate>Mon, 24 Sep 2023 10:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Wirtschaft: DAX steigt]]></title>
      <link><![CDATA[https://www.spiegel.de/2]]></link>
      <pubDate>Mon, 24 Sep 2023 09:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Politik: EU-Gipfel]]></title>
      <link><![CDATA[https://www.spiegel.de/3]]></link>
      <pubDate>Mon, 24 Sep 2023 08:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Sport: Bundesliga]]></title>
      <link><![CDATA[https://www.spiegel.de/4]]></link>
      <pubDate>Mon, 24 Sep 2023 07:00:00 +0000</pubDate>
    </item>
  </channel>
</rss>`

	server := SetupMockServer(mockRSSWithVariedTitles, http.StatusOK)
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	tests := []struct {
		name          string
		format        string
		filter        string
		expectedCount int
	}{
		{
			name:          "JSON export with Politik filter",
			format:        "json",
			filter:        "Politik",
			expectedCount: 2,
		},
		{
			name:          "CSV export with Sport filter",
			format:        "csv",
			filter:        "Sport",
			expectedCount: 1,
		},
		{
			name:          "JSON export no matches",
			format:        "json",
			filter:        "Technology",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.ResetCache()

			url := "/rss/spiegel/export?format=" + tt.format + "&filter=" + tt.filter

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.ExportHeadlines(c)

			assert.Equal(t, http.StatusOK, w.Code)

			if tt.format == "json" {
				var response HeadlinesResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(response.Headlines))
			} else {
				csvContent := w.Body.String()
				lines := strings.Split(csvContent, "\n")
				// -2 for header and empty last line
				actualCount := len(lines) - 2
				if actualCount < 0 {
					actualCount = 0
				}
				assert.Equal(t, tt.expectedCount, actualCount)
			}
		})
	}
}

func TestRSSHandler_ExportHeadlines_NetworkError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = "http://invalid-url-that-does-not-exist.invalid"
	handler.ResetCache()

	req := httptest.NewRequest("GET", "/rss/spiegel/export?format=json", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ExportHeadlines(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Unable to fetch RSS feed", response.Error)
}

func TestRSSHandler_ExportHeadlines_LimitValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := SetupMockServer(MockRSSResponse, http.StatusOK)
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	tests := []struct {
		name          string
		limit         string
		expectedCount int
	}{
		{
			name:          "Negative limit defaults to all",
			limit:         "-1",
			expectedCount: 6,
		},
		{
			name:          "Zero limit returns all",
			limit:         "0",
			expectedCount: 6,
		},
		{
			name:          "Very large limit caps at available",
			limit:         "1000",
			expectedCount: 6,
		},
		{
			name:          "Invalid limit string defaults to all",
			limit:         "abc",
			expectedCount: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.ResetCache()

			url := "/rss/spiegel/export?format=json&limit=" + tt.limit

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.ExportHeadlines(c)

			assert.Equal(t, http.StatusOK, w.Code)

			var response HeadlinesResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(response.Headlines))
		})
	}
}