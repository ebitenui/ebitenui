package widget

import (
	img "image"
	"image/color"

	"github.com/blizzy78/ebitenui/event"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/input"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

type Button struct {
	Image             *ButtonImage
	KeepPressedOnExit bool
	GraphicImage      *ButtonImageImage

	PressedEvent  *event.Event
	ReleasedEvent *event.Event
	ClickedEvent  *event.Event

	widgetOpts               []WidgetOpt
	autoUpdateTextAndGraphic bool
	textColor                *ButtonTextColor

	init      *MultiOnce
	widget    *Widget
	container *Container
	graphic   *Graphic
	text      *Text
	hovering  bool
	pressing  bool
}

type ButtonOpt func(b *Button)

type ButtonImage struct {
	Idle     *image.NineSlice
	Hover    *image.NineSlice
	Pressed  *image.NineSlice
	Disabled *image.NineSlice
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

type ButtonPressedHandlerFunc func(args *ButtonPressedEventArgs)

type ButtonReleasedHandlerFunc func(args *ButtonReleasedEventArgs)

type ButtonClickedHandlerFunc func(args *ButtonClickedEventArgs)

const ButtonOpts = buttonOpts(true)

type buttonOpts bool

func NewButton(opts ...ButtonOpt) *Button {
	b := &Button{
		PressedEvent:  &event.Event{},
		ReleasedEvent: &event.Event{},
		ClickedEvent:  &event.Event{},

		Image:        &ButtonImage{},
		GraphicImage: &ButtonImageImage{},
		textColor:    &ButtonTextColor{},

		init: &MultiOnce{},
	}

	b.init.Append(b.createWidget)

	for _, o := range opts {
		o(b)
	}

	return b
}

func (o buttonOpts) WithWidgetOpts(opts ...WidgetOpt) ButtonOpt {
	return func(b *Button) {
		b.widgetOpts = append(b.widgetOpts, opts...)
	}
}

func (o buttonOpts) WithImage(i *ButtonImage) ButtonOpt {
	return func(b *Button) {
		b.Image = i
	}
}

func (o buttonOpts) WithTextSimpleLeft(label string, face font.Face, color *ButtonTextColor, padding Insets) ButtonOpt {
	return func(b *Button) {
		b.init.Append(func() {
			b.container = NewContainer(
				ContainerOpts.WithLayout(NewRowLayout(
					RowLayoutOpts.WithPadding(padding))),
				ContainerOpts.WithAutoDisableChildren())

			b.text = NewText(TextOpts.WithText(label, face, color.Idle))
			b.container.AddChild(b.text)

			b.autoUpdateTextAndGraphic = true
			b.textColor = color
		})
	}
}

func (o buttonOpts) WithText(label string, face font.Face, color *ButtonTextColor) ButtonOpt {
	return func(b *Button) {
		b.init.Append(func() {
			b.container = NewContainer(
				ContainerOpts.WithLayout(NewFillLayout(
					FillLayoutOpts.WithPadding(Insets{
						Left:   10,
						Right:  10,
						Top:    6,
						Bottom: 6,
					}))),
				ContainerOpts.WithAutoDisableChildren())

			b.text = NewText(
				TextOpts.WithText(label, face, color.Idle),
				TextOpts.WithPosition(TextPositionCenter))
			b.container.AddChild(b.text)

			b.autoUpdateTextAndGraphic = true
			b.textColor = color
		})
	}
}

// TODO: add parameter for image position (start/end)
func (o buttonOpts) WithTextAndImage(label string, face font.Face, image *ButtonImageImage, color *ButtonTextColor) ButtonOpt {
	return func(b *Button) {
		b.init.Append(func() {
			b.container = NewContainer(
				ContainerOpts.WithLayout(NewRowLayout(
					RowLayoutOpts.WithDirection(DirectionVertical),
					RowLayoutOpts.WithPadding(Insets{
						Left:   10,
						Right:  10,
						Top:    6,
						Bottom: 6,
					}))),
				ContainerOpts.WithAutoDisableChildren())

			c := NewContainer(
				ContainerOpts.WithWidgetOpts(WidgetOpts.WithLayoutData(&RowLayoutData{
					Position: RowLayoutPositionCenter,
				})),
				ContainerOpts.WithLayout(NewRowLayout(
					RowLayoutOpts.WithSpacing(10))),
				ContainerOpts.WithAutoDisableChildren())
			b.container.AddChild(c)

			b.text = NewText(
				TextOpts.WithWidgetOpts(WidgetOpts.WithLayoutData(&RowLayoutData{
					Stretch: true,
				})),
				TextOpts.WithText(label, face, color.Idle))
			c.AddChild(b.text)

			b.graphic = NewGraphic(
				GraphicOpts.WithWidgetOpts(WidgetOpts.WithLayoutData(&RowLayoutData{
					Stretch: true,
				})),
				GraphicOpts.WithImage(image.Idle))
			c.AddChild(b.graphic)

			b.autoUpdateTextAndGraphic = true
			b.GraphicImage = image
			b.textColor = color
		})
	}
}

func (o buttonOpts) WithGraphic(i *ebiten.Image) ButtonOpt {
	return o.withGraphic(GraphicOpts.WithImage(i))
}

func (o buttonOpts) WithGraphicNineSlice(i *image.NineSlice) ButtonOpt {
	return o.withGraphic(GraphicOpts.WithImageNineSlice(i))
}

func (o buttonOpts) withGraphic(opt GraphicOpt) ButtonOpt {
	return func(b *Button) {
		b.init.Append(func() {
			b.container = NewContainer(
				ContainerOpts.WithLayout(NewFillLayout(FillLayoutOpts.WithPadding(NewInsetsSimple(4)))),
				ContainerOpts.WithAutoDisableChildren())

			b.graphic = NewGraphic(opt)
			b.container.AddChild(b.graphic)

			b.autoUpdateTextAndGraphic = true
		})
	}
}

func (o buttonOpts) WithKeepPressedOnExit() ButtonOpt {
	return func(b *Button) {
		b.KeepPressedOnExit = true
	}
}

func (o buttonOpts) WithPressedHandler(f ButtonPressedHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.PressedEvent.AddHandler(func(args interface{}) {
			f(args.(*ButtonPressedEventArgs))
		})
	}
}

func (o buttonOpts) WithReleasedHandler(f ButtonReleasedHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.ReleasedEvent.AddHandler(func(args interface{}) {
			f(args.(*ButtonReleasedEventArgs))
		})
	}
}

