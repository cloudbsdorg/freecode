package util

import (
	"testing"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		pattern string
		want    bool
	}{
		{"exact match", "hello", "hello", true},
		{"exact no match", "hello", "world", false},
		{"single star matches all", "hello world", "*", true},
		{"single star matches across slash", "hello/world", "*", true},
		{"double star matches all including slash", "hello/world", "**", true},
		{"question mark single char", "abc", "a?c", true},
		{"question mark no match", "ac", "a?c", false},
		{"pattern with ext", "file.txt", "*.txt", true},
		{"pattern with ext no match", "file.md", "*.txt", false},
		{"nested path", "src/foo/bar.go", "src/**/*.go", true},
		{"nested path no match", "src/foo/bar.txt", "src/**/*.go", false},
		{"single star in path", "src/foo/bar.go", "src/*/bar.go", true},
		{"single star matches across slash", "src/foo/bar/bar.go", "src/*/bar.go", true},
		{"double star prefix", "foo/bar/baz.go", "**/*.go", true},
		{"question mark in path", "a/b/c", "?/?/?", true},
		{"mixed star and question", "ab/cd/ef", "*/?/*", false},
		{"trailing star", "test.go", "*.go", true},
		{"backslash normalized", "a\\b\\c", "a/b/c", true},
		{"windows path", "C:\\Users\\test", "C:/Users/test", true},
		{"empty string", "", "", true},
		{"empty pattern", "hello", "", false},
		{"empty string with star", "", "*", true},
		{"question mark exact", "?", "?", true},
		{"star question star", "aXb", "*X*", true},
		{"path with dots", "dir./file", "dir.*/file", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Match(tt.str, tt.pattern)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.str, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestMatchTrailingStarOptional(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		pattern string
		want    bool
	}{
		{"rm star matches rm alone", "rm", "rm *", true},
		{"rm star matches rm with args", "rm -rf", "rm *", true},
		{"rm star matches rm file", "rm file.txt", "rm *", true},
		{"ls star matches ls alone", "ls", "ls *", true},
		{"ls star matches ls files", "ls -la", "ls *", true},
		{"no trailing space star", "rm -rf", "rm*", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Match(tt.str, tt.pattern)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.str, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestMatchAll(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		patterns map[string]any
		wantVal  any
		wantOk   bool
	}{
		{
			"exact match",
			"hello",
			map[string]any{"hello": "world"},
			"world",
			true,
		},
		{
			"no match",
			"unknown",
			map[string]any{"hello": "world"},
			nil,
			false,
		},
		{
			"longer glob wins",
			"abc",
			map[string]any{"ab*": "wildcard", "abc*": "exact"},
			"exact",
			true,
		},
		{
			"exact match takes precedence",
			"abc",
			map[string]any{"ab*": "wildcard", "abc": "exact"},
			"exact",
			true,
		},
		{
			"star pattern",
			"test.txt",
			map[string]any{"*.txt": "text file", "*.md": "markdown"},
			"text file",
			true,
		},
		{
			"empty patterns",
			"anything",
			map[string]any{},
			nil,
			false,
		},
		{
			"double star pattern",
			"path/to/file.go",
			map[string]any{"**/*.go": "go file"},
			"go file",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := MatchAll(tt.input, tt.patterns)
			if ok != tt.wantOk {
				t.Errorf("MatchAll(%q, _) ok = %v, want %v", tt.input, ok, tt.wantOk)
			}
			if val != tt.wantVal {
				t.Errorf("MatchAll(%q, _) = %v, want %v", tt.input, val, tt.wantVal)
			}
		})
	}
}

func TestMatchAllSorted(t *testing.T) {
	patterns := map[string]any{
		"a":        "a_val",
		"ab":       "ab_val",
		"abc":      "abc_val",
		"ab*":      "ab*_val",
		"abcd":     "abcd_val",
		"abcdefg":  "abcdefg_val",
		"abcdefgh": "abcdefgh_val",
	}

	val, ok := MatchAll("abc", patterns)
	if !ok {
		t.Error("expected match for abc")
	}
	// Last matching pattern wins (sorted by key length, then alphabetically)
	if val != "abc_val" {
		t.Errorf("expected abc_val (last matching pattern), got %v", val)
	}
}

func TestMatchCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		patterns map[string]any
		wantVal  any
		wantOk   bool
	}{
		{
			"exact command match",
			"rm file.txt",
			map[string]any{"rm *": "remove"},
			"remove",
			true,
		},
		{
			"command no match",
			"ls file.txt",
			map[string]any{"rm *": "remove"},
			nil,
			false,
		},
		{
			"star matches multiple args",
			"rm -rf /",
			map[string]any{"rm *": "remove"},
			"remove",
			true,
		},
		{
			"command with glob file",
			"rm *.txt",
			map[string]any{"rm *.txt": "remove text files"},
			"remove text files",
			true,
		},
		{
			"command empty input",
			"",
			map[string]any{"rm *": "remove"},
			nil,
			false,
		},
		{
			"command exact head match",
			"git commit -m",
			map[string]any{"git *": "git cmd"},
			"git cmd",
			true,
		},
		{
			"multiple star segments",
			"cmd a b c",
			map[string]any{"cmd * *": "two args"},
			"two args",
			true,
		},
		{
			"star matches one then rest",
			"cmd foo bar baz",
			map[string]any{"cmd * baz": "ends with baz"},
			"ends with baz",
			true,
		},
		{
			"star at end",
			"npm test",
			map[string]any{"npm test": "run tests", "npm *": "npm other"},
			"run tests",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := MatchCommand(tt.input, tt.patterns)
			if ok != tt.wantOk {
				t.Errorf("MatchCommand(%q, _) ok = %v, want %v", tt.input, ok, tt.wantOk)
			}
			if val != tt.wantVal {
				t.Errorf("MatchCommand(%q, _) = %v, want %v", tt.input, val, tt.wantVal)
			}
		})
	}
}

