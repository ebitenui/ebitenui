package widget

import (
	img "image"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type LabelOrder int

const (
	CHECKBOX_FIRST LabelOrder = iota
	LABEL_FIRST
)

type CheckboxParams struct {
	Image *CheckboxImage
	Label *LabelParams
}

type Checkbox struct {
	definedParams  CheckboxParams
	computedParams CheckboxParams

	init              *MultiOnce
	widget            *Widget
	widgetOpts        []WidgetOpt
	triState          bool
	StateChangedEvent *event.Event

	state    WidgetState
	hovering bool

	labelString string
	label       *Label
	spacing     int
	order       LabelOrder

	// Allows the user to disable space bar and enter automatically triggering a focused checkbox.
	DisableDefaultKeys bool

	tabOrder int
	focused  bool
	focusMap map[FocusDirection]Focuser
}

type CheckboxOpt func(c *Checkbox)

type CheckboxImage struct {
	Unchecked         *image.NineSlice
	UncheckedHovered  *image.NineSlice
	UncheckedDisabled *image.NineSlice
	Checked           *image.NineSlice
	CheckedHovered    *image.NineSlice
	CheckedDisabled   *image.NineSlice
	Greyed            *image.NineSlice
	GreyedHovered     *image.NineSlice
	GreyedDisabled    *image.NineSlice
}

type CheckboxChangedEventArgs struct {
	Active *Checkbox
	State  WidgetState
}

type CheckboxChangedHandlerFunc func(args *CheckboxChangedEventArgs)

type CheckboxOptions struct {
}

var CheckboxOpts CheckboxOptions

func NewCheckbox(opts ...CheckboxOpt) *Checkbox {
	c := &Checkbox{
		StateChangedEvent: &event.Event{},
		spacing:           8,
		order:             CHECKBOX_FIRST,
		init:              &MultiOnce{},

		focusMap: make(map[FocusDirection]Focuser),
	}

	for _, o := range opts {
		o(c)
	}

	c.createWidget()

	return c
}

func (c *Checkbox) Validate() {
	c.init.Do()
	c.populateComputedParams()

	if c.computedParams.Image == nil {
		panic("Checkbox: Image is required.")
	}

	if c.computedParams.Image.Checked == nil {
		panic("Checkbox: Image.Checked is required.")
	}

	if c.computedParams.Image.Unchecked == nil {
		panic("Checkbox: Image.Unchecked is required.")
	}

	if c.triState && c.computedParams.Image.Greyed == nil {
		panic("Checkbox: Image.Greyed is required if this is a tristate checkbox.")
	}
	if c.state == WidgetGreyed && !c.triState {
		panic("Checkbox: non-tri state Checkbox cannot be in greyed state.")
	}
	if c.label != nil {
		c.setChildComputedParams()
	}
}

func (c *Checkbox) populateComputedParams() {
	params := CheckboxParams{
		Image: &CheckboxImage{},
		Label: &LabelParams{},
	}

	theme := c.GetWidget().GetTheme()

	if theme != nil {
		if theme.CheckboxTheme != nil {
			params.Image = theme.CheckboxTheme.Image
			params.Label = theme.CheckboxTheme.Label
			if theme.CheckboxTheme.Label != nil {
				params.Label.Face = theme.LabelTheme.Face
				params.Label.Color = theme.LabelTheme.Color
				params.Label.Padding = theme.LabelTheme.Padding
			}
		}
	}

	if c.definedParams.Image != nil {
		if c.definedParams.Image.Checked != nil {
			params.Image.Checked = c.definedParams.Image.Checked
		}
		if c.definedParams.Image.CheckedDisabled != nil {
			params.Image.CheckedDisabled = c.definedParams.Image.CheckedDisabled
		}
		if c.definedParams.Image.CheckedHovered != nil {
			params.Image.CheckedHovered = c.definedParams.Image.CheckedHovered
		}
		if c.definedParams.Image.Greyed != nil {
			params.Image.Greyed = c.definedParams.Image.Greyed
		}
		if c.definedParams.Image.GreyedDisabled != nil {
			params.Image.GreyedDisabled = c.definedParams.Image.GreyedDisabled
		}
		if c.definedParams.Image.GreyedHovered != nil {
			params.Image.GreyedHovered = c.definedParams.Image.GreyedHovered
		}
		if c.definedParams.Image.Unchecked != nil {
			params.Image.Unchecked = c.definedParams.Image.Unchecked
		}
		if c.definedParams.Image.UncheckedDisabled != nil {
			params.Image.UncheckedDisabled = c.definedParams.Image.UncheckedDisabled
		}
		if c.definedParams.Image.UncheckedHovered != nil {
			params.Image.UncheckedHovered = c.definedParams.Image.UncheckedHovered
		}
	}

	if c.definedParams.Label != nil {
		if params.Label.Color == nil {
			params.Label.Color = c.definedParams.Label.Color
		} else if c.definedParams.Label.Color != nil {
			if c.definedParams.Label.Color.Idle != nil {
				params.Label.Color.Idle = c.definedParams.Label.Color.Idle
			}
			if c.definedParams.Label.Color.Disabled != nil {
				params.Label.Color.Disabled = c.definedParams.Label.Color.Disabled
			}
		}
		if c.definedParams.Label.Face != nil {
			params.Label.Face = c.definedParams.Label.Face
		}
		if c.definedParams.Label.Padding != nil {
			params.Label.Padding = c.definedParams.Label.Padding
		}
	}

	c.computedParams = params
}

func (o CheckboxOptions) WidgetOpts(opts ...WidgetOpt) CheckboxOpt {
	return func(c *Checkbox) {
		c.widgetOpts = append(c.widgetOpts, opts...)
	}
}

// This option allows you to specify a label to be shown before or after the checkbox.
func (o CheckboxOptions) Text(labelString string, face *text.Face, color *LabelColor) CheckboxOpt {
	return func(l *Checkbox) {
		if l.definedParams.Label == nil {
			l.definedParams.Label = &LabelParams{}
		}
		l.definedParams.Label.Color = color
		l.definedParams.Label.Face = face
		l.labelString = labelString
	}
}

func (o CheckboxOptions) TextLabel(label string) CheckboxOpt {
	return func(l *Checkbox) {
		if l.definedParams.Label == nil {
			l.definedParams.Label = &LabelParams{}
		}
		l.labelString = label
	}
}

func (o CheckboxOptions) TextFace(face *text.Face) CheckboxOpt {
	return func(l *Checkbox) {
		if l.definedParams.Label == nil {
			l.definedParams.Label = &LabelParams{}
		}
		l.definedParams.Label.Face = face
	}
}

func (o CheckboxOptions) TextColor(color *LabelColor) CheckboxOpt {
	return func(l *Checkbox) {
		if l.definedParams.Label == nil {
			l.definedParams.Label = &LabelParams{}
		}
		l.definedParams.Label.Color = color
	}
}

// This option defines how far the checkbox and label should be spaced horizontally if there is a label.
func (o CheckboxOptions) Spacing(s int) CheckboxOpt {
	return func(l *Checkbox) {
		l.spacing = s
	}
}

// This option indicates that the label should be before the checkbox.
func (o CheckboxOptions) LabelFirst() CheckboxOpt {
	return func(l *Checkbox) {
		l.order = LABEL_FIRST
	}
}

// This option defines the images to show for the checkbox.
// i.Checked and i.Unchecked are required.
func (o CheckboxOptions) Image(i *CheckboxImage) CheckboxOpt {
	return func(c *Checkbox) {
		c.definedParams.Image = i
	}
}

// This option indicates this checkbox should have 3 states.
// If this option is specified a Greyed image is required.
func (o CheckboxOptions) TriState() CheckboxOpt {
	return func(c *Checkbox) {
		c.triState = true
	}
}

func (o CheckboxOptions) TabOrder(tabOrder int) CheckboxOpt {
	return func(c *Checkbox) {
		c.tabOrder = tabOrder
	}
}

// This option allows you to specify a callback to be called when the checkbox state is changed.
func (o CheckboxOptions) StateChangedHandler(f CheckboxChangedHandlerFunc) CheckboxOpt {
	return func(c *Checkbox) {
		c.StateChangedEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*CheckboxChangedEventArgs); ok {
				f(arg)
			}
		})
	}
}

