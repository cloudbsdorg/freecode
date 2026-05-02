package i18n

import "testing"

func TestIsRTL(t *testing.T) {
	tests := []struct {
		tag      string
		expected bool
	}{
		{"ar", true},
		{"AR", true},
		{"he", true},
		{"HE", true},
		{"ur", true},
		{"en", false},
		{"es", false},
		{"zh", false},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			if got := IsRTLLang(tt.tag); got != tt.expected {
				t.Errorf("IsRTLLang(%q) = %v, want %v", tt.tag, got, tt.expected)
			}
		})
	}
}

func TestDetectLanguage(t *testing.T) {
	t.Setenv("FREECODE_LANG", "")
	t.Setenv("LANG", "")
	t.Setenv("LC_ALL", "")

	if got := DetectLanguage(); got != "en" {
		t.Errorf("DetectLanguage() = %v, want en", got)
	}
}

func TestDetectLanguageEnvOverride(t *testing.T) {
	t.Setenv("FREECODE_LANG", "es")
	t.Setenv("LANG", "fr_FR.UTF-8")
	t.Setenv("LC_ALL", "de_DE.UTF-8")

	if got := DetectLanguage(); got != "es" {
		t.Errorf("DetectLanguage() = %v, want es (FREECODE_LANG should take precedence)", got)
	}
}

func TestDetectLanguageLangNormalization(t *testing.T) {
	t.Setenv("FREECODE_LANG", "")
	t.Setenv("LC_ALL", "")

	t.Setenv("LANG", "zh_CN.UTF-8")
	if got := DetectLanguage(); got != "zh-CN" {
		t.Errorf("DetectLanguage() = %v, want zh-CN", got)
	}

	t.Setenv("LANG", "pt_BR.UTF-8")
	if got := DetectLanguage(); got != "pt-BR" {
		t.Errorf("DetectLanguage() = %v, want pt-BR", got)
	}
}
