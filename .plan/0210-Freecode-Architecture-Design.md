# Freecode — Go Language Design

## 1.0 Purpose

This document specifies the Go language architecture, package structure, concurrency patterns, and design decisions for freecode.

---

## 2.0 Project Structure

```
freecode/
├── cmd/
│   ├── freecode/
│   │   └── main.go              # CLI entry
│   └── freecode-server/
│       └── main.go              # Server entry
├── internal/
│   ├── cli/                      # Cobra commands
│   ├── config/                   # Configuration
│   ├── agent/                    # Agent engine
│   ├── tool/                     # Tool implementations
│   ├── hook/                     # Hook system
│   ├── mcp/                      # MCP client
│   ├── provider/                 # AI providers
│   ├── shell/                    # Shell integration
│   ├── session/                 # Session management
│   ├── server/                   # HTTP server
│   ├── ui/                       # TUI
│   └── platform/                 # Platform-specific
├── pkg/
│   ├── api/                      # SDK packages
│   └── shared/                   # Shared utilities
├── go.mod
├── go.sum
└── Makefile
```

---

## 3.0 Package Design

### 3.1 internal/cli

Cobra command handlers.

```go
package cli

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "freecode",
    Short: "Freecode - AI coding assistant",
    Long:  `Freecode is a platform-independent AI coding assistant.`,
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.AddCommand(runCmd, serveCmd, agentCmd, mcpCmd, ...)
}
```

### 3.2 internal/config

Configuration loading, migration, and validation.

```go
package config

type Config struct {
    Shell       string
    LogLevel    string
    Yolo        bool              // Skip confirmations
    Server      ServerConfig
    Agent       AgentConfig
    // ... full schema
}

func Load(path string) (*Config, error)
func Migrate(from opencode.Config) (*Config, error)
func MergeOMO(omo *OMOConfig) error
```

### 3.3 internal/agent

Agent execution engine with 11 built-in agents.

```go
package agent

type Engine struct {
    config   *config.Config
    tools    *tool.Registry
    hooks    *hook.Registry
    sessions *session.Manager
}

func (e *Engine) Run(ctx context.Context, req Request) (*Response, error)
func (e *Engine) RegisterBuiltinAgents()
```

### 3.4 internal/tool

Tool registry and implementations.

```go
package tool

type Registry struct {
    mu    sync.RWMutex
    tools map[string]Tool
}

type Tool interface {
    Name() string
    Description() string
    Execute(ctx context.Context, req Request) (*Response, error)
    Schema() ToolSchema
}

func (r *Registry) Register(t Tool)
func (r *Registry) Get(name string) (Tool, bool)
```

### 3.5 internal/hook

Hook system with 52 hooks across 5 tiers.

```go
package hook

type Registry struct {
    mu           sync.RWMutex
    sessionHooks map[string][]SessionHook
    toolHooks    map[string][]ToolHook
    transforms   []TransformHook
    continuations []ContinuationHook
}

type SessionHook func(ctx context.Context, evt SessionEvent) error
type ToolHook func(ctx context.Context, evt ToolEvent) (error, bool) // second bool = handled
type TransformHook func(msg *Message) (*Message, error)
type ContinuationHook func(ctx context.Context, session *Session) (*ContinueSignal, error)
```

### 3.6 internal/session

Session management with tab support.

```go
package session

type Manager struct {
    mu       sync.RWMutex
    sessions map[string]*Session
    tabs     map[string]*Tab
    db       *sqlite.DB
}

type Tab struct {
    ID       string
    Name     string
    Sessions []string // Session IDs in this tab
    Active   string  // Active session ID
}

func (m *Manager) CreateTab(name string) (*Tab, error)
func (m *Manager) CloseTab(id string) error
func (m *Manager) MoveSession(sessID, tabID string) error
```

---

## 4.0 Concurrency Patterns

### 4.1 Context Propagation

