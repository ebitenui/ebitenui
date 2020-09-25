package widget

import (
	"image"

	"github.com/blizzy78/ebitenui/input"

	"github.com/hajimehoshi/ebiten"
)

type ComboButton struct {
	ContentVisible bool

	buttonOpts       []ButtonOpt
	maxContentHeight int

	init    *MultiOnce
	button  *Button
	content HasWidget
}

type ComboButtonOpt func(c *ComboButton)

const ComboButtonOpts = comboButtonOpts(true)

type comboButtonOpts bool

func NewComboButton(opts ...ComboButtonOpt) *ComboButton {
	c := &ComboButton{
		init: &MultiOnce{},
	}

	c.init.Append(c.createWidget)

	for _, o := range opts {
		o(c)
	}

	return c
}

func (o comboButtonOpts) ButtonOpts(opts ...ButtonOpt) ComboButtonOpt {
	return func(c *ComboButton) {
		c.buttonOpts = append(c.buttonOpts, opts...)
	}
}

func (o comboButtonOpts) Content(c HasWidget) ComboButtonOpt {
	return func(cb *ComboButton) {
		cb.content = c
	}
}

func (o comboButtonOpts) MaxContentHeight(h int) ComboButtonOpt {
	return func(c *ComboButton) {
		c.maxContentHeight = h
	}
}

func (c *ComboButton) GetWidget() *Widget {
	c.init.Do()
	return c.button.GetWidget()
}

func (c *ComboButton) SetLocation(rect image.Rectangle) {
	c.init.Do()
	c.button.GetWidget().Rect = rect
}

func (c *ComboButton) PreferredSize() (int, int) {
	c.init.Do()
	return c.button.PreferredSize()
}

func (c *ComboButton) SetLabel(l string) {
	c.init.Do()
	c.button.Text().Label = l
	c.button.RequestRelayout()
}

func (c *ComboButton) Label() string {
	c.init.Do()
	return c.button.Text().Label
}

func (c *ComboButton) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	c.init.Do()

	c.button.SetupInputLayer(def)

	if c.content != nil && c.ContentVisible {
		def(func(def input.DeferredSetupInputLayerFunc) {
			c.content.GetWidget().ElevateToNewInputLayer(&input.Layer{
				DebugLabel: "combo button content visible",
				EventTypes: input.LayerEventTypeAll,
				BlockLower: true,
				FullScreen: false,
				RectFunc: func() image.Rectangle {
					return c.content.GetWidget().Rect
				},
			})

			if il, ok := c.content.(input.InputLayerer); ok {
				il.SetupInputLayer(def)
			}
		})
	}
}

func (c *ComboButton) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	c.init.Do()

	c.button.Render(screen, def)

	if c.content != nil && c.ContentVisible {
		def(c.renderContent)
	}
}

func (c *ComboButton) renderContent(screen *ebiten.Image, def DeferredRenderFunc) {
	if input.MouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := input.CursorPosition()
		p := image.Point{x, y}
		if !p.In(c.button.GetWidget().Rect) && !p.In(c.content.GetWidget().Rect) {
			c.ContentVisible = false
			return
		}
	}

	r, ok := c.content.(Renderer)
	if !ok {
		return
	}

	if l, ok := c.content.(Locateable); ok {
		rect := c.button.GetWidget().Rect
		x, y := rect.Min.X, rect.Max.Y+2

		var w int
		var h int
		if p, ok := c.content.(PreferredSizer); ok {
			w, h = p.PreferredSize()
		} else {
			w, h = 50, 50
		}

		if c.maxContentHeight > 0 && h > c.maxContentHeight {
			h = c.maxContentHeight
		}

		l.SetLocation(image.Rect(x, y, x+w, y+h))
		if r, ok := c.content.(Relayoutable); ok {
			r.RequestRelayout()
		}
	}

	r.Render(screen, def)
}

func (c *ComboButton) createWidget() {
	c.button = NewButton(append(c.buttonOpts, []ButtonOpt{
		ButtonOpts.ClickedHandler(func(args *ButtonClickedEventArgs) {
			c.ContentVisible = !c.ContentVisible
		}),
	}...)...)
	c.buttonOpts = nil
}
