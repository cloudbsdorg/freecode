package tool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type WriteTool struct{}

func init() {
	Register("write", func() Tool { return &WriteTool{} })
}

func NewWriteTool() *WriteTool {
	return &WriteTool{}
}

func (t *WriteTool) Name() string {
	return "write"
}

func (t *WriteTool) Description() string {
	return "Write content to a file"
}

func (t *WriteTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "write",
		Description: "Write content to a file",
		Parameters: map[string]Parameter{
			"file_path": {
				Type:        "string",
				Description: "Path to the file to write",
				Required:    true,
			},
			"content": {
				Type:        "string",
				Description: "Content to write to the file",
				Required:    true,
			},
			"append": {
				Type:        "boolean",
				Description: "Append to file instead of overwriting",
				Default:     false,
			},
		},
	}
}

func (t *WriteTool) Execute(ctx context.Context, req Request) (*Response, error) {
	path, ok := req.Arguments["file_path"].(string)
	if !ok {
		return nil, fmt.Errorf("file_path must be a string")
	}

	content, ok := req.Arguments["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content must be a string")
	}

	appendMode := false
	if a, ok := req.Arguments["append"].(bool); ok {
		appendMode = a
	}

	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if appendMode {
		flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	}

	f, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	action := "Wrote"
	if appendMode {
		action = "Appended"
	}

	return &Response{
		Result: fmt.Sprintf("%s %d bytes to %s", action, len(content), filepath.FromSlash(path)),
	}, nil
}
