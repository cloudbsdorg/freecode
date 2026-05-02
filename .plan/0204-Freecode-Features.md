# Freecode — Features

## 1.0 Purpose

This document details all features built into freecode. These are native features, not plugins or integrations—freecode is a single unified product.

---

## 2.0 Built-in Agents (11 Total)

**Implementation Status:** Agents are documented here but currently exist as stubs in `internal/agent/sisyphus.go`. Full agent prompting and behavior needs implementation.

### 2.1 Agent Definitions

| Agent | Mode | Default Model | Purpose | Status |
|-------|------|---------------|---------|--------|
| Sisyphus | primary | claude-opus-4-7 | Main orchestrator | Stub |
| Hephaestus | primary | gpt-5.4 | Code generation | Stub |
| Oracle | subagent | gpt-5.4 | Architecture consultation | Stub |
| Librarian | subagent | gpt-5.4-mini-fast | Research/library | Stub |
| Explore | subagent | gpt-5.4-mini-fast | Exploration | Stub |
| Prometheus | all | claude-opus-4-7 | Planning | Stub |
| Metis | all | claude-opus-4-7 | Plan consultation | Stub |
| Momus | all | gpt-5.4 | Code review | Stub |
| Atlas | primary | claude-sonnet-4-6 | Session tracking | Stub |
| Multimodal-Looker | subagent | gpt-5.4 | Multimodal | Stub |
| Sisyphus-Junior | all | (inherited) | Simpler tasks | Stub |

### 2.2 Agent System Prompts

**Sisyphus (Primary):**
```
You are Sisyphus, the eternal orchestrator. You coordinate all other agents
to accomplish complex coding tasks. You delegate to specialized agents while
maintaining overall project coherence. Your work is never done, but you
find satisfaction in the journey.
```

**Hephaestus (Code Generation):**
```
You are Hephaestus, the craftsman. You write clean, efficient code with
attention to detail. You understand the forge and anvil of software
construction - patterns, performance, and reliability.
```

**Oracle (Architecture):**
```
You are the Oracle, keeper of architectural wisdom. You see the big picture
and guide structural decisions. When asked about architecture, you provide
deep analysis of trade-offs and long-term implications.
```

### 2.3 Agent Configuration

```go
type AgentConfig struct {
    Name        string
    Mode        AgentMode // primary|subagent|all
    Description string
    SystemPrompt string
    DefaultModel string
    Tools       []string
    Categories  []string // Model categories to inherit from
}

var BuiltinAgents = map[string]AgentConfig{
    "sisyphus": {
        Name:        "Sisyphus",
        Mode:        AgentModePrimary,
        Description: "Main orchestrator agent",
        SystemPrompt: "You are Sisyphus, the eternal orchestrator...",
        DefaultModel: "claude-opus-4-7",
        Tools:       []string{"bash", "read", "write", "edit", "glob", "grep", "task", "skill"},
    },
    // ... all 11 agents
}
```

---

## 3.0 Lifecycle Hooks (60+ Total)

### 3.1 Hook Tiers

| Tier | Count | Purpose |
|------|-------|---------|
| Session | 26 | Session lifecycle events |
| Tool-Guard | 14 | Tool execution pre/post |
| Transform | 5 | Message/response transformation |
| Continuation | 10 | Auto-continuation logic |
| Skill | 2 | Skill invocation |
| Ralph | 3 | Ralph loop self-reference |

### 3.2 Session Hooks (26)

| Hook | Description |
|------|-------------|
| `session.start` | Fires when session begins |
| `session.end` | Fires when session ends |
| `session.error` | Fires on session error |
| `session.created` | After new session created |
| `session.deleted` | Before session deleted |
| `session.title_changed` | Session title updated |
| `session.renamed` | Session renamed |
| `session.forked` | Session forked |
| `session.merged` | Sessions merged |
| `session.shared` | Session shared |
| `session.exported` | Session exported |
| `session.imported` | Session imported |
| `session.notification` | Notification received |
| `session.error_recovery` | Error recovery attempt |
| `session.context_exhausted` | Context window near limit |
| `session.compaction_start` | Compaction begun |
| `session.compaction_end` | Compaction completed |
| `session.tab_created` | Tab created |
| `session.tab_closed` | Tab closed |
| `session.tab_changed` | Active tab changed |
| `session.message_added` | New message in session |
| `session.turn_end` | Turn completed |
| `session.ultrawork_start` | Ultrawork mode entered |
| `session.ultrawork_end` | Ultrawork mode exited |

### 3.3 Tool Hooks (14)

| Hook | Description |
|------|-------------|
| `tool.execute.before` | Before tool execution |
| `tool.execute.after` | After tool execution |
| `tool.execute.error` | Tool execution error |
| `tool.execute.timeout` | Tool timeout |
| `tool.register` | Tool registered |
| `tool.unregister` | Tool unregistered |
| `tool.schema` | Tool schema accessed |
| `tool.validate` | Tool input validated |
| `tool.output.truncate` | Tool output truncated |
| `tool.output.expand` | Tool output expanded |
| `tool.confirm` | Tool confirmation requested |
| `tool.confirm.deny` | Tool confirmation denied |
| `tool.confirm.allow` | Tool confirmation allowed |
| `tool.rate_limit` | Tool rate limited |

### 3.4 Transform Hooks (5)

| Hook | Description |
|------|-------------|
| `transform.message` | Transform message content |
| `transform.response` | Transform response |
| `transform.prompt` | Transform prompt before send |
| `transform.tool_result` | Transform tool result |
| `transform.error` | Transform error response |

### 3.5 Continuation Hooks (10)

| Hook | Description |
|------|-------------|
| `continuation.auto` | Auto-continue enabled |
| `continuation.manual` | Manual continue |
| `continuation.think` | Thinking mode |
| `continuation.until_done` | Continue until complete |
| `continuation.max_turns` | Max turns reached |
| `continuation.idle_timeout` | Idle timeout |
| `continuation.budget_exhausted` | Budget exhausted |
| `continuation.ralph_loop` | Ralph self-reference loop |
| `continuation.preemptive_compaction` | Compaction before limit |
| `continuation.atlas` | Master orchestrator for boulder |

### 3.6 Ralph Loop Hooks (3)

| Hook | Description |
|------|-------------|
| `ralph.think` | Self-review thought triggered |
| `ralph.step_back` | Agent stepping back to re-plan |
| `ralph.detected` | Self-contradiction detected |

### 3.7 Skill Hooks (2)

| Hook | Description |
|------|-------------|
| `skill.invoked` | Skill invoked |
| `skill.completed` | Skill completed |

### 3.8 Additional Hooks

| Hook | Description |
|------|-------------|
| `todo.continuation_enforcer` | Forces todo completion |
| `keyword.detector` | Ultrawork/search/analyze mode |
| `preemptive.compaction` | Context before limit hit |

### 3.7 Hook Registry

```go
type HookRegistry struct {
    mu            sync.RWMutex
    sessionHooks  map[string][]SessionHook
    toolHooks     map[string][]ToolHook
    transforms    map[string][]TransformHook
    continuations map[string][]ContinuationHook
    skillHooks    map[string][]SkillHook
    disabled      map[string]bool
}

type SessionHook func(ctx context.Context, evt SessionEvent) error
type ToolHook func(ctx context.Context, evt ToolEvent) (error, bool) // bool=handled
type TransformHook func(msg *Message) (*Message, error)
type ContinuationHook func(ctx context.Context, session *Session) (*ContinueSignal, error)
type SkillHook func(ctx context.Context, skill string, args map[string]any) error

func (r *HookRegistry) Register(hookType, name string, hook interface{}) error
func (r *HookRegistry) Unregister(hookType, name string) error
func (r *HookRegistry) Disable(hookType, name string) error
func (r *HookRegistry) Enable(hookType, name string) error
```

### 3.9 Skills System ✅ DONE

Freecode includes a skills system for specialized task execution. Skills are defined in `.skills/` directory.

| Skill | Description | Location |
|-------|-------------|----------|
| git-master | Expert git operations, history analysis, bisect | `.skills/git-master/SKILL.md` |
| playwright | Browser automation, testing, web scraping | `.skills/playwright/SKILL.md` |
| frontend-ui-ux | UI development, accessibility, design | `.skills/frontend-ui-ux/SKILL.md` |
| review-work | Code review, security audit, quality | `.skills/review-work/SKILL.md` |
| ai-slop-remover | Detect and fix AI-generated code patterns | `.skills/ai-slop-remover/SKILL.md` |
| search-code | Code search with grep, ast-grep, LSP | `.skills/search-code/SKILL.md` |
| architect | System design, architecture patterns | `.skills/architect/SKILL.md` |

**Usage:**
```go
task(
    category="visual-engineering",
    load_skills=["frontend-ui-ux", "playwright"],
    prompt="Build a login form with tests..."
)
```

**Skill Format:**
```markdown
---
name: skill-name
description: One-line description
---

# Skill Title

Detailed skill content...
```

---

## 4.0 Configuration Categories (8)

### 4.1 Category Definitions

| Category | Default Model | Purpose |
|----------|---------------|---------|
| `visual-engineering` | gpt-5.4 | UI/UX work |
| `ultrabrain` | claude-opus-4-7 | Deep reasoning |
| `deep` | claude-opus-4-7 | Complex analysis |
| `artistry` | gpt-5.4 | Creative work |
| `quick` | gpt-5.4-mini-fast | Fast tasks |
| `unspecified-low` | (inherit) | Low priority |
| `unspecified-high` | (inherit) | High priority |
| `writing` | gpt-5.4 | Writing tasks |

### 4.2 Category Config

```go
type CategoryConfig struct {
    Description      string            `yaml:"description"`
    Model            string            `yaml:"model"`
    FallbackModels   []string          `yaml:"fallback_models"`
    Variant          string            `yaml:"variant"`
    Temperature      float64           `yaml:"temperature"`
    TopP             float64           `yaml:"top_p"`
    MaxTokens        int               `yaml:"maxTokens"`
    Thinking         ThinkingConfig    `yaml:"thinking"`
    ReasoningEffort  string            `yaml:"reasoningEffort"`
    TextVerbosity    string            `yaml:"textVerbosity"`
    Tools            map[string]bool   `yaml:"tools"`
    PromptAppend     string            `yaml:"prompt_append"`
    MaxPromptTokens  int               `yaml:"max_prompt_tokens"`
    IsUnstableAgent  bool              `yaml:"is_unstable_agent"`
    Disable          bool              `yaml:"disable"`
}
```

---

## 5.0 Background Task System

### 5.1 Background Task Config

```go
type BackgroundTaskConfig struct {
    DefaultConcurrency int                   `yaml:"defaultConcurrency"`
    ProviderConcurrency map[string]int       `yaml:"providerConcurrency"`
    ModelConcurrency   map[string]int       `yaml:"modelConcurrency"`
    MaxDepth           int                   `yaml:"maxDepth"`
    StaleTimeoutMs     int                  `yaml:"staleTimeoutMs"`
    MessageStalenessTimeoutMs int            `yaml:"messageStalenessTimeoutMs"`
    TaskTtlMs          int                  `yaml:"taskTtlMs"`
    SessionGoneTimeoutMs int                `yaml:"sessionGoneTimeoutMs"`
    SyncPollTimeoutMs  int                  `yaml:"syncPollTimeoutMs"`
    MaxToolCalls       int                  `yaml:"maxToolCalls"`
    CircuitBreaker     CircuitBreakerConfig `yaml:"circuitBreaker"`
}
```

### 5.2 Circuit Breaker

```go
type CircuitBreakerConfig struct {
    Enabled            bool `yaml:"enabled"`
    MaxToolCalls       int  `yaml:"maxToolCalls"`
    ConsecutiveThreshold int `yaml:"consecutiveThreshold"`
}
```

