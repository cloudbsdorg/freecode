package component

import (
	"testing"

	"github.com/freecode/freecode/internal/renderer"
)

type mockRenderer struct {
	lastRender string
	calls      []string
}

func (m mockRenderer) RenderBox(x, y, w, h int, bgColor string) string {
	m.calls = append(m.calls, "RenderBox")
	return m.lastRender
}

func (m mockRenderer) RenderText(text string, x, y int, fgColor string) string {
	m.calls = append(m.calls, "RenderText")
	return m.lastRender
}

func (m mockRenderer) RenderBorder(x, y, w, h int, fgColor string) string {
	m.calls = append(m.calls, "RenderBorder")
	return m.lastRender
}

func (m mockRenderer) RenderSelected(text string, x, y, w int, fg, bg string) string {
	m.calls = append(m.calls, "RenderSelected")
	return m.lastRender
}

func (m mockRenderer) Width() int {
	return 80
}

func (m mockRenderer) Height() int {
	return 24
}

var _ renderer.Renderer = mockRenderer{}

func TestNewButton(t *testing.T) {
	colors := ButtonColors{
		Background:    "#000000",
		Foreground:    "#FFFFFF",
		SelectedBg:    "#0000FF",
		SelectedFg:    "#FFFFFF",
		PressedBg:     "#00FF00",
		PressedFg:     "#000000",
	}
	btn := NewButton[mockRenderer]("Click", 10, 5, colors)

	if btn.Text != "Click" {
		t.Errorf("expected Text=Click, got %s", btn.Text)
	}
	if btn.X != 10 {
		t.Errorf("expected X=10, got %d", btn.X)
	}
	if btn.Y != 5 {
		t.Errorf("expected Y=5, got %d", btn.Y)
	}
	if btn.Width != 7 {
		t.Errorf("expected Width=7 (len(Click)+2), got %d", btn.Width)
	}
	if !btn.Visible {
		t.Error("expected Visible=true")
	}
}

func TestButtonSetText(t *testing.T) {
	colors := ButtonColors{}
	btn := NewButton[mockRenderer]("Hi", 0, 0, colors)
	btn.SetText("Longer Text!")
	if btn.Text != "Longer Text!" {
		t.Errorf("expected Text=Longer Text!, got %s", btn.Text)
	}
	if btn.Width != 14 {
		t.Errorf("expected Width=14, got %d", btn.Width)
	}
}

func TestButtonSelect(t *testing.T) {
	colors := ButtonColors{}
	btn := NewButton[mockRenderer]("Test", 0, 0, colors)
	if btn.Selected {
		t.Error("expected Selected=false initially")
	}
	btn.Select()
	if !btn.Selected {
		t.Error("expected Selected=true after Select()")
	}
}

func TestButtonPress(t *testing.T) {
	colors := ButtonColors{}
	called := false
	btn := NewButton[mockRenderer]("Test", 0, 0, colors)
	btn.OnClick = func() {
		called = true
	}

	btn.Press()
	if !btn.Pressed {
		t.Error("expected Pressed=true after Press()")
	}
	if !called {
		t.Error("expected OnClick to be called")
	}
}

func TestButtonRelease(t *testing.T) {
	colors := ButtonColors{}
	btn := NewButton[mockRenderer]("Test", 0, 0, colors)
	btn.Pressed = true
	btn.Release()
	if btn.Pressed {
		t.Error("expected Pressed=false after Release()")
	}
}

func TestButtonRender(t *testing.T) {
	colors := ButtonColors{
		Background: "#000000",
		Foreground: "#FFFFFF",
	}
	btn := NewButton[mockRenderer]("Test", 5, 10, colors)
	m := mockRenderer{lastRender: "rendered"}
	btn.Render(m)
}

func TestButtonRenderNotVisible(t *testing.T) {
	colors := ButtonColors{}
	btn := NewButton[mockRenderer]("Test", 0, 0, colors)
	btn.Visible = false
	m := mockRenderer{}
	result := btn.Render(m)
	if result != "" {
		t.Error("expected empty string when not visible")
	}
}

func TestButtonRenderSelected(t *testing.T) {
	colors := ButtonColors{
		SelectedBg: "#0000FF",
		SelectedFg: "#FFFFFF",
	}
	btn := NewButton[mockRenderer]("Test", 0, 0, colors)
	btn.Selected = true
	m := mockRenderer{lastRender: "selected"}
	btn.Render(m)
}

func TestButtonRenderPressed(t *testing.T) {
	colors := ButtonColors{
		PressedBg: "#00FF00",
		PressedFg: "#000000",
	}
	btn := NewButton[mockRenderer]("Test", 0, 0, colors)
	btn.Pressed = true
	m := mockRenderer{lastRender: "pressed"}
	btn.Render(m)
}

func TestComponentSetPosition(t *testing.T) {
	c := &Component[mockRenderer]{}
	c.SetPosition(5, 10)
	if c.X != 5 {
		t.Errorf("expected X=5, got %d", c.X)
	}
	if c.Y != 10 {
		t.Errorf("expected Y=10, got %d", c.Y)
	}
}

func TestComponentSetSize(t *testing.T) {
	c := &Component[mockRenderer]{}
	c.SetSize(100, 50)
	if c.Width != 100 {
		t.Errorf("expected Width=100, got %d", c.Width)
	}
	if c.Height != 50 {
		t.Errorf("expected Height=50, got %d", c.Height)
	}
}

func TestComponentShowHide(t *testing.T) {
	c := &Component[mockRenderer]{Visible: false}
	c.Show()
	if !c.Visible {
		t.Error("expected Visible=true after Show()")
	}
	c.Hide()
	if c.Visible {
		t.Error("expected Visible=false after Hide()")
	}
}

func TestComponentIsVisible(t *testing.T) {
	c := &Component[mockRenderer]{Visible: true}
	if !c.IsVisible() {
		t.Error("expected IsVisible()=true")
	}
}
