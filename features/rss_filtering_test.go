package features

import (
	"bytes"
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
	"github.com/f00b455/golang-template/internal/config"
	"github.com/f00b455/golang-template/internal/handlers"
	"github.com/f00b455/golang-template/internal/middleware"
	"github.com/f00b455/golang-template/pkg/shared"
	"github.com/gin-gonic/gin"
)

// Mock RSS transport with filtering test data
type mockFilteringRSSTransport struct {
	mockArticles []mockArticle
}

type mockArticle struct {
	title string
	link  string
}

func (m *mockFilteringRSSTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Build RSS XML from mock articles
	var items strings.Builder
	for _, article := range m.mockArticles {
		items.WriteString(fmt.Sprintf(`
		<item>
			<title>%s</title>
			<link>%s</link>
			<description>Article about %s</description>
			<pubDate>Mon, 25 Sep 2025 12:00:00 +0200</pubDate>
		</item>`, article.title, article.link, article.title))
	}

	mockRSS := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>SPIEGEL ONLINE</title>
		<link>https://www.spiegel.de</link>
		<description>Deutschlands führende Nachrichtenseite</description>
		%s
	</channel>
</rss>`, items.String())

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockRSS)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

type rssFilteringContext struct {
	router           *gin.Engine
	response         *httptest.ResponseRecorder
	responseBody     string
	mockClient       *http.Client
	mockTransport    *mockFilteringRSSTransport
	expectedCount    int
	expectedKeyword  string
}

func (ctx *rssFilteringContext) theRSSFeedHasArticles(count int) error {
	ctx.mockTransport.mockArticles = make([]mockArticle, 0, count)
	// Create diverse mock articles for testing
	templates := []string{
		"Breaking News: %s Update",
		"Analysis: %s Today",
		"Expert Opinion on %s",
		"Latest %s Report",
		"Investigation: %s Revealed",
	}

	topics := []string{"Politik", "Wirtschaft", "Sport", "Kultur", "Wissenschaft", "Technologie"}

	for i := 0; i < count; i++ {
		template := templates[i%len(templates)]
		topic := topics[i%len(topics)]
		ctx.mockTransport.mockArticles = append(ctx.mockTransport.mockArticles, mockArticle{
			title: fmt.Sprintf(template, topic),
			link:  fmt.Sprintf("https://www.spiegel.de/article%d", i+1),
		})
	}
	return nil
}

func (ctx *rssFilteringContext) articlesContainTheWord(count int, keyword string) error {
	// Replace first 'count' articles with ones containing the keyword
	for i := 0; i < count && i < len(ctx.mockTransport.mockArticles); i++ {
		ctx.mockTransport.mockArticles[i] = mockArticle{
			title: fmt.Sprintf("Breaking: %s News Today - Article %d", keyword, i+1),
			link:  fmt.Sprintf("https://www.spiegel.de/%s-%d", strings.ToLower(keyword), i+1),
		}
	}
	ctx.expectedCount = count
	ctx.expectedKeyword = keyword
	return nil
}

func (ctx *rssFilteringContext) noArticlesContainTheWord(keyword string) error {
	// Ensure no articles contain the keyword
	for i := range ctx.mockTransport.mockArticles {
		// Replace any occurrence of the keyword
		ctx.mockTransport.mockArticles[i].title = strings.ReplaceAll(
			ctx.mockTransport.mockArticles[i].title,
			keyword,
			"OtherTopic",
		)
	}
	ctx.expectedCount = 0
	ctx.expectedKeyword = keyword
	return nil
}

func (ctx *rssFilteringContext) theRSSFeedHasArticlesWithVariations(variations ...string) error {
	ctx.mockTransport.mockArticles = make([]mockArticle, 0)
	for i, variation := range variations {
		ctx.mockTransport.mockArticles = append(ctx.mockTransport.mockArticles, mockArticle{
			title: fmt.Sprintf("News about %s today", variation),
			link:  fmt.Sprintf("https://www.spiegel.de/article-%d", i+1),
		})
	}
	return nil
}

func (ctx *rssFilteringContext) iRequestTopArticlesWithFilter(limit int, filter string) error {
	endpoint := fmt.Sprintf("/api/rss/spiegel/top5?limit=%d&filter=%s", limit, filter)
	return ctx.makeRequest(endpoint)
}

func (ctx *rssFilteringContext) iRequestArticlesWithFilter(filter string) error {
	endpoint := fmt.Sprintf("/api/rss/spiegel/top5?filter=%s", filter)
	return ctx.makeRequest(endpoint)
}

func (ctx *rssFilteringContext) makeRequest(endpoint string) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}

	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	ctx.responseBody = ctx.response.Body.String()
	return nil
}

func (ctx *rssFilteringContext) iShouldReceiveExactlyArticles(count int) error {
	var response struct {
		Headlines []shared.RssHeadline `json:"headlines"`
	}

	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid response JSON: %w", err)
	}

	if len(response.Headlines) != count {
		return fmt.Errorf("expected exactly %d articles, got %d", count, len(response.Headlines))
	}

	return nil
}

func (ctx *rssFilteringContext) allReturnedArticlesShouldContainInTheirHeadline(keyword string) error {
	var response struct {
		Headlines []shared.RssHeadline `json:"headlines"`
	}

	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid response JSON: %w", err)
	}

	for i, headline := range response.Headlines {
		if !strings.Contains(strings.ToLower(headline.Title), strings.ToLower(keyword)) {
			return fmt.Errorf("headline %d '%s' does not contain keyword '%s'", i, headline.Title, keyword)
		}
	}

	return nil
}

func (ctx *rssFilteringContext) iShouldReceiveAllArticlesRegardlessOfCase() error {
	var response struct {
		Headlines []shared.RssHeadline `json:"headlines"`
	}

	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid response JSON: %w", err)
	}

	// Check that we got articles with COVID in various cases
	foundVariations := make(map[string]bool)
	for _, headline := range response.Headlines {
		titleLower := strings.ToLower(headline.Title)
		if strings.Contains(titleLower, "covid") {
			if strings.Contains(headline.Title, "COVID") {
				foundVariations["COVID"] = true
			} else if strings.Contains(headline.Title, "covid") {
				foundVariations["covid"] = true
			} else if strings.Contains(headline.Title, "Covid") {
				foundVariations["Covid"] = true
			}
		}
	}

	if len(foundVariations) < 2 {
		return fmt.Errorf("case-insensitive filtering not working, found variations: %v", foundVariations)
	}

	return nil
}

func (ctx *rssFilteringContext) iShouldReceiveAnEmptyResultSet() error {
	var response struct {
		Headlines []shared.RssHeadline `json:"headlines"`
	}

	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid response JSON: %w", err)
	}

	if len(response.Headlines) != 0 {
		return fmt.Errorf("expected empty result set, got %d articles", len(response.Headlines))
	}

	return nil
}

func (ctx *rssFilteringContext) iShouldReceiveAllArticlesContaining(substring string) error {
	var response struct {
		Headlines []shared.RssHeadline `json:"headlines"`
	}

	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid response JSON: %w", err)
	}

	for _, headline := range response.Headlines {
		if !strings.Contains(strings.ToLower(headline.Title), strings.ToLower(substring)) {
			return fmt.Errorf("headline '%s' does not contain '%s'", headline.Title, substring)
		}
	}

	return nil
}

func (ctx *rssFilteringContext) allHeadlinesShouldContainCaseInsensitively(keyword string) error {
	var response struct {
		Headlines []shared.RssHeadline `json:"headlines"`
	}

	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		// For single headline response from /latest endpoint
		var headline shared.RssHeadline
		if err := json.Unmarshal([]byte(ctx.responseBody), &headline); err != nil {
			return fmt.Errorf("invalid response JSON: %w", err)
		}

		if !strings.Contains(strings.ToLower(headline.Title), strings.ToLower(keyword)) {
			return fmt.Errorf("headline '%s' does not contain keyword '%s' (case-insensitive)", headline.Title, keyword)
		}
		return nil
	}

	for i, headline := range response.Headlines {
		if !strings.Contains(strings.ToLower(headline.Title), strings.ToLower(keyword)) {
			return fmt.Errorf("headline %d '%s' does not contain keyword '%s' (case-insensitive)", i, headline.Title, keyword)
		}
	}

	return nil
}

func (ctx *rssFilteringContext) theHeadlineShouldContainCaseInsensitively(keyword string) error {
	var headline shared.RssHeadline
	if err := json.Unmarshal([]byte(ctx.responseBody), &headline); err != nil {
		return fmt.Errorf("invalid response JSON: %w", err)
	}

	if !strings.Contains(strings.ToLower(headline.Title), strings.ToLower(keyword)) {
		return fmt.Errorf("headline '%s' does not contain keyword '%s' (case-insensitive)", headline.Title, keyword)
	}

	return nil
}

func InitializeRSSFilteringScenario(ctx *godog.ScenarioContext) {
	featureCtx := &rssFilteringContext{}

	// Background steps
	ctx.Step(`^the API server is running$`, func() error {
		// Set gin to test mode
		gin.SetMode(gin.TestMode)

		// Setup mock transport
		featureCtx.mockTransport = &mockFilteringRSSTransport{
			mockArticles: []mockArticle{},
		}

		// Setup mock HTTP client
		featureCtx.mockClient = &http.Client{
			Transport: featureCtx.mockTransport,
			Timeout:   5 * time.Second,
		}

		// Create router with middleware
		featureCtx.router = gin.New()
		featureCtx.router.Use(gin.Recovery())
		featureCtx.router.Use(middleware.CORS())

		// Set up routes
		api := featureCtx.router.Group("/api")
		{
			// RSS endpoints with mocked HTTP client
			rssHandler := handlers.NewRSSHandlerWithClient(featureCtx.mockClient)
			api.GET("/rss/spiegel/latest", rssHandler.GetLatest)
			api.GET("/rss/spiegel/top5", rssHandler.GetTop5)
		}

		return nil
	})

	// Given steps
	ctx.Step(`^the RSS feed has (\d+) articles$`, featureCtx.theRSSFeedHasArticles)
	ctx.Step(`^(\d+) articles contain the word "([^"]*)"$`, featureCtx.articlesContainTheWord)
	ctx.Step(`^no articles contain the word "([^"]*)"$`, featureCtx.noArticlesContainTheWord)
	ctx.Step(`^the RSS feed has articles with "([^"]*)", "([^"]*)", and "([^"]*)"$`,
		func(v1, v2, v3 string) error {
			return featureCtx.theRSSFeedHasArticlesWithVariations(v1, v2, v3)
		})
	ctx.Step(`^the RSS feed has articles with "([^"]*)" and "([^"]*)"$`,
		func(v1, v2 string) error {
			return featureCtx.theRSSFeedHasArticlesWithVariations(v1, v2)
		})

	// When steps
	ctx.Step(`^I request top (\d+) articles with filter "([^"]*)"$`, featureCtx.iRequestTopArticlesWithFilter)
	ctx.Step(`^I request articles with filter "([^"]*)"$`, featureCtx.iRequestArticlesWithFilter)
	ctx.Step(`^I make a GET request to "([^"]*)"$`, func(endpoint string) error {
		return featureCtx.makeRequest(endpoint)
	})

	// Then steps
	ctx.Step(`^I should receive exactly (\d+) articles$`, featureCtx.iShouldReceiveExactlyArticles)
	ctx.Step(`^all returned articles should contain "([^"]*)" in their headline$`, featureCtx.allReturnedArticlesShouldContainInTheirHeadline)
	ctx.Step(`^I should receive all articles regardless of case$`, featureCtx.iShouldReceiveAllArticlesRegardlessOfCase)
	ctx.Step(`^I should receive an empty result set$`, featureCtx.iShouldReceiveAnEmptyResultSet)
	ctx.Step(`^I should receive all articles containing "([^"]*)"$`, featureCtx.iShouldReceiveAllArticlesContaining)
	ctx.Step(`^the response status should be (\d+)$`, func(status int) error {
		if featureCtx.response.Code != status {
			return fmt.Errorf("expected status %d, got %d", status, featureCtx.response.Code)
		}
		return nil
	})
	ctx.Step(`^the response should contain a headlines array$`, func() error {
		var response struct {
			Headlines []shared.RssHeadline `json:"headlines"`
		}
		if err := json.Unmarshal([]byte(featureCtx.responseBody), &response); err != nil {
			return fmt.Errorf("invalid response JSON: %w", err)
		}
		return nil
	})
	ctx.Step(`^all headlines should contain "([^"]*)" case-insensitively$`, featureCtx.allHeadlinesShouldContainCaseInsensitively)
	ctx.Step(`^the headline should contain "([^"]*)" case-insensitively$`, featureCtx.theHeadlineShouldContainCaseInsensitively)

	// Initialize config before each scenario
	ctx.Before(func(ctx context.Context, scenario *godog.Scenario) (context.Context, error) {
		config.Load()
		return ctx, nil
	})
}

func TestRSSFilteringFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeRSSFilteringScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"rss-filtering.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run RSS filtering feature tests")
	}
}