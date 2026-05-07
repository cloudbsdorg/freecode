# Freecode — LSP Implementation Plan

**Document ID:** LSP-Implementation-Plan
**Version:** 1.0
**Created:** 2026-05-06
**Author:** Sisyphus
**Status:** ACTIVE
**Classification:** INTERNAL

---

## 1.0 Current State

### 1.1 Implementation Complete (2026-05-06)

All phases completed:

- ✅ Syntax errors fixed (`map[string]any{}{` → `map[string]any{`)
- ✅ Bidirectional handlers using `jsonrpc2.HandlerWithError`
- ✅ stdin/stdout wrapper (`stdinStdout` struct)
- ✅ Diagnostic store with debouncing (`diagnostic.go`)
- ✅ Server lifecycle management (`server.go`)
- ✅ Language detection (`language.go`)
- ✅ Tool integration (`tool/lsp.go`)

### 1.2 Files Status

| File | Lines | Status |
|------|-------|--------|
| `internal/lsp/lsp.go` | ~560 | ✅ COMPLETE |
| `internal/lsp/diagnostic.go` | ~185 | ✅ COMPLETE |
| `internal/lsp/server.go` | ~150 | ✅ COMPLETE |
| `internal/lsp/language.go` | ~120 | ✅ COMPLETE |
| `internal/tool/lsp.go` | ~330 | ✅ COMPLETE |

### 1.3 Reference Implementation

**TypeScript Source:** `/home/mlapointe/secure/git/opencode/packages/opencode/src/lsp/`

| File | Lines | Purpose |
|------|-------|---------|
| `client.ts` | 697 | Bidirectional LSP client with diagnostics |
| `server.ts` | 60k+ | LSP server implementation |
| `lsp.ts` | 517 | LSP types and utilities |
| `language.ts` | 2,559 | Language detection |
| `launch.ts` | 794 | LSP server spawning |
| `diagnostic.ts` | 900 | Diagnostic handling |

---

## 2.0 Implementation Status

### 2.1 ✅ Phase 1: Bidirectional Communication - DONE

**Bidirectional handlers using `jsonrpc2.HandlerWithError`**

| # | Task | File | Status |
|---|------|------|--------|
| 1.2.1 | Add stdinStdout wrapper | `lsp.go` | ✅ DONE |
| 1.2.2 | Implement `window/workDoneProgress/create` | `lsp.go` | ✅ DONE |
| 1.2.3 | Implement `workspace/configuration` | `lsp.go` | ✅ DONE |
| 1.2.4 | Implement `client/registerCapability` | `lsp.go` | ✅ DONE |
| 1.2.5 | Implement `client/unregisterCapability` | `lsp.go` | ✅ DONE |
| 1.2.6 | Implement `workspace/workspaceFolders` | `lsp.go` | ✅ DONE |
| 1.2.7 | Implement `workspace/diagnostic/refresh` | `lsp.go` | ✅ DONE |
| 1.2.8 | Implement `textDocument/publishDiagnostics` | `lsp.go` | ✅ DONE |
| 1.2.9 | Wire handlers to jsonrpc2 | `lsp.go` | ✅ DONE |

---

### 2.2 ✅ Phase 2: Diagnostics System - DONE

**Diagnostic store with 150ms debouncing implemented**

| # | Task | File | Status |
|---|------|------|--------|
| 2.2.1 | Create diagnostic types | `diagnostic.go` | ✅ DONE |
| 2.2.2 | Add diagnostics state | `diagnostic.go` | ✅ DONE |
| 2.2.3 | Implement debouncing | `diagnostic.go` | ✅ DONE |
| 2.2.4 | Implement `textDocument/diagnostic` | `lsp.go` | ✅ DONE |
| 2.2.5 | Add diagnostic callbacks | `diagnostic.go` | ✅ DONE |

---

### 2.3 ✅ Phase 3: Server Management - DONE

**Server lifecycle management with auto-detection**

| # | Task | File | Status |
|---|------|------|--------|
| 3.3.1 | Create server detection | `server.go` | ✅ DONE |
| 3.3.2 | Add server process management | `server.go` | ✅ DONE |
| 3.3.3 | Implement language detection | `language.go` | ✅ DONE |
| 3.3.4 | Add server initialization options | `server.go` | ✅ DONE |
| 3.3.5 | Implement server status | `server.go` | ✅ DONE |

---

### 2.4 ✅ Phase 4: Tool Integration - DONE

**Full tool integration with hover, definition, references, completion, diagnostics**

| # | Task | File | Status |
|---|------|------|--------|
| 4.4.1 | Wire LSP client to tool | `tool/lsp.go` | ✅ DONE |
| 4.4.2 | Implement `lsp_complete` | `tool/lsp.go` | ✅ DONE |
| 4.4.3 | Implement `lsp_definition` | `tool/lsp.go` | ✅ DONE |
| 4.4.4 | Implement `lsp_references` | `tool/lsp.go` | ✅ DONE |
| 4.4.5 | Implement `lsp_hover` | `tool/lsp.go` | ✅ DONE |
| 4.4.6 | Implement `lsp_diagnostics` | `tool/lsp.go` | ✅ DONE |

---

## 3.0 File Structure

```
internal/lsp/
├── lsp.go           # Main LSP client (560 lines)
├── diagnostic.go    # Diagnostic types and debouncing (185 lines)
├── server.go        # Server lifecycle management (150 lines)
└── language.go      # Language detection (120 lines)

internal/tool/
└── lsp.go          # LSP tool for agents (330 lines)
```

---

## 4.0 Verification Checklist

### Build Verification:
```bash
go build ./internal/lsp/
go build ./internal/tool/
go build ./...
```

### Functionality Verification:
```bash
# Start LSP server for Go
./freecode lsp start --language go

# Test hover
./freecode lsp hover --file ./internal/lsp/lsp.go --line 10 --character 1

# Test definition
./freecode lsp definition --file ./internal/lsp/lsp.go --line 10 --character 1

# Test references
./freecode lsp references --file ./internal/lsp/lsp.go --line 10 --character 1

# Test completion
./freecode lsp completion --file ./internal/lsp/lsp.go --line 10 --character 1

# Test diagnostics
./freecode lsp diagnostics --file ./internal/lsp/lsp.go
```

---

## 5.0 Dependencies

| Dependency | Status |
|------------|--------|
| `github.com/sourcegraph/jsonrpc2` | ✅ In go.mod |
| `internal/tool` | ✅ Used for tool registration |

---

## 6.0 Change Log

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 2.0 | 2026-05-06 | Sisyphus | Full implementation complete |
| 1.0 | 2026-05-06 | Sisyphus | Initial implementation plan |

---

**Classification:** INTERNAL

---

**Classification:** INTERNAL
