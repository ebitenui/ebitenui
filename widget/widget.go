package widget

import (
	"image"

	"github.com/blizzy78/ebitenui/event"
	"github.com/blizzy78/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
)

// A Widget is an abstraction of a user interface widget, such as a button. Actual widget implementations
// "have" a Widget in their internal structure.
type Widget struct {
	// Rect specifies the widget's position on screen. It is usually not set directly, but a Layouter is
	// used to set the position in relation to other widgets or the space available.
	Rect image.Rectangle

	// LayoutData specifies additional optional data for a Layouter that is used to layout this widget's
	// parent container. The exact type depends on the layout being used, for example, GridLayout requires
	// GridLayoutData to be used.
	LayoutData interface{}

	// Disabled specifies whether the widget is disabled, whatever that means. Disabled widgets should
	// usually render in some sort of "greyed out" visual state, and not react to user input.
	//
	// Not reacting to user input depends on the actual implementation. For example, List will not allow
	// entry selection via clicking, but the scrollbars will still be usable. The reasoning is that from
	// the user's perspective, scrolling does not change state, but only the display of that state.
	Disabled bool

	// CursorEnterEvent fires an event with *WidgetCursorEnterEventArgs when the cursor enters the widget's Rect.
	CursorEnterEvent *event.Event

	// CursorExitEvent fires an event with *WidgetCursorExitEventArgs when the cursor exits the widget's Rect.
	CursorExitEvent *event.Event

	// MouseButtonPressedEvent fires an event with *WidgetMouseButtonPressedEventArgs when a mouse button is pressed
	// while the cursor is inside the widget's Rect.
	MouseButtonPressedEvent *event.Event

	// MouseButtonReleasedEvent fires an event with *WidgetMouseButtonReleasedEventArgs when a mouse button is released
	// while the cursor is inside the widget's Rect.
	MouseButtonReleasedEvent *event.Event

	// ScrolledEvent fires an event with *WidgetScrolledEventArgs when the mouse wheel is scrolled while
	// the cursor is inside the widget's Rect.
	ScrolledEvent *event.Event

	FocusEvent *event.Event

	parent                     *Widget
	lastUpdateCursorEntered    bool
	lastUpdateMouseLeftPressed bool
	mouseLeftPressedInside     bool
	inputLayer                 *input.Layer
}

// WidgetOpt is a function that configures w.
type WidgetOpt func(w *Widget) //nolint:golint

// HasWidget must be implemented by concrete widget types to get their Widget.
type HasWidget interface {
	GetWidget() *Widget
}

// Renderer may be implemented by concrete widget types that can render onto the screen.
type Renderer interface {
	// Render renders the widget onto screen. def may be called to defer additional rendering.
	Render(screen *ebiten.Image, def DeferredRenderFunc)
}

type Focuser interface {
	Focus(focused bool)
}

// RenderFunc is a function that renders a widget onto screen. def may be called to defer
// additional rendering.
type RenderFunc func(screen *ebiten.Image, def DeferredRenderFunc)

// DeferredRenderFunc is a function that stores r for deferred execution.
type DeferredRenderFunc func(r RenderFunc)

// PreferredSizer may be implemented by concrete widget types that can report a preferred size.
type PreferredSizer interface {
	PreferredSize() (int, int)
}

// WidgetCursorEnterEventArgs are the arguments for cursor enter events.
type WidgetCursorEnterEventArgs struct { //nolint:golint
	Widget *Widget
}

// WidgetCursorExitEventArgs are the arguments for cursor exit events.
type WidgetCursorExitEventArgs struct { //nolint:golint
	Widget *Widget
}

// WidgetMouseButtonPressedEventArgs are the arguments for mouse button press events.
type WidgetMouseButtonPressedEventArgs struct { //nolint:golint
	Widget *Widget
	Button ebiten.MouseButton

	// OffsetX is the x offset relative to the widget's Rect.
	OffsetX int

	// OffsetY is the y offset relative to the widget's Rect.
	OffsetY int
}

// WidgetMouseButtonReleasedEventArgs are the arguments for mouse button release events.
type WidgetMouseButtonReleasedEventArgs struct { //nolint:golint
	Widget *Widget
	Button ebiten.MouseButton

	// Inside specifies whether the button has been released inside the widget's Rect.
	Inside bool

	// OffsetX is the x offset relative to the widget's Rect.
	OffsetX int

	// OffsetY is the y offset relative to the widget's Rect.
	OffsetY int
}

