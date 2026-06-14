package theme

import (
	"image/color"
)

// WarpTheme defines the color palette for Warp's UI.
type WarpTheme struct {
	name   string
	dark   bool
	colors map[string]color.Color
}

// Built-in theme names.
const (
	ThemeStandard   = "Standard"
	ThemeDracula    = "Dracula"
	ThemeMonokai    = "Monokai"
	ThemeSolarized  = "Solarized"
	ThemeNord       = "Nord"
	ThemeCatppuccin = "Catppuccin"
	ThemeOneDark    = "One Dark"
)

// Warp color keys.
const (
	ColorBackground   = "background"
	ColorSurface      = "surface"
	ColorPrimary      = "primary"
	ColorSecondary    = "secondary"
	ColorText         = "text"
	ColorTextDim      = "text_dim"
	ColorAccent       = "accent"
	ColorSuccess      = "success"
	ColorWarning      = "warning"
	ColorError        = "error"
	ColorBlockCommand = "block_command"
	ColorBlockOutput  = "block_output"
	ColorBlockAI      = "block_ai"
	ColorBlockError   = "block_error"
	ColorPrompt       = "prompt"
	ColorCursor       = "cursor"
	ColorSelection    = "selection"
	ColorBorder       = "border"
	ColorTabActive    = "tab_active"
	ColorTabInactive  = "tab_inactive"
	ColorInputBg      = "input_bg"
	ColorToolbar      = "toolbar"
)

// Themes maps theme names to their definitions.
var Themes = map[string]*WarpTheme{
	ThemeStandard:   standardTheme(),
	ThemeDracula:    draculaTheme(),
	ThemeMonokai:    monokaiTheme(),
	ThemeNord:       nordTheme(),
	ThemeOneDark:    oneDarkTheme(),
	ThemeCatppuccin: catppuccinTheme(),
}

// DefaultTheme returns the default theme.
func DefaultTheme() *WarpTheme {
	return Themes[ThemeStandard]
}

// GetTheme returns a theme by name.
func GetTheme(name string) *WarpTheme {
	if t, ok := Themes[name]; ok {
		return t
	}
	return DefaultTheme()
}

// ThemeNames returns all available theme names.
func ThemeNames() []string {
	return []string{ThemeStandard, ThemeDracula, ThemeMonokai, ThemeSolarized, ThemeNord, ThemeCatppuccin, ThemeOneDark}
}

// Color returns a named color from the theme.
func (t *WarpTheme) Color(name string) color.Color {
	if c, ok := t.colors[name]; ok {
		return c
	}
	return color.Black
}

// Name returns the theme's name.
func (t *WarpTheme) Name() string {
	return t.name
}

// IsDark returns whether this is a dark theme.
func (t *WarpTheme) IsDark() bool {
	return t.dark
}

// COLORREF returns a Win32 COLORREF for the named color.
func (t *WarpTheme) COLORREF(name string) uint32 {
	c := t.Color(name)
	return colorToCOLORREF(c)
}

// colorToCOLORREF converts a color.Color to a Win32 COLORREF.
func colorToCOLORREF(c color.Color) uint32 {
	r, g, b, _ := c.RGBA()
	return uint32(uint8(r>>8)) | uint32(uint8(g>>8))<<8 | uint32(uint8(b>>8))<<16
}

func makeTheme(name string, dark bool, colors map[string]color.Color) *WarpTheme {
	return &WarpTheme{name: name, dark: dark, colors: colors}
}

// ─── Theme Definitions ─────────────────────────────────────────

func standardTheme() *WarpTheme {
	return makeTheme(ThemeStandard, true, map[string]color.Color{
		ColorBackground:   color.NRGBA{24, 26, 32, 255},
		ColorSurface:      color.NRGBA{38, 40, 48, 255},
		ColorPrimary:      color.NRGBA{0, 153, 255, 255},
		ColorSecondary:    color.NRGBA{100, 110, 130, 255},
		ColorText:         color.NRGBA{212, 214, 218, 255},
		ColorTextDim:      color.NRGBA{120, 126, 138, 255},
		ColorAccent:       color.NRGBA{0, 200, 255, 255},
		ColorSuccess:      color.NRGBA{80, 200, 120, 255},
		ColorWarning:      color.NRGBA{255, 199, 0, 255},
		ColorError:        color.NRGBA{255, 85, 85, 255},
		ColorBlockCommand: color.NRGBA{0, 153, 255, 255},
		ColorBlockOutput:  color.NRGBA{150, 150, 150, 255},
		ColorBlockAI:      color.NRGBA{160, 80, 255, 255},
		ColorBlockError:   color.NRGBA{255, 85, 85, 255},
		ColorPrompt:       color.NRGBA{0, 200, 255, 255},
		ColorCursor:       color.NRGBA{0, 200, 255, 255},
		ColorSelection:    color.NRGBA{0, 153, 255, 255},
		ColorBorder:       color.NRGBA{60, 64, 74, 255},
		ColorTabActive:    color.NRGBA{0, 153, 255, 255},
		ColorTabInactive:  color.NRGBA{100, 110, 130, 255},
		ColorInputBg:      color.NRGBA{30, 32, 38, 255},
		ColorToolbar:      color.NRGBA{32, 34, 40, 255},
	})
}

