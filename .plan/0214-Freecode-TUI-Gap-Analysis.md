# Freecode TUI Gap Analysis vs OpenCode

**Status**: Research Complete
**Last Updated**: 2026-05-09
**Related Document**: [0212-Freecode-TUI-Analysis.md](./0212-Freecode-TUI-Analysis.md) - Implementation status

## Executive Summary

Freecode's TUI is significantly behind OpenCode's. OpenCode has a sophisticated component library (~100+ components) while Freecode has ~20 basic components. This document catalogs the gaps and provides a path forward.

## OpenCode TUI Structure

### Framework & Architecture
- **Framework**: Solid.js + @opentui/core (TypeScript)
- **TUI Location**: `packages/opencode/src/cli/cmd/tui/`
- **UI Library**: `packages/ui/` (shared component library)
- **Key Files**:
  - `app.tsx` - Main TUI application (27KB)
  - `thread.ts` - Session/thread management
  - `ui/dialog*.tsx` - Dialog system (6 dialog types)
  - `component/` - TUI-specific components

### OpenCode UI Components (`packages/ui/src/components/`)

| Component | Size | Purpose |
|-----------|------|---------|
| `message-part.tsx` | 75KB | **Core message rendering with parts** |
| `file.tsx` | 32KB | File display with syntax highlighting |
| `icon.tsx` | 39KB | Icon library (Lucide icons) |
| `markdown.tsx` | 10KB | Markdown rendering |
| `session-turn.tsx` | 19KB | Session turn/message display |
| `session-review.tsx` | 27KB | Session review interface |
| `timeline-playground.stories.tsx` | 65KB | Timeline visualization |
| `thinking-heading.stories.tsx` | 28KB | Thinking/reasoning display |
| `line-comment.tsx` | 13KB | Code line annotations |
| `basic-tool.tsx` | 8KB | Tool rendering |
| `tool-error-card.tsx` | 5KB | Tool error display |
| `tool-status-title.tsx` | 4KB | Tool status display |
| `text-reveal.tsx` | 4KB | Animated text reveal |
| `spinner.ts` | 12KB | Loading spinner |
| `dialog.tsx` | 5KB | Dialog base component |
| `toast.tsx` | 5KB | Toast notifications |
| `select.tsx` | 5KB | Selection component |
| `dropdown-menu.tsx` | 9KB | Dropdown menus |
| `context-menu.tsx` | 9KB | Context menus |
| `tabs.tsx` | 3KB | Tab component |
| `checkbox.tsx` | 2KB | Checkbox |
| `switch.tsx` | 1KB | Toggle switch |
| `accordion.tsx` | 2KB | Collapsible sections |
| `popover.tsx` | 4KB | Popover component |
| `tooltip.tsx` | 4KB | Tooltips |
| `list.tsx` | 13KB | List component |
| `scroll-view.tsx` | 7KB | Scrollable view |
| `resize-handle.tsx` | 2KB | Resizable panels |

**Total Components**: 100+ React/Solid components

### OpenCode TUI Dialogs (`cli/cmd/tui/ui/`)

| Dialog | Purpose |
|--------|---------|
| `dialog.tsx` | Base dialog with backdrop, keyboard handling |
| `dialog-select.tsx` | Selection dialog (14KB) |
| `dialog-export-options.tsx` | Export options |
| `dialog-help.tsx` | Help dialog |
| `dialog-prompt.tsx` | Prompt dialog |
| `dialog-confirm.tsx` | Confirmation dialog |
| `dialog-alert.tsx` | Alert dialog |
| `spinner.ts` | Loading spinner |
| `toast.tsx` | Toast notifications |
| `link.tsx` | Link component |

### OpenCode TUI Key Features

1. **Message Rendering** (`message-part.tsx`)
   - Multi-part messages (text, tool calls, tool results, images)
   - Syntax highlighting for code blocks
   - File diffs with inline annotations
   - Thinking/reasoning blocks with visibility toggle
   - Copy-to-clipboard on selection

2. **Session Management** (`thread.ts`)
   - Session continuation
   - Fork sessions
   - Session review/retry
   - Turn-based conversation view

3. **Tool Integration**
   - Tool call/result display
   - Error states with retry
   - Progress indicators
   - Permission requests

4. **Dialog System**
   - Modal dialogs with backdrop
   - Keyboard navigation (vim-style)
   - Mouse support (click to dismiss)
   - Selection handling

5. **Context Providers** (from `app.tsx`)
   - Dialog context
   - Toast context
   - Theme context
   - Selection context
   - Keyboard context
   - Renderer context

## Freecode TUI Structure

### Framework & Architecture
- **Framework**: Bubble Tea (Go)
- **TUI Location**: `internal/ui/`
- **Components**: ~20 components

