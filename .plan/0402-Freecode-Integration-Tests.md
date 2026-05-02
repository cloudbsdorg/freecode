# Freecode — Integration Tests

**Document ID:** Freecode-IntegrationTests
**Version:** 1.0
**Last Updated:** 2026-05-02
**Maintainer:** Mark LaPointe <mark@cloudbsd.org>
**Status:** ACTIVE
**Classification:** INTERNAL

---

## End-to-End Scenarios

### Full Lifecycle Tests

| Scenario | Steps | Validation |
|----------|-------|------------|
| CLI startup | `freecode` → prompt | TUI renders |
| Session creation | New tab | Tab appears |
| Tool execution | Run `ls` | Output shown |
| Session export | `freecode session export` | File created |
| Session import | `freecode session import` | Sessions restored |

### Inter-Component Workflows

| Workflow | Components | Validation |
|----------|------------|------------|
| Provider auth | config → auth → provider | Token stored |
| Tool execution | cli → session → tool → output | Result correct |
| Fleet connection | head → agent → client | Tasks assigned |

## Performance and Stress

### Load Testing

| Metric | Target | Method |
|--------|--------|--------|
| Startup time | <2s | Benchmark |
| Memory baseline | <50MB | RSS check |
| Memory max | <200MB | Extended session |
| Concurrent tabs | 10+ | Stress test |

### Longevity Testing

- 24h+ session stability
- Memory leak detection
- Log rotation verification

## Network and Environment

### Test Topology

```
┌─────────────────────────────────────────┐
│         Integration Test Host          │
│  ┌─────────────────────────────────┐    │
│  │        Freecode Instance        │    │
│  │  ┌─────────────────────────┐   │    │
│  │  │   API Server :18792     │   │    │
│  │  │   MCP Server :18793     │   │    │
│  │  │   Web UI    :18791      │   │    │
│  │  └─────────────────────────┘   │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘
              │              │
              ▼              ▼
        ┌──────────┐  ┌──────────┐
        │ Head     │  │ Agent    │
        │ :7842    │  │ :7843    │
        └──────────┘  └──────────┘
```

### External Dependencies

| Dependency | Purpose | Test Strategy |
|------------|---------|----------------|
| OpenAI API | LLM calls | Mock for unit, real for E2E |
| Anthropic API | LLM calls | Mock for unit, real for E2E |
| GitHub API | PR tools | Mock responses |

---

## Change Log

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-05-02 | Mark LaPointe | Initial integration test document |

**Last Updated:** 2026-05-02 07:30 UTC
**Classification:** INTERNAL
