package widget

import (
	img "image"
	"math"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type ScrollContainer struct {
	ScrollLeft float64
	ScrollTop  float64

	// ScrollMovedEvent fires when a drag, or the motion left after one,
	// scrolls the content, with ScrollContainerScrollMovedEventArgs.
	ScrollMovedEvent *event.Event

	widgetOpts          []WidgetOpt
	definedImage        *ScrollContainerImage // given; image is it or the theme's
	image               *ScrollContainerImage
	content             PreferredSizeLocateableWidget
	forwarding          PreferredSizeLocateableWidget // the content whose events it forwards
	padding             *Insets
	stretchContentWidth bool
	dragToScroll        bool
	drag                scrollDrag

	init      *MultiOnce
	widget    *Widget
	validated bool
	renderBuf *image.MaskedRenderBuffer
}

type ScrollContainerOpt func(s *ScrollContainer)

// ScrollContainerScrollMovedEventArgs are the arguments of ScrollMovedEvent:
// where the content is scrolled to, as ScrollLeft and ScrollTop hold it.
type ScrollContainerScrollMovedEventArgs struct {
	ScrollContainer *ScrollContainer
	Left, Top       float64
}

const (
	dragThreshold = 8    // pixels a press moves before it drags the content
	dragFriction  = 0.95 // the share of its speed the content keeps each update after a drag
	dragStop      = 0.5  // pixels an update below which the content stops
)

// scrollDrag is a drag of the content in progress, or the motion left
// after one.
type scrollDrag struct {
	pressed, dragging   bool
	startX, startY      int     // where the press began
	startLeft, startTop float64 // the scroll then, in pixels
	lastX, lastY        int     // the pointer on the last update
	vx, vy              float64 // the scroll's speed, in pixels an update
}

type ScrollContainerImage struct {
	Idle     *image.NineSlice
	Disabled *image.NineSlice
	Mask     *image.NineSlice
}

type ScrollContainerOptions struct {
}

var ScrollContainerOpts ScrollContainerOptions

func NewScrollContainer(opts ...ScrollContainerOpt) *ScrollContainer {
	s := &ScrollContainer{
		ScrollMovedEvent: &event.Event{},

		init: &MultiOnce{},

		renderBuf: image.NewMaskedRenderBuffer(),
	}

	s.init.Append(s.createWidget)

	for _, o := range opts {
		o(s)
	}

	if s.padding == nil {
		s.padding = &Insets{}
	}

	return s
}

func (s *ScrollContainer) Validate() {
	if s.content == nil {
		panic("ScrollContainer: Content is required.")
	}
	// Without images of its own, a scroll container looks like its
	// theme's lists.
	s.image = s.definedImage
	if s.image == nil {
		if theme := s.GetWidget().GetTheme(); theme != nil && theme.ListTheme != nil {
			s.image = theme.ListTheme.ScrollContainerImage
		}
	}
	if s.image == nil {
		panic("ScrollContainer: Image is required, given or in the theme's lists.")
	}
	if s.image.Idle == nil {
		panic("ScrollContainer: Image.Idle is required.")
	}
	if s.image.Mask == nil {
		panic("ScrollContainer: Image.Mask is required.")
	}

	s.content.GetWidget().parent = s.widget
	s.content.GetWidget().self = s.content
	if s.forwarding != s.content {
		s.forwardContentEvents()
	}

	s.content.Validate()
	s.validated = true
}

// forwardContentEvents passes the content's context menu, focus, tooltip
// and drag and drop events on as the scroll container's, once however
// often it is validated.
func (s *ScrollContainer) forwardContentEvents() {
	s.forwarding = s.content
	s.content.GetWidget().ContextMenuEvent.AddHandler(func(args interface{}) {
		if a, ok := args.(*WidgetContextMenuEventArgs); ok {
			s.GetWidget().FireContextMenuEvent(a.Widget, a.Location)
		}
	})
	s.content.GetWidget().FocusEvent.AddHandler(func(args interface{}) {
		if a, ok := args.(*WidgetFocusEventArgs); ok {
			s.GetWidget().FireFocusEvent(a.Widget, a.Focused, a.Location)
		}
	})
	s.content.GetWidget().ToolTipEvent.AddHandler(func(args interface{}) {
		if a, ok := args.(*WidgetToolTipEventArgs); ok {
			s.GetWidget().FireToolTipEvent(a.Window, a.Show)
		}
	})
	s.content.GetWidget().DragAndDropEvent.AddHandler(func(args interface{}) {
		if a, ok := args.(*WidgetDragAndDropEventArgs); ok {
			s.GetWidget().FireDragAndDropEvent(a.Window, a.Show, a.DnD)
		}
	})
}

func (o ScrollContainerOptions) WidgetOpts(opts ...WidgetOpt) ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.widgetOpts = append(s.widgetOpts, opts...)
	}
}

