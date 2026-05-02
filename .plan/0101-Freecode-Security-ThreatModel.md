# Freecode — Security Threat Model

**Document ID:** Freecode-Security-ThreatModel
**Version:** 1.0
**Last Updated:** 2026-05-02
**Maintainer:** Mark LaPointe <mark@cloudbsd.org>
**Status:** ACTIVE
**Classification:** INTERNAL

---

## Executive Summary

Freecode is a local AI coding assistant that runs entirely on localhost. All network services bind exclusively to 127.0.0.1 and ::1. No telemetry, no external data collection, no third-party dependencies for core functionality.

## Assets to Protect

| Asset | Classification | Protection Required |
|-------|---------------|---------------------|
| User's source code | Confidential | Read access only, no exfiltration |
| API keys / credentials | Secret | Encrypted at rest, never logged |
| Session data | Private | Isolated per-session |
| Configuration | Internal | User-controlled, validated |
| Model responses | Private | Not stored beyond session |

## Threat Categories

| Category | Threats | Mitigation |
|----------|---------|------------|
| Data Exfiltration | Malicious instructions, prompt injection | Local-only operation, no network egress |
| Credential Theft | Keylogging, config reading | Encrypted storage, memory protection |
| Code Corruption | Malicious LLM outputs | Sandbox execution, git backup |
| Resource Exhaustion | Fork bombs, infinite loops | Process limits, timeouts |

## Trust Model

| Level | Component | Trust | Reason |
|-------|-----------|-------|--------|
| T0 | Local CLI execution | Full | User's own machine |
| T1 | LLM API calls | Partial | External service, use HTTPS |
| T2 | Plugin/extension code | Low | User-installed, not verified |
| T3 | Fleet connections | Minimal | Requires explicit auth |

## Isolation Architecture

```
┌─────────────────────────────────────────┐
│           User's Machine                │
│  ┌─────────────────────────────────┐    │
│  │         Freecode CLI             │    │
│  │  ┌─────────────────────────┐    │    │
│  │  │   Session Manager       │    │    │
│  │  │  ┌─────┐ ┌─────┐      │    │    │
│  │  │  │Tab 1│ │Tab 2│ ...  │    │    │
│  │  │  └─────┘ └─────┘      │    │    │
│  │  └─────────────────────────┘    │    │
│  │  ┌─────────────────────────┐    │    │
│  │  │   Tool Executor         │    │    │
│  │  │  (sandboxed per tool)  │    │    │
│  │  └─────────────────────────┘    │    │
│  └─────────────────────────────────┘    │
│         │          │           │         │
│         ▼          ▼           ▼         │
│    ┌────────┐ ┌────────┐ ┌────────┐     │
│    │ 127.0.0.1 │ │ 127.0.0.1 │ │ 127.0.0.1 │     │
│    │   :18792  │ │   :18793  │ │   :18791  │     │
│    │   (API)   │ │   (MCP)   │ │   (Web)   │     │
│    └────────┘ └────────┘ └────────┘     │
└─────────────────────────────────────────┘
         NO EXTERNAL NETWORK
```

## Process Isolation

- Each tool execution runs in isolated subprocess
- Memory limits enforced via OS constraints
- No persistent state between tool calls
- Filesystem access via whitelisted paths only

## Multi-Session Isolation

- Sessions are completely isolated
- Each tab operates independently
- No cross-session state sharing
- Fleet mode uses explicit authentication

---

## Change Log

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-05-02 | Mark LaPointe | Initial threat model |

**Last Updated:** 2026-05-02 07:30 UTC
**Classification:** INTERNAL
