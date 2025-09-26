# Issue #22: CRITICAL: Claude Code Improvements workflow fails to detect PR - breaking automation pipeline

**Issue URL**: https://github.com/f00b455/golang-template/issues/22
**Created**: 2025-09-26T06:57:55Z
**Assignee**: Unassigned

## Description
## 🚨 Critical Bug

The Claude Code Improvements workflow is **completely broken** and not implementing requested changes, breaking our entire automated iteration pipeline.

## Problem
The workflow triggers via `workflow_run` events but fails to find the PR number:
```
No PR found for SHA: 60c1b42b2a15c7b0c2921d77bea790753e469cf6
expected an object but got: array ([0])
```

## Root Cause
The workflow uses `github.event.workflow_run.head_sha` to find the PR, but this SHA may not correspond to an open PR because:
1. It's the SHA from the **review workflow** run
2. The review workflow may use a different commit than the current PR HEAD
3. The GitHub API returns an empty array, causing the jq query to fail

## Impact
- ❌ **Automated improvements NOT running** 
- ❌ **Manual fixes required** for all "Request Changes" reviews
- ❌ **Pipeline completely broken**
- ❌ **No iteration happening**

## Immediate Fix Needed
Replace the PR detection logic with a more robust approach:

```yaml
- name: Get PR number from workflow name
  run: |
    WORKFLOW_NAME="${{ github.event.workflow_run.name }}"
    # Extract PR number from workflow name pattern: "Title (#123)"
    PR_NUMBER=$(echo "$WORKFLOW_NAME" | grep -oE '\(#[0-9]+\)' | grep -oE '[0-9]+')
    
    if [[ -z "$PR_NUMBER" ]]; then
      echo "No PR number found in workflow name: $WORKFLOW_NAME"
      exit 1
    fi
    
    echo "PR_NUMBER=$PR_NUMBER" >> $GITHUB_OUTPUT
    echo "Found PR #$PR_NUMBER from workflow name"
```

## Priority
**P0 - Critical** - This is blocking all automated improvements and breaking the core value proposition of the automated pipeline.

## Test Cases
- PR #21 currently has "Changes Requested" but improvements workflow is not running
- All future PRs will have the same issue

## Workflow Run Logs
See failing run: https://github.com/f00b455/golang-template/actions/runs/18030279061

## Work Log
- Branch created: issue-22-critical-claude-code-improvements-workflow-fails-t
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