---

## 6.0 Tmux Integration

### 6.1 Tmux Config

```go
type TmuxConfig struct {
    Enabled          bool   `yaml:"enabled"`
    Layout           string `yaml:"layout"` // main-vertical|horizontal|tiled|even-horizontal|even-vertical
    MainPaneSize     int    `yaml:"main_pane_size"` // 20-80
    MainPaneMinWidth int    `yaml:"main_pane_min_width"`
    AgentPaneMinWidth int    `yaml:"agent_pane_min_width"`
    Isolation        string `yaml:"isolation"` // inline|window|session
}
```

---

## 7.0 MCP Servers

### 7.1 MCP Overview

MCP (Model Context Protocol) allows freecode to connect to external tools and services.

### 7.2 Built-in MCP Servers

| MCP | Purpose | Auth |
|-----|---------|------|
| `stdio` | Standard I/O subprocess | None |
| `http` | HTTP REST API | API Key/Bearer |
| `sse` | Server-Sent Events | API Key/Bearer |

### 7.3 MCP Config

```go
type MCPConfig struct {
    Enabled bool `yaml:"enabled"`
    Command string `yaml:"command"` // e.g., "npx", "python"
    Args    []string `yaml:"args"`
    Env     map[string]string `yaml:"env"`
    URL     string `yaml:"url"` // For HTTP/SSE
    Headers map[string]string `yaml:"headers"`
}
```

### 7.4 MCP Example

```yaml
mcp:
  # Built-in websearch
  websearch:
    enabled: true

  # Custom MCP server
  myServer:
    command: "npx"
    args: ["-y", "@myserver/mcp"]
    env:
      API_KEY: "${MY_API_KEY}"
```

### 7.5 Git MCP (go-git fallback)

Git operations use system `git` by default. If system git is not found, freecode can use `go-git` (pure Go implementation) as a fallback convenience.

**Auto-detection:** Freecode checks for `git` on startup. If found, uses system git. If not found and `go-git` is enabled, uses go-git MCP.

```go
type GitConfig struct {
    Enabled   bool   `yaml:"enabled"`
    UseGoGit  bool   `yaml:"useGoGit"`  // Use go-git when system git not found
    GoGitPath string `yaml:"goGitPath"` // Path to go-git binary (default: "go-git")
}

type MCPServer struct {
    Command string            `yaml:"command"`
    Args    []string          `yaml:"args"`
    Env     map[string]string `yaml:"env"`
}

func (g GitConfig) DetectGit() (MCPServer, error) {
    // Check system git
    if _, err := exec.LookPath("git"); err == nil {
        return MCPServer{Command: "git"}, nil
    }

    // Fallback to go-git if enabled
    if g.UseGoGit {
        return MCPServer{
            Command: g.GoGitPath,
            Args:    []string{"mcp"},
        }, nil
    }

    return MCPServer{}, fmt.Errorf("git not found and go-git fallback disabled")
}
```

**go-git MCP Server:**
```yaml
mcp:
  git:
    enabled: true
    useGoGit: true        # Enable fallback to go-git
    goGitPath: "go-git"   # Command or path to go-git binary
```

**Why use go-git fallback?**
- Useful on systems without git installed (e.g., minimal containers)
- FreeBSD base system doesn't include git by default
- Enables git operations in restricted environments
- Pure Go - no external dependencies

**go-git vs system git:**
| Feature | System git | go-git |
|---------|------------|--------|
| Full repo support | Yes | Yes |
| LFS | Yes | No |
| git hooks | Yes | No |
| diff3 merge | Yes | No |
| ssh remotes | Yes | Limited |

### 7.6 VCS MCP Servers

VCS tools can be exposed as MCP servers for agent access:

```yaml
mcp:
  # Git (default: system git, fallback: go-git)
  git:
    enabled: true
    useGoGit: true

  # Subversion (CLI-based, FreeBSD workflow)
  svn:
    command: "svn"
    args: ["mcp"]
    # No pure Go implementation, uses CLI

  # Mercurial (CLI-based)
  hg:
    command: "hg"
    args: ["serve", "--stdio"}
```

### 7.7 File Transfer MCP Servers

File transfer protocols can be exposed as MCP servers for agent file access:

```yaml
mcp:
  # SMB/CIFS file system
  smb:
    command: "npx"
    args: ["-y", "@aws/aws-smb-mcp"]
    env:
      SMB_HOST: "${SMB_HOST}"
      SMB_USER: "${SMB_USER}"
      SMB_PASS: "${SMB_PASS}"

  # NFS file system
  nfs:
    command: "npx"
    args: ["-y", "@aws/aws-nfs-mcp"]
    env:
      NFS_HOST: "${NFS_HOST}"
      NFS_EXPORT: "${NFS_EXPORT}"

  # AFP (Apple Filing Protocol) for macOS
  afp:
    command: "npx"
    args: ["-y", "@company/afp-mcp"]
    env:
      AFP_HOST: "${AFP_HOST}"
      AFP_USER: "${AFP_USER}"
      AFP_PASS: "${AFP_PASS}"

  # rsync for file synchronization
  rsync:
    command: "python"
    args: ["-m", "rsync_mcp"]
    env:
      RSYNC_SSH_KEY: "${RSYNC_SSH_KEY}"

  # BitTorrent for distributed file transfer
  bittorrent:
    command: "python"
    args: ["-m", "bittorrent_mcp"]
    env:
      BT_SEED_DIR: "${BT_SEED_DIR}"
      BT_PORT: "${BT_PORT}"  # Default: 6881
      BT_TRACKERS: "${BT_TRACKERS}"  # Comma-separated tracker list
    # Multi-agent: enables serving files to other peers
    # In multi-agent setup, agents can share downloaded content
```

---

## 8.0 Remote Access

### 8.1 Remote Access Philosophy

Freecode is primarily a LOCAL tool. Remote access is opt-in and requires explicit configuration.

**All remote access binds to localhost only.** There is no server mode that exposes services externally.

### 8.2 SSH Tool

SSH is configured via the bash tool with key-based auth:

```yaml
tools:
  bash: true

permission:
  allowBash: true
  deniedCommands:
    - "ssh -R *"  # Disable reverse tunnels
```

**YubiKey Integration for SSH:**
- YubiKey can store SSH keys via PIV (Personal Identity Verification)
- Use `ssh -I /usr/lib/x86_64-linux-gnu/opensc-pkcs11.so` to load key from YubiKey
- Configure in `~/.ssh/config`:

```
Host github.com
    PKCS11Provider /usr/lib/x86_64-linux-gnu/opensc-pkcs11.so
```

### 8.3 VNC/RDP Tools

VNC and RDP are NOT built-in tools. They can be invoked via bash if needed:

```bash
# VNC (requires vncviewer installed)
vncviewer hostname:5900

# RDP (requires rdesktop or similar)
rdesktop -u user -p password hostname:3389
```

**Warning:** These expose your screen. Use VPN or tunnel over SSH.

### 8.4 Remote Access Config

```go
type RemoteAccessConfig struct {
    SSH         SSHConfig      `yaml:"ssh"`
    VNC         VNCConfig      `yaml:"vnc"`
    RDP         RDPConfig      `yaml:"rdp"`
    Tunnel      TunnelConfig   `yaml:"tunnel"`
    FileSystems FileSystems    `yaml:"fileSystems"`
}

type SSHConfig struct {
    Enabled      bool   `yaml:"enabled"`
    KeyPath      string `yaml:"keyPath"`      // Path to SSH key
    UseYubiKey   bool   `yaml:"useYubiKey"`   // Use PIV from YubiKey
    ProxyCommand string `yaml:"proxyCommand"` // ProxyJump
}

type VNCConfig struct {
    Enabled bool `yaml:"enabled"`
    Binary  string `yaml:"binary"` // vncviewer path
}

type RDPConfig struct {
    Enabled bool `yaml:"enabled"`
    Binary  string `yaml:"binary"` // rdesktop path
}

type TunnelConfig struct {
    Enabled    bool   `yaml:"enabled"`
    LocalPort  int    `yaml:"localPort"`
    RemoteHost string `yaml:"remoteHost"`
    RemotePort int    `yaml:"remotePort"`
}
```

### 8.5 File Transfer Protocols

Freecode supports remote file systems via bash tool invocation. These can also be exposed as MCP services for agent access.

**MCP Service Pattern:**
Each file system protocol can run as an MCP server, allowing the agent to browse, read, and write files without manual mount commands.

| Protocol | MCP Server | Type | Purpose |
|----------|------------|------|---------|
| SMB/CIFS | `mcp-server-smb` | Mount | Windows shares, NAS devices |
| NFS | `mcp-server-nfs` | Mount | Unix/Linux file shares |
| AFP | `mcp-server-afp` | Mount | macOS file shares |
| rsync | `mcp-server-rsync` | Transfer | One-way file sync |
| BitTorrent | `mcp-server-bittorrent` | Transfer | Distributed file sharing |

#### 8.5.1 SMB/CIFS

SMB (Server Message Block) is used for Windows file shares and NAS devices:

```bash
# Mount SMB share (requires cifs-utils on Linux)
sudo mount -t cifs //hostname/share /mntpoint -o username=user,password=pass

# Or with keychain stored credentials
mount -t cifs //hostname/share /mntpoint -o credentials=~/.smbcred

# macOS
mount_smbfs //user:pass@hostname/share /mntpoint
```

```yaml
# Config for SMB mounts
remoteAccess:
  smb:
    enabled: true
    mounts:
      - name: "nas"
        host: "nas.local"
        share: "documents"
        mountPoint: "/mnt/nas"
        credentials: "~/.smbcred"  # Stored credentials file
```

#### 8.5.2 NFS

NFS (Network File System) is used for Unix/Linux file shares:

```bash
# Mount NFS share (requires nfs-utils)
sudo mount -t nfs hostname:/exported/path /mntpoint

# NFS v4 with Kerberos
sudo mount -t nfs4 hostname:/ /mntpoint -o sec=krb5

# macOS
mount -t nfs hostname:/exported/path /mntpoint
```

```yaml
remoteAccess:
  nfs:
    enabled: true
    mounts:
      - name: "fileserver"
        host: "fileserver.local"
        export: "/home"
        mountPoint: "/mnt/fileserver"
        version: 4  # NFS version (3 or 4)
        options: "rsize=1048576,wsize=1048576"
```

#### 8.5.3 AFP (Apple Filing Protocol)

AFP is the native macOS file share protocol:

```bash
# Mount AFP share (macOS)
mount_afp afp://user:pass@hostname/ShareName /mntpoint
```

```yaml
remoteAccess:
  afp:
    enabled: true
    mounts:
      - name: "macos-share"
        host: "macos-server.local"
        share: "SharedFiles"
        mountPoint: "/mnt/macos"
        credentials: "~/.afpcred"
```

#### 8.5.4 rsync

rsync is useful for one-way file synchronization over SSH:

```bash
# Sync local to remote
rsync -avz --progress ./local/path/ user@hostname:/remote/path/

# Sync with deletion (mirror)
rsync -avz --delete ./local/path/ user@hostname:/remote/path/

# Over SSH with compression
rsync -avz -e ssh ./local/path/ user@hostname:/remote/path/
```

```yaml
remoteAccess:
  rsync:
    enabled: true
    sshKey: "~/.ssh/id_ed25519"
    sshPort: 22
    compression: true
```

#### 8.5.5 BitTorrent

BitTorrent is useful for distributing large files across multiple peers. In multi-agent setups, agents can seed files to each other.

