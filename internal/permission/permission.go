package permission

import (
	"context"
	"sync"
)

type Action string

const (
	ActionRead    Action = "read"
	ActionWrite   Action = "write"
	ActionDelete  Action = "delete"
	ActionExecute Action = "execute"
)

type ResourceKind string

const (
	KindFile  ResourceKind = "file"
	KindHTTP  ResourceKind = "http"
	KindShell ResourceKind = "shell"
	KindAny   ResourceKind = "*"
)

type Permission struct {
	Resource string
	Kind     ResourceKind
	Actions  []Action
}

// matchPattern implements glob-style pattern matching.
// "*" matches any chars except /, "**" matches any chars including /, "?" matches single char.
func matchPattern(resource, pattern string) bool {
	if pattern == "*" {
		return true
	}

	if !containsDoubleWildcard(pattern) {
		return matchSimple(resource, pattern)
	}

	return matchDoubleStar(resource, pattern)
}

func containsDoubleWildcard(pattern string) bool {
	return len(pattern) >= 2 && contains(pattern, "**")
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func matchSimple(resource, pattern string) bool {
	si := 0
	wi := 0
	wildcardStart := -1
	resourceStart := 0

	for si < len(resource) && wi < len(pattern) {
		switch pattern[wi] {
		case '*':
			if wi+1 < len(pattern) {
				wildcardStart = wi + 1
				resourceStart = si
				wi++
			} else {
				return true
			}
		case '?':
			wi++
			si++
		default:
			if resource[si] != pattern[wi] {
				if wildcardStart >= 0 {
					resourceStart++
					si = resourceStart
					wi = wildcardStart
				} else {
					return false
				}
			} else {
				wi++
				si++
			}
		}
	}

	if wi < len(pattern) {
		for wi < len(pattern) && pattern[wi] == '*' {
			wi++
		}
	}

	return wi == len(pattern) && si == len(resource)
}

func matchDoubleStar(resource, pattern string) bool {
	starIdx := indexOf(pattern, "**")
	if starIdx < 0 {
		return false
	}

	before := pattern[:starIdx]
	after := pattern[starIdx+2:]

	if indexOf(after, "**") >= 0 {
		nextStar := indexOf(after, "**")
		middle := after[:nextStar]
		rest := after[nextStar+2:]

		if !startsWith(resource, before) {
			return false
		}

		searchStart := len(before)
		for i := searchStart; i <= len(resource); i++ {
			if i == len(resource) || (i > 0 && resource[i-1] == '/') {
				if matchDoubleStar(resource[i:], middle+"/**/"+rest) {
					return true
				}
			}
		}
		return false
	}

	if !startsWith(resource, before) {
		return false
	}

	if !endsWith(resource, after) {
		return false
	}

	return true
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func startsWith(s, prefix string) bool {
	if len(prefix) > len(s) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if s[i] != prefix[i] {
			return false
		}
	}
	return true
}

func endsWith(s, suffix string) bool {
	if len(suffix) > len(s) {
		return false
	}
	offset := len(s) - len(suffix)
	for i := 0; i < len(suffix); i++ {
		if s[offset+i] != suffix[i] {
			return false
		}
	}
	return true
}

type compiledPattern struct {
	pattern  string
	kind     ResourceKind
	segments []string
	isWild   bool
}

func NewCompiledPattern(pattern string, kind ResourceKind) *compiledPattern {
	return &compiledPattern{
		pattern:  pattern,
		kind:     kind,
		segments: splitPatternSegments(pattern),
		isWild:   pattern == "*" || containsDoubleWildcard(pattern),
	}
}

func splitPatternSegments(pattern string) []string {
	var segments []string
	var current []byte
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == '/' {
			if len(current) > 0 {
				segments = append(segments, string(current))
				current = nil
			}
		} else {
			current = append(current, pattern[i])
		}
	}
	if len(current) > 0 {
		segments = append(segments, string(current))
	}
	return segments
}

func (cp *compiledPattern) Match(resource string) bool {
	if cp.isWild {
		return matchPattern(resource, cp.pattern)
	}
	return matchSimple(resource, cp.pattern)
}

func (cp *compiledPattern) MatchWithKind(resource string, kind ResourceKind) bool {
	if cp.kind != KindAny && cp.kind != kind {
		return false
	}
	return cp.Match(resource)
}

type PatternChecker struct {
	memoryChecker
	compiledPatterns map[string][]*compiledPattern
	mu               sync.RWMutex
}

func NewPatternChecker() *PatternChecker {
	return &PatternChecker{
		memoryChecker:    *NewMemoryChecker(),
		compiledPatterns: make(map[string][]*compiledPattern),
	}
}

func (c *PatternChecker) Grant(ctx context.Context, subject string, permission Permission) error {
	cp := NewCompiledPattern(permission.Resource, permission.Kind)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.memoryChecker.Grant(ctx, subject, permission)
	c.compiledPatterns[subject] = append(c.compiledPatterns[subject], cp)
	return nil
}

func (c *PatternChecker) Check(ctx context.Context, subject string, permission Permission) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	granted, err := c.memoryChecker.Check(ctx, subject, permission)
	if err != nil {
		return false, err
	}
	if granted {
		return true, nil
	}

	patterns, ok := c.compiledPatterns[subject]
	if !ok {
		return false, nil
	}

	for _, cp := range patterns {
		if cp.MatchWithKind(permission.Resource, permission.Kind) {
			for _, action := range permission.Actions {
				if c.hasAction(cp, action) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (c *PatternChecker) hasAction(cp *compiledPattern, action Action) bool {
	return true
}

func (c *PatternChecker) Revoke(ctx context.Context, subject string, permission Permission) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.memoryChecker.Revoke(ctx, subject, permission)

	patterns := c.compiledPatterns[subject]
	var filtered []*compiledPattern
	for _, cp := range patterns {
		if cp.pattern != permission.Resource {
			filtered = append(filtered, cp)
		}
	}
	c.compiledPatterns[subject] = filtered

	return nil
}

type memoryChecker struct {
	mu          sync.RWMutex
	permissions map[string][]Permission
}

func NewMemoryChecker() *memoryChecker {
	return &memoryChecker{permissions: make(map[string][]Permission)}
}

func (c *memoryChecker) Check(ctx context.Context, subject string, permission Permission) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	perms, ok := c.permissions[subject]
	if !ok {
		return false, nil
	}
	for _, p := range perms {
		if p.Resource == permission.Resource {
			for _, a := range p.Actions {
				if a == permission.Actions[0] {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (c *memoryChecker) Grant(ctx context.Context, subject string, permission Permission) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.permissions[subject] = append(c.permissions[subject], permission)
	return nil
}

func (c *memoryChecker) Revoke(ctx context.Context, subject string, permission Permission) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	perms := c.permissions[subject]
	for i, p := range perms {
		if p.Resource == permission.Resource {
			c.permissions[subject] = append(perms[:i], perms[i+1:]...)
			return nil
		}
	}
	return nil
}
