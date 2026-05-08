package patch

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidDiff    = errors.New("invalid diff format")
	ErrMismatch      = errors.New("patch does not match file content")
	ErrFileNotFound  = errors.New("file not found")
)

// Patch represents a unified diff patch
type Patch struct {
	Path     string
	OldLines []string
	NewLines []string
}

// Apply applies the patch to the file at path
func Apply(ctx context.Context, patch Patch) error {
	if patch.Path == "" {
		return ErrInvalidDiff
	}

	// Read the original file
	content, err := os.ReadFile(patch.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotFound
		}
		return err
	}

	// Split into lines
	lines := strings.Split(string(content), "\n")

	// Find and apply the patch
	result, ok := applyPatchToLines(lines, patch.OldLines, patch.NewLines)
	if !ok {
		return ErrMismatch
	}

	// Write the result
	return os.WriteFile(patch.Path, []byte(strings.Join(result, "\n")), 0644)
}

// applyPatchToLines applies the patch and returns the new lines
func applyPatchToLines(original, old, new []string) ([]string, bool) {
	if len(old) == 0 {
		// Edge case: adding new lines at the beginning
		return append(new, original...), true
	}

	// Find the old lines in original
	startIdx := findSequence(original, old)
	if startIdx < 0 {
		return nil, false
	}

	// Build result: lines before + new + lines after
	endIdx := startIdx + len(old)
	result := make([]string, 0, len(original)-len(old)+len(new))
	result = append(result, original[:startIdx]...)
	result = append(result, new...)
	result = append(result, original[endIdx:]...)

	return result, true
}

// findSequence finds the start index of a sequence in lines, returns -1 if not found
func findSequence(lines, seq []string) int {
	if len(seq) == 0 {
		return 0
	}
	if len(seq) > len(lines) {
		return -1
	}

	for i := 0; i <= len(lines)-len(seq); i++ {
		found := true
		for j := 0; j < len(seq); j++ {
			if lines[i+j] != seq[j] {
				found = false
				break
			}
		}
		if found {
			return i
		}
	}
	return -1
}

// Create creates a Patch from old and new content strings
func Create(ctx context.Context, old, new string) (*Patch, error) {
	oldLines := strings.Split(old, "\n")
	newLines := strings.Split(new, "\n")
	return &Patch{OldLines: oldLines, NewLines: newLines}, nil
}

// Parse parses a unified diff string into Patch structs
func Parse(ctx context.Context, diff string) ([]*Patch, error) {
	var patches []*Patch
	var currentPatch *Patch
	var inHunk bool
	_ = inHunk

	scanner := bufio.NewScanner(strings.NewReader(diff))

	for scanner.Scan() {
		line := scanner.Text()

		// Check for new file or deleted file header
		if strings.HasPrefix(line, "--- ") || strings.HasPrefix(line, "+++ ") {
			// Extract filename
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				filename := strings.TrimPrefix(fields[1], "a/") // handle a/ prefix
				filename = strings.TrimPrefix(filename, "b/")
				if currentPatch != nil && len(currentPatch.OldLines) > 0 {
					patches = append(patches, currentPatch)
				}
				currentPatch = &Patch{Path: filename}
			}
			continue
		}

		// Check for hunk header
		hunkHeader := regexp.MustCompile(`@@ -(\d+),?(\d*) \+(\d+),?(\d*) @@`)
		matches := hunkHeader.FindStringSubmatch(line)
		if matches != nil {
			inHunk = true
			_, _ = strconv.Atoi(matches[1])
			_, _ = strconv.Atoi(matches[2])
			_, _ = strconv.Atoi(matches[3])
			_, _ = strconv.Atoi(matches[4])
			continue
		}

		if !inHunk || currentPatch == nil {
			continue
		}

		// Parse hunk content
		if len(line) > 0 {
			switch line[0] {
			case ' ', '+':
				if line[0] == '+' {
					currentPatch.NewLines = append(currentPatch.NewLines, line[1:])
				} else {
					currentPatch.NewLines = append(currentPatch.NewLines, line[1:])
					currentPatch.OldLines = append(currentPatch.OldLines, line[1:])
				}
			case '-':
				currentPatch.OldLines = append(currentPatch.OldLines, line[1:])
			case '\\':
				// Handle "\ No newline at end of file"
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if currentPatch != nil && len(currentPatch.OldLines) > 0 {
		patches = append(patches, currentPatch)
	}

	return patches, nil
}

// ParseFile parses a diff file and returns patches
func ParseFile(ctx context.Context, filename string) ([]*Patch, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Parse(ctx, string(content))
}

// ApplyToFile applies patches from a diff file to the target
func ApplyToFile(ctx context.Context, diffFile, targetDir string) error {
	patches, err := ParseFile(ctx, diffFile)
	if err != nil {
		return fmt.Errorf("failed to parse diff: %w", err)
	}

	for _, p := range patches {
		p.Path = targetDir + "/" + p.Path
		if err := Apply(ctx, *p); err != nil {
			return fmt.Errorf("failed to apply patch to %s: %w", p.Path, err)
		}
	}

	return nil
}

// CreateDiff generates a unified diff string between old and new content
func CreateDiff(path, oldContent, newContent string) string {
	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")

	var diff strings.Builder

	diff.WriteString(fmt.Sprintf("--- %s\n", path))
	diff.WriteString(fmt.Sprintf("+++ %s\n", path))

	// Simple hunk generation (not handling complex cases)
	maxLines := len(oldLines)
	if len(newLines) > maxLines {
		maxLines = len(newLines)
	}

	// Find first difference
	start := 0
	for start < maxLines && start < len(oldLines) && start < len(newLines) {
		if oldLines[start] != newLines[start] {
			break
		}
		start++
	}

	// Find last difference
	end := 0
	for end < maxLines {
		oldIdx := len(oldLines) - 1 - end
		newIdx := len(newLines) - 1 - end
		if oldIdx >= 0 && newIdx >= 0 && oldLines[oldIdx] == newLines[newIdx] {
			break
		}
		end++
	}

	oldCount := len(oldLines) - start - end
	if oldCount < 0 {
		oldCount = 0
	}
	newCount := len(newLines) - start - end
	if newCount < 0 {
		newCount = 0
	}

	diff.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n", start+1, oldCount, start+1, newCount))

	for i := start; i < len(oldLines)-end; i++ {
		diff.WriteString(fmt.Sprintf("-%s\n", oldLines[i]))
	}
	for i := start; i < len(newLines)-end; i++ {
		diff.WriteString(fmt.Sprintf("+%s\n", newLines[i]))
	}

	return diff.String()
}