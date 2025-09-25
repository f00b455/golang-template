# Issue: #7
# URL: https://github.com/f00b455/golang-template/issues/7
@pkg(api) @issue-7
Feature: Greeting API Endpoints
  As a user
  I want to get a greeting message via HTTP API
  So that I can verify the API endpoints are working

  Background:
    Given the API server is running

  @happy-path
  Scenario: Get default greeting via API
    When I make a GET request to "/api/greet"
    Then the response status should be 200
    And the response should contain JSON {"message": "Hello, World!"}

  @happy-path
  Scenario: Get personalized greeting via API
    When I make a GET request to "/api/greet?name=Go"
    Then the response status should be 200
    And the response should contain JSON {"message": "Hello, Go!"}

  @edge-cases
  Scenario: Get greeting with special characters
    When I make a GET request to "/api/greet?name=Go%20Developer"
    Then the response status should be 200
    And the response should contain JSON {"message": "Hello, Go Developer!"}

  @edge-cases
  Scenario: Get greeting with empty name parameter
    When I make a GET request to "/api/greet?name="
    Then the response status should be 200
    And the response should contain JSON {"message": "Error: Name cannot be empty"}