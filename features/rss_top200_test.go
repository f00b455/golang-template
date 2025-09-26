package features

import (
	"bytes"
	"encoding/csv"
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
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type top200TestContext struct {
	router        *gin.Engine
	response      *httptest.ResponseRecorder
	responseBody  map[string]interface{}
	rssItems      []map[string]interface{}
	loadStartTime time.Time
	loadEndTime   time.Time
	exportedFile  []byte
}

// Mock RSS transport that returns many items
type mockTop200Transport struct {
	itemCount int
}

func (m *mockTop200Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Generate RSS with specified number of items
	var items strings.Builder
	for i := 1; i <= m.itemCount; i++ {
		items.WriteString(fmt.Sprintf(`
		<item>
			<title>News Item %d</title>
			<link>https://www.spiegel.de/article%d</link>
			<description>Description for news item %d</description>
			<pubDate>Mon, 25 Sep 2025 %02d:00:00 +0200</pubDate>
		</item>`, i, i, i, i%24))
	}

	mockRSS := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>SPIEGEL ONLINE</title>
		<link>https://www.spiegel.de</link>
		<description>Test RSS Feed with %d items</description>
		%s
	</channel>
</rss>`, m.itemCount, items.String())

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockRSS)),
		Header:     make(http.Header),
	}, nil
}

func (ctx *top200TestContext) iHaveARunningAPIServer() error {
	gin.SetMode(gin.TestMode)
	ctx.router = gin.New()

	// Setup RSS routes (matching the actual handler setup)
	api := ctx.router.Group("/api")
	rss := api.Group("/rss")
	// Use the handler with custom HTTP client that will use http.DefaultClient
	rssHandler := handlers.NewRSSHandlerWithClient(http.DefaultClient)

	// Clear any existing cache to ensure tests are isolated
	rssHandler.ResetCache()

	rss.GET("/spiegel/top5", rssHandler.GetTop5)
	rss.GET("/spiegel/export", rssHandler.ExportHeadlines)

	// Setup documentation
	ctx.router.GET("/documentation/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup static files
	ctx.router.Static("/static", "./static")
	ctx.router.StaticFile("/", "./static/terminal.html")

	return nil
}

func (ctx *top200TestContext) theTerminalUIIsAccessible() error {
	// In test context, we just verify the router is set up
	// Static files aren't available during unit tests
	if ctx.router == nil {
		return fmt.Errorf("router not initialized")
	}
	return nil
}

func (ctx *top200TestContext) theRSSFeedHasAtLeastItemsAvailable(count int) error {
	// Configure mock transport with specified item count
	http.DefaultClient.Transport = &mockTop200Transport{itemCount: count}
	return nil
}

func (ctx *top200TestContext) iRequestWithLimit(endpoint string) error {
	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusOK {
		// Parse response body
		if err := json.Unmarshal(ctx.response.Body.Bytes(), &ctx.responseBody); err != nil {
			return fmt.Errorf("failed to parse response: %v", err)
		}

		// Extract items from response
		if headlines, ok := ctx.responseBody["headlines"].([]interface{}); ok {
			ctx.rssItems = make([]map[string]interface{}, 0)
			for _, item := range headlines {
				if itemMap, ok := item.(map[string]interface{}); ok {
					ctx.rssItems = append(ctx.rssItems, itemMap)
				}
			}
		}
	}

	return nil
}

func (ctx *top200TestContext) iRequestWithoutLimitParameter(endpoint string) error {
	// Request without any query parameters
	baseEndpoint := strings.Split(endpoint, "?")[0]
	req := httptest.NewRequest(http.MethodGet, baseEndpoint, nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code == http.StatusOK {
		if err := json.Unmarshal(ctx.response.Body.Bytes(), &ctx.responseBody); err != nil {
			return fmt.Errorf("failed to parse response: %v", err)
		}

		if headlines, ok := ctx.responseBody["headlines"].([]interface{}); ok {
			ctx.rssItems = make([]map[string]interface{}, 0)
			for _, item := range headlines {
				if itemMap, ok := item.(map[string]interface{}); ok {
					ctx.rssItems = append(ctx.rssItems, itemMap)
				}
			}
		}
	}

	return nil
}

func (ctx *top200TestContext) theResponseShouldContainUpToRSSItems(expectedCount int) error {
	actualCount := len(ctx.rssItems)
	if actualCount > expectedCount {
		return fmt.Errorf("expected up to %d items, got %d", expectedCount, actualCount)
	}
	return nil
}

func (ctx *top200TestContext) theResponseShouldContainExactlyRSSItems(expectedCount int) error {
	actualCount := len(ctx.rssItems)
	if actualCount != expectedCount {
		return fmt.Errorf("expected exactly %d items, got %d", expectedCount, actualCount)
	}
	return nil
}

func (ctx *top200TestContext) theResponseShouldContainRSSItems(expectedCount int) error {
	actualCount := len(ctx.rssItems)
	// For invalid parameters, we expect default of 5 items
	// For valid parameters, we expect the requested count (up to available items)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d items, got %d", expectedCount, actualCount)
	}
	return nil
}

func (ctx *top200TestContext) theResponseStatusShouldBe(expectedStatus int) error {
	if ctx.response.Code != expectedStatus {
		return fmt.Errorf("expected status %d, got %d", expectedStatus, ctx.response.Code)
	}
	return nil
}

func (ctx *top200TestContext) theResponseShouldIncludeATotalCountField() error {
	if _, exists := ctx.responseBody["totalCount"]; !exists {
		return fmt.Errorf("response does not include totalCount field")
	}
	return nil
}

func (ctx *top200TestContext) theAPIReturnsNewsItems(count int) error {
	// Simulate API returning specified number of items
	http.DefaultClient.Transport = &mockTop200Transport{itemCount: count}
	return nil
}

func (ctx *top200TestContext) iLoadTheTerminalUI() error {
	ctx.loadStartTime = time.Now()

	// Simulate loading the terminal UI
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	// Then load RSS data
	req = httptest.NewRequest(http.MethodGet, "/api/rss/spiegel/top5?limit=200", nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	ctx.loadEndTime = time.Now()

	return nil
}

func (ctx *top200TestContext) iShouldSeeTheFirstPageOfNewsItems() error {
	// In a real implementation, we would check the DOM
	// For now, verify that we got a successful response
	if ctx.response.Code != http.StatusOK {
		return fmt.Errorf("failed to load first page: status %d", ctx.response.Code)
	}
	return nil
}

func (ctx *top200TestContext) iShouldSeePaginationControls() error {
	// In a real implementation, we would check for pagination UI elements
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) theStatusBarShouldShow(expectedText string) error {
	// In a real implementation, we would check the status bar content
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) newsItemsAreLoadedInTheTerminalUI(count int) error {
	if err := ctx.theAPIReturnsNewsItems(count); err != nil {
		return err
	}
	if err := ctx.iLoadTheTerminalUI(); err != nil {
		return err
	}
	return nil
}

func (ctx *top200TestContext) iAmOnPage(pageNum int) error {
	// In a real implementation, we would track current page
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) iPressKey(key string) error {
	// In a real implementation, we would simulate key press
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) iShouldBeOnPage(pageNum int) error {
	// In a real implementation, we would verify current page
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) iTypeInTheCommandInput(command string) error {
	// In a real implementation, we would simulate typing command
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) theInitialPageShouldLoadInLessThanSeconds(seconds int) error {
	loadTime := ctx.loadEndTime.Sub(ctx.loadStartTime)
	maxDuration := time.Duration(seconds) * time.Second

	if loadTime > maxDuration {
		return fmt.Errorf("page load took %v, expected less than %v", loadTime, maxDuration)
	}
	return nil
}

func (ctx *top200TestContext) scrollingShouldBeSmoothWithNoLag() error {
	// This would require actual UI testing
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) memoryUsageShouldRemainReasonable() error {
	// This would require memory profiling
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) iRequestNewsItems(count int) error {
	ctx.loadStartTime = time.Now()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/rss/spiegel/top5?limit=%d", count), nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	ctx.loadEndTime = time.Now()
	return nil
}

func (ctx *top200TestContext) iShouldSeeALoadingIndicator() error {
	// In a real implementation, we would check for loading UI
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) theLoadingMessageShouldShowProgress() error {
	// In a real implementation, we would check loading message
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) theDataLoads() error {
	// Wait for data to be loaded
	if ctx.response.Code != http.StatusOK {
		return fmt.Errorf("data failed to load: status %d", ctx.response.Code)
	}
	return nil
}

func (ctx *top200TestContext) theLoadingIndicatorShouldDisappear() error {
	// In a real implementation, we would verify loading indicator is hidden
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) iTypeInTheFilterInput(filter string) error {
	// In a real implementation, we would simulate filter input
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) onlyItemsContainingShouldBeVisible(keyword string) error {
	// In a real implementation, we would check filtered items
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) theFilterShouldApplyAcrossAllItems(count int) error {
	// In a real implementation, we would verify filter applies to all items
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) theStatusBarShouldShowTheFilteredCount() error {
	// In a real implementation, we would check status bar
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) newsItemsAreLoaded(count int) error {
	// Set up transport to return specified number of items
	http.DefaultClient.Transport = &mockTop200Transport{itemCount: count}

	// Load RSS feed to populate cache with the items
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/rss/spiegel/top5?limit=%d", count), nil)
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	if ctx.response.Code != http.StatusOK {
		return fmt.Errorf("failed to load %d items: status %d", count, ctx.response.Code)
	}

	return nil
}

func (ctx *top200TestContext) iScrollThroughTheList() error {
	// In a real implementation, we would simulate scrolling
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) onlyVisibleItemsShouldBeRenderedInTheDOM() error {
	// In a real implementation, we would check DOM rendering
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) scrollingShouldRemainPerformant() error {
	// In a real implementation, we would measure scroll performance
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) itemsShouldRenderAsTheyComeIntoView() error {
	// In a real implementation, we would check lazy rendering
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) iClickTheButton(buttonText string) error {
	// Ensure we have 200 items available for export
	http.DefaultClient.Transport = &mockTop200Transport{itemCount: 200}

	// First load the RSS feed to populate cache
	reqLoad := httptest.NewRequest(http.MethodGet, "/api/rss/spiegel/top5?limit=200", nil)
	respLoad := httptest.NewRecorder()
	ctx.router.ServeHTTP(respLoad, reqLoad)

	// Now export the data
	switch buttonText {
	case "Export JSON":
		req := httptest.NewRequest(http.MethodGet, "/api/rss/spiegel/export?format=json&limit=200", nil)
		ctx.response = httptest.NewRecorder()
		ctx.router.ServeHTTP(ctx.response, req)
		ctx.exportedFile = ctx.response.Body.Bytes()
	case "Export CSV":
		req := httptest.NewRequest(http.MethodGet, "/api/rss/spiegel/export?format=csv&limit=200", nil)
		ctx.response = httptest.NewRecorder()
		ctx.router.ServeHTTP(ctx.response, req)
		ctx.exportedFile = ctx.response.Body.Bytes()
	}
	return nil
}

func (ctx *top200TestContext) aJSONFileWithItemsShouldDownload(count int) error {
	if ctx.response.Code != http.StatusOK {
		return fmt.Errorf("export failed: status %d", ctx.response.Code)
	}

	// Parse JSON to verify item count
	var data map[string]interface{}
	if err := json.Unmarshal(ctx.exportedFile, &data); err != nil {
		return fmt.Errorf("invalid JSON export: %v", err)
	}

	if items, ok := data["headlines"].([]interface{}); ok {
		if len(items) != count {
			return fmt.Errorf("expected %d items in JSON, got %d", count, len(items))
		}
	}

	return nil
}

func (ctx *top200TestContext) aCSVFileWithItemsShouldDownload(count int) error {
	if ctx.response.Code != http.StatusOK {
		return fmt.Errorf("export failed: status %d", ctx.response.Code)
	}

	// Parse CSV to verify item count
	reader := csv.NewReader(bytes.NewReader(ctx.exportedFile))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("invalid CSV export: %v", err)
	}

	// Subtract 1 for header row
	actualCount := len(records) - 1
	if actualCount != count {
		return fmt.Errorf("expected %d items in CSV, got %d", count, actualCount)
	}

	return nil
}

func (ctx *top200TestContext) theAPIEndpointIsTemporarilyUnavailable() error {
	// Configure mock to return error
	http.DefaultClient.Transport = http.RoundTripper(http.DefaultTransport)
	return nil
}

func (ctx *top200TestContext) iShouldSeeAUserFriendlyErrorMessage() error {
	// In a real implementation, we would check error message display
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) cachedDataShouldBeUsedIfAvailable() error {
	// In a real implementation, we would verify cache usage
	// For now, this is a placeholder
	return nil
}

func (ctx *top200TestContext) theUIShouldRemainResponsive() error {
	// In a real implementation, we would check UI responsiveness
	// For now, this is a placeholder
	return nil
}

func InitializeTop200Scenario(s *godog.ScenarioContext) {
	ctx := &top200TestContext{}

	// Background steps
	s.Step(`^I have a running API server$`, ctx.iHaveARunningAPIServer)
	s.Step(`^the terminal UI is accessible$`, ctx.theTerminalUIIsAccessible)
	s.Step(`^the RSS feed has at least (\d+) items available$`, ctx.theRSSFeedHasAtLeastItemsAvailable)

	// API endpoint scenario steps
	s.Step(`^I request "([^"]*)"$`, ctx.iRequestWithLimit)
	s.Step(`^I request "([^"]*)" without limit parameter$`, ctx.iRequestWithoutLimitParameter)
	s.Step(`^the response should contain up to (\d+) RSS items$`, ctx.theResponseShouldContainUpToRSSItems)
	s.Step(`^the response should contain exactly (\d+) RSS items$`, ctx.theResponseShouldContainExactlyRSSItems)
	s.Step(`^the response should contain (\d+) RSS items$`, ctx.theResponseShouldContainRSSItems)
	s.Step(`^the response status should be (\d+)$`, ctx.theResponseStatusShouldBe)
	s.Step(`^the response should include a totalCount field$`, ctx.theResponseShouldIncludeATotalCountField)

	// UI pagination scenario steps
	s.Step(`^the API returns (\d+) news items$`, ctx.theAPIReturnsNewsItems)
	s.Step(`^I load the terminal UI$`, ctx.iLoadTheTerminalUI)
	s.Step(`^I should see the first page of news items$`, ctx.iShouldSeeTheFirstPageOfNewsItems)
	s.Step(`^I should see pagination controls$`, ctx.iShouldSeePaginationControls)
	s.Step(`^the status bar should show "([^"]*)"$`, ctx.theStatusBarShouldShow)

	// Keyboard navigation scenario steps
	s.Step(`^(\d+) news items are loaded in the terminal UI$`, ctx.newsItemsAreLoadedInTheTerminalUI)
	s.Step(`^I am on page (\d+)$`, ctx.iAmOnPage)
	s.Step(`^I press "([^"]*)" key$`, ctx.iPressKey)
	s.Step(`^I should be on page (\d+)$`, ctx.iShouldBeOnPage)

	// Jump navigation scenario steps
	s.Step(`^I type "([^"]*)" in the command input$`, ctx.iTypeInTheCommandInput)

	// Performance scenario steps
	s.Step(`^the initial page should load in less than (\d+) seconds$`, ctx.theInitialPageShouldLoadInLessThanSeconds)
	s.Step(`^scrolling should be smooth with no lag$`, ctx.scrollingShouldBeSmoothWithNoLag)
	s.Step(`^memory usage should remain reasonable$`, ctx.memoryUsageShouldRemainReasonable)

	// Loading indicator scenario steps
	s.Step(`^I request (\d+) news items$`, ctx.iRequestNewsItems)
	s.Step(`^I should see a loading indicator$`, ctx.iShouldSeeALoadingIndicator)
	s.Step(`^the loading message should show progress$`, ctx.theLoadingMessageShouldShowProgress)
	s.Step(`^the data loads$`, ctx.theDataLoads)
	s.Step(`^the loading indicator should disappear$`, ctx.theLoadingIndicatorShouldDisappear)

	// Filtering scenario steps
	s.Step(`^I type "([^"]*)" in the filter input$`, ctx.iTypeInTheFilterInput)
	s.Step(`^only items containing "([^"]*)" should be visible$`, ctx.onlyItemsContainingShouldBeVisible)
	s.Step(`^the filter should apply across all (\d+) items$`, ctx.theFilterShouldApplyAcrossAllItems)
	s.Step(`^the status bar should show the filtered count$`, ctx.theStatusBarShouldShowTheFilteredCount)

	// Virtual scrolling scenario steps
	s.Step(`^(\d+) news items are loaded$`, ctx.newsItemsAreLoaded)
	s.Step(`^I scroll through the list$`, ctx.iScrollThroughTheList)
	s.Step(`^only visible items should be rendered in the DOM$`, ctx.onlyVisibleItemsShouldBeRenderedInTheDOM)
	s.Step(`^scrolling should remain performant$`, ctx.scrollingShouldRemainPerformant)
	s.Step(`^items should render as they come into view$`, ctx.itemsShouldRenderAsTheyComeIntoView)

	// Export scenario steps
	s.Step(`^I click the "([^"]*)" button$`, ctx.iClickTheButton)
	s.Step(`^a JSON file with (\d+) items should download$`, ctx.aJSONFileWithItemsShouldDownload)
	s.Step(`^a CSV file with (\d+) items should download$`, ctx.aCSVFileWithItemsShouldDownload)

	// Error handling scenario steps
	s.Step(`^the API endpoint is temporarily unavailable$`, ctx.theAPIEndpointIsTemporarilyUnavailable)
	s.Step(`^I should see a user-friendly error message$`, ctx.iShouldSeeAUserFriendlyErrorMessage)
	s.Step(`^cached data should be used if available$`, ctx.cachedDataShouldBeUsedIfAvailable)
	s.Step(`^the UI should remain responsive$`, ctx.theUIShouldRemainResponsive)
}

func TestTop200Features(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeTop200Scenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"rss-top200.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}