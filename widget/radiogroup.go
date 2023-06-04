package widget

import (
	"github.com/ebitenui/ebitenui/event"
)

type WidgetState int

const (
	WidgetUnchecked = WidgetState(iota)
	WidgetChecked
	WidgetGreyed
)

type RadioGroupElement interface {
	SetState(state WidgetState)
	getStateChangedEvent() *event.Event
}

type RadioGroup struct {
	ChangedEvent *event.Event

	elements  []RadioGroupElement
	active    RadioGroupElement
	initial   RadioGroupElement
	listen    bool
	doneEvent *event.Event
}

type RadioGroupOpt func(r *RadioGroup)

type RadioGroupOptions struct {
}

type RadioGroupChangedEventArgs struct {
	Active RadioGroupElement
}

type RadioGroupChangedHandlerFunc func(args *RadioGroupChangedEventArgs)

var RadioGroupOpts RadioGroupOptions

func NewRadioGroup(opts ...RadioGroupOpt) *RadioGroup {
	r := &RadioGroup{
		ChangedEvent: &event.Event{},

		listen:    true,
		doneEvent: &event.Event{},
	}

	for _, o := range opts {
		o(r)
	}

	// use deferred event to initialize
	e := &event.Event{}
	event.AddEventHandlerOneShot(e, func(_ interface{}) {
		r.create()
	})
	e.Fire(nil)

	return r
}

func (o RadioGroupOptions) Elements(e ...RadioGroupElement) RadioGroupOpt {
	return func(r *RadioGroup) {
		for idx := range e {
			switch eletype := e[idx].(type) {
			case *Button:
				eletype.ToggleMode = true
			}
		}
		r.elements = e
	}
}

func (o RadioGroupOptions) ChangedHandler(f RadioGroupChangedHandlerFunc) RadioGroupOpt {
	return func(r *RadioGroup) {
		r.ChangedEvent.AddHandler(func(args interface{}) {
			f(args.(*RadioGroupChangedEventArgs))
		})
	}
}

// This function allows you to select which element should be selected initialially.
// Otherwise it will select the first element in the Elements array.
func (o RadioGroupOptions) InitialElement(e RadioGroupElement) RadioGroupOpt {
	return func(r *RadioGroup) {
		r.initial = e
	}
}

func (r *RadioGroup) Active() RadioGroupElement {
	return r.active
}

func (r *RadioGroup) SetActive(a RadioGroupElement) {
	r.listen = false
	oldActive := r.active
	for _, c := range r.elements {
		if c == a {
			r.active = c

			// ignore unchecking and reset to checked
			c.SetState(WidgetChecked)
		} else {
			c.SetState(WidgetUnchecked)
		}
	}

	// SetState() fires deferred events, so we need something *after* those to tell us we should listen again
	event.AddEventHandlerOneShot(r.doneEvent, func(_ interface{}) {
		r.listen = true
	})
	r.doneEvent.Fire(nil)

	if a != oldActive {
		r.ChangedEvent.Fire(&RadioGroupChangedEventArgs{
			Active: a,
		})
	}
}

func (r *RadioGroup) create() {
	for _, c := range r.elements {
		c.getStateChangedEvent().AddHandler(func(args interface{}) {
			if !r.listen {
				return
			}
			switch args := args.(type) {
			case *CheckboxChangedEventArgs:
				r.SetActive(args.Active)
			case *ButtonChangedEventArgs:
				r.SetActive(args.Button)
			}
		})
	}

	if r.initial != nil {
		r.SetActive(r.initial)
	} else if len(r.elements) > 0 {
		r.SetActive(r.elements[0])
	}
}
