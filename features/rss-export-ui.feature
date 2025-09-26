# Issue: #20
# URL: https://github.com/f00b455/golang-template/issues/20
@frontend @issue-20
Feature: RSS Export Download Links in Terminal UI
  As a user
  I want to export RSS headlines from the terminal UI
  So that I can download and analyze news data in different formats

  Background:
    Given I have the terminal-themed frontend running at "http://localhost:8080"
    And the RSS API export endpoints are available
    And RSS feed items are displayed in the UI

  @happy-path
  Scenario: Export RSS data as JSON
    Given I am viewing RSS headlines in the terminal UI
    When I click the "Download as JSON" button
    Then a JSON file should be downloaded to my computer
    And the file should contain the current RSS headlines
    And the filename should include "rss-export" and timestamp

  @csv-export
  Scenario: Export RSS data as CSV
    Given I am viewing RSS headlines in the terminal UI
    When I click the "Download as CSV" button
    Then a CSV file should be downloaded to my computer
    And the file should contain the current RSS headlines in CSV format
    And the filename should include "rss-export" and timestamp

  @export-with-filter
  Scenario: Export filtered RSS data
    Given I am viewing RSS headlines in the terminal UI
    And I have applied a filter "+tech -politics"
    When I click the "Download as JSON" button
    Then the downloaded file should only contain filtered items
    And the filename should include the filter text

  @export-limit
  Scenario: Export with item limit
    Given I am viewing RSS headlines in the terminal UI
    When I select "Export first 10 items" option
    And I click the "Download as JSON" button
    Then the downloaded file should contain exactly 10 items

  @export-via-command
  Scenario: Export using terminal command
    Given I am on the terminal-themed frontend
    When I type ":export json" in the command input
    Then a JSON file should be downloaded
    When I type ":export csv" in the command input
    Then a CSV file should be downloaded

  @export-error-handling
  Scenario: Handle export errors gracefully
    Given I am viewing RSS headlines in the terminal UI
    And the export API is temporarily unavailable
    When I click the "Download as JSON" button
    Then I should see an error message "Export failed. Please try again."
    And no file should be downloaded

  @export-ui-placement
  Scenario: Export buttons are properly positioned
    Given I am viewing RSS headlines in the terminal UI
    Then I should see export buttons near the RSS feed container
    And the buttons should have terminal-style theming (green text, black background)
    And the buttons should show tooltips on hover

  @export-progress
  Scenario: Show export progress for large datasets
    Given I am viewing RSS headlines with more than 100 items
    When I click the "Download as JSON" button
    Then I should see a progress indicator "Exporting..."
    And the progress indicator should disappear when download completes

  @export-keyboard-shortcut
  Scenario: Export using keyboard shortcuts
    Given I am viewing RSS headlines in the terminal UI
    When I press "Ctrl+E" followed by "J"
    Then a JSON export should be triggered
    When I press "Ctrl+E" followed by "C"
    Then a CSV export should be triggered

  @export-mobile
  Scenario: Export works on mobile devices
    Given I am on a mobile device
    When I tap the export button
    Then the file should be downloadable on the mobile browser
    And the UI should remain responsive during export