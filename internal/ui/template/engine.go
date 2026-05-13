package template

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/freecode/freecode/internal/renderer"
)

type Engine[R renderer.Renderer] struct {
	components map[string]*ComponentNode
	theme      *ThemeConfig
}

type ComponentNode struct {
	ID       string
	Type     ElementType
	X        int
	Y        int
	Width    int
	Height   int
	Content  interface{}
	Children []*ComponentNode
	Attrs    map[string]string
	Visible  bool
}

type ThemeConfig struct {
	Styles map[string]StyleConfig
}

type StyleConfig struct {
	Foreground string
	Background string
	Bold       bool
	Italic     bool
}

func NewEngine[R renderer.Renderer]() *Engine[R] {
	return &Engine[R]{
		components: make(map[string]*ComponentNode),
	}
}

func (e *Engine[R]) SetTheme(theme *ThemeConfig) {
	e.theme = theme
}

func (e *Engine[R]) ParseAndRender(src string, width, height int, r R) (string, error) {
	tmpl, err := Parse(src)
	if err != nil {
		return "", err
	}

	nodes, err := e.BuildTree(tmpl.Root, 0, 0, width, height)
	if err != nil {
		return "", err
	}

	e.resolveVisibility(nodes)

	return e.renderNodes(nodes, r), nil
}

func (e *Engine[R]) resolveVisibility(nodes []*ComponentNode) {
	for _, n := range nodes {
		n.Visible = true
		e.resolveVisibility(n.Children)
	}
}

func (e *Engine[R]) IsComponentVisible(id string) bool {
	node := e.components[id]
	if node == nil {
		return true
	}
	return node.Visible
}

