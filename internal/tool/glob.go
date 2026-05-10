package tool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type GlobTool struct{}

func init() {
	Register("glob", func() Tool { return &GlobTool{} })
}

func NewGlobTool() *GlobTool {
	return &GlobTool{}
}

func (t *GlobTool) Name() string {
	return "glob"
}

func (t *GlobTool) Description() string {
	return "Find files matching a glob pattern"
}

func (t *GlobTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "glob",
		Description: "Find files matching a glob pattern",
		Parameters: map[string]Parameter{
			"pattern": {
				Type:        "string",
				Description: "Glob pattern (e.g., **/*.go)",
				Required:    true,
			},
			"path": {
				Type:        "string",
				Description: "Base path to search from",
				Default:     ".",
			},
		},
	}
}

func (t *GlobTool) Execute(ctx context.Context, req Request) (*Response, error) {
	pattern, ok := req.Arguments["pattern"].(string)
	if !ok {
		return nil, fmt.Errorf("pattern must be a string")
	}

	basePath := "."
	if p, ok := req.Arguments["path"].(string); ok {
		basePath = p
	}

	matches, err := filepath.Glob(filepath.Join(basePath, pattern))
	if err != nil {
		return nil, fmt.Errorf("glob failed: %w", err)
	}

	for i, m := range matches {
		matches[i] = filepath.FromSlash(m)
	}

	return &Response{
		Result: strings.Join(matches, "\n"),
	}, nil
}

func globWalkDir(root, pattern string, results *[]string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		path := filepath.Join(root, entry.Name())
		matched, err := filepath.Match(pattern, entry.Name())
		if err != nil {
			continue
		}
		if matched {
			*results = append(*results, filepath.FromSlash(path))
		}
		if entry.IsDir() {
			globWalkDir(path, pattern, results)
		}
	}
	return nil
}