```go
func Run(ctx context.Context, req Request) (*Response, error) {
    ctx, span := tracer.Start(ctx, "agent.Run")
    defer span.End()

    return e.engine.Run(ctx, req)
}
```

### 4.2 Channel-based Events

```go
type Event struct {
    Type    EventType
    Payload any
    Session string
}

events := make(chan Event, 100)
go e.processEvents(ctx, events)

select {
case evt := <-events:
    // Handle event
case <-ctx.Done():
    return ctx.Err()
}
```

### 4.3 RWMutex for Shared State

```go
type Registry struct {
    mu    sync.RWMutex
    items map[string]Item
}

func (r *Registry) Get(key string) (Item, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.items[key]
}

func (r *Registry) Set(key string, item Item) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.items[key] = item
}
```

### 4.4 errgroup for Parallel Operations

```go
import "golang.org/x/sync/errgroup"

func gatherTools(ctx context.Context, tools []string) ([]Result, error) {
    g, ctx := errgroup.WithContext(ctx)
    results := make([]Result, len(tools))

    for i, tool := range tools {
        i, tool := i, tool
        g.Go(func() error {
            result, err := executeTool(ctx, tool)
            results[i] = result
            return err
        })
    }

    return results, g.Wait()
}
```

---

## 5.0 Error Handling

### 5.1 Error Wrapping

```go
import "fmt"

if err != nil {
    return fmt.Errorf("failed to execute tool %s: %w", toolName, err)
}
```

### 5.2 Sentinel Errors

```go
var (
    ErrSessionNotFound = fmt.Errorf("session not found")
    ErrToolNotFound   = fmt.Errorf("tool not found")
    ErrAgentDisabled   = fmt.Errorf("agent is disabled")
)
```

### 5.3 Multi-error Handling

```go
import "golang.org/x/sync/errgroup"

func validateAll(configs []Config) error {
    var errs []error
    for _, cfg := range configs {
        if err := validate(cfg); err != nil {
            errs = append(errs, err)
        }
    }
    if len(errs) > 0 {
        return errors.Join(errs...)
    }
    return nil
}
```

---

## 6.0 Configuration Loading

### 6.1 Config File Locations (in order)

1. `~/.config/freecode/config.yaml`
2. `~/.config/freecode/config.json`
3. `~/.config/opencode/config.json` (read-only migration)
4. `~/.config/opencode/config.toml` (read-only migration)
5. `~/.config/opencode/opencode.json` (read-only migration)
6. `~/.config/opencode/oh-my-opencode.jsonc` (oh-my-openagent merge)
7. `$XDG_CONFIG_HOME/freecode/config.yaml`
8. Project `.freecode/config.yaml`
9. `FREECODE_CONFIG` environment variable

### 6.2 Config Merging

```go
func Load() (*Config, error) {
    cfg := DefaultConfig()

    // Load in order, later sources override
    for _, path := range configPaths() {
        if data, err := os.ReadFile(path); err == nil {
            if partial, err := parseConfig(data); err == nil {
                cfg.Merge(partial)
            }
        }
    }

    // Environment overrides
    cfg.applyEnv()

    return cfg, nil
}
```

---

## 7.0 TUI Architecture (Bubble Tea)

### 7.1 Main Model

```go
import "github.com/charmbracelet/bubbletea"

type Model struct {
    width    int
    height   int
    tabs     []TabModel
    active   int
    session  *session.Session
    yolo     bool
}

func (m Model) Init() tea.Cmd { return nil }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { ... }
func (m Model) View() string { ... }
```

### 7.2 YOLO Toggle

```go
type YoloToggleMsg bool

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case YoloToggleMsg:
        m.yolo = bool(msg)
        return m, nil
    case ToggleYolo:
        m.yolo = !m.yolo
        return m, func() tea.Msg {
            return YoloStatusMsg{Yolo: m.yolo}
        }
    }
}
```

### 7.3 Mouse and Click Support

All interactive elements support both keyboard AND mouse input. Mouse events must be explicitly enabled.

