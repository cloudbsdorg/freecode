# Freecode — Implementation Tasks

## 1.0 Purpose

This document contains the detailed task breakdown for implementing freecode. Tasks are organized by phase and include dependencies.

---

## Phase 1: Core CLI Foundation (Week 1-2)

### 1.1 Project Setup

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Initialize Go module | `go.mod` | ⏳ | - |
| Create directory structure | All | ⏳ | - |
| Setup Makefile | `Makefile` | ⏳ | - |
| Setup goreleaser | `.goreleaser.yaml` | ⏳ | - |

### 1.2 CLI Commands

**Status Note:** Most commands exist as stubs. See [0213-Freecode-Missing-Features.md](./0213-Freecode-Missing-Features.md) for gap analysis.

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Root command | `internal/cli/root.go` | ✅ Done | 1.1 |
| Run command | `internal/cli/run.go` | ✅ Done | 1.1 |
| Serve command | `internal/cli/serve.go` | ✅ Done | 1.1 |
| Agent command | `internal/cli/agent.go` | ✅ Done | 1.1 |
| Session command | `internal/cli/session.go` | ✅ Done | 1.1 |
| Tab command | `internal/cli/tab.go` | ✅ Done | 1.1 |
| MCP command | `internal/cli/mcp.go` | ✅ Done | 1.1 |
| Stats command | `internal/cli/stats.go` | ✅ Done | 1.1 |
| Doctor command | `internal/cli/doctor.go` | ✅ Done | 1.1 |
| Upgrade command | `internal/cli/upgrade.go` | ✅ Done | 1.1 |

### 1.3 Missing CLI Commands (New)

| Task | File | Priority | Status |
|------|------|----------|--------|
| Account command | `internal/cli/account.go` | HIGH | ✅ Done |
| Web command | `internal/cli/web.go` | MEDIUM | ✅ Done |
| Cmd command | `internal/cli/cmd.go` | MEDIUM | ⏳ Planned |
| Plug command | `internal/cli/plug.go` | LOW | ⏳ Planned |
| Generate command | `internal/cli/generate.go` | LOW | ⏳ Planned |

