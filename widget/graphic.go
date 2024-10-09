package widget

import (
	"image/color"
	"image/gif"
	"time"

	"github.com/ebitenui/ebitenui/image"

	img "image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Graphic struct {
	Image          *ebiten.Image
	ImageNineSlice *image.NineSlice
	images         *GraphicImage
	gif            *gif.GIF

	gifImages         []*ebiten.Image
	gifCurrentImage   int
	gifCurrentImageAt time.Time

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
	hasImage := false
	for _, b := range []bool{g.Image != nil, g.ImageNineSlice != nil, g.gif != nil, g.images != nil} {
		if b && !hasImage {
			hasImage = true
		} else if b && hasImage {
			panic("Only one type of image: Image, ImageNineSlice, Images or GIF can be defined")
		}
	}
	if !hasImage {
		panic("One image is required: Image, ImageNineSlice, Images or GIF")
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
		g.images = i
	}
}

func (o GraphicOptions) GIF(gif *gif.GIF) GraphicOpt {
	return func(g *Graphic) {
		g.gif = gif
		g.gifImages = make([]*ebiten.Image, len(gif.Image))
		images := make([]img.Image, len(gif.Image))
		rect := gif.Image[0].Rect
		for i := 0; i < len(gif.Image); i++ {
			if i == 0 {
				g.gifImages[i] = ebiten.NewImageFromImage(gif.Image[i])
				images[i] = gif.Image[i]
				continue
			}
			img := restoreFrame(gif.Image[i], images[i-1], rect)
			g.gifImages[i] = ebiten.NewImageFromImage(img)
			images[i] = img
		}
	}
}

func restoreFrame(current *img.Paletted, prev img.Image, rect img.Rectangle) img.Image {
	img := img.NewRGBA(rect)
	for x := 0; x < rect.Dx(); x++ {
		for y := 0; y < rect.Dy(); y++ {
			if isInRect(x, y, current.Rect) && isOpaque(current.At(x, y)) {
				img.Set(x, y, current.At(x, y))
			} else {
				img.Set(x, y, prev.At(x, y))
			}
		}
	}
	return img
}

func isOpaque(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a > 0
}

func isInRect(x, y int, r img.Rectangle) bool {
	return r.Min.X <= x && x < r.Max.X &&
		r.Min.Y <= y && y < r.Max.Y
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
		s := g.Image.Bounds().Size()
		return s.X, s.Y
	} else if g.gif != nil {
		s := g.gifImages[0].Bounds().Size()
		return s.X, s.Y
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

	if g.gif != nil {
		if g.gifCurrentImageAt.IsZero() {
			g.gifCurrentImageAt = time.Now()
		}

		if time.Now().Sub(g.gifCurrentImageAt) > time.Duration(g.gif.Delay[g.gifCurrentImage]*10)*time.Millisecond {
			g.gifCurrentImageAt = time.Now()
			g.gifCurrentImage += 1
			if g.gifCurrentImage >= len(g.gifImages) {
				g.gifCurrentImage = 0
			}
		}
	}

	g.widget.Update()
}

func (g *Graphic) draw(screen *ebiten.Image) {
	if g.ImageNineSlice != nil {
		g.ImageNineSlice.Draw(screen, g.widget.Rect.Dx(), g.widget.Rect.Dy(), g.widget.drawImageOptions)
		return
	}

	if g.Image != nil {
		g.images = &GraphicImage{
			Idle: g.Image,
		}
	} else if g.gif != nil {
		g.images = &GraphicImage{
			Idle: g.gifImages[g.gifCurrentImage],
		}
	}

	i := g.images.Idle
	if g.widget.Disabled && g.images.Disabled != nil {
		i = g.images.Disabled
	}

	opts := ebiten.DrawImageOptions{}
	ib := i.Bounds()
	w, h := ib.Dx(), ib.Dy()
	opts.GeoM.Translate(float64((g.widget.Rect.Dx()-w)/2), float64((g.widget.Rect.Dy()-h)/2))
	g.widget.drawImageOptions(&opts)
	screen.DrawImage(i, &opts)
}

func (g *Graphic) createWidget() {
	g.widget = NewWidget(g.widgetOpts...)
	g.widgetOpts = nil
}
