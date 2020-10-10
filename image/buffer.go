package image

import "github.com/hajimehoshi/ebiten/v2"

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
}

// DrawFunc is a function that draws something into buf.
type DrawFunc func(buf *ebiten.Image)

// Image returns the internal Ebiten Image. If b.Width or b.Height have changed, a new Image
// will be created and returned, otherwise the cached Image will be returned.
func (b *BufferedImage) Image() *ebiten.Image {
	w, h := -1, -1
	if b.image != nil {
		w, h = b.image.Size()
	}

	if b.image == nil || b.Width != w || b.Height != h {
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
	w, h := screen.Size()

	m.renderBuf.Width, m.renderBuf.Height = w, h
	renderBuf := m.renderBuf.Image()
	renderBuf.Clear()

	m.maskedBuf.Width, m.maskedBuf.Height = w, h
	maskedBuf := m.maskedBuf.Image()
	maskedBuf.Clear()

	d(renderBuf)
	dm(maskedBuf)

	maskedBuf.DrawImage(renderBuf, &ebiten.DrawImageOptions{
		CompositeMode: ebiten.CompositeModeSourceIn,
	})

	screen.DrawImage(maskedBuf, nil)
}
