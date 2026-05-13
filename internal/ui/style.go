package ui

import (
	"github.com/freecode/freecode/internal/style"
)

var (
	ActiveTab = style.NewStyle().
			Foreground("#FFFFFF").
			Background("#007ACC").
			Padding(0, 1)

	InactiveTab = style.NewStyle().
			Foreground("#808080").
			Background("#2D2D2D").
			Padding(0, 1)

	TabBar = style.NewStyle().
		Background("#2D2D2D")

	Content = style.NewStyle().
		Foreground("#E0E0E0").
		Background("#1E1E1E")

	StatusLine = style.NewStyle().
			Foreground("#808080").
			Background("#2D2D2D")

	Error = style.NewStyle().
		Foreground("#F44747")

	Warning = style.NewStyle().
		Foreground("#FFCC00")

	Success = style.NewStyle().
		Foreground("#4EC9B0")
)

type Theme struct {
	Name        string
	ActiveTab   style.Style
	InactiveTab style.Style
	TabBar      style.Style
	Content     style.Style
	StatusLine  style.Style
	Error       style.Style
	Warning     style.Style
	Success     style.Style
	Prompt      style.Style
	UserMsg     style.Style
	AssistantMsg style.Style
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
		Prompt: style.NewStyle().
			Foreground("#FFFFFF").
			Background("#3C3C3C"),
		UserMsg: style.NewStyle().
			Foreground("#4EC9B0"),
		AssistantMsg: style.NewStyle().
			Foreground("#DCDCAA"),
	}
}

func DarkTheme() Theme {
	return DefaultTheme()
}

func LightTheme() Theme {
	return Theme{
		Name:   "light",
		ActiveTab: style.NewStyle().
			Foreground("#FFFFFF").
			Background("#0066CC").
			Padding(0, 1),
		InactiveTab: style.NewStyle().
			Foreground("#333333").
			Background("#E0E0E0").
			Padding(0, 1),
		TabBar: style.NewStyle().
			Background("#E0E0E0"),
		Content: style.NewStyle().
			Foreground("#333333").
			Background("#FFFFFF"),
		StatusLine: style.NewStyle().
			Foreground("#666666").
			Background("#E0E0E0"),
		Error:   Error,
		Warning: Warning,
		Success: Success,
		Prompt: style.NewStyle().
			Foreground("#FFFFFF").
			Background("#0066CC"),
		UserMsg: style.NewStyle().
			Foreground("#0066CC"),
		AssistantMsg: style.NewStyle().
			Foreground("#8B6914"),
	}
}

func DraculaTheme() Theme {
	return Theme{
		Name:   "dracula",
		ActiveTab: style.NewStyle().
			Foreground("#FFFFFF").
			Background("#BD93F9").
			Padding(0, 1),
		InactiveTab: style.NewStyle().
			Foreground("#6272A4").
			Background("#44475A").
			Padding(0, 1),
		TabBar: style.NewStyle().
			Background("#44475A"),
		Content: style.NewStyle().
			Foreground("#F8F8F2").
			Background("#282A36"),
		StatusLine: style.NewStyle().
			Foreground("#F8F8F2").
			Background("#44475A"),
		Error: style.NewStyle().
			Foreground("#FF5555"),
		Warning: style.NewStyle().
			Foreground("#F1FA8C"),
		Success: style.NewStyle().
			Foreground("#50FA7B"),
		Prompt: style.NewStyle().
			Foreground("#FF79C6").
			Background("#44475A"),
		UserMsg: style.NewStyle().
			Foreground("#8BE9FD"),
		AssistantMsg: style.NewStyle().
			Foreground("#F1FA8C"),
	}
}

