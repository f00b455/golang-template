# Issue: #8
# URL: https://github.com/f00b455/golang-template/issues/8
@pkg(api) @issue-8
Feature: SPIEGEL RSS API Endpoints
  As a user
  I want to fetch RSS headlines from SPIEGEL
  So that I can get the latest news via API

  Background:
    Given the API server is running

  @happy-path
  Scenario: Get latest SPIEGEL RSS headline
    When I make a GET request to "/api/rss/spiegel/latest"
    Then the response status should be 200
    And the response should contain a valid RSS headline
    And the headline should have a title
    And the headline should have a link
    And the headline should have a publishedAt timestamp
    And the headline should have source "SPIEGEL"

  @happy-path
  Scenario: Get top 5 SPIEGEL RSS headlines
    When I make a GET request to "/api/rss/spiegel/top5"
    Then the response status should be 200
    And the response should contain a headlines array
    And the headlines array should have 5 or fewer items
    And each headline should have title, link, publishedAt, and source fields

  @parameterized
  Scenario: Get limited number of headlines
    When I make a GET request to "/api/rss/spiegel/top5?limit=3"
    Then the response status should be 200
    And the response should contain a headlines array
    And the headlines array should have exactly 3 items

  @edge-cases
  Scenario: Request with limit of 1
    When I make a GET request to "/api/rss/spiegel/top5?limit=1"
    Then the response status should be 200
    And the response should contain a headlines array
    And the headlines array should have exactly 1 item

  @edge-cases
  Scenario: Request with maximum limit
    When I make a GET request to "/api/rss/spiegel/top5?limit=5"
    Then the response status should be 200
    And the response should contain a headlines array
    And the headlines array should have 5 or fewer items

  @error-handling
  Scenario: Request with invalid limit (too high)
    When I make a GET request to "/api/rss/spiegel/top5?limit=10"
    Then the response status should be 200
    And the response should contain a headlines array
    And the headlines array should have 5 or fewer items

  @error-handling
  Scenario: Request with invalid limit (zero)
    When I make a GET request to "/api/rss/spiegel/top5?limit=0"
    Then the response status should be 200
    And the response should contain a headlines array
    And the headlines array should have 1 or more items