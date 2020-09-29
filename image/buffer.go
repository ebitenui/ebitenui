package image

import "github.com/hajimehoshi/ebiten"

// BufferedImage is a wrapper for an Ebiten Image that helps with caching the Image.
// As long as Width and Height stay the same, no new Image will be created.
type BufferedImage struct {
	Width  int
	Height int

	image *ebiten.Image
}

type MaskedRenderBuffer struct {
	renderBuf *BufferedImage
	maskedBuf *BufferedImage
}

type DrawFunc func(buf *ebiten.Image)

// Image returns the internal Ebiten Image. If b.Width or b.Height have changed, a new Image
// will be created and returned, otherwise the cached Image will be returned.
func (b *BufferedImage) Image() *ebiten.Image {
	w, h := -1, -1
	if b.image != nil {
		w, h = b.image.Size()
	}

	if b.image == nil || b.Width != w || b.Height != h {
		b.image, _ = ebiten.NewImage(b.Width, b.Height, ebiten.FilterDefault)
	}

	return b.image
}

func NewMaskedRenderBuffer() *MaskedRenderBuffer {
	return &MaskedRenderBuffer{
		renderBuf: &BufferedImage{},
		maskedBuf: &BufferedImage{},
	}
}

func (m *MaskedRenderBuffer) Draw(screen *ebiten.Image, w int, h int, d DrawFunc, dm DrawFunc) {
	m.renderBuf.Width, m.renderBuf.Height = w, h
	renderBuf := m.renderBuf.Image()
	_ = renderBuf.Clear()

	m.maskedBuf.Width, m.maskedBuf.Height = w, h
	maskedBuf := m.maskedBuf.Image()
	_ = maskedBuf.Clear()

	d(renderBuf)
	dm(maskedBuf)

	_ = maskedBuf.DrawImage(renderBuf, &ebiten.DrawImageOptions{
		CompositeMode: ebiten.CompositeModeSourceIn,
	})

	_ = screen.DrawImage(maskedBuf, nil)
}
