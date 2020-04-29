package widget

import (
	"image"

	"github.com/blizzy78/ebitenui/event"
	"github.com/blizzy78/ebitenui/input"

	"github.com/hajimehoshi/ebiten"
)

type Widget struct {
	Rect       image.Rectangle
	LayoutData interface{}
	Disabled   bool
	InputLayer *input.Layer

	CursorEnterEvent         *event.Event
	CursorExitEvent          *event.Event
	MouseButtonPressedEvent  *event.Event
	MouseButtonReleasedEvent *event.Event
	ScrolledEvent            *event.Event

	parent                     *Widget
	lastUpdateCursorEntered    bool
	lastUpdateMouseLeftPressed bool
	mouseLeftPressedInside     bool
}

type WidgetOpt func(w *Widget)

type HasWidget interface {
	GetWidget() *Widget
}

type Renderer interface {
	Render(screen *ebiten.Image, def DeferredRenderFunc)
}

type DeferredRenderFunc func(r RenderFunc)

type RenderFunc func(screen *ebiten.Image, def DeferredRenderFunc)

type PreferredSizer interface {
	PreferredSize() (width int, height int)
}

type WidgetCursorEnterEventArgs struct {
	Widget *Widget
}

type WidgetCursorExitEventArgs struct {
	Widget *Widget
}

type WidgetMouseButtonPressedEventArgs struct {
	Widget  *Widget
	Button  ebiten.MouseButton
	OffsetX int
	OffsetY int
}

type WidgetMouseButtonReleasedEventArgs struct {
	Widget  *Widget
	Button  ebiten.MouseButton
	Inside  bool
	OffsetX int
	OffsetY int
}

type WidgetScrolledEventArgs struct {
	Widget *Widget
	X      float64
	Y      float64
}

type WidgetCursorEnterHandlerFunc func(args *WidgetCursorEnterEventArgs)

type WidgetCursorExitHandlerFunc func(args *WidgetCursorExitEventArgs)

type WidgetMouseButtonPressedHandlerFunc func(args *WidgetMouseButtonPressedEventArgs)

type WidgetMouseButtonReleasedHandlerFunc func(args *WidgetMouseButtonReleasedEventArgs)

type WidgetScrolledHandlerFunc func(args *WidgetScrolledEventArgs)

const WidgetOpts = widgetOpts(true)

type widgetOpts bool

var deferredRenders []RenderFunc

func NewWidget(opts ...WidgetOpt) *Widget {
	w := &Widget{
		CursorEnterEvent:         &event.Event{},
		CursorExitEvent:          &event.Event{},
		MouseButtonPressedEvent:  &event.Event{},
		MouseButtonReleasedEvent: &event.Event{},
		ScrolledEvent:            &event.Event{},
	}

	for _, o := range opts {
		o(w)
	}

	return w
}

func (o widgetOpts) WithLayoutData(ld interface{}) WidgetOpt {
	return func(w *Widget) {
		w.LayoutData = ld
	}
}

func (o widgetOpts) WithCursorEnterHandler(f WidgetCursorEnterHandlerFunc) WidgetOpt {
	return func(w *Widget) {
		w.CursorEnterEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetCursorEnterEventArgs))
		})
	}
}

func (o widgetOpts) WithCursorExitHandler(f WidgetCursorExitHandlerFunc) WidgetOpt {
	return func(w *Widget) {
		w.CursorExitEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetCursorExitEventArgs))
		})
	}
}

func (o widgetOpts) WithMouseButtonPressedHandler(f WidgetMouseButtonPressedHandlerFunc) WidgetOpt {
	return func(w *Widget) {
		w.MouseButtonPressedEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetMouseButtonPressedEventArgs))
		})
	}
}

func (o widgetOpts) WithMouseButtonReleasedHandler(f WidgetMouseButtonReleasedHandlerFunc) WidgetOpt {
	return func(w *Widget) {
		w.MouseButtonReleasedEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetMouseButtonReleasedEventArgs))
		})
	}
}

func (o widgetOpts) WithScrolledHandler(f WidgetScrolledHandlerFunc) WidgetOpt {
	return func(w *Widget) {
		w.ScrolledEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetScrolledEventArgs))
		})
	}
}

func (w *Widget) drawImageOptions(opts *ebiten.DrawImageOptions) {
	opts.GeoM.Translate(float64(w.Rect.Min.X), float64(w.Rect.Min.Y))
}

func (w *Widget) EffectiveInputLayer() *input.Layer {
	l := w.InputLayer
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
	if scrollX != 0 || scrollY != 0 {
		if inside {
			w.ScrolledEvent.Fire(&WidgetScrolledEventArgs{
				Widget: w,
				X:      scrollX,
				Y:      scrollY,
			})
		}
	}
}

func (w *Widget) SetLocation(rect image.Rectangle) {
	w.Rect = rect
}

func (w *Widget) ElevateToNewInputLayer(l *input.Layer) {
	input.AddLayer(l)
	w.InputLayer = l
}

func RenderWithDeferred(r Renderer, screen *ebiten.Image) {
	appendToDeferredRenderQueue(r.Render)
	renderDeferredRenderQueue(screen)
}

func renderDeferredRenderQueue(screen *ebiten.Image) {
	for len(deferredRenders) > 0 {
		r := deferredRenders[0]
		deferredRenders = deferredRenders[1:]

		r(screen, appendToDeferredRenderQueue)
	}
}

func appendToDeferredRenderQueue(r RenderFunc) {
	deferredRenders = append(deferredRenders, r)
}
