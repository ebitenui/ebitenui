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

// focusSlider returns a horizontal slider placed at rect, from 0 to 100 at
// 50, moving 10 a page.
func focusSlider(rect image.Rectangle) *widget.Slider {
	img := eimage.NewNineSliceColor(color.White)
	s := widget.NewSlider(
		widget.SliderOpts.Images(&widget.SliderTrackImage{Idle: img}, &widget.ButtonImage{Idle: img, Pressed: img}),
		widget.SliderOpts.MinMax(0, 100), widget.SliderOpts.InitialCurrent(50),
		widget.SliderOpts.PageSizeFunc(func() int { return 10 }))
	s.GetWidget().Rect = rect
	return s
}

// navigationUI returns a UI with GamepadFocus on over widgets, drawn once so
// they are laid out, and devices with one gamepad. press holds a key or a
// gamepad button for two updates, then lets go for one.
func navigationUI(t *testing.T, widgets ...widget.PreferredSizeLocateableWidget) (ui *UI, d *inputtest.Devices, press func(any)) {
	t.Helper()
	ui = focusUI(widgets...)
	ui.GamepadFocus = true
	ui.Draw(ebiten.NewImage(200, 200))
	d = inputtest.Use(t)
	d.Pads = make([]inputtest.Pad, 1)
	press = func(k any) {
		for _, down := range []bool{true, true, false} {
			switch k := k.(type) {
			case ebiten.Key:
				d.Keys[k] = down
			case ebiten.StandardGamepadButton:
				d.Pads[0].Buttons[k] = down
			}
			ui.Update()
			event.ExecuteDeferred()
		}
	}
	return ui, d, press
}

// The focused widget takes the actions it acts on, and the UI acts on those
// left: a gamepad's Right moves a focused slider, and its Down, which the
// slider leaves, moves the focus. The arrow keys move the slider too, but
// never the focus.
func TestUI_FocusedWidgetTakesActions(t *testing.T) {
	is := is.New(t)

	s := focusSlider(image.Rect(0, 0, 100, 20))
	right := focusButton(image.Rect(120, 0, 160, 20))
	below := focusButton(image.Rect(0, 40, 100, 60))
	ui, _, press := navigationUI(t, s, right, below)
	ui.SetFocusedWidget(s)

	press(ebiten.StandardGamepadButtonLeftRight)
	is.Equal(s.Current, 60)             // the slider took Right
	is.True(ui.GetFocusedWidget() == s) // so the focus stays

	press(ebiten.KeyArrowRight)
	is.Equal(s.Current, 70) // and the arrow key's
	press(ebiten.KeyArrowDown)
	is.True(ui.GetFocusedWidget() == s) // a key's Down left: the focus stays

	press(ebiten.StandardGamepadButtonLeftBottom)
	is.True(ui.GetFocusedWidget() == below) // a gamepad's: the focus moves
}
