# Freecode — Testing Framework

## 8.0 Phase 8: Unit Testing

### 8.1 Overview

This document defines the testing strategy and implementation for Freecode, aligned with CloudBSD testing infrastructure guidelines.

**Key Principles (from CloudBSD Testing Infrastructure):**
- **Host Safety**: Kernel testing must never run on the development host directly
- **Reproducibility**: All test environments defined as code with known-clean snapshots
- **Structured Output**: All test output in AI-agent-parseable formats (JSON, TAP, JUnit XML)
- **Incremental Testing**: Fast unit tests first, then integration tests

### 8.2 Current Test Status

| Package | Tests | Coverage | Status |
|---------|-------|----------|--------|
| `agent` | 8 | 15.7% | ✅ PASS |
| `config` | 16 | 45.3% | ✅ PASS |
| `session` | 14 | 41.7% | ✅ PASS |
| `shell` | 6 | 56.5% | ✅ PASS |
| `tool` | 8 | 4.9% | ✅ PASS |
| `i18n` | 2 | 23.1% | ✅ PASS |
| **Total** | **54** | **~15%** | **✅ ALL PASS** |

### 8.3 Test Structure

```
internal/
├── agent/
│   └── engine_test.go        ✅ # Agent engine tests (8 tests)
├── config/
│   ├── config_test.go        ✅ # Config struct tests (11 tests)
│   ├── yaml_test.go         ✅ # YAML tests (4 tests)
│   ├── json_test.go         ✅ # JSON tests (5 tests)
│   ├── env_test.go          ✅ # ENV tests (4 tests)
│   ├── opencode/
│   │   └── migrate_test.go  ⏳ # OpenCode migration tests
│   └── omo/
│       └── merge_test.go    ⏳ # OMO merge tests
├── hook/
│   └── registry_test.go     ⏳ # Hook registration/execution tests
├── mcp/
│   ├── client_test.go       ⏳ # MCP client tests
│   ├── server_test.go       ⏳ # MCP server tests
│   └── builtin_test.go      ⏳ # Builtin MCP tests
├── provider/
│   ├── openai_test.go       ⏳ # OpenAI provider tests
│   └── anthropic_test.go    ⏳ # Anthropic provider tests
├── session/
│   ├── manager_test.go     ✅ # Session manager tests (14 tests)
│   └── store_test.go       ⏳ # Session store tests
├── shell/
│   └── executor_test.go    ✅ # Shell executor tests (6 tests)
├── tool/
│   ├── registry_test.go    ✅ # Tool registry tests (8 tests)
│   ├── bash_test.go        ⏳ # Bash tool tests
│   ├── read_test.go        ⏳ # Read tool tests
│   ├── write_test.go       ⏳ # Write tool tests
│   └── edit_test.go        ⏳ # Edit tool tests
├── ui/
│   ├── model_test.go       ⏳ # TUI model tests
│   ├── input_test.go       ⏳ # Input handler tests
│   ├── style_test.go       ⏳ # Style/theme tests
│   └── tab/
│       └── model_test.go   ⏳ # Tab model tests
└── i18n/
    └── loader_test.go      ✅ # i18n loader tests (2 tests)

cmd/
├── freecode/
│   └── main_test.go        ⏳ # CLI main tests
└── freecode-server/
    └── main_test.go        ⏳ # Server main tests
```

### 8.4 Test Naming Conventions

| Pattern | Description |
|---------|-------------|
| `Test<Function>` | Basic unit test |
| `Test<Function>_withValid` | Test with valid inputs |
| `Test<Function>_withInvalid` | Test with invalid inputs |
| `Test<Function>_withEmpty` | Test with empty inputs |
| `Benchmark<Function>` | Performance benchmark |

### 8.5 Test Categories

| Category | Description | Tools |
|----------|-------------|-------|
| **Unit Tests** | Test individual functions in isolation | `testing` package |
| **Integration Tests** | Test component interactions | `testing` + mocks |
| **Table-Driven Tests** | Test multiple input combinations | `[]struct` test cases |

### 8.6 Table-Driven Test Pattern

```go
func TestLoadYAML(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *Config
        wantErr bool
    }{
        {
            name:    "valid minimal config",
            input:   "shell: /bin/bash\nyolo: true\n",
            want:    &Config{Shell: "/bin/bash", Yolo: true},
            wantErr: false,
        },
        {
            name:    "invalid yaml syntax",
            input:   "shell: [unclosed",
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := LoadYAML([]byte(tt.input))
            if (err != nil) != tt.wantErr {
                t.Errorf("LoadYAML() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("LoadYAML() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 8.7 Mock Pattern

```go
type MockProvider struct {
    GenerateFunc func(ctx context.Context, req *Request) (*Response, error)
}

func (m *MockProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
    return m.GenerateFunc(ctx, req)
}

