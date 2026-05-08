package worktree

import (
	"context"
	"os/exec"
	"strings"
)

type Worktree struct {
	Path   string
	Name   string
	Branch string
}

func Add(ctx context.Context, repoPath, name, branch string) (*Worktree, error) {
	cmd := exec.CommandContext(ctx, "git", "worktree", "add", "-b", name, name, branch)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return &Worktree{Path: name, Name: name, Branch: branch}, nil
}

func List(ctx context.Context, repoPath string) ([]*Worktree, error) {
	cmd := exec.CommandContext(ctx, "git", "worktree", "list")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseList(string(out)), nil
}

func parseList(output string) []*Worktree {
	var worktrees []*Worktree
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		wt := parseWorktreeLine(line)
		if wt != nil {
			worktrees = append(worktrees, wt)
		}
	}
	return worktrees
}

func parseWorktreeLine(line string) *Worktree {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}
	wt := &Worktree{
		Path: parts[0],
	}
	if len(parts) >= 2 {
		wt.Branch = parts[1]
	}
	if strings.Contains(line, "(detached)") {
		wt.Branch = "(detached)"
	}
	if strings.Contains(line, "(pruned)") {
		wt.Branch = "(pruned)"
	}
	if strings.Contains(line, "[") {
		if idx := strings.Index(wt.Branch, "["); idx > 0 {
			wt.Branch = wt.Branch[:idx]
		}
	}
	return wt
}

func Remove(ctx context.Context, repoPath, name string) error {
	cmd := exec.CommandContext(ctx, "git", "worktree", "remove", name)
	cmd.Dir = repoPath
	return cmd.Run()
}
