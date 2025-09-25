# Issue: #7
# URL: https://github.com/f00b455/ts-pure-template/issues/7
@pkg(api) @issue-7
Feature: SPIEGEL RSS Headline Display
  As a User
  I want to see the latest SPIEGEL RSS headline in the website header
  So that I can get the current news at a glance and jump directly to the article

  Background:
    Given the API server is running
    And the SPIEGEL RSS feed is available

  Scenario: Happy Path - Display latest headline
    Given the page loads
    When the header component initializes
    Then it calls the API endpoint "/api/rss/spiegel/latest"
    And displays the title and link in the header
    And clicking the link opens the article in a new tab

  Scenario: Loading State
    Given the API endpoint has not responded yet
    When the header component is mounting
    Then a subtle loading state is displayed
    And no layout shifts occur

  Scenario: Error Handling - API failure
    Given the API call fails or returns empty
    When the header component receives an error
    Then an unobtrusive fallback message "Gerade keine Schlagzeile verfügbar" is displayed
    And no UI errors or layout jumps are caused

  Scenario: Auto-Refresh
    Given the page is open for longer than 5 minutes
    When 5 minutes have passed
    Then the component automatically refreshes the headline
    And the new headline is displayed without page reload

  Scenario: Responsiveness
    Given a small viewport (mobile device)
    When the headline is displayed
    Then the title is cleanly truncated with ellipsis
    And the link remains clickable
    And the time indicator is hidden on very small screens

  Scenario: Accessibility
    When the headline component is rendered
    Then the text is accessible via screen reader
    And updates use aria-live="polite"
    And the link has a descriptive name including "SPIEGEL:"

  Scenario: API Response Time
    When the API endpoint is called
    Then it responds within 600ms for 95th percentile of requests
    And implements caching with 5-minute TTL
    And retries once on timeout after 2 seconds

  Scenario: Dark Mode Support
    Given the user has dark mode enabled
    When the headline component is displayed
    Then it adapts to the dark color scheme
    And maintains readability and contrast