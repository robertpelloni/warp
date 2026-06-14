package theme

import (
	"testing"
)

func TestDefaultTheme(t *testing.T) {
	th := DefaultTheme()
	if th == nil {
		t.Fatal("DefaultTheme() returned nil")
	}
	if th.Name() != "Standard" {
		t.Errorf("Name() = %q, want 'Standard'", th.Name())
	}
	if !th.IsDark() {
		t.Error("DefaultTheme should be dark")
	}
}

func TestGetTheme(t *testing.T) {
	th := GetTheme("Dracula")
	if th.Name() != "Dracula" {
		t.Errorf("GetTheme('Dracula').Name() = %q", th.Name())
	}

	// Unknown theme should return default
	th = GetTheme("NonExistent")
	if th.Name() != "Standard" {
		t.Errorf("GetTheme('NonExistent').Name() = %q, want 'Standard'", th.Name())
	}
}

func TestThemeNames(t *testing.T) {
	names := ThemeNames()
	if len(names) < 6 {
		t.Errorf("ThemeNames() returned %d names, expected >= 6", len(names))
	}
}

func TestAllThemesHaveColors(t *testing.T) {
	requiredColors := []string{
		ColorBackground, ColorSurface, ColorPrimary, ColorSecondary,
		ColorText, ColorTextDim, ColorAccent, ColorSuccess,
		ColorWarning, ColorError, ColorPrompt, ColorBorder,
		ColorInputBg, ColorToolbar,
	}

	for _, name := range ThemeNames() {
		th := GetTheme(name)
		for _, colorName := range requiredColors {
			c := th.Color(colorName)
			if c == nil {
				t.Errorf("Theme %q missing color %q", name, colorName)
			}
		}
	}
}

func TestCOLORREF(t *testing.T) {
	th := DefaultTheme()
	bg := th.COLORREF(ColorBackground)
	if bg == 0 {
		t.Error("COLORREF(Background) should not be 0")
	}
}

func TestAllThemesExist(t *testing.T) {
	expectedThemes := []string{"Standard", "Dracula", "Monokai", "Nord", "One Dark", "Catppuccin"}
	for _, name := range expectedThemes {
		th := GetTheme(name)
		if th.Name() != name {
			t.Errorf("Theme %q not found", name)
		}
	}
}
