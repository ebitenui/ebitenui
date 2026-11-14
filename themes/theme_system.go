package themes

import (
	"sync"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

// The basic themes, each made once, the first time it is asked for.
var (
	basicDarkTheme  = sync.OnceValue(GetBasicDarkTheme)
	basicLightTheme = sync.OnceValue(GetBasicLightTheme)
)

// GetBasicSystemTheme returns the basic theme matching the system's color
// mode: the dark one when the system is dark, the light one when it is
// light or doesn't say. It returns the same theme for the same mode, so
// setting a UI's PrimaryTheme to it every frame makes the UI follow the
// system and re-validate only when the mode changes.
func GetBasicSystemTheme() *widget.Theme {
	return basicThemeFor(ebiten.SystemColorMode())
}

// basicThemeFor returns the basic theme for a color mode.
func basicThemeFor(mode ebiten.ColorMode) *widget.Theme {
	if mode == ebiten.ColorModeDark {
		return basicDarkTheme()
	}
	return basicLightTheme()
}
