# Issue: #7
# URL: https://github.com/f00b455/golang-template/issues/7
@pkg(web) @issue-7
Feature: Header shows top 5 SPIEGEL headlines with dates
  So that I can quickly see multiple current news items,
  I want to see up to 5 headlines with dates in the header.

  Background:
    Given the application is open
    And the header is visible

  @happy-path
  Scenario: Display up to 5 latest headlines with dates
    Given the API returns at least 5 entries
    When the header news component initializes
    Then exactly 5 headlines should be displayed
    And each headline should show title and publication date
    And the list should be sorted by date (newest first)
    And clicking a headline should open the article in a new tab

  @less-than-5
  Scenario: Less than 5 entries available
    Given the API returns 3 entries
    When the component loads
    Then exactly 3 headlines should be displayed
    And no placeholders should be shown for missing entries

  @date-format
  Scenario: Date format and timezone
    Given an entry has UTC date "2025-09-24T08:05:00Z"
    When the component renders
    Then the date should display as "24.09.2025 10:05" in Europe/Berlin timezone

  @refresh
  Scenario: Automatic refresh
    Given the page has been open for a while
    When 5 minutes have passed
    Then the list should be refreshed via API
    And the order and count should remain consistent

  @error
  Scenario: Error fallback
    Given the API call fails or returns empty
    When the component renders
    Then a subtle fallback message should appear
    And there should be no layout jumps