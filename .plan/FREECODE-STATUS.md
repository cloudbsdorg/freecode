# Freecode — Implementation Status & Revised Plan

**Last Updated:** 2026-05-04
**Author:** Mark LaPointe <mark@cloudbsd.org>
**Build Status:** ✅ Builds Successfully | ⚠️ Tests Unknown (timeout)

---

## Executive Summary

**TRUE STATUS: Phase 1基础设施完整，但大多数模块是存根。**

| Category | Claimed Done | Actually Complete | Stub Only |
|----------|-------------|-------------------|-----------|
| CLI Commands | 17 | ~12 | ~5 |
| Core Modules | 16 (Phase 2) | ~4 | ~12 |
| Extended Modules | 12 (Phase 3) | ~3 | ~9 |
| **Total Tasks Claimed** | **~45** | **~19 (42%)** | **~26 (58%)** |

**Critical Finding:** The TASKLISTs mark modules as "Done" when only a single stub file exists (1-2 Go files vs OpenCode's 5-30 source files). True feature parity requires significant additional implementation.

---

## Part 1: TRUE Implementation Status

### 1.1 Package-by-Package Assessment

| Package | OpenCode Files | Freecode Go Files | True Status | Notes |
|---------|---------------|-------------------|-------------|-------|
| **agent** | 1 | 8+6t | ✅ REAL | Engine, prompts, streaming all implemented |
| **auth** | 1 | 1 | ✅ REAL | Simple credential store |
| **bus** | 3 | 2+1t | ⚠️ STUB | EventBus interface exists, wildcard/session not fully implemented |
| **command** | 1 | 1+1t | ⚠️ STUB | Basic registry exists, template system missing |
| **config** | 22 | 11+7t | ⚠️ PARTIAL | Most config formats done, env var merging incomplete |
| **effect** | 10 | 1+1t | 🚨 STUB | Only Effect struct, no runtime |
| **git** | 1 | 1+1t | ✅ REAL | Uses go-git, basic operations work |
| **hook** | 0 (new) | 7+2t | ✅ REAL | Full hook system with 52 triggers |
| **i18n** | 0 (new) | 2+1t | ✅ REAL | RTL support, loader working |
| **ide** | 1 | 1+1t | ⚠️ STUB | Handler stub only, no LSP integration |
| **installation** | 1 | 1 | ✅ REAL | Version detection works |
| **lsp** | 6 | 5 | ✅ REAL | Full LSP implementation: client, diagnostics, server mgmt, language detection, tool |
| **mcp** | 4 | 4+2t | ✅ REAL | Client, server, OAuth, builtin all working |
| **patch** | 1 | 1 | ⚠️ STUB | Basic patch struct only |
| **permission** | 4 | 1+1t | 🚨 STUB | Struct only, pattern matching not implemented |
| **platform** | 0 (new) | 1+1t | ✅ REAL | Platform detection working |
| **plugin** | 10 | 1+1t | 🚨 STUB | Plugin interface only, no loader |
| **project** | 6 | 1+1t | 🚨 STUB | Struct only, detection not implemented |
| **provider** | 30 | 52+4t | ✅ REAL | 48+ providers, excellent coverage |
| **pty** | 6 | 1+1t | 🚨 STUB | PTY interface stub, no terminal handling |
| **question** | 2 | 1 | ⚠️ STUB | Basic struct only |
| **server** | 77 | 6+2t | 🚨 STUB | HTTP server works, routes/handlers minimal |
| **session** | 19 | 3+3t | ⚠️ PARTIAL | Manager works, store/compaction partial |
| **share** | 3 | 1 | 🚨 STUB | Struct only, no sharing protocol |
| **shell** | 1 | 3+2t | ✅ REAL | Executor, tmux integration working |
| **skill** | 2 | 1 | ⚠️ STUB | Skill registry only, discovery missing |
| **snapshot** | 1 | 1 | ⚠️ STUB | Struct only, no snapshot logic |
| **storage** | 7 | 1+1t | 🚨 STUB | DB interface only, no schema/migration |
| **sync** | 3 | 1+1t | 🚨 STUB | Struct only, no sync protocol |
| **tool** | 23 | 15+2t | ⚠️ PARTIAL | Most tools done, some missing (notably LSP tool partial) |
| **ui** | 0 (new) | 8+6t | ✅ REAL | Bubble Tea TUI functional |
| **util** | 33 | 1+1t | 🚨 STUB | Only basic ID generation |
| **v2** | 4 | 1 | 🚨 STUB | Struct only, no API |
| **worktree** | 1 | 1 | ⚠️ STUB | Basic struct only |
| **account** | 5 | 1+1t | 🚨 STUB | Struct only |
| **acp** | 3 | 1 | 🚨 STUB | ACP protocol stub |
| **controlplane** | 9 | 0 | 🚨 MISSING | Fleet control plane not started |
| **env** | 1 | 1 | ✅ REAL | Env var reading works |
| **file** | 5 | 1 | 🚨 STUB | File watcher/ops stub |
| **format** | 2 | 1 | ⚠️ STUB | Basic struct only |
| **id** | 1 | 1+1t | ✅ REAL | ID generation working |
| **cli** | 149 | 28+2t | ⚠️ PARTIAL | Most commands exist, cmd/generate/plug missing |

### 1.2 CLI Commands Status

| Command | Status | Notes |
|---------|--------|-------|
| root | ✅ | Cobra root command |
| run | ✅ | Main session runner |
| serve | ✅ | HTTP server |
| agent | ✅ | Agent mode |
| session | ✅ | Session management |
| tab | ✅ | Session tabbing |
| mcp | ✅ | MCP server/client |
| stats | ✅ | Usage statistics |
| doctor | ✅ | Health checks |
| upgrade | ✅ | Self-upgrade |
| account | ⚠️ STUB | Account management stub |
| web | ⚠️ STUB | Web interface stub |
| db | ✅ | Database commands |
| debug | ✅ | Debug commands |
| export | ✅ | Session export |
| github | ✅ | GitHub integration |
| import | ✅ | Session import |
| mcp | ✅ | MCP commands |
| models | ✅ | Model listing |
| providers | ✅ | Provider management |
| pr | ✅ | PR operations |
| attach | ✅ | Attach to session |
| **cmd** | 🚨 MISSING | Command framework |
| **generate** | 🚨 MISSING | Code generation |
| **plug** | 🚨 MISSING | Plugin system |
| **uninstall** | ✅ | Uninstall command |
| **version** | ✅ | Version info |

---

## Part 2: Invalid Completed Tasks (Need Redo)

These tasks are marked "Done" but are actually stubs or incomplete:

### 2.1 Phase 2 Tasks Marked Done But Are Stubs

| Module | Claimed | Actual | Required Work |
|--------|---------|--------|---------------|
| `internal/bus` | ✅ Done | 🚨 STUB | Implement wildcard subscriptions, session-scoped events |
| `internal/command` | ✅ Done | ⚠️ STUB | Template system, argument parsing, help generation |
| `internal/lsp` | ✅ Done | 🚨 STUB | Full LSP protocol client implementation |
| `internal/pty` | ✅ Done | 🚨 STUB | Terminal/PTY with resize, shell integration |
| `internal/storage` | ✅ Done | 🚨 STUB | SQLite schema, migration from JSON |
| `internal/sync` | ✅ Done | 🚨 STUB | Sync protocol, conflict resolution |
| `internal/project` | ✅ Done | 🚨 STUB | Project detection (git, npm, etc.) |
| `internal/permission` | ✅ Done | 🚨 STUB | Pattern matching for tool access |
| `internal/ide` | ✅ Done | 🚨 STUB | LSP integration, diagnostics |
| `internal/effect` | ✅ Done | 🚨 STUB | Effect runtime with concurrency |
| `internal/patch` | ✅ Done | ⚠️ STUB | Apply/parse patches |
| `internal/share` | ✅ Done | 🚨 STUB | Sharing protocol |
| `internal/snapshot` | ✅ Done | ⚠️ STUB | Snapshot creation/restoration |
| `internal/v2` | ✅ Done | 🚨 STUB | API v2 endpoints |
| `internal/worktree` | ✅ Done | ⚠️ STUB | Git worktree management |

### 2.2 Phase 3 Tasks Marked Done But Are Stubs

| Module | Claimed | Actual | Required Work |
|--------|---------|--------|---------------|
| `internal/account` | ✅ Done | 🚨 STUB | Account CRUD, session linking |
| `internal/acp` | ✅ Done | 🚨 STUB | ACP protocol handshake |
| `internal/control-plane` | ✅ Done | ❌ MISSING | Fleet orchestration |
| `internal/file` | ✅ Done | 🚨 STUB | File watcher, tree |
| `internal/plugin` | ✅ Done | 🚨 STUB | Plugin loader, Sandman integration |
| `internal/skill` | ✅ Done | ⚠️ STUB | Skill discovery, registry |
| `internal/env` | ✅ Done | ✅ REAL | Env var reading (actually done) |
| `internal/format` | ✅ Done | ⚠️ STUB | Code formatting hooks |
| `internal/question` | ✅ Done | ⚠️ STUB | Q&A flow |
| `internal/util` | ✅ Done | 🚨 STUB | Only 1 util file vs 33 in OpenCode |
| `internal/id` | ✅ Done | ✅ REAL | ID generation (actually done) |
| `internal/installation` | ✅ Done | ✅ REAL | Version detection (actually done) |

### 2.3 TASKLIST.md Priority Tasks Marked Done But Are Stubs

| Task | File | Claimed | Actual |
|------|------|---------|--------|
| Session History | `internal/session/compaction.go` | ✅ Done | ⚠️ PARTIAL | History struct exists, compaction logic partial |
| Hashline Tool | `internal/tool/hashline.go` | ✅ Done | ✅ REAL | Actually implemented |
| OAuth Flow | `internal/mcp/oauth.go` | ✅ Done | ✅ REAL | Actually implemented |
| WebSocket Server | `internal/server/websocket.go` | ✅ Done | 🚨 STUB | File doesn't exist in actual codebase |
| Background Tasks | `internal/agent/background.go` | ✅ Done | ⚠️ PARTIAL | Basic implementation |
| Runtime Fallback | `internal/agent/fallback.go` | ✅ Done | ✅ REAL | Model fallback chain works |
| Tmux Integration | `internal/shell/tmux.go` | ✅ Done | ✅ REAL | Actually implemented |
| Transform Hooks | `internal/hook/transform.go` | ✅ Done | ⚠️ PARTIAL | Hook exists, transform logic partial |
| Continuation Hooks | `internal/hook/continuation.go` | ✅ Done | ⚠️ PARTIAL | Hook exists, continuation logic partial |
| CLI: Cmd | `internal/cli/cmd.go` | ✅ Done | 🚨 MISSING | File doesn't exist |
| CLI: Plug | `internal/cli/plug.go` | ✅ Done | 🚨 MISSING | File doesn't exist |
| CLI: Generate | `internal/cli/generate.go` | ✅ Done | 🚨 MISSING | File doesn't exist |

---

## Part 3: Invalid Tasks (Never Started)

These tasks are marked in plans but were never started:

| Task | Reference | Status |
|------|-----------|--------|
| `internal/fleet/` module | Multiple docs | ❌ NOT STARTED - UI stub only in `internal/ui/commands.go` |
| `internal/controlplane/` | TASKLIST-PHASE3 | ❌ NOT STARTED |
| `cmd/generate` | TASKLIST.md | ❌ NOT STARTED |
| `cmd/plug` | TASKLIST.md | ❌ NOT STARTED |
| `cmd/cmd` framework | 0213-Missing-Features.md | ❌ NOT STARTED |
| Web UI | Multiple docs | ❌ NOT STARTED - Only `web.go` stub |

---

## Part 4: New Freecode Features (Not in OpenCode)

These features are planned for Freecode but don't exist in OpenCode:

### 4.1 Fleet/Clustering System

| Feature | Description | Status | Priority |
|---------|-------------|--------|----------|
| Fleet Head Mode | Single instance coordinates others | ❌ NOT STARTED | P0 |
| Fleet Agent Mode | Worker instance, receives tasks | ❌ NOT STARTED | P0 |
| Fleet Client Mode | Thin client for remote monitoring | ❌ NOT STARTED | P1 |
| Fleet Web Panel | TUI panel for fleet management | ⚠️ STUB | `internal/ui/commands.go` has `FleetStatusCmd` stub |
| Fleet TLS Auth | Secure fleet communication | 🔄 IN PROGRESS | SEC-07 tracking |
| BitTorrent Transfer | Distributed file sharing | ❌ NOT STARTED | P1 |

**Fleet Architecture:**
```
┌─────────────────────────────────────────────────────┐
│                   Fleet Head                          │
│  ┌─────────┐  ┌─────────┐  ┌─────────────────────┐   │
│  │ WebSocket│  │  Task   │  │   Instance Registry │   │
│  │ Server  │  │ Queue   │  │                     │   │
│  └─────────┘  └─────────┘  └─────────────────────┘   │
└─────────────────────────────────────────────────────┘
         ▲              ▲              ▲
    ┌────┴────┐   ┌────┴────┐   ┌────┴────┐
    │ Agent 1  │   │ Agent 2  │   │ Agent N  │
    └─────────┘   └─────────┘   └─────────┘
```

### 4.2 Other New Features

| Feature | Description | OpenCode Equivalent | Status |
|---------|-------------|---------------------|--------|
| Bubble Tea TUI | Go-native terminal UI | Solid.js Web TUI | ✅ REAL |
| Hook System | 52 lifecycle hooks | 26 hooks | ✅ REAL |
| Native Providers | 48+ model providers | LiteLLM only | ✅ REAL |
| Agent System | 11 specialized agents | 11 agents | ⚠️ PARTIAL (prompts done, execution partial) |
| Skills System | Extensible skill directory | Skills | ⚠️ PARTIAL (registry done, discovery missing) |
| Shell Integration | tmux, pty, executor | Shell tools | ✅ REAL |
| i18n System | RTL, 17 languages | i18n | ✅ REAL |

---

## Part 5: Revised Implementation Plan

### Phase 1: Foundation (Truly Complete)

| Module | Status | Verification |
|---------|--------|-------------|
| Go module setup | ✅ Done | `go.mod`, `go.sum` |
| Directory structure | ✅ Done | `internal/*` layout |
| Makefile | ✅ Done | `make build` works |
| Build on 5 platforms | ✅ Done | FreeBSD, Linux, macOS, IllumOS |
| Provider system | ✅ Done | 48+ providers |
| Hook system | ✅ Done | 52 triggers |
| Shell integration | ✅ Done | tmux, executor |
| i18n | ✅ Done | RTL, loader |
| Platform detection | ✅ Done | Cross-platform |
| Agent engine | ✅ Done | Engine, streaming |
| UI (Bubble Tea) | ✅ Done | Commands, keys, input |

### Phase 2: Core Module Parity (NEEDS REDO)

**ALL 16 TASKS ARE STUBS. Start over.**

| # | Module | True Status | Real Work Required |
|---|--------|-------------|-------------------|
| 1 | `internal/bus` | 🚨 STUB | Wildcard subscriptions, session events |
| 2 | `internal/storage` | 🚨 STUB | SQLite schema, migration |
| 3 | `internal/command` | ⚠️ STUB | Templates, argument parsing |
| 4 | `internal/pty` | 🚨 STUB | Terminal with resize |
| 5 | `internal/lsp` | ✅ REAL | Full implementation complete |
| 6 | `internal/git` | ✅ REAL | Already complete |
| 7 | `internal/sync` | 🚨 STUB | Sync protocol |
| 8 | `internal/project` | 🚨 STUB | Detection algorithms |
| 9 | `internal/permission` | 🚨 STUB | Pattern matching |
| 10 | `internal/ide` | 🚨 STUB | LSP integration |
| 11 | `internal/effect` | 🚨 STUB | Full runtime |
| 12 | `internal/patch` | ⚠️ STUB | Apply/parse logic |
| 13 | `internal/share` | 🚨 STUB | Sharing protocol |
| 14 | `internal/snapshot` | ⚠️ STUB | Snapshot logic |
| 15 | `internal/v2` | 🚨 STUB | API endpoints |
| 16 | `internal/worktree` | ⚠️ STUB | Worktree operations |

### Phase 3: Extended Module Parity (NEEDS REDO)

**ALL 12 TASKS NEED SIGNIFICANT WORK.**

| # | Module | True Status | Real Work Required |
|---|--------|-------------|-------------------|
| 1 | `internal/account` | 🚨 STUB | Account operations |
| 2 | `internal/acp` | 🚨 STUB | ACP protocol |
| 3 | `internal/controlplane` | ❌ MISSING | Fleet orchestration |
| 4 | `internal/file` | 🚨 STUB | File watcher, ops |
| 5 | `internal/plugin` | 🚨 STUB | Plugin loader |
| 6 | `internal/skill` | ⚠️ STUB | Discovery |
| 7 | `internal/env` | ✅ REAL | Complete |
| 8 | `internal/format` | ⚠️ STUB | Formatting |
| 9 | `internal/question` | ⚠️ STUB | Q&A flow |
| 10 | `internal/util` | 🚨 STUB | 33→1 files |
| 11 | `internal/id` | ✅ REAL | Complete |
| 12 | `internal/installation` | ✅ REAL | Complete |

### Phase 4: Fleet/Clustering (NEW)

| # | Feature | Description | Status |
|---|---------|-------------|--------|
| 1 | Fleet Head | WebSocket server for fleet | ❌ NOT STARTED |
| 2 | Fleet Agent | Worker instance protocol | ❌ NOT STARTED |
| 3 | Fleet Client | Thin monitoring client | ❌ NOT STARTED |
| 4 | Instance Registry | Track fleet members | ❌ NOT STARTED |
| 5 | Task Queue | Distributed task execution | ❌ NOT STARTED |
| 6 | Fleet TLS | Secure communication | 🔄 IN PROGRESS |
| 7 | BitTorrent | Distributed file transfer | ❌ NOT STARTED |
| 8 | Fleet Panel | TUI fleet management | ⚠️ STUB |

### Phase 5: CLI Completeness

| Command | Status | Fix Required |
|---------|--------|-------------|
| `cmd` framework | ❌ MISSING | Implement full framework |
| `generate` | ❌ MISSING | Code generation |
| `plug` | ❌ MISSING | Plugin system |
| `account` | ⚠️ STUB | Full implementation |
| `web` | ⚠️ STUB | Web interface |

---

## Part 6: True Progress Metrics

### Current State

| Metric | Value |
|--------|-------|
| **Builds Successfully** | ✅ Yes |
| **Tests Pass** | ⚠️ Unknown (timeout) |
| **CLI Commands (26 total)** | 21 real, 5 missing/stub |
| **Core Modules (16)** | 4 real, 12 stub |
| **Extended Modules (12)** | 3 real, 9 stub/missing |
| **New Freecode Features** | 5 real, 3 stub/missing |
| **TRUE Feature Parity** | **~35%** |

### By Category

| Category | Total | Complete | Stub | Missing |
|----------|-------|----------|------|---------|
| CLI Commands | 26 | 18 | 3 | 5 |
| Core Modules | 16 | 4 | 10 | 2 |
| Extended Modules | 12 | 3 | 7 | 2 |
| Fleet System | 8 | 0 | 1 | 7 |
| **TOTAL** | **62** | **25 (40%)** | **21 (34%)** | **16 (26%)** |

---

## Part 7: Recommended Priorities

### P0 - UNBLOCK BUILD (Today)

0. **FIX `internal/lsp/lsp.go` SYNTAX** — 7 lines need `map[string]any{}{` → `map[string]any{`
   - **Effort:** 30 minutes
   - **Blocks:** ALL builds (project doesn't compile)
   - **Reference:** [LSP-IMPLEMENTATION.md](./LSP-IMPLEMENTATION.md)

### Immediate (This Sprint)

1. **Complete `internal/bus`** — Event bus is foundational for sync, project, fleet
2. **Complete `internal/storage`** — Database needed for session persistence
3. **Fix TASKLIST.md** — Remove false "Done" markers on stubs

### Short Term (Next Month)

4. **Complete `internal/lsp`** — IDE features blocked on this
   - **Effort:** 16-24 hours total
   - **Reference:** [LSP-IMPLEMENTATION.md](./LSP-IMPLEMENTATION.md)
   - Phases: Fix syntax → Add handlers → Add diagnostics → Add server mgmt → Wire tool

5. **Complete `internal/command`** — Needed for generate/plug commands
6. **Complete `internal/pty`** — Terminal support for shell integration
7. **Start Fleet Head** — Begin fleet architecture

### Medium Term

8. **Complete remaining stubs** — All Phase 2/3 modules
9. **Implement `cmd` framework** — CLI completeness
10. **Fleet Agent + Client** — Full clustering

---

## Part 8: Change Log

| Date | Description |
|------|-------------|
| 2026-05-06 | **LSP FULLY IMPLEMENTED**: lsp.go, diagnostic.go, server.go, language.go, tool/lsp.go complete |
| 2026-05-06 | **LSP FIXED**: lsp.go compiles, bidirectional handlers done (Sisyphus) |
| 2026-05-06 | Added LSP-IMPLEMENTATION.md, updated LSP status in all docs |
| 2026-05-04 | Complete status audit - found 58% of "done" tasks are actually stubs |
| 2026-05-02 | Phase 2 TASKLIST marked all 16 as done (incorrect) |
| 2026-05-02 | Phase 3 TASKLIST marked all 12 as done (incorrect) |
| 2026-05-02 | Initial plan documents created |

---

## Appendix: OpenCode Source Reference

Based on `.discovery/` docs, OpenCode has:

- **35 core modules** in `packages/opencode/src/`
- **149 CLI command files** in `packages/opencode/src/cli/cmd/`
- **23 tool implementations** in `packages/opencode/src/tool/`
- **77 server route files** in `packages/opencode/src/server/`

Freecode's ~289 Go files vs OpenCode's ~2650 source files represents approximately **11% feature coverage** by line count, though functionality coverage is higher due to Go's conciseness and provider system.
