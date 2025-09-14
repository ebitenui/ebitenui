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
	GraphicImage *GraphicImage
	TextColor    *ButtonTextColor

	TextPosition   *TextPositioning
	TextPadding    *Insets
	TextFace       *text.Face
	GraphicPadding *Insets
	MinSize        *img.Point
}

type Button struct {
	definedParams           ButtonParams
	computedParams          ButtonParams
	IgnoreTransparentPixels bool
	KeepPressedOnExit       bool
	ToggleMode              bool
	GraphicImage            *GraphicImage

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
	Idle            *image.NineSlice
	Hover           *image.NineSlice
	Pressed         *image.NineSlice
	PressedHover    *image.NineSlice
	Disabled        *image.NineSlice
	PressedDisabled *image.NineSlice
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
	b.init.Do()
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

	b.initWidget()
}

func (b *Button) populateComputedParams() {
	btnParams := ButtonParams{
		TextPosition: &TextPositioning{
			HTextPosition: TextPositionCenter,
			VTextPosition: TextPositionCenter,
		},
		GraphicPadding: &Insets{},
		TextPadding:    &Insets{},
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
				btnParams.GraphicImage = &GraphicImage{
					Idle:     theme.ButtonTheme.GraphicImage.Idle,
					Disabled: theme.ButtonTheme.GraphicImage.Disabled,
				}
			}
			if theme.ButtonTheme.GraphicPadding != nil {
				btnParams.GraphicPadding = theme.ButtonTheme.GraphicPadding
			}
			if theme.ButtonTheme.TextPosition != nil {
				btnParams.TextPosition.HTextPosition = theme.ButtonTheme.TextPosition.HTextPosition
				btnParams.TextPosition.VTextPosition = theme.ButtonTheme.TextPosition.VTextPosition
			}
			if theme.ButtonTheme.TextPadding != nil {
				btnParams.TextPadding = theme.ButtonTheme.TextPadding
			}
			if theme.ButtonTheme.TextFace != nil {
				btnParams.TextFace = theme.ButtonTheme.TextFace
			}
			if theme.ButtonTheme.MinSize != nil {
				btnParams.MinSize = theme.ButtonTheme.MinSize
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
		if btnParams.Image == nil {
			btnParams.Image = b.definedParams.Image
		} else {
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
	}

	if b.definedParams.GraphicImage != nil {
		if btnParams.GraphicImage == nil {
			btnParams.GraphicImage = b.definedParams.GraphicImage
		} else {
			if b.definedParams.GraphicImage.Idle != nil {
				btnParams.GraphicImage.Idle = b.definedParams.GraphicImage.Idle
			}
			if b.definedParams.GraphicImage.Disabled != nil {
				btnParams.GraphicImage.Disabled = b.definedParams.GraphicImage.Disabled
			}
		}
	}
	if b.definedParams.GraphicPadding != nil {
		btnParams.GraphicPadding = b.definedParams.GraphicPadding
	}
	if b.definedParams.TextPosition != nil {
		btnParams.TextPosition.HTextPosition = b.definedParams.TextPosition.HTextPosition
		btnParams.TextPosition.VTextPosition = b.definedParams.TextPosition.VTextPosition
	}
	if b.definedParams.TextFace != nil {
		btnParams.TextFace = b.definedParams.TextFace
	}
	if b.definedParams.TextPadding != nil {
		btnParams.TextPadding = b.definedParams.TextPadding
	}
	if b.definedParams.MinSize != nil {
		btnParams.MinSize = b.definedParams.MinSize
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
func (o ButtonOptions) TextAndImage(label string, face *text.Face, image *GraphicImage, color *ButtonTextColor) ButtonOpt {
	return func(b *Button) {
		b.autoUpdateTextAndGraphic = true
		b.textLabel = label
		b.definedParams.TextFace = face
		b.definedParams.TextColor = color
		b.definedParams.GraphicImage = image
		b.definedParams.TextColor = color
	}
}

// TextPosition sets the horizontal and vertical position of the text within the button.
// Default is TextPositionCenter for both.
func (o ButtonOptions) TextPosition(h TextPosition, v TextPosition) ButtonOpt {
	return func(b *Button) {
		b.definedParams.TextPosition = &TextPositioning{
			VTextPosition: v,
			HTextPosition: h,
		}
	}
}

func (o ButtonOptions) TextPadding(p *Insets) ButtonOpt {
	return func(b *Button) {
		b.definedParams.TextPadding = p
	}
}

func (o ButtonOptions) Graphic(image *GraphicImage) ButtonOpt {
	return func(b *Button) {
		b.definedParams.GraphicImage = image
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

	// If we are ignoring transparent pixels and the mask isn't set to the current Image/Size, rebuild the mask.
	if b.IgnoreTransparentPixels && (b.GetWidget().mask == nil || b.widget.Rect == img.Rectangle{} || b.widget.Rect.Dx() != rect.Dx() || b.widget.Rect.Dy() != rect.Dy()) {
		maskImage := ebiten.NewImage(rect.Dx(), rect.Dy())
		b.computedParams.Image.Idle.Draw(maskImage, maskImage.Bounds().Dx(), maskImage.Bounds().Dy(), func(_ *ebiten.DrawImageOptions) {})

		wx := maskImage.Bounds().Dx()
		wy := maskImage.Bounds().Dy()
		b.GetWidget().mask = make([]byte, wx*wy*4)
		maskImage.ReadPixels(b.GetWidget().mask)
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

	if b.autoUpdateTextAndGraphic {

		// We set the defaults first and then if needed
		// they'll be overwritten by the other states
		if b.text != nil {
			b.text.SetColor(b.computedParams.TextColor.Idle)
		}

		if b.computedParams.GraphicImage != nil {
			if b.widget.Disabled && b.computedParams.GraphicImage.Disabled != nil {
				b.graphic.Image = b.computedParams.GraphicImage.Disabled
			} else {
				b.graphic.Image = b.computedParams.GraphicImage.Idle
			}
		}

		switch {
		case b.widget.Disabled:
			if b.text != nil && b.computedParams.TextColor.Disabled != nil {
				b.text.SetColor(b.computedParams.TextColor.Disabled)
			}
			if b.GraphicImage != nil && b.GraphicImage.Disabled != nil {
				b.graphic.Image = b.GraphicImage.Disabled
			}

		case b.pressing && (b.hovering || b.KeepPressedOnExit) || (b.ToggleMode && b.state == WidgetChecked) || b.justSubmitted:
			if b.text != nil && b.computedParams.TextColor.Pressed != nil {
				b.text.SetColor(b.computedParams.TextColor.Pressed)
			}
			if b.GraphicImage != nil && b.GraphicImage.Pressed != nil {
				b.graphic.Image = b.GraphicImage.Pressed
			}

		case b.hovering || b.focused:
			if b.computedParams.TextColor.Hover != nil {
				b.text.SetColor(b.computedParams.TextColor.Hover)
			}
			if b.computedParams.GraphicImage != nil && b.computedParams.GraphicImage.Hover != nil {
				b.graphic.Image = b.computedParams.GraphicImage.Hover
			}
		default:
			b.text.SetColor(b.computedParams.TextColor.Idle)
		}
	}

	if b.container != nil {
		b.container.Render(screen)
	}
}

func (b *Button) Update(updObj *UpdateObject) {
	b.init.Do()
	b.widget.Update(updObj)
	if b.container != nil {
		b.container.Update(updObj)
	}

	if !b.DisableDefaultKeys {
		b.handleSubmit()
	} else {
		b.justSubmitted = false
	}

}

func (b *Button) draw(screen *ebiten.Image) {
	i := b.computedParams.Image.Idle
	pressed := (b.pressing && (b.hovering || b.KeepPressedOnExit)) || (b.ToggleMode && b.state == WidgetChecked)
	switch {
	case b.widget.Disabled:
		if b.computedParams.Image.Disabled != nil {
			i = b.computedParams.Image.Disabled
		}
		if pressed {
			if b.computedParams.Image.PressedDisabled != nil {
				i = b.computedParams.Image.PressedDisabled
			}
		}
	case b.focused, b.hovering:
		if pressed {
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
	case pressed, b.justSubmitted:
		if b.computedParams.Image.Pressed != nil {
			i = b.computedParams.Image.Pressed
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

	b.pressing = true
	b.widget.MouseButtonClickedEvent.Fire(&WidgetMouseButtonClickedEventArgs{
		Widget:  b.widget,
		Button:  ebiten.MouseButtonLeft,
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
	if b.GetWidget().mask != nil {
		offx /= 2
		offy /= 2
	}
	b.hovering = true
	b.pressing = true
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
	if b.GetWidget().mask != nil {
		offx /= 2
		offy /= 2
	}
	if !b.ToggleMode {
		b.hovering = false
	}
	b.widget.MouseButtonReleasedEvent.Fire(&WidgetMouseButtonReleasedEventArgs{
		Widget:  b.widget,
		Inside:  true,
		Button:  ebiten.MouseButtonLeft,
		OffsetX: offx,
		OffsetY: offy,
	})
	if b.pressing {
		b.widget.MouseButtonClickedEvent.Fire(&WidgetMouseButtonClickedEventArgs{
			Widget:  b.widget,
			Button:  ebiten.MouseButtonLeft,
			OffsetX: offx,
			OffsetY: offy,
		})
		b.pressing = false
	}
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

func (b *Button) SetText(text string) {
	b.init.Do()
	b.textLabel = text
	if b.text != nil {
		b.text.Label = text
	}
}

func (b *Button) SetGraphicImage(image *GraphicImage) {
	b.init.Do()
	b.definedParams.GraphicImage = image
	b.computedParams.GraphicImage = image
}

func (b *Button) initWidget() {

	if b.computedParams.MinSize != nil {
		b.widget.MinWidth = b.computedParams.MinSize.X
		b.widget.MinHeight = b.computedParams.MinSize.Y
	}

	b.container = NewContainer(
		ContainerOpts.Layout(NewAnchorLayout(AnchorLayoutOpts.Padding(b.computedParams.TextPadding))),
		ContainerOpts.AutoDisableChildren(),
	)
	var textLayoutData any = AnchorLayoutData{
		StretchHorizontal: true,
		StretchVertical:   true,
	}
	if b.computedParams.GraphicImage != nil {
		textLayoutData = RowLayoutData{Stretch: true}
	}
	if b.computedParams.TextColor != nil {
		if b.text != nil {
			b.text.SetFace(b.computedParams.TextFace)
			b.text.SetColor(b.computedParams.TextColor.Idle)
			b.text.SetPosition(b.computedParams.TextPosition)
			b.text.widget.LayoutData = textLayoutData
		} else {
			b.text = NewText(
				TextOpts.WidgetOpts(WidgetOpts.LayoutData(textLayoutData)),
				TextOpts.Text(b.textLabel, b.computedParams.TextFace, b.computedParams.TextColor.Idle),
				TextOpts.ProcessBBCode(b.textProcessBBCode),
				TextOpts.Position(b.computedParams.TextPosition.HTextPosition, b.computedParams.TextPosition.VTextPosition),
			)
		}
		b.autoUpdateTextAndGraphic = true
	}
	if b.computedParams.GraphicImage != nil {
		c := NewContainer(
			ContainerOpts.WidgetOpts(WidgetOpts.LayoutData(AnchorLayoutData{
				StretchHorizontal: true,
				StretchVertical:   true,
			})),
			ContainerOpts.Layout(NewRowLayout(
				RowLayoutOpts.Spacing(10),
				RowLayoutOpts.Padding(b.computedParams.GraphicPadding),
				RowLayoutOpts.Direction(DirectionHorizontal),
			)),
			ContainerOpts.AutoDisableChildren(),
		)
		b.container.AddChild(c)
		b.graphic = NewGraphic(
			GraphicOpts.WidgetOpts(WidgetOpts.LayoutData(RowLayoutData{
				Stretch: false,
			})),
			GraphicOpts.Image(b.computedParams.GraphicImage.Idle),
		)
		if b.text != nil {
			c.AddChild(b.text)
		}
		c.AddChild(b.graphic)
	} else if b.text != nil {
		b.container.AddChild(b.text)
	}
	b.container.Validate()
}

func (b *Button) createWidget() {
	b.widget = NewWidget(append([]WidgetOpt{
		WidgetOpts.TrackHover(true),
		WidgetOpts.CursorEnterHandler(func(args *WidgetCursorEnterEventArgs) {
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

		}),

		WidgetOpts.CursorMoveHandler(func(args *WidgetCursorMoveEventArgs) {
			b.CursorMovedEvent.Fire(&ButtonHoverEventArgs{
				Button:  b,
				Entered: false,
				OffsetX: args.OffsetX,
				OffsetY: args.OffsetY,
				DiffX:   args.DiffX,
				DiffY:   args.DiffY,
			})
		}),

		WidgetOpts.CursorExitHandler(func(args *WidgetCursorExitEventArgs) {
			if b.hovering {
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

			if !b.widget.Disabled && args.Button == ebiten.MouseButtonLeft {
				b.pressing = true
				b.PressedEvent.Fire(&ButtonPressedEventArgs{
					Button:  b,
					OffsetX: args.OffsetX,
					OffsetY: args.OffsetY,
				})
			}

		}),

		WidgetOpts.MouseButtonReleasedHandler(func(args *WidgetMouseButtonReleasedEventArgs) {
			if b.pressing && !b.widget.Disabled && args.Button == ebiten.MouseButtonLeft {

				b.ReleasedEvent.Fire(&ButtonReleasedEventArgs{
					Button:  b,
					Inside:  args.Inside,
					OffsetX: args.OffsetX,
					OffsetY: args.OffsetY,
				})
			}

			b.pressing = false
		}),

		WidgetOpts.MouseButtonClickedHandler(func(args *WidgetMouseButtonClickedEventArgs) {
			if !b.widget.Disabled && args.Button == ebiten.MouseButtonLeft {
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
		}),
	}, b.widgetOpts...)...)
	b.widgetOpts = nil
}
