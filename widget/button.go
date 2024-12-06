package widget

import (
	img "image"
	"image/color"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type ButtonParams struct {
	Image        *ButtonImage
	GraphicImage *ButtonImageImage
	TextColor    *ButtonTextColor

	VTextPosition  TextPosition
	HTextPosition  TextPosition
	TextPadding    *Insets
	TextFace       *text.Face
	GraphicPadding *Insets
}

type Button struct {
	definedParams           ButtonParams
	computedParams          ButtonParams
	IgnoreTransparentPixels bool
	KeepPressedOnExit       bool
	ToggleMode              bool

	// Allows the user to disable space bar and enter automatically triggering a focused button.
	DisableDefaultKeys bool

	PressedEvent       *event.Event
	ReleasedEvent      *event.Event
	ClickedEvent       *event.Event
	CursorEnteredEvent *event.Event
	CursorMovedEvent   *event.Event
	CursorExitedEvent  *event.Event
	StateChangedEvent  *event.Event

	widgetOpts               []WidgetOpt
	autoUpdateTextAndGraphic bool

	init              *MultiOnce
	widget            *Widget
	container         *Container
	graphic           *Graphic
	mask              []byte
	text              *Text
	textLabel         string
	textProcessBBCode bool
	hovering          bool
	pressing          bool
	state             WidgetState

	tabOrder      int
	focused       bool
	justSubmitted bool

	focusMap map[FocusDirection]Focuser
}

type ButtonOpt func(b *Button)

type ButtonImage struct {
	Idle         *image.NineSlice
	Hover        *image.NineSlice
	Pressed      *image.NineSlice
	PressedHover *image.NineSlice
	Disabled     *image.NineSlice
}

type ButtonImageImage struct {
	Idle     *ebiten.Image
	Disabled *ebiten.Image
}

type ButtonTextColor struct {
	Idle     color.Color
	Disabled color.Color
	Hover    color.Color
	Pressed  color.Color
}

type ButtonPressedEventArgs struct {
	Button  *Button
	OffsetX int
	OffsetY int
}

type ButtonReleasedEventArgs struct {
	Button  *Button
	Inside  bool
	OffsetX int
	OffsetY int
}

type ButtonClickedEventArgs struct {
	Button  *Button
	OffsetX int
	OffsetY int
}

type ButtonHoverEventArgs struct {
	Button  *Button
	Entered bool
	OffsetX int
	OffsetY int
	DiffX   int
	DiffY   int
}

type ButtonChangedEventArgs struct {
	Button  *Button
	State   WidgetState
	OffsetX int
	OffsetY int
}

type ButtonPressedHandlerFunc func(args *ButtonPressedEventArgs)

type ButtonReleasedHandlerFunc func(args *ButtonReleasedEventArgs)

type ButtonClickedHandlerFunc func(args *ButtonClickedEventArgs)

type ButtonCursorHoverHandlerFunc func(args *ButtonHoverEventArgs)

type ButtonChangedHandlerFunc func(args *ButtonChangedEventArgs)

type ButtonOptions struct {
}

var ButtonOpts ButtonOptions

func NewButton(opts ...ButtonOpt) *Button {
	b := &Button{
		PressedEvent:       &event.Event{},
		ReleasedEvent:      &event.Event{},
		ClickedEvent:       &event.Event{},
		CursorEnteredEvent: &event.Event{},
		CursorMovedEvent:   &event.Event{},
		CursorExitedEvent:  &event.Event{},
		StateChangedEvent:  &event.Event{},

		init: &MultiOnce{},

		focusMap: make(map[FocusDirection]Focuser),
	}

	b.init.Append(b.createWidget)

	for _, o := range opts {
		o(b)
	}

	return b
}

func (b *Button) Validate() {

	b.populateComputedParams()

	if b.computedParams.Image == nil {
		panic("Button: Image is required.")
	}
	if b.computedParams.Image.Idle == nil {
		panic("Button: Image.Idle is required.")
	}
	if b.computedParams.Image.Pressed == nil {
		panic("Button: Image.Pressed is required.")
	}

	if len(b.textLabel) > 0 {
		if b.computedParams.TextFace == nil {
			panic("Button: TextFace is required if TextLabel is set.")
		}
		if b.computedParams.TextColor == nil {
			panic("Button: TextColor is required if TextLabel is set.")
		}
		if b.computedParams.TextColor.Idle == nil {
			panic("Button: TextColor.Idle is required if TextLabel is set.")
		}
	}

	b.initText()
}

func (b *Button) populateComputedParams() {
	btnParams := ButtonParams{
		HTextPosition: TextPositionCenter,
		VTextPosition: TextPositionCenter,
	}
	theme := b.widget.GetTheme()
	// clone the theme
	if theme != nil {
		btnParams.TextFace = theme.DefaultFace
		if theme.DefaultTextColor != nil {
			btnParams.TextColor = &ButtonTextColor{
				Idle:     theme.DefaultTextColor,
				Disabled: theme.DefaultTextColor,
				Hover:    theme.DefaultTextColor,
				Pressed:  theme.DefaultTextColor,
			}
		}
		if theme.ButtonTheme != nil {
			if theme.ButtonTheme.Image != nil {
				btnParams.Image = &ButtonImage{
					Idle:         theme.ButtonTheme.Image.Idle,
					Hover:        theme.ButtonTheme.Image.Hover,
					Pressed:      theme.ButtonTheme.Image.Pressed,
					PressedHover: theme.ButtonTheme.Image.PressedHover,
					Disabled:     theme.ButtonTheme.Image.Disabled,
				}
			}
			if theme.ButtonTheme.GraphicImage != nil {
				btnParams.GraphicImage = &ButtonImageImage{
					Idle:     theme.ButtonTheme.GraphicImage.Idle,
					Disabled: theme.ButtonTheme.GraphicImage.Disabled,
				}
			}
			btnParams.GraphicPadding = theme.ButtonTheme.GraphicPadding
			btnParams.HTextPosition = theme.ButtonTheme.HTextPosition
			btnParams.VTextPosition = theme.ButtonTheme.VTextPosition
			btnParams.TextPadding = theme.ButtonTheme.TextPadding

			if theme.ButtonTheme.TextFace != nil {
				btnParams.TextFace = theme.ButtonTheme.TextFace
			}

			if theme.ButtonTheme.TextColor != nil {
				if btnParams.TextColor == nil {
					btnParams.TextColor = theme.ButtonTheme.TextColor
				} else {
					if theme.ButtonTheme.TextColor.Disabled != nil {
						btnParams.TextColor.Disabled = theme.ButtonTheme.TextColor.Disabled
					}
					if theme.ButtonTheme.TextColor.Hover != nil {
						btnParams.TextColor.Hover = theme.ButtonTheme.TextColor.Hover
					}
					if theme.ButtonTheme.TextColor.Idle != nil {
						btnParams.TextColor.Idle = theme.ButtonTheme.TextColor.Idle
					}
					if theme.ButtonTheme.TextColor.Pressed != nil {
						btnParams.TextColor.Pressed = theme.ButtonTheme.TextColor.Pressed
					}
				}
			}
		}
	}

	if b.definedParams.Image != nil {
		if b.definedParams.Image.Idle != nil {
			btnParams.Image.Idle = b.definedParams.Image.Idle
		}
		if b.definedParams.Image.Hover != nil {
			btnParams.Image.Hover = b.definedParams.Image.Hover
		}
		if b.definedParams.Image.Pressed != nil {
			btnParams.Image.Pressed = b.definedParams.Image.Pressed
		}
		if b.definedParams.Image.PressedHover != nil {
			btnParams.Image.PressedHover = b.definedParams.Image.PressedHover
		}
		if b.definedParams.Image.Disabled != nil {
			btnParams.Image.Disabled = b.definedParams.Image.Disabled
		}
	}

	if b.definedParams.GraphicImage != nil {
		if b.definedParams.GraphicImage.Idle != nil {
			btnParams.GraphicImage.Idle = b.definedParams.GraphicImage.Idle
		}
		if b.definedParams.GraphicImage.Disabled != nil {
			btnParams.GraphicImage.Disabled = b.definedParams.GraphicImage.Disabled
		}
	}
	if b.definedParams.GraphicPadding != nil {
		btnParams.GraphicPadding = b.definedParams.GraphicPadding
	}
	if b.definedParams.HTextPosition != TextPositionCenter {
		btnParams.HTextPosition = b.definedParams.HTextPosition
	}
	if b.definedParams.VTextPosition != TextPositionCenter {
		btnParams.VTextPosition = b.definedParams.VTextPosition
	}
	if b.definedParams.TextFace != nil {
		btnParams.TextFace = b.definedParams.TextFace
	}
	if b.definedParams.TextPadding != nil {
		btnParams.TextPadding = b.definedParams.TextPadding
	}
	if b.definedParams.TextColor != nil {
		if btnParams.TextColor == nil {
			btnParams.TextColor = b.definedParams.TextColor
		} else {
			if b.definedParams.TextColor.Disabled != nil {
				btnParams.TextColor.Disabled = b.definedParams.TextColor.Disabled
			}
			if b.definedParams.TextColor.Hover != nil {
				btnParams.TextColor.Hover = b.definedParams.TextColor.Hover
			}
			if b.definedParams.TextColor.Idle != nil {
				btnParams.TextColor.Idle = b.definedParams.TextColor.Idle
			}
			if b.definedParams.TextColor.Pressed != nil {
				btnParams.TextColor.Pressed = b.definedParams.TextColor.Pressed
			}
		}
	}

	b.computedParams = btnParams
}

func (o ButtonOptions) WidgetOpts(opts ...WidgetOpt) ButtonOpt {
	return func(b *Button) {
		b.widgetOpts = append(b.widgetOpts, opts...)
	}
}

func (o ButtonOptions) Image(i *ButtonImage) ButtonOpt {
	return func(b *Button) {
		b.definedParams.Image = i
	}
}

// IgnoreTransparentPixels disables mouse events like cursor entered,
// moved and exited if the mouse pointer is over a pixel that is transparent
// (alpha = 0). The source of pixels is Image.Idle. This options is
// especially useful, if your button does not have a rectangular shape.
func (o ButtonOptions) IgnoreTransparentPixels(ignoreTransparentPixels bool) ButtonOpt {
	return func(b *Button) {
		b.IgnoreTransparentPixels = ignoreTransparentPixels
	}
}

// Text combines three options: TextLabel, TextFace and TextColor.
// It can be used for the inline configurations of Text object while
// separate functions are useful for a multi-step configuration.
func (o ButtonOptions) Text(label string, face *text.Face, color *ButtonTextColor) ButtonOpt {
	return func(b *Button) {
		b.textLabel = label
		b.definedParams.TextFace = face
		b.definedParams.TextColor = color
	}
}

func (o ButtonOptions) TextLabel(label string) ButtonOpt {
	return func(b *Button) {
		b.textLabel = label
	}
}

func (o ButtonOptions) TextFace(face *text.Face) ButtonOpt {
	return func(b *Button) {
		b.definedParams.TextFace = face
	}
}

func (o ButtonOptions) TextColor(color *ButtonTextColor) ButtonOpt {
	return func(b *Button) {
		b.definedParams.TextColor = color
	}
}

func (o ButtonOptions) TextProcessBBCode(enabled bool) ButtonOpt {
	return func(b *Button) {
		b.textProcessBBCode = enabled
	}
}

// TODO: add parameter for image position (start/end).
func (o ButtonOptions) TextAndImage(label string, face *text.Face, image *ButtonImageImage, color *ButtonTextColor) ButtonOpt {
	return func(b *Button) {
		b.init.Append(func() {
			b.container = NewContainer(
				ContainerOpts.Layout(NewAnchorLayout(AnchorLayoutOpts.Padding(*b.computedParams.TextPadding))),
				ContainerOpts.AutoDisableChildren(),
			)

			c := NewContainer(
				ContainerOpts.WidgetOpts(WidgetOpts.LayoutData(AnchorLayoutData{
					HorizontalPosition: AnchorLayoutPositionCenter,
					VerticalPosition:   AnchorLayoutPositionCenter,
				})),
				ContainerOpts.Layout(NewRowLayout(RowLayoutOpts.Spacing(10))),
				ContainerOpts.AutoDisableChildren(),
			)
			b.container.AddChild(c)

			b.text = NewText(
				TextOpts.WidgetOpts(WidgetOpts.LayoutData(RowLayoutData{
					Stretch: true,
				})),
				TextOpts.Text(label, face, color.Idle),
				TextOpts.ProcessBBCode(b.textProcessBBCode),
			)
			c.AddChild(b.text)

			b.graphic = NewGraphic(
				GraphicOpts.WidgetOpts(WidgetOpts.LayoutData(RowLayoutData{
					Stretch: true,
				})),
				GraphicOpts.Image(image.Idle))
			c.AddChild(b.graphic)

			b.autoUpdateTextAndGraphic = true
			b.definedParams.GraphicImage = image
			b.definedParams.TextColor = color
		})
	}
}

// TextPosition sets the horizontal and vertical position of the text within the button.
// Default is TextPositionCenter for both.
func (o ButtonOptions) TextPosition(h TextPosition, v TextPosition) ButtonOpt {
	return func(b *Button) {
		b.definedParams.HTextPosition = h
		b.definedParams.VTextPosition = v
	}
}

func (o ButtonOptions) TextPadding(p Insets) ButtonOpt {
	return func(b *Button) {
		b.definedParams.TextPadding = &p
	}
}

func (o ButtonOptions) Graphic(i *ebiten.Image) ButtonOpt {
	return o.withGraphic(GraphicOpts.Image(i))
}

func (o ButtonOptions) GraphicNineSlice(i *image.NineSlice) ButtonOpt {
	return o.withGraphic(GraphicOpts.ImageNineSlice(i))
}

func (o ButtonOptions) withGraphic(opt GraphicOpt) ButtonOpt {
	return func(b *Button) {
		b.init.Append(func() {
			b.container = NewContainer(
				ContainerOpts.Layout(NewAnchorLayout(AnchorLayoutOpts.Padding(*b.computedParams.GraphicPadding))),
				ContainerOpts.AutoDisableChildren())

			b.graphic = NewGraphic(
				opt,
				GraphicOpts.WidgetOpts(WidgetOpts.LayoutData(AnchorLayoutData{
					HorizontalPosition: AnchorLayoutPositionCenter,
					VerticalPosition:   AnchorLayoutPositionCenter,
				})),
			)
			b.container.AddChild(b.graphic)

			b.autoUpdateTextAndGraphic = true
		})
	}
}

func (o ButtonOptions) GraphicPadding(i Insets) ButtonOpt {
	return func(b *Button) {
		b.definedParams.GraphicPadding = &i
	}
}

func (o ButtonOptions) KeepPressedOnExit() ButtonOpt {
	return func(b *Button) {
		b.KeepPressedOnExit = true
	}
}

func (o ButtonOptions) ToggleMode() ButtonOpt {
	return func(b *Button) {
		b.ToggleMode = true
	}
}

// This option will disable enter and space from submitting a focused button.
func (o ButtonOptions) DisableDefaultKeys() ButtonOpt {
	return func(b *Button) {
		b.DisableDefaultKeys = true
	}
}

func (o ButtonOptions) PressedHandler(f ButtonPressedHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.PressedEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*ButtonPressedEventArgs); ok {
				f(arg)
			}
		})
	}
}

