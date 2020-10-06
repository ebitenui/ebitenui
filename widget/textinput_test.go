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

func newTextInput(t *testing.T, opts ...TextInputOpt) *TextInput {
	ti := NewTextInput(append(opts, []TextInputOpt{
		TextInputOpts.Face(loadFont(t)),
		TextInputOpts.Color(&TextInputColor{
			Idle:  color.White,
			Caret: color.White,
		}),
	}...)...)
	event.ExecuteDeferred()
	render(ti, t)
	return ti
}
