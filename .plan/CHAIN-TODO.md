# Phase 2 Implementation Chain

**Last Updated:** 2026-05-02
**Author:** Mark LaPointe <mark@cloudbsd.org>
**Purpose:** Ordered implementation sequence for Phase 2 modules

---

## Chain of Dependencies

```
bus ──────────────┐
  └──────────────┼──────────────┐
storage ─────────┤              │
  └──────────────┼──────────────┤
command ─────────┤              │
  └──────────────┼──────────────┤
pty ─────────────┤              │
  └──────────────┼──┐           │
lsp ─────────────┤  │           │
  └──────────────┤  │           │
git ─────────────┘  │           │
  └───────────────┼──┘           │
sync ─────────────┤              │
  └──────────────┼──────────────┤
project ─────────┤              │
  └──────────────┼──────────────┤
permission ──────┤              │
  └──────────────┼──────────────┤
ide ─────────────┘              │
                               │
effect ─────────────────────────┤
  └────────────────────────────┼──────────────┐
patch ─────────────────────────┤              │
  └────────────────────────────┤              │
share ─────────────────────────┤              │
  └────────────────────────────┤              │
snapshot ──────────────────────┤              │
  └────────────────────────────┤              │
v2 ───────────────────────────┤              │
  └────────────────────────────┤              │
worktree ──────────────────────┘              │
```

---

## Implementation Order

### Phase A: Foundation (No dependencies on each other)

| Step | Module | File | Why First |
|------|--------|------|-----------|
| 1 | `bus` | `internal/bus/` | Event system needed by many others |
| 2 | `storage` | `internal/storage/` | Persistence needed for sync |

### Phase B: Core Framework

| Step | Module | File | Dependencies |
|------|--------|------|-------------|
| 3 | `command` | `internal/command/` | bus (optional) |
| 4 | `pty` | `internal/pty/` | shell integration |
| 5 | `lsp` | `internal/lsp/` | IDE features |

### Phase C: Integration

| Step | Module | File | Dependencies |
|------|--------|------|-------------|
| 6 | `git` | `internal/git/` | pty (for terminal git) |
| 7 | `sync` | `internal/sync/` | storage |
| 8 | `project` | `internal/project/` | git, bus |
| 9 | `permission` | `internal/permission/` | storage |
| 10 | `ide` | `internal/ide/` | lsp, bus |

### Phase D: Advanced Features

| Step | Module | File | Dependencies |
|------|--------|------|-------------|
| 11 | `effect` | `internal/effect/` | bus |
| 12 | `patch` | `internal/patch/` | storage |
| 13 | `share` | `internal/share/` | storage |
| 14 | `snapshot` | `internal/snapshot/` | storage |
| 15 | `v2` | `internal/v2/` | storage |
| 16 | `worktree` | `internal/worktree/` | git |

---

## Implementation Checklist

### Phase A - Foundation

- [ ] 1. `bus` - Event bus system
  - [ ] `internal/bus/bus.go`
  - [ ] `internal/bus/event.go`
  - [ ] `internal/bus/global.go`
  - [ ] Tests pass

- [ ] 2. `storage` - Database persistence
  - [ ] `internal/storage/db.go`
  - [ ] `internal/storage/schema.go`
  - [ ] `internal/storage/migration.go`
  - [ ] Tests pass

### Phase B - Core Framework

- [ ] 3. `command` - Command framework
  - [ ] `internal/command/registry.go`
  - [ ] `internal/command/command.go`
  - [ ] `internal/command/template.go`
  - [ ] Tests pass

- [ ] 4. `pty` - Terminal/PTY
  - [ ] `internal/pty/pty.go`
  - [ ] `internal/pty/terminal.go`
  - [ ] `internal/pty/input.go`
  - [ ] Tests pass

- [ ] 5. `lsp` - Language Server Protocol
  - [ ] `internal/lsp/client.go`
  - [ ] `internal/lsp/diagnostic.go`
  - [ ] `internal/lsp/server.go`
  - [ ] Tests pass

### Phase C - Integration

- [ ] 6. `git` - Git operations
  - [ ] `internal/git/git.go`
  - [ ] `internal/git/status.go`
  - [ ] `internal/git/commit.go`
  - [ ] Tests pass

- [ ] 7. `sync` - Session sync
  - [ ] `internal/sync/sync.go`
  - [ ] `internal/sync/transport.go`
  - [ ] Tests pass

- [ ] 8. `project` - Project management
  - [ ] `internal/project/project.go`
  - [ ] `internal/project/detect.go`
  - [ ] Tests pass

- [ ] 9. `permission` - Permission system
  - [ ] `internal/permission/permission.go`
  - [ ] `internal/permission/check.go`
  - [ ] Tests pass

- [ ] 10. `ide` - IDE integration
  - [ ] `internal/ide/ide.go`
  - [ ] `internal/ide/handler.go`
  - [ ] Tests pass

### Phase D - Advanced Features

- [ ] 11. `effect` - Effects system
  - [ ] `internal/effect/effect.go`
  - [ ] Tests pass

- [ ] 12. `patch` - Patching
  - [ ] `internal/patch/patch.go`
  - [ ] Tests pass

- [ ] 13. `share` - Sharing
  - [ ] `internal/share/share.go`
  - [ ] Tests pass

- [ ] 14. `snapshot` - Snapshots
  - [ ] `internal/snapshot/snapshot.go`
  - [ ] Tests pass

- [ ] 15. `v2` - API v2
  - [ ] `internal/v2/api.go`
  - [ ] Tests pass

- [ ] 16. `worktree` - Git worktree
  - [ ] `internal/worktree/worktree.go`
  - [ ] Tests pass

---

## Progress

| Phase | Tasks | Truly Complete | Notes |
|-------|-------|---------------|-------|
| A - Foundation | 2 | 0 | bus (stub), storage (stub) |
| B - Core | 3 | 0 | command (partial), pty (stub), lsp (stub) |
| C - Integration | 5 | 1 | git (real), sync (stub), project (stub), permission (stub), ide (stub) |
| D - Advanced | 6 | 0 | All stubs |
| **Total** | **16** | **1** | ~6% true completion |

> ⚠️ **NOTE (2026-05-04):** See [FREECODE-STATUS.md](./FREECODE-STATUS.md) for true status audit. Phase 2/3 task lists incorrectly marked modules as "Done" when they are stubs.

---

## Change Log

| Date | Description |
|------|-------------|
| 2026-05-02 | Initial chain creation |
