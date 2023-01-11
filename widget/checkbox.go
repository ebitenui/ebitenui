package widget

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mcarpenter622/ebitenui/event"
	"github.com/mcarpenter622/ebitenui/input"
)

type Checkbox struct {
	StateChangedEvent *event.Event
	buttonOpts        []ButtonOpt
	image             *CheckboxGraphicImage
	triState          bool

	init   *MultiOnce
	button *Button
	state  WidgetState
}

type CheckboxOpt func(c *Checkbox)

type CheckboxGraphicImage struct {
	Unchecked *ButtonImageImage
	Checked   *ButtonImageImage
	Greyed    *ButtonImageImage
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
	}

	c.init.Append(c.createWidget)

	for _, o := range opts {
		o(c)
	}

	return c
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

func (o CheckboxOptions) StateChangedHandler(f CheckboxChangedHandlerFunc) CheckboxOpt {
	return func(c *Checkbox) {
		c.StateChangedEvent.AddHandler(func(args interface{}) {
			f(args.(*CheckboxChangedEventArgs))
		})
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

func (c *Checkbox) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	c.init.Do()

	c.button.GraphicImage = c.state.graphicImage(c.image)

	c.button.Render(screen, def)
}

func (c *Checkbox) createWidget() {
	c.button = NewButton(append(c.buttonOpts, []ButtonOpt{
		ButtonOpts.Graphic(c.image.Unchecked.Idle),

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

func (s WidgetState) graphicImage(i *CheckboxGraphicImage) *ButtonImageImage {
	if s == WidgetChecked {
		return i.Checked
	}

	if s == WidgetGreyed {
		return i.Greyed
	}

	return i.Unchecked
}
