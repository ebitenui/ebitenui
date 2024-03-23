package widget

import (
	img "image"
	"image/color"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type Button struct {
	Image             *ButtonImage
	KeepPressedOnExit bool
	ToggleMode        bool
	GraphicImage      *ButtonImageImage
	TextColor         *ButtonTextColor

	PressedEvent       *event.Event
	ReleasedEvent      *event.Event
	ClickedEvent       *event.Event
	CursorEnteredEvent *event.Event
	CursorExitedEvent  *event.Event
	StateChangedEvent  *event.Event

	widgetOpts               []WidgetOpt
	autoUpdateTextAndGraphic bool
	vTextPosition            TextPosition
	hTextPosition            TextPosition
	textPadding              Insets
	graphicPadding           Insets

	init      *MultiOnce
	widget    *Widget
	container *Container
	graphic   *Graphic
	text      *Text
	textLabel string
	textFace  font.Face
	hovering  bool
	pressing  bool
	state     WidgetState

	tabOrder      int
	focused       bool
	justSubmitted bool
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
	Button *Button
}
type ButtonHoverEventArgs struct {
	Button  *Button
	Entered bool
}
type ButtonChangedEventArgs struct {
	Button *Button
	State  WidgetState
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
		hTextPosition: TextPositionCenter,
		vTextPosition: TextPositionCenter,

		PressedEvent:       &event.Event{},
		ReleasedEvent:      &event.Event{},
		ClickedEvent:       &event.Event{},
		CursorEnteredEvent: &event.Event{},
		CursorExitedEvent:  &event.Event{},
		StateChangedEvent:  &event.Event{},

		init: &MultiOnce{},
	}

	b.init.Append(b.createWidget)

	for _, o := range opts {
		o(b)
	}

	return b
}

func (o ButtonOptions) WidgetOpts(opts ...WidgetOpt) ButtonOpt {
	return func(b *Button) {
		b.widgetOpts = append(b.widgetOpts, opts...)
	}
}

func (o ButtonOptions) Image(i *ButtonImage) ButtonOpt {
	return func(b *Button) {
		b.Image = i
	}
}

// Text combines three options: TextLabel, TextFace and TextColor.
// It can be used for the inline configurations of Text object while
// separate functions are useful for a multi-step configuration.
func (o ButtonOptions) Text(label string, face font.Face, color *ButtonTextColor) ButtonOpt {
	return func(b *Button) {
		b.textLabel = label
		b.textFace = face
		b.TextColor = color
	}
}

func (o ButtonOptions) TextLabel(label string) ButtonOpt {
	return func(b *Button) {
		b.textLabel = label
	}
}

func (o ButtonOptions) TextFace(face font.Face) ButtonOpt {
	return func(b *Button) {
		b.textFace = face
	}
}

func (o ButtonOptions) TextColor(color *ButtonTextColor) ButtonOpt {
	return func(b *Button) {
		b.TextColor = color
	}
}

// TODO: add parameter for image position (start/end)
func (o ButtonOptions) TextAndImage(label string, face font.Face, image *ButtonImageImage, color *ButtonTextColor) ButtonOpt {
	return func(b *Button) {
		b.init.Append(func() {
			b.container = NewContainer(
				ContainerOpts.Layout(NewAnchorLayout(AnchorLayoutOpts.Padding(b.textPadding))),
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
				TextOpts.Text(label, face, color.Idle))
			c.AddChild(b.text)

			b.graphic = NewGraphic(
				GraphicOpts.WidgetOpts(WidgetOpts.LayoutData(RowLayoutData{
					Stretch: true,
				})),
				GraphicOpts.Image(image.Idle))
			c.AddChild(b.graphic)

			b.autoUpdateTextAndGraphic = true
			b.GraphicImage = image
			b.TextColor = color
		})
	}
}

// TextPosition sets the horizontal and vertical position of the text within the button.
// Default is TextPositionCenter for both.
func (o ButtonOptions) TextPosition(h TextPosition, v TextPosition) ButtonOpt {
	return func(b *Button) {
		b.hTextPosition = h
		b.vTextPosition = v
	}
}

func (o ButtonOptions) TextPadding(p Insets) ButtonOpt {
	return func(b *Button) {
		b.textPadding = p
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
				ContainerOpts.Layout(NewAnchorLayout(AnchorLayoutOpts.Padding(b.graphicPadding))),
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
		b.graphicPadding = i
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

func (o ButtonOptions) PressedHandler(f ButtonPressedHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.PressedEvent.AddHandler(func(args interface{}) {
			f(args.(*ButtonPressedEventArgs))
		})
	}
}

func (o ButtonOptions) ReleasedHandler(f ButtonReleasedHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.ReleasedEvent.AddHandler(func(args interface{}) {
			f(args.(*ButtonReleasedEventArgs))
		})
	}
}

func (o ButtonOptions) ClickedHandler(f ButtonClickedHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.ClickedEvent.AddHandler(func(args interface{}) {
			f(args.(*ButtonClickedEventArgs))
		})
	}
}

func (o ButtonOptions) CursorEnteredHandler(f ButtonCursorHoverHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.CursorEnteredEvent.AddHandler(func(args interface{}) {
			f(args.(*ButtonHoverEventArgs))
		})
	}
}

func (o ButtonOptions) CursorExitedHandler(f ButtonCursorHoverHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.CursorExitedEvent.AddHandler(func(args interface{}) {
			f(args.(*ButtonHoverEventArgs))
		})
	}
}

