# AGENTS START HERE — Freecode

> **Purpose:** This is the primary entry point for autonomous agents working on the **Freecode** project.

> **What is Freecode?** A unified, platform-independent AI coding assistant that combines the best of opencode with enhanced agents, hooks, fleet management, and workflow features—all as a single cohesive product.

> **Platform:** Freecode runs on FreeBSD 16, Linux, macOS, and IllumOS (OpenSolaris). Windows is explicitly **NOT** supported.

---

## What We're Building

A **platform-independent AI coding assistant** that provides:

- **Go-based CLI** — Single static binary, no runtime dependencies
- **Migration from opencode** — Reads opencode configs, creates freecode configs
- **11 Built-in Agents** — Sisyphus, Hephaestus, Oracle, Librarian, Explore, Prometheus, Metis, Momus, Atlas, Multimodal-Looker, Sisyphus-Junior
- **60+ Lifecycle Hooks** — Session, tool, transform, continuation, skill, and Ralph hooks
- **8 Task Categories** — visual-engineering, ultrabrain, deep, artistry, quick, writing, etc.
- **Session tabbing** — Multiple concurrent sessions in TUI with split views
- **Fleet Management** — Head/agent/client modes for multi-instance coordination
- **Built-in MCPs** — Exa websearch, Context7 docs, Grep.app
- **Security-first** — All services bound to localhost only (127.0.0.1, ::1)
- **NO TELEMETRY** — Zero analytics or tracking

This is a conversion of opencode (TypeScript) to Go, with all features built directly into freecode—not as plugins.

---

## Document Structure

All plan documents are in the `.plan/` directory following CloudBSD 4-digit numbering:

### Meta (0000-0002)
| File | What It Covers |
|------|----------------|
| `0000-Freecode-TOC.md` | Master table of contents |
| `0001-Freecode-Workflow.md` | Task claiming, completion, merge handling |
| `0002-Freecode-Build-Status.md` | CI/CD pipeline, build artifacts |

### Security (0100-0106)
| File | What It Covers |
|------|----------------|
| `0100-Freecode-Security-Overview.md` | Security strategy |
| `0101-Freecode-Security-ThreatModel.md` | Threat analysis |
| `0102-Freecode-Security-AccessControl.md` | Permissions, credentials |
| `0103-Freecode-Security-Runtime.md` | Sandbox, filesystem |
| `0104-Freecode-Security-Implementation.md` | Security tasks |
| `0105-Freecode-Security-Audit.md` | oh-my-openagent audit |
| `0106-Freecode-Security-Additional.md` | Audit logging, supply chain |

### Overview & Architecture (0200-0212)
| File | What It Covers |
|------|----------------|
| `0200-Freecode-Overview.md` | Executive summary, phases |
| `0201-Freecode-Current-Architecture.md` | TypeScript → Go mapping |
| `0202-Freecode-Platform-Specific.md` | Platform-specific details |
| `0203-Freecode-Feature-Inventory.md` | 88 features tracked |
| `0204-Freecode-Features.md` | All features |
| `0210-Freecode-Architecture-Design.md` | Go architecture |
| `0211-Freecode-LiteLLM-Integration.md` | LiteLLM provider consolidation |
| `0212-Freecode-TUI-Analysis.md` | OpenCode vs Freecode TUI parity |

### Implementation (0300-0301)
| File | What It Covers |
|------|----------------|
| `0300-Freecode-Implementation-Tasks.md` | Phase-by-phase tasks |
| `0301-Freecode-Session-Tabbing.md` | TUI tabs, split view |

### Testing (0400-0403)
| File | What It Covers |
|------|----------------|
| `0400-Freecode-Testing.md` | Test strategy |
| `0401-Freecode-Unit-Tests.md` | Unit testing plan |
| `0402-Freecode-Integration-Tests.md` | Integration testing plan |
| `0403-Freecode-Code-Validation.md` | Linting, fuzzing, security |

### Operations (0500-0504)
| File | What It Covers |
|------|----------------|
| `0501-Freecode-Configuration.md` | Config schema, migration |
| `0502-Freecode-Packaging.md` | FreeBSD, Linux, macOS, IllumOS |
| `0503-Freecode-Dependencies.md` | Build dependencies |
| `0504-Freecode-I18N.md` | Internationalization |
| `0510-Freecode-Tooling.md` | Development guide |

### Risks (0700)
| File | What It Covers |
|------|----------------|
| `0700-Freecode-Risks.md` | Risk register |

### Validation (0900)
| File | What It Covers |
|------|----------------|
| `0900-Freecode-Validation.md` | Task completion, validation |

---

## Primary Directives

### 1. Feature Parity
- **Every opencode/oh-my-openagent feature** must exist in freecode
- **All freecode features are native** — Not plugins, not integrations
- **Backward compatibility** — Read opencode configs, generate freecode configs

### 2. Security First
- **Localhost only** — All services bind to 127.0.0.1 and ::1
- **Fleet can be exposed** — Explicit opt-in with TLS + auth
- **Permission system** — Configurable tool permissions per agent
- **YOLO mode** — Optional skip-all-confirmations for automation
- **NO TELEMETRY** — Zero analytics, tracking, or third-party data collection

