# AGENTS START HERE — Freecode

> **Purpose:** This is the primary entry point for autonomous agents working on the **Freecode** project — a Go-based AI coding assistant converted from opencode (TypeScript).

> **Platform:** Freecode is designed to run on FreeBSD 16, Linux, macOS, and IllumOS (OpenSolaris). Windows is explicitly **NOT** supported.

---

## What We're Building

A **platform-independent AI coding assistant** that provides:

- **Go-based CLI** — Single static binary, no runtime dependencies
- **OpenCode compatible** — Reads opencode configurations and creates freecode configs
- **oh-my-openagent integration** — All 11 agents, 52 hooks, 8 categories
- **Session tabbing** — Multiple concurrent sessions in TUI with split views
- **Security-first** — All services bound to localhost only (127.0.0.1, ::1)
- **MCP protocol support** — Model Context Protocol client
- **Cross-platform** — FreeBSD 16 primary, Linux Flatpak, macOS Homebrew, IllumOS tarball

The project converts the opencode TypeScript monorepo to Go while:
- Preserving ALL features from opencode
- Integrating ALL configurables from oh-my-openagent
- Adding session tabbing for multi-session workflows
- Maintaining security (localhost-only services)

---

## Document Structure

All plan documents are in the `.plan/` directory:

| # | File | What It Covers |
|---|------|----------------|
| 0.0 | [`0.0-TOC.md`](.plan/0.0-TOC.md) | Master table of contents with clickable links |
| 0.1 | [`0.1-Workflow.md`](.plan/0.1-Workflow.md) | Task claiming, completion, merge conflict handling |
| 1.0 | [`1.0-Overview.md`](.plan/1.0-Overview.md) | Executive summary, phases, architecture |
| 1.1 | [`1.1-Architecture.md`](.plan/1.1-Architecture.md) | TypeScript → Go mapping |
| 2.0 | [`2.0-Design.md`](.plan/2.0-Design.md) | Go architecture, packages, concurrency |
| 3.0 | [`3.0-Implementation-Tasks.md`](.plan/3.0-Implementation-Tasks.md) | Phase-by-phase task breakdown |
| 4.0 | [`4.0-Configuration.md`](.plan/4.0-Configuration.md) | Config schema, migration, OMO integration |
| 5.0 | [`5.0-oh-my-openagent-Integration.md`](.plan/5.0-oh-my-openagent-Integration.md) | All 11 agents, 52 hooks, 8 categories |
| 6.0 | [`6.0-Session-Tabbing.md`](.plan/6.0-Session-Tabbing.md) | TUI tabs, split view, YOLO toggle |
| 7.0 | [`7.0-Packaging.md`](.plan/7.0-Packaging.md) | FreeBSD, Linux, macOS, IllumOS packages |
| 8.0 | [`8.0-Dependencies.md`](.plan/8.0-Dependencies.md) | Homebrew dependencies, admin requirements |
| 9.0 | [`9.0-Security.md`](.plan/9.0-Security.md) | Localhost binding, permissions, service security |
| 10.0 | [`10.0-Platform-Specific.md`](.plan/10.0-Platform-Specific.md) | FreeBSD, macOS, Linux, IllumOS specifics |
| 11.0 | [`11.0-Validation.md`](.plan/11.0-Validation.md) | Task completion tracking |
| 12.0 | [`12.0-Risks.md`](.plan/12.0-Risks.md) | Risks, TODO tracker |

---

## Primary Directives

### 1. Feature Parity
- **Every opencode feature** must exist in freecode
- **Every oh-my-openagent configurable** must be represented
- **Backward compatibility** — Read opencode configs, generate freecode configs

### 2. Security First
- **Localhost only** — All services bind to 127.0.0.1 and ::1
- **No remote access** — Freecode is a local-only tool
- **Permission system** — Configurable tool permissions per agent
- **YOLO mode** — Optional skip-all-confirmations for automation

### 3. Platform Independence
- **Same code, all platforms** — Minimize platform-specific code
- **FreeBSD primary** — Designed for FreeBSD 16 first
- **Progressive enhancement** — Use platform features when available

### 4. Go Idioms
- **Standard library preferred** — Use stdlib over external deps
- **Context propagation** — Use `context.Context` everywhere
- **Error wrapping** — Use `fmt.Errorf("...: %w", err)`
- **Structured logging** — Use `log/slog`

