package ui

import (
	"errors"
	"testing"
)

type mockComponent struct {
	shouldPanic bool
	updateCount int
}

func (m *mockComponent) Update(msg interface{}) (TeaComponent, interface{}) {
	m.updateCount++
	if m.shouldPanic {
		panic("test panic")
	}
	return m, nil
}

func (m *mockComponent) View() string {
	return "mock view"
}

type mockHandler struct {
	handledError error
}

func (h *mockHandler) Handle(err error) {
	h.handledError = err
}

func TestNewErrorBoundary(t *testing.T) {
	eb := NewErrorBoundary()
	if eb == nil {
		t.Fatal("NewErrorBoundary() returned nil")
	}
	if eb.IsVisible() {
		t.Error("IsVisible() = true, want false")
	}
	if eb.IsRecovering() {
		t.Error("IsRecovering() = true, want false")
	}
}

func TestErrorBoundaryWrap(t *testing.T) {
	eb := NewErrorBoundary()
	comp := &mockComponent{}
	result := eb.Wrap(comp)
	if result != eb {
		t.Error("Wrap() should return the same ErrorBoundary")
	}
	if eb.component != comp {
		t.Error("Wrap() should set the component")
	}
}

func TestErrorBoundaryViewWithoutError(t *testing.T) {
	eb := NewErrorBoundary()
	comp := &mockComponent{}
	eb.Wrap(comp)

	view := eb.View()
	if view != "mock view" {
		t.Errorf("View() = %q, want %q", view, "mock view")
	}
}

func TestErrorBoundaryPanicRecovery(t *testing.T) {
	eb := NewErrorBoundary()
	comp := &mockComponent{shouldPanic: true}
	eb.Wrap(comp)

	eb.Update("test message")

	if !eb.IsRecovering() {
		t.Error("IsRecovering() = false, want true after panic")
	}
	if !eb.IsVisible() {
		t.Error("IsVisible() = false, want true after panic")
	}
	if eb.Error() == nil {
		t.Error("Error() = nil, want error after panic")
	}
}

func TestErrorBoundaryHandleError(t *testing.T) {
	eb := NewErrorBoundary()
	testErr := errors.New("test error")

	eb.Handle(testErr)

	if !eb.IsRecovering() {
		t.Error("IsRecovering() = false, want true after Handle()")
	}
	if !eb.IsVisible() {
		t.Error("IsVisible() = false, want true after Handle()")
	}
	if eb.Error() != testErr {
		t.Errorf("Error() = %v, want %v", eb.Error(), testErr)
	}
}

func TestErrorBoundaryRecover(t *testing.T) {
	eb := NewErrorBoundary()
	eb.Handle(errors.New("test error"))

	eb.Recover()

	if eb.IsRecovering() {
		t.Error("IsRecovering() = true, want false after Recover()")
	}
	if eb.IsVisible() {
		t.Error("IsVisible() = true, want false after Recover()")
	}
	if eb.Error() != nil {
		t.Errorf("Error() = %v, want nil after Recover()", eb.Error())
	}
}

func TestErrorBoundaryDismiss(t *testing.T) {
	eb := NewErrorBoundary()
	eb.Handle(errors.New("test error"))

	eb.Dismiss()

	if eb.IsVisible() {
		t.Error("IsVisible() = true, want false after Dismiss()")
	}
	if !eb.IsRecovering() {
		t.Error("IsRecovering() = false, want true after Dismiss()")
	}
}

func TestErrorBoundaryRetry(t *testing.T) {
	eb := NewErrorBoundary()
	eb.Handle(errors.New("test error"))

	eb.Retry()

	if eb.IsRecovering() {
		t.Error("IsRecovering() = true, want false after Retry()")
	}
	if eb.IsVisible() {
		t.Error("IsVisible() = true, want false after Retry()")
	}
}

func TestErrorBoundarySetHandler(t *testing.T) {
	eb := NewErrorBoundary()
	handler := &mockHandler{}

	eb.SetHandler(handler)
	eb.Handle(errors.New("test error"))

	if handler.handledError == nil {
		t.Error("handler.Handle() was not called")
	}
}

func TestErrorBoundaryAutoRecoverDisabled(t *testing.T) {
	eb := NewErrorBoundary()
	eb.SetAutoRecover(false)
	comp := &mockComponent{shouldPanic: true}
	eb.Wrap(comp)

	eb.Update("test")

	if !eb.IsRecovering() {
		t.Error("IsRecovering() = false, want true (error was caught)")
	}
	if !eb.IsVisible() {
		t.Error("IsVisible() = false, want true (error dialog should show)")
	}
}

func TestErrorBoundaryHandleKey(t *testing.T) {
	eb := NewErrorBoundary()
	eb.Handle(errors.New("test error"))

	if !eb.HandleKey("r") {
		t.Error("HandleKey(\"r\") = false, want true")
	}
	if eb.IsRecovering() {
		t.Error("IsRecovering() = true, want false after retry key")
	}

	eb.Handle(errors.New("test error"))
	if !eb.HandleKey("d") {
		t.Error("HandleKey(\"d\") = false, want true")
	}
	if eb.IsVisible() {
		t.Error("IsVisible() = true, want false after dismiss key")
	}

	eb.Handle(errors.New("test error"))
	if !eb.HandleKey("escape") {
		t.Error("HandleKey(\"escape\") = false, want true")
	}

	eb.Handle(errors.New("test error"))
	if !eb.HandleKey("c") {
		t.Error("HandleKey(\"c\") = false, want true")
	}
	if eb.IsRecovering() {
		t.Error("IsRecovering() = true, want false after continue key")
	}
}

func TestErrorBoundaryHandleKeyNotVisible(t *testing.T) {
	eb := NewErrorBoundary()

	if eb.HandleKey("r") {
		t.Error("HandleKey(\"r\") = true, want false when not visible")
	}
}

func TestErrorBoundaryRender(t *testing.T) {
	eb := NewErrorBoundary()
	eb.SetWidth(60)
	eb.SetHeight(15)
	eb.Handle(errors.New("test error"))

	rendered := eb.Render()
	if rendered == "" {
		t.Error("Render() returned empty string")
	}
}

func TestErrorBoundaryRenderNotVisible(t *testing.T) {
	eb := NewErrorBoundary()

	rendered := eb.Render()
	if rendered != "" {
		t.Errorf("Render() = %q, want empty string when not visible", rendered)
	}
}