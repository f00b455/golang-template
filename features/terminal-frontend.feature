# Issue: #13
# URL: https://github.com/f00b455/golang-template/issues/13
@frontend @issue-13
Feature: Terminal-themed Hugo frontend with live RSS filtering
  As a developer
  I want a terminal-themed static site frontend
  So that I have a hacker-style news reader with powerful filtering

  Background:
    Given I have a running Hugo static site with terminal theme
    And the RSS API endpoint is available at "/api/rss/spiegel/top5"

  @happy-path
  Scenario: View terminal-themed interface
    Given I am on the terminal-themed frontend
    When the page loads
    Then I should see a terminal-style interface with green text on black background
    And I should see ASCII art headers
    And the page should load in less than 1 second

  @filtering
  Scenario: Filter RSS items in real-time
    Given I am on the terminal-themed frontend
    And RSS feed items are displayed
    When I type "golang" in the command-line filter
    Then I should see only RSS items containing "golang"
    And the filtering should happen in less than 50ms
    And the filter field should look like a terminal prompt with blinking cursor

  @keyboard
  Scenario: Navigate with keyboard shortcuts
    Given I am on the terminal-themed frontend
    And RSS feed items are displayed
    When I press "j" key
    Then the next item should be highlighted
    When I press "k" key
    Then the previous item should be highlighted
    When I press "/" key
    Then the filter input should be focused
    When I press "Escape" key
    Then the filter should be cleared

  @commands
  Scenario: Execute terminal commands
    Given I am on the terminal-themed frontend
    When I type ":help" in the command input
    Then I should see available commands displayed
    When I type ":refresh" in the command input
    Then the RSS feed should reload
    When I type ":clear" in the command input
    Then the screen should be cleared
    When I type ":theme" in the command input
    Then I should see theme options (green, amber, matrix)

  @advanced-filtering
  Scenario: Use advanced filter syntax
    Given I am on the terminal-themed frontend
    And RSS feed items are displayed
    When I type "+tech -politics" in the filter
    Then I should see items containing "tech" but not "politics"
    When I type '"exact phrase"' in the filter
    Then I should see only items with the exact phrase
    When I type "/regex.*pattern/" in the filter
    Then I should see items matching the regex pattern

  @offline
  Scenario: Work offline with cached data
    Given I am on the terminal-themed frontend
    And RSS feed data has been loaded once
    When I go offline
    And I refresh the page
    Then I should still see the cached RSS feed items
    And I should see an offline indicator

  @mobile
  Scenario: Responsive mobile design
    Given I am on a mobile device
    When I visit the terminal-themed frontend
    Then the interface should be responsive
    And touch gestures should work for navigation
    And the terminal theme should adapt to smaller screens

  @performance
  Scenario: Meet performance requirements
    Given I am on the terminal-themed frontend
    When I measure page load time
    Then initial load should be under 1 second
    When I type in the filter field
    Then filtering should respond in under 50ms
    And there should be no visible lag in the typewriter effect