```bash
# Download torrent (requires aria2c, transmission-cli, or similar)
aria2c --seed-time=0 "movie.torrent"

# Download with specific trackers
aria2c --bt-tracker="tracker1,tracker2" "movie.torrent"

# Create torrent (requires mktorrent or similar)
mktorrent -a tracker1,tracker2 -o file.torrent /path/to/content

# Seed a directory (make available to other peers)
mktorrent -a mytracker:6881 -l 21 -o myseed.torrent ./sharedfiles/
```

**Multi-Agent BitTorrent Setup:**
```yaml
remoteAccess:
  bittorrent:
    enabled: true
    seedDir: "/shared/downloads"
    port: 6881
    trackers:
      - "http://tracker1.example.com:6969/announce"
      - "http://tracker2.example.com:6969/announce"
    enableSeeding: true  # Allow other agents to download from this peer
```

**Warning:** BitTorrent is not encrypted by default. Use VPN or Tor for privacy.

#### 8.5.6 File System Config

```go
// FileSystemProvider is implemented by mount-based file system types.
// rsync and BitTorrent are handled via FileTransferProvider (separate interface).
type FileSystemProvider interface {
    Protocol() string  // "smb", "nfs", "afp"
    Enabled() bool
    Mounts() []MountPoint
    Validate() error
}

// MountPoint is the common interface for mount configurations
type MountPoint interface {
    Name() string
    Host() string
    MountPoint() string
    Credentials() string
}

// FileTransferProvider is implemented by file transfer protocols (rsync, BitTorrent)
type FileTransferProvider interface {
    Protocol() string  // "rsync", "bittorrent"
    Enabled() bool
    Validate() error
}

// FileSystems holds all configured file system providers and file transfer providers
type FileSystems struct {
    providers  []FileSystemProvider
    transfers  []FileTransferProvider
}

func (f *FileSystems) AddMountProvider(provider FileSystemProvider) {
    f.providers = append(f.providers, provider)
}

func (f *FileSystems) AddTransferProvider(provider FileTransferProvider) {
    f.transfers = append(f.transfers, provider)
}

func (f *FileSystems) EnabledMounts() []FileSystemProvider {
    var enabled []FileSystemProvider
    for _, p := range f.providers {
        if p.Enabled() {
            enabled = append(enabled, p)
        }
    }
    return enabled
}

func (f *FileSystems) EnabledTransfers() []FileTransferProvider {
    var enabled []FileTransferProvider
    for _, p := range f.transfers {
        if p.Enabled() {
            enabled = append(enabled, p)
        }
    }
    return enabled
}

// SMBConfig implements FileSystemProvider
type SMBConfig struct {
    Enabled bool `yaml:"enabled"`
    Mounts  []SMBMount `yaml:"mounts"`
}

func (s SMBConfig) Protocol() string  { return "smb" }
func (s SMBConfig) Enabled() bool    { return s.Enabled }
func (s SMBConfig) Mounts() []MountPoint {
    mounts := make([]MountPoint, len(s.Mounts))
    for i := range s.Mounts {
        mounts[i] = s.Mounts[i]
    }
    return mounts
}
func (s SMBConfig) Validate() error  { return nil } // TODO

type SMBMount struct {
    Name        string `yaml:"name"`
    Host        string `yaml:"host"`
    Share       string `yaml:"share"`
    MountPoint  string `yaml:"mountPoint"`
    Credentials string `yaml:"credentials"`
}

func (m SMBMount) Name() string        { return m.Name }
func (m SMBMount) Host() string        { return m.Host }
func (m SMBMount) MountPoint() string  { return m.MountPoint }
func (m SMBMount) Credentials() string { return m.Credentials }

// NFSConfig implements FileSystemProvider
type NFSConfig struct {
    Enabled bool `yaml:"enabled"`
    Mounts  []NFSMount `yaml:"mounts"`
}

func (n NFSConfig) Protocol() string  { return "nfs" }
func (n NFSConfig) Enabled() bool     { return n.Enabled }
func (n NFSConfig) Mounts() []MountPoint {
    mounts := make([]MountPoint, len(n.Mounts))
    for i := range n.Mounts {
        mounts[i] = n.Mounts[i]
    }
    return mounts
}
func (n NFSConfig) Validate() error  { return nil } // TODO

type NFSMount struct {
    Name        string `yaml:"name"`
    Host        string `yaml:"host"`
    Export      string `yaml:"export"`
    MountPoint  string `yaml:"mountPoint"`
    Version     int    `yaml:"version"` // 3 or 4
    Options     string `yaml:"options"`
}

func (m NFSMount) Name() string        { return m.Name }
func (m NFSMount) Host() string        { return m.Host }
func (m NFSMount) MountPoint() string   { return m.MountPoint }
func (m NFSMount) Credentials() string { return "" } // NFS doesn't use credentials

// AFPConfig implements FileSystemProvider
type AFPConfig struct {
    Enabled bool `yaml:"enabled"`
    Mounts  []AFPMount `yaml:"mounts"`
}

func (a AFPConfig) Protocol() string  { return "afp" }
func (a AFPConfig) Enabled() bool     { return a.Enabled }
func (a AFPConfig) Mounts() []MountPoint {
    mounts := make([]MountPoint, len(a.Mounts))
    for i := range a.Mounts {
        mounts[i] = a.Mounts[i]
    }
    return mounts
}
func (a AFPConfig) Validate() error  { return nil } // TODO

type AFPMount struct {
    Name        string `yaml:"name"`
    Host        string `yaml:"host"`
    Share       string `yaml:"share"`
    MountPoint  string `yaml:"mountPoint"`
    Credentials string `yaml:"credentials"`
}

func (m AFPMount) Name() string        { return m.Name }
func (m AFPMount) Host() string        { return m.Host }
func (m AFPMount) MountPoint() string  { return m.MountPoint }
func (m AFPMount) Credentials() string { return m.Credentials }

// BitTorrentConfig implements FileSystemProvider
type BitTorrentConfig struct {
    Enabled       bool     `yaml:"enabled"`
    SeedDir       string   `yaml:"seedDir"`
    Port          int      `yaml:"port"`
    Trackers      []string `yaml:"trackers"`
    EnableSeeding bool     `yaml:"enableSeeding"`
}

func (b BitTorrentConfig) Protocol() string { return "bittorrent" }
func (b BitTorrentConfig) Enabled() bool   { return b.Enabled }
func (b BitTorrentConfig) Validate() error  { return nil } // TODO

// BitTorrent does not implement FileSystemProvider (mount-based)
// It implements FileTransferProvider instead

// RsyncConfig implements FileTransferProvider
type RsyncConfig struct {
    Enabled      bool   `yaml:"enabled"`
    SSHKey       string `yaml:"sshKey"`       // Path to SSH key for rsync over SSH
    SSHPort      int    `yaml:"sshPort"`      // SSH port (default 22)
    Compression  bool   `yaml:"compression"`  // Use compression (-z flag)
}

func (r RsyncConfig) Protocol() string { return "rsync" }
func (r RsyncConfig) Enabled() bool   { return r.Enabled }
func (r RsyncConfig) Validate() error { return nil } // TODO
```

---

## 9.0 Security Keys (YubiKey)

### 9.1 YubiKey Integration

YubiKey can be used for:

| Use Case | Protocol | Purpose |
|----------|----------|---------|
| SSH Authentication | PIV/OpenSC | Authenticate to git@github.com |
| GPG Signing | GPG | Sign commits and tags |
| U2F/FIDO2 | WebAuthn | Two-factor authentication |
| PIV | PKCS#11 | Client certificates |

### 9.2 YubiKey Config

```go
type YubiKeyConfig struct {
    Enabled bool `yaml:"enabled"`

    // SSH via PIV
    SSH SSHKeyConfig `yaml:"ssh"`

    // GPG for commit signing
    GPG GPGKeyConfig `yaml:"gpg"`

    // U2F for freecode server auth
    U2F U2FConfig `yaml:"u2f"`
}

type SSHKeyConfig struct {
    Enabled bool `yaml:"enabled"`
    Slot string `yaml:"slot"` // 9a, 9c, 9d, 9e
    PINPrompt bool `yaml:"pinPrompt"` // Ask for PIN
}

type GPGKeyConfig struct {
    Enabled bool `yaml:"enabled"`
    KeyRef string `yaml:"keyRef"` // Key ID or fingerprint
    PINPrompt bool `yaml:"pinPrompt"`
}

type U2FConfig struct {
    Enabled bool `yaml:"enabled"`
    AppID string `yaml:"appId"` // For freecode server
}
```

### 9.3 YubiKey SSH Setup

```yaml
# ~/.config/freecode/config.yaml
yubiKey:
  enabled: true
  ssh:
    enabled: true
    slot: "9a"  # PIV Authentication
    pinPrompt: true
```

```bash
# Test YubiKey SSH
ssh -I /usr/lib/x86_64-linux-gnu/opensc-pkcs11.so git@github.com

# Or with gpg-agent + PKCS11 module
export SSH_AUTH_SOCK=$(gpgconf --list-dirs agent-ssh-socket)
ssh -T git@github.com
```

### 9.4 YubiKey GPG Setup

```yaml
yubiKey:
  enabled: true
  gpg:
    enabled: true
    keyRef: "ABC123DEF456"  # Your GPG key ID
    pinPrompt: true
```

```bash
# Configure gpg-agent to use YubiKey
echo "reader-port SCard" >> ~/.gnupg/scdaemon.conf
gpg --card-status
```

### 9.5 Where YubiKey Makes Sense

**SSH for Git:**
- YES: Store git@github.com key on YubiKey
- YES: Use for SSH to servers via jump host
- NO: Don't put personal SSH keys on YubiKey (harder to rotate)

**GPG Signing:**
- YES: Sign commits with YubiKey
- YES: Sign releases/tags
- YES: Encrypt sensitive config files

**Freecode Server Auth:**
- YES: U2F second factor for freecode server
- MAYBE: Client certificates for mTLS

**NOT Recommended:**
- SSH keys for production servers (use deploy keys instead)
- Keys that need to be shared (YubiKey is per-device)

---

## 10.0 Ralph Loop

### 10.1 Ralph Loop Config

```go
type RalphLoopConfig struct {
    Enabled            bool   `yaml:"enabled"`
    DefaultMaxIterations int  `yaml:"default_max_iterations"`
    StateDir          string `yaml:"state_dir"`
    DefaultStrategy    string `yaml:"default_strategy"` // reset|continue
}
```

---

## 11.0 Runtime Fallback

### 11.1 Runtime Fallback Config

```go
type RuntimeFallbackConfig struct {
    Enabled            bool    `yaml:"enabled"`
    RetryOnErrors      []int   `yaml:"retry_on_errors"` // HTTP codes
    MaxFallbackAttempts int    `yaml:"max_fallback_attempts"`
    CooldownSeconds    int    `yaml:"cooldown_seconds"`
    TimeoutSeconds     int    `yaml:"timeout_seconds"`
    NotifyOnFallback   bool    `yaml:"notify_on_fallback"`
}
```

---

## 12.0 Dynamic Context Pruning

### 12.1 Pruning Config

```go
type DynamicPruningConfig struct {
    Enabled       bool             `yaml:"enabled"`
    Notification  string           `yaml:"notification"` // off|minimal|detailed
    TurnProtection TurnProtection  `yaml:"turn_protection"`
    ProtectedTools []string       `yaml:"protected_tools"`
    Strategies    PruningStrategies `yaml:"strategies"`
}

type TurnProtection struct {
    Enabled bool `yaml:"enabled"`
    Turns   int  `yaml:"turns"` // 1-10
}

type PruningStrategies struct {
    Deduplication  DeduplicationConfig  `yaml:"deduplication"`
    SupersedeWrites SupersedeConfig      `yaml:"supersede_writes"`
    PurgeErrors    PurgeErrorsConfig    `yaml:"purge_errors"`
}
```

