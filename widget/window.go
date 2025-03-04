package widget

import (
	"image"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
)

type WindowCloseMode int

const (
	// The window will not automatically close.
	NONE WindowCloseMode = iota
	// The window will close when you click anywhere.
	CLICK
	// The window will close when you click outside the window.
	CLICK_OUT
)

type RemoveWindowFunc func()

type WindowChangedEventArgs struct {
	Window *Window
	Rect   image.Rectangle
}

type WindowChangedHandlerFunc func(args *WindowChangedEventArgs)

type WindowClosedEventArgs struct {
	Window *Window
}

type WindowClosedHandlerFunc func(args *WindowClosedEventArgs)

type Window struct {
	ResizeEvent *event.Event
	MoveEvent   *event.Event
	ClosedEvent *event.Event

	Modal      bool
	Contents   *Container
	TitleBar   *Container
	Draggable  bool
	Resizeable bool
	MinSize    *image.Point
	MaxSize    *image.Point
	DrawLayer  int
	// Used to indicate this window should close if other windows close.
	Ephemeral bool

	// INTERNAL USE ONLY: Used to indicate that this window is the currently focused window.
	FocusedWindow bool
	// Sets whether the window should move to the top of the draw order when clicked
	// Default: false.
	DisableRelayering bool

	closeMode WindowCloseMode
	closeFunc RemoveWindowFunc
	container *Container

	titleBarHeight int

	startingPoint  image.Point
	dragging       bool
	resizing       bool
	resizingWidth  bool
	resizingHeight bool
	blockLower     bool
	originalSize   image.Point
	init           *MultiOnce
}

type WindowOpt func(w *Window)

type WindowOptions struct {
}

var WindowOpts WindowOptions

func NewWindow(opts ...WindowOpt) *Window {
	w := &Window{
		MoveEvent:   &event.Event{},
		ResizeEvent: &event.Event{},
		ClosedEvent: &event.Event{},
		init:        &MultiOnce{},
		blockLower:  true,
	}
	w.init.Append(w.createWidget)
	for _, o := range opts {
		o(w)
	}

	w.validate()

	w.init.Do()
	return w
}

func (w *Window) validate() {
	if w.Contents == nil {
		panic("Window: Contents is required.")
	}
}

// This is the container with the body of this window.
func (o WindowOptions) Contents(c *Container) WindowOpt {
	return func(w *Window) {
		w.Contents = c
	}
}

// Sets the container for the TitleBar and its fixed height.
func (o WindowOptions) TitleBar(tb *Container, height int) WindowOpt {
	return func(w *Window) {
		w.TitleBar = tb
		w.titleBarHeight = height
	}
}

// Sets the window to be modal. Blocking UI interactions on anything else.
func (o WindowOptions) Modal() WindowOpt {
	return func(w *Window) {
		w.Modal = true
	}
}

// Sets the window to be draggable. The handle for this is the titleBar.
// If you haven't provided a titleBar this option is ignored.
func (o WindowOptions) Draggable() WindowOpt {
	return func(w *Window) {
		w.Draggable = true
	}
}

// Sets the window to be resizeable.
func (o WindowOptions) Resizeable() WindowOpt {
	return func(w *Window) {
		w.Resizeable = true
	}
}

// Sets whether the window should block input beneath the window or not.
// Default: true.
func (o WindowOptions) BlockLower(blockLower bool) WindowOpt {
	return func(w *Window) {
		w.blockLower = blockLower
	}
}

// Sets whether the window should move to the top of the draw order when clicked
// Default: false.
func (o WindowOptions) DisableRelayering(disableRelayering bool) WindowOpt {
	return func(w *Window) {
		w.DisableRelayering = disableRelayering
	}
}

// Sets the minimum size that the window can be reszied to.
func (o WindowOptions) MinSize(width int, height int) WindowOpt {
	return func(w *Window) {
		w.MinSize = &image.Point{X: width, Y: height}
	}
}

// Set the maximum size that the window can be resized to.
func (o WindowOptions) MaxSize(width int, height int) WindowOpt {
	return func(w *Window) {
		w.MaxSize = &image.Point{X: width, Y: height}
	}
}

// Set the way this window should close.
func (o WindowOptions) CloseMode(mode WindowCloseMode) WindowOpt {
	return func(w *Window) {
		w.closeMode = mode
	}
}

// This handler is triggered when a move event is completed.
func (o WindowOptions) MoveHandler(f WindowChangedHandlerFunc) WindowOpt {
	return func(w *Window) {
		w.MoveEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*WindowChangedEventArgs); ok {
				f(arg)
			}
		})
	}
}

