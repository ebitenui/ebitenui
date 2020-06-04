package widget

import (
	"image"
	"image/color"
	"math"

	"github.com/blizzy78/ebitenui/event"
	ebimage "github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/input"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

type List struct {
	EntrySelectedEvent *event.Event

	containerOpts        []ContainerOpt
	scrollContainerOpts  []ScrollContainerOpt
	sliderOpts           []SliderOpt
	entries              []interface{}
	entryLabelFunc       ListEntryLabelFunc
	entryFace            font.Face
	entryUnselectedColor *ButtonImage
	entrySelectedColor   *ButtonImage
	entryTextColor       *ButtonTextColor
	controlWidgetSpacing int
	hideHorizontalSlider bool
	hideVerticalSlider   bool
	allowReselect        bool

	init            *MultiOnce
	container       *Container
	scrollContainer *ScrollContainer
	vSlider         *Slider
	hSlider         *Slider
	buttons         []*Button
	selectedEntry   interface{}
}

type ListOpt func(l *List)

type ListEntryLabelFunc func(e interface{}) string

type ListEntryColor struct {
	Unselected                 color.Color
	Selected                   color.Color
	DisabledUnselected         color.Color
	DisabledSelected           color.Color
	SelectedBackground         color.Color
	DisabledSelectedBackground color.Color
}

type ListEntrySelectedEventArgs struct {
	List          *List
	Entry         interface{}
	PreviousEntry interface{}
}

type ListEntrySelectedHandlerFunc func(args *ListEntrySelectedEventArgs)

const ListOpts = listOpts(true)

type listOpts bool

func NewList(opts ...ListOpt) *List {
	l := &List{
		EntrySelectedEvent: &event.Event{},

		init: &MultiOnce{},
	}

	l.init.Append(l.createWidget)

	for _, o := range opts {
		o(l)
	}

	return l
}

func (o listOpts) WithContainerOpts(opts ...ContainerOpt) ListOpt {
	return func(l *List) {
		l.containerOpts = append(l.containerOpts, opts...)
	}
}

func (o listOpts) WithScrollContainerOpts(opts ...ScrollContainerOpt) ListOpt {
	return func(l *List) {
		l.scrollContainerOpts = append(l.scrollContainerOpts, opts...)
	}
}

func (o listOpts) WithSliderOpts(opts ...SliderOpt) ListOpt {
	return func(l *List) {
		l.sliderOpts = append(l.sliderOpts, opts...)
	}
}

func (o listOpts) WithControlWidgetSpacing(s int) ListOpt {
	return func(l *List) {
		l.controlWidgetSpacing = s
	}
}

func (o listOpts) WithHideHorizontalSlider() ListOpt {
	return func(l *List) {
		l.hideHorizontalSlider = true
	}
}

func (o listOpts) WithHideVerticalSlider() ListOpt {
	return func(l *List) {
		l.hideVerticalSlider = true
	}
}

func (o listOpts) WithEntries(e []interface{}) ListOpt {
	return func(l *List) {
		l.entries = e
	}
}

func (o listOpts) WithEntryLabelFunc(f ListEntryLabelFunc) ListOpt {
	return func(l *List) {
		l.entryLabelFunc = f
	}
}

func (o listOpts) WithEntryFontFace(f font.Face) ListOpt {
	return func(l *List) {
		l.entryFace = f
	}
}

func (o listOpts) WithEntryColor(c *ListEntryColor) ListOpt {
	return func(l *List) {
		l.entryUnselectedColor = &ButtonImage{
			Idle:     ebimage.NewNineSliceColor(color.Transparent),
			Disabled: ebimage.NewNineSliceColor(color.Transparent),
		}

		l.entrySelectedColor = &ButtonImage{
			Idle:     ebimage.NewNineSliceColor(c.SelectedBackground),
			Disabled: ebimage.NewNineSliceColor(c.DisabledSelectedBackground),
		}

		l.entryTextColor = &ButtonTextColor{
			Idle:     c.Unselected,
			Disabled: c.DisabledUnselected,
		}
	}
}

func (o listOpts) WithEntrySelectedHandler(f ListEntrySelectedHandlerFunc) ListOpt {
	return func(l *List) {
		l.EntrySelectedEvent.AddHandler(func(args interface{}) {
			f(args.(*ListEntrySelectedEventArgs))
		})
	}
}

func (o listOpts) WithAllowReselect() ListOpt {
	return func(l *List) {
		l.allowReselect = true
	}
}

func (l *List) GetWidget() *Widget {
	l.init.Do()
	return l.container.GetWidget()
}

func (l *List) PreferredSize() (int, int) {
	l.init.Do()
	return l.container.PreferredSize()
}

func (l *List) SetLocation(rect image.Rectangle) {
	l.init.Do()
	l.container.GetWidget().Rect = rect
}

func (l *List) RequestRelayout() {
	l.init.Do()
	l.container.RequestRelayout()
}

func (l *List) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	l.init.Do()
	l.container.SetupInputLayer(def)
}

func (l *List) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	l.init.Do()

	l.scrollContainer.GetWidget().Disabled = l.container.GetWidget().Disabled

	l.container.Render(screen, def)
}

