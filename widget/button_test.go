package widget

import (
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/matryer/is"
)

func TestButton_PressedEvent_User(t *testing.T) {
	is := is.New(t)

	var eventArgs *ButtonPressedEventArgs

	b := newButton(t,
		ButtonOpts.PressedHandler(func(args *ButtonPressedEventArgs) {
			eventArgs = args
		}))

	leftMouseButtonPress(b, t)

	is.True(eventArgs != nil)
}

func TestButton_ReleasedEvent_User(t *testing.T) {
	is := is.New(t)

	var eventArgs *ButtonReleasedEventArgs

	b := newButton(t,
		ButtonOpts.ReleasedHandler(func(args *ButtonReleasedEventArgs) {
			eventArgs = args
		}))

	leftMouseButtonPress(b, t)
	leftMouseButtonRelease(b, t)

	is.True(eventArgs != nil)
}

func TestButton_ClickedEvent_User(t *testing.T) {
	is := is.New(t)

	var eventArgs *ButtonClickedEventArgs

	b := newButton(t,
		ButtonOpts.ClickedHandler(func(args *ButtonClickedEventArgs) {
			eventArgs = args
		}))

	leftMouseButtonClick(b, t)

	is.True(eventArgs != nil)
}

func TestButton_ManualClick(t *testing.T) {
	is := is.New(t)

	var eventArgs *ButtonClickedEventArgs

	b := newButton(t,
		ButtonOpts.ClickedHandler(func(args *ButtonClickedEventArgs) {
			eventArgs = args
		}))

	b.Click()
	event.ExecuteDeferred()

	is.True(eventArgs != nil)
}

func TestButton_PreferredSizeBeforeValidation(t *testing.T) {
	is := is.New(t)

	buttonImage := &ButtonImage{
		Idle:    newNineSliceEmpty(t),
		Pressed: newNineSliceEmpty(t),
	}
	b := NewButton(
		ButtonOpts.Image(buttonImage),
		ButtonOpts.WidgetOpts(WidgetOpts.MinSize(60, 70)),
	)
	c := NewContainer()
	c.AddChild(b)

	w, h := b.PreferredSize()
	is.Equal(w, 60)
	is.Equal(h, 70)
}

func newButton(t *testing.T, opts ...ButtonOpt) *Button {
	t.Helper()

	b := NewButton(append(opts, ButtonOpts.Image(&ButtonImage{
		Idle:    newNineSliceEmpty(t),
		Pressed: newNineSliceEmpty(t),
	}))...)
	event.ExecuteDeferred()
	render(b, t)
	return b
}
