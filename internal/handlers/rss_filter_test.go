package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/f00b455/golang-template/pkg/shared"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRSSHandler_GetTop5_WithFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock RSS response with different headline titles for filtering
	mockRSSWithVariedTitles := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>SPIEGEL ONLINE</title>
    <item>
      <title><![CDATA[Politik: Neue Gesetzgebung verabschiedet]]></title>
      <link><![CDATA[https://www.spiegel.de/1]]></link>
      <pubDate>Mon, 24 Sep 2023 10:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Wirtschaft: DAX erreicht neues Hoch]]></title>
      <link><![CDATA[https://www.spiegel.de/2]]></link>
      <pubDate>Mon, 24 Sep 2023 09:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Politik: EU-Gipfel in Brüssel]]></title>
      <link><![CDATA[https://www.spiegel.de/3]]></link>
      <pubDate>Mon, 24 Sep 2023 08:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Sport: Bayern München gewinnt]]></title>
      <link><![CDATA[https://www.spiegel.de/4]]></link>
      <pubDate>Mon, 24 Sep 2023 07:00:00 +0000</pubDate>
    </item>
    <item>
      <title><![CDATA[Wirtschaft: Inflation sinkt weiter]]></title>
      <link><![CDATA[https://www.spiegel.de/5]]></link>
      <pubDate>Mon, 24 Sep 2023 06:00:00 +0000</pubDate>
    </item>
  </channel>
</rss>`

	server := SetupMockServer(mockRSSWithVariedTitles, http.StatusOK)
	defer server.Close()

	handler := NewRSSHandler()
	handler.cfg.SpiegelRSSURL = server.URL
	handler.ResetCache()

	tests := []struct {
		name           string
		filter         string
		expectedCount  int
		expectedTitles []string
	}{
		{
			name:           "Filter Politik",
			filter:         "Politik",
			expectedCount:  2,
			expectedTitles: []string{"Politik: Neue Gesetzgebung verabschiedet", "Politik: EU-Gipfel in Brüssel"},
		},
		{
			name:           "Filter Wirtschaft",
			filter:         "Wirtschaft",
			expectedCount:  2,
			expectedTitles: []string{"Wirtschaft: DAX erreicht neues Hoch", "Wirtschaft: Inflation sinkt weiter"},
		},
		{
			name:           "Filter Sport",
			filter:         "Sport",
			expectedCount:  1,
			expectedTitles: []string{"Sport: Bayern München gewinnt"},
		},
		{
			name:          "Filter NonExistent",
			filter:        "Technology",
			expectedCount: 0,
		},
		{
			name:          "Empty filter returns all",
			filter:        "",
			expectedCount: 5,
		},
		{
			name:           "Case insensitive filter",
			filter:         "politik",
			expectedCount:  2,
			expectedTitles: []string{"Politik: Neue Gesetzgebung verabschiedet", "Politik: EU-Gipfel in Brüssel"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset cache before each test
			handler.ResetCache()

			url := "/rss/spiegel/top5"
			if tt.filter != "" {
				url += "?filter=" + tt.filter
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
			assert.Equal(t, tt.expectedCount, len(response.Headlines))

			// Check that the expected titles are present
			if tt.expectedTitles != nil {
				for i, expectedTitle := range tt.expectedTitles {
					if i < len(response.Headlines) {
						assert.Equal(t, expectedTitle, response.Headlines[i].Title)
					}
				}
			}
		})
	}
}

func TestFilterHeadlines(t *testing.T) {
	handler := NewRSSHandler()

	headlines := []shared.RssHeadline{
		{Title: "Politik: Neue Gesetze", Link: "http://example.com/1"},
		{Title: "Wirtschaft: DAX steigt", Link: "http://example.com/2"},
		{Title: "Sport: Bundesliga Highlights", Link: "http://example.com/3"},
		{Title: "Kultur: Neue Ausstellung", Link: "http://example.com/4"},
		{Title: "Politik: EU-Gipfel", Link: "http://example.com/5"},
	}

	tests := []struct {
		name          string
		filter        string
		expectedCount int
		expectedFirst string
	}{
		{
			name:          "Filter Politik",
			filter:        "Politik",
			expectedCount: 2,
			expectedFirst: "Politik: Neue Gesetze",
		},
		{
			name:          "Filter Wirtschaft",
			filter:        "Wirtschaft",
			expectedCount: 1,
			expectedFirst: "Wirtschaft: DAX steigt",
		},
		{
			name:          "Case insensitive filter",
			filter:        "politik",
			expectedCount: 2,
			expectedFirst: "Politik: Neue Gesetze",
		},
		{
			name:          "Partial match",
			filter:        "Liga",
			expectedCount: 1,
			expectedFirst: "Sport: Bundesliga Highlights",
		},
		{
			name:          "No matches",
			filter:        "Technology",
			expectedCount: 0,
		},
		{
			name:          "Empty filter returns all",
			filter:        "",
			expectedCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.filterHeadlines(headlines, tt.filter)

			assert.Equal(t, tt.expectedCount, len(result))
			if tt.expectedCount > 0 && tt.expectedFirst != "" {
				assert.Equal(t, tt.expectedFirst, result[0].Title)
			}
		})
	}
}


func TestRSSHandler_validateFilter(t *testing.T) {
	handler := NewRSSHandler()

	tests := []struct {
		name        string
		filter      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid filter",
			filter:      "Politik",
			expectError: false,
		},
		{
			name:        "Empty filter is valid",
			filter:      "",
			expectError: false,
		},
		{
			name:        "Filter at max length",
			filter:      strings.Repeat("a", 100),
			expectError: false,
		},
		{
			name:        "Filter exceeds max length",
			filter:      strings.Repeat("a", 101),
			expectError: true,
			errorMsg:    "filter parameter too long (max 100 characters)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateFilter(tt.filter)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}