// This handler is triggered when a resize event is completed.
func (o WindowOptions) ResizeHandler(f WindowChangedHandlerFunc) WindowOpt {
	return func(w *Window) {
		w.ResizeEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*WindowChangedEventArgs); ok {
				f(arg)
			}
		})
	}
}

// This handler is triggered when window is closed.
// The window can be closed by either calling [Window.Close] or
// by invoking an [UI.AddWindow] returned close function.
//
// This handler is called after the window is closed.
// The provided Window object is still accessible, but the window
// is already removed from UI.
func (o WindowOptions) ClosedHandler(f WindowClosedHandlerFunc) WindowOpt {
	return func(w *Window) {
		w.ClosedEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*WindowClosedEventArgs); ok {
				f(arg)
			}
		})
	}
}

// This option sets the size and location of the window.
// This method will account for specified MinSize and MaxSize values.
func (o WindowOptions) Location(rect image.Rectangle) WindowOpt {
	return func(w *Window) {
		w.init.Append(func() { w.container.SetLocation(rect) })
	}
}

// This option sets order the window will be drawn
//   - &lt; 0 will have the window drawn before the container
//   - &gt;= 0 will have the window drawn after the container
func (o WindowOptions) DrawLayer(layer int) WindowOpt {
	return func(w *Window) {
		w.DrawLayer = layer
	}
}

// This method is used to be able to close the window.
func (w *Window) Close() {
	if w.closeFunc != nil {
		w.closeFunc()
	}
}

// This method will set the size and location of this window.
// This method will account for specified MinSize and MaxSize values.
func (w *Window) SetLocation(rect image.Rectangle) {
	if rect != w.container.widget.Rect {
		if w.MinSize != nil {
			if rect.Dx() < w.MinSize.X {
				rect.Max.X = rect.Min.X + w.MinSize.X
			}
			if rect.Dy() < w.MinSize.Y {
				rect.Max.Y = rect.Min.Y + w.MinSize.Y
			}
		}

		if w.MaxSize != nil {
			if rect.Dx() > w.MaxSize.X {
				rect.Max.X = rect.Min.X + w.MaxSize.X
			}
			if rect.Dy() > w.MaxSize.Y {
				rect.Max.Y = rect.Min.Y + w.MaxSize.Y
			}
		}

		w.container.SetLocation(rect)
	}
}

// Typically used internally.
//
//	Returns the root container that holds the provided titlebar and contents.
func (w *Window) GetContainer() *Container {
	return w.container
}

// Typically used internally.
func (w *Window) SetCloseFunction(removeWindowFunc RemoveWindowFunc) {
	w.closeFunc = func() {
		removeWindowFunc()
		w.ClosedEvent.Fire(&WindowClosedEventArgs{
			Window: w,
		})
	}
}

// Typically used internally.
func (w *Window) GetCloseFunction() RemoveWindowFunc {
	return w.closeFunc
}

// Typically used internally.
func (w *Window) RequestRelayout() {
	w.container.RequestRelayout()
}

// Typically used internally.
func (w *Window) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	w.container.GetWidget().ElevateToNewInputLayer(&input.Layer{
		DebugLabel: "window",
		EventTypes: input.LayerEventTypeAll,
		BlockLower: w.blockLower && !w.Ephemeral,
		FullScreen: w.Modal,
		RectFunc: func() image.Rectangle {
			return w.container.GetWidget().Rect
		},
	})
	w.container.SetupInputLayer(def)
}

// Typically used internally.
func (w *Window) Render(screen *ebiten.Image) {
	x, y := input.CursorPosition()

	if w.dragging {
		if w.startingPoint.X != x || w.startingPoint.Y != y {
			newRect := w.container.GetWidget().Rect.Add(image.Point{x - w.startingPoint.X, y - w.startingPoint.Y})
			w.SetLocation(newRect)
			w.startingPoint = image.Point{x, y}
		}
	}
	if w.resizing {
		if w.startingPoint.X != x || w.startingPoint.Y != y {
			if w.resizingWidth {
				newRect := w.container.GetWidget().Rect
				newRect.Max.X = w.originalSize.X - (w.startingPoint.X - x)
				w.SetLocation(newRect)
			}
			if w.resizingHeight {
				newRect := w.container.GetWidget().Rect
				newRect.Max.Y = w.originalSize.Y - (w.startingPoint.Y - y)
				w.SetLocation(newRect)
			}
		}
	}

	if w.Resizeable {
		if w.container.GetWidget().inputLayer.ActiveFor(x, y, input.LayerEventTypeAll) {
			xRect := image.Rect(w.container.GetWidget().Rect.Max.X-6, w.container.GetWidget().Rect.Min.Y, w.container.GetWidget().Rect.Max.X, w.container.GetWidget().Rect.Max.Y)
			yRect := image.Rect(w.container.GetWidget().Rect.Min.X, w.container.GetWidget().Rect.Max.Y-6, w.container.GetWidget().Rect.Max.X, w.container.GetWidget().Rect.Max.Y)
			cursorRect := image.Rect(x, y, x+1, y+1)
			switch {
			case cursorRect.Overlaps(xRect):
				input.SetCursorShape(input.CURSOR_EWRESIZE)
				w.resizingWidth = true
				w.resizingHeight = false
			case cursorRect.Overlaps(yRect):
				input.SetCursorShape(input.CURSOR_NSRESIZE)
				w.resizingWidth = false
				w.resizingHeight = true
			case !input.MouseButtonPressed(ebiten.MouseButtonLeft):
				w.resizingWidth = false
				w.resizingHeight = false
			}
		} else if !input.MouseButtonPressed(ebiten.MouseButtonLeft) {
			w.resizingWidth = false
			w.resizingHeight = false
		}
	}
	w.container.Render(screen)
}

