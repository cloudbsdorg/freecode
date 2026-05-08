# Freecode — Consolidated Task List

**Last Updated:** 2026-05-04
**Author:** Mark LaPointe <mark@cloudbsd.org>
**Reference:** [FREECODE-STATUS.md](./FREECODE-STATUS.md) for true status

---

## Legend

| Symbol | Meaning |
|---------|---------|
| ✅ REAL | Fully implemented and tested |
| ⚠️ PARTIAL | Partially implemented, needs completion |
| 🚨 STUB | Interface only, needs significant work |
| ❌ MISSING | Not started |
| 🔄 IN PROGRESS | Work underway |

---

## Section 0: Foundation (VERIFIED COMPLETE)

### 0.1 Project Setup

| Task | File | Status | Verification |
|------|------|--------|---------------|
| Go module | `go.mod` | ✅ REAL | Builds on all platforms |
| Directory structure | `internal/*` | ✅ REAL | All packages created |
| Makefile | `Makefile` | ✅ REAL | `make build` works |
| Cross-compile | GOOS/GOARCH | ✅ REAL | 5 platforms |
| goreleaser | `.goreleaser.yaml` | ✅ REAL | Binary builds |

### 0.2 Core Infrastructure

| Task | File | Status | Verification |
|------|------|--------|---------------|
| Provider system | `internal/provider/*.go` | ✅ REAL | 48+ providers |
| Hook system | `internal/hook/*.go` | ✅ REAL | 52 triggers |
| Shell integration | `internal/shell/*.go` | ✅ REAL | tmux, executor |
| i18n system | `internal/i18n/*.go` | ✅ REAL | RTL, 17 languages |
| Platform detection | `internal/platform/*.go` | ✅ REAL | Cross-platform |
| Agent engine | `internal/agent/*.go` | ✅ REAL | Streaming, tools |
| TUI | `internal/ui/*.go` | ✅ REAL | Bubble Tea functional |
| Auth | `internal/auth/*.go` | ✅ REAL | Credential store |

---

## Section 1: CLI Commands

### 1.1 Completed Commands

| Task | File | Status | Notes |
|------|------|--------|-------|
| Root command | `internal/cli/root.go` | ✅ REAL | Cobra root |
| Run command | `internal/cli/run.go` | ✅ REAL | Main session |
| Serve command | `internal/cli/serve.go` | ✅ REAL | HTTP server |
| Agent command | `internal/cli/agent.go` | ✅ REAL | Agent mode |
| Session command | `internal/cli/session.go` | ✅ REAL | Session mgmt |
| Tab command | `internal/cli/tab.go` | ✅ REAL | Session tabs |
| MCP command | `internal/cli/mcp.go` | ✅ REAL | MCP server |
| Stats command | `internal/cli/stats.go` | ✅ REAL | Usage stats |
| Doctor command | `internal/cli/doctor.go` | ✅ REAL | Health checks |
| Upgrade command | `internal/cli/upgrade.go` | ✅ REAL | Self-upgrade |
| DB command | `internal/cli/db.go` | ✅ REAL | DB operations |
| Debug command | `internal/cli/debug.go` | ✅ REAL | Debug tools |
| Export command | `internal/cli/export.go` | ✅ REAL | Session export |
| GitHub command | `internal/cli/github.go` | ✅ REAL | PR, issues |
| Import command | `internal/cli/import.go` | ✅ REAL | Session import |
| Models command | `internal/cli/models.go` | ✅ REAL | List models |
| Providers command | `internal/cli/providers.go` | ✅ REAL | List providers |
| PR command | `internal/cli/pr.go` | ✅ REAL | PR operations |
| Attach command | `internal/cli/attach.go` | ✅ REAL | Attach to session |
| Uninstall command | `internal/cli/uninstall.go` | ✅ REAL | Cleanup |
| Version command | `internal/cli/version.go` | ✅ REAL | Version info |

### 1.2 Stub Commands (Need Completion)

| Task | File | Status | Required Work |
|------|------|--------|---------------|
| Account command | `internal/cli/account.go` | ⚠️ PARTIAL | Full account operations |
| Web command | `internal/cli/web.go` | ⚠️ PARTIAL | Web interface |

### 1.3 Missing Commands (Not Started)

