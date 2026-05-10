package component

import "github.com/freecode/freecode/internal/renderer"

type Tab struct {
	ID    string
	Title string
	Closeable bool
}

type TabBar[R renderer.Renderer] struct {
	Component[R]
	Tabs      []Tab
	ActiveIdx int
	Colors    TabBarColors
}

type TabBarColors struct {
	Background    string
	Foreground    string
	ActiveBg      string
	ActiveFg      string
	InactiveBg    string
	InactiveFg    string
	CloseColor    string
}

func NewTabBar[R renderer.Renderer](width int, colors TabBarColors) *TabBar[R] {
	return &TabBar[R]{
		Component: Component[R]{
			X:       0,
			Y:       0,
			Width:   width,
			Height:  1,
			Visible: true,
		},
		Tabs:      []Tab{},
		ActiveIdx: 0,
		Colors:    colors,
	}
}

func (t *TabBar[R]) AddTab(id, title string) {
	t.Tabs = append(t.Tabs, Tab{ID: id, Title: title, Closeable: true})
}

func (t *TabBar[R]) CloseTab(idx int) {
	if idx < 0 || idx >= len(t.Tabs) {
		return
	}
	t.Tabs = append(t.Tabs[:idx], t.Tabs[idx+1:]...)
	if t.ActiveIdx >= len(t.Tabs) {
		t.ActiveIdx = len(t.Tabs) - 1
		if t.ActiveIdx < 0 {
			t.ActiveIdx = 0
		}
	}
}

func (t *TabBar[R]) SetActive(idx int) {
	if idx >= 0 && idx < len(t.Tabs) {
		t.ActiveIdx = idx
	}
}

func (t *TabBar[R]) GetActive() *Tab {
	if t.ActiveIdx < 0 || t.ActiveIdx >= len(t.Tabs) {
		return nil
	}
	return &t.Tabs[t.ActiveIdx]
}

func (t *TabBar[R]) Render(r R) string {
	if !t.Visible || len(t.Tabs) == 0 {
		return r.RenderBox(t.X, t.Y, t.Width, 1, t.Colors.Background)
	}

	result := ""
	x := t.X

	for i, tab := range t.Tabs {
		tabWidth := len(tab.Title) + 4
		if x+tabWidth > t.X+t.Width {
			break
		}

		if i == t.ActiveIdx {
			result += r.RenderSelected(" "+tab.Title+" ", x, t.Y, tabWidth, t.Colors.ActiveFg, t.Colors.ActiveBg)
		} else {
			result += r.RenderSelected(" "+tab.Title+" ", x, t.Y, tabWidth, t.Colors.InactiveFg, t.Colors.InactiveBg)
		}
		x += tabWidth
	}

	if x < t.X+t.Width {
		result += r.RenderText("", x, t.Y, t.Colors.Background)
	}

	return r.RenderBox(t.X, t.Y, t.Width, 1, t.Colors.Background) + result
}
