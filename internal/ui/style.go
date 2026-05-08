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
	Name        string
	ActiveTab   lipgloss.Style
	InactiveTab lipgloss.Style
	TabBar      lipgloss.Style
	Content     lipgloss.Style
	StatusLine  lipgloss.Style
	Error       lipgloss.Style
	Warning     lipgloss.Style
	Success     lipgloss.Style
	Prompt      lipgloss.Style
	UserMsg     lipgloss.Style
	AssistantMsg lipgloss.Style
}

func DefaultTheme() Theme {
	return Theme{
		Name:        "default",
		ActiveTab:   ActiveTab,
		InactiveTab: InactiveTab,
		TabBar:      TabBar,
		Content:     Content,
		StatusLine:  StatusLine,
		Error:       Error,
		Warning:     Warning,
		Success:     Success,
		Prompt: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#3C3C3C")),
		UserMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4EC9B0")),
		AssistantMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DCDCAA")),
	}
}

func DarkTheme() Theme {
	return DefaultTheme()
}

func LightTheme() Theme {
	return Theme{
		Name:   "light",
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
		Prompt: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#0066CC")),
		UserMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0066CC")),
		AssistantMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8B6914")),
	}
}

func DraculaTheme() Theme {
	return Theme{
		Name:   "dracula",
		ActiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#BD93F9")).
			Padding(0, 1),
		InactiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4")).
			Background(lipgloss.Color("#44475A")).
			Padding(0, 1),
		TabBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#44475A")),
		Content: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			Background(lipgloss.Color("#282A36")),
		StatusLine: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			Background(lipgloss.Color("#44475A")),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F1FA8C")),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")),
		Prompt: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF79C6")).
			Background(lipgloss.Color("#44475A")),
		UserMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD")),
		AssistantMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F1FA8C")),
	}
}

func NordTheme() Theme {
	return Theme{
		Name:   "nord",
		ActiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ECEFF4")).
			Background(lipgloss.Color("#81A1C1")).
			Padding(0, 1),
		InactiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D8DEE9")).
			Background(lipgloss.Color("#3B4252")).
			Padding(0, 1),
		TabBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#3B4252")),
		Content: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ECEFF4")).
			Background(lipgloss.Color("#2E3440")),
		StatusLine: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D8DEE9")).
			Background(lipgloss.Color("#3B4252")),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BF616A")),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EBCB8B")),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A3BE8C")),
		Prompt: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#88C0D0")).
			Background(lipgloss.Color("#3B4252")),
		UserMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#81A1C1")),
		AssistantMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8FBCBB")),
	}
}

func MonokaiTheme() Theme {
	return Theme{
		Name:   "monokai",
		ActiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			Background(lipgloss.Color("#F92672")).
			Padding(0, 1),
		InactiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#75715E")).
			Background(lipgloss.Color("#3E3D32")).
			Padding(0, 1),
		TabBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#3E3D32")),
		Content: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			Background(lipgloss.Color("#272822")),
		StatusLine: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			Background(lipgloss.Color("#3E3D32")),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F92672")),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E6DB74")),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A6E22E")),
		Prompt: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F92672")).
			Background(lipgloss.Color("#3E3D32")),
		UserMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#66D9EF")),
		AssistantMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A6E22E")),
	}
}

func GruvboxTheme() Theme {
	return Theme{
		Name:   "gruvbox",
		ActiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282828")).
			Background(lipgloss.Color("#FABD2F")).
			Padding(0, 1),
		InactiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A89984")).
			Background(lipgloss.Color("#3C3836")).
			Padding(0, 1),
		TabBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#3C3836")),
		Content: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EBDBB2")).
			Background(lipgloss.Color("#282828")),
		StatusLine: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EBDBB2")).
			Background(lipgloss.Color("#3C3836")),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FB4934")),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FABD2F")),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B8BB26")),
		Prompt: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FB4934")).
			Background(lipgloss.Color("#3C3836")),
		UserMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#83A598")),
		AssistantMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B8BB26")),
	}
}

func SolarizedTheme() Theme {
	return Theme{
		Name:   "solarized",
		ActiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#002B36")).
			Background(lipgloss.Color("#268BD2")).
			Padding(0, 1),
		InactiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#657B83")).
			Background(lipgloss.Color("#073642")).
			Padding(0, 1),
		TabBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#073642")),
		Content: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#839496")).
			Background(lipgloss.Color("#002B36")),
		StatusLine: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#839496")).
			Background(lipgloss.Color("#073642")),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DC322F")),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B58900")),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#859900")),
		Prompt: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#268BD2")).
			Background(lipgloss.Color("#073642")),
		UserMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2AA198")),
		AssistantMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#859900")),
	}
}

var themes = map[string]Theme{
	"default":   DefaultTheme(),
	"dark":      DarkTheme(),
	"light":     LightTheme(),
	"dracula":   DraculaTheme(),
	"nord":      NordTheme(),
	"monokai":   MonokaiTheme(),
	"gruvbox":   GruvboxTheme(),
	"solarized": SolarizedTheme(),
}

func GetTheme(name string) Theme {
	if t, ok := themes[name]; ok {
		return t
	}
	return DefaultTheme()
}

func ListThemes() []string {
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	return names
}