| Task | File | Status | Required Work |
|------|------|--------|---------------|
| **Cmd framework** | `internal/cli/cmd.go` | ❌ MISSING | Full command system |
| **Generate command** | `internal/cli/generate.go` | ❌ MISSING | Code generation |
| **Plug command** | `internal/cli/plug.go` | ❌ MISSING | Plugin system |

---

## Section 2: Core Modules (Phase 2 - NEEDS REDO)

> **⚠️ WARNING:** These were marked "Done" but are actually stubs. True implementation required.

### 2.1 Foundation Modules

| Task | File | Status | True Work Required |
|------|------|--------|-------------------|
| Event Bus | `internal/bus/bus.go` | 🚨 STUB | Wildcard subscriptions, session events |
| Storage | `internal/storage/*.go` | 🚨 STUB | SQLite schema, migration |
| Command | `internal/command/*.go` | ⚠️ PARTIAL | Templates, argument parsing |

### 2.2 Integration Modules

| Task | File | Status | True Work Required |
|------|------|--------|-------------------|
| PTY/Terminal | `internal/pty/*.go` | 🚨 STUB | Terminal with resize, shell |
| LSP Client | `internal/lsp/*.go` | 🚨 STUB | Full protocol implementation |
| Git | `internal/git/*.go` | ✅ REAL | Uses go-git, complete |
| Sync | `internal/sync/*.go` | 🚨 STUB | Sync protocol, conflict resolution |
| Project | `internal/project/*.go` | 🚨 STUB | Detection algorithms |
| Permission | `internal/permission/*.go` | 🚨 STUB | Pattern matching |
| IDE | `internal/ide/*.go` | 🚨 STUB | LSP integration |

### 2.3 Advanced Modules

| Task | File | Status | True Work Required |
|------|------|--------|-------------------|
| Effect | `internal/effect/*.go` | 🚨 STUB | Full effect runtime |
| Patch | `internal/patch/*.go` | ⚠️ PARTIAL | Apply/parse patches |
| Share | `internal/share/*.go` | 🚨 STUB | Sharing protocol |
| Snapshot | `internal/snapshot/*.go` | ⚠️ PARTIAL | Snapshot logic |
| V2 API | `internal/v2/*.go` | 🚨 STUB | API endpoints |
| Worktree | `internal/worktree/*.go` | ⚠️ PARTIAL | Worktree operations |

---

## Section 3: Extended Modules (Phase 3 - NEEDS REDO)

> **⚠️ WARNING:** These were marked "Done" but need significant work.

### 3.1 High Priority

| Task | File | Status | True Work Required |
|------|------|--------|-------------------|
| Account | `internal/account/*.go` | 🚨 STUB | Account CRUD |
| ACP | `internal/acp/*.go` | 🚨 STUB | ACP protocol |
| Control Plane | `internal/controlplane/*.go` | ❌ MISSING | Fleet orchestration |
| File | `internal/file/*.go` | 🚨 STUB | File watcher, ops |
| Plugin | `internal/plugin/*.go` | 🚨 STUB | Plugin loader |

### 3.2 Medium Priority

| Task | File | Status | True Work Required |
|------|------|--------|-------------------|
| Skill | `internal/skill/*.go` | ⚠️ PARTIAL | Discovery, registry |
| Env | `internal/env/*.go` | ✅ REAL | Complete |
| Format | `internal/format/*.go` | ⚠️ PARTIAL | Formatting hooks |
| Question | `internal/question/*.go` | ⚠️ PARTIAL | Q&A flow |
| Util | `internal/util/*.go` | 🚨 STUB | 33 files in OpenCode |

### 3.3 Low Priority

| Task | File | Status | True Work Required |
|------|------|--------|-------------------|
| ID | `internal/id/*.go` | ✅ REAL | Complete |
| Installation | `internal/installation/*.go` | ✅ REAL | Complete |

---

## Section 4: Fleet/Clustering (NEW - PHASE 4)

> **🆕 NEW FEATURE:** Not in OpenCode, requires full implementation.

### 4.1 Fleet Core

| Task | File | Status | Description |
|------|------|--------|-------------|
| Fleet Head | `internal/fleet/head.go` | ❌ MISSING | WebSocket server, task queue |
| Fleet Agent | `internal/fleet/agent.go` | ❌ MISSING | Worker protocol |
| Fleet Client | `internal/fleet/client.go` | ❌ MISSING | Thin monitoring client |
| Instance Registry | `internal/fleet/registry.go` | ❌ MISSING | Track fleet members |
| Task Queue | `internal/fleet/queue.go` | ❌ MISSING | Distributed execution |

