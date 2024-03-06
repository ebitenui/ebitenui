package widget

import (
	img "image"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type Container struct {
	BackgroundImage     *image.NineSlice
	AutoDisableChildren bool

	widgetOpts  []WidgetOpt
	layout      Layouter
	layoutDirty bool

	init     *MultiOnce
	widget   *Widget
	children []PreferredSizeLocateableWidget
}

type ContainerOpt func(c *Container)

type RemoveChildFunc func()

type ContainerOptions struct {
}

var ContainerOpts ContainerOptions

type PreferredSizeLocateableWidget interface {
	HasWidget
	PreferredSizer
	Locateable
}

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

func (o ContainerOptions) WidgetOpts(opts ...WidgetOpt) ContainerOpt {
	return func(c *Container) {
		c.widgetOpts = append(c.widgetOpts, opts...)
	}
}

func (o ContainerOptions) BackgroundImage(i *image.NineSlice) ContainerOpt {
	return func(c *Container) {
		c.BackgroundImage = i
	}
}

func (o ContainerOptions) AutoDisableChildren() ContainerOpt {
	return func(c *Container) {
		c.AutoDisableChildren = true
	}
}

func (o ContainerOptions) Layout(layout Layouter) ContainerOpt {
	return func(c *Container) {
		c.layout = layout
	}
}

func (c *Container) AddChild(child PreferredSizeLocateableWidget) RemoveChildFunc {
	c.init.Do()

	if child == nil {
		panic("cannot add nil child")
	}

	c.children = append(c.children, child)

	child.GetWidget().parent = c.widget
	child.GetWidget().self = child

	child.GetWidget().ContextMenuEvent.AddHandler(func(args interface{}) {
		a := args.(*WidgetContextMenuEventArgs)
		c.GetWidget().FireContextMenuEvent(a.Widget, a.Location)
	})
	child.GetWidget().FocusEvent.AddHandler(func(args interface{}) {
		a := args.(*WidgetFocusEventArgs)
		c.GetWidget().FireFocusEvent(a.Widget, a.Focused, a.Location)
	})
	child.GetWidget().ToolTipEvent.AddHandler(func(args interface{}) {
		a := args.(*WidgetToolTipEventArgs)
		c.GetWidget().FireToolTipEvent(a.Window, a.Show)
	})
	child.GetWidget().DragAndDropEvent.AddHandler(func(args interface{}) {
		a := args.(*WidgetDragAndDropEventArgs)
		c.GetWidget().FireDragAndDropEvent(a.Window, a.Show, a.DnD)
	})
	c.RequestRelayout()

	return func() {
		c.RemoveChild(child)
	}
}

func (c *Container) RemoveChild(child PreferredSizeLocateableWidget) {
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

	if child.GetWidget().ToolTip != nil && child.GetWidget().ToolTip.window != nil {
		child.GetWidget().ToolTip.window.Close()
	}

	if child.GetWidget().DragAndDrop != nil && child.GetWidget().DragAndDrop.window != nil {
		child.GetWidget().DragAndDrop.window.Close()
	}

	if child.GetWidget().ContextMenuWindow != nil {
		child.GetWidget().ContextMenuWindow.Close()
	}
	c.RequestRelayout()
}

func (c *Container) RemoveChildren() {
	for i := range c.children {
		childWidget := c.children[i].GetWidget()
		childWidget.parent = nil

		if childWidget.ToolTip != nil && childWidget.ToolTip.window != nil {
			childWidget.ToolTip.window.Close()
		}

		if childWidget.DragAndDrop != nil && childWidget.DragAndDrop.window != nil {
			childWidget.DragAndDrop.window.Close()
		}

		if childWidget.ContextMenuWindow != nil {
			childWidget.ContextMenuWindow.Close()
		}
	}
	c.children = nil

	c.RequestRelayout()
}

func (c *Container) Children() []PreferredSizeLocateableWidget {
	return c.children
}

func (c *Container) RequestRelayout() {
	c.init.Do()

	c.layoutDirty = true

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
	w, h := c.layout.PreferredSize(c.children)
	if c.widget != nil && h < c.widget.MinHeight {
		h = c.widget.MinHeight
	}
	if c.widget != nil && w < c.widget.MinWidth {
		w = c.widget.MinWidth
	}
	return w, h
}

func (c *Container) SetLocation(rect img.Rectangle) {
	c.init.Do()
	c.widget.Rect = rect
	c.RequestRelayout()
}

func (c *Container) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	c.init.Do()

	if c.widget.Visibility == Visibility_Hide_Blocking || c.widget.Visibility == Visibility_Hide {
		return
	}

	if c.AutoDisableChildren {
		for _, ch := range c.children {
			ch.GetWidget().Disabled = c.widget.Disabled
		}
	}

	c.widget.Render(screen, def)

	c.doLayout()

	c.draw(screen)

	for _, ch := range c.children {
		if cr, ok := ch.(Renderer); ok {
			if ch.GetWidget().Visibility == Visibility_Hide_Blocking || ch.GetWidget().Visibility == Visibility_Hide {
				continue
			}
			cr.Render(screen, def)
		}
	}
}

func (c *Container) doLayout() {
	if c.layout != nil && c.layoutDirty {
		c.layout.Layout(c.children, c.widget.Rect)
		c.layoutDirty = false
	}
}

func (c *Container) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	c.init.Do()

	for _, ch := range c.children {
		if il, ok := ch.(input.Layerer); ok {
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

func (c *Container) GetFocusers() []Focuser {
	var result []Focuser
	for _, child := range c.children {
		switch v := child.(type) {
		case Focuser:
			if v.TabOrder() >= 0 && !v.(HasWidget).GetWidget().Disabled {
				result = append(result, v)
			}
		case *Container:
			result = append(result, v.GetFocusers()...)
		case *FlipBook:
			result = append(result, v.GetFocusers()...)
		case *TabBook:
			result = append(result, v.container.GetFocusers()...)
		case *TabBookTab:
			result = append(result, v.GetFocusers()...)
		case *ScrollContainer:
			result = append(result, v.GetFocusers()...)
		case *TextArea:
			result = append(result, v.GetFocusers()...)
		}
	}
	return result
}

func (c *Container) GetDropTargets() []HasWidget {
	var result []HasWidget
	if c.GetWidget().drop != nil {
		result = append(result, c)
	}
	for _, child := range c.children {
		if v, ok := child.(Dropper); ok {
			result = append(result, v.GetDropTargets()...)
		}
	}

	return result
}

// WidgetAt implements WidgetLocator.
func (c *Container) WidgetAt(x int, y int) HasWidget {
	c.init.Do()

	p := img.Point{x, y}

	if !p.In(c.GetWidget().Rect) {
		return nil
	}

	for _, ch := range c.children {
		if wl, ok := ch.(Locater); ok {
			w := wl.WidgetAt(x, y)
			if w != nil {
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
