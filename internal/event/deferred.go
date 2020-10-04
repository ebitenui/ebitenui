package event

type DeferredAction interface {
	Do()
}

var deferredActions []DeferredAction

func AddDeferred(d DeferredAction) {
	deferredActions = append(deferredActions, d)
}

// ExecuteDeferredActions processes the queue of deferred actions and executes them.
func ExecuteDeferredActions() {
	defer func(d []DeferredAction) {
		deferredActions = d[:0]
	}(deferredActions)

	for len(deferredActions) > 0 {
		a := deferredActions[0]
		deferredActions = deferredActions[1:]

		a.Do()
	}
}
