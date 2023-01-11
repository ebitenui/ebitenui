package widget

import (
	img "image"
	"image/color"
	"math"

	"github.com/mcarpenter622/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type TextArea struct {
	containerOpts        []ContainerOpt
	scrollContainerOpts  []ScrollContainerOpt
	sliderOpts           []SliderOpt
	face                 font.Face
	foregroundColor      color.Color
	textPadding          Insets
	controlWidgetSpacing int
	hideHorizontalSlider bool
	hideVerticalSlider   bool
	allowReselect        bool
	initialText          string

	init            *MultiOnce
	container       *Container
	scrollContainer *ScrollContainer
	vSlider         *Slider
	hSlider         *Slider
	text            *Text
}

type TextAreaOpt func(l *TextArea)

type TextAreaEntrySelectedEventArgs struct {
	TextArea      *TextArea
	Entry         interface{}
	PreviousEntry interface{}
}

type TextAreaEntrySelectedHandlerFunc func(args *TextAreaEntrySelectedEventArgs)

type TextAreaOptions struct {
}

var TextAreaOpts TextAreaOptions

func NewTextArea(opts ...TextAreaOpt) *TextArea {
	l := &TextArea{
		init: &MultiOnce{},
	}

	l.init.Append(l.createWidget)

	for _, o := range opts {
		o(l)
	}

	return l
}

func (o TextAreaOptions) ContainerOpts(opts ...ContainerOpt) TextAreaOpt {
	return func(l *TextArea) {
		l.containerOpts = append(l.containerOpts, opts...)
	}
}

func (o TextAreaOptions) ScrollContainerOpts(opts ...ScrollContainerOpt) TextAreaOpt {
	return func(l *TextArea) {
		l.scrollContainerOpts = append(l.scrollContainerOpts, opts...)
	}
}

func (o TextAreaOptions) SliderOpts(opts ...SliderOpt) TextAreaOpt {
	return func(l *TextArea) {
		l.sliderOpts = append(l.sliderOpts, opts...)
	}
}

func (o TextAreaOptions) ControlWidgetSpacing(s int) TextAreaOpt {
	return func(l *TextArea) {
		l.controlWidgetSpacing = s
	}
}

func (o TextAreaOptions) HideHorizontalSlider() TextAreaOpt {
	return func(l *TextArea) {
		l.hideHorizontalSlider = true
	}
}

func (o TextAreaOptions) HideVerticalSlider() TextAreaOpt {
	return func(l *TextArea) {
		l.hideVerticalSlider = true
	}
}

func (o TextAreaOptions) FontFace(f font.Face) TextAreaOpt {
	return func(l *TextArea) {
		l.face = f
	}
}

func (o TextAreaOptions) FontColor(color color.Color) TextAreaOpt {
	return func(l *TextArea) {
		l.foregroundColor = color
	}
}

func (o TextAreaOptions) TextPadding(i Insets) TextAreaOpt {
	return func(l *TextArea) {
		l.textPadding = i
	}
}

func (o TextAreaOptions) AllowReselect() TextAreaOpt {
	return func(l *TextArea) {
		l.allowReselect = true
	}
}

func (o TextAreaOptions) Text(initialText string) TextAreaOpt {
	return func(l *TextArea) {
		l.initialText = initialText
	}
}

func (l *TextArea) GetWidget() *Widget {
	l.init.Do()
	return l.container.GetWidget()
}

func (l *TextArea) PreferredSize() (int, int) {
	l.init.Do()
	w, h := l.container.PreferredSize()

	if l.container.widget != nil && h < l.container.widget.MinHeight {
		h = l.container.widget.MinHeight
	}
	if l.container.widget != nil && w < l.container.widget.MinWidth {
		w = l.container.widget.MinWidth
	}
	return w, h
}

func (l *TextArea) SetLocation(rect img.Rectangle) {
	l.init.Do()
	l.container.GetWidget().Rect = rect
}

