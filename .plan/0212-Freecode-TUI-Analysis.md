# Freecode Feature Parity Plan

**Status**: Analysis Complete
**Last Updated**: 2026-05-02
**Related Document**: [0213-Freecode-Missing-Features.md](./0213-Freecode-Missing-Features.md) - Comprehensive gap analysis

## Executive Summary

OpenCode and Freecode have fundamentally different architectures that make full feature parity impossible in certain areas. This plan identifies achievable parity goals and a realistic path forward.

**Key Finding:** Freecode leads in provider count (50+ native vs LiteLLM-only) and hook system (60+ vs ~20). OpenCode leads in TUI sophistication and has several missing CLI commands that freecode already has.

## Architectural Mismatch

| Aspect | OpenCode | Freecode | Parity Status |
|--------|----------|----------|---------------|
| **Language** | TypeScript (Node/Bun) | Go | Incompatible - different languages |
| **TUI Framework** | Solid.js + @opentui/* | Bubble Tea (Go) | **Cannot port - requires complete rewrite** |
| **TUI Code Size** | 27,000+ lines (app.tsx alone) | 160 lines (stub) | **Not achievable without years of work** |
| **Commands** | 22 files | 10 files | **11 commands missing** |
| **UI Complexity** | 20+ context providers, 15+ dialog types | None | **Not achievable** |

## TUI Reality

OpenCode's TUI is a sophisticated Solid.js application with:
- 20+ nested context providers
- 15+ dialog types (providers, agents, models, sessions, MCPs, themes, help, etc.)
- Route-based navigation (Home, Session, Plugin)
- 50+ registered commands
- Full keyboard/mouse support via @opentui/core

**This cannot be ported to Go/Bubble Tea** without essentially writing an entirely new UI framework.

### Recommendation: Keep Bubble Tea, Enhance Incrementally

Instead of trying to match OpenCode's Solid.js TUI:
1. Keep the existing Bubble Tea implementation
2. Add missing commands as CLI tools (which work regardless of TUI)
3. Focus on making the CLI commands feature-complete
4. Incrementally enhance the TUI where feasible

## Missing Commands (Accurate as of 2026-05-02)

**OpenCode commands (21 unique):** account, acp, agent, cmd, db, export, generate, github, import, mcp, models, plug, pr, providers, run, serve, session, stats, uninstall, upgrade, web

**Freecode commands (24 unique):** acp, agent, attach, db, debug, doctor, export, github, import, mcp, models, plugin, pr, providers, root, run, serve, session, stats, tab, uninstall, upgrade, version

**Truly missing from freecode (0):**
| Command | Priority | Notes |
|---------|----------|-------|
| `account` | High | Account management ✅ DONE |
| `web` | ✅ DONE | API server + embedded web UI (grows with TUI features) |
| `cmd` | N/A | Not applicable - Cobra provides this |
| `plug` | N/A | Not applicable - Node.js plugin system |
| `generate` | Low | Code generation (check if needed) |

**Extra in freecode (8):** attach, debug, doctor, plugin, root, tab, version (some are stubs)

### Tier 1: High Value, Achievable
1. **`account`** - Account management (7889 bytes in opencode)
   - User authentication
   - Account settings
   - Subscription management

### Tier 2: Medium Value, Achievable
2. **`web`** - Start web interface (2462 bytes in opencode)
   - Start server and open browser
   - Show network/localhost URLs
   - Differs from `serve` (web has browser auto-open)

3. **`cmd`** - CLI command framework (183 bytes in opencode)
   - Command registration
   - Help text generation

### Tier 3: Complex/Optional
4. **`plug`** - Plugin system (6952 bytes in opencode)
   - NPM package installation
   - Plugin manifest parsing
   - Config patching
   - Freecode has `plugin` but different implementation

5. **`generate`** - Code generation (2983 bytes in opencode)
   - Template-based code generation
   - Scaffold new projects

## Implementation Strategy

### Phase 1: Missing Commands (Priority Order)
Implement the 5 missing CLI commands:
1. `account` - Account management (High user-facing impact)
2. `web` - Web interface launcher (Medium, differs from `serve`)
3. `cmd` - CLI framework (Medium, enables command registration)
4. `plug` - Plugin system (Low, freecode has `plugin` with different approach)
5. `generate` - Code generation (Low)

### Phase 2: Agent Implementation (Critical)
**This is the highest priority work:**
1. Implement actual agent prompts for all 11 agents
2. Wire up hook triggers to fire on events
3. Implement skill system with built-in skills

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

## Files to Modify

### Missing Commands to Create
- `internal/cli/account.go` - Account management (HIGH PRIORITY)
- `internal/cli/web.go` - Web interface launcher
- `internal/cli/cmd.go` - CLI command framework
- `internal/cli/plug.go` - Plugin system
- `internal/cli/generate.go` - Code generation

### Existing Commands to Enhance
- `internal/cli/providers.go` - Already exists, enhance OAuth flows
- `internal/cli/models.go` - Already exists, add missing subcommands
- `internal/cli/mcp.go` - Already exists, add missing subcommands

### TUI Files
- `internal/ui/model.go` - Major enhancement needed
- `internal/ui/components/` - New components directory
- `internal/ui/dialogs/` - New dialogs directory

### Agent Files (CRITICAL - Currently Stubs)
- `internal/agent/sisyphus.go` - Implement actual agent prompts
- `internal/agent/engine.go` - Wire up agent execution
- `internal/hook/registry.go` - Implement hook triggers

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