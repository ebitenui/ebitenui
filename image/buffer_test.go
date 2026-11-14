package image

import (
	img "image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

func TestBufferedImage_Image(t *testing.T) {
	is := is.New(t)

	b := &BufferedImage{}
	b.Width, b.Height = 100, 100
	i := b.Image()
	w, h := i.Bounds().Dx(), i.Bounds().Dy()
	is.Equal(w, 100)
	is.Equal(h, 100)

	b.Width, b.Height = 150, 70
	i = b.Image()
	w, h = i.Bounds().Dx(), i.Bounds().Dy()
	is.Equal(w, 150)
	is.Equal(h, 70)
}

func TestMaskedRenderBuffer_Draw(t *testing.T) {
	is := is.New(t)

	b := NewMaskedRenderBuffer()
	screen := newImageEmptySize(100, 100, t)

	draw := false
	drawMask := false

	b.Draw(screen, func(buf *ebiten.Image) {
		draw = true
	}, func(buf *ebiten.Image) {
		drawMask = true
	})

	is.True(draw)
	is.True(drawMask)
}

// DrawIn draws onto images covering only the area masked, in the screen's
// coordinates, kept while the area stays the same and replaced when it
// moves.
func TestMaskedRenderBuffer_DrawIn(t *testing.T) {
	is := is.New(t)

	screen := ebiten.NewImage(200, 200)
	m := NewMaskedRenderBuffer()
	var drawnOn, maskedOn *ebiten.Image
	draw := func(rect img.Rectangle) {
		m.DrawIn(screen, rect, func(buf *ebiten.Image) { drawnOn = buf }, func(buf *ebiten.Image) { maskedOn = buf })
	}

	draw(img.Rect(20, 30, 70, 90))
	is.Equal(drawnOn.Bounds(), img.Rect(20, 30, 70, 90)) // the area, where it is
	is.Equal(maskedOn.Bounds(), img.Rect(20, 30, 70, 90))
	first := drawnOn

	draw(img.Rect(20, 30, 70, 90))
	is.Equal(drawnOn, first) // kept

	draw(img.Rect(150, 150, 250, 250))
	is.Equal(drawnOn.Bounds(), img.Rect(150, 150, 200, 200)) // moved, and cut to the screen
}
