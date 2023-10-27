package widget

import (
	img "image"
	"math"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type Slider struct {
	Min               int
	Max               int
	Current           int
	DrawTrackDisabled bool

	ChangedEvent *event.Event

	widgetOpts         []WidgetOpt
	handleOpts         []ButtonOpt
	direction          Direction
	trackImage         *SliderTrackImage
	trackPadding       Insets
	minHandleSize      int
	fixedHandleSize    int
	trackOffset        int
	pageSizeFunc       SliderPageSizeFunc
	disableDefaultKeys bool

	init                         *MultiOnce
	widget                       *Widget
	handle                       *Button
	lastCurrent                  int
	hovering                     bool
	dragging                     bool
	handlePressedCursorX         int
	handlePressedCursorY         int
	handlePressedOffsetX         int
	handlePressedOffsetY         int
	handlePressedInternalCurrent float64

	tabOrder  int
	justMoved bool
}

type SliderTrackImage struct {
	Idle     *image.NineSlice
	Hover    *image.NineSlice
	Disabled *image.NineSlice
}

type SliderOpt func(s *Slider)

type SliderPageSizeFunc func() int

type SliderChangedEventArgs struct {
	Slider   *Slider
	Current  int
	Dragging bool
}

type SliderChangedHandlerFunc func(args *SliderChangedEventArgs)

type SliderOptions struct {
}

var SliderOpts SliderOptions

func NewSlider(opts ...SliderOpt) *Slider {
	s := &Slider{
		Min:     1,
		Max:     100,
		Current: 1,

		ChangedEvent: &event.Event{},

		trackImage:    &SliderTrackImage{},
		minHandleSize: 16,
		pageSizeFunc: func() int {
			return 10
		},

		lastCurrent: 1,

		init: &MultiOnce{},
	}

	s.init.Append(s.createWidget)

	for _, o := range opts {
		o(s)
	}

	return s
}

func (o SliderOptions) WidgetOpts(opts ...WidgetOpt) SliderOpt {
	return func(s *Slider) {
		s.widgetOpts = append(s.widgetOpts, opts...)
	}
}

func (o SliderOptions) Direction(d Direction) SliderOpt {
	return func(s *Slider) {
		s.direction = d
	}
}

func (o SliderOptions) Images(track *SliderTrackImage, handle *ButtonImage) SliderOpt {
	return func(s *Slider) {
		s.trackImage = track
		s.handleOpts = append(s.handleOpts, ButtonOpts.Image(handle))
	}
}

func (o SliderOptions) TrackPadding(i Insets) SliderOpt {
	return func(s *Slider) {
		s.trackPadding = i
	}
}
func (o SliderOptions) TrackOffset(i int) SliderOpt {
	return func(s *Slider) {
		s.trackOffset = i
	}
}

func (o SliderOptions) MinHandleSize(s int) SliderOpt {
	return func(sl *Slider) {
		sl.minHandleSize = s
	}
}

func (o SliderOptions) FixedHandleSize(s int) SliderOpt {
	return func(sl *Slider) {
		sl.fixedHandleSize = s
	}
}

func (o SliderOptions) MinMax(min int, max int) SliderOpt {
	return func(s *Slider) {
		s.Min = min
		s.Max = max
	}
}

func (o SliderOptions) PageSizeFunc(f SliderPageSizeFunc) SliderOpt {
	return func(s *Slider) {
		s.pageSizeFunc = f
	}
}

func (o SliderOptions) ChangedHandler(f SliderChangedHandlerFunc) SliderOpt {
	return func(s *Slider) {
		s.ChangedEvent.AddHandler(func(args interface{}) {
			f(args.(*SliderChangedEventArgs))
		})
	}
}

func (o SliderOptions) TabOrder(tabOrder int) SliderOpt {
	return func(sl *Slider) {
		sl.tabOrder = tabOrder
	}
}

func (o SliderOptions) DisableDefaultKeys(val bool) SliderOpt {
	return func(sl *Slider) {
		sl.disableDefaultKeys = val
	}
}

func (s *Slider) Focus(focused bool) {
	s.init.Do()
	s.GetWidget().FireFocusEvent(s, focused, img.Point{-1, -1})
	s.handle.focused = focused
}

func (s *Slider) IsFocused() bool {
	return s.handle.focused
}

func (s *Slider) TabOrder() int {
	return s.tabOrder
}

func (s *Slider) GetWidget() *Widget {
	s.init.Do()
	return s.widget
}

func (s *Slider) PreferredSize() (int, int) {
	var w, h int
	if s.direction == DirectionHorizontal {
		w = 0
		h = s.minHandleSize + s.trackPadding.Top + s.trackPadding.Bottom
	} else {
		w = s.minHandleSize + s.trackPadding.Left + s.trackPadding.Right
		h = 0
	}

	if s.widget != nil {
		if w < s.widget.MinWidth {
			w = s.widget.MinWidth
		}
		if h < s.widget.MinHeight {
			h = s.widget.MinHeight
		}
	}
	return w, h
}

func (s *Slider) SetLocation(rect img.Rectangle) {
	s.init.Do()
	s.widget.Rect = rect
}

func (s *Slider) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	s.handle.GetWidget().ElevateToNewInputLayer(&input.Layer{
		DebugLabel: "slider handle",
		EventTypes: input.LayerEventTypeAll,
		BlockLower: true,
		FullScreen: false,
		RectFunc: func() img.Rectangle {
			return s.handle.GetWidget().Rect
		},
	})

	s.handle.SetupInputLayer(def)
}

