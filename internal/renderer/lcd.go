package renderer

type LCDRenderer struct {
	width  int
	height int
}

func NewLCDRenderer(w, h int) *LCDRenderer {
	return &LCDRenderer{width: w, height: h}
}

func (l *LCDRenderer) RenderBox(x, y, w, h int, bgColor string) string {
	return ""
}

func (l *LCDRenderer) RenderText(text string, x, y int, fgColor string) string {
	if x < 0 || x >= l.width || y < 0 || y >= l.height {
		return ""
	}
	if x+len(text) > l.width {
		text = text[:l.width-x]
	}
	return "\x1b[" + itoa(y+1) + ";" + itoa(x+1) + "H" + text
}

func (l *LCDRenderer) RenderBorder(x, y, w, h int, fgColor string) string {
	return ""
}

func (l *LCDRenderer) RenderSelected(text string, x, y, w int, fg, bg string) string {
	return ""
}

func (l *LCDRenderer) Width() int  { return l.width }
func (l *LCDRenderer) Height() int { return l.height }

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b [20]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(b[pos:])
}