**Enable Mouse Support:**
```go
func NewProgram() *tea.Program {
    return tea.NewProgram(
        Model{},
        tea.WithMouse(),
    )
}
```

### 7.4 Clickable Elements Map

| Element | Click Action | Hover | Drag |
|---------|-------------|-------|------|
| Status bar items | Toggle/select | Highlight | - |
| Tab bar tabs | Switch tab | Highlight | Reorder tabs |
| Tab close button | Close tab | Highlight | - |
| New tab button | Create tab | Highlight | - |
| Command palette hint | Open palette | Underline | - |
| Palette items | Execute command | Highlight | Scroll |
| Instance list items | Select instance | Highlight | Multi-select |
| Task list items | Select task | Highlight | - |
| Action buttons | Execute action | Highlight | - |
| Scroll areas | Scroll | - | Scroll |
| Session content | Focus session | - | Selection |
| Text input fields | Focus field | - | Cursor position |
| Checkboxes | Toggle state | Highlight | - |
| Links/URLs | Open link | Underline | - |

### 7.5 Region Registry

All clickable regions are registered in a central registry for hit-testing:

```go
type RegionRegistry struct {
    mu sync.RWMutex
    regions map[string]*Region
}

type Region struct {
    ID       string
    Bounds   Rectangle  // X, Y, Width, Height
    Style    lipgloss.Style
    HoverStyle lipgloss.Style
    Click    func() tea.Cmd
    Hover    func()
    Leave    func()
    Drag     DragHandler
    Cursor   string  // "pointer", "crosshair", "text", etc.
}

type Rectangle struct {
    X, Y       int
    Width, Height int
}
```

**Hit Testing:**
```go
func (r *RegionRegistry) HitTest(x, y int) *Region {
    r.mu.RLock()
    defer r.mu.RUnlock()

    // Z-order: later regions are on top
    for i := len(r.regions) - 1; i >= 0; i-- {
        region := r.regions[i]
        if x >= region.Bounds.X && x < region.Bounds.X+region.Bounds.Width &&
           y >= region.Bounds.Y && y < region.Bounds.Y+region.Bounds.Height {
            return region
        }
    }
    return nil
}
```

### 7.6 Mouse Event Handler

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        return m.handleMouse(msg)
    case tea.WindowSizeMsg:
        m.handleResize(msg)
    }
}

func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
    region := m.regions.HitTest(msg.X, msg.Y)

    switch msg.Type {
    case tea.MouseLeft:
        if region != nil && region.Click != nil {
            return m, region.Click()
        }

    case tea.MouseRight:
        if region != nil {
            return m, m.showContextMenu(msg.X, msg.Y, region.ID)
        }

    case tea.MouseMotion:
        if m.hoveredRegion != region {
            if m.hoveredRegion != nil && m.hoveredRegion.Leave != nil {
                m.hoveredRegion.Leave()
            }
            m.hoveredRegion = region
            if region != nil && region.Hover != nil {
                region.Hover()
            }
        }

    case tea.WheelUp, tea.WheelDown:
        return m, m.handleScroll(msg)
    }

    return m, nil
}
```

### 7.7 Status Bar Clicks

Status bar items are clickable buttons:

```go
type StatusBar struct {
    items []StatusItem
    x     int
}

type StatusItem struct {
    Text   string
    ID     string
    Active bool
    Click  func() tea.Cmd
}

func (m *StatusBar) View() string {
    m.x = 0
    var s strings.Builder

    for _, item := range m.items {
        style := lipgloss.NewStyle().
            Background(lipgloss.Color("57")).
            Foreground(lipgloss.Color("0")).
            Padding(0, 1)

        if item.Active {
            style = style.Background(lipgloss.Color("12"))
        }

        // Register clickable region
        m.regions.Set(item.ID, &Region{
            Bounds: Rectangle{X: m.x, Y: 0, Width: len(item.Text) + 2, Height: 1},
            Style: style,
            Click: item.Click,
            Cursor: "pointer",
        })

        s.WriteString(style.Render(" "+item.Text+" "))
        m.x += len(item.Text) + 2
    }
    return s.String()
}
```

### 7.8 Tab Bar Interactions

**Clickable Tab Elements:**
```go
type TabBar struct {
    tabs     []Tab
    active   int
    regions  map[string]*Region
}

