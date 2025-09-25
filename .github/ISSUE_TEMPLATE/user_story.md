---
name: User Story
about: Describe a feature from the user's perspective
title: '[STORY] As a [user type], I want [goal] so that [benefit]'
labels: user-story, triage
assignees: ''

---

## User Story
**As a** [type of user]
**I want** [goal/desire]
**So that** [benefit/value]

## Background/Context
Provide any relevant background information or context for this story.

## BDD Feature File (REQUIRED - Write FIRST!)
**⚠️ IMPORTANT: Feature files MUST be written BEFORE any implementation code!**

### Red-Green-Refactor Cycle
- [ ] **RED Phase**: Feature file written with failing scenarios
- [ ] **GREEN Phase**: Minimal implementation to pass tests
- [ ] **REFACTOR Phase**: Code quality improvements

### Example Feature File Content
```gherkin
# Issue: #[this issue number]
# URL: https://github.com/f00b455/golang-template/issues/[this issue number]
@pkg([package]) @issue-[this issue number]
Feature: [Feature name from title]
  As a [user type]
  I want [goal]
  So that [benefit]

  Scenario: [Primary happy path]
    Given [precondition]
    When [action]
    Then [expected result]

  Scenario: [Error handling]
    Given [precondition]
    When [invalid action]
    Then [error result]
```

## Acceptance Criteria
```gherkin
Given [precondition]
When [action]
Then [expected result]
```

### Detailed Acceptance Criteria
- [ ] Feature file created in `features/` directory
- [ ] All scenarios written before implementation
- [ ] Tests fail initially (RED phase)
- [ ] Minimal code to pass tests (GREEN phase)
- [ ] Code refactored for quality (REFACTOR phase)
- [ ] All acceptance criteria below met

## Technical Notes
Any technical considerations or constraints:
- Database changes required
- API endpoints affected
- Performance considerations
- Security implications

## Design/UX Notes
- UI/UX requirements
- Mockups or wireframes (attach images if available)
- User flow description

## Definition of Done
- [ ] Code complete and follows project standards
- [ ] Unit tests written and passing
- [ ] Integration tests written and passing
- [ ] Code reviewed and approved
- [ ] Documentation updated
- [ ] Feature tested in staging environment
- [ ] Acceptance criteria verified
- [ ] Performance benchmarks met

## Dependencies
- Depends on: #issue_number
- Blocks: #issue_number

## Story Points/Estimation
- Estimated effort: [XS/S/M/L/XL]
- Story points: [1/2/3/5/8/13]

## Additional Information
Any other relevant information, links, or references.