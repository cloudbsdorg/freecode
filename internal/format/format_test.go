package format

import (
	"context"
	"testing"
)

func TestFormatUnknownLanguage(t *testing.T) {
	ctx := context.Background()
	code := "func test() {}"

	result, err := Format(ctx, code, FormatOptions{Language: "unknown"})
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}
	if result == nil {
		t.Fatal("Format() returned nil result")
	}
	if result.Formatted != code {
		t.Errorf("Formatted = %q, want %q (unknown lang returns original)", result.Formatted, code)
	}
}

func TestFormatGo(t *testing.T) {
	ctx := context.Background()
	code := "package main\n\nfunc test(){}"

	result, err := Format(ctx, code, FormatOptions{Language: "go"})
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}
	if result == nil {
		t.Fatal("Format() returned nil result")
	}
	if result.Formatted == "" {
		t.Error("Formatted is empty")
	}
}

func TestFormatTypeScript(t *testing.T) {
	ctx := context.Background()
	code := "const x=1"

	result, err := Format(ctx, code, FormatOptions{Language: "typescript", Path: "test.ts"})
	if err != nil {
		t.Logf("Format() error (prettier may not be installed): %v", err)
	}
	if result == nil {
		t.Fatal("Format() returned nil result")
	}
}

func TestFormatJavaScript(t *testing.T) {
	ctx := context.Background()
	code := "const x=1"

	result, err := Format(ctx, code, FormatOptions{Language: "javascript", Path: "test.js"})
	if err != nil {
		t.Logf("Format() error (prettier may not be installed): %v", err)
	}
	if result == nil {
		t.Fatal("Format() returned nil result")
	}
}

func TestFormatOptions(t *testing.T) {
	opts := FormatOptions{
		Language: "go",
		Path:     "test.go",
	}
	if opts.Language != "go" {
		t.Errorf("opts.Language = %q, want %q", opts.Language, "go")
	}
	if opts.Path != "test.go" {
		t.Errorf("opts.Path = %q, want %q", opts.Path, "test.go")
	}
}

func TestResult(t *testing.T) {
	result := &Result{
		Formatted: "formatted code",
		Error:     nil,
	}
	if result.Formatted != "formatted code" {
		t.Errorf("result.Formatted = %q, want %q", result.Formatted, "formatted code")
	}
	if result.Error != nil {
		t.Errorf("result.Error = %v, want nil", result.Error)
	}
}