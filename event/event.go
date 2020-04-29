package event

type Event struct {
	idCounter uint32
	handlers  []handler
}

type HandlerFunc func(args interface{})

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