### 5. Traceability
- **Every task must be claimed** — Update the task table before starting
- **Every task must have tests** — No task is done without tests
- **Every change must be committed** — Commit after claiming, commit after completing
- **Fix other agents' code** — If tests fail due to another agent's bugs, fix them

---

## Workflow Summary

### Picking a Task
1. Pull latest: `git pull --rebase`
2. Open [`.plan/0.0-TOC.md`](.plan/0.0-TOC.md) for the master index
3. Find a task in [`.plan/3.0-Implementation-Tasks.md`](.plan/3.0-Implementation-Tasks.md)
4. Check that all `Dependencies` are marked ✅ DONE
5. Mark your task as `in_progress` in the todo list
6. Create a branch: `git checkout -b freecode/<task-name>`
7. Commit: `git add .plan/ && git commit -m "Claim task <name>" && git push`

### Completing a Task
1. Implement the task following the plan document
2. **Run all tests** — `go test ./...`
3. **Run linting** — `golangci-lint run`
4. Mark complete in the todo list and validation report
5. Commit: `git add -A && git commit -m "Complete task <name>: <desc>" && git push`
6. Create pull request

### Handling Merge Conflicts
1. Check if your task was taken by another agent
2. If taken, abandon and pick a different task
3. If not taken, resolve the conflict, keep both changes if they affect different tasks
4. `git add <file> && git rebase --continue && git push`

> **Full details:** See [`.plan/0.1-Workflow.md`](.plan/0.1-Workflow.md)

---

## Reading Order

For a new agent, read the documents in this order:

1. **This file** (`AGENTS_START_HERE.md`) — You are here
2. **[`0.1-Workflow.md`](.plan/0.1-Workflow.md)** — How to work on tasks
3. **[`1.0-Overview.md`](.plan/1.0-Overview.md)** — The big picture
4. **[`1.1-Architecture.md`](.plan/1.1-Architecture.md)** — Current state analysis
5. **[`2.0-Design.md`](.plan/2.0-Design.md)** — Go architecture proposal
6. **[`3.0-Implementation-Tasks.md`](.plan/3.0-Implementation-Tasks.md)** — Task breakdown
7. **[`4.0-Configuration.md`](.plan/4.0-Configuration.md)** — Config system
8. **[`5.0-oh-my-openagent-Integration.md`](.plan/5.0-oh-my-openagent-Integration.md)** — OMO features

Then dive into the specific phase you're working on.

---

## Key Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Language | Go | Single static binary, cross-compilation, fast startup |
| CLI Framework | Cobra | Standard Go CLI, batteries included |
| TUI Framework | Bubble Tea | Charm library, composable, extensible |
| Database | SQLite (modernc.org) | Embedded, no server, cross-platform |
| HTTP Router | chi | Lightweight, stdlib compatible |
| Config Format | YAML | Human-readable, oh-my-openagent uses JSONC |
| Session Storage | SQLite | Persistent, queryable, concurrent |

---

## Quick Reference

### Key Directories

| Directory | Purpose |
|-----------|---------|
| `cmd/freecode/` | CLI entry point |
| `cmd/freecode-server/` | Server mode entry |
| `internal/cli/` | Cobra commands |
| `internal/config/` | Configuration loading/migration |
| `internal/agent/` | Agent engine, 11 built-in agents |
| `internal/tool/` | Tool registry and implementations |
| `internal/hook/` | Hook system (52 hooks) |
| `internal/session/` | Session management, tabs |
| `internal/ui/` | Bubble Tea TUI |
| `internal/server/` | HTTP server (localhost only) |
| `internal/platform/` | Platform-specific code |

### Build Commands

```bash
# Build CLI
go build -o freecode ./cmd/freecode

# Build server
go build -o freecode-server ./cmd/freecode-server

# Run tests
go test ./...

# Run linting
golangci-lint run

# Cross-compile for FreeBSD
GOOS=freebsd GOARCH=amd64 go build -o freecode-freebsd ./cmd/freecode
```

### Ports (All Localhost Only)

| Service | Port | Protocol |
|---------|------|----------|
| CLI Default | - | stdin/stdout |
| API Server | 18792 | TCP (127.0.0.1, ::1) |
| MCP Server | 18793 | TCP (127.0.0.1, ::1) |
| Web UI | 18791 | TCP (127.0.0.1, ::1) |

---

## Implementation Phases

