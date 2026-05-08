# Freecode — Implementation Plan Overview

## 1.0 Executive Summary

**Project:** Convert opencode (TypeScript/JavaScript monorepo) to Go language, creating "freecode"

**What is Freecode?** A unified, platform-independent AI coding assistant that combines the best of opencode with the enhancements from oh-my-openagent—all as a single cohesive product. Freecode is NOT opencode + a plugin; it is a standalone product with all features built in.

**Goals:**
- Runs on FreeBSD, Linux, macOS, and IllumOS (OpenSolaris)
- Migrates from opencode configurations seamlessly
- Provides all features from oh-my-openagent (agents, hooks, categories) built directly in
- Includes session tabbing for multiple concurrent sessions
- Binds all services to localhost only (127.0.0.1 and ::1)

**Method:** Incremental conversion preserving all features while leveraging Go idioms.

---

## 2.0 Why Go?

| Factor | TypeScript (opencode) | Go (freecode) |
|--------|----------------------|---------------|
| Binary distribution | Requires Node.js/Bun | Single static binary |
| Cross-compilation | Complex | Built-in `GOOS`/`GOARCH` |
| Performance | JIT overhead | AOT compiled |
| Memory | Higher | Lower footprint |
| Deployment | Multiple files | Single binary |
| Startup | Slow | Fast |

---

## 3.0 Key Design Principles

1. **Feature Parity:** Every opencode feature exists in freecode
2. **Config Compatibility:** Reads opencode configs, creates freecode configs
3. **Built-in Agents:** 11 agents (Sisyphus, Hephaestus, Oracle, etc.) are native to freecode
4. **Built-in Hooks:** 52 lifecycle hooks are native to freecode, not a plugin
5. **Built-in Categories:** 8 task categories native to freecode
6. **Platform Independence:** Same code, all platforms
7. **Security First:** All services on 127.0.0.1/::1 only
8. **Session Tabs:** Multiple concurrent sessions in TUI
9. **NO TELEMETRY:** Zero analytics, tracking, or third-party data collection

---