### 3. Platform Independence
- **Same code, all platforms** — Minimize platform-specific code
- **FreeBSD primary** — Designed for FreeBSD 16 first (Go 1.25 in ports)
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
2. Open [`.plan/000.0-TOC.md`](.plan/000.0-TOC.md) for the master index
3. Find a task in [`.plan/009.0-Feature-Inventory.md`](.plan/009.0-Feature-Inventory.md)
4. Check that all dependencies are marked ✅ DONE
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
1. Check if your task was cited by another agent
2. If cited, abandon and pick a different task
3. If not cited, resolve the conflict, keep both changes if they affect different tasks
4. `git add <file> && git rebase --continue && git push`

> **Full details:** See [`.plan/001.0-Workflow.md`](.plan/001.0-Workflow.md)

---

## Reading Order

For a new agent, read the documents in this order:

1. **This file** (`AGENTS_START_HERE.md`) — You are here
2. **[`001.0-Workflow.md`](.plan/001.0-Workflow.md)** — How to work on tasks
3. **[`002.0-Overview.md`](.plan/002.0-Overview.md)** — The big picture
4. **[`003.0-Architecture.md`](.plan/003.0-Architecture.md)** — Current state analysis
5. **[`009.0-Feature-Inventory.md`](.plan/009.0-Feature-Inventory.md)** — 88 features with task IDs
6. **[`010.0-Design.md`](.plan/010.0-Design.md)** — Go architecture proposal
7. **[`013.0-Features.md`](.plan/013.0-Features.md)** — All features documented

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
| `internal/hook/` | Hook system (60+ hooks) |
| `internal/session/` | Session management, tabs |
| `internal/ui/` | Bubble Tea TUI |
| `internal/server/` | HTTP server (localhost only) |
| `internal/fleet/` | Fleet head/agent/client |
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

# Cross-compile for FreeBSD (Go 1.25 compatible)
GOOS=freebsd GOARCH=amd64 go build -o freecode-freebsd ./cmd/freecode
```

### Installed Tools (macOS)

```bash
go version        # go1.26.2 darwin/arm64
git --version    # git 2.50.1
gcc --version    # Apple clang 21.0.0
golangci-lint run # Linting
staticcheck      # Static analysis
goreleaser       # Cross-compilation
shellcheck       # Shell validation
```

### Ports (All Localhost Only by Default)

| Service | Port | Protocol |
|---------|------|----------|
| CLI Default | - | stdin/stdout |
| API Server | 18792 | TCP (127.0.0.1, ::1) |
| MCP Server | 18793 | TCP (127.0.0.1, ::1) |
| Web UI | 18791 | TCP (127.0.0.1, ::1) |
| **Fleet Head** | **7842** | **TCP (opt-in, LAN preferred)** |

---

## Freecode Feature Map

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

### 60+ Lifecycle Hooks
| Tier | Count | Examples |
|------|-------|----------|
| Session | 26 | `session.start`, `session.end`, `session.error` |
| Tool | 14 | `tool.execute.before`, `tool.execute.after` |
| Transform | 5 | `transform.message`, `transform.response` |
| Continuation | 10 | `continuation.auto`, `ralph_loop`, `atlas` |
| Ralph | 3 | `ralph.think`, `ralph.step_back` |
| Skill | 2 | `skill.invoked`, `skill.completed` |

### 8 Categories
`visual-engineering`, `ultrabrain`, `deep`, `artistry`, `quick`, `unspecified-low`, `unspecified-high`, `writing`

### Fleet Management
| Mode | Description |
|------|-------------|
| head | Control plane - accepts connections from agents |
| agent | Worker - connects to head, receives tasks |
| both | Hybrid - can be head for some, agent for others |
| client | Thin client - connect from laptop to view/manage |

---

## Session Tabbing & TUI

Freecode supports multiple concurrent sessions with full mouse support:

- `Ctrl+T` — New tab
- `Ctrl+W` — Close tab
- `Ctrl+Tab` — Next tab
- `Ctrl+P` — Command palette
- `Ctrl+\` — Fleet panel
- `Ctrl+Shift+V` — Vertical split
- `Ctrl+Shift+H` — Horizontal split
- `Ctrl+Y` — Toggle YOLO mode

**Mouse support:** All interactive elements are clickable. Right-click for context menus.

See [`.plan/014.0-Session-Tabbing.md`](.plan/014.0-Session-Tabbing.md) for details.

---

## Package Formats

| Platform | Format | Location |
|----------|--------|----------|
| FreeBSD 16 | .pkg or tarball | `packaging/freebsd/` |
| Linux | Flatpak | `packaging/linux/` |
| macOS | Homebrew | `packaging/macos/` |
| IllumOS | tarball | `packaging/illuminos/` |

See [`.plan/015.0-Packaging.md`](.plan/015.0-Packaging.md) for details.

---

## Need Help?

If you encounter issues:
1. Check the relevant plan document for guidance
2. Check the task's Notes column for known issues
3. Mark the task as `blocked` with the reason
4. Commit and push so other agents know
5. Ask for guidance

> **Remember:** The goal is to build freecode—a unified, platform-independent AI coding assistant with all features built in. Every task should bring us closer to that goal.

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
