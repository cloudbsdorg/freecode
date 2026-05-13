package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/freecode/freecode/internal/style"
)

type TimelineNode struct {
	ID           string
	Title        string
	Timestamp    time.Time
	MessageCount int
	ParentID     string
	IsFork       bool
	Model        string
	Agent        string
	Children     []*TimelineNode
}

type TimelineDialog struct {
	width        int
	height       int
	isOpen       bool
	selected     int
	filter       string
	filtered     []*TimelineNode
	allNodes     []*TimelineNode
	onSelect     func(string)
	onFork       func(string)
	scrollOffset int
	preview      *TimelineNode
}

func NewTimelineDialog() *TimelineDialog {
	return &TimelineDialog{
		width:        70,
		height:       25,
		isOpen:       false,
		selected:     0,
		filter:       "",
		filtered:     []*TimelineNode{},
		allNodes:     []*TimelineNode{},
		scrollOffset: 0,
	}
}

func (t *TimelineDialog) Open() {
	t.isOpen = true
	t.selected = 0
	t.scrollOffset = 0
	t.filter = ""
	t.applyFilter()
}

func (t *TimelineDialog) Close() {
	t.isOpen = false
	t.preview = nil
}

func (t *TimelineDialog) IsOpen() bool {
	return t.isOpen
}

func (t *TimelineDialog) SetWidth(w int) {
	t.width = w
}

func (t *TimelineDialog) SetHeight(h int) {
	t.height = h
}

func (t *TimelineDialog) SetOnSelect(fn func(string)) {
	t.onSelect = fn
}

func (t *TimelineDialog) SetOnFork(fn func(string)) {
	t.onFork = fn
}

func (t *TimelineDialog) SetNodes(nodes []*TimelineNode) {
	t.allNodes = nodes
	t.applyFilter()
}

func (t *TimelineDialog) SetFilter(filter string) {
	t.filter = filter
	t.applyFilter()
}

func (t *TimelineDialog) applyFilter() {
	if t.filter == "" {
		t.filtered = t.allNodes
	} else {
		needle := strings.ToLower(t.filter)
		var result []*TimelineNode
		for _, node := range t.allNodes {
			if strings.Contains(strings.ToLower(node.Title), needle) {
				result = append(result, node)
			} else if strings.Contains(strings.ToLower(node.Agent), needle) {
				result = append(result, node)
			} else if strings.Contains(strings.ToLower(node.Model), needle) {
				result = append(result, node)
			}
		}
		t.filtered = result
	}
	if t.selected >= len(t.filtered) {
		t.selected = 0
	}
	t.scrollOffset = 0
}

func (t *TimelineDialog) Next() {
	if len(t.filtered) == 0 {
		return
	}
	t.selected = (t.selected + 1) % len(t.filtered)
	t.ensureVisible()
}

func (t *TimelineDialog) Prev() {
	if len(t.filtered) == 0 {
		return
	}
	t.selected = (t.selected - 1 + len(t.filtered)) % len(t.filtered)
	t.ensureVisible()
}

func (t *TimelineDialog) ensureVisible() {
	if t.selected < t.scrollOffset {
		t.scrollOffset = t.selected
	}
	maxVisible := t.height - 10
	if t.selected >= t.scrollOffset+maxVisible {
		t.scrollOffset = t.selected - maxVisible + 1
	}
}

func (t *TimelineDialog) GetSelected() *TimelineNode {
	if t.selected < 0 || t.selected >= len(t.filtered) {
		return nil
	}
	return t.filtered[t.selected]
}

func (t *TimelineDialog) ScrollDown() {
	maxVisible := t.height - 10
	if t.scrollOffset+maxVisible < len(t.filtered) {
		t.scrollOffset++
	}
}

func (t *TimelineDialog) ScrollUp() {
	if t.scrollOffset > 0 {
		t.scrollOffset--
	}
}

