package ui

import (
	"testing"
)

func TestDefaultTheme(t *testing.T) {
	theme := DefaultTheme()

	if theme.ActiveTab.String() == "" {
		t.Error("DefaultTheme() should return non-empty ActiveTab")
	}
}

func TestDarkTheme(t *testing.T) {
	theme := DarkTheme()

	if theme.ActiveTab.String() == "" {
		t.Error("DarkTheme() should return non-empty ActiveTab")
	}
}

func TestLightTheme(t *testing.T) {
	theme := LightTheme()

	if theme.ActiveTab.String() == "" {
		t.Error("LightTheme() should return non-empty ActiveTab")
	}
}
