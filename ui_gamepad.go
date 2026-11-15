package ebitenui

import (
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
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

// handleGamepadFocus moves the focus with the gamepads' D-pad, and activates
// the focused widget with their bottom face button, while GamepadFocus is
// on. The input package merges the gamepads, so a button pressed on two at
// once acts once.
func (u *UI) handleGamepadFocus() {
	if !u.GamepadFocus {
		return
	}
	for _, g := range gamepadDirections {
		if input.GamepadButtonJustPressed(g.button) {
			u.ChangeFocus(g.direction)
		}
	}
	if input.GamepadButtonJustPressed(ebiten.StandardGamepadButtonRightBottom) {
		u.ActivateFocused()
	}
}
