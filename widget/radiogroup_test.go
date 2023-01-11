package widget

import (
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/matryer/is"
)

func TestRadioGroup_Active_Initial(t *testing.T) {
	is := is.New(t)

	cbs := []*Checkbox{}
	for i := 0; i < 3; i++ {
		c := newCheckbox(t)
		cbs = append(cbs, c)
	}

	var eventArgs *RadioGroupChangedEventArgs

	r := newRadioGroup(t, cbs, RadioGroupOpts.ChangedHandler(func(args *RadioGroupChangedEventArgs) {
		eventArgs = args
	}))

	is.Equal(r.Active(), cbs[0])
	is.Equal(eventArgs.Active, cbs[0])
	is.Equal(cbs[0].State(), WidgetChecked)
}

func TestRadioGroup_ChangedEvent_User(t *testing.T) {
	is := is.New(t)

	cbs := []*Checkbox{}
	for i := 0; i < 3; i++ {
		c := newCheckbox(t)
		cbs = append(cbs, c)
	}

	r := newRadioGroup(t, cbs)

	var eventArgs *RadioGroupChangedEventArgs
	r.ChangedEvent.AddHandler(func(args interface{}) {
		eventArgs = args.(*RadioGroupChangedEventArgs)
	})

	leftMouseButtonClick(cbs[1], t)

	is.Equal(r.Active(), cbs[1])
	is.Equal(eventArgs.Active, cbs[1])
	is.Equal(cbs[0].State(), WidgetUnchecked)
	is.Equal(cbs[1].State(), WidgetChecked)
}

func TestRadioGroup_SetActive(t *testing.T) {
	is := is.New(t)

	cbs := []*Checkbox{}
	for i := 0; i < 3; i++ {
		c := newCheckbox(t)
		cbs = append(cbs, c)
	}

	r := newRadioGroup(t, cbs)

	var eventArgs *RadioGroupChangedEventArgs
	r.ChangedEvent.AddHandler(func(args interface{}) {
		eventArgs = args.(*RadioGroupChangedEventArgs)
	})

	r.SetActive(cbs[1])
	event.ExecuteDeferred()

	is.Equal(r.Active(), cbs[1])
	is.Equal(eventArgs.Active, cbs[1])
	is.Equal(cbs[0].State(), WidgetUnchecked)
	is.Equal(cbs[1].State(), WidgetChecked)
}

func newRadioGroup(t *testing.T, cbs []*Checkbox, opts ...RadioGroupOpt) *RadioGroup {
	t.Helper()

	elements := []RadioGroupElement{}

	for _, cb := range cbs {
		elements = append(elements, cb)
	}
	r := NewRadioGroup(append(opts, RadioGroupOpts.Elements(elements...))...)
	event.ExecuteDeferred()
	for _, c := range cbs {
		render(c, t)
	}
	return r
}
