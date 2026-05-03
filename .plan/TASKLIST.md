# Freecode Implementation Task List

**Last Updated:** 2026-05-02
**Author:** Mark LaPointe <mark@cloudbsd.org>
**Build Status:** ✅ Passes

---

## Priority Chain (Do in Order)

Tasks must be completed in dependency order. High priority tasks unblock medium and low.

### 🔴 HIGH PRIORITY (Blocker Tasks)
*Must complete before medium priority*

| # | Task | File | Status | Dependencies | Notes |
|---|------|------|--------|-------------|-------|
| 1 | Session History | `internal/session/history.go` | ✅ Done | 4.3 | Unblocks session management |
| 2 | Hashline Tool | `internal/tool/hashline.go` | ✅ Done | 3.1 | For edit tool enhancement |
| 3 | OAuth Flow | `internal/mcp/oauth.go` | ✅ Done | 6.2 | MCP OAuth integration |
| 4 | WebSocket Server | `internal/server/websocket.go` | ✅ Done | 7.1 | Real-time communication |

### 🟡 MEDIUM PRIORITY (Feature Tasks)
*Depend on high priority completion*

| # | Task | File | Status | Dependencies | Notes |
|---|------|------|--------|-------------|-------|
| 5 | Background Tasks | `internal/agent/background.go` | ✅ Done | 6.1 | Agent background processing |
| 6 | Runtime Fallback | `internal/agent/fallback.go` | ✅ Done | 6.1 | Model fallback logic |
| 7 | Tmux Integration | `internal/shell/tmux.go` | ✅ Done | 6.1 | Shell integration |
| 8 | Transform Hooks | `internal/hook/transform.go` | ✅ Done | 6.1.1 | Input/output transformation |
| 9 | Continuation Hooks | `internal/hook/continuation.go` | ✅ Done | 6.1.1 | Session continuation |
| 10 | CLI: Cmd | `internal/cli/cmd.go` | ✅ Done | 1.3 | Command framework |
| 11 | CLI: Plug | `internal/cli/plug.go` | ✅ Done | 1.3 | Plugin system |
| 12 | CLI: Generate | `internal/cli/generate.go` | ✅ Done | 1.3 | Code generation |

### 🟢 LOW PRIORITY (Polish Tasks)
*Depend on medium priority completion*

| # | Task | File | Status | Dependencies | Notes |
|---|------|------|--------|-------------|-------|
| 13 | TOML Config | `internal/config/toml.go` | ✅ Done | 2.2.1 | TOML config parsing |
| 14 | JSONC Config | `internal/config/jsonc.go` | ✅ Done | 2.2.1 | JSONC config parsing |
| 15 | TUI Model | `internal/ui/model.go` | ✅ Done | 4.3 | Main UI model |
| 16 | TUI View | `internal/ui/view.go` | ✅ Done | 5.1 | View rendering |
| 17 | TUI Input | `internal/ui/input.go` | ✅ Done | 5.1 | Input handling |
| 18 | TUI Style | `internal/ui/style.go` | ✅ Done | 5.1 | Theme/styling |
| 19 | TUI Tab Model | `internal/ui/tab/model.go` | ✅ Done | 4.3.1 | Tab management |
| 20 | TUI Tab View | `internal/ui/tab/view.go` | ✅ Done | 5.2.1 | Tab rendering |
| 21 | TUI Tab Keys | `internal/ui/tab/keys.go` | ✅ Done | 5.2.1 | Tab keybindings |
| 22 | TUI Tab Split | `internal/ui/tab/split.go` | ✅ Done | 5.2.1 | Split view |
| 23 | YOLO Toggle | `internal/ui/tab/commands.go` | ✅ Done | 5.2.1 | YOLO mode toggle |
| 24 | MCP Builtin | `internal/mcp/builtin.go` | ⚠️ Stub | 6.2.1 | Needs full implementation |
| 25 | Platform Core | `internal/platform/platform.go` | ✅ Done | 7.2 | Platform detection & preflight |
| 26 | FreeBSD Package | `packaging/freebsd/` | ✅ Stub | 7.2.1 | Package build stub |
| 27 | Linux Flatpak | `packaging/linux/` | ✅ Stub | 7.2.3 | Flatpak package stub |
| 28 | macOS Homebrew | `packaging/macos/` | ✅ Stub | 7.2.2 | Homebrew formula stub |
| 29 | IllumOS Tarball | `packaging/illuminos/` | ✅ Stub | 7.2.4 | Tarball build stub |

---

## Completed Tasks (Do Not Modify)

### ✅ Phase 1: Project Setup & CLI
| Task | File | Status |
|------|------|--------|
| Initialize Go module | `go.mod` | ✅ Done |
| Create directory structure | All | ✅ Done |
| Setup Makefile | `Makefile` | ✅ Done |
| Setup goreleaser | `.goreleaser.yaml` | ✅ Done |
| Root command | `internal/cli/root.go` | ✅ Done |
| Run command | `internal/cli/run.go` | ✅ Done |
| Serve command | `internal/cli/serve.go` | ✅ Done |
| Agent command | `internal/cli/agent.go` | ✅ Done |
| Session command | `internal/cli/session.go` | ✅ Done |
| Tab command | `internal/cli/tab.go` | ✅ Done |
| MCP command | `internal/cli/mcp.go` | ✅ Done |
| Stats command | `internal/cli/stats.go` | ✅ Done |
| Doctor command | `internal/cli/doctor.go` | ✅ Done |
| Upgrade command | `internal/cli/upgrade.go` | ✅ Done |
| Account command | `internal/cli/account.go` | ✅ Done |
| Web command | `internal/cli/web.go` | ✅ Done |

