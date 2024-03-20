package widget

import (
	img "image"
	"image/color"
	"math"

	"github.com/ebitenui/ebitenui/input"

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
	showHorizontalSlider bool
	showVerticalSlider   bool
	processBBCode        bool
	initialText          string
	verticalScrollMode   ScrollMode
	horizontalScrollMode ScrollMode

	init            *MultiOnce
	container       *Container
	scrollContainer *ScrollContainer
	vSlider         *Slider
	hSlider         *Slider
	text            *Text
}

type ScrollMode int

const (
	// Default. Scrolling is not automatically handled
	None ScrollMode = iota

	// The TextArea is automatically scrolled to the beginning on change
	ScrollBeginning

	// The TextArea is automatically scrolled to the end on change
	ScrollEnd

	// The TextArea will initially position the text at the end of the scroll area
	PositionAtEnd
)

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

// Specify the Container options for the text area
func (o TextAreaOptions) ContainerOpts(opts ...ContainerOpt) TextAreaOpt {
	return func(l *TextArea) {
		l.containerOpts = append(l.containerOpts, opts...)
	}
}

// Specify the options for the scroll container
func (o TextAreaOptions) ScrollContainerOpts(opts ...ScrollContainerOpt) TextAreaOpt {
	return func(l *TextArea) {
		l.scrollContainerOpts = append(l.scrollContainerOpts, opts...)
	}
}

// Specify the options for the scroll bars
func (o TextAreaOptions) SliderOpts(opts ...SliderOpt) TextAreaOpt {
	return func(l *TextArea) {
		l.sliderOpts = append(l.sliderOpts, opts...)
	}
}

// Specify spacing between the text container and scrollbars
func (o TextAreaOptions) ControlWidgetSpacing(s int) TextAreaOpt {
	return func(l *TextArea) {
		l.controlWidgetSpacing = s
	}
}

// Show the horizontal scrollbar.
func (o TextAreaOptions) ShowHorizontalScrollbar() TextAreaOpt {
	return func(l *TextArea) {
		l.showHorizontalSlider = true
	}
}

// Show the vertical scrollbar.
func (o TextAreaOptions) ShowVerticalScrollbar() TextAreaOpt {
	return func(l *TextArea) {
		l.showVerticalSlider = true
	}
}

// Set how vertical scrolling should be handled
func (o TextAreaOptions) VerticalScrollMode(scrollMode ScrollMode) TextAreaOpt {
	return func(l *TextArea) {
		l.verticalScrollMode = scrollMode
	}
}

// Set how horizontal scrolling should be handled
func (o TextAreaOptions) HorizontalScrollMode(scrollMode ScrollMode) TextAreaOpt {
	return func(l *TextArea) {
		l.horizontalScrollMode = scrollMode
	}
}

// Set the font face for this text area
func (o TextAreaOptions) FontFace(f font.Face) TextAreaOpt {
	return func(l *TextArea) {
		l.face = f
	}
}

// Set the default color for the text area
func (o TextAreaOptions) FontColor(color color.Color) TextAreaOpt {
	return func(l *TextArea) {
		l.foregroundColor = color
	}
}

// Set how far from the edges of the textarea the text should be set
func (o TextAreaOptions) TextPadding(i Insets) TextAreaOpt {
	return func(l *TextArea) {
		l.textPadding = i
	}
}

// Set the initial Text for the text area
func (o TextAreaOptions) Text(initialText string) TextAreaOpt {
	return func(l *TextArea) {
		l.initialText = initialText
	}
}