### Freecode Components

| Component | Size | Purpose |
|-----------|------|---------|
| `model.go` | 31KB | Main TUI model, all logic |
| `messagelist.go` | 5KB | Message display |
| `inputarea.go` | 5KB | Input handling |
| `sidebar.go` | 4KB | Session list |
| `tabbar.go` | 3KB | Tab bar |
| `statusbar.go` | 3KB | Status bar |
| `palette.go` | 5KB | Command palette |
| `tooldialog.go` | 4KB | Tool toggle dialog |
| `question.go` | 11KB | Question dialog |
| `permission.go` | 10KB | Permission dialog |
| `select.go` | 5KB | Selection dialog |
| `autocomplete.go` | 9KB | Autocomplete |
| `timeline.go` | 14KB | Timeline display |
| `fleet.go` | 10KB | Fleet panel |
| `animation.go` | 2KB | Animation manager |
| `sound.go` | 2KB | Sound manager |
| `console.go` | 4KB | Debug console |
| `error.go` | 6KB | Error display |
| `exportdialog.go` | 5KB | Export dialog |
| `mcpdialog.go` | 3KB | MCP dialog |
| `helpdialog.go` | 3KB | Help dialog |
| `statusdialog.go` | 2KB | Status dialog |
| `style.go` | 9KB | Styling |
| `tab/` | - | Tab subcomponent |

**Total Components**: ~20 Go components

## Critical Gaps

### 1. Message Rendering ❌ CRITICAL
**OpenCode**: 75KB `message-part.tsx` with:
- Multi-part message support
- Inline tool calls with results
- Code syntax highlighting
- Diff display with annotations
- Thinking blocks with collapse
- Image inline display

**Freecode**: `messagelist.go` - basic text rendering only
- No tool call rendering
- No syntax highlighting
- No diff display
- No image support

### 2. Tool Call/Result Display ❌ CRITICAL
**OpenCode**: `basic-tool.tsx`, `tool-error-card.tsx`, `tool-status-title.tsx`
- Full tool state visualization
- Error cards with retry
- Progress indicators

**Freecode**: No dedicated tool rendering
- Tool calls shown as plain text
- No error visualization
- No progress display

### 3. Dialog System ❌ HIGH
**OpenCode**: Full modal system with:
- Backdrop with click-to-dismiss
- Keyboard navigation
- Focus management
- Selection integration

**Freecode**: Basic dialogs
- No backdrop
- No focus management
- No selection integration
- Keyboard handling is ad-hoc

### 4. Session Review/Retry ❌ HIGH
**OpenCode**: `session-review.tsx` (27KB)
- View past sessions
- Retry failed turns
- Fork sessions

**Freecode**: No session review UI
- Only current session visible

### 5. Component Library ❌ MEDIUM
**OpenCode**: 100+ components with:
- Storybook documentation
- Consistent styling
- Accessibility

**Freecode**: 20 basic components
- No documentation
- Inconsistent styling
- No accessibility

## Recommended Priority

### Phase 1: Message Rendering (CRITICAL)
1. Add tool call rendering to MessageList
2. Add syntax highlighting for code blocks
3. Add basic diff display
4. Add thinking/collapse support

### Phase 2: Tool Display (CRITICAL)
1. Create ToolCall component
2. Create ToolResult component
3. Create ToolError component
4. Add progress indicators

### Phase 3: Dialog System (HIGH)
1. Add backdrop to dialogs
2. Add focus management
3. Add keyboard navigation
4. Add selection integration

### Phase 4: Session Review (HIGH)
1. Create session list view
2. Add session preview
3. Add fork/continue actions
4. Add retry capability

### Phase 5: Component Library (MEDIUM)
1. Document all components
2. Add consistent styling
3. Add animation support
4. Add accessibility

## Files to Create/Modify

### New Files
- `internal/ui/messagepart.go` - Message part rendering
- `internal/ui/toolcall.go` - Tool call display
- `internal/ui/toolresult.go` - Tool result display
- `internal/ui/sessionreview.go` - Session review UI
- `internal/ui/dialogbase.go` - Base dialog with backdrop

### Files to Modify
- `internal/ui/messagelist.go` - Add part rendering
- `internal/ui/model.go` - Add session review
- `internal/ui/dialog.go` - Add backdrop support

## Conclusion

Freecode's TUI has a solid foundation but is missing most of OpenCode's sophisticated UI features. The most critical gaps are:
1. Message rendering (no tool calls, no syntax highlighting)
2. Tool display (no visual tool state)
3. Dialog system (no backdrop, no focus management)

Given the architectural difference (Solid.js vs Bubble Tea), full parity is not achievable. However, improving message rendering and tool display should be the priority.
