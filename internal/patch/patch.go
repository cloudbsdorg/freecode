package patch

import (
	"context"
	"strings"
)

type Patch struct {
	Path    string
	OldLines []string
	NewLines []string
}

func Apply(ctx context.Context, patch Patch) error {
	return nil
}

func Create(ctx context.Context, old, new string) (*Patch, error) {
	oldLines := strings.Split(old, "\n")
	newLines := strings.Split(new, "\n")
	return &Patch{OldLines: oldLines, NewLines: newLines}, nil
}

func Parse(ctx context.Context, diff string) ([]*Patch, error) {
	return []*Patch{}, nil
}
