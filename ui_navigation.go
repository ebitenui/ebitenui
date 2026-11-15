package ebitenui

import (
	"image"
	"math"
	"time"

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

// The repeat of a direction held, unless the UI sets its own.
const (
	defaultRepeatDelay    = 400 * time.Millisecond
	defaultRepeatInterval = 80 * time.Millisecond
)

// The left stick points a way once pushed past stickPress along it, and
// keeps pointing it until it falls under stickRelease.
const (
	stickPress   = 0.5
	stickRelease = 0.33
)

// navigation turns the keys and the gamepads into each update's actions.
type navigation struct {
	actions widget.Actions
	held    [len(directions)]int  // the updates each direction has been held
	byPad   [len(directions)]bool // whether a gamepad gave each direction when last given
	stick   int                   // the direction the left stick points, its index plus one; 0 for none
}

// giveActions gives the update's actions: the directions from the arrow
// keys, and from the D-pad and the left stick with GamepadFocus on,
// repeating while held, and Activate from a gamepad's bottom face button.
// The input package merges the gamepads, so a button pressed on two at once
// gives one action.
func (u *UI) giveActions() {
	n := &u.navigation
	n.actions.Reset()
	n.stick = stickDirection(input.GamepadAxis(ebiten.StandardGamepadAxisLeftStickHorizontal),
		input.GamepadAxis(ebiten.StandardGamepadAxisLeftStickVertical), n.stick)
	delay, interval := u.repeatUpdates()
	for i, d := range directions {
		byPad := u.GamepadFocus && (input.GamepadButtonPressed(d.button) || n.stick == i+1)
		if !byPad && !input.KeyPressed(d.key) {
			n.held[i] = 0
			continue
		}
		n.held[i]++
		if repeats(n.held[i], delay, interval) {
			n.actions.Give(d.action)
			n.byPad[i] = byPad
		}
	}
	if u.GamepadFocus && input.GamepadButtonJustPressed(ebiten.StandardGamepadButtonRightBottom) {
		n.actions.Give(widget.ActionActivate)
	}
	if input.KeyJustPressed(ebiten.KeyEscape) ||
		u.GamepadFocus && input.GamepadButtonJustPressed(ebiten.StandardGamepadButtonRightRight) {
		n.actions.Give(widget.ActionBack)
	}
}

// stickDirection returns the direction a stick at x, y points, as an index
// into directions plus one, or 0 for none: the way it is pushed furthest,
// past stickPress, or past stickRelease for the way it pointed (last).
func stickDirection(x, y float64, last int) int {
	pointed, along := 0, 0.0
	for i, v := range [len(directions)]float64{-y, y, -x, x} { // up, down, left, right
		threshold := stickPress
		if i+1 == last {
			threshold = stickRelease
		}
		if v >= threshold && v > along {
			pointed, along = i+1, v
		}
	}
	return pointed
}

// repeatUpdates returns the repeat's delay and interval in updates, at
// Ebitengine's TPS, or 60 a second when the TPS follows the display.
func (u *UI) repeatUpdates() (delay, interval int) {
	tps := float64(ebiten.TPS())
	if tps <= 0 {
		tps = ebiten.DefaultTPS
	}
	d, i := u.NavigationRepeatDelay, u.NavigationRepeatInterval
	if d <= 0 {
		d = defaultRepeatDelay
	}
	if i <= 0 {
		i = defaultRepeatInterval
	}
	return int(math.Round(d.Seconds() * tps)), max(1, int(math.Round(i.Seconds()*tps)))
}

// repeats reports whether a direction held for held updates acts in this
// one: as it goes down, then once delay updates have passed, every
// interval updates.
func repeats(held, delay, interval int) bool {
	return held == 1 || held > delay && (held-1-delay)%interval == 0
}

// actOnActionsLeft acts on the actions no widget took: a direction from a
// gamepad moves the focus, Activate opens the on-screen keyboard for a text
// input or activates the focused widget, Back closes the open keyboard or
// calls BackHandler.
func (u *UI) actOnActionsLeft() {
	n := &u.navigation
	for i, d := range directions {
		if n.byPad[i] && n.actions.Left(d.action) {
			u.ChangeFocus(d.focus)
		}
	}
	if n.actions.Left(widget.ActionActivate) && !u.openKeyboard() {
		u.ActivateFocused()
	}
	if n.actions.Left(widget.ActionBack) {
		if k := u.OnScreenKeyboard; k != nil && k.IsOpen() {
			k.Close()
		} else if u.BackHandler != nil {
			u.BackHandler()
		}
	}
}

// openKeyboard opens the OnScreenKeyboard for the focused text input, in a
// modal window across the bottom of the screen, and focuses its first
// key; closing it gives the focus back to the input. It reports whether it
// opened.
func (u *UI) openKeyboard() bool {
	k := u.OnScreenKeyboard
	t, ok := u.focusedWidget.(*widget.TextInput)
	if k == nil || !ok || t.GetWidget().Disabled {
		return false
	}
	if u.keyboardWindow == nil || u.keyboardWindow.Contents != k.GetContainer() {
		u.keyboardWindow = widget.NewWindow(widget.WindowOpts.Contents(k.GetContainer()), widget.WindowOpts.Modal())
	}
	screen := input.GetWindowSize()
	_, h := k.GetContainer().PreferredSize()
	u.keyboardWindow.SetLocation(image.Rect(0, screen.Y-h, screen.X, screen.Y))
	remove := u.AddWindow(u.keyboardWindow)
	k.Open(t, func() {
		remove()
		u.SetFocusedWidget(t)
	})
	u.SetFocusedWidget(k.FirstKey())
	return true
}