**Reference:** See [0212-Freecode-TUI-Analysis.md](./0212-Freecode-TUI-Analysis.md#missing-commands-accurate-as-of-2026-05-02)

### 1.3 Run Command Options

```go
// From opencode run.ts (lines 206-293)
var runCmd = &cobra.Command{
    Use:   "run [message..]",
    RunE:  run,
}
runCmd.Flags().BoolP("continue", "c", false, "Continue last session")
runCmd.Flags().StringP("session", "s", "", "Session ID")
runCmd.Flags().Bool("fork", false, "Fork session")
runCmd.Flags().Bool("share", false, "Share session")
runCmd.Flags().StringP("model", "m", "", "Model (provider/model)")
runCmd.Flags().String("agent", "", "Agent to use")
runCmd.Flags().String("format", "default", "Output format")
runCmd.Flags().StringSliceP("file", "f", nil, "Files to attach")
runCmd.Flags().String("title", "", "Session title")
runCmd.Flags().String("attach", "", "Attach to remote server")
runCmd.Flags().StringP("password", "p", "", "Auth password")
runCmd.Flags().String("dir", "", "Working directory")
runCmd.Flags().Int("port", 0, "Local server port")
runCmd.Flags().String("variant", "", "Model variant")
runCmd.Flags().Bool("thinking", false, "Show thinking blocks")
runCmd.Flags().Bool("dangerously-skip-permissions", false, "Skip permission checks")
runCmd.Flags().Bool("yolo", false, "Skip all confirmations")
```

---

## Phase 2: Configuration System (Week 2-3)

### 2.1 Config Structure

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Config struct | `internal/config/config.go` | ⏳ | 1.1 |
| Config loading | `internal/config/load.go` | ⏳ | 2.1 |
| YAML parsing | `internal/config/yaml.go` | ⏳ | 2.1 |
| JSON parsing | `internal/config/json.go` | ⏳ | 2.1 |
| Env var support | `internal/config/env.go` | ⏳ | 2.2 |

### 2.2 OpenCode Config Migration

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Read opencode configs | `internal/config/opencode/read.go` | ⏳ | 2.1 |
| Convert to freecode | `internal/config/opencode/migrate.go` | ⏳ | 2.2.1 |
| TOML support | `internal/config/opencode/toml.go` | ⏳ | 2.2.1 |
| JSONC support | `internal/config/opencode/jsonc.go` | ⏳ | 2.2.1 |

### 2.3 oh-my-openagent Config Integration

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Read OMO configs | `internal/config/omo/read.go` | ⏳ | 2.1 |
| Merge into freecode | `internal/config/omo/merge.go` | ⏳ | 2.3.1 |

---

## Phase 3: Tool Implementations (Week 3-5)

### 3.1 Core Tools

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Tool registry | `internal/tool/registry.go` | ⏳ | 1.2 |
| Bash tool | `internal/tool/bash.go` | ⏳ | 3.1 |
| Read tool | `internal/tool/read.go` | ⏳ | 3.1 |
| Write tool | `internal/tool/write.go` | ⏳ | 3.1 |
| Edit tool | `internal/tool/edit.go` | ⏳ | 3.1 |
| Glob tool | `internal/tool/glob.go` | ⏳ | 3.1 |
| Grep tool | `internal/tool/grep.go` | ⏳ | 3.1 |

### 3.2 Advanced Tools

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| WebFetch tool | `internal/tool/webfetch.go` | ⏳ | 3.1 |
| WebSearch tool | `internal/tool/websearch.go` | ⏳ | 3.1 |
| Task tool | `internal/tool/task.go` | ⏳ | 3.1 |
| Skill tool | `internal/tool/skill.go` | ⏳ | 3.1 |
| Todo tool | `internal/tool/todo.go` | ⏳ | 3.1 |
| Question tool | `internal/tool/question.go` | ⏳ | 3.1 |
| Plan tool | `internal/tool/plan.go` | ⏳ | 3.1 |
| LSP tool | `internal/tool/lsp.go` | ⏳ | 3.1 |

---

## Phase 4: Agent Engine & Session Management (Week 5-7)

**CRITICAL:** All 11 agents are currently stubs. See [0213-Freecode-Missing-Features.md](./0213-Freecode-Missing-Features.md#61-agent-stubs)

### 4.1 Agent Engine

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Engine struct | `internal/agent/engine.go` | ✅ Done | 3.1 |
| Message handling | `internal/agent/message.go` | ✅ Done | 4.1 |
| Tool calling | `internal/agent/tools.go` | ✅ Done | 4.1 |
| Response streaming | `internal/agent/stream.go` | ✅ Done | 4.1 |

### 4.2 Built-in Agents

**All agents have prompts in `prompts.go` - need execution implementation in `sisyphus.go`**

| Task | Agent | Status | Implementation File |
|------|-------|--------|---------------------|
| Sisyphus | primary | ✅ Prompts Done | `internal/agent/prompts.go` |
| Hephaestus | primary | ✅ Prompts Done | `internal/agent/prompts.go` |
| Oracle | subagent | ✅ Prompts Done | `internal/agent/prompts.go` |
| Librarian | subagent | ✅ Prompts Done | `internal/agent/prompts.go` |
| Explore | subagent | ✅ Prompts Done | `internal/agent/prompts.go` |
| Prometheus | all | ✅ Prompts Done | `internal/agent/prompts.go` |
| Metis | all | ✅ Prompts Done | `internal/agent/prompts.go` |
| Momus | all | ✅ Prompts Done | `internal/agent/prompts.go` |
| Atlas | primary | ✅ Prompts Done | `internal/agent/prompts.go` |
| Multimodal-Looker | subagent | ✅ Prompts Done | `internal/agent/prompts.go` |
| Sisyphus-Junior | all | ✅ Prompts Done | `internal/agent/prompts.go` |

**Next Step:** Implement actual agent execution in `sisyphus.go` using prompts from `prompts.go`

### 4.3 Session Management

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Session manager | `internal/session/manager.go` | ✅ Done | 2.1 |
| Session store | `internal/session/store.go` | ✅ Done | 4.3.1 |
| Session compaction | `internal/session/compaction.go` | ✅ Done | 4.3.1 |
| Message history | `internal/session/history.go` | ⏳ Planned | 4.3.1 |

---

## Phase 5: TUI & Session Tabs (Week 7-9)

### 5.1 TUI Framework

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Main model | `internal/ui/model.go` | ⏳ | 4.3 |
| View rendering | `internal/ui/view.go` | ⏳ | 5.1 |
| Input handling | `internal/ui/input.go` | ⏳ | 5.1 |
| Style/theme | `internal/ui/style.go` | ⏳ | 5.1 |

### 5.2 Session Tabs

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Tab model | `internal/ui/tab/model.go` | ⏳ | 4.3.1 |
| Tab rendering | `internal/ui/tab/view.go` | ⏳ | 5.2.1 |
| Tab keybindings | `internal/ui/tab/keys.go` | ⏳ | 5.2.1 |
| Tab commands | `internal/cli/tab.go` | ⏳ | 5.2.1 |
| Split view | `internal/ui/tab/split.go` | ⏳ | 5.2.1 |

### 5.3 YOLO Toggle

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| YOLO state | `internal/ui/model.go` | ⏳ | 5.1 |
| Toggle command | `internal/ui/tab/commands.go` | ⏳ | 5.2.1 |
| Toggle keybinding | `internal/ui/tab/keys.go` | ⏳ | 5.3.2 |
| Visual indicator | `internal/ui/tab/view.go` | ⏳ | 5.3.2 |

---

## Phase 6: oh-my-openagent Integration (Week 9-11)

### 6.1 Hook System

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Hook registry | `internal/hook/registry.go` | ✅ Done | 4.1 |
| Session hooks (26) | `internal/hook/triggers.go` | ✅ Done | 6.1.1 |
| Tool hooks (9) | `internal/hook/triggers.go` | ✅ Done | 6.1.1 |
| Default implementations | `internal/hook/builtins.go` | ✅ Done | 6.1.1 |
| Transform hooks | `internal/hook/transform.go` | ⏳ Planned | 6.1.1 |
| Continuation hooks | `internal/hook/continuation.go` | ⏳ Planned | 6.1.1 |

### 6.2 MCP Client

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| MCP client | `internal/mcp/client.go` | ⏳ | 4.1 |
| MCP server handling | `internal/mcp/server.go` | ⏳ | 6.2.1 |
| OAuth flow | `internal/mcp/oauth.go` | ⏳ | 6.2.1 |
| Built-in MCPs | `internal/mcp/builtin.go` | ⏳ | 6.2.1 |

### 6.3 Advanced Features

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Background tasks | `internal/agent/background.go` | ⏳ | 6.1 |
| Tmux integration | `internal/shell/tmux.go` | ⏳ | 6.1 |
| Runtime fallback | `internal/agent/fallback.go` | ⏳ | 6.1 |
| Hashline edit | `internal/tool/hashline.go` | ⏳ | 3.1 |

### 6.4 Skills System ✅ DONE

| Skill | File | Status |
|-------|------|--------|
| git-master | `.skills/git-master/SKILL.md` | ✅ Done |
| playwright | `.skills/playwright/SKILL.md` | ✅ Done |
| frontend-ui-ux | `.skills/frontend-ui-ux/SKILL.md` | ✅ Done |
| review-work | `.skills/review-work/SKILL.md` | ✅ Done |
| ai-slop-remover | `.skills/ai-slop-remover/SKILL.md` | ✅ Done |
| search-code | `.skills/search-code/SKILL.md` | ✅ Done |
| architect | `.skills/architect/SKILL.md` | ✅ Done |

---

## Phase 7: Polish & Packaging (Week 11-12)

### 7.1 Server Mode

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| HTTP server | `internal/server/server.go` | ⏳ | 4.3 |
| Routes | `internal/server/routes.go` | ⏳ | 7.1.1 |
| WebSocket | `internal/server/websocket.go` | ⏳ | 7.1.1 |
| Health check | `internal/server/health.go` | ⏳ | 7.1.1 |

### 7.2 Platform-Specific

| Task | Platform | Status | Dependencies |
|------|----------|--------|--------------|
| FreeBSD support | `internal/platform/freebsd.go` | ⏳ | 3.1 |
| macOS support | `internal/platform/darwin.go` | ⏳ | 3.1 |
| Linux support | `internal/platform/linux.go` | ⏳ | 3.1 |
| IllumOS support | `internal/platform/illuminos.go` | ⏳ | 3.1 |

### 7.3 Packaging

| Task | Platform | Status | Dependencies |
|------|----------|--------|--------------|
| FreeBSD pkg | `packaging/freebsd/` | ⏳ | 7.2.1 |
| Linux Flatpak | `packaging/linux/` | ⏳ | 7.2.3 |
| macOS Homebrew | `packaging/macos/` | ⏳ | 7.2.2 |
| IllumOS tarball | `packaging/illuminos/` | ⏳ | 7.2.4 |

### 7.4 Testing & Docs

| Task | Status | Dependencies |
|------|--------|--------------|
| Unit tests | ⏳ | All implementation |
| Integration tests | ⏳ | 7.1 |
| Documentation | ⏳ | All implementation |

---

## Task Dependencies Graph

```
Phase 1 ──────────────────────────────────────────────────────────────
  │
  ├─ 1.1 Project Setup ─────────┐
  │                              │
  └─ 1.2 CLI Commands ───────────┼── Phase 2 ── 2.1 Config Struct ── 2.2 OpenCode ── 2.3 OMO
                                 │       │           │           │
                                 │       └───────────┴───────────┘
Phase 3 ── 3.1 Registry ──────────────────────────────┘
  │
  └─ 3.2 Core Tools ───────────────────────────────────────────────┐
                                                                   │
Phase 4 ── 4.1 Engine ── 4.2 Built-in Agents ─────────────────────┤
  │                                                                │
  └─ 4.3 Session Manager ─────────────────────────────────────────┤
                                                                   │
Phase 5 ── 5.1 TUI Framework ── 5.2 Session Tabs ──────────────────┤
                                                                   │
Phase 6 ── 6.1 Hook System ── 6.2 MCP Client ── 6.3 Advanced ─────┤
                                                                   │
Phase 7 ── 7.1 Server ── 7.2 Platform ── 7.3 Packaging ──────────┘
```

---

## File Ownership

| Directory | Owner | Responsibility |
|-----------|-------|----------------|
| `cmd/` | All | CLI entry points |
| `internal/cli/` | All | Command handlers |
| `internal/config/` | All | Configuration |
| `internal/agent/` | All | Agent engine |
| `internal/tool/` | All | Tools |
| `internal/hook/` | All | Hooks |
| `internal/session/` | All | Sessions |
| `internal/ui/` | All | TUI |
| `internal/platform/` | Platform | Platform-specific |

---

## Change Log

| Date | Version | Changes |
|------|---------|---------|
| 2026-05-01 | 1.0 | Initial task breakdown |

---

**Author:** Mark LaPointe <mark@cloudbsd.org>
**Last Updated:** 2026-05-01

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author
