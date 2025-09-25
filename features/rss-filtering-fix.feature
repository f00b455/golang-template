# Issue: #15
# URL: https://github.com/f00b455/golang-template/issues/15
@api @issue-15
Feature: RSS Filtering searches entire feed
  As an API user
  I want to filter RSS headlines from the entire feed
  So that I can find all matching articles, not just those in the first 5

  Background:
    Given the API server is running
    And the RSS feed contains 50+ headlines with various keywords

  @bug-fix
  Scenario: Filter finds matches beyond the first 5 items
    Given the first 5 RSS headlines do not contain the word "test-keyword-xyz"
    And items 10-15 contain headlines with "test-keyword-xyz"
    When I make a GET request to "/api/rss/spiegel/top5?filter=test-keyword-xyz"
    Then the response status should be 200
    And the response should contain matching headlines
    And the headlines should contain "test-keyword-xyz" in their titles

  @large-dataset
  Scenario: API fetches sufficient items before filtering
    Given the RSS feed has 100 items total
    When I make a GET request to "/api/rss/spiegel/top5?filter=rare-keyword"
    Then the API should fetch at least 50 items from the RSS feed
    And apply the filter to all fetched items
    And return up to 5 matching results

  @performance
  Scenario: Filtering performance with larger dataset
    Given the RSS feed has 100 items
    When I make a GET request to "/api/rss/spiegel/top5?filter=common-word"
    Then the response should be returned within 3 seconds
    And the response should contain up to 5 filtered results

  @edge-case
  Scenario: Filter returns empty when no matches in entire feed
    Given the RSS feed contains 50 items
    And none of the items contain "impossible-keyword-xyz123"
    When I make a GET request to "/api/rss/spiegel/top5?filter=impossible-keyword-xyz123"
    Then the response status should be 200
    And the headlines array should be empty
    And the totalCount should reflect the total fetched items

  @cache-behavior
  Scenario: Cache stores more items for better filtering
    Given the cache is empty
    When I make a GET request to "/api/rss/spiegel/top5"
    Then the cache should store at least 50 headlines
    When I make a subsequent request with filter "specific-word"
    Then the filter should be applied to all cached items
    And no new RSS fetch should occur if within cache TTL