package template

import (
	"testing"

	"github.com/freecode/freecode/internal/renderer"
)

var _ renderer.Renderer = mockRenderer{}

func TestParsePercent(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	tests := []struct {
		val    string
		base   int
		want   int
	}{
		{"50%", 100, 50},
		{"75%", 200, 150},
		{"100%", 80, 80},
		{"25%", 80, 20},
		{"10%", 50, 5},
		{"0%", 100, 0},
		{"120%", 100, 120},
	}

	for _, tt := range tests {
		got := engine.parsePercent(tt.val, tt.base)
		if got != tt.want {
			t.Errorf("parsePercent(%q, %d) = %d, want %d", tt.val, tt.base, got, tt.want)
		}
	}
}

func TestEvaluateCondition(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	tests := []struct {
		cond   string
		width  int
		height int
		want   bool
	}{
		{"width>=80", 80, 24, true},
		{"width>=80", 79, 24, false},
		{"width>80", 81, 24, true},
		{"width>80", 80, 24, false},
		{"width<=80", 80, 24, true},
		{"width<=80", 81, 24, false},
		{"width<80", 79, 24, true},
		{"width<80", 80, 24, false},
		{"width==80", 80, 24, true},
		{"width==80", 81, 24, false},
		{"height>=24", 80, 24, true},
		{"height>=24", 80, 23, false},
		{"height>24", 80, 25, true},
		{"height<24", 80, 23, true},
	}

	for _, tt := range tests {
		got := engine.evaluateCondition(tt.cond, tt.width, tt.height)
		if got != tt.want {
			t.Errorf("evaluateCondition(%q, %d, %d) = %v, want %v", tt.cond, tt.width, tt.height, got, tt.want)
		}
	}
}

func TestParseWidthAttr(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	tests := []struct {
		attrs map[string]string
		key   string
		def   int
		base  int
		want  int
	}{
		{map[string]string{"width": "50%"}, "width", 100, 100, 50},
		{map[string]string{"width": "100"}, "width", 50, 100, 100},
		{map[string]string{"foo": "bar"}, "width", 80, 100, 80},
		{map[string]string{"width": ""}, "width", 80, 100, 80},
		{map[string]string{"width": "25%"}, "width", 0, 200, 50},
	}

	for _, tt := range tests {
		got := engine.parseWidthAttr(tt.attrs, tt.key, tt.def, tt.base)
		if got != tt.want {
			t.Errorf("parseWidthAttr(%v, %q, %d, %d) = %d, want %d", tt.attrs, tt.key, tt.def, tt.base, got, tt.want)
		}
	}
}

type mockObserver struct {
	callCount int
	w         int
	h         int
}

func (m *mockObserver) OnSizeChanged(width, height int) {
	m.callCount++
	m.w = width
	m.h = height
}

func TestResponsiveEngineObserver(t *testing.T) {
	re := NewResponsiveEngine[mockRenderer]()

	obs1 := &mockObserver{callCount: 0}
	obs2 := &mockObserver{callCount: 0}

	re.AddSizeObserver(obs1)
	re.AddSizeObserver(obs2)

	re.NotifySizeChanged(100, 50)

	if obs1.callCount != 1 {
		t.Errorf("obs1 callCount = %d, want 1", obs1.callCount)
	}
	if obs1.w != 100 || obs1.h != 50 {
		t.Errorf("obs1 size = (%d,%d), want (100,50)", obs1.w, obs1.h)
	}

	if obs2.callCount != 1 {
		t.Errorf("obs2 callCount = %d, want 1", obs2.callCount)
	}

	re.RemoveSizeObserver(obs1)
	re.NotifySizeChanged(200, 100)

	if obs1.callCount != 1 {
		t.Error("obs1 should have been removed and not called again")
	}
	if obs2.callCount != 2 {
		t.Errorf("obs2 callCount = %d, want 2", obs2.callCount)
	}
}

func TestResponsiveEngineRenderAt(t *testing.T) {
	re := NewResponsiveEngine[mockRenderer]()

	obs := &mockObserver{}
	re.AddSizeObserver(obs)

	src := `<text value="Hello" />`

	_, err := re.RenderAt(src, 80, 24, mockRenderer{})
	if err != nil {
		t.Fatalf("RenderAt error = %v", err)
	}

	if obs.callCount == 0 {
		t.Error("observer should have been notified")
	}
	if obs.w != 80 || obs.h != 24 {
		t.Errorf("observer size = (%d,%d), want (80,24)", obs.w, obs.h)
	}
}

type visibilityObserver struct {
	id       string
	visible  bool
	called   bool
}

func (v *visibilityObserver) OnVisibilityChanged(id string, visible bool) {
	v.id = id
	v.visible = visible
	v.called = true
}

