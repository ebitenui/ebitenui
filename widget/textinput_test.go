package widget

import (
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/matryer/is"
)

func TestTextInput_ChangedEvent(t *testing.T) {
	is := is.New(t)

	var eventArgs *TextInputChangedEventArgs
	ti := newTextInput(t, TextInputOpts.ChangedHandler(func(args *TextInputChangedEventArgs) {
		eventArgs = args
	}))

	ti.SetText("foo")
	render(ti, t)

	is.Equal(eventArgs.InputText, "foo")
}

func TestTextInput_ChangedEvent_OnlyOnce(t *testing.T) {
	is := is.New(t)

	numEvents := 0
	ti := newTextInput(t, TextInputOpts.ChangedHandler(func(_ *TextInputChangedEventArgs) {
		numEvents++
	}))

	ti.SetText("foo")
	render(ti, t)
	render(ti, t)

	is.Equal(numEvents, 1)
}

func TestTextInput_DoBackspace(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t)
	ti.SetText("foo")
	ti.lastInputText = ti.GetText()
	ti.cursorPosition = 1
	render(ti, t)

	ti.ChangedEvent.AddHandler(func(args interface{}) {
		is.Equal(args.(*TextInputChangedEventArgs).InputText, "oo")
	})

	ti.Backspace()
	render(ti, t)
}

func TestTextInput_DoBackspace_Disabled(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t)
	ti.GetWidget().Disabled = true
	ti.SetText("foo")
	ti.cursorPosition = 1
	render(ti, t)

	ti.ChangedEvent.AddHandler(func(_ interface{}) {
		is.Fail() // received event even though widget is disabled
	})

	ti.Backspace()
	render(ti, t)
}

func TestTextInput_DoDelete(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t)
	ti.SetText("foo")
	render(ti, t)

	ti.ChangedEvent.AddHandler(func(args interface{}) {
		is.Equal(args.(*TextInputChangedEventArgs).InputText, "oo")
	})

	ti.Delete()
	render(ti, t)
}

func TestTextInput_DoDelete_Disabled(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t)
	ti.GetWidget().Disabled = true
	ti.SetText("foo")
	render(ti, t)

	ti.ChangedEvent.AddHandler(func(_ interface{}) {
		is.Fail() // received event even though widget is disabled
	})

	ti.Delete()
	render(ti, t)
}

func TestTextInput_DoInsert(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t)
	ti.SetText("foo")
	ti.lastInputText = ti.GetText()
	ti.cursorPosition = 1
	render(ti, t)

	ti.Insert([]rune("ab€c"))

	is.Equal(ti.GetText(), "fab€coo")
	is.Equal(ti.cursorPosition, 5)
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
