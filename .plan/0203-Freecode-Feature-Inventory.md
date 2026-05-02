# Freecode — Feature Inventory

This document catalogs all features from opencode and oh-my-openagent, indicating their inclusion status in freecode.

---

## Feature Table

| # | Product | Feature | Summary | Included | Task IDs | Dependencies |
|---|---------|---------|---------|----------|----------|-------------|
| 1 | opencode | Multi-Agent System | 11 specialized agents (Sisyphus, Hephaestus, Oracle, etc.) that collaborate | Yes | FREECODE-001-T1, FREECODE-001-T2, FREECODE-001-T3 | - |
| 2 | opencode | Category-Based Routing | Routes tasks to optimal models based on category (visual, ultra, deep, etc.) | Yes | FREECODE-002-T1, FREECODE-002-T2 | FREECODE-001 |
| 3 | opencode | Lifecycle Hooks (60+) | Hooks for session, tool, transform, continuation, skill, ralph events | Yes | FREECODE-003-T1, FREECODE-003-T2, FREECODE-003-T3 | - |
| 4 | oh-my-openagent | Hashline Edit | Uses line content hashes instead of line numbers for reliability (68% vs 6.7% success) | Yes | FREECODE-004-T1, FREECODE-004-T2, FREECODE-004-T3 | - |
| 5 | oh-my-openagent | Slash Commands | /ralph-loop, /init-deep, /handoff, /start-work, /refactor, etc. | Yes | FREECODE-005-T1, FREECODE-005-T2, FREECODE-005-T3 | FREECODE-001 |
| 6 | oh-my-openagent | Background Agents | Multiple agents running in tmux panes simultaneously | Yes | FREECODE-006-T1, FREECODE-006-T2, FREECODE-006-T3, FREECODE-006-T4 | FREECODE-071 |
| 7 | opencode | Session Tools | session_list, session_read, session_search, session_info, session_export, session_import | Yes | FREECODE-007-T1, FREECODE-007-T2, FREECODE-007-T3 | - |
| 8 | oh-my-openagent | Task Tools | task_create, task_list, task_update with dependency graphs | Yes | FREECODE-008-T1, FREECODE-008-T2, FREECODE-008-T3 | - |
| 9 | opencode | AST-grep | Search/replace across 25+ languages understanding code structure | Yes | FREECODE-009-T1, FREECODE-009-T2, FREECODE-009-T3 | - |
| 10 | opencode | Look_at Tool | Analyze PDFs, images, diagrams with vision models | Yes | FREECODE-010-T1, FREECODE-010-T2 | - |
| 11 | opencode | Interactive Bash | Tmux-based terminal for vim, ssh, htop, mysql, psql, etc. | Yes | FREECODE-011-T1, FREECODE-011-T2, FREECODE-011-T3 | FREECODE-006 |
| 12 | opencode | Built-in MCPs | Exa (websearch), Context7 (docs), Grep.app (GitHub search) | Yes | FREECODE-012-T1, FREECODE-012-T2, FREECODE-012-T3 | - |
| 13 | opencode | LSP Support | Auto-detect and manage LSP servers for 30+ languages (gopls, pyright, rust-analyzer, etc.) | Yes | FREECODE-013-T1, FREECODE-013-T2, FREECODE-013-T3 | - |
| 14 | oh-my-openagent | Commands System | Custom command definitions with templates and per-command agent | Yes | FREECODE-014-T1, FREECODE-014-T2 | FREECODE-005 |
| 15 | opencode | gh Integration | PR creation, review, merge, issues, releases via GitHub CLI | Yes | FREECODE-015-T1, FREECODE-015-T2 | - |
| 16 | opencode | glab Integration | MR creation, issues, snippets via GitLab CLI | Yes | FREECODE-016-T1, FREECODE-016-T2 | - |
| 17 | oh-my-openagent | Fleet Management | Head/agent/client modes for coordinating freecode instances | Yes | FREECODE-017-T1, FREECODE-017-T2, FREECODE-017-T3, FREECODE-017-T4, FREECODE-017-T5 | FREECODE-031, FREECODE-035 |
| 18 | opencode | SMB/CIFS Support | Mount and browse Windows file shares | Yes | FREECODE-018-T1, FREECODE-018-T2 | - |
| 19 | opencode | NFS Support | Mount Unix/Linux NFS exports | Yes | FREECODE-019-T1, FREECODE-019-T2 | - |
| 20 | opencode | AFP Support | Apple Filing Protocol for macOS file shares | Yes | FREECODE-020-T1, FREECODE-020-T2 | - |
| 21 | opencode | BitTorrent | Download/seed files via BitTorrent for multi-agent file sharing | Yes | FREECODE-021-T1, FREECODE-021-T2, FREECODE-021-T3 | - |
| 22 | opencode | YubiKey SSH | PIV-based SSH key storage on YubiKey | Yes | FREECODE-022-T1, FREECODE-022-T2 | - |
| 23 | opencode | YubiKey GPG | GPG key on YubiKey for commits and tags | Yes | FREECODE-023-T1, FREECODE-023-T2 | - |
| 24 | opencode | go-git Fallback | Use go-git when system git unavailable (FreeBSD base) | Yes | FREECODE-024-T1, FREECODE-024-T2 | - |
| 25 | opencode | SVN Support | CLI-based Subversion for FreeBSD development workflow | Yes | FREECODE-025-T1, FREECODE-025-T2 | - |
| 26 | oh-my-openagent | git-master Skill | Atomic commits, rebase surgery, branch management | Yes | FREECODE-026-T1, FREECODE-026-T2, FREECODE-026-T3 | FREECODE-015 |
| 27 | oh-my-openagent | playwright Skill | Built-in browser automation with Playwright | Yes | FREECODE-027-T1, FREECODE-027-T2 | FREECODE-012 |
| 28 | oh-my-openagent | frontend-ui-ux Skill | Design-first UI/UX development workflow | Yes | FREECODE-028-T1, FREECODE-028-T2 | - |
| 29 | oh-my-openagent | review-work Skill | 5 parallel subagent code review workflow | Yes | FREECODE-029-T1, FREECODE-029-T2 | FREECODE-001 |
| 30 | oh-my-openagent | ai-slop-remover Skill | Remove AI-generated code smells | Yes | FREECODE-030-T1, FREECODE-030-T2 | - |
| 31 | opencode | TUI Mouse Support | Full mouse interaction with region registry, clickable elements | Yes | FREECODE-031-T1, FREECODE-031-T2, FREECODE-031-T3 | - |
| 32 | opencode | Context Menus | Context-specific right-click menus | Yes | FREECODE-032-T1, FREECODE-032-T2 | FREECODE-031 |
| 33 | opencode | Scroll Regions | Mouse wheel scrolling in defined regions | Yes | FREECODE-033-T1 | FREECODE-031 |
| 34 | opencode | Text Selection | Click and drag to select text | Yes | FREECODE-034-T1 | FREECODE-031 |
| 35 | opencode | Command Palette | Ctrl+P command palette with fuzzy search | Yes | FREECODE-035-T1, FREECODE-035-T2 | FREECODE-031 |
| 36 | opencode | Fleet Panel | Ctrl+\ dedicated fleet panel with instance list, tasks, logs | Yes | FREECODE-036-T1, FREECODE-036-T2 | FREECODE-017, FREECODE-031 |
| 37 | opencode | Session Tabs | Tabbed interface for multiple concurrent sessions | Yes | FREECODE-037-T1, FREECODE-037-T2, FREECODE-037-T3 | FREECODE-031 |
| 38 | opencode | Split View | Vertical/horizontal session splits | Yes | FREECODE-038-T1, FREECODE-038-T2 | FREECODE-037 |
| 39 | opencode | Ralph Loop | Agent reviews own output, detects self-contradictions | Yes | FREECODE-039-T1, FREECODE-039-T2, FREECODE-039-T3 | FREECODE-001 |
| 40 | opencode | Ultrawork Mode | Extended execution for complex long-running tasks | Yes | FREECODE-040-T1, FREECODE-040-T2 | - |
| 41 | opencode | Dynamic Pruning | Auto-prune old tool outputs, protect recent turns | Yes | FREECODE-041-T1, FREECODE-041-T2 | - |
| 42 | opencode | Runtime Fallback | Retry with different model on failure | Yes | FREECODE-042-T1, FREECODE-042-T2 | - |
| 43 | opencode | Model Fallback Chains | Ordered model lists with per-model settings | Yes | FREECODE-043-T1, FREECODE-043-T2 | FREECODE-042 |
| 44 | opencode | Git Master Integration | Commit footer, co-author, env prefix conventions | Yes | FREECODE-044-T1, FREECODE-044-T2 | FREECODE-015 |
| 45 | oh-my-openagent | MCP OAuth | OAuth 2.0 + PKCE for MCP server authentication | Partially | FREECODE-045-T1, FREECODE-045-T2 | - |
| 46 | opencode | Slack Integration | Bot posts updates to Slack threads | No | - | - |
| 47 | opencode | Discord Integration | Bot posts updates to Discord channels | No | - | - |
| 48 | opencode | PostHog Telemetry | Anonymous usage tracking | No | - | - |
| 49 | opencode | NPM Update Check | Check for npm package updates | No | - | - |
| 50 | opencode | Desktop App (Electron) | Electron-based standalone desktop application | Future | FREECODE-050-T1, FREECODE-050-T2 | - |
| 51 | opencode | Code Interpreter | Sandboxed code execution (OpenAI style) | Future | FREECODE-051-T1, FREECODE-051-T2, FREECODE-051-T3 | - |
| 52 | opencode | Image Generation | AI image generation (GPT-image-1) | Future | FREECODE-052-T1 | - |
| 53 | opencode | Enterprise Share | Share sessions via URL | Future | FREECODE-053-T1, FREECODE-053-T2 | FREECODE-017 |
| 54 | opencode | MDM Managed Config | Enterprise macOS MDM configuration | Future | FREECODE-054-T1 | - |
| 55 | oh-my-openagent | Todo Continuation | Enforces todo completion before session end | Yes | FREECODE-055-T1 | FREECODE-008 |
| 56 | oh-my-openagent | Keyword Detector | Detect ultrawork/search/analyze mode triggers | Yes | FREECODE-056-T1 | FREECODE-003 |
| 57 | oh-my-openagent | Preemptive Compaction | Compact context before hitting limits | Yes | FREECODE-057-T1, FREECODE-057-T2 | FREECODE-041 |
| 58 | oh-my-openagent | Atlas Orchestrator | Master orchestrator for boulder sessions | Yes | FREECODE-058-T1, FREECODE-058-T2 | FREECODE-001, FREECODE-007 |
| 59 | opencode | File Prompt Loading | Load prompts from files with file:// URI support | Yes | FREECODE-059-T1 | - |
| 60 | opencode | Extended Providers | GitLab, Mistral, OpenRouter, Bedrock, Azure, Cloudflare models | Yes | FREECODE-060-T1, FREECODE-060-T2 | - |
| 61 | opencode | Model Capabilities Cache | Cache model capabilities, refresh on schedule | Yes | FREECODE-061-T1, FREECODE-061-T2 | - |
| 62 | opencode | Browser Automation Engine | Configurable browser backend (playwright vs agent-browser) | Yes | FREECODE-062-T1 | FREECODE-027 |
| 63 | oh-my-openagent | Comment Checker | Check comments are accurate and up-to-date | Yes | FREECODE-063-T1 | - |
| 64 | opencode | Parallel Review | 5 parallel subagent review workflow (see review-work) | Yes | FREECODE-029 | FREECODE-001 |
| 65 | opencode | Agent Browser Skill | Alternative browser automation skill (puppeteer) | Yes | FREECODE-027 | FREECODE-012 |
| 66 | opencode | Websearch Exa | Exa web search without API key (rate limited) | Yes | FREECODE-012 | - |
| 67 | opencode | Context7 Docs | Search official library documentation | Yes | FREECODE-012 | - |
| 68 | opencode | Grep.app | Search code across GitHub repos | Yes | FREECODE-012 | - |
| 69 | opencode | Skill MCP | Invoke MCP from skill-embedded servers | Yes | FREECODE-064-T1, FREECODE-064-T2 | FREECODE-012 |
| 70 | opencode | Call Agent Tool | Spawn explore/librarian agents from skill | Yes | FREECODE-065-T1 | FREECODE-001 |
| 71 | opencode | Background Output | Stream background task output | Yes | FREECODE-066-T1 | FREECODE-006 |
| 72 | opencode | Background Cancel | Cancel running background tasks | Yes | FREECODE-066-T1 | FREECODE-006 |
| 73 | opencode | rsync Support | One-way file synchronization over SSH | Yes | FREECODE-067-T1 | - |
| 74 | opencode | VNC Tool | View remote VNC screens | Yes | FREECODE-068-T1 | - |
| 75 | opencode | RDP Tool | Connect to Windows RDP sessions | Yes | FREECODE-069-T1 | - |
| 76 | opencode | SSH Tunnels | SSH local/remote port forwarding | Yes | FREECODE-070-T1 | - |
| 77 | opencode | Tmux Integration | Built-in tmux session and pane management | Yes | FREECODE-071-T1, FREECODE-071-T2, FREECODE-071-T3 | - |
| 78 | opencode | YOLO Mode | Execute dangerous commands without confirmation | Yes | FREECODE-072-T1, FREECODE-072-T2 | - |
| 79 | opencode | Permission Patterns | Pattern-based tool access control (glob matching) | Yes | FREECODE-073-T1, FREECODE-073-T2 | - |
| 80 | opencode | External Directory Access | Limit file access to specific directories | Yes | FREECODE-074-T1 | - |
| 81 | opencode | Doom Loop Detection | Detect repeated tool call patterns | Yes | FREECODE-075-T1 | - |
| 82 | oh-my-openagent | Handoff Command | Generate context summary for new sessions | Yes | FREECODE-076-T1, FREECODE-076-T2 | FREECODE-007 |
| 83 | oh-my-openagent | Init Deep Command | Generate hierarchical AGENTS.md file | Yes | FREECODE-077-T1, FREECODE-077-T2 | - |
| 84 | oh-my-openagent | Start Work Command | Execute Prometheus plans via Atlas | Yes | FREECODE-078-T1, FREECODE-078-T2 | FREECODE-058 |
| 85 | oh-my-openagent | Refactor Command | Full refactoring workflow | Yes | FREECODE-079-T1, FREECODE-079-T2 | FREECODE-004, FREECODE-009 |
| 86 | oh-my-openagent | Explainshell Command | Explain shell commands | Yes | FREECODE-080-T1 | - |
| 87 | oh-my-openagent | Test Command | Generate and run tests | Yes | FREECODE-081-T1, FREECODE-081-T2 | - |
| 88 | oh-my-openagent | Docs Command | Generate documentation | Yes | FREECODE-082-T1, FREECODE-082-T2 | - |

