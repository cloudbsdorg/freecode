# Freecode — Unit Tests

**Document ID:** Freecode-UnitTests
**Version:** 1.0
**Last Updated:** 2026-05-02
**Maintainer:** Mark LaPointe <mark@cloudbsd.org>
**Status:** ACTIVE
**Classification:** INTERNAL

---

## Testing Scope

### Core Logic Identification

| Package | Module | Coverage Target |
|---------|--------|-----------------|
| internal/auth | store, credentials | 85%+ |
| internal/config | loading, migration | 85%+ |
| internal/cli | commands, flags | 80%+ |
| internal/session | manager, tabs | 85%+ |
| internal/provider | catalog, auth | 80%+ |
| internal/tool | registry, executor | 85%+ |

### Boundary Analysis

| Scenario | Test Type | Method |
|----------|-----------|--------|
| Empty provider list | Edge | Table-driven |
| Malformed config | Error | Unit |
| Concurrent session ops | Concurrency | Goroutine test |
| Path traversal attempts | Security | Fuzzing |

## Mocking Strategy

### When to Mock

- **External APIs**: Mock HTTP responses
- **File system**: Use `os.MkdirTemp`
- **Database**: Use in-memory SQLite
- **Time**: Use clock interfaces

### When NOT to Mock

- Business logic with no deps
- Error handling paths
- Concurrency primitives

## Validation Metrics

| Metric | Target | Current |
|--------|--------|---------|
| Line coverage | 85% | 68% |
| Branch coverage | 75% | 52% |
| Critical path | 100% | 95% |
| Regression tests | 100 new tests | 156 existing |

## Environment

```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/cli/...

# Run with race detector
go test -race ./...

# Run verbose
go test -v ./...
```

---

## Change Log

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-05-02 | Mark LaPointe | Initial unit test document |

**Last Updated:** 2026-05-02 07:30 UTC
**Classification:** INTERNAL
