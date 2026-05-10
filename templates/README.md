# TUI Template System

Declarative TUI layout system inspired by HTML/XML.

## Directory Structure

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

## Available Components

### Containers

| Component | Description | Attributes |
|-----------|-------------|------------|
| `<window>` | Framed window with border and optional title | `title`, `width`, `height`, `x`, `y`, `center`, `padding`, `border` |
| `<vbox>` | Vertical box layout (top-to-bottom) | `gap`, `align` (left/center/right) |
| `<hbox>` | Horizontal box layout (left-to-right) | `gap`, `align` (top/center/bottom) |
| `<grid>` | Grid layout with cols/rows | `cols`, `rows`, `gap` |

### Elements

| Component | Description | Attributes |
|-----------|-------------|------------|
| `<text>` | Static text display | `value`, `color`, `bold` |
| `<list>` | Non-selectable list with bullets | `items` (comma-separated) |
| `<button>` | Clickable button | `id`, `label`, `primary` |
| `<input>` | Text input field | `id`, `placeholder`, `hidden` |
| `<spacer>` | Flexible spacing | `width`, `height`, `flex` |
| `<divider>` | Horizontal line | `char` |
| `<progress>` | Progress bar | `value` (0-100), `width` |
| `<toast>` | Notification message | `message`, `type` (info/success/warning/error) |
| `<tabbar>` | Tab navigation | `tabs`, `active` |
| `<messagelist>` | Chat message list | `messages` |
| `<selectionlist>` | Selectable list | `items`, `selected` |
| `<dialog>` | Modal dialog | `title`, `content` |

## Usage

### Parsing a Template

```go
engine := template.NewEngine[BubbleRenderer]()
src := `<window title="Setup" width="60" height="20" center="true">
  <vbox gap="1">
    <text value="Hello ${name}!" />
    <button id="next" label="Next" primary="true" />
  </vbox>
</window>`

result, err := engine.ParseAndRender(src, 80, 24, renderer)
```

### Setting Variables

```go
tmpl, err := template.Parse(src)
tmpl.Set("name", "Freecode")
tmpl.Set("count", 42)
```

### Getting Components by ID

```go
engine := template.NewEngine[BubbleRenderer]()
engine.ParseAndRender(src, 80, 24, renderer)

btn := engine.GetComponent("next")
if btn != nil {
    fmt.Printf("Button: %s at (%d,%d)\n", btn.Content, btn.X, btn.Y)
}
```

### Variable Interpolation

Variables in attributes use `${varName}` syntax:

```xml
<text value="${greeting}" color="${textColor:#ffffff}" />
```

Variables in content use the same syntax:

```xml
<text>Hello ${name}!</text>
```

### Flex Layout

Use `flex` on spacers to fill remaining space:

```xml
<vbox>
  <text value="Header" />
  <spacer flex="1" />  <!-- Takes all remaining space -->
  <text value="Footer" />
</vbox>
```

### Centering

```xml
<window center="true">
  <text value="Centered" />
</window>
```

### Nesting

```xml
<vbox>
  <hbox>
    <text value="Left" />
    <spacer width="1" />
    <text value="Right" />
  </hbox>
  <divider />
  <text value="Bottom" />
</vbox>
```

## Example: Setup Dialog

```xml
<window title="Setup" width="60" height="20" center="true" padding="2">
  <vbox gap="1">
    <text value="Welcome to Freecode Setup" bold="true" />
    <divider />
    <vbox gap="0">
      <text value="Provider:" />
      <input id="provider" placeholder="Select a provider..." />
    </vbox>
    <spacer flex="1" />
    <hbox gap="2">
      <spacer flex="1" />
      <button id="cancel" label="Cancel" />
      <button id="next" label="Next" primary="true" />
    </hbox>
  </vbox>
</window>
```

## Rendering

The engine renders to any type implementing the `Renderer` interface:

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

Implementations:
- `BubbleRenderer` - lipgloss-based styling
- `LCDRenderer` - simple terminal
- `HeadlessRenderer` - testing/no-op
