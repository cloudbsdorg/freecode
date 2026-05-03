# Freecode Implementation Task List - Phase 2

**Last Updated:** 2026-05-02
**Author:** Mark LaPointe <mark@cloudbsd.org>
**Build Status:** ✅ Passes
**Phase:** 2 - OpenCode Parity Gap

---

## Overview

This document tracks implementation of missing modules to achieve OpenCode parity. Phase 1 (project setup, CLI commands, hooks, skills) is complete. Phase 2 focuses on core module parity.

---

## Missing Modules (OpenCode → Freecode)

### High Priority (Blocker for Core Functionality)

| # | Module | Description | OpenCode Path | Status |
|---|--------|-------------|---------------|--------|
| 1 | `internal/bus` | Event bus pub/sub system | `packages/opencode/src/bus/` | ❌ Missing |
| 2 | `internal/command` | Command framework | `packages/opencode/src/command/` | ❌ Missing |
| 3 | `internal/lsp` | Language Server Protocol client | `packages/opencode/src/lsp/` | ❌ Missing |
| 4 | `internal/pty` | Terminal/PTY handling | `packages/opencode/src/pty/` | ❌ Missing |
| 5 | `internal/storage` | Database persistence | `packages/opencode/src/storage/` | ❌ Missing |

### Medium Priority (Feature Parity)

| # | Module | Description | OpenCode Path | Status |
|---|--------|-------------|---------------|--------|
| 6 | `internal/sync` | Session sync | `packages/opencode/src/sync/` | ❌ Missing |
| 7 | `internal/project` | Project management | `packages/opencode/src/project/` | ❌ Missing |
| 8 | `internal/git` | Git operations | `packages/opencode/src/git/` | ❌ Missing |
| 9 | `internal/permission` | Permission system | `packages/opencode/src/permission/` | ❌ Missing |
| 10 | `internal/ide` | IDE integration | `packages/opencode/src/ide/` | ❌ Missing |

### Low Priority (Advanced Features)

| # | Module | Description | OpenCode Path | Status |
|---|--------|-------------|---------------|--------|
| 11 | `internal/effect` | Effects system | `packages/opencode/src/effect/` | ❌ Missing |
| 12 | `internal/patch` | Patching | `packages/opencode/src/patch/` | ❌ Missing |
| 13 | `internal/share` | Sharing | `packages/opencode/src/share/` | ❌ Missing |
| 14 | `internal/snapshot` | Snapshots | `packages/opencode/src/snapshot/` | ❌ Missing |
| 15 | `internal/v2` | API v2 | `packages/opencode/src/v2/` | ❌ Missing |
| 16 | `internal/worktree` | Git worktree | `packages/opencode/src/worktree/` | ❌ Missing |

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
- `internal/lsp/client.go` - LSP client
- `internal/lsp/diagnostic.go` - Diagnostics
- `internal/lsp/server.go` - Server management
- `internal/lsp/handler.go` - Request/response handlers

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

## Progress Summary

| Priority | Total | Done | Remaining |
|----------|-------|------|-----------|
| HIGH | 5 | 0 | 5 |
| MEDIUM | 5 | 0 | 5 |
| LOW | 6 | 0 | 6 |
| **Total** | **16** | **0** | **16** |

---

## Change Log

| Date | Commit | Description |
|------|--------|-------------|
| 2026-05-02 | - | Created Phase 2 TASKLIST.md |
