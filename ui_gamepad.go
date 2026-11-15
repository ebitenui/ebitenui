package ebitenui

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Ebitengine's gamepad state, reached through variables so tests can stand
// in for a gamepad.
var (
	appendGamepadIDs   = ebiten.AppendGamepadIDs
	isStandardGamepad  = ebiten.IsStandardGamepadLayoutAvailable
	gamepadJustPressed = inpututil.IsStandardGamepadButtonJustPressed
)

// gamepadDirections are the D-pad's buttons and where each moves the focus.
var gamepadDirections = [...]struct {
	button    ebiten.StandardGamepadButton
	direction widget.FocusDirection
}{
	{ebiten.StandardGamepadButtonLeftTop, widget.FOCUS_NORTH},
	{ebiten.StandardGamepadButtonLeftBottom, widget.FOCUS_SOUTH},
	{ebiten.StandardGamepadButtonLeftLeft, widget.FOCUS_WEST},
	{ebiten.StandardGamepadButtonLeftRight, widget.FOCUS_EAST},
}

// handleGamepadFocus moves the focus with the D-pad of each gamepad with the
// standard layout, and activates the focused widget with its bottom face
// button, while GamepadFocus is on.
func (u *UI) handleGamepadFocus() {
	if !u.GamepadFocus {
		return
	}
	u.gamepads = appendGamepadIDs(u.gamepads[:0])
	for _, id := range u.gamepads {
		if !isStandardGamepad(id) {
			continue
		}
		for _, g := range gamepadDirections {
			if gamepadJustPressed(id, g.button) {
				u.ChangeFocus(g.direction)
			}
		}
		if gamepadJustPressed(id, ebiten.StandardGamepadButtonRightBottom) {
			u.ActivateFocused()
		}
	}
}
