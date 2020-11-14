package demo

import (
	"image"

	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type sizedPanel struct {
	width     int
	height    int
	container *widget.Container
}

func newSizedPanel(w int, h int, opts ...widget.ContainerOpt) *sizedPanel {
	return &sizedPanel{
		width:     w,
		height:    h,
		container: widget.NewContainer(opts...),
	}
}

func (p *sizedPanel) GetWidget() *widget.Widget {
	return p.container.GetWidget()
}

func (p *sizedPanel) PreferredSize() (int, int) {
	return p.width, p.height
}

func (p *sizedPanel) SetLocation(rect image.Rectangle) {
	p.container.SetLocation(rect)
}

func (p *sizedPanel) Render(screen *ebiten.Image, def widget.DeferredRenderFunc) {
	p.container.Render(screen, def)
}

func (p *sizedPanel) Container() *widget.Container {
	return p.container
}