func (e *Engine[R]) BuildTree(elem *Element, x, y, w, h int) ([]*ComponentNode, error) {
	if elem == nil {
		return nil, nil
	}

	switch elem.Type {
	case "var":
		return nil, nil

	case TypeWindow:
		winW := getAttrInt(elem.Attributes, "width", w)
		winH := getAttrInt(elem.Attributes, "height", h)
		winX := getAttrInt(elem.Attributes, "x", x)
		winY := getAttrInt(elem.Attributes, "y", y)

		if getAttrBool(elem.Attributes, "center", false) {
			winX = x + (w-winW)/2
			winY = y + (h-winH)/2
		}

		padding := getAttrInt(elem.Attributes, "padding", 0)
		children, err := e.BuildTreeFromElements(elem.Children, winX+padding, winY+padding, winW-2*padding, winH-2*padding)
		if err != nil {
			return nil, err
		}

		node := &ComponentNode{
			ID:       elem.ID,
			Type:     TypeWindow,
			X:        winX,
			Y:        winY,
			Width:    winW,
			Height:   winH,
			Children: children,
			Attrs:    elem.Attributes,
		}
		if elem.ID != "" {
			e.components[elem.ID] = node
		}
		return []*ComponentNode{node}, nil

	case TypeVBox:
		padding := getAttrInt(elem.Attributes, "padding", 0)
		gap := getAttrInt(elem.Attributes, "gap", 0)
		align := elem.Attributes["align"]

		children, err := e.BuildTreeFromElements(elem.Children, x+padding, y+padding, w-2*padding, h-2*padding)
		if err != nil {
			return nil, err
		}

		e.layoutVBox(children, w-2*padding, gap, align)

		return children, nil

	case TypeHBox:
		padding := getAttrInt(elem.Attributes, "padding", 0)
		gap := getAttrInt(elem.Attributes, "gap", 0)
		align := elem.Attributes["align"]

		children, err := e.BuildTreeFromElements(elem.Children, x+padding, y+padding, w-2*padding, h-2*padding)
		if err != nil {
			return nil, err
		}

		e.layoutHBox(children, h-2*padding, gap, align)

		return children, nil

	case TypeGrid:
		cols := getAttrInt(elem.Attributes, "cols", 2)
		rows := getAttrInt(elem.Attributes, "rows", 0)
		gap := getAttrInt(elem.Attributes, "gap", 0)

		children, err := e.BuildTreeFromElements(elem.Children, x, y, w, h)
		if err != nil {
			return nil, err
		}

		e.layoutGrid(children, cols, rows, w, h, gap)

		return children, nil

	case TypeSpacer:
		spacerW := getAttrInt(elem.Attributes, "width", w)
		spacerH := getAttrInt(elem.Attributes, "height", 1)
		flex := getAttrInt(elem.Attributes, "flex", 1)

		return []*ComponentNode{{
			ID:     elem.ID,
			Type:   TypeSpacer,
			X:      x,
			Y:      y,
			Width:  spacerW,
			Height: spacerH,
			Attrs:  map[string]string{"flex": strconv.Itoa(flex)},
		}}, nil

	case TypeDivider:
		char := elem.Attributes["char"]
		if char == "" {
			char = "-"
		}
		return []*ComponentNode{{
			ID:      elem.ID,
			Type:    TypeDivider,
			X:       x,
			Y:       y,
			Width:   w,
			Height:  1,
			Content: char,
			Attrs:   elem.Attributes,
		}}, nil

	case TypeText:
		value := elem.Content
		if value == "" {
			value = elem.Attributes["value"]
		}
		color := elem.Attributes["color"]
		bold := getAttrBool(elem.Attributes, "bold", false)

		return []*ComponentNode{{
			ID:      elem.ID,
			Type:    TypeText,
			X:       x,
			Y:       y,
			Width:   w,
			Height:  1,
			Content: value,
			Attrs:   map[string]string{"color": color, "bold": strconv.FormatBool(bold)},
		}}, nil

	case TypeList:
		items := elem.Content
		if items == "" {
			items = elem.Attributes["items"]
		}
		itemList := strings.Split(items, ",")

		return []*ComponentNode{{
			ID:      elem.ID,
			Type:    TypeList,
			X:       x,
			Y:       y,
			Width:   w,
			Height:  h,
			Content: itemList,
			Attrs:   elem.Attributes,
		}}, nil

	case TypeButton:
		label := elem.Attributes["label"]
		if label == "" {
			label = elem.Content
		}
		primary := getAttrBool(elem.Attributes, "primary", false)

		node := &ComponentNode{
			ID:      elem.ID,
			Type:    TypeButton,
			X:       x,
			Y:       y,
			Width:   len(label) + 4,
			Height:  1,
			Content: label,
			Attrs:   map[string]string{"primary": strconv.FormatBool(primary)},
		}
		if elem.ID != "" {
			e.components[elem.ID] = node
		}
		return []*ComponentNode{node}, nil

	case TypeInput:
		placeholder := elem.Attributes["placeholder"]
		hidden := getAttrBool(elem.Attributes, "hidden", false)

		node := &ComponentNode{
			ID:      elem.ID,
			Type:    TypeInput,
			X:       x,
			Y:       y,
			Width:   w,
			Height:  1,
			Content: placeholder,
			Attrs:   map[string]string{"hidden": strconv.FormatBool(hidden)},
		}
		if elem.ID != "" {
			e.components[elem.ID] = node
		}
		return []*ComponentNode{node}, nil

	case TypeProgress:
		value := getAttrInt(elem.Attributes, "value", 0)
		width := getAttrInt(elem.Attributes, "width", w)

		return []*ComponentNode{{
			ID:      elem.ID,
			Type:    TypeProgress,
			X:       x,
			Y:       y,
			Width:   width,
			Height:  1,
			Content: value,
			Attrs:   elem.Attributes,
		}}, nil

	case TypeTabbar:
		tabs := getAttrString(elem.Attributes, "tabs", "")
		active := getAttrInt(elem.Attributes, "active", 0)
		return []*ComponentNode{{
			ID:      elem.ID,
			Type:    TypeTabbar,
			X:       x,
			Y:       y,
			Width:   w,
			Height:  1,
			Content: tabs,
			Attrs:   map[string]string{"tabs": tabs, "active": strconv.Itoa(active)},
		}}, nil

	case TypeStatusbar:
		model := getAttrString(elem.Attributes, "model", "")
		agent := getAttrString(elem.Attributes, "agent", "")
		provider := getAttrString(elem.Attributes, "provider", "")
		yolo := getAttrString(elem.Attributes, "yolo", "off")
		return []*ComponentNode{{
			ID:      elem.ID,
			Type:    TypeStatusbar,
			X:       x,
			Y:       y,
			Width:   w,
			Height:  1,
			Content: map[string]string{"model": model, "agent": agent, "provider": provider, "yolo": yolo},
			Attrs:   elem.Attributes,
		}}, nil

	case TypeMessageList:
		messages := getAttrString(elem.Attributes, "messages", "")
		return []*ComponentNode{{
			ID:      elem.ID,
			Type:    TypeMessageList,
			X:       x,
			Y:       y,
			Width:   w,
			Height:  h,
			Content: messages,
			Attrs:   elem.Attributes,
		}}, nil

	case TypeSelectionList:
		items := getAttrString(elem.Attributes, "items", "")
		selected := getAttrInt(elem.Attributes, "selected", 0)
		return []*ComponentNode{{
			ID:      elem.ID,
			Type:    TypeSelectionList,
			X:       x,
			Y:       y,
			Width:   w,
			Height:  h,
			Content: strings.Split(items, ","),
			Attrs:   map[string]string{"items": items, "selected": strconv.Itoa(selected)},
		}}, nil

	case TypeToast:
		message := getAttrString(elem.Attributes, "message", "")
		toastType := getAttrString(elem.Attributes, "type", "info")
		return []*ComponentNode{{
			ID:      elem.ID,
			Type:    TypeToast,
			X:       x,
			Y:       y,
			Width:   w,
			Height:  1,
			Content: message,
			Attrs:   map[string]string{"type": toastType},
		}}, nil

	case TypeDialog:
		title := getAttrString(elem.Attributes, "title", "")
		content := getAttrString(elem.Attributes, "content", "")
		dialogType := getAttrString(elem.Attributes, "dialog-type", "alert")
		return []*ComponentNode{{
			ID:      elem.ID,
			Type:    TypeDialog,
			X:       x,
			Y:       y,
			Width:   w,
			Height:  h,
			Content: content,
			Attrs:   map[string]string{"title": title, "dialog-type": dialogType},
		}}, nil

	default:
		return e.BuildTreeFromElements(elem.Children, x, y, w, h)
	}
}

