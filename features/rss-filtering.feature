# Issue: #5
# URL: https://github.com/f00b455/golang-template/issues/5
@pkg(api) @issue-5
Feature: RSS Headline Filtering
  As a web application user
  I want to filter RSS headlines by keywords
  So that I can see only the news articles that match my interests

  Background:
    Given the API server is running

  @happy-path
  Scenario: Filter headlines with matching keyword
    Given the RSS feed has these articles:
      | title                              |
      | Politik: Neue Gesetzesänderungen   |
      | Wirtschaft: DAX erreicht Rekordhoch|
      | Politik: Wahlen in Deutschland     |
      | Sport: Bundesliga Ergebnisse      |
      | Politik: EU-Gipfel in Brüssel      |
      | Kultur: Neue Ausstellung eröffnet  |
    When I make a GET request to "/api/rss/spiegel/top5?filter=Politik"
    Then the response status should be 200
    And the response should contain a headlines array
    And the headlines array should have exactly 3 items
    And all headlines should contain "Politik" in their title

  @happy-path
  Scenario: Filter returns exact count when more matches exist
    Given the RSS feed has these articles:
      | title                                      |
      | Wirtschaft: DAX steigt                     |
      | Politik: Neue Regierung                    |
      | Wirtschaft: Euro-Kurs                      |
      | Wirtschaft: Arbeitslosenzahlen            |
      | Sport: Champions League                    |
      | Wirtschaft: Inflation steigt              |
      | Wirtschaft: Börse schließt höher          |
      | Wirtschaft: Exporte gestiegen             |
      | Wirtschaft: BIP-Wachstum                  |
      | Politik: Opposition kritisiert            |
    When I make a GET request to "/api/rss/spiegel/top5?filter=Wirtschaft&limit=5"
    Then the response status should be 200
    And the response should contain a headlines array
    And the headlines array should have exactly 5 items
    And all headlines should contain "Wirtschaft" in their title

  @case-insensitive
  Scenario: Case-insensitive filtering
    Given the RSS feed has these articles:
      | title                               |
      | COVID-19: Neue Variante entdeckt    |
      | Sport: covid unterbricht Saison     |
      | Covid: Impfstoff-Update             |
      | CORONA: Maßnahmen verschärft        |
    When I make a GET request to "/api/rss/spiegel/top5?filter=covid"
    Then the response status should be 200
    And the response should contain a headlines array
    And the headlines array should have exactly 3 items

  @edge-cases
  Scenario: No matches for filter
    Given the RSS feed has these articles:
      | title                               |
      | Politik: Neue Gesetze               |
      | Sport: Bundesliga                   |
      | Wirtschaft: DAX                     |
    When I make a GET request to "/api/rss/spiegel/top5?filter=Blockchain"
    Then the response status should be 200
    And the response should contain a headlines array
    And the headlines array should be empty

  @edge-cases
  Scenario: Empty filter returns all articles
    Given the RSS feed has these articles:
      | title                               |
      | Politik: Nachrichten                |
      | Sport: Ergebnisse                   |
      | Wirtschaft: Börse                   |
      | Kultur: Events                      |
      | Wissenschaft: Forschung             |
      | Technik: Innovation                 |
    When I make a GET request to "/api/rss/spiegel/top5?filter="
    Then the response status should be 200
    And the response should contain a headlines array
    And the headlines array should have exactly 5 items

  @edge-cases
  Scenario: Filter with limit parameter interaction
    Given the RSS feed has these articles:
      | title                               |
      | Sport: Fußball-Bundesliga          |
      | Sport: Tennis-Turnier              |
      | Politik: Bundestag                 |
      | Sport: Handball-WM                 |
      | Sport: Formel 1                    |
    When I make a GET request to "/api/rss/spiegel/top5?filter=Sport&limit=2"
    Then the response status should be 200
    And the response should contain a headlines array
    And the headlines array should have exactly 2 items
    And all headlines should contain "Sport" in their title

  @latest-endpoint
  Scenario: Filter latest endpoint
    Given the RSS feed has these articles:
      | title                               |
      | Sport: Aktuell                      |
      | Politik: Breaking News              |
      | Wirtschaft: Börse geschlossen       |
    When I make a GET request to "/api/rss/spiegel/latest?filter=Politik"
    Then the response status should be 200
    And the headline title should contain "Politik"

  @latest-endpoint
  Scenario: Latest endpoint with no matching filter
    Given the RSS feed has these articles:
      | title                               |
      | Sport: Fußball                      |
      | Politik: Nachrichten                |
      | Wirtschaft: DAX                     |
    When I make a GET request to "/api/rss/spiegel/latest?filter=Technologie"
    Then the response status should be 404
    And the response should contain an error message