// This option sets the initial state for the checkbox.
func (o CheckboxOptions) InitialState(state WidgetState) CheckboxOpt {
	return func(c *Checkbox) {
		c.state = state
	}
}

// This option will disable enter and space from submitting a focused checkbox.
func (o CheckboxOptions) DisableDefaultKeys() CheckboxOpt {
	return func(c *Checkbox) {
		c.DisableDefaultKeys = true
	}
}

// This function will return the internal Text object if this checkbox has a label, otherwise nil.
func (tw *Checkbox) Text() *Text {
	if tw.label == nil {
		return nil
	}
	return tw.label.text
}

// This function will return the current state the checkbox is in.
func (tw *Checkbox) State() WidgetState {
	return tw.state
}

// This function will allow you to update the checkbox's current state.
func (tw *Checkbox) SetState(state WidgetState) {
	if state == WidgetGreyed && !tw.triState {
		panic("non-tri state Checkbox cannot be in greyed state")
	}

	if state != tw.state {
		tw.state = state

		tw.StateChangedEvent.Fire(&CheckboxChangedEventArgs{
			Active: tw,
			State:  tw.state,
		})
	}
}

// This method is required for this to be part of a radio button group.
func (tw *Checkbox) getStateChangedEvent() *event.Event {
	return tw.StateChangedEvent
}

