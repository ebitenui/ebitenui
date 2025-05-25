package widget

import (
	img "image"
	"math"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type SliderParams struct {
	Orientation     *Direction
	TrackImage      *SliderTrackImage
	TrackPadding    *Insets
	MinHandleSize   *int
	FixedHandleSize *int
	TrackOffset     *int
	HandleImage     *ButtonImage
	PageSizeFunc    SliderPageSizeFunc
}

type Slider struct {
	definedParams  SliderParams
	computedParams SliderParams

	Min                int
	Max                int
	Current            int
	DrawTrackDisabled  bool
	disableDefaultKeys bool

	widgetOpts []WidgetOpt

	ChangedEvent *event.Event

	init   *MultiOnce
	widget *Widget
	handle *Button

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
	focusMap  map[FocusDirection]Focuser
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

		lastCurrent: 1,

		init:     &MultiOnce{},
		focusMap: make(map[FocusDirection]Focuser),
	}

	s.init.Append(s.createWidget)

	for _, o := range opts {
		o(s)
	}

	return s
}

func (s *Slider) Validate() {
	s.init.Do()
	s.populateComputedParams()

	s.setChildComputedParams()
}

func (s *Slider) populateComputedParams() {
	params := SliderParams{}

	// Set Theme
	theme := s.widget.GetTheme()
	if theme != nil {
		if theme.SliderTheme != nil {
			params.FixedHandleSize = theme.SliderTheme.FixedHandleSize
			params.HandleImage = theme.SliderTheme.HandleImage
			params.MinHandleSize = theme.SliderTheme.MinHandleSize
			params.Orientation = theme.SliderTheme.Orientation
			params.PageSizeFunc = theme.SliderTheme.PageSizeFunc
			params.TrackImage = theme.SliderTheme.TrackImage
			params.TrackOffset = theme.SliderTheme.TrackOffset
			params.TrackPadding = theme.SliderTheme.TrackPadding
		}
	}

	// Set defined params
	if s.definedParams.FixedHandleSize != nil {
		params.FixedHandleSize = s.definedParams.FixedHandleSize
	}
	if s.definedParams.MinHandleSize != nil {
		params.MinHandleSize = s.definedParams.MinHandleSize
	}
	if s.definedParams.HandleImage != nil {
		params.HandleImage = s.definedParams.HandleImage
	}
	if s.definedParams.Orientation != nil {
		params.Orientation = s.definedParams.Orientation
	}
	if s.definedParams.PageSizeFunc != nil {
		params.PageSizeFunc = s.definedParams.PageSizeFunc
	}
	if s.definedParams.TrackImage != nil {
		params.TrackImage = s.definedParams.TrackImage
	}
	if s.definedParams.TrackOffset != nil {
		params.TrackOffset = s.definedParams.TrackOffset
	}
	if s.definedParams.TrackPadding != nil {
		params.TrackPadding = s.definedParams.TrackPadding
	}

	// Set defaults
	if params.MinHandleSize == nil {
		size := 16
		params.MinHandleSize = &size
	}
	if params.PageSizeFunc == nil {
		params.PageSizeFunc = func() int {
			return 10
		}
	}
	if params.Orientation == nil {
		o := DirectionHorizontal
		params.Orientation = &o
	}
	if params.TrackPadding == nil {
		params.TrackPadding = &Insets{}
	}
	if params.TrackOffset == nil {
		o := 0
		params.TrackOffset = &o
	}
	if params.TrackImage == nil {
		params.TrackImage = &SliderTrackImage{}
	}
	if params.FixedHandleSize == nil {
		o := 0
		params.FixedHandleSize = &o
	}

	s.computedParams = params
}

func (o SliderOptions) WidgetOpts(opts ...WidgetOpt) SliderOpt {
	return func(s *Slider) {
		s.widgetOpts = append(s.widgetOpts, opts...)
	}
}

func (so SliderOptions) Orientation(o Direction) SliderOpt {
	return func(s *Slider) {
		s.definedParams.Orientation = &o
	}
}

// Deprecated: Use Orientation(o *Direction) instead.
func (o SliderOptions) Direction(d Direction) SliderOpt {
	return func(s *Slider) {
		s.definedParams.Orientation = &d
	}
}

// This sets the track images (not required) and the handle images (required).
func (o SliderOptions) Images(track *SliderTrackImage, handle *ButtonImage) SliderOpt {
	return func(s *Slider) {
		s.definedParams.TrackImage = track
		s.definedParams.HandleImage = handle
	}
}

