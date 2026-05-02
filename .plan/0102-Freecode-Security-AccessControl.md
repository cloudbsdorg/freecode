# Freecode — Security Access Control

**Document ID:** Freecode-Security-AccessControl
**Version:** 1.0
**Last Updated:** 2026-05-02
**Maintainer:** Mark LaPointe <mark@cloudbsd.org>
**Status:** ACTIVE
**Classification:** INTERNAL

---

## Permission System

Freecode implements a tiered permission system for tool execution.

### Permission Tiers

| Tier | Name | Description | Example Tools |
|------|------|-------------|---------------|
| P0 | Critical | System-level operations | `rm -rf`, `kill`, `sudo` |
| P1 | High | File system modifications | `write`, `edit`, `delete` |
| P2 | Medium | File system reading | `read`, `grep`, `glob` |
| P3 | Low | Read-only information | `stats`, `list`, `search` |
| P4 | Minimal | No side effects | `echo`, `date`, `version` |

### Default Permissions by Agent

| Agent | Default Tier | Can Override |
|-------|-------------|--------------|
| Sisyphus (primary) | P1 | Yes, with confirmation |
| Hephaestus | P2 | Yes, with confirmation |
| Oracle | P3 | No |
| Librarian | P3 | No |
| Explore | P2 | Yes, with confirmation |
| Prometheus | P1 | Yes, with confirmation |
| Metis | P3 | No |
| Momus | P3 | No |
| Atlas | P2 | Yes, with confirmation |
| Multimodal-Looker | P3 | No |
| Sisyphus-Junior | P3 | No |

### YOLO Mode

When YOLO mode is enabled, all permission checks are bypassed:

```bash
# Enable YOLO mode
export FREECODE_YOLO=true
./freecode

# Or via hotkey during session
Ctrl+Y
```

**Warning:** YOLO mode skips all confirmations. Use only in trusted environments.

## Credential Storage

### Auth Store

Credentials are stored encrypted using AES-256-GCM:

```
internal/auth/
├── store.go          # Encryption/decryption
├── credentials.go    # Credential types
└── keyring/          # Platform keyring integration
```

### Supported Keyring Backends

| Platform | Backend |
|----------|---------|
| macOS | Keychain |
| Linux | libsecret |
| FreeBSD | gnome-keyring or pass |
| Windows | Credential Manager |

### Credential Types

| Type | Storage | Transmission |
|------|---------|--------------|
| API Keys | Encrypted | Never logged |
| Tokens | Encrypted | HTTPS only |
| Passwords | Encrypted | Never plaintext |

---

## Change Log

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-05-02 | Mark LaPointe | Initial access control document |

**Last Updated:** 2026-05-02 07:30 UTC
**Classification:** INTERNAL
