package themes

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

// Each color mode gets one theme, the same every time: the dark one for a
// dark system, the light one for a light system or one that doesn't say.
func TestBasicThemeFor(t *testing.T) {
	is := is.New(t)

	dark, light := basicThemeFor(ebiten.ColorModeDark), basicThemeFor(ebiten.ColorModeLight)
	is.True(dark != light)
	is.Equal(basicThemeFor(ebiten.ColorModeDark), dark)     // made once
	is.Equal(basicThemeFor(ebiten.ColorModeLight), light)   // made once
	is.Equal(basicThemeFor(ebiten.ColorModeUnknown), light) // light when the system doesn't say
	is.Equal(GetBasicSystemTheme(), basicThemeFor(ebiten.SystemColorMode()))
}