---

## Summary Statistics

| Category | Count |
|----------|-------|
| **Total Features** | 88 |
| **Included (Yes)** | 80 |
| **Partially Included** | 1 |
| **Not Included** | 3 |
| **Future** | 4 |

### By Origin

| Origin | Count | Included |
|--------|-------|----------|
| opencode only | 58 | 51 |
| oh-my-openagent only | 22 | 22 |
| Both | 8 | 7 |

### Not Included Features

| Feature | Reason |
|---------|--------|
| Slack Integration | External dependency, not core to freecode |
| Discord Integration | External dependency, not core to freecode |
| PostHog Telemetry | Privacy-first: NO telemetry in freecode |
| NPM Update Check | Not relevant to freecode |

### Future Features

| Feature | Blocking Dependencies |
|---------|----------------------|
| Desktop App (Electron) | None - separate effort |
| Code Interpreter | Security sandboxing needed |
| Image Generation | Depends on model support |
| Enterprise Share | After fleet is mature |

---

## Dependency Graph

Features grouped by dependency chain:

### Core Infrastructure
```
FREECODE-001 (Multi-Agent System)
├── FREECODE-002 (Category-Based Routing)
├── FREECODE-005 (Slash Commands)
├── FREECODE-029 (Parallel Review)
├── FREECODE-039 (Ralph Loop)
├── FREECODE-058 (Atlas Orchestrator)
│   └── FREECODE-078 (Start Work Command)
├── FREECODE-065 (Call Agent Tool)
└── FREECODE-076 (Handoff Command)

FREECODE-003 (Lifecycle Hooks)
└── FREECODE-056 (Keyword Detector)

FREECODE-007 (Session Tools)
├── FREECODE-055 (Todo Continuation)
└── FREECODE-076 (Handoff Command)
```