func (o ButtonOptions) ReleasedHandler(f ButtonReleasedHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.ReleasedEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*ButtonReleasedEventArgs); ok {
				f(arg)
			}
		})
	}
}

func (o ButtonOptions) ClickedHandler(f ButtonClickedHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.ClickedEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*ButtonClickedEventArgs); ok {
				f(arg)
			}
		})
	}
}

func (o ButtonOptions) CursorEnteredHandler(f ButtonCursorHoverHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.CursorEnteredEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*ButtonHoverEventArgs); ok {
				f(arg)
			}
		})
	}
}

func (o ButtonOptions) CursorMovedHandler(f ButtonCursorHoverHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.CursorMovedEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*ButtonHoverEventArgs); ok {
				f(arg)
			}
		})
	}
}

func (o ButtonOptions) CursorExitedHandler(f ButtonCursorHoverHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.CursorExitedEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*ButtonHoverEventArgs); ok {
				f(arg)
			}
		})
	}
}

func (o ButtonOptions) StateChangedHandler(f ButtonChangedHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.StateChangedEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*ButtonChangedEventArgs); ok {
				f(arg)
			}
		})
	}
}

func (o ButtonOptions) TabOrder(tabOrder int) ButtonOpt {
	return func(b *Button) {
		b.tabOrder = tabOrder
	}
}