func NordTheme() Theme {
	return Theme{
		Name:   "nord",
		ActiveTab: style.NewStyle().
			Foreground("#ECEFF4").
			Background("#81A1C1").
			Padding(0, 1),
		InactiveTab: style.NewStyle().
			Foreground("#D8DEE9").
			Background("#3B4252").
			Padding(0, 1),
		TabBar: style.NewStyle().
			Background("#3B4252"),
		Content: style.NewStyle().
			Foreground("#ECEFF4").
			Background("#2E3440"),
		StatusLine: style.NewStyle().
			Foreground("#D8DEE9").
			Background("#3B4252"),
		Error: style.NewStyle().
			Foreground("#BF616A"),
		Warning: style.NewStyle().
			Foreground("#EBCB8B"),
		Success: style.NewStyle().
			Foreground("#A3BE8C"),
		Prompt: style.NewStyle().
			Foreground("#88C0D0").
			Background("#3B4252"),
		UserMsg: style.NewStyle().
			Foreground("#81A1C1"),
		AssistantMsg: style.NewStyle().
			Foreground("#8FBCBB"),
	}
}

func MonokaiTheme() Theme {
	return Theme{
		Name:   "monokai",
		ActiveTab: style.NewStyle().
			Foreground("#F8F8F2").
			Background("#F92672").
			Padding(0, 1),
		InactiveTab: style.NewStyle().
			Foreground("#75715E").
			Background("#3E3D32").
			Padding(0, 1),
		TabBar: style.NewStyle().
			Background("#3E3D32"),
		Content: style.NewStyle().
			Foreground("#F8F8F2").
			Background("#272822"),
		StatusLine: style.NewStyle().
			Foreground("#F8F8F2").
			Background("#3E3D32"),
		Error: style.NewStyle().
			Foreground("#F92672"),
		Warning: style.NewStyle().
			Foreground("#E6DB74"),
		Success: style.NewStyle().
			Foreground("#A6E22E"),
		Prompt: style.NewStyle().
			Foreground("#F92672").
			Background("#3E3D32"),
		UserMsg: style.NewStyle().
			Foreground("#66D9EF"),
		AssistantMsg: style.NewStyle().
			Foreground("#A6E22E"),
	}
}

func GruvboxTheme() Theme {
	return Theme{
		Name:   "gruvbox",
		ActiveTab: style.NewStyle().
			Foreground("#282828").
			Background("#FABD2F").
			Padding(0, 1),
		InactiveTab: style.NewStyle().
			Foreground("#A89984").
			Background("#3C3836").
			Padding(0, 1),
		TabBar: style.NewStyle().
			Background("#3C3836"),
		Content: style.NewStyle().
			Foreground("#EBDBB2").
			Background("#282828"),
		StatusLine: style.NewStyle().
			Foreground("#EBDBB2").
			Background("#3C3836"),
		Error: style.NewStyle().
			Foreground("#FB4934"),
		Warning: style.NewStyle().
			Foreground("#FABD2F"),
		Success: style.NewStyle().
			Foreground("#B8BB26"),
		Prompt: style.NewStyle().
			Foreground("#FB4934").
			Background("#3C3836"),
		UserMsg: style.NewStyle().
			Foreground("#83A598"),
		AssistantMsg: style.NewStyle().
			Foreground("#B8BB26"),
	}
}

func SolarizedTheme() Theme {
	return Theme{
		Name:   "solarized",
		ActiveTab: style.NewStyle().
			Foreground("#002B36").
			Background("#268BD2").
			Padding(0, 1),
		InactiveTab: style.NewStyle().
			Foreground("#657B83").
			Background("#073642").
			Padding(0, 1),
		TabBar: style.NewStyle().
			Background("#073642"),
		Content: style.NewStyle().
			Foreground("#839496").
			Background("#002B36"),
		StatusLine: style.NewStyle().
			Foreground("#839496").
			Background("#073642"),
		Error: style.NewStyle().
			Foreground("#DC322F"),
		Warning: style.NewStyle().
			Foreground("#B58900"),
		Success: style.NewStyle().
			Foreground("#859900"),
		Prompt: style.NewStyle().
			Foreground("#268BD2").
			Background("#073642"),
		UserMsg: style.NewStyle().
			Foreground("#2AA198"),
		AssistantMsg: style.NewStyle().
			Foreground("#859900"),
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
