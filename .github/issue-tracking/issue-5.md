# Issue #5: [STORY] As a user, I want to filter RSS headlines so that I see only relevant content

**Issue URL**: https://github.com/f00b455/golang-template/issues/5
**Created**: 2025-09-25T13:07:54Z
**Assignee**: Unassigned

## Description
## User Story
**As a** web application user
**I want** to filter RSS headlines by keywords
**So that** I can see only the news articles that match my interests

## Background/Context
Currently, the RSS feed displays all headlines from Spiegel. Users need a way to filter these headlines based on keywords (similar to grep functionality) to focus on specific topics. The filter should be applied on the API side to ensure the correct number of filtered results is always returned.

## BDD Feature File (REQUIRED - Write FIRST!)
**⚠️ IMPORTANT: Feature files MUST be written BEFORE any implementation code!**

### Red-Green-Refactor Cycle
- [ ] **RED Phase**: Feature file written with failing scenarios
- [ ] **GREEN Phase**: Minimal implementation to pass tests
- [ ] **REFACTOR Phase**: Code quality improvements

### Example Feature File Content
```gherkin
# Issue: #8
# URL: https://github.com/f00b455/golang-template/issues/8
@pkg(api) @pkg(web) @issue-8
Feature: RSS Headline Filtering
  As a web application user
  I want to filter RSS headlines by keywords
  So that I can see only the news articles that match my interests

  Scenario: Filter headlines with matching keyword
    Given the RSS feed has 10 articles
    And 3 articles contain the word "Politik"
    When I request top 5 articles with filter "Politik"
    Then I should receive exactly 3 articles
    And all returned articles should contain "Politik" in their headline

  Scenario: Filter returns exact count when more matches exist
    Given the RSS feed has 20 articles
    And 8 articles contain the word "Wirtschaft"
    When I request top 5 articles with filter "Wirtschaft"
    Then I should receive exactly 5 articles
    And all returned articles should contain "Wirtschaft" in their headline

  Scenario: Case-insensitive filtering
    Given the RSS feed has articles with "COVID", "covid", and "Covid"
    When I request articles with filter "covid"
    Then I should receive all articles regardless of case

  Scenario: No matches for filter
    Given the RSS feed has 10 articles
    And no articles contain the word "Blockchain"
    When I request top 5 articles with filter "Blockchain"
    Then I should receive an empty result set

  Scenario: Empty filter returns all articles
    Given the RSS feed has 10 articles
    When I request top 5 articles with filter ""
    Then I should receive exactly 5 articles
```

## Acceptance Criteria
```gherkin
Given I am viewing the RSS feed page
When I enter a filter keyword "Politik"
Then only headlines containing "Politik" are displayed
```

### Detailed Acceptance Criteria
- [ ] Feature file created in `features/` directory
- [ ] All scenarios written before implementation
- [ ] Tests fail initially (RED phase)
- [ ] Minimal code to pass tests (GREEN phase)
- [ ] Code refactored for quality (REFACTOR phase)
- [ ] Filter input field added to web interface
- [ ] Filter parameter passed to API endpoint
- [ ] API applies filter before limiting results (e.g., filter then take top 5)
- [ ] Case-insensitive filtering implemented
- [ ] Empty filter shows all results
- [ ] Correct count always returned (e.g., exactly 5 when requested, or less if fewer matches)

## Technical Notes
- **API Changes**: Modify `/api/rss/spiegel/top5` endpoint to accept `filter` query parameter
- **Filtering Logic**: Apply filter BEFORE limiting to ensure correct count
- **Implementation approach**:
  1. Filter all RSS items by keyword (case-insensitive)
  2. Take the first N items from filtered results
  3. Return the limited set
- **Web Changes**: Add input field with real-time or button-triggered filtering
- **Performance**: Consider caching filtered results

## Design/UX Notes
- Filter input field above the RSS feed list
- Placeholder text: "Filter headlines... (e.g., Politik, Wirtschaft)"
- Clear filter button or 'x' in input field
- Show count of filtered results (e.g., "Showing 5 of 12 matching articles")
- Consider debouncing for real-time filtering

## Definition of Done
- [ ] Code complete and follows project standards
- [ ] Unit tests written and passing
- [ ] Integration tests written and passing
- [ ] BDD scenarios all passing
- [ ] Code reviewed and approved
- [ ] Documentation updated
- [ ] Feature tested in staging environment
- [ ] Acceptance criteria verified
- [ ] Performance benchmarks met

## Dependencies
- Current RSS feed implementation
- Existing API endpoint structure

## Story Points/Estimation
- Estimated effort: M (Medium)
- Story points: 5

## Additional Information
This feature enhances the RSS feed functionality by allowing users to focus on specific topics. The grep-like filtering should be intuitive and fast, with the API ensuring the correct number of results is always returned.

Example API calls:
- `/api/rss/spiegel/top5` - Returns top 5 articles (current behavior)
- `/api/rss/spiegel/top5?filter=Politik` - Returns top 5 articles containing "Politik"
- `/api/rss/spiegel/latest?filter=Wirtschaft` - Returns latest articles containing "Wirtschaft"

## Work Log
- Branch created: issue-5-story-as-a-user-i-want-to-filter-rss-headlines-so-
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
