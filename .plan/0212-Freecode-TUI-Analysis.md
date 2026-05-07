# Freecode Feature Parity Plan

**Status**: TUI Implementation In Progress
**Last Updated**: 2026-05-06
**Related Document**: [0213-Freecode-Missing-Features.md](./0213-Freecode-Missing-Features.md) - Comprehensive gap analysis

## Executive Summary

OpenCode and Freecode have fundamentally different architectures that make full feature parity impossible in certain areas. This plan identifies achievable parity goals and a realistic path forward.

**Key Finding:** Freecode leads in provider count (50+ native vs LiteLLM-only) and hook system (60+ vs ~20). OpenCode leads in TUI sophistication and has several missing CLI commands that freecode already has.

## Architectural Mismatch

| Aspect | OpenCode | Freecode | Parity Status |
|--------|----------|----------|---------------|
| **Language** | TypeScript (Node/Bun) | Go | Incompatible - different languages |
| **TUI Framework** | Solid.js + @opentui/* | Bubble Tea (Go) | **Requires complete rewrite** |
| **TUI Code Size** | 27,000+ lines (app.tsx alone) | ~800 lines (in progress) | **In progress** |
| **Commands** | 21 files | 27 files | **1 command missing** |
| **UI Complexity** | 20+ context providers, 15+ dialog types | Components being built | **In progress** |

## TUI Implementation (2026-05-06)

Freecode's TUI is now being built with Bubble Tea:

### Implemented Components
- `TabBarComponent` - Horizontal tab bar with active/inactive styling
- `StatusBar` - Status bar with model, agent, YOLO indicator
- `MessageList` - Scrollable message history display
- `InputArea` - Input with cursor, history navigation
- `CommandPalette` - Fuzzy-searchable command palette
- `Sidebar` - Session list with selection

### TUI Features
- Tab management: Ctrl+T (new), Ctrl+W (close), Tab/Shift+Tab (switch)
- Command palette: Ctrl+P with fuzzy search
- Sidebar toggle: Ctrl+B
- YOLO mode: Ctrl+Y
- Split views: Ctrl+Shift+V/H (placeholder)
- Message scrolling: j/k or arrow keys

### Key Files
| File | Purpose |
|------|---------|
| `internal/ui/model.go` | Main TUI model, wires all components |
| `internal/ui/tabbar.go` | Tab bar component |
| `internal/ui/statusbar.go` | Status bar component |
| `internal/ui/messagelist.go` | Message display |
| `internal/ui/inputarea.go` | Input handling |
| `internal/ui/palette.go` | Command palette |
| `internal/ui/sidebar.go` | Session sidebar |

### OpenCode TUI Reference (for parity)
OpenCode's TUI uses @opentui/core which is TypeScript-only. Key areas to implement:
- Route-based navigation (Home, Session)
- Dialog system (providers, agents, models, sessions, MCPs, themes, help)
- Session/thread view with message rendering
- Permission prompts
- Timeline/fork dialogs

## Missing Commands (Accurate as of 2026-05-02)

**OpenCode commands (21 unique):** account, acp, agent, cmd, db, export, generate, github, import, mcp, models, plug, pr, providers, run, serve, session, stats, uninstall, upgrade, web

**Freecode commands (27 unique):** acp, agent, attach, account, db, debug, doctor, export, github, import, mcp, models, plugin, pr, providers, root, run, serve, session, stats, tab, uninstall, upgrade, version, web

**Completed (2):**
| Command | Status |
|---------|--------|
| `account` | ✅ Done - Login/logout/org switching |
| `web` | ✅ Done - API server + embedded web UI + browser auto-open |

**Still Missing (1):**
| Command | Priority | Notes |
|---------|----------|-------|
| `generate` | Low | Code generation |

**Extra in freecode (9):** attach, debug, doctor, plugin, root, tab, version (some are stubs)

## Implementation Strategy

### Phase 1: Completed Commands ✅
1. `account` - Account management ✅ DONE
2. `web` - Web interface launcher ✅ DONE

### Phase 2: Remaining Missing Commands
1. `generate` - Code generation (Low priority)

### Phase 3: Agent Implementation ✅ (Done)
1. ✅ Implement actual agent prompts for all 11 agents
2. ✅ Wire up hook triggers to fire on events
3. ✅ Implement skill system with built-in skills

### Phase 3: TUI Enhancement (Ongoing)
Incrementally enhance the Bubble Tea TUI:
1. Add proper session management UI
2. Add provider connection dialog
3. Add command palette
4. Add tab management
5. Add status bar with model/agent info

### Phase 4: Advanced Features (Future)
- ACP/fleet implementation
- Plugin system
- GitHub integration

## Files Status

### Completed Files ✅
- `internal/cli/account.go` - Account management ✅ DONE
- `internal/cli/web.go` - Web interface launcher ✅ DONE
- `internal/agent/prompts.go` - Agent prompts ✅ DONE
- `internal/hook/triggers.go` - Hook triggers ✅ DONE
- `internal/hook/builtins.go` - Hook defaults ✅ DONE
- `internal/session/manager.go` - Session manager ✅ DONE
- `internal/session/store.go` - Session store ✅ DONE
- `internal/session/compaction.go` - Session compaction ✅ DONE

### Existing Commands to Enhance
- `internal/cli/providers.go` - Enhance OAuth flows
- `internal/cli/models.go` - Add missing subcommands
- `internal/cli/mcp.go` - Add missing subcommands

### TUI Files (In Progress)
- `internal/ui/model.go` - Enhancement in progress
- `internal/ui/components/` - New components directory
- `internal/ui/dialogs/` - New dialogs directory

### Reference: OpenCode Missing Commands
| Command | Size | Path |
|---------|------|------|
| account | 7889 bytes | `/opencode/packages/opencode/src/cli/cmd/account.ts` |
| web | 2462 bytes | `/opencode/packages/opencode/src/cli/cmd/web.ts` |
| cmd | 183 bytes | `/opencode/packages/opencode/src/cli/cmd/cmd.ts` |
| plug | 6952 bytes | `/opencode/packages/opencode/src/cli/cmd/plug.ts` |
| generate | 2983 bytes | `/opencode/packages/opencode/src/cli/cmd/generate.ts` |

## Key OpenCode Files for Reference

| Command | Path |
|---------|------|
| providers | `/opencode/packages/opencode/src/cli/cmd/providers.ts` |
| acp | `/opencode/packages/opencode/src/cli/cmd/acp.ts` |
| web | `/opencode/packages/opencode/src/cli/cmd/web.ts` |
| plugin | `/opencode/packages/opencode/src/cli/cmd/plug.ts` |
| db | `/opencode/packages/opencode/src/cli/cmd/db.ts` |
| debug | `/opencode/packages/opencode/src/cli/cmd/debug/` |
| export | `/opencode/packages/opencode/src/cli/cmd/export.ts` |
| import | `/opencode/packages/opencode/src/cli/cmd/import.ts` |
| github | `/opencode/packages/opencode/src/cli/cmd/github.ts` |
| pr | `/opencode/packages/opencode/src/cli/cmd/pr.ts` |
| uninstall | `/opencode/packages/opencode/src/cli/cmd/uninstall.ts` |
| account | `/opencode/packages/opencode/src/cli/cmd/account.ts` |
| TUI app | `/opencode/packages/opencode/src/cli/cmd/tui/app.tsx` |
| TUI thread | `/opencode/packages/opencode/src/cli/cmd/tui/thread.ts` |

## Risks & Limitations

1. **TUI Parity**: OpenCode's Solid.js TUI will NEVER be matched by Bubble Tea
2. **TypeScript→Go**: Commands must be rewritten, not ported
3. **Effect System**: OpenCode uses Effect (fp-ts alternative) - hard to replicate in Go
4. **Plugin System**: Requires npm registry access and Node.js integration

## Success Criteria

| # | Criteria | Status | Priority |
|---|----------|--------|----------|
| 1 | `account` command implemented | ✅ Done | HIGH |
| 2 | `web` command implemented | ✅ Done | MEDIUM |
| 3 | All 11 agents have real prompts | ✅ Done | CRITICAL |
| 4 | Hook triggers fire on events | ✅ Done | CRITICAL |
| 5 | Web UI embedded | ✅ Done (basic) | MEDIUM |
| 6 | TUI enhancements (tabs, dialogs) | ⏳ Ongoing | MEDIUM |
| 7 | Full TUI parity | ❌ Not achievable | N/A |

## Related Documents

- [0213-Freecode-Missing-Features.md](./0213-Freecode-Missing-Features.md) - Comprehensive gap analysis
- [0204-Freecode-Features.md](./0204-Freecode-Features.md) - All documented features (agents marked as stubs)
- [0300-Freecode-Implementation-Tasks.md](./0300-Freecode-Implementation-Tasks.md) - Detailed implementation tasks