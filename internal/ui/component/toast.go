package component

import "github.com/freecode/freecode/internal/renderer"

type Toast[R renderer.Renderer] struct {
	Component[R]
	Message   string
	Type      ToastType
	Colors    ToastColors
	AutoHide  bool
	HideAfter int
}

type ToastType int

const (
	ToastInfo ToastType = iota
	ToastSuccess
	ToastWarning
	ToastError
)

type ToastColors struct {
	Background string
	Foreground string
	InfoColor  string
	SuccessColor string
	WarningColor string
	ErrorColor   string
}

func NewToast[R renderer.Renderer](message string, toastType ToastType, colors ToastColors) *Toast[R] {
	return &Toast[R]{
		Component: Component[R]{
			X:       0,
			Y:       0,
			Width:   len(message) + 4,
			Height:  1,
			Visible: true,
		},
		Message: message,
		Type:    toastType,
		Colors:  colors,
	}
}

func (t *Toast[R]) SetMessage(msg string) {
	t.Message = msg
	t.Width = len(msg) + 4
}

func (t *Toast[R]) SetType(toastType ToastType) {
	t.Type = toastType
}

func (t *Toast[R]) Render(r R) string {
	if !t.Visible {
		return ""
	}

	var color string
	switch t.Type {
	case ToastSuccess:
		color = t.Colors.SuccessColor
	case ToastWarning:
		color = t.Colors.WarningColor
	case ToastError:
		color = t.Colors.ErrorColor
	default:
		color = t.Colors.InfoColor
	}

	prefix := " "
	if t.Type == ToastSuccess {
		prefix = "✓ "
	} else if t.Type == ToastError {
		prefix = "✗ "
	} else if t.Type == ToastWarning {
		prefix = "⚠ "
	}

	return r.RenderSelected(prefix+t.Message+" ", t.X, t.Y, t.Width, color, t.Colors.Background)
}
