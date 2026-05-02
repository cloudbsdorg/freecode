# Freecode — Missing Features vs OpenCode

**Document ID:** Freecode-Missing-Features
**Version:** 1.0
**Last Updated:** 2026-05-02
**Maintainer:** Mark LaPointe <mark@cloudbsd.org>
**Status:** ACTIVE
**Classification:** INTERNAL

---

## 1.0 Purpose

This document catalogs features that exist in opencode but are missing or stubbed in freecode. It serves as a roadmap for achieving full feature parity.

---

## 2.0 Missing CLI Commands

### 2.1 High Priority

| Command | opencode | freecode | Gap | Reference |
|---------|----------|----------|-----|-----------|
| `account` | ✅ account.ts (7889 bytes) | ❌ Missing | Full command missing | [opencode account.ts](file:///Users/mlapointe/git/opencode/packages/opencode/src/cli/cmd/account.ts) |

### 2.2 Medium Priority

| Command | opencode | freecode | Gap | Reference |
|---------|----------|----------|-----|-----------|
| `cmd` | ✅ cmd.ts (183 bytes) | ❌ Missing | CLI command framework | [opencode cmd.ts](file:///Users/mlapointe/git/opencode/packages/opencode/src/cli/cmd/cmd.ts) |
| `web` | ✅ web.ts (2462 bytes) | ❌ Missing (has `serve`) | Web interface launcher | [opencode web.ts](file:///Users/mlapointe/git/opencode/packages/opencode/src/cli/cmd/web.ts) |
| `plug` | ✅ plug.ts (6952 bytes) | ⚠️ Has `plugin` | Plugin system | [opencode plug.ts](file:///Users/mlapointe/git/opencode/packages/opencode/src/cli/cmd/plug.ts) |
| `generate` | ✅ generate.ts (2983 bytes) | ❌ Missing | Code generation | [opencode generate.ts](file:///Users/mlapointe/git/opencode/packages/opencode/src/cli/cmd/generate.ts) |

### 2.3 CLI Comparison

**Opencode commands (25):**
```
account, acp, agent, cmd, db, export, generate, github, import, mcp,
models, plug, pr, providers, run, serve, session, stats, uninstall,
upgrade, web
```

**Freecode commands (26 files, 25 unique):**
```
acp, agent, attach, cli_test, db, debug, doctor, export, github,
import, mcp, models, plugin, pr, providers, root, run, serve,
session, stats, tab, uninstall, upgrade, version
```

**Missing:** `account`, `cmd`, `generate`, `plug`, `web`

**Extra in freecode:** `attach`, `cli_test`, `debug`, `doctor`, `plugin`, `tab`, `version` (but some may be stubs)

---

## 3.0 Missing UI Components

### 3.1 App UI Components

**Opencode app components (39 files):**

| Component | opencode | freecode | Status |
|-----------|----------|----------|--------|
| Debug bar | ✅ debug-bar.tsx | ❌ | Missing |
| Dialog: Connect provider | ✅ dialog-connect-provider.tsx | ❌ | Missing |
| Dialog: Custom provider | ✅ dialog-custom-provider.tsx | ❌ | Missing |
| Dialog: Edit project | ✅ dialog-edit-project.tsx | ❌ | Missing |
| Dialog: Fork | ✅ dialog-fork.tsx | ❌ | Missing |
| Dialog: Manage models | ✅ dialog-manage-models.tsx | ❌ | Missing |
| Dialog: Release notes | ✅ dialog-release-notes.tsx | ❌ | Missing |
| Dialog: Select directory | ✅ dialog-select-directory.tsx | ❌ | Missing |
| Dialog: Select file | ✅ dialog-select-file.tsx | ❌ | Missing |
| Dialog: Select MCP | ✅ dialog-select-mcp.tsx | ❌ | Missing |
| Dialog: Select model unpaid | ✅ dialog-select-model-unpaid.tsx | ❌ | Missing |
| Dialog: Select model | ✅ dialog-select-model.tsx | ❌ | Missing |
| Dialog: Select provider | ✅ dialog-select-provider.tsx | ❌ | Missing |
| Dialog: Select server | ✅ dialog-select-server.tsx | ❌ | Missing |
| Dialog: Settings | ✅ dialog-settings.tsx | ❌ | Missing |
| File tree | ✅ file-tree.tsx | ❌ | Missing |
| Model tooltip | ✅ model-tooltip.tsx | ❌ | Missing |
| Prompt input | ✅ prompt-input.tsx (54635 bytes) | ❌ | Missing |
| Session context usage | ✅ session-context-usage.tsx | ❌ | Missing |
| Settings general | ✅ settings-general.tsx | ❌ | Missing |
| Settings keybinds | ✅ settings-keybinds.tsx | ❌ | Missing |
| Settings list | ✅ settings-list.tsx | ❌ | Missing |
| Settings models | ✅ settings-models.tsx | ❌ | Missing |
| Settings providers | ✅ settings-providers.tsx | ❌ | Missing |
| Status popover body | ✅ status-popover-body.tsx | ❌ | Missing |
| Status popover | ✅ status-popover.tsx | ❌ | Missing |
| Terminal | ✅ terminal.tsx | ❌ | Missing |
| Titlebar history | ✅ titlebar-history.ts | ❌ | Missing |
| Titlebar | ✅ titlebar.tsx | ❌ | Missing |

### 3.2 Session Components

**Opencode session components (13 files):**

| Component | opencode | freecode | Status |
|-----------|----------|----------|--------|
| Session context breakdown | ✅ session-context-breakdown.tsx | ❌ | Missing |
| Session context format | ✅ session-context-format.ts | ❌ | Missing |
| Session context metrics | ✅ session-context-metrics.tsx | ❌ | Missing |
| Session context tab | ✅ session-context-tab.tsx | ❌ | Missing |
| Session header | ✅ session-header.tsx | ❌ | Missing |
| Session new view | ✅ session-new-view.tsx | ❌ | Missing |
| Session sortable tab | ✅ session-sortable-tab.tsx | ❌ | Missing |
| Session sortable terminal tab | ✅ session-sortable-terminal-tab.tsx | ❌ | Missing |

### 3.3 Prompt Input Components

**Opencode prompt-input components (20 files):**

| Component | opencode | freecode | Status |
|-----------|----------|----------|--------|
| Attachments | ✅ attachments.ts | ❌ | Missing |
| Build request parts | ✅ build-request-parts.ts | ❌ | Missing |
| Context items | ✅ context-items.tsx | ❌ | Missing |
| Drag overlay | ✅ drag-overlay.tsx | ❌ | Missing |
| Editor DOM | ✅ editor-dom.ts | ❌ | Missing |
| Files | ✅ files.ts | ❌ | Missing |
| History | ✅ history.ts | ❌ | Missing |
| Image attachments | ✅ image-attachments.tsx | ❌ | Missing |
| Paste | ✅ paste.ts | ❌ | Missing |
| Placeholder | ✅ placeholder.ts | ❌ | Missing |
| Slash popover | ✅ slash-popover.tsx | ❌ | Missing |
| Submit | ✅ submit.ts | ❌ | Missing |

### 3.4 App Hooks

| Hook | opencode | freecode | Status |
|------|----------|----------|--------|
| use-providers | ✅ use-providers.ts | ❌ | Missing |

### 3.5 App Context Providers

Opencode has 20+ context providers in `packages/app/src/context/`:

| Provider | opencode | freecode | Status |
|---------|----------|----------|--------|
| File context | ✅ | ❌ | Missing |
| Global sync | ✅ | ❌ | Missing |

---

## 4.0 Missing TUI Components

### 4.1 TUI Components

Opencode TUI components in `packages/opencode/src/cli/cmd/tui/component/`:

| Component | opencode | freecode | Status |
|-----------|----------|----------|--------|
| Dialog: Agent | ✅ dialog-agent.tsx | ❌ | Missing |
| Dialog: Command | ✅ dialog-command.tsx | ❌ | Missing |
| Dialog: Console org | ✅ dialog-console-org.tsx | ❌ | Missing |
| Dialog: Go upsell | ✅ dialog-go-upsell.tsx | ❌ | Missing |
| Dialog: MCP | ✅ dialog-mcp.tsx | ❌ | Missing |
| Dialog: Model | ✅ dialog-model.tsx | ❌ | Missing |
| Dialog: Provider | ✅ dialog-provider.tsx | ❌ | Missing |
| Dialog: Session delete failed | ✅ dialog-session-delete-failed.tsx | ❌ | Missing |
| Dialog: Session list | ✅ dialog-session-list.tsx | ❌ | Missing |
| Dialog: Session rename | ✅ dialog-session-rename.tsx | ❌ | Missing |
| Dialog: Skill | ✅ dialog-skill.tsx | ❌ | Missing |
| Dialog: Stash | ✅ dialog-stash.tsx | ❌ | Missing |
| Dialog: Status | ✅ dialog-status.tsx | ❌ | Missing |
| Dialog: Tag | ✅ dialog-tag.tsx | ❌ | Missing |
| Dialog: Theme list | ✅ dialog-theme-list.tsx | ❌ | Missing |
| Dialog: Variant | ✅ dialog-variant.tsx | ❌ | Missing |
| Dialog: Workspace create | ✅ dialog-workspace-create.tsx | ❌ | Missing |
| Dialog: Workspace unavailable | ✅ dialog-workspace-unavailable.tsx | ❌ | Missing |
| Logo | ✅ logo.tsx (27759 bytes) | ❌ | Missing |
| Prompt component | ✅ prompt/ | ❌ | Missing |

### 4.2 TUI UI Components

| Component | opencode | freecode | Status |
|-----------|----------|----------|--------|
| Dialog alert | ✅ dialog-alert.tsx | ❌ | Missing |
| Dialog confirm | ✅ dialog-confirm.tsx | ❌ | Missing |
| Dialog export options | ✅ dialog-export-options.tsx | ❌ | Missing |
| Dialog help | ✅ dialog-help.tsx | ❌ | Missing |
| Dialog prompt | ✅ dialog-prompt.tsx | ❌ | Missing |
| Dialog select | ✅ dialog-select.tsx | ❌ | Missing |
| Dialog | ✅ dialog.tsx | ❌ | Missing |
| Link | ✅ link.tsx | ❌ | Missing |
| Spinner | ✅ spinner.ts | ❌ | Missing |
| Toast | ✅ toast.tsx | ❌ | Missing |

### 4.3 TUI Context

Opencode has extensive context management in `packages/opencode/src/cli/cmd/tui/context/`:

| Context | opencode | freecode | Status |
|---------|----------|----------|--------|
| Args | ✅ args.tsx | ❌ | Missing |
| Directory | ✅ directory.ts | ❌ | Missing |
| Editor (zed) | ✅ editor-zed.ts | ❌ | Missing |
| Editor | ✅ editor.ts | ❌ | Missing |
| Event | ✅ event.ts | ❌ | Missing |
| Exit | ✅ exit.tsx | ❌ | Missing |
| Helper | ✅ helper.tsx | ❌ | Missing |
| Keybind | ✅ keybind.tsx | ❌ | Missing |
| KV | ✅ kv.tsx | ❌ | Missing |
| Local | ✅ local.tsx | ❌ | Missing |
| Plugin keybinds | ✅ plugin-keybinds.ts | ❌ | Missing |
| Project | ✅ project.tsx | ❌ | Missing |
| Prompt | ✅ prompt.tsx | ❌ | Missing |
| Route | ✅ route.tsx | ❌ | Missing |
| SDK | ✅ sdk.tsx | ❌ | Missing |
| Sync | ✅ sync.tsx | ❌ | Missing |
| Theme | ✅ theme.tsx (31002 bytes) | ❌ | Missing |
| TUI Config | ✅ tui-config.tsx | ❌ | Missing |

---

## 5.0 Missing Modules

### 5.1 Opencode Modules with No Freecode Equivalent

Opencode has 45+ modules in `packages/opencode/src/`:

| Module | opencode | freecode | Status |
|--------|----------|----------|--------|
| `account` | ✅ | ❌ | Missing |
| `acp` | ✅ | ⚠️ Stub | ACP protocol stub |
| `audio` | ✅ | ❌ | Missing |
| `bus` | ✅ | ❌ | Missing |
| `command` | ✅ | ❌ | Missing |
| `control-plane` | ✅ | ❌ | Missing |
| `effect` | ✅ | ❌ | Missing |
| `env` | ✅ | ❌ | Missing |
| `file` | ✅ | ❌ | Missing |
| `format` | ✅ | ❌ | Missing |
| `git` | ✅ | ❌ | Missing |
| `id` | ✅ | ❌ | Missing |
| `ide` | ✅ | ❌ | Missing |
| `installation` | ✅ | ❌ | Missing |
| `lsp` | ✅ | ❌ | Missing |
| `patch` | ✅ | ❌ | Missing |
| `permission` | ✅ | ❌ | Missing |
| `plugin` | ✅ | ❌ | Missing |
| `project` | ✅ | ❌ | Missing |
| `pty` | ✅ | ❌ | Missing |
| `question` | ✅ | ❌ | Missing |
| `share` | ✅ | ❌ | Missing |
| `skill` | ✅ | ❌ | Missing |
| `snapshot` | ✅ | ❌ | Missing |
| `storage` | ✅ | ❌ | Missing |
| `sync` | ✅ | ❌ | Missing |
| `temporary` | ✅ | ❌ | Missing |
| `util` | ✅ (36 files) | ⚠️ Partial | Some utils present |
| `v2` | ✅ | ❌ | Missing |
| `worktree` | ✅ | ❌ | Missing |

### 5.2 Freecode Modules with No Opencode Equivalent

| Module | freecode | opencode | Notes |
|--------|----------|----------|-------|
| `auth` | ✅ | ❌ | Credential storage |
| `config` | ✅ | ❌ | Configuration |
| `fleet` | ✅ | ❌ | Fleet management |
| `hook` | ✅ | ❌ | Lifecycle hooks |
| `platform` | ✅ | ❌ | Platform-specific |
| `provider` | ✅ | ❌ | 50+ native providers |
| `shell` | ✅ | ❌ | Shell integration |
| `tool` | ✅ | ❌ | Tool registry |

---

## 6.0 Stub Implementations

### 6.1 Agent Stubs

All 11 agents exist as stubs in `internal/agent/sisyphus.go`. Each agent returns a placeholder response without any actual prompting or behavior.

| Agent | File | Status |
|-------|------|--------|
| Sisyphus | sisyphus.go | Stub |
| Hephaestus | sisyphus.go | Stub |
| Oracle | sisyphus.go | Stub |
| Librarian | sisyphus.go | Stub |
| Explore | sisyphus.go | Stub |
| Prometheus | ❌ | Missing |
| Metis | ❌ | Missing |
| Momus | ❌ | Missing |
| Atlas | ❌ | Missing |
| Multimodal-Looker | ❌ | Missing |
| Sisyphus-Junior | ❌ | Missing |

### 6.2 Hook Implementation

| Hook Tier | Registry | Status |
|-----------|----------|--------|
| Session (26) | hook/registry.go | Registry only, no triggers |
| Tool (14) | hook/registry.go | Registry only, no triggers |
| Transform (5) | hook/registry.go | Registry only, no triggers |
| Continuation (10) | hook/registry.go | Registry only, no triggers |
| Skill (2) | hook/registry.go | Registry only, no triggers |
| Ralph (3) | hook/registry.go | Registry only, no triggers |

### 6.3 Skills System

| Skill | Documented | Implemented |
|-------|------------|-------------|
| git-master | Yes | ❌ Missing |
| playwright | Yes | ❌ Missing |
| frontend-ui-ux | Yes | ❌ Missing |
| review-work | Yes | ❌ Missing |
| ai-slop-remover | Yes | ❌ Missing |

---

## 7.0 Implementation Priority

### Tier 1: Critical (Must Have)

1. **Agent Prompts** - Implement actual agent prompting for all 11 agents
2. **Hook Triggers** - Wire up hook system to actually fire on events
3. **Skills System** - Implement skill registry and execution

### Tier 2: High Value

4. **Missing CLI Commands** - `account`, `web`, `cmd`, `plug`, `generate`
5. **TUI Basic Components** - Dialog system, status bar, session tabs
6. **Provider System** - Already excellent (50+ native providers)

### Tier 3: Medium Value

7. **TUI Advanced Components** - Command palette, context menus, theme system
8. **Context Providers** - File context, global sync
9. **Missing Modules** - skill, lsp, pty, sync

### Tier 4: Nice to Have

10. **Full TUI Parity** - Logo, complex dialogs (not achievable without significant work)

---

## 8.0 Reference File Paths

### OpenCode Reference Files

| Feature | Path |
|---------|------|
| CLI commands | `/Users/mlapointe/git/opencode/packages/opencode/src/cli/cmd/` |
| App components | `/Users/mlapointe/git/opencode/packages/app/src/components/` |
| App pages | `/Users/mlapointe/git/opencode/packages/app/src/pages/` |
| App hooks | `/Users/mlapointe/git/opencode/packages/app/src/hooks/` |
| TUI components | `/Users/mlapointe/git/opencode/packages/opencode/src/cli/cmd/tui/component/` |
| TUI ui | `/Users/mlapointe/git/opencode/packages/opencode/src/cli/cmd/tui/ui/` |
| TUI context | `/Users/mlapointe/git/opencode/packages/opencode/src/cli/cmd/tui/context/` |
| Agent prompts | `/Users/mlapointe/git/opencode/packages/opencode/src/agent/prompt/` |

### Freecode Implementation Files

| Feature | Path |
|---------|------|
| Agent stubs | `internal/agent/sisyphus.go` |
| Agent engine | `internal/agent/engine.go` |
| Hook registry | `internal/hook/registry.go` |
| CLI commands | `internal/cli/` |
| TUI | `internal/ui/` |

---

## 9.0 Architectural Notes

### 9.1 Why TUI Parity Is Difficult

OpenCode's TUI is a sophisticated Solid.js application:
- 27,000+ lines in `app.tsx` alone
- 20+ nested context providers
- 15+ dialog types
- Full keyboard/mouse support via @opentui/core
- Effect (fp-ts alternative) for async operations

**Freecode's Bubble Tea TUI** cannot achieve full parity without essentially writing a new UI framework.

### 9.2 Provider Strategy Success

Freecode's native provider implementation (50+ providers) is **superior** to opencode's LiteLLM-only approach. This strategy should be replicated for other features.

### 9.3 Recommended Approach

Instead of trying to port TypeScript/Solid.js:
1. Keep Go/Bubble Tea for TUI
2. Implement missing features as CLI tools
3. Focus on agent prompting and hooks
4. Accept that full TUI parity is not achievable

---

**Author:** Mark LaPointe <mark@cloudbsd.org>

**Last Updated:** 2026-05-02

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author
