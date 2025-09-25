package testutil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// MockRSSTransport provides a mock HTTP transport for RSS feed testing
type MockRSSTransport struct {
	ItemCount        int
	SpecialKeyword   string
	KeywordStartItem int
	KeywordEndItem   int
}

// RoundTrip implements the http.RoundTripper interface
func (m *MockRSSTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rssContent := m.GenerateMockRSS()
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(rssContent)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

// GenerateMockRSS generates a mock RSS feed with configurable content
func (m *MockRSSTransport) GenerateMockRSS() string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>SPIEGEL ONLINE</title>
    <link>https://www.spiegel.de</link>
    <description>Deutschlands führende Nachrichtenseite</description>`)

	for i := 1; i <= m.ItemCount; i++ {
		title := m.generateTitle(i)
		builder.WriteString(fmt.Sprintf(`
    <item>
      <title><![CDATA[%s]]></title>
      <link><![CDATA[https://www.spiegel.de/%d]]></link>
      <pubDate>Mon, 24 Sep 2023 %02d:00:00 +0000</pubDate>
    </item>`, title, i, 23-(i%24)))
	}

	builder.WriteString(`
  </channel>
</rss>`)

	return builder.String()
}

func (m *MockRSSTransport) generateTitle(itemNum int) string {
	if m.SpecialKeyword != "" && itemNum >= m.KeywordStartItem && itemNum <= m.KeywordEndItem {
		return fmt.Sprintf("Article with %s %d", m.SpecialKeyword, itemNum)
	}
	return fmt.Sprintf("Regular Article %d", itemNum)
}

// NewLargeMockRSSTransport creates a mock transport with 60 items where a keyword appears in specific range
func NewLargeMockRSSTransport(keyword string, startItem, endItem int) *MockRSSTransport {
	return &MockRSSTransport{
		ItemCount:        60,
		SpecialKeyword:   keyword,
		KeywordStartItem: startItem,
		KeywordEndItem:   endItem,
	}
}