func (o buttonOpts) WithClickedHandler(f ButtonClickedHandlerFunc) ButtonOpt {
	return func(b *Button) {
		b.ClickedEvent.AddHandler(func(args interface{}) {
			f(args.(*ButtonClickedEventArgs))
		})
	}
}

func (b *Button) GetWidget() *Widget {
	b.init.Do()
	return b.widget
}

func (b *Button) PreferredSize() (int, int) {
	b.init.Do()

	if b.container == nil || len(b.container.children) == 0 {
		return 50, 50
	}

	return b.container.PreferredSize()
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

	if b.pressing {
		def(func(def input.DeferredSetupInputLayerFunc) {
			b.widget.ElevateToNewInputLayer(&input.Layer{
				DebugLabel: "button pressed",
				EventTypes: input.LayerEventTypeAll,
				BlockLower: true,
				FullScreen: true,
			})
		})
	}
}

func (b *Button) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	b.init.Do()

	if b.container != nil {
		w := b.container.GetWidget()
		w.Rect = b.widget.Rect
		w.Disabled = b.widget.Disabled
	}

	b.widget.Render(screen, def)

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
				b.text.Color = b.textColor.Disabled
			} else {
				b.text.Color = b.textColor.Idle
			}
		}
	}

	if b.container != nil {
		b.container.Render(screen, def)
	}
}

func (b *Button) draw(screen *ebiten.Image) {
	i := b.Image.Idle
	if b.widget.Disabled {
		if b.Image.Disabled != nil {
			i = b.Image.Disabled
		}
	} else if b.pressing && (b.hovering || b.KeepPressedOnExit) {
		if b.Image.Pressed != nil {
			i = b.Image.Pressed
		}
	} else if b.hovering {
		if b.Image.Hover != nil {
			i = b.Image.Hover
		}
	}

	if i != nil {
		i.Draw(screen, b.widget.Rect.Dx(), b.widget.Rect.Dy(), func(opts *ebiten.DrawImageOptions) {
			b.widget.drawImageOptions(opts)
			b.drawImageOptions(opts)
		})
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

func (b *Button) createWidget() {
	b.widget = NewWidget(
		append(b.widgetOpts, []WidgetOpt{
			WidgetOpts.WithCursorEnterHandler(func(args *WidgetCursorEnterEventArgs) {
				if !b.widget.Disabled {
					b.hovering = true
				}
			}),

			WidgetOpts.WithCursorExitHandler(func(args *WidgetCursorExitEventArgs) {
				b.hovering = false
			}),

			WidgetOpts.WithMouseButtonPressedHandler(func(args *WidgetMouseButtonPressedEventArgs) {
				if !b.widget.Disabled {
					b.pressing = true

					b.PressedEvent.Fire(&ButtonPressedEventArgs{
						Button:  b,
						OffsetX: args.OffsetX,
						OffsetY: args.OffsetY,
					})
				}
			}),

			WidgetOpts.WithMouseButtonReleasedHandler(func(args *WidgetMouseButtonReleasedEventArgs) {
				b.pressing = false

				if !b.widget.Disabled {
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
					}
				}
			}),
		}...)...)
	b.widgetOpts = nil
}