func (w *Window) Update() {
	w.init.Do()
	w.container.Update()
}

func (w *Window) createWidget() {
	if w.TitleBar != nil {
		w.container = NewContainer(
			ContainerOpts.Layout(NewGridLayout(
				GridLayoutOpts.Columns(1),
				GridLayoutOpts.Stretch([]bool{true}, []bool{false, true}),
			)),
			ContainerOpts.WidgetOpts(WidgetOpts.TrackHover(true)),
		)
		w.TitleBar.GetWidget().LayoutData = GridLayoutData{MaxHeight: w.titleBarHeight}
		w.TitleBar.GetWidget().MinHeight = w.titleBarHeight
		if w.Draggable {
			w.TitleBar.GetWidget().MouseButtonPressedEvent.AddHandler(func(a any) {
				if args, ok := a.(*WidgetMouseButtonPressedEventArgs); ok {
					if args.Button == ebiten.MouseButtonLeft {
						x, y := input.CursorPosition()
						w.startingPoint = image.Point{x, y}
						w.dragging = true
					}
				}
			})
			w.TitleBar.GetWidget().MouseButtonReleasedEvent.AddHandler(func(a any) {
				if args, ok := a.(*WidgetMouseButtonReleasedEventArgs); ok {
					if w.dragging && args.Button == ebiten.MouseButtonLeft {
						w.dragging = false
						w.MoveEvent.Fire(&WindowChangedEventArgs{
							Window: w,
							Rect:   w.container.GetWidget().Rect,
						})
					}
				}
			})
		}
		w.container.AddChild(w.TitleBar)
		w.container.AddChild(w.Contents)
	} else {
		w.container = NewContainer(
			ContainerOpts.Layout(NewGridLayout(
				GridLayoutOpts.Columns(1),
				GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
			)),
			ContainerOpts.WidgetOpts(WidgetOpts.TrackHover(true)),
		)

		w.container.AddChild(w.Contents)
	}

	if w.Resizeable {
		w.Contents.GetWidget().MouseButtonPressedEvent.AddHandler(func(a any) {
			if args, ok := a.(*WidgetMouseButtonPressedEventArgs); ok {
				if args.Button == ebiten.MouseButtonLeft {
					x, y := input.CursorPosition()
					w.startingPoint = image.Point{x, y}
					w.originalSize.X = w.container.GetWidget().Rect.Max.X
					w.originalSize.Y = w.container.GetWidget().Rect.Max.Y
					w.resizing = true
				}
			}
		})
		w.Contents.GetWidget().MouseButtonReleasedEvent.AddHandler(func(a any) {
			if args, ok := a.(*WidgetMouseButtonReleasedEventArgs); ok {
				if w.resizing && args.Button == ebiten.MouseButtonLeft {
					w.resizing = false
					w.ResizeEvent.Fire(&WindowChangedEventArgs{
						Window: w,
						Rect:   w.container.GetWidget().Rect,
					})
				}
			}
		})
	}

	if w.closeMode == CLICK || w.closeMode == CLICK_OUT {
		w.container.GetWidget().CustomData = "Window"
		w.container.GetWidget().MouseButtonReleasedEvent.AddHandler(func(args interface{}) {
			if a, ok := args.(*WidgetMouseButtonReleasedEventArgs); ok {
				if w.closeMode == CLICK || (w.closeMode == CLICK_OUT && !a.Inside) {
					if w.closeFunc != nil {
						w.closeFunc()
					}
				}
			}
		})
	}

	w.container.GetWidget().MouseButtonPressedEvent.AddHandler(func(_ any) {
		if !w.DisableRelayering {
			w.FocusedWindow = true
		}
	})

}
