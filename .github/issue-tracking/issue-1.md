# Issue #1: Change CLI greeting from 'Hello' to 'Hi'

**Issue URL**: https://github.com/f00b455/golang-template/issues/1
**Created**: 2025-09-25T11:51:57Z
**Assignee**: Unassigned

## Description
## User Story

As a user of the CLI tool,
I want the greeting message to say 'Hi' instead of 'Hello',
So that the greeting feels more casual and friendly.

## Acceptance Criteria

- [ ] When running `./bin/cli-tool` without parameters, it should display "Hi, World!" instead of "Hello, World!"
- [ ] When running `./bin/cli-tool --name Alice`, it should display "Hi, Alice!" instead of "Hello, Alice!"
- [ ] The greeting function in `pkg/shared/greet.go` should be updated to use "Hi" prefix
- [ ] All unit tests should be updated to expect "Hi" instead of "Hello"
- [ ] All BDD tests should be updated to expect the new greeting format
- [ ] The error message for empty names should remain unchanged: "Error: Name cannot be empty"

## Technical Details

Files that need to be updated:
- `pkg/shared/greet.go` - Change the greeting prefix
- `pkg/shared/greet_test.go` - Update test expectations
- `internal/handlers/greet.go` - May need updates if it has hardcoded strings
- `internal/handlers/greet_test.go` - Update test expectations
- `features/greet.feature` - Update BDD scenarios
- `features/greet_test.go` - Update step definitions
- `features/cli.feature` - Update CLI BDD scenarios
- `features/cli_test.go` - Update CLI step definitions
- `features/api-greet.feature` - Update API BDD scenarios
- `features/api_test.go` - Update API step definitions

## Definition of Done

- [ ] Code changes implemented
- [ ] All unit tests passing
- [ ] All BDD tests passing
- [ ] `make validate` passes successfully
- [ ] Documentation updated if needed

## Work Log
- Branch created: issue-1-change-cli-greeting-from-hello-to-hi
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
