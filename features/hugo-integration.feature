# Issue: #26
# URL: https://github.com/f00b455/golang-template/issues/26
@pkg(hugo) @issue-26
Feature: Hugo Static Site Integration
  As a developer
  I want minimal Hugo integration with basic functionality
  So that I can have a clean static site generator foundation

  Background:
    Given Hugo is installed and available
    And the Hugo site directory exists at "site/"

  @happy-path
  Scenario: Creating a basic Hugo site
    Given I have no existing Hugo site
    When I run the Hugo site creation command
    Then a new Hugo site should be created in "site/" directory
    And it should have basic directory structure
    And it should have no themes installed

  @content
  Scenario: Adding story content in markdown
    Given I have a Hugo site initialized
    When I add a story content file "content/stories/first-story.md"
    Then the markdown file should be created
    And it should contain valid frontmatter
    And it should contain story content in markdown format

  @api-integration
  Scenario: Displaying RSS data from API
    Given the Go API is running on port 3002
    And Hugo site has a template for RSS display
    When I fetch RSS data from the API endpoint "/api/rss"
    Then the data should be displayed in plain HTML
    And no CSS styling should be applied
    And the data should be readable without styles

  @search
  Scenario: Simple search functionality
    Given I have multiple story content pages
    And the site has a search feature
    When I search for a specific term
    Then matching content should be displayed
    And results should be in plain HTML format

  @build
  Scenario: Building the Hugo site
    Given I have a configured Hugo site with content
    When I run the Hugo build command
    Then the site should build successfully
    And static files should be generated in "public/" directory
    And the build should complete without errors

  @serve
  Scenario: Running Hugo development server
    Given I have a built Hugo site
    When I start the Hugo server on port 1313
    Then the server should start successfully
    And the site should be accessible at "http://localhost:1313"
    And both API (port 3002) and Hugo (port 1313) should run simultaneously