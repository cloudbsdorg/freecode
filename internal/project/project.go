package project

import (
	"context"
	"os"
	"path/filepath"
)

type Project struct {
	ID       string
	Name     string
	Path     string
	Remote   string
	VCS      string
	Created  int64
	Modified int64
}

type Detector interface {
	Detect(ctx context.Context, dir string) (*Project, error)
}

type projectDetector struct{}

func NewDetector() *projectDetector {
	return &projectDetector{}
}

func (d *projectDetector) Detect(ctx context.Context, dir string) (*Project, error) {
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		return &Project{
			ID:   dir,
			Path: dir,
			VCS:  "git",
		}, nil
	}
	return nil, nil
}

func List(ctx context.Context, dirs []string) ([]*Project, error) {
	var projects []*Project
	for _, dir := range dirs {
		detector := NewDetector()
		if proj, err := detector.Detect(ctx, dir); err == nil && proj != nil {
			projects = append(projects, proj)
		}
	}
	return projects, nil
}
