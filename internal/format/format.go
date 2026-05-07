package format

import (
	"context"
	"os/exec"
	"strings"
)

type FormatOptions struct {
	Language string
	Path     string
}

type Result struct {
	Formatted string
	Error     error
}

func Format(ctx context.Context, code string, opts FormatOptions) (*Result, error) {
	switch opts.Language {
	case "go":
		return formatGo(ctx, code)
	case "typescript", "javascript":
		return formatPrettier(ctx, code, opts)
	default:
		return &Result{Formatted: code}, nil
	}
}

func formatGo(ctx context.Context, code string) (*Result, error) {
	cmd := exec.CommandContext(ctx, "gofmt", "-w")
	cmd.Stdin = strings.NewReader(code)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return &Result{Formatted: string(out)}, nil
}

func formatPrettier(ctx context.Context, code string, opts FormatOptions) (*Result, error) {
	cmd := exec.CommandContext(ctx, "npx", "prettier", "--stdin-filepath", opts.Path)
	cmd.Stdin = strings.NewReader(code)
	out, err := cmd.Output()
	if err != nil {
		return &Result{Formatted: code, Error: err}, nil
	}
	return &Result{Formatted: string(out)}, nil
}