func (b *Button) State() WidgetState {
	return b.state
}

func (b *Button) SetState(state WidgetState) {
	if state != b.state {
		b.state = state

		b.StateChangedEvent.Fire(&ButtonChangedEventArgs{
			Button:  b,
			State:   b.state,
			OffsetX: -1,
			OffsetY: -1,
		})
	}
}

func (b *Button) getStateChangedEvent() *event.Event {
	return b.StateChangedEvent
}

func (b *Button) Configure(opts ...ButtonOpt) {
	for _, o := range opts {
		o(b)
	}
}

/** Focuser Interface - Start **/

func (b *Button) Focus(focused bool) {
	b.init.Do()
	b.GetWidget().FireFocusEvent(b, focused, img.Point{-1, -1})
	b.focused = focused
}

func (b *Button) IsFocused() bool {
	return b.focused
}

func (b *Button) TabOrder() int {
	return b.tabOrder
}

func (b *Button) GetFocus(direction FocusDirection) Focuser {
	return b.focusMap[direction]
}

func (b *Button) AddFocus(direction FocusDirection, focus Focuser) {
	b.focusMap[direction] = focus
}

/** Focuser Interface - End **/

func (b *Button) GetWidget() *Widget {
	b.init.Do()
	return b.widget
}

func (b *Button) PreferredSize() (int, int) {
	b.init.Do()

	w, h := 50, 50

	if b.container != nil && len(b.container.children) > 0 {
		w, h = b.container.PreferredSize()
	}

	if b.widget != nil && h < b.widget.MinHeight {
		h = b.widget.MinHeight
	}
	if b.widget != nil && w < b.widget.MinWidth {
		w = b.widget.MinWidth
	}

	iw, ih := b.computedParams.Image.Idle.MinSize()
	if w < iw {
		w = iw
	}
	if h < ih {
		h = ih
	}

	return w, h
}

