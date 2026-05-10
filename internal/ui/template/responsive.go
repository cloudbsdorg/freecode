package template

import (
	"strconv"
	"strings"

	"github.com/freecode/freecode/internal/renderer"
)

type Constraint struct {
	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int
}

type Responsive interface {
	GetConstraint() Constraint
}

func (e *Engine[R]) shouldRender(elem *Element, width, height int) bool {
	attrs := elem.Attributes

	if minW := attrs["min-width"]; minW != "" {
		if w := e.parsePercent(minW, width); w > width {
			return false
		}
	}
	if maxW := attrs["max-width"]; maxW != "" {
		if w := e.parsePercent(maxW, width); w < width {
			return false
		}
	}
	if minH := attrs["min-height"]; minH != "" {
		if h := e.parsePercent(minH, height); h > height {
			return false
		}
	}
	if maxH := attrs["max-height"]; maxH != "" {
		if h := e.parsePercent(maxH, height); h < height {
			return false
		}
	}

	if showIf := attrs["show-if"]; showIf != "" {
		if !e.evaluateCondition(showIf, width, height) {
			return false
		}
	}
	if hideIf := attrs["hide-if"]; hideIf != "" {
		if e.evaluateCondition(hideIf, width, height) {
			return false
		}
	}

	return true
}

func (e *Engine[R]) parsePercent(val string, base int) int {
	val = strings.TrimSpace(val)
	if strings.HasSuffix(val, "%") {
		pct, err := strconv.ParseFloat(val[:len(val)-1], 64)
		if err != nil {
			return base
		}
		return int(float64(base) * pct / 100)
	}
	w, err := strconv.Atoi(val)
	if err != nil {
		return base
	}
	return w
}

func (e *Engine[R]) evaluateCondition(cond string, width, height int) bool {
	cond = strings.TrimSpace(cond)

	parts := strings.Split(cond, ">=")
	if len(parts) == 2 {
		val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		if strings.TrimSpace(parts[0]) == "width" {
			return width >= val
		}
		if strings.TrimSpace(parts[0]) == "height" {
			return height >= val
		}
	}

	parts = strings.Split(cond, "<=")
	if len(parts) == 2 {
		val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		if strings.TrimSpace(parts[0]) == "width" {
			return width <= val
		}
		if strings.TrimSpace(parts[0]) == "height" {
			return height <= val
		}
	}

	parts = strings.Split(cond, ">")
	if len(parts) == 2 {
		val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		if strings.TrimSpace(parts[0]) == "width" {
			return width > val
		}
		if strings.TrimSpace(parts[0]) == "height" {
			return height > val
		}
	}

	parts = strings.Split(cond, "<")
	if len(parts) == 2 {
		val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		if strings.TrimSpace(parts[0]) == "width" {
			return width < val
		}
		if strings.TrimSpace(parts[0]) == "height" {
			return height < val
		}
	}

	parts = strings.Split(cond, "==")
	if len(parts) == 2 {
		val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		if strings.TrimSpace(parts[0]) == "width" {
			return width == val
		}
		if strings.TrimSpace(parts[0]) == "height" {
			return height == val
		}
	}

	return false
}

func (e *Engine[R]) parseWidthAttr(attrs map[string]string, key string, def int, parentWidth int) int {
	if val := attrs[key]; val != "" {
		if strings.HasSuffix(val, "%") {
			pct, err := strconv.ParseFloat(val[:len(val)-1], 64)
			if err == nil {
				return int(float64(parentWidth) * pct / 100)
			}
		}
		if w, err := strconv.Atoi(val); err == nil {
			return w
		}
	}
	return def
}

type SizeObserver interface {
	OnSizeChanged(width, height int)
}

type VisibilityObserver interface {
	OnVisibilityChanged(id string, visible bool)
}

type ResponsiveEngine[R renderer.Renderer] struct {
	*Engine[R]
	sizeObservers      []SizeObserver
	visibilityObs      []VisibilityObserver
	hiddenComponents   map[string]bool
	hiddenByCondition  map[string]string
	OnDialogShown      func(dialog any)
	OnDialogHidden     func(dialog any)
}

