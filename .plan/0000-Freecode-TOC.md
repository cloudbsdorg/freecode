# Freecode — Master Table of Contents

**Document ID:** Freecode-TOC
**Version:** 2.0
**Last Updated:** 2026-05-02
**Maintainer:** Mark LaPointe <mark@cloudbsd.org>
**Status:** ACTIVE
**Classification:** INTERNAL

---

## Document Map

### Meta (0000-0002)

| File | Title | Status | Description |
|------|-------|--------|-------------|
| `0000-Freecode-TOC.md` | Master Table of Contents | ✅ ACTIVE | This document |
| `0001-Freecode-Workflow.md` | Workflow | ✅ ACTIVE | Task claiming, completion, merge handling |
| `0002-Freecode-Build-Status.md` | Build Status | ✅ ACTIVE | CI/CD pipeline, build artifacts |

### Security (0100-0106)

| File | Title | Status | Description |
|------|-------|--------|-------------|
| `0100-Freecode-Security-Overview.md` | Security Overview | ✅ ACTIVE | Security strategy summary |
| `0101-Freecode-Security-ThreatModel.md` | Threat Model | ✅ ACTIVE | Threat analysis, trust model |
| `0102-Freecode-Security-AccessControl.md` | Access Control | ✅ ACTIVE | Permission tiers, credential storage |
| `0103-Freecode-Security-Runtime.md` | Runtime Security | ✅ ACTIVE | Sandbox, filesystem, crash containment |
| `0104-Freecode-Security-Implementation.md` | Security Tasks | ✅ ACTIVE | Implementation checklist |
| `0105-Freecode-Security-Audit.md` | Security Audit | ✅ ACTIVE | oh-my-openagent security analysis |
| `0106-Freecode-Security-Additional.md` | Additional Analysis | ✅ ACTIVE | Audit logging, supply chain |

### Overview & Architecture (0200-0213)

| File | Title | Status | Description |
|------|-------|--------|-------------|
| `0200-Freecode-Overview.md` | Overview | ✅ ACTIVE | Executive summary, phases |
| `0201-Freecode-Current-Architecture.md` | Current Architecture | ✅ ACTIVE | TypeScript → Go mapping |
| `0202-Freecode-Platform-Specific.md` | Platform-Specific | ✅ ACTIVE | FreeBSD, macOS, Linux, IllumOS |
| `0203-Freecode-Feature-Inventory.md` | Feature Inventory | ✅ ACTIVE | 88 features tracked |
| `0204-Freecode-Features.md` | Features | ✅ ACTIVE | All features: agents, hooks, MCP |
| `0210-Freecode-Architecture-Design.md` | Architecture Design | ✅ ACTIVE | Go architecture, packages |
| `0211-Freecode-LiteLLM-Integration.md` | LiteLLM Integration | ✅ ACTIVE | Provider consolidation |
| `0212-Freecode-TUI-Analysis.md` | TUI Analysis | ✅ ACTIVE | OpenCode vs Freecode TUI parity |
| `0213-Freecode-Missing-Features.md` | Missing Features | ✅ ACTIVE | Gap analysis vs opencode |

### Implementation (0300-0301)

| File | Title | Status | Description |
|------|-------|--------|-------------|
| `0300-Freecode-Implementation-Tasks.md` | Implementation Tasks | ✅ ACTIVE | Phase-by-phase breakdown |
| `0301-Freecode-Session-Tabbing.md` | Session Tabbing | ✅ ACTIVE | TUI tabs, split view |

### Testing (0400-0403)

| File | Title | Status | Description |
|------|-------|--------|-------------|
| `0400-Freecode-Testing.md` | Testing Overview | ✅ ACTIVE | Test strategy |
| `0401-Freecode-Unit-Tests.md` | Unit Tests | ✅ ACTIVE | Unit testing plan |
| `0402-Freecode-Integration-Tests.md` | Integration Tests | ✅ ACTIVE | Integration testing plan |
| `0403-Freecode-Code-Validation.md` | Code Validation | ✅ ACTIVE | Linting, fuzzing, security |

### Operations (0500-0504)

| File | Title | Status | Description |
|------|-------|--------|-------------|
| `0501-Freecode-Configuration.md` | Configuration | ✅ ACTIVE | Config schema, migration |
| `0502-Freecode-Packaging.md` | Packaging | ✅ ACTIVE | FreeBSD, Linux, macOS, IllumOS |
| `0503-Freecode-Dependencies.md` | Dependencies | ✅ ACTIVE | Build dependencies |
| `0504-Freecode-I18N.md` | Internationalization | ✅ ACTIVE | i18n support |
| `0510-Freecode-Tooling.md` | Tooling | ✅ ACTIVE | Development guide |

### Risks (0700)

| File | Title | Status | Description |
|------|-------|--------|-------------|
| `0700-Freecode-Risks.md` | Risks | ✅ ACTIVE | Risk register |

