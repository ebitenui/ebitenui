package event

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

type firedEvent struct {
	event *Event
	args  interface{}
}

type deferredAddHandler struct {
	event   *Event
	handler handler
}

var deferredEvents []*firedEvent
var deferredAddHandlers []deferredAddHandler

// AddHandler registers handler h with e. It returns a function to remove h from e if desired.
func (e *Event) AddHandler(h HandlerFunc) RemoveHandlerFunc {
	e.idCounter++

	id := e.idCounter
	handler := handler{
		id: id,
		h:  h,
	}

	deferredAddHandlers = append(deferredAddHandlers, deferredAddHandler{
		event:   e,
		handler: handler,
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
	deferredEvents = append(deferredEvents, &firedEvent{
		event: e,
		args:  args,
	})
}

func (e *Event) handle(args interface{}) {
	for _, h := range e.handlers {
		h.h(args)
	}
}

// FireDeferredEvents processes the deferred queue of events and fires those events.
// This function should not be called directly.
func FireDeferredEvents() {
	for len(deferredEvents) > 0 {
		f := deferredEvents[0]
		deferredEvents = deferredEvents[1:]

		f.event.handle(f.args)
	}

	for _, d := range deferredAddHandlers {
		d.event.handlers = append(d.event.handlers, d.handler)
	}
	deferredAddHandlers = deferredAddHandlers[:0]
}
