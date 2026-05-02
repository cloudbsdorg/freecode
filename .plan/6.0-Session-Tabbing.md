# Freecode — Session Tabbing

## 1.0 Purpose

This document describes the session tabbing feature for freecode's TUI.

---

## 2.0 Tab Model

### 2.1 Tab Structure

```go
type Tab struct {
    ID        string    // Unique tab ID
    Name      string    // Tab name (editable)
    Sessions  []string  // Session IDs in this tab
    Active    string    // Active session ID
    Layout    TabLayout // single|vertical|horizontal
    CreatedAt time.Time
    UpdatedAt time.Time
}

type TabLayout int

const (
    LayoutSingle TabLayout = iota
    LayoutVertical
    LayoutHorizontal
)
```

### 2.2 Tab Manager

```go
type TabManager struct {
    mu       sync.RWMutex
    tabs     map[string]*Tab
    order    []string // Tab order
    active   string   // Active tab ID
    db       *sql.DB
}
```

---

## 3.0 Tab Commands

### 3.1 CLI Commands

```go
// Create tab
var tabCreateCmd = &cobra.Command{
    Use:   "create [name]",
    Short: "Create a new tab",
    Args:  cobra.RangeArgs(0, 1),
    RunE: func(cmd *cobra.Command, args []string) error {
        name := "New Tab"
        if len(args) > 0 {
            name = args[0]
        }
        return tabManager.CreateTab(name)
    },
}

// Close tab
var tabCloseCmd = &cobra.Command{
    Use:   "close [tab-id]",
    Short: "Close a tab",
    RunE: func(cmd *cobra.Command, args []string) error {
        tabID := args[0]
        return tabManager.CloseTab(tabID)
    },
}

// Rename tab
var tabRenameCmd = &cobra.Command{
    Use:   "rename [tab-id] [name]",
    Short: "Rename a tab",
    RunE: func(cmd *cobra.Command, args []string) error {
        return tabManager.RenameTab(args[0], args[1])
    },
}

// List tabs
var tabListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all tabs",
    RunE: func(cmd *cobra.Command, args []string) error {
        tabs, err := tabManager.ListTabs()
        // Display
    },
}

// Move session between tabs
var tabMoveCmd = &cobra.Command{
    Use:   "move [session-id] [tab-id]",
    Short: "Move session to tab",
    RunE: func(cmd *cobra.Command, args []string) error {
        return tabManager.MoveSession(args[0], args[1])
    },
}

// Detach tab to new window
var tabDetachCmd = &cobra.Command{
    Use:   "detach [tab-id]",
    Short: "Detach tab to new window",
    RunE: func(cmd *cobra.Command, args []string) error {
        return tabManager.DetachTab(args[0])
    },
}
```

---

## 4.0 TUI Tab Interface

### 4.1 Tab Bar

```
┌─[1] main ───────[2] feature ───────[3] debug ───────[+]─┐
│                                                             │
│  Session content here...                                    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 Tab Keybindings

| Key | Action |
|-----|--------|
| `Ctrl+T` | New tab |
| `Ctrl+W` | Close current tab |
| `Ctrl+Tab` | Next tab |
| `Ctrl+Shift+Tab` | Previous tab |
| `Alt+1-9` | Switch to tab 1-9 |
| `Ctrl+R` | Rename current tab |
| `Ctrl+Shift+V` | Vertical split |
| `Ctrl+Shift+H` | Horizontal split |
| `Ctrl+Shift+D` | Detach tab |
| `Ctrl+Y` | Toggle YOLO mode |

### 4.3 Tab Model (Bubble Tea)

```go
type TabModel struct {
    tabs    []Tab
    active  int
    width   int
    height  int
}

func (m TabModel) Init() tea.Cmd {
    return nil
}