func (b *Button) SetLocation(rect img.Rectangle) {
	b.init.Do()

	if b.IgnoreTransparentPixels && (b.mask == nil || b.widget.Rect == img.Rectangle{} || b.widget.Rect.Dx() != rect.Dx() || b.widget.Rect.Dy() != rect.Dy()) {
		maskImage := ebiten.NewImage(rect.Dx(), rect.Dy())
		b.computedParams.Image.Idle.Draw(maskImage, maskImage.Bounds().Dx(), maskImage.Bounds().Dy(), func(_ *ebiten.DrawImageOptions) {})

		wx := maskImage.Bounds().Dx()
		wy := maskImage.Bounds().Dy()
		b.mask = make([]byte, wx*wy*4)
		maskImage.ReadPixels(b.mask)
	}

	b.widget.Rect = rect
}

func (b *Button) RequestRelayout() {
	b.init.Do()

	if b.container != nil {
		b.container.RequestRelayout()
	}
}

func (b *Button) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	b.init.Do()

	if b.container != nil {
		b.container.SetupInputLayer(def)
	}
}

func (b *Button) Render(screen *ebiten.Image) {
	b.init.Do()

	if b.container != nil {
		w := b.container.GetWidget()
		w.Rect = b.widget.Rect
		w.Disabled = b.widget.Disabled
		b.container.RequestRelayout()
	}

	b.widget.Render(screen)
	b.draw(screen)

	if !b.DisableDefaultKeys {
		b.handleSubmit()
	} else {
		b.justSubmitted = false
	}

	if b.autoUpdateTextAndGraphic {
		if b.computedParams.GraphicImage != nil {
			if b.widget.Disabled && b.computedParams.GraphicImage.Disabled != nil {
				b.graphic.Image = b.computedParams.GraphicImage.Disabled
			} else {
				b.graphic.Image = b.computedParams.GraphicImage.Idle
			}
		}

		if b.text != nil {
			switch {
			case b.widget.Disabled && b.computedParams.TextColor.Disabled != nil:
				b.text.SetColor(b.computedParams.TextColor.Disabled)

			case (b.pressing && (b.hovering || b.KeepPressedOnExit) || (b.ToggleMode && b.state == WidgetChecked) || b.justSubmitted) && b.computedParams.TextColor.Pressed != nil:
				b.text.SetColor(b.computedParams.TextColor.Pressed)

			case (b.hovering || b.focused) && b.computedParams.TextColor.Hover != nil:
				b.text.SetColor(b.computedParams.TextColor.Hover)

			default:
				b.text.SetColor(b.computedParams.TextColor.Idle)
			}
		}
	}

	if b.container != nil {
		b.container.Render(screen)
	}
}

