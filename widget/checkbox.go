package widget

import (
	"image"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
)

type CheckboxParams struct {
	Image *CheckboxGraphicImage
}

type Checkbox struct {
	definedParams  CheckboxParams
	computedParams CheckboxParams

	StateChangedEvent *event.Event
	widgetOpts        []WidgetOpt
	triState          bool

	init   *MultiOnce
	button *Button
	state  WidgetState

	tabOrder int

	focusMap map[FocusDirection]Focuser
}

type CheckboxOpt func(c *Checkbox)

type CheckboxGraphicImage struct {
	Unchecked *ButtonImage
	Checked   *ButtonImage
	Greyed    *ButtonImage
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

		init: &MultiOnce{},

		focusMap: make(map[FocusDirection]Focuser),
	}

	for _, o := range opts {
		o(c)
	}

	c.createWidget()

	return c
}

func (c *Checkbox) Validate() {
	c.populateComputedParams()

	if c.computedParams.Image == nil {
		panic("Checkbox: Image is required.")
	}
	if c.computedParams.Image.Checked == nil {
		panic("Checkbox: Image.Checked is required.")
	}
	if c.computedParams.Image.Checked.Idle == nil {
		panic("Checkbox: Image.Checked.Idle is required.")
	}

	if c.computedParams.Image.Unchecked == nil {
		panic("Checkbox: Image.Unchecked is required.")
	}
	if c.computedParams.Image.Unchecked.Idle == nil {
		panic("Checkbox: Image.Unchecked.Idle is required.")
	}
	if c.state == WidgetGreyed && !c.triState {
		panic("Checkbox: non-tri state Checkbox cannot be in greyed state.")
	}
	if c.triState {
		if c.computedParams.Image.Greyed == nil {
			panic("Checkbox: Image.Greyed is required for tri-state checkboxes.")
		} else if c.computedParams.Image.Greyed.Idle == nil {
			panic("Checkbox: Image.Greyed.Idle is required for tri-state checkboxes.")
		}
	}
}

func (c *Checkbox) populateComputedParams() {
	checkboxParams := CheckboxParams{
		Image: &CheckboxGraphicImage{},
	}
	theme := c.GetWidget().GetTheme()
	// clone the theme
	if theme != nil {
		if theme.CheckboxTheme != nil {
			if theme.CheckboxTheme.Image != nil {
				checkboxParams.Image = theme.CheckboxTheme.Image
			}
		}
	}

	if c.definedParams.Image != nil {
		if checkboxParams.Image == nil {
			checkboxParams.Image = c.definedParams.Image
		} else {
			if c.definedParams.Image.Checked != nil {
				checkboxParams.Image.Checked = c.definedParams.Image.Checked
			}
			if c.definedParams.Image.Unchecked != nil {
				checkboxParams.Image.Unchecked = c.definedParams.Image.Unchecked
			}
			if c.definedParams.Image.Greyed != nil {
				checkboxParams.Image.Greyed = c.definedParams.Image.Greyed
			}

		}
	}

	c.computedParams = checkboxParams
	c.button.computedParams.Image = c.state.graphicImage(c.computedParams.Image)
}

func (o CheckboxOptions) WidgetOpts(opts ...WidgetOpt) CheckboxOpt {
	return func(c *Checkbox) {
		c.widgetOpts = append(c.widgetOpts, opts...)
	}
}

func (o CheckboxOptions) Image(i *CheckboxGraphicImage) CheckboxOpt {
	return func(c *Checkbox) {
		c.definedParams.Image = i
	}
}

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

func (o CheckboxOptions) StateChangedHandler(f CheckboxChangedHandlerFunc) CheckboxOpt {
	return func(c *Checkbox) {
		c.StateChangedEvent.AddHandler(func(args interface{}) {
			if arg, ok := args.(*CheckboxChangedEventArgs); ok {
				f(arg)
			}
		})
	}
}

func (o CheckboxOptions) InitialState(state WidgetState) CheckboxOpt {
	return func(c *Checkbox) {
		c.state = state
	}
}

func (tw *Checkbox) State() WidgetState {
	return tw.state
}

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

func (tw *Checkbox) getStateChangedEvent() *event.Event {
	return tw.StateChangedEvent
}

func (c *Checkbox) GetWidget() *Widget {
	c.init.Do()
	return c.button.GetWidget()
}

func (c *Checkbox) PreferredSize() (int, int) {
	c.init.Do()
	return c.button.PreferredSize()
}

func (c *Checkbox) SetLocation(rect image.Rectangle) {
	c.init.Do()
	c.button.SetLocation(rect)
}

func (c *Checkbox) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	c.init.Do()
	c.button.SetupInputLayer(def)
}

func (c *Checkbox) Render(screen *ebiten.Image) {
	c.init.Do()

	c.button.computedParams.Image = c.state.graphicImage(c.computedParams.Image)

	c.button.Render(screen)
}

func (c *Checkbox) Update() {
	c.init.Do()

	c.button.Update()
}

/** Focuser Interface - Start **/

func (c *Checkbox) Focus(focused bool) {
	c.init.Do()
	c.GetWidget().FireFocusEvent(c, focused, image.Point{-1, -1})
	c.button.focused = focused
}

func (c *Checkbox) IsFocused() bool {
	return c.button.focused
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
	c.button.Click()
}

func (c *Checkbox) createWidget() {
	c.button = NewButton([]ButtonOpt{
		ButtonOpts.ClickedHandler(func(_ *ButtonClickedEventArgs) {
			c.SetState(c.state.Advance(c.triState))
		}),
		ButtonOpts.WidgetOpts(c.widgetOpts...),
	}...)
	c.widgetOpts = nil
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

func (s WidgetState) graphicImage(i *CheckboxGraphicImage) *ButtonImage {
	if s == WidgetChecked {
		return i.Checked
	}

	if s == WidgetGreyed {
		return i.Greyed
	}

	return i.Unchecked
}
