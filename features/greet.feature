# Issue: #1
# URL: https://github.com/f00b455/golang-template/issues/1
@pkg(shared) @issue-1
Feature: Shared Package Greet Function
  As a developer
  I want to use a greet function
  So that I can create personalized greeting messages

  Background:
    Given I am using the shared greet function

  @happy-path
  Scenario: Greet with valid name
    Given I have the name "World"
    When I call the greet function
    Then I should receive "Hello, World!"

  @happy-path
  Scenario: Greet with another name
    Given I have the name "Alice"
    When I call the greet function
    Then I should receive "Hello, Alice!"

  @error-handling
  Scenario: Greet with empty name
    Given I have the name ""
    When I call the greet function
    Then I should receive "Error: Name cannot be empty"

  @error-handling
  Scenario: Greet with whitespace only name
    Given I have the name "   "
    When I call the greet function
    Then I should receive "Error: Name cannot be empty"