func (m *TabBar) View() string {
    var s strings.Builder
    x := 0

    for i, tab := range m.tabs {
        // Tab label
        style := m.tabStyle(i)
        if i == m.active {
            style = style.Background(m.theme.Active)
        }

        // Tab region (click to select)
        m.regions["tab_"+tab.ID] = &Region{
            Bounds: Rectangle{X: x, Y: 0, Width: len(tab.Name) + 4, Height: 1},
            Click: func() tea.Cmd { return SwitchTabMsg{Tab: i} },
            Cursor: "pointer",
        }

        s.WriteString(style.Render("[" + tab.Name + "]"))
        x += len(tab.Name) + 4

        // Close button (X)
        if i != m.active {
            closeStyle := style.Foreground(lipgloss.Color("9"))
            m.regions["close_"+tab.ID] = &Region{
                Bounds: Rectangle{X: x - 2, Y: 0, Width: 1, Height: 1},
                Click: func() tea.Cmd { return CloseTabMsg{Tab: i} },
                Cursor: "pointer",
            }
        }

        s.WriteString(" ")
        x++
    }

    // New tab button (+)
    m.regions["new_tab"] = &Region{
        Bounds: Rectangle{X: x, Y: 0, Width: 3, Height: 1},
        Click: func() tea.Cmd { return NewTabMsg{} },
        Cursor: "pointer",
    }
    s.WriteString(m.newTabStyle.Render("[+]"))

    return s.String()
}
```

### 7.9 Command Palette Clicks

Command palette items are fully clickable:

```go
type CommandPalette struct {
    items     []PaletteItem
    selected  int
    filter    string
    regions   map[string]*Region
}

func (m *CommandPalette) View() string {
    var s strings.Builder

    for i, item := range m.filteredItems() {
        style := m.itemStyle
        if i == m.selected {
            style = style.Background(m.theme.Selected)
        }

        // Register item region
        m.regions["palette_item_"+item.ID] = &Region{
            Bounds: Rectangle{X: 0, Y: i + 1, Width: m.width, Height: 1},
            Click: func() tea.Cmd { return ExecuteCommandMsg{ID: item.ID} },
            Hover: func() { m.selected = i },
            Cursor: "pointer",
        }

        // Render icon + name + key hint
        s.WriteString(style.Render(
            fmt.Sprintf(" %s %s %s",
                item.Icon,
                item.Name,
                m.keyHintStyle.Render(item.KeyHint),
            ),
        ))
        s.WriteString("\n")
    }
    return s.String()
}
```

### 7.10 Instance List (Fleet Panel)

Fleet instance list supports click to select, double-click to connect:

```go
type InstanceList struct {
    instances []Instance
    selected  int
    regions   map[string]*Region
}

func (m *InstanceList) View() string {
    var s strings.Builder
    y := 0

    for i, inst := range m.instances {
        style := m.itemStyle
        if i == m.selected {
            style = style.Background(m.theme.Selected)
        }

        // Status indicator (colored dot)
        statusIcon := m.statusIcon(inst.Status)

        // Instance row region
        m.regions["instance_"+inst.ID] = &Region{
            Bounds: Rectangle{X: 0, Y: y, Width: m.width, Height: 1},
            Click:    func() tea.Cmd { return SelectInstanceMsg{ID: inst.ID} },
            DblClick: func() tea.Cmd { return ConnectInstanceMsg{ID: inst.ID} },
            Hover:    func() { m.selected = i },
            Cursor:   "pointer",
        }

        // Status indicator separate region
        m.regions["status_"+inst.ID] = &Region{
            Bounds: Rectangle{X: 0, Y: y, Width: 2, Height: 1},
            Click: func() tea.Cmd { return ToggleInstanceStatusMsg{ID: inst.ID} },
            Cursor: "pointer",
        }

        s.WriteString(style.Render(
            fmt.Sprintf("%s %s (%s) [%s]",
                statusIcon,
                inst.Name,
                inst.Platform,
                strings.Join(inst.Tags, ","),
            ),
        ))
        s.WriteString("\n")
        y++
    }
    return s.String()
}
```

### 7.11 Context Menus (Right-Click)

Right-click shows context menu relevant to the clicked element:

```go
type ContextMenu struct {
    items []MenuItem
    x, y int
}

