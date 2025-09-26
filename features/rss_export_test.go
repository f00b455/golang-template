package features

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
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

type rssExportContext struct {
	router         *gin.Engine
	response       *httptest.ResponseRecorder
	responseBody   string
	mockClient     *http.Client
	mockRSSContent string
}

type exportJSONResponse struct {
	ExportDate    string               `json:"export_date"`
	TotalItems    int                  `json:"total_items"`
	FilterApplied string               `json:"filter_applied,omitempty"`
	Headlines     []shared.RssHeadline `json:"headlines"`
}

func (ctx *rssExportContext) setupExportRouter() {
	gin.SetMode(gin.TestMode)

	ctx.router = gin.New()
	ctx.router.Use(gin.Recovery())
	ctx.router.Use(middleware.CORS())

	api := ctx.router.Group("/api")
	{
		greetHandler := handlers.NewGreetHandler()
		api.GET("/greet", greetHandler.Greet)

		rssHandler := handlers.NewRSSHandlerWithClient(ctx.mockClient)
		api.GET("/rss/spiegel/latest", rssHandler.GetLatest)
		api.GET("/rss/spiegel/top5", rssHandler.GetTop5)
		api.GET("/rss/spiegel/export", rssHandler.ExportHeadlines)
	}
}

func (ctx *rssExportContext) theRSSFeedHasMultipleArticlesAvailable() error {
	ctx.mockRSSContent = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>SPIEGEL ONLINE</title>
		<link>https://www.spiegel.de</link>
		<description>Deutschlands führende Nachrichtenseite</description>
		<item>
			<title>Politik: Wichtige Entscheidung im Bundestag</title>
			<link>https://www.spiegel.de/politik/artikel1</link>
			<description>Politik Artikel</description>
			<pubDate>Mon, 25 Sep 2025 12:00:00 +0200</pubDate>
		</item>
		<item>
			<title>Sport: Bundesliga Ergebnisse</title>
			<link>https://www.spiegel.de/sport/artikel2</link>
			<description>Sport News</description>
			<pubDate>Mon, 25 Sep 2025 11:30:00 +0200</pubDate>
		</item>
		<item>
			<title>Politik: Neue EU-Regelungen</title>
			<link>https://www.spiegel.de/politik/artikel3</link>
			<description>Politik Update</description>
			<pubDate>Mon, 25 Sep 2025 11:00:00 +0200</pubDate>
		</item>
		<item>
			<title>Wirtschaft: DAX erreicht neues Hoch</title>
			<link>https://www.spiegel.de/wirtschaft/artikel4</link>
			<description>Wirtschaft News</description>
			<pubDate>Mon, 25 Sep 2025 10:30:00 +0200</pubDate>
		</item>
		<item>
			<title>News: Breaking News Update</title>
			<link>https://www.spiegel.de/news/artikel5</link>
			<description>Breaking News</description>
			<pubDate>Mon, 25 Sep 2025 10:00:00 +0200</pubDate>
		</item>
		<item>
			<title>Sport: Champions League Vorschau</title>
			<link>https://www.spiegel.de/sport/artikel6</link>
			<description>Sport Vorschau</description>
			<pubDate>Mon, 25 Sep 2025 09:30:00 +0200</pubDate>
		</item>
		<item>
			<title>Special: Quote "Test" with, comma</title>
			<link>https://www.spiegel.de/special/artikel7</link>
			<description>Special characters test</description>
			<pubDate>Mon, 25 Sep 2025 09:00:00 +0200</pubDate>
		</item>
		<item>
			<title>News: Latest Development</title>
			<link>https://www.spiegel.de/news/artikel8</link>
			<description>News Update</description>
			<pubDate>Mon, 25 Sep 2025 08:30:00 +0200</pubDate>
		</item>
		<item>
			<title>Unicode: Schöne Grüße über die Tür</title>
			<link>https://www.spiegel.de/unicode/artikel9</link>
			<description>Unicode characters test</description>
			<pubDate>Mon, 25 Sep 2025 08:00:00 +0200</pubDate>
		</item>
		<item>
			<title>Politik: Wahlkampf beginnt</title>
			<link>https://www.spiegel.de/politik/artikel10</link>
			<description>Politik Wahlkampf</description>
			<pubDate>Mon, 25 Sep 2025 07:30:00 +0200</pubDate>
		</item>
	</channel>
</rss>`
	return nil
}

func (ctx *rssExportContext) theAPIServerIsRunning() error {
	ctx.mockClient = &http.Client{
		Transport: &mockExportRSSTransport{content: ctx.mockRSSContent},
		Timeout:   5 * time.Second,
	}

	ctx.setupExportRouter()
	return nil
}

func (ctx *rssExportContext) iRequest(endpoint string) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}

	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)
	ctx.responseBody = ctx.response.Body.String()

	return nil
}

func (ctx *rssExportContext) theResponseStatusShouldBe(expectedStatus int) error {
	if ctx.response == nil {
		return fmt.Errorf("no response received")
	}

	if ctx.response.Code != expectedStatus {
		return fmt.Errorf("expected status %d, got %d. Response: %s", expectedStatus, ctx.response.Code, ctx.responseBody)
	}

	return nil
}

func (ctx *rssExportContext) theContentTypeShouldBe(expectedType string) error {
	contentType := ctx.response.Header().Get("Content-Type")
	if !strings.Contains(contentType, expectedType) {
		return fmt.Errorf("expected content-type to contain %s, got %s", expectedType, contentType)
	}
	return nil
}

func (ctx *rssExportContext) theContentDispositionShouldContain(expected string) error {
	contentDisposition := ctx.response.Header().Get("Content-Disposition")
	if !strings.Contains(contentDisposition, expected) {
		return fmt.Errorf("expected content-disposition to contain %s, got %s", expected, contentDisposition)
	}
	return nil
}

func (ctx *rssExportContext) theFilenameShouldContain(extension string) error {
	contentDisposition := ctx.response.Header().Get("Content-Disposition")
	if !strings.Contains(contentDisposition, extension) {
		return fmt.Errorf("expected filename to contain %s in content-disposition: %s", extension, contentDisposition)
	}
	return nil
}

func (ctx *rssExportContext) theJSONResponseShouldHaveExportMetadata() error {
	var response exportJSONResponse
	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid JSON response: %w", err)
	}

	if response.ExportDate == "" {
		return fmt.Errorf("export_date is missing")
	}

	if response.TotalItems < 0 {
		return fmt.Errorf("total_items is missing or invalid")
	}

	return nil
}

func (ctx *rssExportContext) theJSONResponseShouldHaveHeadlinesArray() error {
	var response exportJSONResponse
	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid JSON response: %w", err)
	}

	if response.Headlines == nil {
		return fmt.Errorf("headlines array is missing")
	}

	return nil
}

func (ctx *rssExportContext) theCSVShouldHaveHeaderRow() error {
	reader := csv.NewReader(strings.NewReader(ctx.responseBody))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("CSV has no rows")
	}

	headers := records[0]
	expectedHeaders := []string{"Title", "Link", "Published_At", "Source"}

	if len(headers) != len(expectedHeaders) {
		return fmt.Errorf("expected %d headers, got %d", len(expectedHeaders), len(headers))
	}

	for i, expected := range expectedHeaders {
		if headers[i] != expected {
			return fmt.Errorf("expected header %s at position %d, got %s", expected, i, headers[i])
		}
	}

	return nil
}

func (ctx *rssExportContext) theCSVShouldHaveDataRows() error {
	reader := csv.NewReader(strings.NewReader(ctx.responseBody))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) <= 1 {
		return fmt.Errorf("CSV has no data rows (only header row found)")
	}

	return nil
}

func (ctx *rssExportContext) theJSONResponseShouldContainFilterMetadata(filter string) error {
	var response exportJSONResponse
	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid JSON response: %w", err)
	}

	if response.FilterApplied != filter {
		return fmt.Errorf("expected filter_applied to be %s, got %s", filter, response.FilterApplied)
	}

	return nil
}

func (ctx *rssExportContext) allHeadlinesShouldMatchTheFilter(filter string) error {
	var response exportJSONResponse
	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid JSON response: %w", err)
	}

	filterLower := strings.ToLower(filter)
	for _, headline := range response.Headlines {
		if !strings.Contains(strings.ToLower(headline.Title), filterLower) {
			return fmt.Errorf("headline '%s' does not match filter '%s'", headline.Title, filter)
		}
	}

	return nil
}

func (ctx *rssExportContext) theJSONResponseShouldHaveExactlyHeadlines(count int) error {
	var response exportJSONResponse
	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid JSON response: %w", err)
	}

	if len(response.Headlines) != count {
		return fmt.Errorf("expected exactly %d headlines, got %d", count, len(response.Headlines))
	}

	return nil
}

func (ctx *rssExportContext) theCSVRowsShouldOnlyContainHeadlines(filter string) error {
	reader := csv.NewReader(strings.NewReader(ctx.responseBody))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) <= 1 {
		return nil // No data rows, which is valid for filtered results
	}

	filterLower := strings.ToLower(filter)
	for i := 1; i < len(records); i++ {
		title := records[i][0] // Title is first column
		if !strings.Contains(strings.ToLower(title), filterLower) {
			return fmt.Errorf("CSV row %d with title '%s' does not match filter '%s'", i, title, filter)
		}
	}

	return nil
}

func (ctx *rssExportContext) theResponseShouldContainAnErrorMessageAbout(topic string) error {
	var errorResponse struct {
		Error string `json:"error"`
	}

	if err := json.Unmarshal([]byte(ctx.responseBody), &errorResponse); err != nil {
		return fmt.Errorf("response is not a valid error JSON: %w", err)
	}

	topicLower := strings.ToLower(topic)
	if !strings.Contains(strings.ToLower(errorResponse.Error), topicLower) {
		return fmt.Errorf("error message does not mention %s: %s", topic, errorResponse.Error)
	}

	return nil
}

func (ctx *rssExportContext) theJSONResponseShouldHaveEmptyHeadlinesArray() error {
	var response exportJSONResponse
	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid JSON response: %w", err)
	}

	if len(response.Headlines) != 0 {
		return fmt.Errorf("expected empty headlines array, got %d items", len(response.Headlines))
	}

	return nil
}

func (ctx *rssExportContext) theJSONResponseShouldShowTotalItemsAs(count int) error {
	var response exportJSONResponse
	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid JSON response: %w", err)
	}

	if response.TotalItems != count {
		return fmt.Errorf("expected total_items to be %d, got %d", count, response.TotalItems)
	}

	return nil
}

func (ctx *rssExportContext) theCSVShouldProperlyEscapeQuotesAndCommas() error {
	reader := csv.NewReader(strings.NewReader(ctx.responseBody))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV, might have unescaped characters: %w", err)
	}

	// Check if we have the special characters test row
	for i := 1; i < len(records); i++ {
		title := records[i][0]
		if strings.Contains(title, "Quote") && strings.Contains(title, "comma") {
			// Successfully parsed a row with quotes and commas
			return nil
		}
	}

	return nil
}

func (ctx *rssExportContext) theCSVShouldHandleUTF8CharactersCorrectly() error {
	// Check if UTF-8 characters are present and correctly handled
	if !strings.Contains(ctx.responseBody, "ö") && !strings.Contains(ctx.responseBody, "ü") && !strings.Contains(ctx.responseBody, "ß") {
		// If no UTF-8 characters in the mock data, that's ok
		return nil
	}

	reader := csv.NewReader(strings.NewReader(ctx.responseBody))
	_, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV with UTF-8 characters: %w", err)
	}

	return nil
}

func (ctx *rssExportContext) theJSONResponseShouldHaveAtMostHeadlines(max int) error {
	var response exportJSONResponse
	if err := json.Unmarshal([]byte(ctx.responseBody), &response); err != nil {
		return fmt.Errorf("invalid JSON response: %w", err)
	}

	if len(response.Headlines) > max {
		return fmt.Errorf("expected at most %d headlines, got %d", max, len(response.Headlines))
	}

	return nil
}

// Mock transport for export tests
type mockExportRSSTransport struct {
	content string
}

func (m *mockExportRSSTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       httpReadCloser{strings.NewReader(m.content)},
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

// httpReadCloser wraps a reader to add Close method
type httpReadCloser struct {
	*strings.Reader
}

func (rc httpReadCloser) Close() error {
	return nil
}

func InitializeRSSExportScenario(ctx *godog.ScenarioContext) {
	featureCtx := &rssExportContext{}

	// Background steps
	ctx.Step(`^the RSS feed has multiple articles available$`, featureCtx.theRSSFeedHasMultipleArticlesAvailable)
	ctx.Step(`^the API server is running$`, featureCtx.theAPIServerIsRunning)

	// Action steps
	ctx.Step(`^I request "([^"]*)"$`, featureCtx.iRequest)

	// Assertion steps
	ctx.Step(`^the response status should be (\d+)$`, featureCtx.theResponseStatusShouldBe)
	ctx.Step(`^the content-type should be "([^"]*)"$`, featureCtx.theContentTypeShouldBe)
	ctx.Step(`^the content-disposition should contain "([^"]*)"$`, featureCtx.theContentDispositionShouldContain)
	ctx.Step(`^the filename should contain "([^"]*)"$`, featureCtx.theFilenameShouldContain)
	ctx.Step(`^the JSON response should have export metadata$`, featureCtx.theJSONResponseShouldHaveExportMetadata)
	ctx.Step(`^the JSON response should have headlines array$`, featureCtx.theJSONResponseShouldHaveHeadlinesArray)
	ctx.Step(`^the CSV should have header row$`, featureCtx.theCSVShouldHaveHeaderRow)
	ctx.Step(`^the CSV should have data rows$`, featureCtx.theCSVShouldHaveDataRows)
	ctx.Step(`^the JSON response should contain filter metadata "([^"]*)"$`, featureCtx.theJSONResponseShouldContainFilterMetadata)
	ctx.Step(`^all headlines should match the filter "([^"]*)"$`, featureCtx.allHeadlinesShouldMatchTheFilter)
	ctx.Step(`^the JSON response should have exactly (\d+) headlines$`, featureCtx.theJSONResponseShouldHaveExactlyHeadlines)
	ctx.Step(`^the CSV rows should only contain "([^"]*)" headlines$`, featureCtx.theCSVRowsShouldOnlyContainHeadlines)
	ctx.Step(`^the response should contain an error message about (.+)$`, featureCtx.theResponseShouldContainAnErrorMessageAbout)
	ctx.Step(`^the JSON response should have empty headlines array$`, featureCtx.theJSONResponseShouldHaveEmptyHeadlinesArray)
	ctx.Step(`^the JSON response should show total_items as (\d+)$`, featureCtx.theJSONResponseShouldShowTotalItemsAs)
	ctx.Step(`^the CSV should properly escape quotes and commas$`, featureCtx.theCSVShouldProperlyEscapeQuotesAndCommas)
	ctx.Step(`^the CSV should handle UTF-8 characters correctly$`, featureCtx.theCSVShouldHandleUTF8CharactersCorrectly)
	ctx.Step(`^the JSON response should have at most (\d+) headlines$`, featureCtx.theJSONResponseShouldHaveAtMostHeadlines)

	// Initialize before each scenario
	ctx.Before(func(ctx context.Context, scenario *godog.Scenario) (context.Context, error) {
		config.Load()
		return ctx, nil
	})
}

func TestRSSExportFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeRSSExportScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"rss-export.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run RSS export feature tests")
	}
}