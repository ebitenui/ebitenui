package widget

import (
	"image"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
)

type Checkbox struct {
	StateChangedEvent *event.Event
	buttonOpts        []ButtonOpt
	image             *CheckboxGraphicImage
	triState          bool

	init   *MultiOnce
	button *Button
	state  WidgetState

	tabOrder int

	focusMap map[FocusDirection]Focuser
}

type CheckboxOpt func(c *Checkbox)

type CheckboxGraphicImage struct {
	Unchecked *GraphicImage
	Checked   *GraphicImage
	Greyed    *GraphicImage
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

	c.init.Append(c.createWidget)

	for _, o := range opts {
		o(c)
	}

	c.validate()

	return c
}

func (c *Checkbox) validate() {
	if len(c.buttonOpts) == 0 {
		panic("Checkbox: ButtonOpts are required.")
	}
	if c.image == nil {
		panic("Checkbox: Image is required.")
	}
	if c.image.Checked == nil {
		panic("Checkbox: Image.Checked is required.")
	}
	if c.image.Checked.Idle == nil {
		panic("Checkbox: Image.Checked.Idle is required.")
	}

	if c.image.Unchecked == nil {
		panic("Checkbox: Image.Unchecked is required.")
	}
	if c.image.Unchecked.Idle == nil {
		panic("Checkbox: Image.Unchecked.Idle is required.")
	}
	if c.state == WidgetGreyed && !c.triState {
		panic("Checkbox: non-tri state Checkbox cannot be in greyed state")
	}
}

func (o CheckboxOptions) ButtonOpts(opts ...ButtonOpt) CheckboxOpt {
	return func(c *Checkbox) {
		c.buttonOpts = append(c.buttonOpts, opts...)
	}
}

func (o CheckboxOptions) Image(i *CheckboxGraphicImage) CheckboxOpt {
	return func(c *Checkbox) {
		c.image = i
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

	c.button.GraphicImage = c.state.graphicImage(c.image)

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
	c.button = NewButton(append(c.buttonOpts, []ButtonOpt{
		ButtonOpts.Graphic(c.image.Unchecked),
		ButtonOpts.ClickedHandler(func(_ *ButtonClickedEventArgs) {
			c.SetState(c.state.Advance(c.triState))
		}),
	}...)...)
	c.buttonOpts = nil
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

func (s WidgetState) graphicImage(i *CheckboxGraphicImage) *GraphicImage {
	if s == WidgetChecked {
		return i.Checked
	}

	if s == WidgetGreyed {
		return i.Greyed
	}

	return i.Unchecked
}