func (m TabModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+t":
            m.tabs = append(m.tabs, Tab{
                ID:   uuid.New().String(),
                Name: "New Tab",
            })
            m.active = len(m.tabs) - 1
        case "ctrl+w":
            if len(m.tabs) > 1 {
                m.tabs = append(m.tabs[:m.active], m.tabs[m.active+1:]...)
                if m.active >= len(m.tabs) {
                    m.active = len(m.tabs) - 1
                }
            }
        case "ctrl+tab":
            m.active = (m.active + 1) % len(m.tabs)
        case "ctrl+shift+tab":
            m.active = (m.active - 1 + len(m.tabs)) % len(m.tabs)
        case "alt+1", "alt+2", "alt+3", "alt+4", "alt+5", "alt+6", "alt+7", "alt+8", "alt+9":
            idx := int(msg.String()[3] - '1')
            if idx < len(m.tabs) {
                m.active = idx
            }
        case "ctrl+r":
            return m, func() tea.Msg { return RenameTabMsg{} }
        }
    }
    return m, nil
}

func (m TabModel) View() string {
    s := renderTabBar(m.tabs, m.active)
    s += "\n"
    s += renderTabContent(m.tabs[m.active])
    return s
}
```

---

## 5.0 Command Palette

### 5.1 Command Palette Overview

The command palette (`Ctrl+P`) provides quick access to all freecode commands including fleet operations.

```
┌─ Command Palette ──────────────────────────────────────┐
│ > fleet                                              │
│                                                     │
│   Fleet                                             │
│   ├─ connect     Connect to fleet head...            │
│   ├─ status      View fleet status                   │
│   ├─ instances   List connected instances            │
│   ├─ exec        Dispatch command to fleet...        │
│   ├─ ssh         Shell into instance...              │
│   ├─ cp          Copy files to/from instances...    │
│   ├─ tasks       View running tasks                  │
│   └─ logs        Stream instance logs                │
│                                                     │
│   Session                                           │
│   ├─ new tab     Create new tab                      │
│   ├─ close tab   Close current tab                   │
│   ├─ rename     Rename tab...                        │
│   └─ split      Split view...                       │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### 5.2 Menu Keybindings

| Key | Action |
|-----|--------|
| `Ctrl+P` | Open command palette |
| `Ctrl+,` | Open settings |
| `Ctrl+\` | Toggle fleet panel |
| `Esc` | Close palette |

### 5.3 Fleet Menu Items

**Fleet Panel (`Ctrl+\`):**
```
┌─ Fleet ──────────────────────────────────────────────────┐
│ ● Connected to: fleet.example.com:7842                    │
│                                                          │
│ Instances (5)                                            │
│ ├─ 🟢 freebsd-build-01  (freebsd/amd64)  [build-server]│
│ ├─ 🟢 macos-workstation   (darwin/arm64)  [dev]        │
│ ├─ 🟢 linux-server-01   (linux/amd64)   [prod,web]    │
│ ├─ 🟡 windows-ci         (windows/amd64) [ci]          │
│ └─ ⚠ freebsd-nas        (freebsd/amd64)  [OFFLINE]    │
│                                                          │
│ Tasks (2 running)                                       │
│ ├─ Building freebsd-01...  [████████░░] 80%           │
│ └─ Testing macos...       [██░░░░░░░░░] 20%           │
│                                                          │
│ Actions: [Connect] [Exec] [SSH] [CP] [Tasks]           │
└──────────────────────────────────────────────────────────┘
```

### 5.4 Fleet Connect Flow

**Connect to Fleet Head:**
```
┌─ Connect to Fleet ──────────────────────────────────────┐
│                                                          │
│ Fleet URL:                                               │
│ [ https://fleet.example.com:7842________________]       │
│                                                          │
│ Authentication:                                          │
│ ○ API Key                                                │
│ ● mTLS (Certificate)                                    │
│                                                          │
│ [ Connect ]  [ Cancel ]                                  │
└──────────────────────────────────────────────────────────┘
```

### 5.5 Command Palette Config

```go
type CommandPaletteConfig struct {
    Enabled     bool              `yaml:"enabled"`
    Keybinding  string            `yaml:"keybinding"` // Default: "ctrl+p"
    FuzzyMatch  bool              `yaml:"fuzzyMatch"`
    MaxResults  int              `yaml:"maxResults"`
    Categories  []string          `yaml:"categories"`  // Order of categories
}

type FleetPaletteConfig struct {
    Enabled        bool `yaml:"enabled"`
    PanelKeybinding string `yaml:"panelKeybinding"` // Default: "ctrl+\"
    RefreshSeconds int   `yaml:"refreshSeconds"`
}
```

### 5.6 Fleet Panel Model

```go
type FleetPanelModel struct {
    head       *FleetHead
    instances  []InstanceInfo
    tasks      []Task
    selected   int
    view       FleetPanelView
}

type FleetPanelView int

const (
    ViewInstances FleetPanelView = iota
    ViewTasks
    ViewLogs
    ViewConnect
)

func (m FleetPanelModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "j", "down":
            m.selected++
        case "k", "up":
            m.selected--
        case "enter":
            return m, m.handleSelect()
        case "d":
            return m, m.handleDisconnect()
        case "e":
            return m, m.handleExec()
        case "s":
            return m, m.handleSSH()
        case "c":
            return m, m.handleCopy()
        case "l":
            m.view = ViewLogs
        case "t":
            m.view = ViewTasks
        case "i":
            m.view = ViewInstances
        }
    }
    return m, nil
}
```

---

## 6.0 Split View

### 5.1 Split Layouts

**Vertical Split:**
```
┌─────────────────────┬─────────────────────┐
│     Session A       │     Session B       │
│                     │                     │
└─────────────────────┴─────────────────────┘
```

**Horizontal Split:**
```
┌─────────────────────┐
│     Session A       │
├─────────────────────┤
│     Session B       │
└─────────────────────┘
```

### 5.2 Split Implementation

```go
type SplitView struct {
    Layout   SplitLayout
    Panes    []Pane
    Sizes    []int // Percentage for each pane
}

