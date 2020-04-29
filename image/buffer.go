package image

import "github.com/hajimehoshi/ebiten"

type BufferedImage struct {
	Width  int
	Height int

	image *ebiten.Image
}

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