func (t *TimelineDialog) HandleKey(msg string) bool {
	if !t.isOpen {
		return false
	}

	switch msg {
	case "escape":
		t.Close()
		return true
	case "enter":
		sel := t.GetSelected()
		if sel != nil && t.onSelect != nil {
			t.onSelect(sel.ID)
		}
		return true
	case "f":
		sel := t.GetSelected()
		if sel != nil && t.onFork != nil {
			t.onFork(sel.ID)
		}
		return true
	case "up", "k":
		t.Prev()
		return true
	case "down", "j":
		t.Next()
		return true
	case "pgup":
		t.ScrollUp()
		return true
	case "pgdown":
		t.ScrollDown()
		return true
	case "g":
		if len(t.filtered) > 0 {
			t.selected = 0
			t.scrollOffset = 0
		}
		return true
	case "G":
		if len(t.filtered) > 0 {
			t.selected = len(t.filtered) - 1
			t.ensureVisible()
		}
		return true
	}
	return false
}

func (t *TimelineDialog) Render() string {
	if !t.isOpen {
		return ""
	}

	dialogStyle := style.NewStyle().
		Background(style.Color("#1E1E1E")).
		BorderStyle(style.RoundedBorder()).
		Width(t.width).
		Height(t.height)

	return dialogStyle.Render(t.renderContent())
}

func (t *TimelineDialog) renderContent() string {
	var lines []string

	lines = append(lines, t.renderHeader())
	lines = append(lines, "")
	lines = append(lines, t.renderFilter())
	lines = append(lines, "")
	lines = append(lines, t.renderTimeline()...)
	lines = append(lines, "")
	lines = append(lines, t.renderFooter())

	return strings.Join(lines, "\n")
}

func (t *TimelineDialog) renderHeader() string {
	headerStyle := style.NewStyle().
		Foreground(style.Color("#FFFFFF")).
		Bold(true)
	return headerStyle.Render("Session Timeline")
}

func (t *TimelineDialog) renderFilter() string {
	filterStyle := style.NewStyle().
		Background(style.Color("#3C3C3C")).
		Foreground(style.Color("#E0E0E0")).
		Padding(0, 1)
	placeholder := "Filter by title, agent, or model..."
	if t.filter != "" {
		placeholder = t.filter
	}
	return filterStyle.Render("Search: " + placeholder + "_")
}

func (t *TimelineDialog) renderTimeline() []string {
	var lines []string

	if len(t.filtered) == 0 {
		noSessionsStyle := style.NewStyle().Foreground(style.Color("#808080"))
		lines = append(lines, "  "+noSessionsStyle.Render("No sessions found"))
		return lines
	}

	maxVisible := t.height - 12
	endIdx := t.scrollOffset + maxVisible
	if endIdx > len(t.filtered) {
		endIdx = len(t.filtered)
	}

	visibleNodes := t.filtered[t.scrollOffset:endIdx]

	for i, node := range visibleNodes {
		actualIdx := t.scrollOffset + i
		lines = append(lines, t.renderNode(node, actualIdx == t.selected))
	}

	totalVisible := t.scrollOffset + maxVisible
	if totalVisible < len(t.filtered) {
		moreStyle := style.NewStyle().Foreground(style.Color("#606060"))
		lines = append(lines, "  "+moreStyle.Render(fmt.Sprintf("▼ %d more sessions", len(t.filtered)-totalVisible)))
	} else if t.scrollOffset > 0 {
		moreStyle := style.NewStyle().Foreground(style.Color("#606060"))
		lines = append(lines, "  "+moreStyle.Render(fmt.Sprintf("▲ %d hidden", t.scrollOffset)))
	}

	return lines
}

func (t *TimelineDialog) renderNode(node *TimelineNode, selected bool) string {
	var b strings.Builder

	var branchChar string
	if node.IsFork {
		branchChar = "├"
	} else {
		branchChar = "│"
	}

	var contentStyle style.Style
	if selected {
		contentStyle = style.NewStyle().
			Background(style.Color("#007ACC")).
			Foreground(style.Color("#FFFFFF"))
	} else {
		contentStyle = style.NewStyle().
			Foreground(style.Color("#E0E0E0"))
	}

	ts := formatTimestamp(node.Timestamp)
	tsStyle := style.NewStyle().Foreground(style.Color("#606060"))

	title := node.Title
	if len(title) > 40 {
		title = title[:37] + "..."
	}

	msgCount := fmt.Sprintf("%d msgs", node.MessageCount)
	msgStyle := style.NewStyle().Foreground(style.Color("#808080"))

	metaStyle := style.NewStyle().Foreground(style.Color("#606060"))
	meta := fmt.Sprintf("%s · %s", node.Agent, node.Model)

	forkIndicator := ""
	if node.IsFork {
		forkStyle := style.NewStyle().
			Foreground(style.Color("#FFCC00")).
			Bold(true)
		forkIndicator = " " + forkStyle.Render("FORK")
	}

	b.WriteString(" ")
	b.WriteString(tsStyle.Render(branchChar))
	b.WriteString(" ")
	b.WriteString(contentStyle.Render(fmt.Sprintf("[%s]", ts)))
	b.WriteString(" ")
	b.WriteString(contentStyle.Render(title))
	b.WriteString(forkIndicator)
	b.WriteString(" ")
	b.WriteString(msgStyle.Render(msgCount))
	b.WriteString("\n")

	b.WriteString(" ")
	b.WriteString(tsStyle.Render(" "))
	b.WriteString(" ")
	b.WriteString(metaStyle.Render(meta))

	return b.String()
}