// This sets the track images (not required).
func (o SliderOptions) TrackImage(track *SliderTrackImage) SliderOpt {
	return func(s *Slider) {
		s.definedParams.TrackImage = track
	}
}

// This sets the handle images (required).
func (o SliderOptions) HandleImage(handle *ButtonImage) SliderOpt {
	return func(s *Slider) {
		s.definedParams.HandleImage = handle
	}
}

func (o SliderOptions) TrackPadding(i *Insets) SliderOpt {
	return func(s *Slider) {
		s.definedParams.TrackPadding = i
	}
}
func (o SliderOptions) TrackOffset(i int) SliderOpt {
	return func(s *Slider) {
		s.definedParams.TrackOffset = &i
	}
}

func (o SliderOptions) MinHandleSize(s int) SliderOpt {
	return func(sl *Slider) {
		sl.definedParams.MinHandleSize = &s
	}
}

func (o SliderOptions) FixedHandleSize(s int) SliderOpt {
	return func(sl *Slider) {
		sl.definedParams.FixedHandleSize = &s
	}
}

func (o SliderOptions) MinMax(minValue int, maxValue int) SliderOpt {
	return func(s *Slider) {
		s.Min = minValue
		s.Max = maxValue
	}
}

func (o SliderOptions) InitialCurrent(value int) SliderOpt {
	return func(s *Slider) {
		s.Current = value
		s.lastCurrent = value
	}
}

func (o SliderOptions) PageSizeFunc(f SliderPageSizeFunc) SliderOpt {
	return func(s *Slider) {
		s.definedParams.PageSizeFunc = f
	}
}

