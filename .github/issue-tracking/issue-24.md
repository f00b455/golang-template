# Issue #24: Add support for displaying top 200 news items in terminal UI

**Issue URL**: https://github.com/f00b455/golang-template/issues/24
**Created**: 2025-09-26T08:46:15Z
**Assignee**: Unassigned

## Description
## User Story
As a user, I want to see more than 5 news items in the terminal interface, so that I can browse through a larger selection of current headlines.

## Acceptance Criteria
- [x] Increase the default news display limit from 5 to 200 items
- [x] Update the terminal UI to handle pagination for large result sets
- [x] Add keyboard navigation (Page Up/Page Down) for browsing long lists
- [x] Ensure performance remains good with 200 items loaded
- [x] Update API endpoint to support larger limits (currently max 5)
- [x] Add loading indicators for large data sets

## Technical Requirements
- Modify the `/api/rss/spiegel/top5` endpoint to support higher limits
- Update frontend JavaScript to handle 200 items efficiently
- Implement virtual scrolling or pagination to maintain performance
- Add proper loading states and error handling

## Priority
Medium - Enhancement to improve user experience with more content browsing options.

## Work Log
- Branch created: issue-24-add-support-for-displaying-top-200-news-items-in-t
- [x] Implementation - Completed backend and frontend features
- [x] Tests - Unit and BDD tests added for 200 item support
- [x] Documentation - Swagger API docs updated
