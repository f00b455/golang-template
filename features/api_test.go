package features

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/f00b455/golang-template/internal/config"
	"github.com/f00b455/golang-template/internal/handlers"
	"github.com/f00b455/golang-template/internal/middleware"
	"github.com/f00b455/golang-template/pkg/shared"
	"github.com/gin-gonic/gin"
)

// Mock RSS transport for intercepting HTTP requests
type mockRSSTransport struct{}

func (m *mockRSSTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Return mock RSS XML for SPIEGEL feed
	mockRSS := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>SPIEGEL ONLINE</title>
		<link>https://www.spiegel.de</link>
		<description>Deutschlands führende Nachrichtenseite</description>
		<item>
			<title>Mock Headline 1</title>
			<link>https://www.spiegel.de/article1</link>
			<description>First mock article</description>
			<pubDate>Mon, 25 Sep 2025 12:00:00 +0200</pubDate>
		</item>
		<item>
			<title>Mock Headline 2</title>
			<link>https://www.spiegel.de/article2</link>
			<description>Second mock article</description>
			<pubDate>Mon, 25 Sep 2025 11:00:00 +0200</pubDate>
		</item>
		<item>
			<title>Mock Headline 3</title>
			<link>https://www.spiegel.de/article3</link>
			<description>Third mock article</description>
			<pubDate>Mon, 25 Sep 2025 10:00:00 +0200</pubDate>
		</item>
		<item>
			<title>Mock Headline 4</title>
			<link>https://www.spiegel.de/article4</link>
			<description>Fourth mock article</description>
			<pubDate>Mon, 25 Sep 2025 09:00:00 +0200</pubDate>
		</item>
		<item>
			<title>Mock Headline 5</title>
			<link>https://www.spiegel.de/article5</link>
			<description>Fifth mock article</description>
			<pubDate>Mon, 25 Sep 2025 08:00:00 +0200</pubDate>
		</item>
	</channel>
</rss>`

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockRSS)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

type apiMockContext struct {
	router       *gin.Engine
	response     *httptest.ResponseRecorder
	responseBody string
	lastError    error
	mockClient   *http.Client
}

func (ctx *apiMockContext) setupRouter() {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create router with middleware
	ctx.router = gin.New()
	ctx.router.Use(gin.Recovery())
	ctx.router.Use(middleware.CORS())

	// Set up routes
	api := ctx.router.Group("/api")
	{
		// Greet endpoints
		greetHandler := handlers.NewGreetHandler()
		api.GET("/greet", greetHandler.Greet)

		// RSS endpoints with mocked HTTP client
		rssHandler := handlers.NewRSSHandlerWithClient(ctx.mockClient)
		api.GET("/rss/spiegel/latest", rssHandler.GetLatest)
		api.GET("/rss/spiegel/top5", rssHandler.GetTop5)
	}
}

func (ctx *apiMockContext) theAPIServerIsRunning() error {
	// Setup mock HTTP client for RSS feeds
	ctx.mockClient = &http.Client{
		Transport: &mockRSSTransport{},
		Timeout:   5 * time.Second,
	}

	// Setup router with mocked dependencies
	ctx.setupRouter()
	return nil
}

func (ctx *apiMockContext) iMakeAGETRequestTo(endpoint string) error {
	// Create test request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		ctx.lastError = err
		return nil
	}

	// Record response
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	ctx.responseBody = ctx.response.Body.String()

	return nil
}

func (ctx *apiMockContext) theResponseStatusShouldBe(expectedStatus int) error {
	if ctx.response == nil {
		return fmt.Errorf("no response received")
	}

	if ctx.response.Code != expectedStatus {
		return fmt.Errorf("expected status %d, got %d. Response: %s", expectedStatus, ctx.response.Code, ctx.responseBody)
	}

	return nil
}

func (ctx *apiMockContext) theResponseShouldContainJSON(expectedJSON string) error {
	if ctx.responseBody == "" {
		return fmt.Errorf("no response body")
	}

	var expected, actual any

	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		return fmt.Errorf("invalid expected JSON: %w", err)
	}

	if err := json.Unmarshal([]byte(ctx.responseBody), &actual); err != nil {
		return fmt.Errorf("invalid response JSON: %w", err)
	}

	expectedStr, _ := json.Marshal(expected)
	actualStr, _ := json.Marshal(actual)

	if string(expectedStr) != string(actualStr) {
		return fmt.Errorf("expected JSON %s, got %s", expectedStr, actualStr)
	}

	return nil
}

func (ctx *apiMockContext) theResponseShouldContainAValidRSSHeadline() error {
	if ctx.responseBody == "" {
		return fmt.Errorf("no response body")
	}

	var headline shared.RssHeadline
	if err := json.Unmarshal([]byte(ctx.responseBody), &headline); err != nil {
		return fmt.Errorf("invalid RSS headline JSON: %w", err)
	}

	return nil
}

func (ctx *apiMockContext) theHeadlineShouldHaveATitle() error {
	var headline shared.RssHeadline
	if err := json.Unmarshal([]byte(ctx.responseBody), &headline); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if headline.Title == "" {
		return fmt.Errorf("headline title is empty")
	}

	return nil
}

func (ctx *apiMockContext) theHeadlineShouldHaveALink() error {
	var headline shared.RssHeadline
	if err := json.Unmarshal([]byte(ctx.responseBody), &headline); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if headline.Link == "" {
		return fmt.Errorf("headline link is empty")
	}

	return nil
}

func (ctx *apiMockContext) theHeadlineShouldHaveAPublishedAtTimestamp() error {
	var headline shared.RssHeadline
	if err := json.Unmarshal([]byte(ctx.responseBody), &headline); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if headline.PublishedAt == "" {
		return fmt.Errorf("headline publishedAt is empty")
	}

	// Try to parse the timestamp
	if _, err := time.Parse(time.RFC3339, headline.PublishedAt); err != nil {
		return fmt.Errorf("invalid publishedAt timestamp format: %s", headline.PublishedAt)
	}

	return nil
}

func (ctx *apiMockContext) theHeadlineShouldHaveSource(expectedSource string) error {
	var headline shared.RssHeadline
	if err := json.Unmarshal([]byte(ctx.responseBody), &headline); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if headline.Source != expectedSource {
		return fmt.Errorf("expected source %s, got %s", expectedSource, headline.Source)
	}

	return nil
}

func (ctx *apiMockContext) theResponseShouldContainAHeadlinesArray() error {
	var response struct {
		Headlines []shared.RssHeadline `json:"headlines"`
	}

	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid headlines response JSON: %w", err)
	}

	if response.Headlines == nil {
		return fmt.Errorf("headlines array is null")
	}

	return nil
}

func (ctx *apiMockContext) theHeadlinesArrayShouldHaveOrFewerItems(maxItems int) error {
	var response struct {
		Headlines []shared.RssHeadline `json:"headlines"`
	}

	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid headlines response JSON: %w", err)
	}

	if len(response.Headlines) > maxItems {
		return fmt.Errorf("expected %d or fewer headlines, got %d", maxItems, len(response.Headlines))
	}

	return nil
}

func (ctx *apiMockContext) theHeadlinesArrayShouldHaveExactlyItems(expectedItems int) error {
	var response struct {
		Headlines []shared.RssHeadline `json:"headlines"`
	}

	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid headlines response JSON: %w", err)
	}

	if len(response.Headlines) != expectedItems {
		return fmt.Errorf("expected exactly %d headlines, got %d", expectedItems, len(response.Headlines))
	}

	return nil
}

func (ctx *apiMockContext) theHeadlinesArrayShouldHaveOrMoreItems(minItems int) error {
	var response struct {
		Headlines []shared.RssHeadline `json:"headlines"`
	}

	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid headlines response JSON: %w", err)
	}

	if len(response.Headlines) < minItems {
		return fmt.Errorf("expected %d or more headlines, got %d", minItems, len(response.Headlines))
	}

	return nil
}

func (ctx *apiMockContext) eachHeadlineShouldHaveTitleLinkPublishedAtAndSourceFields() error {
	var response struct {
		Headlines []shared.RssHeadline `json:"headlines"`
	}

	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid headlines response JSON: %w", err)
	}

	for i, headline := range response.Headlines {
		if headline.Title == "" {
			return fmt.Errorf("headline %d missing title", i)
		}
		if headline.Link == "" {
			return fmt.Errorf("headline %d missing link", i)
		}
		if headline.PublishedAt == "" {
			return fmt.Errorf("headline %d missing publishedAt", i)
		}
		if headline.Source == "" {
			return fmt.Errorf("headline %d missing source", i)
		}
	}

	return nil
}

func InitializeMockAPIScenario(ctx *godog.ScenarioContext) {
	featureCtx := &apiMockContext{}

	// Background steps
	ctx.Step(`^the API server is running$`, featureCtx.theAPIServerIsRunning)

	// Action steps
	ctx.Step(`^I make a GET request to "([^"]*)"$`, featureCtx.iMakeAGETRequestTo)

	// Assertion steps
	ctx.Step(`^the response status should be (\d+)$`, featureCtx.theResponseStatusShouldBe)
	ctx.Step(`^the response should contain JSON (.+)$`, featureCtx.theResponseShouldContainJSON)
	ctx.Step(`^the response should contain a valid RSS headline$`, featureCtx.theResponseShouldContainAValidRSSHeadline)
	ctx.Step(`^the headline should have a title$`, featureCtx.theHeadlineShouldHaveATitle)
	ctx.Step(`^the headline should have a link$`, featureCtx.theHeadlineShouldHaveALink)
	ctx.Step(`^the headline should have a publishedAt timestamp$`, featureCtx.theHeadlineShouldHaveAPublishedAtTimestamp)
	ctx.Step(`^the headline should have source "([^"]*)"$`, featureCtx.theHeadlineShouldHaveSource)
	ctx.Step(`^the response should contain a headlines array$`, featureCtx.theResponseShouldContainAHeadlinesArray)
	ctx.Step(`^the headlines array should have (\d+) or fewer items$`, featureCtx.theHeadlinesArrayShouldHaveOrFewerItems)
	ctx.Step(`^the headlines array should have exactly (\d+) items?$`, featureCtx.theHeadlinesArrayShouldHaveExactlyItems)
	ctx.Step(`^the headlines array should have (\d+) or more items?$`, featureCtx.theHeadlinesArrayShouldHaveOrMoreItems)
	ctx.Step(`^each headline should have title, link, publishedAt, and source fields$`, featureCtx.eachHeadlineShouldHaveTitleLinkPublishedAtAndSourceFields)

	// Initialize router before each scenario
	ctx.Before(func(ctx context.Context, scenario *godog.Scenario) (context.Context, error) {
		// Initialize config to avoid nil pointer issues
		config.Load()
		return ctx, nil
	})
}

func TestAPIFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeMockAPIScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"api-greet.feature", "api-rss.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run API feature tests with mocks")
	}
}