func (t *TimelineDialog) renderFooter() string {
	hintStyle := style.NewStyle().Foreground(style.Color("#808080"))

	if t.preview != nil {
		previewStyle := style.NewStyle().Foreground(style.Color("#4EC9B0"))
		return hintStyle.Render("↑↓ navigate  enter select  f fork  esc close  | Preview: ") +
			previewStyle.Render(t.preview.Title)
	}

	return hintStyle.Render("↑↓ navigate  enter select  f fork  esc close")
}

type ForkDialog struct {
	width          int
	height         int
	isOpen         bool
	parentSession  *TimelineNode
	newSessionName string
	cursorPos      int
	onConfirm      func(string, string)
	onCancel       func()
	confirmFocused bool
}

func NewForkDialog() *ForkDialog {
	return &ForkDialog{
		width:          60,
		height:         15,
		isOpen:         false,
		newSessionName: "",
		cursorPos:      0,
		confirmFocused: false,
	}
}

func (f *ForkDialog) Open(session *TimelineNode) {
	f.isOpen = true
	f.parentSession = session
	f.newSessionName = session.Title + " (fork)"
	f.cursorPos = len(f.newSessionName)
	f.confirmFocused = false
}

func (f *ForkDialog) Close() {
	f.isOpen = false
	f.parentSession = nil
	f.newSessionName = ""
}

func (f *ForkDialog) IsOpen() bool {
	return f.isOpen
}

func (f *ForkDialog) SetWidth(w int) {
	f.width = w
}

func (f *ForkDialog) SetHeight(h int) {
	f.height = h
}

func (f *ForkDialog) SetOnConfirm(fn func(string, string)) {
	f.onConfirm = fn
}

func (f *ForkDialog) SetOnCancel(fn func()) {
	f.onCancel = fn
}

func (f *ForkDialog) HandleKey(msg string) bool {
	if !f.isOpen {
		return false
	}

	switch msg {
	case "escape":
		f.Close()
		if f.onCancel != nil {
			f.onCancel()
		}
		return true
	case "enter":
		if f.onConfirm != nil && f.parentSession != nil {
			f.onConfirm(f.parentSession.ID, f.newSessionName)
		}
		f.Close()
		return true
	case "tab":
		f.confirmFocused = !f.confirmFocused
		return true
	case "left", "h":
		if f.cursorPos > 0 {
			f.cursorPos--
		}
		return true
	case "right", "l":
		if f.cursorPos < len(f.newSessionName) {
			f.cursorPos++
		}
		return true
	case "backspace":
		if f.cursorPos > 0 && len(f.newSessionName) > 0 {
			f.newSessionName = f.newSessionName[:f.cursorPos-1] + f.newSessionName[f.cursorPos:]
			f.cursorPos--
		}
		return true
	case "delete":
		if f.cursorPos < len(f.newSessionName) {
			f.newSessionName = f.newSessionName[:f.cursorPos] + f.newSessionName[f.cursorPos+1:]
		}
		return true
	case "home":
		f.cursorPos = 0
		return true
	case "end":
		f.cursorPos = len(f.newSessionName)
		return true
	}

	if len(msg) == 1 {
		ch := msg[0]
		if ch >= 32 && ch < 127 {
			f.newSessionName = f.newSessionName[:f.cursorPos] + string(ch) + f.newSessionName[f.cursorPos:]
			f.cursorPos++
		}
		return true
	}

	return false
}

