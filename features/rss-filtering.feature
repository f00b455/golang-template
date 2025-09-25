# Issue: #5
# URL: https://github.com/f00b455/golang-template/issues/5
@pkg(api) @pkg(web) @issue-5
Feature: RSS Headline Filtering
  As a web application user
  I want to filter RSS headlines by keywords
  So that I can see only the news articles that match my interests

  Background:
    Given the API server is running

  @happy-path
  Scenario: Filter headlines with matching keyword
    Given the RSS feed has 10 articles
    And 3 articles contain the word "Politik"
    When I request top 5 articles with filter "Politik"
    Then I should receive exactly 3 articles
    And all returned articles should contain "Politik" in their headline

  @happy-path
  Scenario: Filter returns exact count when more matches exist
    Given the RSS feed has 20 articles
    And 8 articles contain the word "Wirtschaft"
    When I request top 5 articles with filter "Wirtschaft"
    Then I should receive exactly 5 articles
    And all returned articles should contain "Wirtschaft" in their headline

  @case-sensitivity
  Scenario: Case-insensitive filtering
    Given the RSS feed has articles with "COVID", "covid", and "Covid"
    When I request articles with filter "covid"
    Then I should receive all articles regardless of case

  @edge-cases
  Scenario: No matches for filter
    Given the RSS feed has 10 articles
    And no articles contain the word "Blockchain"
    When I request top 5 articles with filter "Blockchain"
    Then I should receive an empty result set

  @edge-cases
  Scenario: Empty filter returns all articles
    Given the RSS feed has 10 articles
    When I request top 5 articles with filter ""
    Then I should receive exactly 5 articles

  @integration
  Scenario: Filter with special characters
    Given the RSS feed has articles with "EU-Parlament" and "EU Parliament"
    When I request articles with filter "EU"
    Then I should receive all articles containing "EU"

  @performance
  Scenario: Filter with limit parameter
    Given the RSS feed has 15 articles
    And 10 articles contain the word "Deutschland"
    When I request top 3 articles with filter "Deutschland"
    Then I should receive exactly 3 articles
    And all returned articles should contain "Deutschland" in their headline

  @api-endpoint
  Scenario: Filter via API endpoint with query parameter
    When I make a GET request to "/api/rss/spiegel/top5?filter=Politik"
    Then the response status should be 200
    And the response should contain a headlines array
    And all headlines should contain "Politik" case-insensitively

  @api-endpoint
  Scenario: Latest endpoint with filter
    When I make a GET request to "/api/rss/spiegel/latest?filter=Wirtschaft"
    Then the response status should be 200
    And the headline should contain "Wirtschaft" case-insensitively