func (b *Button) Update() {
	b.init.Do()
	b.widget.Update()
	if b.container != nil {
		b.container.Update()
	}
}

func (b *Button) draw(screen *ebiten.Image) {
	i := b.computedParams.Image.Idle
	switch {
	case b.widget.Disabled:
		if b.computedParams.Image.Disabled != nil {
			i = b.computedParams.Image.Disabled
		}

	case b.pressing && (b.hovering || b.KeepPressedOnExit) || (b.ToggleMode && b.state == WidgetChecked) || b.justSubmitted:
		if b.computedParams.Image.Pressed != nil {
			i = b.computedParams.Image.Pressed
		}

	case b.focused, b.hovering:
		if b.ToggleMode && b.state == WidgetChecked || b.pressing && (b.hovering || b.KeepPressedOnExit) {
			if b.computedParams.Image.PressedHover != nil {
				i = b.computedParams.Image.PressedHover
			} else {
				i = b.computedParams.Image.Pressed
			}
		} else {
			if b.computedParams.Image.Hover != nil {
				i = b.computedParams.Image.Hover
			}
		}
	}

	if i != nil {
		i.Draw(screen, b.widget.Rect.Dx(), b.widget.Rect.Dy(), func(opts *ebiten.DrawImageOptions) {
			b.widget.drawImageOptions(opts)
			b.drawImageOptions(opts)
		})
	}
}