func TestShowHide(t *testing.T) {
	re := NewResponsiveEngine[mockRenderer]()

	re.Hide("btn1")
	if re.IsVisible("btn1") {
		t.Error("btn1 should be hidden after Hide()")
	}

	re.Show("btn1")
	if !re.IsVisible("btn1") {
		t.Error("btn1 should be visible after Show()")
	}
}

func TestToggle(t *testing.T) {
	re := NewResponsiveEngine[mockRenderer]()

	re.Hide("myid")
	re.Toggle("myid")
	if !re.IsVisible("myid") {
		t.Error("myid should be visible after first toggle (was hidden)")
	}

	re.Toggle("myid")
	if re.IsVisible("myid") {
		t.Error("myid should be hidden after second toggle")
	}
}

func TestShowAllHideAll(t *testing.T) {
	re := NewResponsiveEngine[mockRenderer]()

	re.Hide("a")
	re.Hide("b")
	re.Hide("c")

	re.HideAll()

	if re.IsVisible("a") || re.IsVisible("b") || re.IsVisible("c") {
		t.Error("all should be hidden after HideAll()")
	}

	re.ShowAll()

	if !re.IsVisible("a") || !re.IsVisible("b") || !re.IsVisible("c") {
		t.Error("all should be visible after ShowAll()")
	}
}

func TestVisibilityObserver(t *testing.T) {
	re := NewResponsiveEngine[mockRenderer]()

	obs := &visibilityObserver{}
	re.AddVisibilityObserver(obs)

	re.Hide("myid")

	if !obs.called {
		t.Error("visibility observer should be called")
	}
	if obs.id != "myid" {
		t.Errorf("id = %q, want %q", obs.id, "myid")
	}
	if obs.visible {
		t.Error("visible should be false after Hide()")
	}

	obs.called = false
	re.Show("myid")

	if !obs.called {
		t.Error("visibility observer should be called after Show()")
	}
	if !obs.visible {
		t.Error("visible should be true after Show()")
	}
}

func TestSetVisible(t *testing.T) {
	re := NewResponsiveEngine[mockRenderer]()

	re.SetVisible("testid", false)
	if re.IsVisible("testid") {
		t.Error("testid should be hidden after SetVisible(false)")
	}

	re.SetVisible("testid", true)
	if !re.IsVisible("testid") {
		t.Error("testid should be visible after SetVisible(true)")
	}
}

func TestShouldRenderMinWidth(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	elem := &Element{
		Type:       TypeText,
		Attributes: map[string]string{"min-width": "100"},
	}

	if !engine.shouldRender(elem, 100, 24) {
		t.Error("shouldRender(100) should be true")
	}
	if !engine.shouldRender(elem, 150, 24) {
		t.Error("shouldRender(150) should be true")
	}
	if engine.shouldRender(elem, 50, 24) {
		t.Error("shouldRender(50) should be false")
	}
}

func TestShouldRenderMaxWidth(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	elem := &Element{
		Type:       TypeText,
		Attributes: map[string]string{"max-width": "100"},
	}

	if !engine.shouldRender(elem, 100, 24) {
		t.Error("shouldRender(100) should be true")
	}
	if !engine.shouldRender(elem, 50, 24) {
		t.Error("shouldRender(50) should be true")
	}
	if engine.shouldRender(elem, 150, 24) {
		t.Error("shouldRender(150) should be false")
	}
}

func TestShouldRenderShowIf(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	elem := &Element{
		Type:       TypeText,
		Attributes: map[string]string{"show-if": "width>=80"},
	}

	if !engine.shouldRender(elem, 80, 24) {
		t.Error("shouldRender(80) should be true for width>=80")
	}
	if !engine.shouldRender(elem, 100, 24) {
		t.Error("shouldRender(100) should be true for width>=80")
	}
	if engine.shouldRender(elem, 79, 24) {
		t.Error("shouldRender(79) should be false for width>=80")
	}
}

func TestShouldRenderHideIf(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	elem := &Element{
		Type:       TypeText,
		Attributes: map[string]string{"hide-if": "width<80"},
	}

	if engine.shouldRender(elem, 79, 24) {
		t.Error("shouldRender(79) should be false for hide-if width<80")
	}
	if !engine.shouldRender(elem, 80, 24) {
		t.Error("shouldRender(80) should be true for hide-if width<80")
	}
}

func TestShouldRenderPercentConstraint(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	elem := &Element{
		Type:       TypeText,
		Attributes: map[string]string{"min-width": "50"},
	}

	if !engine.shouldRender(elem, 100, 24) {
		t.Error("shouldRender(100) should be true when min-width=50")
	}
	if engine.shouldRender(elem, 49, 24) {
		t.Error("shouldRender(49) should be false when min-width=50")
	}
}
