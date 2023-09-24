package widget

import (
	"image"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"

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

	//The minimum width for this Widget
	MinWidth int

	//The minimum height for this Widget
	MinHeight int

	// Disabled specifies whether the widget is disabled, whatever that means. Disabled widgets should
	// usually render in some sort of "greyed out" visual state, and not react to user input.
	//
	// Not reacting to user input depends on the actual implementation. For example, List will not allow
	// entry selection via clicking, but the scrollbars will still be usable. The reasoning is that from
	// the user's perspective, scrolling does not change state, but only the display of that state.
	Disabled bool

	// Hidden specifies whether the widget is visible. Hidden widgets should
	// not render anything or react to user input.
	Visibility Visibility

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

	ContextMenuEvent *event.Event

	ToolTipEvent *event.Event

	DragAndDropEvent *event.Event

	//Custom Data is a field to allow users to attach data to any widget
	CustomData interface{}

	canDrop CanDropFunc
	drop    DropFunc

	parent                      *Widget
	self                        HasWidget
	lastUpdateCursorEntered     bool
	lastUpdateMouseLeftPressed  bool
	lastUpdateMouseRightPressed bool
	mouseLeftPressedInside      bool
	inputLayer                  *input.Layer
	focusable                   Focuser

	ContextMenu          *Container
	ContextMenuWindow    *Window
	ContextMenuCloseMode WindowCloseMode

	ToolTip *ToolTip

	DragAndDrop *DragAndDrop

	CursorHovered string
	CursorPressed string
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
	IsFocused() bool
	TabOrder() int
}

type Dropper interface {
	GetDropTargets() []HasWidget
}

type Visibility int

const (
	Visibility_Show          Visibility = iota
	Visibility_Hide_Blocking            // Hide widget, but take up space
	Visibility_Hide                     // Hide widget, but don't take up space
)

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
	Widget   HasWidget
	Focused  bool
	Location image.Point
}

type WidgetContextMenuEventArgs struct { //nolint:golint
	Widget   *Widget
	Location image.Point
}

type WidgetToolTipEventArgs struct { //nolint:golint
	Window *Window
	Show   bool
}

type WidgetDragAndDropEventArgs struct { //nolint:golint
	Window *Window
	Show   bool
	DnD    *DragAndDrop
}

type DragAndDropDroppedEventArgs struct { //nolint:golint
	Source  HasWidget
	SourceX int
	SourceY int
	Target  HasWidget
	TargetX int
	TargetY int
	Data    interface{}
}

type CanDropFunc func(args *DragAndDropDroppedEventArgs) bool
type DropFunc func(args *DragAndDropDroppedEventArgs)

type DragAndDropDroppedHandlerFunc func(args *DragAndDropDroppedEventArgs)

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
		ContextMenuEvent:         &event.Event{},
		ToolTipEvent:             &event.Event{},
		DragAndDropEvent:         &event.Event{},
		ContextMenuCloseMode:     CLICK,
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

// WithCustomData configures a Widget with custom data cd.
func (o WidgetOptions) CustomData(cd interface{}) WidgetOpt {
	return func(w *Widget) {
		w.CustomData = cd
	}
}

func (o WidgetOptions) MinSize(minWidth int, minHeight int) WidgetOpt {
	return func(w *Widget) {
		w.MinWidth = minWidth
		w.MinHeight = minHeight
	}
}

func (o WidgetOptions) ContextMenu(contextMenu *Container) WidgetOpt {
	return func(w *Widget) {
		w.ContextMenu = contextMenu
	}
}

func (o WidgetOptions) ContextMenuCloseMode(contextMenuCloseMode WindowCloseMode) WidgetOpt {
	return func(w *Widget) {
		w.ContextMenuCloseMode = contextMenuCloseMode
	}
}

func (o WidgetOptions) ToolTip(toolTip *ToolTip) WidgetOpt {
	return func(w *Widget) {
		w.ToolTip = toolTip
	}
}

// This sets the source of a Drag and Drop action
func (o WidgetOptions) EnableDragAndDrop(d *DragAndDrop) WidgetOpt {
	return func(w *Widget) {
		w.DragAndDrop = d
	}
}

// This sets the widget as a target of a Drag and Drop action
//
//	The Drop function must return true if it accepts this drop and false if it does not accept the drop
func (o WidgetOptions) CanDrop(candropFunc CanDropFunc) WidgetOpt {
	return func(w *Widget) {
		w.canDrop = candropFunc
	}
}

