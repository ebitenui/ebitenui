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

// A window added before the UI's first update takes the UI's theme: its
// contents are validated as it is added.
func TestUI_WindowAddedBeforeTheFirstUpdate(t *testing.T) {
	is := is.New(t)
	t.Cleanup(func() { widget.SetUITheme(nil) })

	theme := textTheme(t, color.White)
	ui := &UI{Container: widget.NewContainer(), PrimaryTheme: theme}
	w, label := newTextWindow()
	ui.AddWindow(w)

	is.Equal(label.GetWidget().GetTheme(), theme)
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

// A window with a transition stays in the UI on its way out and leaves at
// the end; added again on its way out, it comes back instead.
func TestUI_WindowTransition(t *testing.T) {
	is := is.New(t)

	ui := newThemedUI(t, textTheme(t, color.White))
	w := widget.NewWindow(
		widget.WindowOpts.Contents(widget.NewContainer()),
		widget.WindowOpts.Transition(widget.WindowTransition{Duration: 4, In: widget.WindowMotion{Offset: image.Pt(0, 40)}}),
	)
	w.SetLocation(image.Rect(0, 0, 100, 50))
	update := func(n int) {
		for range n {
			ui.Update()
		}
	}

	closeWindow := ui.AddWindow(w)
	update(4)
	closeWindow()
	update(2)
	is.True(ui.IsWindowOpen(w)) // on its way out

	ui.AddWindow(w)
	update(4)
	is.True(ui.IsWindowOpen(w)) // came back
	is.True(!w.Closing())

	w.Close()
	update(4)
	is.True(!ui.IsWindowOpen(w)) // out
}
