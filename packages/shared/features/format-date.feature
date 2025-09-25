# Issue: #7
# URL: https://github.com/f00b455/ts-pure-template/issues/7
@pkg(shared) @issue-7
Feature: Shared Date Formatting Function
  As a developer using the shared utilities
  I want to use a pure date formatting function
  So that I can format dates consistently across applications

  Background:
    Given I have access to the shared formatDate function

  Scenario: Format date to YYYY-MM-DD
    Given I have a date "2024-03-15T10:30:00Z"
    When I call formatDate with that date
    Then the result should be "2024-03-15"

  Scenario: Format date with different timezone
    Given I have a date "2024-12-31T23:59:59Z"
    When I call formatDate with that date
    Then the result should be "2024-12-31"

  Scenario: Function purity - readonly input parameter
    Given I have a date object
    When I call formatDate with that date
    Then the original date object should remain unchanged
    And the function should not modify its input

  Scenario: Function purity - same input produces same output
    Given I have a date "2024-01-01T00:00:00Z"
    When I call formatDate with that date multiple times
    Then all results should be identical
    And the result should be "2024-01-01"