func (b *Button) Click() {
	b.init.Do()

	b.justSubmitted = true
	b.ClickedEvent.Fire(&ButtonClickedEventArgs{
		Button:  b,
		OffsetX: -1,
		OffsetY: -1,
	})
	if b.ToggleMode {
		if b.state == WidgetUnchecked {
			b.state = WidgetChecked
		} else {
			b.state = WidgetUnchecked
		}
		b.StateChangedEvent.Fire(&ButtonChangedEventArgs{
			Button:  b,
			State:   b.state,
			OffsetX: -1,
			OffsetY: -1,
		})
	}
}

// Press presses the button emulating a Mouse Left click.
func (b *Button) Press() {
	b.init.Do()

	offx := b.widget.Rect.Dx()
	offy := b.widget.Rect.Dy()

	// This means that there are some pixels that are not clickable.
	if b.mask != nil {
		offx /= 2
		offy /= 2
	}
	b.hovering = true
	b.widget.MouseButtonPressedEvent.Fire(&WidgetMouseButtonPressedEventArgs{
		Widget:  b.widget,
		Button:  ebiten.MouseButtonLeft,
		OffsetX: offx,
		OffsetY: offy,
	})
}

// Release releases the button emulating a Mouse Left release.
func (b *Button) Release() {
	b.init.Do()

	offx := b.widget.Rect.Dx()
	offy := b.widget.Rect.Dy()

	// This means that there are some pixels that are not clickable.
	if b.mask != nil {
		offx /= 2
		offy /= 2
	}
	b.hovering = false
	b.widget.MouseButtonReleasedEvent.Fire(&WidgetMouseButtonReleasedEventArgs{
		Widget:  b.widget,
		Inside:  true,
		Button:  ebiten.MouseButtonLeft,
		OffsetX: offx,
		OffsetY: offy,
	})
}

