# Issue: #7
# URL: https://github.com/f00b455/ts-pure-template/issues/7
@pkg(lib-foo) @issue-7
Feature: Foo Processing Functions
  As a developer using the lib-foo library
  I want to process data with configurable prefix and suffix
  So that I can transform data consistently with functional patterns

  Background:
    Given I have access to the foo processing functions

  Scenario: Process data with prefix only
    Given I have a config with prefix "FOO: "
    When I process the input "test"
    Then the result should be "FOO: test"

  Scenario: Process data with prefix and suffix
    Given I have a config with prefix "FOO: " and suffix " :BAR"
    When I process the input "test"
    Then the result should be "FOO: test :BAR"

  Scenario: Process empty input
    Given I have a config with prefix "PREFIX: "
    When I process the input ""
    Then the result should be "PREFIX: "

  Scenario: Function purity - same config produces same processor
    Given I have a config with prefix "TEST: "
    When I create processors with this config multiple times
    Then all processors should behave identically
    And processing "data" should always return "TEST: data"