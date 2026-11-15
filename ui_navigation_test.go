package ebitenui

import (
	"image"
	"image/color"
	"testing"
	"time"

	"github.com/ebitenui/ebitenui/event"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/internal/input/inputtest"
	"github.com/ebitenui/ebitenui/themes"
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
// they are laid out, with its devices (useDevices).
func navigationUI(t *testing.T, widgets ...widget.PreferredSizeLocateableWidget) (ui *UI, d *inputtest.Devices, press func(any)) {
	t.Helper()
	ui = focusUI(widgets...)
	ui.GamepadFocus = true
	d, press = useDevices(t, ui)
	return ui, d, press
}

// useDevices draws ui once, so its widgets are laid out, and returns
// devices with one gamepad. press holds a key or a gamepad button for two
// updates of ui, then lets go for one.
func useDevices(t *testing.T, ui *UI) (d *inputtest.Devices, press func(any)) {
	t.Helper()
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
	return d, press
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

// A direction held repeats, from a key as from the D-pad: once as it goes
// down, then after the delay once every interval, counted in updates at
// Ebitengine's TPS (60 when the TPS follows the display); let go, it stops.
// Unset, the delay is 0.4 s and the interval 0.08 s.
func TestUI_DirectionsRepeatWhileHeld(t *testing.T) {
	is := is.New(t)

	s := focusSlider(image.Rect(0, 0, 100, 20))
	s.Max = 1000
	ui, d, _ := navigationUI(t, s)
	ui.SetFocusedWidget(s)
	// steps holds a direction for some updates, lets go, and returns how
	// many steps the slider took.
	steps := func(hold func(bool), updates int) int {
		start := s.Current
		hold(true)
		for range updates {
			ui.Update()
		}
		hold(false)
		ui.Update()
		return (s.Current - start) / 10
	}
	key := func(down bool) { d.Keys[ebiten.KeyArrowRight] = down }
	pad := func(down bool) { d.Pads[0].Buttons[ebiten.StandardGamepadButtonLeftRight] = down }

	ui.NavigationRepeatDelay = 100 * time.Millisecond   // 6 updates
	ui.NavigationRepeatInterval = 50 * time.Millisecond // 3

	is.Equal(steps(key, 6), 1)  // within the delay
	is.Equal(steps(key, 7), 2)  // past it
	is.Equal(steps(pad, 13), 4) // and twice the interval more

	ui.NavigationRepeatInterval = time.Millisecond // under an update: every update
	is.Equal(steps(key, 9), 4)

	ui.NavigationRepeatDelay, ui.NavigationRepeatInterval = 0, 0 // 24 updates, then 5
	is.Equal(steps(key, 24), 1)
	is.Equal(steps(key, 25), 2)
	is.Equal(steps(key, 30), 3)

	ebiten.SetTPS(ebiten.SyncWithFPS)
	t.Cleanup(func() { ebiten.SetTPS(ebiten.DefaultTPS) })
	is.Equal(steps(key, 24), 1)
}

// With GamepadFocus on, the left stick is a D-pad: it points the way it is
// pushed furthest, once past half its travel, and keeps pointing until it
// falls under a third, so a stick resting near half doesn't press again and
// again.
func TestUI_LeftStickNavigates(t *testing.T) {
	is := is.New(t)

	above := focusButton(image.Rect(0, 0, 100, 20))
	s := focusSlider(image.Rect(0, 40, 100, 60))
	s.Max = 1000
	below := focusButton(image.Rect(0, 80, 100, 100))
	ui, d, _ := navigationUI(t, above, s, below)
	ui.SetFocusedWidget(s)
	// stick holds the stick at x, y for an update and returns the steps
	// the slider took.
	stick := func(x, y float64) int {
		start := s.Current
		d.Pads[0].Axes[ebiten.StandardGamepadAxisLeftStickHorizontal] = x
		d.Pads[0].Axes[ebiten.StandardGamepadAxisLeftStickVertical] = y
		ui.Update()
		return (s.Current - start) / 10
	}

	for i, c := range []struct {
		x     float64
		steps int
	}{
		{0.45, 0}, // under half
		{0.55, 1}, // past it: Right
		{0.4, 0},  // back over a third: held
		{0.55, 0}, // so no new press
		{0.3, 0},  // under a third: let go
		{0.55, 1},
		{0, 0},
		{-0.9, -1}, // Left
		{0, 0},
	} {
		is.Equal([]any{i, stick(c.x, 0)}, []any{i, c.steps})
	}

	is.Equal(stick(0.8, 0.6), 1) // a diagonal: the way it is pushed furthest
	stick(0, 0)
	stick(0.6, 0.8) // Down, which the slider leaves
	is.True(ui.GetFocusedWidget() == below)

	ui.SetFocusedWidget(s)
	stick(0, 0)
	stick(0, -0.9)
	is.True(ui.GetFocusedWidget() == above) // Up

	ui.SetFocusedWidget(s)
	ui.GamepadFocus = false
	stick(0, 0)
	is.Equal(stick(0.9, 0), 0) // GamepadFocus off: nothing
}

// Back, from Escape or, with GamepadFocus on, a gamepad's right face
// button, calls BackHandler once a press, unless the focused widget takes
// it: an open list combo button closes instead.
func TestUI_Back(t *testing.T) {
	is := is.New(t)

	combo := widget.NewListComboButton(
		widget.ListComboButtonOpts.Entries([]any{"a", "b"}),
		widget.ListComboButtonOpts.EntryLabelFunc(func(e any) string { return e.(string) }, func(e any) string { return e.(string) }), //nolint:forcetypeassert // strings
	)
	combo.GetWidget().Rect = image.Rect(0, 0, 100, 20)
	ui := focusUI(combo)
	ui.PrimaryTheme = themes.GetBasicDarkTheme() // the combo's images and faces
	backs := 0
	ui.BackHandler = func() { backs++ }
	_, press := useDevices(t, ui)

	press(ebiten.KeyEscape)
	is.Equal(backs, 1) // once, though held
	press(ebiten.StandardGamepadButtonRightRight)
	is.Equal(backs, 1) // GamepadFocus off
	ui.GamepadFocus = true
	press(ebiten.StandardGamepadButtonRightRight)
	is.Equal(backs, 2)

	ui.SetFocusedWidget(combo)
	combo.SetContentVisible(true)
	press(ebiten.KeyEscape)
	is.True(!combo.ContentVisible()) // the open combo took it
	is.Equal(backs, 2)
}
