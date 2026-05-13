package style

type Style interface {
	Render(text string) string
	Foreground(color string) Style
	Background(color string) Style
	Bold(v ...bool) Style
	Italic(v ...bool) Style
	Width(w int) Style
	Height(h int) Style
	Padding(values ...int) Style
	Margin(values ...int) Style
	MarginTop(v int) Style
	MarginLeft(v int) Style
	BorderStyle(b Border) Style
	BorderForeground(color string) Style
}

type Border interface{}

type NormalBorder struct{}
type RoundedBorder struct{}
type HiddenBorder struct{}

var NormalBorderType NormalBorder
var RoundedBorderType RoundedBorder
var HiddenBorderType HiddenBorder

func Normal() Border   { return NormalBorderType }
func Rounded() Border  { return RoundedBorderType }
func Hidden() Border   { return HiddenBorderType }