type Pane struct {
    SessionID string
    Border   bool
    Focused  bool
}

func (m *TabModel) SplitCurrent(direction SplitDirection) error {
    tab := &m.tabs[m.active]
    if len(tab.Sessions) == 0 {
        return fmt.Errorf("no session to split")
    }

    activeSession := tab.Active

    switch direction {
    case SplitVertical:
        tab.Layout = LayoutVertical
        tab.Sessions = append(tab.Sessions, activeSession)
    case SplitHorizontal:
        tab.Layout = LayoutHorizontal
        tab.Sessions = append(tab.Sessions, activeSession)
    }
    return nil
}
```

---

## 7.0 YOLO Toggle

### 6.1 YOLO State

```go
type YOLOState struct {
    Enabled bool
    Config  YOLOConfig
}

type YOLOConfig struct {
    SkipEditConfirmations   bool
    SkipBashConfirmations  bool
    SkipDeleteConfirmations bool
    SkipPermissionPrompts  bool
    SkipToolConfirmations  bool
}
```

### 6.2 Commands Menu

```
┌─ Commands ────────────────────────────────┐
│                                             │
│  [✓] YOLO Mode                             │
│                                             │
│  Tab:                                       │
│    [ ] New Tab        Ctrl+T               │
│    [ ] Close Tab      Ctrl+W               │
│    [ ] Rename Tab     Ctrl+R               │
│    [ ] Split Vertical Ctrl+Shift+V         │
│    [ ] Split Horizontal Ctrl+Shift+H       │
│    [ ] Detach Tab     Ctrl+Shift+D         │
│                                             │
│  Session:                                   │
│    [ ] New Session                          │
│    [ ] Fork Session                         │
│    [ ] Share Session                        │
│                                             │
│  General:                                   │
│    [ ] Toggle Theme                         │
│    [ ] Open Config                          │
│    [ ] Run Doctor                           │
│                                             │
│                              [Esc] Close    │
└─────────────────────────────────────────────┘
```

### 6.3 YOLO Keybinding

```go
func (m TabModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+y":
            m.yolo = !m.yolo
            return m, func() tea.Msg {
                return YOLOStatusMsg{Enabled: m.yolo}
            }
        }
    }
    return m, nil
}
```

---

## 8.0 Tab Persistence

### 7.1 Database Schema

```sql
CREATE TABLE tabs (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    layout TEXT NOT NULL DEFAULT 'single',
    active_session TEXT,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE tab_sessions (
    tab_id TEXT REFERENCES tabs(id) ON DELETE CASCADE,
    session_id TEXT REFERENCES sessions(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    PRIMARY KEY (tab_id, session_id)
);
```

### 7.2 Tab Persistence Operations

```go
func (m *TabManager) SaveTab(tab *Tab) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.Exec(`
        INSERT OR REPLACE INTO tabs (id, name, layout, active_session, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `, tab.ID, tab.Name, string(tab.Layout), tab.Active, tab.CreatedAt.Unix(), tab.UpdatedAt.Unix())

    _, err = tx.Exec(`DELETE FROM tab_sessions WHERE tab_id = ?`, tab.ID)
    for i, sessID := range tab.Sessions {
        _, err = tx.Exec(`
            INSERT INTO tab_sessions (tab_id, session_id, position)
            VALUES (?, ?, ?)
        `, tab.ID, sessID, i)
    }

    return tx.Commit()
}
```

---

## 9.0 Tab Lifecycle

### 8.1 Create Tab

```go
func (m *TabManager) CreateTab(name string) (*Tab, error) {
    tab := &Tab{
        ID:        uuid.New().String(),
        Name:      name,
        Sessions:  []string{},
        Layout:    LayoutSingle,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := m.SaveTab(tab); err != nil {
        return nil, err
    }

    m.mu.Lock()
    m.tabs[tab.ID] = tab
    m.order = append(m.order, tab.ID)
    m.active = tab.ID
    m.mu.Unlock()

    return tab, nil
}
```

### 8.2 Close Tab

```go
func (m *TabManager) CloseTab(tabID string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    tab, ok := m.tabs[tabID]
    if !ok {
        return ErrTabNotFound
    }

    // Delete tab sessions
    for _, sessID := range tab.Sessions {
        if err := m.sessions.DeleteSession(sessID); err != nil {
            // Log but continue
        }
    }

    // Delete tab
    _, err := m.db.Exec(`DELETE FROM tabs WHERE id = ?`, tabID)
    if err != nil {
        return err
    }

    delete(m.tabs, tabID)
    m.order = slices.DeleteFunc(m.order, func(id string) bool {
        return id == tabID
    })

    // Set new active
    if len(m.order) > 0 {
        if m.active == tabID {
            m.active = m.order[len(m.order)-1]
        }
    }

    return nil
}
```

### 8.3 Move Session Between Tabs

```go
func (m *TabManager) MoveSession(sessionID, targetTabID string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Find current tab
    var currentTabID string
    for tabID, tab := range m.tabs {
        for _, sessID := range tab.Sessions {
            if sessID == sessionID {
                currentTabID = tabID
                break
            }
        }
    }

    if currentTabID == "" {
        return ErrSessionNotFound
    }

    targetTab, ok := m.tabs[targetTabID]
    if !ok {
        return ErrTabNotFound
    }

    // Remove from current tab
    currentTab := m.tabs[currentTabID]
    currentTab.Sessions = slices.DeleteFunc(currentTab.Sessions, func(id string) bool {
        return id == sessionID
    })
    if currentTab.Active == sessionID {
        if len(currentTab.Sessions) > 0 {
            currentTab.Active = currentTab.Sessions[len(currentTab.Sessions)-1]
        } else {
            currentTab.Active = ""
        }
    }

    // Add to target tab
    targetTab.Sessions = append(targetTab.Sessions, sessionID)
    targetTab.Active = sessionID
    targetTab.UpdatedAt = time.Now()

    // Save both
    if err := m.SaveTab(currentTab); err != nil {
        return err
    }
    return m.SaveTab(targetTab)
}
```

---

## 10.0 Default Tab Behavior

### 9.1 Startup

- Create default tab "main" if no tabs exist
- Load last active tab from database
- Restore sessions in each tab

### 9.2 New Session in Tab

```go
func (m *TabManager) NewSessionInTab(tabID string, msg string) (*Session, error) {
    tab, ok := m.tabs[tabID]
    if !ok {
        return nil, ErrTabNotFound
    }

    sess, err := m.sessions.Create(msg)
    if err != nil {
        return nil, err
    }

    tab.Sessions = append(tab.Sessions, sess.ID)
    tab.Active = sess.ID
    tab.UpdatedAt = time.Now()

    if err := m.SaveTab(tab); err != nil {
        return nil, err
    }

    return sess, nil
}
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
