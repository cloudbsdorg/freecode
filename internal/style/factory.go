package style

type Type string

const (
	TypeLipgloss Type = "lipgloss"
	TypeGlyph    Type = "glyph"
)

var defaultType = TypeLipgloss

func SetDefault(t Type) {
	defaultType = t
}

func GetDefault() Type {
	return defaultType
}

func NewStyle() Style {
	switch defaultType {
	case TypeGlyph:
		return NewGlyphStyle()
	default:
		return NewLipglossStyle()
	}
}

func Parse(s string) Type {
	switch s {
	case "lipgloss":
		return TypeLipgloss
	case "glyph":
		return TypeGlyph
	default:
		return defaultType
	}
}
