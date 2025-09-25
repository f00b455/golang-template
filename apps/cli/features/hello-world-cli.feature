# Issue: #14
# URL: https://github.com/f00b455/ts-pure-template/issues/14
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
    Then I should see a spinner for 5 seconds
    And I should see a progress bar that runs for 3 seconds
    And I should see a greeting message for "World"
    And the greeting should be displayed with colorful gradient
    And the greeting should be in a decorative box

  @happy-path
  Scenario: Generate custom greeting with name parameter
    When I run hello-cli with name "Alice"
    Then I should see a spinner for 5 seconds
    And I should see a progress bar that runs for 3 seconds
    And I should see a greeting message for "Alice"
    And the greeting should contain the prefix "✨"
    And the greeting should contain the suffix "✨"

  @edge-case
  Scenario: Handle special characters in name
    When I run hello-cli with name "O'Brien"
    Then the command should complete successfully
    And I should see a greeting message for "O'Brien"

  @integration
  Scenario: Integration with lib-foo
    When I run hello-cli with name "Test"
    Then the greeting should use the fooGreet function
    And the output should match the lib-foo format

  @error-handling
  Scenario: Handle invalid options gracefully
    When I run hello-cli with an invalid option
    Then I should see a help message
    And the process should exit with a non-zero code