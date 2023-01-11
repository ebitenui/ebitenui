package widget

import (
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/matryer/is"
)

func TestCheckbox_State_Initial(t *testing.T) {
	is := is.New(t)

	c := newCheckbox(t,
		CheckboxOpts.StateChangedHandler(func(_ *CheckboxChangedEventArgs) {
			is.Fail() // event fired without previous action
		}))

	is.Equal(c.State(), WidgetUnchecked)
}

func TestCheckbox_ChangedEvent_User(t *testing.T) {
	is := is.New(t)

	var eventArgs *CheckboxChangedEventArgs

	c := newCheckbox(t,
		CheckboxOpts.StateChangedHandler(func(args *CheckboxChangedEventArgs) {
			eventArgs = args
		}))

	leftMouseButtonClick(c, t)

	is.Equal(eventArgs.State, WidgetChecked)
	is.Equal(c.State(), WidgetChecked)
}

func TestCheckbox_SetState(t *testing.T) {
	is := is.New(t)

	var eventArgs *CheckboxChangedEventArgs
	numEvents := 0

	c := newCheckbox(t,
		CheckboxOpts.StateChangedHandler(func(args *CheckboxChangedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	c.SetState(WidgetChecked)
	event.ExecuteDeferred()

	is.Equal(eventArgs.State, WidgetChecked)
	is.Equal(c.State(), WidgetChecked)

	c.SetState(WidgetChecked)
	event.ExecuteDeferred()

	is.Equal(numEvents, 1)
}

func TestCheckbox_State_Cycle(t *testing.T) {
	is := is.New(t)

	c := newCheckbox(t)
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), WidgetChecked)
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), WidgetUnchecked)
}

func TestCheckbox_State_Cycle_TriState(t *testing.T) {
	is := is.New(t)

	c := newCheckbox(t, CheckboxOpts.TriState())
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), WidgetChecked)
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), WidgetGreyed)
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), WidgetUnchecked)
}

func newCheckbox(t *testing.T, opts ...CheckboxOpt) *Checkbox {
	t.Helper()

	c := NewCheckbox(append(opts, []CheckboxOpt{
		CheckboxOpts.ButtonOpts(ButtonOpts.Image(&ButtonImage{
			Idle: newNineSliceEmpty(t),
		})),

		CheckboxOpts.Image(&CheckboxGraphicImage{
			Unchecked: &ButtonImageImage{
				Idle: newImageEmpty(t),
			},
			Checked: &ButtonImageImage{
				Idle: newImageEmpty(t),
			},
		}),
	}...)...)
	event.ExecuteDeferred()
	render(c, t)
	return c
}