func (c *Checkbox) GetWidget() *Widget {
	c.init.Do()
	return c.widget
}

func (c *Checkbox) PreferredSize() (int, int) {
	c.init.Do()

	w, h := 0, 0

	iw, ih := c.computedParams.Image.Unchecked.MinSize()
	if w < iw {
		w = iw
	}
	if h < ih {
		h = ih
	}

	if c.label != nil {
		labelX, labelY := c.label.PreferredSize()
		w = w + labelX + c.spacing
		h = max(h, labelY)
	}

	if c.widget != nil && h < c.widget.MinHeight {
		h = c.widget.MinHeight
	}
	if c.widget != nil && w < c.widget.MinWidth {
		w = c.widget.MinWidth
	}
	return w, h
}

func (c *Checkbox) checkboxPreferredSize() (int, int) {
	c.init.Do()

	w, h := 0, 0

	iw, ih := c.computedParams.Image.Unchecked.MinSize()
	if w < iw {
		w = iw
	}
	if h < ih {
		h = ih
	}

	if c.widget != nil && h < c.widget.MinHeight {
		h = c.widget.MinHeight
	}
	if c.widget != nil && w < c.widget.MinWidth {
		w = c.widget.MinWidth
	}

	return w, h
}

func (c *Checkbox) SetLocation(rect img.Rectangle) {
	c.init.Do()
	c.widget.Rect = rect

}

func (c *Checkbox) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	c.init.Do()
}

func (c *Checkbox) RequestRelayout() {
	c.init.Do()
}

func (c *Checkbox) Render(screen *ebiten.Image) {
	c.init.Do()

	c.widget.Render(screen)
	if c.label == nil {
		c.currentImage().Draw(screen, c.widget.Rect.Dx(), c.widget.Rect.Dy(), func(opts *ebiten.DrawImageOptions) {
			opts.GeoM.Translate(float64(c.widget.Rect.Min.X), float64(c.widget.Rect.Min.Y))
		})
	} else {
		lx, _ := c.label.PreferredSize()
		cx, cy := c.checkboxPreferredSize()

		// If the label is larger than the checkbox, center the checkbox
		checkboxStartY := c.widget.Rect.Min.Y
		if cy < c.widget.Rect.Dy() {
			checkboxStartY += ((c.widget.Rect.Dy() - cy) / 2)
		}

		if c.order == CHECKBOX_FIRST {
			c.currentImage().Draw(screen, cx, cy, func(opts *ebiten.DrawImageOptions) {
				opts.GeoM.Translate(float64(c.widget.Rect.Min.X), float64(checkboxStartY))
			})
			c.label.SetLocation(img.Rect(c.widget.Rect.Min.X+c.spacing+cx, c.widget.Rect.Min.Y, c.widget.Rect.Min.X+c.spacing+cx+lx, c.widget.Rect.Min.Y+c.widget.Rect.Dy()))
			c.label.Render(screen)
		} else {
			c.label.SetLocation(img.Rect(c.widget.Rect.Min.X, c.widget.Rect.Min.Y, c.widget.Rect.Min.X+lx, c.widget.Rect.Min.Y+c.widget.Rect.Dy()))
			c.label.Render(screen)
			c.currentImage().Draw(screen, cx, cy, func(opts *ebiten.DrawImageOptions) {
				opts.GeoM.Translate(float64(c.widget.Rect.Min.X+lx+c.spacing), float64(checkboxStartY))
			})

		}
	}
}

