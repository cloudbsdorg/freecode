package syntax

import "github.com/charmbracelet/lipgloss"

type Style struct {
	Keyword     lipgloss.Style
	String      lipgloss.Style
	Number      lipgloss.Style
	Comment     lipgloss.Style
	Function    lipgloss.Style
	Variable    lipgloss.Style
	Type        lipgloss.Style
	Operator    lipgloss.Style
	Punctuation lipgloss.Style
	Error       lipgloss.Style
	Warning     lipgloss.Style
	Info        lipgloss.Style
	Debug       lipgloss.Style
}

func NewStyle(cfg func(*Style)) Style {
	s := Style{}
	if cfg != nil {
		cfg(&s)
	}
	return s
}

type Theme struct {
	Name      string
	Styles    Style
	initialized bool
}

func DefaultTheme() Theme {
	return Theme{
		Name: "default",
		Styles: Style{
			Keyword:     lipgloss.NewStyle().Foreground(lipgloss.Color("#569CD6")).Bold(true),
			String:      lipgloss.NewStyle().Foreground(lipgloss.Color("#CE9178")),
			Number:      lipgloss.NewStyle().Foreground(lipgloss.Color("#B5CEA8")),
			Comment:     lipgloss.NewStyle().Foreground(lipgloss.Color("#6A9955")).Italic(true),
			Function:    lipgloss.NewStyle().Foreground(lipgloss.Color("#DCDCAA")),
			Variable:    lipgloss.NewStyle().Foreground(lipgloss.Color("#9CDCFE")),
			Type:        lipgloss.NewStyle().Foreground(lipgloss.Color("#4EC9B0")),
			Operator:    lipgloss.NewStyle().Foreground(lipgloss.Color("#D4D4D4")),
			Punctuation: lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")),
			Error:       lipgloss.NewStyle().Foreground(lipgloss.Color("#F14C4C")),
			Warning:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")),
			Info:        lipgloss.NewStyle().Foreground(lipgloss.Color("#4EC9B0")),
			Debug:       lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")),
		},
		initialized: true,
	}
}

func DarkPlusTheme() Theme {
	return DefaultTheme()
}


