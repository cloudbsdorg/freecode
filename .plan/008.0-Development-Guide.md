# Freecode — Development Guide

## 1.0 Purpose

This document provides guidance for developers working on freecode.

---

## 2.0 Getting Started

### 2.1 Prerequisites

**Required:**
- Go 1.22+
- git
- make

**Optional:**
- golangci-lint (for linting)
- goreleaser (for cross-platform builds)
- staticcheck (for static analysis)

### 2.2 Initial Setup

```bash
# Clone the repository
git clone https://github.com/cloudbsdorg/freecode.git
cd freecode

# Download dependencies
go mod download

# Build the CLI
go build -o freecode ./cmd/freecode

# Build the server
go build -o freecode-server ./cmd/freecode-server
```

### 2.3 Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/config/...
go test ./internal/tool/...
```

### 2.4 Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Run linter on specific package
golangci-lint run ./internal/config/...
```

---

## 3.0 Project Structure

```
freecode/
├── cmd/
│   ├── freecode/           # CLI entry point
│   └── freecode-server/   # Server entry point
├── internal/
│   ├── cli/                # Cobra commands
│   ├── config/             # Configuration
│   ├── agent/              # Agent engine
│   ├── tool/               # Tools
│   ├── hook/               # Hooks
│   ├── session/            # Session management
│   ├── ui/                 # TUI
│   ├── server/              # HTTP server
│   ├── mcp/                # MCP client
│   └── platform/           # Platform-specific
├── pkg/
│   └── api/                # SDK
├── .plan/                  # Plan documents
└── packaging/              # Distribution
```

---

## 4.0 Code Style

### 4.1 Formatting

```bash
# Format code
go fmt ./...

# Run before commit
```

### 4.2 Error Handling

```go
// Good - wrap errors with context
if err != nil {
    return fmt.Errorf("failed to execute tool %s: %w", toolName, err)
}

// Bad - bare errors
if err != nil {
    return err
}
```

### 4.3 Context Usage

```go
// Every function that does I/O should accept context
func ExecuteTool(ctx context.Context, req *Request) (*Response, error) {
    // Use ctx for cancellation and timeouts
}
```

### 4.4 Logging

```go
import "log/slog"

// Use structured logging
slog.Info("starting tool", "tool", toolName, "session", sessionID)
slog.Error("tool failed", "error", err)
```

---

## 5.0 Adding a New Tool

### 5.1 Tool Interface

```go
// internal/tool/tool.go
type Tool interface {
    Name() string
    Description() string
    Execute(ctx context.Context, req *Request) (*Response, error)
    Schema() ToolSchema
}
```

### 5.2 Example Tool

```go
// internal/tool/example.go
type exampleTool struct{}

func (t *exampleTool) Name() string { return "example" }
func (t *exampleTool) Description() string { return "An example tool" }

func (t *exampleTool) Execute(ctx context.Context, req *Request) (*Response, error) {
    // Implementation
    return &Response{Output: "result"}, nil
}

func (t *exampleTool) Schema() ToolSchema {
    return ToolSchema{
        Name:        "example",
        Description: "An example tool",
        InputSchema:  {...},
    }
}

// Register in tool registry
func init() {
    RegisterTool(&exampleTool{})
}
```

### 5.3 Registering in CLI

```go
// internal/cli/root.go
func init() {
    // Tools are registered via init() in each tool file
}
```

---

## 6.0 Adding a New Agent

### 6.1 Agent Structure

```go
// internal/agent/agent.go
type Agent interface {
    Name() string
    Mode() AgentMode // primary|subagent|all
    SystemPrompt() string
    Execute(ctx context.Context, req *Request) (*Response, error)
}
```

### 6.2 Built-in Agent

```go
// internal/agent/sisyphus.go
type sisyphusAgent struct{}

func (a *sisyphusAgent) Name() string { return "sisyphus" }
func (a *sisyphusAgent) Mode() AgentMode { return AgentModePrimary }

func (a *sisyphusAgent) SystemPrompt() string {
    return `You are Sisyphus, the eternal orchestrator...`
}

func (a *sisyphusAgent) Execute(ctx context.Context, req *Request) (*Response, error) {
    // Implementation
}
```

---

## 7.0 Adding a New Hook

### 7.1 Hook Types

```go
// Session hooks - called on session events
type SessionHook func(ctx context.Context, evt SessionEvent) error

// Tool hooks - called before/after tool execution
type ToolHook func(ctx context.Context, evt ToolEvent) (error, bool) // bool=handled

// Transform hooks - modify messages
type TransformHook func(msg *Message) (*Message, error)
```

### 7.2 Registering Hooks

```go
// internal/hook/registry.go
func (r *Registry) Register(hookType, name string, hook interface{}) error {
    // Registration implementation
}
```

---

## 8.0 Configuration Schema

### 8.1 Adding Config Fields

```go
// internal/config/config.go
type Config struct {
    // Add new field
    NewFeature string `yaml:"newFeature" json:"newFeature"`

    // Existing fields...
}
```

### 8.2 Config Validation

```go
func (c *Config) Validate() error {
    if c.NewFeature != "" && !validFeature(c.NewFeature) {
        return fmt.Errorf("invalid newFeature: %s", c.NewFeature)
    }
    return nil
}
```

---

## 9.0 Debugging

### 9.1 Verbose Logging

```bash
FREECODE_LOG_LEVEL=DEBUG ./freecode run "hello"
```

### 9.2 Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug CLI
dlv debug ./cmd/freecode

# Debug tests
dlv test ./internal/tool/...
```

### 9.3 Testing Specific Functions

```bash
# Run test with verbose output
go test -v ./internal/config/...

# Run test with debugger
dlv test ./internal/config/...
```

---

## 10.0 Building Releases

### 10.1 Local Build

```bash
# Build for current platform
go build -o freecode ./cmd/freecode

# Cross-compile for FreeBSD
GOOS=freebsd GOARCH=amd64 go build -o freecode-freebsd-amd64 ./cmd/freecode

# Cross-compile for all platforms
goreleaser build --snapshot --clean
```

### 10.2 Package Managers

```bash
# FreeBSD
cd packaging/freebsd && make package

# macOS (Homebrew)
cd packaging/macos
brew formula ./Formula/freecode.rb

# Linux (Flatpak)
flatpak-builder --user --install build-dir com.freecode.Freecode.yml
```

---

## 11.0 Performance Profiling

### 11.1 CPU Profile

```go
import _ "net/http/pprof"

// In your server:
go http.ListenAndServe(":6060", nil)

// Then:
go tool pprof http://localhost:6060/debug/pprof/profile
```

### 11.2 Memory Profile

```bash
go tool pprof http://localhost:6060/debug/pprof/heap
```

### 11.3 Trace

```bash
go tool trace http://localhost:6060/debug/pprof/trace?seconds=10
```

---

## 12.0 Common Issues

### 12.1 "cannot find package"

```bash
go mod tidy
go mod download
```

### 12.2 Tests failing

```bash
# Clean test cache
go clean -testcache
go test ./...
```

### 12.3 Linter errors

```bash
# Update linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run with auto-fix where possible
golangci-lint run ./... --fix
```

---

**Author:** Mark LaPointe <mark@cloudbsd.org>
**Last Updated:** 2026-05-01

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author