### Validation (0900)

| File | Title | Status | Description |
|------|-------|--------|-------------|
| `0900-Freecode-Validation.md` | Validation | ✅ ACTIVE | Task completion, validation checklists |

---

## Document Dependencies

```
0000 (TOC) ──┬── 0001 (Workflow)
             └── 0002 (Build Status)
                      │
                      ├──► 0100 (Security Overview)
                      │         │
                      │         ├──► 0101 (Threat Model)
                      │         ├──► 0102 (Access Control)
                      │         ├──► 0103 (Runtime Security)
                      │         ├──► 0104 (Security Implementation)
                      │         ├──► 0105 (Security Audit)
                      │         └──► 0106 (Security Additional)
                      │
                       ├──► 0200 (Overview)
                       │         │
                       │         ├──► 0201 (Current Architecture)
                       │         ├──► 0202 (Platform-Specific)
                       │         ├──► 0203 (Feature Inventory)
                       │         ├──► 0204 (Features)
                       │         ├──► 0210 (Architecture Design)
                       │         ├──► 0211 (LiteLLM Integration)
                       │         ├──► 0212 (TUI Analysis)
                       │         └──► 0213 (Missing Features)
                      │
                      ├──► 0300 (Implementation Tasks)
                      │         │
                      │         └──► 0301 (Session Tabbing)
                      │
                      ├──► 0400 (Testing)
                      │         ├──► 0401 (Unit Tests)
                      │         ├──► 0402 (Integration Tests)
                      │         └──► 0403 (Code Validation)
                      │
                      ├──► 0501 (Configuration)
                      │         ├──► 0502 (Packaging)
                      │         ├──► 0503 (Dependencies)
                      │         ├──► 0504 (I18N)
                      │         └──► 0510 (Tooling)
                      │
                      ├──► 0700 (Risks)
                      │
                      └──► 0900 (Validation)

Legend: ──┬── = references, ──► = depends on
```

---

## Reading Order

For new contributors:

1. **[`AGENTS_START_HERE.md`](../AGENTS_START_HERE.md)** — Read this first
2. **[`0000-Freecode-TOC.md`](./0000-Freecode-TOC.md)** — Master index
3. **[`0001-Freecode-Workflow.md`](./0001-Freecode-Workflow.md)** — How to work
4. **[`0200-Freecode-Overview.md`](./0200-Freecode-Overview.md)** — The big picture
5. **[`0100-Freecode-Security-Overview.md`](./0100-Freecode-Security-Overview.md)** — Security strategy
6. **[`0201-Freecode-Current-Architecture.md`](./0201-Freecode-Current-Architecture.md)** — Current state
7. **[`0210-Freecode-Architecture-Design.md`](./0210-Freecode-Architecture-Design.md)** — Target design
8. **[`0213-Freecode-Missing-Features.md`](./0213-Freecode-Missing-Features.md)** — Gap analysis vs opencode
9. **[`0300-Freecode-Implementation-Tasks.md`](./0300-Freecode-Implementation-Tasks.md)** — What to build

---

## Quick Reference

### Key Directory Structure

```
freecode/
├── cmd/
│   ├── freecode/              # Main CLI entry
│   └── freecode-server/       # Server mode
├── internal/
│   ├── agent/                 # 11 built-in agents (prompts done)
│   ├── auth/                  # Credential storage
│   ├── cli/                   # 27 CLI commands (account, web done)
│   ├── config/                # Configuration
│   ├── hook/                  # 60+ lifecycle hooks (26+9 done)
│   ├── mcp/                   # MCP protocol
│   ├── provider/              # 50+ providers
│   ├── session/               # Session management
│   ├── shell/                 # Shell integration
│   ├── tool/                  # Tool registry
│   ├── ui/                    # Bubble Tea TUI
│   └── fleet/                 # Fleet management
├── .skills/                   # Skills (7 defined)
├── .plan/                     # Plan documents
├── make/                      # Build makefiles
└── go.mod
```

### Build Commands

```bash
# Build CLI
go build -o freecode ./cmd/freecode

# Build server
go build -o freecode-server ./cmd/freecode-server

# Or use make
make build

# Run tests
make test

# Run all quality gates
make test && make fmt && make tidy
```

---

## Change Log

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 2.2 | 2026-05-02 | Mark LaPointe | Add skills, update status (agents prompts done, hooks done, account/web done) |
| 2.1 | 2026-05-02 | Mark LaPointe | Add 0213 Missing Features; update 0204 agents as stubs, 0212 TUI Analysis, 0300 tasks |
| 2.0 | 2026-05-02 | Mark LaPointe | Migrate to CloudBSD 4-digit numbering |
| 1.0 | 2026-05-01 | Mark LaPointe | Initial TOC |

---

**Author:** Mark LaPointe <mark@cloudbsd.org>

**Last Updated:** 2026-05-02

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author
