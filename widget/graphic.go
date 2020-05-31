package widget

import (
	"github.com/blizzy78/ebitenui/image"

	img "image"

	"github.com/hajimehoshi/ebiten"
)

type Graphic struct {
	Image          *ebiten.Image
	ImageNineSlice *image.NineSlice

	widget *Widget
}

type GraphicOpt func(g *Graphic)

const GraphicOpts = graphicOpts(true)

type graphicOpts bool

func NewGraphic(opts ...GraphicOpt) *Graphic {
	g := &Graphic{
		widget: NewWidget(),
	}

	for _, o := range opts {
		o(g)
	}

	return g
}

func (o graphicOpts) WithLayoutData(ld interface{}) GraphicOpt {
	return func(g *Graphic) {
		WidgetOpts.WithLayoutData(ld)(g.widget)
	}
}

func (o graphicOpts) WithImage(i *ebiten.Image) GraphicOpt {
	return func(g *Graphic) {
		g.Image = i
	}
}

func (o graphicOpts) WithImageNineSlice(i *image.NineSlice) GraphicOpt {
	return func(g *Graphic) {
		g.ImageNineSlice = i
	}
}

func (g *Graphic) GetWidget() *Widget {
	return g.widget
}

func (g *Graphic) SetLocation(rect img.Rectangle) {
	g.widget.Rect = rect
}

func (g *Graphic) PreferredSize() (int, int) {
	if g.Image != nil {
		return g.Image.Size()
	}
	return 50, 50
}

func (g *Graphic) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	g.widget.Render(screen, def)
	g.draw(screen)
}

func (g *Graphic) draw(screen *ebiten.Image) {
	if g.Image != nil {
		opts := ebiten.DrawImageOptions{}
		w, h := g.Image.Size()
		opts.GeoM.Translate(float64((g.widget.Rect.Dx()-w)/2), float64((g.widget.Rect.Dy()-h)/2))
		g.widget.drawImageOptions(&opts)
		_ = screen.DrawImage(g.Image, &opts)
	} else if g.ImageNineSlice != nil {
		g.ImageNineSlice.Draw(screen, g.widget.Rect.Dx(), g.widget.Rect.Dy(), g.widget.drawImageOptions)
	}
}
