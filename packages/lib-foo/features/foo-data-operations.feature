# Issue: #7
# URL: https://github.com/f00b455/ts-pure-template/issues/7
@pkg(lib-foo) @issue-7
Feature: Foo Data Transformation and Filtering
  As a developer using the lib-foo library
  I want to transform and filter data using pure functions
  So that I can process collections without side effects

  Background:
    Given I have access to the foo data operation functions

  Scenario: Transform string array to uppercase
    Given I have an array '["a", "b", "c"]'
    When I transform it with uppercase function
    Then the result should be '["A", "B", "C"]'

  Scenario: Transform with custom function
    Given I have an array '["hello", "world"]'
    When I transform it with a function that adds "!"
    Then the result should be '["hello!", "world!"]'

  Scenario: Filter numbers for even values
    Given I have a number array '[1, 2, 3, 4, 5]'
    When I filter it for even numbers
    Then the result should be '[2, 4]'

  Scenario: Filter strings containing specific text
    Given I have a string array '["foo", "bar", "foobar"]'
    When I filter it for items containing "foo"
    Then the result should be '["foo", "foobar"]'

  Scenario: Function purity - original arrays unchanged
    Given I have an array '["original", "data"]'
    When I transform and filter the array
    Then the original array should remain unchanged
    And the functions should return new arrays

  Scenario: Function purity - empty array handling
    Given I have an empty array
    When I transform it with any function
    Then the result should be an empty array
    When I filter it with any predicate
    Then the result should be an empty array