package event

// A DeferredAction is an action that is executed at a later time.
type DeferredAction interface {
	// Do executes the action.
	Do()
}

var deferredActions []DeferredAction

// AddDeferred adds d to the queue of deferred actions.
func AddDeferred(d DeferredAction) {
	deferredActions = append(deferredActions, d)
}

// ExecuteDeferred processes the queue of deferred actions and executes them.
func ExecuteDeferred() {
	defer func(d []DeferredAction) {
		deferredActions = d[:0]
	}(deferredActions)

	for len(deferredActions) > 0 {
		a := deferredActions[0]
		deferredActions = deferredActions[1:]

		a.Do()
	}
}
