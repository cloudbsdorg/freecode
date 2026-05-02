package i18n

import (
	"testing"
	"testing/fstest"
)

func TestIsRTLLang(t *testing.T) {
	tests := []struct {
		tag      string
		expected bool
	}{
		{"ar", true},
		{"AR", true},
		{"he", true},
		{"HE", true},
		{"ur", true},
		{"UR", true},
		{"fa", true},
		{"FA", true},
		{"yi", true},
		{"YI", true},
		{"en", false},
		{"EN", false},
		{"es", false},
		{"ES", false},
		{"zh", false},
		{"ZH", false},
		{"ja", false},
		{"JA", false},
		{"fr", false},
		{"de", false},
		{"pt", false},
		{"ru", false},
		{"", false},
		{"xyz", false},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			if got := IsRTLLang(tt.tag); got != tt.expected {
				t.Errorf("IsRTLLang(%q) = %v, want %v", tt.tag, got, tt.expected)
			}
		})
	}
}

func TestReorderForRTL(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{"empty", "", ""},
		{"ltr only", "Hello", "Hello"},
		{"ltr with spaces", "Hello World", "Hello World"},
		{"arabic word", "مرحبا", "مرحبا"},
		{"mixed ltr rtl", "Hello مرحبا", "Hello مرحبا"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReorderForRTL(tt.input); got != tt.output {
				t.Errorf("ReorderForRTL(%q) = %q, want %q", tt.input, got, tt.output)
			}
		})
	}
}

func TestDetectLanguageFREECODELANG(t *testing.T) {
	t.Setenv("FREECODE_LANG", "es")
	t.Setenv("LANG", "en_US.UTF-8")
	t.Setenv("LC_ALL", "de_DE.UTF-8")

	if got := DetectLanguage(); got != "es" {
		t.Errorf("DetectLanguage() = %v, want es (FREECODE_LANG takes precedence)", got)
	}
}

func TestDetectLanguageLANG(t *testing.T) {
	t.Setenv("FREECODE_LANG", "")
	t.Setenv("LC_ALL", "")

	t.Setenv("LANG", "fr_FR.UTF-8")
	if got := DetectLanguage(); got != "fr-FR" {
		t.Errorf("DetectLanguage() = %v, want fr-FR", got)
	}

	t.Setenv("LANG", "zh_CN.UTF-8")
	if got := DetectLanguage(); got != "zh-CN" {
		t.Errorf("DetectLanguage() = %v, want zh-CN", got)
	}

	t.Setenv("LANG", "pt_BR.UTF-8")
	if got := DetectLanguage(); got != "pt-BR" {
		t.Errorf("DetectLanguage() = %v, want pt-BR", got)
	}

	t.Setenv("LANG", "ja_JP.UTF-8")
	if got := DetectLanguage(); got != "ja-JP" {
		t.Errorf("DetectLanguage() = %v, want ja-JP", got)
	}

	t.Setenv("LANG", "en_US.UTF-8")
	if got := DetectLanguage(); got != "en-US" {
		t.Errorf("DetectLanguage() = %v, want en-US", got)
	}
}

func TestDetectLanguageLCALL(t *testing.T) {
	t.Setenv("FREECODE_LANG", "")
	t.Setenv("LANG", "")

	t.Setenv("LC_ALL", "es_ES.UTF-8")
	if got := DetectLanguage(); got != "es-ES" {
		t.Errorf("DetectLanguage() = %v, want es-ES", got)
	}

	t.Setenv("LC_ALL", "ru_RU.UTF-8")
	if got := DetectLanguage(); got != "ru-RU" {
		t.Errorf("DetectLanguage() = %v, want ru-RU", got)
	}
}

func TestDetectLanguageDefault(t *testing.T) {
	t.Setenv("FREECODE_LANG", "")
	t.Setenv("LANG", "")
	t.Setenv("LC_ALL", "")

	if got := DetectLanguage(); got != "en" {
		t.Errorf("DetectLanguage() = %v, want en", got)
	}
}

func TestDetectLanguageNoEncoding(t *testing.T) {
	t.Setenv("FREECODE_LANG", "")
	t.Setenv("LC_ALL", "")

	t.Setenv("LANG", "fr_FR")
	if got := DetectLanguage(); got != "fr-FR" {
		t.Errorf("DetectLanguage() = %v, want fr-FR (no encoding suffix)", got)
	}
}

func TestLoaderLoadLocalesError(t *testing.T) {
	badFS := fstest.MapFS{}
	_, err := newLoaderWithFS("en", badFS)
	if err == nil {
		t.Error("expected error with non-existent locales dir")
	}
}
