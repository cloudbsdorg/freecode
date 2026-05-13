package dialog

import (
	"testing"
)

func TestTextInputWithMaxLen(t *testing.T) {
	ti := NewTextInput(TextInputWithMaxLen(10))
	if ti.MaxLen != 10 {
		t.Errorf("expected MaxLen=10, got %d", ti.MaxLen)
	}
}

func TestTextInputWithHidden(t *testing.T) {
	ti := NewTextInput(TextInputWithHidden(true))
	if !ti.Hidden {
		t.Error("expected Hidden=true")
	}
}

func TestTextInputWithPlaceholder(t *testing.T) {
	ti := NewTextInput(TextInputWithPlaceholder("enter text"))
	if ti.Placeholder != "enter text" {
		t.Errorf("expected placeholder, got %s", ti.Placeholder)
	}
}

func TestTextInputWithColors(t *testing.T) {
	ti := NewTextInput(TextInputWithColors(Light))
	if ti.Colors != Light {
		t.Error("expected Light colors")
	}
}

func TestTextInputWithOnChange(t *testing.T) {
	called := false
	ti := NewTextInput(TextInputWithOnChange(func(s string) {
		called = true
	}))
	ti.SetValue("test")
	if !called {
		t.Error("expected OnChange to be called")
	}
}

func TestTextInputWithOnSubmit(t *testing.T) {
	called := false
	ti := NewTextInput(TextInputWithOnSubmit(func(s string) {
		called = true
	}))
	ti.Submit()
	if !called {
		t.Error("expected OnSubmit to be called")
	}
}

func TestTextInputWithOnCancel(t *testing.T) {
	called := false
	ti := NewTextInput(TextInputWithOnCancel(func() {
		called = true
	}))
	ti.Cancel()
	if !called {
		t.Error("expected OnCancel to be called")
	}
}

func TestTextInputWithWidth(t *testing.T) {
	ti := NewTextInput(TextInputWithWidth(100))
	if ti.Width != 100 {
		t.Errorf("expected Width=100, got %d", ti.Width)
	}
}

func TestNewTextInput(t *testing.T) {
	ti := NewTextInput()
	if ti.Value != "" {
		t.Errorf("expected empty Value, got %s", ti.Value)
	}
	if ti.MaxLen != 0 {
		t.Errorf("expected MaxLen=0, got %d", ti.MaxLen)
	}
	if ti.Width != 50 {
		t.Errorf("expected Width=50, got %d", ti.Width)
	}
	if ti.Hidden {
		t.Error("expected Hidden=false")
	}
	if ti.Cursor != 0 {
		t.Errorf("expected Cursor=0, got %d", ti.Cursor)
	}
	if !ti.Focused {
		t.Error("expected Focused=true")
	}
}

func TestTextInputSetValue(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("hello")
	if ti.Value != "hello" {
		t.Errorf("expected Value=hello, got %s", ti.Value)
	}
}

func TestTextInputSetValueWithMaxLen(t *testing.T) {
	ti := NewTextInput(TextInputWithMaxLen(5))
	ti.SetValue("longer string")
	if ti.Value != "longe" {
		t.Errorf("expected Value truncated to 5 chars, got %s", ti.Value)
	}
}

func TestTextInputSetValueCursorAdjustment(t *testing.T) {
	ti := NewTextInput()
	ti.Cursor = 10
	ti.SetValue("hi")
	if ti.Cursor != 2 {
		t.Errorf("expected Cursor=2, got %d", ti.Cursor)
	}
}

func TestTextInputGetValue(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("test")
	if ti.GetValue() != "test" {
		t.Error("expected GetValue to return test")
	}
}

func TestTextInputClear(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("test")
	ti.Cursor = 4

	called := false
	ti.OnChange = func(s string) {
		called = true
	}

	ti.Clear()
	if ti.Value != "" {
		t.Error("expected empty Value")
	}
	if ti.Cursor != 0 {
		t.Error("expected Cursor=0")
	}
	if !called {
		t.Error("expected OnChange to be called")
	}
}

func TestTextInputAppend(t *testing.T) {
	ti := NewTextInput()
	ti.Append('a')
	if ti.Value != "a" {
		t.Errorf("expected Value=a, got %s", ti.Value)
	}
	if ti.Cursor != 1 {
		t.Errorf("expected Cursor=1, got %d", ti.Cursor)
	}
}

func TestTextInputAppendAtMiddle(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("abc")
	ti.Cursor = 1
	ti.Append('x')
	if ti.Value != "axbc" {
		t.Errorf("expected Value=axbc, got %s", ti.Value)
	}
}

