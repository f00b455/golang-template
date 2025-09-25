# Issue: #7
# URL: https://github.com/f00b455/ts-pure-template/issues/7
@pkg(lib-foo) @issue-7
Feature: Foo Greeting Integration
  As a developer using the lib-foo library
  I want to create greeting functions with foo processing
  So that I can combine greeting functionality with custom formatting

  Background:
    Given I have access to the foo greeting functions

  Scenario: Create greeter with prefix
    Given I have a greeting config with prefix "[FOO] "
    When I greet "World"
    Then the result should be "[FOO] Hello, World!"

  Scenario: Create greeter with different prefix
    Given I have a greeting config with prefix "[TEST] "
    When I greet "Alice"
    Then the result should be "[TEST] Hello, Alice!"

  Scenario: Function purity - greeting with empty name handling
    Given I have a greeting config with prefix "[PREFIX] "
    When I greet ""
    Then the result should be "[PREFIX] Error: Name cannot be empty"

  Scenario: Function composition - greeting integrates with shared utilities
    Given I have a greeting config with prefix "[COMPOSED] "
    When I greet "Integration"
    Then the result should contain the shared greeting format
    And the result should start with "[COMPOSED] "