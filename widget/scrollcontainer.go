package widget

import (
	img "image"
	"math"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type ScrollContainer struct {
	ScrollLeft float64
	ScrollTop  float64

	widgetOpts          []WidgetOpt
	image               *ScrollContainerImage
	content             HasWidget
	padding             Insets
	stretchContentWidth bool

	init      *MultiOnce
	widget    *Widget
	renderBuf *image.MaskedRenderBuffer
}

type ScrollContainerOpt func(s *ScrollContainer)

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
		init: &MultiOnce{},

		renderBuf: image.NewMaskedRenderBuffer(),
	}

	s.init.Append(s.createWidget)

	for _, o := range opts {
		o(s)
	}

	s.content.GetWidget().ContextMenuEvent.AddHandler(func(args interface{}) {
		a := args.(*WidgetContextMenuEventArgs)
		s.GetWidget().FireContextMenuEvent(a.Widget, a.Location)
	})
	s.content.GetWidget().FocusEvent.AddHandler(func(args interface{}) {
		a := args.(*WidgetFocusEventArgs)
		s.GetWidget().FireFocusEvent(a.Widget, a.Focused, a.Location)
	})
	s.content.GetWidget().ToolTipEvent.AddHandler(func(args interface{}) {
		a := args.(*WidgetToolTipEventArgs)
		s.GetWidget().FireToolTipEvent(a.Window, a.Show)
	})
	s.content.GetWidget().DragAndDropEvent.AddHandler(func(args interface{}) {
		a := args.(*WidgetDragAndDropEventArgs)
		s.GetWidget().FireDragAndDropEvent(a.Window, a.Show, a.DnD)
	})
	return s
}

func (o ScrollContainerOptions) WidgetOpts(opts ...WidgetOpt) ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.widgetOpts = append(s.widgetOpts, opts...)
	}
}

func (o ScrollContainerOptions) Image(i *ScrollContainerImage) ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.image = i
	}
}

func (o ScrollContainerOptions) Content(c HasWidget) ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.content = c
	}
}

func (o ScrollContainerOptions) Padding(p Insets) ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.padding = p
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

func (s *ScrollContainer) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	s.init.Do()

	s.clampScroll()
	s.content.GetWidget().Disabled = s.widget.Disabled

	s.widget.Render(screen, def)

	s.draw(screen)

	s.renderContent(screen, def)
}
func (s *ScrollContainer) GetFocusers() []Focuser {
	result := []Focuser{}
	switch v := s.content.(type) {
	case Focuser:
		if v.TabOrder() >= 0 && !v.(HasWidget).GetWidget().Disabled {
			result = append(result, v)
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
		opts.ColorM.Scale(1, 1, 1, 0.35)
	}
}

func (s *ScrollContainer) renderContent(screen *ebiten.Image, def DeferredRenderFunc) {
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

	s.renderBuf.Draw(screen,
		func(buf *ebiten.Image) {
			r.Render(buf, def)
		},
		func(buf *ebiten.Image) {
			s.image.Mask.Draw(buf, s.widget.Rect.Dx()-s.padding.Dx(), s.widget.Rect.Dy()-s.padding.Dy(), func(opts *ebiten.DrawImageOptions) {
				opts.GeoM.Translate(float64(s.widget.Rect.Min.X+s.padding.Left), float64(s.widget.Rect.Min.Y+s.padding.Top))
				opts.CompositeMode = ebiten.CompositeModeCopy
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
}
