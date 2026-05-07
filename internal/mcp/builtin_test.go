package mcp

import (
	"testing"
)

func TestBuiltinMCP_ExaSearch(t *testing.T) {
	m := NewBuiltinMCP()
	if m == nil {
		t.Fatal("NewBuiltinMCP() returned nil")
	}

	res, err := m.ExaSearch("query", 5)
	if err != nil {
		t.Fatalf("ExaSearch() error: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("ExaSearch() returned %d results, want 1", len(res))
	}
	if res[0].Title != "Example Result" {
		t.Errorf("ExaSearch result title = %q, want %q", res[0].Title, "Example Result")
	}
}

func TestBuiltinMCP_Context7Docs(t *testing.T) {
	m := NewBuiltinMCP()
	if m == nil {
		t.Fatal("NewBuiltinMCP() returned nil")
	}

	docs, err := m.Context7Docs("test-ctx")
	if err != nil {
		t.Fatalf("Context7Docs() error: %v", err)
	}
	if docs != "Context7 documentation for: test-ctx" {
		t.Errorf("Context7Docs() = %q, want %q", docs, "Context7 documentation for: test-ctx")
	}
}

func TestBuiltinMCP_GrepApp(t *testing.T) {
	m := NewBuiltinMCP()
	if m == nil {
		t.Fatal("NewBuiltinMCP() returned nil")
	}

	urls, err := m.GrepApp("query")
	if err != nil {
		t.Fatalf("GrepApp() error: %v", err)
	}
	if len(urls) != 2 {
		t.Fatalf("GrepApp() returned %d urls, want 2", len(urls))
	}
	if urls[0] == urls[1] {
		t.Errorf("GrepApp() returned duplicate URLs: %q", urls)
	}
}
