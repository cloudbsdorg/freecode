package ui

import (
	"strings"

	"github.com/freecode/freecode/internal/style"
	"github.com/common-nighthawk/go-figure"
)

var BannerStyle = style.NewStyle().
	Foreground(style.Color("#22C55E")).
	Bold(true)

func GetBanner() string {
	f := figure.NewFigure("FREECODE", "cosmike", true)
	lines := strings.Split(f.String(), "\n")
	for i, line := range lines {
		lines[i] = BannerStyle.Render(line)
	}
	return strings.Join(lines, "\n")
}
