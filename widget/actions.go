package widget

// An Action is a navigation command, the same from keys and from gamepads.
// The UI gives the actions of each update (UpdateObject.Actions); the
// focused widget takes those it acts on, a list moving through its entries
// or a slider its handle, and the UI acts on those left.
type Action int

const (
	// The directions: arrow keys, and a gamepad's D-pad with
	// ebitenui.UI.GamepadFocus on. Left by the focused widget, a gamepad's
	// move the focus.
	ActionUp Action = iota
	ActionDown
	ActionLeft
	ActionRight
	// ActionActivate comes from a gamepad's bottom face button. Left by the
	// focused widget, it clicks, toggles or submits it
	// (ebitenui.UI.ActivateFocused). Enter and Space activate widgets
	// themselves.
	ActionActivate

	actionCount
)

// Actions are the navigation actions of one update, each taken by one
// widget at most. A nil *Actions gives none.
type Actions struct {
	given, taken [actionCount]bool
}

// Reset forgets the actions given and taken.
func (a *Actions) Reset() { *a = Actions{} }

// Give gives action act.
func (a *Actions) Give(act Action) { a.given[act] = true }

// Take reports whether act is given and not taken yet, and takes it.
func (a *Actions) Take(act Action) bool {
	if a == nil || !a.given[act] || a.taken[act] {
		return false
	}
	a.taken[act] = true
	return true
}

// Left reports whether act is given and no widget took it.
func (a *Actions) Left(act Action) bool {
	return a != nil && a.given[act] && !a.taken[act]
}