func (f *ForkDialog) Render() string {
	if !f.isOpen {
		return ""
	}

	dialogStyle := style.NewStyle().
		Background(style.Color("#1E1E1E")).
		BorderStyle(style.RoundedBorder()).
		Width(f.width).
		Height(f.height)

	return dialogStyle.Render(f.renderContent())
}

func (f *ForkDialog) renderContent() string {
	var lines []string

	lines = append(lines, f.renderHeader())
	lines = append(lines, "")
	lines = append(lines, f.renderSessionInfo())
	lines = append(lines, "")
	lines = append(lines, f.renderNameInput())
	lines = append(lines, "")
	lines = append(lines, f.renderButtons())
	lines = append(lines, "")
	lines = append(lines, f.renderHints())

	return strings.Join(lines, "\n")
}

func (f *ForkDialog) renderHeader() string {
	headerStyle := style.NewStyle().
		Foreground(style.Color("#FFFFFF")).
		Bold(true)
	forkIcon := style.NewStyle().
		Foreground(style.Color("#FFCC00")).
		Render("⎇")
	return headerStyle.Render(forkIcon + " Fork Session")
}

func (f *ForkDialog) renderSessionInfo() string {
	if f.parentSession == nil {
		return ""
	}

	infoStyle := style.NewStyle().Foreground(style.Color("#808080"))
	contentStyle := style.NewStyle().Foreground(style.Color("#E0E0E0"))

	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(infoStyle.Render("Parent:"))
	b.WriteString(" ")
	b.WriteString(contentStyle.Render(f.parentSession.Title))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(infoStyle.Render("Messages:"))
	b.WriteString(" ")
	b.WriteString(contentStyle.Render(fmt.Sprintf("%d", f.parentSession.MessageCount)))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(infoStyle.Render("Created:"))
	b.WriteString(" ")
	b.WriteString(contentStyle.Render(formatTimestamp(f.parentSession.Timestamp)))

	return b.String()
}

func (f *ForkDialog) renderNameInput() string {
	labelStyle := style.NewStyle().Foreground(style.Color("#808080"))
	inputBgStyle := style.NewStyle().
		Background(style.Color("#3C3C3C")).
		Foreground(style.Color("#E0E0E0"))

	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(labelStyle.Render("New session name:"))
	b.WriteString("\n")
	b.WriteString("  ")

	inputWidth := f.width - 6

	b.WriteString(inputBgStyle.Render(" " + f.newSessionName))

	if f.confirmFocused {
		b.WriteString(inputBgStyle.Render(" "))
	} else {
		cursorBg := style.NewStyle().
			Background(style.Color("#007ACC")).
			Foreground(style.Color("#FFFFFF"))
		if f.cursorPos >= len(f.newSessionName) {
			b.WriteString(cursorBg.Render(" "))
		} else {
			char := string(f.newSessionName[f.cursorPos])
			b.WriteString(cursorBg.Render(char))
		}
	}

	for i := len(f.newSessionName) + 1; i < inputWidth-1; i++ {
		b.WriteString(inputBgStyle.Render(" "))
	}
	b.WriteString(inputBgStyle.Render(" "))

	return b.String()
}

func (f *ForkDialog) renderButtons() string {
	cancelStyle := style.NewStyle().
		Foreground(style.Color("#E0E0E0"))
	confirmStyle := style.NewStyle().
		Foreground(style.Color("#FFFFFF"))

	if f.confirmFocused {
		cancelStyle = style.NewStyle().
			Foreground(style.Color("#808080"))
		confirmStyle = style.NewStyle().
			Background(style.Color("#007ACC")).
			Foreground(style.Color("#FFFFFF")).
			Bold(true)
	}

	space := strings.Repeat(" ", (f.width-40)/2)
	cancelRendered := cancelStyle.Render("[ Cancel ]")
	confirmRendered := confirmStyle.Render("[ Fork Session ]")
	return "  " + cancelRendered + space + confirmRendered
}

func (f *ForkDialog) renderHints() string {
	hintStyle := style.NewStyle().Foreground(style.Color("#606060"))
	return hintStyle.Render("←→ move cursor  backspace delete  tab toggle  enter confirm  esc cancel")
}
