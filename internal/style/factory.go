package style

var defaultType = TypeLipgloss

type Type string

const TypeLipgloss Type = "lipgloss"

func SetDefault(t Type) {
	defaultType = t
}

func GetDefault() Type {
	return defaultType
}

func NewStyle() Style {
	return NewLipglossStyle()
}

func Parse(s string) Type {
	if s == "lipgloss" {
		return TypeLipgloss
	}
	return defaultType
}
