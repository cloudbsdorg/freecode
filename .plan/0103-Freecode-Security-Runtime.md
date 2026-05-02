# Freecode — Security Runtime

**Document ID:** Freecode-Security-Runtime
**Version:** 1.0
**Last Updated:** 2026-05-02
**Maintainer:** Mark LaPointe <mark@cloudbsd.org>
**Status:** ACTIVE
**Classification:** INTERNAL

---

## Filesystem Security

### Path Validation

All file paths are validated before access:

- **Realpath validation** — Resolves symlinks, ensures no `../` escapes
- **Allowlist prefixes** — Only configured directories are accessible
- **Temporary file isolation** — All temp files in `$TMPDIR/freecode-*`

### Path Restrictions

```yaml
security:
  allowed_paths:
    - "$HOME"           # User's home directory
    - "/tmp"           # Temporary files
    - "$WORKSPACE"     # Current working directory

  blocked_paths:
    - "/etc/sudoers"   # System credentials
    - "$HOME/.ssh"     # SSH keys (explicit read-only needed)
    - "$HOME/.aws"     # Cloud credentials
```

### Tool Execution Sandbox

Each tool runs in an isolated subprocess with:

| Constraint | Value | Purpose |
|------------|-------|---------|
| Max memory | 512MB | Prevent OOM attacks |
| Max CPU | 30s | Prevent infinite loops |
| Max files | 100 | Prevent file descriptor exhaustion |
| Network | none | No network access during execution |

## Service Security

### Localhost Binding

All services bind exclusively to localhost:

| Service | Port | Bind Address |
|---------|------|--------------|
| API Server | 18792 | 127.0.0.1, ::1 |
| MCP Server | 18793 | 127.0.0.1, ::1 |
| Web UI | 18791 | 127.0.0.1, ::1 |
| Fleet Head | 7842 | 0.0.0.0 (opt-in) |

### Fleet Exposure

Fleet mode can be exposed to LAN for multi-machine coordination:

```bash
# Expose fleet on LAN (requires TLS + auth)
./freecode serve --fleet --fleet-listen 0.0.0.0:7842 --fleet-tls

# Verify binding
netstat -an | grep 7842
```

## Crash Containment

### Process Limits

| Resource | Limit | Enforcement |
|----------|-------|-------------|
| Max child processes | 10 | OS rlimit |
| Max open files | 256 | OS rlimit |
| Max memory | 1GB | OS cgroup |
| Max CPU time | 5min | OS scheduler |

### Graceful Degradation

On crash detection:
1. Terminate all child processes
2. Save session state to disk
3. Log incident (no sensitive data)
4. Restart in clean state

---

## Change Log

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-05-02 | Mark LaPointe | Initial runtime security document |

**Last Updated:** 2026-05-02 07:30 UTC
**Classification:** INTERNAL
