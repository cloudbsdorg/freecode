# Freecode — Implementation Tasks

**Document ID:** IMPL-001
**Version:** 4.0
**Last Updated:** 2026-05-10
**Author:** Mark LaPointe <mark@cloudbsd.org>

---

## TODO Tracker Summary

| Phase | Focus | Tasks | Completed | Total | Progress |
|-------|-------|-------|-----------|-------|----------|
| Phase 1 | Core CLI Foundation | Project setup, CLI commands | 35 | 40 | 88% |
| Phase 2 | Configuration | Config system, migration | 8 | 12 | 67% |
| Phase 3 | Core Modules | Agent, hooks, shell, i18n | 12 | 15 | 80% |
| Phase 4 | Platform | Platform detection, IDE, LSP | 6 | 8 | 75% |
| Phase 5 | UI/TUI | TUI, components, dialogs | 18 | 25 | 72% |
| Phase 6 | Session/Storage | Session mgmt, storage | 4 | 8 | 50% |
| Phase 7 | Advanced | MCP, plugins, sync | 5 | 12 | 42% |
| Phase 8 | Module Parity | OpenCode feature parity | 2 | 20 | 10% |
| **Total** | | | **90** | **140** | **64%** |

---

## 1.0 Purpose

This document contains the detailed task breakdown for implementing freecode. Tasks are organized by phase and include dependencies.