func (l *TextArea) RequestRelayout() {
	l.init.Do()
	l.container.RequestRelayout()
}

func (l *TextArea) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	l.init.Do()
	l.container.SetupInputLayer(def)
}

func (l *TextArea) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	l.init.Do()

	d := l.container.GetWidget().Disabled

	if l.vSlider != nil {
		l.vSlider.DrawTrackDisabled = d
	}
	if l.hSlider != nil {
		l.hSlider.DrawTrackDisabled = d
	}
	l.text.MaxWidth = float64(l.container.GetWidget().Rect.Dx())
	l.scrollContainer.GetWidget().Disabled = d
	l.container.Render(screen, def)
}

func (l *TextArea) createWidget() {
	var cols int
	if l.hideVerticalSlider {
		cols = 1
	} else {
		cols = 2
	}

	l.container = NewContainer(
		append(l.containerOpts,
			ContainerOpts.Layout(NewGridLayout(
				GridLayoutOpts.Columns(cols),
				GridLayoutOpts.Stretch([]bool{true, false}, []bool{true, false}),
				GridLayoutOpts.Spacing(l.controlWidgetSpacing, l.controlWidgetSpacing))))...)
	l.containerOpts = nil

	content := NewContainer(
		ContainerOpts.Layout(NewRowLayout(
			RowLayoutOpts.Direction(DirectionVertical))),
		ContainerOpts.AutoDisableChildren())

	l.text = NewText(
		TextOpts.Text(l.initialText, l.face, l.foregroundColor),
		TextOpts.Insets(l.textPadding),
	)
	content.AddChild(l.text)

	l.scrollContainer = NewScrollContainer(append(l.scrollContainerOpts, []ScrollContainerOpt{
		ScrollContainerOpts.Content(content),
		ScrollContainerOpts.StretchContentWidth(),
	}...)...)
	l.scrollContainerOpts = nil
	l.container.AddChild(l.scrollContainer)

	if !l.hideVerticalSlider {
		pageSizeFunc := func() int {
			return int(math.Round(float64(l.scrollContainer.ContentRect().Dy()) / float64(content.GetWidget().Rect.Dy()) * 1000))
		}

		l.vSlider = NewSlider(append(l.sliderOpts, []SliderOpt{
			SliderOpts.Direction(DirectionVertical),
			SliderOpts.MinMax(0, 1000),
			SliderOpts.PageSizeFunc(pageSizeFunc),
			SliderOpts.ChangedHandler(func(args *SliderChangedEventArgs) {
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
		l.hSlider = NewSlider(append(l.sliderOpts, []SliderOpt{
			SliderOpts.Direction(DirectionHorizontal),
			SliderOpts.MinMax(0, 1000),
			SliderOpts.PageSizeFunc(func() int {
				return int(math.Round(float64(l.scrollContainer.ContentRect().Dx()) / float64(content.GetWidget().Rect.Dx()) * 1000))
			}),
			SliderOpts.ChangedHandler(func(args *SliderChangedEventArgs) {
				l.scrollContainer.ScrollLeft = float64(args.Slider.Current) / 1000
			}),
		}...)...)
		l.container.AddChild(l.hSlider)
	}

	l.sliderOpts = nil
}

func (l *TextArea) SetScrollTop(t float64) {
	l.init.Do()
	if l.vSlider != nil {
		l.vSlider.Current = int(math.Round(t * 1000))
	}
	l.scrollContainer.ScrollTop = t
}

func (l *TextArea) SetScrollLeft(left float64) {
	l.init.Do()
	if l.hSlider != nil {
		l.hSlider.Current = int(math.Round(left * 1000))
	}
	l.scrollContainer.ScrollLeft = left
}

func (l *TextArea) PrependText(value string) {
	l.text.Label = value + l.text.Label
}

func (l *TextArea) AppendText(value string) {
	l.text.Label = l.text.Label + value
}

func (l *TextArea) SetText(value string) {
	l.text.Label = value
}
func (l *TextArea) GetText() string {
	return l.text.Label
}
