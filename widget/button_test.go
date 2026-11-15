package widget

import (
	img "image"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
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

// A widget in a button's content is under the button: it takes the
// button's input layer, and the theme of what holds the button.
func TestButton_ContentUnderTheButton(t *testing.T) {
	is := is.New(t)

	inner := NewButton(ButtonOpts.Image(&ButtonImage{Idle: newNineSliceEmpty(t), Pressed: newNineSliceEmpty(t)}))
	b := NewButton(ButtonOpts.Image(&ButtonImage{Idle: newNineSliceEmpty(t), Pressed: newNineSliceEmpty(t)}), ButtonOpts.Content(inner))
	theme := &Theme{}
	holder := NewContainer()
	holder.GetWidget().SetTheme(theme)
	holder.AddChild(b)
	holder.Validate()
	layer := &input.Layer{DebugLabel: "button", EventTypes: input.LayerEventTypeAll, RectFunc: func() img.Rectangle { return b.GetWidget().Rect }}
	t.Cleanup(func() { input.SetupInputLayersWithDeferred(nil) }) // the layers are the input package's
	b.GetWidget().ElevateToNewInputLayer(layer)

	is.Equal(inner.GetWidget().EffectiveInputLayer(), layer) // the button's layer
	is.Equal(inner.GetWidget().GetTheme(), theme)            // its holder's theme
}
