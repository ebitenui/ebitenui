package event

import internalevent "github.com/ebitenui/ebitenui/internal/event"

// ExecuteDeferred processes the queue of deferred actions and executes them. This should only be called by UI.
// Additionally, it can be called in unit tests to process events programmatically.
func ExecuteDeferred() {
	internalevent.ExecuteDeferred()
}