func (b *Button) handleSubmit() {
	if input.KeyPressed(ebiten.KeyEnter) || input.KeyPressed(ebiten.KeySpace) {
		if !b.justSubmitted && b.focused {
			b.justSubmitted = true
			b.Press()
		}
	} else if b.justSubmitted {
		b.Release()
		b.justSubmitted = false
	}
}

func (b *Button) drawImageOptions(opts *ebiten.DrawImageOptions) {
	if b.widget.Disabled && b.computedParams.Image.Disabled == nil {
		opts.ColorM.Scale(1, 1, 1, 0.35)
	}
}

func (b *Button) Text() *Text {
	b.init.Do()
	return b.text
}

func (b *Button) initText() {
	if b.computedParams.TextColor == nil {
		return // Nothing to do.
	}

	if b.text != nil {
		b.text.SetFace(b.computedParams.TextFace)
		b.text.SetColor(b.computedParams.TextColor.Idle)
		b.text.horizontalPosition = b.computedParams.HTextPosition
		b.text.verticalPosition = b.computedParams.VTextPosition
		b.text.widget.LayoutData = AnchorLayoutData{
			HorizontalPosition: AnchorLayoutPosition(b.computedParams.HTextPosition),
			VerticalPosition:   AnchorLayoutPosition(b.computedParams.VTextPosition),
		}
		if aLayout, ok := b.container.layout.(*AnchorLayout); ok {
			aLayout.padding = *b.computedParams.TextPadding
		}
	} else {

		// We're expecting all 3 options to be present: label, font face and color.
		// TODO: add some sort of the error checking/reporting here.
		// Even if users use a Text() 3-in-one API, they can pass nil or something.

		b.container = NewContainer(
			ContainerOpts.Layout(NewAnchorLayout(AnchorLayoutOpts.Padding(*b.computedParams.TextPadding))),
			ContainerOpts.AutoDisableChildren(),
		)

		b.text = NewText(
			TextOpts.WidgetOpts(WidgetOpts.LayoutData(AnchorLayoutData{
				HorizontalPosition: AnchorLayoutPosition(b.computedParams.HTextPosition),
				VerticalPosition:   AnchorLayoutPosition(b.computedParams.VTextPosition),
			})),
			TextOpts.Text(b.textLabel, b.computedParams.TextFace, b.computedParams.TextColor.Idle),
			TextOpts.ProcessBBCode(b.textProcessBBCode),
			TextOpts.Position(b.computedParams.HTextPosition, b.computedParams.VTextPosition),
		)
		b.container.AddChild(b.text)
		b.container.Validate()

		b.autoUpdateTextAndGraphic = true
	}
}