func TestAgentWithMockProvider(t *testing.T) {
    mock := &MockProvider{
        GenerateFunc: func(ctx context.Context, req *Request) (*Response, error) {
            return &Response{Content: "mock response"}, nil
        },
    }
    // test agent behavior with mock
}
```

### 8.8 Task List

| # | Task | Status | File | Coverage |
|---|------|--------|------|---------|
| 8.1 | Config struct tests | ✅ DONE | `config/config_test.go` | 45.3% |
| 8.2 | YAML parsing tests | ✅ DONE | `config/yaml_test.go` | included |
| 8.3 | JSON parsing tests | ✅ DONE | `config/json_test.go` | included |
| 8.4 | Env var tests | ✅ DONE | `config/env_test.go` | included |
| 8.5 | OpenCode migration tests | ⏳ TODO | `config/opencode/migrate_test.go` | 0% |
| 8.6 | OMO merge tests | ⏳ TODO | `config/omo/merge_test.go` | 0% |
| 8.7 | Tool registry tests | ✅ DONE | `tool/registry_test.go` | 4.9% |
| 8.8 | Bash tool tests | ⏳ TODO | `tool/bash_test.go` | 0% |
| 8.9 | Read tool tests | ⏳ TODO | `tool/read_test.go` | 0% |
| 8.10 | Write tool tests | ⏳ TODO | `tool/write_test.go` | 0% |
| 8.11 | Edit tool tests | ⏳ TODO | `tool/edit_test.go` | 0% |
| 8.12 | Glob tool tests | ⏳ TODO | `tool/glob_test.go` | 0% |
| 8.13 | Grep tool tests | ⏳ TODO | `tool/grep_test.go` | 0% |
| 8.14 | Agent engine tests | ✅ DONE | `agent/engine_test.go` | 15.7% |
| 8.15 | Session manager tests | ✅ DONE | `session/manager_test.go` | 41.7% |
| 8.16 | Session store tests | ⏳ TODO | `session/store_test.go` | 0% |
| 8.17 | Session compaction tests | ⏳ TODO | `session/compaction_test.go` | 0% |
| 8.18 | Hook registry tests | ⏳ TODO | `hook/registry_test.go` | 0% |
| 8.19 | Hook session tests | ⏳ TODO | `hook/session_test.go` | 0% |
| 8.20 | Hook tool tests | ⏳ TODO | `hook/tool_test.go` | 0% |
| 8.21 | Hook continuation tests | ⏳ TODO | `hook/continuation_test.go` | 0% |
| 8.22 | MCP client tests | ⏳ TODO | `mcp/client_test.go` | 0% |
| 8.23 | MCP server tests | ⏳ TODO | `mcp/server_test.go` | 0% |
| 8.24 | MCP builtin tests | ⏳ TODO | `mcp/builtin_test.go` | 0% |
| 8.25 | OpenAI provider tests | ⏳ TODO | `provider/openai_test.go` | 0% |
| 8.26 | Anthropic provider tests | ⏳ TODO | `provider/anthropic_test.go` | 0% |
| 8.27 | Shell executor tests | ✅ DONE | `shell/executor_test.go` | 56.5% |
| 8.28 | Shell PTY tests | ⏳ TODO | `shell/pty_test.go` | 0% |
| 8.29 | UI model tests | ⏳ TODO | `ui/model_test.go` | 0% |
| 8.30 | UI input tests | ⏳ TODO | `ui/input_test.go` | 0% |
| 8.31 | UI style tests | ⏳ TODO | `ui/style_test.go` | 0% |
| 8.32 | UI commands tests | ⏳ TODO | `ui/commands_test.go` | 0% |
| 8.33 | UI tab model tests | ⏳ TODO | `ui/tab/model_test.go` | 0% |
| 8.34 | UI tab view tests | ⏳ TODO | `ui/tab/view_test.go` | 0% |
| 8.35 | UI tab keys tests | ⏳ TODO | `ui/tab/keys_test.go` | 0% |
| 8.36 | UI tab split tests | ⏳ TODO | `ui/tab/split_test.go` | 0% |
| 8.37 | i18n loader tests | ✅ DONE | `i18n/loader_test.go` | 23.1% |
| 8.38 | Platform tests | ⏳ TODO | `platform/platform_test.go` | 0% |
| 8.39 | Server tests | ⏳ TODO | `server/server_test.go` | 0% |
| 8.40 | CLI root tests | ⏳ TODO | `cli/root_test.go` | 0% |

### 8.9 Test Coverage Goals

| Package | Current | Target | Critical |
|---------|---------|--------|----------|
| `config` | 45.3% | 80% | Yes |
| `tool` | 4.9% | 85% | Yes |
| `agent` | 15.7% | 75% | Yes |
| `session` | 41.7% | 80% | Yes |
| `shell` | 56.5% | 70% | No |
| `i18n` | 23.1% | 50% | No |
| `hook` | 0% | 70% | No |
| `mcp` | 0% | 70% | No |
| `provider` | 0% | 60% | No |
| `ui` | 0% | 50% | No |
| `cli` | 0% | 60% | No |
| `server` | 0% | 50% | No |
| `platform` | 0% | 40% | No |
| **Overall** | **~15%** | **65%** | |

### 8.10 Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package
go test ./internal/config/...

# Run with verbose
go test -v ./...

# Run benchmarks
go test -bench=. ./...

# Run specific test
go test -v -run TestLoadYAML ./internal/config/...
```

### 8.11 CI Integration

```yaml
# .github/workflows/test.yml
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Run tests
        run: go test -cover ./...
      - name: Upload coverage
        uses: codecov/codecov-action@v4
```

### 8.12 Test Data

Test fixtures should be in `testdata/` directories:

```
internal/config/
├── config_test.go
└── testdata/
    ├── valid.yaml
    ├── valid.json
    └── invalid.yaml

internal/tool/
├── bash_test.go
└── testdata/
    └── script.sh
```

Load test data using `os.ReadFile`:

```go
func TestLoadYAMLFile(t *testing.T) {
    data, err := os.ReadFile("testdata/valid.yaml")
    if err != nil {
        t.Skip("testdata not found")
    }

    cfg, err := LoadYAML(data)
    if err != nil {
        t.Errorf("LoadYAML() error = %v", err)
    }
    // assertions...
}
```

---

**Author:** Mark LaPointe <mark@cloudbsd.org>
**Last Updated:** 2026-05-02