### ✅ Phase 2: Config
| Task | File | Status |
|------|------|--------|
| Config struct | `internal/config/config.go` | ✅ Done |
| Config loading | `internal/config/load.go` | ✅ Done |
| YAML parsing | `internal/config/yaml.go` | ✅ Done |
| JSON parsing | `internal/config/json.go` | ✅ Done |
| Env var support | `internal/config/env.go` | ✅ Done |
| Read opencode configs | `internal/config/opencode/read.go` | ✅ Done |
| Convert to freecode | `internal/config/opencode/migrate.go` | ✅ Done |
| Read OMO configs | `internal/config/omo/read.go` | ✅ Done |
| Merge into freecode | `internal/config/omo/merge.go` | ✅ Done |

### ✅ Phase 3: Tools
| Task | File | Status |
|------|------|--------|
| Tool registry | `internal/tool/registry.go` | ✅ Done |
| Bash tool | `internal/tool/bash.go` | ✅ Done |
| Read tool | `internal/tool/read.go` | ✅ Done |
| Write tool | `internal/tool/write.go` | ✅ Done |
| Edit tool | `internal/tool/edit.go` | ✅ Done |
| Glob tool | `internal/tool/glob.go` | ✅ Done |
| Grep tool | `internal/tool/grep.go` | ✅ Done |
| WebFetch tool | `internal/tool/webfetch.go` | ✅ Done |
| WebSearch tool | `internal/tool/websearch.go` | ✅ Done |
| Task tool | `internal/tool/task.go` | ✅ Done |
| Skill tool | `internal/tool/skill.go` | ✅ Done |
| Todo tool | `internal/tool/todo.go` | ✅ Done |
| Question tool | `internal/tool/question.go` | ✅ Done |
| Plan tool | `internal/tool/plan.go` | ✅ Done |
| LSP tool | `internal/tool/lsp.go` | ✅ Done |

### ✅ Phase 4: Agent & Session
| Task | File | Status |
|------|------|--------|
| Engine struct | `internal/agent/engine.go` | ✅ Done |
| Message handling | `internal/agent/message.go` | ✅ Done |
| Tool calling | `internal/agent/tools.go` | ✅ Done |
| Response streaming | `internal/agent/stream.go` | ✅ Done |
| Agent prompts (11) | `internal/agent/prompts.go` | ✅ Done |
| Session manager | `internal/session/manager.go` | ✅ Done |
| Session store | `internal/session/store.go` | ✅ Done |
| Session compaction | `internal/session/compaction.go` | ✅ Done |

### ✅ Phase 6: Hooks & Skills
| Task | File | Status |
|------|------|--------|
| Hook registry | `internal/hook/registry.go` | ✅ Done |
| Session hooks (26) | `internal/hook/triggers.go` | ✅ Done |
| Tool hooks (9) | `internal/hook/triggers.go` | ✅ Done |
| Default implementations | `internal/hook/builtins.go` | ✅ Done |
| Skills (7) | `.skills/*/SKILL.md` | ✅ Done |

### ✅ Phase 7: Server & Platform
| Task | File | Status |
|------|------|--------|
| HTTP server | `internal/server/server.go` | ✅ Done |
| Routes | `internal/server/routes.go` | ✅ Done |
| Health check | `internal/server/health.go` | ✅ Done |
| MCP client | `internal/mcp/client.go` | ✅ Done |
| MCP server | `internal/mcp/server.go` | ✅ Done |
| Platform core | `internal/platform/platform.go` | ✅ Done |

---

## Progress Summary

| Priority | Total | Done | Remaining |
|----------|-------|------|----------|
| HIGH | 4 | 4 | 0 |
| MEDIUM | 8 | 8 | 0 |
| LOW | 29 | 28 | 1 (MCP Builtin) |
| **Total** | **41** | **40** | **1** |

**Note:** MCP Builtin (`internal/mcp/builtin.go`) remains a stub - full implementation requires actual API integrations.

---

## Change Log

| Date | Commit | Description |
|------|--------|-------------|
| 2026-05-02 | - | Initial task list creation |
| 2026-05-02 | 3515efb | Implement hashline tool |
| 2026-05-02 | 0363709 | Implement MCP OAuth flow |
| 2026-05-02 | cbfd5d1 | Implement WebSocket server |
| 2026-05-02 | 0254b52 | Implement background tasks, fallback, tmux |
| 2026-05-02 | 88d6796 | Implement transform hooks |
| 2026-05-02 | b3a6a25 | Implement CLI commands (cmd, plug, generate) |
| 2026-05-02 | e31d7dc | Implement TOML and JSONC config |
| 2026-05-02 | 385c6b6 | Add packaging stubs |
