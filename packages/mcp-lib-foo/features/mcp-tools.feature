# Issue: #19
# URL: https://github.com/f00b455/ts-pure-template/issues/19
@pkg-mcp-lib-foo @issue-19
Feature: MCP Server Tools for lib-foo
  As a developer
  I want to use lib-foo functions through MCP tools
  So that I can leverage foo processing in Claude Code

  Background:
    Given the MCP server is initialized

  @happy-path
  Scenario: Process text with fooProcess tool
    Given I have input "test"
    And I have prefix "pre-"
    And I have suffix "-post"
    When I call the fooProcess tool
    Then the result should be "pre-test-post"

  Scenario: Create greeting with fooGreet tool
    Given I have a name "Alice"
    And I have prefix "***"
    And I have suffix "!!!"
    When I call the fooGreet tool
    Then the result should contain "Hello, Alice"
    And the result should start with "***"
    And the result should end with "!!!"

  Scenario: Use configured processor
    Given I have a configured processor with prefix "[" and suffix "]"
    And I select action "process"
    And I have input "data"
    When I call the createFooProcessor tool
    Then the result should be "[data]"

  Scenario: Transform array of strings
    Given I have an array of strings ["hello", "world"]
    And I select transform type "uppercase"
    When I call the fooTransform tool
    Then the result should be ["HELLO", "WORLD"]

  Scenario: Filter array with predicate
    Given I have an array [1, "test", 0, "", "valid", null]
    And I select filter type "truthy"
    When I call the fooFilter tool
    Then the result should contain exactly [1, "test", "valid"]