func (b *Button) createWidget() {
	b.widget = NewWidget(append([]WidgetOpt{
		WidgetOpts.TrackHover(true),
		WidgetOpts.CursorEnterHandler(func(args *WidgetCursorEnterEventArgs) {
			if b.mask == nil {
				if !b.widget.Disabled {
					b.hovering = true
				}
				if b.hovering {
					b.CursorEnteredEvent.Fire(&ButtonHoverEventArgs{
						Button:  b,
						Entered: true,
						OffsetX: args.OffsetX,
						OffsetY: args.OffsetY,
						DiffX:   0,
						DiffY:   0,
					})
				}
			}
		}),

		WidgetOpts.CursorMoveHandler(func(args *WidgetCursorMoveEventArgs) {
			if b.onMask(args.OffsetX, args.OffsetY) {
				if !b.hovering {
					b.CursorEnteredEvent.Fire(&ButtonHoverEventArgs{
						Button:  b,
						Entered: true,
						OffsetX: args.OffsetX,
						OffsetY: args.OffsetY,
						DiffX:   0,
						DiffY:   0,
					})
				}
				if !b.widget.Disabled {
					b.hovering = true
				}
				b.CursorMovedEvent.Fire(&ButtonHoverEventArgs{
					Button:  b,
					Entered: false,
					OffsetX: args.OffsetX,
					OffsetY: args.OffsetY,
					DiffX:   args.DiffX,
					DiffY:   args.DiffY,
				})
			} else if b.hovering {
				b.hovering = false
				b.CursorExitedEvent.Fire(&ButtonHoverEventArgs{
					Button:  b,
					Entered: false,
					OffsetX: args.OffsetX,
					OffsetY: args.OffsetY,
					DiffX:   0,
					DiffY:   0,
				})
			}
		}),

		WidgetOpts.CursorExitHandler(func(args *WidgetCursorExitEventArgs) {
			if b.hovering || b.mask == nil {
				b.hovering = false
				b.CursorExitedEvent.Fire(&ButtonHoverEventArgs{
					Button:  b,
					Entered: false,
					OffsetX: args.OffsetX,
					OffsetY: args.OffsetY,
					DiffX:   0,
					DiffY:   0,
				})
			}
		}),

		WidgetOpts.MouseButtonPressedHandler(func(args *WidgetMouseButtonPressedEventArgs) {
			if b.onMask(args.OffsetX, args.OffsetY) {
				if !b.widget.Disabled && args.Button == ebiten.MouseButtonLeft {
					b.pressing = true
					b.PressedEvent.Fire(&ButtonPressedEventArgs{
						Button:  b,
						OffsetX: args.OffsetX,
						OffsetY: args.OffsetY,
					})
				}
			}
		}),

		WidgetOpts.MouseButtonReleasedHandler(func(args *WidgetMouseButtonReleasedEventArgs) {
			if b.pressing && !b.widget.Disabled && args.Button == ebiten.MouseButtonLeft {
				inside := args.Inside && b.onMask(args.OffsetX, args.OffsetY)

				b.ReleasedEvent.Fire(&ButtonReleasedEventArgs{
					Button:  b,
					Inside:  inside,
					OffsetX: args.OffsetX,
					OffsetY: args.OffsetY,
				})
				if inside {
					b.ClickedEvent.Fire(&ButtonClickedEventArgs{
						Button:  b,
						OffsetX: args.OffsetX,
						OffsetY: args.OffsetY,
					})
					if b.ToggleMode {
						if b.state == WidgetUnchecked {
							b.state = WidgetChecked
						} else {
							b.state = WidgetUnchecked
						}
						b.StateChangedEvent.Fire(&ButtonChangedEventArgs{
							Button:  b,
							State:   b.state,
							OffsetX: args.OffsetX,
							OffsetY: args.OffsetY,
						})
					}
				}
			}

			b.pressing = false
		}),
	}, b.widgetOpts...)...)
	b.widgetOpts = nil
}

func (b *Button) onMask(x, y int) bool {
	if b.mask == nil {
		return true
	}
	i := ((x * 4) + (y * b.widget.Rect.Dx() * 4) + 3)
	if len(b.mask)-1 < i {
		return false
	}
	return (b.mask[i] > 0)
}
