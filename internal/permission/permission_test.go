package permission

import (
	"context"
	"testing"
)

func TestChecker(t *testing.T) {
	c := NewMemoryChecker()
	ctx := context.Background()

	perm := Permission{Resource: "file.txt", Actions: []Action{ActionRead}}
	granted, err := c.Check(ctx, "user1", perm)
	if err != nil {
		t.Errorf("Check error: %v", err)
	}
	if granted {
		t.Error("expected not granted initially")
	}

	if err := c.Grant(ctx, "user1", perm); err != nil {
		t.Errorf("Grant error: %v", err)
	}

	granted, _ = c.Check(ctx, "user1", perm)
	if !granted {
		t.Error("expected granted after Grant")
	}

	if err := c.Revoke(ctx, "user1", perm); err != nil {
		t.Errorf("Revoke error: %v", err)
	}

	granted, _ = c.Check(ctx, "user1", perm)
	if granted {
		t.Error("expected not granted after Revoke")
	}
}

func TestMatchPattern(t *testing.T) {
	cases := []struct {
		resource string
		pattern  string
		expected bool
	}{
		{"file.txt", "*", true},
		{"file.txt", "file.txt", true},
		{"file.txt", "*.txt", true},
		{"file.txt", "*.go", false},
		{"file.txt", "file?", false},
		{"file.txt", "file??", false},
		{"file.txt", "????", false},
		{"file.txt", "other", false},
		{"fileA", "file?", true},
		{"fileAB", "file??", true},
		{"fileABC", "file???", true},
		{"file.txt", "file*", true},
		{"file.txt", "file*.txt", true},
		{"file.txt", "*.txt", true},
		{"other.txt", "*.txt", true},
		{"other.tx", "*.txt", false},
		{"https://api.example.com/users", "https://api.example.com/**", true},
		{"https://api.example.com/users/123", "https://api.example.com/**", true},
		{"https://other.com/users", "https://api.example.com/**", false},
		{"https://api.example.com/users/file.go", "https://api.example.com/**", true},
	}

	for _, tc := range cases {
		result := matchPattern(tc.resource, tc.pattern)
		if result != tc.expected {
			t.Errorf("matchPattern(%q, %q) = %v, want %v", tc.resource, tc.pattern, result, tc.expected)
		}
	}
}

func TestCompiledPattern(t *testing.T) {
	cp := NewCompiledPattern("*.go", KindFile)

	if !cp.Match("file.go") {
		t.Error("expected file.go to match")
	}
	if !cp.Match("src/file.go") {
		t.Error("expected src/file.go to match")
	}
	if cp.Match("file.txt") {
		t.Error("expected file.txt not to match")
	}

	if !cp.MatchWithKind("file.go", KindFile) {
		t.Error("expected MatchWithKind to match with KindFile")
	}
	if cp.MatchWithKind("file.go", KindHTTP) {
		t.Error("expected MatchWithKind to not match with KindHTTP")
	}
}

func TestPatternChecker(t *testing.T) {
	c := NewPatternChecker()
	ctx := context.Background()

	perm := Permission{Resource: "*.go", Kind: KindFile, Actions: []Action{ActionRead}}
	if err := c.Grant(ctx, "user1", perm); err != nil {
		t.Errorf("Grant error: %v", err)
	}

	checkPerm := Permission{Resource: "file.go", Kind: KindFile, Actions: []Action{ActionRead}}
	granted, _ := c.Check(ctx, "user1", checkPerm)
	if !granted {
		t.Error("expected granted for matching pattern")
	}

	checkPerm = Permission{Resource: "file.txt", Kind: KindFile, Actions: []Action{ActionRead}}
	granted, _ = c.Check(ctx, "user1", checkPerm)
	if granted {
		t.Error("expected not granted for non-matching pattern")
	}
}

func TestPatternCheckerRevoke(t *testing.T) {
	c := NewPatternChecker()
	ctx := context.Background()

	perm := Permission{Resource: "*.go", Kind: KindFile, Actions: []Action{ActionRead}}
	if err := c.Grant(ctx, "user1", perm); err != nil {
		t.Errorf("Grant error: %v", err)
	}

	checkPerm := Permission{Resource: "file.go", Kind: KindFile, Actions: []Action{ActionRead}}
	granted, _ := c.Check(ctx, "user1", checkPerm)
	if !granted {
		t.Error("expected granted before revoke")
	}

	if err := c.Revoke(ctx, "user1", perm); err != nil {
		t.Errorf("Revoke error: %v", err)
	}

	granted, _ = c.Check(ctx, "user1", checkPerm)
	if granted {
		t.Error("expected not granted after revoke")
	}
}
