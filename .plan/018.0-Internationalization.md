# Internationalization (i18n) Plan

## Overview

Freecode follows the CloudBSD i18n guidelines with full internationalization support. This document covers the technical approach, language support, and implementation tasks.

**References:**
- [CloudBSD Internationalization Guidelines](https://github.com/cloudbsdorg/application_guidelines/blob/main/Internationalization/INTERNATIONALIZATION.md)
- [go-i18n/v2](https://github.com/nicksnyder/go-i18n)
- [golang.org/x/text/unicode/bidi](https://pkg.go.dev/golang.org/x/text/unicode/bidi)

---

## 1. Core Principles

### Text Extraction
- Never hardcode user-facing strings in source code
- All strings stored in external locale files (`locales/*.toml`)
- Use `goi18n` CLI for extract/merge workflow

### Standardized Tools
- **gettext** style: `leonelquinteros/gotext` (PO/MO files)
- **Catalog style**: `go-i18n/v2` (TOML/JSON/YAML) — **PREFERRED**
- **Language tags**: `golang.org/x/text/language` (BCP 47)
- **RTL/Bidi**: `golang.org/x/text/unicode/bidi`

### Encoding
- UTF-8 everywhere (source, config, locale files)
- Handle UTF-8 in all inputs/outputs

### Date/Time
- ISO 8601 (`YYYY-MM-DDTHH:MM:SSZ`) for machine-readable
- Locale-aware formatting via `golang.org/x/text/message`

---

## 2. Language Support

### Required Languages (per CloudBSD guidelines, ordered by native name)

English first, then alphabetically by native name:

| # | Language | Native Name | BCP 47 | RTL |
|---|---------|-------------|--------|-----|
| 1 | English | English | en | ❌ |
| 2 | Spanish | Español | es | ❌ |
| 3 | French | Français | fr | ❌ |
| 4 | Esperanto | Esperanto | eo | ❌ |
| 5 | Italian | Italiano | it | ❌ |
| 6 | Norwegian | Norsk | no | ❌ |
| 7 | Swedish | Svenska | sv | ❌ |
| 8 | Punjabi | ਪੰਜਾਬੀ | pa | ❌ |
| 9 | Klingon | tlhIngan | tlh | ❌ |
| 10 | Elvish | Eledh | el | ❌ |
| 11 | German | Deutsch | de | ❌ |
| 12 | Chinese | 中文 | zh | ❌ |
| 13 | Japanese | 日本語 | ja | ❌ |
| 14 | Arabic | العربية | ar | ✅ |
| 15 | Kiswahili | Kiswahili | sw | ❌ |
| 16 | Yorùbá | Yorùbá | yo | ❌ |
| 17 | Hindi | हिन्दी | hi | ❌ |
| 18 | Korean | 한국어 | ko | ❌ |
| 19 | Finnish | Suomi | fi | ❌ |
| 20 | Russian | Русский | ru | ❌ |
| 21 | Polish | Polski | pl | ❌ |
| 22 | Dothraki | Dothraki | dr | ❌ |
| 23 | Valyrian | Valyrian | va | ❌ |
| 24 | Na'vi | Na'vi | nv | ❌ |
| 25 | Atlantean | Atlantean | at | ❌ |
| 26 | Turkish | Türkçe | tr | ❌ |
| 27 | Catalan | Català | ca | ❌ |
| 28 | Czech | Čeština | cs | ❌ |
| 29 | Greek | Ελληνικά | el | ❌ |
| 30 | Hebrew | עברית | he | ✅ |
| 31 | Ukrainian | Українська | uk | ❌ |
| 32 | Serbian | Српски | sr | ❌ |
| 33 | Slovak | Slovenčina | sk | ❌ |
| 34 | Slovenian | Slovenščina | sl | ❌ |
| 35 | Urdu | اردو | ur | ✅ |
| 36 | Bulgarian | Български | bg | ❌ |
| 37 | Croatian | Hrvatski | hr | ❌ |
| 38 | Hungarian | Magyar | hu | ❌ |
| 39 | Lithuanian | Lietuvių | lt | ❌ |
| 40 | Latvian | Latviešu | lv | ❌ |
| 41 | Indonesian | Bahasa Indonesia | id | ❌ |
| 42 | Portuguese (Brazil) | Português (Brasil) | pt-BR | ❌ |
| 43 | Portuguese (Portugal) | Português (Portugal) | pt-PT | ❌ |
| 44 | Romanian | Română | ro | ❌ |

### Fictional Language Resources

| Language | Resource | Format |
|----------|----------|--------|
| Klingon | [Vaporjawn/Klingon-Translator](https://github.com/Vaporjawn/Klingon-Translator) | TypeScript |
| Klingon | [monkeytypegame/monkeytype](https://github.com/monkeytypegame/monkeytype) | JSON |
| Elvish | [pfstrack/eldamo](https://github.com/pfstrack/eldamo) | XML |
| Elvish | [Omikhleia/sindict](https://github.com/Omikhleia/sindict) | TEI |
| Elvish | [galadhremmin/Parf-Edhellen](https://github.com/galadhremmin/Parf-Edhellen) | DB |
| Dothraki | [chandan-kj/Dothraki-translator](https://github.com/chandan-kj/Dothraki-translator) | JS |
| Valyrian | [shivlloyd/markSeven-valyrian-translation-app](https://github.com/shivlloyd/markSeven-valyrian-translation-app) | JS |
| Na'vi | [Willem3141/navi-reykunyu](https://github.com/Willem3141/navi-reykunyu) | JSON |
| Na'vi | [fwew/fwew](https://github.com/fwew/fwew) | Go CLI |
| Atlantean | [LangMaker](https://langmaker.github.io/atlantean.htm) | Web |

---

## 3. Technical Architecture

### Directory Structure

```
freecode/
  locales/
    active.en.toml           # English source (always present)
    active.es.toml           # Spanish
    active.fr.toml           # French
    ...
    translate.*.toml         # Translation work files
  internal/
    i18n/
      loader.go              # Translation loader
      locale.go              # Locale detection
      rtl.go                 # RTL/Bidi support
  cmd/
    freecode/
      main.go                # Wire i18n
```

### Locale File Format (TOML)

```toml
# locales/active.en.toml
[id = "greeting"]
other = "Hello, {name}!"

[id = "file_saved"]
other = "File saved successfully."

[id = "items_count"]
one = "{n} item"
other = "{n} items"
```

### Go i18n Integration

```go
// internal/i18n/loader.go
package i18n

import (
    "embed"
    "golang.org/x/text/language"
    "github.com/nicksnyder/go-i18n/v2/i18n"
)

// LocaleFS embeds all locale files
//go:embed locales/*.toml
var LocaleFS embed.FS

type Loader struct {
    bundle *i18n.Bundle
    localizer *i18n.Localizer
}

func NewLoader(defaultLang string) (*Loader, error) {
    bundle := i18n.NewBundle(language.English)
    bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

    // Load all locale files from embedded FS
    entries, _ := LocaleFS.ReadDir("locales")
    for _, e := range entries {
        if !e.IsDir() && strings.HasSuffix(e.Name(), ".toml") {
            data, _ := LocaleFS.ReadFile("locales/" + e.Name())
            bundle.Parse([]byte(data))
        }
    }

    localizer := i18n.NewLocalizer(bundle, defaultLang)
    return &Loader{bundle, localizer}, nil
}

func (l *Loader) T(templateID string, args ...interface{}) string {
    return l.localizer.MustLocalize(&i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{ID: templateID},
        TemplateData: args,
    })
}
```

### RTL Support

```go
// internal/i18n/rtl.go
package i18n

import (
    "golang.org/x/text/unicode/bidi"
)

var rtlLanguages = map[string]bool{
    "ar": true,  // Arabic
    "he": true,  // Hebrew
    "ur": true,  // Urdu
}

// IsRTL returns true if the language is right-to-left
func IsRTL(tag string) bool {
    return rtlLanguages[tag]
}

// ReorderRTL reorders a string for RTL display
func ReorderRTL(s string) string {
    // Use bidi algorithm to compute visual order
    p := bidi.NewParagraph(nil)
    p.SetString(s, bidi.LeftToRight)
    return p.String()
}
```

### Locale Detection

```go
// internal/i18n/locale.go

// DetectOrder returns user's preferred languages
func DetectOrder() []language.Tag {
    // 1. Check FREECODE_LANG env var
    if lang := os.Getenv("FREECODE_LANG"); lang != "" {
        tag, _ := language.Parse(lang)
        return []language.Tag{tag}
    }

    // 2. Check LANG/LC_ALL env vars
    if lang := os.Getenv("LANG"); lang != "" {
        tag, _ := language.Parse(lang)
        return []language.Tag{tag}
    }

    // 3. Fall back to English
    return []language.Tag{language.English}
}
```

---

## 4. Configuration

### Config Schema

```yaml
i18n:
  # Explicit language override
  language: "en"

  # RTL mode: auto, force-ltr, force-rtl
  rtl: "auto"

  # Fallback languages (in order)
  fallback:
    - "en"
    - "es"
```

---

## 5. README Translations

Structure for translated README files:

```
README.md              # English (default)
README.es.md           # Spanish
README.fr.md           # French
README.zh.md           # Chinese
README.ja.md           # Japanese
...
```

Each translation file:
- Same content, translated
- Language name in native script in header
- Link back to English version

---

## 6. TTY/RTL Considerations

### Terminal Support (2026)

| Terminal | RTL Status |
|----------|------------|
| GNOME Terminal | ✅ Full support |
| Konsole | ✅ |
| WezTerm | ⚠️ Experimental |
| tmux/screen | ❌ Not supported |
| byobu | ❌ (inherits tmux) |

### Detection Strategy

1. Check `TERM_PROGRAM` for known terminals
2. Check terminfo for `Bidi` capability
3. Check `FREECODE_RTL` env var (explicit override)
4. Fall back to config setting

### Bubble Tea TTY Considerations

For Bubble Tea TUI:
- Use `lipgloss` for styling (supports RTL via Unicode)
- Detect RTL and apply `lipgloss.RTL` style when needed
- For terminals without RTL: show warning, use LTR fallback

---

## 7. Implementation Tasks

| Task | Description | Status |
|------|-------------|--------|
| 7.1 | Create `internal/i18n/` package structure | ⏳ |
| 7.2 | Implement `Loader` with embed.FS for locale files | ⏳ |
| 7.3 | Add `golang.org/x/text` dependencies to go.mod | ⏳ |
| 7.4 | Create `active.en.toml` with base strings (~100 strings) | ⏳ |
| 7.5 | Implement locale detection (env vars, config) | ⏳ |
| 7.6 | Add RTL support (`IsRTL`, `ReorderRTL`) | ⏳ |
| 7.7 | Create translated README files (ES, FR, ZH, JA) | ⏳ |
| 7.8 | Add `i18n` config section to schema | ⏳ |
| 7.9 | Create `goi18n` Makefile targets (extract, merge) | ⏳ |
| 7.10 | Add fictional language support (Klingon, Elvish, etc.) | ⏳ |
| 7.11 | Add i18n tests | ⏳ |
| 7.12 | Document translation workflow in CONTRIBUTING.md | ⏳ |

---

## 8. Translation Workflow

### Adding New Strings

```bash
# 1. Mark strings in code with message IDs
// 2. Run extract to update active.en.toml
make i18n-extract

# 3. Run merge to create translate.*.toml files
make i18n-merge

# 4. Translate each translate.*.toml file
# (manual or via Weblate/Transifex)

# 5. Activate translations
make i18n-activate
```

### Continuous Translation

- Use [Weblate](https://weblate.org/) or [Transifex](https://transifex.com/) for collaborative translation
- CI integration to detect new strings and notify translators
- Community PRs for translation updates

---

## 9. Dependencies

```go
// go.mod additions
require (
    github.com/nicksnyder/go-i18n/v2 v2.1.3
    golang.org/x/text v0.27.0
    github.com/BurntSushi/toml v1.4.0
)
```

---

## 10. Status

**Current phase:** Planning

**Next action:** Implement `internal/i18n/` package structure and create base English locale file.