func (o ButtonOptions) StateChangedHandler(f ButtonChangedHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.StateChangedEvent.AddHandler(func(args interface{}) {
			f(args.(*ButtonChangedEventArgs))
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
			Button: b,
			State:  b.state,
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

	iw, ih := b.Image.Idle.MinSize()
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

func (b *Button) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	b.init.Do()

	if b.container != nil {
		w := b.container.GetWidget()
		w.Rect = b.widget.Rect
		w.Disabled = b.widget.Disabled
		b.container.RequestRelayout()
	}

	b.widget.Render(screen, def)
	b.handleSubmit()
	b.draw(screen)

	if b.autoUpdateTextAndGraphic {
		if b.graphic != nil {
			if b.widget.Disabled {
				b.graphic.Image = b.GraphicImage.Disabled
			} else {
				b.graphic.Image = b.GraphicImage.Idle
			}
		}

		if b.text != nil {
			if b.widget.Disabled {
				b.text.Color = b.TextColor.Disabled
			} else {
				b.text.Color = b.TextColor.Idle
			}
		}
	}

	if b.container != nil {
		b.container.Render(screen, def)
	}
}

func (b *Button) draw(screen *ebiten.Image) {
	i := b.Image.Idle
	switch {
	case b.widget.Disabled:
		if b.Image.Disabled != nil {
			i = b.Image.Disabled
		}
	case b.focused, b.hovering:
		if b.ToggleMode && b.state == WidgetChecked || b.pressing && (b.hovering || b.KeepPressedOnExit) {
			if b.Image.PressedHover != nil {
				i = b.Image.PressedHover
			} else {
				i = b.Image.Pressed
			}
		} else {
			if b.Image.Hover != nil {
				i = b.Image.Hover
			}
		}
	case b.pressing && (b.hovering || b.KeepPressedOnExit) || (b.ToggleMode && b.state == WidgetChecked):
		if b.Image.Pressed != nil {
			i = b.Image.Pressed
		}

	}

	if i != nil {
		i.Draw(screen, b.widget.Rect.Dx(), b.widget.Rect.Dy(), func(opts *ebiten.DrawImageOptions) {
			b.widget.drawImageOptions(opts)
			b.drawImageOptions(opts)
		})
	}
}

func (b *Button) handleSubmit() {
	if input.KeyPressed(ebiten.KeyEnter) || input.KeyPressed(ebiten.KeySpace) {
		if !b.justSubmitted && b.focused {
			b.ClickedEvent.Fire(&ButtonClickedEventArgs{
				Button: b,
			})
			if b.ToggleMode {
				if b.state == WidgetUnchecked {
					b.state = WidgetChecked
				} else {
					b.state = WidgetUnchecked
				}
				b.StateChangedEvent.Fire(&ButtonChangedEventArgs{
					Button: b,
					State:  b.state,
				})
			}
			b.justSubmitted = true
		}
	} else {
		b.justSubmitted = false
	}
}

func (b *Button) drawImageOptions(opts *ebiten.DrawImageOptions) {
	if b.widget.Disabled && b.Image.Disabled == nil {
		opts.ColorM.Scale(1, 1, 1, 0.35)
	}
}

func (b *Button) Text() *Text {
	b.init.Do()
	return b.text
}

func (b *Button) initText() {
	if b.TextColor == nil {
		return // Nothing to do
	}

	// We're expecting all 3 options to be present: label, font face and color.
	// TODO: add some sort of the error checking/reporting here.
	// Even if users use a Text() 3-in-one API, they can pass nil or something.

	b.container = NewContainer(
		ContainerOpts.Layout(NewAnchorLayout(AnchorLayoutOpts.Padding(b.textPadding))),
		ContainerOpts.AutoDisableChildren(),
	)

	b.text = NewText(
		TextOpts.WidgetOpts(WidgetOpts.LayoutData(AnchorLayoutData{
			HorizontalPosition: AnchorLayoutPosition(b.hTextPosition),
			VerticalPosition:   AnchorLayoutPosition(b.vTextPosition),
		})),
		TextOpts.Text(b.textLabel, b.textFace, b.TextColor.Idle),
	)
	b.container.AddChild(b.text)

	b.autoUpdateTextAndGraphic = true
}

func (b *Button) createWidget() {
	b.widget = NewWidget(append(b.widgetOpts, []WidgetOpt{
		WidgetOpts.CursorEnterHandler(func(_ *WidgetCursorEnterEventArgs) {
			if !b.widget.Disabled {
				b.hovering = true
			}
			b.CursorEnteredEvent.Fire(&ButtonHoverEventArgs{
				Button:  b,
				Entered: true,
			})
		}),

		WidgetOpts.CursorExitHandler(func(_ *WidgetCursorExitEventArgs) {
			b.hovering = false
			b.CursorExitedEvent.Fire(&ButtonHoverEventArgs{
				Button:  b,
				Entered: false,
			})
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

				if args.Inside {
					b.ClickedEvent.Fire(&ButtonClickedEventArgs{
						Button: b,
					})
					if b.ToggleMode {
						if b.state == WidgetUnchecked {
							b.state = WidgetChecked
						} else {
							b.state = WidgetUnchecked
						}
						b.StateChangedEvent.Fire(&ButtonChangedEventArgs{
							Button: b,
							State:  b.state,
						})
					}
				}
			}

			b.pressing = false
		}),
	}...)...)
	b.widgetOpts = nil

	b.initText()
}
