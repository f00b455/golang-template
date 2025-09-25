# Issue #15: 🐛 RSS filtering only searches first 5 items instead of entire feed

**Issue URL**: https://github.com/f00b455/golang-template/issues/15
**Created**: 2025-09-25T16:54:19Z
**Assignee**: Unassigned

## Description
## Bug Description
When filtering RSS headlines, the API only fetches 5 items from the RSS feed and then filters those 5 items. This means if the search term doesn't appear in the first 5 items, no results are returned even though matching items may exist further in the feed.

## Current Behavior
1. User enters filter term (e.g., 'Trump')
2. API fetches only 5 RSS items
3. Filter is applied to those 5 items
4. If no matches in first 5, user sees no results

## Expected Behavior
1. User enters filter term
2. API fetches a large number of RSS items (50-100)
3. Filter is applied to all fetched items
4. Top 5 filtered results are returned

## Root Cause
In `internal/handlers/rss.go` line 160:
```go
headlines, err := h.fetchMultipleHeadlines(5) // BUG: Should fetch many more
```

## Reproduction Steps
1. Start the API and web servers
2. Open the web UI
3. Enter a filter term that doesn't appear in the first 5 RSS items
4. Observe no results even though the term exists in later items

## Fix Required
- Modify `GetTop5` handler to fetch 50-100 items before filtering
- Ensure filtering happens on the full dataset
- Return top 5 results after filtering

## Work Log
- Branch created: issue-15-rss-filtering-only-searches-first-5-items-instead-
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
