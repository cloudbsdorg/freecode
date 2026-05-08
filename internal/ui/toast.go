package ui

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

type ToastVariant string

const (
	ToastVariantInfo    ToastVariant = "info"
	ToastVariantSuccess ToastVariant = "success"
	ToastVariantWarning ToastVariant = "warning"
	ToastVariantError   ToastVariant = "error"
)

type Toast struct {
	ID        string
	Title     string
	Message   string
	Variant   ToastVariant
	Duration  time.Duration
	CreatedAt time.Time
}

type ToastManager struct {
	toasts     []Toast
	maxToasts  int
	width      int
	height     int
	isVisible  bool
}

func NewToastManager() *ToastManager {
	return &ToastManager{
		toasts:    make([]Toast, 0),
		maxToasts: 5,
		width:     50,
		height:    3,
		isVisible: false,
	}
}

func (t *ToastManager) Show(title, message string, variant ToastVariant, duration time.Duration) {
	toast := Toast{
		ID:        time.Now().String(),
		Title:     title,
		Message:   message,
		Variant:   variant,
		Duration:  duration,
		CreatedAt: time.Now(),
	}
	t.toasts = append(t.toasts, toast)
	if len(t.toasts) > t.maxToasts {
		t.toasts = t.toasts[1:]
	}
	t.isVisible = true
}

func (t *ToastManager) ShowInfo(message string) {
	t.Show("", message, ToastVariantInfo, 3000*time.Millisecond)
}

func (t *ToastManager) ShowSuccess(message string) {
	t.Show("", message, ToastVariantSuccess, 3000*time.Millisecond)
}

func (t *ToastManager) ShowWarning(message string) {
	t.Show("", message, ToastVariantWarning, 5000*time.Millisecond)
}

func (t *ToastManager) ShowError(message string) {
	t.Show("", message, ToastVariantError, 5000*time.Millisecond)
}

func (t *ToastManager) Dismiss() {
	if len(t.toasts) > 0 {
		t.toasts = t.toasts[1:]
	}
	if len(t.toasts) == 0 {
		t.isVisible = false
	}
}

func (t *ToastManager) IsVisible() bool {
	return t.isVisible && len(t.toasts) > 0
}

func (t *ToastManager) CurrentToast() *Toast {
	if len(t.toasts) == 0 {
		return nil
	}
	return &t.toasts[0]
}

func (t *ToastManager) SetWidth(w int) {
	t.width = w
}

func (t *ToastManager) SetHeight(h int) {
	t.height = h
}

func (t *ToastManager) Render() string {
	toast := t.CurrentToast()
	if toast == nil {
		return ""
	}

	variantColor := lipgloss.Color("#007ACC")
	switch toast.Variant {
	case ToastVariantSuccess:
		variantColor = lipgloss.Color("#4EC9B0")
	case ToastVariantWarning:
		variantColor = lipgloss.Color("#DCDCAA")
	case ToastVariantError:
		variantColor = lipgloss.Color("#F44747")
	}

	borderStyle := lipgloss.NewStyle().
		BorderForeground(variantColor).
		Width(t.width).
		Padding(1, 2)

	content := ""
	if toast.Title != "" {
		titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
		content += titleStyle.Render(toast.Title) + "\n"
	}
	content += toast.Message

	return borderStyle.Render(content)
}