func (s *Slider) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	s.init.Do()

	s.handleDirection()
	s.clampCurrentMinMax()

	s.handle.GetWidget().Disabled = s.widget.Disabled

	s.widget.Render(screen, def)

	s.draw(screen)

	hl, tl := s.handleLengthAndTrackLength()
	s.updateHandleLocation(hl, tl)
	s.updateHandleSize(hl)

	s.handle.Render(screen, def)

	s.fireEvents()

	s.lastCurrent = s.Current
}

func (s *Slider) draw(screen *ebiten.Image) {
	i := s.trackImage.Idle
	if s.widget.Disabled || s.DrawTrackDisabled {
		if s.trackImage.Disabled != nil {
			i = s.trackImage.Disabled
		}
	} else if s.hovering {
		if s.trackImage.Hover != nil {
			i = s.trackImage.Hover
		}
	}

	if i != nil {
		i.Draw(screen, s.widget.Rect.Dx(), s.widget.Rect.Dy(), func(opts *ebiten.DrawImageOptions) {
			if s.direction == DirectionHorizontal {
				opts.GeoM.Translate(float64(s.widget.Rect.Min.X), float64(s.widget.Rect.Min.Y+s.trackOffset))
			} else {
				opts.GeoM.Translate(float64(s.widget.Rect.Min.X+s.trackOffset), float64(s.widget.Rect.Min.Y))
			}
		})
	}
}

func (s *Slider) handleDirection() {
	if !s.disableDefaultKeys {
		if s.direction == DirectionHorizontal {
			if input.KeyPressed(ebiten.KeyLeft) || input.KeyPressed(ebiten.KeyRight) {
				if !s.justMoved && s.handle.focused {
					changeDir := 1
					if input.KeyPressed(ebiten.KeyLeft) {
						changeDir = -1
					}
					s.Current = s.Current + (changeDir * s.pageSizeFunc())
					s.justMoved = true
				}
			} else {
				s.justMoved = false
			}
		} else {
			if input.KeyPressed(ebiten.KeyUp) || input.KeyPressed(ebiten.KeyDown) {
				if !s.justMoved && s.handle.focused {
					changeDir := 1
					if input.KeyPressed(ebiten.KeyUp) {
						changeDir = -1
					}
					s.Current = s.Current + (changeDir * s.pageSizeFunc())
					s.justMoved = true
				}
			} else {
				s.justMoved = false
			}
		}
	}
}

func (s *Slider) fireEvents() {
	if s.Current != s.lastCurrent {
		s.ChangedEvent.Fire(&SliderChangedEventArgs{
			Slider:  s,
			Current: s.Current,
		})
	}
}

func (s *Slider) updateHandleSize(handleLength float64) {
	l := int(math.Round(handleLength))
	if l < s.minHandleSize {
		l = s.minHandleSize
	}

	rect := s.widget.Rect

	var p img.Point
	if s.direction == DirectionHorizontal {
		p = img.Point{l, rect.Dy() - s.trackPadding.Top - s.trackPadding.Bottom}
	} else {
		p = img.Point{rect.Dx() - s.trackPadding.Left - s.trackPadding.Right, l}
	}

	s.handle.GetWidget().Rect.Max = s.handle.GetWidget().Rect.Min.Add(p)
}

func (s *Slider) updateHandleLocation(handleLength float64, trackLength float64) {
	internalTrackLength := int(math.Ceil(trackLength - handleLength))
	internalTrackStart := int(math.Floor(handleLength / 2))
	internalTrackEnd := internalTrackStart + internalTrackLength

	var i float64
	if s.dragging {
		x, y := input.CursorPosition()

		var dragOffset int
		if s.direction == DirectionHorizontal {
			dragOffset = x - s.handlePressedCursorX
		} else {
			dragOffset = y - s.handlePressedCursorY
		}
		var internalDragOffset float64
		if internalTrackLength > 0 {
			internalDragOffset = float64(dragOffset) / float64(internalTrackLength)
		} else {
			internalDragOffset = 0
		}

		i = s.handlePressedInternalCurrent + internalDragOffset
		if i < 0 {
			i = 0
		} else if i > 1 {
			i = 1
		}
		s.Current = s.internalToCurrent(i)
	} else {
		i = s.currentToInternal(s.Current)
	}

	off := int(math.Round(float64(internalTrackStart)*(1-i)+float64(internalTrackEnd)*i) - handleLength/2)

	rect := s.widget.Rect
	if s.direction == DirectionHorizontal {
		rect.Min = rect.Min.Add(img.Point{off + s.trackPadding.Left, s.trackPadding.Top})
	} else {
		rect.Min = rect.Min.Add(img.Point{s.trackPadding.Left, off + s.trackPadding.Top})
	}
	s.handle.GetWidget().Rect = rect
}

