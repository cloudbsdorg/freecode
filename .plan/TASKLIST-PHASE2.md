# Freecode Implementation Task List - Phase 2

**⚠️ DEPRECATED:** See [TASKLIST.md](./TASKLIST.md) and [FREECODE-STATUS.md](./FREECODE-STATUS.md) for accurate status.

**Last Updated:** 2026-05-02
**Author:** Mark LaPointe <mark@cloudbsd.org>
**Build Status:** ✅ Passes
**Phase:** 2 - OpenCode Parity Gap

> **⚠️ DEPRECATED (2026-05-04):** This document incorrectly marked modules as "Done" when they are actually stubs. The true status is documented in FREECODE-STATUS.md. Phase 2 has ~1 truly complete module out of 16.

---

## Overview

This document tracks implementation of missing modules to achieve OpenCode parity. Phase 1 (project setup, CLI commands, hooks, skills) is complete. Phase 2 focuses on core module parity.

---

## Missing Modules (OpenCode → Freecode)

### High Priority (Blocker for Core Functionality)

| # | Module | Description | OpenCode Path | Status |
|---|--------|-------------|---------------|--------|
| 1 | `internal/bus` | Event bus pub/sub system | `packages/opencode/src/bus/` | ✅ Done |
| 2 | `internal/command` | Command framework | `packages/opencode/src/command/` | ✅ Done |
| 3 | `internal/lsp` | Language Server Protocol client | `packages/opencode/src/lsp/` | ✅ DONE |
| 4 | `internal/pty` | Terminal/PTY handling | `packages/opencode/src/pty/` | ✅ Done |
| 5 | `internal/storage` | Database persistence | `packages/opencode/src/storage/` | ✅ Done |

### Medium Priority (Feature Parity)

| # | Module | Description | OpenCode Path | Status |
|---|--------|-------------|---------------|--------|
| 6 | `internal/sync` | Session sync | `packages/opencode/src/sync/` | ✅ Done |
| 7 | `internal/project` | Project management | `packages/opencode/src/project/` | ✅ Done |
| 8 | `internal/git` | Git operations | `packages/opencode/src/git/` | ✅ Done |
| 9 | `internal/permission` | Permission system | `packages/opencode/src/permission/` | ✅ Done |
| 10 | `internal/ide` | IDE integration | `packages/opencode/src/ide/` | ✅ Done |

### Low Priority (Advanced Features)

| # | Module | Description | OpenCode Path | Status |
|---|--------|-------------|---------------|--------|
| 11 | `internal/effect` | Effects system | `packages/opencode/src/effect/` | ✅ Done |
| 12 | `internal/patch` | Patching | `packages/opencode/src/patch/` | ✅ Done |
| 13 | `internal/share` | Sharing | `packages/opencode/src/share/` | ✅ Done |
| 14 | `internal/snapshot` | Snapshots | `packages/opencode/src/snapshot/` | ✅ Done |
| 15 | `internal/v2` | API v2 | `packages/opencode/src/v2/` | ✅ Done |
| 16 | `internal/worktree` | Git worktree | `packages/opencode/src/worktree/` | ✅ Done |

---

## Task Details

### 🔴 HIGH PRIORITY

#### Task 1: Event Bus (`internal/bus`)

**Reference:** `packages/opencode/src/bus/`

**Features:**
- Pub/Sub event system
- Typed event definitions
- Wildcard subscriptions
- Global bus for cross-instance events

**Files to create:**
- `internal/bus/bus.go` - Main bus service
- `internal/bus/event.go` - Event definitions
- `internal/bus/global.go` - Global event bus

**Implementation:**
```go
type EventBus interface {
    Publish(eventType string, payload interface{}) error
    Subscribe(eventType string, handler EventHandler) error
    SubscribeAll(handler EventHandler) error
    Unsubscribe(eventType string, handler EventHandler) error
}
```

---

