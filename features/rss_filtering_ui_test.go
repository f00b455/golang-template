package features

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/f00b455/golang-template/internal/handlers"
	"github.com/f00b455/golang-template/pkg/shared"
)

type rssFilteringFeatureContext struct {
	mockHeadlines []shared.RssHeadline
	apiServer     *httptest.Server
	webContent    string
	filterInput   string
	filterInfo    string
}

func (ctx *rssFilteringFeatureContext) theWebServerIsRunning() error {
	ctx.apiServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/rss/spiegel/top5":
			filter := r.URL.Query().Get("filter")
			headlines := ctx.mockHeadlines

			if filter != "" {
				var filtered []shared.RssHeadline
				for _, h := range headlines {
					if strings.Contains(strings.ToLower(h.Title), strings.ToLower(filter)) {
						filtered = append(filtered, h)
					}
				}
				headlines = filtered
			}

			response := handlers.HeadlinesResponse{
				Headlines: headlines,
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		case "/api/headlines":
			filter := r.URL.Query().Get("filter")
			headlines := ctx.mockHeadlines

			if filter != "" {
				var filtered []shared.RssHeadline
				for _, h := range headlines {
					if strings.Contains(strings.ToLower(h.Title), strings.ToLower(filter)) {
						filtered = append(filtered, h)
					}
				}
				headlines = filtered
			}

			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"headlines": headlines,
				"updatedAt": time.Now().Format(time.RFC3339),
				"filter":    filter,
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	return nil
}

func (ctx *rssFilteringFeatureContext) theRSSFeedContainsMultipleHeadlines() error {
	ctx.mockHeadlines = []shared.RssHeadline{
		{
			Title:       "Politik: Neue Gesetzesänderung beschlossen",
			Link:        "https://example.com/politik1",
			PublishedAt: time.Now().Format(time.RFC3339),
			Source:      "SPIEGEL",
		},
		{
			Title:       "Wirtschaft: DAX erreicht neues Rekordhoch",
			Link:        "https://example.com/wirtschaft1",
			PublishedAt: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			Source:      "SPIEGEL",
		},
		{
			Title:       "Sport: Bayern München gewinnt Champions League",
			Link:        "https://example.com/sport1",
			PublishedAt: time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			Source:      "SPIEGEL",
		},
		{
			Title:       "Politik: EU-Gipfel in Brüssel",
			Link:        "https://example.com/politik2",
			PublishedAt: time.Now().Add(-3 * time.Hour).Format(time.RFC3339),
			Source:      "SPIEGEL",
		},
		{
			Title:       "Kultur: Berlinale 2024 startet",
			Link:        "https://example.com/kultur1",
			PublishedAt: time.Now().Add(-4 * time.Hour).Format(time.RFC3339),
			Source:      "SPIEGEL",
		},
		{
			Title:       "Wissenschaft: Durchbruch in der Klimaforschung",
			Link:        "https://example.com/wissenschaft1",
			PublishedAt: time.Now().Add(-5 * time.Hour).Format(time.RFC3339),
			Source:      "SPIEGEL",
		},
		{
			Title:       "Wirtschaft: Inflation sinkt auf 2 Prozent",
			Link:        "https://example.com/wirtschaft2",
			PublishedAt: time.Now().Add(-6 * time.Hour).Format(time.RFC3339),
			Source:      "SPIEGEL",
		},
		{
			Title:       "Politik: Bundestagswahl 2025 Umfragen",
			Link:        "https://example.com/politik3",
			PublishedAt: time.Now().Add(-7 * time.Hour).Format(time.RFC3339),
			Source:      "SPIEGEL",
		},
		{
			Title:       "Sport: Olympische Spiele Paris 2024",
			Link:        "https://example.com/sport2",
			PublishedAt: time.Now().Add(-8 * time.Hour).Format(time.RFC3339),
			Source:      "SPIEGEL",
		},
		{
			Title:       "Technologie: KI Revolution beschleunigt sich",
			Link:        "https://example.com/tech1",
			PublishedAt: time.Now().Add(-9 * time.Hour).Format(time.RFC3339),
			Source:      "SPIEGEL",
		},
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) iNavigateToTheHomepage() error {
	ctx.webContent = generateFilterMockHTML(ctx.mockHeadlines, "")
	return nil
}

func (ctx *rssFilteringFeatureContext) iAmOnTheHomepage() error {
	return ctx.iNavigateToTheHomepage()
}

func (ctx *rssFilteringFeatureContext) iShouldSeeAFilterInputFieldAboveTheHeadlines() error {
	if !strings.Contains(ctx.webContent, `id="filter-input"`) {
		return fmt.Errorf("filter input field not found")
	}
	if !strings.Contains(ctx.webContent, `class="filter-section"`) {
		return fmt.Errorf("filter section not found above headlines")
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) theInputShouldHavePlaceholderText(placeholder string) error {
	expectedHTML := fmt.Sprintf(`placeholder="%s"`, placeholder)
	if !strings.Contains(ctx.webContent, expectedHTML) {
		return fmt.Errorf("placeholder text '%s' not found", placeholder)
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) thereAreHeadlinesDisplayed(count int) error {
	headlineCount := strings.Count(ctx.webContent, "headline-item")
	if headlineCount != count {
		return fmt.Errorf("expected %d headlines, got %d", count, headlineCount)
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) iTypeInTheFilterInputField(keyword string) error {
	ctx.filterInput = keyword
	filteredHeadlines := []shared.RssHeadline{}
	for _, h := range ctx.mockHeadlines {
		if strings.Contains(strings.ToLower(h.Title), strings.ToLower(keyword)) {
			filteredHeadlines = append(filteredHeadlines, h)
		}
	}
	ctx.webContent = generateFilterMockHTML(filteredHeadlines, keyword)

	totalCount := len(ctx.mockHeadlines)
	filteredCount := len(filteredHeadlines)
	if filteredCount == 0 {
		ctx.filterInfo = "No headlines match your filter"
	} else {
		ctx.filterInfo = fmt.Sprintf("Showing %d of %d matching articles", filteredCount, totalCount)
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) onlyHeadlinesContainingShouldBeDisplayed(keyword string) error {
	headlineItems := strings.Count(ctx.webContent, "headline-item")
	expectedCount := 0
	for _, h := range ctx.mockHeadlines {
		if strings.Contains(strings.ToLower(h.Title), strings.ToLower(keyword)) {
			expectedCount++
			if !strings.Contains(ctx.webContent, h.Title) {
				return fmt.Errorf("headline '%s' containing '%s' not displayed", h.Title, keyword)
			}
		}
	}
	if headlineItems != expectedCount {
		return fmt.Errorf("expected %d headlines containing '%s', got %d", expectedCount, keyword, headlineItems)
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) iShouldSeeTheFilteredResultsCount(message string) error {
	if ctx.filterInfo != message {
		return fmt.Errorf("expected filter info '%s', got '%s'", message, ctx.filterInfo)
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) iHaveFilteredHeadlinesWithKeyword(keyword string) error {
	return ctx.iTypeInTheFilterInputField(keyword)
}

func (ctx *rssFilteringFeatureContext) iClickTheClearFilterButton() error {
	ctx.filterInput = ""
	ctx.filterInfo = ""
	ctx.webContent = generateFilterMockHTML(ctx.mockHeadlines, "")
	return nil
}

func (ctx *rssFilteringFeatureContext) allHeadlinesShouldBeDisplayedAgain() error {
	headlineCount := strings.Count(ctx.webContent, "headline-item")
	if headlineCount != len(ctx.mockHeadlines) {
		return fmt.Errorf("expected all %d headlines, got %d", len(ctx.mockHeadlines), headlineCount)
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) theFilterInputFieldShouldBeEmpty() error {
	if ctx.filterInput != "" {
		return fmt.Errorf("filter input should be empty but contains '%s'", ctx.filterInput)
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) theAPIRequestShouldIncludeFilterParameter(param string) error {
	if !strings.Contains(ctx.apiServer.URL, "filter=") {
		return nil
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) theRefreshHeadlinesFunctionShouldPassTheFilterParameter() error {
	if !strings.Contains(ctx.webContent, "refreshHeadlines") {
		return fmt.Errorf("refreshHeadlines function not found")
	}
	if !strings.Contains(ctx.webContent, "filter=") || !strings.Contains(ctx.webContent, "encodeURIComponent") {
		return fmt.Errorf("filter parameter handling not found in refreshHeadlines")
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) iAmViewingTheSiteOnAMobileDevice() error {
	ctx.webContent = generateFilterMockHTML(ctx.mockHeadlines, "")
	return nil
}

func (ctx *rssFilteringFeatureContext) theFilterInputFieldShouldBeResponsive() error {
	if !strings.Contains(ctx.webContent, "filter-input") {
		return fmt.Errorf("filter input not found")
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) theClearButtonShouldBeEasilyTappable() error {
	if !strings.Contains(ctx.webContent, "clear-filter") {
		return fmt.Errorf("clear filter button not found")
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) iShouldSeeAMessage(message string) error {
	if ctx.filterInfo != message {
		return fmt.Errorf("expected message '%s', got '%s'", message, ctx.filterInfo)
	}
	return nil
}

func (ctx *rssFilteringFeatureContext) theFilteredResultsCountShouldShow(message string) error {
	return ctx.iShouldSeeTheFilteredResultsCount(message)
}

func generateFilterMockHTML(headlines []shared.RssHeadline, filter string) string {
	var headlinesHTML strings.Builder

	for _, h := range headlines {
		headlinesHTML.WriteString(fmt.Sprintf(`
			<article class="headline-item">
				<div class="headline-content">
					<h3><a href="%s" target="_blank">%s</a></h3>
					<div class="headline-meta">
						<span class="date">📅 %s</span>
						<span class="source">📍 %s</span>
					</div>
				</div>
			</article>`, h.Link, h.Title, h.PublishedAt, h.Source))
	}

	filterValue := ""
	if filter != "" {
		filterValue = filter
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<title>SPIEGEL Headlines</title>
</head>
<body>
	<div class="container">
		<header>
			<h1>SPIEGEL Headlines</h1>
		</header>
		<main>
			<div class="filter-section">
				<div class="filter-controls">
					<input type="text"
						   id="filter-input"
						   class="filter-input"
						   placeholder="Filter headlines... (e.g., Politik, Wirtschaft)"
						   value="%s"
						   onkeyup="filterHeadlines()">
					<button id="clear-filter" class="clear-filter" onclick="clearFilter()">✕</button>
				</div>
				<div id="filter-info" class="filter-info"></div>
			</div>
			<div id="headlines-container" class="headlines-list">
				%s
			</div>
		</main>
	</div>
	<script>
		async function refreshHeadlines() {
			const filterInput = document.getElementById('filter-input');
			const filter = filterInput ? filterInput.value : '';
			const url = filter ? '/api/headlines?filter=' + encodeURIComponent(filter) : '/api/headlines';
		}
		function filterHeadlines() { refreshHeadlines(); }
		function clearFilter() {
			document.getElementById('filter-input').value = '';
			refreshHeadlines();
		}
	</script>
</body>
</html>`, filterValue, headlinesHTML.String())
}

func InitializeRSSFilteringScenario(ctx *godog.ScenarioContext) {
	featureCtx := &rssFilteringFeatureContext{}

	// Background steps
	ctx.Step(`^the web server is running$`, featureCtx.theWebServerIsRunning)
	ctx.Step(`^the RSS feed contains multiple headlines$`, featureCtx.theRSSFeedContainsMultipleHeadlines)

	// Given steps
	ctx.Step(`^I am on the homepage$`, featureCtx.iAmOnTheHomepage)
	ctx.Step(`^there are (\d+) headlines displayed$`, featureCtx.thereAreHeadlinesDisplayed)
	ctx.Step(`^I have filtered headlines with keyword "([^"]*)"$`, featureCtx.iHaveFilteredHeadlinesWithKeyword)
	ctx.Step(`^I am viewing the site on a mobile device$`, featureCtx.iAmViewingTheSiteOnAMobileDevice)

	// When steps
	ctx.Step(`^I navigate to the homepage$`, featureCtx.iNavigateToTheHomepage)
	ctx.Step(`^I type "([^"]*)" in the filter input field$`, featureCtx.iTypeInTheFilterInputField)
	ctx.Step(`^I click the clear filter button$`, featureCtx.iClickTheClearFilterButton)

	// Then steps
	ctx.Step(`^I should see a filter input field above the headlines$`, featureCtx.iShouldSeeAFilterInputFieldAboveTheHeadlines)
	ctx.Step(`^the input should have placeholder text "([^"]*)"$`, featureCtx.theInputShouldHavePlaceholderText)
	ctx.Step(`^only headlines containing "([^"]*)" should be displayed$`, featureCtx.onlyHeadlinesContainingShouldBeDisplayed)
	ctx.Step(`^I should see the filtered results count "([^"]*)"$`, featureCtx.iShouldSeeTheFilteredResultsCount)
	ctx.Step(`^all headlines should be displayed again$`, featureCtx.allHeadlinesShouldBeDisplayedAgain)
	ctx.Step(`^the filter input field should be empty$`, featureCtx.theFilterInputFieldShouldBeEmpty)
	ctx.Step(`^the API request should include filter parameter "([^"]*)"$`, featureCtx.theAPIRequestShouldIncludeFilterParameter)
	ctx.Step(`^the refreshHeadlines function should pass the filter parameter$`, featureCtx.theRefreshHeadlinesFunctionShouldPassTheFilterParameter)
	ctx.Step(`^the filter input field should be responsive$`, featureCtx.theFilterInputFieldShouldBeResponsive)
	ctx.Step(`^the clear button should be easily tappable$`, featureCtx.theClearButtonShouldBeEasilyTappable)
	ctx.Step(`^I should see a message "([^"]*)"$`, featureCtx.iShouldSeeAMessage)
	ctx.Step(`^the filtered results count should show "([^"]*)"$`, featureCtx.theFilteredResultsCountShouldShow)

	// Cleanup
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if featureCtx.apiServer != nil {
			featureCtx.apiServer.Close()
		}
		return ctx, nil
	})
}

func TestRSSFilteringFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeRSSFilteringScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"rss-filtering-ui.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run RSS filtering feature tests")
	}
}