func (s *Slider) handleLengthAndTrackLength() (float64, float64) {
	var trackLength float64
	if s.direction == DirectionHorizontal {
		trackLength = float64(s.widget.Rect.Dx()) - float64(s.trackPadding.Left) - float64(s.trackPadding.Right)
	} else {
		trackLength = float64(s.widget.Rect.Dy()) - float64(s.trackPadding.Top) - float64(s.trackPadding.Bottom)
	}

	handleLength := 0.0
	if s.fixedHandleSize != 0 {
		handleLength = float64(s.fixedHandleSize)
	} else {
		ps := float64(s.pageSizeFunc())
		length := float64(s.Max - s.Min + 1)
		handleLength = ps / length * trackLength
		if handleLength < float64(s.minHandleSize) {
			handleLength = float64(s.minHandleSize)
		}
	}

	if handleLength > trackLength {
		handleLength = trackLength
	}

	return handleLength, trackLength
}

func (s *Slider) currentToInternal(c int) float64 {
	if s.Max <= s.Min {
		return 0
	}

	return float64(c-s.Min) / float64(s.Max-s.Min)
}

func (s *Slider) internalToCurrent(i float64) int {
	return int(math.Round(float64(s.Min)*(1-i) + float64(s.Max)*i))
}

func (s *Slider) clampCurrentMinMax() {
	if s.Current < s.Min {
		s.Current = s.Min
	} else if s.Current > s.Max {
		s.Current = s.Max
	}
}

func (s *Slider) createWidget() {
	s.widget = NewWidget(append(s.widgetOpts, []WidgetOpt{
		WidgetOpts.CursorEnterHandler(func(_ *WidgetCursorEnterEventArgs) {
			if !s.widget.Disabled {
				s.hovering = true
			}
		}),

		WidgetOpts.CursorExitHandler(func(_ *WidgetCursorExitEventArgs) {
			s.hovering = false
		}),
		WidgetOpts.ScrolledHandler(func(args *WidgetScrolledEventArgs) {
			if !s.widget.Disabled {
				ps := s.pageSizeFunc()
				if s.direction == DirectionHorizontal {
					s.Current += ps * int(args.Y)
				} else {
					s.Current += ps * int(args.X)
				}
				s.clampCurrentMinMax()
			}
		}),

		// TODO: keeping the mouse button pressed should move the handle repeatedly (in PageSize steps) until it stops under the cursor
		WidgetOpts.MouseButtonPressedHandler(func(args *WidgetMouseButtonPressedEventArgs) {
			if !s.widget.Disabled && args.Button == ebiten.MouseButtonLeft {
				x, y := input.CursorPosition()
				ps := s.pageSizeFunc()
				rect := s.handle.GetWidget().Rect
				if s.direction == DirectionHorizontal {
					if x < rect.Min.X {
						s.Current -= ps
					} else if x >= rect.Max.X {
						s.Current += ps
					}
				} else {
					if y < rect.Min.Y {
						s.Current -= ps
					} else if y >= rect.Max.Y {
						s.Current += ps
					}
				}

				s.clampCurrentMinMax()
			}
		}),
	}...)...)
	s.widgetOpts = nil

	s.handle = NewButton(append(s.handleOpts, []ButtonOpt{
		ButtonOpts.KeepPressedOnExit(),

		ButtonOpts.PressedHandler(func(args *ButtonPressedEventArgs) {
			s.dragging = true
			s.handlePressedCursorX, s.handlePressedCursorY = input.CursorPosition()
			s.handlePressedOffsetX = args.OffsetX
			s.handlePressedOffsetY = args.OffsetY
			s.handlePressedInternalCurrent = s.currentToInternal(s.Current)
		}),

		ButtonOpts.ReleasedHandler(func(_ *ButtonReleasedEventArgs) {
			s.dragging = false
		}),

		ButtonOpts.WidgetOpts(WidgetOpts.ScrolledHandler(func(args *WidgetScrolledEventArgs) {
			if !s.widget.Disabled {
				ps := s.pageSizeFunc()
				if s.direction == DirectionHorizontal {
					s.Current += ps * int(args.Y)
				} else {
					s.Current += ps * int(args.X)
				}
				s.clampCurrentMinMax()
			}
		})),
	}...)...)
}
