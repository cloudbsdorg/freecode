---
name: ai-slop-remover
description: Detect and fix AI-generated code patterns that reduce quality, readability, or maintainability
---

# AI Slop Remover Skill

Identify and refactor code that exhibits AI-generated patterns - overly verbose, unnecessarily complex, or following trendy patterns without justification.

## When to Use

- Cleaning up AI-generated boilerplate
- Simplifying over-engineered solutions
- Removing redundant comments/explanations
- Fixing "hallucinated" abstractions
- Streamlining generated test files

## Detection Patterns

### Over-Commenting
```go
// Bad: Obvious comments that add noise
// Iterate through the list
for i := range items {
    // Process each item
    process(items[i])
}

// Good: No unnecessary comments
for i := range items {
    process(items[i])
}
```

### Unnecessary Abstractions
```go
// Bad: Factory for simple constructor
func NewUser(name string) *User { return &User{Name: name} }

// Good: Direct struct literal
user := &User{Name: name}
```

### Verbose Error Handling
```go
// Bad: Wrapping everything
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething failed: %w", err)
}

// Good: Only wrap when context needed
result, err := doSomething()
if err != nil {
    return err  // Caller adds context if needed
}
```

### Empty Error Handlers
```go
// Bad
if err != nil {}

// Good
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

### Redundant Type Annotations
```go
// Bad: Redundant type
var users map[string]*User = make(map[string]*User)

// Good: Infer type
users := make(map[string]*User)
```

## Refactoring Patterns

| Pattern | Problem | Fix |
|---------|---------|-----|
| `err != nil {}` | Silent failure | Log or return |
| Excessive interfaces | Over-abstraction | Remove if single impl |
| Gold-plating | Future-proofing | YAGNI |
| Commented code | Dead code | Delete |

## Code Smells Checklist

- [ ] Comments explaining obvious code
- [ ] Classes with single methods
- [ ] Interfaces for concrete types
- [ ] Unused parameters/imports
- [ ] Magic numbers without constants
- [ ] Deeply nested conditionals

## Integration

Use `task(category='quick', load_skills=['ai-slop-remover'])` for cleanup tasks. Target specific files or directories.

## Philosophy

> Debugging is twice as hard as writing code. If you write code as cleverly as possible, you are, by definition, not smart enough to debug it.

Keep code simple. Remove what doesn't add value. Let code breathe.