// Set whether or not the text area should process BBCodes. e.g. [color=FF0000]red[/color]
func (o TextAreaOptions) ProcessBBCode(processBBCode bool) TextAreaOpt {
	return func(l *TextArea) {
		l.processBBCode = processBBCode
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

func (l *TextArea) GetFocusers() []Focuser {
	l.init.Do()
	var result []Focuser
	if l.hSlider != nil && l.hSlider.tabOrder != -1 {
		result = append(result, l.hSlider)
	}
	if l.vSlider != nil && l.vSlider.tabOrder != -1 {
		result = append(result, l.vSlider)
	}
	return result
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
	if l.showVerticalSlider {
		cols = 2
	} else {
		cols = 1
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
		TextOpts.ProcessBBCode(l.processBBCode),
	)
	content.AddChild(l.text)
	l.text.widget.parent = l.container.GetWidget()

	l.scrollContainer = NewScrollContainer(append(l.scrollContainerOpts, []ScrollContainerOpt{
		ScrollContainerOpts.Content(content),
		ScrollContainerOpts.StretchContentWidth(),
	}...)...)
	l.scrollContainerOpts = nil
	l.container.AddChild(l.scrollContainer)

	if l.showVerticalSlider {
		pageSizeFunc := func() int {
			return int(math.Round(float64(l.scrollContainer.ViewRect().Dy()) / float64(content.GetWidget().Rect.Dy()) * 1000))
		}

		l.vSlider = NewSlider(append(l.sliderOpts, []SliderOpt{
			SliderOpts.Direction(DirectionVertical),
			SliderOpts.MinMax(0, 1000),
			SliderOpts.PageSizeFunc(pageSizeFunc),
			SliderOpts.ChangedHandler(func(args *SliderChangedEventArgs) {
				current := args.Slider.Current
				if pageSizeFunc() >= 1000 {
					current = 0
					if l.verticalScrollMode == ScrollEnd || l.verticalScrollMode == PositionAtEnd {
						current = 1000
					}
				}
				l.scrollContainer.ScrollTop = float64(current) / 1000
			}),
		}...)...)
		if l.verticalScrollMode == ScrollEnd || l.verticalScrollMode == PositionAtEnd {
			l.vSlider.Current = l.vSlider.Max
		}
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

	if l.showHorizontalSlider {
		pageSizeFunc := func() int {
			return int(math.Round(float64(l.scrollContainer.ViewRect().Dx()) / float64(content.GetWidget().Rect.Dx()) * 1000))
		}

		l.hSlider = NewSlider(append(l.sliderOpts, []SliderOpt{
			SliderOpts.Direction(DirectionHorizontal),
			SliderOpts.MinMax(0, 1000),
			SliderOpts.PageSizeFunc(pageSizeFunc),
			SliderOpts.ChangedHandler(func(args *SliderChangedEventArgs) {
				current := args.Slider.Current
				if pageSizeFunc() >= 1000 {
					current = 0
					if l.horizontalScrollMode == ScrollEnd || l.horizontalScrollMode == PositionAtEnd {
						current = 1000
					}
				}
				l.scrollContainer.ScrollLeft = float64(current) / 1000
			}),
		}...)...)
		if l.horizontalScrollMode == ScrollEnd || l.horizontalScrollMode == PositionAtEnd {
			l.hSlider.Current = l.hSlider.Max
		}
		l.container.AddChild(l.hSlider)
	}

	l.sliderOpts = nil
}

func (l *TextArea) PrependText(value string) {
	l.text.Label = value + l.text.Label

	if l.showHorizontalSlider {
		if l.horizontalScrollMode == ScrollBeginning {
			l.hSlider.Current = 0
		} else if l.horizontalScrollMode == ScrollEnd {
			l.hSlider.Current = l.hSlider.Max
		}
	}
	if l.showVerticalSlider {
		if l.verticalScrollMode == ScrollBeginning {
			l.vSlider.Current = 0
		} else if l.verticalScrollMode == ScrollEnd {
			l.vSlider.Current = l.vSlider.Max
		}
	}
}

func (l *TextArea) AppendText(value string) {
	l.init.Do()
	l.text.Label = l.text.Label + value

	if l.showHorizontalSlider {
		if l.horizontalScrollMode == ScrollBeginning {
			l.hSlider.Current = 0
		} else if l.horizontalScrollMode == ScrollEnd {
			l.hSlider.Current = l.hSlider.Max
		}
	}
	if l.showVerticalSlider {
		if l.verticalScrollMode == ScrollBeginning {
			l.vSlider.Current = 0
		} else if l.verticalScrollMode == ScrollEnd {
			l.vSlider.Current = l.vSlider.Max
		}
	}
}

func (l *TextArea) SetText(value string) {
	l.init.Do()
	l.text.Label = value

	if l.showHorizontalSlider {
		if l.horizontalScrollMode == ScrollBeginning {
			l.hSlider.Current = 0
		} else if l.horizontalScrollMode == ScrollEnd {
			l.hSlider.Current = l.hSlider.Max
		}
	}
	if l.showVerticalSlider {
		if l.verticalScrollMode == ScrollBeginning {
			l.vSlider.Current = 0
		} else if l.verticalScrollMode == ScrollEnd {
			l.vSlider.Current = l.vSlider.Max
		}
	}
}

func (l *TextArea) GetText() string {
	l.init.Do()
	return l.text.Label
}