#### Task 2: Command Framework (`internal/command`)

**Reference:** `packages/opencode/src/command/`

**Features:**
- Command registration
- Template system
- Argument parsing
- Command help generation

**Files to create:**
- `internal/command/command.go` - Command interface
- `internal/command/registry.go` - Command registry
- `internal/command/template.go` - Template engine

---

#### Task 3: LSP Client (`internal/lsp`)

**Reference:** `packages/opencode/src/lsp/`

**Features:**
- LSP protocol implementation
- Diagnostic handling
- Language server launch
- Server management

**Files to create:**
- `internal/lsp/lsp.go` - **EXISTS BUT BROKEN** (7 syntax errors)
- `internal/lsp/diagnostic.go` - Diagnostics (MISSING)
- `internal/lsp/server.go` - Server management (MISSING)
- `internal/lsp/handler.go` - Request/response handlers (MISSING)

**Current Status:** 🚨 BROKEN
- File `internal/lsp/lsp.go` exists (511 lines)
- Has syntax error: `map[string]any{}{` should be `map[string]any{` on 7 lines
- Missing bidirectional handlers (onNotification/onRequest)
- Missing diagnostics push/pull
- Reference TS implementation: `packages/opencode/src/lsp/client.ts` (697 lines)

---

#### Task 4: PTY/Terminal (`internal/pty`)

**Reference:** `packages/opencode/src/pty/`

**Features:**
- PTY creation and management
- Terminal input/output
- Window resize handling
- Shell integration

**Files to create:**
- `internal/pty/pty.go` - PTY interface
- `internal/pty/terminal.go` - Terminal handling
- `internal/pty/input.go` - Input handling

---

#### Task 5: Storage (`internal/storage`)

**Reference:** `packages/opencode/src/storage/`

**Features:**
- Database schema
- SQLite/BadgerDB integration
- JSON migration
- Data persistence

**Files to create:**
- `internal/storage/db.go` - Database interface
- `internal/storage/schema.go` - Schema definitions
- `internal/storage/migration.go` - Data migration

---

### 🟡 MEDIUM PRIORITY

#### Task 6: Sync (`internal/sync`)

Session synchronization for multi-instance coordination.

#### Task 7: Project (`internal/project`)

Project detection and management.

#### Task 8: Git (`internal/git`)

Git operations wrapper.

#### Task 9: Permission (`internal/permission`)

Permission checking system.

#### Task 10: IDE (`internal/ide`)

IDE integration hooks.

---

### 🟢 LOW PRIORITY

#### Tasks 11-16: Advanced Modules

Effect, patch, share, snapshot, v2, worktree implementations.

---

## LSP Implementation Plan

**Reference:** [LSP-IMPLEMENTATION.md](./LSP-IMPLEMENTATION.md)

The LSP module is currently **BROKEN** and needs comprehensive implementation.

**Quick Fix (30 min):** Change `map[string]any{}{` → `map[string]any{` on lines 193, 214, 232, 256, 280, 300, 327

**Full Implementation (16-24 hours):**
1. Fix syntax → Restore build
2. Add bidirectional handlers → Handle server notifications
3. Add diagnostics → Push/pull with debouncing
4. Add server management → Auto-detect, start/stop servers
5. Wire to tool → LSP tools for agents

---

## Progress Summary

| Priority | Total | Truly Complete | Broken | Remaining |
|----------|-------|-----------------|--------|-----------|
| HIGH | 5 | 0 | 1 (lsp) | 4 (all stubs) |
| MEDIUM | 5 | 1 (git) | 0 | 4 |
| LOW | 6 | 0 | 0 | 6 (all stubs) |
| **Total** | **16** | **1** | **1** | **14** |

> True completion: ~6% (but 1 module is BROKEN, not just stub)

---

## Change Log

| Date | Commit | Description |
|------|--------|-------------|
| 2026-05-02 | - | Created Phase 2 TASKLIST.md |
