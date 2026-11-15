package image

import (
	img "image"

	"github.com/hajimehoshi/ebiten/v2"
)

// BufferedImage is a wrapper for an Ebiten Image that helps with caching the Image.
// As long as Width and Height stay the same, no new Image will be created.
type BufferedImage struct {
	Width  int
	Height int

	image *ebiten.Image
}

// MaskedRenderBuffer is a helper to draw images using a mask.
type MaskedRenderBuffer struct {
	renderBuf *BufferedImage
	maskedBuf *BufferedImage

	// For DrawIn: images covering the masked area, in screen coordinates.
	renderIn, maskedIn *ebiten.Image
	op                 ebiten.DrawImageOptions
}

// coveringMargin is how far, in pixels, an image covering an area reaches
// past it once the area has moved, so that it can keep sliding inside.
const coveringMargin = 64

// DrawFunc is a function that draws something into buf.
type DrawFunc func(buf *ebiten.Image)

// Image returns the internal Ebiten Image. If b.Width or b.Height have changed, a new Image
// will be created and returned, otherwise the cached Image will be returned.
func (b *BufferedImage) Image() *ebiten.Image {
	w, h := -1, -1
	if b.image != nil {
		w, h = b.image.Bounds().Dx(), b.image.Bounds().Dy()
	}

	if b.image == nil || b.Width != w || b.Height != h {
		if b.image != nil {
			b.image.Deallocate()
		}
		b.image = ebiten.NewImage(b.Width, b.Height)
	}

	return b.image
}

// NewMaskedRenderBuffer returns a new MaskedRenderBuffer.
func NewMaskedRenderBuffer() *MaskedRenderBuffer {
	return &MaskedRenderBuffer{
		renderBuf: &BufferedImage{},
		maskedBuf: &BufferedImage{},
	}
}

// Draw calls d to draw onto screen, using the mask drawn by dm. The buffer images passed
// to d and dm are of the same size as screen.
func (m *MaskedRenderBuffer) Draw(screen *ebiten.Image, d DrawFunc, dm DrawFunc) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

	m.renderBuf.Width, m.renderBuf.Height = w, h
	renderBuf := m.renderBuf.Image()
	renderBuf.Clear()

	m.maskedBuf.Width, m.maskedBuf.Height = w, h
	maskedBuf := m.maskedBuf.Image()
	maskedBuf.Clear()

	d(renderBuf)
	dm(maskedBuf)

	maskedBuf.DrawImage(renderBuf, &ebiten.DrawImageOptions{
		Blend: ebiten.BlendSourceIn,
	})

	screen.DrawImage(maskedBuf, nil)
}

// DrawIn calls d to draw, and dm to draw the mask, onto images covering only
// rect, in the same coordinates as screen, then draws what d drew through
// the mask onto screen. Widgets drawn by d that lie outside rect can skip
// drawing. The images are as big as rect, not screen, while rect stays put;
// once it moves they reach coveringMargin past it, and are kept while it
// slides inside them, as in a window easing in.
func (m *MaskedRenderBuffer) DrawIn(screen *ebiten.Image, rect img.Rectangle, d DrawFunc, dm DrawFunc) {
	rect = rect.Intersect(screen.Bounds())
	if rect.Empty() {
		return
	}
	renderBuf := imageCovering(&m.renderIn, rect, screen.Bounds())
	maskedBuf := imageCovering(&m.maskedIn, rect, screen.Bounds())

	d(renderBuf)
	dm(maskedBuf)

	m.op.GeoM.Reset()
	m.op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	m.op.Blend = ebiten.BlendSourceIn
	maskedBuf.DrawImage(renderBuf, &m.op)
	m.op.Blend = ebiten.Blend{}
	screen.DrawImage(maskedBuf, &m.op)
}

// imageCovering returns a cleared image whose bounds are rect: the part of
// *i covering rect while rect lies inside it (Ebitengine keeps one sub-image
// per area, so an area that stays gets the same one), else part of a new
// image, exactly rect the first time and reaching coveringMargin past it,
// within screen, once rect has moved.
func imageCovering(i **ebiten.Image, rect, screen img.Rectangle) *ebiten.Image {
	if *i == nil || !rect.In((*i).Bounds()) {
		margin := 0
		if *i != nil {
			margin = coveringMargin
			(*i).Deallocate()
		}
		*i = ebiten.NewImageWithOptions(rect.Inset(-margin).Intersect(screen), nil)
	}
	part, ok := (*i).SubImage(rect).(*ebiten.Image)
	if !ok {
		part = *i
	}
	part.Clear()
	return part
}
