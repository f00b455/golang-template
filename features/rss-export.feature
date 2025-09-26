# Issue: #11
# URL: https://github.com/f00b455/golang-template/issues/11
@pkg(api) @pkg(web) @issue-11
Feature: Export RSS Headlines
  As a data analyst
  I want to export RSS headlines in various formats
  So that I can analyze news data offline

  Background:
    Given the RSS feed has multiple articles available
    And the API server is running

  @export-json @happy-path
  Scenario: Export headlines as JSON
    When I request "/api/rss/spiegel/export?format=json"
    Then the response status should be 200
    And the content-type should be "application/json"
    And the content-disposition should contain "attachment"
    And the filename should contain ".json"
    And the JSON response should have export metadata
    And the JSON response should have headlines array

  @export-csv @happy-path
  Scenario: Export headlines as CSV
    When I request "/api/rss/spiegel/export?format=csv"
    Then the response status should be 200
    And the content-type should be "text/csv"
    And the content-disposition should contain "attachment"
    And the filename should contain ".csv"
    And the CSV should have header row
    And the CSV should have data rows

  @export-filtered
  Scenario: Export filtered headlines as JSON
    When I request "/api/rss/spiegel/export?format=json&filter=Politik"
    Then the response status should be 200
    And the content-type should be "application/json"
    And the JSON response should contain filter metadata "Politik"
    And all headlines should match the filter "Politik"

  @export-limited
  Scenario: Export limited number of headlines
    When I request "/api/rss/spiegel/export?format=json&limit=3"
    Then the response status should be 200
    And the content-type should be "application/json"
    And the JSON response should have exactly 3 headlines

  @export-csv-filtered
  Scenario: Export filtered headlines as CSV
    When I request "/api/rss/spiegel/export?format=csv&filter=Sport"
    Then the response status should be 200
    And the content-type should be "text/csv"
    And the CSV rows should only contain "Sport" headlines

  @error-handling
  Scenario: Export with invalid format
    When I request "/api/rss/spiegel/export?format=xml"
    Then the response status should be 400
    And the response should contain an error message about invalid format

  @error-handling
  Scenario: Export without format parameter
    When I request "/api/rss/spiegel/export"
    Then the response status should be 400
    And the response should contain an error message about missing format

  @edge-case
  Scenario: Export empty result set
    When I request "/api/rss/spiegel/export?format=json&filter=NONEXISTENT12345"
    Then the response status should be 200
    And the content-type should be "application/json"
    And the JSON response should have empty headlines array
    And the JSON response should show total_items as 0

  @export-csv-special-chars
  Scenario: Export CSV with special characters
    When I request "/api/rss/spiegel/export?format=csv"
    Then the response status should be 200
    And the CSV should properly escape quotes and commas
    And the CSV should handle UTF-8 characters correctly

  @export-combined-filters
  Scenario: Export with filter and limit combined
    When I request "/api/rss/spiegel/export?format=json&filter=News&limit=2"
    Then the response status should be 200
    And the JSON response should have at most 2 headlines
    And all headlines should match the filter "News"