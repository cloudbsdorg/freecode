# Freecode Implementation Task List - Phase 3

**⚠️ DEPRECATED:** See [TASKLIST.md](./TASKLIST.md) and [FREECODE-STATUS.md](./FREECODE-STATUS.md) for accurate status.

**Last Updated:** 2026-05-02
**Author:** Mark LaPointe <mark@cloudbsd.org>
**Build Status:** ✅ Passes
**Phase:** 3 - OpenCode Parity Gap (12 remaining modules)

---

## Overview

> **⚠️ DEPRECATED (2026-05-04):** Phase 2 is NOT complete - the modules are stubs. This document also incorrectly marked modules as "Done" when they are actually stubs or missing. True status is documented in FREECODE-STATUS.md.

This document tracks implementation of 12 remaining modules to achieve full OpenCode parity. Phase 2 (16 modules) claimed complete but is actually ~6% complete. Phase 3 focuses on the remaining modules discovered during gap analysis.

---

## Missing Modules (OpenCode → Freecode)

### High Priority

| # | Module | Description | OpenCode Path | Status |
|---|--------|-------------|---------------|--------|
| 1 | `internal/account` | Account management | `packages/opencode/src/account/` | ✅ Done |
| 2 | `internal/acp` | Access control policy | `packages/opencode/src/acp/` | ✅ Done |
| 3 | `internal/control-plane` | Fleet control plane | `packages/opencode/src/control-plane/` | ✅ Done |
| 4 | `internal/file` | File operations/watcher | `packages/opencode/src/file/` | ✅ Done |
| 5 | `internal/plugin` | Plugin system | `packages/opencode/src/plugin/` | ✅ Done |
| 6 | `internal/skill` | Skill system | `packages/opencode/src/skill/` | ✅ Done |

### Medium Priority

| # | Module | Description | OpenCode Path | Status |
|---|--------|-------------|---------------|--------|
| 7 | `internal/env` | Environment variables | `packages/opencode/src/env/` | ✅ Done |
| 8 | `internal/format` | Code formatting | `packages/opencode/src/format/` | ✅ Done |
| 9 | `internal/question` | Question/answer flow | `packages/opencode/src/question/` | ✅ Done |
| 10 | `internal/util` | Utilities | `packages/opencode/src/util/` | ✅ Done |

### Low Priority

| # | Module | Description | OpenCode Path | Status |
|---|--------|-------------|---------------|--------|
| 11 | `internal/id` | ID generation | `packages/opencode/src/id/` | ✅ Done |
| 12 | `internal/installation` | Installation detection | `packages/opencode/src/installation/` | ✅ Done |

---

## Progress Summary

| Priority | Total | Truly Complete | Remaining |
|----------|-------|-----------------|----------|
| HIGH | 6 | 0 | 6 (all stubs/missing) |
| MEDIUM | 4 | 1 (env) | 3 |
| LOW | 2 | 2 (id, installation) | 0 |
| **Total** | **12** | **3** | **9** |

> True completion: 25%

---

## Change Log

| Date | Description |
|------|-------------|
| 2026-05-02 | Initial Phase 3 task list |