// This is the function that is run if an item is dropped on this widget
func (o WidgetOptions) Dropped(dropFunc DropFunc) WidgetOpt {
	return func(w *Widget) {
		w.drop = dropFunc
	}
}

func (o WidgetOptions) CursorHovered(cursorHovered string) WidgetOpt {
	return func(w *Widget) {
		w.CursorHovered = cursorHovered
	}
}

func (o WidgetOptions) CursorPressed(cursorPressed string) WidgetOpt {
	return func(w *Widget) {
		w.CursorPressed = cursorPressed
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

// always call this method first before rendering themselves.
func (w *Widget) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	w.fireEvents()
	if w.ToolTip != nil {
		w.ToolTip.Render(w, screen, def)
	}
	if w.DragAndDrop != nil {
		w.DragAndDrop.Render(w.self, screen, def)
	}
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

	if entered && len(w.CursorHovered) > 0 {
		input.SetCursorShape(w.CursorHovered)
	}

	if entered && len(w.CursorPressed) > 0 && input.MouseButtonPressedLayer(ebiten.MouseButtonLeft, layer) {
		input.SetCursorShape(w.CursorPressed)
	}

	if input.MouseButtonJustPressedLayer(ebiten.MouseButtonRight, layer) {
		w.lastUpdateMouseRightPressed = true
		if inside {
			off := p.Sub(w.Rect.Min)
			w.MouseButtonPressedEvent.Fire(&WidgetMouseButtonPressedEventArgs{
				Widget:  w,
				Button:  ebiten.MouseButtonRight,
				OffsetX: off.X,
				OffsetY: off.Y,
			})
			if w.ContextMenu != nil {
				w.FireContextMenuEvent(nil, p)
			}
		}
	}

	if w.lastUpdateMouseRightPressed && !input.MouseButtonPressed(ebiten.MouseButtonRight) {
		w.lastUpdateMouseRightPressed = false

		off := p.Sub(w.Rect.Min)
		w.MouseButtonReleasedEvent.Fire(&WidgetMouseButtonReleasedEventArgs{
			Widget:  w,
			Button:  ebiten.MouseButtonRight,
			Inside:  inside,
			OffsetX: off.X,
			OffsetY: off.Y,
		})
	}

	if input.MouseButtonJustPressed(ebiten.MouseButtonLeft) {
		w.lastUpdateMouseLeftPressed = true
	}

	if inside && input.MouseButtonJustPressedLayer(ebiten.MouseButtonLeft, layer) {
		if inside {
			w.mouseLeftPressedInside = inside

			if w.focusable != nil && !w.Disabled {
				w.focusable.Focus(true)
			} else {
				w.FireFocusEvent(nil, false, p)
			}

			off := p.Sub(w.Rect.Min)
			w.MouseButtonPressedEvent.Fire(&WidgetMouseButtonPressedEventArgs{
				Widget:  w,
				Button:  ebiten.MouseButtonLeft,
				OffsetX: off.X,
				OffsetY: off.Y,
			})
		}
	}

	if w.lastUpdateMouseLeftPressed && !input.MouseButtonPressed(ebiten.MouseButtonLeft) {
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

func (widget *Widget) FireFocusEvent(w HasWidget, focused bool, location image.Point) { //nolint:golint
	widget.FocusEvent.Fire(&WidgetFocusEventArgs{
		Widget:   w,
		Focused:  focused,
		Location: location,
	})
}

func (widget *Widget) FireContextMenuEvent(w *Widget, l image.Point) { //nolint:golint
	if w == nil {
		w = widget
	}
	if w.ContextMenu != nil {
		widget.ContextMenuEvent.Fire(&WidgetContextMenuEventArgs{
			Widget:   w,
			Location: l,
		})
	}
}

func (widget *Widget) FireToolTipEvent(w *Window, show bool) { //nolint:golint
	widget.ToolTipEvent.Fire(&WidgetToolTipEventArgs{
		Window: w,
		Show:   show,
	})
}

func (widget *Widget) FireDragAndDropEvent(w *Window, show bool, dnd *DragAndDrop) { //nolint:golint
	widget.DragAndDropEvent.Fire(&WidgetDragAndDropEventArgs{
		Window: w,
		Show:   show,
		DnD:    dnd,
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
