# Issue #20: Add RSS export download links to terminal UI

**Issue URL**: https://github.com/f00b455/golang-template/issues/20
**Created**: 2025-09-26T06:37:36Z
**Assignee**: Unassigned

## Description
## Problem
The RSS export API endpoints are working (JSON and CSV formats), but there's no UI integration in the terminal frontend to access these export features.

## Current State
- ✅ Backend API endpoints working: `/api/rss/spiegel/export?format=json` and `/api/rss/spiegel/export?format=csv`
- ❌ No download buttons or links in the terminal UI
- ❌ Users cannot easily export data without knowing the API endpoints

## Required Changes
Add download/export functionality to the terminal UI at http://localhost:8080:

### 1. Export Buttons
- Add "Download as JSON" button
- Add "Download as CSV" button
- Place buttons near the RSS headlines display

### 2. Export Options
- Allow users to specify export limit (number of items)
- Allow users to apply current filter to export
- Show export progress/confirmation

### 3. Download Behavior
- Trigger browser download with proper filename
- Include current filter in filename if applied
- Handle errors gracefully

### 4. UI Integration
- Style buttons to match terminal theme
- Position logically near RSS content
- Add tooltips/help text for export options

## Expected Result
Users can export RSS headlines directly from the web UI without needing to know API endpoints.

## Priority
Medium - Feature is complete on backend but missing user-facing integration

## Work Log
- Branch created: issue-20-add-rss-export-download-links-to-terminal-ui
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
