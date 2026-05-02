# Freecode

**Unified, platform-independent AI coding assistant.**

Freecode is a Go-based CLI that combines the best of opencode with enhanced agents, hooks, fleet management, and workflow features—all as a single cohesive product.

## Features

- **11 Built-in Agents** — Sisyphus (orchestrator), Hephaestus (code gen), Oracle, Librarian, Explore, Prometheus, Metis, Momus, Atlas, Multimodal-Looker, Sisyphus-Junior
- **60+ Lifecycle Hooks** — Session, tool, transform, continuation, Ralph, skill hooks
- **8 Task Categories** — visual-engineering, ultrabrain, deep, artistry, quick, writing, and more
- **Session Tabbing** — Multiple concurrent sessions with split views
- **Fleet Management** — Head/agent/client modes for multi-instance coordination
- **Built-in MCPs** — Exa websearch, Context7 docs, Grep.app, GitHub/GitLab CLI
- **TUI with Mouse Support** — Full interactive interface with clickable elements
- **Security-first** — All services bound to localhost only, no telemetry

## Platform Support

| Platform | Status | Notes |
|----------|--------|-------|
| FreeBSD 16 | Primary | Go 1.25 in ports |
| Linux | Supported | Flatpak packaging |
| macOS | Supported | Homebrew |
| IllumOS | Supported | tarball |
| Windows | ❌ NOT supported | |

## Quick Start

```bash
# Build
go build -o freecode ./cmd/freecode

# Run
./freecode

# Or install via Homebrew (macOS)
brew install freecode
```

## Architecture

- **Go-based CLI** — Single static binary, no runtime dependencies
- **Cobra CLI framework** — Standard Go CLI patterns
- **Bubble Tea TUI** — Composable terminal UI
- **SQLite** — Embedded persistent storage
- **chi router** — Lightweight HTTP API

## Key Directories

```
cmd/freecode/          # CLI entry point
cmd/freecode-server/  # Server mode entry
internal/cli/          # Cobra commands
internal/agent/        # 11 built-in agents
internal/hook/         # 60+ lifecycle hooks
internal/session/      # Session management, tabs
internal/ui/           # Bubble Tea TUI
internal/fleet/        # Fleet head/agent/client
internal/platform/     # Platform-specific code
```

## Ports (Localhost Only)

| Service | Port |
|---------|------|
| API Server | 18792 |
| MCP Server | 18793 |
| Web UI | 18791 |
| Fleet Head | 7842 |

## Comparison to opencode

Freecode is a Go conversion of opencode with all oh-my-openagent features integrated natively:

| Feature | opencode | freecode |
|---------|----------|----------|
| Language | TypeScript | Go |
| Distribution | NPM | Static binary |
| Agents | 7 | 11 |
| Hooks | ~20 | 60+ |
| Fleet Mode | ❌ | ✅ |
| Built-in MCPs | ❌ | ✅ |

## Security

- **Localhost only** — All services bind to 127.0.0.1 and ::1
- **No telemetry** — Zero analytics or tracking
- **Permission system** — Configurable tool permissions per agent
- **YOLO mode** — Optional skip-all-confirmations (off by default)

## Documentation

See [AGENTS_START_HERE.md](AGENTS_START_HERE.md) for autonomous agent guidance, or browse the [.plan/](.plan/) directory for detailed planning documents.

## Author

Mark LaPointe <mark@cloudbsd.org>

All commits are authored by Mark LaPointe. No co-authors, no sponsorships.

## License

Unlimited license granted to the project.