func (c *Checkbox) Update() {
	c.init.Do()

	c.widget.Update()
	if c.label != nil {
		c.label.GetWidget().Disabled = c.widget.Disabled
	}
	c.handleDefaultInput()
}

func (c *Checkbox) handleDefaultInput() {
	if !c.DisableDefaultKeys && c.focused &&
		(inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace)) {
		c.Click()
	}
}

/** Focuser Interface - Start **/

func (c *Checkbox) Focus(focused bool) {
	c.init.Do()
	c.GetWidget().FireFocusEvent(c, focused, img.Point{-1, -1})
	c.focused = focused
}

func (c *Checkbox) IsFocused() bool {
	return c.focused
}

func (c *Checkbox) TabOrder() int {
	return c.tabOrder
}

func (c *Checkbox) GetFocus(direction FocusDirection) Focuser {
	return c.focusMap[direction]
}

func (c *Checkbox) AddFocus(direction FocusDirection, focus Focuser) {
	c.focusMap[direction] = focus
}

/** Focuser Interface - End **/

func (c *Checkbox) Click() {
	c.init.Do()
	if !c.widget.Disabled {
		c.SetState(c.state.Advance(c.triState))
	}
}

func (c *Checkbox) createWidget() {
	c.widget = NewWidget(append([]WidgetOpt{
		WidgetOpts.TrackHover(true),
		WidgetOpts.CursorEnterHandler(func(args *WidgetCursorEnterEventArgs) {
			if !c.widget.Disabled {
				c.hovering = true
			}
		}),

		WidgetOpts.CursorExitHandler(func(args *WidgetCursorExitEventArgs) {
			c.hovering = false
		}),

		WidgetOpts.MouseButtonClickedHandler(func(args *WidgetMouseButtonClickedEventArgs) {
			if !c.widget.Disabled && args.Button == ebiten.MouseButtonLeft {
				c.Click()
			}
		}),
	}, c.widgetOpts...)...)

	if len(c.labelString) > 0 {
		c.label = NewLabel(
			LabelOpts.TextOpts(
				TextOpts.Position(TextPositionStart, TextPositionCenter),
				TextOpts.WidgetOpts(
					WidgetOpts.MouseButtonClickedHandler(func(args *WidgetMouseButtonClickedEventArgs) {
						c.Click()
					}),
				),
			),
			LabelOpts.LabelText(c.labelString),
		)
	}
}

func (s WidgetState) Advance(triState bool) WidgetState {
	if s == WidgetUnchecked {
		return WidgetChecked
	}

	if s == WidgetChecked {
		if triState {
			return WidgetGreyed
		}

		return WidgetUnchecked
	}

	return WidgetUnchecked
}

func (c *Checkbox) currentImage() *image.NineSlice {
	// Handle disabled
	if c.widget.Disabled {
		switch c.state {
		case WidgetChecked:
			if c.computedParams.Image.CheckedDisabled != nil {
				return c.computedParams.Image.CheckedDisabled
			}
		case WidgetUnchecked:
			if c.computedParams.Image.UncheckedDisabled != nil {
				return c.computedParams.Image.UncheckedDisabled
			}
		case WidgetGreyed:
			if c.computedParams.Image.UncheckedDisabled != nil {
				return c.computedParams.Image.GreyedDisabled
			}
		}
	}
	// Handle hovered
	if c.hovering || c.focused {
		switch c.state {
		case WidgetChecked:
			if c.computedParams.Image.CheckedHovered != nil {
				return c.computedParams.Image.CheckedHovered
			}
		case WidgetUnchecked:
			if c.computedParams.Image.UncheckedHovered != nil {
				return c.computedParams.Image.UncheckedHovered
			}
		case WidgetGreyed:
			if c.computedParams.Image.GreyedHovered != nil {
				return c.computedParams.Image.GreyedHovered
			}
		}
	}
	// Fallback to default images
	switch c.state {
	case WidgetChecked:
		return c.computedParams.Image.Checked

	case WidgetUnchecked:
		return c.computedParams.Image.Unchecked

	case WidgetGreyed:
		return c.computedParams.Image.Greyed
	}

	return c.computedParams.Image.Unchecked
}

func (c *Checkbox) setChildComputedParams() {
	c.label.definedParams.Color = c.computedParams.Label.Color
	c.label.definedParams.Face = c.computedParams.Label.Face
	c.label.definedParams.Padding = c.computedParams.Label.Padding
	c.label.Validate()
}