func draculaTheme() *WarpTheme {
	return makeTheme(ThemeDracula, true, map[string]color.Color{
		ColorBackground:   color.NRGBA{40, 42, 54, 255},
		ColorSurface:      color.NRGBA{50, 52, 64, 255},
		ColorPrimary:      color.NRGBA{189, 147, 249, 255},
		ColorSecondary:    color.NRGBA{98, 114, 164, 255},
		ColorText:         color.NRGBA{248, 248, 242, 255},
		ColorTextDim:      color.NRGBA{98, 114, 164, 255},
		ColorAccent:       color.NRGBA{255, 121, 198, 255},
		ColorSuccess:      color.NRGBA{80, 250, 123, 255},
		ColorWarning:      color.NRGBA{241, 250, 140, 255},
		ColorError:        color.NRGBA{255, 85, 85, 255},
		ColorBlockCommand: color.NRGBA{189, 147, 249, 255},
		ColorBlockOutput:  color.NRGBA{150, 150, 150, 255},
		ColorBlockAI:      color.NRGBA{255, 121, 198, 255},
		ColorBlockError:   color.NRGBA{255, 85, 85, 255},
		ColorPrompt:       color.NRGBA{80, 250, 123, 255},
		ColorCursor:       color.NRGBA{248, 248, 242, 255},
		ColorSelection:    color.NRGBA{189, 147, 249, 255},
		ColorBorder:       color.NRGBA{68, 71, 90, 255},
		ColorTabActive:    color.NRGBA{189, 147, 249, 255},
		ColorTabInactive:  color.NRGBA{98, 114, 164, 255},
		ColorInputBg:      color.NRGBA{34, 36, 46, 255},
		ColorToolbar:      color.NRGBA{44, 46, 58, 255},
	})
}

func monokaiTheme() *WarpTheme {
	return makeTheme(ThemeMonokai, true, map[string]color.Color{
		ColorBackground:   color.NRGBA{39, 40, 34, 255},
		ColorSurface:      color.NRGBA{49, 50, 44, 255},
		ColorPrimary:      color.NRGBA{102, 217, 239, 255},
		ColorSecondary:    color.NRGBA{117, 113, 94, 255},
		ColorText:         color.NRGBA{248, 248, 242, 255},
		ColorTextDim:      color.NRGBA{117, 113, 94, 255},
		ColorAccent:       color.NRGBA{249, 38, 114, 255},
		ColorSuccess:      color.NRGBA{166, 226, 46, 255},
		ColorWarning:      color.NRGBA{230, 219, 116, 255},
		ColorError:        color.NRGBA{249, 38, 114, 255},
		ColorBlockCommand: color.NRGBA{102, 217, 239, 255},
		ColorBlockOutput:  color.NRGBA{150, 150, 150, 255},
		ColorBlockAI:      color.NRGBA{249, 38, 114, 255},
		ColorBlockError:   color.NRGBA{249, 38, 114, 255},
		ColorPrompt:       color.NRGBA{166, 226, 46, 255},
		ColorCursor:       color.NRGBA{248, 248, 242, 255},
		ColorSelection:    color.NRGBA{102, 217, 239, 255},
		ColorBorder:       color.NRGBA{73, 72, 62, 255},
		ColorTabActive:    color.NRGBA{102, 217, 239, 255},
		ColorTabInactive:  color.NRGBA{117, 113, 94, 255},
		ColorInputBg:      color.NRGBA{33, 34, 28, 255},
		ColorToolbar:      color.NRGBA{43, 44, 38, 255},
	})
}