func (e *Engine[R]) BuildTreeFromElements(elems []*Element, x, y, w, h int) ([]*ComponentNode, error) {
	var nodes []*ComponentNode
	cury := y

	for _, elem := range elems {
		elemNodes, err := e.BuildTree(elem, x, cury, w, h-cury+y)
		if err != nil {
			return nil, err
		}
		for _, n := range elemNodes {
			if n.Type != TypeSpacer || getAttrInt(n.Attrs, "flex", 1) == 0 {
				n.X = x
				n.Y = cury
			}
		}
		nodes = append(nodes, elemNodes...)
		if len(elemNodes) > 0 {
			last := elemNodes[len(elemNodes)-1]
			cury = last.Y + last.Height
		}
	}

	return nodes, nil
}

func (e *Engine[R]) layoutVBox(nodes []*ComponentNode, width, gap int, align string) {
	cury := 0
	flexUnit := e.calculateFlexUnit(nodes, width, gap, true)

	for _, n := range nodes {
		n.Y = cury
		if n.Type == TypeSpacer {
			flex := getAttrInt(n.Attrs, "flex", 1)
			n.Height = flex * flexUnit
			if n.Height < 1 {
				n.Height = 1
			}
		}
		cury += n.Height + gap

		switch align {
		case "center":
			n.X = (width - n.Width) / 2
		case "right":
			n.X = width - n.Width
		default:
			n.X = 0
		}
	}
}

func (e *Engine[R]) layoutHBox(nodes []*ComponentNode, height, gap int, align string) {
	curx := 0
	flexUnit := e.calculateFlexUnit(nodes, height, gap, false)

	for _, n := range nodes {
		n.X = curx
		if n.Type == TypeSpacer {
			flex := getAttrInt(n.Attrs, "flex", 1)
			n.Width = flex * flexUnit
			if n.Width < 1 {
				n.Width = 1
			}
		}
		curx += n.Width + gap

		switch align {
		case "center":
			n.Y = (height - n.Height) / 2
		case "bottom":
			n.Y = height - n.Height
		default:
			n.Y = 0
		}
	}
}

