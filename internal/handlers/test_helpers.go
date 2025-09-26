package handlers

import (
	"net/http"
	"net/http/httptest"
)

// Mock RSS response for testing
const MockRSSResponse = `<?xml version="1.0" encoding="UTF-8"?>
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
    <item>
      <title><![CDATA[Headline 6]]></title>
      <link><![CDATA[https://www.spiegel.de/6]]></link>
      <pubDate>Mon, 24 Sep 2023 05:00:00 +0000</pubDate>
    </item>
  </channel>
</rss>`

// Mock RSS response with fewer items for testing
const MockRSSResponseFewItems = `<?xml version="1.0" encoding="UTF-8"?>
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

// SetupMockServer creates a test HTTP server that returns mock RSS data
func SetupMockServer(response string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(response))
	}))
}