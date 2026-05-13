package ui

import (
	"strings"

	"github.com/freecode/freecode/internal/style"
)

type ToolState struct {
	Status  string
	Title   string
	Subtitle string
	Args    map[string]string
	Error   string
}

type ToolCallComponent struct {
	width    int
	tools    []ToolState
	selected int
}

func NewToolCallComponent() *ToolCallComponent {
	return &ToolCallComponent{
		width:    80,
		tools:    make([]ToolState, 0),
		selected: -1,
	}
}

func (t *ToolCallComponent) SetWidth(w int) {
	t.width = w
}

func (t *ToolCallComponent) SetTools(tools []ToolState) {
	t.tools = tools
	if t.selected >= len(tools) {
		t.selected = len(tools) - 1
	}
}

func (t *ToolCallComponent) SelectNext() {
	if t.selected < len(t.tools)-1 {
		t.selected++
	}
}

func (t *ToolCallComponent) SelectPrev() {
	if t.selected > 0 {
		t.selected--
	}
}

func (t *ToolCallComponent) Selected() *ToolState {
	if t.selected >= 0 && t.selected < len(t.tools) {
		return &t.tools[t.selected]
	}
	return nil
}

func (t *ToolCallComponent) Render() string {
	if len(t.tools) == 0 {
		return ""
	}

	var result strings.Builder

	for i, tool := range t.tools {
		if i > 0 {
			result.WriteString("\n")
		}

		isSelected := i == t.selected
		result.WriteString(t.renderTool(tool, isSelected))
	}

	return result.String()
}

func (t *ToolCallComponent) renderTool(tool ToolState, selected bool) string {
	var result strings.Builder

	borderStyle := style.NewStyle().
		Foreground(style.Color("#404040"))

	if selected {
		borderStyle = style.NewStyle().
			Foreground(style.Color("#4EC9B0"))
	}

	result.WriteString(borderStyle.Render("┌─"))

	nameStyle := style.NewStyle().
		Foreground(style.Color("#DCDCAA"))
	if selected {
		nameStyle = nameStyle.Bold(true)
	}

	statusIcon := "⏳"
	statusStyle := ToolRunningStyle

	switch tool.Status {
	case ToolStatusCompleted:
		statusIcon = "✅"
		statusStyle = ToolSuccessStyle
	case ToolStatusError:
		statusIcon = "❌"
		statusStyle = ToolErrorStyle
	case ToolStatusRunning:
		statusIcon = "⚡"
		statusStyle = ToolRunningStyle
	}

	result.WriteString(nameStyle.Render(" "+tool.Title))
	result.WriteString(" ")
	result.WriteString(statusStyle.Render(statusIcon))
	result.WriteString(" ")
	result.WriteString(statusStyle.Render(tool.Status))

	result.WriteString(borderStyle.Render(" ─"))
	result.WriteString("\n")

	if tool.Subtitle != "" {
		result.WriteString(borderStyle.Render("│ "))
		result.WriteString(ToolArgsStyle.Render(tool.Subtitle))
		result.WriteString("\n")
	}

	for key, value := range tool.Args {
		result.WriteString(borderStyle.Render("│ "))
		result.WriteString(ToolCallStyle.Render(key + ": "))
		result.WriteString(ToolArgsStyle.Render(value))
		result.WriteString("\n")
	}

	if tool.Error != "" {
		result.WriteString(borderStyle.Render("│ "))
		result.WriteString(ToolErrorStyle.Render("Error: " + tool.Error))
		result.WriteString("\n")
	}

	result.WriteString(borderStyle.Render("└─"))
	for i := 0; i < t.width-4; i++ {
		result.WriteString("─")
	}

	return result.String()
}

type ToolCallView struct {
	collapsible *CollapsibleView
	info        ToolInfo
}

type CollapsibleView struct {
	open       bool
	header     string
	content    string
	onOpen     func()
	onClose    func()
}

func NewCollapsibleView() *CollapsibleView {
	return &CollapsibleView{
		open:    false,
		header:  "",
		content: "",
	}
}

func (c *CollapsibleView) SetHeader(header string) {
	c.header = header
}

func (c *CollapsibleView) SetContent(content string) {
	c.content = content
}

func (c *CollapsibleView) Toggle() {
	c.open = !c.open
}

func (c *CollapsibleView) IsOpen() bool {
	return c.open
}

func (c *CollapsibleView) Render() string {
	if c.open {
		return c.renderOpen()
	}
	return c.renderCollapsed()
}

func (c *CollapsibleView) renderCollapsed() string {
	return style.NewStyle().
		Foreground(style.Color("#808080")).
		Render("▶ " + c.header + " (click to expand)")
}

func (c *CollapsibleView) renderOpen() string {
	var result strings.Builder
	result.WriteString(style.NewStyle().
		Foreground(style.Color("#808080")).
		Render("▼ " + c.header))
	result.WriteString("\n")
	result.WriteString(c.content)
	return result.String()
}
