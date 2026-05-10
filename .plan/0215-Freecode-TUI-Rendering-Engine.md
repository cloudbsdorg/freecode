# Freecode TUI Rendering Engine

## Overview

The TUI rendering engine provides a declarative, component-based UI system with HTML/XML-like templates, responsive sizing, and dynamic show/hide capabilities. It supports multiple renderers (Bubble/lipgloss, LCD, Headless) through a generic interface.

## Architecture

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

## Lipgloss Compatibility

The template engine is designed to work with lipgloss-style positioning. Multi-line text (like banners) is properly rendered by splitting into individual lines at successive y positions:

```go
func (e *Engine[R]) renderMultiLineText(text string, x int, y int, color string, r R) string {
    lines := strings.Split(text, "\n")
    var result string
    for i, line := range lines {
        result += r.RenderText(line, x, y+i, color)
    }
    return result
}
```

## Core Interfaces

### Renderer Interface

```go
type Renderer interface {
    RenderBox(x, y, w, h int, bgColor string) string
    RenderText(text string, x, y int, fgColor string) string
    RenderBorder(x, y, w, h int, fgColor string) string
    RenderSelected(text string, x, y, w int, fg, bg string) string
    Width() int
    Height() int
}
```

**Implementations:**
- `BubbleRenderer` - lipgloss-based styling for rich terminals
- `LCDRenderer` - simple terminal rendering
- `HeadlessRenderer` - testing/no-op

### Component Interface

```go
type Component[R renderer.Renderer] struct {
    X, Y       int
    Width, Height int
    Visible    bool
}

func (c *Component[R]) Render(r R) string
```

## Template System

### Syntax

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

### Available Components

#### Containers

| Component | Description | Key Attributes |
|-----------|-------------|----------------|
| `<window>` | Framed window with border | `title`, `width`, `height`, `x`, `y`, `center`, `padding`, `border` |
| `<vbox>` | Vertical box layout | `gap`, `align` (left/center/right) |
| `<hbox>` | Horizontal box layout | `gap`, `align` (top/center/bottom) |
| `<grid>` | Grid layout | `cols`, `rows`, `gap` |

#### Elements

| Component | Description | Key Attributes |
|-----------|-------------|----------------|
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

### Constraints & Conditions

```xml
<!-- Size constraints -->
<text min-width="50" max-width="100" />

<!-- Percentage of parent -->
<text width="50%" />

<!-- Conditional visibility -->
<sidebar show-if="width>=100" />
<panel hide-if="height<40" />
```

Supported conditions: `width>=N`, `width<=N`, `width>N`, `width<N`, `width==N` (same for `height`)

## Engine Classes

### Engine[R Renderer]

Core template parser and renderer:

```go
engine := NewEngine[BubbleRenderer]()

// Parse and render in one step
result, err := engine.ParseAndRender(src, 80, 24, renderer)

// Get component by ID after rendering
btn := engine.GetComponent("submit")
```

### ResponsiveEngine[R Renderer]

Adds responsive sizing and visibility management:

```go
re := NewResponsiveEngine[BubbleRenderer]()

// Size observation
re.AddSizeObserver(SizeObserverFunc(func(w, h int) {
    fmt.Printf("Size: %dx%d\n", w, h)
}))

// Visibility control
re.Hide("sidebar")
re.Show("dialog")
re.Toggle("panel")

// Check visibility
if re.IsVisible("dialog") { ... }

// Conditional rendering
re.RenderAt(src, width, height, renderer)
```

### DialogPresenter

For dynamic dialog content:

```go
// Alert dialog
alert := NewAlertPresenter(re, "alert", "msg")
alert.ShowError("File not found!")
alert.ShowSuccess("Saved!")
alert.ShowWarning("Low memory")
alert.ShowInfo("Update available")

// Confirm dialog
confirm := NewConfirmPresenter(re, "confirm", "msg", "ok", "cancel")
confirm.SetConfirmHandler(func() { fmt.Println("OK!") })
confirm.SetCancelHandler(func() { fmt.Println("Cancelled") })
confirm.ShowConfirm("Delete?", "Remove this file?", "Delete", "Cancel")
```

## Usage Patterns

### Basic View

```go
engine := template.NewEngine[BubbleRenderer]()
src := `
<window title="My App" width="80" height="24">
  <vbox padding="1">
    <text value="Hello World" bold="true" />
    <spacer height="1" />
    <list items="Option A,Option B,Option C" />
  </vbox>
</window>`

output, err := engine.ParseAndRender(src, 80, 24, renderer)
```

### Dynamic Content

```go
engine := template.NewResponsiveEngine[BubbleRenderer]()

// Set variables before rendering
tmpl, _ := template.Parse(src)
tmpl.Set("username", "Alice")
tmpl.Set("items", "a,b,c,d")

// Or use component manipulation
engine.SetComponentAttr("title", "value", "New Title")
engine.SetComponentContent("list", []string{"x", "y", "z"})
```

### Integration with Bubble Tea

```go
type Model struct {
    engine *template.ResponsiveEngine[BubbleRenderer]
    view   string
}

func (m Model) View() string {
    output, _ := m.engine.RenderAt(m.view, m.width, m.height, renderer)
    return output
}
```

## File Structure

```
internal/ui/
├── renderer/
│   ├── renderer.go       # Renderer interface
│   ├── bubble.go         # BubbleRenderer (lipgloss)
│   ├── lcd.go           # LCDRenderer
│   └── headless.go       # HeadlessRenderer
├── component/
│   ├── component.go      # Base Component[R]
│   ├── statusbar.go     # StatusBar
│   ├── button.go        # Button
│   ├── textinput.go     # TextInput
│   ├── selectionlist.go  # SelectionList
│   ├── list.go          # List
│   ├── dialog.go        # Dialog
│   ├── window.go        # Window
│   ├── toast.go         # Toast
│   ├── messagelist.go   # MessageList
│   ├── tabbar.go        # TabBar
│   └── inputarea.go     # InputArea
└── template/
    ├── parser.go        # Template parser
    ├── engine.go        # Engine[R]
    ├── responsive.go    # ResponsiveEngine
    ├── dialog_presenter.go # Dialog/Alert/Confirm
    └── README.md        # Template documentation
```

## Testing

```bash
go test ./internal/ui/template/... -v
go test ./internal/ui/component/... -v
```

Mock renderer for testing:

```go
type mockRenderer struct{}
func (m mockRenderer) RenderBox(x, y, w, h int, bgColor string) string { return "" }
func (m mockRenderer) RenderText(text string, x, y int, fgColor string) string { return text }
func (m mockRenderer) RenderBorder(x, y, w, h int, fgColor string) string { return "" }
func (m mockRenderer) RenderSelected(text string, x, y, w int, fg, bg string) string { return text }
func (m mockRenderer) Width() int { return 80 }
func (m mockRenderer) Height() int { return 24 }
```

## Design Decisions

1. **Generic over Renderer** - Components and engines are generic `type Component[R Renderer]` allowing same code to work with different renderers
2. **Interface-based DialogEngine** - Allows presenters to work with any engine implementing Show/Hide/SetComponentAttr
3. **String-based attributes** - Template attributes stored as strings, parsed at render time for flexibility
4. **Responsive via observers** - Size changes propagate through observer pattern, not direct coupling
5. **ID-based component access** - Components accessed by string ID for loose coupling between template and code
