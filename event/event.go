package event

import internalevent "github.com/blizzy78/ebitenui/internal/event"

// Event encapsulates an arbitrary event that event handlers may be interested in.
type Event struct {
	idCounter uint32
	handlers  []handler
}

// A HandlerFunc is a function that receives and handles an event.
type HandlerFunc func(args interface{})

// RemoveHandlerFunc is a function that removes a handler from an event.
type RemoveHandlerFunc func()

type handler struct {
	id uint32
	h  HandlerFunc
}

type deferredEvent struct {
	event *Event
	args  interface{}
}

type deferredAddHandler struct {
	event   *Event
	handler handler
}

// AddHandler registers handler h with e. It returns a function to remove h from e if desired.
func (e *Event) AddHandler(h HandlerFunc) RemoveHandlerFunc {
	e.idCounter++

	id := e.idCounter

	internalevent.AddDeferred(&deferredAddHandler{
		event: e,
		handler: handler{
			id: id,
			h:  h,
		},
	})

	return func() {
		e.removeHandler(id)
	}
}

func (e *Event) removeHandler(id uint32) {
	index := -1
	for i, h := range e.handlers {
		if h.id == id {
			index = i
			break
		}
	}

	if index < 0 {
		return
	}

	e.handlers = append(e.handlers[:index], e.handlers[index+1:]...)
}

// Fire fires an event to all registered handlers.
//
// Events are not fired directly, but are put into a deferred queue. This queue is then
// processed by the UI system.
func (e *Event) Fire(args interface{}) {
	internalevent.AddDeferred(&deferredEvent{
		event: e,
		args:  args,
	})
}

func (e *Event) handle(args interface{}) {
	for _, h := range e.handlers {
		h.h(args)
	}
}

// Do implements DeferredAction.
func (e *deferredEvent) Do() {
	e.event.handle(e.args)
}

// Do implements DeferredAction.
func (a *deferredAddHandler) Do() {
	a.event.handlers = append(a.event.handlers, a.handler)
}