func (l *List) createWidget() {
	var cols int
	if l.hideVerticalSlider {
		cols = 1
	} else {
		cols = 2
	}

	l.container = NewContainer(
		append(l.containerOpts,
			ContainerOpts.WithLayout(NewGridLayout(
				GridLayoutOpts.WithColumns(cols),
				GridLayoutOpts.WithStretch([]bool{true, false}, []bool{true, false}),
				GridLayoutOpts.WithSpacing(l.controlWidgetSpacing, l.controlWidgetSpacing))))...)
	l.containerOpts = nil

	content := NewContainer(
		ContainerOpts.WithLayout(NewRowLayout(
			RowLayoutOpts.WithDirection(DirectionVertical))),
		ContainerOpts.WithAutoDisableChildren())

	l.buttons = make([]*Button, 0, len(l.entries))
	for _, e := range l.entries {
		e := e
		but := NewButton(
			ButtonOpts.WithWidgetOpts(WidgetOpts.WithLayoutData(&RowLayoutData{
				Stretch: true,
			})),
			ButtonOpts.WithImage(l.entryUnselectedColor),
			ButtonOpts.WithTextSimpleLeft(l.entryLabelFunc(e), l.entryFace, l.entryTextColor, Insets{
				Left:   6,
				Right:  6,
				Top:    2,
				Bottom: 2,
			}),

			ButtonOpts.WithClickedHandler(func(args *ButtonClickedEventArgs) {
				l.setSelectedEntry(e, true)
			}))

		l.buttons = append(l.buttons, but)

		content.AddChild(but)
	}

	l.scrollContainer = NewScrollContainer(
		append(l.scrollContainerOpts, []ScrollContainerOpt{
			ScrollContainerOpts.WithContent(content),
			ScrollContainerOpts.WithStretchContentWidth(),
		}...)...)
	l.scrollContainerOpts = nil
	l.container.AddChild(l.scrollContainer)

	if !l.hideVerticalSlider {
		pageSizeFunc := func() int {
			return int(math.Round(float64(l.scrollContainer.ContentRect().Dy()) / float64(content.GetWidget().Rect.Dy()) * 1000))
		}

		l.vSlider = NewSlider(
			append(l.sliderOpts, []SliderOpt{
				SliderOpts.WithDirection(DirectionVertical),
				SliderOpts.WithMinMax(0, 1000),
				SliderOpts.WithPageSizeFunc(pageSizeFunc),
				SliderOpts.WithChangedHandler(func(args *SliderChangedEventArgs) {
					l.scrollContainer.ScrollTop = float64(args.Slider.Current) / 1000
				}),
			}...)...)
		l.container.AddChild(l.vSlider)

		l.scrollContainer.widget.ScrolledEvent.AddHandler(func(args interface{}) {
			a := args.(*WidgetScrolledEventArgs)
			p := pageSizeFunc() / 3
			if p < 1 {
				p = 1
			}
			l.vSlider.Current -= int(math.Round(a.Y * float64(p)))
		})
	}

	if !l.hideHorizontalSlider {
		l.hSlider = NewSlider(
			append(l.sliderOpts, []SliderOpt{
				SliderOpts.WithDirection(DirectionHorizontal),
				SliderOpts.WithMinMax(0, 1000),
				SliderOpts.WithPageSizeFunc(func() int {
					return int(math.Round(float64(l.scrollContainer.ContentRect().Dx()) / float64(content.GetWidget().Rect.Dx()) * 1000))
				}),
				SliderOpts.WithChangedHandler(func(args *SliderChangedEventArgs) {
					l.scrollContainer.ScrollLeft = float64(args.Slider.Current) / 1000
				}),
			}...)...,
		)
		l.container.AddChild(l.hSlider)
	}

	l.sliderOpts = nil
}

func (l *List) SetSelectedEntry(e interface{}) {
	l.setSelectedEntry(e, false)
}

func (l *List) setSelectedEntry(e interface{}, user bool) {
	if e != l.selectedEntry || (user && l.allowReselect) {
		l.init.Do()

		prev := l.selectedEntry
		l.selectedEntry = e

		for i, b := range l.buttons {
			if l.entries[i] == e {
				b.Image = l.entrySelectedColor
			} else {
				b.Image = l.entryUnselectedColor
			}
		}

		l.EntrySelectedEvent.Fire(&ListEntrySelectedEventArgs{
			Entry:         e,
			PreviousEntry: prev,
		})
	}
}

func (l *List) SelectedEntry() interface{} {
	l.init.Do()
	return l.selectedEntry
}

func (l *List) SetScrollTop(t float64) {
	l.init.Do()
	if l.vSlider != nil {
		l.vSlider.Current = int(math.Round(t * 1000))
	}
	l.scrollContainer.ScrollTop = t
}

func (li *List) SetScrollLeft(l float64) {
	li.init.Do()
	if li.hSlider != nil {
		li.hSlider.Current = int(math.Round(l * 1000))
	}
	li.scrollContainer.ScrollLeft = l
}