---

## 13.0 Websearch Configuration

### 13.1 Websearch Config

```go
type WebsearchConfig struct {
    Provider string `yaml:"provider"` // exa|tavily
}
```

---

## 14.0 Skills Configuration

### 14.1 Skills Overview

Skills are reusable behavior packages that agents can invoke. They contain prompts, MCP configurations, and tool definitions.

### 14.2 Built-in Skills

| Skill | Purpose | MCP |
|-------|---------|-----|
| `git-master` | Atomic commits, rebase surgery | git |
| `playwright` | Browser automation | playwright |
| `agent-browser` | Alternative browser automation | puppeteer |
| `frontend-ui-ux` | Design-first UI development | - |
| `review-work` | 5 parallel subagent review | - |
| `ai-slop-remover` | Remove AI-generated code smells | - |

### 14.3 Git Master Skill

```yaml
skill git-master:
  description: Atomic commits and rebase surgery
  agent: sisyphus
  tools:
    - bash
    - read
    - edit
  prompt: |
    You are Git Master, an expert in Git workflows.

    Your specialties:
    - Atomic commits (one logical change per commit)
    - Rebase surgery (interactive rebase, fixup, squash)
    - Commit message conventions
    - Branch management
    - Conflict resolution

    Rules:
    1. Always verify changes before committing
    2. Use clear commit messages: "type: subject"
    3. Squash related commits before merge
    4. Never force push to shared branches
```

### 14.4 Playwright Skill

```yaml
skill playwright:
  description: Browser automation with Playwright
  agent: sisyphus
  tools:
    - bash
    - read
    - write
    - playwright  # Built-in MCP
  prompt: |
    You are Playwright Master, expert in browser automation.

    Capabilities:
    - Navigate and interact with web pages
    - Fill forms and submit
    - Take screenshots and videos
    - Extract data (scraping)
    - Test UI interactions

    Usage:
    /playwright open https://example.com
    /playwright screenshot https://example.com
    /playwright click "Login" button
```

### 14.5 Skills Config

```go
type SkillsConfig struct {
    Sources   []SkillSource       `yaml:"sources"`
    Enable    []string            `yaml:"enable"`
    Disable   []string            `yaml:"disable"`
    PerSkill  map[string]SkillOverride `yaml:",remain"`
}

type SkillSource struct {
    Path      string `yaml:"path"`
    Recursive bool   `yaml:"recursive"`
    Glob      string `yaml:"glob"`
}

type SkillOverride struct {
    Description   string            `yaml:"description"`
    Template      string            `yaml:"template"`
    From          string            `yaml:"from"`
    Model         string            `yaml:"model"`
    Agent         string            `yaml:"agent"`
    Subtask       bool              `yaml:"subtask"`
    ArgumentHint  string            `yaml:"argument-hint"`
    License       string            `yaml:"license"`
    Compatibility string            `yaml:"compatibility"`
    Metadata      map[string]string `yaml:"metadata"`
    AllowedTools  []string          `yaml:"allowed-tools"`
    Disable       bool              `yaml:"disable"`
}
```

### 14.6 Skills Discovery Paths

```go
var SkillPaths = []string{
    ".freecode/skills/",           // Project level
    ".opencode/skills/",           // Legacy project
    "~/.config/freecode/skills/",  // User level
    "~/.config/opencode/skills/",  // Legacy user
    ".claude/skills/",             // Claude Code compat
    ".agents/skills/",             // Claude Code compat
}
```

---

## 15.0 Git Master Integration

### 15.1 Git Master Config

```go
type GitMasterConfig struct {
    CommitFooter       bool   `yaml:"commit_footer"`
    IncludeCoAuthoredBy bool  `yaml:"include_co_authored_by"`
    GitEnvPrefix       string `yaml:"git_env_prefix"`
}
```

---

## 16.0 Experimental Features

### 16.1 Experimental Config

```go
type ExperimentalConfig struct {
    AggressiveTruncation     bool                `yaml:"aggressive_truncation"`
    AutoResume               bool                `yaml:"auto_resume"`
    PreemptiveCompaction     bool                `yaml:"preemptive_compaction"`
    TruncateAllToolOutputs   bool                `yaml:"truncate_all_tool_outputs"`
    DynamicContextPruning    DynamicPruningConfig `yaml:"dynamic_context_pruning"`
    TaskSystem               bool                `yaml:"task_system"`
    PluginLoadTimeoutMs      int                 `yaml:"plugin_load_timeout_ms"`
    SafeHookCreation         bool                `yaml:"safe_hook_creation"`
    DisableOMOEnv            bool                `yaml:"disable_omo_env"`
    HashlineEdit             bool                `yaml:"hashline_edit"`
    ModelFallbackTitle       bool                `yaml:"model_fallback_title"`
    MaxTools                 int                 `yaml:"max_tools"`
    NewTaskSystem            bool                `yaml:"new_task_system"`
    RalphLoop                RalphLoopConfig     `yaml:"ralph_loop"`
    RuntimeFallback          RuntimeFallbackConfig `yaml:"runtime_fallback"`
}
```

---

## 17.0 Language Server Protocol (LSP)

### 17.1 LSP Overview

LSP provides language-aware features: autocomplete, goto definition, find references, rename, hover docs, diagnostics, and more. Freecode auto-detects and manages LSP servers per language.

### 17.2 Supported Languages

| Language | LSP Server | Auto-install | Extension |
|----------|------------|--------------|-----------|
| Go | gopls | If `go` available | .go |
| Python | pyright | If pyright in project | .py, .pyi |
| Rust | rust-analyzer | If `rust-analyzer` available | .rs |
| TypeScript | typescript | If in project | .ts, .tsx |
| JavaScript | typescript | If in project | .js, .jsx |
| C/C++ | clangd | Auto-install | .c, .cpp, .h, .hpp |
| C# | csharp | If .NET SDK available | .cs |
| Java | jdtls | If Java 21+ available | .java |
| Ruby | ruby-lsp | If `ruby` available | .rb, .rake |
| PHP | intelephense | Auto for PHP | .php |
| Dart | dart | If `dart` available | .dart |
| Swift | sourcekit-lsp | If Xcode on macOS | .swift |
| Kotlin | kotlin-ls | Auto for Kotlin | .kt, .kts |
| Zig | zls | If `zig` available | .zig |
| Lua | lua-ls | Auto for Lua | .lua |
| Svelte | svelte | Auto for Svelte | .svelte |
| Vue | vue | Auto for Vue | .vue |
| Astro | astro | Auto for Astro | .astro |
| Terraform | terraform | Auto-install | .tf, .tfvars |
| YAML | yaml-ls | Auto-install | .yaml, .yml |
| Bash | bash | Auto-install | .sh, .bash, .zsh |
| Elixir | elixir-ls | If `elixir` available | .ex, .exs |
| Deno | deno | If deno.json exists | .ts, .tsx, .js |
| Prisma | prisma | If `prisma` available | .prisma |
| Typst | tinymist | Auto-install | .typ, .typc |

### 17.3 LSP Config

```go
type LSPConfig struct {
    Enabled        bool              `yaml:"enabled"`
    AutoInstall    bool              `yaml:"autoInstall"`    // Auto-download LSP servers
    DownloadPath   string            `yaml:"downloadPath"`   // Where to store LSP binaries
    Servers        map[string]Server `yaml:"servers"`       // Per-language config
    Diagnostics    DiagnosticsConfig `yaml:"diagnostics"`    // Error display
}

type Server struct {
    Command  string   `yaml:"command"`  // LSP server command
    Args     []string `yaml:"args"`
    Required bool     `yaml:"required"` // Fail if not available
}

type DiagnosticsConfig struct {
    Enabled     bool `yaml:"enabled"`
    MaxProblems int  `yaml:"maxProblems"` // Problems per file
}
```

### 17.4 LSP Example

```yaml
lsp:
  enabled: true
  autoInstall: true
  downloadPath: ~/.cache/freecode/lsp
  servers:
    gopls:
      command: "gopls"
      required: false
    pyright:
      command: "pyright-langserver"
      args: ["--stdio"]
      required: false
```

### 17.5 Using LSP in Code

```go
// LSP tool for agent access
type LSPTool struct {
    server *LSPServer
}

func (l *LSPTool) Tools() []tool.Definition {
    return []tool.Definition{
        {
            Name: "lsp_complete",
            Input: CompletionRequest{},
        },
        {
            Name: "lsp_definition",
            Input: DefinitionRequest{},
        },
        {
            Name: "lsp_references",
            Input: ReferencesRequest{},
        },
        {
            Name: "lsp_hover",
            Input: HoverRequest{},
        },
        {
            Name: "lsp_diagnostics",
            Input: DiagnosticsRequest{},
        },
    }
}
```

---

## 18.0 Commands System

### 18.1 Commands Overview

Commands are reusable task templates that can be invoked with custom arguments. They provide shortcuts for common workflows.

### 18.2 Built-in Commands

| Command | Description |
|---------|-------------|
| `init` | Guided setup of AGENTS.md |
| `review` | Review changes [commit\|branch\|pr] |
| `github` | GitHub integration commands |

### 18.3 Custom Commands

Define custom commands in config:

```yaml
commands:
  build:
    description: "Build the project"
    agent: sisyphus
    model: gpt-5.4
    template: |
      Build the project using the appropriate build system.
      Project type: {{.projectType}}
      Target: {{.target | default "debug"}}
    subtask: true

  deploy:
    description: "Deploy to environment"
    agent: sisyphus
    model: claude-opus-4-7
    template: |
      Deploy {{.service}} to {{.environment}}.
      Region: {{.region | default "us-east-1"}}
    subtask: false
```

### 18.4 Command Config

```go
type CommandConfig struct {
    Commands map[string]Command `yaml:"commands"`
}

type Command struct {
    Name        string            `yaml:"name"`
    Description string            `yaml:"description"`
    Agent      string            `yaml:"agent"`      // Which agent to use
    Model      string            `yaml:"model"`       // Override default model
    Template   string            `yaml:"template"`    // Prompt template
    Subtask    bool              `yaml:"subtask"`     // Run as subtask
    Tools      []string          `yaml:"tools"`       // Restrict tools
    Env        map[string]string `yaml:"env"`         // Environment vars
}
```

### 18.5 Command Invocation

```bash
# Built-in
freecode review --commit abc123
freecode github install

# Custom
freecode run build --projectType go --target release
freecode run deploy --service api --environment prod --region eu-west-1
```

---

## 19.0 GitHub/GitLab CLI Integration

### 19.1 gh (GitHub CLI)

gh provides GitHub operations via CLI:

```yaml
tools:
  bash: true

github:
  enabled: true
  binary: "gh"  # or path to gh
```

**gh Commands:**
```bash
# Authentication
gh auth login --hostname github.com

# PR operations
gh pr create --title "Fix bug" --body "Description"
gh pr review abc123 --approve
gh pr merge abc123 --squash
gh pr list --state open

# Issues
gh issue create --title "Bug" --body "Description"
gh issue list --label bug

# Releases
gh release create v1.0.0 --title "v1.0.0" --notes "Release notes"

# Repository
gh repo clone owner/repo
gh repo view --web
```

### 19.2 glab (GitLab CLI)

glab is the GitLab equivalent:

```yaml
gitlab:
  enabled: true
  binary: "glab"  # or path to glab
```

**glab Commands:**
```bash
# Authentication
glab auth login --hostname gitlab.com

# MR operations
glab mr create --title "Fix" --description "Desc"
glab mr merge 123
glab mr list --state opened

# Issues
glab issue create --title "Bug" --description "Desc"
glab issue list --label bug

# Snippets
glab snippet create --title "Snippet" --file path
```

