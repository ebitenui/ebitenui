package widget

import "github.com/blizzy78/ebitenui/event"

type RadioGroup struct {
	ChangedEvent *event.Event

	checkboxes []*Checkbox
	active     *Checkbox
	listen     bool
	doneEvent  *event.Event
}

type RadioGroupOpt func(r *RadioGroup)

const RadioGroupOpts = radioGroupOpts(true)

type radioGroupOpts bool

type RadioGroupChangedEventArgs struct {
	Active *Checkbox
}

type RadioGroupChangedHandlerFunc func(args *RadioGroupChangedEventArgs)

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

func (o radioGroupOpts) Checkboxes(cb ...*Checkbox) RadioGroupOpt {
	return func(r *RadioGroup) {
		r.checkboxes = cb
	}
}

func (o radioGroupOpts) ChangedHandler(f RadioGroupChangedHandlerFunc) RadioGroupOpt {
	return func(r *RadioGroup) {
		r.ChangedEvent.AddHandler(func(args interface{}) {
			f(args.(*RadioGroupChangedEventArgs))
		})
	}
}

func (r *RadioGroup) Active() *Checkbox {
	return r.active
}

func (r *RadioGroup) SetActive(a *Checkbox) {
	r.listen = false

	oldActive := r.active
	for _, c := range r.checkboxes {
		if c == a {
			r.active = c

			// ignore unchecking and reset to checked
			c.SetState(CheckboxChecked)
		} else {
			c.SetState(CheckboxUnchecked)
		}
	}

	// SetState() fires deferred events, so we need something *after* those to tell us we should listen again
	event.AddEventHandlerOneShot(r.doneEvent, func(_ interface{}) {
		r.listen = true
	})
	r.doneEvent.Fire(nil)

	if r.active != oldActive {
		r.ChangedEvent.Fire(&RadioGroupChangedEventArgs{
			Active: r.active,
		})
	}
}

func (r *RadioGroup) create() {
	for _, c := range r.checkboxes {
		c.ChangedEvent.AddHandler(func(args interface{}) {
			if !r.listen {
				return
			}

			a := args.(*CheckboxChangedEventArgs)
			r.SetActive(a.Checkbox)
		})
	}

	if r.active == nil {
		r.SetActive(r.checkboxes[0])
	}
}