### UI/TUI Layer
```
FREECODE-031 (TUI Mouse Support)
├── FREECODE-032 (Context Menus)
├── FREECODE-033 (Scroll Regions)
├── FREECODE-034 (Text Selection)
├── FREECODE-035 (Command Palette)
└── FREECODE-037 (Session Tabs)
    └── FREECODE-038 (Split View)

FREECODE-035 (Command Palette)
└── FREECODE-017 (Fleet Management)
    └── FREECODE-036 (Fleet Panel)
```

### Agents & Background
```
FREECODE-071 (Tmux Integration)
└── FREECODE-006 (Background Agents)
    ├── FREECODE-011 (Interactive Bash)
    └── FREECODE-066 (Background Output/Cancel)
```

### Skills & Commands
```
FREECODE-012 (Built-in MCPs)
├── FREECODE-027 (Playwright Skill)
│   └── FREECODE-062 (Browser Automation Engine)
├── FREECODE-064 (Skill MCP)
└── FREECODE-026 (git-master Skill)
    └── FREECODE-044 (Git Master Integration)
        └── FREECODE-015 (gh Integration)

FREECODE-005 (Slash Commands)
└── FREECODE-014 (Commands System)
```

### Context Management
```
FREECODE-041 (Dynamic Pruning)
└── FREECODE-057 (Preemptive Compaction)

FREECODE-042 (Runtime Fallback)
└── FREECODE-043 (Model Fallback Chains)
```