### 19.3 VCS Config

```go
type VCSConfig struct {
    GitHub GitHubConfig `yaml:"github"`
    GitLab GitLabConfig `yaml:"gitlab"`
}

type GitHubConfig struct {
    Enabled bool   `yaml:"enabled"`
    Binary  string `yaml:"binary"` // "gh" or path
    Host    string `yaml:"host"`   // github.com or enterprise URL
}

type GitLabConfig struct {
    Enabled bool   `yaml:"enabled"`
    Binary  string `yaml:"binary"` // "glab" or path
    Host    string `yaml:"host"`   // gitlab.com or self-hosted
}
```

---

## 20.0 Fleet Management

### 20.1 Fleet Overview

Fleet management enables a single freecode instance ("head") to coordinate multiple freecode instances ("agents") across different machines, platforms, and environments.

**Use Cases:**
- "Build this on FreeBSD" → delegates to FreeBSD agent
- "Run tests on macOS" → delegates to macOS agent
- "Deploy to Linux server" → delegates to Linux agent
- Direct specific instance for platform-specific tasks

### 20.2 Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Fleet Head (Control)                     │
│                    TLS + Auth + API                          │
│              ┌────────────────────────────────┐             │
│              │  Instance Registry              │             │
│              │  - Platform, OS, version        │             │
│              │  - Available tools             │             │
│              │  - Capabilities, tags           │             │
│              │  - Status (online/offline/busy) │             │
│              └────────────────────────────────┘             │
└─────────────────────────────────────────────────────────────┘
           │                    │                    │
    ┌──────┴──────┐      ┌──────┴──────┐      ┌──────┴──────┐
    │ FreeBSD     │      │ macOS       │      │ Linux       │
    │ Agent       │      │ Agent       │      │ Agent       │
    │ Instance    │      │ Instance    │      │ Instance    │
    └─────────────┘      └─────────────┘      └─────────────┘
```

### 20.3 Fleet Modes

| Mode | Description |
|------|-------------|
| `head` | Control plane - accepts connections from agents |
| `agent` | Worker - connects to head, receives tasks |
| `both` | Hybrid - can be head for some, agent for others |
| `client` | Thin client - connects to head to view/manage fleet |

### 20.4 Fleet Client Mode

Client mode allows you to connect to any fleet head from anywhere - your laptop at home, phone, or another workstation - to monitor and manage the fleet remotely.

```
┌──────────────────────────────────────────────────────────────────────┐
│                           Fleet Client (You)                          │
│                      freecode --mode client                           │
│                    Your laptop/phone/tablet                           │
└──────────────────────────────────────────────────────────────────────┘
                                   │
                                   │ TLS + Auth
                                   ▼
┌─────────────────────────────────────────────────────────────┐
│                     Fleet Head (Control)                     │
│              TLS + Auth + API + WebSocket                    │
│              ┌────────────────────────────────┐             │
│              │  Instance Registry              │             │
│              │  - Platform, OS, version      │             │
│              │  - Available tools             │             │
│              │  - Capabilities, tags           │             │
│              │  - Status (online/offline/busy) │             │
│              └────────────────────────────────┘             │
└─────────────────────────────────────────────────────────────┘
           │                    │                    │
    ┌──────┴──────┐      ┌──────┴──────┐      ┌──────┴──────┐
    │ FreeBSD     │      │ macOS       │      │ Linux       │
    │ Agent       │      │ Agent       │      │ Agent       │
    └─────────────┘      └─────────────┘      └─────────────┘
```

**Client Configuration:**
```yaml
fleet:
  enabled: true
  mode: client
  client:
    headURL: "https://fleet.example.com:7842"
    apiKey: "fleet-key-1-abc123"
    # Or use client certificate for mTLS
    # certPath: ~/.config/freecode/client.crt
    # keyPath: ~/.config/freecode/client.key

  # Optional: TUI mode (default when running interactively)
  # Set to false to use REST API only
  tui:
    enabled: true
    refreshInterval: 5s
```

```go
type ClientConfig struct {
    HeadURL        string        `yaml:"headURL"`
    APIKey         string        `yaml:"apiKey"`
    CertPath       string        `yaml:"certPath"`       // For mTLS
    KeyPath        string        `yaml:"keyPath"`        // For mTLS
    TUI            TUIConfig     `yaml:"tui"`
    RefreshInterval time.Duration `yaml:"refreshInterval"`
}

type TUIConfig struct {
    Enabled        bool `yaml:"enabled"`
    RefreshSeconds int  `yaml:"refreshSeconds"`
}
```

### 20.5 Interface Selection

When multiple network interfaces exist, freecode prefers LAN over WAN:

```go
type FleetConfig struct {
    Enabled        bool          `yaml:"enabled"`
    Mode           FleetMode     `yaml:"mode"`  // head, agent, both
    BindInterface  string        `yaml:"bindInterface"` // auto, lan, wan, <IP>
    Port           int           `yaml:"port"`   // Default: 7842
    TLS            TLSConfig     `yaml:"tls"`
    Auth           AuthConfig    `yaml:"auth"`
    Agent          AgentConfig   `yaml:"agent"`
}

type FleetMode string

const (
    FleetModeHead  FleetMode = "head"
    FleetModeAgent FleetMode = "agent"
    FleetModeBoth  FleetMode = "both"
)
```

**Interface Selection Logic:**
```go
func selectInterface(pref string) (net.IP, error) {
    interfaces, _ := net.Interfaces()
    var lanIPs, wanIPs []net.IP

    for _, iface := range interfaces {
        if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
            continue
        }
        addrs, _ := iface.Addrs()
        for _, addr := range addrs {
            ip, _, _ := net.ParseCIDR(addr.String())
            if ip.IsPrivate() {
                lanIPs = append(lanIPs, ip)
            } else {
                wanIPs = append(wanIPs, ip)
            }
        }
    }

    switch pref {
    case "lan":
        return lanIPs[0], nil
    case "wan":
        return wanIPs[0], nil
    case "auto":
        if len(lanIPs) > 0 {
            return lanIPs[0], nil // Prefer LAN
        }
        return wanIPs[0], nil
    default:
        return net.ParseIP(pref), nil // Specific IP
    }
}
```

### 20.6 TLS Configuration

**Option A: Self-signed (quick setup)**
```yaml
fleet:
  tls:
    autoGenerate: true  # Creates self-signed cert
    certPath: ~/.config/freecode/fleet.crt
    keyPath: ~/.config/freecode/fleet.key
```

**Option B: Certbot/Let's Encrypt**
```yaml
fleet:
  tls:
    certPath: /etc/letsencrypt/live/fleet.example.com/fullchain.pem
    keyPath: /etc/letsencrypt/live/fleet.example.com/privkey.pem
```

**Option C: Custom CA**
```yaml
fleet:
  tls:
    certPath: /path/to/ca.crt
    keyPath: /path/to/ca.key
    clientCAs:
      - /path/to/client-ca.crt
```

### 20.7 Authentication

```go
type AuthConfig struct {
    Type AuthType `yaml:"type"` // apikey, cert, mtls
    APIKeys []string `yaml:"apiKeys"` // For apikey auth
    CertPath string `yaml:"certPath"` // For cert/mTLS
}

type AuthType string

const (
    AuthAPIKey AuthType = "apikey"
    AuthCert   AuthType = "cert"
    AuthMTLS  AuthType = "mtls"
)
```

**API Key Auth:**
```yaml
fleet:
  auth:
    type: apikey
    apiKeys:
      - "fleet-key-1-abc123"
      - "fleet-key-2-def456"
```

### 20.8 Instance Registration

When an agent connects, it registers its info:

```go
type InstanceInfo struct {
    ID        string            `yaml:"id"`         // Unique instance ID
    Name      string            `yaml:"name"`       // Human-readable name
    Platform  Platform          `yaml:"platform"`    // freebsd, darwin, linux, windows
    Arch      string            `yaml:"arch"`        // amd64, arm64
    Version   string            `yaml:"version"`     // freecode version
    OSVersion string            `yaml:"osVersion"`   // OS version string
    Tools     []string          `yaml:"tools"`       // Available tools (git, gh, go, etc.)
    Tags      []string          `yaml:"tags"`        // Custom tags (build-server, prod, etc.)
    Labels    map[string]string `yaml:"labels"`      // Key-value labels
    LastSeen  time.Time         `yaml:"lastSeen"`
    Status    InstanceStatus    `yaml:"status"`      // online, offline, busy
}

type Platform string

const (
    PlatformFreeBSD Platform = "freebsd"
    PlatformDarwin  Platform = "darwin"
    PlatformLinux   Platform = "linux"
    PlatformWindows Platform = "windows"
)

type InstanceStatus int

const (
    InstanceOnline InstanceStatus = iota
    InstanceOffline
    InstanceBusy
)
```

### 20.9 Agent Configuration

```yaml
fleet:
  enabled: true
  mode: agent
  agent:
    headURL: "https://fleet.example.com:7842"
    instanceName: "freebsd-build-01"
    apiKey: "fleet-key-1-abc123"
    registerInterval: 30s
    reconnectDelay: 5s
    tags:
      - build-server
      - freebsd-14
    labels:
      environment: production
      role: ci-builder
```

```go
type AgentConfig struct {
    HeadURL          string            `yaml:"headURL"`
    InstanceName     string            `yaml:"instanceName"`
    APIKey           string            `yaml:"apiKey"`
    RegisterInterval time.Duration     `yaml:"registerInterval"`
    ReconnectDelay   time.Duration     `yaml:"reconnectDelay"`
    Tags             []string          `yaml:"tags"`
    Labels           map[string]string `yaml:"labels"`
}
```

### 20.10 Head Configuration

```yaml
fleet:
  enabled: true
  mode: head
  bindInterface: auto  # lan, wan, or specific IP
  port: 7842
  tls:
    autoGenerate: true
  auth:
    type: apikey
    apiKeys:
      - "fleet-key-1-abc123"
      - "fleet-key-2-def456"
  instances:
    autoApprove: false  # Require manual approval for new agents
    heartbeatTimeout: 60s  # Mark offline if no heartbeat
```

```go
type HeadConfig struct {
    BindInterface  string        `yaml:"bindInterface"`
    Port           int           `yaml:"port"`
    AutoApprove    bool          `yaml:"autoApprove"`
    HeartbeatTimeout time.Duration `yaml:"heartbeatTimeout"`
}
```

### 20.11 Fleet Commands

**Head commands:**
```bash
freecode fleet serve                    # Start fleet head
freecode fleet list                     # List connected instances
freecode fleet status                   # Show fleet health
freecode fleet approve <instance-id>    # Approve pending instance
freecode fleet revoke <instance-id>     # Remove instance
```

**Agent commands:**
```bash
freecode fleet join <head-url>          # Join a fleet as agent
freecode fleet leave                    # Leave current fleet
```

**Client commands (connect from your laptop):**
```bash
# Connect to fleet head (interactive)
freecode fleet connect                  # Opens connection dialog

# Connect with URL and key
freecode fleet connect https://fleet.example.com:7842 --api-key fleet-key-1

# Disconnect from current fleet
freecode fleet disconnect

# View fleet (TUI panel)
freecode fleet panel                   # Open fleet panel (Ctrl+\)
```

**Instance commands:**
```bash
# Shell into instance
freecode fleet ssh freebsd-01

# Copy files to/from instance
freecode fleet cp --from freebsd-01:/var/log/app.log --to ./logs/
freecode fleet cp --to linux-01:/tmp/script.sh --from ./script.sh
```
```bash
# Connect to fleet head
freecode fleet connect https://fleet.example.com:7842

# View fleet status
freecode fleet status                   # Show all instances, health
freecode fleet list                     # List instances with details
freecode fleet tail                     # Stream instance logs in real-time

