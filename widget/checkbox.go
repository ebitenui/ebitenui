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

type Checkbox struct {
	init              *MultiOnce
	widget            *Widget
	widgetOpts        []WidgetOpt
	image             *CheckboxImage
	triState          bool
	StateChangedEvent *event.Event

	state    WidgetState
	hovering bool

	label    *Label
	labelOpt LabelOpt
	spacing  int
	order    LabelOrder

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

		spacing: 8,
		order:   CHECKBOX_FIRST,
		init:    &MultiOnce{},

		focusMap: make(map[FocusDirection]Focuser),
	}

	c.init.Append(c.createWidget)

	for _, o := range opts {
		o(c)
	}

	c.validate()

	return c
}

func (c *Checkbox) validate() {

	if c.image == nil {
		panic("Checkbox: Image is required.")
	}

	if c.image.Checked == nil {
		panic("Checkbox: Image.Checked is required.")
	}

	if c.image.Unchecked == nil {
		panic("Checkbox: Image.Unchecked is required.")
	}

	if c.triState && c.image.Greyed == nil {
		panic("Checkbox: Image.Greyed is required if this is a tristate checkbox.")
	}
	if c.state == WidgetGreyed && !c.triState {
		panic("Checkbox: non-tri state Checkbox cannot be in greyed state")
	}
}

func (o CheckboxOptions) WidgetOpts(opts ...WidgetOpt) CheckboxOpt {
	return func(c *Checkbox) {
		c.widgetOpts = append(c.widgetOpts, opts...)
	}
}

// This option allows you to specify a label to be shown before or after the checkbox.
func (o CheckboxOptions) Text(label string, face text.Face, color *LabelColor) CheckboxOpt {
	return func(l *Checkbox) {
		l.labelOpt = LabelOpts.Text(label, face, color)
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
		c.image = i
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

	iw, ih := c.image.Unchecked.MinSize()
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

	iw, ih := c.image.Unchecked.MinSize()
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

	if c.labelOpt != nil {
		c.label = NewLabel([]LabelOpt{c.labelOpt, LabelOpts.TextOpts(
			TextOpts.Position(TextPositionStart, TextPositionCenter),
			TextOpts.WidgetOpts(
				WidgetOpts.MouseButtonClickedHandler(func(args *WidgetMouseButtonClickedEventArgs) {
					c.Click()
				}),
			),
		)}...)
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
			if c.image.CheckedDisabled != nil {
				return c.image.CheckedDisabled
			}
		case WidgetUnchecked:
			if c.image.UncheckedDisabled != nil {
				return c.image.UncheckedDisabled
			}
		case WidgetGreyed:
			if c.image.UncheckedDisabled != nil {
				return c.image.GreyedDisabled
			}
		}
	}
	// Handle hovered
	if c.hovering || c.focused {
		switch c.state {
		case WidgetChecked:
			if c.image.CheckedHovered != nil {
				return c.image.CheckedHovered
			}
		case WidgetUnchecked:
			if c.image.UncheckedHovered != nil {
				return c.image.UncheckedHovered
			}
		case WidgetGreyed:
			if c.image.GreyedHovered != nil {
				return c.image.GreyedHovered
			}
		}
	}
	// Fallback to default images
	switch c.state {
	case WidgetChecked:
		return c.image.Checked

	case WidgetUnchecked:
		return c.image.Unchecked

	case WidgetGreyed:
		return c.image.Greyed
	}

	return c.image.Unchecked
}
