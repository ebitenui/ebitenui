// Package inputtest stands in for the keyboard and the gamepads in tests.
package inputtest

import (
	"testing"

	internalinput "github.com/ebitenui/ebitenui/internal/input"
	"github.com/hajimehoshi/ebiten/v2"
)

// Devices is a keyboard and gamepads a test presses.
type Devices struct {
	Keys [ebiten.KeyMax + 1]bool
	// Pads are the gamepads connected, all with the standard layout; pad i
	// has the ID i.
	Pads []Pad
}

// Pad is a gamepad with the standard layout.
type Pad struct {
	Buttons [ebiten.StandardGamepadButtonMax + 1]bool
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