# Dispatch commands to specific instances or groups
freecode fleet exec --instance freebsd-01 -- git status
freecode fleet exec --instance freebsd-01 -- "cd /usr/src && svn up"
freecode fleet exec --platform freebsd -- uname -a
freecode fleet exec --tag build-server -- ./run-tests.sh
freecode fleet exec --any-tag ubuntu,debian -- apt update

# Interactive shell on remote instance
freecode fleet ssh freebsd-01           # Shell into instance

# File commands
freecode fleet cp --from freebsd-01:/var/log/app.log --to ./logs/
freecode fleet cp --to linux-01:/tmp/script.sh --from ./script.sh
freecode fleet sync --src ./dist --dest linux-01:/var/www/html

# Task monitoring
freecode fleet tasks                    # List running tasks
freecode fleet tasks --watch            # Watch task output
freecode fleet logs <task-id>           # Stream task logs
```

**Dispatch from client with targeting:**
```bash
# Target by platform
freecode fleet exec --platform freebsd -- make build

# Target by tag (AND - must have all tags)
freecode fleet exec --tag prod --tag database -- pg_dumpall

# Target by any tag (OR - has any of the tags)
freecode fleet exec --any-tag web,api,worker -- ./restart.sh

# Target by label
freecode fleet exec --label environment=prod -- ./health-check.sh

# Target specific instance
freecode fleet exec --instance macos-workstation-01 -- xcodebuild

