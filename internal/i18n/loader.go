package i18n

import (
	"embed"
	"os"
	"strings"

	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var LocaleFS embed.FS

type Loader struct {
	language  language.Tag
	supported []language.Tag
}

func NewLoader(defaultLang string) (*Loader, error) {
	tag, err := language.Parse(defaultLang)
	if err != nil {
		tag = language.English
	}
	return &Loader{language: tag}, nil
}

func (l *Loader) T(templateID string, args ...interface{}) string {
	return templateID
}

func (l *Loader) Language() language.Tag {
	return l.language
}

func (l *Loader) SupportedLanguages() []language.Tag {
	return l.supported
}

func IsRTL(tag string) bool {
	switch strings.ToLower(tag) {
	case "ar", "he", "ur", "fa", "yi":
		return true
	default:
		return false
	}
}

func DetectLanguage() string {
	if lang := os.Getenv("FREECODE_LANG"); lang != "" {
		return lang
	}
	if lang := os.Getenv("LANG"); lang != "" {
		if idx := strings.Index(lang, "."); idx > 0 {
			lang = lang[:idx]
		}
		return strings.ReplaceAll(lang, "_", "-")
	}
	if lang := os.Getenv("LC_ALL"); lang != "" {
		if idx := strings.Index(lang, "."); idx > 0 {
			lang = lang[:idx]
		}
		return strings.ReplaceAll(lang, "_", "-")
	}
	return "en"
}
