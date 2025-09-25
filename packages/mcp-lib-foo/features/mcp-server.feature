# Issue: #19
# URL: https://github.com/f00b455/ts-pure-template/issues/19
@pkg-mcp-lib-foo @issue-19
Feature: MCP Server Infrastructure
  As a developer
  I want a properly functioning MCP server
  So that I can connect to it from Claude Code

  Background:
    Given the MCP server package is built

  Scenario: Server starts successfully
    When I start the MCP server
    Then the server should be running
    And it should respond to list tools request
    And it should expose 5 tools

  Scenario: List available tools
    Given the MCP server is running
    When I request the list of tools
    Then I should see "fooProcess" tool
    And I should see "fooGreet" tool
    And I should see "createFooProcessor" tool
    And I should see "fooTransform" tool
    And I should see "fooFilter" tool

  Scenario: Handle invalid tool call
    Given the MCP server is running
    When I call a non-existent tool "invalidTool"
    Then I should receive an error message
    And the error should contain "Unknown tool"

  Scenario: Handle malformed input
    Given the MCP server is running
    When I call fooProcess with missing required parameters
    Then I should receive an error response
    And the server should remain running