package widget

import (
	"image"

	"github.com/blizzy78/ebitenui/input"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

type ComboButton struct {
	ContentVisible bool

	maxContentHeight int

	button  *Button
	content HasWidget
}

type ComboButtonOpt func(c *ComboButton)

type ComboButtonClickedEventArgs struct {
	ComboButton *ComboButton
}

type ComboButtonClickedHandlerFunc func(args *ComboButtonClickedEventArgs)

const ComboButtonOpts = comboButtonOpts(true)

type comboButtonOpts bool

func NewComboButton(opts ...ComboButtonOpt) *ComboButton {
	var c *ComboButton
	c = &ComboButton{
		button: NewButton(
			ButtonOpts.WithClickedHandler(func(args *ButtonClickedEventArgs) {
				c.ContentVisible = !c.ContentVisible
			}),
		),
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

func (o comboButtonOpts) WithLayoutData(ld interface{}) ComboButtonOpt {
	return func(c *ComboButton) {
		ButtonOpts.WithLayoutData(ld)(c.button)
	}
}

func (o comboButtonOpts) WithImage(i *ButtonImage) ComboButtonOpt {
	return func(c *ComboButton) {
		ButtonOpts.WithImage(i)(c.button)
	}
}

func (o comboButtonOpts) WithTextAndImage(label string, face font.Face, image *ButtonImageImage, color *ButtonTextColor) ComboButtonOpt {
	return func(c *ComboButton) {
		ButtonOpts.WithTextAndImage(label, face, image, color)(c.button)
	}
}

func (o comboButtonOpts) WithContent(c HasWidget) ComboButtonOpt {
	return func(cb *ComboButton) {
		cb.content = c
	}
}

func (o comboButtonOpts) WithClickedHandler(f ComboButtonClickedHandlerFunc) ComboButtonOpt {
	return func(c *ComboButton) {
		ButtonOpts.WithClickedHandler(func(args *ButtonClickedEventArgs) {
			f(&ComboButtonClickedEventArgs{
				ComboButton: c,
			})
		})(c.button)
	}
}

func (o comboButtonOpts) WithMaxContentHeight(h int) ComboButtonOpt {
	return func(c *ComboButton) {
		c.maxContentHeight = h
	}
}

func (c *ComboButton) GetWidget() *Widget {
	return c.button.GetWidget()
}

func (c *ComboButton) SetLocation(rect image.Rectangle) {
	c.button.GetWidget().Rect = rect
}

func (c *ComboButton) PreferredSize() (int, int) {
	return c.button.PreferredSize()
}

func (c *ComboButton) SetLabel(l string) {
	c.button.Text().Label = l
	c.button.RequestRelayout()
}

func (c *ComboButton) Label() string {
	return c.button.Text().Label
}

func (c *ComboButton) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
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
