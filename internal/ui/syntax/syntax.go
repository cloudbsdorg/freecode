package syntax

type StatementType int

const (
	TypeKeyword StatementType = iota
	TypeString
	TypeNumber
	TypeComment
	TypeFunction
	TypeVariable
	TypeType
	TypeOperator
	TypePunctuation
	TypeError
	TypeWarning
	TypeInfo
	TypeDebug
)

var typeNames = map[StatementType]string{
	TypeKeyword:     "keyword",
	TypeString:       "string",
	TypeNumber:       "number",
	TypeComment:      "comment",
	TypeFunction:     "function",
	TypeVariable:     "variable",
	TypeType:         "type",
	TypeOperator:     "operator",
	TypePunctuation:  "punctuation",
	TypeError:        "error",
	TypeWarning:      "warning",
	TypeInfo:         "info",
	TypeDebug:        "debug",
}

func (t StatementType) String() string {
	return typeNames[t]
}

type Language int

const (
	LangGo Language = iota
	LangTypeScript
	LangJavaScript
	LangPython
	LangJSON
	LangYAML
	LangMarkdown
	LangPlain
)

var langNames = map[Language]string{
	LangGo:         "go",
	LangTypeScript: "typescript",
	LangJavaScript: "javascript",
	LangPython:     "python",
	LangJSON:       "json",
	LangYAML:       "yaml",
	LangMarkdown:   "markdown",
	LangPlain:      "plain",
}

func (l Language) String() string {
	return langNames[l]
}

type Config struct {
	Theme         Theme
	Language      Language
	LineNumbers   bool
	TabWidth      int
	MaxLineWidth  int
	WordWrap      bool
	ShowHidden    bool
	HighlightOnly bool
}

func DefaultConfig() Config {
	return Config{
		Theme:       DefaultTheme(),
		Language:    LangPlain,
		LineNumbers: false,
		TabWidth:    8,
		MaxLineWidth: 0,
		WordWrap:    false,
		ShowHidden:  false,
	}
}

func ForLanguage(lang Language) HighlightFunc {
	switch lang {
	case LangJSON:
		return func(code string, theme Theme) string {
			return HighlightJSON(code, theme)
		}
	default:
		return func(code string, theme Theme) string {
			return Highlight(code, theme)
		}
	}
}
