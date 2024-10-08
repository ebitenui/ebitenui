package widget

import (
	"github.com/ebitenui/ebitenui/image"

	img "image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Graphic struct {
	Image          *ebiten.Image
	ImageNineSlice *image.NineSlice
	Images         *GraphicImage

	widgetOpts []WidgetOpt

	init   *MultiOnce
	widget *Widget
}

type GraphicImage struct {
	Idle     *ebiten.Image
	Disabled *ebiten.Image
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

	g.validate()

	return g
}

func (g *Graphic) validate() {
	if (g.Image != nil && g.ImageNineSlice != nil) || (g.Image != nil && g.Images != nil) || (g.ImageNineSlice != nil && g.Images != nil) {
		panic("Only one type of image: Image, ImageNineSlice or Images can be defined")
	}
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

func (o GraphicOptions) Images(i *GraphicImage) GraphicOpt {
	return func(g *Graphic) {
		g.Images = i
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

func (g *Graphic) Render(screen *ebiten.Image) {
	g.init.Do()
	g.widget.Render(screen)
	g.draw(screen)
}

func (g *Graphic) Update() {
	g.init.Do()

	g.widget.Update()
}

func (g *Graphic) draw(screen *ebiten.Image) {
	if g.ImageNineSlice != nil {
		g.ImageNineSlice.Draw(screen, g.widget.Rect.Dx(), g.widget.Rect.Dy(), g.widget.drawImageOptions)
		return
	}
	if g.Image != nil {
		g.Images = &GraphicImage{
			Idle: g.Image,
		}
	}

	i := g.Images.Idle
	if g.widget.Disabled && g.Images.Disabled != nil {
		i = g.Images.Disabled
	}

	opts := ebiten.DrawImageOptions{}
	w, h := i.Size()
	opts.GeoM.Translate(float64((g.widget.Rect.Dx()-w)/2), float64((g.widget.Rect.Dy()-h)/2))
	g.widget.drawImageOptions(&opts)
	screen.DrawImage(i, &opts)
}

func (g *Graphic) createWidget() {
	g.widget = NewWidget(g.widgetOpts...)
	g.widgetOpts = nil
}