func nordTheme() *WarpTheme {
	return makeTheme(ThemeNord, true, map[string]color.Color{
		ColorBackground:   color.NRGBA{46, 52, 64, 255},
		ColorSurface:      color.NRGBA{59, 66, 82, 255},
		ColorPrimary:      color.NRGBA{136, 192, 208, 255},
		ColorSecondary:    color.NRGBA{76, 86, 106, 255},
		ColorText:         color.NRGBA{216, 222, 233, 255},
		ColorTextDim:      color.NRGBA{76, 86, 106, 255},
		ColorAccent:       color.NRGBA{180, 142, 173, 255},
		ColorSuccess:      color.NRGBA{163, 190, 140, 255},
		ColorWarning:      color.NRGBA{235, 203, 139, 255},
		ColorError:        color.NRGBA{191, 97, 106, 255},
		ColorBlockCommand: color.NRGBA{136, 192, 208, 255},
		ColorBlockOutput:  color.NRGBA{150, 150, 150, 255},
		ColorBlockAI:      color.NRGBA{180, 142, 173, 255},
		ColorBlockError:   color.NRGBA{191, 97, 106, 255},
		ColorPrompt:       color.NRGBA{163, 190, 140, 255},
		ColorCursor:       color.NRGBA{216, 222, 233, 255},
		ColorSelection:    color.NRGBA{136, 192, 208, 255},
		ColorBorder:       color.NRGBA{67, 76, 94, 255},
		ColorTabActive:    color.NRGBA{136, 192, 208, 255},
		ColorTabInactive:  color.NRGBA{76, 86, 106, 255},
		ColorInputBg:      color.NRGBA{40, 45, 56, 255},
		ColorToolbar:      color.NRGBA{52, 58, 72, 255},
	})
}

func oneDarkTheme() *WarpTheme {
	return makeTheme(ThemeOneDark, true, map[string]color.Color{
		ColorBackground:   color.NRGBA{40, 44, 52, 255},
		ColorSurface:      color.NRGBA{50, 54, 62, 255},
		ColorPrimary:      color.NRGBA{97, 175, 239, 255},
		ColorSecondary:    color.NRGBA{92, 99, 112, 255},
		ColorText:         color.NRGBA{171, 178, 191, 255},
		ColorTextDim:      color.NRGBA{92, 99, 112, 255},
		ColorAccent:       color.NRGBA{198, 120, 221, 255},
		ColorSuccess:      color.NRGBA{152, 195, 121, 255},
		ColorWarning:      color.NRGBA{229, 192, 123, 255},
		ColorError:        color.NRGBA{224, 108, 117, 255},
		ColorBlockCommand: color.NRGBA{97, 175, 239, 255},
		ColorBlockOutput:  color.NRGBA{150, 150, 150, 255},
		ColorBlockAI:      color.NRGBA{198, 120, 221, 255},
		ColorBlockError:   color.NRGBA{224, 108, 117, 255},
		ColorPrompt:       color.NRGBA{152, 195, 121, 255},
		ColorCursor:       color.NRGBA{171, 178, 191, 255},
		ColorSelection:    color.NRGBA{97, 175, 239, 255},
		ColorBorder:       color.NRGBA{60, 64, 72, 255},
		ColorTabActive:    color.NRGBA{97, 175, 239, 255},
		ColorTabInactive:  color.NRGBA{92, 99, 112, 255},
		ColorInputBg:      color.NRGBA{34, 38, 46, 255},
		ColorToolbar:      color.NRGBA{44, 48, 56, 255},
	})
}

func catppuccinTheme() *WarpTheme {
	return makeTheme(ThemeCatppuccin, true, map[string]color.Color{
		ColorBackground:   color.NRGBA{30, 30, 46, 255},
		ColorSurface:      color.NRGBA{40, 40, 60, 255},
		ColorPrimary:      color.NRGBA{137, 180, 250, 255},
		ColorSecondary:    color.NRGBA{88, 91, 112, 255},
		ColorText:         color.NRGBA{205, 214, 244, 255},
		ColorTextDim:      color.NRGBA{88, 91, 112, 255},
		ColorAccent:       color.NRGBA{203, 166, 247, 255},
		ColorSuccess:      color.NRGBA{166, 227, 161, 255},
		ColorWarning:      color.NRGBA{249, 226, 175, 255},
		ColorError:        color.NRGBA{243, 139, 168, 255},
		ColorBlockCommand: color.NRGBA{137, 180, 250, 255},
		ColorBlockOutput:  color.NRGBA{150, 150, 150, 255},
		ColorBlockAI:      color.NRGBA{203, 166, 247, 255},
		ColorBlockError:   color.NRGBA{243, 139, 168, 255},
		ColorPrompt:       color.NRGBA{166, 227, 161, 255},
		ColorCursor:       color.NRGBA{205, 214, 244, 255},
		ColorSelection:    color.NRGBA{137, 180, 250, 255},
		ColorBorder:       color.NRGBA{54, 58, 79, 255},
		ColorTabActive:    color.NRGBA{137, 180, 250, 255},
		ColorTabInactive:  color.NRGBA{88, 91, 112, 255},
		ColorInputBg:      color.NRGBA{24, 24, 38, 255},
		ColorToolbar:      color.NRGBA{34, 34, 52, 255},
	})
}
