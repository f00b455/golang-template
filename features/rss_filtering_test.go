package features

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/f00b455/golang-template/internal/handlers"
	"github.com/f00b455/golang-template/pkg/shared"
	"github.com/gin-gonic/gin"
)

type contextKey string

const rssFilteringContextKey contextKey = "rssFilteringContext"

type rssFilteringContext struct {
	server       *httptest.Server
	responseBody []byte
	statusCode   int
	mockTransport *mockFilteringRSSTransport
	articles      []mockArticle
}

type mockArticle struct {
	title string
}

type mockFilteringRSSTransport struct {
	articles []mockArticle
}

func (m *mockFilteringRSSTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.Contains(req.URL.Path, "rss") || strings.Contains(req.URL.Host, "spiegel") {
		rssContent := m.generateRSS()
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(rssContent)),
			Header:     http.Header{"Content-Type": []string{"application/rss+xml"}},
		}, nil
	}
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader("Not found")),
	}, nil
}

func (m *mockFilteringRSSTransport) generateRSS() string {
	var items strings.Builder
	pubDate := time.Now()

	for _, article := range m.articles {
		items.WriteString(fmt.Sprintf(`
		<item>
			<title><![CDATA[%s]]></title>
			<link>https://www.spiegel.de/article/%d</link>
			<pubDate>%s</pubDate>
			<description>Test description</description>
		</item>`, article.title, len(article.title), pubDate.Format(time.RFC1123Z)))
		pubDate = pubDate.Add(-1 * time.Hour)
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>SPIEGEL ONLINE</title>
		<link>https://www.spiegel.de</link>
		<description>Test RSS Feed</description>
		%s
	</channel>
</rss>`, items.String())
}

func (ctx *rssFilteringContext) theRSSFeedHasTheseArticles(articles *godog.Table) error {
	ctx.articles = []mockArticle{}

	// Skip header row
	for i := 1; i < len(articles.Rows); i++ {
		if len(articles.Rows[i].Cells) > 0 {
			ctx.articles = append(ctx.articles, mockArticle{
				title: articles.Rows[i].Cells[0].Value,
			})
		}
	}

	// Set up mock transport with these articles
	ctx.mockTransport = &mockFilteringRSSTransport{articles: ctx.articles}

	return nil
}

func (ctx *rssFilteringContext) allHeadlinesShouldContainInTheirTitle(keyword string) error {
	var response handlers.HeadlinesResponse
	if err := json.Unmarshal(ctx.responseBody, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	lowercaseKeyword := strings.ToLower(keyword)
	for _, headline := range response.Headlines {
		if !strings.Contains(strings.ToLower(headline.Title), lowercaseKeyword) {
			return fmt.Errorf("headline '%s' does not contain keyword '%s'", headline.Title, keyword)
		}
	}

	return nil
}

func (ctx *rssFilteringContext) theHeadlinesArrayShouldHaveExactlyItems(count int) error {
	var response handlers.HeadlinesResponse
	if err := json.Unmarshal(ctx.responseBody, &response); err != nil {
		// Try single headline response for latest endpoint
		var singleHeadline shared.RssHeadline
		if err := json.Unmarshal(ctx.responseBody, &singleHeadline); err == nil {
			if count != 1 {
				return fmt.Errorf("expected %d items, but got 1 item", count)
			}
			return nil
		}
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Headlines) != count {
		return fmt.Errorf("expected %d headlines, but got %d", count, len(response.Headlines))
	}

	return nil
}

func (ctx *rssFilteringContext) theHeadlinesArrayShouldBeEmpty() error {
	var response handlers.HeadlinesResponse
	if err := json.Unmarshal(ctx.responseBody, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Headlines) != 0 {
		return fmt.Errorf("expected empty array, but got %d headlines", len(response.Headlines))
	}

	return nil
}

func (ctx *rssFilteringContext) theHeadlineTitleShouldContain(keyword string) error {
	var headline shared.RssHeadline
	if err := json.Unmarshal(ctx.responseBody, &headline); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !strings.Contains(strings.ToLower(headline.Title), strings.ToLower(keyword)) {
		return fmt.Errorf("headline title '%s' does not contain keyword '%s'", headline.Title, keyword)
	}

	return nil
}

func (ctx *rssFilteringContext) theResponseShouldContainAnErrorMessage() error {
	var errorResponse handlers.ErrorResponse
	if err := json.Unmarshal(ctx.responseBody, &errorResponse); err != nil {
		return fmt.Errorf("failed to parse error response: %w", err)
	}

	if errorResponse.Error == "" {
		return fmt.Errorf("expected error message in response, but got none")
	}

	return nil
}

func (ctx *rssFilteringContext) setupFilteringServer() {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create RSS handler with mock HTTP client
	httpClient := &http.Client{
		Transport: ctx.mockTransport,
		Timeout:   2 * time.Second,
	}
	rssHandler := handlers.NewRSSHandlerWithClient(httpClient)

	// Register routes
	api := router.Group("/api")
	rss := api.Group("/rss")
	spiegel := rss.Group("/spiegel")
	spiegel.GET("/latest", rssHandler.GetLatest)
	spiegel.GET("/top5", rssHandler.GetTop5)

	ctx.server = httptest.NewServer(router)
}

func TestRSSFilteringFeatures(t *testing.T) {
	suite := godog.TestSuite{
		TestSuiteInitializer: InitializeRSSFilteringTestSuite,
		ScenarioInitializer:  InitializeRSSFilteringScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"rss-filtering.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeRSSFilteringTestSuite(ctx *godog.TestSuiteContext) {
	// Suite-level initialization if needed
}

func InitializeRSSFilteringScenario(sctx *godog.ScenarioContext) {
	ctx := &rssFilteringContext{}

	// Common API context steps
	sctx.Step(`^the API server is running$`, ctx.theAPIServerIsRunning)
	sctx.Step(`^I make a GET request to "([^"]*)"$`, ctx.iMakeAGETRequestTo)
	sctx.Step(`^the response status should be (\d+)$`, ctx.theResponseStatusShouldBe)
	sctx.Step(`^the response should contain a headlines array$`, ctx.theResponseShouldContainAHeadlinesArray)

	// RSS filtering specific steps
	sctx.Step(`^the RSS feed has these articles:$`, ctx.theRSSFeedHasTheseArticles)
	sctx.Step(`^all headlines should contain "([^"]*)" in their title$`, ctx.allHeadlinesShouldContainInTheirTitle)
	sctx.Step(`^the headlines array should have exactly (\d+) items?$`, ctx.theHeadlinesArrayShouldHaveExactlyItems)
	sctx.Step(`^the headlines array should be empty$`, ctx.theHeadlinesArrayShouldBeEmpty)
	sctx.Step(`^the headline title should contain "([^"]*)"$`, ctx.theHeadlineTitleShouldContain)
	sctx.Step(`^the response should contain an error message$`, ctx.theResponseShouldContainAnErrorMessage)

	sctx.Before(func(goCtx context.Context, sc *godog.Scenario) (context.Context, error) {
		// Reset state before each scenario
		goCtx = context.WithValue(goCtx, rssFilteringContextKey, ctx)
		return goCtx, nil
	})

	sctx.After(func(goCtx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		// Cleanup after each scenario
		if ctx.server != nil {
			ctx.server.Close()
		}
		return goCtx, nil
	})
}

func (ctx *rssFilteringContext) theAPIServerIsRunning() error {
	// Setup server with mock RSS transport
	ctx.setupFilteringServer()
	return nil
}

func (ctx *rssFilteringContext) iMakeAGETRequestTo(path string) error {
	resp, err := http.Get(ctx.server.URL + path)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	ctx.responseBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ctx.statusCode = resp.StatusCode
	return nil
}

func (ctx *rssFilteringContext) theResponseStatusShouldBe(expectedStatus int) error {
	if ctx.statusCode != expectedStatus {
		return fmt.Errorf("expected status %d, but got %d", expectedStatus, ctx.statusCode)
	}
	return nil
}

func (ctx *rssFilteringContext) theResponseShouldContainAHeadlinesArray() error {
	var response handlers.HeadlinesResponse
	if err := json.Unmarshal(ctx.responseBody, &response); err != nil {
		return fmt.Errorf("failed to parse response as headlines array: %w", err)
	}
	return nil
}