package tool

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type EditTool struct{}

func NewEditTool() *EditTool {
	return &EditTool{}
}

func (t *EditTool) Name() string {
	return "edit"
}

func (t *EditTool) Description() string {
	return "Edit a file by replacing lines"
}

func (t *EditTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "edit",
		Description: "Edit a file by replacing lines",
		Parameters: map[string]Parameter{
			"file_path": {
				Type:        "string",
				Description: "Path to the file to edit",
				Required:    true,
			},
			"old_string": {
				Type:        "string",
				Description: "String to replace",
				Required:    true,
			},
			"new_string": {
				Type:        "string",
				Description: "Replacement string",
				Required:    true,
			},
		},
	}
}

func (t *EditTool) Execute(ctx context.Context, req Request) (*Response, error) {
	path, ok := req.Arguments["file_path"].(string)
	if !ok {
		return nil, fmt.Errorf("file_path must be a string")
	}

	oldStr, ok := req.Arguments["old_string"].(string)
	if !ok {
		return nil, fmt.Errorf("old_string must be a string")
	}

	newStr, ok := req.Arguments["new_string"].(string)
	if !ok {
		return nil, fmt.Errorf("new_string must be a string")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	content := strings.Join(lines, "\n")

	if !strings.Contains(content, oldStr) {
		return nil, fmt.Errorf("old_string not found in file")
	}

	newContent := strings.Replace(content, oldStr, newStr, 1)

	if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &Response{
		Result: fmt.Sprintf("Edited %s", path),
	}, nil
}

type GrepTool struct{}

func NewGrepTool() *GrepTool {
	return &GrepTool{}
}

func (t *GrepTool) Name() string {
	return "grep"
}

func (t *GrepTool) Description() string {
	return "Search for patterns in files"
}

func (t *GrepTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "grep",
		Description: "Search for patterns in files",
		Parameters: map[string]Parameter{
			"pattern": {
				Type:        "string",
				Description: "Regular expression pattern to search for",
				Required:    true,
			},
			"path": {
				Type:        "string",
				Description: "Path to search in",
				Required:    true,
			},
			"recursive": {
				Type:        "boolean",
				Description: "Search recursively",
				Default:     false,
			},
			"ignore_case": {
				Type:        "boolean",
				Description: "Case insensitive search",
				Default:     false,
			},
		},
	}
}

func (t *GrepTool) Execute(ctx context.Context, req Request) (*Response, error) {
	pattern, ok := req.Arguments["pattern"].(string)
	if !ok {
		return nil, fmt.Errorf("pattern must be a string")
	}

	path, ok := req.Arguments["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path must be a string")
	}

	recursive, _ := req.Arguments["recursive"].(bool)
	ignoreCase, _ := req.Arguments["ignore_case"].(bool)

	searchPattern := pattern
	if ignoreCase {
		searchPattern = "(?i)" + pattern
	}

	_, err := regexp.Compile(searchPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	return &Response{
		Result: fmt.Sprintf("Grep: pattern=%s path=%s recursive=%v", pattern, path, recursive),
	}, nil
}
