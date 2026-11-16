package ebitenui

import (
	"image"
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/internal/input/inputtest"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
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

// Without links of its own, a widget passes the focus to the nearest one in
// the direction asked, diagonals included; at an edge it keeps it.
func TestUI_ChangeFocusToTheNearestInADirection(t *testing.T) {
	is := is.New(t)

	ui, g := focusGrid()
	for _, step := range []struct {
		from *widget.Button
		dir  widget.FocusDirection
		want *widget.Button
	}{
		{g[1][1], widget.FOCUS_EAST, g[1][2]},
		{g[1][0], widget.FOCUS_EAST, g[1][1]}, // the nearer of two in the row
		{g[1][1], widget.FOCUS_WEST, g[1][0]},
		{g[1][1], widget.FOCUS_NORTH, g[0][1]},
		{g[1][1], widget.FOCUS_SOUTH, g[2][1]},
		{g[1][1], widget.FOCUS_NORTHEAST, g[0][2]},
		{g[1][1], widget.FOCUS_SOUTHWEST, g[2][0]},
		{g[1][2], widget.FOCUS_EAST, g[1][2]}, // nothing further east
		{g[0][0], widget.FOCUS_NORTHWEST, g[0][0]},
	} {
		ui.SetFocusedWidget(step.from)
		ui.ChangeFocus(step.dir)
		is.Equal(ui.GetFocusedWidget(), step.want)
		is.True(step.want.IsFocused())
	}
}

// A widget in the same row, or column, wins over a nearer one beside it.
func TestUI_ChangeFocusKeepsToTheRow(t *testing.T) {
	is := is.New(t)

	from := focusButton(image.Rect(0, 0, 40, 20))
	far := focusButton(image.Rect(300, 0, 340, 20))
	beside := focusButton(image.Rect(50, 30, 90, 50)) // nearer, one row down
	ui := focusUI(from, far, beside)

	ui.SetFocusedWidget(from)
	ui.ChangeFocus(widget.FOCUS_EAST)
	is.Equal(ui.GetFocusedWidget(), far)

	bottom := focusButton(image.Rect(0, 300, 40, 320))
	above := focusButton(image.Rect(0, 0, 40, 20))
	aside := focusButton(image.Rect(50, 250, 90, 270)) // nearer, one column across
	ui = focusUI(bottom, above, aside)
	ui.SetFocusedWidget(bottom)
	ui.ChangeFocus(widget.FOCUS_NORTH)
	is.Equal(ui.GetFocusedWidget(), above)
}

// A link set for a direction wins; a disabled widget is passed over.
func TestUI_ChangeFocusFollowsLinksAndSkipsDisabled(t *testing.T) {
	is := is.New(t)

	ui, g := focusGrid()
	g[1][1].AddFocus(widget.FOCUS_EAST, g[2][0])
	ui.SetFocusedWidget(g[1][1])
	ui.ChangeFocus(widget.FOCUS_EAST)
	is.Equal(ui.GetFocusedWidget(), g[2][0]) // the link

	g[1][0].GetWidget().Disabled = true
	g[1][1].AddFocus(widget.FOCUS_WEST, g[1][0])
	ui.SetFocusedWidget(g[1][1])
	ui.ChangeFocus(widget.FOCUS_WEST)
	is.True(ui.GetFocusedWidget() != g[1][0]) // a link to a disabled widget gives way to the search

	g[0][1].GetWidget().Disabled = true
	ui.SetFocusedWidget(g[1][1])
	ui.ChangeFocus(widget.FOCUS_NORTH)
	is.True(ui.GetFocusedWidget() != g[0][1])                                     // not the disabled one
	is.True(ui.GetFocusedWidget() == g[0][0] || ui.GetFocusedWidget() == g[0][2]) // the row above
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

// ActivateFocused clicks the focused button, as Enter does; it does nothing
// without a focused widget or for a disabled one.
func TestUI_ActivateFocused(t *testing.T) {
	is := is.New(t)

	clicks := 0
	b := focusButton(image.Rect(0, 0, 40, 20), widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) { clicks++ }))
	ui := focusUI(b)

	is.True(!ui.ActivateFocused()) // nothing focused
	ui.SetFocusedWidget(b)
	is.True(ui.ActivateFocused())
	event.ExecuteDeferred()
	is.Equal(clicks, 1)

	b.GetWidget().Disabled = true
	is.True(!ui.ActivateFocused())
	event.ExecuteDeferred()
	is.Equal(clicks, 1)
}

// With GamepadFocus on, the gamepads' D-pad moves the focus and their bottom
// face button activates the focused widget; off, gamepads do nothing. Two
// gamepads pressing together act once.
func TestUI_GamepadFocus(t *testing.T) {
	is := is.New(t)

	clicks := 0
	left := focusButton(image.Rect(0, 0, 40, 20))
	middle := focusButton(image.Rect(50, 0, 90, 20), widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) { clicks++ }))
	right := focusButton(image.Rect(100, 0, 140, 20))
	ui := focusUI(left, middle, right)
	ui.SetFocusedWidget(left)
	d := inputtest.Use(t)
	d.Pads = make([]inputtest.Pad, 2)
	// tap holds b on the pads given for two updates, then lets go.
	tap := func(b ebiten.StandardGamepadButton, pads ...int) {
		for _, down := range []bool{true, true, false} {
			for _, p := range pads {
				d.Pads[p].Buttons[b] = down
			}
			ui.Update()
			event.ExecuteDeferred()
		}
	}

	tap(ebiten.StandardGamepadButtonLeftRight, 0)
	is.True(ui.GetFocusedWidget() == left) // off by default

	ui.GamepadFocus = true
	tap(ebiten.StandardGamepadButtonLeftRight, 0, 1)
	is.True(ui.GetFocusedWidget() == middle) // both pads at once: one step

	tap(ebiten.StandardGamepadButtonRightBottom, 1)
	is.Equal(clicks, 1) // the bottom face button, on either pad
}
