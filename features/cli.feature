# Issue: #14
# URL: https://github.com/f00b455/golang-template/issues/14
@pkg(cli) @issue-14
Feature: Hello-World CLI Application
  As a developer
  I want to use a colorful CLI to generate greeting messages
  So that I can have an engaging command-line experience

  Background:
    Given I have the hello-cli command available

  @happy-path
  Scenario: Generate default greeting without name parameter
    When I run hello-cli without parameters
    Then I should see a spinner message "Ready!"
    And I should see a progress message "Progress completed: 100%"
    And I should see a greeting message for "World"
    And the greeting should contain the prefix "✨"
    And the greeting should contain the suffix "✨"
    And the output should contain a decorative box

  @happy-path
  Scenario: Generate custom greeting with name parameter
    When I run hello-cli with name "Alice"
    Then I should see a spinner message "Ready!"
    And I should see a progress message "Progress completed: 100%"
    And I should see a greeting message for "Alice"
    And the greeting should contain the prefix "✨"
    And the greeting should contain the suffix "✨"

  @edge-case
  Scenario: Handle special characters in name
    When I run hello-cli with name "O'Brien"
    Then the command should complete successfully
    And I should see a greeting message for "O'Brien"

  @integration
  Scenario: Integration with core package
    When I run hello-cli with name "Test"
    Then the greeting should use the FooGreet function
    And the output should match the core package format

  @error-handling
  Scenario: Display help with --help flag
    When I run hello-cli with "--help" flag
    Then I should see a help message
    And the help should contain "A colorful Hello-World CLI application"