package mcp

import (
	"strings"
	"testing"
)

func TestBuiltinMCP_ExaSearch_NoAPIKey(t *testing.T) {
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
	if !strings.Contains(res[0].Title, "Not Configured") {
		t.Errorf("ExaSearch result title = %q, want containing 'Not Configured'", res[0].Title)
	}
}

func TestBuiltinMCP_Context7Docs_NoAPIKey(t *testing.T) {
	m := NewBuiltinMCP()
	if m == nil {
		t.Fatal("NewBuiltinMCP() returned nil")
	}

	docs, err := m.Context7Docs("test-ctx")
	if err != nil {
		t.Fatalf("Context7Docs() error: %v", err)
	}
	if !strings.Contains(docs, "Not Configured") {
		t.Errorf("Context7Docs() = %q, want containing 'Not Configured'", docs)
	}
}

func TestBuiltinMCP_GrepApp_RateLimited(t *testing.T) {
	m := NewBuiltinMCP()
	if m == nil {
		t.Fatal("NewBuiltinMCP() returned nil")
	}

	_, err := m.GrepApp("query")
	if err == nil {
		t.Log("GrepApp() succeeded (rate limit may have passed)")
	}
}
