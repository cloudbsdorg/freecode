package git

import (
	"context"
	"os/exec"
	"strings"
)

type Repository struct {
	dir string
}

type Status struct {
	Staged    []string
	Modified  []string
	Untracked []string
	Deleted   []string
}

func Open(dir string) (*Repository, error) {
	return &Repository{dir: dir}, nil
}

func (r *Repository) Status(ctx context.Context) (*Status, error) {
	out, err := exec.CommandContext(ctx, "git", "status", "--porcelain").Output()
	if err != nil {
		return nil, err
	}

	status := &Status{}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		index := line[0]
		worktree := line[1]
		file := strings.TrimSpace(line[3:])

		switch index {
		case 'A', 'M':
			status.Staged = append(status.Staged, file)
		}
		switch worktree {
		case 'M':
			status.Modified = append(status.Modified, file)
		case '?':
			status.Untracked = append(status.Untracked, file)
		case 'D':
			status.Deleted = append(status.Deleted, file)
		}
	}
	return status, nil
}

func (r *Repository) Add(ctx context.Context, files ...string) error {
	args := append([]string{"add"}, files...)
	return exec.CommandContext(ctx, "git", args...).Run()
}

func (r *Repository) Commit(ctx context.Context, message string) error {
	cmd := exec.CommandContext(ctx, "git", "commit", "-m", message)
	cmd.Dir = r.dir
	return cmd.Run()
}

func (r *Repository) Log(ctx context.Context, count int) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "log", "--oneline", "-n", string(rune(count)))
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}

func (r *Repository) Diff(ctx context.Context, file string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "diff", file)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (r *Repository) Push(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "push")
	cmd.Dir = r.dir
	return cmd.Run()
}

func (r *Repository) Pull(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "pull")
	cmd.Dir = r.dir
	return cmd.Run()
}