// Image sets the scroll container's images; without them, it takes its
// theme's list images.
func (o ScrollContainerOptions) Image(i *ScrollContainerImage) ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.definedImage = i
	}
}

func (o ScrollContainerOptions) Content(c PreferredSizeLocateableWidget) ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.content = c
	}
}

func (o ScrollContainerOptions) Padding(p *Insets) ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.padding = p
	}
}

// DragToScroll lets the pointer scroll by dragging the content: a press in
// the view that moves more than a few pixels drags the content along and
// clicks nothing when released. Let go while moving, the content keeps
// going and slows to a stop; a press catches it.
func (o ScrollContainerOptions) DragToScroll() ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.dragToScroll = true
	}
}

func (o ScrollContainerOptions) StretchContentWidth() ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.stretchContentWidth = true
	}
}

func (s *ScrollContainer) GetWidget() *Widget {
	s.init.Do()
	return s.widget
}

func (s *ScrollContainer) SetLocation(rect img.Rectangle) {
	s.init.Do()
	s.widget.Rect = rect
}

func (s *ScrollContainer) PreferredSize() (int, int) {
	s.init.Do()
	if !s.validated {
		s.Validate()
	}
	if s.content == nil {
		return 50, 50
	}

	p, ok := s.content.(PreferredSizer)
	if !ok {
		return 50, 50
	}

	w, h := p.PreferredSize()
	return w + s.padding.Dx(), h + s.padding.Dy()
}

func (s *ScrollContainer) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	s.init.Do()

	if s.widget.IsVisible() {
		s.content.GetWidget().ElevateToNewInputLayer(&input.Layer{
			DebugLabel: "scroll container content",
			EventTypes: input.LayerEventTypeAll ^ input.LayerEventTypeWheel,
			BlockLower: true,
			FullScreen: false,
			RectFunc:   s.ViewRect,
		})

		if il, ok := s.content.(input.Layerer); ok {
			il.SetupInputLayer(def)
		}
	}
}

func (s *ScrollContainer) Render(screen *ebiten.Image) {
	s.init.Do()

	s.clampScroll()
	s.content.GetWidget().Disabled = s.widget.Disabled

	s.widget.Render(screen)

	s.draw(screen)

	s.renderContent(screen)
}

func (s *ScrollContainer) Update(updObj *UpdateObject) {
	s.init.Do()

	s.widget.Update(updObj)

	if s.content == nil {
		return
	}
	s.updateDrag()

	r, ok := s.content.(Updater)
	if !ok {
		return
	}
	r.Update(updObj)
}

func (s *ScrollContainer) GetFocusers() []Focuser {
	result := []Focuser{}
	switch v := s.content.(type) {
	case Focuser:
		if widget, ok := v.(HasWidget); ok {
			if v.TabOrder() >= 0 && !widget.GetWidget().Disabled && widget.GetWidget().IsVisible() {
				result = append(result, v)
			}
		}
	case *Container:
		result = append(result, v.GetFocusers()...)
	case *FlipBook:
		result = append(result, v.GetFocusers()...)
	case *TabBook:
		result = append(result, v.container.GetFocusers()...)
	case *TabBookTab:
		result = append(result, v.GetFocusers()...)
	case *ScrollContainer:
		result = append(result, v.GetFocusers()...)
	case *TextArea:
		result = append(result, v.GetFocusers()...)
	}
	return result
}