func (e *Engine[R]) calculateFlexUnit(nodes []*ComponentNode, total, gap int, isVertical bool) int {
	flexTotal := 0
	fixedTotal := 0

	for _, n := range nodes {
		if n.Type == TypeSpacer {
			flexTotal += getAttrInt(n.Attrs, "flex", 1)
		} else if isVertical {
			fixedTotal += n.Height
		} else {
			fixedTotal += n.Width
		}
	}

	if flexTotal == 0 {
		return 0
	}
	flexUnit := (total - fixedTotal - gap*(len(nodes)-1)) / flexTotal
	if flexUnit < 0 {
		return 0
	}
	return flexUnit
}

func (e *Engine[R]) layoutGrid(nodes []*ComponentNode, cols, rows, width, height, gap int) {
	if cols <= 0 {
		cols = 2
	}
	if rows <= 0 {
		rows = (len(nodes) + cols - 1) / cols
	}

	cellW := width / cols
	cellH := height / rows

	for i, n := range nodes {
		col := i % cols
		row := i / cols
		n.X = col * (cellW + gap)
		n.Y = row * (cellH + gap)
		n.Width = cellW
		n.Height = cellH
	}
}

func (e *Engine[R]) renderNodes(nodes []*ComponentNode, r R) string {
	result := ""
	for _, n := range nodes {
		result += e.renderNode(n, r)
	}
	return result
}

func (e *Engine[R]) renderNode(n *ComponentNode, r R) string {
	switch n.Type {
	case TypeWindow:
		border := r.RenderBorder(n.X, n.Y, n.Width, n.Height, getAttrColor(n.Attrs, "border", ""))
		content := e.renderNodes(n.Children, r)
		title := n.Attrs["title"]
		if title != "" {
			titleStr := r.RenderText(title, n.X+1, n.Y, getAttrColor(n.Attrs, "title-color", ""))
			return border + titleStr + content
		}
		return border + content

	case TypeText:
		text := fmt.Sprintf("%v", n.Content)
		return e.renderMultiLineText(text, n.X, n.Y, getAttrColor(n.Attrs, "color", ""), r)

	case TypeList:
		items, ok := n.Content.([]string)
		if !ok {
			return ""
		}
		content := strings.Join(items, "\n")
		return e.renderMultiLineText(content, n.X, n.Y, getAttrColor(n.Attrs, "color", ""), r)

	case TypeButton:
		label := fmt.Sprintf("%v", n.Content)
		primary := n.Attrs["primary"] == "true"
		fg := getAttrColor(n.Attrs, "color", "")
		bg := getAttrColor(n.Attrs, "bg", "")
		if primary {
			bg = getAttrColor(n.Attrs, "primary-bg", "#0099ff")
			fg = getAttrColor(n.Attrs, "primary-fg", "#ffffff")
		}
		return r.RenderSelected(" "+label+" ", n.X, n.Y, n.Width, fg, bg)

	case TypeInput:
		placeholder := fmt.Sprintf("%v", n.Content)
		hidden := n.Attrs["hidden"] == "true"
		text := placeholder
		if hidden {
			text = strings.Repeat("*", len(placeholder))
		}
		return r.RenderSelected(" "+text+"_", n.X, n.Y, n.Width, "", "")

	case TypeSpacer:
		return ""

	case TypeDivider:
		char := "-"
		if c, ok := n.Attrs["char"]; ok {
			char = c
		}
		line := strings.Repeat(char, n.Width)
		return r.RenderText(line, n.X, n.Y, getAttrColor(n.Attrs, "color", ""))

	case TypeProgress:
		val := 0
		if v, ok := n.Content.(int); ok {
			val = v
		}
		barLen := n.Width - 2
		filled := (barLen * val) / 100
		if filled > barLen {
			filled = barLen
		}
		bar := "[" + strings.Repeat("=", filled) + strings.Repeat(" ", barLen-filled) + "]"
		return r.RenderText(bar, n.X, n.Y, getAttrColor(n.Attrs, "color", ""))

	case TypeTabbar:
		tabs := n.Attrs["tabs"]
		active := 0
		if a, err := strconv.Atoi(n.Attrs["active"]); err == nil {
			active = a
		}
		return e.renderTabbar(tabs, active, r)

	case TypeStatusbar:
		content, ok := n.Content.(map[string]string)
		if !ok {
			content = make(map[string]string)
		}
		return e.renderStatusbar(content, r)

	case TypeMessageList:
		return e.renderMessageList(n, r)

	case TypeSelectionList:
		items, ok := n.Content.([]string)
		if !ok {
			items = []string{}
		}
		selected := 0
		if s, err := strconv.Atoi(n.Attrs["selected"]); err == nil {
			selected = s
		}
		return e.renderSelectionList(items, selected, n, r)

	case TypeToast:
		message := fmt.Sprintf("%v", n.Content)
		toastType := n.Attrs["type"]
		return e.renderToast(message, toastType, r)

	case TypeDialog:
		title := n.Attrs["title"]
		content := fmt.Sprintf("%v", n.Content)
		dialogType := n.Attrs["dialog-type"]
		return e.renderDialog(title, content, dialogType, n, r)

	default:
		return e.renderNodes(n.Children, r)
	}
}

