package image

import "github.com/hajimehoshi/ebiten"

// BufferedImage is a wrapper for an Ebiten Image that helps with caching the Image.
// As long as Width and Height stay the same, no new Image will be created.
type BufferedImage struct {
	Width  int
	Height int

	image *ebiten.Image
}

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
