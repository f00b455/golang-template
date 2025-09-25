# Issue #7: Implement RSS headline filtering functionality

**Issue URL**: https://github.com/f00b455/golang-template/issues/7
**Created**: 2025-09-25T14:39:10Z
**Assignee**: Unassigned

## Description
## Problem
The RSS filtering feature described in issue #5 is not yet implemented. Users cannot filter RSS headlines to see only relevant content.

## Current State
- ✅ API server running on port 3002
- ✅ Web server running on port 8080
- ❌ No filter parameter in API endpoints
- ❌ No filter input in web UI
- ❌ No filtering logic in RSS handlers

## Expected Behavior
Based on issue #5, the filtering should:
1. Add a `filter` query parameter to RSS API endpoints
2. Filter headlines by title/content (case-insensitive)
3. Apply filtering BEFORE limiting results (important for correct count)
4. Add filter input field to web interface
5. Maintain the requested number of results after filtering

## API Example
```
GET /api/rss/spiegel/top5?filter=Taylor
# Should return up to 5 headlines containing "Taylor"
```

## Web UI Example
- Add text input field above headline list
- Filter headlines in real-time as user types
- Show filtered count vs total count

## Acceptance Criteria
- [ ] API endpoints accept `filter` query parameter
- [ ] Filtering is case-insensitive
- [ ] Results are filtered before applying count limit
- [ ] Web UI has filter input field
- [ ] Filtering works in real-time
- [ ] BDD tests cover filtering scenarios
- [ ] All existing tests still pass

## Priority
High - This is the main feature request for improved RSS browsing experience.

## Work Log
- Branch created: issue-7-implement-rss-headline-filtering-functionality
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
