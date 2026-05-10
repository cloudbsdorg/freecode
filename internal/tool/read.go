package tool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ReadTool struct{}

func init() {
	Register("read", func() Tool { return &ReadTool{} })
}

func NewReadTool() *ReadTool {
	return &ReadTool{}
}

func (t *ReadTool) Name() string {
	return "read"
}

func (t *ReadTool) Description() string {
	return "Read file contents"
}

func (t *ReadTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "read",
		Description: "Read file contents",
		Parameters: map[string]Parameter{
			"file_path": {
				Type:        "string",
				Description: "Path to the file to read",
				Required:    true,
			},
			"offset": {
				Type:        "integer",
				Description: "Line offset to start reading from",
				Default:     0,
			},
			"limit": {
				Type:        "integer",
				Description: "Maximum number of lines to read",
			},
		},
	}
}

func (t *ReadTool) Execute(ctx context.Context, req Request) (*Response, error) {
	path, ok := req.Arguments["file_path"].(string)
	if !ok {
		return nil, fmt.Errorf("file_path must be a string")
	}

	offset := 0
	if o, ok := req.Arguments["offset"].(int); ok {
		offset = o
	}

	limit := -1
	if l, ok := req.Arguments["limit"].(int); ok {
		limit = l
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	if offset > 0 && offset < len(lines) {
		lines = lines[offset:]
	}

	if limit > 0 && limit < len(lines) {
		lines = lines[:limit]
	}

	return &Response{
		Result: filepath.FromSlash(path) + ":\n" + strings.Join(lines, "\n"),
	}, nil
}