### 4.2 Fleet Security

| Task | File | Status | Description |
|------|------|--------|-------------|
| Fleet TLS | `internal/fleet/tls.go` | 🔄 IN PROGRESS | SEC-07 tracking |
| API Key Auth | `internal/fleet/auth.go` | ❌ MISSING | Fleet authentication |

### 4.3 Fleet Tools

| Task | File | Status | Description |
|------|------|--------|-------------|
| Fleet Status | `internal/fleet/status.go` | ❌ MISSING | Health monitoring |
| Fleet Exec | `internal/fleet/exec.go` | ❌ MISSING | Remote command execution |
| Fleet SCP | `internal/fleet/scp.go` | ❌ MISSING | File transfer |
| Fleet SSH | `internal/fleet/ssh.go` | ❌ MISSING | Shell into instances |

### 4.4 Distributed Storage

| Task | File | Status | Description |
|------|------|--------|-------------|
| BitTorrent | `internal/fleet/bt/*.go` | ❌ MISSING | Distributed file sharing |
| Sync Protocol | `internal/sync/*.go` | 🚨 STUB | Session sync |

---

## Section 5: TUI Components

### 5.1 Completed TUI

| Task | File | Status |
|------|------|--------|
| Commands | `internal/ui/commands.go` | ✅ REAL |
| Keys | `internal/ui/keys.go` | ✅ REAL |
| Input | `internal/ui/input.go` | ✅ REAL |
| Model | `internal/ui/model.go` | ✅ REAL |

### 5.2 Fleet Panel (TUI)

| Task | File | Status | Description |
|------|------|--------|-------------|
| Fleet Status Cmd | `internal/ui/commands.go` | ⚠️ PARTIAL | `FleetStatusCmd` stub |
| Fleet Panel | `internal/ui/fleet.go` | ❌ MISSING | Full fleet management |

---

## Section 6: P3 Nice to Have (Enhancements)

### 6.1 P3 Enhancements

| Task | File | Status | Description |
|------|------|--------|-------------|
| Sound effects | `internal/ui/sound.go` | ✅ REAL | Terminal bell/beep on events |
| Prompt autocomplete | `internal/ui/autocomplete.go` | ✅ REAL | Full frecency/history support |
| Plugin Runtime | `internal/plugin/runtime.go` | ✅ REAL | Complete plugin loader |
| Timeline/fork dialogs | `internal/ui/timeline.go` | ✅ REAL | Advanced session timeline |
| Error boundary | `internal/ui/error.go` | ✅ REAL | Error recovery component |
| Diff wrap toggle | `internal/ui/diff.go` | ✅ REAL | Toggle diff line wrapping |
| Animation toggle | `internal/ui/animation.go` | ✅ REAL | Enable/disable animations |

---

## Progress Summary

### By Section

| Section | Total | ✅ Real | ⚠️ Partial | 🚨 Stub | ❌ Missing |
|---------|-------|---------|-------------|---------|-----------|
| 0. Foundation | 9 | 9 | 0 | 0 | 0 |
| 1. CLI Commands | 26 | 21 | 2 | 0 | 3 |
| 2. Core Modules | 16 | 1 | 3 | 11 | 1 |
| 3. Extended Modules | 12 | 3 | 4 | 4 | 1 |
| 4. Fleet (NEW) | 14 | 0 | 1 | 0 | 13 |
| 5. TUI Components | 6 | 4 | 1 | 0 | 1 |
| 6. P3 Nice to Have | 7 | 7 | 0 | 0 | 0 |
| **TOTAL** | **90** | **45 (50%)** | **4 (4%)** | **15 (17%)** | **26 (29%)** |

### True Completion Rate

**VERIFIED COMPLETE (working implementation):** 45 tasks (50%)
**NEEDS COMPLETION:** 45 tasks (50%)

---

## Change Log

| Date | Description |
|------|-------------|
| 2026-05-07 | Added Section 6: P3 Nice to Have with 7 enhancement tasks |
| 2026-05-04 | Complete rewrite - found 58% of "done" tasks were actually stubs |
| 2026-05-02 | Phase 2/3 TASKLISTs created with incorrect "Done" status |
| 2026-05-02 | Initial task tracking started |
