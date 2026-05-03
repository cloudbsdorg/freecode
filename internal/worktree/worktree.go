package worktree

import (
	"context"
	"os/exec"
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
	return []*Worktree{}
}

func Remove(ctx context.Context, repoPath, name string) error {
	cmd := exec.CommandContext(ctx, "git", "worktree", "remove", name)
	cmd.Dir = repoPath
	return cmd.Run()
}
