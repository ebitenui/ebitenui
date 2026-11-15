package input

import (
	"testing"

	"github.com/ebitenui/ebitenui/internal/input/inputtest"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

// A key reads as just pressed in the update it goes down only, and as
// pressed while it is held; keys out of range read as up.
func TestKeyJustPressed(t *testing.T) {
	is := is.New(t)

	d := inputtest.Use(t)
	for i, s := range []struct{ down, pressed, just bool }{
		{false, false, false},
		{true, true, true},  // goes down
		{true, true, false}, // held
		{false, false, false},
		{true, true, true}, // down again
	} {
		d.Keys[ebiten.KeyEnter] = s.down
		Update()
		is.Equal([]any{i, KeyPressed(ebiten.KeyEnter), KeyJustPressed(ebiten.KeyEnter), AnyKeyPressed()}, []any{i, s.pressed, s.just, s.down})
		AfterUpdate()
	}
	for _, k := range []ebiten.Key{-1, ebiten.KeyMax + 1} {
		is.True(!KeyPressed(k) && !KeyJustPressed(k))
	}
}

// A gamepad button is down while any gamepad holds it, so two gamepads
// pressing it together press it once; a gamepad unplugged lets go of its
// buttons.
func TestGamepadButtons(t *testing.T) {
	is := is.New(t)

	d := inputtest.Use(t)
	d.Pads = make([]inputtest.Pad, 2)
	b := ebiten.StandardGamepadButtonRightBottom
	for i, s := range []struct{ first, second, pressed, just bool }{
		{false, false, false, false},
		{true, true, true, true},   // both at once: one press
		{false, true, true, false}, // still held by the second
		{false, false, false, false},
		{false, true, true, true}, // the second alone
	} {
		d.Pads[0].Buttons[b], d.Pads[1].Buttons[b] = s.first, s.second
		Update()
		is.Equal([]any{i, GamepadButtonPressed(b), GamepadButtonJustPressed(b)}, []any{i, s.pressed, s.just})
		AfterUpdate()
	}

	d.Pads = d.Pads[:1] // the second, holding b, unplugged
	Update()
	is.True(!GamepadButtonPressed(b))
	for _, b := range []ebiten.StandardGamepadButton{-1, ebiten.StandardGamepadButtonMax + 1} {
		is.True(!GamepadButtonPressed(b) && !GamepadButtonJustPressed(b))
	}
}

// An axis reads as the gamepad pushing it furthest, either way, so a pad
// at rest or drifting a little doesn't hide another one in use; a gamepad
// unplugged leaves its axes.
func TestGamepadAxes(t *testing.T) {
	is := is.New(t)

	d := inputtest.Use(t)
	d.Pads = make([]inputtest.Pad, 2)
	a := ebiten.StandardGamepadAxisLeftStickHorizontal
	for _, v := range [][2]float64{{0.2, -0.9}, {-0.9, 0.2}} {
		d.Pads[0].Axes[a], d.Pads[1].Axes[a] = v[0], v[1]
		Update()
		is.Equal(GamepadAxis(a), -0.9)
	}
	is.Equal(GamepadAxis(ebiten.StandardGamepadAxisLeftStickVertical), 0.0)

	d.Pads = d.Pads[1:] // the first, at -0.9, unplugged
	Update()
	is.Equal(GamepadAxis(a), 0.2)
	for _, a := range []ebiten.StandardGamepadAxis{-1, ebiten.StandardGamepadAxisMax + 1} {
		is.Equal(GamepadAxis(a), 0.0)
	}
}

// Taking the keyboard and the gamepads allocates nothing, however many
// updates run (a list growing by a little each update would only show over
// many).
func TestUpdate_DoesNotAllocate(t *testing.T) {
	is := is.New(t)

	d := inputtest.Use(t)
	d.Pads = make([]inputtest.Pad, 2)
	d.Keys[ebiten.KeyA] = true
	is.Equal(testing.AllocsPerRun(1, func() {
		for range 10_000 {
			Update()
			AfterUpdate()
		}
	}), 0.0)
}