**Note:** This document uses simplified task tables for readability. See [CloudBSD Planning Standards](https://github.com/cloudbsdorg/application_guidelines/tree/main/Planning) for full format with ID, Priority, Assigned To, Owner, Phase, Start, End, Spec, and Notes columns.

---

## Phase 1: Core CLI Foundation (Week 1-2)

### 1.1 Project Setup

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Initialize Go module | `go.mod` | ✅ Done | - |
| Create directory structure | All | ✅ Done | - |
| Setup Makefile | `Makefile` | ✅ Done | - |
| Setup goreleaser | `.goreleaser.yaml` | ✅ Done | - |

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
| Config struct | `internal/config/config.go` | ✅ Done | 1.1 |
| Config loading | `internal/config/load.go` | ✅ Done | 2.1 |
| YAML parsing | `internal/config/yaml.go` | ✅ Done | 2.1 |
| JSON parsing | `internal/config/json.go` | ✅ Done | 2.1 |
| Env var support | `internal/config/env.go` | ✅ Done | 2.2 |

### 2.2 OpenCode Config Migration

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Read opencode configs | `internal/config/opencode/read.go` | ✅ Done | 2.1 |
| Convert to freecode | `internal/config/opencode/migrate.go` | ✅ Done | 2.2.1 |
| TOML support | `internal/config/opencode/toml.go` | ⏳ Planned | 2.2.1 |
| JSONC support | `internal/config/opencode/jsonc.go` | ⏳ Planned | 2.2.1 |

### 2.3 oh-my-openagent Config Integration

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Read OMO configs | `internal/config/omo/read.go` | ✅ Done | 2.1 |
| Merge into freecode | `internal/config/omo/merge.go` | ✅ Done | 2.3.1 |

---

## Phase 3: Tool Implementations (Week 3-5)

### 3.1 Core Tools

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Tool registry | `internal/tool/registry.go` | ✅ Done | 1.2 |
| Bash tool | `internal/tool/bash.go` | ✅ Done | 3.1 |
| Read tool | `internal/tool/read.go` | ✅ Done | 3.1 |
| Write tool | `internal/tool/write.go` | ✅ Done | 3.1 |
| Edit tool | `internal/tool/edit.go` | ✅ Done | 3.1 |
| Glob tool | `internal/tool/glob.go` | ✅ Done | 3.1 |
| Grep tool | `internal/tool/grep.go` | ✅ Done | 3.1 |

### 3.2 Advanced Tools

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| WebFetch tool | `internal/tool/webfetch.go` | ✅ Done | 3.1 |
| WebSearch tool | `internal/tool/websearch.go` | ✅ Done | 3.1 |
| Task tool | `internal/tool/task.go` | ✅ Done | 3.1 |
| Skill tool | `internal/tool/skill.go` | ✅ Done | 3.1 |
| Todo tool | `internal/tool/todo.go` | ✅ Done | 3.1 |
| Question tool | `internal/tool/question.go` | ✅ Done | 3.1 |
| Plan tool | `internal/tool/plan.go` | ✅ Done | 3.1 |
| LSP tool | `internal/tool/lsp.go` | ✅ Done | 3.1 |

---

## Phase 4: Agent Engine & Session Management (Week 5-7)

**Note:** Agent prompts are defined in `prompts.go`. See [0213-Freecode-Missing-Features.md](./0213-Freecode-Missing-Features.md#61-agent-prompts-done)

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
| MCP client | `internal/mcp/client.go` | ✅ Done | 4.1 |
| MCP server handling | `internal/mcp/server.go` | ✅ Done | 6.2.1 |
| OAuth flow | `internal/mcp/oauth.go` | ❌ Missing | 6.2.1 |
| Built-in MCPs | `internal/mcp/builtin.go` | ⚠️ Stub | 6.2.1 |

### 6.3 Advanced Features

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Background tasks | `internal/agent/background.go` | ❌ Missing | 6.1 |
| Tmux integration | `internal/shell/tmux.go` | ❌ Missing | 6.1 |
| Runtime fallback | `internal/agent/fallback.go` | ❌ Missing | 6.1 |
| Hashline edit | `internal/tool/hashline.go` | ❌ Missing | 3.1 |

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
| HTTP server | `internal/server/server.go` | ✅ Done | 4.3 |
| Routes | `internal/server/routes.go` | ✅ Done | 7.1.1 |
| WebSocket | `internal/server/websocket.go` | ❌ Missing | 7.1.1 |
| Health check | `internal/server/health.go` | ✅ Done | 7.1.1 |

### 7.2 Platform-Specific

| Task | Platform | Status | Dependencies |
|------|----------|--------|--------------|
| FreeBSD support | `internal/platform/freebsd.go` | ⏳ Planned | 3.1 |
| macOS support | `internal/platform/darwin.go` | ⏳ Planned | 3.1 |
| Linux support | `internal/platform/linux.go` | ⏳ Planned | 3.1 |
| IllumOS support | `internal/platform/illuminos.go` | ⏳ Planned | 3.1 |
| **Current** | `internal/platform/platform.go` | ✅ Done | 3.1 |

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

---

## Phase 8.5: TUI Gap Analysis Tasks (2026-05-10)

*Remaining work from [0214-Freecode-TUI-Gap-Analysis.md](./0214-Freecode-TUI-Gap-Analysis.md)*

### 8.5.1 Message Rendering (CRITICAL)

| ID | Task | Priority | Status | Phase | Dependencies | Files | Spec |
|----|------|----------|--------|-------|--------------|-------|------|
| 8.5.1.1 | Add tool call rendering to MessageList | P0 | 🔄 IN PROGRESS | Phase 8.5 | - | `internal/ui/messagelist.go` | [Spec](#8501) |
| 8.5.1.2 | Add syntax highlighting for code blocks | P1 | ⬜ PENDING | Phase 8.5 | 8.5.1.1 | `internal/ui/messagepart.go` | [Spec](#8501) |
| 8.5.1.3 | Add basic diff display | P1 | ⬜ PENDING | Phase 8.5 | 8.5.1.2 | `internal/ui/messagepart.go` | [Spec](#8501) |
| 8.5.1.4 | Add thinking/collapse support | P2 | ⬜ PENDING | Phase 8.5 | 8.5.1.1 | `internal/ui/messagelist.go` | [Spec](#8501) |

### 8.5.2 Tool Display (CRITICAL)

| ID | Task | Priority | Status | Phase | Dependencies | Files | Spec |
|----|------|----------|--------|-------|--------------|-------|------|
| 8.5.2.1 | Create ToolCall component | P0 | ⬜ PENDING | Phase 8.5 | - | `internal/ui/toolcall.go` | [Spec](#8502) |
| 8.5.2.2 | Create ToolResult component | P1 | ⬜ PENDING | Phase 8.5 | 8.5.2.1 | `internal/ui/toolresult.go` | [Spec](#8502) |
| 8.5.2.3 | Create ToolError component | P1 | ⬜ PENDING | Phase 8.5 | 8.5.2.2 | `internal/ui/toolerror.go` | [Spec](#8502) |
| 8.5.2.4 | Add progress indicators | P2 | ⬜ PENDING | Phase 8.5 | 8.5.2.1 | `internal/ui/toolcall.go` | [Spec](#8502) |

### 8.5.3 Dialog System (HIGH)

| ID | Task | Priority | Status | Phase | Dependencies | Files | Spec |
|----|------|----------|--------|-------|--------------|-------|------|
| 8.5.3.1 | Add backdrop to dialogs | P1 | ⬜ PENDING | Phase 8.5 | - | `internal/ui/dialog.go` | [Spec](#8503) |
| 8.5.3.2 | Add focus management | P1 | ⬜ PENDING | Phase 8.5 | 8.5.3.1 | `internal/ui/dialog.go` | [Spec](#8503) |
| 8.5.3.3 | Add keyboard navigation | P2 | ⬜ PENDING | Phase 8.5 | 8.5.3.2 | `internal/ui/` | [Spec](#8503) |
| 8.5.3.4 | Add selection integration | P2 | ⬜ PENDING | Phase 8.5 | 8.5.3.3 | `internal/ui/` | [Spec](#8503) |

### 8.5.4 Session Review (HIGH)

| ID | Task | Priority | Status | Phase | Dependencies | Files | Spec |
|----|------|----------|--------|-------|--------------|-------|------|
| 8.5.4.1 | Create session list view | P1 | ⬜ PENDING | Phase 8.5 | - | `internal/ui/sessionlist.go` | [Spec](#8504) |
| 8.5.4.2 | Add session preview | P1 | ⬜ PENDING | Phase 8.5 | 8.5.4.1 | `internal/ui/sessionreview.go` | [Spec](#8504) |
| 8.5.4.3 | Add fork/continue actions | P2 | ⬜ PENDING | Phase 8.5 | 8.5.4.2 | `internal/ui/sessionreview.go` | [Spec](#8504) |
| 8.5.4.4 | Add retry capability | P2 | ⬜ PENDING | Phase 8.5 | 8.5.4.3 | `internal/ui/sessionreview.go` | [Spec](#8504) |

---

## Detailed Specifications

### 8.5.1 Message Rendering {#8501}

**Detailed Specification:**

- Create `messagepart.go` with multi-part message support ✅ DONE
- Add `ToolCallPart`, `ToolResultPart`, `CodeBlockPart`, `ThinkingPart` types ✅ DONE
- Integrate syntax highlighting using `charmbracelet/lipgloss` colorization ✅ DONE
- Add diff display with inline annotations ✅ DONE (renderDiffBlock)
- Support collapse/expand for thinking blocks ✅ DONE (RenderThinkingBlock)

**Acceptance Criteria:**
- [x] Messages display with proper formatting
- [x] Tool calls shown with distinct styling
- [x] Code blocks have syntax highlighting
- [x] Thinking blocks can be collapsed (RenderThinkingBlock with collapsed param)

### 8.5.2 Tool Display {#8502}

**Detailed Specification:**

- Create `toolcall.go` with visual tool state
- Create `toolresult.go` for tool output display
- Create `toolerror.go` with retry capability
- Add progress indicators for long-running tools

**Acceptance Criteria:**
- [ ] Tool calls show visual state (running, success, error)
- [ ] Error cards offer retry button
- [ ] Progress indicators during tool execution

### 8.5.3 Dialog System {#8503}

**Detailed Specification:**

- Add backdrop overlay to dialogs ✅ DONE (RenderBackdrop, RenderDialogWithBackdrop)
- Implement focus trapping within dialog ✅ DONE (focusDialog enum)
- Add vim-style keyboard navigation ✅ DONE (j/k/Enter/Esc in HelpDialog)
- Connect selection context to dialog results ⚠️ PARTIAL

**Acceptance Criteria:**
- [x] Dialogs have visible backdrop
- [x] Vim-style navigation (j/k/Enter/Esc) works
- [ ] Tab/Shift+Tab cycles focus
- [x] Escape cancels

### 8.5.4 Session Review {#8504}

**Detailed Specification:**

- Create session list view with history ✅ DONE (sessionreview.go)
- Add session preview panel ✅ DONE (renderSessionPreview)
- Implement fork session action ⚠️ PARTIAL (SessionReviewDialog exists, fork action needs wiring)
- Add retry for failed turns ⚠️ PARTIAL (UI exists, retry logic needs wiring)

**Acceptance Criteria:**
- [x] Can browse past sessions
- [x] Can preview session content
- [ ] Can fork a session (needs model integration)
- [ ] Can retry failed turns (needs model integration)
| `internal/session/` | All | Sessions |
| `internal/ui/` | All | TUI |
| `internal/platform/` | Platform | Platform-specific |

---

## Phase 8: OpenCode Module Parity (Week 12-16)

### 8.1 Event Bus

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Event bus core | `internal/bus/bus.go` | ⏳ Planned | 7.1 |
| Event definitions | `internal/bus/event.go` | ⏳ Planned | 8.1.1 |
| Global bus | `internal/bus/global.go` | ⏳ Planned | 8.1.1 |

### 8.2 Command Framework

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Command registry | `internal/command/registry.go` | ⏳ Planned | 8.1 |
| Command interface | `internal/command/command.go` | ⏳ Planned | 8.2.1 |
| Template engine | `internal/command/template.go` | ⏳ Planned | 8.2.1 |

### 8.3 LSP Client

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| LSP client core | `internal/lsp/lsp.go` | ✅ DONE | 8.1 |
| Diagnostics | `internal/lsp/diagnostic.go` | ✅ DONE | 8.3.1 |
| Server management | `internal/lsp/server.go` | ✅ DONE | 8.3.1 |
| Language detection | `internal/lsp/language.go` | ✅ DONE | 8.3.1 |
| Tool integration | `internal/tool/lsp.go` | ✅ DONE | 8.3.1 |

**Implementation Status (2026-05-06):**
- ✅ FIXED: `map[string]any{}{` → `map[string]any{` on 7 lines (syntax errors resolved)
- ✅ DONE: Bidirectional handlers using `jsonrpc2.HandlerWithError`
- ✅ DONE: stdin/stdout wrapper (`stdinStdout` struct)
- ✅ DONE: LSP types (Client, Server, textDocument, Range, Position, etc.)
- ✅ DONE: Connect, Initialize, Shutdown
- ✅ DONE: DidOpen, DidChange notifications
- ✅ DONE: Hover, Definition, References, Completion, DocumentSymbol, WorkspaceSymbol
- ✅ DONE: Diagnostic store with 150ms debouncing
- ✅ DONE: Server lifecycle management with auto-detection
- ✅ DONE: Tool integration with all LSP operations

**Reference:** `packages/opencode/src/lsp/client.ts` (697 lines) - bidirectional LSP client with push/pull diagnostics

### 8.4 PTY/Terminal

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| PTY interface | `internal/pty/pty.go` | ⏳ Planned | 8.1 |
| Terminal handling | `internal/pty/terminal.go` | ⏳ Planned | 8.4.1 |
| Input handling | `internal/pty/input.go` | ⏳ Planned | 8.4.1 |

### 8.5 Storage

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Database interface | `internal/storage/db.go` | ⏳ Planned | 7.1 |
| Schema | `internal/storage/schema.go` | ⏳ Planned | 8.5.1 |
| Migration | `internal/storage/migration.go` | ⏳ Planned | 8.5.1 |

### 8.6 Medium Priority Modules

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Sync | `internal/sync/` | ⏳ Planned | 8.5 |
| Project | `internal/project/` | ⏳ Planned | 8.3, 8.7 |
| Git | `internal/git/` | ⏳ Planned | 8.6 |
| Permission | `internal/permission/` | ⏳ Planned | 8.5 |
| IDE | `internal/ide/` | ⏳ Planned | 8.3 |

### 8.7 Low Priority Modules

| Task | File | Status | Dependencies |
|------|------|--------|--------------|
| Effect | `internal/effect/` | ⏳ Planned | 8.1 |
| Patch | `internal/patch/` | ⏳ Planned | 8.5 |
| Share | `internal/share/` | ⏳ Planned | 8.5 |
| Snapshot | `internal/snapshot/` | ⏳ Planned | 8.5 |
| V2 API | `internal/v2/` | ⏳ Planned | 8.5 |
| Worktree | `internal/worktree/` | ⏳ Planned | 8.7 |

---

## LSP Implementation Reference

**Detailed Plan:** [LSP-IMPLEMENTATION.md](./LSP-IMPLEMENTATION.md)

---

## Change Log

| Date | Version | Changes |
|------|---------|---------|
| 2026-05-12 | 5.0 | TUI test coverage added: dialog 88%, tab 82%, component 7%, syntax 5% |
| 2026-05-06 | 3.0 | Updated 8.3 LSP status to BROKEN, added implementation plan reference |
| 2026-05-01 | 1.0 | Initial task breakdown |
| 2026-05-02 | 2.0 | Phase 8 added - Module parity plan |

---

## TUI Test Coverage (2026-05-12)

| Package | Coverage | Test Files |
|---------|----------|------------|
| `internal/ui` | 12.9% | messagepart_test.go, sessionreview_test.go |
| `internal/ui/dialog` | **88.3%** | window_test.go, selection_test.go, textinput_test.go |
| `internal/ui/tab` | 82.1% | (pre-existing) |
| `internal/ui/template` | 48.4% | (pre-existing) |
| `internal/ui/syntax` | **4.9%** | styles_test.go |
| `internal/ui/component` | **6.9%** | button_test.go |

**Note:** PTY-based headless TUI testing was explored but abandoned - Bubble Tea's signal/signal handling doesn't work with raw PTYs in test environments. Components render output and state machine logic are tested via unit tests.

---

**Author:** Mark LaPointe <mark@cloudbsd.org>
**Last Updated:** 2026-05-12

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author
