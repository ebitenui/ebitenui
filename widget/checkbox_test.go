package widget

import (
	"testing"

	internalevent "github.com/blizzy78/ebitenui/internal/event"
	"github.com/matryer/is"
)

func TestCheckbox_State_Initial(t *testing.T) {
	is := is.New(t)

	c := newCheckbox(t,
		CheckboxOpts.ChangedHandler(func(args *CheckboxChangedEventArgs) {
			is.Fail() // event fired without previous action
		}))

	is.Equal(c.State(), CheckboxUnchecked)
}

func TestCheckbox_ChangedEvent_User(t *testing.T) {
	is := is.New(t)

	var eventArgs *CheckboxChangedEventArgs

	c := newCheckbox(t,
		CheckboxOpts.ChangedHandler(func(args *CheckboxChangedEventArgs) {
			eventArgs = args
		}))

	leftMouseButtonClick(c, t)

	is.Equal(eventArgs.State, CheckboxChecked)
	is.Equal(c.State(), CheckboxChecked)
}

func TestCheckbox_SetState(t *testing.T) {
	is := is.New(t)

	var eventArgs *CheckboxChangedEventArgs
	numEvents := 0

	c := newCheckbox(t,
		CheckboxOpts.ChangedHandler(func(args *CheckboxChangedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	c.SetState(CheckboxChecked)
	internalevent.ExecuteDeferredActions()

	is.Equal(eventArgs.State, CheckboxChecked)
	is.Equal(c.State(), CheckboxChecked)

	c.SetState(CheckboxChecked)
	internalevent.ExecuteDeferredActions()

	is.Equal(numEvents, 1)
}

func TestCheckbox_State_Cycle(t *testing.T) {
	is := is.New(t)

	c := newCheckbox(t)
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), CheckboxChecked)
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), CheckboxUnchecked)
}

func TestCheckbox_State_Cycle_TriState(t *testing.T) {
	is := is.New(t)

	c := newCheckbox(t, CheckboxOpts.TriState())
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), CheckboxChecked)
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), CheckboxGreyed)
	leftMouseButtonClick(c, t)
	is.Equal(c.State(), CheckboxUnchecked)
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
		}),
	}...)...)
	internalevent.ExecuteDeferredActions()
	render(c, t)
	return c
}
