package widget

import (
	img "image"
	"image/color"
	"math"

	"github.com/ebitenui/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TextAreaParams struct {
	Face                   *text.Face
	ForegroundColor        *color.Color
	TextPadding            *Insets
	ControlWidgetSpacing   *int
	StripBBCode            *bool
	LinkColor              *TextLinkColor
	Slider                 *SliderParams
	ScrollContainerImage   *ScrollContainerImage
	ScrollContainerPadding *Insets
}

type TextArea struct {
	definedParams  TextAreaParams
	computedParams TextAreaParams

	containerOpts []ContainerOpt

	processBBCode        bool
	initialText          string
	verticalScrollMode   ScrollMode
	horizontalScrollMode ScrollMode
	showHorizontalSlider bool
	showVerticalSlider   bool

	linkClickedFunc       LinkHandlerFunc
	linkCursorEnteredFunc LinkHandlerFunc
	linkCursorExitedFunc  LinkHandlerFunc

	init            *MultiOnce
	container       *Container
	layout          *GridLayout
	scrollContainer *ScrollContainer
	vSlider         *Slider
	hSlider         *Slider
	text            *Text
}

type ScrollMode int

const (
	// Default. Scrolling is not automatically handled.
	None ScrollMode = iota

	// The TextArea is automatically scrolled to the beginning on change.
	ScrollBeginning

	// The TextArea is automatically scrolled to the end on change.
	ScrollEnd

	// The TextArea will initially position the text at the end of the scroll area.
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

func (t *TextArea) Validate() {
	t.init.Do()
	t.populateComputedParams()
	if t.computedParams.ForegroundColor == nil {
		panic("TextArea: FontColor is required.")
	}
	if t.computedParams.Face == nil {
		panic("TextArea: FontFace is required.")
	}
	t.initWidget()
}

func (t *TextArea) populateComputedParams() {
	params := TextAreaParams{}

	theme := t.GetWidget().GetTheme()

	// Set theme values
	if theme != nil {
		if theme.TextAreaTheme != nil {
			params.ControlWidgetSpacing = theme.TextAreaTheme.ControlWidgetSpacing
			if theme.TextAreaTheme.Face != nil {
				params.Face = theme.TextAreaTheme.Face
			} else {
				params.Face = theme.DefaultFace
			}
			if theme.TextAreaTheme.ForegroundColor != nil {
				params.ForegroundColor = theme.TextAreaTheme.ForegroundColor
			} else {
				params.ForegroundColor = &theme.DefaultTextColor
			}
			params.LinkColor = theme.TextAreaTheme.LinkColor
			params.ScrollContainerImage = theme.TextAreaTheme.ScrollContainerImage
			params.ScrollContainerPadding = theme.TextAreaTheme.ScrollContainerPadding
			params.Slider = theme.TextAreaTheme.Slider
			params.StripBBCode = theme.TextAreaTheme.StripBBCode
			params.TextPadding = theme.TextAreaTheme.TextPadding
		}
	}

	// Set definedParam values
	if t.definedParams.ControlWidgetSpacing != nil {
		params.ControlWidgetSpacing = t.definedParams.ControlWidgetSpacing
	}
	if t.definedParams.Face != nil {
		params.Face = t.definedParams.Face
	}
	if t.definedParams.ForegroundColor != nil {
		params.ForegroundColor = t.definedParams.ForegroundColor
	}
	if t.definedParams.LinkColor != nil {
		params.LinkColor = t.definedParams.LinkColor
	}
	if t.definedParams.ScrollContainerImage != nil {
		params.ScrollContainerImage = t.definedParams.ScrollContainerImage
	}
	if t.definedParams.ScrollContainerPadding != nil {
		params.ScrollContainerPadding = t.definedParams.ScrollContainerPadding
	}
	if t.definedParams.Slider != nil {
		if params.Slider == nil {
			params.Slider = &SliderParams{}
		}
		if t.definedParams.Slider.FixedHandleSize != nil {
			params.Slider.FixedHandleSize = t.definedParams.Slider.FixedHandleSize
		}
		if t.definedParams.Slider.HandleImage != nil {
			params.Slider.HandleImage = t.definedParams.Slider.HandleImage
		}
		if t.definedParams.Slider.MinHandleSize != nil {
			params.Slider.MinHandleSize = t.definedParams.Slider.MinHandleSize
		}
		if t.definedParams.Slider.TrackImage != nil {
			params.Slider.TrackImage = t.definedParams.Slider.TrackImage
		}
		if t.definedParams.Slider.TrackOffset != nil {
			params.Slider.TrackOffset = t.definedParams.Slider.TrackOffset
		}
		if t.definedParams.Slider.TrackPadding != nil {
			params.Slider.TrackPadding = t.definedParams.Slider.TrackPadding
		}
	}
	if t.definedParams.StripBBCode != nil {
		params.StripBBCode = t.definedParams.StripBBCode
	}
	if t.definedParams.TextPadding != nil {
		params.TextPadding = t.definedParams.TextPadding
	}

	// Set defaults

	if params.TextPadding == nil {
		params.TextPadding = &Insets{}
	}
	if params.ControlWidgetSpacing == nil {
		spacing := 0
		params.ControlWidgetSpacing = &spacing
	}
	if params.StripBBCode == nil {
		FALSE := false
		params.StripBBCode = &FALSE
	}
	if params.ScrollContainerPadding == nil {
		params.ScrollContainerPadding = &Insets{}
	}

	t.computedParams = params
}

// Specify the Container options for the text area.
func (o TextAreaOptions) ContainerOpts(opts ...ContainerOpt) TextAreaOpt {
	return func(l *TextArea) {
		l.containerOpts = append(l.containerOpts, opts...)
	}
}

// Specify the options for the scroll container.
func (o TextAreaOptions) ScrollContainerImage(image *ScrollContainerImage) TextAreaOpt {
	return func(l *TextArea) {
		l.definedParams.ScrollContainerImage = image
	}
}

// Specify the options for the scroll bars.
func (o TextAreaOptions) SliderParams(sliderParams *SliderParams) TextAreaOpt {
	return func(l *TextArea) {
		l.definedParams.Slider = sliderParams
	}
}

// Specify spacing between the text container and scrollbars.
func (o TextAreaOptions) ControlWidgetSpacing(s int) TextAreaOpt {
	return func(l *TextArea) {
		l.definedParams.ControlWidgetSpacing = &s
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

// Set how vertical scrolling should be handled.
func (o TextAreaOptions) VerticalScrollMode(scrollMode ScrollMode) TextAreaOpt {
	return func(l *TextArea) {
		l.verticalScrollMode = scrollMode
	}
}

// Set how horizontal scrolling should be handled.
func (o TextAreaOptions) HorizontalScrollMode(scrollMode ScrollMode) TextAreaOpt {
	return func(l *TextArea) {
		l.horizontalScrollMode = scrollMode
	}
}

// Set the font face for this text area.
func (o TextAreaOptions) FontFace(f *text.Face) TextAreaOpt {
	return func(l *TextArea) {
		l.definedParams.Face = f
	}
}

// Set the default color for the text area.
func (o TextAreaOptions) FontColor(color color.Color) TextAreaOpt {
	return func(l *TextArea) {
		l.definedParams.ForegroundColor = &color
	}
}

// Set how far from the edges of the textarea the text should be set.
func (o TextAreaOptions) TextPadding(i Insets) TextAreaOpt {
	return func(l *TextArea) {
		l.definedParams.TextPadding = &i
	}
}

// Set the initial Text for the text area.
func (o TextAreaOptions) Text(initialText string) TextAreaOpt {
	return func(l *TextArea) {
		l.initialText = initialText
	}
}

// This option tells the textarea object to process BBCodes.
//
// Currently the system supports the following BBCodes:
//   - color - [color=#FFFFFF] text [/color] - defines a color code for the enclosed text
//   - link - [link=id arg1:value1 ... argX:valueX] text [/link] - defines a clickable section of text,
//     that will trigger a callback.
func (o TextAreaOptions) ProcessBBCode(processBBCode bool) TextAreaOpt {
	return func(l *TextArea) {
		l.processBBCode = processBBCode
	}
}

// Set whether or not the text area should automatically strip out BBCodes from being displayed.
func (o TextAreaOptions) StripBBCode(stripBBCode bool) TextAreaOpt {
	return func(l *TextArea) {
		l.definedParams.StripBBCode = &stripBBCode
	}
}

// This option sets the idle and hover color for text that is wrapped in a
// [link][/link] bbcode.
//
// Note: this is only used if ProcessBBCode is true.
func (o TextAreaOptions) LinkColor(linkColor *TextLinkColor) TextAreaOpt {
	return func(t *TextArea) {
		t.definedParams.LinkColor = linkColor
	}
}

// Defines the handler to be called when a BBCode defined link is clicked.
//
// Note: this is only used if ProcessBBCode is true.
func (o TextAreaOptions) LinkClickedEvent(linkClickedFunc LinkHandlerFunc) TextAreaOpt {
	return func(l *TextArea) {
		l.linkClickedFunc = linkClickedFunc
	}
}

// Defines the handler to be called when the cursor enters a BBCode defined link.
//
// Note: this is only used if ProcessBBCode is true.
func (o TextAreaOptions) LinkCursorEnteredEvent(linkCursorEnteredFunc LinkHandlerFunc) TextAreaOpt {
	return func(l *TextArea) {
		l.linkCursorEnteredFunc = linkCursorEnteredFunc
	}
}

// Defines the handler to be called when the cursor enters a BBCode defined link.
//
// Note: this is only used if ProcessBBCode is true.
func (o TextAreaOptions) LinkCursorExitedEvent(linkCursorExitedFunc LinkHandlerFunc) TextAreaOpt {
	return func(l *TextArea) {
		l.linkCursorExitedFunc = linkCursorExitedFunc
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

func (l *TextArea) Render(screen *ebiten.Image) {
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
	l.container.Render(screen)
}

func (l *TextArea) Update() {
	l.init.Do()
	if l.container != nil {
		l.container.Update()
	}
}

func (l *TextArea) createWidget() {
	var cols int
	if l.showVerticalSlider {
		cols = 2
	} else {
		cols = 1
	}
	l.layout = NewGridLayout(
		GridLayoutOpts.Columns(cols),
		GridLayoutOpts.Stretch([]bool{true, false}, []bool{true, false}))

	l.container = NewContainer(
		append([]ContainerOpt{
			ContainerOpts.WidgetOpts(WidgetOpts.TrackHover(true)),
			ContainerOpts.Layout(l.layout),
		}, l.containerOpts...,
		)...)

	l.text = NewText(TextOpts.TextLabel(l.initialText))
}

func (l *TextArea) initWidget() {
	l.layout.columnSpacing = *l.computedParams.ControlWidgetSpacing
	l.layout.rowSpacing = *l.computedParams.ControlWidgetSpacing

	content := NewContainer(
		ContainerOpts.Layout(NewRowLayout(
			RowLayoutOpts.Direction(DirectionVertical))),
		ContainerOpts.AutoDisableChildren())

	l.text = NewText(
		TextOpts.Text(l.initialText, l.computedParams.Face, *l.computedParams.ForegroundColor),
		TextOpts.Padding(l.computedParams.TextPadding),
		TextOpts.ProcessBBCode(l.processBBCode),
		TextOpts.StripBBCode(*l.computedParams.StripBBCode),
		TextOpts.LinkColor(l.computedParams.LinkColor),
		TextOpts.LinkClickedHandler(l.linkClickedFunc),
		TextOpts.LinkCursorEnteredHandler(l.linkCursorEnteredFunc),
		TextOpts.LinkCursorExitedHandler(l.linkCursorExitedFunc),
	)
	content.AddChild(l.text)
	l.text.widget.parent = l.container.GetWidget()

	l.scrollContainer = NewScrollContainer(
		ScrollContainerOpts.Content(content),
		ScrollContainerOpts.StretchContentWidth(),
		ScrollContainerOpts.Image(l.computedParams.ScrollContainerImage),
		ScrollContainerOpts.Padding(*l.computedParams.ScrollContainerPadding),
	)
	l.container.AddChild(l.scrollContainer)

	var sliderOpts []SliderOpt
	if l.computedParams.Slider != nil {
		if l.computedParams.Slider.FixedHandleSize != nil {
			sliderOpts = append(sliderOpts, SliderOpts.FixedHandleSize(*l.computedParams.Slider.FixedHandleSize))
		}
		if l.computedParams.Slider.HandleImage != nil {
			sliderOpts = append(sliderOpts, SliderOpts.HandleImage(l.computedParams.Slider.HandleImage))
		}
		if l.computedParams.Slider.TrackImage != nil {
			sliderOpts = append(sliderOpts, SliderOpts.TrackImage(l.computedParams.Slider.TrackImage))
		}
		if l.computedParams.Slider.MinHandleSize != nil {
			sliderOpts = append(sliderOpts, SliderOpts.MinHandleSize(*l.computedParams.Slider.MinHandleSize))
		}
		if l.computedParams.Slider.TrackOffset != nil {
			sliderOpts = append(sliderOpts, SliderOpts.TrackOffset(*l.computedParams.Slider.TrackOffset))
		}
		if l.computedParams.Slider.TrackPadding != nil {
			sliderOpts = append(sliderOpts, SliderOpts.TrackPadding(l.computedParams.Slider.TrackPadding))
		}
	}

	if l.showVerticalSlider {
		pageSizeFunc := func() int {
			return int(math.Round(float64(l.scrollContainer.ViewRect().Dy()) / float64(content.GetWidget().Rect.Dy()) * 1000))
		}

		l.vSlider = NewSlider(append(sliderOpts,
			SliderOpts.Orientation(DirectionVertical),
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
		)...)

		if l.verticalScrollMode == ScrollEnd || l.verticalScrollMode == PositionAtEnd {
			l.vSlider.Current = l.vSlider.Max
		}
		l.container.AddChild(l.vSlider)

		l.scrollContainer.widget.ScrolledEvent.AddHandler(func(args interface{}) {
			if a, ok := args.(*WidgetScrolledEventArgs); ok {
				p := pageSizeFunc() / 3
				if p < 1 {
					p = 1
				}
				l.vSlider.Current -= int(math.Round(a.Y * float64(p)))
			}
		})
	}

	if l.showHorizontalSlider {
		pageSizeFunc := func() int {
			return int(math.Round(float64(l.scrollContainer.ViewRect().Dx()) / float64(content.GetWidget().Rect.Dx()) * 1000))
		}

		l.hSlider = NewSlider(append(sliderOpts,
			SliderOpts.Orientation(DirectionHorizontal),
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
		)...)

		if l.horizontalScrollMode == ScrollEnd || l.horizontalScrollMode == PositionAtEnd {
			l.hSlider.Current = l.hSlider.Max
		}
		l.container.AddChild(l.hSlider)
	}

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
	l.text.Label += value

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
