package widget

import (
	img "image"

	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/input"

	"github.com/hajimehoshi/ebiten"
)

type Container struct {
	BackgroundImage     *image.NineSlice
	AutoDisableChildren bool

	widgetOpts []WidgetOpt
	layout     Layouter

	init     *MultiOnce
	widget   *Widget
	children []HasWidget
}

type ContainerOpt func(c *Container)

type RemoveChildFunc func()

const ContainerOpts = containerOpts(true)

type containerOpts bool

func NewContainer(opts ...ContainerOpt) *Container {
	c := &Container{
		init: &MultiOnce{},
	}

	c.init.Append(c.createWidget)

	for _, o := range opts {
		o(c)
	}

	return c
}

func (o containerOpts) WithWidgetOpts(opts ...WidgetOpt) ContainerOpt {
	return func(c *Container) {
		c.widgetOpts = append(c.widgetOpts, opts...)
	}
}

func (o containerOpts) WithBackgroundImage(i *image.NineSlice) ContainerOpt {
	return func(c *Container) {
		c.BackgroundImage = i
	}
}

func (o containerOpts) WithAutoDisableChildren() ContainerOpt {
	return func(c *Container) {
		c.AutoDisableChildren = true
	}
}

func (o containerOpts) WithLayout(layout Layouter) ContainerOpt {
	return func(c *Container) {
		c.layout = layout
	}
}

func (c *Container) AddChild(child HasWidget) RemoveChildFunc {
	c.init.Do()

	if child == nil {
		panic("cannot add nil child")
	}

	c.children = append(c.children, child)

	child.GetWidget().parent = c.widget

	c.RequestRelayout()

	return func() {
		c.removeChild(child)
	}
}

func (c *Container) removeChild(child HasWidget) {
	index := -1
	for i, ch := range c.children {
		if ch == child {
			index = i
			break
		}
	}

	if index < 0 {
		return
	}

	c.children = append(c.children[:index], c.children[index+1:]...)

	child.GetWidget().parent = nil

	c.RequestRelayout()
}

func (c *Container) RequestRelayout() {
	c.init.Do()

	if c.layout != nil {
		if d, ok := c.layout.(Dirtyable); ok {
			d.MarkDirty()
		}
	}

	for _, ch := range c.children {
		if r, ok := ch.(Relayoutable); ok {
			r.RequestRelayout()
		}
	}
}

func (c *Container) GetWidget() *Widget {
	c.init.Do()
	return c.widget
}

func (c *Container) PreferredSize() (int, int) {
	c.init.Do()

	if c.layout == nil {
		return 50, 50
	}

	return c.layout.PreferredSize(c.children)
}

func (c *Container) SetLocation(rect img.Rectangle) {
	c.init.Do()
	c.widget.Rect = rect
}

func (c *Container) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	c.init.Do()

	if c.AutoDisableChildren {
		for _, ch := range c.children {
			ch.GetWidget().Disabled = c.widget.Disabled
		}
	}

	c.doLayout()

	c.widget.Render(screen, def)

	c.draw(screen)

	for _, ch := range c.children {
		if cr, ok := ch.(Renderer); ok {
			cr.Render(screen, def)
		}
	}
}

func (c *Container) doLayout() {
	if c.layout != nil {
		c.layout.Layout(c.children, c.widget.Rect)
	}
}

func (c *Container) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	c.init.Do()

	for _, ch := range c.children {
		if il, ok := ch.(input.InputLayerer); ok {
			il.SetupInputLayer(def)
		}
	}
}

func (c *Container) draw(screen *ebiten.Image) {
	if c.BackgroundImage != nil {
		c.BackgroundImage.Draw(screen, c.widget.Rect.Dx(), c.widget.Rect.Dy(), c.widget.drawImageOptions)
	}
}

func (c *Container) createWidget() {
	c.widget = NewWidget(c.widgetOpts...)
	c.widgetOpts = nil
}

// WidgetAt implements WidgetLocator.
func (c *Container) WidgetAt(x int, y int) HasWidget {
	c.init.Do()

	p := img.Point{x, y}

	if !p.In(c.GetWidget().Rect) {
		return nil
	}

	for _, ch := range c.children {
		if wl, ok := ch.(WidgetLocator); ok {
			if w := wl.WidgetAt(x, y); w != nil {
				return w
			}

			continue
		}

		if p.In(ch.GetWidget().Rect) {
			return ch
		}
	}

	return c
}