type MenuItem struct {
    Text   string
    ID     string
    Action func() tea.Cmd
    Disable bool
}

func (m Model) showContextMenu(x, y int, regionID string) tea.Cmd {
    var items []MenuItem

    switch regionID {
    case "instance_freebsd-01":
        items = []MenuItem{
            {Text: "Connect", ID: "connect", Action: m.connectInstance},
            {Text: "SSH", ID: "ssh", Action: m.sshInstance},
            {Text: "Copy Files", ID: "cp", Action: m.copyFiles},
            {Text: "---", ID: "sep1", Disable: true},
            {Text: "Remove", ID: "remove", Action: m.removeInstance},
        }
    case "tab":
        items = []MenuItem{
            {Text: "Close Tab", ID: "close", Action: m.closeTab},
            {Text: "Close Other Tabs", ID: "close_others", Action: m.closeOtherTabs},
            {Text: "Rename", ID: "rename", Action: m.renameTab},
        }
    default:
        items = []MenuItem{
            {Text: "Copy", ID: "copy", Action: m.copy},
            {Text: "Paste", ID: "paste", Action: m.paste},
        }
    }

    return tea.Cmd(func() tea.Msg {
        return ShowContextMenuMsg{X: x, Y: y, Items: items}
    })
}
```

### 7.12 Scroll Regions

Mouse wheel scrolls within defined scroll regions:

```go
type ScrollRegion struct {
    ID       string
    Content  []string
    Top      int      // First visible line
    Height   int      // Visible height
    OnScroll func(int) // Callback with new top line
}

func (m *Model) handleScroll(msg tea.MouseMsg) tea.Cmd {
    region := m.regions.ScrollRegionAt(msg.X, msg.Y)
    if region == nil {
        return nil
    }

    var delta int
    if msg.Type == tea.WheelUp {
        delta = -3
    } else {
        delta = 3
    }

    newTop := region.Top + delta
    if newTop < 0 {
        newTop = 0
    }
    if newTop > len(region.Content)-region.Height {
        newTop = len(region.Content) - region.Height
    }

    if region.OnScroll != nil {
        return region.OnScroll(newTop)
    }
    return nil
}
```

### 7.13 Text Selection

Click and drag to select text in session output:

```go
type Selection struct {
    StartLine, EndLine int
    StartCol, EndCol   int
    Active             bool
}

func (m *Model) handleTextSelection(msg tea.MouseMsg) {
    if msg.Type == tea.MouseLeft {
        if msg.Modifiers == 0 {
            // Start selection
            m.selection = Selection{
                StartLine: msg.Y,
                StartCol: msg.X,
                Active: true,
            }
        } else if msg.Modifiers == tea.ModifierShift && m.selection.Active {
            // Extend selection
            m.selection.EndLine = msg.Y
            m.selection.EndCol = msg.X
        }
    } else if msg.Type == tea.MouseLeft && msg.Modifiers == tea.ModifierShift {
        // Copy selection on shift-click
        if m.selection.Active {
            return m.copySelection()
        }
    }
}
```

### 7.14 Cursor Rendering

Cursor changes based on what's under the mouse:

```go
func (m Model) cursor() string {
    region := m.regions.HitTest(m.mouseX, m.mouseY)
    if region != nil {
        return region.Cursor
    }
    return "default"
}