func TestTextInputAppendMaxLen(t *testing.T) {
	ti := NewTextInput(TextInputWithMaxLen(3))
	ti.SetValue("abc")
	ti.Append('d')
	if ti.Value != "abc" {
		t.Errorf("expected Value=abc (unchanged), got %s", ti.Value)
	}
}

func TestTextInputBackspace(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("abc")
	ti.Cursor = 3
	ti.Backspace()
	if ti.Value != "ab" {
		t.Errorf("expected Value=ab, got %s", ti.Value)
	}
	if ti.Cursor != 2 {
		t.Errorf("expected Cursor=2, got %d", ti.Cursor)
	}
}

func TestTextInputBackspaceAtStart(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("abc")
	ti.Cursor = 0
	ti.Backspace()
	if ti.Value != "abc" {
		t.Errorf("expected Value=abc (unchanged), got %s", ti.Value)
	}
}

func TestTextInputDelete(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("abc")
	ti.Cursor = 1
	ti.Delete()
	if ti.Value != "ac" {
		t.Errorf("expected Value=ac, got %s", ti.Value)
	}
}

func TestTextInputDeleteAtEnd(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("abc")
	ti.Cursor = 3
	ti.Delete()
	if ti.Value != "abc" {
		t.Errorf("expected Value=abc (unchanged), got %s", ti.Value)
	}
}

func TestTextInputMoveLeft(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("abc")
	ti.Cursor = 2
	ti.MoveLeft()
	if ti.Cursor != 1 {
		t.Errorf("expected Cursor=1, got %d", ti.Cursor)
	}
}

func TestTextInputMoveLeftAtStart(t *testing.T) {
	ti := NewTextInput()
	ti.Cursor = 0
	ti.MoveLeft()
	if ti.Cursor != 0 {
		t.Errorf("expected Cursor=0, got %d", ti.Cursor)
	}
}

func TestTextInputMoveRight(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("abc")
	ti.Cursor = 1
	ti.MoveRight()
	if ti.Cursor != 2 {
		t.Errorf("expected Cursor=2, got %d", ti.Cursor)
	}
}

func TestTextInputMoveRightAtEnd(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("abc")
	ti.Cursor = 3
	ti.MoveRight()
	if ti.Cursor != 3 {
		t.Errorf("expected Cursor=3 (unchanged), got %d", ti.Cursor)
	}
}

func TestTextInputMoveToStart(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("abc")
	ti.Cursor = 3
	ti.MoveToStart()
	if ti.Cursor != 0 {
		t.Errorf("expected Cursor=0, got %d", ti.Cursor)
	}
}

func TestTextInputMoveToEnd(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("abc")
	ti.Cursor = 0
	ti.MoveToEnd()
	if ti.Cursor != 3 {
		t.Errorf("expected Cursor=3, got %d", ti.Cursor)
	}
}

func TestTextInputSubmit(t *testing.T) {
	called := false
	ti := NewTextInput(TextInputWithOnSubmit(func(s string) {
		called = true
	}))
	ti.SetValue("test")
	ti.Submit()
	if !called {
		t.Error("expected OnSubmit to be called")
	}
}

func TestTextInputCancel(t *testing.T) {
	called := false
	ti := NewTextInput(TextInputWithOnCancel(func() {
		called = true
	}))
	ti.Cancel()
	if !called {
		t.Error("expected OnCancel to be called")
	}
}

func TestTextInputRenderDisplay(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("secret")
	if ti.RenderDisplay() != "secret" {
		t.Error("expected RenderDisplay to return secret")
	}
}

func TestTextInputRenderDisplayHidden(t *testing.T) {
	ti := NewTextInput(TextInputWithHidden(true))
	ti.SetValue("secret")
	if ti.RenderDisplay() != "••••••" {
		t.Error("expected RenderDisplay to return bullets")
	}
}

func TestTextInputRender(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("test")
	result := ti.Render()
	if result == "" {
		t.Error("expected non-empty render")
	}
}

func TestTextInputRenderWithPrefix(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("test")
	result := ti.RenderWithPrefix("Label:")
	if result == "" {
		t.Error("expected non-empty render")
	}
}

func TestTextInputRenderLabeled(t *testing.T) {
	ti := NewTextInput()
	ti.SetValue("test")
	result := ti.RenderLabeled("Enter:")
	if result == "" {
		t.Error("expected non-empty render")
	}
}