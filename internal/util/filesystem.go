// Package util provides filesystem utilities.
package util

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// Exists checks if the given path exists.
func Exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// IsDir checks if the given path is a directory.
func IsDir(p string) bool {
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Stat returns file info for the given path.
// Returns nil error if the path does not exist (instead of an error).
func Stat(p string) (os.FileInfo, error) {
	info, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return info, nil
}

// ReadText reads the file at p and returns its contents as a string.
func ReadText(p string) (string, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return "", fmt.Errorf("failed to read text from %s: %w", p, err)
	}
	return string(data), nil
}

// ReadBytes reads the file at p and returns its contents as bytes.
func ReadBytes(p string) ([]byte, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes from %s: %w", p, err)
	}
	return data, nil
}

// ReadJson reads and parses JSON from the file at p into type T.
func ReadJson[T any](p string) (T, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to read JSON from %s: %w", p, err)
	}
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("failed to parse JSON from %s: %w", p, err)
	}
	return result, nil
}

// Write writes the given content to the file at p, creating parent directories as needed.
// If the file does not exist, it is created with the specified mode.
func Write(p string, content []byte, mode os.FileMode) error {
	// Ensure parent directory exists
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write the file
	if err := os.WriteFile(p, content, mode); err != nil {
		return fmt.Errorf("failed to write to %s: %w", p, err)
	}

	// Set the final mode
	if err := os.Chmod(p, mode); err != nil {
		return fmt.Errorf("failed to set mode on %s: %w", p, err)
	}

	return nil
}

// WriteJson writes data as indented JSON to the file at p.
func WriteJson(p string, data any, mode os.FileMode) error {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	content = append(content, '\n')
	return Write(p, content, mode)
}

// Resolve resolves symlinks and normalizes the path.
// Returns the canonical path.
func Resolve(p string) string {
	// First resolve any Windows-style paths
	p = WindowsPath(p)

	// Resolve the path (handles .., . etc) and symlinks
	resolved, err := filepath.EvalSymlinks(p)
	if err != nil {
		// If the path doesn't exist, just normalize it
		resolved = filepath.Clean(p)
	}

	// On Windows, also normalize casing
	if runtime.GOOS == "windows" {
		resolved = normalizePathWindows(resolved)
	}

	return resolved
}

// PathContains checks if parent contains child (child is at or under parent).
// Returns true if child is equal to or nested within parent.
func PathContains(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return !filepath.IsAbs(rel) && rel != ".." && !hasDotDotPrefix(rel)
}

// hasDotDotPrefix checks if the path starts with ".." (parent dir).
// Returns false for filenames like "..foo".
func hasDotDotPrefix(p string) bool {
	if len(p) >= 2 && p[0] == '.' && p[1] == '.' {
		if len(p) == 2 {
			return true
		}
		if p[2] == filepath.Separator || p[2] == '/' {
			return true
		}
	}
	return false
}

// FindUp walks up from start to stop looking for files matching target.
// Returns all matching paths found.
func FindUp(target string, start string, stop string) ([]string, error) {
	return FindUpMulti([]string{target}, start, stop, false)
}

// FindUpMulti is the internal implementation that handles multiple targets.
// The stop directory itself is not searched - search stops before reaching it.
func FindUpMulti(targets []string, start string, stop string, rootFirst bool) ([]string, error) {
	dirs := []string{start}

	// Walk up the directory tree
	current := start
	for {
		parent := filepath.Dir(current)
		if parent == current {
			// Reached the root
			break
		}
		if stop != "" && parent == stop {
			// Stop before adding the stop directory
			break
		}
		dirs = append(dirs, parent)
		current = parent
	}

	if rootFirst {
		// Reverse to search from root down
		for i, j := 0, len(dirs)-1; i < j; i, j = i+1, j-1 {
			dirs[i], dirs[j] = dirs[j], dirs[i]
		}
	}

	var results []string
	for _, dir := range dirs {
		for _, target := range targets {
			candidate := filepath.Join(dir, target)
			if Exists(candidate) {
				results = append(results, candidate)
			}
		}
	}

	return results, nil
}

