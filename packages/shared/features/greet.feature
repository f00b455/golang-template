# Issue: #7
# URL: https://github.com/f00b455/ts-pure-template/issues/7
@pkg(shared) @issue-7
Feature: Shared Greeting Function
  As a developer using the shared utilities
  I want to use a pure greeting function
  So that I can generate consistent greeting messages across applications

  Background:
    Given I have access to the shared greet function

  Scenario: Generate greeting with valid name
    When I call greet with "TypeScript"
    Then the result should be "Hello, TypeScript!"

  Scenario: Generate greeting with empty name
    When I call greet with ""
    Then the result should be "Error: Name cannot be empty"

  Scenario: Generate greeting with whitespace-only name
    When I call greet with "   "
    Then the result should be "Error: Name cannot be empty"

  Scenario: Function purity - same input produces same output
    When I call greet with "Alice" multiple times
    Then all results should be identical
    And the result should be "Hello, Alice!"