func (s *ScrollContainer) GetDropTargets() []HasWidget {
	result := []HasWidget{}
	if s.GetWidget().drop != nil {
		result = append(result, s)
	}

	if v, ok := s.content.(Dropper); ok {
		result = append(result, v.GetDropTargets()...)
	}

	return result
}

func (s *ScrollContainer) draw(screen *ebiten.Image) {
	i := s.image.Idle
	if s.widget.Disabled {
		if s.image.Disabled != nil {
			i = s.image.Disabled
		}
	}

	if i != nil {
		i.Draw(screen, s.widget.Rect.Dx(), s.widget.Rect.Dy(), func(opts *ebiten.DrawImageOptions) {
			s.widget.drawImageOptions(opts)
			s.drawImageOptions(opts)
		})
	}
}

func (s *ScrollContainer) drawImageOptions(opts *ebiten.DrawImageOptions) {
	if s.widget.Disabled && s.image.Disabled == nil {
		opts.ColorScale.ScaleAlpha(0.35)
	}
}

func (s *ScrollContainer) renderContent(screen *ebiten.Image) {
	if s.content == nil {
		return
	}

	r, ok := s.content.(Renderer)
	if !ok {
		return
	}

	if l, ok := s.content.(Locateable); ok {
		cw, ch := 50, 50
		if p, ok := s.content.(PreferredSizer); ok {
			cw, ch = p.PreferredSize()
		}

		vrect := s.ViewRect()
		if s.stretchContentWidth && cw < vrect.Dx() {
			cw = vrect.Dx()
		}

		rect := img.Rect(0, 0, cw, ch)
		rect = rect.Add(s.widget.Rect.Min)
		rect = rect.Add(img.Point{s.padding.Left, s.padding.Top})

		rect = rect.Sub(img.Point{int(math.Round(float64(cw-vrect.Dx()) * s.ScrollLeft)), int(math.Round(float64(ch-vrect.Dy()) * s.ScrollTop))})

		if rect != s.content.GetWidget().Rect {
			l.SetLocation(rect)

			if r, ok := s.content.(Relayoutable); ok {
				r.RequestRelayout()
			}
		}
	}

	s.renderBuf.DrawIn(screen, s.ViewRect(),
		func(buf *ebiten.Image) {
			r.Render(buf)
		},
		func(buf *ebiten.Image) {
			s.image.Mask.Draw(buf, s.widget.Rect.Dx()-s.padding.Dx(), s.widget.Rect.Dy()-s.padding.Dy(), func(opts *ebiten.DrawImageOptions) {
				opts.GeoM.Translate(float64(s.widget.Rect.Min.X+s.padding.Left), float64(s.widget.Rect.Min.Y+s.padding.Top))
				opts.Blend = ebiten.BlendCopy
			})
		})
}

func (s *ScrollContainer) ViewRect() img.Rectangle {
	s.init.Do()
	return s.padding.Apply(s.widget.Rect)
}

func (s *ScrollContainer) ContentRect() img.Rectangle {
	return s.content.GetWidget().Rect
}

// ScrollVisible scrolls vertically so w, a widget inside the content, shows
// a fraction p of the way down the view: at its top for 0, in its middle for
// 0.5, at its bottom for 1. Content shorter than the view doesn't scroll, and
// near the content's ends the scroll stops at them. It uses the content's
// layout from when it was last drawn; a w outside the content changes
// nothing.
func (s *ScrollContainer) ScrollVisible(w HasWidget, p float64) {
	s.init.Do()
	if s.content == nil || !s.contains(w.GetWidget()) {
		return
	}
	view, content, target := s.ViewRect(), s.content.GetWidget().Rect, w.GetWidget().Rect
	scrollable := content.Dy() - view.Dy()
	if scrollable <= 0 {
		s.ScrollTop = 0
		return
	}
	offset := float64(target.Min.Y-content.Min.Y) - p*float64(view.Dy()-target.Dy())
	s.ScrollTop = offset / float64(scrollable)
	s.clampScroll()
}

