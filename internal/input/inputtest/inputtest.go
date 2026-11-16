// Package inputtest stands in for the keyboard, the gamepads and the
// touches in tests.
package inputtest

import (
	"image"
	"testing"

	internalinput "github.com/ebitenui/ebitenui/internal/input"
	"github.com/hajimehoshi/ebiten/v2"
)

// Devices is a keyboard, gamepads and a screen a test presses.
type Devices struct {
	Keys [ebiten.KeyMax + 1]bool
	// Pads are the gamepads connected, all with the standard layout; pad i
	// has the ID i.
	Pads []Pad
	// Touches are the fingers on the screen, where each is; their IDs count
	// up from FirstTouchID, so changing it puts down other fingers.
	// Touchscreen makes the platform one without a mouse pointer.
	Touches      []image.Point
	FirstTouchID ebiten.TouchID
	Touchscreen  bool
}

// Pad is a gamepad with the standard layout.
type Pad struct {
	Buttons [ebiten.StandardGamepadButtonMax + 1]bool
	Axes    [ebiten.StandardGamepadAxisMax + 1]float64
}

// Use makes the input handler read the returned devices until the test
// ends, and then updates it once with nothing pressed.
func Use(t testing.TB) *Devices {
	t.Helper()
	d := &Devices{}
	h := internalinput.InputHandler
	old := h.Devices
	h.Devices = d
	t.Cleanup(func() {
		*d = Devices{}
		h.Update()
		h.AfterUpdate()
		h.Devices = old
	})
	return d
}

func (d *Devices) IsKeyPressed(k ebiten.Key) bool { return d.Keys[k] }

func (d *Devices) AppendGamepadIDs(ids []ebiten.GamepadID) []ebiten.GamepadID {
	for i := range d.Pads {
		ids = append(ids, ebiten.GamepadID(i))
	}
	return ids
}

func (d *Devices) IsStandardGamepadButtonPressed(id ebiten.GamepadID, b ebiten.StandardGamepadButton) bool {
	return int(id) < len(d.Pads) && d.Pads[id].Buttons[b]
}

func (d *Devices) StandardGamepadAxisValue(id ebiten.GamepadID, a ebiten.StandardGamepadAxis) float64 {
	if int(id) >= len(d.Pads) {
		return 0
	}
	return d.Pads[id].Axes[a]
}

func (d *Devices) AppendTouchIDs(ids []ebiten.TouchID) []ebiten.TouchID {
	for i := range d.Touches {
		ids = append(ids, d.FirstTouchID+ebiten.TouchID(i))
	}
	return ids
}

func (d *Devices) TouchPosition(id ebiten.TouchID) (int, int) {
	i := int(id - d.FirstTouchID)
	if i < 0 || i >= len(d.Touches) {
		return 0, 0
	}
	return d.Touches[i].X, d.Touches[i].Y
}

func (d *Devices) IsTouchscreen() bool { return d.Touchscreen }