| Phase | Description | Tasks | Status |
|-------|-------------|-------|--------|
| 1 | Core CLI Foundation | 20 | ⏳ PENDING |
| 2 | Configuration System | 15 | ⏳ PENDING |
| 3 | Tool Implementations | 25 | ⏳ PENDING |
| 4 | Agent Engine & Sessions | 20 | ⏳ PENDING |
| 5 | TUI & Session Tabs | 15 | ⏳ PENDING |
| 6 | oh-my-openagent Integration | 20 | ⏳ PENDING |
| 7 | Polish & Packaging | 15 | ⏳ PENDING |
| **Total** | | **130** | **0%** |

---

## oh-my-openagent Feature Map

### 11 Built-in Agents
| Agent | Mode | Purpose |
|-------|------|---------|
| Sisyphus | primary | Main orchestrator |
| Hephaestus | primary | Code generation |
| Oracle | subagent | Architecture consultation |
| Librarian | subagent | Research/library |
| Explore | subagent | Exploration |
| Prometheus | all | Planning |
| Metis | all | Plan consultation |
| Momus | all | Code review |
| Atlas | primary | Session tracking |
| Multimodal-Looker | subagent | Multimodal |
| Sisyphus-Junior | all | Simpler tasks |

### 52 Lifecycle Hooks
| Tier | Count | Examples |
|------|-------|----------|
| Session | 24 | `session.start`, `session.end`, `session.error` |
| Tool | 14 | `tool.execute.before`, `tool.execute.after` |
| Transform | 5 | `transform.message`, `transform.response` |
| Continuation | 7 | `continuation.auto`, `continuation.until_done` |
| Skill | 2 | `skill.invoked`, `skill.completed` |

### 8 Categories
`visual-engineering`, `ultrabrain`, `deep`, `artistry`, `quick`, `unspecified-low`, `unspecified-high`, `writing`

---

## Session Tabbing

Freecode supports multiple concurrent sessions organized in tabs:

- `Ctrl+T` — New tab
- `Ctrl+W` — Close tab
- `Ctrl+Tab` — Next tab
- `Ctrl+Shift+V` — Vertical split
- `Ctrl+Shift+H` — Horizontal split
- `Ctrl+Y` — Toggle YOLO mode

See [`.plan/6.0-Session-Tabbing.md`](.plan/6.0-Session-Tabbing.md) for details.

---

## YOLO Mode

YOLO mode (`yolo: true`) skips all confirmations:

```yaml
yolo:
  skipEditConfirmations: true
  skipBashConfirmations: true
  skipDeleteConfirmations: true
  skipPermissionPrompts: true
  skipToolConfirmations: true
```

Toggle in TUI with `Ctrl+Y` or via commands menu.

---

## Package Formats

| Platform | Format | Location |
|----------|--------|----------|
| FreeBSD 16 | .pkg or tarball | `packaging/freebsd/` |
| Linux | Flatpak | `packaging/linux/` |
| macOS | Homebrew | `packaging/macos/` |
| IllumOS | tarball | `packaging/illuminos/` |

See [`.plan/7.0-Packaging.md`](.plan/7.0-Packaging.md) for details.

---

## Need Help?

If you encounter issues:
1. Check the relevant plan document for guidance
2. Check the task's Notes column for known issues
3. Mark the task as `blocked` with the reason
4. Commit and push so other agents know
5. Ask for guidance

> **Remember:** The goal is to build a platform-independent AI coding assistant that rivals opencode in capability while being written in Go with better cross-platform support. Every task should bring us closer to that goal.

---

**Author:** Mark LaPointe <mark@cloudbsd.org>

**Platform:** FreeBSD 16, Linux, macOS, IllumOS
**Language:** Go
**Windows:** NOT supported

---

## Author and Contribution Policy

### Author

**Mark LaPointe** <mark@cloudbsd.org>

All commits are authored by Mark LaPointe. No trailers, no co-authors, no sponsorships.

### Commit Policy

- **Single author**: All commits are made by Mark LaPointe
- **No co-authors**: No `Co-authored-by` trailers
- **No sponsorships**: No `Sponsored-by` or similar trailers
- **No funding mentions**: No grant numbers, funding sources, or sponsorship acknowledgments

### Contribution Guidelines

- Issues and pull requests welcome but final commit authority rests with the author
- All contributions are reviewed before merge
- No obligation to accept any contribution
- Contributor retains copyright but grants unlimited license to the project
