---
name: review-work
description: Comprehensive code review covering correctness, security, performance, and maintainability
---

# Code Review Skill

Conduct thorough code reviews that balance quality with pragmatism. Focus on issues that matter and provide actionable feedback.

## When to Use

- Reviewing pull requests
- Pre-commit code reviews
- Security-focused audits
- Performance optimization reviews
- Architecture and design discussions

## Review Framework

### Correctness
- Does the code do what it's supposed to?
- Are edge cases handled?
- Are there potential runtime errors?
- Does it handle null/undefined appropriately?

### Security
- SQL injection vulnerabilities?
- XSS attack vectors?
- Authentication/authorization checks?
- Sensitive data exposure?
- Input validation?

### Performance
- N+1 query patterns?
- Unnecessary re-renders?
- Memory leaks?
- Inefficient algorithms?

### Maintainability
- Clear variable/function names?
- DRY principles followed?
- Appropriate abstractions?
- Tests included?

## Review Output Format

```markdown
## Summary
Brief overview of changes

## Issues Found
### Critical
- Issue with fix suggestion

### Minor
- Nitpick or style preference

## Recommendations
- Optional improvements

## LGTM
Approved with minor notes
```

## Giving Feedback

- **Be specific** - Point to exact line numbers
- **Be constructive** - Suggest improvements, don't just criticize
- **Explain why** - Help author understand the concern
- **Prioritize** - Distinguish blocking issues from suggestions

## Automated Checks (CI)

Ensure these pass before review:
- [ ] Linting/formatting
- [ ] Type checking
- [ ] Unit tests
- [ ] Build succeeds

## Integration

Use `task(category='unspecified-high', load_skills=['review-work'])` for thorough code review tasks. Provide PR diff or file paths as input.

## Focus Areas by Language

| Language | Common Issues |
|----------|---------------|
| Go | Goroutine leaks, error handling |
| TypeScript | Type safety, async patterns |
| Python | Indentation, PEP8 compliance |
| SQL | Query performance, injection |
