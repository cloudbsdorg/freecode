package syntax

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestNewStyle(t *testing.T) {
	s := NewStyle(nil)
	_ = s.Keyword
}

func TestNewStyleWithConfig(t *testing.T) {
	s := NewStyle(func(st *Style) {
		st.Keyword = lipgloss.NewStyle().Bold(true)
	})
	_ = s.Keyword
}

func TestDefaultTheme(t *testing.T) {
	theme := DefaultTheme()
	if theme.Name != "default" {
		t.Errorf("expected name=default, got %s", theme.Name)
	}
	if !theme.initialized {
		t.Error("expected initialized=true")
	}
}

func TestDarkPlusTheme(t *testing.T) {
	theme := DarkPlusTheme()
	if theme.Name != "default" {
		t.Errorf("expected name=default, got %s", theme.Name)
	}
}