package syntax

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type HighlightFunc func(string, Theme) string

var highlighters = map[string]HighlightFunc{}

func RegisterHighlighter(lang string, fn HighlightFunc) {
	highlighters[lang] = fn
}

func GetHighlighter(lang string) HighlightFunc {
	return highlighters[lang]
}

func Highlight(code string, theme Theme) string {
	if !theme.initialized {
		theme = DefaultTheme()
	}
	return highlight(code, theme)
}

func highlight(code string, theme Theme) string {
	s := theme.Styles

	code = applyPattern(code, `//.*$`, s.Comment)
	code = applyPattern(code, `/\*[\s\S]*?\*/`, s.Comment)
	code = applyPattern(code, `#.*$`, s.Comment)

	code = applyPattern(code, `"[^"\\]*(\\.[^"\\]*)*"`, s.String)
	code = applyPattern(code, `'[^'\\]*(\\.[^'\\]*)*'`, s.String)
	code = applyPattern(code, "`[^`]*`", s.String)

	code = applyPattern(code, `\b\d+\.?\d*\b`, s.Number)
	code = applyPattern(code, `\b0x[0-9a-fA-F]+\b`, s.Number)

	code = highlightKeywords(code, theme.Styles)

	code = applyPattern(code, `\b([a-zA-Z_][a-zA-Z0-9_]*)\(`,
		theme.Styles.Function)

	return code
}

func applyPattern(code string, pattern string, style lipgloss.Style) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllStringFunc(code, func(match string) string {
		return style.Render(match)
	})
}

func highlightKeywords(code string, style Style) string {
	keywords := []string{
		"func", "struct", "type", "interface", "package", "import",
		"var", "const", "if", "else", "for", "range", "return",
		"switch", "case", "default", "break", "continue", "go", "defer",
		"select", "chan", "map", "make", "new", "true", "false", "nil",
	}

	for _, kw := range keywords {
		pattern := `\b` + kw + `\b`
		re := regexp.MustCompile(pattern)
		code = re.ReplaceAllStringFunc(code, func(match string) string {
			return style.Keyword.Render(match)
		})
	}

	return code
}

func HighlightJSON(code string, theme Theme) string {
	if !theme.initialized {
		theme = DefaultTheme()
	}

	result := code

	keyRe := regexp.MustCompile(`"([^"]+)":`)
	result = keyRe.ReplaceAllString(result, `$1:`)

	stringRe := regexp.MustCompile(`:\s*"([^"]*)"`)
	result = stringRe.ReplaceAllString(result, `: "$1"`)

	numberRe := regexp.MustCompile(`:\s*(\d+\.?\d*)`)
	result = numberRe.ReplaceAllString(result, `: $1`)

	boolRe := regexp.MustCompile(`:\s*(true|false|null)`)
	result = boolRe.ReplaceAllString(result, `: $1`)

	return result
}

func HighlightWithLineNumbers(code string, theme Theme) string {
	if !theme.initialized {
		theme = DefaultTheme()
	}

	lines := strings.Split(code, "\n")
	width := 3
	if len(lines) > 999 {
		width = 4
	}
	if len(lines) > 9999 {
		width = 5
	}

	var result strings.Builder

	for i, line := range lines {
		lineNum := i + 1
		result.WriteString(fmt.Sprintf("%*d", width, lineNum))
		result.WriteString(" │ ")
		result.WriteString(Highlight(line, theme))
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func ExtensibleHighlight(code string, theme Theme, cfg func(*HighlightConfig)) string {
	if !theme.initialized {
		theme = DefaultTheme()
	}

	config := &HighlightConfig{
		ShowLineNumbers: false,
		TabWidth:        8,
		MaxLineWidth:    0,
	}
	if cfg != nil {
		cfg(config)
	}

	result := code

	result = highlightComments(result, theme)
	result = highlightStrings(result, theme)
	result = highlightNumbers(result, theme)
	result = highlightKeywords(result, theme.Styles)
	result = highlightFunctions(result, theme)

	if config.ShowLineNumbers {
		result = addLineNumbers(result, config.TabWidth)
	}

	return result
}

type HighlightConfig struct {
	ShowLineNumbers bool
	TabWidth        int
	MaxLineWidth    int
}

func highlightComments(code string, theme Theme) string {
	s := theme.Styles
	patterns := []string{`//.*$`, `/\*[\s\S]*?\*/`, `#.*$`}
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		code = re.ReplaceAllString(code, "")
	}
	_ = s
	return code
}

func highlightStrings(code string, theme Theme) string {
	s := theme.Styles
	patterns := []string{
		`"[^"\\]*(\\.[^"\\]*)*"`,
		`'[^'\\]*(\\.[^'\\]*)*'`,
		"`[^`]*`",
	}
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		code = re.ReplaceAllString(code, "")
	}
	_ = s
	return code
}

func highlightNumbers(code string, theme Theme) string {
	s := theme.Styles
	patterns := []string{
		`\b\d+\.?\d*\b`,
		`\b0x[0-9a-fA-F]+\b`,
	}
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		code = re.ReplaceAllString(code, "")
	}
	_ = s
	return code
}

func highlightFunctions(code string, theme Theme) string {
	s := theme.Styles
	re := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\(`)
	code = re.ReplaceAllString(code, "$1(")
	_ = s
	return code
}

func addLineNumbers(code string, tabWidth int) string {
	lines := strings.Split(code, "\n")
	width := 3
	if len(lines) > 999 {
		width = 4
	}

	var result strings.Builder
	for i, line := range lines {
		lineNum := i + 1
		result.WriteString(fmt.Sprintf("%*d", width, lineNum))
		result.WriteString(" │ ")
		result.WriteString(line)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}