// WindowsPath converts various Windows path formats to standard Windows paths.
// Handles Git Bash (/c/...), Cygwin (/cygdrive/c/...), and WSL (/mnt/c/...) formats.
func WindowsPath(p string) string {
	if runtime.GOOS == "windows" {
		return p
	}

	// Git Bash paths: /c/... or /C/...
	if match, n := parseGitBashPath(p); n > 0 {
		return match
	}

	// Cygwin paths: /cygdrive/c/...
	if match, n := parseCygwinPath(p); n > 0 {
		return match
	}

	// WSL paths: /mnt/c/...
	if match, n := parseWSLPath(p); n > 0 {
		return match
	}

	return p
}

// parseGitBashPath matches /<drive>/... or /<drive> patterns
func parseGitBashPath(p string) (string, int) {
	// Match /<drive>/... where drive is a single letter
	if len(p) >= 3 && p[0] == '/' {
		drive := p[1]
		if (drive >= 'a' && drive <= 'z') || (drive >= 'A' && drive <= 'Z') {
			if p[2] == '/' || p[2] == '\\' {
				return fmt.Sprintf("%c:/%s", uppercaseChar(drive), p[3:]), len(p)
			}
		}
	}
	return "", 0
}

// parseCygwinPath matches /cygdrive/<drive>/...
func parseCygwinPath(p string) (string, int) {
	const prefix = "/cygdrive/"
	if len(p) > len(prefix) {
		if p[:len(prefix)] == prefix {
			drive := p[len(prefix)]
			if (drive >= 'a' && drive <= 'z') || (drive >= 'A' && drive <= 'Z') {
				if len(p) > len(prefix)+1 && (p[len(prefix)+1] == '/' || p[len(prefix)+1] == '\\') {
					return fmt.Sprintf("%c:/%s", uppercaseChar(drive), p[len(prefix)+2:]), len(p)
				}
			}
		}
	}
	return "", 0
}

// parseWSLPath matches /mnt/<drive>/...
func parseWSLPath(p string) (string, int) {
	const prefix = "/mnt/"
	if len(p) > len(prefix) {
		if p[:len(prefix)] == prefix {
			drive := p[len(prefix)]
			if (drive >= 'a' && drive <= 'z') || (drive >= 'A' && drive <= 'Z') {
				if len(p) > len(prefix)+1 && (p[len(prefix)+1] == '/' || p[len(prefix)+1] == '\\') {
					return fmt.Sprintf("%c:/%s", uppercaseChar(drive), p[len(prefix)+2:]), len(p)
				}
			}
		}
	}
	return "", 0
}

// uppercaseChar converts a byte to uppercase if it's a letter.
func uppercaseChar(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 'a' + 'A'
	}
	return b
}

// normalizePathWindows normalizes a Windows path to canonical casing.
// On Windows, paths are case-insensitive but the filesystem may have
// different casing than what we receive from LSP servers.
func normalizePathWindows(p string) string {
	// Normalize separators
	p = filepath.Clean(p)

	// Use the native filesystem to resolve the real casing
	// Try realpath which handles this on Windows
	if resolved, err := filepath.EvalSymlinks(p); err == nil {
		return resolved
	}

	// If eval symlinks fails, just return the cleaned path
	return p
}

// WriteStream writes a stream to the file at p.
func WriteStream(p string, r io.Reader, mode os.FileMode) error {
	// Ensure parent directory exists
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Create the file
	f, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", p, err)
	}
	defer f.Close()

	// Copy the data
	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("failed to write to %s: %w", p, err)
	}

	// Set the mode
	if err := f.Chmod(mode); err != nil {
		return fmt.Errorf("failed to set mode on %s: %w", p, err)
	}

	return nil
}

// Size returns the size of the file at p, or 0 if it doesn't exist.
func Size(p string) (int64, error) {
	info, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	return info.Size(), nil
}

// Copy copies a file from src to dst.
func Copy(dst, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source %s: %w", src, err)
	}
	defer srcFile.Close()

	return WriteStream(dst, srcFile, 0644)
}

// isEnoent checks if the error is an ENOENT (file not found) error.
func isEnoent(err error) bool {
	return os.IsNotExist(err)
}
