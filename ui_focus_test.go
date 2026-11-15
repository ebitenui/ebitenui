package ebitenui

import (
	"image"
	"image/color"
	"testing"

	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/matryer/is"
)

// focusButton returns a button placed at rect.
func focusButton(rect image.Rectangle, opts ...widget.ButtonOpt) *widget.Button {
	img := eimage.NewNineSliceColor(color.White)
	b := widget.NewButton(append(opts, widget.ButtonOpts.Image(&widget.ButtonImage{Idle: img, Pressed: img}))...)
	b.GetWidget().Rect = rect
	return b
}

// focusUI returns a UI holding widgets where they are placed; its container
// has no layout, so they stay there.
func focusUI(widgets ...widget.PreferredSizeLocateableWidget) *UI {
	ui := &UI{Container: widget.NewContainer()}
	ui.Container.AddChild(widgets...)
	return ui
}

// focusGrid lays out three rows of three buttons, 40 by 20, 10 apart, and
// returns them by row and column.
func focusGrid() (*UI, [3][3]*widget.Button) {
	var g [3][3]*widget.Button
	all := make([]widget.PreferredSizeLocateableWidget, 0, 9)
	for r := range 3 {
		for c := range 3 {
			g[r][c] = focusButton(image.Rect(c*50, r*30, c*50+40, r*30+20))
			all = append(all, g[r][c])
		}
	}
	return focusUI(all...), g
}

// Shift-Tab from the first widget wraps to the last and keeps track of it.
func TestUI_ChangeFocusPreviousWraps(t *testing.T) {
	is := is.New(t)

	ui, g := focusGrid()
	ui.SetFocusedWidget(g[0][0])
	ui.ChangeFocus(widget.FOCUS_PREVIOUS)
	is.Equal(ui.GetFocusedWidget(), g[2][2])
	is.True(!g[0][0].IsFocused())
}
