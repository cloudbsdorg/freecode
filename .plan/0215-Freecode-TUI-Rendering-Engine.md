# Freecode TUI Rendering Engine

**Document ID:** TUI-001
**Version:** 2.0
**Last Updated:** 2026-05-10
**Status:** IMPLEMENTED

## Overview

The TUI rendering engine provides a reactive, component-based UI system with external templates, state-driven rendering, and dynamic interpolation. It supports multiple renderers (Bubble/lipgloss, LCD, Headless) through a generic interface.

## Implemented Architecture

```
templates/
├── components/          # Reusable component templates
│   ├── text/
│   ├── button/
│   ├── input/
│   ├── list/
│   ├── window/
│   ├── vbox/
│   ├── hbox/
│   ├── grid/
│   ├── spacer/
│   ├── divider/
│   ├── progress/
│   ├── toast/
│   ├── tabbar/
│   ├── messagelist/
│   ├── selectionlist/
│   └── dialog/
└── views/              # Full page/view templates
    ├── setup/
    ├── home/
    ├── session/
    └── commandpalette/
```

## Core Implementation

### ReactiveEngine[R Renderer]

The `ReactiveEngine` provides React-like state management:

```go
engine := template.NewReactiveEngine[BubbleRenderer]()

// Set state - automatically marks components dirty
engine.Set("username", "Alice")
engine.Set("messages", []string{"Hi", "Hello"})
engine.Set("showDialog", true)

// Interpolate template with state
result := engine.Interpolate(`<text value="Hello ${username}!" />`)

// Full render with layout calculation
output := engine.MustRender(templateSrc, width, height, renderer)

// Or render from external template file
output := engine.RenderTemplate("home", width, height, renderer)
```

### State-Driven Rendering

```go
// Subscribe to state changes
for change := range engine.Subscribe() {
    fmt.Printf("State changed: %s = %v\n", change.Key, change.Value)
}

// Render only dirty components for efficiency
output := engine.RenderDirty(src, width, height, renderer)
```

### Template Loading

```go
// Load templates from external files
engine := template.NewReactiveEngineWithLoader[BubbleRenderer]("./templates")

// Or manually load specific templates
engine.LoadTemplate("home")
engine.LoadComponent("button")

// Get raw template content
src := engine.GetTemplate("home")
```

## Template Syntax

Templates use HTML/XML-like syntax with `${variable}` interpolation:

```xml
<window title="Setup" width="60" height="20" center="true" padding="2">
  <vbox gap="1">
    <text value="Welcome ${name}!" bold="true" />
    <divider />
    <input id="provider" placeholder="Select provider..." />
    <spacer flex="1" />
    <hbox gap="2">
      <spacer flex="1" />
      <button id="cancel" label="Cancel" />
      <button id="next" label="Next" primary="true" />
    </hbox>
  </vbox>
</window>
```

### Interpolation Features

- `${key}` - Simple variable substitution
- `${key:default}` - With default value if empty
- Arrays joined with commas: `${messages}` → "msg1,msg2,msg3"

## Implemented Components

| Component | Description | Key Attributes |
|-----------|-------------|----------------|
| `<window>` | Framed window with border | `title`, `width`, `height`, `x`, `y`, `center`, `padding`, `border` |
| `<vbox>` | Vertical box layout | `gap`, `align` (left/center/right) |
| `<hbox>` | Horizontal box layout | `gap`, `align` (top/center/bottom) |
| `<grid>` | Grid layout | `cols`, `rows`, `gap` |
| `<text>` | Static text | `value`, `color`, `bold` |
| `<list>` | Bulleted list | `items` (comma-separated) |
| `<button>` | Clickable button | `id`, `label`, `primary` |
| `<input>` | Text input | `id`, `placeholder`, `hidden` |
| `<spacer>` | Flexible spacing | `width`, `height`, `flex` |
| `<divider>` | Horizontal line | `char` |
| `<progress>` | Progress bar | `value` (0-100), `width` |
| `<toast>` | Notification | `message`, `type`, `autoHide` |
| `<tabbar>` | Tab navigation | `tabs`, `active` |
| `<messagelist>` | Chat messages | `messages` |
| `<selectionlist>` | Selectable list | `items`, `selected` |
| `<dialog>` | Modal dialog | `title`, `content` |

## ResponsiveEngine[R Renderer]

Adds responsive sizing and visibility management:

```go
re := template.NewResponsiveEngine[BubbleRenderer]()

// Size observation
re.AddSizeObserver(func(w, h int) {
    fmt.Printf("Size: %dx%d\n", w, h)
})

// Visibility control
re.Hide("sidebar")
re.Show("dialog")
re.Toggle("panel")

// Conditional rendering with constraints
// <text show-if="width>=80" />
// <panel hide-if="height<30" />
```

## Integration with Model

The TUI model uses the reactive engine:

```go
type Model struct {
    engine *template.ReactiveEngine[BubbleRenderer]
    view   string
    width  int
    height int
}

func (m Model) View() string {
    // Sync state to engine
    m.engine.Set("messages", m.messages)
    m.engine.Set("input", m.inputText)

    // Render current view template
    return m.engine.RenderTemplate(m.view, m.width, m.height, renderer)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // On state change, mark engine dirty
    m.engine.Set("session", m.sessionName)
    return m, nil
}
```

## File Structure

```
internal/ui/
├── template/
│   ├── parser.go        # Template parser, interpolates ${var} syntax
│   ├── engine.go        # Engine[R] - core template rendering
│   ├── responsive.go    # ResponsiveEngine - sizing, visibility constraints
│   ├── reactive_engine.go # ReactiveEngine - state management, dirty tracking
│   ├── loader.go        # Loads templates from filesystem
│   └── reactive_engine_test.go
├── renderer/
│   ├── renderer.go      # Renderer interface
│   ├── bubble.go        # BubbleRenderer (lipgloss)
│   ├── lcd.go           # LCDRenderer
│   └── headless.go      # HeadlessRenderer
└── model.go             # Main TUI model, uses ReactiveEngine
```

## Testing

```bash
go test ./internal/ui/template/... -v
```

## Design Decisions

1. **Generic over Renderer** - Components work with any renderer implementation
2. **State-driven** - React-like `Set(key, value)` triggers re-render
3. **Dirty tracking** - Only re-render changed components for performance
4. **External templates** - Templates stored as files, not hardcoded
5. **Thread-safe** - RWMutex protects state access