### Fleet
```
FREECODE-017 (Fleet Management)
├── FREECODE-036 (Fleet Panel)
└── FREECODE-053 (Enterprise Share) [Future]
```

---

## Task ID Reference

### FREECODE-001: Multi-Agent System
| Task | Description |
|------|-------------|
| FREECODE-001-T1 | Implement agent registry and config |
| FREECODE-001-T2 | Implement agent communication bus |
| FREECODE-001-T3 | Implement agent spawning and lifecycle |

### FREECODE-004: Hashline Edit
| Task | Description |
|------|-------------|
| FREECODE-004-T1 | Implement line hashing (sha256, 8 bytes) |
| FREECODE-004-T2 | Implement hash-based line lookup |
| FREECODE-004-T3 | Implement hash conflict resolution |

### FREECODE-006: Background Agents
| Task | Description |
|------|-------------|
| FREECODE-006-T1 | Implement tmux pane manager |
| FREECODE-006-T2 | Implement agent spawner |
| FREECODE-006-T3 | Implement output streaming |
| FREECODE-006-T4 | Implement task cancellation |

### FREECODE-017: Fleet Management
| Task | Description |
|------|-------------|
| FREECODE-017-T1 | Implement head mode (accept agents) |
| FREECODE-017-T2 | Implement agent mode (connect to head) |
| FREECODE-017-T3 | Implement client mode (view/manage) |
| FREECODE-017-T4 | Implement task dispatch |
| FREECODE-017-T5 | Implement file transfer |

