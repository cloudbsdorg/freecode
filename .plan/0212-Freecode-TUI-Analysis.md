# Freecode Feature Parity Plan

**Status**: Analysis Complete
**Last Updated**: 2026-05-02

## Executive Summary

OpenCode and Freecode have fundamentally different architectures that make full feature parity impossible in certain areas. This plan identifies achievable parity goals and a realistic path forward.

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

## Missing Commands (Priority Order)

### Tier 1: High Value, Achievable
1. **`providers`** - Provider management with OAuth/API key flows (526 lines)
   - `providers list` - Show configured credentials
   - `providers login` - Interactive login with multiple auth methods
   - `providers logout` - Remove credentials

2. **`models`** - Enhance existing (already exists, needs enhancement)
   - Already has `models list` and `models delete`

3. **`mcp`** - Enhance existing (already exists, needs enhancement)
   - Already has basic functionality

### Tier 2: Medium Value, Achievable
4. **`web`** - Start web interface (81 lines)
   - Start server and open browser
   - Show network/localhost URLs
   - Already partially exists as `serve` command

5. **`session`** - Enhance existing
   - Already has basic functionality

6. **`stats`** - Enhance existing
   - Already has basic functionality

7. **`agent`** - Enhance existing
   - Already has basic functionality

### Tier 3: Complex/Optional
8. **`acp`** - Agent Client Protocol server (70 lines)
   - Requires ACP protocol implementation
   - Fleet/cluster functionality

9. **`plugin`** / `plug` - Plugin system (233 lines)
   - NPM package installation
   - Plugin manifest parsing
   - Config patching

10. **`db`** - Database operations (3852 bytes)
    - SQLite management

11. **`debug`** - Debugging tools (8 subcommands)
    - `debug agent` - Agent debugging
    - `debug config` - Config inspection
    - `debug file` - File operations
    - `debug lsp` - LSP debugging
    - `debug ripgrep` - Search debugging
    - `debug snapshot` - Heap snapshots

12. **`export`** / **`import`** - Data portability (10103 + 6868 bytes)
    - Session export/import
    - Configuration portability

13. **`github`** - GitHub integration (59095 bytes)
    - PR creation, review, merge
    - Issue management
    - Repository operations

14. **`pr`** - Pull request management (5122 bytes)
    - Dedicated PR command

15. **`uninstall`** - Clean removal (10338 bytes)
    - Config cleanup
    - Cache removal

16. **`account`** - Account management (7889 bytes)
    - User account operations

## Implementation Strategy

### Phase 1: Missing Commands (Weeks 1-4)
Implement the 11 missing CLI commands in priority order:
1. `providers` - Most user-facing impact
2. `web` - Easy win (similar to existing `serve`)
3. `plugin` - Extensibility
4. `acp` - Fleet functionality
5. Others as needed

### Phase 2: TUI Enhancement (Ongoing)
Incrementally enhance the Bubble Tea TUI:
1. Add proper session management UI
2. Add provider connection dialog
3. Add command palette (matching OpenCode's 50+ commands)
4. Add tab management
5. Add status bar with model/agent info

### Phase 3: Advanced Features (Future)
- ACP/fleet implementation
- Plugin system
- GitHub integration

## Files to Modify

### New Commands
- `internal/cli/providers.go` - Provider management
- `internal/cli/web.go` - Web interface
- `internal/cli/plugin.go` - Plugin system
- `internal/cli/acp.go` - ACP server
- `internal/cli/debug.go` - Debugging tools
- `internal/cli/export.go` - Export functionality
- `internal/cli/import.go` - Import functionality
- `internal/cli/github.go` - GitHub integration
- `internal/cli/pr.go` - PR management
- `internal/cli/uninstall.go` - Uninstall
- `internal/cli/account.go` - Account management

### Existing Commands to Enhance
- `internal/cli/models.go` - Add missing subcommands
- `internal/cli/mcp.go` - Add missing subcommands
- `internal/cli/session.go` - Add missing subcommands

### TUI Files
- `internal/ui/model.go` - Major enhancement needed
- `internal/ui/components/` - New components directory
- `internal/ui/dialogs/` - New dialogs directory

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

1. ✅ All 11 missing commands implemented in Go
2. ✅ `freecode providers login` works with at least 3 providers
3. ✅ `freecode web` opens the web interface
4. ✅ `freecode plugin install <pkg>` works
5. ⏳ TUI enhancements (continuous improvement)
6. ❌ Full TUI parity (not achievable)