package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type ErrorHandler interface {
	Handle(error)
}

type TeaComponent interface {
	Update(msg interface{}) (TeaComponent, interface{})
	View() string
}

type ErrorBoundary struct {
	component     TeaComponent
	err           error
	recovering    bool
	width         int
	height        int
	isVisible     bool
	handler       ErrorHandler
	autoRecover   bool
	autoRecoverAt time.Time
}

func NewErrorBoundary() *ErrorBoundary {
	return &ErrorBoundary{
		width:       60,
		height:      15,
		isVisible:   false,
		autoRecover: true,
	}
}

func (e *ErrorBoundary) Wrap(component TeaComponent) *ErrorBoundary {
	e.component = component
	return e
}

func (e *ErrorBoundary) SetHandler(handler ErrorHandler) *ErrorBoundary {
	e.handler = handler
	return e
}

func (e *ErrorBoundary) SetAutoRecover(autoRecover bool) *ErrorBoundary {
	e.autoRecover = autoRecover
	return e
}

func (e *ErrorBoundary) SetWidth(width int) {
	e.width = width
}

func (e *ErrorBoundary) SetHeight(height int) {
	e.height = height
}

func (e *ErrorBoundary) Update(msg interface{}) (*ErrorBoundary, interface{}) {
	if e.component == nil {
		return e, nil
	}

	if e.recovering && e.autoRecover {
		if time.Now().After(e.autoRecoverAt) {
			e.Recover()
		}
	}

	var nextComponent TeaComponent
	var updateResult interface{}

	func() {
		defer func() {
			if r := recover(); r != nil {
				e.handlePanic(r)
			}
		}()
		nextComponent, updateResult = e.component.Update(msg)
	}()

	if e.err != nil && !e.isVisible {
		e.isVisible = true
	}

	if nextComponent != nil {
		e.component = nextComponent
	}

	return e, updateResult
}

func (e *ErrorBoundary) View() string {
	if e.component != nil && !e.isVisible {
		return e.component.View()
	}

	return e.Render()
}

func (e *ErrorBoundary) handlePanic(r interface{}) {
	var err error
	switch v := r.(type) {
	case error:
		err = v
	case string:
		err = fmt.Errorf("panic: %s", v)
	default:
		err = fmt.Errorf("panic: %v", v)
	}

	e.err = err
	e.recovering = true
	e.autoRecoverAt = time.Now().Add(5 * time.Second)

	if e.handler != nil {
		e.handler.Handle(err)
	}
}

func (e *ErrorBoundary) Handle(err error) {
	e.err = err
	e.recovering = true
	e.autoRecoverAt = time.Now().Add(5 * time.Second)
	e.isVisible = true

	if e.handler != nil {
		e.handler.Handle(err)
	}
}

func (e *ErrorBoundary) IsRecovering() bool {
	return e.recovering
}

func (e *ErrorBoundary) IsVisible() bool {
	return e.isVisible
}

func (e *ErrorBoundary) Error() error {
	return e.err
}

func (e *ErrorBoundary) Recover() {
	e.err = nil
	e.recovering = false
	e.isVisible = false
}

func (e *ErrorBoundary) Dismiss() {
	e.isVisible = false
}

func (e *ErrorBoundary) Retry() {
	e.Recover()
}

func (e *ErrorBoundary) Render() string {
	if !e.isVisible {
		return ""
	}

	if e.component != nil {
		return e.component.View() + "\n" + e.renderErrorDialog()
	}

	return e.renderErrorDialog()
}

func (e *ErrorBoundary) renderErrorDialog() string {
	dialogStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1E1E1E")).
		Border(lipgloss.HiddenBorder()).
		Width(e.width).
		Height(e.height)

	return dialogStyle.Render(e.renderErrorContent())
}

func (e *ErrorBoundary) renderErrorContent() string {
	var lines []string

	lines = append(lines, e.renderHeader())
	lines = append(lines, "")
	lines = append(lines, e.renderErrorInfo()...)
	lines = append(lines, "")
	lines = append(lines, e.renderActions()...)
	lines = append(lines, "")
	lines = append(lines, e.renderHints())

	return strings.Join(lines, "\n")
}

func (e *ErrorBoundary) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F44747")).
		Bold(true)
	return headerStyle.Render("⚠ Error")
}

func (e *ErrorBoundary) renderErrorInfo() []string {
	var lines []string

	if e.err == nil {
		return lines
	}

	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	lines = append(lines, "  An error occurred:")
	lines = append(lines, "")

	errorMsg := e.err.Error()
	if len(errorMsg) > e.width-10 {
		words := strings.Split(errorMsg, " ")
		var currentLine string
		for _, word := range words {
			if len(currentLine)+len(word) > e.width-10 {
				if currentLine != "" {
					lines = append(lines, "    "+errorStyle.Render(currentLine))
				}
				currentLine = word
			} else {
				if currentLine == "" {
					currentLine = word
				} else {
					currentLine += " " + word
				}
			}
		}
		if currentLine != "" {
			lines = append(lines, "    "+errorStyle.Render(currentLine))
		}
	} else {
		lines = append(lines, "    "+errorStyle.Render(errorMsg))
	}

	return lines
}

func (e *ErrorBoundary) renderActions() []string {
	var lines []string

	if e.autoRecover && e.recovering {
		recoveringStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4EC9B0"))
		remaining := time.Until(e.autoRecoverAt)
		if remaining > 0 {
			lines = append(lines, recoveringStyle.Render(fmt.Sprintf("  Auto-recovering in %ds...", int(remaining.Seconds()))))
		}
	}

	lines = append(lines, "  [R]etry   [D]ismiss   [C]ontinue")

	return lines
}

func (e *ErrorBoundary) renderHints() string {
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	return hintStyle.Render("  r: retry  d: dismiss  c: continue")
}

func (e *ErrorBoundary) HandleKey(key string) bool {
	if !e.isVisible {
		return false
	}

	switch key {
	case "r", "R", "enter":
		e.Retry()
		return true
	case "d", "D", "escape":
		e.Dismiss()
		return true
	case "c", "C":
		e.Recover()
		return true
	}

	return false
}