// contains reports whether w is the content or lies inside it.
func (s *ScrollContainer) contains(w *Widget) bool {
	for ; w != nil; w = w.parent {
		if w == s.content.GetWidget() {
			return true
		}
	}
	return false
}

// updateDrag follows a drag of the content, or moves it on after one.
func (s *ScrollContainer) updateDrag() {
	if !s.dragToScroll {
		return
	}
	d := &s.drag
	x, y := input.CursorPosition()
	spanX, spanY := s.scrollSpan()
	if input.MouseButtonJustPressedLayer(ebiten.MouseButtonLeft, s.content.GetWidget().EffectiveInputLayer()) &&
		img.Pt(x, y).In(s.ViewRect()) {
		*d = scrollDrag{pressed: true, startX: x, startY: y, lastX: x, lastY: y,
			startLeft: s.ScrollLeft * float64(spanX), startTop: s.ScrollTop * float64(spanY)}
		return
	}
	if d.pressed && input.MouseButtonPressed(ebiten.MouseButtonLeft) {
		dx, dy := x-d.startX, y-d.startY
		if !d.dragging && (spanX > 0 && abs(dx) > dragThreshold || spanY > 0 && abs(dy) > dragThreshold) {
			d.dragging = true
			cancelPresses()
		}
		if d.dragging {
			d.vx, d.vy = float64(d.lastX-x), float64(d.lastY-y)
			s.scrollTo(d.startLeft-float64(dx), d.startTop-float64(dy), spanX, spanY)
		}
		d.lastX, d.lastY = x, y
		return
	}
	if d.pressed {
		if !d.dragging {
			d.vx, d.vy = 0, 0
		}
		d.pressed, d.dragging = false, false
	}
	if d.vx == 0 && d.vy == 0 {
		return
	}
	left, top := s.ScrollLeft*float64(spanX)+d.vx, s.ScrollTop*float64(spanY)+d.vy
	if !s.scrollTo(left, top, spanX, spanY) {
		d.vx, d.vy = 0, 0 // at an end
		return
	}
	d.vx, d.vy = d.vx*dragFriction, d.vy*dragFriction
	if math.Abs(d.vx) < dragStop {
		d.vx = 0
	}
	if math.Abs(d.vy) < dragStop {
		d.vy = 0
	}
}

// scrollSpan returns how far the content scrolls each way, in pixels.
func (s *ScrollContainer) scrollSpan() (int, int) {
	view, content := s.ViewRect(), s.content.GetWidget().Rect
	return max(0, content.Dx()-view.Dx()), max(0, content.Dy()-view.Dy())
}

// scrollTo scrolls left and top pixels into the content, as far as it
// goes, firing ScrollMovedEvent; it reports whether the scroll moved.
func (s *ScrollContainer) scrollTo(left, top float64, spanX, spanY int) bool {
	l, t := s.ScrollLeft, s.ScrollTop
	if spanX > 0 {
		s.ScrollLeft = min(1, max(0, left/float64(spanX)))
	}
	if spanY > 0 {
		s.ScrollTop = min(1, max(0, top/float64(spanY)))
	}
	if s.ScrollLeft == l && s.ScrollTop == t {
		return false
	}
	s.ScrollMovedEvent.Fire(&ScrollContainerScrollMovedEventArgs{ScrollContainer: s, Left: s.ScrollLeft, Top: s.ScrollTop})
	return true
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func (s *ScrollContainer) clampScroll() {
	if s.ScrollTop < 0 {
		s.ScrollTop = 0
	} else if s.ScrollTop > 1 {
		s.ScrollTop = 1
	}

	if s.ScrollLeft < 0 {
		s.ScrollLeft = 0
	} else if s.ScrollLeft > 1 {
		s.ScrollLeft = 1
	}
}

func (s *ScrollContainer) createWidget() {
	s.widget = NewWidget(s.widgetOpts...)
	s.widgetOpts = nil
	s.widget.self = s
	s.content.GetWidget().self = s.content
}
