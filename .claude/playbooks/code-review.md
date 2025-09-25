# Go Code Review Playbook

## When reviewing Go code in pull requests:

### 1. **Check BDD Compliance**:
   - ✅ Feature file exists for the issue
   - ✅ Feature file has issue reference: `# Issue: #<number>`
   - ✅ Tests were written BEFORE implementation
   - ✅ All BDD scenarios pass

### 2. **Go Code Quality**:
   - ✅ Functions under 60 lines (enforced by golangci-lint)
   - ✅ Cyclomatic complexity < 10
   - ✅ All errors handled explicitly
   - ✅ No naked returns in long functions
   - ✅ Interfaces used for dependency injection

### 3. **Testing**:
   - ✅ Table-driven tests for multiple cases
   - ✅ Test coverage adequate (aim for >80%)
   - ✅ Mocks used for external dependencies
   - ✅ BDD tests cover user scenarios

### 4. **Performance & Security**:
   - ✅ No SQL injection vulnerabilities
   - ✅ Proper input validation
   - ✅ No sensitive data in logs
   - ✅ Efficient database queries (no N+1)
   - ✅ Proper goroutine management

### 5. **Documentation**:
   - ✅ Swagger docs updated (if API changes)
   - ✅ README updated (if setup changes)
   - ✅ CLAUDE.md followed

## Review Response Format:
```markdown
## Code Review Summary

### ✅ Strengths
- [List positive aspects]

### ⚠️ Suggestions
- [Non-critical improvements]

### ❌ Issues to Address
- [Must-fix problems]

### 📊 Coverage Report
- Tests: X/Y passing
- Coverage: XX%
```