func (o SliderOptions) ChangedHandler(f SliderChangedHandlerFunc) SliderOpt {
	return func(s *Slider) {
		s.ChangedEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*SliderChangedEventArgs); ok {
				f(arg)
			}
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

/** Focuser Interface - Start **/

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

func (s *Slider) GetFocus(direction FocusDirection) Focuser {
	return s.focusMap[direction]
}

func (s *Slider) AddFocus(direction FocusDirection, focus Focuser) {
	s.focusMap[direction] = focus
}

/** Focuser Interface - End **/

func (s *Slider) GetWidget() *Widget {
	s.init.Do()
	return s.widget
}

func (s *Slider) PreferredSize() (int, int) {
	var w, h int
	if *s.computedParams.Orientation == DirectionHorizontal {
		w = 0
		h = *s.computedParams.MinHandleSize + s.computedParams.TrackPadding.Top + s.computedParams.TrackPadding.Bottom
	} else {
		w = *s.computedParams.MinHandleSize + s.computedParams.TrackPadding.Left + s.computedParams.TrackPadding.Right
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

func (s *Slider) Render(screen *ebiten.Image) {
	s.init.Do()

	s.handleOrientation()
	s.clampCurrentMinMax()

	s.handle.GetWidget().Disabled = s.widget.Disabled

	s.widget.Render(screen)

	s.drawTrack(screen)

	hl, tl := s.handleLengthAndTrackLength()
	s.updateHandleLocation(hl, tl)
	s.updateHandleSize(hl)

	s.handle.Render(screen)

	if s.Current != s.lastCurrent {
		s.fireEvents()
	}

	s.lastCurrent = s.Current
}

func (s *Slider) Update() {
	s.init.Do()

	s.widget.Update()
	s.handle.Update()
}

func (s *Slider) drawTrack(screen *ebiten.Image) {
	if s.computedParams.TrackImage != nil {
		i := s.computedParams.TrackImage.Idle
		if s.widget.Disabled || s.DrawTrackDisabled {
			if s.computedParams.TrackImage.Disabled != nil {
				i = s.computedParams.TrackImage.Disabled
			}
		} else if s.hovering {
			if s.computedParams.TrackImage.Hover != nil {
				i = s.computedParams.TrackImage.Hover
			}
		}

		if i != nil {
			i.Draw(screen, s.widget.Rect.Dx(), s.widget.Rect.Dy(), func(opts *ebiten.DrawImageOptions) {
				if *s.computedParams.Orientation == DirectionHorizontal {
					opts.GeoM.Translate(float64(s.widget.Rect.Min.X), float64(s.widget.Rect.Min.Y+*s.computedParams.TrackOffset))
				} else {
					opts.GeoM.Translate(float64(s.widget.Rect.Min.X+*s.computedParams.TrackOffset), float64(s.widget.Rect.Min.Y))
				}
			})
		}
	}
}

func (s *Slider) handleOrientation() {
	if !s.disableDefaultKeys {
		if *s.computedParams.Orientation == DirectionHorizontal {
			if input.KeyPressed(ebiten.KeyLeft) || input.KeyPressed(ebiten.KeyRight) {
				if !s.justMoved && s.handle.focused {
					changeDir := 1
					if input.KeyPressed(ebiten.KeyLeft) {
						changeDir = -1
					}
					s.Current += (changeDir * s.computedParams.PageSizeFunc())
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
					s.Current += (changeDir * s.computedParams.PageSizeFunc())
					s.justMoved = true
				}
			} else {
				s.justMoved = false
			}
		}
	}
}

func (s *Slider) fireEvents() {
	s.ChangedEvent.Fire(&SliderChangedEventArgs{
		Slider:   s,
		Current:  s.Current,
		Dragging: s.dragging,
	})
}

func (s *Slider) updateHandleSize(handleLength float64) {
	l := int(math.Round(handleLength))
	if l < *s.computedParams.MinHandleSize {
		l = *s.computedParams.MinHandleSize
	}

	rect := s.widget.Rect

	var p img.Point
	if *s.computedParams.Orientation == DirectionHorizontal {
		p = img.Point{l, rect.Dy() - s.computedParams.TrackPadding.Top - s.computedParams.TrackPadding.Bottom}
	} else {
		p = img.Point{rect.Dx() - s.computedParams.TrackPadding.Left - s.computedParams.TrackPadding.Right, l}
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
		if *s.computedParams.Orientation == DirectionHorizontal {
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
	if *s.computedParams.Orientation == DirectionHorizontal {
		rect.Min = rect.Min.Add(img.Point{off + s.computedParams.TrackPadding.Left, s.computedParams.TrackPadding.Top})
	} else {
		rect.Min = rect.Min.Add(img.Point{s.computedParams.TrackPadding.Left, off + s.computedParams.TrackPadding.Top})
	}
	s.handle.GetWidget().Rect = rect
}

func (s *Slider) handleLengthAndTrackLength() (float64, float64) {
	var trackLength float64
	if *s.computedParams.Orientation == DirectionHorizontal {
		trackLength = float64(s.widget.Rect.Dx()) - float64(s.computedParams.TrackPadding.Left) - float64(s.computedParams.TrackPadding.Right)
	} else {
		trackLength = float64(s.widget.Rect.Dy()) - float64(s.computedParams.TrackPadding.Top) - float64(s.computedParams.TrackPadding.Bottom)
	}

	handleLength := 0.0
	if *s.computedParams.FixedHandleSize != 0 {
		handleLength = float64(*s.computedParams.FixedHandleSize)
	} else {
		ps := float64(s.computedParams.PageSizeFunc())
		length := float64(s.Max - s.Min + 1)
		handleLength = ps / length * trackLength
		if handleLength < float64(*s.computedParams.MinHandleSize) {
			handleLength = float64(*s.computedParams.MinHandleSize)
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
	s.widget = NewWidget(append([]WidgetOpt{
		WidgetOpts.TrackHover(true),
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
				ps := s.computedParams.PageSizeFunc()
				if *s.computedParams.Orientation == DirectionHorizontal {
					s.Current += ps * int(args.Y)
				} else {
					s.Current -= ps * int(args.Y)
				}
				s.clampCurrentMinMax()
			}
		}),

		// TODO: keeping the mouse button pressed should move the handle repeatedly (in PageSize steps) until it stops under the cursor
		WidgetOpts.MouseButtonPressedHandler(func(args *WidgetMouseButtonPressedEventArgs) {
			if !s.widget.Disabled && args.Button == ebiten.MouseButtonLeft {
				x, y := input.CursorPosition()
				ps := s.computedParams.PageSizeFunc()
				rect := s.handle.GetWidget().Rect
				if *s.computedParams.Orientation == DirectionHorizontal {
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
	}, s.widgetOpts...)...)
	s.widgetOpts = nil

	s.handle = NewButton([]ButtonOpt{
		ButtonOpts.Image(s.computedParams.HandleImage),
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
			s.fireEvents()
		}),

		ButtonOpts.WidgetOpts(WidgetOpts.ScrolledHandler(func(args *WidgetScrolledEventArgs) {
			if !s.widget.Disabled {
				ps := s.computedParams.PageSizeFunc()
				if *s.computedParams.Orientation == DirectionHorizontal {
					s.Current += ps * int(args.Y)
				} else {
					s.Current -= ps * int(args.Y)
				}
				s.clampCurrentMinMax()
			}
		})),
	}...)
}

func (s *Slider) setChildComputedParams() {
	s.handle.definedParams.Image = s.computedParams.HandleImage
	s.handle.Validate()
}
