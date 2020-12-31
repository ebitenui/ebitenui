package widget

import (
	"image/color"
	"testing"

	"github.com/blizzy78/ebitenui/event"
	"github.com/matryer/is"
)

func TestTextInput_ChangedEvent(t *testing.T) {
	is := is.New(t)

	var eventArgs *TextInputChangedEventArgs
	ti := newTextInput(t, TextInputOpts.ChangedHandler(func(args *TextInputChangedEventArgs) {
		eventArgs = args
	}))

	ti.InputText = "foo"
	render(ti, t)

	is.Equal(eventArgs.InputText, "foo")
}

func TestTextInput_ChangedEvent_OnlyOnce(t *testing.T) {
	is := is.New(t)

	numEvents := 0
	ti := newTextInput(t, TextInputOpts.ChangedHandler(func(args *TextInputChangedEventArgs) {
		numEvents++
	}))

	ti.InputText = "foo"
	render(ti, t)
	render(ti, t)

	is.Equal(numEvents, 1)
}

func TestTextInput_DoBackspace(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t)
	ti.InputText = "foo"
	ti.cursorPosition = 1
	render(ti, t)

	ti.ChangedEvent.AddHandler(func(args interface{}) {
		is.Equal(args.(*TextInputChangedEventArgs).InputText, "oo")
	})

	ti.doBackspace()
	render(ti, t)
}

func TestTextInput_DoBackspace_Disabled(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t)
	ti.GetWidget().Disabled = true
	ti.InputText = "foo"
	ti.cursorPosition = 1
	render(ti, t)

	ti.ChangedEvent.AddHandler(func(args interface{}) {
		is.Fail() // received event even though widget is disabled
	})

	ti.doBackspace()
	render(ti, t)
}

func TestTextInput_DoDelete(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t)
	ti.InputText = "foo"
	render(ti, t)

	ti.ChangedEvent.AddHandler(func(args interface{}) {
		is.Equal(args.(*TextInputChangedEventArgs).InputText, "oo")
	})

	ti.doDelete()
	render(ti, t)
}

func TestTextInput_DoDelete_Disabled(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t)
	ti.GetWidget().Disabled = true
	ti.InputText = "foo"
	render(ti, t)

	ti.ChangedEvent.AddHandler(func(args interface{}) {
		is.Fail() // received event even though widget is disabled
	})

	ti.doDelete()
	render(ti, t)
}

func newTextInput(t *testing.T, opts ...TextInputOpt) *TextInput {
	ti := NewTextInput(append(opts, []TextInputOpt{
		TextInputOpts.Face(loadFont(t)),
		TextInputOpts.Color(&TextInputColor{
			Idle:     color.White,
			Disabled: color.White,
			Caret:    color.White,
		}),
		TextInputOpts.CaretOpts(
			CaretOpts.Size(loadFont(t), 1)),
	}...)...)
	event.ExecuteDeferred()
	render(ti, t)
	return ti
}
