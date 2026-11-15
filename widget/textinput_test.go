package widget

import (
	"image/color"
	"math"
	"testing"
	"unicode/utf8"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/internal/input/inputtest"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/matryer/is"
)

// The caret before each rune sits where that rune is drawn, kerning against
// the rune before it included.
func TestTextInput_CaretAtDrawnGlyphs(t *testing.T) {
	is := is.New(t)

	face := loadFont(t)
	s := "To Wa AV"
	for _, g := range text.AppendLazyGlyphs(nil, s, *face, nil) {
		p := utf8.RuneCountInString(s[:g.StartIndexInBytes])
		is.Equal(fontAdvanceAt(s, p, face), int(math.Round(g.OriginX))) // caret before the rune
	}
}

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

// Typed characters go in, except while Control is held, as the input package
// reports the keys: then they belong to a shortcut.
func TestTextInput_TypingWithControlHeld(t *testing.T) {
	is := is.New(t)

	d := inputtest.Use(t)
	ti := newTextInput(t)
	ti.Focus(true)
	for i, s := range []struct {
		control bool
		text    string
	}{
		{true, ""},
		{false, "a"},
	} {
		d.Keys[ebiten.KeyControl] = s.control
		input.Update()
		typeChars(t, "a")
		ti.Update(&UpdateObject{})
		input.AfterUpdate()
		is.Equal([]any{i, ti.GetText()}, []any{i, s.text})
	}
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
