# Issue: #24
# URL: https://github.com/f00b455/golang-template/issues/24
@rss @frontend @issue-24
Feature: Display top 200 news items in terminal UI
  As a user
  I want to see more than 5 news items in the terminal interface
  So that I can browse through a larger selection of current headlines

  Background:
    Given I have a running API server
    And the terminal UI is accessible
    And the RSS feed has at least 200 items available

  @api-limit
  Scenario: API endpoint supports up to 200 items
    When I request "/api/rss/spiegel/top5?limit=200"
    Then the response should contain up to 200 RSS items
    And the response status should be 200
    And the response should include a totalCount field

  @api-backward-compatibility
  Scenario: API maintains backward compatibility
    When I request "/api/rss/spiegel/top5" without limit parameter
    Then the response should contain exactly 5 RSS items
    And the response status should be 200

  @api-validation
  Scenario Outline: API validates limit parameter
    When I request "/api/rss/spiegel/top5?limit=<limit>"
    Then the response should contain <expected_items> RSS items
    And the response status should be <status>

    Examples:
      | limit | expected_items | status |
      | 1     | 1              | 200    |
      | 50    | 50             | 200    |
      | 100   | 100            | 200    |
      | 200   | 200            | 200    |
      | 201   | 200            | 200    |
      | -1    | 5              | 200    |
      | abc   | 5              | 200    |

  @ui-pagination
  Scenario: Terminal UI displays paginated results
    Given the API returns 200 news items
    When I load the terminal UI
    Then I should see the first page of news items
    And I should see pagination controls
    And the status bar should show "1-20 of 200"

  @ui-keyboard-navigation
  Scenario: Navigate pages with keyboard shortcuts
    Given 200 news items are loaded in the terminal UI
    And I am on page 1
    When I press "Page Down" key
    Then I should be on page 2
    And the status bar should show "21-40 of 200"
    When I press "Page Up" key
    Then I should be on page 1
    And the status bar should show "1-20 of 200"

  @ui-jump-navigation
  Scenario: Jump to specific pages
    Given 200 news items are loaded in the terminal UI
    When I type ":page 5" in the command input
    Then I should be on page 5
    And the status bar should show "81-100 of 200"

  @performance
  Scenario: Performance with 200 items
    Given the API returns 200 news items
    When I load the terminal UI
    Then the initial page should load in less than 2 seconds
    And scrolling should be smooth with no lag
    And memory usage should remain reasonable

  @loading-indicator
  Scenario: Show loading indicator for large data sets
    When I request 200 news items
    Then I should see a loading indicator
    And the loading message should show progress
    When the data loads
    Then the loading indicator should disappear

  @filtering-large-dataset
  Scenario: Filter works with 200 items
    Given 200 news items are loaded in the terminal UI
    When I type "technology" in the filter input
    Then only items containing "technology" should be visible
    And the filter should apply across all 200 items
    And the status bar should show the filtered count

  @virtual-scrolling
  Scenario: Virtual scrolling for performance
    Given 200 news items are loaded
    When I scroll through the list
    Then only visible items should be rendered in the DOM
    And scrolling should remain performant
    And items should render as they come into view

  @export-large-dataset
  Scenario: Export 200 items
    Given 200 news items are loaded
    When I click the "Export JSON" button
    Then a JSON file with 200 items should download
    When I click the "Export CSV" button
    Then a CSV file with 200 items should download

  @error-handling
  Scenario: Handle API errors gracefully
    Given the API endpoint is temporarily unavailable
    When I request 200 items
    Then I should see a user-friendly error message
    And cached data should be used if available
    And the UI should remain responsive