// WidgetScrolledEventArgs are the arguments for mouse wheel scroll events.
type WidgetScrolledEventArgs struct { //nolint:golint
	Widget *Widget
	X      float64
	Y      float64
}

type WidgetFocusEventArgs struct { //nolint:golint
	Widget  *Widget
	Focused bool
}

// WidgetCursorEnterHandlerFunc is a function that handles cursor enter events.
type WidgetCursorEnterHandlerFunc func(args *WidgetCursorEnterEventArgs) //nolint:golint

// WidgetCursorExitHandlerFunc is a function that handles cursor exit events.
type WidgetCursorExitHandlerFunc func(args *WidgetCursorExitEventArgs) //nolint:golint

// WidgetMouseButtonPressedHandlerFunc is a function that handles mouse button press events.
type WidgetMouseButtonPressedHandlerFunc func(args *WidgetMouseButtonPressedEventArgs) //nolint:golint

// WidgetMouseButtonReleasedHandlerFunc is a function that handles mouse button release events.
type WidgetMouseButtonReleasedHandlerFunc func(args *WidgetMouseButtonReleasedEventArgs) //nolint:golint

// WidgetScrolledHandlerFunc is a function that handles mouse wheel scroll events.
type WidgetScrolledHandlerFunc func(args *WidgetScrolledEventArgs) //nolint:golint

type WidgetOptions struct { //nolint:golint
}

// WidgetOpts contains functions that configure a Widget.
var WidgetOpts WidgetOptions

var deferredRenders []RenderFunc

// NewWidget constructs a new Widget configured with opts.
func NewWidget(opts ...WidgetOpt) *Widget {
	w := &Widget{
		CursorEnterEvent:         &event.Event{},
		CursorExitEvent:          &event.Event{},
		MouseButtonPressedEvent:  &event.Event{},
		MouseButtonReleasedEvent: &event.Event{},
		ScrolledEvent:            &event.Event{},
		FocusEvent:               &event.Event{},
	}

	for _, o := range opts {
		o(w)
	}

	return w
}

// WithLayoutData configures a Widget with layout data ld.
func (o WidgetOptions) LayoutData(ld interface{}) WidgetOpt {
	return func(w *Widget) {
		w.LayoutData = ld
	}
}

// WithCursorEnterHandler configures a Widget with cursor enter event handler f.
func (o WidgetOptions) CursorEnterHandler(f WidgetCursorEnterHandlerFunc) WidgetOpt {
	return func(w *Widget) {
		w.CursorEnterEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetCursorEnterEventArgs))
		})
	}
}

// WithCursorExitHandler configures a Widget with cursor exit event handler f.
func (o WidgetOptions) CursorExitHandler(f WidgetCursorExitHandlerFunc) WidgetOpt {
	return func(w *Widget) {
		w.CursorExitEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetCursorExitEventArgs))
		})
	}
}

// WithMouseButtonPressedHandler configures a Widget with mouse button press event handler f.
func (o WidgetOptions) MouseButtonPressedHandler(f WidgetMouseButtonPressedHandlerFunc) WidgetOpt {
	return func(w *Widget) {
		w.MouseButtonPressedEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetMouseButtonPressedEventArgs))
		})
	}
}

// WithMouseButtonReleasedHandler configures a Widget with mouse button release event handler f.
func (o WidgetOptions) MouseButtonReleasedHandler(f WidgetMouseButtonReleasedHandlerFunc) WidgetOpt {
	return func(w *Widget) {
		w.MouseButtonReleasedEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetMouseButtonReleasedEventArgs))
		})
	}
}

// WithScrolledHandler configures a Widget with mouse wheel scroll event handler f.
func (o WidgetOptions) ScrolledHandler(f WidgetScrolledHandlerFunc) WidgetOpt {
	return func(w *Widget) {
		w.ScrolledEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetScrolledEventArgs))
		})
	}
}

func (w *Widget) drawImageOptions(opts *ebiten.DrawImageOptions) {
	opts.GeoM.Translate(float64(w.Rect.Min.X), float64(w.Rect.Min.Y))
}

