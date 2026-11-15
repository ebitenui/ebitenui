package widget

import (
	img "image"
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
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

// placedWidget counts the times a layout places it.
type placedWidget struct {
	widget *Widget
	placed int
}

func (p *placedWidget) GetWidget() *Widget        { return p.widget }
func (p *placedWidget) PreferredSize() (int, int) { return 10, 10 }
func (p *placedWidget) Validate()                 {}

func (p *placedWidget) SetLocation(rect img.Rectangle) {
	p.widget.Rect = rect
	p.placed++
}

// A button lays its content out when it moves or resizes, or when asked to
// (a window's relayout asks every widget in it), not on every frame.
func TestButton_ContentLaidOutWhenNeeded(t *testing.T) {
	is := is.New(t)

	content := &placedWidget{widget: NewWidget()}
	b := newButton(t, ButtonOpts.Content(content))
	screen := ebiten.NewImage(100, 100)
	b.SetLocation(img.Rect(0, 0, 40, 20))
	b.Render(screen)
	placed, at := content.placed, content.widget.Rect

	b.Render(screen)
	b.Render(screen)
	is.Equal(content.placed, placed) // the button stayed put

	b.SetLocation(img.Rect(10, 5, 50, 25))
	b.Render(screen)
	is.Equal(content.placed, placed+1)                      // laid out again
	is.Equal(content.widget.Rect, at.Add(img.Point{10, 5})) // following the button

	b.RequestRelayout()
	b.Render(screen)
	is.Equal(content.placed, placed+2) // asked to
}
