package widget

import (
	"github.com/ebitenui/ebitenui/image"

	img "image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Graphic struct {
	Image          *ebiten.Image
	ImageNineSlice *image.NineSlice

	widgetOpts []WidgetOpt

	init   *MultiOnce
	widget *Widget
}

type GraphicOpt func(g *Graphic)

type GraphicOptions struct {
}

var GraphicOpts GraphicOptions

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

func (o GraphicOptions) WidgetOpts(opts ...WidgetOpt) GraphicOpt {
	return func(g *Graphic) {
		g.widgetOpts = append(g.widgetOpts, opts...)
	}
}

func (o GraphicOptions) Image(i *ebiten.Image) GraphicOpt {
	return func(g *Graphic) {
		g.Image = i
	}
}

func (o GraphicOptions) ImageNineSlice(i *image.NineSlice) GraphicOpt {
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
		screen.DrawImage(g.Image, &opts)
	} else if g.ImageNineSlice != nil {
		g.ImageNineSlice.Draw(screen, g.widget.Rect.Dx(), g.widget.Rect.Dy(), g.widget.drawImageOptions)
	}
}

func (g *Graphic) createWidget() {
	g.widget = NewWidget(g.widgetOpts...)
	g.widgetOpts = nil
}
