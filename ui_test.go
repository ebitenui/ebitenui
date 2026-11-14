package ebitenui

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/matryer/is"
	"golang.org/x/image/font/gofont/goregular"
)

func textTheme(t *testing.T, c color.Color) *widget.Theme {
	t.Helper()
	src, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatal(err)
	}
	var face text.Face = &text.GoTextFace{Source: src, Size: 16}
	return &widget.Theme{TextTheme: &widget.TextParams{Face: &face, Color: c}}
}

func newThemedUI(t *testing.T, theme *widget.Theme) *UI {
	t.Helper()
	t.Cleanup(func() { widget.SetUITheme(nil) })
	ui := &UI{Container: widget.NewContainer(), PrimaryTheme: theme}
	ui.Update()
	return ui
}

func newTextWindow() (*widget.Window, *widget.Text) {
	label := widget.NewText(widget.TextOpts.TextLabel("Hello")) // no face or color of its own
	contents := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout()))
	contents.AddChild(label)
	return widget.NewWindow(widget.WindowOpts.Contents(contents)), label
}

// A window's contents take the UI's theme even before the window is added,
// so they can be measured to place it (#307).
func TestUI_WindowContentsTakeTheThemeBeforeBeingAdded(t *testing.T) {
	is := is.New(t)

	newThemedUI(t, textTheme(t, color.White))
	w, _ := newTextWindow()

	x, y := w.Contents.PreferredSize()
	is.True(x > 0 && y > 0) // measured with the theme's face
}

// An open window follows the UI's theme when it changes (#307).
func TestUI_WindowsFollowThemeChanges(t *testing.T) {
	is := is.New(t)

	light, dark := textTheme(t, color.White), textTheme(t, color.Black)
	ui := newThemedUI(t, light)
	w, label := newTextWindow()
	w.SetLocation(image.Rect(0, 0, 100, 50))
	ui.AddWindow(w)

	ui.PrimaryTheme = dark
	ui.Update()

	is.Equal(w.GetContainer().GetWidget().GetTheme(), dark) // the window's
	is.Equal(label.GetWidget().GetTheme(), dark)            // its contents'
}
