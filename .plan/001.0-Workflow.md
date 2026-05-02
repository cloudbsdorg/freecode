# Freecode — Agent Workflow Instructions

## 1.0 Task Claiming Protocol

### 1.1 Before Starting Any Task

1. **Read the relevant plan documents** in full before touching any code
2. **Check the [Validation Report](./11.0-Validation.md)** for current task status
3. **Mark your task as `in_progress`** in the todo list
4. **Create a branch** for your work: `git checkout -b freecode/<task-name>`

### 1.2 Task States

| State | Description |
|-------|-------------|
| `pending` | Not yet started |
| `in_progress` | Currently being worked on |
| `completed` | Finished and tested |
| `cancelled` | No longer needed |

### 1.3 Completion Protocol

When completing a task:

1. **Ensure all tests pass** (`go test ./...`)
2. **Run linting** (`golangci-lint run`)
3. **Update the [Validation Report](./11.0-Validation.md)** to mark task complete
4. **Update any relevant plan documents** if architecture changed
5. **Create a pull request** with clear description of changes

---

## 2.0 Merge Conflict Handling

### 2.1 If Conflicts Occur

1. **Do NOT resolve conflicts blindly** — understand the intent of both changes
2. **Consult the original plan documents** to understand the intended design
3. **If unclear, ask** before making decisions
4. **Prefer the Go idiomatic solution** when TypeScript patterns conflict

### 2.2 Cross-Document Consistency

If you update one plan document, check others for consistency:
- 3.0 (Implementation Tasks) should match 2.0 (Design)
- 4.0 (Configuration) should match 1.1 (Architecture)
- 9.0 (Security) constraints must be reflected in implementation

---

## 3.0 Multi-Agent Coordination

### 3.1 Independent Tasks

The following tasks can be done in parallel:
- Core CLI foundation (Phase 1)
- Configuration system skeleton (Phase 2.1)
- Package structure setup (Phase 1.1)

### 3.2 Dependent Tasks

Wait for these before starting:
- Phase 2.2 (Config schema) requires Phase 2.1 complete
- Phase 3.1 (Tool implementation) requires core CLI
- Phase 4.1 (Server) requires tools

### 3.3 Weekly Sync

Every Friday:
- Update validation report
- Review open pull requests
- Update task statuses

---

## 4.0 Code Style Guidelines

### 4.1 Go Idioms

- Use `go fmt` before commits
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `context.Context` for cancellation
- Prefer `errors.Wrap` over bare `error` returns
- Use `log/slog` for structured logging

### 4.2 Error Handling

```go
// Good
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Bad
if err != nil {
    return err
}
```

### 4.3 Testing

```go
func TestFeature(t *testing.T) {
    // Use table-driven tests
    tests := []struct {
        name    string
        input   string
        want    string
    }{
        {"basic", "hello", "hello"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := feature(tt.input)
            if got != tt.want {
                t.Errorf("feature() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

## 5.0 Commit Message Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Example:
```
feat(config): add YAML config file support

Adds support for reading configuration from YAML files alongside
the existing JSON support. Maintains backward compatibility.

Closes #42
```

---

## 6.0 Pull Request Checklist

- [ ] Tests added/updated
- [ ] `go fmt` applied
- [ ] `golangci-lint` passes
- [ ] Documentation updated (if needed)
- [ ] Plan documents updated (if needed)
- [ ] Validation report updated

---

**Author:** Mark LaPointe <mark@cloudbsd.org>
**Last Updated:** 2026-05-01

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author
