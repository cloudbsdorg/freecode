# Freecode — Target Architecture Analysis

## 1.0 Purpose

This document maps the current opencode TypeScript architecture to the target Go architecture for freecode. It serves as a reference for implementation decisions.

---

## 2.0 OpenCode Package Structure

### 2.1 Monorepo Layout (Current)

```
opencode/
├── packages/
│   ├── opencode/           # Main CLI (TypeScript)
│   │   ├── src/
│   │   │   ├── index.ts          # CLI entry (yargs)
│   │   │   ├── config/          # Config service
│   │   │   ├── cli/             # CLI commands
│   │   │   ├── tool/            # Tools (bash, read, write, etc.)
│   │   │   ├── agent/           # Agent implementation
│   │   │   ├── provider/       # AI provider integration
│   │   │   ├── server/          # Hono server
│   │   │   ├── storage/         # SQLite storage
│   │   │   ├── effect/          # Effect system layer
│   │   │   └── shell/           # Shell utilities
│   │   └── package.json
│   ├── core/              # Shared utilities
│   ├── plugin/            # Plugin interface
│   ├── ui/                 # UI components
│   └── app/                # Web application
└── package.json           # Turborepo root
```

### 2.2 Key File Locations

| Component | TypeScript File | Lines |
|-----------|----------------|-------|
| CLI Entry | `/packages/opencode/src/index.ts` | 247 |
| Config Service | `/packages/opencode/src/config/config.ts` | 807 |
| Config Schema | `/packages/opencode/src/config/config.ts:102-261` | ~160 |
| Run Command | `/packages/opencode/src/cli/cmd/run.ts` | ~300 |
| Agent Service | `/packages/opencode/src/agent/agent.ts` | ~400 |
| Tool Implementations | `/packages/opencode/src/tool/*.ts` | Various |
| Shell Utils | `/packages/opencode/src/shell/shell.ts` | ~300 |
| Server | `/packages/opencode/src/server/server.ts` | ~500 |

---

## 3.0 Go Target Structure

### 3.1 Package Layout

```
freecode/
├── cmd/
│   └── freecode/
│       └── main.go              # CLI entry point
├── internal/
│   ├── cli/
│   │   ├── root.go              # Root cobra command
│   │   ├── run.go               # run command
│   │   ├── serve.go             # server command
│   │   ├── agent.go             # agent command
│   │   └── ...                  # Other commands
│   ├── config/
│   │   ├── config.go            # Config struct and loading
│   │   ├── schema.go           # Config schema definitions
│   │   ├── migration.go        # Config migration from opencode
│   │   └── omo_config.go       # oh-my-openagent config
│   ├── agent/
│   │   ├── engine.go           # Agent execution engine
│   │   ├── router.go          # Agent routing
│   │   └── builtin.go          # Built-in agents (11)
│   ├── tool/
│   │   ├── registry.go        # Tool registry
│   │   ├── bash.go            # Bash tool
│   │   ├── read.go            # Read tool
│   │   ├── write.go           # Write tool
│   │   ├── edit.go            # Edit tool
│   │   ├── glob.go            # Glob tool
│   │   ├── grep.go            # Grep tool
│   │   ├── task.go            # Task/subagent tool
│   │   └── ...                # Other tools
│   ├── hook/
│   │   ├── registry.go        # Hook registry
│   │   ├── session.go         # Session hooks (24)
│   │   ├── tool.go            # Tool hooks (14)
│   │   ├── transform.go       # Transform hooks (5)
│   │   └── continuation.go    # Continuation hooks (7)
│   ├── mcp/
│   │   ├── client.go         # MCP client
│   │   └── oauth.go           # OAuth handling
│   ├── provider/
│   │   ├── anthropic.go      # Anthropic provider
│   │   ├── openai.go         # OpenAI provider
│   │   └── ...               # Other providers
│   ├── shell/
│   │   ├── shell.go          # Shell detection/selection
│   │   ├── pty.go            # PTY handling
│   │   └── unix.go           # Unix-specific shell handling
│   ├── session/
│   │   ├── manager.go        # Session manager
│   │   ├── tab.go            # Session tabs
│   │   ├── store.go          # SQLite storage
│   │   └── compaction.go     # Context compaction
│   ├── server/
│   │   ├── server.go         # HTTP API server
│   │   ├── routes/           # Route handlers
│   │   └── adapter.go        # Platform adapters
│   ├── ui/
│   │   ├── tui.go            # Bubbletea TUI
│   │   ├── tabs.go           # Tab management
│   │   └── ...               # UI components
│   └── platform/
│       ├── freebsd.go        # FreeBSD-specific
│       ├── darwin.go         # macOS-specific
│       ├── linux.go          # Linux-specific
│       └── illuminos.go      # IllumOS-specific
├── pkg/
│   ├── api/
│   │   └── v1/               # SDK for external clients
│   └── shared/
│       └── ...               # Shared utilities
└── go.mod
```

