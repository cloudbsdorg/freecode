package i18n

import (
	"embed"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"golang.org/x/text/language"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var LocaleFS embed.FS

type Loader struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
	language  language.Tag
}

func NewLoader(defaultLang string) (*Loader, error) {
	tag, err := language.Parse(defaultLang)
	if err != nil {
		tag = language.English
	}

	bundle := i18n.NewBundle(tag)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	l := &Loader{
		bundle:  bundle,
		language: tag,
	}

	if err := l.loadLocales(); err != nil {
		return nil, err
	}

	l.localizer = i18n.NewLocalizer(bundle, defaultLang)
	return l, nil
}

func (l *Loader) loadLocales() error {
	entries, err := LocaleFS.ReadDir("locales")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}
		if strings.HasPrefix(entry.Name(), "translate.") {
			continue
		}

		data, err := LocaleFS.ReadFile("locales/" + entry.Name())
		if err != nil {
			continue
		}

		if _, err := l.bundle.ParseMessageFileBytes(data, entry.Name()); err != nil {
			continue
		}
	}
	return nil
}

func (l *Loader) T(templateID string, args ...interface{}) string {
	return l.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: templateID},
		TemplateData:   args,
	})
}

func (l *Loader) Language() language.Tag {
	return l.language
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
