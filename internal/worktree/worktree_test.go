package worktree

import (
	"testing"
)

func TestWorktree(t *testing.T) {
	wt := &Worktree{
		Path:   "/path/to/worktree",
		Name:   "feature-branch",
		Branch: "main",
	}

	if wt.Path != "/path/to/worktree" {
		t.Errorf("Path = %q, want %q", wt.Path, "/path/to/worktree")
	}
	if wt.Name != "feature-branch" {
		t.Errorf("Name = %q, want %q", wt.Name, "feature-branch")
	}
	if wt.Branch != "main" {
		t.Errorf("Branch = %q, want %q", wt.Branch, "main")
	}
}

func TestParseList(t *testing.T) {
	output := `/path/to/main e8f3c2a [main]
/path/to/feature /path/to/feature (detached)
		`

	worktrees := parseList(output)
	if len(worktrees) != 2 {
		t.Errorf("parseList() returned %d worktrees, want 2", len(worktrees))
	}
}

func TestParseListEmpty(t *testing.T) {
	output := ""

	worktrees := parseList(output)
	if len(worktrees) != 0 {
		t.Errorf("parseList() on empty returned %d worktrees, want 0", len(worktrees))
	}
}

func TestParseWorktreeLine(t *testing.T) {
	line := "/path/to/worktree feature-branch"

	wt := parseWorktreeLine(line)
	if wt == nil {
		t.Fatal("parseWorktreeLine() returned nil")
	}
	if wt.Path != "/path/to/worktree" {
		t.Errorf("Path = %q, want %q", wt.Path, "/path/to/worktree")
	}
	if wt.Branch != "feature-branch" {
		t.Errorf("Branch = %q, want %q", wt.Branch, "feature-branch")
	}
}

func TestParseWorktreeLineDetached(t *testing.T) {
	line := "/path/to/worktree (detached)"

	wt := parseWorktreeLine(line)
	if wt == nil {
		t.Fatal("parseWorktreeLine() returned nil")
	}
	if wt.Path != "/path/to/worktree" {
		t.Errorf("Path = %q, want %q", wt.Path, "/path/to/worktree")
	}
}

func TestParseWorktreeLineEmpty(t *testing.T) {
	line := ""

	wt := parseWorktreeLine(line)
	if wt != nil {
		t.Error("parseWorktreeLine() on empty line should return nil")
	}
}

func TestParseWorktreeLineTooShort(t *testing.T) {
	line := "/only/one/part"

	wt := parseWorktreeLine(line)
	if wt != nil {
		t.Error("parseWorktreeLine() on short line should return nil")
	}
}