### FREECODE-031: TUI Mouse Support
| Task | Description |
|------|-------------|
| FREECODE-031-T1 | Implement region registry |
| FREECODE-031-T2 | Implement mouse event handler |
| FREECODE-031-T3 | Implement cursor rendering |

### FREECODE-037: Session Tabs
| Task | Description |
|------|-------------|
| FREECODE-037-T1 | Implement tab bar |
| FREECODE-037-T2 | Implement tab content rendering |
| FREECODE-037-T3 | Implement tab persistence |

(Full task reference continues...)

---

## Feature Categories

### Agents & Orchestration (12 features)
- Multi-Agent System, Category-Based Routing, Atlas Orchestrator, Ralph Loop, Ultrawork Mode, Parallel Review, Background Agents, Call Agent Tool, Start Work Command, Handoff Command, Init Deep Command, Todo Continuation

### Tools (15 features)
- Hashline Edit, AST-grep, Look_at, Interactive Bash, Session Tools, Task Tools, Background Output, Background Cancel, VNC Tool, RDP Tool, SSH Tunnels, rsync Support, doom Loop Detection, Comment Checker, Explainshell Command

### Commands & Skills (14 features)
- Slash Commands, Commands System, git-master Skill, playwright Skill, frontend-ui-ux Skill, review-work Skill, ai-slop-remover Skill, Agent Browser Skill, Refactor Command, Test Command, Docs Command, Keyword Detector, Preemptive Compaction, File Prompt Loading

### VCS & Integrations (10 features)
- gh Integration, glab Integration, go-git Fallback, SVN Support, Git Master Integration, Websearch Exa, Context7 Docs, Grep.app, Built-in MCPs, BitTorrent

### Network & File (8 features)
- SMB/CIFS, NFS, AFP, BitTorrent, rsync, VNC, RDP, SSH Tunnels

### Fleet (3 features)
- Fleet Management, Fleet Panel, Enterprise Share

### Security (7 features)
- YubiKey SSH, YubiKey GPG, Permission Patterns, External Directory Access, doom Loop Detection, YOLO Mode, MCP OAuth

### UI/UX (7 features)
- TUI Mouse Support, Context Menus, Scroll Regions, Text Selection, Command Palette, Session Tabs, Split View

### Model Management (5 features)
- LSP Support, Extended Providers, Model Fallback Chains, Model Capabilities Cache, Runtime Fallback

### Context Management (4 features)
- Dynamic Pruning, Preemptive Compaction, Lifecycle Hooks (60+), Atlas Orchestrator

---

**Author:** Mark LaPointe <mark@cloudbsd.org>
**Last Updated:** 2026-05-01

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author
