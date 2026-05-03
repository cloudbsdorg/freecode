# Skill: TASKLIST.md Creation Workflow

## Purpose

Create comprehensive TASKLIST.md files that follow a proper planning-before-implementation workflow.

## Core Principle

**DOCUMENT FIRST, IMPLEMENT LAST**

Never start implementing until all planning documents are updated.

---

## Workflow Steps

### Step 1: Discovery Phase

1. **Audit Current State**
   - List all existing files/modules in the project
   - Compare against reference implementation (e.g., opencode)
   - Identify what's MISSING

2. **Categorize Missing Items**
   - HIGH priority: Blockers, foundational
   - MEDIUM priority: Feature parity
   - LOW priority: Nice-to-have

3. **Determine Dependencies**
   - Which modules depend on which?
   - What must be implemented first?

### Step 2: Documentation Phase

1. **Update Existing Plan Documents**
   - Add new tasks to 0213 (Missing Features)
   - Add new phase to 0300 (Implementation Tasks)
   - Update 0204 (Features) if needed

2. **Create TASKLIST.md**
   - Separate file for detailed tracking
   - Include status table with columns:
     - Task name
     - File path
     - Status (❌ Missing → ⏳ Planned → ✅ Done)
     - Dependencies
     - Notes

3. **Create CHAIN-TODO.md** (optional)
   - Ordered sequence of implementation
   - Dependency chain visualization

### Step 3: Validation Phase

1. **Verify All Files Exist**
   - Do NOT assume - use `ls` or `test -f`
   - Check each planned file

2. **Run Tests**
   - Unit tests must pass
   - Integration tests if applicable

3. **Build Check**
   - `go build ./...` must pass

### Step 4: Implementation Phase (LAST)

1. **Implement in Dependency Order**
   - Never implement a module that another unimplemented module depends on

2. **Mark Complete Immediately**
   - Update TASKLIST.md as you go
   - Commit after each module

3. **Validate After Each**
   - Build after each file
   - Test after each module

---

## TASKLIST.md Template

```markdown
# Project Implementation Task List

**Last Updated:** YYYY-MM-DD
**Author:** Name <email>
**Build Status:** ✅ Passes / ❌ Fails

---

## Overview

Brief description of what this task list covers.

---

## Priority Chain

### 🔴 HIGH PRIORITY

| # | Task | File | Status | Dependencies | Notes |
|---|------|------|--------|-------------|-------|
| 1 | Task name | path/to/file.go | ❌ Missing | dep1, dep2 | Description |

### 🟡 MEDIUM PRIORITY

| # | Task | File | Status | Dependencies | Notes |
|---|------|------|--------|-------------|-------|
| 10 | Task name | path/to/file.go | ❌ Missing | dep1 | Description |

### 🟢 LOW PRIORITY

| # | Task | File | Status | Dependencies | Notes |
|---|------|------|--------|-------------|-------|
| 20 | Task name | path/to/file.go | ❌ Missing | dep1 | Description |

---

## Progress Summary

| Priority | Total | Done | Remaining |
|----------|-------|------|----------|
| HIGH | X | 0 | X |
| MEDIUM | X | 0 | X |
| LOW | X | 0 | X |
| **Total** | **X** | **0** | **X** |

---

## Change Log

| Date | Commit | Description |
|------|--------|-------------|
| YYYY-MM-DD | - | Initial task list |
```

---

## Status Values

| Status | Meaning |
|--------|---------|
| ❌ Missing | Does not exist |
| ⏳ Planned | Will create |
| ⚠️ Stub | Exists but incomplete |
| ✅ Done | Implemented and tested |

---

## Anti-Patterns

❌ **NEVER** start implementing before documentation is complete
❌ **NEVER** mark something done that doesn't compile
❌ **NEVER** skip the validation phase

---

## Files Created by This Skill

- `.skills/tasklist-workflow/SKILL.md` - This file
- `.plan/TASKLIST-XXXX.md` - Main task list
- `.plan/CHAIN-TODO.md` - Ordered implementation chain (optional)

---

## Author

Mark LaPointe <mark@cloudbsd.org>

All commits authored solely by Mark LaPointe. No co-authors.