// EffectiveInputLayer returns w's effective input layer. If w does not have an input layer,
// or if the input layer is no longer valid, it returns w's parent widget's effective input layer.
// If w does not have a parent widget, it returns input.DefaultLayer.
func (w *Widget) EffectiveInputLayer() *input.Layer {
	l := w.inputLayer
	if l != nil && !l.Valid() {
		l = nil
	}

	if l == nil {
		if w.parent == nil {
			return &input.DefaultLayer
		}

		return w.parent.EffectiveInputLayer()
	}

	return l
}

// Render renders w onto screen. Since Widget is only an abstraction, it does not actually draw
// anything, but it is still responsible for firing events. Concrete widget implementations should
// always call this method first before rendering themselves.
func (w *Widget) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	w.fireEvents()
}

func (w *Widget) fireEvents() {
	x, y := input.CursorPosition()
	p := image.Point{x, y}
	layer := w.EffectiveInputLayer()
	inside := p.In(w.Rect)

	entered := inside && layer.ActiveFor(x, y, input.LayerEventTypeAny)
	if entered != w.lastUpdateCursorEntered {
		if entered {
			w.CursorEnterEvent.Fire(&WidgetCursorEnterEventArgs{
				Widget: w,
			})
		} else {
			w.CursorExitEvent.Fire(&WidgetCursorExitEventArgs{
				Widget: w,
			})
		}

		w.lastUpdateCursorEntered = entered
	}

	if inside && input.MouseButtonJustPressedLayer(ebiten.MouseButtonLeft, layer) {
		w.lastUpdateMouseLeftPressed = true
		w.mouseLeftPressedInside = inside

		off := p.Sub(w.Rect.Min)
		w.MouseButtonPressedEvent.Fire(&WidgetMouseButtonPressedEventArgs{
			Widget:  w,
			Button:  ebiten.MouseButtonLeft,
			OffsetX: off.X,
			OffsetY: off.Y,
		})
	}

	if w.lastUpdateMouseLeftPressed && !input.MouseButtonPressedLayer(ebiten.MouseButtonLeft, layer) {
		w.lastUpdateMouseLeftPressed = false

		off := p.Sub(w.Rect.Min)
		w.MouseButtonReleasedEvent.Fire(&WidgetMouseButtonReleasedEventArgs{
			Widget:  w,
			Button:  ebiten.MouseButtonLeft,
			Inside:  inside,
			OffsetX: off.X,
			OffsetY: off.Y,
		})
	}

	scrollX, scrollY := input.WheelLayer(layer)
	if inside && (scrollX != 0 || scrollY != 0) {
		w.ScrolledEvent.Fire(&WidgetScrolledEventArgs{
			Widget: w,
			X:      scrollX,
			Y:      scrollY,
		})
	}
}

// SetLocation sets w's position to rect. This is usually not called directly, but by a layout.
func (w *Widget) SetLocation(rect image.Rectangle) {
	w.Rect = rect
}

// ElevateToNewInputLayer adds l to the top of the input layer stack, then sets w's input layer to l.
func (w *Widget) ElevateToNewInputLayer(l *input.Layer) {
	input.AddLayer(l)
	w.inputLayer = l
}

func (w *Widget) Parent() *Widget {
	return w.parent
}

func WidgetFireFocusEvent(w *Widget, focused bool) { //nolint:golint
	w.FocusEvent.Fire(&WidgetFocusEventArgs{
		Widget:  w,
		Focused: focused,
	})
}

// RenderWithDeferred renders r to screen. This function should not be called directly.
func RenderWithDeferred(screen *ebiten.Image, rs []Renderer) {
	for _, r := range rs {
		appendToDeferredRenderQueue(r.Render)
	}

	renderDeferredRenderQueue(screen)
}

func renderDeferredRenderQueue(screen *ebiten.Image) {
	defer func(d []RenderFunc) {
		deferredRenders = d[:0]
	}(deferredRenders)

	for len(deferredRenders) > 0 {
		r := deferredRenders[0]
		deferredRenders = deferredRenders[1:]

		r(screen, appendToDeferredRenderQueue)
	}
}

func appendToDeferredRenderQueue(r RenderFunc) {
	deferredRenders = append(deferredRenders, r)
}
