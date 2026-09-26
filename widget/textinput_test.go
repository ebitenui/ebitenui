package widget

import (
	"image/color"
	"testing"
	"time"

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
		if arg, ok := args.(*TextInputChangedEventArgs); ok {
			is.Equal(arg.InputText, "oo")
		}
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
		if arg, ok := args.(*TextInputChangedEventArgs); ok {
			is.Equal(arg.InputText, "oo")
		}
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

	ti.Insert("ab€c")

	is.Equal(ti.GetText(), "fab€coo")
	is.Equal(ti.cursorPosition, 5)
}

func TestTextInput_WordBounds(t *testing.T) {
	tests := []struct {
		text          string
		index         int
		expectedStart int
		expectedEnd   int
	}{
		{"hello brave world", 8, 6, 11},
		{"hello brave world", 6, 6, 11},
		{"hello  world", 6, 5, 7},
		{"a...b", 2, 1, 4},
		{"snake_case name", 3, 0, 10},
		{"hello world", 11, 6, 11},
		{"", 0, 0, 0},
	}

	for _, test := range tests {
		t.Run(test.text, func(t *testing.T) {
			is := is.New(t)

			start, end := wordBounds([]rune(test.text), test.index)

			is.Equal(start, test.expectedStart)
			is.Equal(end, test.expectedEnd)
		})
	}
}

func TestTextInput_DoubleClick(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t)
	ti.SetText("hello brave world")
	now := time.Now()

	ti.press(8, false, now)
	ti.press(8, false, now.Add(100*time.Millisecond))

	is.Equal(ti.SelectedText(), "brave")
}

func TestTextInput_DoubleClick_TooSlow(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t)
	ti.SetText("hello brave world")
	now := time.Now()

	ti.press(8, false, now)
	ti.press(8, false, now.Add(time.Second))

	is.Equal(ti.SelectedText(), "")
	is.Equal(ti.cursorPosition, 8)
}

func TestTextInput_ShiftClick(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t)
	ti.SetText("hello brave world")
	ti.cursorPosition = 6

	ti.press(11, true, time.Now())

	is.Equal(ti.SelectedText(), "brave")
}

func TestTextInput_ShiftArrows(t *testing.T) {
	is := is.New(t)

	ti := newTextInput(t)
	ti.SetText("hello")
	ti.cursorPosition = 1

	ti.moveCursor(ti.CursorMoveRight, false, true)
	ti.moveCursor(ti.CursorMoveRight, false, true)
	is.Equal(ti.SelectedText(), "el")

	ti.moveCursor(ti.CursorMoveLeft, true, false)
	is.Equal(ti.SelectedText(), "")
	is.Equal(ti.cursorPosition, 1) // cursor lands on the start of the old selection

	ti.moveCursor(ti.CursorMoveRight, false, false)
	is.Equal(ti.cursorPosition, 2)
}

func newTextInput(t *testing.T, opts ...TextInputOpt) *TextInput {
	ti := NewTextInput(append(opts, []TextInputOpt{
		TextInputOpts.Face(loadFont(t)),
		TextInputOpts.Color(&TextInputColor{
			Idle:     color.White,
			Disabled: color.White,
			Caret:    color.White,
		}),
	}...)...)
	event.ExecuteDeferred()
	render(ti, t)
	return ti
}
