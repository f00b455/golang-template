# Issue #13: [STORY] As a developer, I want a terminal-themed Hugo static site frontend with live RSS filtering so that I have a hacker-style news reader

**Issue URL**: https://github.com/f00b455/golang-template/issues/13
**Created**: 2025-09-25T16:31:26Z
**Assignee**: Unassigned

## Description
## User Story
**As a** developer/power user
**I want** a terminal-themed static site frontend using Hugo
**So that** I can have a retro, hacker-style news reader with powerful filtering capabilities

## Background/Context
Replace the current basic HTML frontend with a sophisticated Hugo-generated static site using the terminal theme. This creates a unique, developer-focused news reading experience that looks like a terminal/console interface while providing modern functionality.

## Visual Design Concept
- **Terminal aesthetic**: Green/amber text on black background
- **ASCII art headers**: Old-school BBS-style decorations
- **Typewriter effect**: Text appears character by character
- **Command-line style input**: Filter field looks like terminal prompt
- **Matrix-style animations**: Subtle background effects
- **Retro sound effects**: Optional keyboard clicks, modem sounds

## Core Features

### 1. Terminal-Style Filter Input
- Real-time filtering as you type (like Claude's chat input)
- Command-line appearance with blinking cursor
- Support for advanced filter syntax (+include, -exclude, "exact", /regex/)
- Command history with arrow keys
- Autocomplete suggestions

### 2. Terminal Commands
```
:help     - Show available commands
:export   - Export data (csv/json)  
:theme    - Switch color themes
:refresh  - Reload RSS feed
:clear    - Clear screen
:stats    - Show statistics
:vim      - Enable vim keybindings
```

### 3. Keyboard Navigation
- j/k - Navigate up/down (vim-style)
- / - Focus search
- Enter - Open article
- Escape - Clear filter
- Tab - Autocomplete

## Technical Architecture

### Hugo Setup
- Install Hugo and terminal theme
- Configure for dynamic content integration
- Set up build pipeline

### Frontend Components
1. **Command-line filter** - Terminal-style input with real-time filtering
2. **RSS display** - Styled like terminal output
3. **Loading animations** - ASCII spinners, progress bars
4. **CRT effects** - Scan lines, flicker, glow

### API Integration
- Connect to existing Go backend (/api/rss/spiegel/top5)
- Real-time filter updates
- Cached responses for offline viewing

## Acceptance Criteria
- [ ] Hugo site builds and deploys successfully
- [ ] Terminal theme properly integrated
- [ ] Filter input works with <50ms response time
- [ ] All keyboard shortcuts functional
- [ ] Terminal commands execute correctly
- [ ] Mobile responsive design
- [ ] Page loads in <1 second
- [ ] Works offline with cached data

## Definition of Done
- [ ] Hugo site fully functional
- [ ] Terminal theme customized
- [ ] Filter system complete
- [ ] Keyboard navigation working
- [ ] Documentation complete
- [ ] Cross-browser tested
- [ ] Performance optimized

## Story Points: 13 (Large)
## Timeline: 2-3 sprints

## Work Log
- Branch created: issue-13-story-as-a-developer-i-want-a-terminal-themed-hugo
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
