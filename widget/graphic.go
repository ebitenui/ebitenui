package widget

import (
	"github.com/blizzy78/ebitenui/image"

	img "image"

	"github.com/hajimehoshi/ebiten"
)

type Graphic struct {
	Image          *ebiten.Image
	ImageNineSlice *image.NineSlice

	widgetOpts []WidgetOpt

	init   *MultiOnce
	widget *Widget
}

type GraphicOpt func(g *Graphic)

const GraphicOpts = graphicOpts(true)

type graphicOpts bool

func NewGraphic(opts ...GraphicOpt) *Graphic {
	g := &Graphic{
		init: &MultiOnce{},
	}

	g.init.Append(g.createWidget)

	for _, o := range opts {
		o(g)
	}

	return g
}

func (o graphicOpts) WithWidgetOpt(opt WidgetOpt) GraphicOpt {
	return func(g *Graphic) {
		g.widgetOpts = append(g.widgetOpts, opt)
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
	g.init.Do()
	return g.widget
}

func (g *Graphic) SetLocation(rect img.Rectangle) {
	g.init.Do()
	g.widget.Rect = rect
}

func (g *Graphic) PreferredSize() (int, int) {
	g.init.Do()
	if g.Image != nil {
		return g.Image.Size()
	}
	return 50, 50
}

func (g *Graphic) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	g.init.Do()
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

func (g *Graphic) createWidget() {
	g.widget = NewWidget(g.widgetOpts...)
	g.widgetOpts = nil
}