# Run on first available instance matching criteria
freecode fleet exec --platform linux --once -- ./script.sh
```

### 20.12 Fleet API

```go
// Head REST API (on port 7842)
type FleetAPI struct {
    // Instance management
    POST   /instances/register    // Agent registration
    DELETE /instances/:id         // Remove instance
    GET    /instances             // List all instances
    GET    /instances/:id         // Get instance info
    PUT    /instances/:id/tags    // Update tags

    // Task dispatch
    POST   /tasks                 // Create task
    GET    /tasks/:id             // Get task status
    POST   /tasks/:id/cancel      // Cancel task

    // File operations
    GET    /instances/:id/files/* // Read file from instance
    PUT    /instances/:id/files/* // Write file to instance
    POST   /instances/:id/exec    // Execute command
}
```

### 20.13 Task Dispatch

```go
type Task struct {
    ID          string            `yaml:"id"`
    Command     string            `yaml:"command"`      // Command to run
    Args        []string          `yaml:"args"`         // Arguments
    Target      TargetSelector     `yaml:"target"`       // Which instances
    Env         map[string]string `yaml:"env"`          // Environment
    WorkingDir  string            `yaml:"workingDir"`   // Working directory
    Timeout     time.Duration     `yaml:"timeout"`      // Max runtime
    Status      TaskStatus        `yaml:"status"`
    Result      *TaskResult       `yaml:"result"`
}

type TargetSelector struct {
    InstanceIDs []string `yaml:"instanceIds"`  // Specific instances
    Platform    Platform `yaml:"platform"`      // freebsd, darwin, linux, windows
    Tags        []string `yaml:"tags"`          // Has these tags (AND)
    AnyTag      []string `yaml:"anyTag"`       // Has any tag (OR)
    LabelSelector map[string]string `yaml:"labelSelector"` // Label match
}
```

### 20.14 File Transfer Between Instances

Files can be transferred between head and agents, or agent-to-agent via the head:

```go
// Via rsync (existing)
type FileTransfer struct {
    SourceInstance string `yaml:"sourceInstance"`
    SourcePath     string `yaml:"sourcePath"`
    DestInstance   string `yaml:"destInstance"`
    DestPath       string `yaml:"destPath"`
    Mode           string `yaml:"mode"` // copy, sync
}

// Via BitTorrent (existing)
type BTTransfer struct {
    TorrentPath string `yaml:"torrentPath"`
    SeedDir    string `yaml:"seedDir"`
    Peers      []string `yaml:"peers"` // Instance IDs
}
```

### 20.15 Multi-Agent Communication

Agents can communicate through the head for coordination:

```go
// Agent-to-Agent messaging
type AgentMessage struct {
    From      string `yaml:"from"`
    To        string `yaml:"to"`        // Instance ID or "*" for broadcast
    Type      string `yaml:"type"`      // "file-ready", "task-complete", etc.
    Payload   any    `yaml:"payload"`
}

// Use cases:
// - Head tells agent "build complete, now deploy to linux-02"
// - Agent notifies head "new binary ready at /path"
// - Agent broadcasts "cache invalidated" to all web servers
```

### 20.16 Fleet Config Complete

```go
type FleetConfig struct {
    Enabled       bool         `yaml:"enabled"`
    Mode          FleetMode    `yaml:"mode"`
    BindInterface string       `yaml:"bindInterface"`
    Port          int          `yaml:"port"`
    TLS           TLSConfig    `yaml:"tls"`
    Auth          AuthConfig   `yaml:"auth"`
    Agent         AgentConfig  `yaml:"agent"`
    Head          HeadConfig   `yaml:"head"`
    Instances     InstancesConfig `yaml:"instances"`
}

type TLSConfig struct {
    AutoGenerate bool     `yaml:"autoGenerate"`
    CertPath     string   `yaml:"certPath"`
    KeyPath      string   `yaml:"keyPath"`
    ClientCAs    []string `yaml:"clientCAs"`
}

type InstancesConfig struct {
    AutoApprove     bool          `yaml:"autoApprove"`
    HeartbeatTimeout time.Duration `yaml:"heartbeatTimeout"`
    MaxInstances    int           `yaml:"maxInstances"`
}
```

---

## 21.0 Hashline Edit Tool

### 21.1 Overview

The hashline edit tool uses content-addressed line identification instead of line numbers. This prevents "stale line" errors when files change during editing.

**Problem:** Traditional line-based editing fails when:
- Another process inserts lines above the target
- A linter reformats the file
- Git merge conflicts shift line numbers

**Solution:** Hash the content of each line to create stable identifiers.

### 21.2 How It Works

```
Line 1: func hello() {           → HASH_A1B2C3
Line 2:     return "hello";      → HASH_D4E5F6
Line 3: }                         → HASH_G7H8I9
```

When editing, reference lines by their hash, not position:
```
edit {
  file: "example.go"
  hash: "HASH_D4E5F6"
  new_content: "    return \"hello, world\";"
}
```

### 21.3 Hashline Edit Tool

```go
type HashlineTool struct{}

func (h *HashlineTool) Definition() tool.Definition {
    return tool.Definition{
        Name: "hashline_edit",
        Input: HashlineEditInput{},
    }
}

type HashlineEditInput struct {
    File       string `json:"file"`
    Hash       string `json:"hash"`        // Line content hash
    NewContent string `json:"new_content"` // New content for this line
    AfterHash string  `json:"after_hash"`  // Optional: insert after this hash
    BeforeHash string `json:"before_hash"` // Optional: insert before this hash
}

func (h *HashlineTool) Execute(ctx context.Context, input HashlineEditInput) (*tool.Result, error) {
    // 1. Read file
    // 2. Compute hashes for each line
    // 3. Find target line by hash
    // 4. Replace or insert as requested
    // 5. Write back
}
```

### 21.4 Hash Computation

```go
import "crypto/sha256"

func hashLine(line string) string {
    h := sha256.Sum256([]byte(line))
    return base64.RawURLEncoding.EncodeToString(h[:8]) // First 8 bytes
}

// Example:
// "    return \"hello\";" → "a3f4b7c2d1e8"
```

### 21.5 Hashline Config

```go
type HashlineConfig struct {
    Enabled     bool   `yaml:"enabled"`
    HashLength  int    `yaml:"hashLength"`  // Default: 8 bytes
    Algorithm   string `yaml:"algorithm"`    // sha256 (default), md5
}
```

```yaml
tools:
  hashline_edit: true

hashline:
  enabled: true
  hashLength: 8
  algorithm: sha256
```

### 21.6 Comparison

| Feature | Line-based Edit | Hashline Edit |
|---------|-----------------|---------------|
| Stability | Fails on file changes | Immune to insertions |
| Success rate | ~6.7% | ~68.3% |
| Speed | Fast | Slightly slower (hash computation) |
| Undo support | Yes | Yes |
| Git blame | Works | Works (with hash mapping) |

---

## 22.0 Slash Commands

### 22.1 Overview

Slash commands are shortcuts prefixed with `/` that trigger predefined workflows. They appear in the command palette and can be typed directly.

### 22.2 Built-in Slash Commands

| Command | Description |
|---------|-------------|
| `/init-deep` | Generate hierarchical AGENTS.md |
| `/ralph-loop` | Self-referential development loop |
| `/ulw-loop` | Ultra-work loop for long tasks |
| `/refactor` | Full refactoring workflow |
| `/start-work` | Execute Prometheus plans via Atlas |
| `/handoff` | Context summary for new sessions |
| `/review` | Multi-agent code review |
| `/test` | Generate and run tests |
| `/docs` | Generate documentation |
| `/explainshell` | Explain shell commands |

### 22.3 Slash Command Flow

```
User types: /ralph-loop

    ↓
Parse command (/ralph-loop)

    ↓
Load command definition from registry

    ↓
Execute command handler
  - Spawn RalphLoop agent
  - Set up continuation hooks
  - Configure ultrawork mode

    ↓
Return to normal mode when done
```

### 22.4 Slash Command Definition

```go
type SlashCommand struct {
    Name        string            `yaml:"name"`         // "ralph-loop"
    Description string            `yaml:"description"`  // "Self-referential dev loop"
    Usage       string            `yaml:"usage"`       // "/ralph-loop [max-iterations]"
    Agent       string            `yaml:"agent"`       // Agent to use
    Model       string            `yaml:"model"`        // Override default model
    Hooks       []string          `yaml:"hooks"`       // Hooks to enable
    Config      map[string]any   `yaml:"config"`      // Command-specific config
}

type SlashCommandRegistry struct {
    commands map[string]*SlashCommand
}
```

### 22.5 Example: Ralph Loop

```go
type RalphLoopCommand struct{}

func (c *RalphLoopCommand) Execute(ctx context.Context, args []string) error {
    // 1. Enable continuation hooks
    // 2. Spawn ralph-loop agent
    // 3. Set max iterations (default: 50)
    // 4. Configure self-reference detection
    // 5. Start loop
}
```

**Ralph Loop behavior:**
- Agent reviews its own recent output
- Looks for patterns: repeated failures, contradictions, regressions
- Can "step back" and re-plan
- Stops on success or max iterations

### 22.6 Custom Slash Commands

```yaml
commands:
  - name: deploy
    description: "Deploy service to environment"
    usage: "/deploy <service> <environment>"
    agent: sisyphus
    model: claude-opus-4-7
    hooks:
      - session.ultrawork_start
    config:
      default_env: production
      require_confirmation: true

  - name: pr
    description: "Create and review pull request"
    usage: "/pr [title]"
    agent: sisyphus
    template: |
      Create a pull request for the current changes.
      Title: {{.title}}
```

### 22.7 Slash Command Config

```go
type SlashCommandsConfig struct {
    Enabled    bool              `yaml:"enabled"`
    Commands   []*SlashCommand   `yaml:"commands"`
    Aliases    map[string]string `yaml:"aliases"`  // "r" -> "review"
}

type AliasesConfig struct {
    // "r" -> "review"
    // "h" -> "handoff"
}
```

---

## 23.0 Background Agents

### 23.1 Overview

Background agents run in parallel with the main session, enabling concurrent task execution. Each agent runs in its own tmux pane.

### 23.2 Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Main Session                            │
│                     (Sisyphus)                              │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Background Agent 1 (tmux: freecode-agent-1)        │  │
│  │  - Running task: "Run tests"                        │  │
│  └──────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Background Agent 2 (tmux: freecode-agent-2)        │  │
│  │  - Running task: "Build docs"                        │  │
│  └──────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Background Agent 3 (tmux: freecode-agent-3)        │  │
│  │  - Running task: "Lint files"                        │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 23.3 Background Agent Tool

```go
type BackgroundAgentTool struct {
    tmux *tmux.Manager
}

func (b *BackgroundAgentTool) Spawn(ctx context.Context, input SpawnInput) (*AgentResult, error) {
    name := fmt.Sprintf("freecode-agent-%d", b.nextID())

    // Create tmux pane
    pane, err := b.tmux.NewPane(tmux.PaneConfig{
        Name:    name,
        Command: "freecode agent --task " + input.TaskID,
    })
    if err != nil {
        return nil, err
    }

    // Start agent in pane
    go func() {
        b.runAgent(pane, input)
    }()

    return &AgentResult{
        PaneID:  pane.ID,
        TaskID:  input.TaskID,
        Status:  "running",
    }, nil
}

func (b *BackgroundAgentTool) List() []*AgentInfo {
    panes := b.tmux.ListPanes()
    return transformPanes(panes)
}

func (b *BackgroundAgentTool) Output(paneID string) (string, error) {
    return b.tmux.CapturePane(paneID)
}

func (b *BackgroundAgentTool) Cancel(paneID string) error {
    return b.tmux.KillPane(paneID)
}
```

### 23.4 Background Agent Config

```go
type BackgroundAgentsConfig struct {
    Enabled        bool   `yaml:"enabled"`
    MaxAgents      int    `yaml:"maxAgents"`      // Default: 5
    TmuxSocket     string `yaml:"tmuxSocket"`      // tmux socket path
    DefaultTimeout time.Duration `yaml:"defaultTimeout"`
    KeepAlive      bool   `yaml:"keepAlive"`       // Don't auto-kill on exit
}

type SpawnInput struct {
    TaskID      string            `json:"task_id"`
    Agent       string            `json:"agent"`        // "explore", "librarian"
    Model       string            `json:"model"`
    Prompt      string            `json:"prompt"`
    Tools       []string          `json:"tools"`         // Tool whitelist
    Timeout     time.Duration     `json:"timeout"`
}
```

```yaml
background_agents:
  enabled: true
  maxAgents: 5
  defaultTimeout: 30m
  keepAlive: false
```

### 23.5 Background Tasks View

```
┌─ Background Tasks ──────────────────────────────────────────┐
│                                                              │
│  Agent 1: Run tests                    [████████░░] 80%   │
│  Agent 2: Build docs                       [████░░░░░░] 40% │
│  Agent 3: Lint files                      [░░░░░░░░░░]  0%   │
│                                                              │
│  [1] Select  [c] Cancel  [l] Logs  [Enter] Attach          │
└──────────────────────────────────────────────────────────────┘
```

### 23.6 Parallel Execution Patterns

**Map-reduce pattern:**
```go
// Spawn multiple agents for parallel work
task := task.Create("build-all-platforms")
for _, platform := range ["freebsd", "linux", "darwin"] {
    task.AddSubtask(fmt.Sprintf("build-%s", platform), &SpawnInput{
        Agent:  "sisyphus",
        Prompt: fmt.Sprintf("Build for %s", platform),
    })
}
task.StartAll()
task.Wait()
```

**Pipeline pattern:**
```go
// Output of one agent feeds into next
agent1 := spawn("extract-data")
agent2 := spawn("transform-data", withInput(agent1.Output()))
agent3 := spawn("load-data", withInput(agent2.Output()))
```

---

## 24.0 Session Tools

### 24.1 Overview

Session tools manage session history, enabling search, resume, and cross-session context.

### 24.2 Session Tool Definitions

```go
type SessionTools struct {
    db *SessionDB
}

func (s *SessionTools) Definitions() []tool.Definition {
    return []tool.Definition{
        {
            Name: "session_list",
            Input: SessionListInput{},
            Description: "List all sessions",
        },
        {
            Name: "session_read",
            Input: SessionReadInput{},
            Description: "Read session content",
        },
        {
            Name: "session_search",
            Input: SessionSearchInput{},
            Description: "Search across sessions",
        },
        {
            Name: "session_info",
            Input: SessionInfoInput{},
            Description: "Get session metadata",
        },
        {
            Name: "session_export",
            Input: SessionExportInput{},
            Description: "Export session to file",
        },
        {
            Name: "session_import",
            Input: SessionImportInput{},
            Description: "Import session from file",
        },
    }
}

type SessionListInput struct {
    Filter   string `json:"filter"`    // Filter by tag, project
    Limit    int    `json:"limit"`     // Max results (default: 20)
    Offset   int    `json:"offset"`
}

type SessionReadInput struct {
    ID     string `json:"id"`      // Session ID
    Since  int    `json:"since"`    // Start at message N
    Limit  int    `json:"limit"`    // Max messages
}

type SessionSearchInput struct {
    Query    string `json:"query"`     // Search query
    Sessions []string `json:"sessions"` // Limit to these sessions
    Since    time.Time `json:"since"`   // Time range
}
```

### 24.3 Session Search Example

```bash
# Find sessions about PostgreSQL
freecode session search "postgresql connection pool"

# Find sessions in last week
freecode session search --since 7d "refactor"

# List sessions with tag
freecode session list --tag production
```

### 24.4 Session Config

```go
type SessionToolsConfig struct {
    Enabled     bool   `yaml:"enabled"`
    MaxResults  int    `yaml:"maxResults"`   // Search result limit
    AutoArchive bool   `yaml:"autoArchive"`  // Archive old sessions
    ArchiveAfter time.Duration `yaml:"archiveAfter"` // Default: 90 days
}
```

---

## 25.0 Task Tools

### 25.1 Overview

Task tools provide persistent task tracking with dependency graphs and background execution.

### 25.2 Task Tool Definitions

```go
type TaskTools struct {
    db *TaskDB
}

func (t *TaskTools) Definitions() []tool.Definition {
    return []tool.Definition{
        {
            Name: "task_create",
            Input: TaskCreateInput{},
            Description: "Create a new task",
        },
        {
            Name: "task_get",
            Input: TaskGetInput{},
            Description: "Get task details",
        },
        {
            Name: "task_list",
            Input: TaskListInput{},
            Description: "List tasks",
        },
        {
            Name: "task_update",
            Input: TaskUpdateInput{},
            Description: "Update task status",
        },
        {
            Name: "task_depend",
            Input: TaskDependInput{},
            Description: "Add task dependency",
        },
    }
}

type TaskCreateInput struct {
    Title       string            `json:"title"`
    Description string            `json:"description"`
    Tags        []string          `json:"tags"`
    Priority    int               `json:"priority"`    // 1-5
    ParentID    string            `json:"parent_id"`  // Subtask of
    Agent       string            `json:"agent"`      // Assigned agent
    Due         *time.Time       `json:"due"`        // Due date
}

type TaskDependInput struct {
    TaskID      string   `json:"task_id"`
    DependsOn   string   `json:"depends_on"`  // Task ID this depends on
}
```

### 25.3 Task Dependency Graph

```go
type TaskGraph struct {
    tasks map[string]*Task
    edges map[string][]string  // task -> dependencies
}

func (g *TaskGraph) CanStart(taskID string) bool {
    for _, dep := range g.edges[taskID] {
        if g.tasks[dep].Status != "done" {
            return false
        }
    }
    return true
}

func (g *TaskGraph) TopologicalSort() []string {
    // Kahn's algorithm for execution order
}
```

### 25.4 Task Board View

```
┌─ Tasks ─────────────────────────────────────────────────────┐
│ Project: freecode-backend                      [All ▼]      │
│                                                           │
│ TODO (3)              IN PROGRESS (2)      DONE (12)       │
│ ┌─────────────────┐   ┌─────────────────┐  ┌────────────┐ │
│ │ □ Auth refactor │   │ ▣ API v2       │  │ ✓ DB migr │ │
│ │ □ Add tests     │   │   [Agent-2]    │  │ ✓ Config   │ │
│ │ □ Docs update   │   │ □ User mgmt    │  │ ✓ Tests    │ │
│ └─────────────────┘   └─────────────────┘  └────────────┘ │
│                                                            │
│ [n] New  [d] Details  [Enter] Start  [r] Refresh          │
└────────────────────────────────────────────────────────────┘
```

### 25.5 Task Config

```go
type TaskToolsConfig struct {
    Enabled       bool   `yaml:"enabled"`
    StoragePath   string `yaml:"storagePath"`   // ~/.config/freecode/tasks/
    AutoCleanup   bool   `yaml:"autoCleanup"`   // Auto-close completed
    CleanupAfter  time.Duration `yaml:"cleanupAfter"` // After 30 days
}
```

---

## 26.0 AST-grep Tool

### 26.1 Overview

AST-grep performs AST-aware code search and replacement. Unlike regex, it understands code structure.

### 26.2 Supported Languages

| Language | Support |
|----------|---------|
| TypeScript/JavaScript | Full |
| Go | Full |
| Python | Full |
| Rust | Full |
| Java | Full |
| C/C++ | Full |
| C# | Full |
| Ruby | Full |
| PHP | Full |
| ... | 25+ languages |

### 26.3 AST-grep Tool

```go
type ASTGrepTool struct {
    runners map[string]*astgrep.Runner  // Per-language runner
}

func (a *ASTGrepTool) Definitions() []tool.Definition {
    return []tool.Definition{
        {
            Name: "ast_grep_search",
            Input: ASTSearchInput{},
            Description: "AST-aware code search",
        },
        {
            Name: "ast_grep_replace",
            Input: ASTReplaceInput{},
            Description: "AST-aware code replacement",
        },
    }
}

type ASTSearchInput struct {
    Language string `json:"language"`  // "typescript", "go", etc.
    Pattern  string `json:"pattern"`    // AST pattern
    Path     string `json:"path"`      // Directory to search
    FileFilter string `json:"file_filter"` // Glob pattern
}

type ASTReplaceInput struct {
    Language string `json:"language"`
    Pattern  string `json:"pattern"`     // Search pattern
    Fix      string `json:"fix"`         // Replacement pattern
    Path     string `json:"path"`
}
```

### 26.4 Pattern Examples

```yaml
# Find all react useState hooks
language: typescript
pattern: 'useState($TYPE)'

# Find all unhandled promises
language: typescript
pattern: 'await $EXPR'
where:
  expr_not_handled: true

# Replace fetch with axios
language: typescript
pattern: 'fetch($URL, $OPTS)'
fix: 'await axios($URL, $OPTS)'
```

### 26.5 AST-grep Config

```go
type ASTGrepConfig struct {
    Enabled  bool              `yaml:"enabled"`
    Language map[string]bool   `yaml:"language"`  // Enable per language
    CacheDir string            `yaml:"cacheDir"`  // AST cache
}
```

---

## 27.0 Look_at Tool (Multimodal)

### 27.1 Overview

Look_at analyzes files visually: PDFs, images, diagrams, screenshots. Uses multimodal model.

### 27.2 Look_at Tool

```go
type LookAtTool struct {
    model multimodal.Model
}

func (l *LookAtTool) Definition() tool.Definition {
    return tool.Definition{
        Name: "look_at",
        Input: LookAtInput{},
        Description: "Analyze files visually (PDFs, images, diagrams)",
    }
}

type LookAtInput struct {
    Path      string `json:"path"`       // File to analyze
    Query     string `json:"query"`      // "What does this diagram show?"
    Detail    string `json:"detail"`     // "low", "high"
}

func (l *LookAtTool) Execute(ctx context.Context, input LookAtInput) (*tool.Result, error) {
    data, err := os.ReadFile(input.Path)
    if err != nil {
        return nil, err
    }

    mime := mime.TypeByExtension(filepath.Ext(input.Path))

    response, err := l.model.Analyze(ctx, multimodal.Input{
        Content: data,
        MimeType: mime,
        Query:   input.Query,
        Detail:  input.Detail,
    })

    return &tool.Result{Content: response}, nil
}
```

### 27.3 Use Cases

| File Type | Example Query |
|-----------|---------------|
| PDF | "Summarize this architecture diagram" |
| Screenshot | "What error is shown?" |
| UML | "Explain this class diagram" |
| Video frame | "Describe what's happening in this frame" |
| Log output | "Extract any error messages" |

### 27.4 Look_at Config

```go
type LookAtConfig struct {
    Enabled     bool     `yaml:"enabled"`
    Model       string   `yaml:"model"`        // Multimodal model
    MaxFileSize int64    `yaml:"maxFileSize"`  // Max file size (default: 50MB)
    AllowedExts []string `yaml:"allowedExts"`  // ["pdf", "png", "jpg", ...]
}
```

---

## 28.0 Interactive Bash

### 28.1 Overview

Interactive bash runs commands that require a terminal: vim, nano, htop, ssh, mysql, psql, etc.

### 28.2 Interactive Bash Tool

```go
type InteractiveBashTool struct {
    tmux *tmux.Manager
}

func (i *InteractiveBashTool) Definition() tool.Definition {
    return tool.Definition{
        Name: "interactive_bash",
        Input: InteractiveBashInput{},
        Description: "Run interactive terminal command",
    }
}

type InteractiveBashInput struct {
    Command  string `json:"command"`   // Command to run
    WorkDir  string `json:"work_dir"`
    Height   int    `json:"height"`  // Tmux pane height
    Width    int    `json:"width"`   // Tmux pane width
}

func (i *InteractiveBashTool) Execute(ctx context.Context, input InteractiveBashInput) (*tool.Result, error) {
    pane, err := i.tmux.NewPane(tmux.PaneConfig{
        Command: input.Command,
        WorkDir: input.WorkDir,
        Height:  input.Height,
        Width:   input.Width,
    })
    if err != nil {
        return nil, err
    }

    // Wait for completion or user detach
    result, err := pane.Wait()
    return &tool.Result{Content: result.Output}, err
}
```

### 28.3 Detachable Sessions

Commands run in their own tmux pane and can be detached to run in background:

```bash
# Run vim in detachable pane
interactive_bash --command vim --detach

# Re-attach to running command
tmux attach -t freecode-cmd-1

# List running interactive commands
freecode cmd list
```

### 28.4 Interactive Bash Config

```go
type InteractiveBashConfig struct {
    Enabled        bool   `yaml:"enabled"`
    AllowedCommands []string `yaml:"allowedCommands"` // Whitelist
    BlockedCommands []string `yaml:"blockedCommands"` // Blacklist
    DefaultHeight  int    `yaml:"defaultHeight"`  // Default pane height
    DefaultWidth   int    `yaml:"defaultWidth"`   // Default pane width
}
```

```yaml
interactive_bash:
  enabled: true
  allowedCommands:
    - vim
    - nano
    - htop
    - ssh
    - mysql
    - psql
    - docker
  blockedCommands:
    - rm -rf /
    - dd
```

---

## 29.0 Built-in MCP Servers

### 29.1 Overview

Freecode includes built-in MCP servers that are always available, no configuration required.

### 29.2 Built-in MCPs

| MCP | Purpose | Auth |
|-----|---------|------|
| `exa` | Web search | None (rate limited) |
| `context7` | Official documentation | API key optional |
| `grep-app` | GitHub code search | None |
| `filesystem` | Local file access | Scope-limited |
| `git` | Git operations | Via git CLI or go-git |

### 29.3 Exa Web Search

```go
type ExaMCP struct {
    apiKey string
}

func (e *ExaMCP) Handle(ctx context.Context, req MCPRequest) (*MCPResponse, error) {
    switch req.Method {
    case "search":
        query := req.Params["query"].(string)
        results, err := e.search(ctx, query)
        return results, err
    }
}
```

**Usage:**
```bash
# Via MCP tool
mcp__exa__search --query "golang best practices 2024"
```

### 29.4 Context7 (Documentation)

```go
type Context7MCP struct {
    apiKey string
}

func (c *Context7MCP) Handle(ctx context.Context, req MCPRequest) (*MCPResponse, error) {
    switch req.Method {
    case "docs_search":
        library := req.Params["library"].(string)  // "react", "golang", etc.
        query := req.Params["query"].(string)
        return c.searchDocs(ctx, library, query)
    }
}
```

**Usage:**
```bash
mcp__context7__docs_search --library golang --query "context timeout"
```

### 29.5 Grep.app (GitHub Search)

```go
type GrepAppMCP struct{}

func (g *GrepAppMCP) Handle(ctx context.Context, req MCPRequest) (*MCPResponse, error) {
    switch req.Method {
    case "search":
        return g.githubSearch(ctx, req.Params)
    }
}
```

**Usage:**
```bash
mcp__grep_app__search --query "auth middleware language:go" --org golang --repo go
```

### 29.6 Built-in MCP Config

```go
type BuiltInMCPsConfig struct {
    Enabled    bool          `yaml:"enabled"`
    Exa        ExaConfig     `yaml:"exa"`
    Context7   Context7Config `yaml:"context7"`
    GrepApp    GrepAppConfig `yaml:"grepApp"`
    FileSystem FileSystemMCPConfig `yaml:"filesystem"`
}

type ExaConfig struct {
    Enabled    bool    `yaml:"enabled"`
    RateLimit  int     `yaml:"rateLimit"`  // Requests per minute
}

type Context7Config struct {
    Enabled bool   `yaml:"enabled"`
    APIKey  string `yaml:"apiKey"`  // Optional
}
```

---

## 30.0 Extended Model Providers

### 30.1 Overview

Freecode supports additional AI providers beyond the standard ones.

### 30.2 Provider List

| Provider | Models | Auth |
|----------|--------|------|
| `gitlab` | GitLab Duo | GitLab OAuth |
| `mistral` | Mistral AI | API key |
| `openrouter` | 100+ models | API key |
| `amazon-bedrock` | AWS models | AWS credentials |
| `azure` | Azure OpenAI | Azure AD |
| `cloudflare` | Workers AI | Cloudflare token |
| `deepinfra` | Various | API key |
| `cerebras` | Cerebras models | API key |
| `cohere` | Command R+ | API key |

### 30.3 Provider Config

```go
type ProviderConfig struct {
    GitLab GitLabProvider `yaml:"gitlab"`
    Mistral MistralProvider `yaml:"mistral"`
    OpenRouter OpenRouterProvider `yaml:"openrouter"`
    Bedrock BedrockProvider `yaml:"bedrock"`
    Azure AzureProvider `yaml:"azure"`
}

type GitLabProvider struct {
    Enabled  bool    `yaml:"enabled"`
    SiteURL  string  `yaml:"siteURL"`  // GitLab instance URL
    Token    string  `yaml:"token"`    // GitLab personal token
}

type MistralProvider struct {
    Enabled bool   `yaml:"enabled"`
    APIKey  string `yaml:"apiKey"`
    BaseURL string `yaml:"baseURL"`  // Optional custom endpoint
}
```

### 30.4 Provider Example

```yaml
providers:
  gitlab:
    enabled: true
    siteURL: https://gitlab.com
    token: glpat-xxxx

  mistral:
    enabled: true
    apiKey: ${MISTRAL_API_KEY}
```

---

## 31.0 Model Fallback Chains

### 31.1 Overview

Model fallback chains define ordered lists of models to try when primary fails.

### 31.2 Fallback Config

```go
type ModelFallbackConfig struct {
    Enabled     bool            `yaml:"enabled"`
    Chains      []FallbackChain `yaml:"chains"`
    MaxAttempts int             `yaml:"maxAttempts"`  // Global max
}

type FallbackChain struct {
    Name      string          `yaml:"name"`      // "coding", "reasoning"
    Models    []FallbackModel `yaml:"models"`   // Ordered list
    Timeout   time.Duration   `yaml:"timeout"`  // Per-attempt timeout
}

type FallbackModel struct {
    Model   string `yaml:"model"`    // "claude-opus-4-7"
    Weight  int    `yaml:"weight"`   // Selection weight
    Timeout int    `yaml:"timeout"`   // Override global timeout
    Enabled bool   `yaml:"enabled"`  // Toggle this model
}
```

### 31.3 Fallback Example

```yaml
model_fallback:
  enabled: true
  maxAttempts: 3
  chains:
    - name: coding
      timeout: 30s
      models:
        - model: claude-opus-4-7
          weight: 10
        - model: gpt-5.4
          weight: 5
        - model: claude-sonnet-4-6
          weight: 3

    - name: reasoning
      timeout: 60s
      models:
        - model: claude-opus-4-7
          weight: 10
        - model: deepseek-v3
          weight: 5
```

---

## 32.0 Additional Configuration

### 32.1 Background Task Config

```go
type BackgroundTaskConfig struct {
    DefaultConcurrency int                   `yaml:"defaultConcurrency"`
    ProviderConcurrency map[string]int       `yaml:"providerConcurrency"`
    ModelConcurrency   map[string]int       `yaml:"modelConcurrency"`
    MaxDepth           int                   `yaml:"maxDepth"`
    StaleTimeoutMs     int                  `yaml:"staleTimeoutMs"`
    MessageStalenessTimeoutMs int            `yaml:"messageStalenessTimeoutMs"`
    TaskTtlMs          int                  `yaml:"taskTtlMs"`
    SessionGoneTimeoutMs int                `yaml:"sessionGoneTimeoutMs"`
    SyncPollTimeoutMs  int                  `yaml:"syncPollTimeoutMs"`
    MaxToolCalls       int                  `yaml:"maxToolCalls"`
    CircuitBreaker     CircuitBreakerConfig `yaml:"circuitBreaker"`
}
```

### 32.2 Model Capabilities Cache

```go
type ModelCapabilitiesConfig struct {
    CachePath    string `yaml:"cachePath"`    // ~/.cache/freecode/model-capabilities.json
    RefreshOnStart bool `yaml:"refreshOnStart"`
    RefreshInterval time.Duration `yaml:"refreshInterval"` // Default: 24h
}
```

```yaml
model_capabilities:
  cachePath: ~/.cache/freecode/model-capabilities.json
  refreshOnStart: false
  refreshInterval: 24h
```

### 32.3 Browser Automation Engine

```go
type BrowserAutomationConfig struct {
    Engine    string `yaml:"engine"`  // "playwright", "agent-browser"
    Playwright PlaywrightConfig `yaml:"playwright"`
}

type PlaywrightConfig struct {
    Browser string `yaml:"browser"`  // "chromium", "firefox", "webkit"
    Headless bool `yaml:"headless"`
}
```

```yaml
browser_automation:
  engine: playwright
  playwright:
    browser: chromium
    headless: true
```

### 32.4 Comment Checker

```go
type CommentCheckerConfig struct {
    Enabled    bool   `yaml:"enabled"`
    Prompt     string `yaml:"prompt"`  // Custom check prompt
    AutoFix    bool   `yaml:"autoFix"`
}
```

```yaml
comment_checker:
  enabled: true
  prompt: |
    Check if comments are:
    1. Accurate and up-to-date
    2. Not redundant with code
    3. Properly formatted
  autoFix: true
```

---

**Author:** Mark LaPointe <mark@cloudbsd.org>
**Last Updated:** 2026-05-02 (Added skills section)

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author
