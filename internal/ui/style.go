package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ActiveTab = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#007ACC")).
			Padding(0, 1)

	InactiveTab = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			Background(lipgloss.Color("#2D2D2D")).
			Padding(0, 1)

	TabBar = lipgloss.NewStyle().
		Background(lipgloss.Color("#2D2D2D"))

	Content = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E0E0E0")).
		Background(lipgloss.Color("#1E1E1E"))

	StatusLine = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			Background(lipgloss.Color("#2D2D2D"))

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F44747"))

	Warning = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFCC00"))

	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4EC9B0"))
)

type Theme struct {
	ActiveTab   lipgloss.Style
	InactiveTab lipgloss.Style
	TabBar      lipgloss.Style
	Content     lipgloss.Style
	StatusLine  lipgloss.Style
	Error       lipgloss.Style
	Warning     lipgloss.Style
	Success     lipgloss.Style
}

func DefaultTheme() Theme {
	return Theme{
		ActiveTab:   ActiveTab,
		InactiveTab: InactiveTab,
		TabBar:      TabBar,
		Content:     Content,
		StatusLine:  StatusLine,
		Error:       Error,
		Warning:     Warning,
		Success:     Success,
	}
}

func DarkTheme() Theme {
	return DefaultTheme()
}

func LightTheme() Theme {
	return Theme{
		ActiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#0066CC")).
			Padding(0, 1),
		InactiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#333333")).
			Background(lipgloss.Color("#E0E0E0")).
			Padding(0, 1),
		TabBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#E0E0E0")),
		Content: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#333333")).
			Background(lipgloss.Color("#FFFFFF")),
		StatusLine: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Background(lipgloss.Color("#E0E0E0")),
		Error:   Error,
		Warning: Warning,
		Success: Success,
	}
}