---

## 4.0 TypeScript → Go Mapping

### 4.1 Effect System (TypeScript) → Context + Errors (Go)

**TypeScript Effect:**
```typescript
export const Service = Context.Service<Service, Interface>()("@opencode/Config")

export const layer = Layer.effect(Service, Effect.gen(function* () {
  const config = yield* Config
  return { ... }
}))
```

**Go Equivalent:**
```go
type ConfigService interface {
    GetConfig() (*Config, error)
    SetConfig(*Config) error
}

type configService struct {
    config *Config
    mu     sync.RWMutex
}
```

### 4.2 yargs → Cobra

**TypeScript (yargs):**
```typescript
yargs
  .command('run [message..]', 'Run opencode', (y) => {
    y.option('continue', { alias: 'c' })
    y.option('session', { alias: 's' })
  })
  .parse()
```

**Go (Cobra):**
```go
var runCmd = &cobra.Command{
    Use:   "run [message..]",
    Short: "Run freecode",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
    },
}
runCmd.Flags().BoolP("continue", "c", false, "Continue last session")
runCmd.Flags().StringP("session", "s", "", "Session ID")
```

### 4.3 Hono Server → net/http + chi

**TypeScript (Hono):**
```typescript
const app = new Hono()
app.get('/api/session', (c) => c.json(sessions))
app.post('/api/session', async (c) => {
  const body = await c.req.json()
  return c.json(createSession(body))
})
```

**Go (chi):**
```go
r := chi.NewRouter()
r.Get("/api/session", handleSessionList)
r.Post("/api/session", handleSessionCreate)
```

### 4.4 SQLite (Drizzle) → SQLite (modernc.org/sqlite)

**TypeScript (Drizzle):**
```typescript
import { sql } from 'drizzle-orm'
const sessions = await db.select().from(sessionTable)
```

**Go (modernc.org/sqlite):**
```go
rows, err := db.Query("SELECT * FROM sessions")
```

Or using sqlx for convenience.

---

## 5.0 Key Configuration Schema

### 5.1 Main Config (from opencode)

```go
type Config struct {
    Shell       string            `json:"shell"`
    LogLevel    string            `json:"logLevel"` // DEBUG|INFO|WARN|ERROR
    Server      ServerConfig      `json:"server"`
    Commands    map[string]CommandConfig `json:"command"`
    Skills      SkillsConfig      `json:"skills"`
    Watcher     WatcherConfig     `json:"watcher"`
    Snapshot    bool              `json:"snapshot"`
    Plugins     []PluginSpec      `json:"plugin"`
    Share       string            `json:"share"` // manual|auto|disabled
    AutoUpdate  bool             `json:"autoupdate"`
    Model       string            `json:"model"`
    SmallModel  string            `json:"small_model"`
    DefaultAgent string           `json:"default_agent"`
    Agent       AgentConfig       `json:"agent"`
    Providers   map[string]ProviderConfig `json:"provider"`
    MCP         map[string]MCPConfig `json:"mcp"`
    Formatter   FormatterConfig   `json:"formatter"`
    LSP         LSPConfig        `json:"lsp"`
    Instructions []string        `json:"instructions"`
    Permission  PermissionConfig `json:"permission"`
    Tools       map[string]bool  `json:"tools"`
    Compaction  CompactionConfig `json:"compaction"`
    Experimental ExperimentalConfig `json:"experimental"`
    Yolo        bool              `json:"yolo"` // Skip confirmations
}
```