## 4.0 Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        freecode CLI                          │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Config    │  │   Agent     │  │       Tools         │ │
│  │   System    │  │   Engine    │  │  (bash, read, etc.) │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Hooks     │  │    MCP      │  │      Session        │ │
│  │   System    │  │   Client    │  │      Manager        │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Tmux      │  │ Background  │  │    Categories       │ │
│  │ Integration │  │   Tasks     │  │                     │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                    TUI (Bubble Tea)                         │
├─────────────────────────────────────────────────────────────┤
│                      SQLite Storage                          │
└─────────────────────────────────────────────────────────────┘
```

---

## 5.0 Built-in Agents (11 Total)

Freecode includes 11 built-in agents as a native feature—not a plugin:

| Agent | Mode | Default Model | Purpose |
|-------|------|---------------|---------|
| Sisyphus | primary | claude-opus-4-7 | Main orchestrator |
| Hephaestus | primary | gpt-5.4 | Code generation |
| Oracle | subagent | gpt-5.4 | Architecture advisor (read-only) |
| Librarian | subagent | gpt-5.4-mini-fast | Library research (read-only) |
| Explore | subagent | gpt-5.4-mini-fast | Codebase exploration (read-only) |
| Prometheus | all | claude-opus-4-7 | Planning |
| Metis | all | claude-opus-4-7 | Pre-planning consultation |
| Momus | all | gpt-5.4 | Code review |
| Atlas | primary | claude-sonnet-4-6 | Todo completion |
| Multimodal-Looker | subagent | gpt-5.4 | Image analysis |
| Sisyphus-Junior | all | (inherited) | Simple tasks |

---

## 6.0 Built-in Hooks (52 Total)

Lifecycle hooks are native to freecode:

| Tier | Count | Examples |
|------|-------|----------|
| Session | 24 | `session.start`, `session.end`, `session.error` |
| Tool Guard | 14 | `tool.execute.before`, `tool.execute.after` |
| Transform | 5 | `transform.message`, `transform.response` |
| Continuation | 7 | `continuation.auto`, `continuation.until_done` |
| Skill | 2 | `skill.invoked`, `skill.completed` |

---

## 7.0 Built-in Categories (8 Total)

Task categories are native to freecode:

`visual-engineering`, `ultrabrain`, `deep`, `artistry`, `quick`, `unspecified-low`, `unspecified-high`, `writing`

---

## 8.0 Session Tabbing

Sessions can be organized in tabs:
- Create, close, rename tabs
- Move session between tabs
- Split view (vertical/horizontal)
- Detach tab to new window

---

## 9.0 Platform Support

| Platform | Package Format | Status |
|---------|---------------|--------|
| FreeBSD | .pkg or tarball | Primary target |
| Linux | Flatpak | Secondary |
| macOS | Homebrew | Secondary |
| IllumOS | tarball | Experimental |

**Note:** Windows explicitly NOT supported ("dead platform")

---

## 10.0 Security Model

All network services bind to localhost only:
- `127.0.0.1` (IPv4 loopback)
- `::1` (IPv6 loopback)

No remote access. Services:
- Local API server (port 18792)
- MCP server (port 18793)
- Web UI (port 18791)

---

## 11.0 Configuration Compatibility

### 11.1 Legacy (Read-Only Migration)

On first run, freecode reads opencode configs to migrate settings:

```
~/.config/opencode/config.json
~/.config/opencode/config.toml
~/.config/opencode/opencode.json
~/.config/opencode/opencode.jsonc
~/.config/opencode/tui.json
```

### 11.2 Freecode Native

Freecode writes to its own unified config:

```
~/.config/freecode/config.yaml
~/.config/freecode/config.json
~/.config/freecode/profiles/
```

---

## 12.0 Implementation Phases

| Phase | Focus | Duration |
|-------|-------|----------|
| Phase 1 | Core CLI foundation, Go project setup | Week 1-2 |
| Phase 2 | Configuration system, config migration | Week 2-3 |
| Phase 3 | Tool implementations (bash, read, write, edit, etc.) | Week 3-5 |
| Phase 4 | Agent engine, hooks, session management | Week 5-7 |
| Phase 5 | TUI, session tabs, server mode | Week 7-9 |
| Phase 6 | MCP client, background tasks, Tmux | Week 9-11 |
| Phase 7 | Polish, testing, packaging | Week 11-12 |

---

## 13.0 TODO Tracker

| Phase | Tasks | Completed | Total | Progress |
|-------|-------|-----------|-------|----------|
| Phase 1 | Core CLI | 0 | 20 | 0% |
| Phase 2 | Config | 0 | 15 | 0% |
| Phase 3 | Tools | 0 | 25 | 0% |
| Phase 4 | Agent/Session | 0 | 20 | 0% |
| Phase 5 | TUI | 0 | 15 | 0% |
| Phase 6 | MCP/Tmux | 0 | 20 | 0% |
| Phase 7 | Polish/Package | 0 | 15 | 0% |
| **Total** | | **0** | **130** | **0%** |

---

## 14.0 Related Documents

| Document | Purpose |
|----------|---------|
| [`1.1-Architecture.md`](./1.1-Architecture.md) | Detailed architecture analysis |
| [`2.0-Design.md`](./2.0-Design.md) | Go language design |
| [`3.0-Implementation-Tasks.md`](./3.0-Implementation-Tasks.md) | Detailed task list |
| [`4.0-Configuration.md`](./4.0-Configuration.md) | Unified configuration system |
| [`5.0-Features.md`](./5.0-Features.md) | All freecode features |
| [`7.0-Packaging.md`](./7.0-Packaging.md) | Packaging details |
| [`8.0-Dependencies.md`](./8.0-Dependencies.md) | Dependency installation |
| [`9.0-Security.md`](./9.0-Security.md) | Security model |
| [`13.0-Audit-Report.md`](./13.0-oh-my-openagent-Audit.md) | Audit of oh-my-openagent |

---

**Author:** Mark LaPointe <mark@cloudbsd.org>
**Last Updated:** 2026-05-01

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author
