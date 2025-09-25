# Issue #27: Add ESLint rules to enforce function and file length limits

**Issue URL**: https://github.com/f00b455/ts-pure-template/issues/27
**Created**: 2025-09-25T07:24:24Z
**Assignee**: Unassigned

## Description
## Description
As discussed in PR #26, we should enforce code organization limits from our Clean Code principles using ESLint rather than relying on manual code review.

## Proposed Solution
Add ESLint rules to enforce both function and file length limits across the codebase.

### Suggested Configuration
```javascript
rules: {
  // Function length limit
  'max-lines-per-function': ['error', {
    max: 20,
    skipBlankLines: true,
    skipComments: true,
    IIFEs: false
  }],
  
  // File length limit for production code
  'max-lines': ['error', {
    max: 300,           // Maximum lines per file
    skipBlankLines: true,
    skipComments: true
  }]
}
```

### Override for Test Files
Test files and step definitions can be longer:
```javascript
overrides: [
  {
    files: ['**/*.test.ts', '**/*.spec.ts', '**/*.steps.ts'],
    rules: {
      'max-lines': ['error', { max: 500 }],
      'max-lines-per-function': ['error', { max: 50 }]  // Tests can be longer
    }
  }
]
```

## Considerations
⚠️ **Warning**: This will likely trigger a massive refactoring session as many existing files and functions probably exceed these limits.

## Suggested Implementation Approach
1. **Phase 1 - Assessment**
   - Run rules as `warn` to assess impact
   - Generate report of violations
   - Estimate refactoring effort

2. **Phase 2 - Gradual Refactoring**
   - Start with the most critical violations
   - Refactor package by package
   - Create separate PRs for each package

3. **Phase 3 - Enforcement**
   - Switch rules to `error`
   - Add to CI/CD pipeline
   - Document in CLAUDE.md

## Benefits
- **Better Code Organization**: Forces modular design
- **Single Responsibility**: Smaller files = focused purpose
- **Easier Testing**: Smaller units are easier to test
- **Better Readability**: Shorter files are easier to understand
- **Reduced Complexity**: Natural limit on complexity
- **Faster Code Reviews**: Smaller chunks to review

## Potential Refactoring Strategies
- Extract helper functions to separate modules
- Split large classes into smaller ones
- Use composition over inheritance
- Create utility modules for shared logic
- Split complex components into sub-components

## Related
- PR #26 where this was initially discussed
- CLAUDE.md Clean Code principles section
- Robert C. Martin's Clean Code book recommendations

---
_This issue was created during the implementation of test-reporter step definitions where we manually refactored functions to comply with the 20-line limit._

## Work Log
- Branch created: issue-27-add-eslint-rules-to-enforce-function-and-file-leng
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
