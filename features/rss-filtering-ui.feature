# Issue: #9
# URL: https://github.com/f00b455/golang-template/issues/9
@web @issue-9
Feature: RSS Filtering UI in Web Interface
  As a web application user
  I want to filter RSS headlines in the web interface
  So that I can quickly find articles that interest me

  Background:
    Given the web server is running
    And the RSS feed contains multiple headlines

  @filter-input
  Scenario: Filter input field is present on the web page
    When I navigate to the homepage
    Then I should see a filter input field above the headlines
    And the input should have placeholder text "Filter headlines... (e.g., Politik, Wirtschaft)"

  @realtime-filtering
  Scenario: Real-time filtering of headlines
    Given I am on the homepage
    And there are 10 headlines displayed
    When I type "Politik" in the filter input field
    Then only headlines containing "Politik" should be displayed
    And I should see the filtered results count "Showing 3 of 10 matching articles"

  @clear-filter
  Scenario: Clear filter functionality
    Given I am on the homepage
    And I have filtered headlines with keyword "Wirtschaft"
    When I click the clear filter button
    Then all headlines should be displayed again
    And the filter input field should be empty

  @api-integration
  Scenario: Filter parameter is sent to API endpoint
    Given I am on the homepage
    When I type "Sport" in the filter input field
    Then the API request should include filter parameter "Sport"
    And the refreshHeadlines function should pass the filter parameter

  @mobile-responsive
  Scenario: Filter UI is mobile-responsive
    Given I am viewing the site on a mobile device
    When I navigate to the homepage
    Then the filter input field should be responsive
    And the clear button should be easily tappable

  @empty-filter-results
  Scenario: Handle empty filter results
    Given I am on the homepage
    When I type "NonExistentKeyword" in the filter input field
    Then I should see a message "No headlines match your filter"
    And the filtered results count should show "Showing 0 of 10 matching articles"