func (e *Engine[R]) renderMultiLineText(text string, x int, y int, color string, r R) string {
	lines := strings.Split(text, "\n")
	var result string
	for i, line := range lines {
		result += r.RenderText(line, x, y+i, color)
	}
	return result
}

func (e *Engine[R]) renderTabbar(tabs string, active int, r R) string {
	tabList := strings.Split(tabs, ",")
	var result string
	for i, tab := range tabList {
		tab = strings.TrimSpace(tab)
		if i == active {
			result += r.RenderSelected("["+tab+"]", 0, 0, 0, "", "")
		} else {
			result += r.RenderText("["+tab+"]", 0, 0, "")
		}
	}
	return result
}

func (e *Engine[R]) renderStatusbar(content map[string]string, r R) string {
	model := content["model"]
	agent := content["agent"]
	provider := content["provider"]
	yolo := content["yolo"]

	var result string
	if provider != "" {
		result += "● " + provider + "  "
	}
	if yolo == "on" {
		result += "YOLO: ON  "
	}
	if model != "" {
		result += "Model: " + model + "  "
	}
	if agent != "" {
		result += "Agent: " + agent
	}
	return result
}

func (e *Engine[R]) renderMessageList(n *ComponentNode, r R) string {
	messages := getAttrString(n.Attrs, "messages", "")
	return r.RenderText(messages, n.X, n.Y, "")
}

func (e *Engine[R]) renderSelectionList(items []string, selected int, n *ComponentNode, r R) string {
	var result string
	for i, item := range items {
		item = strings.TrimSpace(item)
		if i == selected {
			result += r.RenderSelected(item, n.X, n.Y+i, n.Width, "", "")
		} else {
			result += r.RenderText(item, n.X, n.Y+i, "")
		}
	}
	return result
}

func (e *Engine[R]) renderToast(message string, toastType string, r R) string {
	var color string
	switch toastType {
	case "success":
		color = "#4EC9B0"
	case "warning":
		color = "#FFCC00"
	case "error":
		color = "#F44747"
	default:
		color = "#007ACC"
	}
	return r.RenderText(message, 0, 0, color)
}

func (e *Engine[R]) renderDialog(title, content, dialogType string, n *ComponentNode, r R) string {
	border := r.RenderBorder(n.X, n.Y, n.Width, n.Height, "#3D3D3D")
	titleStr := r.RenderText(title, n.X+1, n.Y, "#007ACC")
	contentStr := r.RenderText(content, n.X+1, n.Y+1, "#E0E0E0")
	return border + titleStr + contentStr
}

func (e *Engine[R]) GetComponent(id string) *ComponentNode {
	return e.components[id]
}

func (e *Engine[R]) RegisterComponent(id string, node *ComponentNode) {
	e.components[id] = node
}

func getAttr(attrs map[string]string, key, def string) string {
	if v, ok := attrs[key]; ok {
		return v
	}
	return def
}

func getAttrInt(attrs map[string]string, key string, def int) int {
	if v, ok := attrs[key]; ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func getAttrBool(attrs map[string]string, key string, def bool) bool {
	if v, ok := attrs[key]; ok {
		return v == "true" || v == "1" || v == "yes"
	}
	return def
}

func getAttrString(attrs map[string]string, key, def string) string {
	if v, ok := attrs[key]; ok {
		return v
	}
	return def
}

func getAttrColor(attrs map[string]string, key, def string) string {
	if v, ok := attrs[key]; ok {
		return v
	}
	return def
}
