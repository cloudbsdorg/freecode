package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

var BannerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#22C55E")).
	Bold(true)

func GetBanner() string {
	f := figure.NewFigure("FREECODE", "cosmike", true)
	lines := strings.Split(f.String(), "\n")
	for i, line := range lines {
		lines[i] = BannerStyle.Render(line)
	}
	return strings.Join(lines, "\n")
}
