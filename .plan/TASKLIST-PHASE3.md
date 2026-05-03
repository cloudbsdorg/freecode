# Freecode Implementation Task List - Phase 3

**Last Updated:** 2026-05-02
**Author:** Mark LaPointe <mark@cloudbsd.org>
**Build Status:** ✅ Passes
**Phase:** 3 - OpenCode Parity Gap (12 remaining modules)

---

## Overview

This document tracks implementation of 12 remaining modules to achieve full OpenCode parity. Phase 2 (16 modules) is complete. Phase 3 focuses on the remaining modules discovered during gap analysis.

---

## Missing Modules (OpenCode → Freecode)

### High Priority

| # | Module | Description | OpenCode Path | Status |
|---|--------|-------------|---------------|--------|
| 1 | `internal/account` | Account management | `packages/opencode/src/account/` | ❌ Missing |
| 2 | `internal/acp` | Access control policy | `packages/opencode/src/acp/` | ❌ Missing |
| 3 | `internal/control-plane` | Fleet control plane | `packages/opencode/src/control-plane/` | ❌ Missing |
| 4 | `internal/file` | File operations/watcher | `packages/opencode/src/file/` | ❌ Missing |
| 5 | `internal/plugin` | Plugin system | `packages/opencode/src/plugin/` | ❌ Missing |
| 6 | `internal/skill` | Skill system | `packages/opencode/src/skill/` | ⚠️ Partial |

### Medium Priority

| # | Module | Description | OpenCode Path | Status |
|---|--------|-------------|---------------|--------|
| 7 | `internal/env` | Environment variables | `packages/opencode/src/env/` | ❌ Missing |
| 8 | `internal/format` | Code formatting | `packages/opencode/src/format/` | ❌ Missing |
| 9 | `internal/question` | Question/answer flow | `packages/opencode/src/question/` | ❌ Missing |
| 10 | `internal/util` | Utilities | `packages/opencode/src/util/` | ❌ Missing |

### Low Priority

| # | Module | Description | OpenCode Path | Status |
|---|--------|-------------|---------------|--------|
| 11 | `internal/id` | ID generation | `packages/opencode/src/id/` | ❌ Missing |
| 12 | `internal/installation` | Installation detection | `packages/opencode/src/installation/` | ❌ Missing |

---

## Progress Summary

| Priority | Total | Done | Remaining |
|----------|-------|------|----------|
| HIGH | 6 | 0 | 6 |
| MEDIUM | 4 | 0 | 4 |
| LOW | 2 | 0 | 2 |
| **Total** | **12** | **0** | **12** |

---

## Change Log

| Date | Description |
|------|-------------|
| 2026-05-02 | Initial Phase 3 task list |