### 5.2 oh-my-openagent Config

```go
type OMOConfig struct {
    NewTaskSystemEnabled bool              `json:"new_task_system_enabled"`
    DefaultRunAgent      string             `json:"default_run_agent"`
    AgentDefinitions     []string           `json:"agent_definitions"`
    DisabledMCps         []string           `json:"disabled_mcps"`
    DisabledAgents      []string           `json:"disabled_agents"`
    DisabledSkills      []string           `json:"disabled_skills"`
    DisabledHooks       []string           `json:"disabled_hooks"`
    DisabledCommands    []string           `json:"disabled_commands"`
    DisabledTools       []string           `json:"disabled_tools"`
    MCPEnvAllowlist     []string           `json:"mcp_env_allowlist"`
    HashlineEdit        bool               `json:"hashline_edit"`
    ModelFallback       bool               `json:"model_fallback"`
    Agents              map[string]AgentConfig `json:"agents"`
    Categories          map[string]CategoryConfig `json:"categories"`
    ClaudeCode          ClaudeCodeConfig   `json:"claude_code"`
    Tmux                TmuxConfig         `json:"tmux"`
    RalphLoop           RalphLoopConfig    `json:"ralph_loop"`
    RuntimeFallback     RuntimeFallbackConfig `json:"runtime_fallback"`
    BackgroundTask      BackgroundTaskConfig `json:"background_task"`
    Experimental        ExperimentalConfig `json:"experimental"`
    Skills              SkillsConfig      `json:"skills"`
    Websearch           WebsearchConfig   `json:"websearch"`
}
```

---

## 6.0 Platform-Specific Code Locations

| Platform | File | Purpose |
|----------|------|---------|
| All | `internal/shell/shell.go` | Shell detection |
| All | `internal/shell/pty.go` | PTY management |
| win32 | `internal/platform/win32.go` | Windows-specific |
| darwin | `internal/platform/darwin.go` | macOS-specific |
| freebsd | `internal/platform/freebsd.go` | FreeBSD-specific |
| linux | `internal/platform/linux.go` | Linux-specific |
| illuminos | `internal/platform/illuminos.go` | IllumOS-specific |

---

## 7.0 Tool Implementations

| Tool | File | Purpose |
|------|------|---------|
| bash | `internal/tool/bash.go` | Execute shell commands |
| read | `internal/tool/read.go` | Read files |
| write | `internal/tool/write.go` | Write files |
| edit | `internal/tool/edit.go` | Edit files (diff-based) |
| glob | `internal/tool/glob.go` | File pattern matching |
| grep | `internal/tool/grep.go` | Search file contents |
| webfetch | `internal/tool/webfetch.go` | Fetch web content |
| websearch | `internal/tool/websearch.go` | Web search |
| task | `internal/tool/task.go` | Subagent execution |
| skill | `internal/tool/skill.go` | Skill invocation |
| todo | `internal/tool/todo.go` | Todo list management |
| question | `internal/tool/question.go` | Ask user questions |
| plan | `internal/tool/plan.go` | Plan mode |
| lsp | `internal/tool/lsp.go` | LSP integration |

---

**Author:** Mark LaPointe <mark@cloudbsd.org>
**Last Updated:** 2026-05-01

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author
