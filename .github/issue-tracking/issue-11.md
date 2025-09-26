# Issue #11: [STORY] As a user, I want to export/download RSS headlines as CSV/JSON so that I can analyze or archive news data

**Issue URL**: https://github.com/f00b455/golang-template/issues/11
**Created**: 2025-09-25T16:07:50Z
**Assignee**: Unassigned

## Description
## User Story
**As a** data analyst or researcher
**I want** to export RSS headlines in different formats (CSV, JSON)
**So that** I can analyze news trends offline or integrate with other tools

## Background/Context
Users currently can only view RSS headlines in the web interface. They need the ability to export this data for:
- Offline analysis in Excel/Google Sheets
- Integration with data analysis tools
- Creating personal news archives
- Sharing curated news collections
- Academic research on news trends

## BDD Feature File (REQUIRED - Write FIRST!)
**⚠️ IMPORTANT: Feature files MUST be written BEFORE any implementation code!**

### Red-Green-Refactor Cycle
- [ ] **RED Phase**: Feature file written with failing scenarios
- [ ] **GREEN Phase**: Minimal implementation to pass tests
- [ ] **REFACTOR Phase**: Code quality improvements

### Example Feature File Content
```gherkin
# Issue: #10
# URL: https://github.com/f00b455/golang-template/issues/10
@pkg(api) @pkg(web) @issue-10
Feature: Export RSS Headlines
  As a data analyst
  I want to export RSS headlines in various formats
  So that I can analyze news data offline

  Background:
    Given the RSS feed has 10 articles available
    And the API server is running

  @export-json
  Scenario: Export headlines as JSON
    When I request "/api/rss/spiegel/export?format=json"
    Then I should receive a JSON file download
    And the content-type should be "application/json"
    And the filename should contain timestamp
    And all headline fields should be included

  @export-csv
  Scenario: Export headlines as CSV
    When I request "/api/rss/spiegel/export?format=csv"
    Then I should receive a CSV file download
    And the content-type should be "text/csv"
    And the CSV should have headers
    And special characters should be properly escaped
    And dates should be in ISO format

  @export-filtered
  Scenario: Export filtered headlines
    Given articles contain both "Politik" and "Sport" topics
    When I request "/api/rss/spiegel/export?format=csv&filter=Politik"
    Then only "Politik" articles should be in the export
    And the filename should include the filter term
    And the export should maintain proper formatting

  @export-limited
  Scenario: Export with limit parameter
    When I request "/api/rss/spiegel/export?format=json&limit=5"
    Then exactly 5 articles should be exported
    And they should be the most recent articles
    And the export should be properly formatted

  @export-ui
  Scenario: Export button in web interface
    Given I am on the RSS headlines page
    When I click the "Export" button
    Then I should see format options (CSV, JSON)
    And current filters should be applied to export
    And the download should start immediately
```

## Technical Requirements

### API Implementation:
- [ ] New endpoint: `GET /api/rss/spiegel/export`
- [ ] Query parameters:
  - `format`: "json" | "csv" (required)
  - `filter`: keyword filter (optional)
  - `limit`: number of items (optional, default: all)
  - `sort`: "date" | "title" (optional, default: date)

### Export Formats:

**CSV Format:**
```csv
Title,Link,Published_At,Source,Category
"Article Title","https://example.com","2024-01-01T10:00:00Z","SPIEGEL","Politics"
```

**JSON Format:**
```json
{
  "export_date": "2024-01-01T12:00:00Z",
  "total_items": 10,
  "filter_applied": "Politik",
  "headlines": [
    {
      "title": "Article Title",
      "link": "https://example.com",
      "publishedAt": "2024-01-01T10:00:00Z",
      "source": "SPIEGEL",
      "category": "Politics"
    }
  ]
}
```

### Frontend Changes:
- [ ] Add export button/dropdown to UI
- [ ] Format selection (CSV/JSON)
- [ ] Apply current filters to export
- [ ] Loading indicator during export generation
- [ ] Success/error notifications

### Implementation Details:
- [ ] Proper CSV escaping (quotes, commas, newlines)
- [ ] UTF-8 encoding for international characters
- [ ] Content-Disposition headers for downloads
- [ ] Filename generation with timestamp
- [ ] Memory-efficient streaming for large exports

## Acceptance Criteria
- [ ] Export functionality accessible from web UI
- [ ] CSV exports open correctly in Excel/Google Sheets
- [ ] JSON exports are valid and parseable
- [ ] Filtered exports contain only matching items
- [ ] Export includes metadata (timestamp, filter, count)
- [ ] Large exports (100+ items) work efficiently
- [ ] Proper error handling for edge cases

## Edge Cases to Handle
- [ ] Empty result sets (no matching filters)
- [ ] Special characters in headlines (quotes, unicode)
- [ ] Very long headlines or links
- [ ] Network interruptions during download
- [ ] Invalid format parameter handling

## Definition of Done
- [ ] BDD scenarios all passing
- [ ] Unit tests for export logic
- [ ] Integration tests for API endpoint
- [ ] UI export button functional
- [ ] Documentation updated
- [ ] Code reviewed and approved
- [ ] Performance tested with 100+ items

## Dependencies
- Existing RSS feed infrastructure
- Current filtering functionality
- Web UI framework

## Story Points/Estimation
- Estimated effort: M (Medium)
- Story points: 5
- Timeline: 1 sprint

## Success Metrics
- Export functionality used by >10% of users
- Zero data corruption in exports
- Export generation time <2 seconds for 100 items

## Work Log
- Branch created: issue-11-story-as-a-user-i-want-to-export-download-rss-head
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
