package tool

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type HashlineTool struct{}

func init() {
	Register("hashline", func() Tool { return &HashlineTool{} })
}

func NewHashlineTool() *HashlineTool {
	return &HashlineTool{}
}

func (t *HashlineTool) Name() string {
	return "hashline"
}

func (t *HashlineTool) Description() string {
	return "Edit a file using hash/line references for precise changes"
}

func (t *HashlineTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "hashline",
		Description: "Edit a file by specifying line numbers or patterns with optional hash verification",
		Parameters: map[string]Parameter{
			"file_path": {
				Type:        "string",
				Description: "Path to the file to edit",
				Required:    true,
			},
			"operation": {
				Type:        "string",
				Description: "Operation to perform: replace, insert, delete, extract",
				Required:    true,
			},
			"start_line": {
				Type:        "number",
				Description: "Starting line number (1-indexed)",
				Required:    false,
			},
			"end_line": {
				Type:        "number",
				Description: "Ending line number (1-indexed, inclusive)",
				Required:    false,
			},
			"content": {
				Type:        "string",
				Description: "New content for replace/insert operations",
				Required:    false,
			},
			"pattern": {
				Type:        "string",
				Description: "Pattern to match for line-based operations",
				Required:    false,
			},
			"hash": {
				Type:        "string",
				Description: "Expected hash of the lines being edited (for verification)",
				Required:    false,
			},
		},
	}
}

func (t *HashlineTool) Execute(ctx context.Context, req Request) (*Response, error) {
	path, ok := req.Arguments["file_path"].(string)
	if !ok {
		return nil, fmt.Errorf("file_path must be a string")
	}

	operation, ok := req.Arguments["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation must be a string")
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

	switch operation {
	case "extract":
		return t.extract(lines, req.Arguments)
	case "replace":
		return t.replace(path, lines, req.Arguments)
	case "insert":
		return t.insert(path, lines, req.Arguments)
	case "delete":
		return t.delete(path, lines, req.Arguments)
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}
}

func (t *HashlineTool) extract(lines []string, args map[string]interface{}) (*Response, error) {
	startLine, _ := args["start_line"].(float64)
	endLine, _ := args["end_line"].(float64)
	pattern, _ := args["pattern"].(string)

	if pattern != "" {
		return t.extractByPattern(lines, pattern)
	}

	return t.extractByLineRange(lines, int(startLine), int(endLine))
}

func (t *HashlineTool) extractByLineRange(lines []string, start, end int) (*Response, error) {
	if start < 1 || start > len(lines) {
		return nil, fmt.Errorf("start_line %d out of range (1-%d)", start, len(lines))
	}

	if end == 0 || end > len(lines) {
		end = len(lines)
	}

	if start > end {
		return nil, fmt.Errorf("start_line %d is greater than end_line %d", start, end)
	}

	extracted := lines[start-1 : end]
	return &Response{
		Result: fmt.Sprintf("Extracted lines %d-%d:\n%s", start, end, strings.Join(extracted, "\n")),
	}, nil
}

func (t *HashlineTool) extractByPattern(lines []string, pattern string) (*Response, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	var matchingLines []string
	var lineNums []int
	for i, line := range lines {
		if re.MatchString(line) {
			matchingLines = append(matchingLines, line)
			lineNums = append(lineNums, i+1)
		}
	}

	if len(matchingLines) == 0 {
		return nil, fmt.Errorf("no lines match pattern: %s", pattern)
	}

	return &Response{
		Result: fmt.Sprintf("Found %d matches at lines %v:\n%s", len(matchingLines), lineNums, strings.Join(matchingLines, "\n")),
	}, nil
}

func (t *HashlineTool) replace(path string, lines []string, args map[string]interface{}) (*Response, error) {
	startLine, _ := args["start_line"].(float64)
	endLine, _ := args["end_line"].(float64)
	content, _ := args["content"].(string)
	hash, _ := args["hash"].(string)

	start := int(startLine)
	end := int(endLine)

	if start < 1 || start > len(lines) {
		return nil, fmt.Errorf("start_line %d out of range (1-%d)", start, len(lines))
	}

	if end == 0 || end > len(lines) {
		end = len(lines)
	}

	if start > end {
		return nil, fmt.Errorf("start_line %d is greater than end_line %d", start, end)
	}

	if hash != "" {
		actualHash := hashLines(lines[start-1 : end])
		if actualHash != hash {
			return nil, fmt.Errorf("hash mismatch: expected %s, got %s", hash, actualHash)
		}
	}

	newLines := make([]string, 0, len(lines))
	newLines = append(newLines, lines[:start-1]...)
	newLines = append(newLines, strings.Split(content, "\n")...)
	newLines = append(newLines, lines[end:]...)

	if err := os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &Response{
		Result: fmt.Sprintf("Replaced lines %d-%d in %s", start, end, path),
	}, nil
}

func (t *HashlineTool) insert(path string, lines []string, args map[string]interface{}) (*Response, error) {
	lineNum, _ := args["start_line"].(float64)
	content, _ := args["content"].(string)

	pos := int(lineNum)
	if pos < 1 || pos > len(lines)+1 {
		return nil, fmt.Errorf("start_line %d out of range (1-%d)", pos, len(lines)+1)
	}

	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:pos-1]...)
	newLines = append(newLines, strings.Split(content, "\n")...)
	newLines = append(newLines, lines[pos-1:]...)

	if err := os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &Response{
		Result: fmt.Sprintf("Inserted content at line %d in %s", pos, path),
	}, nil
}

func (t *HashlineTool) delete(path string, lines []string, args map[string]interface{}) (*Response, error) {
	startLine, _ := args["start_line"].(float64)
	endLine, _ := args["end_line"].(float64)
	hash, _ := args["hash"].(string)

	start := int(startLine)
	end := int(endLine)

	if start < 1 || start > len(lines) {
		return nil, fmt.Errorf("start_line %d out of range (1-%d)", start, len(lines))
	}

	if end == 0 || end > len(lines) {
		end = len(lines)
	}

	if start > end {
		return nil, fmt.Errorf("start_line %d is greater than end_line %d", start, end)
	}

	if hash != "" {
		actualHash := hashLines(lines[start-1 : end])
		if actualHash != hash {
			return nil, fmt.Errorf("hash mismatch: expected %s, got %s", hash, actualHash)
		}
	}

	newLines := make([]string, 0, len(lines))
	newLines = append(newLines, lines[:start-1]...)
	newLines = append(newLines, lines[end:]...)

	if err := os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &Response{
		Result: fmt.Sprintf("Deleted lines %d-%d from %s", start, end, path),
	}, nil
}

func hashLines(lines []string) string {
	content := strings.Join(lines, "\n")
	h := 0
	for _, c := range content {
		h = h*31 + int(c)
	}
	return strconv.FormatInt(int64(h), 16)
}