// In View(), output cursor escape code
func (m Model) View() string {
    s := m.renderContent()
    cursor := m.cursor()
    if cursor == "pointer" {
        s += "\033[6 q" // Block cursor
    }
    return s
}
```

### 7.15 Command Palette Clicks

```go
type CommandPaletteItem struct {
    Category string
    Name     string
    Action   string  // Command ID
    KeyHint  string  // Keyboard shortcut
}

func (m CommandPalette) View() string {
    for i, item := range m.items {
        // Highlight on hover
        style := m.styles.Item
        if i == m.selected {
            style = style.Background(lipgloss.Color("57"))
        }

        // Click handling via region registration
        m.regions["cmd_"+item.Action] = ClickableRegion{
            X: 0, Y: i, Width: m.width, Height: 1,
            Action: func() { m.execute(item.Action) },
        }
    }
}
```

---

## 8.0 Server Architecture

### 8.1 Routes (chi router)

```go
import "github.com/go-chi/chi/v5"

func setupRouter(s *Server) *chi.Mux {
    r := chi.NewRouter()

    r.Route("/api/v1", func(r chi.Router) {
        r.Get("/health", s.handleHealth)
        r.Get("/session", s.handleSessionList)
        r.Post("/session", s.handleSessionCreate)
        r.Get("/session/{id}", s.handleSessionGet)
        r.Delete("/session/{id}", s.handleSessionDelete)
    })

    // Mount MCP at /mcp
    mcp.Mount(chi.NewRouter(), r)

    return r
}
```

### 8.2 Localhost Binding

```go
func (s *Server) ListenAndServe() error {
    addr := "127.0.0.1:18792" // or [::1]:18792 for IPv6
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return fmt.Errorf("failed to bind to %s: %w", addr, err)
    }
    return s.http.Serve(listener)
}
```

---

## 9.0 Database (SQLite)

### 9.1 Schema

```sql
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    title TEXT,
    created_at INTEGER,
    updated_at INTEGER,
    model TEXT,
    agent TEXT,
    tab_id TEXT,
    metadata TEXT
);

CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    session_id TEXT REFERENCES sessions(id),
    role TEXT,
    content TEXT,
    created_at INTEGER
);

CREATE TABLE tabs (
    id TEXT PRIMARY KEY,
    name TEXT,
    created_at INTEGER,
    active_session TEXT
);
```

---

## 10.0 Testing

### 10.1 Table-Driven Tests

```go
func TestToolExecute(t *testing.T) {
    tests := []struct {
        name    string
        tool    string
        input   string
        want    string
        wantErr bool
    }{
        {"bash echo", "bash", "echo hello", "hello\n", false},
        {"read missing", "read", "/nonexistent", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := execute(tt.tool, tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("execute() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("execute() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 10.2 Integration Tests

```go
func TestSessionLifecycle(t *testing.T) {
    // Create session
    sess, err := session.Create("test-session")
    if err != nil {
        t.Fatalf("failed to create session: %v", err)
    }

    // Run agent
    resp, err := agent.Run(sess.ID, "Hello")
    if err != nil {
        t.Fatalf("failed to run agent: %v", err)
    }

    if len(resp.Messages) == 0 {
        t.Error("expected messages in response")
    }
}
```

---

## 11.0 Dependencies

### 11.1 Required Packages

```go
require (
    github.com/spf13/cobra v1.8.0
    github.com/go-chi/chi/v5 v5.0.12
    github.com/charmbracelet/bubbletea v0.25.0
    modernc.org/sqlite v1.28.0
    github.com/google/uuid v1.6.0
    github.com/spf13/viper v1.18.0
    golang.org/x/sync v0.6.0
    go.opentelemetry.io/otel v1.22.0
)
```

---

**Author:** Mark LaPointe <mark@cloudbsd.org>
**Last Updated:** 2026-05-01

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author