func NewResponsiveEngine[R renderer.Renderer]() *ResponsiveEngine[R] {
	return &ResponsiveEngine[R]{
		Engine:           NewEngine[R](),
		sizeObservers:     make([]SizeObserver, 0),
		visibilityObs:     make([]VisibilityObserver, 0),
		hiddenComponents:  make(map[string]bool),
		hiddenByCondition: make(map[string]string),
	}
}

func (re *ResponsiveEngine[R]) AddSizeObserver(o SizeObserver) {
	re.sizeObservers = append(re.sizeObservers, o)
}

func (re *ResponsiveEngine[R]) RemoveSizeObserver(o SizeObserver) {
	for i, obs := range re.sizeObservers {
		if obs == o {
			re.sizeObservers = append(re.sizeObservers[:i], re.sizeObservers[i+1:]...)
			break
		}
	}
}

func (re *ResponsiveEngine[R]) AddVisibilityObserver(o VisibilityObserver) {
	re.visibilityObs = append(re.visibilityObs, o)
}

func (re *ResponsiveEngine[R]) RemoveVisibilityObserver(o VisibilityObserver) {
	for i, obs := range re.visibilityObs {
		if obs == o {
			re.visibilityObs = append(re.visibilityObs[:i], re.visibilityObs[i+1:]...)
			break
		}
	}
}

func (re *ResponsiveEngine[R]) NotifySizeChanged(width, height int) {
	for _, o := range re.sizeObservers {
		o.OnSizeChanged(width, height)
	}
}

func (re *ResponsiveEngine[R]) NotifyVisibilityChanged(id string, visible bool) {
	for _, o := range re.visibilityObs {
		o.OnVisibilityChanged(id, visible)
	}
}

func (re *ResponsiveEngine[R]) Show(id string) {
	if re.hiddenComponents[id] {
		re.hiddenComponents[id] = false
		re.NotifyVisibilityChanged(id, true)
	}
}

func (re *ResponsiveEngine[R]) Hide(id string) {
	if !re.hiddenComponents[id] {
		re.hiddenComponents[id] = true
		re.NotifyVisibilityChanged(id, false)
	}
}

func (re *ResponsiveEngine[R]) Toggle(id string) {
	if re.hiddenComponents[id] {
		re.Show(id)
	} else {
		re.Hide(id)
	}
}

func (re *ResponsiveEngine[R]) IsVisible(id string) bool {
	if re.hiddenComponents[id] {
		return false
	}
	node := re.GetComponent(id)
	if node == nil {
		return true
	}
	return node.Visible
}

func (re *ResponsiveEngine[R]) SetVisible(id string, visible bool) {
	if visible {
		re.Show(id)
	} else {
		re.Hide(id)
	}
}

func (re *ResponsiveEngine[R]) ShowAll() {
	for id := range re.hiddenComponents {
		re.Show(id)
	}
}

func (re *ResponsiveEngine[R]) HideAll() {
	for id := range re.hiddenComponents {
		re.Hide(id)
	}
}

func (re *ResponsiveEngine[R]) RenderAt(src string, width, height int, r R) (string, error) {
	re.NotifySizeChanged(width, height)
	return re.Engine.ParseAndRender(src, width, height, r)
}

func (re *ResponsiveEngine[R]) IsHiddenByCondition(id string) bool {
	_, exists := re.hiddenByCondition[id]
	return exists
}

func (re *ResponsiveEngine[R]) SetComponentAttr(id, key, value string) {
	node := re.GetComponent(id)
	if node != nil {
		if node.Attrs == nil {
			node.Attrs = make(map[string]string)
		}
		node.Attrs[key] = value
	}
}

func (re *ResponsiveEngine[R]) SetComponentContent(id string, content interface{}) {
	node := re.GetComponent(id)
	if node != nil {
		node.Content = content
	}
}

func (re *ResponsiveEngine[R]) GetComponentAttr(id, key string) string {
	node := re.GetComponent(id)
	if node != nil {
		if val, ok := node.Attrs[key]; ok {
			return val
		}
	}
	return ""
}

func (re *ResponsiveEngine[R]) GetComponentContent(id string) interface{} {
	node := re.GetComponent(id)
	if node != nil {
		return node.Content
	}
	return nil
}
