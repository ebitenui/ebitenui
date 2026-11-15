package ebitenui

import (
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

// directions are the direction actions: the arrow key and the D-pad button
// that give each, and where a gamepad's moves the focus when the focused
// widget leaves it.
var directions = [...]struct {
	action widget.Action
	key    ebiten.Key
	button ebiten.StandardGamepadButton
	focus  widget.FocusDirection
}{
	{widget.ActionUp, ebiten.KeyArrowUp, ebiten.StandardGamepadButtonLeftTop, widget.FOCUS_NORTH},
	{widget.ActionDown, ebiten.KeyArrowDown, ebiten.StandardGamepadButtonLeftBottom, widget.FOCUS_SOUTH},
	{widget.ActionLeft, ebiten.KeyArrowLeft, ebiten.StandardGamepadButtonLeftLeft, widget.FOCUS_WEST},
	{widget.ActionRight, ebiten.KeyArrowRight, ebiten.StandardGamepadButtonLeftRight, widget.FOCUS_EAST},
}

// navigation turns the keys and the gamepads into each update's actions.
type navigation struct {
	actions widget.Actions
	byPad   [len(directions)]bool // the directions a gamepad gave this update
}

// giveActions gives the update's actions: the directions from the arrow
// keys, and from the D-pad with GamepadFocus on, and Activate from a
// gamepad's bottom face button. The input package merges the gamepads, so
// a button pressed on two at once gives one action.
func (u *UI) giveActions() {
	n := &u.navigation
	n.actions.Reset()
	for i, d := range directions {
		n.byPad[i] = u.GamepadFocus && input.GamepadButtonJustPressed(d.button)
		if n.byPad[i] || input.KeyJustPressed(d.key) {
			n.actions.Give(d.action)
		}
	}
	if u.GamepadFocus && input.GamepadButtonJustPressed(ebiten.StandardGamepadButtonRightBottom) {
		n.actions.Give(widget.ActionActivate)
	}
}

// actOnActionsLeft acts on the actions no widget took: a direction from a
// gamepad moves the focus, Activate activates the focused widget.
func (u *UI) actOnActionsLeft() {
	n := &u.navigation
	for i, d := range directions {
		if n.byPad[i] && n.actions.Left(d.action) {
			u.ChangeFocus(d.focus)
		}
	}
	if n.actions.Left(widget.ActionActivate) {
		u.ActivateFocused()
	}
}