func TestMatchCommandSorted(t *testing.T) {
	patterns := map[string]any{
		"rm *":      "remove anything",
		"rm":        "remove nothing",
		"rm -rf *":  "force remove",
		"rm *.txt":  "remove txt",
	}

	val, ok := MatchCommand("rm -rf /", patterns)
	if !ok {
		t.Error("expected match")
	}
	if val != "force remove" {
		t.Errorf("expected 'force remove' (last matching pattern), got %v", val)
	}
}

func TestMatchSequence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		pattern  string
		wantOk   bool
	}{
		{"star matches all", "a b c", "*", true},
		{"star matches one", "a", "*", true},
		{"exact match", "a", "a", true},
		{"exact no match", "a", "b", false},
		{"star then exact tail", "a b", "* b", true},
		{"star then exact tail fails", "a c", "* b", false},
		{"exact head then star", "a b", "a *", true},
		{"two stars", "a b c", "* *", true},
		{"star consumes one", "a b c", "* b c", true},
		{"star consumes none", "a b c", "a * c", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := []string{}
			if tt.input != "" {
				words = splitWords(tt.input)
			}
			parts := splitWords(tt.pattern)
			ok := matchSequence(words, parts)
			if ok != tt.wantOk {
				t.Errorf("matchSequence(%v, %v) = %v, want %v", words, parts, ok, tt.wantOk)
			}
		})
	}
}

func splitWords(s string) []string {
	if s == "" {
		return []string{}
	}
	words := make([]string, 0)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			if start < i {
				words = append(words, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		words = append(words, s[start:])
	}
	return words
}

func BenchmarkMatch(b *testing.B) {
	patterns := []string{
		"**/*.go",
		"src/**/*.ts",
		"*/build/*",
		"test/**/*.js",
		"*.json",
	}
	inputs := []string{
		"path/to/file.go",
		"src/components/App.ts",
		"dist/build/app.js",
		"test/utils/helper.js",
		"config.json",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range patterns {
			for _, inp := range inputs {
				Match(inp, p)
			}
		}
	}
}

func BenchmarkMatchAll(b *testing.B) {
	patterns := map[string]any{
		"*.txt":      "text",
		"src/**/*.go": "go source",
		"test/**":    "test files",
		"docs/*":     "docs",
		"**/*.md":    "markdown",
		"bin/*":      "binaries",
	}
	input := "src/util/wildcard.go"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MatchAll(input, patterns)
	}
}

func BenchmarkMatchCommand(b *testing.B) {
	patterns := map[string]any{
		"rm *":     "remove",
		"rm -rf *": "force remove",
		"ls *":     "list",
		"git *":    "git",
		"npm *":    "npm",
	}
	input := "rm -rf /path/to/dir"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MatchCommand(input, patterns)
	}
}
