# Issue: #7
# URL: https://github.com/f00b455/golang-template/issues/7
@pkg(handlers) @issue-7
Feature: RSS Headline Filtering
  As a user
  I want to filter RSS headlines by text
  So that I can find relevant news quickly

  Background:
    Given the API server is running
    And the RSS feed contains multiple headlines

  @happy-path
  Scenario: Filter latest headline with matching text
    When I make a GET request to "/api/rss/spiegel/latest?filter=headline"
    Then the response status should be 200
    And the headline title should contain "headline" case-insensitively

  @happy-path
  Scenario: Filter top 5 headlines with matching text
    When I make a GET request to "/api/rss/spiegel/top5?filter=news"
    Then the response status should be 200
    And all headlines should contain "news" case-insensitively
    And the headlines array should have 5 or fewer items

  @happy-path
  Scenario: Filter with limit parameter
    When I make a GET request to "/api/rss/spiegel/top5?filter=headline&limit=3"
    Then the response status should be 200
    And all headlines should contain "headline" case-insensitively
    And the headlines array should have exactly 3 items or fewer

  @edge-cases
  Scenario: Filter with no matches
    When I make a GET request to "/api/rss/spiegel/top5?filter=xyznonexistent"
    Then the response status should be 200
    And the headlines array should be empty

  @edge-cases
  Scenario: Filter latest with no match
    When I make a GET request to "/api/rss/spiegel/latest?filter=xyznonexistent"
    Then the response status should be 200
    And the response should be an empty object or null headline

  @case-insensitive
  Scenario: Case-insensitive filtering
    When I make a GET request to "/api/rss/spiegel/top5?filter=HEADLINE"
    Then the response status should be 200
    And all headlines should contain "headline" case-insensitively

  @filter-before-limit
  Scenario: Filter applied before limit
    Given the RSS feed has 10 headlines with "tech" in title
    And the RSS feed has 10 headlines without "tech" in title
    When I make a GET request to "/api/rss/spiegel/top5?filter=tech&limit=5"
    Then the response status should be 200
    And the headlines array should have exactly 5 items
    And all headlines should contain "tech" case-insensitively

  @empty-filter
  Scenario: Empty filter parameter returns all results
    When I make a GET request to "/api/rss/spiegel/top5?filter="
    Then the response status should be 200
    And the headlines array should have 5 or fewer items
    And no filtering should be applied